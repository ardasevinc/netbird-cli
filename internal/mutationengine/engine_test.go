package mutationengine

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/ledger"
	"github.com/ardasevinc/netbird-cli/internal/mutation"
	"github.com/ardasevinc/netbird-cli/internal/transport"
)

type fakeRemote struct {
	identity     string
	account      string
	before       json.RawMessage
	after        json.RawMessage
	policyBefore json.RawMessage
	policyAfter  json.RawMessage
	updateErr    error
	updates      int
}

func (f *fakeRemote) ServerIdentity() string { return f.identity }

func (f *fakeRemote) AccountScope(_ context.Context, account string) error {
	if account != f.account {
		return errors.New("wrong account")
	}
	return nil
}

func (f *fakeRemote) GetGroup(_ context.Context, _ string) (json.RawMessage, error) {
	return append(json.RawMessage(nil), f.before...), nil
}

func (f *fakeRemote) UpdateGroup(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.before = append(json.RawMessage(nil), f.after...)
	return f.before, nil
}

func (f *fakeRemote) GetPolicyRaw(_ context.Context, _ string) (json.RawMessage, error) {
	return append(json.RawMessage(nil), f.policyBefore...), nil
}

func (f *fakeRemote) UpdatePolicy(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.policyBefore = append(json.RawMessage(nil), f.policyAfter...)
	return f.policyBefore, nil
}

func stageForTest(t *testing.T, store *ledger.Store, before, after string) ledger.Stage {
	t.Helper()
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "groups.update",
		Request:        json.RawMessage(`{"id":"g1","name":"new"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
	})
	if err != nil {
		t.Fatal(err)
	}
	return stage
}

func TestApplyJournalsAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"g1","name":"old","peers_count":0,"resources_count":0}`
	after := `{"id":"g1","name":"new","peers_count":0,"resources_count":0}`
	stage := stageForTest(t, store, before, after)
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", before: []byte(before), after: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1"})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 || result.AttemptID == "" {
		t.Fatalf("unexpected result: %+v updates=%d", result, remote.updates)
	}
	receipt, err := store.GetReceipt(context.Background(), result.AttemptID)
	if err != nil || receipt.State != string(mutation.ConfirmedSuccess) {
		t.Fatalf("receipt=%+v err=%v", receipt, err)
	}
}

func TestApplyRefusesDriftBeforeJournalingOrDispatch(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	stage := stageForTest(t, store, `{"id":"g1","name":"old"}`, `{"id":"g1","name":"new"}`)
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", before: []byte(`{"id":"g1","name":"someone-else"}`), after: []byte(`{"id":"g1","name":"new"}`)}
	_, err = Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1"})
	if err == nil || remote.updates != 0 {
		t.Fatalf("drift was not refused: err=%v updates=%d", err, remote.updates)
	}
	if _, err := store.GetAttempt(context.Background(), "missing"); err == nil {
		t.Fatal("unexpected attempt")
	}
}

func TestApplyDoesNotReplayAfterAmbiguousDispatch(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"g1","name":"old"}`
	after := `{"id":"g1","name":"new"}`
	stage := stageForTest(t, store, before, after)
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", before: []byte(before), after: []byte(after), updateErr: &transport.RequestError{Dispatched: true, Description: "connection lost"}}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1"})
	if err == nil || result.State != mutation.Unknown {
		t.Fatalf("unexpected ambiguous result: %+v err=%v", result, err)
	}
	attempt, err := store.GetAttempt(context.Background(), result.AttemptID)
	if err != nil || attempt.State != string(mutation.Unknown) {
		t.Fatalf("attempt=%+v err=%v", attempt, err)
	}
}

func TestApplyReturnsAlreadySatisfiedWithoutWrite(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	after := `{"id":"g1","name":"new"}`
	stage := stageForTest(t, store, after, after)
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", before: []byte(after), after: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1"})
	if err != nil || result.State != mutation.AlreadySatisfied || remote.updates != 0 {
		t.Fatalf("unexpected no-op result: %+v updates=%d err=%v", result, remote.updates, err)
	}
}

func TestApplyClassifiesDefinitiveRemoteRejection(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"g1","name":"old"}`
	after := `{"id":"g1","name":"new"}`
	stage := stageForTest(t, store, before, after)
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", before: []byte(before), after: []byte(after), updateErr: &transport.RequestError{Dispatched: true, StatusCode: 403, Description: "remote rejected request"}}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1"})
	if err == nil || result.State != mutation.DefinitivelyRejected {
		t.Fatalf("unexpected rejection result: %+v err=%v", result, err)
	}
}

func TestApplyRefusesChangedImpactEvidence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "groups.update",
		Request:        json.RawMessage(`{"id":"g1","name":"new"}`),
		Before:         json.RawMessage(`{"id":"g1","name":"old"}`),
		IntendedAfter:  json.RawMessage(`{"id":"g1","name":"new"}`),
		Impact:         json.RawMessage(`{"classification":"unknown"}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", before: []byte(`{"id":"g1","name":"old"}`), after: []byte(`{"id":"g1","name":"new"}`)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1"})
	if err == nil || result.AttemptID != "" || remote.updates != 0 {
		t.Fatalf("changed impact was not refused before dispatch: result=%+v err=%v updates=%d", result, err, remote.updates)
	}
}

func TestApplyDispatchesPolicyUpdateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"p1","name":"old","rules":[]}`
	after := `{"id":"p1","name":"new","rules":[]}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "policies.update",
		Request:        json.RawMessage(`{"id":"p1","name":"new","rules":[]}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
		Impact:         json.RawMessage(`{"classification":"metadata_only","reachability":"unchanged","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"high","evidence":["policy metadata changed without changing policy rules"],"completeness":{"state":"complete","reason":null}}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", policyBefore: []byte(before), policyAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected policy result: %+v updates=%d", result, remote.updates)
	}
}
