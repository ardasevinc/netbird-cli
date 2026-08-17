package netbird

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/transport"
)

func TestIdentityAndPostureReadsAreGETOnly(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		switch r.URL.Path {
		case "/api/identity-providers":
			_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "idp-1", "name": "zitadel", "type": "oidc", "issuer": "https://idp.example", "client_id": "client-1", "client_secret": "must-not-leak"}})
		case "/api/identity-providers/idp-1":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "idp-1", "name": "zitadel", "type": "oidc", "issuer": "https://idp.example", "client_id": "client-1"})
		case "/api/posture-checks":
			_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "pc-1", "name": "managed", "description": nil, "checks": map[string]any{"os_version_check": map[string]any{"linux": map[string]any{"min_kernel_version": "6.1"}}}}})
		case "/api/posture-checks/pc-1":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "pc-1", "name": "managed", "checks": map[string]any{}})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	transportClient, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	client := NewClient(transportClient)
	providers, err := client.ListIdentityProviders(context.Background())
	if err != nil || len(providers) != 1 || providers[0].Name != "zitadel" {
		t.Fatalf("providers=%+v err=%v", providers, err)
	}
	provider, err := client.GetIdentityProvider(context.Background(), "idp-1")
	if err != nil || provider.ClientID != "client-1" {
		t.Fatalf("provider=%+v err=%v", provider, err)
	}
	checks, err := client.ListPostureChecks(context.Background())
	if err != nil || len(checks) != 1 || checks[0].Name != "managed" || checks[0].Checks == nil {
		t.Fatalf("checks=%+v err=%v", checks, err)
	}
	check, err := client.GetPostureCheck(context.Background(), "pc-1")
	if err != nil || check.ID != "pc-1" {
		t.Fatalf("check=%+v err=%v", check, err)
	}
}

func TestCreatePostureCheckUsesPOST(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/posture-checks" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var request map[string]any
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		if request["name"] != "managed" {
			t.Fatalf("unexpected request: %#v", request)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "pc-2", "name": "managed", "checks": map[string]any{}})
	}))
	defer server.Close()
	transportClient, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	response, err := NewClient(transportClient).CreatePostureCheck(context.Background(), json.RawMessage(`{"name":"managed","checks":{}}`))
	if err != nil {
		t.Fatal(err)
	}
	var created map[string]any
	if err := json.Unmarshal(response, &created); err != nil || created["id"] != "pc-2" {
		t.Fatalf("created=%s err=%v", response, err)
	}
}

func TestUpdatePostureCheckUsesPUT(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/api/posture-checks/pc-1" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var request map[string]any
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		if request["name"] != "managed-v2" {
			t.Fatalf("unexpected request: %#v", request)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "pc-1", "name": "managed-v2", "checks": map[string]any{}})
	}))
	defer server.Close()
	transportClient, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	response, err := NewClient(transportClient).UpdatePostureCheck(context.Background(), "pc-1", json.RawMessage(`{"name":"managed-v2","checks":{}}`))
	if err != nil {
		t.Fatal(err)
	}
	var updated map[string]any
	if err := json.Unmarshal(response, &updated); err != nil || updated["name"] != "managed-v2" {
		t.Fatalf("updated=%s err=%v", response, err)
	}
}

func TestDeletePostureCheckUsesDELETE(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/api/posture-checks/pc-1" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	transportClient, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewClient(transportClient).DeletePostureCheck(context.Background(), "pc-1"); err != nil {
		t.Fatal(err)
	}
}
