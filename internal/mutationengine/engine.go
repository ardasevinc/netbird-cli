package mutationengine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/ardasevinc/netbird-cli/internal/analysis"
	"github.com/ardasevinc/netbird-cli/internal/ledger"
	"github.com/ardasevinc/netbird-cli/internal/mutation"
	"github.com/ardasevinc/netbird-cli/internal/operations"
)

// Remote is the smallest adapter surface needed by the first consequential
// operation. It deliberately exposes no arbitrary method-and-URL dispatcher.
type Remote interface {
	ServerIdentity() string
	AccountScope(context.Context, string) error
	GetGroup(context.Context, string) (json.RawMessage, error)
	ListGroupsRaw(context.Context) (json.RawMessage, error)
	CreateGroup(context.Context, json.RawMessage) (json.RawMessage, error)
	UpdateGroup(context.Context, string, json.RawMessage) (json.RawMessage, error)
	DeleteGroup(context.Context, string) (json.RawMessage, error)
	GetPolicyRaw(context.Context, string) (json.RawMessage, error)
	ListPoliciesRaw(context.Context) (json.RawMessage, error)
	CreatePolicy(context.Context, json.RawMessage) (json.RawMessage, error)
	UpdatePolicy(context.Context, string, json.RawMessage) (json.RawMessage, error)
	DeletePolicy(context.Context, string) (json.RawMessage, error)
	GetRouteRaw(context.Context, string) (json.RawMessage, error)
	ListRoutesRaw(context.Context) (json.RawMessage, error)
	CreateRoute(context.Context, json.RawMessage) (json.RawMessage, error)
	UpdateRoute(context.Context, string, json.RawMessage) (json.RawMessage, error)
	DeleteRoute(context.Context, string) (json.RawMessage, error)
	GetPeerRaw(context.Context, string) (json.RawMessage, error)
	UpdatePeer(context.Context, string, json.RawMessage) (json.RawMessage, error)
	CreateTemporaryAccessPeer(context.Context, string, json.RawMessage) (json.RawMessage, error)
	CreatePeerJob(context.Context, string, json.RawMessage) (json.RawMessage, error)
	DeletePeer(context.Context, string) (json.RawMessage, error)
	ListEDRBypassedPeersRaw(context.Context) (json.RawMessage, error)
	BypassPeerEDR(context.Context, string) (json.RawMessage, error)
	RevokePeerEDRBypass(context.Context, string) (json.RawMessage, error)
	GetNetworkRaw(context.Context, string) (json.RawMessage, error)
	ListNetworksRaw(context.Context) (json.RawMessage, error)
	CreateNetwork(context.Context, json.RawMessage) (json.RawMessage, error)
	UpdateNetwork(context.Context, string, json.RawMessage) (json.RawMessage, error)
	DeleteNetwork(context.Context, string) (json.RawMessage, error)
	GetNetworkResourceRaw(context.Context, string, string) (json.RawMessage, error)
	ListNetworkResourcesRaw(context.Context, string) (json.RawMessage, error)
	CreateNetworkResource(context.Context, string, json.RawMessage) (json.RawMessage, error)
	UpdateNetworkResource(context.Context, string, string, json.RawMessage) (json.RawMessage, error)
	DeleteNetworkResource(context.Context, string, string) (json.RawMessage, error)
	GetNetworkRouterRaw(context.Context, string, string) (json.RawMessage, error)
	ListNetworkRoutersRaw(context.Context, string) (json.RawMessage, error)
	CreateNetworkRouter(context.Context, string, json.RawMessage) (json.RawMessage, error)
	UpdateNetworkRouter(context.Context, string, string, json.RawMessage) (json.RawMessage, error)
	DeleteNetworkRouter(context.Context, string, string) (json.RawMessage, error)
	ListDNSZonesRaw(context.Context) (json.RawMessage, error)
	CreateDNSZone(context.Context, json.RawMessage) (json.RawMessage, error)
	GetDNSZoneRaw(context.Context, string) (json.RawMessage, error)
	UpdateDNSZone(context.Context, string, json.RawMessage) (json.RawMessage, error)
	DeleteDNSZone(context.Context, string) (json.RawMessage, error)
	ListDNSRecordsRaw(context.Context, string) (json.RawMessage, error)
	CreateDNSRecord(context.Context, string, json.RawMessage) (json.RawMessage, error)
	GetDNSRecordRaw(context.Context, string, string) (json.RawMessage, error)
	UpdateDNSRecord(context.Context, string, string, json.RawMessage) (json.RawMessage, error)
	DeleteDNSRecord(context.Context, string, string) (json.RawMessage, error)
	ListNameserverGroupsRaw(context.Context) (json.RawMessage, error)
	CreateNameserverGroup(context.Context, json.RawMessage) (json.RawMessage, error)
	GetNameserverGroupRaw(context.Context, string) (json.RawMessage, error)
	UpdateNameserverGroup(context.Context, string, json.RawMessage) (json.RawMessage, error)
	DeleteNameserverGroup(context.Context, string) (json.RawMessage, error)
	GetDNSSettingsRaw(context.Context) (json.RawMessage, error)
	UpdateDNSSettings(context.Context, json.RawMessage) (json.RawMessage, error)
	GetAccountRaw(context.Context, string) (json.RawMessage, error)
	UpdateAccount(context.Context, string, json.RawMessage) (json.RawMessage, error)
	DeleteAccount(context.Context, string) (json.RawMessage, error)
	ListPostureChecksRaw(context.Context) (json.RawMessage, error)
	CreatePostureCheck(context.Context, json.RawMessage) (json.RawMessage, error)
	GetPostureCheckRaw(context.Context, string) (json.RawMessage, error)
	UpdatePostureCheck(context.Context, string, json.RawMessage) (json.RawMessage, error)
	DeletePostureCheck(context.Context, string) (json.RawMessage, error)
	ListIngressPeersRaw(context.Context) (json.RawMessage, error)
	CreateIngressPeer(context.Context, json.RawMessage) (json.RawMessage, error)
	GetIngressPeerRaw(context.Context, string) (json.RawMessage, error)
	UpdateIngressPeer(context.Context, string, json.RawMessage) (json.RawMessage, error)
	DeleteIngressPeer(context.Context, string) (json.RawMessage, error)
	ListIngressPortAllocationsRaw(context.Context, string) (json.RawMessage, error)
	GetIngressPortAllocationRaw(context.Context, string, string) (json.RawMessage, error)
	CreateIngressPortAllocation(context.Context, string, json.RawMessage) (json.RawMessage, error)
	UpdateIngressPortAllocation(context.Context, string, string, json.RawMessage) (json.RawMessage, error)
	DeleteIngressPortAllocation(context.Context, string, string) (json.RawMessage, error)
	GetAgentNetworkSettingsRaw(context.Context) (json.RawMessage, error)
	UpdateAgentNetworkSettings(context.Context, json.RawMessage) (json.RawMessage, error)
	CreateAgentNetworkSettings(context.Context, json.RawMessage) (json.RawMessage, error)
	DeleteAgentNetworkSettings(context.Context) (json.RawMessage, error)
	ListAgentNetworkBudgetRulesRaw(context.Context) (json.RawMessage, error)
	GetAgentNetworkBudgetRuleRaw(context.Context, string) (json.RawMessage, error)
	CreateAgentNetworkBudgetRule(context.Context, json.RawMessage) (json.RawMessage, error)
	UpdateAgentNetworkBudgetRule(context.Context, string, json.RawMessage) (json.RawMessage, error)
	DeleteAgentNetworkBudgetRule(context.Context, string) (json.RawMessage, error)
	ListAgentNetworkGuardrailsRaw(context.Context) (json.RawMessage, error)
	GetAgentNetworkGuardrailRaw(context.Context, string) (json.RawMessage, error)
	CreateAgentNetworkGuardrail(context.Context, json.RawMessage) (json.RawMessage, error)
	UpdateAgentNetworkGuardrail(context.Context, string, json.RawMessage) (json.RawMessage, error)
	DeleteAgentNetworkGuardrail(context.Context, string) (json.RawMessage, error)
	ListAgentNetworkPoliciesRaw(context.Context) (json.RawMessage, error)
	GetAgentNetworkPolicyRaw(context.Context, string) (json.RawMessage, error)
	CreateAgentNetworkPolicy(context.Context, json.RawMessage) (json.RawMessage, error)
	UpdateAgentNetworkPolicy(context.Context, string, json.RawMessage) (json.RawMessage, error)
	DeleteAgentNetworkPolicy(context.Context, string) (json.RawMessage, error)
	ListAgentNetworkProvidersRaw(context.Context) (json.RawMessage, error)
	GetAgentNetworkProviderRaw(context.Context, string) (json.RawMessage, error)
	CreateAgentNetworkProvider(context.Context, json.RawMessage) (json.RawMessage, error)
	UpdateAgentNetworkProvider(context.Context, string, json.RawMessage) (json.RawMessage, error)
	DeleteAgentNetworkProvider(context.Context, string) (json.RawMessage, error)
	ListUsersRaw(context.Context) (json.RawMessage, error)
	GetUserRaw(context.Context, string) (json.RawMessage, error)
	CreateUser(context.Context, json.RawMessage) (json.RawMessage, error)
	UpdateUser(context.Context, string, json.RawMessage) (json.RawMessage, error)
	DeleteUser(context.Context, string) (json.RawMessage, error)
	ApproveUser(context.Context, string) (json.RawMessage, error)
	RejectUser(context.Context, string) (json.RawMessage, error)
	ChangeUserPassword(context.Context, string, json.RawMessage) (json.RawMessage, error)
	ResendUserInvite(context.Context, string) (json.RawMessage, error)
	GetPersonalAccessTokenRaw(context.Context, string, string) (json.RawMessage, error)
	ListPersonalAccessTokensRaw(context.Context, string) (json.RawMessage, error)
	CreatePersonalAccessToken(context.Context, string, json.RawMessage) (json.RawMessage, error)
	DeletePersonalAccessToken(context.Context, string, string) (json.RawMessage, error)
	GetSetupKeyRaw(context.Context, string) (json.RawMessage, error)
	ListSetupKeysRaw(context.Context) (json.RawMessage, error)
	CreateSetupKey(context.Context, json.RawMessage) (json.RawMessage, error)
	UpdateSetupKey(context.Context, string, json.RawMessage) (json.RawMessage, error)
	DeleteSetupKey(context.Context, string) (json.RawMessage, error)
	ListInvitesRaw(context.Context) (json.RawMessage, error)
	GetInviteRaw(context.Context, string) (json.RawMessage, error)
	CreateInvite(context.Context, json.RawMessage) (json.RawMessage, error)
	DeleteInvite(context.Context, string) (json.RawMessage, error)
	RegenerateInvite(context.Context, string, json.RawMessage) (json.RawMessage, error)
	GetPublicInviteRaw(context.Context, string) (json.RawMessage, error)
	AcceptInvite(context.Context, string, json.RawMessage) (json.RawMessage, error)
	ListEventStreamingIntegrationsRaw(context.Context) (json.RawMessage, error)
	GetEventStreamingIntegrationRaw(context.Context, string) (json.RawMessage, error)
	CreateEventStreamingIntegration(context.Context, json.RawMessage) (json.RawMessage, error)
	UpdateEventStreamingIntegration(context.Context, string, json.RawMessage) (json.RawMessage, error)
	DeleteEventStreamingIntegration(context.Context, string) (json.RawMessage, error)
	ListIdentityProvidersRaw(context.Context) (json.RawMessage, error)
	GetIdentityProviderRaw(context.Context, string) (json.RawMessage, error)
	CreateIdentityProvider(context.Context, json.RawMessage) (json.RawMessage, error)
	UpdateIdentityProvider(context.Context, string, json.RawMessage) (json.RawMessage, error)
	DeleteIdentityProvider(context.Context, string) (json.RawMessage, error)
	ListReverseProxyTokensRaw(context.Context) (json.RawMessage, error)
	CreateReverseProxyToken(context.Context, json.RawMessage) (json.RawMessage, error)
	DeleteReverseProxyToken(context.Context, string) (json.RawMessage, error)
	ListReverseProxyDomainsRaw(context.Context) (json.RawMessage, error)
	CreateReverseProxyDomain(context.Context, json.RawMessage) (json.RawMessage, error)
	DeleteReverseProxyDomain(context.Context, string) (json.RawMessage, error)
	ListReverseProxyClustersRaw(context.Context) (json.RawMessage, error)
	DeleteReverseProxyCluster(context.Context, string) (json.RawMessage, error)
	ListReverseProxyServicesRaw(context.Context) (json.RawMessage, error)
	GetReverseProxyServiceRaw(context.Context, string) (json.RawMessage, error)
	CreateReverseProxyService(context.Context, json.RawMessage) (json.RawMessage, error)
	UpdateReverseProxyService(context.Context, string, json.RawMessage) (json.RawMessage, error)
	DeleteReverseProxyService(context.Context, string) (json.RawMessage, error)
	ListNotificationChannelsRaw(context.Context) (json.RawMessage, error)
	GetNotificationChannelRaw(context.Context, string) (json.RawMessage, error)
	CreateNotificationChannel(context.Context, json.RawMessage) (json.RawMessage, error)
	UpdateNotificationChannel(context.Context, string, json.RawMessage) (json.RawMessage, error)
	DeleteNotificationChannel(context.Context, string) (json.RawMessage, error)
	ListAzureIDPsRaw(context.Context) (json.RawMessage, error)
	GetAzureIDPRaw(context.Context, string) (json.RawMessage, error)
	CreateAzureIDP(context.Context, json.RawMessage) (json.RawMessage, error)
	UpdateAzureIDP(context.Context, string, json.RawMessage) (json.RawMessage, error)
	DeleteAzureIDP(context.Context, string) (json.RawMessage, error)
	SyncAzureIDP(context.Context, string) (json.RawMessage, error)
	ListGoogleIDPsRaw(context.Context) (json.RawMessage, error)
	GetGoogleIDPRaw(context.Context, string) (json.RawMessage, error)
	CreateGoogleIDP(context.Context, json.RawMessage) (json.RawMessage, error)
	UpdateGoogleIDP(context.Context, string, json.RawMessage) (json.RawMessage, error)
	DeleteGoogleIDP(context.Context, string) (json.RawMessage, error)
	SyncGoogleIDP(context.Context, string) (json.RawMessage, error)
	GetEDRIntegrationRaw(context.Context, string) (json.RawMessage, error)
	CreateEDRIntegration(context.Context, string, json.RawMessage) (json.RawMessage, error)
	UpdateEDRIntegration(context.Context, string, json.RawMessage) (json.RawMessage, error)
	DeleteEDRIntegration(context.Context, string) (json.RawMessage, error)
	ListSCIMIntegrationsRaw(context.Context, string) (json.RawMessage, error)
	GetSCIMIntegrationRaw(context.Context, string, string) (json.RawMessage, error)
	CreateSCIMIntegration(context.Context, string, json.RawMessage) (json.RawMessage, error)
	UpdateSCIMIntegration(context.Context, string, string, json.RawMessage) (json.RawMessage, error)
	DeleteSCIMIntegration(context.Context, string, string) (json.RawMessage, error)
	RegenerateSCIMToken(context.Context, string, string) (json.RawMessage, error)
}

