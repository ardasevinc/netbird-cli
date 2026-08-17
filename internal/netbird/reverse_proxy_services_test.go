package netbird

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/transport"
)

func TestReverseProxyServiceAndClusterRoutes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/reverse-proxies/clusters":
			_ = json.NewEncoder(w).Encode([]map[string]any{{"address": "proxy.example.com", "type": "account"}})
		case r.Method == http.MethodDelete && r.URL.EscapedPath() == "/api/reverse-proxies/clusters/proxy.example.com":
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodGet && r.URL.Path == "/api/reverse-proxies/services":
			_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "service-1", "domain": "app.example.com"}})
		case r.Method == http.MethodGet && r.URL.EscapedPath() == "/api/reverse-proxies/services/service%2Fone":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "service/one", "domain": "app.example.com"})
		case r.Method == http.MethodPost && r.URL.Path == "/api/reverse-proxies/services":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "service-1", "domain": "app.example.com"})
		case r.Method == http.MethodPut && r.URL.EscapedPath() == "/api/reverse-proxies/services/service%2Fone":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "service/one", "domain": "new.example.com"})
		case r.Method == http.MethodDelete && r.URL.EscapedPath() == "/api/reverse-proxies/services/service%2Fone":
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
	if _, err := nb.ListReverseProxyClustersRaw(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := nb.DeleteReverseProxyCluster(context.Background(), "proxy.example.com"); err != nil {
		t.Fatal(err)
	}
	if _, err := nb.ListReverseProxyServicesRaw(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := nb.GetReverseProxyServiceRaw(context.Background(), "service/one"); err != nil {
		t.Fatal(err)
	}
	if _, err := nb.CreateReverseProxyService(context.Background(), json.RawMessage(`{"name":"app","domain":"app.example.com"}`)); err != nil {
		t.Fatal(err)
	}
	if _, err := nb.UpdateReverseProxyService(context.Background(), "service/one", json.RawMessage(`{"name":"app-2"}`)); err != nil {
		t.Fatal(err)
	}
	if _, err := nb.DeleteReverseProxyService(context.Background(), "service/one"); err != nil {
		t.Fatal(err)
	}
}
