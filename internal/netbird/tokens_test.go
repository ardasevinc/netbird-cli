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

func TestPersonalAccessTokenReadsOmitSecretMaterial(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		payload := map[string]any{"id": "token-1", "name": "agent", "created_by": "user-1", "created_at": "2026-08-17T00:00:00Z", "expiration_date": "2026-09-17T00:00:00Z", "last_used": nil, "plain_token": "secret-token"}
		if r.URL.Path == "/api/users/user-1/tokens" {
			_ = json.NewEncoder(w).Encode([]map[string]any{payload})
			return
		}
		if r.URL.Path == "/api/users/user-1/tokens/token-1" {
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
	tokens, err := client.ListPersonalAccessTokens(context.Background(), "user-1")
	if err != nil || len(tokens) != 1 {
		t.Fatalf("tokens=%+v err=%v", tokens, err)
	}
	token, err := client.GetPersonalAccessToken(context.Background(), "user-1", "token-1")
	if err != nil || token.ID != "token-1" {
		t.Fatalf("token=%+v err=%v", token, err)
	}
	encoded, err := json.Marshal(token)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), "secret-token") {
		t.Fatalf("secret leaked: %s", encoded)
	}
}
