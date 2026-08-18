package cli

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/ardasevinc/netbird-cli/internal/config"
	"github.com/ardasevinc/netbird-cli/internal/netbird"
	"github.com/ardasevinc/netbird-cli/internal/transport"
	"github.com/spf13/cobra"
)

func groupsCommand(state *commandState, stdout io.Writer) *cobra.Command {
	command := &cobra.Command{Use: "groups", Short: "inspect NetBird groups"}
	command.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "list all groups with complete retrieval semantics",
		RunE: func(cmd *cobra.Command, _ []string) error {
			file, err := config.Load(state.configPath)
			if err != nil {
				return fail(3, err)
			}
			profile, err := file.Profile(state.profileName)
			if err != nil {
				return fail(2, err)
			}
			token, err := config.ResolveCredential(profile.CredentialRef)
			if err != nil {
				return fail(3, err)
			}
			transportConfig, err := profileTransportConfig(state, profile, token)
			if err != nil {
				return fail(3, err)
			}
			transportClient, err := transport.New(transportConfig)
			if err != nil {
				return fail(3, err)
			}
			groups, err := netbird.NewClient(transportClient).ListGroups(cmd.Context())
			if err != nil {
				return fail(3, err)
			}
			data := map[string]any{"groups": groups, "completeness": map[string]any{"state": "complete", "reason": nil}}
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/groups-list-result", "ok": true, "operation": "groups.list", "data": data})
			}
			for _, group := range groups {
				if _, err := fmt.Fprintf(stdout, "%s\t%s\tpeers=%d\tresources=%d\n", group.ID, group.Name, group.PeersCount, group.ResourcesCount); err != nil {
					return err
				}
			}
			return nil
		},
	})
	command.AddCommand(&cobra.Command{
		Use:   "get <group-id>",
		Short: "show one group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			file, err := config.Load(state.configPath)
			if err != nil {
				return fail(3, err)
			}
			profile, err := file.Profile(state.profileName)
			if err != nil {
				return fail(2, err)
			}
			token, err := config.ResolveCredential(profile.CredentialRef)
			if err != nil {
				return fail(3, err)
			}
			transportConfig, err := profileTransportConfig(state, profile, token)
			if err != nil {
				return fail(3, err)
			}
			transportClient, err := transport.New(transportConfig)
			if err != nil {
				return fail(3, err)
			}
			group, err := netbird.NewClient(transportClient).GetGroup(cmd.Context(), args[0])
			if err != nil {
				return fail(3, err)
			}
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/groups-get-result", "ok": true, "operation": "groups.get", "data": map[string]any{"group": json.RawMessage(group), "completeness": map[string]any{"state": "complete", "reason": nil}}})
			}
			_, err = fmt.Fprintf(stdout, "%s\n", group)
			return err
		},
	})
	return command
}
