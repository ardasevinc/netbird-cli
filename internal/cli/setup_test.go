package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/version"
)

func TestSetupBootstrapJSONResolvesPasswordAndReturnsPATOnce(t *testing.T) {
	t.Setenv("NB_BOOTSTRAP_PASSWORD", "password123")
	var setupBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/instance":
			_ = json.NewEncoder(w).Encode(map[string]bool{"setup_required": true})
		case "/api/setup":
			if err := json.NewDecoder(r.Body).Decode(&setupBody); err != nil {
				t.Errorf("decode setup body: %v", err)
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]string{"personal_access_token": "nbpat-once"})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	var stdout, stderr bytes.Buffer
	root := newRoot(&commandState{json: true}, &stdout, &stderr, version.Current())
	root.SetArgs([]string{"setup", "bootstrap", "--from-json"})
	input := `{"url":"` + server.URL + `","email":"admin@example.com","name":"Admin","password_ref":"env:NB_BOOTSTRAP_PASSWORD","create_pat":true,"pat_expire_in":30}`
	root.SetIn(strings.NewReader(input))
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if setupBody["password"] != "password123" || setupBody["create_pat"] != true {
		t.Fatalf("unexpected setup body: %v", setupBody)
	}
	if !strings.Contains(stdout.String(), "nbpat-once") || strings.Contains(stdout.String(), "password123") {
		t.Fatalf("unexpected output: %s", stdout.String())
	}
}

func TestSetupBootstrapRejectsLiteralPasswordBeforeNetwork(t *testing.T) {
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		http.Error(w, "unexpected request", http.StatusInternalServerError)
	}))
	defer server.Close()
	var stdout, stderr bytes.Buffer
	root := newRoot(&commandState{json: true}, &stdout, &stderr, version.Current())
	root.SetArgs([]string{"setup", "bootstrap", "--from-json"})
	root.SetIn(strings.NewReader(`{"url":"` + server.URL + `","email":"admin@example.com","name":"Admin","password":"password123"}`))
	if err := root.ExecuteContext(context.Background()); err == nil || !strings.Contains(err.Error(), "unknown field") {
		t.Fatalf("err=%v", err)
	}
	if called {
		t.Fatal("literal password input reached the network")
	}
}

func TestSetupBootstrapRefusesInitializedServer(t *testing.T) {
	calledSetup := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/setup" {
			calledSetup = true
		}
		_ = json.NewEncoder(w).Encode(map[string]bool{"setup_required": false})
	}))
	defer server.Close()
	t.Setenv("NB_BOOTSTRAP_PASSWORD", "")
	var stdout, stderr bytes.Buffer
	root := newRoot(&commandState{json: true}, &stdout, &stderr, version.Current())
	root.SetArgs([]string{"setup", "bootstrap", "--url", server.URL, "--email", "admin@example.com", "--name", "Admin", "--password-ref", "env:NB_BOOTSTRAP_PASSWORD"})
	if err := root.ExecuteContext(context.Background()); err == nil || !strings.Contains(err.Error(), "does not require setup") {
		t.Fatalf("err=%v", err)
	}
	if calledSetup {
		t.Fatal("bootstrap POST was sent to an initialized server")
	}
}
