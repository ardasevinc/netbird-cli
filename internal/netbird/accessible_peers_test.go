package netbird

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/transport"
)

func TestListAccessiblePeersIsGETOnly(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/peers/p1/accessible-peers" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "p2", "name": "db", "ip": "10.0.0.3", "groups": []any{}, "connected": true}})
	}))
	defer server.Close()
	transportClient, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	peers, err := NewClient(transportClient).ListAccessiblePeers(context.Background(), "p1")
	if err != nil || len(peers) != 1 || peers[0].ID != "p2" {
		t.Fatalf("peers=%+v err=%v", peers, err)
	}
}
