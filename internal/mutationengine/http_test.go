package mutationengine

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/ledger"
	"github.com/ardasevinc/netbird-cli/internal/mutation"
	"github.com/ardasevinc/netbird-cli/internal/netbird"
	"github.com/ardasevinc/netbird-cli/internal/transport"
)

func TestApplyAgainstDisposableHTTPServer(t *testing.T) {
	const (
		accountID = "account-1"
		groupID   = "g1"
	)
	var mu sync.Mutex
	group := map[string]any{"id": groupID, "name": "old", "peers_count": 0, "resources_count": 0}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer disposable-token" {
			http.Error(w, "missing auth", http.StatusUnauthorized)
			return
		}
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/accounts":
			_ = json.NewEncoder(w).Encode([]map[string]string{{"id": accountID}})
		case r.Method == http.MethodGet && r.URL.Path == "/api/groups/"+groupID:
			mu.Lock()
			defer mu.Unlock()
			_ = json.NewEncoder(w).Encode(group)
		case r.Method == http.MethodPut && r.URL.Path == "/api/groups/"+groupID:
			var request map[string]any
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				http.Error(w, "bad json", http.StatusBadRequest)
				return
			}
			mu.Lock()
			group["name"] = request["name"]
			mu.Unlock()
			mu.Lock()
			defer mu.Unlock()
			_ = json.NewEncoder(w).Encode(group)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"g1","name":"old","peers_count":0,"resources_count":0}`
	after := `{"id":"g1","name":"new","peers_count":0,"resources_count":0}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: server.URL,
		AccountID:      accountID,
		Operation:      "groups.update",
		Request:        json.RawMessage(`{"id":"g1","name":"new"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
	})
	if err != nil {
		t.Fatal(err)
	}
	client, err := transport.New(transport.Config{BaseURL: server.URL, Token: "disposable-token", HTTP: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	result, err := Apply(context.Background(), store, netbird.NewClient(client), ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: server.URL, AccountID: accountID})
	if err != nil || result.State != mutation.ConfirmedSuccess {
		t.Fatalf("result=%+v err=%v", result, err)
	}
}