type Ledger interface {
	Get(context.Context, string, int) (ledger.Stage, error)
	BeginAttempt(context.Context, string, int) (ledger.Attempt, error)
	SetAttemptState(context.Context, string, string) error
	RecordReceipt(context.Context, ledger.Receipt) error
}

type ApplyInput struct {
	StageID          string
	Revision         int
	Profile          string
	ServerIdentity   string
	AccountID        string
	Acknowledgements []string
	AckAllBlocking   bool
	SecretResolver   func(string) (string, error)
}

type requestTarget struct {
	ID             string `json:"id"`
	NetworkID      string `json:"network_id"`
	ZoneID         string `json:"zone_id"`
	PeerID         string `json:"peer_id"`
	UserID         string `json:"user_id"`
	TokenID        string `json:"token_id"`
	InviteID       string `json:"invite_id"`
	InviteTokenRef string `json:"invite_token_ref"`
	InviteToken    string `json:"-"`
}

type Result struct {
	StageID       string                 `json:"stage_id"`
	Revision      int                    `json:"revision"`
	AttemptID     string                 `json:"attempt_id,omitempty"`
	State         mutation.DispatchState `json:"state"`
	Reason        string                 `json:"reason,omitempty"`
	OneTimeSecret string                 `json:"-"`
}

type ApplyError struct {
	Result Result
	Err    error
}

func (e *ApplyError) Error() string { return e.Err.Error() }

func (e *ApplyError) Unwrap() error { return e.Err }

