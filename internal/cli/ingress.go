package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func ingressCommand(state *commandState, stdout io.Writer) *cobra.Command {
	command := &cobra.Command{Use: "ingress", Short: "inspect ingress peers"}
	command.AddCommand(&cobra.Command{Use: "list", Short: "list ingress peers", RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := managementClient(state)
		if err != nil {
			return err
		}
		items, err := client.ListIngressPeers(cmd.Context())
		if err != nil {
			return fail(3, err)
		}
		data := map[string]any{"ingress_peers": items, "completeness": completeCompleteness()}
		if state.json {
			return writeJSON(stdout, map[string]any{"schema": "nb/v1/ingress-peers-list-result", "ok": true, "operation": "ingress.peers.list", "data": data})
		}
		for _, item := range items {
			if _, err := fmt.Fprintf(stdout, "%s\t%s\t%s\tconnected=%t\tenabled=%t\n", item.ID, item.PeerID, item.Region, item.Connected, item.Enabled); err != nil {
				return err
			}
		}
		return nil
	}}, &cobra.Command{Use: "get <ingress-peer-id>", Args: cobra.ExactArgs(1), Short: "show one ingress peer", RunE: func(cmd *cobra.Command, args []string) error {
		client, err := managementClient(state)
		if err != nil {
			return err
		}
		item, err := client.GetIngressPeer(cmd.Context(), args[0])
		if err != nil {
			return fail(3, err)
		}
		data := map[string]any{"ingress_peer": item, "completeness": completeCompleteness()}
		if state.json {
			return writeJSON(stdout, map[string]any{"schema": "nb/v1/ingress-peers-get-result", "ok": true, "operation": "ingress.peers.get", "data": data})
		}
		_, err = fmt.Fprintf(stdout, "%s\t%s\t%s\tconnected=%t\tenabled=%t\n", item.ID, item.PeerID, item.Region, item.Connected, item.Enabled)
		return err
	}})
	return command
}

func peerIngressPortsCommand(state *commandState, stdout io.Writer) *cobra.Command {
	command := &cobra.Command{Use: "ingress-ports", Short: "inspect peer ingress port allocations"}
	command.AddCommand(&cobra.Command{Use: "list <peer-id>", Args: cobra.ExactArgs(1), Short: "list port allocations", RunE: func(cmd *cobra.Command, args []string) error {
		client, err := managementClient(state)
		if err != nil {
			return err
		}
		items, err := client.ListIngressPortAllocations(cmd.Context(), args[0])
		if err != nil {
			return fail(3, err)
		}
		data := map[string]any{"peer_id": args[0], "allocations": items, "completeness": completeCompleteness()}
		if state.json {
			return writeJSON(stdout, map[string]any{"schema": "nb/v1/peer-ingress-ports-list-result", "ok": true, "operation": "peers.ingress.ports.list", "data": data})
		}
		for _, item := range items {
			if _, err := fmt.Fprintf(stdout, "%s\t%s\t%s\tenabled=%t\tmappings=%d\n", item.ID, item.Name, item.Region, item.Enabled, len(item.PortRangeMappings)); err != nil {
				return err
			}
		}
		return nil
	}}, &cobra.Command{Use: "get <peer-id> <allocation-id>", Args: cobra.ExactArgs(2), Short: "show one port allocation", RunE: func(cmd *cobra.Command, args []string) error {
		client, err := managementClient(state)
		if err != nil {
			return err
		}
		item, err := client.GetIngressPortAllocation(cmd.Context(), args[0], args[1])
		if err != nil {
			return fail(3, err)
		}
		data := map[string]any{"peer_id": args[0], "allocation": item, "completeness": completeCompleteness()}
		if state.json {
			return writeJSON(stdout, map[string]any{"schema": "nb/v1/peer-ingress-ports-get-result", "ok": true, "operation": "peers.ingress.ports.get", "data": data})
		}
		_, err = fmt.Fprintf(stdout, "%s\t%s\t%s\tenabled=%t\tmappings=%d\n", item.ID, item.Name, item.Region, item.Enabled, len(item.PortRangeMappings))
		return err
	}})
	return command
}
