package netbird

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/transport"
)

func TestAzureIDPRoutes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/integrations/azure-idp":
			_ = json.NewEncoder(w).Encode([]map[string]any{{"id": 1, "enabled": true}})
		case r.Method == http.MethodGet && r.URL.EscapedPath() == "/api/integrations/azure-idp/one%2Fone":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": 1, "enabled": true})
		case r.Method == http.MethodPost && r.URL.Path == "/api/integrations/azure-idp":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": 1, "enabled": true})
		case r.Method == http.MethodPut && r.URL.EscapedPath() == "/api/integrations/azure-idp/one%2Fone":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": 1, "enabled": false})
		case r.Method == http.MethodDelete && r.URL.EscapedPath() == "/api/integrations/azure-idp/one%2Fone":
			_ = json.NewEncoder(w).Encode(map[string]any{})
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
	if _, err := nb.ListAzureIDPsRaw(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := nb.GetAzureIDPRaw(context.Background(), "one/one"); err != nil {
		t.Fatal(err)
	}
	if _, err := nb.CreateAzureIDP(context.Background(), json.RawMessage(`{"client_id":"client"}`)); err != nil {
		t.Fatal(err)
	}
	if _, err := nb.UpdateAzureIDP(context.Background(), "one/one", json.RawMessage(`{"enabled":false}`)); err != nil {
		t.Fatal(err)
	}
	if _, err := nb.DeleteAzureIDP(context.Background(), "one/one"); err != nil {
		t.Fatal(err)
	}
}
