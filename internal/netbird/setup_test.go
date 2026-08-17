package netbird

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/transport"
)

func TestBootstrapUsesInstanceGuardAndDeclaredRoute(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/instance":
			if r.Method != http.MethodGet {
				t.Errorf("instance method = %s", r.Method)
				return
			}
			_ = json.NewEncoder(w).Encode(Instance{SetupRequired: true})
		case "/api/setup":
			if r.Method != http.MethodPost {
				t.Errorf("setup method = %s", r.Method)
				return
			}
			var body SetupRequest
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Errorf("decode setup body: %v", err)
				return
			}
			if body.Email != "admin@example.com" || body.Password != "password123" || body.Name != "Admin" || !body.CreatePAT || body.PATExpireIn != 30 {
				t.Errorf("unexpected setup body: %+v", body)
			}
			_ = json.NewEncoder(w).Encode(map[string]string{"personal_access_token": "nbpat-test"})
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
	instance, err := client.GetInstance(context.Background())
	if err != nil || !instance.SetupRequired {
		t.Fatalf("instance=%+v err=%v", instance, err)
	}
	result, err := client.Bootstrap(context.Background(), SetupRequest{Email: "admin@example.com", Password: "password123", Name: "Admin", CreatePAT: true, PATExpireIn: 30})
	if err != nil || result.PersonalAccessToken != "nbpat-test" {
		t.Fatalf("result=%+v err=%v", result, err)
	}
}
