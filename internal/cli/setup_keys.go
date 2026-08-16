package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func setupKeysCommand(state *commandState, stdout io.Writer) *cobra.Command {
	command := &cobra.Command{Use: "setup-keys", Short: "inspect NetBird setup keys without exposing secrets"}
	command.AddCommand(setupKeysListCommand(state, stdout), setupKeysGetCommand(state, stdout))
	return command
}

func setupKeysListCommand(state *commandState, stdout io.Writer) *cobra.Command {
	return &cobra.Command{Use: "list", Short: "list setup keys", RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := managementClient(state)
		if err != nil {
			return err
		}
		keys, err := client.ListSetupKeys(cmd.Context())
		if err != nil {
			return fail(3, err)
		}
		data := map[string]any{"setup_keys": keys, "completeness": completeCompleteness()}
		if state.json {
			return writeJSON(stdout, map[string]any{"schema": "nb/v1/setup-keys-list-result", "ok": true, "operation": "setup_keys.list", "data": data})
		}
		for _, key := range keys {
			if _, err := fmt.Fprintf(stdout, "%s\t%s\t%s\tvalid=%t\tused=%d/%d\n", key.ID, key.Name, key.State, key.Valid, key.UsedTimes, key.UsageLimit); err != nil {
				return err
			}
		}
		return nil
	}}
}

func setupKeysGetCommand(state *commandState, stdout io.Writer) *cobra.Command {
	return &cobra.Command{Use: "get <setup-key-id>", Args: cobra.ExactArgs(1), Short: "show setup-key metadata without the secret", RunE: func(cmd *cobra.Command, args []string) error {
		client, err := managementClient(state)
		if err != nil {
			return err
		}
		key, err := client.GetSetupKey(cmd.Context(), args[0])
		if err != nil {
			return fail(3, err)
		}
		data := map[string]any{"setup_key": key, "completeness": completeCompleteness()}
		if state.json {
			return writeJSON(stdout, map[string]any{"schema": "nb/v1/setup-keys-get-result", "ok": true, "operation": "setup_keys.get", "data": data})
		}
		_, err = fmt.Fprintf(stdout, "%s\t%s\t%s\tvalid=%t\tused=%d/%d\n", key.ID, key.Name, key.State, key.Valid, key.UsedTimes, key.UsageLimit)
		return err
	}}
}
