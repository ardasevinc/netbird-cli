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

func TestAccessiblePeersJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/peers/p1/accessible-peers" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "p2", "name": "db", "ip": "10.0.0.3", "groups": []any{}, "connected": true}})
	}))
	defer server.Close()
	t.Setenv("NB_ACCESSIBLE_TOKEN", "token")
	temp := t.TempDir()
	configPath := filepath.Join(temp, "config.toml")
	if err := os.WriteFile(configPath, []byte("[profiles.default]\nurl = \""+server.URL+"\"\ncredential_ref = \"env:NB_ACCESSIBLE_TOKEN\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	state := &commandState{json: true, configPath: configPath, profileName: "default", statePath: filepath.Join(temp, "ledger.db")}
	var stdout, stderr bytes.Buffer
	root := newRoot(state, &stdout, &stderr, version.Current())
	root.SetArgs([]string{"peers", "accessible", "p1"})
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	var response struct {
		Schema string `json:"schema"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.Schema != "nb/v1/peers-accessible-result" {
		t.Fatalf("unexpected response: %s", stdout.String())
	}
	assertJSONKeys(t, stdout.Bytes(), "data", "ok", "operation", "schema")
	var envelope struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &envelope); err != nil {
		t.Fatal(err)
	}
	assertJSONKeys(t, envelope.Data, "completeness", "peer_id", "peers")
}

func TestNetworkMapPeersJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/peers/p1/accessible-peers" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "p2", "name": "db", "ip": "10.0.0.3", "groups": []any{}, "connected": true}})
	}))
	defer server.Close()
	t.Setenv("NB_NETWORK_MAP_TOKEN", "token")
	temp := t.TempDir()
	configPath := filepath.Join(temp, "config.toml")
	if err := os.WriteFile(configPath, []byte("[profiles.default]\nurl = \""+server.URL+"\"\ncredential_ref = \"env:NB_NETWORK_MAP_TOKEN\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	state := &commandState{json: true, configPath: configPath, profileName: "default", statePath: filepath.Join(temp, "ledger.db")}
	var stdout, stderr bytes.Buffer
	root := newRoot(state, &stdout, &stderr, version.Current())
	root.SetArgs([]string{"peers", "network-map", "p1"})
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	var response struct {
		Schema    string `json:"schema"`
		Operation string `json:"operation"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.Schema != "nb/v1/peers-network-map-result" || response.Operation != "peers.network_map" {
		t.Fatalf("unexpected response: %s", stdout.String())
	}
}

func assertJSONKeys(t *testing.T, document []byte, expected ...string) {
	t.Helper()
	var object map[string]json.RawMessage
	if err := json.Unmarshal(document, &object); err != nil {
		t.Fatal(err)
	}
	if len(object) != len(expected) {
		t.Fatalf("unexpected JSON keys: got=%v expected=%v", object, expected)
	}
	for _, key := range expected {
		if _, ok := object[key]; !ok {
			t.Fatalf("missing JSON key %q in %v", key, object)
		}
	}
}
