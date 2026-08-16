package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/version"
)

func TestStageCreateShowAndCancel(t *testing.T) {
	temp := t.TempDir()
	configPath := filepath.Join(temp, "config.toml")
	statePath := filepath.Join(temp, "ledger.db")
	if err := os.WriteFile(configPath, []byte("[profiles.default]\nurl = \"https://netbird.example.test\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	state := &commandState{json: true, configPath: configPath, profileName: "default", statePath: statePath}
	var stdout, stderr bytes.Buffer
	root := newRoot(state, &stdout, &stderr, version.Current())
	root.SetArgs([]string{"stage", "create", "--from-json"})
	root.SetIn(strings.NewReader(`{"operation":"groups.update","request":{"id":"g1"},"before":{"name":"old"},"intended_after":{"name":"new"}}`))
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	var created struct {
		Data struct {
			StageID  string `json:"stage_id"`
			Revision int    `json:"revision"`
		} `json:"data"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	if created.Data.StageID == "" || created.Data.Revision != 1 {
		t.Fatalf("unexpected stage result: %s", stdout.String())
	}
	if !strings.Contains(stdout.String(), `"classification":"metadata_only"`) {
		t.Fatalf("stage result omitted impact evidence: %s", stdout.String())
	}
	stdout.Reset()
	root.SetArgs([]string{"stage", "show", created.Data.StageID + "@1"})
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	stdout.Reset()
	root.SetArgs([]string{"stage", "cancel", created.Data.StageID + "@1"})
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), `"cancelled":true`) {
		t.Fatalf("unexpected cancel result: %s", stdout.String())
	}
}

func TestStageCreatePolicyRuleChangeRequiresAcknowledgement(t *testing.T) {
	temp := t.TempDir()
	configPath := filepath.Join(temp, "config.toml")
	statePath := filepath.Join(temp, "ledger.db")
	if err := os.WriteFile(configPath, []byte("[profiles.default]\nurl = \"https://netbird.example.test\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	state := &commandState{json: true, configPath: configPath, profileName: "default", statePath: statePath}
	var stdout, stderr bytes.Buffer
	root := newRoot(state, &stdout, &stderr, version.Current())
	root.SetArgs([]string{"stage", "create", "--from-json"})
	root.SetIn(strings.NewReader(`{"operation":"policies.update","request":{"id":"p1"},"before":{"id":"p1","name":"policy","rules":[]},"intended_after":{"id":"p1","name":"policy","rules":[{"action":"accept"}]}}`))
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), `"code":"impact.policy_rule_change"`) || !strings.Contains(stdout.String(), `"severity":"blocking"`) {
		t.Fatalf("policy impact acknowledgement finding missing: %s", stdout.String())
	}
}
