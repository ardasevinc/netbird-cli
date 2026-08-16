package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func eventsCommand(state *commandState, stdout io.Writer) *cobra.Command {
	return &cobra.Command{Use: "events", Short: "inspect NetBird events", RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := managementClient(state)
		if err != nil {
			return err
		}
		events, err := client.ListAuditEvents(cmd.Context())
		if err != nil {
			return fail(3, err)
		}
		data := map[string]any{"events": events, "completeness": completeCompleteness()}
		if state.json {
			return writeJSON(stdout, map[string]any{"schema": "nb/v1/events-audit-result", "ok": true, "operation": "events.audit", "data": data})
		}
		for _, event := range events {
			if _, err := fmt.Fprintf(stdout, "%s\t%s\t%s\t%s\n", event.ID, event.Timestamp, event.ActivityCode, event.Activity); err != nil {
				return err
			}
		}
		return nil
	}}
}
