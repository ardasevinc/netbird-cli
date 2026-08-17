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

func TestStageCreateAccountUpdateRequiresAcknowledgement(t *testing.T) {
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
	root.SetIn(strings.NewReader(`{"operation":"accounts.update","request":{"id":"account-1","settings":{"peer_login_expiration_enabled":false}},"before":{"id":"account-1","settings":{"peer_login_expiration_enabled":true}},"intended_after":{"id":"account-1","settings":{"peer_login_expiration_enabled":false}}}`))
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), `"code":"impact.account_change"`) || !strings.Contains(stdout.String(), `"severity":"blocking"`) {
		t.Fatalf("account impact acknowledgement finding missing: %s", stdout.String())
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

func TestStageCreatePolicyCreateRequiresAcknowledgement(t *testing.T) {
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
	root.SetIn(strings.NewReader(`{"operation":"policies.create","request":{"name":"allow-office","enabled":true,"rules":[{"action":"accept"}]},"before":[],"intended_after":{"name":"allow-office","enabled":true,"rules":[{"action":"accept"}]}}`))
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), `"code":"impact.policy_create"`) || !strings.Contains(stdout.String(), `"severity":"blocking"`) {
		t.Fatalf("policy create acknowledgement finding missing: %s", stdout.String())
	}
}

func TestStageCreateDNSZoneCreateRequiresAcknowledgement(t *testing.T) {
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
	root.SetIn(strings.NewReader(`{"operation":"dns.zones.create","request":{"name":"office","domain":"office.internal","enabled":true},"before":[],"intended_after":{"name":"office","domain":"office.internal","enabled":true}}`))
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), `"code":"impact.dns_zone_create"`) || !strings.Contains(stdout.String(), `"severity":"blocking"`) {
		t.Fatalf("dns zone create acknowledgement finding missing: %s", stdout.String())
	}
}

func TestStageCreateDNSZoneDeleteRequiresAcknowledgement(t *testing.T) {
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
	root.SetIn(strings.NewReader(`{"operation":"dns.zones.delete","request":{"id":"zone-1"},"before":{"id":"zone-1","domain":"office.internal","distribution_groups":["g1"]},"intended_after":{}}`))
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), `"code":"impact.dns_zone_delete"`) || !strings.Contains(stdout.String(), `"severity":"blocking"`) {
		t.Fatalf("dns zone delete acknowledgement finding missing: %s", stdout.String())
	}
}

func TestStageCreateDNSZoneUpdateRequiresAcknowledgement(t *testing.T) {
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
	root.SetIn(strings.NewReader(`{"operation":"dns.zones.update","request":{"id":"zone-1","domain":"corp.internal","enabled":true},"before":{"id":"zone-1","domain":"office.internal","enabled":true},"intended_after":{"id":"zone-1","domain":"corp.internal","enabled":true}}`))
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), `"code":"impact.dns_zone_change"`) || !strings.Contains(stdout.String(), `"severity":"blocking"`) {
		t.Fatalf("dns zone update acknowledgement finding missing: %s", stdout.String())
	}
}

func TestStageCreateDNSRecordCreateRequiresAcknowledgement(t *testing.T) {
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
	root.SetIn(strings.NewReader(`{"operation":"dns.records.create","request":{"zone_id":"zone-1","name":"db","type":"A","content":"10.0.0.5","ttl":60},"before":[],"intended_after":{"zone_id":"zone-1","name":"db","type":"A","content":"10.0.0.5","ttl":60}}`))
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), `"code":"impact.dns_record_create"`) || !strings.Contains(stdout.String(), `"severity":"blocking"`) {
		t.Fatalf("dns record create acknowledgement finding missing: %s", stdout.String())
	}
}

func TestStageCreateDNSRecordUpdateRequiresAcknowledgement(t *testing.T) {
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
	root.SetIn(strings.NewReader(`{"operation":"dns.records.update","request":{"zone_id":"zone-1","id":"record-1","content":"10.0.0.6"},"before":{"id":"record-1","name":"db","type":"A","content":"10.0.0.5","ttl":60},"intended_after":{"id":"record-1","name":"db","type":"A","content":"10.0.0.6","ttl":60}}`))
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), `"code":"impact.dns_record_change"`) || !strings.Contains(stdout.String(), `"severity":"blocking"`) {
		t.Fatalf("dns record update acknowledgement finding missing: %s", stdout.String())
	}
}

