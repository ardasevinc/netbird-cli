package mutation

import "testing"

func TestAcknowledgementsBindToBlockingCodes(t *testing.T) {
	findings := []Finding{{Code: "impact.high", Severity: Blocking}, {Code: "note", Severity: Info}}
	if err := ValidateAcknowledgements(findings, nil, false); err == nil {
		t.Fatal("expected missing blocker acknowledgement")
	}
	if err := ValidateAcknowledgements(findings, []string{"impact.high"}, false); err != nil {
		t.Fatal(err)
	}
}

func TestAmbiguousDispatchCannotRetry(t *testing.T) {
	if RetryAllowed(Unknown) || RetryAllowed(Partial) {
		t.Fatal("uncertain dispatch must not be replayed")
	}
	if !RetryAllowed(NotDispatched) {
		t.Fatal("definitely undispatched request should be retryable")
	}
}
