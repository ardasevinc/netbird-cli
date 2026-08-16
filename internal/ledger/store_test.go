package ledger

import (
	"context"
	"encoding/json"
	"testing"
)

func TestCreateGetAndCancelImmutableStage(t *testing.T) {
	store, err := Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	input := StageInput{
		Profile:        "default",
		ServerIdentity: "server-1",
		AccountID:      "account-1",
		Operation:      "groups.update",
		Request:        json.RawMessage(`{"name":"ops"}`),
		Before:         json.RawMessage(`{"name":"old"}`),
		IntendedAfter:  json.RawMessage(`{"name":"ops"}`),
		Findings:       []Finding{{Code: "info", Severity: "info", Message: "safe"}},
	}
	created, err := store.Create(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	fetched, err := store.Get(context.Background(), created.ID, created.Revision)
	if err != nil {
		t.Fatal(err)
	}
	if fetched.Digest != created.Digest || fetched.Cancelled {
		t.Fatalf("unexpected stage: %+v", fetched)
	}
	if err := store.Cancel(context.Background(), created.ID, created.Revision); err != nil {
		t.Fatal(err)
	}
	fetched, err = store.Get(context.Background(), created.ID, created.Revision)
	if err != nil {
		t.Fatal(err)
	}
	if !fetched.Cancelled {
		t.Fatal("stage was not cancelled")
	}
}
