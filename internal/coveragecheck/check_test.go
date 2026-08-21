package coveragecheck

import (
	"fmt"
	"strings"
	"testing"
)

func TestValidateCountsImplementationAndVerificationStates(t *testing.T) {
	inventoryJSON := inventoryFixture(4)
	manifestJSON := []byte(`{
  "operations": [
    {"method":"GET","path":"/0","implementation":"implemented","verification":"contract_verified"},
    {"method":"GET","path":"/1","implementation":"classified","verification":"disposable_verified"},
    {"method":"GET","path":"/2","implementation":"discovered","verification":"live_verified"},
    {"method":"GET","path":"/3","implementation":"implemented","verification":"unverified_live"}
  ],
  "summary": {"implemented":2,"classified":1,"discovered":1,"contract_verified":1,"disposable_verified":1,"live_verified":1,"unverified_live":1}
}`)
	count, err := Validate(inventoryJSON, manifestJSON)
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Fatalf("unexpected operation count: %d", count)
	}
}

func TestValidateRejectsStaleVerificationSummary(t *testing.T) {
	manifestJSON := []byte(`{
  "operations": [{"method":"GET","path":"/0","implementation":"implemented","verification":"live_verified"}],
  "summary": {"implemented":1,"classified":0,"discovered":0,"contract_verified":1,"disposable_verified":0,"live_verified":0,"unverified_live":0}
}`)
	_, err := Validate(inventoryFixture(1), manifestJSON)
	if err == nil || !strings.Contains(err.Error(), "coverage summary does not match operation rows") {
		t.Fatalf("expected stale verification summary failure, got %v", err)
	}
}

func TestValidateRejectsUnknownAndMissingVerificationStates(t *testing.T) {
	for name, verification := range map[string]string{"unknown": `"mystery"`, "missing": `""`} {
		t.Run(name, func(t *testing.T) {
			manifestJSON := []byte(fmt.Sprintf(`{
  "operations": [{"method":"GET","path":"/0","implementation":"implemented","verification":%s}],
  "summary": {"implemented":1,"classified":0,"discovered":0,"contract_verified":0,"disposable_verified":0,"live_verified":0,"unverified_live":0}
}`, verification))
			if _, err := Validate(inventoryFixture(1), manifestJSON); err == nil {
				t.Fatal("expected invalid verification state to fail")
			}
		})
	}
}

func inventoryFixture(count int) []byte {
	var rows strings.Builder
	for i := range count {
		if i > 0 {
			rows.WriteByte(',')
		}
		fmt.Fprintf(&rows, `{"method":"GET","path":"/%d"}`, i)
	}
	return []byte(`{"operations":[` + rows.String() + `]}`)
}
