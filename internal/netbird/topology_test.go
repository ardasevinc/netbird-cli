package netbird

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/transport"
)

func TestListAndGetTopologyReadsNormalizeOptionalFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		switch r.URL.Path {
		case "/api/routes":
			_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "route-1", "description": "private subnet", "enabled": true, "groups": []string{}, "keep_route": false, "masquerade": true, "metric": 10, "network": "10.0.0.0/24", "network_id": "", "network_type": "static", "peer": nil, "peer_groups": nil, "domains": nil}})
		case "/api/routes/route-1":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "route-1", "description": "private subnet", "enabled": true, "groups": []string{}, "keep_route": false, "masquerade": true, "metric": 10, "network": "10.0.0.0/24", "network_id": "", "network_type": "static", "peer": nil, "peer_groups": nil, "domains": nil})
		case "/api/networks":
			_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "network-1", "name": "office", "policies": []string{"policy-1"}, "resources": []string{}, "routers": []string{}, "routing_peers_count": 1}})
		case "/api/networks/network-1":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "network-1", "name": "office", "description": "office network", "policies": []string{"policy-1"}, "resources": []string{}, "routers": []string{}, "routing_peers_count": 1})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	client, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	adapter := NewClient(client)
	routes, err := adapter.ListRoutes(context.Background())
	if err != nil || len(routes) != 1 || routes[0].ID != "route-1" || routes[0].Domains == nil || routes[0].PeerGroups == nil {
		t.Fatalf("routes=%+v err=%v", routes, err)
	}
	route, err := adapter.GetRoute(context.Background(), "route-1")
	if err != nil || route.Network == nil || *route.Network != "10.0.0.0/24" {
		t.Fatalf("route=%+v err=%v", route, err)
	}
	networks, err := adapter.ListNetworks(context.Background())
	if err != nil || len(networks) != 1 || networks[0].Name != "office" {
		t.Fatalf("networks=%+v err=%v", networks, err)
	}
	network, err := adapter.GetNetwork(context.Background(), "network-1")
	if err != nil || network.Description == nil || *network.Description != "office network" {
		t.Fatalf("network=%+v err=%v", network, err)
	}
}
