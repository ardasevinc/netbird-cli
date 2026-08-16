package cli

import (
	"fmt"
	"io"

	"github.com/ardasevinc/netbird-cli/internal/config"
	"github.com/ardasevinc/netbird-cli/internal/exit"
	"github.com/ardasevinc/netbird-cli/internal/ledger"
	"github.com/ardasevinc/netbird-cli/internal/mutation"
	"github.com/ardasevinc/netbird-cli/internal/mutationengine"
	"github.com/ardasevinc/netbird-cli/internal/netbird"
	"github.com/ardasevinc/netbird-cli/internal/transport"
	"github.com/spf13/cobra"
)

func applyCommand(state *commandState, stdout io.Writer) *cobra.Command {
	var acknowledgements []string
	var ackAll bool
	command := &cobra.Command{
		Use:   "apply <stage-id>@<revision>",
		Short: "dispatch one exact staged mutation revision",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stageID, revision, err := parseRevision(args[0])
			if err != nil {
				return fail(int(exit.InvalidInput), err)
			}
			file, err := config.Load(state.configPath)
			if err != nil {
				return fail(int(exit.PreDispatch), err)
			}
			profile, err := file.Profile(state.profileName)
			if err != nil {
				return fail(int(exit.InvalidInput), err)
			}
			if profile.ServerIdentity == "" || profile.AccountID == "" {
				return fail(int(exit.SafetyConflict), fmt.Errorf("apply requires profile server_identity and account_id"))
			}
			token, err := config.ResolveCredential(profile.CredentialRef)
			if err != nil {
				return fail(int(exit.PreDispatch), err)
			}
			transportClient, err := transport.New(transport.Config{BaseURL: profile.URL, Token: token, CAFile: profile.CAFile})
			if err != nil {
				return fail(int(exit.PreDispatch), err)
			}
			store, err := ledger.Open(state.statePath)
			if err != nil {
				return fail(int(exit.PreDispatch), err)
			}
			defer store.Close()
			result, err := mutationengine.Apply(cmd.Context(), store, netbird.NewClient(transportClient), mutationengine.ApplyInput{
				StageID:          stageID,
				Revision:         revision,
				Profile:          state.profileName,
				ServerIdentity:   profile.ServerIdentity,
				AccountID:        profile.AccountID,
				Acknowledgements: acknowledgements,
				AckAllBlocking:   ackAll,
			})
			if err != nil {
				return fail(applyExit(result.State), err)
			}
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/apply-result", "ok": true, "operation": "apply", "data": result})
			}
			_, err = fmt.Fprintf(stdout, "apply %s@%d\nstate: %s\nattempt: %s\nreason: %s\n", result.StageID, result.Revision, result.State, result.AttemptID, result.Reason)
			return err
		},
	}
	command.Flags().StringArrayVar(&acknowledgements, "ack", nil, "acknowledge one exact blocking finding code")
	command.Flags().BoolVar(&ackAll, "ack-all-blocking", false, "acknowledge all blocking findings on this exact stage revision")
	return command
}

func applyExit(state mutation.DispatchState) int {
	switch state {
	case mutation.NotDispatched:
		return int(exit.PreDispatch)
	case mutation.DefinitivelyRejected:
		return int(exit.Rejected)
	case mutation.Partial, mutation.Unknown:
		return int(exit.Uncertain)
	case mutation.EffectConfirmedReceiptFail:
		return int(exit.ReceiptDelivery)
	default:
		return int(exit.SafetyConflict)
	}
}
