package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func usersCommand(state *commandState, stdout io.Writer) *cobra.Command {
	command := &cobra.Command{Use: "users", Short: "inspect NetBird users and invites"}
	var serviceUser string
	list := &cobra.Command{
		Use:   "list",
		Short: "list users with an optional service-user filter",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := managementClient(state)
			if err != nil {
				return err
			}
			filter, err := optionalBool(serviceUser)
			if err != nil {
				return fail(2, err)
			}
			users, err := client.ListUsers(cmd.Context(), filter)
			if err != nil {
				return fail(3, err)
			}
			data := map[string]any{"users": users, "completeness": completeCompleteness()}
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/users-list-result", "ok": true, "operation": "users.list", "data": data})
			}
			for _, user := range users {
				if _, err := fmt.Fprintf(stdout, "%s\t%s\t%s\t%s\n", user.ID, user.Email, user.Role, user.Status); err != nil {
					return err
				}
			}
			return nil
		},
	}
	list.Flags().StringVar(&serviceUser, "service-user", "", "filter by service user: true or false")
	command.AddCommand(list)
	command.AddCommand(userTokensCommand(state, stdout))
	command.AddCommand(&cobra.Command{
		Use:   "invites",
		Short: "list pending user invites without secret tokens",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := managementClient(state)
			if err != nil {
				return err
			}
			invites, err := client.ListInvites(cmd.Context())
			if err != nil {
				return fail(3, err)
			}
			data := map[string]any{"invites": invites, "completeness": completeCompleteness()}
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/users-invites-list-result", "ok": true, "operation": "users.invites", "data": data})
			}
			for _, invite := range invites {
				if _, err := fmt.Fprintf(stdout, "%s\t%s\t%s\t%s\texpired=%t\n", invite.ID, invite.Email, invite.Name, invite.Role, invite.Expired); err != nil {
					return err
				}
			}
			return nil
		},
	})
	return command
}

func optionalBool(value string) (*bool, error) {
	if value == "" {
		return nil, nil
	}
	if value != "true" && value != "false" {
		return nil, fmt.Errorf("service-user must be true or false")
	}
	parsed := value == "true"
	return &parsed, nil
}
