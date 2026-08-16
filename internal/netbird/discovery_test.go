package netbird

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/transport"
)

func TestDiscoverUsesReadOnlyEndpoints(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s", r.Method)
		}
		switch r.URL.Path {
		case "/api/instance/version":
			_ = json.NewEncoder(w).Encode(Version{ManagementCurrentVersion: "0.77.0"})
		case "/api/instance":
			_ = json.NewEncoder(w).Encode(Instance{SetupRequired: false})
		case "/api/users/current":
			_ = json.NewEncoder(w).Encode(User{ID: "user-1", Role: "admin"})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	client, err := transport.New(transport.Config{BaseURL: server.URL, HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	result, err := NewClient(client).Discover(context.Background(), true)
	if err != nil {
		t.Fatal(err)
	}
	if result.Version.ManagementCurrentVersion != "0.77.0" || result.User == nil || result.User.ID != "user-1" {
		t.Fatalf("unexpected discovery: %+v", result)
	}
}
