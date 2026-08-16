package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func accessiblePeersCommand(state *commandState, stdout io.Writer) *cobra.Command {
	return &cobra.Command{Use: "accessible <peer-id>", Args: cobra.ExactArgs(1), Short: "list peers reachable from a peer", RunE: func(cmd *cobra.Command, args []string) error {
		client, err := managementClient(state)
		if err != nil {
			return err
		}
		peers, err := client.ListAccessiblePeers(cmd.Context(), args[0])
		if err != nil {
			return fail(3, err)
		}
		data := map[string]any{"peer_id": args[0], "peers": peers, "completeness": completeCompleteness()}
		if state.json {
			return writeJSON(stdout, map[string]any{"schema": "nb/v1/peers-accessible-result", "ok": true, "operation": "peers.accessible", "data": data})
		}
		for _, peer := range peers {
			if _, err := fmt.Fprintf(stdout, "%s\t%s\t%s\tconnected=%t\n", peer.ID, peer.Name, peer.IP, peer.Connected); err != nil {
				return err
			}
		}
		return nil
	}}
}