func TestStageCreateDNSRecordDeleteRequiresAcknowledgement(t *testing.T) {
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
	root.SetIn(strings.NewReader(`{"operation":"dns.records.delete","request":{"zone_id":"zone-1","id":"record-1"},"before":{"id":"record-1","name":"db","type":"A","content":"10.0.0.5","ttl":60},"intended_after":{}}`))
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), `"code":"impact.dns_record_delete"`) || !strings.Contains(stdout.String(), `"severity":"blocking"`) {
		t.Fatalf("dns record delete acknowledgement finding missing: %s", stdout.String())
	}
}

func TestStageCreateNameserverGroupCreateRequiresAcknowledgement(t *testing.T) {
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
	root.SetIn(strings.NewReader(`{"operation":"dns.nameservers.create","request":{"name":"office","domains":["office.internal"],"enabled":true,"nameservers":[{"ip":"10.0.0.53","ns_type":"udp","port":53}]},"before":[],"intended_after":{"name":"office","domains":["office.internal"],"enabled":true,"nameservers":[{"ip":"10.0.0.53","ns_type":"udp","port":53}]}}`))
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), `"code":"impact.dns_nameserver_create"`) || !strings.Contains(stdout.String(), `"severity":"blocking"`) {
		t.Fatalf("nameserver create acknowledgement finding missing: %s", stdout.String())
	}
}

func TestStageCreateNameserverGroupUpdateRequiresAcknowledgement(t *testing.T) {
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
	root.SetIn(strings.NewReader(`{"operation":"dns.nameservers.update","request":{"id":"ns-1","description":"new"},"before":{"id":"ns-1","description":"old","enabled":true},"intended_after":{"id":"ns-1","description":"new","enabled":true}}`))
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), `"code":"impact.dns_nameserver_change"`) || !strings.Contains(stdout.String(), `"severity":"blocking"`) {
		t.Fatalf("nameserver update acknowledgement finding missing: %s", stdout.String())
	}
}

func TestStageCreateNameserverGroupDeleteRequiresAcknowledgement(t *testing.T) {
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
	root.SetIn(strings.NewReader(`{"operation":"dns.nameservers.delete","request":{"id":"ns-1"},"before":{"id":"ns-1","domains":["office.internal"],"enabled":true},"intended_after":{}}`))
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), `"code":"impact.dns_nameserver_delete"`) || !strings.Contains(stdout.String(), `"severity":"blocking"`) {
		t.Fatalf("nameserver delete acknowledgement finding missing: %s", stdout.String())
	}
}

