package netbird

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/transport"
)

func TestMSPIntegrationRoutes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		if r.URL.Path == "/api/integrations/msp/tenants" && r.Method == http.MethodGet {
			_, _ = w.Write([]byte(`[{"id":"tenant/one","name":"before"}]`))
			return
		}
		if r.URL.Path == "/api/integrations/msp/tenants" && r.Method == http.MethodPost {
			_, _ = w.Write([]byte(`{"id":"tenant/created","name":"new"}`))
			return
		}
		if r.URL.EscapedPath() == "/api/integrations/msp/tenants/tenant%2Fone" && r.Method == http.MethodPut {
			_, _ = w.Write([]byte(`{"id":"tenant/one","name":"after"}`))
			return
		}
		if r.URL.EscapedPath() == "/api/integrations/msp/tenants/tenant%2Fone/dns" && r.Method == http.MethodPost {
			_, _ = w.Write([]byte(`{"ok":true}`))
			return
		}
		if r.URL.EscapedPath() == "/api/integrations/msp/tenants/tenant%2Fone/invite" && r.Method == http.MethodPost {
			_, _ = w.Write([]byte(`{"id":"tenant/one","status":"invited"}`))
			return
		}
		if r.URL.EscapedPath() == "/api/integrations/msp/tenants/tenant%2Fone/invite" && r.Method == http.MethodPut {
			_, _ = w.Write([]byte(`{"id":"tenant/one","status":"active"}`))
			return
		}
		if r.URL.EscapedPath() == "/api/integrations/msp/tenants/tenant%2Fone/subscription" && r.Method == http.MethodPost {
			_, _ = w.Write([]byte(`{"ok":true}`))
			return
		}
		if r.URL.EscapedPath() == "/api/integrations/msp/tenants/tenant%2Fone/unlink" && r.Method == http.MethodPost {
			_, _ = w.Write([]byte(`{"ok":true}`))
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	client, err := transport.New(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	api := NewClient(client)
	ctx := context.Background()
	if _, err := api.GetMSPTenantRaw(ctx, "tenant/one"); err != nil {
		t.Fatal(err)
	}
	if _, err := api.CreateMSPTenant(ctx, json.RawMessage(`{"name":"new"}`)); err != nil {
		t.Fatal(err)
	}
	if _, err := api.UpdateMSPTenant(ctx, "tenant/one", json.RawMessage(`{"name":"after"}`)); err != nil {
		t.Fatal(err)
	}
	if _, err := api.VerifyMSPTenantDNS(ctx, "tenant/one"); err != nil {
		t.Fatal(err)
	}
	if _, err := api.InviteMSPTenant(ctx, "tenant/one"); err != nil {
		t.Fatal(err)
	}
	if _, err := api.RespondMSPTenantInvite(ctx, "tenant/one", json.RawMessage(`{"value":"accept"}`)); err != nil {
		t.Fatal(err)
	}
	if _, err := api.CreateMSPTenantSubscription(ctx, "tenant/one", json.RawMessage(`{"priceID":"price"}`)); err != nil {
		t.Fatal(err)
	}
	if _, err := api.UnlinkMSPTenant(ctx, "tenant/one", json.RawMessage(`{"owner":"owner"}`)); err != nil {
		t.Fatal(err)
	}
}
