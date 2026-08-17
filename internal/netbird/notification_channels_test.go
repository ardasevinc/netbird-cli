package netbird

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/transport"
)

func TestNotificationChannelRoutes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/integrations/notifications/channels":
			_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "channel-1", "type": "webhook", "enabled": true}})
		case r.Method == http.MethodGet && r.URL.EscapedPath() == "/api/integrations/notifications/channels/channel%2Fone":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "channel/one", "type": "webhook", "enabled": true})
		case r.Method == http.MethodPost && r.URL.Path == "/api/integrations/notifications/channels":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "channel-1", "type": "webhook", "enabled": true})
		case r.Method == http.MethodPut && r.URL.EscapedPath() == "/api/integrations/notifications/channels/channel%2Fone":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "channel/one", "type": "webhook", "enabled": false})
		case r.Method == http.MethodDelete && r.URL.EscapedPath() == "/api/integrations/notifications/channels/channel%2Fone":
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
	if _, err := nb.ListNotificationChannelsRaw(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := nb.GetNotificationChannelRaw(context.Background(), "channel/one"); err != nil {
		t.Fatal(err)
	}
	if _, err := nb.CreateNotificationChannel(context.Background(), json.RawMessage(`{"type":"webhook","enabled":true}`)); err != nil {
		t.Fatal(err)
	}
	if _, err := nb.UpdateNotificationChannel(context.Background(), "channel/one", json.RawMessage(`{"enabled":false}`)); err != nil {
		t.Fatal(err)
	}
	if _, err := nb.DeleteNotificationChannel(context.Background(), "channel/one"); err != nil {
		t.Fatal(err)
	}
}