func TestStageCreateDNSSettingsUpdateRequiresAcknowledgement(t *testing.T) {
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
	root.SetIn(strings.NewReader(`{"operation":"dns.settings.update","request":{"disabled_management_groups":["g1","g2"]},"before":{"disabled_management_groups":["g1"]},"intended_after":{"disabled_management_groups":["g1","g2"]}}`))
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), `"code":"impact.dns_settings_change"`) || !strings.Contains(stdout.String(), `"severity":"blocking"`) {
		t.Fatalf("dns settings update acknowledgement finding missing: %s", stdout.String())
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

func TestStageCreateGroupCreateRequiresAcknowledgement(t *testing.T) {
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
	root.SetIn(strings.NewReader(`{"operation":"groups.create","request":{"name":"engineering"},"before":[],"intended_after":{"name":"engineering"}}`))
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), `"code":"impact.group_create"`) || !strings.Contains(stdout.String(), `"severity":"blocking"`) {
		t.Fatalf("group create acknowledgement finding missing: %s", stdout.String())
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

func TestStageCreateRouteCreateRequiresAcknowledgement(t *testing.T) {
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
	root.SetIn(strings.NewReader(`{"operation":"routes.create","request":{"description":"private subnet","enabled":true,"network":"10.0.0.0/24","groups":["g1"]},"before":[],"intended_after":{"description":"private subnet","enabled":true,"network":"10.0.0.0/24","groups":["g1"]}}`))
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), `"code":"impact.route_create"`) || !strings.Contains(stdout.String(), `"severity":"blocking"`) {
		t.Fatalf("route create acknowledgement finding missing: %s", stdout.String())
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

func TestStageCreateNetworkCreateRequiresAcknowledgement(t *testing.T) {
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
	root.SetIn(strings.NewReader(`{"operation":"networks.create","request":{"name":"office","description":"primary"},"before":[],"intended_after":{"name":"office","description":"primary"}}`))
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), `"code":"impact.network_create"`) || !strings.Contains(stdout.String(), `"severity":"blocking"`) {
		t.Fatalf("network create acknowledgement finding missing: %s", stdout.String())
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

func TestStageCreateNetworkResourceChangeRequiresAcknowledgement(t *testing.T) {
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
	root.SetIn(strings.NewReader(`{"operation":"networks.resources.update","request":{"network_id":"n1","id":"r1","address":"10.0.1.0/24"},"before":{"id":"r1","name":"db","address":"10.0.0.0/24","enabled":true,"groups":["g1"]},"intended_after":{"id":"r1","name":"db","address":"10.0.1.0/24","enabled":true,"groups":["g1"]}}`))
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), `"code":"impact.network_resource_change"`) || !strings.Contains(stdout.String(), `"severity":"blocking"`) {
		t.Fatalf("network resource impact acknowledgement finding missing: %s", stdout.String())
	}
}

func TestStageCreateNetworkResourceCreateRequiresAcknowledgement(t *testing.T) {
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
	root.SetIn(strings.NewReader(`{"operation":"networks.resources.create","request":{"network_id":"n1","name":"db","address":"10.0.0.0/24","enabled":true,"groups":["g1"]},"before":[],"intended_after":{"name":"db","address":"10.0.0.0/24","enabled":true,"groups":["g1"]}}`))
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), `"code":"impact.network_resource_create"`) || !strings.Contains(stdout.String(), `"severity":"blocking"`) {
		t.Fatalf("network resource create acknowledgement finding missing: %s", stdout.String())
	}
}

func TestStageCreateNetworkRouterCreateRequiresAcknowledgement(t *testing.T) {
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
	root.SetIn(strings.NewReader(`{"operation":"networks.routers.create","request":{"network_id":"n1","enabled":true,"masquerade":true,"metric":10,"peer":"p1"},"before":[],"intended_after":{"enabled":true,"masquerade":true,"metric":10,"peer":"p1"}}`))
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), `"code":"impact.network_router_create"`) || !strings.Contains(stdout.String(), `"severity":"blocking"`) {
		t.Fatalf("network router create acknowledgement finding missing: %s", stdout.String())
	}
}

func TestStageCreateNetworkRouterDeleteRequiresAcknowledgement(t *testing.T) {
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
	root.SetIn(strings.NewReader(`{"operation":"networks.routers.delete","request":{"network_id":"n1","id":"rt1"},"before":{"id":"rt1","enabled":true,"masquerade":true,"metric":10,"peer":"p1"},"intended_after":{}}`))
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), `"code":"impact.network_router_delete"`) || !strings.Contains(stdout.String(), `"severity":"blocking"`) {
		t.Fatalf("network router delete acknowledgement finding missing: %s", stdout.String())
	}
}

func TestStageCreateNetworkRouterChangeRequiresAcknowledgement(t *testing.T) {
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
	root.SetIn(strings.NewReader(`{"operation":"networks.routers.update","request":{"network_id":"n1","id":"rt1","enabled":false},"before":{"id":"rt1","enabled":true,"masquerade":true,"metric":10,"peer":"p1","peer_groups":["g1"]},"intended_after":{"id":"rt1","enabled":false,"masquerade":true,"metric":10,"peer":"p1","peer_groups":["g1"]}}`))
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), `"code":"impact.network_router_change"`) || !strings.Contains(stdout.String(), `"severity":"blocking"`) {
		t.Fatalf("network router impact acknowledgement finding missing: %s", stdout.String())
	}
}
