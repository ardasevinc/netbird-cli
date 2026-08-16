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

	"github.com/ardasevinc/netbird-cli/internal/ledger"
	"github.com/ardasevinc/netbird-cli/internal/version"
)

func TestApplyCommandUsesExactStageAndReturnsMachineResult(t *testing.T) {
	group := map[string]any{"id": "g1", "name": "old", "peers_count": 0, "resources_count": 0}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer cli-token" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/accounts":
			_ = json.NewEncoder(w).Encode([]map[string]string{{"id": "account-1"}})
		case r.Method == http.MethodGet && r.URL.Path == "/api/groups/g1":
			_ = json.NewEncoder(w).Encode(group)
		case r.Method == http.MethodPut && r.URL.Path == "/api/groups/g1":
			var request map[string]any
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				t.Errorf("decode request: %v", err)
			}
			group["name"] = request["name"]
			_ = json.NewEncoder(w).Encode(group)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	t.Setenv("NB_APPLY_TOKEN", "cli-token")
	temp := t.TempDir()
	configPath := filepath.Join(temp, "config.toml")
	statePath := filepath.Join(temp, "ledger.db")
	config := "[profiles.default]\nurl = \"" + server.URL + "\"\naccount_id = \"account-1\"\nserver_identity = \"" + server.URL + "\"\ncredential_ref = \"env:NB_APPLY_TOKEN\"\n"
	if err := os.WriteFile(configPath, []byte(config), 0o600); err != nil {
		t.Fatal(err)
	}
	store, err := ledger.Open(statePath)
	if err != nil {
		t.Fatal(err)
	}
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: server.URL,
		AccountID:      "account-1",
		Operation:      "groups.update",
		Request:        json.RawMessage(`{"id":"g1","name":"new"}`),
		Before:         json.RawMessage(`{"id":"g1","name":"old","peers_count":0,"resources_count":0}`),
		IntendedAfter:  json.RawMessage(`{"id":"g1","name":"new","peers_count":0,"resources_count":0}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	_ = store.Close()

	state := &commandState{json: true, configPath: configPath, profileName: "default", statePath: statePath}
	var stdout, stderr bytes.Buffer
	root := newRoot(state, &stdout, &stderr, version.Current())
	root.SetArgs([]string{"apply", stage.ID + "@1"})
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("apply: %v stderr=%s", err, stderr.String())
	}
	var response struct {
		Schema string `json:"schema"`
		Data   struct {
			State string `json:"state"`
		} `json:"data"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		t.Fatalf("response: %v: %s", err, stdout.String())
	}
	if response.Schema != "nb/v1/apply-result" || response.Data.State != "confirmed_success" {
		t.Fatalf("unexpected response: %s", stdout.String())
	}
}
