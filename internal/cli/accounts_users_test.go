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

func TestUsersListJSONFiltersServiceUsers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/users" || r.URL.Query().Get("service_user") != "true" {
			t.Fatalf("unexpected request: %s %s?%s", r.Method, r.URL.Path, r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "u1", "email": "svc@example.com", "name": "svc", "role": "admin", "status": "active", "auto_groups": []string{}, "is_blocked": false, "pending_approval": false, "is_service_user": true}})
	}))
	defer server.Close()
	t.Setenv("NB_USERS_TOKEN", "users-token")
	temp := t.TempDir()
	configPath := filepath.Join(temp, "config.toml")
	if err := os.WriteFile(configPath, []byte("[profiles.default]\nurl = \""+server.URL+"\"\ncredential_ref = \"env:NB_USERS_TOKEN\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	state := &commandState{json: true, configPath: configPath, profileName: "default", statePath: filepath.Join(temp, "ledger.db")}
	var stdout, stderr bytes.Buffer
	root := newRoot(state, &stdout, &stderr, version.Current())
	root.SetArgs([]string{"users", "list", "--service-user", "true"})
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("users list: %v stderr=%s", err, stderr.String())
	}
	var response struct {
		Schema string `json:"schema"`
		Data   struct {
			Users []struct {
				ID string `json:"id"`
			} `json:"users"`
		} `json:"data"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v: %s", err, stdout.String())
	}
	if response.Schema != "nb/v1/users-list-result" || len(response.Data.Users) != 1 || response.Data.Users[0].ID != "u1" {
		t.Fatalf("unexpected response: %s", stdout.String())
	}
}
