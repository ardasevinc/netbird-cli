package cli

import (
	"fmt"
	"io"

	"github.com/ardasevinc/netbird-cli/internal/config"
	"github.com/ardasevinc/netbird-cli/internal/netbird"
	"github.com/ardasevinc/netbird-cli/internal/transport"
	"github.com/spf13/cobra"
)

func capabilitiesCommand(state *commandState, stdout io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "capabilities",
		Short: "discover server identity, version, account, and core capabilities",
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
			client, err := transport.New(transport.Config{BaseURL: profile.URL, Token: token, CAFile: profile.CAFile})
			if err != nil {
				return fail(3, err)
			}
			discovery, err := netbird.NewClient(client).Discover(cmd.Context(), token != "")
			if err != nil {
				return fail(3, err)
			}
			data := map[string]any{
				"profile":         state.profileName,
				"server_identity": profile.ServerIdentity,
				"account_id":      profile.AccountID,
				"version":         discovery.Version,
				"instance":        discovery.Instance,
				"user":            discovery.User,
				"identity_status": discovery.IdentityStatus,
				"completeness":    "complete",
			}
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/capabilities-result", "ok": true, "operation": "capabilities", "data": data})
			}
			_, err = fmt.Fprintf(stdout, "profile: %s\nserver version: %s\nsetup required: %t\naccount id: %s\nidentity status: %s\n", state.profileName, discovery.Version.ManagementCurrentVersion, discovery.Instance.SetupRequired, profile.AccountID, discovery.IdentityStatus)
			if discovery.User != nil {
				_, err = fmt.Fprintf(stdout, "user: %s <%s> (%s)\n", discovery.User.Name, discovery.User.Email, discovery.User.Role)
			}
			return err
		},
	}
}
