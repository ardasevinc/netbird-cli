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
