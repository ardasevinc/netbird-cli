package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/version"
)

func TestIdentityProvidersListJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/identity-providers" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "idp-1", "name": "zitadel", "type": "oidc", "issuer": "https://idp.example", "client_id": "client-1"}})
	}))
	defer server.Close()
	t.Setenv("NB_IDP_TOKEN", "token")
	temp := t.TempDir()
	configPath := filepath.Join(temp, "config.toml")
	if err := os.WriteFile(configPath, []byte("[profiles.default]\nurl = \""+server.URL+"\"\ncredential_ref = \"env:NB_IDP_TOKEN\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	state := &commandState{json: true, configPath: configPath, profileName: "default", statePath: filepath.Join(temp, "ledger.db")}
	var stdout, stderr bytes.Buffer
	root := newRoot(state, &stdout, &stderr, version.Current())
	root.SetArgs([]string{"identity-providers", "list"})
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("identity providers list: %v stderr=%s", err, stderr.String())
	}
	var response struct {
		Schema string `json:"schema"`
		Data   struct {
			Providers []struct {
				Name string `json:"name"`
			} `json:"identity_providers"`
		} `json:"data"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Schema != "nb/v1/identity-providers-list-result" || len(response.Data.Providers) != 1 || response.Data.Providers[0].Name != "zitadel" {
		t.Fatalf("unexpected response: %s", stdout.String())
	}
}
