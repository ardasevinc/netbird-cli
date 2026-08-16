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
