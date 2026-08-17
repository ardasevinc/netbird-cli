package netbird

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/transport"
)

func TestGoogleIDPRoutes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/integrations/google-idp":
			_ = json.NewEncoder(w).Encode([]map[string]any{{"id": 1, "enabled": true}})
		case r.Method == http.MethodGet && r.URL.EscapedPath() == "/api/integrations/google-idp/one%2Fone":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": 1, "enabled": true})
		case r.Method == http.MethodPost && r.URL.Path == "/api/integrations/google-idp":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": 1, "enabled": true})
		case r.Method == http.MethodPut && r.URL.EscapedPath() == "/api/integrations/google-idp/one%2Fone":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": 1, "enabled": false})
		case r.Method == http.MethodDelete && r.URL.EscapedPath() == "/api/integrations/google-idp/one%2Fone":
			_ = json.NewEncoder(w).Encode(map[string]any{})
		case r.Method == http.MethodPost && r.URL.EscapedPath() == "/api/integrations/google-idp/one%2Fone/sync":
			_ = json.NewEncoder(w).Encode(map[string]any{"result": "ok"})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	client, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	nb := NewClient(client)
	if _, err := nb.ListGoogleIDPsRaw(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := nb.GetGoogleIDPRaw(context.Background(), "one/one"); err != nil {
		t.Fatal(err)
	}
	if _, err := nb.CreateGoogleIDP(context.Background(), json.RawMessage(`{"customer_id":"customer"}`)); err != nil {
		t.Fatal(err)
	}
	if _, err := nb.UpdateGoogleIDP(context.Background(), "one/one", json.RawMessage(`{"enabled":false}`)); err != nil {
		t.Fatal(err)
	}
	if _, err := nb.DeleteGoogleIDP(context.Background(), "one/one"); err != nil {
		t.Fatal(err)
	}
	if _, err := nb.SyncGoogleIDP(context.Background(), "one/one"); err != nil {
		t.Fatal(err)
	}
}
