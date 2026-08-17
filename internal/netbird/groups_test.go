package netbird

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/transport"
)

func TestDeleteGroupUsesDELETE(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/api/groups/group-1" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	transportClient, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewClient(transportClient).DeleteGroup(context.Background(), "group-1"); err != nil {
		t.Fatal(err)
	}
}

func TestCreateGroupUsesPOST(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/groups" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
		var request map[string]any
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		if request["name"] != "engineering" {
			t.Fatalf("unexpected request body: %+v", request)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "g1", "name": "engineering", "peers_count": 0, "resources_count": 0})
	}))
	defer server.Close()
	transportClient, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewClient(transportClient).CreateGroup(context.Background(), json.RawMessage(`{"name":"engineering"}`)); err != nil {
		t.Fatal(err)
	}
}
