package cli

import (
	"fmt"
	"io"

	"github.com/ardasevinc/netbird-cli/internal/config"
	"github.com/ardasevinc/netbird-cli/internal/netbird"
	"github.com/ardasevinc/netbird-cli/internal/transport"
	"github.com/spf13/cobra"
)

func peersCommand(state *commandState, stdout io.Writer) *cobra.Command {
	command := &cobra.Command{Use: "peers", Short: "inspect NetBird peers"}
	command.AddCommand(accessiblePeersCommand(state, stdout))
	var name, ip string
	command.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "list peers with bounded completeness semantics",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := managementClient(state)
			if err != nil {
				return err
			}
			peers, err := client.ListPeers(cmd.Context(), name, ip)
			if err != nil {
				return fail(3, err)
			}
			data := map[string]any{"peers": peers, "completeness": map[string]any{"state": "complete", "reason": nil}}
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/peers-list-result", "ok": true, "operation": "peers.list", "data": data})
			}
			for _, peer := range peers {
				if _, err := fmt.Fprintf(stdout, "%s\t%s\t%s\tconnected=%t\n", peer.ID, peer.Name, peer.IP, peer.Connected); err != nil {
					return err
				}
			}
			return nil
		},
	})
	command.AddCommand(&cobra.Command{
		Use:   "get <peer-id>",
		Short: "show one peer",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := managementClient(state)
			if err != nil {
				return err
			}
			peer, err := client.GetPeer(cmd.Context(), args[0])
			if err != nil {
				return fail(3, err)
			}
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/peers-get-result", "ok": true, "operation": "peers.get", "data": map[string]any{"peer": peer, "completeness": map[string]any{"state": "complete", "reason": nil}}})
			}
			_, err = fmt.Fprintf(stdout, "%s\t%s\t%s\tconnected=%t\n", peer.ID, peer.Name, peer.IP, peer.Connected)
			return err
		},
	})
	command.PersistentFlags().StringVar(&name, "name", "", "filter peers by name")
	command.PersistentFlags().StringVar(&ip, "ip", "", "filter peers by IP address")
	return command
}

func managementClient(state *commandState) (*netbird.Client, error) {
	file, err := config.Load(state.configPath)
	if err != nil {
		return nil, fail(3, err)
	}
	profile, err := file.Profile(state.profileName)
	if err != nil {
		return nil, fail(2, err)
	}
	token, err := config.ResolveCredential(profile.CredentialRef)
	if err != nil {
		return nil, fail(3, err)
	}
	client, err := transport.New(transport.Config{BaseURL: profile.URL, Token: token, CAFile: profile.CAFile})
	if err != nil {
		return nil, fail(3, err)
	}
	return netbird.NewClient(client), nil
}
