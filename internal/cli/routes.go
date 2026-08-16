package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func routesCommand(state *commandState, stdout io.Writer) *cobra.Command {
	command := &cobra.Command{Use: "routes", Short: "inspect NetBird network routes"}
	command.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "list routes with bounded completeness semantics",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := managementClient(state)
			if err != nil {
				return err
			}
			routes, err := client.ListRoutes(cmd.Context())
			if err != nil {
				return fail(3, err)
			}
			data := map[string]any{"routes": routes, "completeness": completeCompleteness()}
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/routes-list-result", "ok": true, "operation": "routes.list", "data": data})
			}
			for _, route := range routes {
				if _, err := fmt.Fprintf(stdout, "%s\t%s\tenabled=%t\tmetric=%d\n", route.ID, pointerValue(route.Network), route.Enabled, route.Metric); err != nil {
					return err
				}
			}
			return nil
		},
	})
	command.AddCommand(&cobra.Command{
		Use:   "get <route-id>",
		Short: "show one route",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := managementClient(state)
			if err != nil {
				return err
			}
			route, err := client.GetRoute(cmd.Context(), args[0])
			if err != nil {
				return fail(3, err)
			}
			data := map[string]any{"route": route, "completeness": completeCompleteness()}
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/routes-get-result", "ok": true, "operation": "routes.get", "data": data})
			}
			_, err = fmt.Fprintf(stdout, "%s\t%s\tenabled=%t\tmetric=%d\n", route.ID, pointerValue(route.Network), route.Enabled, route.Metric)
			return err
		},
	})
	return command
}
