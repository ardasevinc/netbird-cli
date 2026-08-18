package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/ardasevinc/netbird-cli/internal/analysis"
	"github.com/ardasevinc/netbird-cli/internal/config"
	"github.com/ardasevinc/netbird-cli/internal/ledger"
	"github.com/ardasevinc/netbird-cli/internal/operations"
	"github.com/spf13/cobra"
)

type stagePlan struct {
	Operation     string           `json:"operation"`
	Request       json.RawMessage  `json:"request"`
	Before        json.RawMessage  `json:"before"`
	IntendedAfter json.RawMessage  `json:"intended_after"`
	Findings      []ledger.Finding `json:"findings"`
}

func validatePersistedSecretSafety(operation string, request json.RawMessage) error {
	if operation == "users.invites.accept" {
		var object map[string]any
		if err := json.Unmarshal(request, &object); err != nil {
			return fmt.Errorf("decode invite acceptance request: %w", err)
		}
		for _, field := range []string{"invite_token", "password"} {
			if _, ok := object[field]; ok {
				return fmt.Errorf("invite acceptance field %s cannot be persisted in a stage; use %s_ref", field, field)
			}
		}
		return nil
	}
	if operation == "users.password.update" {
		var object map[string]any
		if err := json.Unmarshal(request, &object); err != nil {
			return fmt.Errorf("decode password request: %w", err)
		}
		for _, field := range []string{"old_password", "new_password"} {
			if _, ok := object[field]; ok {
				return fmt.Errorf("password field %s cannot be persisted in a stage; use %s_ref", field, field)
			}
		}
		return nil
	}
	if operation == "event_streaming.create" || operation == "event_streaming.update" {
		var object map[string]any
		if err := json.Unmarshal(request, &object); err != nil {
			return fmt.Errorf("decode event-streaming request: %w", err)
		}
		if _, ok := object["config"]; ok {
			return errors.New("event-streaming config cannot be persisted in a stage; use config_ref")
		}
		if operation == "event_streaming.create" {
			ref, ok := object["config_ref"].(string)
			if !ok || strings.TrimSpace(ref) == "" {
				return errors.New("event-streaming create requires config_ref")
			}
		}
		return nil
	}
	if operation == "identity_providers.create" || operation == "identity_providers.update" {
		var object map[string]any
		if err := json.Unmarshal(request, &object); err != nil {
			return fmt.Errorf("decode identity provider request: %w", err)
		}
		if _, ok := object["client_secret"]; ok {
			return errors.New("identity provider client_secret cannot be persisted in a stage; use client_secret_ref")
		}
		if operation == "identity_providers.create" {
			ref, ok := object["client_secret_ref"].(string)
			if !ok || strings.TrimSpace(ref) == "" {
				return errors.New("identity provider create requires client_secret_ref")
			}
		}
		return nil
	}
	if operation == "azure_idp.create" || operation == "azure_idp.update" {
		var object map[string]any
		if err := json.Unmarshal(request, &object); err != nil {
			return fmt.Errorf("decode azure IDP request: %w", err)
		}
		if _, ok := object["client_secret"]; ok {
			return errors.New("azure IDP client_secret cannot be persisted in a stage; use client_secret_ref")
		}
		if operation == "azure_idp.create" {
			ref, ok := object["client_secret_ref"].(string)
			if !ok || strings.TrimSpace(ref) == "" {
				return errors.New("azure IDP create requires client_secret_ref")
			}
		}
		return nil
	}
	if operation == "google_idp.create" || operation == "google_idp.update" {
		var object map[string]any
		if err := json.Unmarshal(request, &object); err != nil {
			return fmt.Errorf("decode google IDP request: %w", err)
		}
		if _, ok := object["service_account_key"]; ok {
			return errors.New("google IDP service_account_key cannot be persisted in a stage; use service_account_key_ref")
		}
		if operation == "google_idp.create" {
			ref, ok := object["service_account_key_ref"].(string)
			if !ok || strings.TrimSpace(ref) == "" {
				return errors.New("google IDP create requires service_account_key_ref")
			}
		}
		return nil
	}
	if strings.HasPrefix(operation, "edr.") && (strings.HasSuffix(operation, ".create") || strings.HasSuffix(operation, ".update")) {
		var object map[string]any
		if err := json.Unmarshal(request, &object); err != nil {
			return fmt.Errorf("decode EDR request: %w", err)
		}
		provider := strings.TrimSuffix(strings.TrimPrefix(operation, "edr."), ".create")
		if strings.HasSuffix(operation, ".update") {
			provider = strings.TrimSuffix(strings.TrimPrefix(operation, "edr."), ".update")
		}
		fields := map[string]string{}
		switch provider {
		case "intune", "falcon":
			fields["secret"] = "secret_ref"
		case "sentinelone", "fleetdm":
			fields["api_token"] = "api_token_ref"
		case "huntress":
			fields["api_key"] = "api_key_ref"
			fields["api_secret"] = "api_secret_ref"
		default:
			return fmt.Errorf("unsupported EDR operation %q", operation)
		}
		for field, refField := range fields {
			if _, ok := object[field]; ok {
				return fmt.Errorf("EDR %s cannot be persisted in a stage; use %s", field, refField)
			}
			ref, ok := object[refField].(string)
			if !ok || strings.TrimSpace(ref) == "" {
				return fmt.Errorf("EDR %s requires %s", provider, refField)
			}
		}
		return nil
	}
	if operation == "reverse_proxy_services.create" || operation == "reverse_proxy_services.update" {
		var object map[string]any
		if err := json.Unmarshal(request, &object); err != nil {
			return fmt.Errorf("decode reverse proxy service request: %w", err)
		}
		if _, ok := object["auth"]; ok {
			return errors.New("reverse proxy service auth cannot be persisted in a stage; use auth_ref")
		}
		return nil
	}
	if operation == "notification_channels.create" || operation == "notification_channels.update" {
		var object map[string]any
		if err := json.Unmarshal(request, &object); err != nil {
			return fmt.Errorf("decode notification channel request: %w", err)
		}
		if _, ok := object["target"]; ok {
			return errors.New("notification channel target cannot be persisted in a stage; use target_ref")
		}
		if operation == "notification_channels.create" {
			ref, ok := object["target_ref"].(string)
			if !ok || strings.TrimSpace(ref) == "" {
				return errors.New("notification channel create requires target_ref")
			}
		}
		return nil
	}
	if operation != "agent_network.providers.create" && operation != "agent_network.providers.update" {
		return nil
	}
	var object map[string]any
	if err := json.Unmarshal(request, &object); err != nil {
		return fmt.Errorf("decode provider request: %w", err)
	}
	if _, ok := object["api_key"]; ok {
		return errors.New("provider api_key cannot be persisted in a stage; use api_key_ref")
	}
	return nil
}

func validateInvitePreimageSafety(before json.RawMessage) error {
	var object map[string]any
	if err := json.Unmarshal(before, &object); err != nil {
		return fmt.Errorf("decode invite preimage: %w", err)
	}
	for _, field := range []string{"token", "invite_token", "password"} {
		if _, ok := object[field]; ok {
			return fmt.Errorf("invite preimage field %s cannot be persisted", field)
		}
	}
	return nil
}

func validateDNSZoneUpdate(before, intendedAfter json.RawMessage) error {
	var beforeObject, afterObject map[string]any
	if err := json.Unmarshal(before, &beforeObject); err != nil {
		return fmt.Errorf("decode DNS zone preimage: %w", err)
	}
	if err := json.Unmarshal(intendedAfter, &afterObject); err != nil {
		return fmt.Errorf("decode DNS zone intended state: %w", err)
	}
	beforeDomain, beforeOK := beforeObject["domain"].(string)
	afterDomain, afterOK := afterObject["domain"].(string)
	if !beforeOK || !afterOK {
		return errors.New("DNS zone update requires a domain in both preimage and intended state")
	}
	if beforeDomain != afterDomain {
		return errors.New("DNS zone domain cannot be changed; create a new zone instead")
	}
	return nil
}

