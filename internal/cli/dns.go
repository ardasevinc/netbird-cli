package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func dnsCommand(state *commandState, stdout io.Writer) *cobra.Command {
	command := &cobra.Command{Use: "dns", Short: "inspect NetBird DNS configuration"}
	nameservers := &cobra.Command{Use: "nameservers", Short: "inspect nameserver groups"}
	nameservers.AddCommand(dnsNameserverListCommand(state, stdout), dnsNameserverGetCommand(state, stdout))
	command.AddCommand(nameservers, &cobra.Command{Use: "settings", Short: "show DNS settings", RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := managementClient(state)
		if err != nil {
			return err
		}
		settings, err := client.GetDNSSettings(cmd.Context())
		if err != nil {
			return fail(3, err)
		}
		data := map[string]any{"settings": settings, "completeness": completeCompleteness()}
		if state.json {
			return writeJSON(stdout, map[string]any{"schema": "nb/v1/dns-settings-result", "ok": true, "operation": "dns.settings", "data": data})
		}
		_, err = fmt.Fprintf(stdout, "disabled-management-groups=%d\n", len(settings.DisabledManagementGroups))
		return err
	}}, dnsZonesCommand(state, stdout))
	return command
}

func dnsNameserverListCommand(state *commandState, stdout io.Writer) *cobra.Command {
	return &cobra.Command{Use: "list", Short: "list nameserver groups", RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := managementClient(state)
		if err != nil {
			return err
		}
		groups, err := client.ListNameserverGroups(cmd.Context())
		if err != nil {
			return fail(3, err)
		}
		data := map[string]any{"nameserver_groups": groups, "completeness": completeCompleteness()}
		if state.json {
			return writeJSON(stdout, map[string]any{"schema": "nb/v1/dns-nameservers-list-result", "ok": true, "operation": "dns.nameservers.list", "data": data})
		}
		for _, group := range groups {
			if _, err := fmt.Fprintf(stdout, "%s\t%s\tnameservers=%d\tenabled=%t\n", group.ID, group.Name, len(group.Nameservers), group.Enabled); err != nil {
				return err
			}
		}
		return nil
	}}
}

func dnsNameserverGetCommand(state *commandState, stdout io.Writer) *cobra.Command {
	return &cobra.Command{Use: "get <nameserver-group-id>", Args: cobra.ExactArgs(1), Short: "show one nameserver group", RunE: func(cmd *cobra.Command, args []string) error {
		client, err := managementClient(state)
		if err != nil {
			return err
		}
		group, err := client.GetNameserverGroup(cmd.Context(), args[0])
		if err != nil {
			return fail(3, err)
		}
		data := map[string]any{"nameserver_group": group, "completeness": completeCompleteness()}
		if state.json {
			return writeJSON(stdout, map[string]any{"schema": "nb/v1/dns-nameservers-get-result", "ok": true, "operation": "dns.nameservers.get", "data": data})
		}
		_, err = fmt.Fprintf(stdout, "%s\t%s\tnameservers=%d\tenabled=%t\n", group.ID, group.Name, len(group.Nameservers), group.Enabled)
		return err
	}}
}

func dnsZonesCommand(state *commandState, stdout io.Writer) *cobra.Command {
	command := &cobra.Command{Use: "zones", Short: "inspect DNS zones"}
	command.AddCommand(&cobra.Command{Use: "list", Short: "list DNS zones", RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := managementClient(state)
		if err != nil {
			return err
		}
		zones, err := client.ListDNSZones(cmd.Context())
		if err != nil {
			return fail(3, err)
		}
		data := map[string]any{"zones": zones, "completeness": completeCompleteness()}
		if state.json {
			return writeJSON(stdout, map[string]any{"schema": "nb/v1/dns-zones-list-result", "ok": true, "operation": "dns.zones.list", "data": data})
		}
		for _, zone := range zones {
			if _, err := fmt.Fprintf(stdout, "%s\t%s\trecords=%d\tenabled=%t\n", zone.ID, zone.Name, len(zone.Records), zone.Enabled); err != nil {
				return err
			}
		}
		return nil
	}}, &cobra.Command{Use: "get <zone-id>", Args: cobra.ExactArgs(1), Short: "show one DNS zone", RunE: func(cmd *cobra.Command, args []string) error {
		client, err := managementClient(state)
		if err != nil {
			return err
		}
		zone, err := client.GetDNSZone(cmd.Context(), args[0])
		if err != nil {
			return fail(3, err)
		}
		data := map[string]any{"zone": zone, "completeness": completeCompleteness()}
		if state.json {
			return writeJSON(stdout, map[string]any{"schema": "nb/v1/dns-zones-get-result", "ok": true, "operation": "dns.zones.get", "data": data})
		}
		_, err = fmt.Fprintf(stdout, "%s\t%s\trecords=%d\tenabled=%t\n", zone.ID, zone.Name, len(zone.Records), zone.Enabled)
		return err
	}})
	return command
}
