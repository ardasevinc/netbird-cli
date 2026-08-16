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

func TestAPIGetUsesManifestPathAndQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/events/proxy" || r.URL.Query().Get("page") != "2" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.String())
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"rows": []map[string]string{{"id": "log-1"}}})
	}))
	defer server.Close()
	t.Setenv("NB_API_TOKEN", "token")
	temp := t.TempDir()
	configPath := filepath.Join(temp, "config.toml")
	if err := os.WriteFile(configPath, []byte("[profiles.default]\nurl = \""+server.URL+"\"\ncredential_ref = \"env:NB_API_TOKEN\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	state := &commandState{json: true, configPath: configPath, profileName: "default", statePath: filepath.Join(temp, "ledger.db")}
	var stdout, stderr bytes.Buffer
	root := newRoot(state, &stdout, &stderr, version.Current())
	root.SetArgs([]string{"api", "get", "events.proxy", "--query", "page=2"})
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("api get: %v stderr=%s", err, stderr.String())
	}
	var response struct {
		Schema string `json:"schema"`
		Data   struct {
			OperationID string `json:"operation_id"`
			Payload     struct {
				Rows []struct {
					ID string `json:"id"`
				} `json:"rows"`
			} `json:"payload"`
		} `json:"data"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.Schema != "nb/v1/raw-get-result" || response.Data.OperationID != "events.proxy" || len(response.Data.Payload.Rows) != 1 || response.Data.Payload.Rows[0].ID != "log-1" {
		t.Fatalf("unexpected response: %s", stdout.String())
	}
}

func TestResolveOperationPathEscapesValues(t *testing.T) {
	path, err := resolveOperationPath("/api/users/{userId}/tokens/{tokenId}", []string{"user/1", "token 1"})
	if err != nil || path != "/api/users/user%2F1/tokens/token%201" {
		t.Fatalf("path=%q err=%v", path, err)
	}
}
