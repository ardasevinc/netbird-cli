package netbird

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/transport"
)

func TestUserMutationMethodsUseDeclaredRoutes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/users":
			_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "user-1"}})
		case r.Method == http.MethodGet && r.URL.EscapedPath() == "/api/users/user%2Fone":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "user/one"})
		case r.Method == http.MethodPost && r.URL.Path == "/api/users":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "user-1"})
		case r.Method == http.MethodPut && r.URL.EscapedPath() == "/api/users/user%2Fone":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "user/one", "role": "admin"})
		case r.Method == http.MethodDelete && r.URL.EscapedPath() == "/api/users/user%2Fone":
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodPost && r.URL.EscapedPath() == "/api/users/user%2Fone/approve":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "user/one", "pending_approval": false})
		case r.Method == http.MethodDelete && r.URL.EscapedPath() == "/api/users/user%2Fone/reject":
			w.WriteHeader(http.StatusNoContent)
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
	ctx := context.Background()
	if _, err := client.ListUsersRaw(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := client.GetUserRaw(ctx, "user/one"); err != nil {
		t.Fatal(err)
	}
	if _, err := client.CreateUser(ctx, json.RawMessage(`{"email":"a@example.com"}`)); err != nil {
		t.Fatal(err)
	}
	if _, err := client.UpdateUser(ctx, "user/one", json.RawMessage(`{"role":"admin"}`)); err != nil {
		t.Fatal(err)
	}
	if _, err := client.ApproveUser(ctx, "user/one"); err != nil {
		t.Fatal(err)
	}
	if _, err := client.RejectUser(ctx, "user/one"); err != nil {
		t.Fatal(err)
	}
	if _, err := client.DeleteUser(ctx, "user/one"); err != nil {
		t.Fatal(err)
	}
}

func TestInviteMutationMethodsUseDeclaredRoutes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/users/invites":
			_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "invite-1"}})
		case r.Method == http.MethodPost && r.URL.Path == "/api/users/invites":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "invite-1", "invite_token": "one-time-invite"})
		case r.Method == http.MethodDelete && r.URL.EscapedPath() == "/api/users/invites/invite%2Fone":
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodPost && r.URL.EscapedPath() == "/api/users/invites/invite%2Fone/regenerate":
			_ = json.NewEncoder(w).Encode(map[string]any{"invite_token": "replacement-invite"})
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
	ctx := context.Background()
	if _, err := client.ListInvitesRaw(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := client.CreateInvite(ctx, json.RawMessage(`{"email":"a@example.com"}`)); err != nil {
		t.Fatal(err)
	}
	if _, err := client.DeleteInvite(ctx, "invite/one"); err != nil {
		t.Fatal(err)
	}
	if _, err := client.RegenerateInvite(ctx, "invite/one", json.RawMessage(`{"expires_in":3600}`)); err != nil {
		t.Fatal(err)
	}
}

func TestAccountReadsAndUpdatesUseAccountScopedEndpoints(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/accounts/account-1" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		switch r.Method {
		case http.MethodGet:
			_, _ = w.Write([]byte(`{"id":"account-1","settings":{"peer_login_expiration_enabled":true}}`))
		case http.MethodPut:
			var request map[string]any
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				t.Fatal(err)
			}
			settings, ok := request["settings"].(map[string]any)
			if !ok || settings["peer_login_expiration_enabled"] != false {
				t.Fatalf("unexpected request: %#v", request)
			}
			_, _ = w.Write([]byte(`{"id":"account-1","settings":{"peer_login_expiration_enabled":false}}`))
		case http.MethodDelete:
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
	before, err := client.GetAccountRaw(context.Background(), "account-1")
	if err != nil || string(before) == "" {
		t.Fatalf("before=%s err=%v", before, err)
	}
	after, err := client.UpdateAccount(context.Background(), "account-1", json.RawMessage(`{"settings":{"peer_login_expiration_enabled":false}}`))
	if err != nil {
		t.Fatal(err)
	}
	var result map[string]any
	if err := json.Unmarshal(after, &result); err != nil || result["id"] != "account-1" {
		t.Fatalf("after=%s err=%v", after, err)
	}
	if _, err := client.DeleteAccount(context.Background(), "account-1"); err != nil {
		t.Fatal(err)
	}
}
