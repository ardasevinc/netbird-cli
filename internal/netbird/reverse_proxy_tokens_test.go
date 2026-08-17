package netbird

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/transport"
)

func TestReverseProxyTokenRoutes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/reverse-proxies/proxy-tokens":
			_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "token-1", "name": "byop"}})
		case r.Method == http.MethodPost && r.URL.Path == "/api/reverse-proxies/proxy-tokens":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "token-1", "name": "byop", "plain_token": "one-time-token"})
		case r.Method == http.MethodDelete && r.URL.EscapedPath() == "/api/reverse-proxies/proxy-tokens/token%2Fone":
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
	if _, err := nb.ListReverseProxyTokensRaw(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := nb.CreateReverseProxyToken(context.Background(), json.RawMessage(`{"name":"byop","expires_in":0}`)); err != nil {
		t.Fatal(err)
	}
	if _, err := nb.DeleteReverseProxyToken(context.Background(), "token/one"); err != nil {
		t.Fatal(err)
	}
}
