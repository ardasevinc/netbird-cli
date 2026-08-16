package netbird

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/transport"
)

func TestNetworkChildReadsNormalizeTopology(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		switch r.URL.Path {
		case "/api/networks/n1/resources":
			_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "r1", "name": "db", "address": "10.0.0.0/24", "type": "subnet", "enabled": true, "groups": []any{}}})
		case "/api/networks/n1/resources/r1":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "r1", "name": "db", "address": "10.0.0.0/24", "type": "subnet", "enabled": true, "groups": []any{}})
		case "/api/networks/n1/routers":
			_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "router-1", "enabled": true, "masquerade": true, "metric": 10, "peer": nil, "peer_groups": nil}})
		case "/api/networks/n1/routers/router-1":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "router-1", "enabled": true, "masquerade": true, "metric": 10, "peer": "peer-1", "peer_groups": nil})
		case "/api/networks/routers":
			_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "router-1", "enabled": true, "masquerade": true, "metric": 10, "peer": "peer-1", "peer_groups": nil}})
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
	resources, err := client.ListNetworkResources(context.Background(), "n1")
	if err != nil || len(resources) != 1 {
		t.Fatalf("resources=%+v err=%v", resources, err)
	}
	resource, err := client.GetNetworkResource(context.Background(), "n1", "r1")
	if err != nil || resource.Address != "10.0.0.0/24" {
		t.Fatalf("resource=%+v err=%v", resource, err)
	}
	routers, err := client.ListNetworkRouters(context.Background(), "n1")
	if err != nil || len(routers) != 1 || routers[0].PeerGroups == nil {
		t.Fatalf("routers=%+v err=%v", routers, err)
	}
	router, err := client.GetNetworkRouter(context.Background(), "n1", "router-1")
	if err != nil || router.Peer == nil {
		t.Fatalf("router=%+v err=%v", router, err)
	}
	all, err := client.ListAllNetworkRouters(context.Background())
	if err != nil || len(all) != 1 {
		t.Fatalf("all=%+v err=%v", all, err)
	}
}

func TestDeleteNetworkResourceUsesDELETE(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/api/networks/n1/resources/r1" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	client, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewClient(client).DeleteNetworkResource(context.Background(), "n1", "r1"); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateNetworkResourceUsesPUT(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/api/networks/n1/resources/r1" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
		var request map[string]any
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		if request["name"] != "new" {
			t.Fatalf("unexpected request body: %+v", request)
		}
		if _, ok := request["network_id"]; ok {
			t.Fatalf("target metadata leaked into request body: %+v", request)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "r1", "name": "new", "address": "10.0.0.0/24", "enabled": true})
	}))
	defer server.Close()
	client, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewClient(client).UpdateNetworkResource(context.Background(), "n1", "r1", json.RawMessage(`{"name":"new"}`)); err != nil {
		t.Fatal(err)
	}
}

func TestCreateNetworkResourceUsesPOST(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/networks/n1/resources" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
		var request map[string]any
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		if request["name"] != "new" {
			t.Fatalf("unexpected request body: %+v", request)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "r1", "name": "new", "address": "10.0.0.0/24", "enabled": true})
	}))
	defer server.Close()
	client, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewClient(client).CreateNetworkResource(context.Background(), "n1", json.RawMessage(`{"name":"new","address":"10.0.0.0/24","enabled":true,"groups":[]}`)); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteNetworkRouterUsesDELETE(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/api/networks/n1/routers/rt1" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	client, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewClient(client).DeleteNetworkRouter(context.Background(), "n1", "rt1"); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateNetworkRouterUsesPUT(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/api/networks/n1/routers/rt1" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
		var request map[string]any
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		if request["enabled"] != false {
			t.Fatalf("unexpected request body: %+v", request)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "rt1", "enabled": false, "masquerade": true, "metric": 10, "peer_groups": []string{}})
	}))
	defer server.Close()
	client, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewClient(client).UpdateNetworkRouter(context.Background(), "n1", "rt1", json.RawMessage(`{"enabled":false}`)); err != nil {
		t.Fatal(err)
	}
}
