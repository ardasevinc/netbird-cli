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

func TestEDRIntegrationRoutes(t *testing.T) {
	providers := []string{"intune", "sentinelone", "falcon", "huntress", "fleetdm"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, provider := range providers {
			if r.URL.Path != "/api/integrations/edr/"+provider {
				continue
			}
			if r.Method == http.MethodDelete {
				_ = json.NewEncoder(w).Encode(map[string]any{})
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"id": 1, "enabled": true})
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()
	client, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	nb := NewClient(client)
	for _, provider := range providers {
		if _, err := nb.GetEDRIntegrationRaw(context.Background(), provider); err != nil {
			t.Fatalf("get %s: %v", provider, err)
		}
		if _, err := nb.CreateEDRIntegration(context.Background(), provider, json.RawMessage(`{"enabled":true}`)); err != nil {
			t.Fatalf("create %s: %v", provider, err)
		}
		if _, err := nb.UpdateEDRIntegration(context.Background(), provider, json.RawMessage(`{"enabled":false}`)); err != nil {
			t.Fatalf("update %s: %v", provider, err)
		}
		if _, err := nb.DeleteEDRIntegration(context.Background(), provider); err != nil {
			t.Fatalf("delete %s: %v", provider, err)
		}
	}
	if _, err := nb.GetEDRIntegrationRaw(context.Background(), "unknown"); err == nil {
		t.Fatal("expected unsupported EDR provider error")
	}
}
