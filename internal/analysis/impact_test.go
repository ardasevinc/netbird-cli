package analysis

import "testing"

func TestGroupUpdateImpactTreatsNameChangeAsReachabilityNeutral(t *testing.T) {
	report, err := GroupUpdateImpact([]byte(`{"id":"g1","name":"old","peers_count":2}`), []byte(`{"id":"g1","name":"new","peers_count":2}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "metadata_only" || report.Reachability != "unchanged" || report.Confidence != "high" || report.Completeness["state"] != "complete" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestGroupUpdateImpactRefusesToOverclaimUnsupportedChanges(t *testing.T) {
	report, err := GroupUpdateImpact([]byte(`{"id":"g1","name":"old"}`), []byte(`{"id":"g1","name":"old","members":["p1"]}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "unknown" || report.Reachability != "unknown" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestPolicyUpdateImpactMarksRuleChangesAsPotentialReachabilityChange(t *testing.T) {
	report, err := PolicyUpdateImpact([]byte(`{"name":"p","rules":[{"id":"r1","action":"accept"}]}`), []byte(`{"name":"p","rules":[{"id":"r1","action":"drop"}]}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "policy_rule_change" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestPolicyDeleteImpactIsConservative(t *testing.T) {
	report, err := PolicyDeleteImpact([]byte(`{"id":"p1","name":"policy","rules":[]}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "policy_delete" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestRouteUpdateImpactMarksDescriptionOnlyChangeAsNeutral(t *testing.T) {
	report, err := RouteUpdateImpact(
		[]byte(`{"id":"r1","description":"old","enabled":true,"metric":10,"groups":["g1"]}`),
		[]byte(`{"id":"r1","description":"new","enabled":true,"metric":10,"groups":["g1"]}`),
	)
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "metadata_only" || report.Reachability != "unchanged" || report.Confidence != "high" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestRouteDeleteImpactIsConservative(t *testing.T) {
	report, err := RouteDeleteImpact([]byte(`{"id":"r1","description":"route","enabled":true,"network":"10.0.0.0/24"}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "route_delete" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestRouteUpdateImpactBlocksRoutingChanges(t *testing.T) {
	report, err := RouteUpdateImpact(
		[]byte(`{"id":"r1","description":"route","enabled":true,"metric":10,"groups":["g1"]}`),
		[]byte(`{"id":"r1","description":"route","enabled":false,"metric":10,"groups":["g1"]}`),
	)
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "route_change" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestPeerUpdateImpactMarksNameOnlyChangeAsNeutral(t *testing.T) {
	report, err := PeerUpdateImpact(
		[]byte(`{"id":"p1","name":"old","approval_required":false,"connected":true}`),
		[]byte(`{"id":"p1","name":"new","approval_required":false,"connected":true}`),
	)
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "metadata_only" || report.Reachability != "unchanged" || report.Confidence != "high" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestGroupDeleteImpactIsConservative(t *testing.T) {
	report, err := GroupDeleteImpact([]byte(`{"id":"g1","name":"group","peers_count":2,"resources_count":1}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "group_delete" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestPeerUpdateImpactBlocksAccessChanges(t *testing.T) {
	report, err := PeerUpdateImpact(
		[]byte(`{"id":"p1","name":"peer","approval_required":false,"connected":true}`),
		[]byte(`{"id":"p1","name":"peer","approval_required":true,"connected":true}`),
	)
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "peer_change" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestPeerDeleteImpactIsConservative(t *testing.T) {
	report, err := PeerDeleteImpact([]byte(`{"id":"p1","name":"peer","connected":true,"approval_required":false}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "peer_delete" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestNetworkUpdateImpactMarksMetadataChangesAsNeutral(t *testing.T) {
	report, err := NetworkUpdateImpact(
		[]byte(`{"id":"n1","name":"old","description":"office","policies":["p1"],"resources":["r1"],"routers":["rt1"]}`),
		[]byte(`{"id":"n1","name":"new","description":"office","policies":["p1"],"resources":["r1"],"routers":["rt1"]}`),
	)
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "metadata_only" || report.Reachability != "unchanged" || report.Confidence != "high" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestNetworkUpdateImpactBlocksTopologyChanges(t *testing.T) {
	report, err := NetworkUpdateImpact(
		[]byte(`{"id":"n1","name":"office","policies":["p1"],"resources":["r1"],"routers":["rt1"]}`),
		[]byte(`{"id":"n1","name":"office","policies":["p2"],"resources":["r1"],"routers":["rt1"]}`),
	)
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "network_change" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestNetworkDeleteImpactIsConservative(t *testing.T) {
	report, err := NetworkDeleteImpact([]byte(`{"id":"n1","name":"office","policies":["p1"],"resources":["r1"],"routers":["rt1"]}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "network_delete" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestNetworkResourceDeleteImpactIsConservative(t *testing.T) {
	report, err := NetworkResourceDeleteImpact([]byte(`{"id":"r1","name":"db","address":"10.0.0.0/24","enabled":true}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "network_resource_delete" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestNetworkResourceUpdateImpactTreatsMetadataAsNeutral(t *testing.T) {
	report, err := NetworkResourceUpdateImpact(
		[]byte(`{"id":"r1","name":"old","description":"db","address":"10.0.0.0/24","enabled":true,"groups":["g1"]}`),
		[]byte(`{"id":"r1","name":"new","description":"db","address":"10.0.0.0/24","enabled":true,"groups":["g1"]}`),
	)
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "metadata_only" || report.Reachability != "unchanged" || report.Confidence != "high" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestNetworkResourceUpdateImpactBlocksTopologyChanges(t *testing.T) {
	report, err := NetworkResourceUpdateImpact(
		[]byte(`{"id":"r1","name":"db","address":"10.0.0.0/24","enabled":true,"groups":["g1"]}`),
		[]byte(`{"id":"r1","name":"db","address":"10.0.1.0/24","enabled":true,"groups":["g1"]}`),
	)
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "network_resource_change" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestNetworkRouterDeleteImpactIsConservative(t *testing.T) {
	report, err := NetworkRouterDeleteImpact([]byte(`{"id":"rt1","enabled":true,"masquerade":true,"metric":10,"peer":"p1"}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "network_router_delete" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}
