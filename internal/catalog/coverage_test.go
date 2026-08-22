package catalog

import "testing"

func TestReadOperationAdmitsGETAndRejectsWrites(t *testing.T) {
	operation, err := ReadOperation("events.proxy")
	if err != nil || operation.Method != "GET" || operation.Path != "/api/events/proxy" {
		t.Fatalf("operation=%+v err=%v", operation, err)
	}
	if _, err := ReadOperation("post.api.groups"); err == nil {
		t.Fatal("write operation was admitted")
	}
}

func TestReadOperationMarksIngressAsRuntimeUnverified(t *testing.T) {
	operation, err := ReadOperation("ingress.peers.list")
	if err != nil {
		t.Fatal(err)
	}
	if operation.Verification != "unverified_live" {
		t.Fatalf("ingress verification = %q, want unverified_live", operation.Verification)
	}
}
