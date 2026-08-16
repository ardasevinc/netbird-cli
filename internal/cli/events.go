package cli

import (
	"fmt"
	"io"
	"time"

	"github.com/ardasevinc/netbird-cli/internal/netbird"
	"github.com/spf13/cobra"
)

func eventsCommand(state *commandState, stdout io.Writer) *cobra.Command {
	command := &cobra.Command{Use: "events", Short: "inspect NetBird events", RunE: func(cmd *cobra.Command, _ []string) error {
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
	command.AddCommand(networkTrafficEventsCommand(state, stdout), proxyEventsCommand(state, stdout))
	return command
}

func networkTrafficEventsCommand(state *commandState, stdout io.Writer) *cobra.Command {
	var page, pageSize int
	command := &cobra.Command{Use: "network-traffic", Short: "list network traffic events", RunE: func(cmd *cobra.Command, _ []string) error {
		if page < 0 || pageSize < 0 || pageSize > 50000 {
			return fmt.Errorf("page must be positive and page size must be between 1 and 50000")
		}
		client, err := managementClient(state)
		if err != nil {
			return err
		}
		result, err := client.ListNetworkTrafficEvents(cmd.Context(), netbird.EventPageOptions{Page: page, PageSize: pageSize})
		if err != nil {
			return fail(3, err)
		}
		data := map[string]any{"events": result.Data, "page": result.Page, "page_size": result.PageSize, "total_pages": result.TotalPages, "total_records": result.TotalRecords, "completeness": paginatedCompleteness(result.Page, result.TotalPages)}
		if state.json {
			return writeJSON(stdout, map[string]any{"schema": "nb/v1/events-network-traffic-result", "ok": true, "operation": "events.network_traffic", "data": data})
		}
		for _, event := range result.Data {
			if _, err := fmt.Fprintf(stdout, "%s\t%s\t%s\t%s\n", event.FlowID, event.WindowStart.Format(time.RFC3339), event.Source.Name, event.Destination.Name); err != nil {
				return err
			}
		}
		return nil
	}}
	command.Flags().IntVar(&page, "page", 0, "1-indexed page number")
	command.Flags().IntVar(&pageSize, "page-size", 0, "items per page (maximum 50000)")
	return command
}

func proxyEventsCommand(state *commandState, stdout io.Writer) *cobra.Command {
	var page, pageSize int
	command := &cobra.Command{Use: "proxy", Short: "list reverse proxy access logs", RunE: func(cmd *cobra.Command, _ []string) error {
		if page < 0 || pageSize < 0 || pageSize > 100 {
			return fmt.Errorf("page must be positive and page size must be between 1 and 100")
		}
		client, err := managementClient(state)
		if err != nil {
			return err
		}
		result, err := client.ListProxyAccessLogs(cmd.Context(), netbird.EventPageOptions{Page: page, PageSize: pageSize})
		if err != nil {
			return fail(3, err)
		}
		data := map[string]any{"logs": result.Data, "page": result.Page, "page_size": result.PageSize, "total_pages": result.TotalPages, "total_records": result.TotalRecords, "completeness": paginatedCompleteness(result.Page, result.TotalPages)}
		if state.json {
			return writeJSON(stdout, map[string]any{"schema": "nb/v1/events-proxy-result", "ok": true, "operation": "events.proxy", "data": data})
		}
		for _, log := range result.Data {
			if _, err := fmt.Fprintf(stdout, "%s\t%s\t%s\t%d\n", log.Timestamp.Format(time.RFC3339), log.Method, log.Path, log.StatusCode); err != nil {
				return err
			}
		}
		return nil
	}}
	command.Flags().IntVar(&page, "page", 0, "1-indexed page number")
	command.Flags().IntVar(&pageSize, "page-size", 0, "items per page (maximum 100)")
	return command
}

func paginatedCompleteness(page, totalPages int) map[string]any {
	if totalPages <= 0 || page >= totalPages {
		return completeCompleteness()
	}
	return map[string]any{"state": "partial", "reason": fmt.Sprintf("page %d of %d returned", page, totalPages)}
}