func Apply(ctx context.Context, store Ledger, remote Remote, input ApplyInput) (Result, error) {
	if input.StageID == "" || input.Revision < 1 {
		return Result{}, &ApplyError{Err: errors.New("apply requires an exact stage id and positive revision")}
	}
	stage, err := store.Get(ctx, input.StageID, input.Revision)
	if err != nil {
		return Result{}, &ApplyError{Err: err}
	}
	result := Result{StageID: stage.ID, Revision: stage.Revision}
	if stage.Cancelled {
		return result, &ApplyError{Result: result, Err: errors.New("stage revision is cancelled")}
	}
	definition, err := operations.Lookup(stage.Operation)
	if err != nil {
		return result, &ApplyError{Result: result, Err: err}
	}
	if definition.Mutation != operations.Consequential || !definition.DispatcherAdmitted {
		return result, &ApplyError{Result: result, Err: fmt.Errorf("operation %q is not admitted for dispatch", stage.Operation)}
	}
	if input.Profile == "" || stage.Profile != input.Profile {
		return result, &ApplyError{Result: result, Err: errors.New("stage profile does not match the selected profile")}
	}
	if input.ServerIdentity == "" || stage.ServerIdentity == "" || stage.ServerIdentity != input.ServerIdentity || remote.ServerIdentity() != input.ServerIdentity {
		return result, &ApplyError{Result: result, Err: errors.New("server identity does not match the staged identity")}
	}
	if input.AccountID == "" || stage.AccountID == "" || stage.AccountID != input.AccountID {
		return result, &ApplyError{Result: result, Err: errors.New("account scope does not match the staged account")}
	}
	if err := remote.AccountScope(ctx, input.AccountID); err != nil {
		return result, &ApplyError{Result: result, Err: err}
	}
	var request requestTarget
	if err := json.Unmarshal(stage.Request, &request); err != nil || (request.ID == "" && !isCreateOperation(stage.Operation) && !isTargetlessOperation(stage.Operation) && !isUserTokenDeleteOperation(stage.Operation) && stage.Operation != "users.tokens.create" && stage.Operation != "users.invites.delete" && stage.Operation != "users.invites.regenerate" && stage.Operation != "users.invites.accept") {
		return result, &ApplyError{Result: result, Err: fmt.Errorf("%s stage request requires a target id", stage.Operation)}
	}
	if (isUserTokenDeleteOperation(stage.Operation) || stage.Operation == "users.tokens.create") && request.UserID == "" {
		return result, &ApplyError{Result: result, Err: fmt.Errorf("%s stage request requires user_id", stage.Operation)}
	}
	if strings.HasPrefix(stage.Operation, "peers.ingress.ports.") && request.PeerID == "" {
		return result, &ApplyError{Result: result, Err: fmt.Errorf("%s stage request requires peer_id", stage.Operation)}
	}
	if stage.Operation == "peers.jobs.create" && request.PeerID == "" {
		return result, &ApplyError{Result: result, Err: fmt.Errorf("%s stage request requires peer_id", stage.Operation)}
	}
	if isUserTokenDeleteOperation(stage.Operation) && request.TokenID == "" {
		return result, &ApplyError{Result: result, Err: fmt.Errorf("%s stage request requires token_id", stage.Operation)}
	}
	if (stage.Operation == "users.invites.delete" || stage.Operation == "users.invites.regenerate") && request.InviteID == "" {
		return result, &ApplyError{Result: result, Err: fmt.Errorf("%s stage request requires invite_id", stage.Operation)}
	}
	if stage.Operation == "users.invites.accept" {
		if strings.TrimSpace(request.InviteTokenRef) == "" {
			return result, &ApplyError{Result: result, Err: errors.New("users.invites.accept stage request requires invite_token_ref")}
		}
		if input.SecretResolver == nil {
			return result, &ApplyError{Result: result, Err: errors.New("users.invites.accept requires a configured secret resolver")}
		}
		token, err := input.SecretResolver(request.InviteTokenRef)
		if err != nil || strings.TrimSpace(token) == "" {
			return result, &ApplyError{Result: result, Err: errors.New("users.invites.accept invite_token_ref could not be resolved")}
		}
		request.InviteToken = token
	}
	if (stage.Operation == "networks.resources.create" || stage.Operation == "networks.resources.update" || stage.Operation == "networks.resources.delete" || stage.Operation == "networks.routers.create" || stage.Operation == "networks.routers.update" || stage.Operation == "networks.routers.delete") && request.NetworkID == "" {
		return result, &ApplyError{Result: result, Err: fmt.Errorf("%s stage request requires network_id", stage.Operation)}
	}
	if strings.HasPrefix(stage.Operation, "dns.records.") && request.ZoneID == "" {
		return result, &ApplyError{Result: result, Err: fmt.Errorf("%s stage request requires zone_id", stage.Operation)}
	}
	findings := make([]mutation.Finding, 0, len(stage.Findings))
	for _, finding := range stage.Findings {
		findings = append(findings, mutation.Finding{Code: finding.Code, Severity: mutation.Severity(finding.Severity), Message: finding.Message})
	}
	if err := mutation.ValidateAcknowledgements(findings, input.Acknowledgements, input.AckAllBlocking); err != nil {
		return result, &ApplyError{Result: result, Err: err}
	}
	liveBefore, err := readPreimage(ctx, remote, stage.Operation, request)
	if err != nil {
		return result, &ApplyError{Result: result, Err: fmt.Errorf("re-read %s preimage: %w", stage.Operation, err)}
	}
	preimage := mutation.PreimageMatches
	if stage.Operation == "azure_idp.sync" || stage.Operation == "google_idp.sync" {
		equal, err := mutation.Equivalent(stage.Before, liveBefore)
		if err != nil {
			return result, &ApplyError{Result: result, Err: fmt.Errorf("classify %s preimage: %w", stage.Operation, err)}
		}
		if !equal {
			return result, &ApplyError{Result: result, Err: fmt.Errorf("staged %s preimage drifted; create a new revision", stage.Operation)}
		}
		preimage = mutation.PreimageMatches
	} else if isCreateOperation(stage.Operation) {
		equal, err := mutation.Equivalent(stage.Before, liveBefore)
		if err != nil {
			return result, &ApplyError{Result: result, Err: fmt.Errorf("classify %s collection preimage: %w", stage.Operation, err)}
		}
		if !equal {
			return result, &ApplyError{Result: result, Err: fmt.Errorf("staged %s preimage drifted; create a new revision", stage.Operation)}
		}
		if stage.Operation == "peers.temporary_access.create" || stage.Operation == "peers.jobs.create" {
			preimage = mutation.PreimageMatches
		} else {
			already, err := collectionContainsIntent(liveBefore, stage.IntendedAfter)
			if err != nil {
				return result, &ApplyError{Result: result, Err: fmt.Errorf("check %s existing state: %w", stage.Operation, err)}
			}
			if already {
				preimage = mutation.PreimageAlreadySatisfied
			}
		}
	} else {
		preimage, err = mutation.ClassifyPreimage(stage.Before, liveBefore, stage.IntendedAfter)
		if err != nil {
			return result, &ApplyError{Result: result, Err: fmt.Errorf("classify %s preimage: %w", stage.Operation, err)}
		}
		if preimage == mutation.PreimageDrifted {
			return result, &ApplyError{Result: result, Err: fmt.Errorf("staged %s preimage drifted; create a new revision", stage.Operation)}
		}
	}
	impact, err := mutationImpact(stage.Operation, liveBefore, stage.IntendedAfter)
	if err != nil {
		return result, &ApplyError{Result: result, Err: fmt.Errorf("recompute %s mutation impact: %w", stage.Operation, err)}
	}
	if len(stage.Impact) != 0 && string(stage.Impact) != "{}" {
		liveImpact, err := json.Marshal(impact)
		if err != nil {
			return result, &ApplyError{Result: result, Err: fmt.Errorf("encode %s mutation impact: %w", stage.Operation, err)}
		}
		equal, err := mutation.Equivalent(stage.Impact, liveImpact)
		if err != nil || !equal {
			return result, &ApplyError{Result: result, Err: errors.New("staged mutation impact changed; create a new revision")}
		}
	}
	attempt, err := store.BeginAttempt(ctx, stage.ID, stage.Revision)
	if err != nil {
		return result, &ApplyError{Result: result, Err: err}
	}
	result.AttemptID = attempt.ID
	if preimage == mutation.PreimageAlreadySatisfied && stage.Operation != "users.password.update" && stage.Operation != "users.invite.resend" && stage.Operation != "scim.token" && stage.Operation != "okta_scim.token" {
		return finish(ctx, store, result, mutation.AlreadySatisfied, "remote state already equals intended state")
	}
	dispatchRequest, err := prepareSecretRequest(stage.Operation, stage.Request, input.SecretResolver)
	if err != nil {
		return finish(ctx, store, result, mutation.NotDispatched, "request secret could not be resolved; mutation was not dispatched")
	}
	dispatchResult, err := dispatch(ctx, remote, stage.Operation, request, dispatchRequest)
	if err != nil {
		state := classifyDispatchError(err)
		return finish(ctx, store, result, state, stage.Operation+" did not produce a confirmed success")
	}
	if stage.Operation == "azure_idp.sync" || stage.Operation == "google_idp.sync" {
		var response struct {
			Result string `json:"result"`
		}
		if err := json.Unmarshal(dispatchResult, &response); err != nil || response.Result != "ok" {
			return finish(ctx, store, result, mutation.Unknown, "IDP sync response did not confirm success")
		}
		return finish(ctx, store, result, mutation.ConfirmedSuccess, "IDP sync endpoint confirmed synchronization")
	}
	if stage.Operation == "scim.token" || stage.Operation == "okta_scim.token" {
		secret, err := responseSecret(dispatchResult)
		if err != nil {
			return finish(ctx, store, result, mutation.EffectConfirmedReceiptFail, "SCIM token regenerated but its one-time value could not be delivered")
		}
		result.OneTimeSecret = secret
		return finish(ctx, store, result, mutation.ConfirmedSuccess, "SCIM token regeneration endpoint confirmed success")
	}
	if isCreateOperation(stage.Operation) {
		if stage.Operation == "peers.temporary_access.create" {
			if _, err := responseID(dispatchResult); err != nil {
				return finish(ctx, store, result, mutation.Unknown, "temporary access may have applied, but the created peer id could not be confirmed")
			}
			matches, err := objectContains(dispatchResult, stage.IntendedAfter)
			if err != nil {
				return finish(ctx, store, result, mutation.Unknown, "temporary access response could not be compared with intent")
			}
			if !matches {
				return finish(ctx, store, result, mutation.Partial, "temporary access response differs from intended state")
			}
			return finish(ctx, store, result, mutation.ConfirmedSuccess, "temporary access peer response matches intended state")
		}
		if stage.Operation == "peers.jobs.create" {
			if _, err := responseID(dispatchResult); err != nil {
				return finish(ctx, store, result, mutation.Unknown, "remote job may have applied, but the job id could not be confirmed")
			}
			matches, err := objectContains(dispatchResult, stage.IntendedAfter)
			if err != nil {
				return finish(ctx, store, result, mutation.Unknown, "remote job response could not be compared with intent")
			}
			if !matches {
				return finish(ctx, store, result, mutation.Partial, "remote job response differs from intended state")
			}
			return finish(ctx, store, result, mutation.ConfirmedSuccess, "remote job response matches intended state")
		}
		if stage.Operation == "azure_idp.create" {
			matches, err := objectContains(dispatchResult, stage.IntendedAfter)
			if err != nil {
				return finish(ctx, store, result, mutation.Unknown, "Azure IDP create response could not be compared with intent")
			}
			if !matches {
				return finish(ctx, store, result, mutation.Partial, "Azure IDP create response differs from intended state")
			}
			return finish(ctx, store, result, mutation.ConfirmedSuccess, "Azure IDP create response matches intended state")
		}
		if stage.Operation == "google_idp.create" {
			matches, err := objectContains(dispatchResult, stage.IntendedAfter)
			if err != nil {
				return finish(ctx, store, result, mutation.Unknown, "google IDP create response could not be compared with intent")
			}
			if !matches {
				return finish(ctx, store, result, mutation.Partial, "google IDP create response differs from intended state")
			}
			return finish(ctx, store, result, mutation.ConfirmedSuccess, "google IDP create response matches intended state")
		}
		if (strings.HasPrefix(stage.Operation, "edr.") || strings.HasPrefix(stage.Operation, "scim.") || strings.HasPrefix(stage.Operation, "okta_scim.")) && strings.HasSuffix(stage.Operation, ".create") {
			matches, err := objectContains(dispatchResult, stage.IntendedAfter)
			if err != nil {
				return finish(ctx, store, result, mutation.Unknown, "integration create response could not be compared with intent")
			}
			if !matches {
				return finish(ctx, store, result, mutation.Partial, "integration create response differs from intended state")
			}
			liveAfter, err := readPreimage(ctx, remote, stage.Operation, request)
			if err != nil {
				return finish(ctx, store, result, mutation.Unknown, "integration create may have applied, but read-back was inconclusive")
			}
			if strings.HasPrefix(stage.Operation, "scim.") || strings.HasPrefix(stage.Operation, "okta_scim.") {
				matches, err = collectionContainsIntent(liveAfter, stage.IntendedAfter)
			} else {
				matches, err = objectContains(liveAfter, stage.IntendedAfter)
			}
			if err != nil || !matches {
				return finish(ctx, store, result, mutation.Partial, "integration create read-back differs from intended state")
			}
			return finish(ctx, store, result, mutation.ConfirmedSuccess, "integration create response and read-back match intended state")
		}
		if stage.Operation == "peers.edr.bypass.create" {
			liveAfter, err := readPreimage(ctx, remote, stage.Operation, request)
			if err != nil {
				return finish(ctx, store, result, mutation.Unknown, "EDR bypass may have applied, but bypass state could not be confirmed")
			}
			present, err := collectionContainsPeerID(liveAfter, request.ID)
			if err != nil {
				return finish(ctx, store, result, mutation.Unknown, "EDR bypass response could not be compared with the peer state")
			}
			if !present {
				return finish(ctx, store, result, mutation.Partial, "EDR bypass dispatch completed but the peer is not present in bypassed state")
			}
			return finish(ctx, store, result, mutation.ConfirmedSuccess, "peer is present in the EDR-bypassed state")
		}
		createdID, err := responseID(dispatchResult)
		if err != nil {
			return finish(ctx, store, result, mutation.Unknown, "create may have applied, but the created resource id could not be confirmed")
		}
		liveAfter, err := readPreimage(ctx, remote, stage.Operation, requestTarget{NetworkID: request.NetworkID, ZoneID: request.ZoneID, UserID: request.UserID, ID: createdID})
		if err != nil {
			return finish(ctx, store, result, mutation.Unknown, "create may have applied, but read-back was inconclusive")
		}
		actual, err := collectionFindID(liveAfter, createdID)
		if err != nil {
			return finish(ctx, store, result, mutation.Unknown, "created resource could not be found in read-back")
		}
		matches, err := objectContains(actual, stage.IntendedAfter)
		if err != nil {
			return finish(ctx, store, result, mutation.Unknown, "created resource response could not be compared with intent")
		}
		if !matches {
			return finish(ctx, store, result, mutation.Partial, "created resource differs from intended state after create")
		}
		if stage.Operation == "users.tokens.create" || stage.Operation == "setup_keys.create" || stage.Operation == "users.invites.create" || stage.Operation == "reverse_proxy_tokens.create" {
			secret, err := responseSecret(dispatchResult)
			if err != nil {
				return finish(ctx, store, result, mutation.EffectConfirmedReceiptFail, "created resource applied but its one-time value could not be delivered")
			}
			result.OneTimeSecret = secret
		}
		return finish(ctx, store, result, mutation.ConfirmedSuccess, "created resource matches intended state")
	}
	if isDeleteOperation(stage.Operation) {
		if stage.Operation == "peers.edr.bypass.delete" {
			liveAfter, err := readPreimage(ctx, remote, stage.Operation, request)
			if err != nil {
				return finish(ctx, store, result, mutation.Unknown, "EDR bypass may have applied, but bypass state could not be confirmed")
			}
			present, err := collectionContainsPeerID(liveAfter, request.ID)
			if err != nil {
				return finish(ctx, store, result, mutation.Unknown, "EDR bypass absence could not be compared with the peer state")
			}
			if present {
				return finish(ctx, store, result, mutation.Partial, "EDR bypass remains present after revoke")
			}
			return finish(ctx, store, result, mutation.ConfirmedSuccess, "peer is absent from the EDR-bypassed state after revoke")
		}
		if stage.Operation == "reverse_proxy_domains.delete" {
			liveAfter, err := readPreimage(ctx, remote, stage.Operation, request)
			if err != nil {
				return finish(ctx, store, result, mutation.Unknown, "reverse proxy domain may have been deleted, but domain absence could not be confirmed")
			}
			present, err := collectionContainsID(liveAfter, request.ID)
			if err != nil {
				return finish(ctx, store, result, mutation.Unknown, "reverse proxy domain absence could not be compared with the domain collection")
			}
			if present {
				return finish(ctx, store, result, mutation.Partial, "reverse proxy domain remains present after delete")
			}
			return finish(ctx, store, result, mutation.ConfirmedSuccess, "reverse proxy domain is absent after delete")
		}
		if stage.Operation == "reverse_proxy_clusters.delete" {
			liveAfter, err := readPreimage(ctx, remote, stage.Operation, request)
			if err != nil {
				return finish(ctx, store, result, mutation.Unknown, "reverse proxy cluster may have been deleted, but cluster absence could not be confirmed")
			}
			present, err := collectionContainsAddress(liveAfter, request.ID)
			if err != nil {
				return finish(ctx, store, result, mutation.Unknown, "reverse proxy cluster absence could not be compared with the cluster collection")
			}
			if present {
				return finish(ctx, store, result, mutation.Partial, "reverse proxy cluster remains present after delete")
			}
			return finish(ctx, store, result, mutation.ConfirmedSuccess, "reverse proxy cluster is absent after delete")
		}
		if err := confirmDeleted(ctx, remote, stage.Operation, request); err != nil && !isNotFound(err) {
			return finish(ctx, store, result, mutation.Unknown, "delete may have applied, but absence could not be confirmed")
		}
		return finish(ctx, store, result, mutation.ConfirmedSuccess, "remote "+strings.TrimSuffix(stage.Operation, ".delete")+" is absent after delete")
	}
	if stage.Operation == "users.invites.accept" {
		var response struct {
			Success bool `json:"success"`
		}
		if err := json.Unmarshal(dispatchResult, &response); err != nil || !response.Success {
			return finish(ctx, store, result, mutation.Unknown, "invite acceptance response did not confirm success")
		}
		return finish(ctx, store, result, mutation.ConfirmedSuccess, "invite acceptance endpoint confirmed account creation")
	}
	liveAfter, err := readPreimage(ctx, remote, stage.Operation, request)
	if err != nil {
		return finish(ctx, store, result, mutation.Unknown, "update may have applied, but read-back was inconclusive")
	}
	equal, err := mutation.Equivalent(liveAfter, stage.IntendedAfter)
	if err != nil {
		return finish(ctx, store, result, mutation.Unknown, "update response could not be compared with intended state")
	}
	if !equal {
		return finish(ctx, store, result, mutation.Partial, "remote state differs from intended state after update")
	}
	if stage.Operation == "users.invites.regenerate" {
		secret, err := responseSecret(dispatchResult)
		if err != nil {
			return finish(ctx, store, result, mutation.EffectConfirmedReceiptFail, "invite was regenerated but its one-time value could not be delivered")
		}
		result.OneTimeSecret = secret
	}
	return finish(ctx, store, result, mutation.ConfirmedSuccess, "remote state matches intended state")
}

