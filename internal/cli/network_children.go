package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func networkResourcesCommand(state *commandState, stdout io.Writer) *cobra.Command {
	command := &cobra.Command{Use: "resources", Short: "inspect network resources"}
	command.AddCommand(&cobra.Command{Use: "list <network-id>", Args: cobra.ExactArgs(1), Short: "list resources", RunE: func(cmd *cobra.Command, args []string) error {
		client, err := managementClient(state)
		if err != nil {
			return err
		}
		items, err := client.ListNetworkResources(cmd.Context(), args[0])
		if err != nil {
			return fail(3, err)
		}
		data := map[string]any{"network_id": args[0], "resources": items, "completeness": completeCompleteness()}
		if state.json {
			return writeJSON(stdout, map[string]any{"schema": "nb/v1/network-resources-list-result", "ok": true, "operation": "networks.resources.list", "data": data})
		}
		for _, item := range items {
			if _, err := fmt.Fprintf(stdout, "%s\t%s\t%s\tenabled=%t\n", item.ID, item.Name, item.Address, item.Enabled); err != nil {
				return err
			}
		}
		return nil
	}}, &cobra.Command{Use: "get <network-id> <resource-id>", Args: cobra.ExactArgs(2), Short: "show one resource", RunE: func(cmd *cobra.Command, args []string) error {
		client, err := managementClient(state)
		if err != nil {
			return err
		}
		item, err := client.GetNetworkResource(cmd.Context(), args[0], args[1])
		if err != nil {
			return fail(3, err)
		}
		data := map[string]any{"network_id": args[0], "resource": item, "completeness": completeCompleteness()}
		if state.json {
			return writeJSON(stdout, map[string]any{"schema": "nb/v1/network-resources-get-result", "ok": true, "operation": "networks.resources.get", "data": data})
		}
		_, err = fmt.Fprintf(stdout, "%s\t%s\t%s\tenabled=%t\n", item.ID, item.Name, item.Address, item.Enabled)
		return err
	}})
	return command
}

func networkRoutersCommand(state *commandState, stdout io.Writer) *cobra.Command {
	command := &cobra.Command{Use: "routers", Short: "inspect network routers"}
	command.AddCommand(&cobra.Command{Use: "list <network-id>", Args: cobra.ExactArgs(1), Short: "list routers", RunE: func(cmd *cobra.Command, args []string) error {
		client, err := managementClient(state)
		if err != nil {
			return err
		}
		items, err := client.ListNetworkRouters(cmd.Context(), args[0])
		if err != nil {
			return fail(3, err)
		}
		data := map[string]any{"network_id": args[0], "routers": items, "completeness": completeCompleteness()}
		if state.json {
			return writeJSON(stdout, map[string]any{"schema": "nb/v1/network-routers-list-result", "ok": true, "operation": "networks.routers.list", "data": data})
		}
		for _, item := range items {
			if _, err := fmt.Fprintf(stdout, "%s\tenabled=%t\tmetric=%d\n", item.ID, item.Enabled, item.Metric); err != nil {
				return err
			}
		}
		return nil
	}}, &cobra.Command{Use: "get <network-id> <router-id>", Args: cobra.ExactArgs(2), Short: "show one router", RunE: func(cmd *cobra.Command, args []string) error {
		client, err := managementClient(state)
		if err != nil {
			return err
		}
		item, err := client.GetNetworkRouter(cmd.Context(), args[0], args[1])
		if err != nil {
			return fail(3, err)
		}
		data := map[string]any{"network_id": args[0], "router": item, "completeness": completeCompleteness()}
		if state.json {
			return writeJSON(stdout, map[string]any{"schema": "nb/v1/network-routers-get-result", "ok": true, "operation": "networks.routers.get", "data": data})
		}
		_, err = fmt.Fprintf(stdout, "%s\tenabled=%t\tmetric=%d\n", item.ID, item.Enabled, item.Metric)
		return err
	}}, &cobra.Command{Use: "list-all", Short: "list routers across all networks", RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := managementClient(state)
		if err != nil {
			return err
		}
		items, err := client.ListAllNetworkRouters(cmd.Context())
		if err != nil {
			return fail(3, err)
		}
		data := map[string]any{"routers": items, "completeness": completeCompleteness()}
		if state.json {
			return writeJSON(stdout, map[string]any{"schema": "nb/v1/network-routers-all-result", "ok": true, "operation": "networks.routers.list_all", "data": data})
		}
		for _, item := range items {
			if _, err := fmt.Fprintf(stdout, "%s\tenabled=%t\tmetric=%d\n", item.ID, item.Enabled, item.Metric); err != nil {
				return err
			}
		}
		return nil
	}})
	return command
}
