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

func TestEventsAuditJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/events/audit" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "event-1", "activity": "peer connected", "activity_code": "peer_connected", "initiator_email": "a@example.com", "initiator_id": "user-1", "initiator_name": "Arda", "meta": map[string]string{}, "target_id": "peer-1", "timestamp": "2026-08-17T10:00:00Z"}})
	}))
	defer server.Close()
	t.Setenv("NB_EVENTS_TOKEN", "token")
	temp := t.TempDir()
	configPath := filepath.Join(temp, "config.toml")
	if err := os.WriteFile(configPath, []byte("[profiles.default]\nurl = \""+server.URL+"\"\ncredential_ref = \"env:NB_EVENTS_TOKEN\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	state := &commandState{json: true, configPath: configPath, profileName: "default", statePath: filepath.Join(temp, "ledger.db")}
	var stdout, stderr bytes.Buffer
	root := newRoot(state, &stdout, &stderr, version.Current())
	root.SetArgs([]string{"events"})
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("events: %v stderr=%s", err, stderr.String())
	}
	var response struct {
		Schema string `json:"schema"`
		Data   struct {
			Events []struct {
				ID string `json:"id"`
			} `json:"events"`
		} `json:"data"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.Schema != "nb/v1/events-audit-result" || len(response.Data.Events) != 1 || response.Data.Events[0].ID != "event-1" {
		t.Fatalf("unexpected response: %s", stdout.String())
	}
}

func TestEventsNetworkTrafficJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/events/network-traffic" || r.URL.Query().Get("page") != "2" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.String())
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"data": []map[string]any{{"flow_id": "flow-1", "window_start": "2026-08-17T10:00:00Z", "window_end": "2026-08-17T10:01:00Z", "source": map[string]any{"name": "a"}, "destination": map[string]any{"name": "b"}}}, "page": 2, "page_size": 1, "total_pages": 3, "total_records": 3})
	}))
	defer server.Close()
	t.Setenv("NB_EVENTS_NETWORK_TOKEN", "token")
	temp := t.TempDir()
	configPath := filepath.Join(temp, "config.toml")
	if err := os.WriteFile(configPath, []byte("[profiles.default]\nurl = \""+server.URL+"\"\ncredential_ref = \"env:NB_EVENTS_NETWORK_TOKEN\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	state := &commandState{json: true, configPath: configPath, profileName: "default", statePath: filepath.Join(temp, "ledger.db")}
	var stdout, stderr bytes.Buffer
	root := newRoot(state, &stdout, &stderr, version.Current())
	root.SetArgs([]string{"events", "network-traffic", "--page", "2"})
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("events network-traffic: %v stderr=%s", err, stderr.String())
	}
	var response struct {
		Schema string `json:"schema"`
		Data   struct {
			Events []struct {
				FlowID string `json:"flow_id"`
			} `json:"events"`
			Completeness struct {
				State string `json:"state"`
			} `json:"completeness"`
		} `json:"data"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.Schema != "nb/v1/events-network-traffic-result" || len(response.Data.Events) != 1 || response.Data.Events[0].FlowID != "flow-1" || response.Data.Completeness.State != "partial" {
		t.Fatalf("unexpected response: %s", stdout.String())
	}
}