func readPreimage(ctx context.Context, remote Remote, operation string, target requestTarget) (json.RawMessage, error) {
	switch operation {
	case "groups.update":
		return remote.GetGroup(ctx, target.ID)
	case "groups.delete":
		return remote.GetGroup(ctx, target.ID)
	case "groups.create":
		return remote.ListGroupsRaw(ctx)
	case "policies.update":
		return remote.GetPolicyRaw(ctx, target.ID)
	case "policies.delete":
		return remote.GetPolicyRaw(ctx, target.ID)
	case "policies.create":
		return remote.ListPoliciesRaw(ctx)
	case "dns.zones.create":
		return remote.ListDNSZonesRaw(ctx)
	case "dns.zones.delete":
		return remote.GetDNSZoneRaw(ctx, target.ID)
	case "dns.zones.update":
		return remote.GetDNSZoneRaw(ctx, target.ID)
	case "dns.records.create":
		return remote.ListDNSRecordsRaw(ctx, target.ZoneID)
	case "dns.records.update":
		return remote.GetDNSRecordRaw(ctx, target.ZoneID, target.ID)
	case "dns.records.delete":
		return remote.GetDNSRecordRaw(ctx, target.ZoneID, target.ID)
	case "dns.nameservers.create":
		return remote.ListNameserverGroupsRaw(ctx)
	case "dns.nameservers.update":
		return remote.GetNameserverGroupRaw(ctx, target.ID)
	case "dns.nameservers.delete":
		return remote.GetNameserverGroupRaw(ctx, target.ID)
	case "dns.settings.update":
		return remote.GetDNSSettingsRaw(ctx)
	case "accounts.update":
		return remote.GetAccountRaw(ctx, target.ID)
	case "accounts.delete":
		return remote.GetAccountRaw(ctx, target.ID)
	case "posture_checks.create":
		return remote.ListPostureChecksRaw(ctx)
	case "posture_checks.update":
		return remote.GetPostureCheckRaw(ctx, target.ID)
	case "posture_checks.delete":
		return remote.GetPostureCheckRaw(ctx, target.ID)
	case "ingress.peers.create":
		return remote.ListIngressPeersRaw(ctx)
	case "ingress.peers.update":
		return remote.GetIngressPeerRaw(ctx, target.ID)
	case "ingress.peers.delete":
		return remote.GetIngressPeerRaw(ctx, target.ID)
	case "peers.ingress.ports.create":
		return remote.ListIngressPortAllocationsRaw(ctx, target.PeerID)
	case "peers.ingress.ports.update", "peers.ingress.ports.delete":
		return remote.GetIngressPortAllocationRaw(ctx, target.PeerID, target.ID)
	case "agent_network.settings.update":
		return remote.GetAgentNetworkSettingsRaw(ctx)
	case "agent_network.settings.create":
		return remote.GetAgentNetworkSettingsRaw(ctx)
	case "agent_network.settings.delete":
		return remote.GetAgentNetworkSettingsRaw(ctx)
	case "agent_network.budget_rules.create":
		return remote.ListAgentNetworkBudgetRulesRaw(ctx)
	case "agent_network.budget_rules.update":
		return remote.GetAgentNetworkBudgetRuleRaw(ctx, target.ID)
	case "agent_network.budget_rules.delete":
		return remote.GetAgentNetworkBudgetRuleRaw(ctx, target.ID)
	case "agent_network.guardrails.create":
		return remote.ListAgentNetworkGuardrailsRaw(ctx)
	case "agent_network.guardrails.update":
		return remote.GetAgentNetworkGuardrailRaw(ctx, target.ID)
	case "agent_network.guardrails.delete":
		return remote.GetAgentNetworkGuardrailRaw(ctx, target.ID)
	case "agent_network.policies.create":
		return remote.ListAgentNetworkPoliciesRaw(ctx)
	case "agent_network.policies.update":
		return remote.GetAgentNetworkPolicyRaw(ctx, target.ID)
	case "agent_network.policies.delete":
		return remote.GetAgentNetworkPolicyRaw(ctx, target.ID)
	case "agent_network.providers.create":
		return remote.ListAgentNetworkProvidersRaw(ctx)
	case "agent_network.providers.update":
		return remote.GetAgentNetworkProviderRaw(ctx, target.ID)
	case "agent_network.providers.delete":
		return remote.GetAgentNetworkProviderRaw(ctx, target.ID)
	case "users.create":
		return remote.ListUsersRaw(ctx)
	case "users.update", "users.delete", "users.approve", "users.reject", "users.password.update", "users.invite.resend":
		return remote.GetUserRaw(ctx, target.ID)
	case "users.tokens.delete":
		return remote.GetPersonalAccessTokenRaw(ctx, target.UserID, target.TokenID)
	case "users.tokens.create":
		return remote.ListPersonalAccessTokensRaw(ctx, target.UserID)
	case "setup_keys.update", "setup_keys.delete":
		return remote.GetSetupKeyRaw(ctx, target.ID)
	case "setup_keys.create":
		return remote.ListSetupKeysRaw(ctx)
	case "event_streaming.create":
		return remote.ListEventStreamingIntegrationsRaw(ctx)
	case "event_streaming.update", "event_streaming.delete":
		return remote.GetEventStreamingIntegrationRaw(ctx, target.ID)
	case "identity_providers.create":
		return remote.ListIdentityProvidersRaw(ctx)
	case "identity_providers.update", "identity_providers.delete":
		return remote.GetIdentityProviderRaw(ctx, target.ID)
	case "reverse_proxy_tokens.create":
		return remote.ListReverseProxyTokensRaw(ctx)
	case "reverse_proxy_tokens.delete":
		return remote.ListReverseProxyTokensRaw(ctx)
	case "reverse_proxy_domains.create", "reverse_proxy_domains.delete":
		return remote.ListReverseProxyDomainsRaw(ctx)
	case "reverse_proxy_clusters.delete":
		return remote.ListReverseProxyClustersRaw(ctx)
	case "reverse_proxy_services.create":
		return remote.ListReverseProxyServicesRaw(ctx)
	case "reverse_proxy_services.update", "reverse_proxy_services.delete":
		return remote.GetReverseProxyServiceRaw(ctx, target.ID)
	case "notification_channels.create":
		return remote.ListNotificationChannelsRaw(ctx)
	case "notification_channels.update", "notification_channels.delete":
		return remote.GetNotificationChannelRaw(ctx, target.ID)
	case "azure_idp.create":
		return remote.ListAzureIDPsRaw(ctx)
	case "azure_idp.update", "azure_idp.delete", "azure_idp.sync":
		return remote.GetAzureIDPRaw(ctx, target.ID)
	case "google_idp.create":
		return remote.ListGoogleIDPsRaw(ctx)
	case "google_idp.update", "google_idp.delete", "google_idp.sync":
		return remote.GetGoogleIDPRaw(ctx, target.ID)
	case "edr.intune.create":
		return readOptionalEDRIntegration(ctx, remote, "intune")
	case "edr.intune.update", "edr.intune.delete":
		return remote.GetEDRIntegrationRaw(ctx, "intune")
	case "edr.sentinelone.create":
		return readOptionalEDRIntegration(ctx, remote, "sentinelone")
	case "edr.sentinelone.update", "edr.sentinelone.delete":
		return remote.GetEDRIntegrationRaw(ctx, "sentinelone")
	case "edr.falcon.create":
		return readOptionalEDRIntegration(ctx, remote, "falcon")
	case "edr.falcon.update", "edr.falcon.delete":
		return remote.GetEDRIntegrationRaw(ctx, "falcon")
	case "edr.huntress.create":
		return readOptionalEDRIntegration(ctx, remote, "huntress")
	case "edr.huntress.update", "edr.huntress.delete":
		return remote.GetEDRIntegrationRaw(ctx, "huntress")
	case "edr.fleetdm.create":
		return readOptionalEDRIntegration(ctx, remote, "fleetdm")
	case "edr.fleetdm.update", "edr.fleetdm.delete":
		return remote.GetEDRIntegrationRaw(ctx, "fleetdm")
	case "scim.create":
		return remote.ListSCIMIntegrationsRaw(ctx, "scim")
	case "scim.update", "scim.delete", "scim.token":
		return remote.GetSCIMIntegrationRaw(ctx, "scim", target.ID)
	case "okta_scim.create":
		return remote.ListSCIMIntegrationsRaw(ctx, "okta_scim")
	case "okta_scim.update", "okta_scim.delete", "okta_scim.token":
		return remote.GetSCIMIntegrationRaw(ctx, "okta_scim", target.ID)
	case "users.invites.create":
		return remote.ListInvitesRaw(ctx)
	case "users.invites.delete", "users.invites.regenerate":
		return remote.GetInviteRaw(ctx, target.InviteID)
	case "users.invites.accept":
		return remote.GetPublicInviteRaw(ctx, target.InviteToken)
	case "routes.update":
		return remote.GetRouteRaw(ctx, target.ID)
	case "routes.delete":
		return remote.GetRouteRaw(ctx, target.ID)
	case "routes.create":
		return remote.ListRoutesRaw(ctx)
	case "peers.update":
		return remote.GetPeerRaw(ctx, target.ID)
	case "peers.temporary_access.create":
		return remote.GetPeerRaw(ctx, target.ID)
	case "peers.jobs.create":
		return remote.GetPeerRaw(ctx, target.PeerID)
	case "peers.delete":
		return remote.GetPeerRaw(ctx, target.ID)
	case "peers.edr.bypass.create", "peers.edr.bypass.delete":
		return remote.ListEDRBypassedPeersRaw(ctx)
	case "networks.update":
		return remote.GetNetworkRaw(ctx, target.ID)
	case "networks.delete":
		return remote.GetNetworkRaw(ctx, target.ID)
	case "networks.create":
		return remote.ListNetworksRaw(ctx)
	case "networks.resources.delete":
		return remote.GetNetworkResourceRaw(ctx, target.NetworkID, target.ID)
	case "networks.resources.update":
		return remote.GetNetworkResourceRaw(ctx, target.NetworkID, target.ID)
	case "networks.resources.create":
		return remote.ListNetworkResourcesRaw(ctx, target.NetworkID)
	case "networks.routers.create":
		return remote.ListNetworkRoutersRaw(ctx, target.NetworkID)
	case "networks.routers.delete":
		return remote.GetNetworkRouterRaw(ctx, target.NetworkID, target.ID)
	case "networks.routers.update":
		return remote.GetNetworkRouterRaw(ctx, target.NetworkID, target.ID)
	default:
		return nil, fmt.Errorf("operation %q has no preimage reader", operation)
	}
}

