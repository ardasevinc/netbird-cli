package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetJSONSendsBearerAndBoundsBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer token" {
			t.Errorf("authorization = %q", r.Header.Get("Authorization"))
		}
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()
	client, err := New(Config{BaseURL: server.URL, Token: "token", HTTP: server.Client(), MaxBody: 64})
	if err != nil {
		t.Fatal(err)
	}
	var result map[string]bool
	if err := client.GetJSON(context.Background(), "/api/instance", &result); err != nil {
		t.Fatal(err)
	}
	if !result["ok"] {
		t.Fatalf("result = %v", result)
	}
}

func TestGetJSONRejectsOversizedBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"long":"this exceeds the configured bound"}`))
	}))
	defer server.Close()
	client, err := New(Config{BaseURL: server.URL, HTTP: server.Client(), MaxBody: 8})
	if err != nil {
		t.Fatal(err)
	}
	if err := client.GetJSON(context.Background(), "/api/instance", &json.RawMessage{}); err == nil {
		t.Fatal("expected oversized response to fail")
	}
}

func TestRequestPreservesEscapedPathSegments(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/api/policies/policy%2Fone" {
			t.Errorf("escaped path = %q", r.URL.EscapedPath())
		}
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()
	client, err := New(Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	var result map[string]bool
	if err := client.GetJSON(context.Background(), "/api/policies/policy%2Fone", &result); err != nil {
		t.Fatal(err)
	}
}

func TestDebugLoggingDoesNotIncludeBearerToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()
	var logs bytes.Buffer
	client, err := New(Config{BaseURL: server.URL, Token: "secret-token", HTTP: server.Client(), LogLevel: "debug", LogWriter: &logs})
	if err != nil {
		t.Fatal(err)
	}
	if err := client.GetJSON(context.Background(), "/api/instance", &map[string]bool{}); err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(logs.Bytes(), []byte("/api/instance")) || !bytes.Contains(logs.Bytes(), []byte("status=200")) {
		t.Fatalf("missing request diagnostics: %s", logs.String())
	}
	if bytes.Contains(logs.Bytes(), []byte("secret-token")) {
		t.Fatalf("bearer token leaked in diagnostics: %s", logs.String())
	}
}
