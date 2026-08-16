package netbird

import (
	"context"
	"encoding/json"
	"io"
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

func TestListNetworkTrafficEvents(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/events/network-traffic" || r.URL.Query().Get("page") != "2" || r.URL.Query().Get("page_size") != "10" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.String())
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"data":[{"flow_id":"flow-1","window_start":"2026-08-17T10:00:00Z","window_end":"2026-08-17T10:01:00Z","source":{"id":"peer-1"},"destination":{"id":"peer-2"}}],"page":2,"page_size":10,"total_pages":3,"total_records":21}`)
	}))
	defer server.Close()

	transportClient, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	result, err := NewClient(transportClient).ListNetworkTrafficEvents(context.Background(), EventPageOptions{Page: 2, PageSize: 10})
	if err != nil || result.Page != 2 || len(result.Data) != 1 || result.Data[0].FlowID != "flow-1" {
		t.Fatalf("result=%+v err=%v", result, err)
	}
}

func TestListProxyAccessLogs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/events/proxy" || r.URL.Query().Get("page_size") != "5" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.String())
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"data":[{"id":"log-1","timestamp":"2026-08-17T10:00:00Z","method":"GET","path":"/health","status_code":200}],"page":1,"page_size":5,"total_pages":1,"total_records":1}`)
	}))
	defer server.Close()

	transportClient, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	result, err := NewClient(transportClient).ListProxyAccessLogs(context.Background(), EventPageOptions{PageSize: 5})
	if err != nil || result.TotalRecords != 1 || len(result.Data) != 1 || result.Data[0].ID != "log-1" {
		t.Fatalf("result=%+v err=%v", result, err)
	}
}
