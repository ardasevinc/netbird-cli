package netbird

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/transport"
)

func TestUpdateAgentNetworkSettingsUsesPUT(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/api/agent-network/settings" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var request map[string]any
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		if request["enabled"] != false {
			t.Fatalf("unexpected request: %#v", request)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"enabled": false})
	}))
	defer server.Close()
	transportClient, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	response, err := NewClient(transportClient).UpdateAgentNetworkSettings(context.Background(), json.RawMessage(`{"enabled":false}`))
	if err != nil {
		t.Fatal(err)
	}
	var updated map[string]any
	if err := json.Unmarshal(response, &updated); err != nil || updated["enabled"] != false {
		t.Fatalf("updated=%s err=%v", response, err)
	}
}

func TestCreateAgentNetworkSettingsUsesPOST(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/agent-network/settings" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var request map[string]any
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		if request["enabled"] != true {
			t.Fatalf("unexpected request: %#v", request)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"enabled": true})
	}))
	defer server.Close()
	transportClient, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	response, err := NewClient(transportClient).CreateAgentNetworkSettings(context.Background(), json.RawMessage(`{"enabled":true}`))
	if err != nil {
		t.Fatal(err)
	}
	var created map[string]any
	if err := json.Unmarshal(response, &created); err != nil || created["enabled"] != true {
		t.Fatalf("created=%s err=%v", response, err)
	}
}

func TestDeleteAgentNetworkSettingsUsesDELETE(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/api/agent-network/settings" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	transportClient, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewClient(transportClient).DeleteAgentNetworkSettings(context.Background()); err != nil {
		t.Fatal(err)
	}
}