func stageCommand(state *commandState, stdout io.Writer) *cobra.Command {
	command := &cobra.Command{Use: "stage", Short: "create and inspect local mutation plans"}
	command.AddCommand(stageCreateCommand(state, stdout))
	command.AddCommand(stageShowCommand(state, stdout))
	command.AddCommand(stageCancelCommand(state, stdout))
	return command
}

func stageCreateCommand(state *commandState, stdout io.Writer) *cobra.Command {
	var fromJSON bool
	command := &cobra.Command{
		Use:   "create",
		Short: "persist a local plan from versioned structured input",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if !fromJSON {
				return fail(2, fmt.Errorf("stage create requires --from-json"))
			}
			file, err := config.Load(state.configPath)
			if err != nil {
				return fail(3, err)
			}
			profile, err := file.Profile(state.profileName)
			if err != nil {
				return fail(2, err)
			}
			var plan stagePlan
			decoder := json.NewDecoder(cmd.InOrStdin())
			decoder.DisallowUnknownFields()
			if err := decoder.Decode(&plan); err != nil {
				return fail(2, fmt.Errorf("decode stage plan: %w", err))
			}
			if plan.Operation == "" {
				return fail(2, fmt.Errorf("stage plan operation is required"))
			}
			if err := validatePersistedSecretSafety(plan.Operation, plan.Request); err != nil {
				return fail(2, err)
			}
			if plan.Operation == "dns.zones.update" {
				if err := validateDNSZoneUpdate(plan.Before, plan.IntendedAfter); err != nil {
					return fail(2, err)
				}
			}
			if plan.Operation == "users.invites.accept" {
				if err := validateInvitePreimageSafety(plan.Before); err != nil {
					return fail(2, err)
				}
			}
			impact := json.RawMessage(`{}`)
			findings := append([]ledger.Finding(nil), plan.Findings...)
			switch plan.Operation {
			case "groups.create", "groups.update", "groups.delete", "policies.create", "policies.update", "policies.delete", "routes.create", "routes.update", "routes.delete", "peers.update", "peers.delete", "peers.temporary_access.create", "peers.jobs.create", "peers.edr.bypass.create", "peers.edr.bypass.delete", "event_streaming.create", "event_streaming.update", "event_streaming.delete", "identity_providers.create", "identity_providers.update", "identity_providers.delete", "reverse_proxy_tokens.create", "reverse_proxy_tokens.delete", "reverse_proxy_domains.create", "reverse_proxy_domains.delete", "reverse_proxy_clusters.delete", "reverse_proxy_services.create", "reverse_proxy_services.update", "reverse_proxy_services.delete", "notification_channels.create", "notification_channels.update", "notification_channels.delete", "azure_idp.create", "azure_idp.update", "azure_idp.delete", "azure_idp.sync", "google_idp.create", "google_idp.update", "google_idp.delete", "google_idp.sync", "edr.intune.create", "edr.intune.update", "edr.intune.delete", "edr.sentinelone.create", "edr.sentinelone.update", "edr.sentinelone.delete", "edr.falcon.create", "edr.falcon.update", "edr.falcon.delete", "edr.huntress.create", "edr.huntress.update", "edr.huntress.delete", "edr.fleetdm.create", "edr.fleetdm.update", "edr.fleetdm.delete", "scim.create", "scim.update", "scim.delete", "scim.token", "okta_scim.create", "okta_scim.update", "okta_scim.delete", "okta_scim.token", "msp.tenants.create", "msp.tenants.update", "msp.tenants.dns", "msp.tenants.invite", "msp.tenants.invite.respond", "msp.tenants.subscription", "msp.tenants.unlink", "billing.aws_marketplace.activate", "billing.aws_marketplace.enrich", "billing.checkout.create", "billing.subscription.update", "networks.create", "networks.update", "networks.delete", "networks.resources.create", "networks.resources.update", "networks.resources.delete", "networks.routers.create", "networks.routers.update", "networks.routers.delete", "dns.zones.create", "dns.zones.update", "dns.zones.delete", "dns.records.create", "dns.records.update", "dns.records.delete", "dns.nameservers.create", "dns.nameservers.update", "dns.nameservers.delete", "dns.settings.update", "accounts.update", "accounts.delete", "posture_checks.create", "posture_checks.update", "posture_checks.delete", "ingress.peers.create", "ingress.peers.update", "ingress.peers.delete", "peers.ingress.ports.create", "peers.ingress.ports.update", "peers.ingress.ports.delete", "agent_network.settings.update", "agent_network.settings.create", "agent_network.settings.delete", "agent_network.budget_rules.create", "agent_network.budget_rules.update", "agent_network.budget_rules.delete", "agent_network.guardrails.create", "agent_network.guardrails.update", "agent_network.guardrails.delete", "agent_network.policies.create", "agent_network.policies.update", "agent_network.policies.delete", "agent_network.providers.create", "agent_network.providers.update", "agent_network.providers.delete", "users.create", "users.update", "users.delete", "users.approve", "users.reject", "users.password.update", "users.invite.resend", "users.tokens.create", "users.tokens.delete", "setup_keys.create", "setup_keys.update", "setup_keys.delete", "users.invites.create", "users.invites.delete", "users.invites.regenerate", "users.invites.accept":
				var report analysis.ImpactReport
				var err error
				switch plan.Operation {
				case "groups.update":
					report, err = analysis.GroupUpdateImpact(plan.Before, plan.IntendedAfter)
				case "groups.create":
					report, err = analysis.GroupCreateImpact(plan.IntendedAfter)
				case "groups.delete":
					report, err = analysis.GroupDeleteImpact(plan.Before)
				case "policies.update":
					report, err = analysis.PolicyUpdateImpact(plan.Before, plan.IntendedAfter)
				case "policies.create":
					report, err = analysis.PolicyCreateImpact(plan.IntendedAfter)
				case "dns.zones.create":
					report, err = analysis.DNSZoneCreateImpact(plan.IntendedAfter)
				case "dns.zones.delete":
					report, err = analysis.DNSZoneDeleteImpact(plan.Before)
				case "dns.zones.update":
					report, err = analysis.DNSZoneUpdateImpact(plan.Before, plan.IntendedAfter)
				case "dns.records.create":
					report, err = analysis.DNSRecordCreateImpact(plan.IntendedAfter)
				case "dns.records.update":
					report, err = analysis.DNSRecordUpdateImpact(plan.Before, plan.IntendedAfter)
				case "dns.records.delete":
					report, err = analysis.DNSRecordDeleteImpact(plan.Before)
				case "dns.nameservers.create":
					report, err = analysis.DNSNameserverCreateImpact(plan.IntendedAfter)
				case "dns.nameservers.update":
					report, err = analysis.DNSNameserverUpdateImpact(plan.Before, plan.IntendedAfter)
				case "dns.nameservers.delete":
					report, err = analysis.DNSNameserverDeleteImpact(plan.Before)
				case "dns.settings.update":
					report, err = analysis.DNSSettingsUpdateImpact(plan.Before, plan.IntendedAfter)
				case "accounts.update":
					report, err = analysis.AccountUpdateImpact(plan.Before, plan.IntendedAfter)
				case "accounts.delete":
					report, err = analysis.AccountDeleteImpact(plan.Before)
				case "posture_checks.create":
					report, err = analysis.PostureCheckCreateImpact(plan.IntendedAfter)
				case "posture_checks.update":
					report, err = analysis.PostureCheckUpdateImpact(plan.Before, plan.IntendedAfter)
				case "posture_checks.delete":
					report, err = analysis.PostureCheckDeleteImpact(plan.Before)
				case "ingress.peers.create":
					report, err = analysis.IngressPeerCreateImpact(plan.IntendedAfter)
				case "ingress.peers.update":
					report, err = analysis.IngressPeerUpdateImpact(plan.Before, plan.IntendedAfter)
				case "ingress.peers.delete":
					report, err = analysis.IngressPeerDeleteImpact(plan.Before)
				case "peers.ingress.ports.create":
					report, err = analysis.IngressPortAllocationCreateImpact(plan.IntendedAfter)
				case "peers.ingress.ports.update":
					report, err = analysis.IngressPortAllocationUpdateImpact(plan.Before, plan.IntendedAfter)
				case "peers.ingress.ports.delete":
					report, err = analysis.IngressPortAllocationDeleteImpact(plan.Before)
				case "peers.edr.bypass.create":
					report, err = analysis.EDRBypassCreateImpact(plan.Before, plan.IntendedAfter)
				case "peers.edr.bypass.delete":
					report, err = analysis.EDRBypassDeleteImpact(plan.Before)
				case "agent_network.settings.update":
					report, err = analysis.AgentNetworkSettingsUpdateImpact(plan.Before, plan.IntendedAfter)
				case "agent_network.settings.create":
					report, err = analysis.AgentNetworkSettingsCreateImpact(plan.IntendedAfter)
				case "agent_network.settings.delete":
					report, err = analysis.AgentNetworkSettingsDeleteImpact(plan.Before)
				case "agent_network.budget_rules.create":
					report, err = analysis.AgentNetworkBudgetRuleCreateImpact(plan.IntendedAfter)
				case "agent_network.budget_rules.update":
					report, err = analysis.AgentNetworkBudgetRuleUpdateImpact(plan.Before, plan.IntendedAfter)
				case "agent_network.budget_rules.delete":
					report, err = analysis.AgentNetworkBudgetRuleDeleteImpact(plan.Before)
				case "agent_network.guardrails.create":
					report, err = analysis.AgentNetworkGuardrailCreateImpact(plan.IntendedAfter)
				case "agent_network.guardrails.update":
					report, err = analysis.AgentNetworkGuardrailUpdateImpact(plan.Before, plan.IntendedAfter)
				case "agent_network.guardrails.delete":
					report, err = analysis.AgentNetworkGuardrailDeleteImpact(plan.Before)
				case "agent_network.policies.create":
					report, err = analysis.AgentNetworkPolicyCreateImpact(plan.IntendedAfter)
				case "agent_network.policies.update":
					report, err = analysis.AgentNetworkPolicyUpdateImpact(plan.Before, plan.IntendedAfter)
				case "agent_network.policies.delete":
					report, err = analysis.AgentNetworkPolicyDeleteImpact(plan.Before)
				case "agent_network.providers.create":
					report, err = analysis.AgentNetworkProviderCreateImpact(plan.IntendedAfter)
				case "agent_network.providers.update":
					report, err = analysis.AgentNetworkProviderUpdateImpact(plan.Before, plan.IntendedAfter)
				case "agent_network.providers.delete":
					report, err = analysis.AgentNetworkProviderDeleteImpact(plan.Before)
				case "users.create":
					report, err = analysis.UserCreateImpact(plan.IntendedAfter)
				case "users.update":
					report, err = analysis.UserUpdateImpact(plan.Before, plan.IntendedAfter)
				case "users.delete":
					report, err = analysis.UserDeleteImpact(plan.Before)
				case "users.approve":
					report, err = analysis.UserApproveImpact(plan.Before)
				case "users.reject":
					report, err = analysis.UserRejectImpact(plan.Before)
				case "users.password.update":
					report, err = analysis.UserPasswordUpdateImpact(plan.Before, plan.IntendedAfter)
				case "users.invite.resend":
					report, err = analysis.UserInviteResendImpact(plan.Before, plan.IntendedAfter)
				case "users.tokens.delete":
					report, err = analysis.UserTokenDeleteImpact(plan.Before)
				case "users.tokens.create":
					report, err = analysis.UserTokenCreateImpact(plan.IntendedAfter)
				case "setup_keys.delete":
					report, err = analysis.SetupKeyDeleteImpact(plan.Before)
				case "setup_keys.create":
					report, err = analysis.SetupKeyCreateImpact(plan.IntendedAfter)
				case "setup_keys.update":
					report, err = analysis.SetupKeyUpdateImpact(plan.Before, plan.IntendedAfter)
				case "users.invites.create":
					report, err = analysis.InviteCreateImpact(plan.IntendedAfter)
				case "users.invites.delete":
					report, err = analysis.InviteDeleteImpact(plan.Before)
				case "users.invites.regenerate":
					report, err = analysis.InviteRegenerateImpact(plan.Before, plan.IntendedAfter)
				case "users.invites.accept":
					report, err = analysis.InviteAcceptImpact(plan.Before, plan.IntendedAfter)
				case "policies.delete":
					report, err = analysis.PolicyDeleteImpact(plan.Before)
				case "routes.update":
					report, err = analysis.RouteUpdateImpact(plan.Before, plan.IntendedAfter)
				case "routes.create":
					report, err = analysis.RouteCreateImpact(plan.IntendedAfter)
				case "routes.delete":
					report, err = analysis.RouteDeleteImpact(plan.Before)
				case "peers.update":
					report, err = analysis.PeerUpdateImpact(plan.Before, plan.IntendedAfter)
				case "peers.delete":
					report, err = analysis.PeerDeleteImpact(plan.Before)
				case "peers.temporary_access.create":
					report, err = analysis.TemporaryAccessCreateImpact(plan.Before, plan.IntendedAfter)
				case "peers.jobs.create":
					report, err = analysis.PeerJobCreateImpact(plan.Before, plan.IntendedAfter)
				case "event_streaming.create":
					report, err = analysis.EventStreamingCreateImpact(plan.IntendedAfter)
				case "event_streaming.update":
					report, err = analysis.EventStreamingUpdateImpact(plan.Before, plan.IntendedAfter)
				case "event_streaming.delete":
					report, err = analysis.EventStreamingDeleteImpact(plan.Before)
				case "identity_providers.create":
					report, err = analysis.IdentityProviderCreateImpact(plan.IntendedAfter)
				case "identity_providers.update":
					report, err = analysis.IdentityProviderUpdateImpact(plan.Before, plan.IntendedAfter)
				case "identity_providers.delete":
					report, err = analysis.IdentityProviderDeleteImpact(plan.Before)
				case "reverse_proxy_tokens.create":
					report, err = analysis.ReverseProxyTokenCreateImpact(plan.IntendedAfter)
				case "reverse_proxy_tokens.delete":
					report, err = analysis.ReverseProxyTokenDeleteImpact(plan.Before)
				case "reverse_proxy_domains.create":
					report, err = analysis.ReverseProxyDomainCreateImpact(plan.IntendedAfter)
				case "reverse_proxy_domains.delete":
					report, err = analysis.ReverseProxyDomainDeleteImpact(plan.Before)
				case "reverse_proxy_clusters.delete":
					report, err = analysis.ReverseProxyClusterDeleteImpact(plan.Before)
				case "reverse_proxy_services.create":
					report, err = analysis.ReverseProxyServiceCreateImpact(plan.IntendedAfter)
				case "reverse_proxy_services.update":
					report, err = analysis.ReverseProxyServiceUpdateImpact(plan.Before, plan.IntendedAfter)
				case "reverse_proxy_services.delete":
					report, err = analysis.ReverseProxyServiceDeleteImpact(plan.Before)
				case "notification_channels.create":
					report, err = analysis.NotificationChannelCreateImpact(plan.IntendedAfter)
				case "notification_channels.update":
					report, err = analysis.NotificationChannelUpdateImpact(plan.Before, plan.IntendedAfter)
				case "notification_channels.delete":
					report, err = analysis.NotificationChannelDeleteImpact(plan.Before)
				case "azure_idp.create":
					report, err = analysis.AzureIDPCreateImpact(plan.IntendedAfter)
				case "azure_idp.update":
					report, err = analysis.AzureIDPUpdateImpact(plan.Before, plan.IntendedAfter)
				case "azure_idp.delete":
					report, err = analysis.AzureIDPDeleteImpact(plan.Before)
				case "azure_idp.sync":
					report, err = analysis.AzureIDPSyncImpact(plan.Before, plan.IntendedAfter)
				case "google_idp.create":
					report, err = analysis.GoogleIDPCreateImpact(plan.IntendedAfter)
				case "google_idp.update":
					report, err = analysis.GoogleIDPUpdateImpact(plan.Before, plan.IntendedAfter)
				case "google_idp.delete":
					report, err = analysis.GoogleIDPDeleteImpact(plan.Before)
				case "google_idp.sync":
					report, err = analysis.GoogleIDPSyncImpact(plan.Before, plan.IntendedAfter)
				case "edr.intune.create", "edr.sentinelone.create", "edr.falcon.create", "edr.huntress.create", "edr.fleetdm.create":
					report, err = analysis.EDRIntegrationCreateImpact(strings.Split(plan.Operation, ".")[1], plan.IntendedAfter)
				case "edr.intune.update", "edr.sentinelone.update", "edr.falcon.update", "edr.huntress.update", "edr.fleetdm.update":
					report, err = analysis.EDRIntegrationUpdateImpact(strings.Split(plan.Operation, ".")[1], plan.Before, plan.IntendedAfter)
				case "edr.intune.delete", "edr.sentinelone.delete", "edr.falcon.delete", "edr.huntress.delete", "edr.fleetdm.delete":
					report, err = analysis.EDRIntegrationDeleteImpact(strings.Split(plan.Operation, ".")[1], plan.Before)
				case "scim.create", "okta_scim.create":
					report, err = analysis.SCIMCreateImpact(strings.Split(plan.Operation, ".")[0], plan.IntendedAfter)
				case "scim.update", "okta_scim.update":
					report, err = analysis.SCIMUpdateImpact(strings.Split(plan.Operation, ".")[0], plan.Before, plan.IntendedAfter)
				case "scim.delete", "okta_scim.delete":
					report, err = analysis.SCIMDeleteImpact(strings.Split(plan.Operation, ".")[0], plan.Before)
				case "scim.token", "okta_scim.token":
					report, err = analysis.SCIMTokenImpact(strings.Split(plan.Operation, ".")[0], plan.Before)
				case "msp.tenants.create":
					report, err = analysis.MSPTenantCreateImpact(plan.IntendedAfter)
				case "msp.tenants.update":
					report, err = analysis.MSPTenantUpdateImpact(plan.Before, plan.IntendedAfter)
				case "msp.tenants.dns":
					report, err = analysis.MSPTenantActionImpact("dns", plan.Before)
				case "msp.tenants.invite":
					report, err = analysis.MSPTenantActionImpact("invite", plan.Before)
				case "msp.tenants.invite.respond":
					report, err = analysis.MSPTenantActionImpact("invite_respond", plan.Before)
				case "msp.tenants.subscription":
					report, err = analysis.MSPTenantActionImpact("subscription", plan.Before)
				case "msp.tenants.unlink":
					report, err = analysis.MSPTenantActionImpact("unlink", plan.Before)
				case "billing.aws_marketplace.activate":
					report, err = analysis.BillingAWSMarketplaceImpact("activate", plan.Before, plan.IntendedAfter)
				case "billing.aws_marketplace.enrich":
					report, err = analysis.BillingAWSMarketplaceImpact("enrich", plan.Before, plan.IntendedAfter)
				case "billing.checkout.create":
					report, err = analysis.BillingCheckoutImpact(plan.IntendedAfter)
				case "billing.subscription.update":
					report, err = analysis.BillingSubscriptionUpdateImpact(plan.Before, plan.IntendedAfter)
				case "networks.update":
					report, err = analysis.NetworkUpdateImpact(plan.Before, plan.IntendedAfter)
				case "networks.create":
					report, err = analysis.NetworkCreateImpact(plan.IntendedAfter)
				case "networks.delete":
					report, err = analysis.NetworkDeleteImpact(plan.Before)
				case "networks.resources.update":
					report, err = analysis.NetworkResourceUpdateImpact(plan.Before, plan.IntendedAfter)
				case "networks.resources.create":
					report, err = analysis.NetworkResourceCreateImpact(plan.IntendedAfter)
				case "networks.routers.create":
					report, err = analysis.NetworkRouterCreateImpact(plan.IntendedAfter)
				case "networks.routers.update":
					report, err = analysis.NetworkRouterUpdateImpact(plan.Before, plan.IntendedAfter)
				case "networks.resources.delete":
					report, err = analysis.NetworkResourceDeleteImpact(plan.Before)
				case "networks.routers.delete":
					report, err = analysis.NetworkRouterDeleteImpact(plan.Before)
				}
				if err != nil {
					return fail(2, err)
				}
				impact, err = json.Marshal(report)
				if err != nil {
					return fail(1, fmt.Errorf("encode mutation impact: %w", err))
				}
				findingCode := ""
				findingMessage := ""
				switch {
				case plan.Operation == "groups.update" && report.Classification == "unknown":
					findingCode = "impact.unknown"
					findingMessage = "the proposed group change may affect reachability, but its impact cannot be calculated"
				case plan.Operation == "groups.create" && report.Classification == "group_create":
					findingCode = "impact.group_create"
					findingMessage = "creating the group may alter policy or resource membership and requires exact acknowledgement"
				case plan.Operation == "groups.delete" && report.Classification == "group_delete":
					findingCode = "impact.group_delete"
					findingMessage = "deleting the group may alter policy membership and requires exact acknowledgement"
				case plan.Operation == "policies.update" && report.Classification == "policy_rule_change":
					findingCode = "impact.policy_rule_change"
					findingMessage = "the proposed policy rule change may alter reachability and requires exact acknowledgement"
				case plan.Operation == "policies.create" && report.Classification == "policy_create":
					findingCode = "impact.policy_create"
					findingMessage = "creating the policy may add reachability and requires exact acknowledgement"
				case plan.Operation == "dns.zones.create" && report.Classification == "dns_zone_create":
					findingCode = "impact.dns_zone_create"
					findingMessage = "creating the DNS zone may alter name resolution and requires exact acknowledgement"
				case plan.Operation == "dns.zones.delete" && report.Classification == "dns_zone_delete":
					findingCode = "impact.dns_zone_delete"
					findingMessage = "deleting the DNS zone may alter name resolution and requires exact acknowledgement"
				case plan.Operation == "dns.zones.update" && report.Classification == "dns_zone_change":
					findingCode = "impact.dns_zone_change"
					findingMessage = "the proposed DNS zone change may alter name resolution and requires exact acknowledgement"
				case plan.Operation == "dns.records.create" && report.Classification == "dns_record_create":
					findingCode = "impact.dns_record_create"
					findingMessage = "creating the DNS record may alter name resolution and requires exact acknowledgement"
				case plan.Operation == "dns.records.update" && report.Classification == "dns_record_change":
					findingCode = "impact.dns_record_change"
					findingMessage = "the proposed DNS record change may alter name resolution and requires exact acknowledgement"
				case plan.Operation == "dns.records.delete" && report.Classification == "dns_record_delete":
					findingCode = "impact.dns_record_delete"
					findingMessage = "deleting the DNS record may alter name resolution and requires exact acknowledgement"
				case plan.Operation == "dns.nameservers.create" && report.Classification == "dns_nameserver_create":
					findingCode = "impact.dns_nameserver_create"
					findingMessage = "creating the nameserver group may alter resolver behavior and requires exact acknowledgement"
				case plan.Operation == "dns.nameservers.update" && report.Classification == "dns_nameserver_change":
					findingCode = "impact.dns_nameserver_change"
					findingMessage = "the proposed nameserver group change may alter resolver behavior and requires exact acknowledgement"
				case plan.Operation == "dns.nameservers.delete" && report.Classification == "dns_nameserver_delete":
					findingCode = "impact.dns_nameserver_delete"
					findingMessage = "deleting the nameserver group may alter resolver behavior and requires exact acknowledgement"
				case plan.Operation == "dns.settings.update" && report.Classification == "dns_settings_change":
					findingCode = "impact.dns_settings_change"
					findingMessage = "the proposed DNS settings change may alter resolver behavior and requires exact acknowledgement"
				case plan.Operation == "accounts.update" && report.Classification == "account_change":
					findingCode = "impact.account_change"
					findingMessage = "the proposed account change may alter management-plane behavior and requires exact acknowledgement"
				case plan.Operation == "accounts.delete" && report.Classification == "account_delete":
					findingCode = "impact.account_delete"
					findingMessage = "deleting the account may remove management access and resources and requires exact acknowledgement"
				case plan.Operation == "posture_checks.create" && report.Classification == "posture_check_create":
					findingCode = "impact.posture_check_create"
					findingMessage = "creating the posture check may alter policy admission and requires exact acknowledgement"
				case plan.Operation == "posture_checks.update" && report.Classification == "posture_check_change":
					findingCode = "impact.posture_check_change"
					findingMessage = "the proposed posture check change may alter policy admission and requires exact acknowledgement"
				case plan.Operation == "posture_checks.delete" && report.Classification == "posture_check_delete":
					findingCode = "impact.posture_check_delete"
					findingMessage = "deleting the posture check may alter policy admission and requires exact acknowledgement"
				case plan.Operation == "ingress.peers.create" && report.Classification == "ingress_peer_create":
					findingCode = "impact.ingress_peer_create"
					findingMessage = "creating the ingress peer may add external reachability and requires exact acknowledgement"
				case plan.Operation == "ingress.peers.update" && report.Classification == "ingress_peer_change":
					findingCode = "impact.ingress_peer_change"
					findingMessage = "the proposed ingress peer change may alter external reachability and requires exact acknowledgement"
				case plan.Operation == "ingress.peers.delete" && report.Classification == "ingress_peer_delete":
					findingCode = "impact.ingress_peer_delete"
					findingMessage = "deleting the ingress peer may remove external reachability and requires exact acknowledgement"
				case plan.Operation == "agent_network.settings.update" && report.Classification == "agent_network_settings_change":
					findingCode = "impact.agent_network_settings_change"
					findingMessage = "the proposed agent-network settings change may alter Cloud agent behavior and requires exact acknowledgement"
				case plan.Operation == "agent_network.settings.create" && report.Classification == "agent_network_settings_create":
					findingCode = "impact.agent_network_settings_create"
					findingMessage = "creating agent-network settings may enable Cloud agent behavior and requires exact acknowledgement"
				case plan.Operation == "agent_network.settings.delete" && report.Classification == "agent_network_settings_delete":
					findingCode = "impact.agent_network_settings_delete"
					findingMessage = "deleting agent-network settings may disable Cloud agent behavior and requires exact acknowledgement"
				case plan.Operation == "agent_network.budget_rules.create" && report.Classification == "agent_network_budget_rule_create":
					findingCode = "impact.agent_network_budget_rule_create"
					findingMessage = "creating the agent-network budget rule may change spend and token admission and requires exact acknowledgement"
				case plan.Operation == "agent_network.budget_rules.update" && report.Classification == "agent_network_budget_rule_change":
					findingCode = "impact.agent_network_budget_rule_change"
					findingMessage = "the proposed agent-network budget rule change may change spend and token admission and requires exact acknowledgement"
				case plan.Operation == "agent_network.budget_rules.delete" && report.Classification == "agent_network_budget_rule_delete":
					findingCode = "impact.agent_network_budget_rule_delete"
					findingMessage = "deleting the agent-network budget rule may remove spend and token limits and requires exact acknowledgement"
				case plan.Operation == "agent_network.guardrails.create" && report.Classification == "agent_network_guardrail_create":
					findingCode = "impact.agent_network_guardrail_create"
					findingMessage = "creating the agent-network guardrail may alter policy admission and requires exact acknowledgement"
				case plan.Operation == "agent_network.guardrails.update" && report.Classification == "agent_network_guardrail_change":
					findingCode = "impact.agent_network_guardrail_change"
					findingMessage = "the proposed agent-network guardrail change may alter policy admission and requires exact acknowledgement"
				case plan.Operation == "agent_network.guardrails.delete" && report.Classification == "agent_network_guardrail_delete":
					findingCode = "impact.agent_network_guardrail_delete"
					findingMessage = "deleting the agent-network guardrail may remove policy admission controls and requires exact acknowledgement"
				case plan.Operation == "agent_network.policies.create" && report.Classification == "agent_network_policy_create":
					findingCode = "impact.agent_network_policy_create"
					findingMessage = "creating the agent-network policy may add provider reachability and requires exact acknowledgement"
				case plan.Operation == "agent_network.policies.update" && report.Classification == "agent_network_policy_change":
					findingCode = "impact.agent_network_policy_change"
					findingMessage = "the proposed agent-network policy change may alter provider reachability and requires exact acknowledgement"
				case plan.Operation == "agent_network.policies.delete" && report.Classification == "agent_network_policy_delete":
					findingCode = "impact.agent_network_policy_delete"
					findingMessage = "deleting the agent-network policy may remove provider reachability and requires exact acknowledgement"
				case plan.Operation == "agent_network.providers.create" && report.Classification == "agent_network_provider_create":
					findingCode = "impact.agent_network_provider_create"
					findingMessage = "creating the agent-network provider may add upstream reachability and requires exact acknowledgement"
				case plan.Operation == "agent_network.providers.update" && report.Classification == "agent_network_provider_change":
					findingCode = "impact.agent_network_provider_change"
					findingMessage = "the proposed agent-network provider change may alter upstream reachability and requires exact acknowledgement"
				case plan.Operation == "agent_network.providers.delete" && report.Classification == "agent_network_provider_delete":
					findingCode = "impact.agent_network_provider_delete"
					findingMessage = "deleting the agent-network provider may remove upstream reachability and requires exact acknowledgement"
				case plan.Operation == "users.create" && report.Classification == "user_create":
					findingCode = "impact.user_create"
					findingMessage = "creating the user may grant account access and requires exact acknowledgement"
				case plan.Operation == "users.update" && report.Classification == "user_change":
					findingCode = "impact.user_change"
					findingMessage = "the proposed user change may alter account access or peer assignment and requires exact acknowledgement"
				case plan.Operation == "users.delete" && report.Classification == "user_delete":
					findingCode = "impact.user_delete"
					findingMessage = "deleting the user may revoke account access and requires exact acknowledgement"
				case plan.Operation == "users.approve" && report.Classification == "user_approve":
					findingCode = "impact.user_approve"
					findingMessage = "approving the user may grant account access and requires exact acknowledgement"
				case plan.Operation == "users.reject" && report.Classification == "user_reject":
					findingCode = "impact.user_reject"
					findingMessage = "rejecting the user removes a pending account-access edge and requires exact acknowledgement"
				case plan.Operation == "users.password.update" && report.Classification == "user_password_change":
					findingCode = "impact.user_password_change"
					findingMessage = "changing the user password changes an authentication credential and requires exact acknowledgement"
				case plan.Operation == "users.invite.resend" && report.Classification == "user_invite_resend":
					findingCode = "impact.user_invite_resend"
					findingMessage = "resending the user invitation triggers external enrollment delivery and requires exact acknowledgement"
				case plan.Operation == "peers.ingress.ports.create" && report.Classification == "ingress_port_allocation_create":
					findingCode = "impact.ingress_port_allocation_create"
					findingMessage = "creating the ingress port allocation may expose peer services externally and requires exact acknowledgement"
				case plan.Operation == "peers.ingress.ports.update" && report.Classification == "ingress_port_allocation_change":
					findingCode = "impact.ingress_port_allocation_change"
					findingMessage = "the proposed ingress port allocation change may alter external peer exposure and requires exact acknowledgement"
				case plan.Operation == "peers.ingress.ports.delete" && report.Classification == "ingress_port_allocation_delete":
					findingCode = "impact.ingress_port_allocation_delete"
					findingMessage = "deleting the ingress port allocation may remove external peer exposure and requires exact acknowledgement"
				case plan.Operation == "peers.edr.bypass.create" && report.Classification == "edr_bypass_create":
					findingCode = "impact.edr_bypass_create"
					findingMessage = "bypassing EDR compliance immediately grants peer network access and requires exact acknowledgement"
				case plan.Operation == "peers.edr.bypass.delete" && report.Classification == "edr_bypass_delete":
					findingCode = "impact.edr_bypass_delete"
					findingMessage = "revoking the EDR bypass restores compliance gating and may remove peer network access; exact acknowledgement is required"
				case plan.Operation == "users.tokens.delete" && report.Classification == "user_token_delete":
					findingCode = "impact.user_token_delete"
					findingMessage = "deleting the personal access token revokes a credential and requires exact acknowledgement"
				case plan.Operation == "users.tokens.create" && report.Classification == "user_token_create":
					findingCode = "impact.user_token_create"
					findingMessage = "creating the personal access token returns a one-time secret and requires exact acknowledgement"
				case plan.Operation == "setup_keys.delete" && report.Classification == "setup_key_delete":
					findingCode = "impact.setup_key_delete"
					findingMessage = "deleting the setup key may stop new peer enrollment and requires exact acknowledgement"
				case plan.Operation == "setup_keys.create" && report.Classification == "setup_key_create":
					findingCode = "impact.setup_key_create"
					findingMessage = "creating the setup key expands peer enrollment authority and requires exact acknowledgement"
				case plan.Operation == "setup_keys.update" && report.Classification == "setup_key_change":
					findingCode = "impact.setup_key_change"
					findingMessage = "changing the setup key may alter peer enrollment authority and requires exact acknowledgement"
				case plan.Operation == "users.invites.create" && report.Classification == "invite_create":
					findingCode = "impact.invite_create"
					findingMessage = "creating the invite expands enrollment authority and returns a one-time token; exact acknowledgement is required"
				case plan.Operation == "users.invites.delete" && report.Classification == "invite_delete":
					findingCode = "impact.invite_delete"
					findingMessage = "deleting the invite removes a pending enrollment edge and requires exact acknowledgement"
				case plan.Operation == "users.invites.regenerate" && report.Classification == "invite_regenerate":
					findingCode = "impact.invite_regenerate"
					findingMessage = "regenerating the invite invalidates its previous token and returns a new one-time token; exact acknowledgement is required"
				case plan.Operation == "users.invites.accept" && report.Classification == "invite_accept":
					findingCode = "impact.invite_accept"
					findingMessage = "accepting the invite creates account access from a public token and requires exact acknowledgement"
				case plan.Operation == "policies.delete" && report.Classification == "policy_delete":
					findingCode = "impact.policy_delete"
					findingMessage = "deleting the policy may remove access and requires exact acknowledgement"
				case plan.Operation == "routes.update" && report.Classification == "route_change":
					findingCode = "impact.route_change"
					findingMessage = "the proposed route change may alter reachability and requires exact acknowledgement"
				case plan.Operation == "routes.create" && report.Classification == "route_create":
					findingCode = "impact.route_create"
					findingMessage = "creating the route may alter reachability and requires exact acknowledgement"
				case plan.Operation == "routes.delete" && report.Classification == "route_delete":
					findingCode = "impact.route_delete"
					findingMessage = "deleting the route may alter reachability and requires exact acknowledgement"
				case plan.Operation == "peers.update" && report.Classification == "peer_change":
					findingCode = "impact.peer_change"
					findingMessage = "the proposed peer change may alter access or connectivity and requires exact acknowledgement"
				case plan.Operation == "peers.delete" && report.Classification == "peer_delete":
					findingCode = "impact.peer_delete"
					findingMessage = "deleting the peer may remove access or connectivity and requires exact acknowledgement"
				case plan.Operation == "peers.temporary_access.create" && report.Classification == "temporary_access_create":
					findingCode = "impact.temporary_access_create"
					findingMessage = "creating a temporary access peer grants scoped network access and requires exact acknowledgement"
				case plan.Operation == "peers.jobs.create" && report.Classification == "peer_job_create":
					findingCode = "impact.peer_job_create"
					findingMessage = "creating the remote job executes a workload on the peer and may collect sensitive diagnostics; exact acknowledgement is required"
				case plan.Operation == "event_streaming.create" && report.Classification == "event_streaming_create":
					findingCode = "impact.event_streaming_create"
					findingMessage = "creating the event-streaming integration exports account activity and requires exact acknowledgement"
				case plan.Operation == "event_streaming.update" && report.Classification == "event_streaming_change":
					findingCode = "impact.event_streaming_change"
					findingMessage = "changing the event-streaming integration may alter external activity delivery and requires exact acknowledgement"
				case plan.Operation == "event_streaming.delete" && report.Classification == "event_streaming_delete":
					findingCode = "impact.event_streaming_delete"
					findingMessage = "deleting the event-streaming integration stops external activity delivery and requires exact acknowledgement"
				case plan.Operation == "identity_providers.create" && report.Classification == "identity_provider_create":
					findingCode = "impact.identity_provider_create"
					findingMessage = "creating the identity provider changes authentication ingress and requires exact acknowledgement"
				case plan.Operation == "identity_providers.update" && report.Classification == "identity_provider_change":
					findingCode = "impact.identity_provider_change"
					findingMessage = "changing the identity provider may alter authentication behavior and requires exact acknowledgement"
				case plan.Operation == "identity_providers.delete" && report.Classification == "identity_provider_delete":
					findingCode = "impact.identity_provider_delete"
					findingMessage = "deleting the identity provider may strand users at authentication and requires exact acknowledgement"
				case plan.Operation == "reverse_proxy_tokens.create" && report.Classification == "reverse_proxy_token_create":
					findingCode = "impact.reverse_proxy_token_create"
					findingMessage = "creating the reverse proxy token grants an external proxy credential and requires exact acknowledgement"
				case plan.Operation == "reverse_proxy_tokens.delete" && report.Classification == "reverse_proxy_token_delete":
					findingCode = "impact.reverse_proxy_token_delete"
					findingMessage = "deleting the reverse proxy token disconnects its proxies and requires exact acknowledgement"
				case plan.Operation == "reverse_proxy_domains.create" && report.Classification == "reverse_proxy_domain_create":
					findingCode = "impact.reverse_proxy_domain_create"
					findingMessage = "creating the reverse proxy domain changes public DNS and ingress exposure and requires exact acknowledgement"
				case plan.Operation == "reverse_proxy_domains.delete" && report.Classification == "reverse_proxy_domain_delete":
					findingCode = "impact.reverse_proxy_domain_delete"
					findingMessage = "deleting the reverse proxy domain removes public DNS and ingress exposure and requires exact acknowledgement"
				case plan.Operation == "reverse_proxy_clusters.delete" && report.Classification == "reverse_proxy_cluster_delete":
					findingCode = "impact.reverse_proxy_cluster_delete"
					findingMessage = "deleting the reverse proxy cluster removes public ingress infrastructure and requires exact acknowledgement"
				case plan.Operation == "reverse_proxy_services.create" && report.Classification == "reverse_proxy_service_create":
					findingCode = "impact.reverse_proxy_service_create"
					findingMessage = "creating the reverse proxy service publishes public ingress to internal targets and requires exact acknowledgement"
				case plan.Operation == "reverse_proxy_services.update" && report.Classification == "reverse_proxy_service_change":
					findingCode = "impact.reverse_proxy_service_change"
					findingMessage = "changing the reverse proxy service may alter public routing or authentication and requires exact acknowledgement"
				case plan.Operation == "reverse_proxy_services.delete" && report.Classification == "reverse_proxy_service_delete":
					findingCode = "impact.reverse_proxy_service_delete"
					findingMessage = "deleting the reverse proxy service removes public ingress and requires exact acknowledgement"
				case plan.Operation == "notification_channels.create" && report.Classification == "notification_channel_create":
					findingCode = "impact.notification_channel_create"
					findingMessage = "creating the notification channel changes external account delivery and requires exact acknowledgement"
				case plan.Operation == "notification_channels.update" && report.Classification == "notification_channel_change":
					findingCode = "impact.notification_channel_change"
					findingMessage = "changing the notification channel may alter external account delivery and requires exact acknowledgement"
				case plan.Operation == "notification_channels.delete" && report.Classification == "notification_channel_delete":
					findingCode = "impact.notification_channel_delete"
					findingMessage = "deleting the notification channel stops external account delivery and requires exact acknowledgement"
				case plan.Operation == "azure_idp.create" && report.Classification == "azure_idp_create":
					findingCode = "impact.azure_idp_create"
					findingMessage = "creating the Azure identity integration changes external authentication and directory synchronization and requires exact acknowledgement"
				case plan.Operation == "azure_idp.update" && report.Classification == "azure_idp_change":
					findingCode = "impact.azure_idp_change"
					findingMessage = "changing the Azure identity integration may alter authentication or directory synchronization and requires exact acknowledgement"
				case plan.Operation == "azure_idp.delete" && report.Classification == "azure_idp_delete":
					findingCode = "impact.azure_idp_delete"
					findingMessage = "deleting the Azure identity integration may strand users or stop directory synchronization and requires exact acknowledgement"
				case plan.Operation == "azure_idp.sync" && report.Classification == "azure_idp_sync":
					findingCode = "impact.azure_idp_sync"
					findingMessage = "triggering Azure directory synchronization may create or update account users and groups and requires exact acknowledgement"
				case plan.Operation == "google_idp.create" && report.Classification == "google_idp_create":
					findingCode = "impact.google_idp_create"
					findingMessage = "creating the Google identity integration changes external authentication and directory synchronization and requires exact acknowledgement"
				case plan.Operation == "google_idp.update" && report.Classification == "google_idp_change":
					findingCode = "impact.google_idp_change"
					findingMessage = "changing the Google identity integration may alter authentication or directory synchronization and requires exact acknowledgement"
				case plan.Operation == "google_idp.delete" && report.Classification == "google_idp_delete":
					findingCode = "impact.google_idp_delete"
					findingMessage = "deleting the Google identity integration may strand users or stop directory synchronization and requires exact acknowledgement"
				case plan.Operation == "google_idp.sync" && report.Classification == "google_idp_sync":
					findingCode = "impact.google_idp_sync"
					findingMessage = "triggering Google directory synchronization may create or update account users and groups and requires exact acknowledgement"
				case strings.HasPrefix(plan.Operation, "edr.") && (strings.HasSuffix(plan.Operation, ".create") || strings.HasSuffix(plan.Operation, ".update") || strings.HasSuffix(plan.Operation, ".delete")):
					parts := strings.Split(plan.Operation, ".")
					if len(parts) == 3 {
						kind := parts[2]
						if kind == "update" {
							kind = "change"
						}
						if report.Classification == "edr_"+parts[1]+"_"+kind {
							findingCode = "impact.edr_" + parts[1] + "_" + kind
							switch kind {
							case "create":
								findingMessage = "creating the EDR integration changes device-compliance enforcement and peer access and requires exact acknowledgement"
							case "change":
								findingMessage = "changing the EDR integration may alter device-compliance enforcement and peer access and requires exact acknowledgement"
							case "delete":
								findingMessage = "deleting the EDR integration removes a device-compliance gate and requires exact acknowledgement"
							}
						}
					}
				case strings.HasPrefix(plan.Operation, "scim.") || strings.HasPrefix(plan.Operation, "okta_scim."):
					parts := strings.Split(plan.Operation, ".")
					if len(parts) == 2 {
						kind := parts[1]
						classificationKind := kind
						if kind == "update" {
							classificationKind = "change"
						}
						findingCode = "impact.scim_" + parts[0] + "_" + classificationKind
						if report.Classification == "scim_"+parts[0]+"_"+classificationKind {
							switch kind {
							case "create":
								findingMessage = "creating the SCIM integration changes external user and group provisioning and requires exact acknowledgement"
							case "update":
								findingCode = "impact.scim_" + parts[0] + "_change"
								findingMessage = "changing the SCIM integration may alter external user and group provisioning and requires exact acknowledgement"
							case "delete":
								findingMessage = "deleting the SCIM integration stops external user and group provisioning and requires exact acknowledgement"
							case "token":
								findingMessage = "regenerating the SCIM token revokes the prior provisioning credential and requires exact acknowledgement"
							}
						}
					}
				case plan.Operation == "msp.tenants.create" && report.Classification == "msp_tenant_create":
					findingCode = "impact.msp_tenant_create"
					findingMessage = "creating the MSP tenant provisions a customer account and requires exact acknowledgement"
				case plan.Operation == "msp.tenants.update" && report.Classification == "msp_tenant_change":
					findingCode = "impact.msp_tenant_change"
					findingMessage = "changing the MSP tenant alters delegated access or customer-account metadata and requires exact acknowledgement"
				case plan.Operation == "msp.tenants.dns" && report.Classification == "msp_tenant_dns":
					findingCode = "impact.msp_tenant_dns"
					findingMessage = "verifying the MSP tenant DNS challenge changes customer-account activation state and requires exact acknowledgement"
				case plan.Operation == "msp.tenants.invite" && report.Classification == "msp_tenant_invite":
					findingCode = "impact.msp_tenant_invite"
					findingMessage = "inviting the MSP tenant changes external customer-account delivery and requires exact acknowledgement"
				case plan.Operation == "msp.tenants.invite.respond" && report.Classification == "msp_tenant_invite_respond":
					findingCode = "impact.msp_tenant_invite_respond"
					findingMessage = "responding to the MSP tenant invitation changes account delegation and requires exact acknowledgement"
				case plan.Operation == "msp.tenants.subscription" && report.Classification == "msp_tenant_subscription":
					findingCode = "impact.msp_tenant_subscription"
					findingMessage = "changing the MSP tenant subscription changes billing entitlement and requires exact acknowledgement"
				case plan.Operation == "msp.tenants.unlink" && report.Classification == "msp_tenant_unlink":
					findingCode = "impact.msp_tenant_unlink"
					findingMessage = "unlinking the MSP tenant removes delegated management and requires exact acknowledgement"
				case plan.Operation == "billing.aws_marketplace.activate" && report.Classification == "billing_aws_marketplace_activate":
					findingCode = "impact.billing_aws_marketplace_activate"
					findingMessage = "activating AWS Marketplace billing changes external subscription entitlement and requires exact acknowledgement"
				case plan.Operation == "billing.aws_marketplace.enrich" && report.Classification == "billing_aws_marketplace_enrich":
					findingCode = "impact.billing_aws_marketplace_enrich"
					findingMessage = "enriching AWS Marketplace billing changes external subscription linkage and requires exact acknowledgement"
				case plan.Operation == "billing.checkout.create" && report.Classification == "billing_checkout_create":
					findingCode = "impact.billing_checkout_create"
					findingMessage = "creating a payment checkout session starts an external billing flow and requires exact acknowledgement"
				case plan.Operation == "billing.subscription.update" && report.Classification == "billing_subscription_change":
					findingCode = "impact.billing_subscription_change"
					findingMessage = "changing the billing subscription alters external entitlement and requires exact acknowledgement"
				case plan.Operation == "networks.update" && report.Classification == "network_change":
					findingCode = "impact.network_change"
					findingMessage = "the proposed network change may alter topology and requires exact acknowledgement"
				case plan.Operation == "networks.create" && report.Classification == "network_create":
					findingCode = "impact.network_create"
					findingMessage = "creating the network may add topology and requires exact acknowledgement"
				case plan.Operation == "networks.delete" && report.Classification == "network_delete":
					findingCode = "impact.network_delete"
					findingMessage = "deleting the network may remove attached topology and requires exact acknowledgement"
				case plan.Operation == "networks.resources.update" && report.Classification == "network_resource_change":
					findingCode = "impact.network_resource_change"
					findingMessage = "the proposed network resource change may alter reachability and requires exact acknowledgement"
				case plan.Operation == "networks.resources.create" && report.Classification == "network_resource_create":
					findingCode = "impact.network_resource_create"
					findingMessage = "creating the network resource may add reachability and requires exact acknowledgement"
				case plan.Operation == "networks.routers.create" && report.Classification == "network_router_create":
					findingCode = "impact.network_router_create"
					findingMessage = "creating the network router may alter reachability and requires exact acknowledgement"
				case plan.Operation == "networks.routers.update" && report.Classification == "network_router_change":
					findingCode = "impact.network_router_change"
					findingMessage = "the proposed network router change may alter reachability and requires exact acknowledgement"
				case plan.Operation == "networks.resources.delete" && report.Classification == "network_resource_delete":
					findingCode = "impact.network_resource_delete"
					findingMessage = "deleting the network resource may remove reachability and requires exact acknowledgement"
				case plan.Operation == "networks.routers.delete" && report.Classification == "network_router_delete":
					findingCode = "impact.network_router_delete"
					findingMessage = "deleting the network router may alter reachability and requires exact acknowledgement"
				}
				if findingCode != "" && !hasFinding(findings, findingCode) {
					findings = append(findings, ledger.Finding{Code: findingCode, Severity: "blocking", Message: findingMessage})
				}
			}
			definition, err := operations.Lookup(plan.Operation)
			if err != nil {
				return fail(2, err)
			}
			store, err := ledger.Open(state.statePath)
			if err != nil {
				return fail(3, err)
			}
			defer store.Close()
			stage, err := store.Create(cmd.Context(), ledger.StageInput{
				Profile:        state.profileName,
				ServerIdentity: profile.ServerIdentity,
				AccountID:      profile.AccountID,
				Operation:      plan.Operation,
				Request:        plan.Request,
				Before:         plan.Before,
				IntendedAfter:  plan.IntendedAfter,
				Impact:         impact,
				Findings:       findings,
			})
			if err != nil {
				return fail(2, err)
			}
			data := map[string]any{"stage_id": stage.ID, "revision": stage.Revision, "digest": stage.Digest, "applicability": "local_stage_only", "operation": stage.Operation, "mutation": definition.Mutation, "availability": definition.Availability, "dispatcher_admitted": definition.DispatcherAdmitted, "impact": json.RawMessage(stage.Impact), "findings": stage.Findings}
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/stage-result", "ok": true, "operation": "stage.create", "data": data})
			}
			_, err = fmt.Fprintf(stdout, "staged %s@%d\ndigest: %s\napplicability: local_stage_only\n", stage.ID, stage.Revision, stage.Digest)
			return err
		},
	}
	command.Flags().BoolVar(&fromJSON, "from-json", false, "read a structured stage plan from stdin")
	return command
}

