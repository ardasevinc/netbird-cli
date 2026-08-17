package netbird

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/transport"
)

func TestEDRBypassRoutesAndReadModel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/peers/edr/bypassed":
			_ = json.NewEncoder(w).Encode([]map[string]string{{"peer_id": "peer-1"}})
		case r.Method == http.MethodPost && r.URL.Path == "/api/peers/peer-1/edr/bypass":
			_ = json.NewEncoder(w).Encode(map[string]string{"peer_id": "peer-1"})
		case r.Method == http.MethodDelete && r.URL.Path == "/api/peers/peer-1/edr/bypass":
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()
	transportClient, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	client := NewClient(transportClient)
	peers, err := client.ListEDRBypassedPeers(context.Background())
	if err != nil || len(peers) != 1 || peers[0].PeerID != "peer-1" {
		t.Fatalf("peers=%+v err=%v", peers, err)
	}
	if _, err := client.BypassPeerEDR(context.Background(), "peer-1"); err != nil {
		t.Fatal(err)
	}
	if _, err := client.RevokePeerEDRBypass(context.Background(), "peer-1"); err != nil {
		t.Fatal(err)
	}
}
