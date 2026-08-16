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

func TestUserTokensListJSONNeverEmitsTokenValue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/users/user-1/tokens" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "token-1", "name": "agent", "created_by": "user-1", "created_at": "", "expiration_date": "", "plain_token": "secret-token"}})
	}))
	defer server.Close()
	t.Setenv("NB_USER_TOKENS_TOKEN", "token")
	temp := t.TempDir()
	configPath := filepath.Join(temp, "config.toml")
	if err := os.WriteFile(configPath, []byte("[profiles.default]\nurl = \""+server.URL+"\"\ncredential_ref = \"env:NB_USER_TOKENS_TOKEN\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	state := &commandState{json: true, configPath: configPath, profileName: "default", statePath: filepath.Join(temp, "ledger.db")}
	var stdout, stderr bytes.Buffer
	root := newRoot(state, &stdout, &stderr, version.Current())
	root.SetArgs([]string{"users", "tokens", "list", "user-1"})
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(stdout.String(), "secret-token") {
		t.Fatalf("secret leaked: %s", stdout.String())
	}
}
