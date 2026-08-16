package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func identityProvidersCommand(state *commandState, stdout io.Writer) *cobra.Command {
	command := &cobra.Command{Use: "identity-providers", Short: "inspect identity providers"}
	command.AddCommand(identityProvidersListCommand(state, stdout), identityProvidersGetCommand(state, stdout))
	return command
}

func identityProvidersListCommand(state *commandState, stdout io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "list identity providers",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := managementClient(state)
			if err != nil {
				return err
			}
			items, err := client.ListIdentityProviders(cmd.Context())
			if err != nil {
				return fail(3, err)
			}
			data := map[string]any{"identity_providers": items, "completeness": completeCompleteness()}
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/identity-providers-list-result", "ok": true, "operation": "identity_providers.list", "data": data})
			}
			for _, item := range items {
				if _, err := fmt.Fprintf(stdout, "%s\t%s\t%s\n", pointerValue(item.ID), item.Name, item.Type); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func identityProvidersGetCommand(state *commandState, stdout io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "get <identity-provider-id>",
		Short: "show one identity provider",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := managementClient(state)
			if err != nil {
				return err
			}
			item, err := client.GetIdentityProvider(cmd.Context(), args[0])
			if err != nil {
				return fail(3, err)
			}
			data := map[string]any{"identity_provider": item, "completeness": completeCompleteness()}
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/identity-providers-get-result", "ok": true, "operation": "identity_providers.get", "data": data})
			}
			_, err = fmt.Fprintf(stdout, "%s\t%s\t%s\n", pointerValue(item.ID), item.Name, item.Type)
			return err
		},
	}
}

func postureChecksCommand(state *commandState, stdout io.Writer) *cobra.Command {
	command := &cobra.Command{Use: "posture-checks", Short: "inspect posture checks"}
	command.AddCommand(postureChecksListCommand(state, stdout), postureChecksGetCommand(state, stdout))
	return command
}

func postureChecksListCommand(state *commandState, stdout io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "list posture checks",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := managementClient(state)
			if err != nil {
				return err
			}
			items, err := client.ListPostureChecks(cmd.Context())
			if err != nil {
				return fail(3, err)
			}
			data := map[string]any{"posture_checks": items, "completeness": completeCompleteness()}
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/posture-checks-list-result", "ok": true, "operation": "posture_checks.list", "data": data})
			}
			for _, item := range items {
				if _, err := fmt.Fprintf(stdout, "%s\t%s\tchecks=%d\n", item.ID, item.Name, len(item.Checks)); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func postureChecksGetCommand(state *commandState, stdout io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "get <posture-check-id>",
		Short: "show one posture check",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := managementClient(state)
			if err != nil {
				return err
			}
			item, err := client.GetPostureCheck(cmd.Context(), args[0])
			if err != nil {
				return fail(3, err)
			}
			data := map[string]any{"posture_check": item, "completeness": completeCompleteness()}
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/posture-checks-get-result", "ok": true, "operation": "posture_checks.get", "data": data})
			}
			_, err = fmt.Fprintf(stdout, "%s\t%s\tchecks=%d\n", item.ID, item.Name, len(item.Checks))
			return err
		},
	}
}
