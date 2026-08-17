package cli

import (
	"encoding/json"
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
			impact := json.RawMessage(`{}`)
			findings := append([]ledger.Finding(nil), plan.Findings...)
			switch plan.Operation {
			case "groups.create", "groups.update", "groups.delete", "policies.create", "policies.update", "policies.delete", "routes.create", "routes.update", "routes.delete", "peers.update", "peers.delete", "networks.create", "networks.update", "networks.delete", "networks.resources.create", "networks.resources.update", "networks.resources.delete", "networks.routers.create", "networks.routers.update", "networks.routers.delete", "dns.zones.create", "dns.zones.update", "dns.zones.delete", "dns.records.create", "dns.records.update", "dns.records.delete":
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
			data := map[string]any{"stage_id": stage.ID, "revision": stage.Revision, "digest": stage.Digest, "applicability": "local_stage_only", "operation": stage.Operation, "mutation": definition.Mutation, "dispatcher_admitted": definition.DispatcherAdmitted, "impact": json.RawMessage(stage.Impact), "findings": stage.Findings}
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
			data := map[string]any{"stage_id": stage.ID, "revision": stage.Revision, "profile": stage.Profile, "server_identity": stage.ServerIdentity, "account_id": stage.AccountID, "operation": stage.Operation, "request": json.RawMessage(stage.Request), "before": json.RawMessage(stage.Before), "intended_after": json.RawMessage(stage.IntendedAfter), "impact": json.RawMessage(stage.Impact), "digest": stage.Digest, "findings": stage.Findings, "cancelled": stage.Cancelled, "created_at": stage.CreatedAt}
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
