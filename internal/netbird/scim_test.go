package netbird

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/transport"
)

func TestSCIMIntegrationRoutes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/integrations/scim-idp/one/token" || r.URL.Path == "/api/integrations/okta-scim-idp/one/token" {
			_ = json.NewEncoder(w).Encode(map[string]string{"auth_token": "nbs_secret"})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"id": 1, "enabled": true})
	}))
	defer server.Close()
	client, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	nb := NewClient(client)
	for _, provider := range []string{"scim", "okta_scim"} {
		if _, err := nb.ListSCIMIntegrationsRaw(context.Background(), provider); err != nil {
			t.Fatal(err)
		}
		if _, err := nb.GetSCIMIntegrationRaw(context.Background(), provider, "one"); err != nil {
			t.Fatal(err)
		}
		if _, err := nb.CreateSCIMIntegration(context.Background(), provider, json.RawMessage(`{"prefix":"one"}`)); err != nil {
			t.Fatal(err)
		}
		if _, err := nb.UpdateSCIMIntegration(context.Background(), provider, "one", json.RawMessage(`{"enabled":false}`)); err != nil {
			t.Fatal(err)
		}
		if _, err := nb.DeleteSCIMIntegration(context.Background(), provider, "one"); err != nil {
			t.Fatal(err)
		}
		if _, err := nb.RegenerateSCIMToken(context.Background(), provider, "one"); err != nil {
			t.Fatal(err)
		}
	}
}