func stageShowCommand(state *commandState, stdout io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "show <stage-id>@<revision>",
		Short: "show an exact immutable stage revision",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, revision, err := parseRevision(args[0])
			if err != nil {
				return fail(2, err)
			}
			store, err := ledger.Open(state.statePath)
			if err != nil {
				return fail(3, err)
			}
			defer store.Close()
			stage, err := store.Get(cmd.Context(), id, revision)
			if err != nil {
				return fail(2, err)
			}
			definition, err := operations.Lookup(stage.Operation)
			if err != nil {
				return fail(2, err)
			}
			data := map[string]any{"stage_id": stage.ID, "revision": stage.Revision, "profile": stage.Profile, "server_identity": stage.ServerIdentity, "account_id": stage.AccountID, "operation": stage.Operation, "availability": definition.Availability, "request": json.RawMessage(stage.Request), "before": json.RawMessage(stage.Before), "intended_after": json.RawMessage(stage.IntendedAfter), "impact": json.RawMessage(stage.Impact), "digest": stage.Digest, "findings": stage.Findings, "cancelled": stage.Cancelled, "created_at": stage.CreatedAt}
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/stage-result", "ok": true, "operation": "stage.show", "data": data})
			}
			_, err = fmt.Fprintf(stdout, "stage %s@%d\noperation: %s\ndigest: %s\ncancelled: %t\n", stage.ID, stage.Revision, stage.Operation, stage.Digest, stage.Cancelled)
			return err
		},
	}
}

