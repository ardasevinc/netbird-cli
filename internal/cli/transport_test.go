package cli

import (
	"testing"
	"time"

	"github.com/ardasevinc/netbird-cli/internal/config"
)

func TestProfileTransportConfigUsesProfileTimeout(t *testing.T) {
	t.Setenv("NB_TIMEOUT", "")
	state := &commandState{timeout: config.DefaultTimeout, logLevel: config.DefaultLogLevel}
	profile := config.Profile{URL: "https://netbird.example.test", Timeout: "7s"}
	settings, err := profileTransportConfig(state, profile, "")
	if err != nil {
		t.Fatal(err)
	}
	if settings.Timeout != 7*time.Second {
		t.Fatalf("timeout=%s, want 7s", settings.Timeout)
	}
}

func TestExplicitTimeoutBeatsProfileTimeout(t *testing.T) {
	t.Setenv("NB_TIMEOUT", "")
	state := &commandState{timeout: 5 * time.Second, timeoutExplicit: true, logLevel: config.DefaultLogLevel}
	profile := config.Profile{URL: "https://netbird.example.test", Timeout: "7s"}
	settings, err := profileTransportConfig(state, profile, "")
	if err != nil {
		t.Fatal(err)
	}
	if settings.Timeout != 5*time.Second {
		t.Fatalf("timeout=%s, want 5s", settings.Timeout)
	}
}
