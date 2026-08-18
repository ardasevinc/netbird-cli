package cli

import (
	"fmt"
	"io"
	"sort"

	"github.com/ardasevinc/netbird-cli/internal/config"
	"github.com/spf13/cobra"
)

func profileCommand(state *commandState, stdout io.Writer) *cobra.Command {
	command := &cobra.Command{Use: "profile", Short: "inspect named connection profiles"}
	command.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "list configured profile names",
		RunE: func(cmd *cobra.Command, _ []string) error {
			file, err := config.Load(state.configPath)
			if err != nil {
				return fail(3, err)
			}
			names := make([]string, 0, len(file.Profiles))
			for name := range file.Profiles {
				names = append(names, name)
			}
			sort.Strings(names)
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/profile-list-result", "ok": true, "operation": "profile.list", "data": map[string]any{"config": state.configPath, "profiles": names}})
			}
			for _, name := range names {
				if _, err := fmt.Fprintln(stdout, name); err != nil {
					return err
				}
			}
			return nil
		},
	})
	command.AddCommand(&cobra.Command{
		Use:   "show [name]",
		Short: "show a profile without resolving its credential",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := state.profileName
			if len(args) == 1 {
				name = args[0]
			}
			file, err := config.Load(state.configPath)
			if err != nil {
				return fail(3, err)
			}
			profile, err := file.Profile(name)
			if err != nil {
				return fail(2, err)
			}
			data := map[string]any{"name": name, "url": profile.URL, "account_id": profile.AccountID, "credential_ref": profile.CredentialRef, "ca_file": profile.CAFile, "server_identity": profile.ServerIdentity, "timeout": profile.Timeout, "read_only": profile.ReadOnly}
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/profile-result", "ok": true, "operation": "profile.show", "data": data})
			}
			_, err = fmt.Fprintf(stdout, "profile %s\nurl: %s\naccount_id: %s\ncredential_ref: %s\nca_file: %s\nserver_identity: %s\ntimeout: %s\nread_only: %t\n", name, profile.URL, profile.AccountID, profile.CredentialRef, profile.CAFile, profile.ServerIdentity, profile.Timeout, profile.ReadOnly)
			return err
		},
	})
	command.AddCommand(&cobra.Command{
		Use:   "validate [name]",
		Short: "validate a profile without network access",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := state.profileName
			if len(args) == 1 {
				name = args[0]
			}
			file, err := config.Load(state.configPath)
			if err != nil {
				return fail(3, err)
			}
			if _, err := file.Profile(name); err != nil {
				return fail(2, err)
			}
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/profile-validate-result", "ok": true, "operation": "profile.validate", "data": map[string]any{"name": name, "valid": true}})
			}
			_, err = fmt.Fprintf(stdout, "profile %s is valid\n", name)
			return err
		},
	})
	return command
}
