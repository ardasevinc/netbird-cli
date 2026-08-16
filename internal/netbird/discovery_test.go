package netbird

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestListAndGetPoliciesUseReadOnlyEndpoints(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		switch r.URL.Path {
		case "/api/policies":
			_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "policy-1", "name": "allow-ssh", "enabled": true, "rules": []any{}, "source_posture_checks": []string{}}})
		case "/api/policies/policy-1":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "policy-1", "name": "allow-ssh", "enabled": true, "rules": []any{}, "source_posture_checks": []string{}})
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
	policies, err := adapter.ListPolicies(context.Background())
	if err != nil || len(policies) != 1 || policies[0].Name != "allow-ssh" {
		t.Fatalf("policies=%+v err=%v", policies, err)
	}
	policy, err := adapter.GetPolicy(context.Background(), "policy-1")
	if err != nil || policy.ID == nil || *policy.ID != "policy-1" {
		t.Fatalf("policy=%+v err=%v", policy, err)
	}
}

func TestListAccountsUsersAndInvitesNeverExposeInviteTokens(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		switch r.URL.Path {
		case "/api/accounts":
			_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "acct-1", "domain": "example.com", "domain_category": "private", "created_at": "2026-01-01T00:00:00Z", "created_by": "user-1", "settings": map[string]any{}, "onboarding": map[string]any{"signup_form_pending": false, "onboarding_flow_pending": false}}})
		case "/api/users":
			if r.URL.Query().Get("service_user") != "true" {
				t.Fatalf("unexpected user filter: %s", r.URL.RawQuery)
			}
			_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "user-1", "email": "svc@example.com", "name": "svc", "role": "admin", "status": "active", "auto_groups": []string{}, "is_blocked": false, "pending_approval": false, "is_service_user": true}})
		case "/api/users/invites":
			_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "invite-1", "email": "new@example.com", "name": "New", "role": "user", "auto_groups": []string{}, "expires_at": "2026-01-02T00:00:00Z", "created_at": "2026-01-01T00:00:00Z", "expired": false, "invite_token": "nbi_secret"}})
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
	accounts, err := adapter.ListAccounts(context.Background())
	if err != nil || len(accounts) != 1 || accounts[0].ID != "acct-1" {
		t.Fatalf("accounts=%+v err=%v", accounts, err)
	}
	serviceUser := true
	users, err := adapter.ListUsers(context.Background(), &serviceUser)
	if err != nil || len(users) != 1 || users[0].ID != "user-1" {
		t.Fatalf("users=%+v err=%v", users, err)
	}
	invites, err := adapter.ListInvites(context.Background())
	if err != nil || len(invites) != 1 || invites[0].ID != "invite-1" {
		t.Fatalf("invites=%+v err=%v", invites, err)
	}
	encoded, err := json.Marshal(invites[0])
	if err != nil {
		t.Fatal(err)
	}
	if string(encoded) == "" || strings.Contains(string(encoded), "nbi_secret") {
		t.Fatalf("invite token leaked: %s", encoded)
	}
}
