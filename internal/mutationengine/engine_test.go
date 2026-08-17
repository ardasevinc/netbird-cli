package mutationengine

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/ledger"
	"github.com/ardasevinc/netbird-cli/internal/mutation"
	"github.com/ardasevinc/netbird-cli/internal/transport"
)

type fakeRemote struct {
	identity             string
	account              string
	accountBefore        json.RawMessage
	accountAfter         json.RawMessage
	before               json.RawMessage
	after                json.RawMessage
	groupCollection      json.RawMessage
	policyBefore         json.RawMessage
	policyAfter          json.RawMessage
	policyCollection     json.RawMessage
	routeBefore          json.RawMessage
	routeAfter           json.RawMessage
	routeCollection      json.RawMessage
	peerBefore           json.RawMessage
	peerAfter            json.RawMessage
	networkBefore        json.RawMessage
	networkAfter         json.RawMessage
	networkCollection    json.RawMessage
	dnsZoneCollection    json.RawMessage
	nameserverCollection json.RawMessage
	nameserverBefore     json.RawMessage
	nameserverAfter      json.RawMessage
	dnsSettingsBefore    json.RawMessage
	dnsSettingsAfter     json.RawMessage
	resourceBefore       json.RawMessage
	resourceAfter        json.RawMessage
	resourceCollection   json.RawMessage
	routerBefore         json.RawMessage
	routerAfter          json.RawMessage
	routerCollection     json.RawMessage
	updateErr            error
	updates              int
}

func (f *fakeRemote) ServerIdentity() string { return f.identity }

func (f *fakeRemote) AccountScope(_ context.Context, account string) error {
	if account != f.account {
		return errors.New("wrong account")
	}
	return nil
}

func (f *fakeRemote) GetAccountRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.accountBefore == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.accountBefore...), nil
}

func (f *fakeRemote) UpdateAccount(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.accountBefore = append(json.RawMessage(nil), f.accountAfter...)
	return append(json.RawMessage(nil), f.accountAfter...), nil
}

func (f *fakeRemote) ListDNSZonesRaw(_ context.Context) (json.RawMessage, error) {
	if f.dnsZoneCollection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.dnsZoneCollection...), nil
}

func (f *fakeRemote) ListNameserverGroupsRaw(_ context.Context) (json.RawMessage, error) {
	if f.nameserverCollection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.nameserverCollection...), nil
}

func (f *fakeRemote) CreateNameserverGroup(_ context.Context, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	if f.after == nil {
		return nil, errors.New("missing created nameserver group")
	}
	f.nameserverCollection = json.RawMessage("[" + string(f.after) + "]")
	return append(json.RawMessage(nil), f.after...), nil
}

func (f *fakeRemote) GetNameserverGroupRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.nameserverBefore == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.nameserverBefore...), nil
}

func (f *fakeRemote) UpdateNameserverGroup(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.nameserverBefore = append(json.RawMessage(nil), f.nameserverAfter...)
	return append(json.RawMessage(nil), f.nameserverAfter...), nil
}

func (f *fakeRemote) DeleteNameserverGroup(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.nameserverBefore = nil
	return nil, nil
}

func (f *fakeRemote) GetDNSSettingsRaw(_ context.Context) (json.RawMessage, error) {
	if f.dnsSettingsBefore == nil {
		return json.RawMessage(`{"disabled_management_groups":[]}`), nil
	}
	return append(json.RawMessage(nil), f.dnsSettingsBefore...), nil
}

func (f *fakeRemote) UpdateDNSSettings(_ context.Context, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.dnsSettingsBefore = append(json.RawMessage(nil), f.dnsSettingsAfter...)
	return append(json.RawMessage(nil), f.dnsSettingsAfter...), nil
}

func (f *fakeRemote) CreateDNSZone(_ context.Context, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	if f.after == nil {
		return nil, errors.New("missing created dns zone")
	}
	f.dnsZoneCollection = json.RawMessage("[" + string(f.after) + "]")
	return append(json.RawMessage(nil), f.after...), nil
}

func (f *fakeRemote) GetDNSZoneRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.before == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.before...), nil
}

func (f *fakeRemote) UpdateDNSZone(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.before = append(json.RawMessage(nil), f.after...)
	return append(json.RawMessage(nil), f.after...), nil
}

func (f *fakeRemote) DeleteDNSZone(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.before = nil
	return nil, nil
}

func (f *fakeRemote) ListDNSRecordsRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.dnsZoneCollection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.dnsZoneCollection...), nil
}

