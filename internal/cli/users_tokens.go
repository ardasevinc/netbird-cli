package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func userTokensCommand(state *commandState, stdout io.Writer) *cobra.Command {
	command := &cobra.Command{Use: "tokens", Short: "inspect user token metadata without token values"}
	command.AddCommand(&cobra.Command{Use: "list <user-id>", Args: cobra.ExactArgs(1), Short: "list token metadata", RunE: func(cmd *cobra.Command, args []string) error {
		client, err := managementClient(state)
		if err != nil {
			return err
		}
		tokens, err := client.ListPersonalAccessTokens(cmd.Context(), args[0])
		if err != nil {
			return fail(3, err)
		}
		data := map[string]any{"user_id": args[0], "tokens": tokens, "completeness": completeCompleteness()}
		if state.json {
			return writeJSON(stdout, map[string]any{"schema": "nb/v1/users-tokens-list-result", "ok": true, "operation": "users.tokens.list", "data": data})
		}
		for _, token := range tokens {
			if _, err := fmt.Fprintf(stdout, "%s\t%s\texpires=%s\n", token.ID, token.Name, token.ExpirationDate); err != nil {
				return err
			}
		}
		return nil
	}}, &cobra.Command{Use: "get <user-id> <token-id>", Args: cobra.ExactArgs(2), Short: "show token metadata", RunE: func(cmd *cobra.Command, args []string) error {
		client, err := managementClient(state)
		if err != nil {
			return err
		}
		token, err := client.GetPersonalAccessToken(cmd.Context(), args[0], args[1])
		if err != nil {
			return fail(3, err)
		}
		data := map[string]any{"user_id": args[0], "token": token, "completeness": completeCompleteness()}
		if state.json {
			return writeJSON(stdout, map[string]any{"schema": "nb/v1/users-tokens-get-result", "ok": true, "operation": "users.tokens.get", "data": data})
		}
		_, err = fmt.Fprintf(stdout, "%s\t%s\texpires=%s\n", token.ID, token.Name, token.ExpirationDate)
		return err
	}})
	return command
}
