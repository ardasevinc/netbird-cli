package analysis

import (
	"context"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/netbird"
)

func TestAccessSeparatesSymmetricAdjacencyFromOneWayFlows(t *testing.T) {
	policyID, ruleID := "policy-1", "rule-1"
	peers := []netbird.Peer{
		{ID: "source", Name: "source", Groups: []netbird.PeerGroup{{ID: "sources"}}},
		{ID: "destination", Name: "destination", Groups: []netbird.PeerGroup{{ID: "destinations"}}},
	}
	policies := []netbird.Policy{{
		ID: &policyID, Name: "one-way", Enabled: true,
		Rules: []netbird.PolicyRule{{
			ID: &ruleID, Name: "https", Enabled: true, Action: "accept", Protocol: "tcp",
			Sources: []netbird.PolicyGroup{{ID: "sources"}}, Destinations: []netbird.PolicyGroup{{ID: "destinations"}},
			Ports: []string{"443"},
		}},
	}}

	sourceReport, err := Access(context.Background(), fakeReader{
		source: peers[0], peers: peers, reachable: []netbird.Peer{{ID: "destination"}}, policies: policies,
	}, "source")
	if err != nil {
		t.Fatal(err)
	}
	sourceRelation := requireRelation(t, sourceReport, "destination")
	if !sourceRelation.NetworkMapAdjacent || !sourceRelation.HasOutboundFlows || sourceRelation.HasInboundFlows {
		t.Fatalf("unexpected source relation: %+v", sourceRelation)
	}
	if len(sourceRelation.OutboundFlows) != 1 || sourceRelation.OutboundFlows[0].Ports[0] != "443" {
		t.Fatalf("unexpected source flows: %+v", sourceRelation.OutboundFlows)
	}

	destinationReport, err := Access(context.Background(), fakeReader{
		source: peers[1], peers: peers, reachable: []netbird.Peer{{ID: "source"}}, policies: policies,
	}, "destination")
	if err != nil {
		t.Fatal(err)
	}
	destinationRelation := requireRelation(t, destinationReport, "source")
	if !destinationRelation.NetworkMapAdjacent || destinationRelation.HasOutboundFlows || !destinationRelation.HasInboundFlows {
		t.Fatalf("unexpected destination relation: %+v", destinationRelation)
	}
	if len(destinationReport.MapOnlyPeerIDs) != 0 {
		t.Fatalf("one-way inbound adjacency was treated as map-only: %+v", destinationReport.MapOnlyPeerIDs)
	}
}

func TestAccessKeepsFlowsCanonicalWhenDirectionsDifferByProtocol(t *testing.T) {
	peerA := netbird.Peer{ID: "a", Groups: []netbird.PeerGroup{{ID: "a-group"}}}
	peerB := netbird.Peer{ID: "b", Groups: []netbird.PeerGroup{{ID: "b-group"}}}
	policies := []netbird.Policy{
		policyWithRule("out", "out", "tcp", false, []string{"443"}, "a-group", "b-group"),
		policyWithRule("in", "in", "tcp", false, []string{"22"}, "b-group", "a-group"),
		policyWithRule("icmp", "icmp", "icmp", true, nil, "a-group", "b-group"),
	}
	report, err := Access(context.Background(), fakeReader{
		source: peerA, peers: []netbird.Peer{peerA, peerB}, reachable: []netbird.Peer{{ID: "b"}}, policies: policies,
	}, "a")
	if err != nil {
		t.Fatal(err)
	}
	relation := requireRelation(t, report, "b")
	if len(relation.OutboundFlows) != 2 || len(relation.InboundFlows) != 2 {
		t.Fatalf("flows were collapsed into a peer-level classification: %+v", relation)
	}
	if len(relation.Directions) != 2 || relation.Directions[0] != DirectionOutbound || relation.Directions[1] != DirectionInbound {
		t.Fatalf("unexpected derived directions: %+v", relation.Directions)
	}
}

func TestAccessDeduplicatesBidirectionalAllGroupFlows(t *testing.T) {
	peerA := netbird.Peer{ID: "a", Groups: []netbird.PeerGroup{{ID: "all"}}}
	peerB := netbird.Peer{ID: "b", Groups: []netbird.PeerGroup{{ID: "all"}}}
	policy := policyWithRule("all", "all", "icmp", true, nil, "all", "all")
	policy.Rules[0].Action = "ACCEPT"
	report, err := Access(context.Background(), fakeReader{
		source: peerA, peers: []netbird.Peer{peerA, peerB}, reachable: []netbird.Peer{{ID: "b"}}, policies: []netbird.Policy{policy},
	}, "a")
	if err != nil {
		t.Fatal(err)
	}
	relation := requireRelation(t, report, "b")
	if len(relation.OutboundFlows) != 1 || len(relation.InboundFlows) != 1 {
		t.Fatalf("bidirectional all-group rule was duplicated: %+v", relation)
	}
	if relation.OutboundFlows[0].Action != "accept" || relation.InboundFlows[0].Action != "accept" {
		t.Fatalf("accept action was not normalized to the schema contract: %+v", relation)
	}
}

