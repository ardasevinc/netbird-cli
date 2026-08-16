package netbird

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/transport"
)

func TestDNSRecordReadsAreGETOnly(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		payload := map[string]any{"id": "record-1", "name": "db", "type": "A", "content": "10.0.0.5", "ttl": 60}
		if r.URL.Path == "/api/dns/zones/zone-1/records" {
			_ = json.NewEncoder(w).Encode([]map[string]any{payload})
			return
		}
		if r.URL.Path == "/api/dns/zones/zone-1/records/record-1" {
			_ = json.NewEncoder(w).Encode(payload)
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()
	transportClient, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	client := NewClient(transportClient)
	records, err := client.ListDNSRecords(context.Background(), "zone-1")
	if err != nil || len(records) != 1 || records[0].Content != "10.0.0.5" {
		t.Fatalf("records=%+v err=%v", records, err)
	}
	record, err := client.GetDNSRecord(context.Background(), "zone-1", "record-1")
	if err != nil || record.TTL != 60 {
		t.Fatalf("record=%+v err=%v", record, err)
	}
}
