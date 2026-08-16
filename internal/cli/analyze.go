package cli

import (
	"fmt"
	"io"

	"github.com/ardasevinc/netbird-cli/internal/analysis"
	"github.com/spf13/cobra"
)

func analyzeCommand(state *commandState, stdout io.Writer) *cobra.Command {
	command := &cobra.Command{Use: "analyze", Short: "explain management-plane topology and reachability"}
	command.AddCommand(&cobra.Command{
		Use:   "reachability <peer-id>",
		Short: "explain a peer's current reachable peers and policy evidence",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := managementClient(state)
			if err != nil {
				return err
			}
			report, err := analysis.Reachability(cmd.Context(), client, args[0])
			if err != nil {
				return fail(3, err)
			}
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/reachability-analysis-result", "ok": true, "operation": "analyze.reachability", "data": report})
			}
			if _, err := fmt.Fprintf(stdout, "source: %s\t%s\nreachable peers: %d\n", report.SourcePeer.ID, report.SourcePeer.Name, report.Summary.ReachablePeerCount); err != nil {
				return err
			}
			for _, peer := range report.ReachablePeers {
				if _, err := fmt.Fprintf(stdout, "%s\t%s\t%s\tconnected=%t\n", peer.ID, peer.Name, peer.IP, peer.Connected); err != nil {
					return err
				}
			}
			if _, err := fmt.Fprintf(stdout, "policy evidence: %d\nunexplained reachable peers: %d\n", report.Summary.PolicyEvidenceCount, report.Summary.UnexplainedReachablePeerCount); err != nil {
				return err
			}
			return nil
		},
	})
	return command
}
