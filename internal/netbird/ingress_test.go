package netbird

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/transport"
)

func TestIngressReadsAreGETOnly(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		peer := map[string]any{"id": "ing-1", "peer_id": "peer-1", "ingress_ip": "203.0.113.10", "region": "eu", "connected": true, "enabled": true, "fallback": false, "available_ports": map[string]any{"tcp": 10, "udp": 20}}
		allocation := map[string]any{"id": "alloc-1", "name": "web", "ingress_ip": "203.0.113.10", "ingress_peer_id": "ing-1", "region": "eu", "enabled": true, "port_range_mappings": nil}
		switch r.URL.Path {
		case "/api/ingress/peers":
			_ = json.NewEncoder(w).Encode([]map[string]any{peer})
		case "/api/ingress/peers/ing-1":
			_ = json.NewEncoder(w).Encode(peer)
		case "/api/peers/peer-1/ingress/ports":
			_ = json.NewEncoder(w).Encode([]map[string]any{allocation})
		case "/api/peers/peer-1/ingress/ports/alloc-1":
			_ = json.NewEncoder(w).Encode(allocation)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	transportClient, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	client := NewClient(transportClient)
	peers, err := client.ListIngressPeers(context.Background())
	if err != nil || len(peers) != 1 || peers[0].PeerID != "peer-1" {
		t.Fatalf("peers=%+v err=%v", peers, err)
	}
	peer, err := client.GetIngressPeer(context.Background(), "ing-1")
	if err != nil || peer.IngressIP == "" {
		t.Fatalf("peer=%+v err=%v", peer, err)
	}
	allocations, err := client.ListIngressPortAllocations(context.Background(), "peer-1")
	if err != nil || len(allocations) != 1 || allocations[0].PortRangeMappings == nil {
		t.Fatalf("allocations=%+v err=%v", allocations, err)
	}
	allocation, err := client.GetIngressPortAllocation(context.Background(), "peer-1", "alloc-1")
	if err != nil || allocation.ID != "alloc-1" {
		t.Fatalf("allocation=%+v err=%v", allocation, err)
	}
}

func TestCreateIngressPeerUsesPOST(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/ingress/peers" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var request map[string]any
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		if request["peer_id"] != "peer-1" {
			t.Fatalf("unexpected request: %#v", request)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "ing-2", "peer_id": "peer-1", "region": "eu", "enabled": true})
	}))
	defer server.Close()
	transportClient, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	response, err := NewClient(transportClient).CreateIngressPeer(context.Background(), json.RawMessage(`{"peer_id":"peer-1","region":"eu","enabled":true}`))
	if err != nil {
		t.Fatal(err)
	}
	var created map[string]any
	if err := json.Unmarshal(response, &created); err != nil || created["id"] != "ing-2" {
		t.Fatalf("created=%s err=%v", response, err)
	}
}

func TestUpdateIngressPeerUsesPUT(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/api/ingress/peers/ing-1" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var request map[string]any
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		if request["enabled"] != false {
			t.Fatalf("unexpected request: %#v", request)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "ing-1", "peer_id": "peer-1", "region": "eu", "enabled": false})
	}))
	defer server.Close()
	transportClient, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	response, err := NewClient(transportClient).UpdateIngressPeer(context.Background(), "ing-1", json.RawMessage(`{"enabled":false}`))
	if err != nil {
		t.Fatal(err)
	}
	var updated map[string]any
	if err := json.Unmarshal(response, &updated); err != nil || updated["enabled"] != false {
		t.Fatalf("updated=%s err=%v", response, err)
	}
}

func TestDeleteIngressPeerUsesDELETE(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/api/ingress/peers/ing-1" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	transportClient, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewClient(transportClient).DeleteIngressPeer(context.Background(), "ing-1"); err != nil {
		t.Fatal(err)
	}
}

func TestIngressPortAllocationMutationsUsePeerScopedRoutes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/peers/peer-1/ingress/ports/alloc-1" && r.URL.Path != "/api/peers/peer-1/ingress/ports" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		switch r.Method {
		case http.MethodPost:
			if r.URL.Path != "/api/peers/peer-1/ingress/ports" {
				t.Fatalf("unexpected create path: %s", r.URL.Path)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "alloc-1", "name": "web", "enabled": true})
		case http.MethodPut:
			if r.URL.Path != "/api/peers/peer-1/ingress/ports/alloc-1" {
				t.Fatalf("unexpected update path: %s", r.URL.Path)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "alloc-1", "name": "web-updated", "enabled": false})
		case http.MethodDelete:
			if r.URL.Path != "/api/peers/peer-1/ingress/ports/alloc-1" {
				t.Fatalf("unexpected delete path: %s", r.URL.Path)
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected method: %s", r.Method)
		}
	}))
	defer server.Close()
	transportClient, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	client := NewClient(transportClient)
	if _, err := client.CreateIngressPortAllocation(context.Background(), "peer-1", json.RawMessage(`{"name":"web","enabled":true}`)); err != nil {
		t.Fatal(err)
	}
	if _, err := client.UpdateIngressPortAllocation(context.Background(), "peer-1", "alloc-1", json.RawMessage(`{"enabled":false}`)); err != nil {
		t.Fatal(err)
	}
	if _, err := client.DeleteIngressPortAllocation(context.Background(), "peer-1", "alloc-1"); err != nil {
		t.Fatal(err)
	}
}
