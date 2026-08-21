package cli

import (
	"fmt"
	"io"

	"github.com/ardasevinc/netbird-cli/internal/analysis"
	"github.com/spf13/cobra"
)

func analyzeCommand(state *commandState, stdout io.Writer) *cobra.Command {
	command := &cobra.Command{Use: "analyze", Short: "explain management-plane topology and configured access"}
	command.AddCommand(&cobra.Command{
		Use:   "reachability <peer-id>",
		Short: "deprecated: show the legacy network-map and policy-evidence projection",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !state.json {
				_, _ = fmt.Fprintln(state.logWriter, "deprecated: use 'nb analyze access <peer-id>'; the legacy result does not prove directional initiation or packet reachability")
			}
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
			if _, err := fmt.Fprintf(stdout, "source: %s\t%s\nnetwork-map peers (legacy reachable_peers): %d\n", report.SourcePeer.ID, report.SourcePeer.Name, report.Summary.ReachablePeerCount); err != nil {
				return err
			}
			for _, peer := range report.ReachablePeers {
				if _, err := fmt.Fprintf(stdout, "%s\t%s\t%s\tconnected=%t\n", peer.ID, peer.Name, peer.IP, peer.Connected); err != nil {
					return err
				}
			}
			if _, err := fmt.Fprintf(stdout, "legacy policy evidence: %d\nunattributed network-map peers (legacy unexplained): %d\n", report.Summary.PolicyEvidenceCount, report.Summary.UnexplainedReachablePeerCount); err != nil {
				return err
			}
			return nil
		},
	})
	command.AddCommand(&cobra.Command{
		Use:   "access <peer-id>",
		Short: "separate network-map adjacency from configured directional policy flows",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := managementClient(state)
			if err != nil {
				return err
			}
			report, err := analysis.Access(cmd.Context(), client, args[0])
			if err != nil {
				return fail(3, err)
			}
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/access-analysis-result", "ok": true, "operation": "analyze.access", "data": report})
			}
			if _, err := fmt.Fprintf(stdout, "subject: %s\t%s\nnetwork-map peers: %d\n", report.SubjectPeer.ID, report.SubjectPeer.Name, report.Summary.NetworkMapPeerCount); err != nil {
				return err
			}
			for _, relation := range report.Relations {
				if _, err := fmt.Fprintf(stdout, "%s\t%s\tadjacent=%t\toutbound_flows=%d\tinbound_flows=%d\n", relation.Peer.ID, relation.Peer.Name, relation.NetworkMapAdjacent, len(relation.OutboundFlows), len(relation.InboundFlows)); err != nil {
					return err
				}
			}
			_, err = fmt.Fprintf(stdout, "map-only peers: %d\nconfigured peers missing from map: %d\nobservations: %d (%s)\n", report.Summary.MapOnlyPeerCount, report.Summary.ConfiguredPeerMissingFromNetworkMapCount, report.Summary.ObservationCount, report.Completeness.Observations.State)
			return err
		},
	})
	return command
}
