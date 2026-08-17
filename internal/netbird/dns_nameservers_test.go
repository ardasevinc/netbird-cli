package netbird

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/transport"
)

func TestCreateNameserverGroupUsesPOST(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/dns/nameservers" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var request map[string]any
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		if request["name"] != "office" {
			t.Fatalf("unexpected request: %#v", request)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "ns-1", "name": "office", "domains": []string{"office.internal"}, "enabled": true,
			"nameservers": []map[string]any{{"ip": "10.0.0.53", "ns_type": "udp", "port": 53}},
		})
	}))
	defer server.Close()
	transportClient, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	response, err := NewClient(transportClient).CreateNameserverGroup(context.Background(), json.RawMessage(`{"name":"office","domains":["office.internal"],"enabled":true}`))
	if err != nil {
		t.Fatal(err)
	}
	var created map[string]any
	if err := json.Unmarshal(response, &created); err != nil || created["id"] != "ns-1" {
		t.Fatalf("created=%s err=%v", response, err)
	}
}

func TestUpdateNameserverGroupUsesPUT(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/api/dns/nameservers/ns-1" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var request map[string]any
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		if request["description"] != "new" {
			t.Fatalf("unexpected request: %#v", request)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "ns-1", "name": "office", "description": "new"})
	}))
	defer server.Close()
	transportClient, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	response, err := NewClient(transportClient).UpdateNameserverGroup(context.Background(), "ns-1", json.RawMessage(`{"description":"new"}`))
	if err != nil {
		t.Fatal(err)
	}
	var updated map[string]any
	if err := json.Unmarshal(response, &updated); err != nil || updated["description"] != "new" {
		t.Fatalf("updated=%s err=%v", response, err)
	}
}

func TestDeleteNameserverGroupUsesDELETE(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/api/dns/nameservers/ns-1" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	transportClient, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewClient(transportClient).DeleteNameserverGroup(context.Background(), "ns-1"); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateDNSSettingsUsesPUT(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/api/dns/settings" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var request map[string]any
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		groups, ok := request["disabled_management_groups"].([]any)
		if !ok || len(groups) != 2 {
			t.Fatalf("unexpected request: %#v", request)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"disabled_management_groups": []string{"g1", "g2"}})
	}))
	defer server.Close()
	transportClient, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	response, err := NewClient(transportClient).UpdateDNSSettings(context.Background(), json.RawMessage(`{"disabled_management_groups":["g1","g2"]}`))
	if err != nil {
		t.Fatal(err)
	}
	var updated map[string]any
	if err := json.Unmarshal(response, &updated); err != nil || updated["disabled_management_groups"] == nil {
		t.Fatalf("updated=%s err=%v", response, err)
	}
}
