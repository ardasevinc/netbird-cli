package netbird

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/transport"
)

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
