package cli

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/ardasevinc/netbird-cli/internal/config"
	"github.com/ardasevinc/netbird-cli/internal/exit"
	"github.com/ardasevinc/netbird-cli/internal/netbird"
	"github.com/ardasevinc/netbird-cli/internal/transport"
	"github.com/spf13/cobra"
)

type bootstrapInput struct {
	URL         string `json:"url"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	PasswordRef string `json:"password_ref"`
	CreatePAT   bool   `json:"create_pat"`
	PATExpireIn int    `json:"pat_expire_in,omitempty"`
}

func setupCommand(state *commandState, stdout io.Writer) *cobra.Command {
	command := &cobra.Command{Use: "setup", Short: "bootstrap an uninitialized NetBird instance"}
	command.AddCommand(setupBootstrapCommand(state, stdout))
	return command
}

func setupBootstrapCommand(state *commandState, stdout io.Writer) *cobra.Command {
	var input bootstrapInput
	var fromJSON bool
	command := &cobra.Command{
		Use:   "bootstrap",
		Short: "create the initial administrator on an instance that still requires setup",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if fromJSON {
				if err := decodeBootstrapInput(cmd.InOrStdin(), &input); err != nil {
					return fail(int(exit.InvalidInput), err)
				}
			}
			if err := validateBootstrapInput(input); err != nil {
				return fail(int(exit.InvalidInput), err)
			}
			profile := config.Profile{URL: input.URL}
			if err := profile.Validate(); err != nil {
				return fail(int(exit.InvalidInput), fmt.Errorf("bootstrap server: %w", err))
			}
			clientTransport, err := transport.New(bootstrapTransportConfig(state, input.URL))
			if err != nil {
				return fail(int(exit.PreDispatch), err)
			}
			client := netbird.NewClient(clientTransport)
			instance, err := client.GetInstance(cmd.Context())
			if err != nil {
				return fail(int(exit.PreDispatch), err)
			}
			if !instance.SetupRequired {
				return fail(int(exit.PreDispatch), errors.New("server does not require setup; refusing bootstrap"))
			}
			password, err := config.ResolveCredential(input.PasswordRef)
			if err != nil {
				return fail(int(exit.PreDispatch), fmt.Errorf("resolve bootstrap password: %w", err))
			}
			if len(password) < 8 {
				return fail(int(exit.InvalidInput), errors.New("bootstrap password must be at least 8 characters"))
			}
			result, err := client.Bootstrap(cmd.Context(), netbird.SetupRequest{Email: input.Email, Password: password, Name: input.Name, CreatePAT: input.CreatePAT, PATExpireIn: input.PATExpireIn})
			if err != nil {
				var requestErr *transport.RequestError
				if errors.As(err, &requestErr) && requestErr.Dispatched {
					return fail(int(exit.Uncertain), errors.New("bootstrap request may have reached the server; do not retry without reconciling instance state"))
				}
				return fail(int(exit.PreDispatch), err)
			}
			if input.CreatePAT && result.PersonalAccessToken == "" {
				return fail(int(exit.ReceiptDelivery), errors.New("bootstrap succeeded but the requested personal access token was not returned"))
			}
			data := map[string]any{"server_identity": client.ServerIdentity(), "setup_required_before": true, "created": true}
			if result.PersonalAccessToken != "" {
				data["personal_access_token"] = result.PersonalAccessToken
			}
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/setup-bootstrap-result", "ok": true, "operation": "setup.bootstrap", "data": data})
			}
			if _, err := fmt.Fprintf(stdout, "bootstrapped %s\n", client.ServerIdentity()); err != nil {
				return err
			}
			if token, ok := data["personal_access_token"].(string); ok {
				_, err = fmt.Fprintf(stdout, "personal_access_token=%s (save it now; it will not be shown again)\n", token)
			}
			return err
		},
	}
	command.Flags().BoolVar(&fromJSON, "from-json", false, "read bootstrap parameters from stdin")
	command.Flags().StringVar(&input.URL, "url", "", "instance origin")
	command.Flags().StringVar(&input.Email, "email", "", "initial administrator email")
	command.Flags().StringVar(&input.Name, "name", "", "initial administrator name")
	command.Flags().StringVar(&input.PasswordRef, "password-ref", "", "external password reference (env:NAME or file:/path)")
	command.Flags().BoolVar(&input.CreatePAT, "create-pat", false, "request a personal access token in the one-time response")
	command.Flags().IntVar(&input.PATExpireIn, "pat-expire-in", 0, "personal access token lifetime in days (1-365)")
	return command
}

func decodeBootstrapInput(reader io.Reader, input *bootstrapInput) error {
	decoder := json.NewDecoder(bufio.NewReader(reader))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(input); err != nil {
		return fmt.Errorf("decode bootstrap input: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return errors.New("bootstrap input must contain exactly one JSON object")
		}
		return fmt.Errorf("decode bootstrap input: %w", err)
	}
	return nil
}

func validateBootstrapInput(input bootstrapInput) error {
	if strings.TrimSpace(input.URL) == "" {
		return errors.New("bootstrap url is required")
	}
	if strings.TrimSpace(input.Email) == "" {
		return errors.New("bootstrap email is required")
	}
	if strings.TrimSpace(input.Name) == "" {
		return errors.New("bootstrap name is required")
	}
	if strings.TrimSpace(input.PasswordRef) == "" {
		return errors.New("bootstrap password_ref is required; literal passwords are not accepted")
	}
	if input.CreatePAT && (input.PATExpireIn < 1 || input.PATExpireIn > 365) {
		return errors.New("pat_expire_in must be between 1 and 365 when create_pat is true")
	}
	if !input.CreatePAT && input.PATExpireIn != 0 {
		return errors.New("pat_expire_in requires create_pat")
	}
	return nil
}