func (f *fakeRemote) CreateDNSRecord(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	if f.after == nil {
		return nil, errors.New("missing created dns record")
	}
	f.dnsZoneCollection = json.RawMessage("[" + string(f.after) + "]")
	return append(json.RawMessage(nil), f.after...), nil
}

func (f *fakeRemote) GetDNSRecordRaw(_ context.Context, _, _ string) (json.RawMessage, error) {
	if f.before == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.before...), nil
}

func (f *fakeRemote) UpdateDNSRecord(_ context.Context, _, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.before = append(json.RawMessage(nil), f.after...)
	return append(json.RawMessage(nil), f.after...), nil
}

func (f *fakeRemote) DeleteDNSRecord(_ context.Context, _, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.before = nil
	return nil, nil
}

func (f *fakeRemote) GetGroup(_ context.Context, _ string) (json.RawMessage, error) {
	if f.before == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.before...), nil
}

func (f *fakeRemote) ListGroupsRaw(_ context.Context) (json.RawMessage, error) {
	if f.groupCollection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.groupCollection...), nil
}

func (f *fakeRemote) CreateGroup(_ context.Context, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	if f.after == nil {
		return nil, errors.New("missing created group")
	}
	f.groupCollection = json.RawMessage("[" + string(f.after) + "]")
	return append(json.RawMessage(nil), f.after...), nil
}

func (f *fakeRemote) UpdateGroup(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.before = append(json.RawMessage(nil), f.after...)
	return f.before, nil
}

func (f *fakeRemote) DeleteGroup(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.before = nil
	return nil, nil
}

func (f *fakeRemote) GetPolicyRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.policyBefore == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.policyBefore...), nil
}

func (f *fakeRemote) ListPoliciesRaw(_ context.Context) (json.RawMessage, error) {
	if f.policyCollection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.policyCollection...), nil
}

func (f *fakeRemote) CreatePolicy(_ context.Context, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	if f.policyAfter == nil {
		return nil, errors.New("missing created policy")
	}
	f.policyCollection = json.RawMessage("[" + string(f.policyAfter) + "]")
	return append(json.RawMessage(nil), f.policyAfter...), nil
}

func (f *fakeRemote) UpdatePolicy(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.policyBefore = append(json.RawMessage(nil), f.policyAfter...)
	return f.policyBefore, nil
}

func (f *fakeRemote) DeletePolicy(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.policyBefore = nil
	return nil, nil
}

func (f *fakeRemote) GetRouteRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.routeBefore == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.routeBefore...), nil
}

func (f *fakeRemote) ListRoutesRaw(_ context.Context) (json.RawMessage, error) {
	if f.routeCollection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.routeCollection...), nil
}

func (f *fakeRemote) CreateRoute(_ context.Context, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	if f.routeAfter == nil {
		return nil, errors.New("missing created route")
	}
	f.routeCollection = json.RawMessage("[" + string(f.routeAfter) + "]")
	return append(json.RawMessage(nil), f.routeAfter...), nil
}

func (f *fakeRemote) UpdateRoute(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.routeBefore = append(json.RawMessage(nil), f.routeAfter...)
	return f.routeBefore, nil
}

func (f *fakeRemote) DeleteRoute(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.routeBefore = nil
	return nil, nil
}

func (f *fakeRemote) GetPeerRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.peerBefore == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.peerBefore...), nil
}

func (f *fakeRemote) UpdatePeer(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.peerBefore = append(json.RawMessage(nil), f.peerAfter...)
	return f.peerBefore, nil
}

func (f *fakeRemote) DeletePeer(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.peerBefore = nil
	return nil, nil
}

func (f *fakeRemote) GetNetworkRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.networkBefore == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.networkBefore...), nil
}

func (f *fakeRemote) ListNetworksRaw(_ context.Context) (json.RawMessage, error) {
	if f.networkCollection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.networkCollection...), nil
}

func (f *fakeRemote) CreateNetwork(_ context.Context, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	if f.networkAfter == nil {
		return nil, errors.New("missing created network")
	}
	f.networkCollection = json.RawMessage("[" + string(f.networkAfter) + "]")
	return append(json.RawMessage(nil), f.networkAfter...), nil
}

func (f *fakeRemote) UpdateNetwork(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.networkBefore = append(json.RawMessage(nil), f.networkAfter...)
	return f.networkBefore, nil
}

func (f *fakeRemote) DeleteNetwork(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.networkBefore = nil
	return nil, nil
}

func (f *fakeRemote) GetNetworkResourceRaw(_ context.Context, _, _ string) (json.RawMessage, error) {
	if f.resourceBefore == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.resourceBefore...), nil
}

func (f *fakeRemote) ListNetworkResourcesRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.resourceCollection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.resourceCollection...), nil
}

