package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func networksCommand(state *commandState, stdout io.Writer) *cobra.Command {
	command := &cobra.Command{Use: "networks", Short: "inspect NetBird routed networks"}
	command.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "list networks with bounded completeness semantics",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := managementClient(state)
			if err != nil {
				return err
			}
			networks, err := client.ListNetworks(cmd.Context())
			if err != nil {
				return fail(3, err)
			}
			data := map[string]any{"networks": networks, "completeness": completeCompleteness()}
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/networks-list-result", "ok": true, "operation": "networks.list", "data": data})
			}
			for _, network := range networks {
				if _, err := fmt.Fprintf(stdout, "%s\t%s\tresources=%d\trouters=%d\n", network.ID, network.Name, len(network.Resources), len(network.Routers)); err != nil {
					return err
				}
			}
			return nil
		},
	})
	command.AddCommand(&cobra.Command{
		Use:   "get <network-id>",
		Short: "show one network",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := managementClient(state)
			if err != nil {
				return err
			}
			network, err := client.GetNetwork(cmd.Context(), args[0])
			if err != nil {
				return fail(3, err)
			}
			data := map[string]any{"network": network, "completeness": completeCompleteness()}
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/networks-get-result", "ok": true, "operation": "networks.get", "data": data})
			}
			_, err = fmt.Fprintf(stdout, "%s\t%s\tresources=%d\trouters=%d\n", network.ID, network.Name, len(network.Resources), len(network.Routers))
			return err
		},
	})
	return command
}
