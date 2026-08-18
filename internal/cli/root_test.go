package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/version"
)

func TestRootUsesNBProfileWhenNoExplicitProfileIsProvided(t *testing.T) {
	t.Setenv("NB_PROFILE", "staging")
	state := &commandState{}
	root := newRoot(state, &bytes.Buffer{}, &bytes.Buffer{}, version.Current())
	if state.profileName != "staging" {
		t.Fatalf("profile=%q, want staging", state.profileName)
	}
	if got := root.PersistentFlags().Lookup("profile").DefValue; got != "staging" {
		t.Fatalf("profile flag default=%q, want staging", got)
	}
}

func TestRootProfileFlagBeatsNBProfile(t *testing.T) {
	t.Setenv("NB_PROFILE", "staging")
	state := &commandState{}
	root := newRoot(state, &bytes.Buffer{}, &bytes.Buffer{}, version.Current())
	root.SetArgs([]string{"--profile", "production", "version"})
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if state.profileName != "production" {
		t.Fatalf("profile=%q, want production", state.profileName)
	}
}

func TestRootExplicitStateAndProfileValuesBeatEnvironment(t *testing.T) {
	t.Setenv("NB_PROFILE", "staging")
	envState := filepath.Join(t.TempDir(), "env-ledger.db")
	t.Setenv("NB_STATE", envState)
	explicitState := filepath.Join(t.TempDir(), "explicit-ledger.db")
	state := &commandState{profileName: "production", statePath: explicitState}
	newRoot(state, &bytes.Buffer{}, &bytes.Buffer{}, version.Current())
	if state.profileName != "production" {
		t.Fatalf("profile=%q, want production", state.profileName)
	}
	if state.statePath != explicitState {
		t.Fatalf("state path=%q, want explicit %q", state.statePath, explicitState)
	}
}

func TestExecuteVersionJSON(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute(context.Background(), []string{"--json", "version"}, &stdout, &stderr, version.Info{Version: "1.2.3", Commit: "abc", Date: "now"})
	if code != 0 {
		t.Fatalf("exit code = %d, stderr = %q", code, stderr.String())
	}
	var result struct {
		Schema string `json:"schema"`
		OK     bool   `json:"ok"`
		Data   struct {
			Version string `json:"version"`
		} `json:"data"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode result: %v", err)
	}
	if result.Schema != "nb/v1/version-result" || !result.OK || result.Data.Version != "1.2.3" {
		t.Fatalf("unexpected result: %s", stdout.String())
	}
}

func TestExecuteSchemaList(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute(context.Background(), []string{"schema", "list"}, &stdout, &stderr, version.Current())
	if code != 0 || stderr.Len() != 0 {
		t.Fatalf("code=%d stderr=%q", code, stderr.String())
	}
	if stdout.String() == "" {
		t.Fatal("schema list was empty")
	}
}

func TestExecuteVersionJSONL(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute(context.Background(), []string{"--jsonl", "version"}, &stdout, &stderr, version.Info{Version: "1.2.3", Commit: "abc", Date: "now"})
	if code != 0 || stderr.Len() != 0 {
		t.Fatalf("code=%d stderr=%q", code, stderr.String())
	}
	lines := bytes.Split(bytes.TrimSpace(stdout.Bytes()), []byte("\n"))
	if len(lines) != 2 {
		t.Fatalf("expected record and complete lines, got %d: %s", len(lines), stdout.String())
	}
	var record struct {
		Schema    string `json:"schema"`
		Type      string `json:"type"`
		Operation string `json:"operation"`
		Data      struct {
			Version string `json:"version"`
		} `json:"data"`
	}
	if err := json.Unmarshal(lines[0], &record); err != nil {
		t.Fatal(err)
	}
	var complete struct {
		Type string `json:"type"`
		Data struct {
			Count int `json:"count"`
		} `json:"data"`
	}
	if err := json.Unmarshal(lines[1], &complete); err != nil {
		t.Fatal(err)
	}
	if record.Schema != streamSchema || record.Type != "record" || record.Operation != "version" || record.Data.Version != "1.2.3" || complete.Type != "complete" || complete.Data.Count != 1 {
		t.Fatalf("unexpected JSONL: %s", stdout.String())
	}
}

func TestExecuteJSONLErrorStream(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute(context.Background(), []string{"--jsonl", "does-not-exist"}, &stdout, &stderr, version.Current())
	if code == 0 || stderr.Len() == 0 {
		t.Fatalf("expected invalid input, code=%d stderr=%q", code, stderr.String())
	}
	var event struct {
		Schema string `json:"schema"`
		Type   string `json:"type"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &event); err != nil {
		t.Fatal(err)
	}
	if event.Schema != streamSchema || event.Type != "error" {
		t.Fatalf("unexpected JSONL error: %s", stdout.String())
	}
}