func hasFinding(findings []ledger.Finding, code string) bool {
	for _, finding := range findings {
		if finding.Code == code {
			return true
		}
	}
	return false
}

func stageCancelCommand(state *commandState, stdout io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "cancel <stage-id>@<revision>",
		Short: "cancel an exact immutable stage revision",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, revision, err := parseRevision(args[0])
			if err != nil {
				return fail(2, err)
			}
			store, err := ledger.Open(state.statePath)
			if err != nil {
				return fail(3, err)
			}
			defer store.Close()
			if err := store.Cancel(cmd.Context(), id, revision); err != nil {
				return fail(2, err)
			}
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/stage-result", "ok": true, "operation": "stage.cancel", "data": map[string]any{"stage_id": id, "revision": revision, "cancelled": true}})
			}
			_, err = fmt.Fprintf(stdout, "cancelled %s@%d\n", id, revision)
			return err
		},
	}
}

func parseRevision(value string) (string, int, error) {
	separator := strings.LastIndexByte(value, '@')
	if separator <= 0 || separator == len(value)-1 {
		return "", 0, fmt.Errorf("revision must use exact <stage-id>@<revision> syntax")
	}
	revision, err := strconv.Atoi(value[separator+1:])
	if err != nil || revision < 1 {
		return "", 0, fmt.Errorf("revision must be a positive integer")
	}
	return value[:separator], revision, nil
}
