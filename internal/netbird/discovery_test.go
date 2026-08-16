package netbird

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/transport"
)

func TestDiscoverUsesReadOnlyEndpoints(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s", r.Method)
		}
		switch r.URL.Path {
		case "/api/instance/version":
			_ = json.NewEncoder(w).Encode(Version{ManagementCurrentVersion: "0.77.0"})
		case "/api/instance":
			_ = json.NewEncoder(w).Encode(Instance{SetupRequired: false})
		case "/api/users/current":
			_ = json.NewEncoder(w).Encode(User{ID: "user-1", Role: "admin"})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	client, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	result, err := NewClient(client).Discover(context.Background(), true)
	if err != nil {
		t.Fatal(err)
	}
	if result.Version.ManagementCurrentVersion != "0.77.0" || result.User == nil || result.User.ID != "user-1" {
		t.Fatalf("unexpected discovery: %+v", result)
	}
}

func TestListGroupsUsesBoundedRead(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/groups" || r.Method != http.MethodGet {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode([]Group{{ID: "g1", Name: "ops", PeersCount: 2}})
	}))
	defer server.Close()
	client, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	groups, err := NewClient(client).ListGroups(context.Background())
	if err != nil || len(groups) != 1 || groups[0].Name != "ops" {
		t.Fatalf("groups=%+v err=%v", groups, err)
	}
}

func TestListAndGetPeersNormalizeReadOnlyResponses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		switch r.URL.Path {
		case "/api/peers":
			if r.URL.Query().Get("name") != "laptop" || r.URL.Query().Get("ip") != "10.0.0.2" {
				t.Fatalf("unexpected query: %s", r.URL.RawQuery)
			}
			_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "p1", "name": "laptop", "ip": "10.0.0.2", "connected": true, "groups": []any{}}})
		case "/api/peers/p1":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "p1", "name": "laptop", "ip": "10.0.0.2", "connected": true, "version": "0.1", "groups": []any{}})
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
	peers, err := adapter.ListPeers(context.Background(), "laptop", "10.0.0.2")
	if err != nil || len(peers) != 1 || peers[0].ID != "p1" {
		t.Fatalf("peers=%+v err=%v", peers, err)
	}
	peer, err := adapter.GetPeer(context.Background(), "p1")
	if err != nil || peer.Version != "0.1" {
		t.Fatalf("peer=%+v err=%v", peer, err)
	}
}