func (f *fakeRemote) CreateNetworkResource(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	if f.resourceAfter == nil {
		return nil, errors.New("missing created resource")
	}
	f.resourceCollection = json.RawMessage("[" + string(f.resourceAfter) + "]")
	return append(json.RawMessage(nil), f.resourceAfter...), nil
}

func (f *fakeRemote) UpdateNetworkResource(_ context.Context, _, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.resourceBefore = append(json.RawMessage(nil), f.resourceAfter...)
	return f.resourceBefore, nil
}

func (f *fakeRemote) DeleteNetworkResource(_ context.Context, _, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.resourceBefore = nil
	return nil, nil
}

func (f *fakeRemote) GetNetworkRouterRaw(_ context.Context, _, _ string) (json.RawMessage, error) {
	if f.routerBefore == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.routerBefore...), nil
}

func (f *fakeRemote) ListNetworkRoutersRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.routerCollection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.routerCollection...), nil
}

func (f *fakeRemote) CreateNetworkRouter(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	if f.routerAfter == nil {
		return nil, errors.New("missing created router")
	}
	f.routerCollection = json.RawMessage("[" + string(f.routerAfter) + "]")
	return append(json.RawMessage(nil), f.routerAfter...), nil
}

func (f *fakeRemote) UpdateNetworkRouter(_ context.Context, _, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.routerBefore = append(json.RawMessage(nil), f.routerAfter...)
	return f.routerBefore, nil
}

func (f *fakeRemote) DeleteNetworkRouter(_ context.Context, _, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.routerBefore = nil
	return nil, nil
}

func stageForTest(t *testing.T, store *ledger.Store, before, after string) ledger.Stage {
	t.Helper()
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "groups.update",
		Request:        json.RawMessage(`{"id":"g1","name":"new"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
	})
	if err != nil {
		t.Fatal(err)
	}
	return stage
}

func TestApplyJournalsAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"g1","name":"old","peers_count":0,"resources_count":0}`
	after := `{"id":"g1","name":"new","peers_count":0,"resources_count":0}`
	stage := stageForTest(t, store, before, after)
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", before: []byte(before), after: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1"})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 || result.AttemptID == "" {
		t.Fatalf("unexpected result: %+v updates=%d", result, remote.updates)
	}
	receipt, err := store.GetReceipt(context.Background(), result.AttemptID)
	if err != nil || receipt.State != string(mutation.ConfirmedSuccess) {
		t.Fatalf("receipt=%+v err=%v", receipt, err)
	}
}

func TestApplyDispatchesGroupCreateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	intended := `{"name":"engineering"}`
	created := `{"id":"g1","name":"engineering","peers_count":0,"resources_count":0}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "groups.create",
		Request:        json.RawMessage(`{"name":"engineering"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(intended),
		Impact:         json.RawMessage(`{"classification":"group_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["creating a group can create a new policy or resource membership principal; affected peers and resources require live topology analysis"],"completeness":{"state":"unknown","reason":"group_create_requires_topology"}}`),
		Findings:       []ledger.Finding{{Code: "impact.group_create", Severity: "blocking", Message: "creating the group may alter policy or resource membership and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", groupCollection: []byte(before), after: []byte(created)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected group create result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesDNSZoneCreateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	intended := `{"name":"office","domain":"office.internal","enabled":true}`
	created := `{"id":"zone-1","name":"office","domain":"office.internal","enabled":true,"distribution_groups":[],"records":[]}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "dns.zones.create",
		Request:        json.RawMessage(intended),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(intended),
		Impact:         json.RawMessage(`{"classification":"dns_zone_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["creating a DNS zone can change name-resolution behavior for distributed peers; affected peers and records require live analysis"],"completeness":{"state":"unknown","reason":"dns_zone_create_requires_dns_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.dns_zone_create", Severity: "blocking", Message: "creating the DNS zone may alter name resolution and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", dnsZoneCollection: []byte(before), after: []byte(created)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected dns zone create result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesDNSZoneDeleteAndConfirmsAbsence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"zone-1","domain":"office.internal","distribution_groups":["g1"]}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "dns.zones.delete",
		Request:        json.RawMessage(`{"id":"zone-1"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(`{}`),
		Impact:         json.RawMessage(`{"classification":"dns_zone_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["deleting a DNS zone can change name-resolution behavior for distributed peers; affected peers and records require live analysis"],"completeness":{"state":"unknown","reason":"dns_zone_delete_requires_dns_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.dns_zone_delete", Severity: "blocking", Message: "deleting the DNS zone may alter name resolution and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", before: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected dns zone delete result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesDNSZoneUpdateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"zone-1","domain":"office.internal","enabled":true}`
	after := `{"id":"zone-1","domain":"corp.internal","enabled":true}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "dns.zones.update",
		Request:        json.RawMessage(after),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
		Impact:         json.RawMessage(`{"classification":"dns_zone_change","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["updating a DNS zone can change name-resolution behavior for distributed peers; affected peers and records require live analysis"],"completeness":{"state":"unknown","reason":"dns_zone_update_requires_dns_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.dns_zone_change", Severity: "blocking", Message: "the proposed DNS zone change may alter name resolution and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", before: []byte(before), after: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected dns zone update result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesDNSRecordCreateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	request := `{"zone_id":"zone-1","name":"db","type":"A","content":"10.0.0.5","ttl":60}`
	intended := `{"name":"db","type":"A","content":"10.0.0.5","ttl":60}`
	created := `{"id":"record-2","name":"db","type":"A","content":"10.0.0.5","ttl":60}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "dns.records.create",
		Request:        json.RawMessage(request),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(intended),
		Impact:         json.RawMessage(`{"classification":"dns_record_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["creating a DNS record can change name-resolution behavior for distributed peers; affected peers and records require live analysis"],"completeness":{"state":"unknown","reason":"dns_record_create_requires_dns_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.dns_record_create", Severity: "blocking", Message: "creating the DNS record may alter name resolution and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", dnsZoneCollection: []byte(before), after: []byte(created)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected dns record create result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesDNSRecordUpdateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"record-1","name":"db","type":"A","content":"10.0.0.5","ttl":60}`
	after := `{"id":"record-1","name":"db","type":"A","content":"10.0.0.6","ttl":60}`
	request := `{"zone_id":"zone-1","id":"record-1","name":"db","type":"A","content":"10.0.0.6","ttl":60}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "dns.records.update",
		Request:        json.RawMessage(request),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
		Impact:         json.RawMessage(`{"classification":"dns_record_change","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["updating a DNS record can change name-resolution behavior for distributed peers; affected peers and records require live analysis"],"completeness":{"state":"unknown","reason":"dns_record_update_requires_dns_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.dns_record_change", Severity: "blocking", Message: "the proposed DNS record change may alter name resolution and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", before: []byte(before), after: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected dns record update result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesDNSRecordDeleteAndConfirmsAbsence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"record-1","name":"db","type":"A","content":"10.0.0.5","ttl":60}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "dns.records.delete",
		Request:        json.RawMessage(`{"zone_id":"zone-1","id":"record-1"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(`{}`),
		Impact:         json.RawMessage(`{"classification":"dns_record_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["deleting a DNS record can change name-resolution behavior for distributed peers; affected peers and records require live analysis"],"completeness":{"state":"unknown","reason":"dns_record_delete_requires_dns_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.dns_record_delete", Severity: "blocking", Message: "deleting the DNS record may alter name resolution and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", before: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected dns record delete result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesNameserverGroupCreateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	intended := `{"name":"office","description":"office resolvers","domains":["office.internal"],"enabled":true,"groups":["g1"],"nameservers":[{"ip":"10.0.0.53","ns_type":"udp","port":53}],"primary":false,"search_domains_enabled":true}`
	created := `{"id":"ns-1","name":"office","description":"office resolvers","domains":["office.internal"],"enabled":true,"groups":["g1"],"nameservers":[{"ip":"10.0.0.53","ns_type":"udp","port":53}],"primary":false,"search_domains_enabled":true}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "dns.nameservers.create",
		Request:        json.RawMessage(intended),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(intended),
		Impact:         json.RawMessage(`{"classification":"dns_nameserver_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["creating a nameserver group can change resolver behavior for distributed peers; affected peers and domains require live analysis"],"completeness":{"state":"unknown","reason":"dns_nameserver_create_requires_dns_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.dns_nameserver_create", Severity: "blocking", Message: "creating the nameserver group may alter resolver behavior and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", nameserverCollection: []byte(before), after: []byte(created)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected nameserver create result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesNameserverGroupUpdateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"ns-1","name":"office","description":"old","domains":["office.internal"],"enabled":true,"groups":["g1"],"nameservers":[{"ip":"10.0.0.53","ns_type":"udp","port":53}],"primary":false,"search_domains_enabled":true}`
	after := `{"id":"ns-1","name":"office","description":"new","domains":["office.internal"],"enabled":true,"groups":["g1"],"nameservers":[{"ip":"10.0.0.53","ns_type":"udp","port":53}],"primary":false,"search_domains_enabled":true}`
	request := `{"id":"ns-1","description":"new"}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "dns.nameservers.update",
		Request:        json.RawMessage(request),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
		Impact:         json.RawMessage(`{"classification":"dns_nameserver_change","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["updating a nameserver group can change resolver behavior for distributed peers; affected peers and domains require live analysis"],"completeness":{"state":"unknown","reason":"dns_nameserver_update_requires_dns_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.dns_nameserver_change", Severity: "blocking", Message: "the proposed nameserver group change may alter resolver behavior and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", nameserverBefore: []byte(before), nameserverAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected nameserver update result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesNameserverGroupDeleteAndConfirmsAbsence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"ns-1","name":"office","description":"office resolvers","domains":["office.internal"],"enabled":true,"groups":["g1"],"nameservers":[{"ip":"10.0.0.53","ns_type":"udp","port":53}],"primary":false,"search_domains_enabled":true}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "dns.nameservers.delete",
		Request:        json.RawMessage(`{"id":"ns-1"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(`{}`),
		Impact:         json.RawMessage(`{"classification":"dns_nameserver_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["deleting a nameserver group can change resolver behavior for distributed peers; affected peers and domains require live analysis"],"completeness":{"state":"unknown","reason":"dns_nameserver_delete_requires_dns_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.dns_nameserver_delete", Severity: "blocking", Message: "deleting the nameserver group may alter resolver behavior and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", nameserverBefore: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected nameserver delete result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesDNSSettingsUpdateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"disabled_management_groups":["g1"]}`
	after := `{"disabled_management_groups":["g1","g2"]}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "dns.settings.update",
		Request:        json.RawMessage(after),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
		Impact:         json.RawMessage(`{"classification":"dns_settings_change","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["updating DNS settings can change resolver behavior for distributed peers; affected peers and domains require live analysis"],"completeness":{"state":"unknown","reason":"dns_settings_update_requires_dns_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.dns_settings_change", Severity: "blocking", Message: "the proposed DNS settings change may alter resolver behavior and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", dnsSettingsBefore: []byte(before), dnsSettingsAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected dns settings update result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesAccountUpdateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"account-1","settings":{"peer_login_expiration_enabled":true}}`
	after := `{"id":"account-1","settings":{"peer_login_expiration_enabled":false}}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "accounts.update",
		Request:        json.RawMessage(`{"id":"account-1","settings":{"peer_login_expiration_enabled":false}}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
		Impact:         json.RawMessage(`{"classification":"account_change","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["updating account settings can change authentication, posture, DNS, or network behavior; affected peers and resources require live analysis"],"completeness":{"state":"unknown","reason":"account_update_requires_management_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.account_change", Severity: "blocking", Message: "the proposed account change may alter management-plane behavior and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", accountBefore: []byte(before), accountAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected account update result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesGroupDeleteAndConfirmsAbsence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"g1","name":"group","peers_count":2,"resources_count":1}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "groups.delete",
		Request:        json.RawMessage(`{"id":"g1"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(`{}`),
		Impact:         json.RawMessage(`{"classification":"group_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["deleting a group can change policy membership and peer access; affected peers and resources require live topology analysis"],"completeness":{"state":"unknown","reason":"group_delete_requires_topology"}}`),
		Findings:       []ledger.Finding{{Code: "impact.group_delete", Severity: "blocking", Message: "deleting the group may alter policy membership and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", before: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected group delete result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyRefusesDriftBeforeJournalingOrDispatch(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	stage := stageForTest(t, store, `{"id":"g1","name":"old"}`, `{"id":"g1","name":"new"}`)
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", before: []byte(`{"id":"g1","name":"someone-else"}`), after: []byte(`{"id":"g1","name":"new"}`)}
	_, err = Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1"})
	if err == nil || remote.updates != 0 {
		t.Fatalf("drift was not refused: err=%v updates=%d", err, remote.updates)
	}
	if _, err := store.GetAttempt(context.Background(), "missing"); err == nil {
		t.Fatal("unexpected attempt")
	}
}

func TestApplyDoesNotReplayAfterAmbiguousDispatch(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"g1","name":"old"}`
	after := `{"id":"g1","name":"new"}`
	stage := stageForTest(t, store, before, after)
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", before: []byte(before), after: []byte(after), updateErr: &transport.RequestError{Dispatched: true, Description: "connection lost"}}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1"})
	if err == nil || result.State != mutation.Unknown {
		t.Fatalf("unexpected ambiguous result: %+v err=%v", result, err)
	}
	attempt, err := store.GetAttempt(context.Background(), result.AttemptID)
	if err != nil || attempt.State != string(mutation.Unknown) {
		t.Fatalf("attempt=%+v err=%v", attempt, err)
	}
}

func TestApplyReturnsAlreadySatisfiedWithoutWrite(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	after := `{"id":"g1","name":"new"}`
	stage := stageForTest(t, store, after, after)
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", before: []byte(after), after: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1"})
	if err != nil || result.State != mutation.AlreadySatisfied || remote.updates != 0 {
		t.Fatalf("unexpected no-op result: %+v updates=%d err=%v", result, remote.updates, err)
	}
}

func TestApplyClassifiesDefinitiveRemoteRejection(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"g1","name":"old"}`
	after := `{"id":"g1","name":"new"}`
	stage := stageForTest(t, store, before, after)
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", before: []byte(before), after: []byte(after), updateErr: &transport.RequestError{Dispatched: true, StatusCode: 403, Description: "remote rejected request"}}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1"})
	if err == nil || result.State != mutation.DefinitivelyRejected {
		t.Fatalf("unexpected rejection result: %+v err=%v", result, err)
	}
}

func TestApplyRefusesChangedImpactEvidence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "groups.update",
		Request:        json.RawMessage(`{"id":"g1","name":"new"}`),
		Before:         json.RawMessage(`{"id":"g1","name":"old"}`),
		IntendedAfter:  json.RawMessage(`{"id":"g1","name":"new"}`),
		Impact:         json.RawMessage(`{"classification":"unknown"}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", before: []byte(`{"id":"g1","name":"old"}`), after: []byte(`{"id":"g1","name":"new"}`)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1"})
	if err == nil || result.AttemptID != "" || remote.updates != 0 {
		t.Fatalf("changed impact was not refused before dispatch: result=%+v err=%v updates=%d", result, err, remote.updates)
	}
}

func TestApplyDispatchesPolicyUpdateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"p1","name":"old","rules":[]}`
	after := `{"id":"p1","name":"new","rules":[]}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "policies.update",
		Request:        json.RawMessage(`{"id":"p1","name":"new","rules":[]}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
		Impact:         json.RawMessage(`{"classification":"metadata_only","reachability":"unchanged","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"high","evidence":["policy metadata changed without changing policy rules"],"completeness":{"state":"complete","reason":null}}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", policyBefore: []byte(before), policyAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected policy result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesPolicyCreateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	intended := `{"name":"allow-office","enabled":true,"rules":[{"action":"accept"}]}`
	created := `{"id":"p1","name":"allow-office","enabled":true,"rules":[{"action":"accept"}],"source_posture_checks":[]}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "policies.create",
		Request:        json.RawMessage(`{"name":"allow-office","enabled":true,"rules":[{"action":"accept"}]}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(intended),
		Impact:         json.RawMessage(`{"classification":"policy_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["creating a policy can add access edges; affected peers and resources require live topology analysis"],"completeness":{"state":"unknown","reason":"policy_create_requires_topology"}}`),
		Findings:       []ledger.Finding{{Code: "impact.policy_create", Severity: "blocking", Message: "creating the policy may add reachability and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", policyCollection: []byte(before), policyAfter: []byte(created)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected policy create result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesPolicyDeleteAndConfirmsAbsence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"p1","name":"policy","rules":[]}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "policies.delete",
		Request:        json.RawMessage(`{"id":"p1"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(`{}`),
		Impact:         json.RawMessage(`{"classification":"policy_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["deleting a policy can remove access edges; affected peers and resources require live topology analysis"],"completeness":{"state":"unknown","reason":"policy_delete_requires_topology"}}`),
		Findings:       []ledger.Finding{{Code: "impact.policy_delete", Severity: "blocking", Message: "deleting the policy may remove access and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", policyBefore: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected policy delete result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesRouteUpdateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"r1","description":"old","enabled":true,"metric":10,"groups":["g1"]}`
	after := `{"id":"r1","description":"new","enabled":true,"metric":10,"groups":["g1"]}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "routes.update",
		Request:        json.RawMessage(`{"id":"r1","description":"new","enabled":true,"metric":10,"groups":["g1"]}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
		Impact:         json.RawMessage(`{"classification":"metadata_only","reachability":"unchanged","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"high","evidence":["route description changed without changing routing behavior"],"completeness":{"state":"complete","reason":null}}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", routeBefore: []byte(before), routeAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected route result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesRouteCreateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	intended := `{"description":"private subnet","enabled":true,"network":"10.0.0.0/24","groups":["g1"]}`
	created := `{"id":"r1","description":"private subnet","enabled":true,"network":"10.0.0.0/24","groups":["g1"],"metric":10}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "routes.create",
		Request:        json.RawMessage(`{"description":"private subnet","enabled":true,"network":"10.0.0.0/24","groups":["g1"]}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(intended),
		Impact:         json.RawMessage(`{"classification":"route_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["creating a route can add a reachable destination or alter path selection; affected peers and resources require live topology analysis"],"completeness":{"state":"unknown","reason":"route_create_requires_topology"}}`),
		Findings:       []ledger.Finding{{Code: "impact.route_create", Severity: "blocking", Message: "creating the route may alter reachability and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", routeCollection: []byte(before), routeAfter: []byte(created)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected route create result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesRouteDeleteAndConfirmsAbsence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"r1","description":"route","enabled":true,"network":"10.0.0.0/24"}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "routes.delete",
		Request:        json.RawMessage(`{"id":"r1"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(`{}`),
		Impact:         json.RawMessage(`{"classification":"route_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["deleting a route can change network reachability; affected peers and resources require live topology analysis"],"completeness":{"state":"unknown","reason":"route_delete_requires_topology"}}`),
		Findings:       []ledger.Finding{{Code: "impact.route_delete", Severity: "blocking", Message: "deleting the route may alter reachability and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", routeBefore: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected route delete result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesPeerUpdateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"p1","name":"old","approval_required":false,"connected":true}`
	after := `{"id":"p1","name":"new","approval_required":false,"connected":true}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "peers.update",
		Request:        json.RawMessage(`{"id":"p1","name":"new"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
		Impact:         json.RawMessage(`{"classification":"metadata_only","reachability":"unchanged","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"high","evidence":["peer name changed without changing peer access or connectivity state"],"completeness":{"state":"complete","reason":null}}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", peerBefore: []byte(before), peerAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected peer result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesPeerDeleteAndConfirmsAbsence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"p1","name":"peer","connected":true,"approval_required":false}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "peers.delete",
		Request:        json.RawMessage(`{"id":"p1"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(`{}`),
		Impact:         json.RawMessage(`{"classification":"peer_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["deleting a peer removes an access and connectivity principal; affected peers and resources require live topology analysis"],"completeness":{"state":"unknown","reason":"peer_delete_requires_topology"}}`),
		Findings:       []ledger.Finding{{Code: "impact.peer_delete", Severity: "blocking", Message: "deleting the peer may remove access or connectivity and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", peerBefore: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected peer delete result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesNetworkUpdateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"n1","name":"old","description":"office","policies":["p1"],"resources":["r1"],"routers":["rt1"]}`
	after := `{"id":"n1","name":"new","description":"office","policies":["p1"],"resources":["r1"],"routers":["rt1"]}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "networks.update",
		Request:        json.RawMessage(`{"id":"n1","name":"new"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
		Impact:         json.RawMessage(`{"classification":"metadata_only","reachability":"unchanged","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"high","evidence":["network metadata changed without changing attached policies, resources, or routers"],"completeness":{"state":"complete","reason":null}}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", networkBefore: []byte(before), networkAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected network result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesNetworkCreateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	intended := `{"name":"office","description":"primary"}`
	created := `{"id":"n1","name":"office","description":"primary","policies":[],"resources":[],"routers":[],"routing_peers_count":0}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "networks.create",
		Request:        json.RawMessage(`{"name":"office","description":"primary"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(intended),
		Impact:         json.RawMessage(`{"classification":"network_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["creating a network can add a topology scope for policies, resources, and routers; affected peers and resources require live topology analysis"],"completeness":{"state":"unknown","reason":"network_create_requires_topology"}}`),
		Findings:       []ledger.Finding{{Code: "impact.network_create", Severity: "blocking", Message: "creating the network may add topology and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", networkCollection: []byte(before), networkAfter: []byte(created)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected network create result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesNetworkDeleteAndConfirmsAbsence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"n1","name":"office","policies":["p1"],"resources":["r1"],"routers":["rt1"]}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "networks.delete",
		Request:        json.RawMessage(`{"id":"n1"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(`{}`),
		Impact:         json.RawMessage(`{"classification":"network_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["deleting a network can remove attached policies, resources, and routers; affected peers and resources require live topology analysis"],"completeness":{"state":"unknown","reason":"network_delete_requires_topology"}}`),
		Findings:       []ledger.Finding{{Code: "impact.network_delete", Severity: "blocking", Message: "deleting the network may remove attached topology and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", networkBefore: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected network delete result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesNetworkResourceDeleteAndConfirmsAbsence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"r1","name":"db","address":"10.0.0.0/24","enabled":true}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "networks.resources.delete",
		Request:        json.RawMessage(`{"network_id":"n1","id":"r1"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(`{}`),
		Impact:         json.RawMessage(`{"classification":"network_resource_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["deleting a network resource can remove a reachable destination; affected peers and routes require live topology analysis"],"completeness":{"state":"unknown","reason":"network_resource_delete_requires_topology"}}`),
		Findings:       []ledger.Finding{{Code: "impact.network_resource_delete", Severity: "blocking", Message: "deleting the network resource may remove reachability and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", resourceBefore: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected network resource delete result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesNetworkResourceUpdateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"r1","name":"old","description":"db","address":"10.0.0.0/24","enabled":true,"groups":["g1"]}`
	after := `{"id":"r1","name":"new","description":"db","address":"10.0.0.0/24","enabled":true,"groups":["g1"]}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "networks.resources.update",
		Request:        json.RawMessage(`{"network_id":"n1","id":"r1","name":"new"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
		Impact:         json.RawMessage(`{"classification":"metadata_only","reachability":"unchanged","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"high","evidence":["network resource metadata changed without changing its address, enablement, or group assignment"],"completeness":{"state":"complete","reason":null}}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", resourceBefore: []byte(before), resourceAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected network resource update result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesNetworkResourceCreateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	intended := `{"name":"db","address":"10.0.0.0/24","enabled":true,"groups":["g1"]}`
	created := `{"id":"r1","name":"db","address":"10.0.0.0/24","enabled":true,"groups":["g1"],"type":"subnet"}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "networks.resources.create",
		Request:        json.RawMessage(`{"network_id":"n1","name":"db","address":"10.0.0.0/24","enabled":true,"groups":["g1"]}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(intended),
		Impact:         json.RawMessage(`{"classification":"network_resource_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["creating a network resource can add a reachable destination; affected peers and routes require live topology analysis"],"completeness":{"state":"unknown","reason":"network_resource_create_requires_topology"}}`),
		Findings:       []ledger.Finding{{Code: "impact.network_resource_create", Severity: "blocking", Message: "creating the network resource may add reachability and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", resourceCollection: []byte(before), resourceAfter: []byte(created)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected network resource create result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesNetworkRouterCreateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	intended := `{"enabled":true,"masquerade":true,"metric":10,"peer":"p1"}`
	created := `{"id":"rt1","enabled":true,"masquerade":true,"metric":10,"peer":"p1","peer_groups":[]}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "networks.routers.create",
		Request:        json.RawMessage(`{"network_id":"n1","enabled":true,"masquerade":true,"metric":10,"peer":"p1"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(intended),
		Impact:         json.RawMessage(`{"classification":"network_router_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["creating a network router can add route availability and peer reachability; affected peers and resources require live topology analysis"],"completeness":{"state":"unknown","reason":"network_router_create_requires_topology"}}`),
		Findings:       []ledger.Finding{{Code: "impact.network_router_create", Severity: "blocking", Message: "creating the network router may alter reachability and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", routerCollection: []byte(before), routerAfter: []byte(created)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected network router create result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesNetworkRouterDeleteAndConfirmsAbsence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"rt1","enabled":true,"masquerade":true,"metric":10,"peer":"p1"}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "networks.routers.delete",
		Request:        json.RawMessage(`{"network_id":"n1","id":"rt1"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(`{}`),
		Impact:         json.RawMessage(`{"classification":"network_router_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["deleting a network router can change route availability and peer reachability; affected peers and resources require live topology analysis"],"completeness":{"state":"unknown","reason":"network_router_delete_requires_topology"}}`),
		Findings:       []ledger.Finding{{Code: "impact.network_router_delete", Severity: "blocking", Message: "deleting the network router may alter reachability and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", routerBefore: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected network router delete result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesNetworkRouterUpdateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"rt1","enabled":true,"masquerade":true,"metric":10,"peer":"p1","peer_groups":["g1"]}`
	after := `{"id":"rt1","enabled":false,"masquerade":true,"metric":10,"peer":"p1","peer_groups":["g1"]}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "networks.routers.update",
		Request:        json.RawMessage(`{"network_id":"n1","id":"rt1","enabled":false}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
		Impact:         json.RawMessage(`{"classification":"network_router_change","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["network router update changes reachability fields: [enabled]; affected peers and resources require live topology analysis"],"completeness":{"state":"unknown","reason":"network_router_change_requires_topology"}}`),
		Findings:       []ledger.Finding{{Code: "impact.network_router_change", Severity: "blocking", Message: "the proposed network router change may alter reachability and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", routerBefore: []byte(before), routerAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected network router update result: %+v updates=%d", result, remote.updates)
	}
}
