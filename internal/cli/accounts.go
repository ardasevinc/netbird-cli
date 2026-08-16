package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func accountsCommand(state *commandState, stdout io.Writer) *cobra.Command {
	command := &cobra.Command{Use: "accounts", Short: "inspect NetBird accounts"}
	command.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "list visible accounts",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := managementClient(state)
			if err != nil {
				return err
			}
			accounts, err := client.ListAccounts(cmd.Context())
			if err != nil {
				return fail(3, err)
			}
			data := map[string]any{"accounts": accounts, "completeness": completeCompleteness()}
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/accounts-list-result", "ok": true, "operation": "accounts.list", "data": data})
			}
			for _, account := range accounts {
				if _, err := fmt.Fprintf(stdout, "%s\t%s\t%s\n", account.ID, account.Domain, account.DomainCategory); err != nil {
					return err
				}
			}
			return nil
		},
	})
	return command
}

func completeCompleteness() map[string]any {
	return map[string]any{"state": "complete", "reason": nil}
}
