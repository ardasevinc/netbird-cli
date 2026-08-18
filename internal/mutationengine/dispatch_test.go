package mutationengine

import (
	"context"
	"encoding/json"
	"testing"
)

func TestDispatchDNSZoneUpdateStripsTargetID(t *testing.T) {
	remote := &fakeRemote{after: json.RawMessage(`{"id":"zone-1","name":"updated"}`)}
	_, err := dispatch(context.Background(), remote, "dns.zones.update", requestTarget{ID: "zone-1"}, json.RawMessage(`{"id":"zone-1","name":"updated"}`))
	if err != nil {
		t.Fatal(err)
	}
	if string(remote.dnsZoneBody) != `{"name":"updated"}` {
		t.Fatalf("unexpected update body: %s", remote.dnsZoneBody)
	}
}
