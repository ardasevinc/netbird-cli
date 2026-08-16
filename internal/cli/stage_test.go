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

func TestStageCreateRouteChangeRequiresAcknowledgement(t *testing.T) {
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
	root.SetIn(strings.NewReader(`{"operation":"routes.update","request":{"id":"r1","enabled":false},"before":{"id":"r1","description":"route","enabled":true,"metric":10,"groups":["g1"]},"intended_after":{"id":"r1","description":"route","enabled":false,"metric":10,"groups":["g1"]}}`))
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), `"code":"impact.route_change"`) || !strings.Contains(stdout.String(), `"severity":"blocking"`) {
		t.Fatalf("route impact acknowledgement finding missing: %s", stdout.String())
	}
}

func TestStageCreatePeerChangeRequiresAcknowledgement(t *testing.T) {
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
	root.SetIn(strings.NewReader(`{"operation":"peers.update","request":{"id":"p1","approval_required":true},"before":{"id":"p1","name":"peer","approval_required":false,"connected":true},"intended_after":{"id":"p1","name":"peer","approval_required":true,"connected":true}}`))
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), `"code":"impact.peer_change"`) || !strings.Contains(stdout.String(), `"severity":"blocking"`) {
		t.Fatalf("peer impact acknowledgement finding missing: %s", stdout.String())
	}
}

func TestStageCreateNetworkChangeRequiresAcknowledgement(t *testing.T) {
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
	root.SetIn(strings.NewReader(`{"operation":"networks.update","request":{"id":"n1","name":"office"},"before":{"id":"n1","name":"office","policies":["p1"],"resources":["r1"],"routers":["rt1"]},"intended_after":{"id":"n1","name":"office","policies":["p2"],"resources":["r1"],"routers":["rt1"]}}`))
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), `"code":"impact.network_change"`) || !strings.Contains(stdout.String(), `"severity":"blocking"`) {
		t.Fatalf("network impact acknowledgement finding missing: %s", stdout.String())
	}
}

func TestStageCreatePolicyDeleteRequiresAcknowledgement(t *testing.T) {
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
	root.SetIn(strings.NewReader(`{"operation":"policies.delete","request":{"id":"p1"},"before":{"id":"p1","name":"policy","rules":[]},"intended_after":{}}`))
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), `"code":"impact.policy_delete"`) || !strings.Contains(stdout.String(), `"severity":"blocking"`) {
		t.Fatalf("policy delete acknowledgement finding missing: %s", stdout.String())
	}
}

func TestStageCreateGroupDeleteRequiresAcknowledgement(t *testing.T) {
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
	root.SetIn(strings.NewReader(`{"operation":"groups.delete","request":{"id":"g1"},"before":{"id":"g1","name":"group","peers_count":2,"resources_count":1},"intended_after":{}}`))
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), `"code":"impact.group_delete"`) || !strings.Contains(stdout.String(), `"severity":"blocking"`) {
		t.Fatalf("group delete acknowledgement finding missing: %s", stdout.String())
	}
}

func TestStageCreateRouteDeleteRequiresAcknowledgement(t *testing.T) {
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
	root.SetIn(strings.NewReader(`{"operation":"routes.delete","request":{"id":"r1"},"before":{"id":"r1","description":"route","enabled":true,"network":"10.0.0.0/24"},"intended_after":{}}`))
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), `"code":"impact.route_delete"`) || !strings.Contains(stdout.String(), `"severity":"blocking"`) {
		t.Fatalf("route delete acknowledgement finding missing: %s", stdout.String())
	}
}

func TestStageCreateNetworkDeleteRequiresAcknowledgement(t *testing.T) {
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
	root.SetIn(strings.NewReader(`{"operation":"networks.delete","request":{"id":"n1"},"before":{"id":"n1","name":"office","policies":["p1"],"resources":["r1"],"routers":["rt1"]},"intended_after":{}}`))
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), `"code":"impact.network_delete"`) || !strings.Contains(stdout.String(), `"severity":"blocking"`) {
		t.Fatalf("network delete acknowledgement finding missing: %s", stdout.String())
	}
}

func TestStageCreatePeerDeleteRequiresAcknowledgement(t *testing.T) {
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
	root.SetIn(strings.NewReader(`{"operation":"peers.delete","request":{"id":"p1"},"before":{"id":"p1","name":"peer","connected":true,"approval_required":false},"intended_after":{}}`))
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), `"code":"impact.peer_delete"`) || !strings.Contains(stdout.String(), `"severity":"blocking"`) {
		t.Fatalf("peer delete acknowledgement finding missing: %s", stdout.String())
	}
}

func TestStageCreateNetworkResourceDeleteRequiresAcknowledgement(t *testing.T) {
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
	root.SetIn(strings.NewReader(`{"operation":"networks.resources.delete","request":{"network_id":"n1","id":"r1"},"before":{"id":"r1","name":"db","address":"10.0.0.0/24","enabled":true},"intended_after":{}}`))
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), `"code":"impact.network_resource_delete"`) || !strings.Contains(stdout.String(), `"severity":"blocking"`) {
		t.Fatalf("network resource delete acknowledgement finding missing: %s", stdout.String())
	}
}
