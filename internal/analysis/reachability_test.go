package analysis

import (
	"context"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/netbird"
)

type fakeReader struct {
	source    netbird.Peer
	reachable []netbird.Peer
	policies  []netbird.Policy
}

func (f fakeReader) GetPeer(context.Context, string) (netbird.Peer, error) {
	return f.source, nil
}

func (f fakeReader) ListAccessiblePeers(context.Context, string) ([]netbird.Peer, error) {
	return f.reachable, nil
}

func (f fakeReader) ListPolicies(context.Context) ([]netbird.Policy, error) {
	return f.policies, nil
}

func TestReachabilityUsesServerPeersAndAddsPolicyEvidence(t *testing.T) {
	policyID := "policy-1"
	ruleID := "rule-1"
	report, err := Reachability(context.Background(), fakeReader{
		source: netbird.Peer{ID: "p1", Name: "source", Groups: []netbird.PeerGroup{{ID: "source-group"}}},
		reachable: []netbird.Peer{
			{ID: "p2", Name: "target", Groups: []netbird.PeerGroup{{ID: "target-group"}}},
			{ID: "p3", Name: "unexplained"},
		},
		policies: []netbird.Policy{{
			ID:      &policyID,
			Name:    "allow-source-target",
			Enabled: true,
			Rules: []netbird.PolicyRule{{
				ID:           &ruleID,
				Name:         "allow",
				Action:       "accept",
				Protocol:     "all",
				Enabled:      true,
				Sources:      []netbird.PolicyGroup{{ID: "source-group"}},
				Destinations: []netbird.PolicyGroup{{ID: "target-group"}},
			}},
		}},
	}, "p1")
	if err != nil {
		t.Fatal(err)
	}
	if report.Summary.ReachablePeerCount != 2 || report.Summary.PolicyEvidenceCount != 1 || report.Summary.UnexplainedReachablePeerCount != 1 {
		t.Fatalf("unexpected summary: %+v", report.Summary)
	}
	if got := report.ReachablePeers[0].ID; got != "p2" {
		t.Fatalf("reachable peers were not sorted: %s", got)
	}
	if got := report.PolicyEvidence[0]; got.PolicyID != policyID || got.RuleID != ruleID || got.Direction != "source_to_destination" {
		t.Fatalf("unexpected evidence: %+v", got)
	}
	if len(report.UnexplainedReachablePeerIDs) != 1 || report.UnexplainedReachablePeerIDs[0] != "p3" {
		t.Fatalf("unexpected unexplained peers: %+v", report.UnexplainedReachablePeerIDs)
	}
}
