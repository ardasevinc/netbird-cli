package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/version"
)

func TestSetupKeysListJSONNeverEmitsSecret(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/setup-keys" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "key-1", "name": "bootstrap", "type": "reusable", "state": "valid", "valid": true, "revoked": false, "ephemeral": false, "allow_extra_dns_labels": false, "auto_groups": []string{}, "usage_limit": 5, "used_times": 1, "expires": "", "last_used": "", "updated_at": "", "key": "super-secret"}})
	}))
	defer server.Close()
	t.Setenv("NB_SETUP_KEYS_TOKEN", "token")
	temp := t.TempDir()
	configPath := filepath.Join(temp, "config.toml")
	if err := os.WriteFile(configPath, []byte("[profiles.default]\nurl = \""+server.URL+"\"\ncredential_ref = \"env:NB_SETUP_KEYS_TOKEN\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	state := &commandState{json: true, configPath: configPath, profileName: "default", statePath: filepath.Join(temp, "ledger.db")}
	var stdout, stderr bytes.Buffer
	root := newRoot(state, &stdout, &stderr, version.Current())
	root.SetArgs([]string{"setup-keys", "list"})
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(stdout.String(), "super-secret") {
		t.Fatalf("secret leaked: %s", stdout.String())
	}
}
