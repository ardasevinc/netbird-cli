package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func policiesCommand(state *commandState, stdout io.Writer) *cobra.Command {
	command := &cobra.Command{Use: "policies", Short: "inspect NetBird access policies"}
	command.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "list policies with bounded completeness semantics",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := peerClient(state)
			if err != nil {
				return err
			}
			policies, err := client.ListPolicies(cmd.Context())
			if err != nil {
				return fail(3, err)
			}
			data := map[string]any{"policies": policies, "completeness": map[string]any{"state": "complete", "reason": nil}}
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/policies-list-result", "ok": true, "operation": "policies.list", "data": data})
			}
			for _, policy := range policies {
				if _, err := fmt.Fprintf(stdout, "%s\t%s\tenabled=%t\trules=%d\n", pointerValue(policy.ID), policy.Name, policy.Enabled, len(policy.Rules)); err != nil {
					return err
				}
			}
			return nil
		},
	})
	command.AddCommand(&cobra.Command{
		Use:   "get <policy-id>",
		Short: "show one policy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := peerClient(state)
			if err != nil {
				return err
			}
			policy, err := client.GetPolicy(cmd.Context(), args[0])
			if err != nil {
				return fail(3, err)
			}
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/policies-get-result", "ok": true, "operation": "policies.get", "data": map[string]any{"policy": policy, "completeness": map[string]any{"state": "complete", "reason": nil}}})
			}
			_, err = fmt.Fprintf(stdout, "%s\t%s\tenabled=%t\trules=%d\n", pointerValue(policy.ID), policy.Name, policy.Enabled, len(policy.Rules))
			return err
		},
	})
	return command
}

func pointerValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
