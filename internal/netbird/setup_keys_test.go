package netbird

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/transport"
)

func TestSetupKeyInventoryOmitsSecret(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		payload := map[string]any{"id": "key-1", "name": "bootstrap", "type": "reusable", "state": "valid", "valid": true, "revoked": false, "ephemeral": false, "allow_extra_dns_labels": false, "auto_groups": nil, "usage_limit": 5, "used_times": 1, "expires": "2026-09-01T00:00:00Z", "last_used": "2026-08-17T00:00:00Z", "updated_at": "2026-08-17T00:00:00Z", "key": "super-secret"}
		if r.URL.Path == "/api/setup-keys" {
			_ = json.NewEncoder(w).Encode([]map[string]any{payload})
			return
		}
		if r.URL.Path == "/api/setup-keys/key-1" {
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
	keys, err := client.ListSetupKeys(context.Background())
	if err != nil || len(keys) != 1 || keys[0].AutoGroups == nil {
		t.Fatalf("keys=%+v err=%v", keys, err)
	}
	key, err := client.GetSetupKey(context.Background(), "key-1")
	if err != nil || key.ID != "key-1" {
		t.Fatalf("key=%+v err=%v", key, err)
	}
	encoded, err := json.Marshal(key)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), "super-secret") {
		t.Fatalf("secret leaked: %s", encoded)
	}
}
