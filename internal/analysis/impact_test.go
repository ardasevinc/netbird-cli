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
