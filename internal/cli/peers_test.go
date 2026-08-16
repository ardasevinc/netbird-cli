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

func TestPeersListCommandEmitsNormalizedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/peers" || r.URL.Query().Get("name") != "laptop" {
			t.Fatalf("unexpected request: %s %s?%s", r.Method, r.URL.Path, r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "p1", "name": "laptop", "ip": "10.0.0.2", "connected": true, "groups": []any{}}})
	}))
	defer server.Close()
	t.Setenv("NB_PEERS_TOKEN", "peer-token")
	temp := t.TempDir()
	configPath := filepath.Join(temp, "config.toml")
	if err := os.WriteFile(configPath, []byte("[profiles.default]\nurl = \""+server.URL+"\"\ncredential_ref = \"env:NB_PEERS_TOKEN\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	state := &commandState{json: true, configPath: configPath, profileName: "default", statePath: filepath.Join(temp, "ledger.db")}
	var stdout, stderr bytes.Buffer
	root := newRoot(state, &stdout, &stderr, version.Current())
	root.SetArgs([]string{"peers", "list", "--name", "laptop"})
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("peers list: %v stderr=%s", err, stderr.String())
	}
	var response struct {
		Schema string `json:"schema"`
		Data   struct {
			Peers []struct {
				ID string `json:"id"`
			} `json:"peers"`
		} `json:"data"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v: %s", err, stdout.String())
	}
	if response.Schema != "nb/v1/peers-list-result" || len(response.Data.Peers) != 1 || response.Data.Peers[0].ID != "p1" {
		t.Fatalf("unexpected response: %s", stdout.String())
	}
}
