package netbird

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/transport"
)

func TestLocationReadsNormalizeEmptyArrays(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		switch r.URL.Path {
		case "/api/locations/countries":
			_ = json.NewEncoder(w).Encode([]string{"DE", "TR"})
		case "/api/locations/countries/TR/cities":
			_ = json.NewEncoder(w).Encode([]string{"Istanbul"})
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
	countries, err := client.ListCountries(context.Background())
	if err != nil || len(countries) != 2 {
		t.Fatalf("countries=%v err=%v", countries, err)
	}
	cities, err := client.ListCountryCities(context.Background(), "TR")
	if err != nil || len(cities) != 1 || cities[0] != "Istanbul" {
		t.Fatalf("cities=%v err=%v", cities, err)
	}
}
