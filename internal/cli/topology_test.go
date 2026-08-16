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

func TestRoutesListJSONEmitsCompleteTopologyResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/routes" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "route-1", "description": "office", "enabled": true, "groups": []string{}, "keep_route": false, "masquerade": true, "metric": 10, "network": "10.0.0.0/24", "network_id": "", "network_type": "static"}})
	}))
	defer server.Close()
	t.Setenv("NB_ROUTES_TOKEN", "routes-token")
	temp := t.TempDir()
	configPath := filepath.Join(temp, "config.toml")
	if err := os.WriteFile(configPath, []byte("[profiles.default]\nurl = \""+server.URL+"\"\ncredential_ref = \"env:NB_ROUTES_TOKEN\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	state := &commandState{json: true, configPath: configPath, profileName: "default", statePath: filepath.Join(temp, "ledger.db")}
	var stdout, stderr bytes.Buffer
	root := newRoot(state, &stdout, &stderr, version.Current())
	root.SetArgs([]string{"routes", "list"})
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("routes list: %v stderr=%s", err, stderr.String())
	}
	var response struct {
		Schema string `json:"schema"`
		Data   struct {
			Routes []struct {
				ID string `json:"id"`
			} `json:"routes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v: %s", err, stdout.String())
	}
	if response.Schema != "nb/v1/routes-list-result" || len(response.Data.Routes) != 1 || response.Data.Routes[0].ID != "route-1" {
		t.Fatalf("unexpected response: %s", stdout.String())
	}
}
