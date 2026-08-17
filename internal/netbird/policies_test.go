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

func TestUpdatePolicyUsesEscapedIDAndPUT(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.EscapedPath() != "/api/policies/policy%2Fone" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		if string(body) != `{"id":"policy/one","name":"updated"}` {
			t.Fatalf("unexpected body: %s", body)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "policy/one", "name": "updated", "rules": []any{}})
	}))
	defer server.Close()
	transportClient, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	result, err := NewClient(transportClient).UpdatePolicy(context.Background(), "policy/one", json.RawMessage(`{"id":"policy/one","name":"updated"}`))
	if err != nil {
		t.Fatal(err)
	}
	if string(result) != `{"id":"policy/one","name":"updated","rules":[]}` {
		t.Fatalf("unexpected result: %s", result)
	}
}

func TestDeletePolicyUsesDELETE(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/api/policies/policy-1" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	transportClient, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewClient(transportClient).DeletePolicy(context.Background(), "policy-1"); err != nil {
		t.Fatal(err)
	}
}

func TestCreatePolicyUsesPOST(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/policies" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
		var request map[string]any
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		if request["name"] != "allow-office" {
			t.Fatalf("unexpected request body: %+v", request)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "p1", "name": "allow-office", "enabled": true, "rules": []any{map[string]any{"action": "accept"}}})
	}))
	defer server.Close()
	transportClient, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewClient(transportClient).CreatePolicy(context.Background(), json.RawMessage(`{"name":"allow-office","enabled":true,"rules":[{"action":"accept"}]}`)); err != nil {
		t.Fatal(err)
	}
}