func readOptionalEDRIntegration(ctx context.Context, remote Remote, provider string) (json.RawMessage, error) {
	result, err := remote.GetEDRIntegrationRaw(ctx, provider)
	if isNotFound(err) {
		return json.RawMessage(`null`), nil
	}
	return result, err
}

func dispatch(ctx context.Context, remote Remote, operation string, target requestTarget, request json.RawMessage) (json.RawMessage, error) {
	switch operation {
	case "groups.update":
		return remote.UpdateGroup(ctx, target.ID, request)
	case "groups.delete":
		return remote.DeleteGroup(ctx, target.ID)
	case "groups.create":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.CreateGroup(ctx, body)
	case "policies.update":
		return remote.UpdatePolicy(ctx, target.ID, request)
	case "policies.delete":
		return remote.DeletePolicy(ctx, target.ID)
	case "policies.create":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.CreatePolicy(ctx, body)
	case "dns.zones.create":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.CreateDNSZone(ctx, body)
	case "dns.zones.delete":
		return remote.DeleteDNSZone(ctx, target.ID)
	case "dns.zones.update":
		return remote.UpdateDNSZone(ctx, target.ID, request)
	case "dns.records.create":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.CreateDNSRecord(ctx, target.ZoneID, body)
	case "dns.records.update":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.UpdateDNSRecord(ctx, target.ZoneID, target.ID, body)
	case "dns.records.delete":
		return remote.DeleteDNSRecord(ctx, target.ZoneID, target.ID)
	case "dns.nameservers.create":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.CreateNameserverGroup(ctx, body)
	case "dns.nameservers.update":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.UpdateNameserverGroup(ctx, target.ID, body)
	case "dns.nameservers.delete":
		return remote.DeleteNameserverGroup(ctx, target.ID)
	case "dns.settings.update":
		return remote.UpdateDNSSettings(ctx, request)
	case "accounts.update":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.UpdateAccount(ctx, target.ID, body)
	case "accounts.delete":
		return remote.DeleteAccount(ctx, target.ID)
	case "posture_checks.create":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.CreatePostureCheck(ctx, body)
	case "posture_checks.update":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.UpdatePostureCheck(ctx, target.ID, body)
	case "posture_checks.delete":
		return remote.DeletePostureCheck(ctx, target.ID)
	case "ingress.peers.create":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.CreateIngressPeer(ctx, body)
	case "ingress.peers.update":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.UpdateIngressPeer(ctx, target.ID, body)
	case "ingress.peers.delete":
		return remote.DeleteIngressPeer(ctx, target.ID)
	case "peers.ingress.ports.create":
		body, err := stripIngressPortTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.CreateIngressPortAllocation(ctx, target.PeerID, body)
	case "peers.ingress.ports.update":
		body, err := stripIngressPortTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.UpdateIngressPortAllocation(ctx, target.PeerID, target.ID, body)
	case "peers.ingress.ports.delete":
		return remote.DeleteIngressPortAllocation(ctx, target.PeerID, target.ID)
	case "agent_network.settings.update":
		return remote.UpdateAgentNetworkSettings(ctx, request)
	case "agent_network.settings.create":
		return remote.CreateAgentNetworkSettings(ctx, request)
	case "agent_network.settings.delete":
		return remote.DeleteAgentNetworkSettings(ctx)
	case "agent_network.budget_rules.create":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.CreateAgentNetworkBudgetRule(ctx, body)
	case "agent_network.budget_rules.update":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.UpdateAgentNetworkBudgetRule(ctx, target.ID, body)
	case "agent_network.budget_rules.delete":
		return remote.DeleteAgentNetworkBudgetRule(ctx, target.ID)
	case "agent_network.guardrails.create":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.CreateAgentNetworkGuardrail(ctx, body)
	case "agent_network.guardrails.update":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.UpdateAgentNetworkGuardrail(ctx, target.ID, body)
	case "agent_network.guardrails.delete":
		return remote.DeleteAgentNetworkGuardrail(ctx, target.ID)
	case "agent_network.policies.create":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.CreateAgentNetworkPolicy(ctx, body)
	case "agent_network.policies.update":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.UpdateAgentNetworkPolicy(ctx, target.ID, body)
	case "agent_network.policies.delete":
		return remote.DeleteAgentNetworkPolicy(ctx, target.ID)
	case "agent_network.providers.create":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.CreateAgentNetworkProvider(ctx, body)
	case "agent_network.providers.update":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.UpdateAgentNetworkProvider(ctx, target.ID, body)
	case "agent_network.providers.delete":
		return remote.DeleteAgentNetworkProvider(ctx, target.ID)
	case "users.create":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.CreateUser(ctx, body)
	case "users.update":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.UpdateUser(ctx, target.ID, body)
	case "users.delete":
		return remote.DeleteUser(ctx, target.ID)
	case "users.approve":
		return remote.ApproveUser(ctx, target.ID)
	case "users.reject":
		return remote.RejectUser(ctx, target.ID)
	case "users.password.update":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.ChangeUserPassword(ctx, target.ID, body)
	case "users.invite.resend":
		return remote.ResendUserInvite(ctx, target.ID)
	case "users.tokens.delete":
		return remote.DeletePersonalAccessToken(ctx, target.UserID, target.TokenID)
	case "users.tokens.create":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.CreatePersonalAccessToken(ctx, target.UserID, body)
	case "setup_keys.delete":
		return remote.DeleteSetupKey(ctx, target.ID)
	case "setup_keys.update":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.UpdateSetupKey(ctx, target.ID, body)
	case "setup_keys.create":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.CreateSetupKey(ctx, body)
	case "event_streaming.create":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.CreateEventStreamingIntegration(ctx, body)
	case "event_streaming.update":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.UpdateEventStreamingIntegration(ctx, target.ID, body)
	case "event_streaming.delete":
		return remote.DeleteEventStreamingIntegration(ctx, target.ID)
	case "identity_providers.create":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.CreateIdentityProvider(ctx, body)
	case "identity_providers.update":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.UpdateIdentityProvider(ctx, target.ID, body)
	case "identity_providers.delete":
		return remote.DeleteIdentityProvider(ctx, target.ID)
	case "reverse_proxy_tokens.create":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.CreateReverseProxyToken(ctx, body)
	case "reverse_proxy_tokens.delete":
		return remote.DeleteReverseProxyToken(ctx, target.ID)
	case "reverse_proxy_domains.create":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.CreateReverseProxyDomain(ctx, body)
	case "reverse_proxy_domains.delete":
		return remote.DeleteReverseProxyDomain(ctx, target.ID)
	case "reverse_proxy_clusters.delete":
		return remote.DeleteReverseProxyCluster(ctx, target.ID)
	case "reverse_proxy_services.create":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.CreateReverseProxyService(ctx, body)
	case "reverse_proxy_services.update":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.UpdateReverseProxyService(ctx, target.ID, body)
	case "reverse_proxy_services.delete":
		return remote.DeleteReverseProxyService(ctx, target.ID)
	case "notification_channels.create":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.CreateNotificationChannel(ctx, body)
	case "notification_channels.update":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.UpdateNotificationChannel(ctx, target.ID, body)
	case "notification_channels.delete":
		return remote.DeleteNotificationChannel(ctx, target.ID)
	case "azure_idp.create":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.CreateAzureIDP(ctx, body)
	case "azure_idp.update":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.UpdateAzureIDP(ctx, target.ID, body)
	case "azure_idp.delete":
		return remote.DeleteAzureIDP(ctx, target.ID)
	case "azure_idp.sync":
		return remote.SyncAzureIDP(ctx, target.ID)
	case "google_idp.create":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.CreateGoogleIDP(ctx, body)
	case "google_idp.update":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.UpdateGoogleIDP(ctx, target.ID, body)
	case "google_idp.delete":
		return remote.DeleteGoogleIDP(ctx, target.ID)
	case "google_idp.sync":
		return remote.SyncGoogleIDP(ctx, target.ID)
	case "edr.intune.create":
		return remote.CreateEDRIntegration(ctx, "intune", request)
	case "edr.intune.update":
		return remote.UpdateEDRIntegration(ctx, "intune", request)
	case "edr.intune.delete":
		return remote.DeleteEDRIntegration(ctx, "intune")
	case "edr.sentinelone.create":
		return remote.CreateEDRIntegration(ctx, "sentinelone", request)
	case "edr.sentinelone.update":
		return remote.UpdateEDRIntegration(ctx, "sentinelone", request)
	case "edr.sentinelone.delete":
		return remote.DeleteEDRIntegration(ctx, "sentinelone")
	case "edr.falcon.create":
		return remote.CreateEDRIntegration(ctx, "falcon", request)
	case "edr.falcon.update":
		return remote.UpdateEDRIntegration(ctx, "falcon", request)
	case "edr.falcon.delete":
		return remote.DeleteEDRIntegration(ctx, "falcon")
	case "edr.huntress.create":
		return remote.CreateEDRIntegration(ctx, "huntress", request)
	case "edr.huntress.update":
		return remote.UpdateEDRIntegration(ctx, "huntress", request)
	case "edr.huntress.delete":
		return remote.DeleteEDRIntegration(ctx, "huntress")
	case "edr.fleetdm.create":
		return remote.CreateEDRIntegration(ctx, "fleetdm", request)
	case "edr.fleetdm.update":
		return remote.UpdateEDRIntegration(ctx, "fleetdm", request)
	case "edr.fleetdm.delete":
		return remote.DeleteEDRIntegration(ctx, "fleetdm")
	case "scim.create":
		return remote.CreateSCIMIntegration(ctx, "scim", request)
	case "scim.update":
		return remote.UpdateSCIMIntegration(ctx, "scim", target.ID, request)
	case "scim.delete":
		return remote.DeleteSCIMIntegration(ctx, "scim", target.ID)
	case "scim.token":
		return remote.RegenerateSCIMToken(ctx, "scim", target.ID)
	case "okta_scim.create":
		return remote.CreateSCIMIntegration(ctx, "okta_scim", request)
	case "okta_scim.update":
		return remote.UpdateSCIMIntegration(ctx, "okta_scim", target.ID, request)
	case "okta_scim.delete":
		return remote.DeleteSCIMIntegration(ctx, "okta_scim", target.ID)
	case "okta_scim.token":
		return remote.RegenerateSCIMToken(ctx, "okta_scim", target.ID)
	case "users.invites.create":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.CreateInvite(ctx, body)
	case "users.invites.delete":
		return remote.DeleteInvite(ctx, target.InviteID)
	case "users.invites.regenerate":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.RegenerateInvite(ctx, target.InviteID, body)
	case "users.invites.accept":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.AcceptInvite(ctx, target.InviteToken, body)
	case "routes.update":
		return remote.UpdateRoute(ctx, target.ID, request)
	case "routes.delete":
		return remote.DeleteRoute(ctx, target.ID)
	case "routes.create":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.CreateRoute(ctx, body)
	case "peers.update":
		return remote.UpdatePeer(ctx, target.ID, request)
	case "peers.temporary_access.create":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.CreateTemporaryAccessPeer(ctx, target.ID, body)
	case "peers.jobs.create":
		body, err := stripPeerJobTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.CreatePeerJob(ctx, target.PeerID, body)
	case "peers.delete":
		return remote.DeletePeer(ctx, target.ID)
	case "peers.edr.bypass.create":
		return remote.BypassPeerEDR(ctx, target.ID)
	case "peers.edr.bypass.delete":
		return remote.RevokePeerEDRBypass(ctx, target.ID)
	case "networks.update":
		return remote.UpdateNetwork(ctx, target.ID, request)
	case "networks.delete":
		return remote.DeleteNetwork(ctx, target.ID)
	case "networks.create":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.CreateNetwork(ctx, body)
	case "networks.resources.delete":
		return remote.DeleteNetworkResource(ctx, target.NetworkID, target.ID)
	case "networks.resources.update":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.UpdateNetworkResource(ctx, target.NetworkID, target.ID, body)
	case "networks.resources.create":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.CreateNetworkResource(ctx, target.NetworkID, body)
	case "networks.routers.create":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.CreateNetworkRouter(ctx, target.NetworkID, body)
	case "networks.routers.delete":
		return remote.DeleteNetworkRouter(ctx, target.NetworkID, target.ID)
	case "networks.routers.update":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.UpdateNetworkRouter(ctx, target.NetworkID, target.ID, body)
	default:
		return nil, fmt.Errorf("operation %q has no dispatcher", operation)
	}
}