func TestAccessSupportsDirectPeerResourcesAndScopedPostureMetadata(t *testing.T) {
	policyID, ruleID := "direct", "direct-rule"
	peerA := netbird.Peer{ID: "a"}
	peerB := netbird.Peer{ID: "b"}
	report, err := Access(context.Background(), fakeReader{
		source:    peerA,
		peers:     []netbird.Peer{peerA, peerB},
		reachable: []netbird.Peer{{ID: "b"}},
		policies: []netbird.Policy{{
			ID: &policyID, Name: "direct", Enabled: true, SourcePostureChecks: []string{"linux"},
			Rules: []netbird.PolicyRule{{
				ID: &ruleID, Name: "direct", Enabled: true, Action: "accept", Protocol: "tcp",
				SourceResource:      &netbird.PolicyResource{ID: "a", Type: "peer"},
				DestinationResource: &netbird.PolicyResource{ID: "b", Type: "peer"},
				PortRanges:          []netbird.PolicyPortRange{{Start: 8000, End: 8010}},
			}},
		}},
	}, "a")
	if err != nil {
		t.Fatal(err)
	}
	flow := requireRelation(t, report, "b").OutboundFlows[0]
	if flow.SourceMatch != "peer_resource" || flow.DestinationMatch != "peer_resource" || len(flow.PortRanges) != 1 {
		t.Fatalf("direct peer flow lost provenance: %+v", flow)
	}
	if len(flow.PolicySourcePostureChecks) != 1 || flow.PolicySourcePostureChecks[0] != "linux" {
		t.Fatalf("posture scope was lost: %+v", flow.PolicySourcePostureChecks)
	}
	if report.Completeness.PolicyProjection.State != "complete" {
		t.Fatalf("configured flow projection should remain complete: %+v", report.Completeness.PolicyProjection)
	}
}

func TestAccessReportsMapOnlyMissingAndUnsupportedActionWithoutDenyInference(t *testing.T) {
	peerA := netbird.Peer{ID: "a", Groups: []netbird.PeerGroup{{ID: "a-group"}}}
	peerB := netbird.Peer{ID: "b", Groups: []netbird.PeerGroup{{ID: "b-group"}}}
	peerC := netbird.Peer{ID: "c"}
	peerD := netbird.Peer{ID: "d", Groups: []netbird.PeerGroup{{ID: "d-group"}}}
	accept := policyWithRule("allow", "allow", "tcp", false, nil, "a-group", "b-group")
	unsupported := policyWithRule("reject", "reject", "tcp", false, nil, "a-group", "d-group")
	unsupported.Rules[0].Action = "reject"
	disabledPolicy := policyWithRule("disabled-policy", "disabled-policy", "tcp", false, nil, "a-group", "d-group")
	disabledPolicy.Enabled = false
	disabledRule := policyWithRule("disabled-rule", "disabled-rule", "tcp", false, nil, "a-group", "d-group")
	disabledRule.Rules[0].Enabled = false
	report, err := Access(context.Background(), fakeReader{
		source: peerA, peers: []netbird.Peer{peerA, peerB, peerC, peerD},
		reachable: []netbird.Peer{{ID: "c"}}, policies: []netbird.Policy{accept, unsupported, disabledPolicy, disabledRule},
	}, "a")
	if err != nil {
		t.Fatal(err)
	}
	if len(report.MapOnlyPeerIDs) != 1 || report.MapOnlyPeerIDs[0] != "c" {
		t.Fatalf("unexpected map-only peers: %+v", report.MapOnlyPeerIDs)
	}
	if len(report.ConfiguredPeerMissingFromNetworkMapIDs) != 1 || report.ConfiguredPeerMissingFromNetworkMapIDs[0] != "b" {
		t.Fatalf("unexpected configured peers missing from map: %+v", report.ConfiguredPeerMissingFromNetworkMapIDs)
	}
	if len(report.ProjectionFindings) != 1 || report.ProjectionFindings[0].Code != "unsupported_action" {
		t.Fatalf("unsupported action was not surfaced: %+v", report.ProjectionFindings)
	}
	if report.Completeness.PolicyProjection.State != "partial" || len(report.Observations) != 0 || report.Completeness.Observations.State != "unknown" {
		t.Fatalf("unexpected completeness: %+v", report.Completeness)
	}
	if _, ok := relationByID(report, "d"); ok {
		t.Fatal("unsupported action was interpreted as a deny relation")
	}
}

func policyWithRule(policyID, ruleID, protocol string, bidirectional bool, ports []string, sourceGroup, destinationGroup string) netbird.Policy {
	return netbird.Policy{
		ID: &policyID, Name: policyID, Enabled: true,
		Rules: []netbird.PolicyRule{{
			ID: &ruleID, Name: ruleID, Enabled: true, Action: "accept", Protocol: protocol, Bidirectional: bidirectional,
			Sources: []netbird.PolicyGroup{{ID: sourceGroup}}, Destinations: []netbird.PolicyGroup{{ID: destinationGroup}}, Ports: ports,
		}},
	}
}

func requireRelation(t *testing.T, report AccessReport, peerID string) AccessRelation {
	t.Helper()
	relation, ok := relationByID(report, peerID)
	if !ok {
		t.Fatalf("missing relation for %s: %+v", peerID, report.Relations)
	}
	return relation
}

func relationByID(report AccessReport, peerID string) (AccessRelation, bool) {
	for _, relation := range report.Relations {
		if relation.Peer.ID == peerID {
			return relation, true
		}
	}
	return AccessRelation{}, false
}
