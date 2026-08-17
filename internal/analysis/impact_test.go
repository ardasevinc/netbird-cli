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

func TestPolicyUpdateImpactMarksRuleChangesAsPotentialReachabilityChange(t *testing.T) {
	report, err := PolicyUpdateImpact([]byte(`{"name":"p","rules":[{"id":"r1","action":"accept"}]}`), []byte(`{"name":"p","rules":[{"id":"r1","action":"drop"}]}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "policy_rule_change" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestPolicyDeleteImpactIsConservative(t *testing.T) {
	report, err := PolicyDeleteImpact([]byte(`{"id":"p1","name":"policy","rules":[]}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "policy_delete" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestPolicyCreateImpactIsConservative(t *testing.T) {
	report, err := PolicyCreateImpact([]byte(`{"name":"allow-office","enabled":true,"rules":[{"action":"accept"}]}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "policy_create" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestDNSZoneCreateImpactIsConservative(t *testing.T) {
	report, err := DNSZoneCreateImpact([]byte(`{"name":"office","domain":"office.internal","enabled":true}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "dns_zone_create" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestDNSZoneDeleteImpactIsConservative(t *testing.T) {
	report, err := DNSZoneDeleteImpact([]byte(`{"id":"zone-1","domain":"office.internal","distribution_groups":["g1"]}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "dns_zone_delete" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestDNSZoneUpdateImpactIsConservative(t *testing.T) {
	report, err := DNSZoneUpdateImpact(
		[]byte(`{"id":"zone-1","domain":"office.internal","enabled":true}`),
		[]byte(`{"id":"zone-1","domain":"corp.internal","enabled":true}`),
	)
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "dns_zone_change" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestDNSRecordCreateImpactIsConservative(t *testing.T) {
	report, err := DNSRecordCreateImpact([]byte(`{"name":"db","type":"A","content":"10.0.0.5","ttl":60}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "dns_record_create" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestDNSRecordUpdateImpactIsConservative(t *testing.T) {
	report, err := DNSRecordUpdateImpact(
		[]byte(`{"id":"record-1","name":"db","type":"A","content":"10.0.0.5","ttl":60}`),
		[]byte(`{"id":"record-1","name":"db","type":"A","content":"10.0.0.6","ttl":60}`),
	)
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "dns_record_change" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestDNSRecordDeleteImpactIsConservative(t *testing.T) {
	report, err := DNSRecordDeleteImpact([]byte(`{"id":"record-1","name":"db","type":"A","content":"10.0.0.5","ttl":60}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "dns_record_delete" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestAccountUpdateImpactIsConservative(t *testing.T) {
	report, err := AccountUpdateImpact(
		[]byte(`{"id":"account-1","settings":{"peer_login_expiration_enabled":true}}`),
		[]byte(`{"id":"account-1","settings":{"peer_login_expiration_enabled":false}}`),
	)
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "account_change" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestAccountDeleteImpactIsConservative(t *testing.T) {
	report, err := AccountDeleteImpact([]byte(`{"id":"account-1","domain":"example.test"}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "account_delete" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestPostureCheckCreateImpactIsConservative(t *testing.T) {
	report, err := PostureCheckCreateImpact([]byte(`{"name":"managed","checks":{"os_version_check":{}}}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "posture_check_create" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestPostureCheckUpdateImpactIsConservative(t *testing.T) {
	report, err := PostureCheckUpdateImpact(
		[]byte(`{"id":"pc-1","name":"managed","checks":{}}`),
		[]byte(`{"id":"pc-1","name":"managed-v2","checks":{}}`),
	)
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "posture_check_change" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestPostureCheckDeleteImpactIsConservative(t *testing.T) {
	report, err := PostureCheckDeleteImpact([]byte(`{"id":"pc-1","name":"managed","checks":{}}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "posture_check_delete" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestIngressPeerCreateImpactIsConservative(t *testing.T) {
	report, err := IngressPeerCreateImpact([]byte(`{"peer_id":"peer-1","region":"eu","enabled":true}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "ingress_peer_create" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestIngressPeerUpdateImpactIsConservative(t *testing.T) {
	report, err := IngressPeerUpdateImpact(
		[]byte(`{"id":"ing-1","peer_id":"peer-1","enabled":true}`),
		[]byte(`{"id":"ing-1","peer_id":"peer-1","enabled":false}`),
	)
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "ingress_peer_change" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestIngressPeerDeleteImpactIsConservative(t *testing.T) {
	report, err := IngressPeerDeleteImpact([]byte(`{"id":"ing-1","peer_id":"peer-1","enabled":true}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "ingress_peer_delete" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestAgentNetworkSettingsUpdateImpactIsConservative(t *testing.T) {
	report, err := AgentNetworkSettingsUpdateImpact(
		[]byte(`{"enabled":true}`),
		[]byte(`{"enabled":false}`),
	)
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "agent_network_settings_change" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestAgentNetworkSettingsCreateImpactIsConservative(t *testing.T) {
	report, err := AgentNetworkSettingsCreateImpact([]byte(`{"enabled":true}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "agent_network_settings_create" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestDNSNameserverCreateImpactIsConservative(t *testing.T) {
	report, err := DNSNameserverCreateImpact([]byte(`{"name":"office","domains":["office.internal"],"enabled":true,"nameservers":[{"ip":"10.0.0.53","ns_type":"udp","port":53}]}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "dns_nameserver_create" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestDNSNameserverUpdateImpactIsConservative(t *testing.T) {
	report, err := DNSNameserverUpdateImpact(
		[]byte(`{"id":"ns-1","description":"old","enabled":true}`),
		[]byte(`{"id":"ns-1","description":"new","enabled":true}`),
	)
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "dns_nameserver_change" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestDNSNameserverDeleteImpactIsConservative(t *testing.T) {
	report, err := DNSNameserverDeleteImpact([]byte(`{"id":"ns-1","domains":["office.internal"],"enabled":true}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "dns_nameserver_delete" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestDNSSettingsUpdateImpactIsConservative(t *testing.T) {
	report, err := DNSSettingsUpdateImpact(
		[]byte(`{"disabled_management_groups":["g1"]}`),
		[]byte(`{"disabled_management_groups":["g1","g2"]}`),
	)
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "dns_settings_change" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestRouteUpdateImpactMarksDescriptionOnlyChangeAsNeutral(t *testing.T) {
	report, err := RouteUpdateImpact(
		[]byte(`{"id":"r1","description":"old","enabled":true,"metric":10,"groups":["g1"]}`),
		[]byte(`{"id":"r1","description":"new","enabled":true,"metric":10,"groups":["g1"]}`),
	)
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "metadata_only" || report.Reachability != "unchanged" || report.Confidence != "high" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestRouteDeleteImpactIsConservative(t *testing.T) {
	report, err := RouteDeleteImpact([]byte(`{"id":"r1","description":"route","enabled":true,"network":"10.0.0.0/24"}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "route_delete" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestRouteCreateImpactIsConservative(t *testing.T) {
	report, err := RouteCreateImpact([]byte(`{"description":"private subnet","enabled":true,"network":"10.0.0.0/24","groups":["g1"]}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "route_create" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestRouteUpdateImpactBlocksRoutingChanges(t *testing.T) {
	report, err := RouteUpdateImpact(
		[]byte(`{"id":"r1","description":"route","enabled":true,"metric":10,"groups":["g1"]}`),
		[]byte(`{"id":"r1","description":"route","enabled":false,"metric":10,"groups":["g1"]}`),
	)
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "route_change" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestPeerUpdateImpactMarksNameOnlyChangeAsNeutral(t *testing.T) {
	report, err := PeerUpdateImpact(
		[]byte(`{"id":"p1","name":"old","approval_required":false,"connected":true}`),
		[]byte(`{"id":"p1","name":"new","approval_required":false,"connected":true}`),
	)
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "metadata_only" || report.Reachability != "unchanged" || report.Confidence != "high" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestGroupDeleteImpactIsConservative(t *testing.T) {
	report, err := GroupDeleteImpact([]byte(`{"id":"g1","name":"group","peers_count":2,"resources_count":1}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "group_delete" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestGroupCreateImpactIsConservative(t *testing.T) {
	report, err := GroupCreateImpact([]byte(`{"name":"engineering"}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "group_create" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestPeerUpdateImpactBlocksAccessChanges(t *testing.T) {
	report, err := PeerUpdateImpact(
		[]byte(`{"id":"p1","name":"peer","approval_required":false,"connected":true}`),
		[]byte(`{"id":"p1","name":"peer","approval_required":true,"connected":true}`),
	)
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "peer_change" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestPeerDeleteImpactIsConservative(t *testing.T) {
	report, err := PeerDeleteImpact([]byte(`{"id":"p1","name":"peer","connected":true,"approval_required":false}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "peer_delete" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestNetworkUpdateImpactMarksMetadataChangesAsNeutral(t *testing.T) {
	report, err := NetworkUpdateImpact(
		[]byte(`{"id":"n1","name":"old","description":"office","policies":["p1"],"resources":["r1"],"routers":["rt1"]}`),
		[]byte(`{"id":"n1","name":"new","description":"office","policies":["p1"],"resources":["r1"],"routers":["rt1"]}`),
	)
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "metadata_only" || report.Reachability != "unchanged" || report.Confidence != "high" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestNetworkUpdateImpactBlocksTopologyChanges(t *testing.T) {
	report, err := NetworkUpdateImpact(
		[]byte(`{"id":"n1","name":"office","policies":["p1"],"resources":["r1"],"routers":["rt1"]}`),
		[]byte(`{"id":"n1","name":"office","policies":["p2"],"resources":["r1"],"routers":["rt1"]}`),
	)
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "network_change" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestNetworkDeleteImpactIsConservative(t *testing.T) {
	report, err := NetworkDeleteImpact([]byte(`{"id":"n1","name":"office","policies":["p1"],"resources":["r1"],"routers":["rt1"]}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "network_delete" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestNetworkCreateImpactIsConservative(t *testing.T) {
	report, err := NetworkCreateImpact([]byte(`{"name":"office","description":"primary"}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "network_create" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestNetworkResourceDeleteImpactIsConservative(t *testing.T) {
	report, err := NetworkResourceDeleteImpact([]byte(`{"id":"r1","name":"db","address":"10.0.0.0/24","enabled":true}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "network_resource_delete" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestNetworkResourceUpdateImpactTreatsMetadataAsNeutral(t *testing.T) {
	report, err := NetworkResourceUpdateImpact(
		[]byte(`{"id":"r1","name":"old","description":"db","address":"10.0.0.0/24","enabled":true,"groups":["g1"]}`),
		[]byte(`{"id":"r1","name":"new","description":"db","address":"10.0.0.0/24","enabled":true,"groups":["g1"]}`),
	)
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "metadata_only" || report.Reachability != "unchanged" || report.Confidence != "high" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestNetworkResourceUpdateImpactBlocksTopologyChanges(t *testing.T) {
	report, err := NetworkResourceUpdateImpact(
		[]byte(`{"id":"r1","name":"db","address":"10.0.0.0/24","enabled":true,"groups":["g1"]}`),
		[]byte(`{"id":"r1","name":"db","address":"10.0.1.0/24","enabled":true,"groups":["g1"]}`),
	)
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "network_resource_change" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestNetworkRouterDeleteImpactIsConservative(t *testing.T) {
	report, err := NetworkRouterDeleteImpact([]byte(`{"id":"rt1","enabled":true,"masquerade":true,"metric":10,"peer":"p1"}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "network_router_delete" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestNetworkRouterUpdateImpactBlocksReachabilityChanges(t *testing.T) {
	report, err := NetworkRouterUpdateImpact(
		[]byte(`{"id":"rt1","enabled":true,"masquerade":true,"metric":10,"peer":"p1","peer_groups":["g1"]}`),
		[]byte(`{"id":"rt1","enabled":false,"masquerade":true,"metric":10,"peer":"p1","peer_groups":["g1"]}`),
	)
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "network_router_change" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestNetworkRouterCreateImpactIsConservative(t *testing.T) {
	report, err := NetworkRouterCreateImpact([]byte(`{"enabled":true,"masquerade":true,"metric":10,"peer":"p1"}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "network_router_create" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}