func isCreateOperation(operation string) bool {
	return operation == "groups.create" || operation == "networks.create" || operation == "networks.resources.create" || operation == "networks.routers.create" || operation == "routes.create" || operation == "policies.create" || operation == "dns.zones.create" || operation == "dns.records.create" || operation == "dns.nameservers.create" || operation == "posture_checks.create" || operation == "ingress.peers.create" || operation == "peers.ingress.ports.create" || operation == "peers.edr.bypass.create" || operation == "peers.temporary_access.create" || operation == "peers.jobs.create" || operation == "event_streaming.create" || operation == "identity_providers.create" || operation == "reverse_proxy_tokens.create" || operation == "reverse_proxy_domains.create" || operation == "reverse_proxy_services.create" || operation == "notification_channels.create" || operation == "azure_idp.create" || operation == "google_idp.create" || operation == "agent_network.budget_rules.create" || operation == "agent_network.guardrails.create" || operation == "agent_network.policies.create" || operation == "agent_network.providers.create" || operation == "users.create" || operation == "users.tokens.create" || operation == "setup_keys.create" || operation == "users.invites.create" || operation == "scim.create" || operation == "okta_scim.create" || (strings.HasPrefix(operation, "edr.") && strings.HasSuffix(operation, ".create"))
}

func isTargetlessOperation(operation string) bool {
	return operation == "dns.settings.update" || operation == "agent_network.settings.update" || operation == "agent_network.settings.create" || operation == "agent_network.settings.delete" || strings.HasPrefix(operation, "edr.")
}

func isUserTokenDeleteOperation(operation string) bool {
	return operation == "users.tokens.delete"
}

func responseID(response json.RawMessage) (string, error) {
	var object struct {
		ID                  string `json:"id"`
		PersonalAccessToken struct {
			ID string `json:"id"`
		} `json:"personal_access_token"`
	}
	if err := json.Unmarshal(response, &object); err != nil {
		return "", errors.New("create response has no id")
	}
	if object.ID != "" {
		return object.ID, nil
	}
	if object.PersonalAccessToken.ID != "" {
		return object.PersonalAccessToken.ID, nil
	}
	return "", errors.New("create response has no id")
}

func responseSecret(response json.RawMessage) (string, error) {
	var object struct {
		PlainToken  string `json:"plain_token"`
		Key         string `json:"key"`
		InviteToken string `json:"invite_token"`
		AuthToken   string `json:"auth_token"`
	}
	if err := json.Unmarshal(response, &object); err != nil {
		return "", errors.New("create response has no one-time secret")
	}
	for _, secret := range []string{object.PlainToken, object.Key, object.InviteToken, object.AuthToken} {
		if strings.TrimSpace(secret) != "" {
			return secret, nil
		}
	}
	return "", errors.New("create response has no one-time secret")
}

func collectionContainsIntent(collection, intent json.RawMessage) (bool, error) {
	var objects []json.RawMessage
	if err := json.Unmarshal(collection, &objects); err != nil {
		return false, err
	}
	for _, object := range objects {
		matches, err := objectContains(object, intent)
		if err != nil {
			return false, err
		}
		if matches {
			return true, nil
		}
	}
	return false, nil
}

func collectionContainsPeerID(collection json.RawMessage, peerID string) (bool, error) {
	intent, err := json.Marshal(map[string]string{"peer_id": peerID})
	if err != nil {
		return false, err
	}
	return collectionContainsIntent(collection, intent)
}

func collectionFindID(collection json.RawMessage, id string) (json.RawMessage, error) {
	var objects []json.RawMessage
	if err := json.Unmarshal(collection, &objects); err != nil {
		return nil, err
	}
	for _, object := range objects {
		var candidate struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(object, &candidate); err != nil {
			return nil, err
		}
		if candidate.ID == id {
			return object, nil
		}
	}
	return nil, errors.New("created id is absent from collection")
}

func collectionContainsID(collection json.RawMessage, id string) (bool, error) {
	var objects []json.RawMessage
	if err := json.Unmarshal(collection, &objects); err != nil {
		return false, err
	}
	for _, object := range objects {
		var candidate struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(object, &candidate); err != nil {
			return false, err
		}
		if candidate.ID == id {
			return true, nil
		}
	}
	return false, nil
}

func collectionContainsAddress(collection json.RawMessage, address string) (bool, error) {
	var objects []json.RawMessage
	if err := json.Unmarshal(collection, &objects); err != nil {
		return false, err
	}
	for _, object := range objects {
		var candidate struct {
			Address        string `json:"address"`
			ClusterAddress string `json:"cluster_address"`
		}
		if err := json.Unmarshal(object, &candidate); err != nil {
			return false, err
		}
		if candidate.Address == address || candidate.ClusterAddress == address {
			return true, nil
		}
	}
	return false, nil
}

func objectContains(actual, expected json.RawMessage) (bool, error) {
	var actualObject, expectedObject map[string]any
	if err := json.Unmarshal(actual, &actualObject); err != nil {
		return false, err
	}
	if err := json.Unmarshal(expected, &expectedObject); err != nil {
		return false, err
	}
	for key, expectedValue := range expectedObject {
		actualValue, ok := actualObject[key]
		if !ok || !jsonValuesEqual(actualValue, expectedValue) {
			return false, nil
		}
	}
	return true, nil
}

func jsonValuesEqual(left, right any) bool {
	leftJSON, _ := json.Marshal(left)
	rightJSON, _ := json.Marshal(right)
	return string(leftJSON) == string(rightJSON)
}

func stripTargetFields(request json.RawMessage) (json.RawMessage, error) {
	var object map[string]any
	if err := json.Unmarshal(request, &object); err != nil {
		return nil, err
	}
	delete(object, "id")
	delete(object, "network_id")
	delete(object, "zone_id")
	delete(object, "user_id")
	delete(object, "token_id")
	delete(object, "invite_id")
	delete(object, "invite_token_ref")
	delete(object, "invite_token")
	delete(object, "password_ref")
	return json.Marshal(object)
}

func stripIngressPortTargetFields(request json.RawMessage) (json.RawMessage, error) {
	body, err := stripTargetFields(request)
	if err != nil {
		return nil, err
	}
	var object map[string]any
	if err := json.Unmarshal(body, &object); err != nil {
		return nil, err
	}
	delete(object, "peer_id")
	return json.Marshal(object)
}

func stripPeerJobTargetFields(request json.RawMessage) (json.RawMessage, error) {
	body, err := stripTargetFields(request)
	if err != nil {
		return nil, err
	}
	var object map[string]any
	if err := json.Unmarshal(body, &object); err != nil {
		return nil, err
	}
	delete(object, "peer_id")
	return json.Marshal(object)
}

