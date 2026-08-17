package netbird

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/transport"
)

func TestReverseProxyDomainRoutes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/reverse-proxies/domains":
			_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "domain-1", "domain": "app.example.com"}})
		case r.Method == http.MethodPost && r.URL.Path == "/api/reverse-proxies/domains":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "domain-1", "domain": "app.example.com"})
		case r.Method == http.MethodDelete && r.URL.EscapedPath() == "/api/reverse-proxies/domains/domain%2Fone":
			w.WriteHeader(http.StatusNoContent)
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
	if _, err := nb.ListReverseProxyDomainsRaw(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := nb.CreateReverseProxyDomain(context.Background(), json.RawMessage(`{"domain":"app.example.com"}`)); err != nil {
		t.Fatal(err)
	}
	if _, err := nb.DeleteReverseProxyDomain(context.Background(), "domain/one"); err != nil {
		t.Fatal(err)
	}
}
