package netbird

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/transport"
)

func TestEventStreamingIntegrationRoutes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/event-streaming":
			_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "stream-1", "enabled": true}})
		case r.Method == http.MethodGet && r.URL.EscapedPath() == "/api/event-streaming/stream%2Fone":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "stream/one", "enabled": true})
		case r.Method == http.MethodPost && r.URL.Path == "/api/event-streaming":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "stream-1", "enabled": true})
		case r.Method == http.MethodPut && r.URL.EscapedPath() == "/api/event-streaming/stream%2Fone":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "stream/one", "enabled": false})
		case r.Method == http.MethodDelete && r.URL.EscapedPath() == "/api/event-streaming/stream%2Fone":
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
	if _, err := nb.ListEventStreamingIntegrationsRaw(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := nb.GetEventStreamingIntegrationRaw(context.Background(), "stream/one"); err != nil {
		t.Fatal(err)
	}
	if _, err := nb.CreateEventStreamingIntegration(context.Background(), json.RawMessage(`{"platform":"s3","enabled":true,"config":{"bucket":"logs"}}`)); err != nil {
		t.Fatal(err)
	}
	if _, err := nb.UpdateEventStreamingIntegration(context.Background(), "stream/one", json.RawMessage(`{"enabled":false}`)); err != nil {
		t.Fatal(err)
	}
	if _, err := nb.DeleteEventStreamingIntegration(context.Background(), "stream/one"); err != nil {
		t.Fatal(err)
	}
}