func prepareSecretRequest(operation string, request json.RawMessage, resolve func(string) (string, error)) (json.RawMessage, error) {
	if operation == "users.invites.accept" {
		var object map[string]any
		if err := json.Unmarshal(request, &object); err != nil {
			return nil, fmt.Errorf("decode invite acceptance request: %w", err)
		}
		if _, ok := object["invite_token"]; ok {
			return nil, errors.New("invite token cannot be persisted; use invite_token_ref")
		}
		if _, ok := object["password"]; ok {
			return nil, errors.New("invite password cannot be persisted; use password_ref")
		}
		ref, ok := object["password_ref"].(string)
		delete(object, "password_ref")
		delete(object, "invite_token_ref")
		if !ok || strings.TrimSpace(ref) == "" {
			return nil, errors.New("invite acceptance requires password_ref")
		}
		if resolve == nil {
			return nil, errors.New("invite password_ref requires a configured secret resolver")
		}
		secret, err := resolve(ref)
		if err != nil || strings.TrimSpace(secret) == "" {
			return nil, errors.New("invite password_ref could not be resolved")
		}
		object["password"] = secret
		return json.Marshal(object)
	}
	if operation == "users.password.update" {
		var object map[string]any
		if err := json.Unmarshal(request, &object); err != nil {
			return nil, fmt.Errorf("decode password request: %w", err)
		}
		for _, field := range []string{"old_password", "new_password"} {
			if _, ok := object[field]; ok {
				return nil, fmt.Errorf("password field %s cannot be persisted; use %s_ref", field, field)
			}
			ref, ok := object[field+"_ref"].(string)
			delete(object, field+"_ref")
			if !ok || strings.TrimSpace(ref) == "" {
				return nil, fmt.Errorf("password request requires %s_ref", field)
			}
			if resolve == nil {
				return nil, errors.New("password refs require a configured secret resolver")
			}
			secret, err := resolve(ref)
			if err != nil {
				return nil, fmt.Errorf("password ref %s could not be resolved", field)
			}
			if strings.TrimSpace(secret) == "" {
				return nil, fmt.Errorf("password ref %s resolved to an empty secret", field)
			}
			object[field] = secret
		}
		return json.Marshal(object)
	}
	if operation == "event_streaming.create" || operation == "event_streaming.update" {
		var object map[string]any
		if err := json.Unmarshal(request, &object); err != nil {
			return nil, fmt.Errorf("decode event-streaming request: %w", err)
		}
		if _, ok := object["config"]; ok {
			return nil, errors.New("event-streaming config cannot be persisted; use config_ref")
		}
		ref, hasRef := object["config_ref"].(string)
		delete(object, "config_ref")
		if !hasRef || strings.TrimSpace(ref) == "" {
			if operation == "event_streaming.create" {
				return nil, errors.New("event-streaming create requires config_ref")
			}
			return json.Marshal(object)
		}
		if resolve == nil {
			return nil, errors.New("event-streaming config_ref requires a configured secret resolver")
		}
		configJSON, err := resolve(ref)
		if err != nil || strings.TrimSpace(configJSON) == "" {
			return nil, errors.New("event-streaming config_ref could not be resolved")
		}
		var config any
		if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
			return nil, errors.New("event-streaming config_ref must resolve to JSON")
		}
		object["config"] = config
		return json.Marshal(object)
	}
	if operation == "identity_providers.create" || operation == "identity_providers.update" {
		var object map[string]any
		if err := json.Unmarshal(request, &object); err != nil {
			return nil, fmt.Errorf("decode identity provider request: %w", err)
		}
		if _, ok := object["client_secret"]; ok {
			return nil, errors.New("identity provider client_secret cannot be persisted; use client_secret_ref")
		}
		ref, hasRef := object["client_secret_ref"].(string)
		delete(object, "client_secret_ref")
		if !hasRef || strings.TrimSpace(ref) == "" {
			if operation == "identity_providers.create" {
				return nil, errors.New("identity provider create requires client_secret_ref")
			}
			return json.Marshal(object)
		}
		if resolve == nil {
			return nil, errors.New("identity provider client_secret_ref requires a configured secret resolver")
		}
		secret, err := resolve(ref)
		if err != nil || strings.TrimSpace(secret) == "" {
			return nil, errors.New("identity provider client_secret_ref could not be resolved")
		}
		object["client_secret"] = secret
		return json.Marshal(object)
	}
	if operation == "azure_idp.create" || operation == "azure_idp.update" {
		var object map[string]any
		if err := json.Unmarshal(request, &object); err != nil {
			return nil, fmt.Errorf("decode azure IDP request: %w", err)
		}
		if _, ok := object["client_secret"]; ok {
			return nil, errors.New("azure IDP client_secret cannot be persisted; use client_secret_ref")
		}
		ref, hasRef := object["client_secret_ref"].(string)
		delete(object, "client_secret_ref")
		if !hasRef || strings.TrimSpace(ref) == "" {
			if operation == "azure_idp.create" {
				return nil, errors.New("azure IDP create requires client_secret_ref")
			}
			return json.Marshal(object)
		}
		if resolve == nil {
			return nil, errors.New("azure IDP client_secret_ref requires a configured secret resolver")
		}
		secret, err := resolve(ref)
		if err != nil || strings.TrimSpace(secret) == "" {
			return nil, errors.New("azure IDP client_secret_ref could not be resolved")
		}
		object["client_secret"] = secret
		return json.Marshal(object)
	}
	if operation == "google_idp.create" || operation == "google_idp.update" {
		var object map[string]any
		if err := json.Unmarshal(request, &object); err != nil {
			return nil, fmt.Errorf("decode google IDP request: %w", err)
		}
		if _, ok := object["service_account_key"]; ok {
			return nil, errors.New("google IDP service_account_key cannot be persisted; use service_account_key_ref")
		}
		ref, hasRef := object["service_account_key_ref"].(string)
		delete(object, "service_account_key_ref")
		if !hasRef || strings.TrimSpace(ref) == "" {
			if operation == "google_idp.create" {
				return nil, errors.New("google IDP create requires service_account_key_ref")
			}
			return json.Marshal(object)
		}
		if resolve == nil {
			return nil, errors.New("google IDP service_account_key_ref requires a configured secret resolver")
		}
		secret, err := resolve(ref)
		if err != nil || strings.TrimSpace(secret) == "" {
			return nil, errors.New("google IDP service_account_key_ref could not be resolved")
		}
		object["service_account_key"] = secret
		return json.Marshal(object)
	}
	if strings.HasPrefix(operation, "edr.") && (strings.HasSuffix(operation, ".create") || strings.HasSuffix(operation, ".update")) {
		var object map[string]any
		if err := json.Unmarshal(request, &object); err != nil {
			return nil, fmt.Errorf("decode EDR request: %w", err)
		}
		provider := strings.TrimSuffix(strings.TrimPrefix(operation, "edr."), ".create")
		if strings.HasSuffix(operation, ".update") {
			provider = strings.TrimSuffix(strings.TrimPrefix(operation, "edr."), ".update")
		}
		fields := map[string]string{}
		switch provider {
		case "intune":
			fields["secret"] = "secret_ref"
		case "sentinelone", "fleetdm":
			fields["api_token"] = "api_token_ref"
		case "falcon":
			fields["secret"] = "secret_ref"
		case "huntress":
			fields["api_key"] = "api_key_ref"
			fields["api_secret"] = "api_secret_ref"
		default:
			return nil, fmt.Errorf("unsupported EDR operation %q", operation)
		}
		for field, refField := range fields {
			if _, ok := object[field]; ok {
				return nil, fmt.Errorf("EDR %s cannot be persisted; use %s", field, refField)
			}
			ref, ok := object[refField].(string)
			delete(object, refField)
			if !ok || strings.TrimSpace(ref) == "" {
				return nil, fmt.Errorf("EDR %s requires %s", provider, refField)
			}
			if resolve == nil {
				return nil, errors.New("EDR credential refs require a configured secret resolver")
			}
			secret, err := resolve(ref)
			if err != nil || strings.TrimSpace(secret) == "" {
				return nil, fmt.Errorf("EDR %s could not resolve %s", provider, refField)
			}
			object[field] = secret
		}
		return json.Marshal(object)
	}
	if operation == "reverse_proxy_services.create" || operation == "reverse_proxy_services.update" {
		var object map[string]any
		if err := json.Unmarshal(request, &object); err != nil {
			return nil, fmt.Errorf("decode reverse proxy service request: %w", err)
		}
		if _, ok := object["auth"]; ok {
			return nil, errors.New("reverse proxy service auth cannot be persisted; use auth_ref")
		}
		ref, hasRef := object["auth_ref"].(string)
		delete(object, "auth_ref")
		if !hasRef || strings.TrimSpace(ref) == "" {
			return json.Marshal(object)
		}
		if resolve == nil {
			return nil, errors.New("reverse proxy service auth_ref requires a configured secret resolver")
		}
		authJSON, err := resolve(ref)
		if err != nil || strings.TrimSpace(authJSON) == "" {
			return nil, errors.New("reverse proxy service auth_ref could not be resolved")
		}
		var auth any
		if err := json.Unmarshal([]byte(authJSON), &auth); err != nil {
			return nil, errors.New("reverse proxy service auth_ref must resolve to JSON")
		}
		object["auth"] = auth
		return json.Marshal(object)
	}
	if operation == "notification_channels.create" || operation == "notification_channels.update" {
		var object map[string]any
		if err := json.Unmarshal(request, &object); err != nil {
			return nil, fmt.Errorf("decode notification channel request: %w", err)
		}
		if _, ok := object["target"]; ok {
			return nil, errors.New("notification channel target cannot be persisted; use target_ref")
		}
		ref, hasRef := object["target_ref"].(string)
		delete(object, "target_ref")
		if !hasRef || strings.TrimSpace(ref) == "" {
			if operation == "notification_channels.create" {
				return nil, errors.New("notification channel create requires target_ref")
			}
			return json.Marshal(object)
		}
		if resolve == nil {
			return nil, errors.New("notification channel target_ref requires a configured secret resolver")
		}
		targetJSON, err := resolve(ref)
		if err != nil || strings.TrimSpace(targetJSON) == "" {
			return nil, errors.New("notification channel target_ref could not be resolved")
		}
		var target any
		if err := json.Unmarshal([]byte(targetJSON), &target); err != nil {
			return nil, errors.New("notification channel target_ref must resolve to JSON")
		}
		object["target"] = target
		return json.Marshal(object)
	}
	if operation != "agent_network.providers.create" && operation != "agent_network.providers.update" {
		return request, nil
	}
	var object map[string]any
	if err := json.Unmarshal(request, &object); err != nil {
		return nil, fmt.Errorf("decode provider request: %w", err)
	}
	if _, ok := object["api_key"]; ok {
		return nil, errors.New("provider api_key cannot be persisted; use api_key_ref")
	}
	ref, ok := object["api_key_ref"].(string)
	delete(object, "api_key_ref")
	if !ok || strings.TrimSpace(ref) == "" {
		if operation == "agent_network.providers.create" {
			return nil, errors.New("provider create requires api_key_ref")
		}
		return json.Marshal(object)
	}
	if resolve == nil {
		return nil, errors.New("provider api_key_ref requires a configured secret resolver")
	}
	secret, err := resolve(ref)
	if err != nil {
		return nil, errors.New("provider api_key_ref could not be resolved")
	}
	if strings.TrimSpace(secret) == "" {
		return nil, errors.New("provider api_key_ref resolved to an empty secret")
	}
	object["api_key"] = secret
	return json.Marshal(object)
}

