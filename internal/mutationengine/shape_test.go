package mutationengine

import (
	"encoding/json"
	"testing"
)

func TestObjectContainsTreatsNullAndEmptyArrayAsEquivalent(t *testing.T) {
	actual := json.RawMessage(`{"id":"n1","resources":null,"routers":null}`)
	expected := json.RawMessage(`{"id":"n1","resources":[],"routers":[]}`)
	equal, err := objectContains(actual, expected)
	if err != nil || !equal {
		t.Fatalf("equal=%t err=%v", equal, err)
	}
}
