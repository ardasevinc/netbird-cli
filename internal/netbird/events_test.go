package netbird

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/transport"
)

func TestListAuditEventsIsGETOnly(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/events/audit" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "event-1", "activity": "peer connected", "activity_code": "peer_connected", "initiator_email": "a@example.com", "initiator_id": "user-1", "initiator_name": "Arda", "meta": map[string]string{"peer": "peer-1"}, "target_id": "peer-1", "timestamp": "2026-08-17T10:00:00Z"}})
	}))
	defer server.Close()
	transportClient, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	events, err := NewClient(transportClient).ListAuditEvents(context.Background())
	if err != nil || len(events) != 1 || events[0].ID != "event-1" || events[0].Meta["peer"] != "peer-1" {
		t.Fatalf("events=%+v err=%v", events, err)
	}
}
