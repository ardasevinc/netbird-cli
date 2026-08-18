package cli

import (
	"os"

	"github.com/ardasevinc/netbird-cli/internal/config"
	"github.com/ardasevinc/netbird-cli/internal/transport"
)

func profileTransportConfig(state *commandState, profile config.Profile, token string) (transport.Config, error) {
	timeout := state.timeout
	if !state.timeoutExplicit && os.Getenv("NB_TIMEOUT") == "" && profile.Timeout != "" {
		parsed, err := config.ParseTimeout(profile.Timeout)
		if err != nil {
			return transport.Config{}, err
		}
		timeout = parsed
	}
	return transport.Config{
		BaseURL:   profile.URL,
		Token:     token,
		CAFile:    profile.CAFile,
		Timeout:   timeout,
		LogLevel:  state.logLevel,
		LogWriter: state.logWriter,
	}, nil
}

func bootstrapTransportConfig(state *commandState, baseURL string) transport.Config {
	return transport.Config{
		BaseURL:   baseURL,
		Timeout:   state.timeout,
		LogLevel:  state.logLevel,
		LogWriter: state.logWriter,
	}
}
