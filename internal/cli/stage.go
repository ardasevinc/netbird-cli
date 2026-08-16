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
			if plan.Operation == "groups.update" {
				report, err := analysis.GroupUpdateImpact(plan.Before, plan.IntendedAfter)
				if err != nil {
					return fail(2, err)
				}
				impact, err = json.Marshal(report)
				if err != nil {
					return fail(1, fmt.Errorf("encode mutation impact: %w", err))
				}
				if report.Classification == "unknown" && !hasFinding(findings, "impact.unknown") {
					findings = append(findings, ledger.Finding{Code: "impact.unknown", Severity: "blocking", Message: "the proposed group change may affect reachability, but its impact cannot be calculated"})
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
