package netbird

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/transport"
)

func TestUpdatePeerUsesPUTAndReturnsRawDocument(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/api/peers/p1" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
		var request map[string]any
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		if request["name"] != "updated" {
			t.Fatalf("unexpected request body: %+v", request)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "p1", "name": "updated", "approval_required": false})
	}))
	defer server.Close()
	client, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	result, err := NewClient(client).UpdatePeer(context.Background(), "p1", json.RawMessage(`{"name":"updated"}`))
	if err != nil {
		t.Fatal(err)
	}
	var response map[string]any
	if err := json.Unmarshal(result, &response); err != nil || response["name"] != "updated" {
		t.Fatalf("unexpected result: %s err=%v", result, err)
	}
}

func TestDeletePeerUsesDELETE(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/api/peers/p1" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	client, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewClient(client).DeletePeer(context.Background(), "p1"); err != nil {
		t.Fatal(err)
	}
}
