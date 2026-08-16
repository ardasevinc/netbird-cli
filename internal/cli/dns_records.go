package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func dnsZoneRecordsCommand(state *commandState, stdout io.Writer) *cobra.Command {
	command := &cobra.Command{Use: "records", Short: "inspect DNS zone records"}
	command.AddCommand(&cobra.Command{Use: "list <zone-id>", Args: cobra.ExactArgs(1), Short: "list records", RunE: func(cmd *cobra.Command, args []string) error {
		client, err := managementClient(state)
		if err != nil {
			return err
		}
		records, err := client.ListDNSRecords(cmd.Context(), args[0])
		if err != nil {
			return fail(3, err)
		}
		data := map[string]any{"zone_id": args[0], "records": records, "completeness": completeCompleteness()}
		if state.json {
			return writeJSON(stdout, map[string]any{"schema": "nb/v1/dns-records-list-result", "ok": true, "operation": "dns.records.list", "data": data})
		}
		for _, record := range records {
			if _, err := fmt.Fprintf(stdout, "%s\t%s\t%s\t%s\tttl=%d\n", record.ID, record.Name, record.Type, record.Content, record.TTL); err != nil {
				return err
			}
		}
		return nil
	}}, &cobra.Command{Use: "get <zone-id> <record-id>", Args: cobra.ExactArgs(2), Short: "show one record", RunE: func(cmd *cobra.Command, args []string) error {
		client, err := managementClient(state)
		if err != nil {
			return err
		}
		record, err := client.GetDNSRecord(cmd.Context(), args[0], args[1])
		if err != nil {
			return fail(3, err)
		}
		data := map[string]any{"zone_id": args[0], "record": record, "completeness": completeCompleteness()}
		if state.json {
			return writeJSON(stdout, map[string]any{"schema": "nb/v1/dns-records-get-result", "ok": true, "operation": "dns.records.get", "data": data})
		}
		_, err = fmt.Fprintf(stdout, "%s\t%s\t%s\t%s\tttl=%d\n", record.ID, record.Name, record.Type, record.Content, record.TTL)
		return err
	}})
	return command
}
