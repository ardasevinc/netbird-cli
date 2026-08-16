package operations

import "testing"

func TestConsequentialOperationRequiresSafetyProof(t *testing.T) {
	definition, err := Lookup("groups.update")
	if err != nil {
		t.Fatal(err)
	}
	if definition.Mutation != Consequential || !definition.RequiresPreimage || !definition.RequiresReadBack || !definition.DispatcherAdmitted {
		t.Fatalf("unexpected operation definition: %+v", definition)
	}
}

func TestUnknownOperationIsNotAdmitted(t *testing.T) {
	if _, err := Lookup("anything.delete"); err == nil {
		t.Fatal("unknown operation was admitted")
	}
}
