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

func TestCreateDNSZoneUsesPOST(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/dns/zones" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var request map[string]any
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		if request["domain"] != "office.internal" {
			t.Fatalf("unexpected request: %#v", request)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "zone-1", "name": "office", "domain": "office.internal", "enabled": true,
			"distribution_groups": []string{},
		})
	}))
	defer server.Close()
	transportClient, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	client := NewClient(transportClient)
	response, err := client.CreateDNSZone(context.Background(), json.RawMessage(`{"name":"office","domain":"office.internal","enabled":true}`))
	if err != nil {
		t.Fatal(err)
	}
	var created map[string]any
	if err := json.Unmarshal(response, &created); err != nil || created["id"] != "zone-1" {
		t.Fatalf("created=%s err=%v", response, err)
	}
}

func TestDeleteDNSZoneUsesDELETE(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/api/dns/zones/zone-1" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	transportClient, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewClient(transportClient).DeleteDNSZone(context.Background(), "zone-1"); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateDNSZoneUsesPUT(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/api/dns/zones/zone-1" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var request map[string]any
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		if request["domain"] != "corp.internal" {
			t.Fatalf("unexpected request: %#v", request)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "zone-1", "domain": "corp.internal", "enabled": true})
	}))
	defer server.Close()
	transportClient, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	response, err := NewClient(transportClient).UpdateDNSZone(context.Background(), "zone-1", json.RawMessage(`{"domain":"corp.internal","enabled":true}`))
	if err != nil {
		t.Fatal(err)
	}
	var updated map[string]any
	if err := json.Unmarshal(response, &updated); err != nil || updated["domain"] != "corp.internal" {
		t.Fatalf("updated=%s err=%v", response, err)
	}
}