func mutationImpact(operation string, before, intendedAfter json.RawMessage) (analysis.ImpactReport, error) {
	switch operation {
	case "groups.update":
		return analysis.GroupUpdateImpact(before, intendedAfter)
	case "groups.delete":
		return analysis.GroupDeleteImpact(before)
	case "groups.create":
		return analysis.GroupCreateImpact(intendedAfter)
	case "policies.update":
		return analysis.PolicyUpdateImpact(before, intendedAfter)
	case "policies.delete":
		return analysis.PolicyDeleteImpact(before)
	case "policies.create":
		return analysis.PolicyCreateImpact(intendedAfter)
	case "dns.zones.create":
		return analysis.DNSZoneCreateImpact(intendedAfter)
	case "dns.zones.delete":
		return analysis.DNSZoneDeleteImpact(before)
	case "dns.zones.update":
		return analysis.DNSZoneUpdateImpact(before, intendedAfter)
	case "dns.records.create":
		return analysis.DNSRecordCreateImpact(intendedAfter)
	case "dns.records.update":
		return analysis.DNSRecordUpdateImpact(before, intendedAfter)
	case "dns.records.delete":
		return analysis.DNSRecordDeleteImpact(before)
	case "dns.nameservers.create":
		return analysis.DNSNameserverCreateImpact(intendedAfter)
	case "dns.nameservers.update":
		return analysis.DNSNameserverUpdateImpact(before, intendedAfter)
	case "dns.nameservers.delete":
		return analysis.DNSNameserverDeleteImpact(before)
	case "dns.settings.update":
		return analysis.DNSSettingsUpdateImpact(before, intendedAfter)
	case "accounts.update":
		return analysis.AccountUpdateImpact(before, intendedAfter)
	case "accounts.delete":
		return analysis.AccountDeleteImpact(before)
	case "posture_checks.create":
		return analysis.PostureCheckCreateImpact(intendedAfter)
	case "posture_checks.update":
		return analysis.PostureCheckUpdateImpact(before, intendedAfter)
	case "posture_checks.delete":
		return analysis.PostureCheckDeleteImpact(before)
	case "ingress.peers.create":
		return analysis.IngressPeerCreateImpact(intendedAfter)
	case "ingress.peers.update":
		return analysis.IngressPeerUpdateImpact(before, intendedAfter)
	case "ingress.peers.delete":
		return analysis.IngressPeerDeleteImpact(before)
	case "peers.ingress.ports.create":
		return analysis.IngressPortAllocationCreateImpact(intendedAfter)
	case "peers.ingress.ports.update":
		return analysis.IngressPortAllocationUpdateImpact(before, intendedAfter)
	case "peers.ingress.ports.delete":
		return analysis.IngressPortAllocationDeleteImpact(before)
	case "agent_network.settings.update":
		return analysis.AgentNetworkSettingsUpdateImpact(before, intendedAfter)
	case "agent_network.settings.create":
		return analysis.AgentNetworkSettingsCreateImpact(intendedAfter)
	case "agent_network.settings.delete":
		return analysis.AgentNetworkSettingsDeleteImpact(before)
	case "agent_network.budget_rules.create":
		return analysis.AgentNetworkBudgetRuleCreateImpact(intendedAfter)
	case "agent_network.budget_rules.update":
		return analysis.AgentNetworkBudgetRuleUpdateImpact(before, intendedAfter)
	case "agent_network.budget_rules.delete":
		return analysis.AgentNetworkBudgetRuleDeleteImpact(before)
	case "agent_network.guardrails.create":
		return analysis.AgentNetworkGuardrailCreateImpact(intendedAfter)
	case "agent_network.guardrails.update":
		return analysis.AgentNetworkGuardrailUpdateImpact(before, intendedAfter)
	case "agent_network.guardrails.delete":
		return analysis.AgentNetworkGuardrailDeleteImpact(before)
	case "agent_network.policies.create":
		return analysis.AgentNetworkPolicyCreateImpact(intendedAfter)
	case "agent_network.policies.update":
		return analysis.AgentNetworkPolicyUpdateImpact(before, intendedAfter)
	case "agent_network.policies.delete":
		return analysis.AgentNetworkPolicyDeleteImpact(before)
	case "agent_network.providers.create":
		return analysis.AgentNetworkProviderCreateImpact(intendedAfter)
	case "agent_network.providers.update":
		return analysis.AgentNetworkProviderUpdateImpact(before, intendedAfter)
	case "agent_network.providers.delete":
		return analysis.AgentNetworkProviderDeleteImpact(before)
	case "users.create":
		return analysis.UserCreateImpact(intendedAfter)
	case "users.update":
		return analysis.UserUpdateImpact(before, intendedAfter)
	case "users.delete":
		return analysis.UserDeleteImpact(before)
	case "users.approve":
		return analysis.UserApproveImpact(before)
	case "users.reject":
		return analysis.UserRejectImpact(before)
	case "users.password.update":
		return analysis.UserPasswordUpdateImpact(before, intendedAfter)
	case "users.invite.resend":
		return analysis.UserInviteResendImpact(before, intendedAfter)
	case "users.tokens.delete":
		return analysis.UserTokenDeleteImpact(before)
	case "users.tokens.create":
		return analysis.UserTokenCreateImpact(intendedAfter)
	case "setup_keys.update":
		return analysis.SetupKeyUpdateImpact(before, intendedAfter)
	case "setup_keys.delete":
		return analysis.SetupKeyDeleteImpact(before)
	case "setup_keys.create":
		return analysis.SetupKeyCreateImpact(intendedAfter)
	case "users.invites.create":
		return analysis.InviteCreateImpact(intendedAfter)
	case "users.invites.regenerate":
		return analysis.InviteRegenerateImpact(before, intendedAfter)
	case "users.invites.delete":
		return analysis.InviteDeleteImpact(before)
	case "users.invites.accept":
		return analysis.InviteAcceptImpact(before, intendedAfter)
	case "routes.update":
		return analysis.RouteUpdateImpact(before, intendedAfter)
	case "routes.delete":
		return analysis.RouteDeleteImpact(before)
	case "routes.create":
		return analysis.RouteCreateImpact(intendedAfter)
	case "peers.update":
		return analysis.PeerUpdateImpact(before, intendedAfter)
	case "peers.temporary_access.create":
		return analysis.TemporaryAccessCreateImpact(before, intendedAfter)
	case "peers.jobs.create":
		return analysis.PeerJobCreateImpact(before, intendedAfter)
	case "event_streaming.create":
		return analysis.EventStreamingCreateImpact(intendedAfter)
	case "event_streaming.update":
		return analysis.EventStreamingUpdateImpact(before, intendedAfter)
	case "event_streaming.delete":
		return analysis.EventStreamingDeleteImpact(before)
	case "identity_providers.create":
		return analysis.IdentityProviderCreateImpact(intendedAfter)
	case "identity_providers.update":
		return analysis.IdentityProviderUpdateImpact(before, intendedAfter)
	case "identity_providers.delete":
		return analysis.IdentityProviderDeleteImpact(before)
	case "reverse_proxy_tokens.create":
		return analysis.ReverseProxyTokenCreateImpact(intendedAfter)
	case "reverse_proxy_tokens.delete":
		return analysis.ReverseProxyTokenDeleteImpact(before)
	case "reverse_proxy_domains.create":
		return analysis.ReverseProxyDomainCreateImpact(intendedAfter)
	case "reverse_proxy_domains.delete":
		return analysis.ReverseProxyDomainDeleteImpact(before)
	case "reverse_proxy_clusters.delete":
		return analysis.ReverseProxyClusterDeleteImpact(before)
	case "reverse_proxy_services.create":
		return analysis.ReverseProxyServiceCreateImpact(intendedAfter)
	case "reverse_proxy_services.update":
		return analysis.ReverseProxyServiceUpdateImpact(before, intendedAfter)
	case "reverse_proxy_services.delete":
		return analysis.ReverseProxyServiceDeleteImpact(before)
	case "notification_channels.create":
		return analysis.NotificationChannelCreateImpact(intendedAfter)
	case "notification_channels.update":
		return analysis.NotificationChannelUpdateImpact(before, intendedAfter)
	case "notification_channels.delete":
		return analysis.NotificationChannelDeleteImpact(before)
	case "azure_idp.create":
		return analysis.AzureIDPCreateImpact(intendedAfter)
	case "azure_idp.update":
		return analysis.AzureIDPUpdateImpact(before, intendedAfter)
	case "azure_idp.delete":
		return analysis.AzureIDPDeleteImpact(before)
	case "azure_idp.sync":
		return analysis.AzureIDPSyncImpact(before, intendedAfter)
	case "google_idp.create":
		return analysis.GoogleIDPCreateImpact(intendedAfter)
	case "google_idp.update":
		return analysis.GoogleIDPUpdateImpact(before, intendedAfter)
	case "google_idp.delete":
		return analysis.GoogleIDPDeleteImpact(before)
	case "google_idp.sync":
		return analysis.GoogleIDPSyncImpact(before, intendedAfter)
	case "edr.intune.create", "edr.sentinelone.create", "edr.falcon.create", "edr.huntress.create", "edr.fleetdm.create":
		return analysis.EDRIntegrationCreateImpact(strings.Split(operation, ".")[1], intendedAfter)
	case "edr.intune.update", "edr.sentinelone.update", "edr.falcon.update", "edr.huntress.update", "edr.fleetdm.update":
		return analysis.EDRIntegrationUpdateImpact(strings.Split(operation, ".")[1], before, intendedAfter)
	case "edr.intune.delete", "edr.sentinelone.delete", "edr.falcon.delete", "edr.huntress.delete", "edr.fleetdm.delete":
		return analysis.EDRIntegrationDeleteImpact(strings.Split(operation, ".")[1], before)
	case "scim.create", "okta_scim.create":
		return analysis.SCIMCreateImpact(strings.Split(operation, ".")[0], intendedAfter)
	case "scim.update", "okta_scim.update":
		return analysis.SCIMUpdateImpact(strings.Split(operation, ".")[0], before, intendedAfter)
	case "scim.delete", "okta_scim.delete":
		return analysis.SCIMDeleteImpact(strings.Split(operation, ".")[0], before)
	case "scim.token", "okta_scim.token":
		return analysis.SCIMTokenImpact(strings.Split(operation, ".")[0], before)
	case "peers.delete":
		return analysis.PeerDeleteImpact(before)
	case "peers.edr.bypass.create":
		return analysis.EDRBypassCreateImpact(before, intendedAfter)
	case "peers.edr.bypass.delete":
		return analysis.EDRBypassDeleteImpact(before)
	case "networks.update":
		return analysis.NetworkUpdateImpact(before, intendedAfter)
	case "networks.delete":
		return analysis.NetworkDeleteImpact(before)
	case "networks.create":
		return analysis.NetworkCreateImpact(intendedAfter)
	case "networks.resources.delete":
		return analysis.NetworkResourceDeleteImpact(before)
	case "networks.resources.update":
		return analysis.NetworkResourceUpdateImpact(before, intendedAfter)
	case "networks.resources.create":
		return analysis.NetworkResourceCreateImpact(intendedAfter)
	case "networks.routers.create":
		return analysis.NetworkRouterCreateImpact(intendedAfter)
	case "networks.routers.delete":
		return analysis.NetworkRouterDeleteImpact(before)
	case "networks.routers.update":
		return analysis.NetworkRouterUpdateImpact(before, intendedAfter)
	default:
		return analysis.ImpactReport{}, fmt.Errorf("operation %q has no impact analyzer", operation)
	}
}

func confirmDeleted(ctx context.Context, remote Remote, operation string, target requestTarget) error {
	_, err := readPreimage(ctx, remote, operation, target)
	if err == nil {
		return errors.New("resource still exists after delete")
	}
	return err
}

func isNotFound(err error) bool {
	var status interface{ StatusCodeState() int }
	return errors.As(err, &status) && status.StatusCodeState() == 404
}

func isDeleteOperation(operation string) bool {
	return operation == "groups.delete" || operation == "policies.delete" || operation == "routes.delete" || operation == "peers.delete" || operation == "peers.edr.bypass.delete" || operation == "networks.delete" || operation == "networks.resources.delete" || operation == "networks.routers.delete" || operation == "dns.zones.delete" || operation == "dns.records.delete" || operation == "dns.nameservers.delete" || operation == "accounts.delete" || operation == "posture_checks.delete" || operation == "ingress.peers.delete" || operation == "peers.ingress.ports.delete" || operation == "agent_network.settings.delete" || operation == "agent_network.budget_rules.delete" || operation == "agent_network.guardrails.delete" || operation == "agent_network.policies.delete" || operation == "agent_network.providers.delete" || operation == "users.delete" || operation == "users.reject" || operation == "users.tokens.delete" || operation == "setup_keys.delete" || operation == "event_streaming.delete" || operation == "identity_providers.delete" || operation == "reverse_proxy_tokens.delete" || operation == "reverse_proxy_domains.delete" || operation == "reverse_proxy_clusters.delete" || operation == "reverse_proxy_services.delete" || operation == "notification_channels.delete" || operation == "azure_idp.delete" || operation == "google_idp.delete" || operation == "users.invites.delete" || operation == "scim.delete" || operation == "okta_scim.delete" || (strings.HasPrefix(operation, "edr.") && strings.HasSuffix(operation, ".delete"))
}

func classifyDispatchError(err error) mutation.DispatchState {
	var dispatched interface{ DispatchedState() bool }
	if errors.As(err, &dispatched) {
		if !dispatched.DispatchedState() {
			return mutation.NotDispatched
		}
		var status interface{ StatusCodeState() int }
		if errors.As(err, &status) && status.StatusCodeState() >= 400 && status.StatusCodeState() < 500 {
			return mutation.DefinitivelyRejected
		}
	}
	return mutation.Unknown
}

func finish(ctx context.Context, store Ledger, result Result, state mutation.DispatchState, reason string) (Result, error) {
	result.State = state
	result.Reason = reason
	if err := store.SetAttemptState(ctx, result.AttemptID, string(state)); err != nil {
		return result, &ApplyError{Result: result, Err: err}
	}
	encoded, err := json.Marshal(result)
	if err != nil {
		return result, &ApplyError{Result: result, Err: err}
	}
	if err := store.RecordReceipt(ctx, ledger.Receipt{AttemptID: result.AttemptID, StageID: result.StageID, Revision: result.Revision, State: string(state), Result: encoded}); err != nil {
		if state == mutation.ConfirmedSuccess || state == mutation.AlreadySatisfied {
			result.State = mutation.EffectConfirmedReceiptFail
		}
		return result, &ApplyError{Result: result, Err: fmt.Errorf("persist mutation receipt: %w", err)}
	}
	if state == mutation.ConfirmedSuccess || state == mutation.AlreadySatisfied {
		return result, nil
	}
	return result, &ApplyError{Result: result, Err: errors.New(reason)}
}
