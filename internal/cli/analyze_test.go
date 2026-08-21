package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/version"
)

func TestReachabilityAnalysisJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		switch r.URL.Path {
		case "/api/peers/p1":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "p1", "name": "source", "groups": []map[string]any{{"id": "source-group"}}})
		case "/api/peers/p1/accessible-peers":
			_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "p2", "name": "target", "ip": "10.0.0.2", "connected": true, "groups": []map[string]any{}}})
		case "/api/peers":
			_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "p2", "name": "target", "ip": "10.0.0.2", "connected": true, "groups": []map[string]any{{"id": "target-group"}}}})
		case "/api/policies":
			_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "policy-1", "name": "allow", "enabled": true, "rules": []map[string]any{{"id": "rule-1", "name": "rule", "action": "accept", "protocol": "all", "enabled": true, "sources": []map[string]any{{"id": "source-group"}}, "destinations": []map[string]any{{"id": "target-group"}}}}}})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	t.Setenv("NB_ANALYZE_TOKEN", "token")
	temp := t.TempDir()
	configPath := filepath.Join(temp, "config.toml")
	if err := os.WriteFile(configPath, []byte("[profiles.default]\nurl = \""+server.URL+"\"\ncredential_ref = \"env:NB_ANALYZE_TOKEN\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	state := &commandState{json: true, configPath: configPath, profileName: "default", statePath: filepath.Join(temp, "ledger.db")}
	var stdout, stderr bytes.Buffer
	root := newRoot(state, &stdout, &stderr, version.Current())
	root.SetArgs([]string{"analyze", "reachability", "p1"})
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("reachability analysis: %v stderr=%s", err, stderr.String())
	}
	var response struct {
		Schema string `json:"schema"`
		Data   struct {
			Summary struct {
				ReachablePeerCount            int `json:"reachable_peer_count"`
				PolicyEvidenceCount           int `json:"policy_evidence_count"`
				UnexplainedReachablePeerCount int `json:"unexplained_reachable_peer_count"`
			} `json:"summary"`
		} `json:"data"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v: %s", err, stdout.String())
	}
	if response.Schema != "nb/v1/reachability-analysis-result" || response.Data.Summary.ReachablePeerCount != 1 || response.Data.Summary.PolicyEvidenceCount != 1 || response.Data.Summary.UnexplainedReachablePeerCount != 0 {
		t.Fatalf("unexpected response: %s", stdout.String())
	}
	assertJSONKeys(t, stdout.Bytes(), "data", "ok", "operation", "schema")
	var envelope struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &envelope); err != nil {
		t.Fatal(err)
	}
	assertJSONKeys(t, envelope.Data, "completeness", "policy_evidence", "reachable_peers", "source_peer", "summary", "unexplained_reachable_peer_ids")
}

func TestAccessAnalysisClassifiesOneWayDestinationAsInboundOnly(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		switch r.URL.Path {
		case "/api/peers/p2":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "p2", "name": "destination", "groups": []map[string]any{{"id": "destinations"}}})
		case "/api/peers/p2/accessible-peers":
			_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "p1", "name": "source", "ip": "10.0.0.1", "connected": true}})
		case "/api/peers":
			_ = json.NewEncoder(w).Encode([]map[string]any{
				{"id": "p1", "name": "source", "ip": "10.0.0.1", "connected": true, "groups": []map[string]any{{"id": "sources"}}},
				{"id": "p2", "name": "destination", "ip": "10.0.0.2", "connected": true, "groups": []map[string]any{{"id": "destinations"}}},
			})
		case "/api/policies":
			_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "policy-1", "name": "one-way", "enabled": true, "rules": []map[string]any{{"id": "rule-1", "name": "https", "action": "accept", "protocol": "tcp", "ports": []string{"443"}, "enabled": true, "bidirectional": false, "sources": []map[string]any{{"id": "sources"}}, "destinations": []map[string]any{{"id": "destinations"}}}}}})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	t.Setenv("NB_ACCESS_ANALYZE_TOKEN", "token")
	temp := t.TempDir()
	configPath := filepath.Join(temp, "config.toml")
	if err := os.WriteFile(configPath, []byte("[profiles.default]\nurl = \""+server.URL+"\"\ncredential_ref = \"env:NB_ACCESS_ANALYZE_TOKEN\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	state := &commandState{json: true, configPath: configPath, profileName: "default", statePath: filepath.Join(temp, "ledger.db")}
	var stdout, stderr bytes.Buffer
	root := newRoot(state, &stdout, &stderr, version.Current())
	root.SetArgs([]string{"analyze", "access", "p2"})
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("access analysis: %v stderr=%s", err, stderr.String())
	}
	var response struct {
		Schema    string `json:"schema"`
		Operation string `json:"operation"`
		Data      struct {
			Relations []struct {
				HasOutboundFlows bool `json:"has_outbound_flows"`
				HasInboundFlows  bool `json:"has_inbound_flows"`
				InboundFlows     []struct {
					Protocol string   `json:"protocol"`
					Ports    []string `json:"ports"`
				} `json:"inbound_flows"`
			} `json:"relations"`
			Observations []json.RawMessage `json:"observations"`
		} `json:"data"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v: %s", err, stdout.String())
	}
	if response.Schema != "nb/v1/access-analysis-result" || response.Operation != "analyze.access" || len(response.Data.Relations) != 1 {
		t.Fatalf("unexpected response: %s", stdout.String())
	}
	relation := response.Data.Relations[0]
	if relation.HasOutboundFlows || !relation.HasInboundFlows || len(relation.InboundFlows) != 1 || relation.InboundFlows[0].Protocol != "tcp" || relation.InboundFlows[0].Ports[0] != "443" {
		t.Fatalf("one-way destination was not inbound-only: %s", stdout.String())
	}
	if len(response.Data.Observations) != 0 {
		t.Fatalf("observations must be scoped records, initially empty: %s", stdout.String())
	}
}
