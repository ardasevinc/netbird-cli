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

func TestBillingOperationsAreCloudEntitlementGated(t *testing.T) {
	for _, name := range []string{"billing.aws_marketplace.activate", "billing.aws_marketplace.enrich", "billing.checkout.create", "billing.subscription.update"} {
		definition, err := Lookup(name)
		if err != nil {
			t.Fatal(err)
		}
		if definition.Availability != "cloud_entitled" {
			t.Fatalf("%s availability = %q", name, definition.Availability)
		}
	}
}
