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

func TestAgentNetworkSettingsDeleteImpactIsConservative(t *testing.T) {
	report, err := AgentNetworkSettingsDeleteImpact([]byte(`{"enabled":true}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "agent_network_settings_delete" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestAgentNetworkBudgetRuleImpactsAreConservative(t *testing.T) {
	create, err := AgentNetworkBudgetRuleCreateImpact([]byte(`{"name":"monthly","enabled":true}`))
	if err != nil {
		t.Fatal(err)
	}
	update, err := AgentNetworkBudgetRuleUpdateImpact([]byte(`{"id":"rule-1","enabled":true}`), []byte(`{"id":"rule-1","enabled":false}`))
	if err != nil {
		t.Fatal(err)
	}
	remove, err := AgentNetworkBudgetRuleDeleteImpact([]byte(`{"id":"rule-1","enabled":true}`))
	if err != nil {
		t.Fatal(err)
	}
	if create.Classification != "agent_network_budget_rule_create" || update.Classification != "agent_network_budget_rule_change" || remove.Classification != "agent_network_budget_rule_delete" {
		t.Fatalf("unexpected classifications: create=%+v update=%+v delete=%+v", create, update, remove)
	}
	for _, report := range []ImpactReport{create, update, remove} {
		if report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
			t.Fatalf("unexpected report: %+v", report)
		}
	}
}

func TestAgentNetworkGuardrailImpactsAreConservative(t *testing.T) {
	create, err := AgentNetworkGuardrailCreateImpact([]byte(`{"name":"strict","checks":{}}`))
	if err != nil {
		t.Fatal(err)
	}
	update, err := AgentNetworkGuardrailUpdateImpact([]byte(`{"id":"guard-1","checks":{}}`), []byte(`{"id":"guard-1","checks":{"model_allowlist":{"enabled":true}}}`))
	if err != nil {
		t.Fatal(err)
	}
	remove, err := AgentNetworkGuardrailDeleteImpact([]byte(`{"id":"guard-1","checks":{}}`))
	if err != nil {
		t.Fatal(err)
	}
	if create.Classification != "agent_network_guardrail_create" || update.Classification != "agent_network_guardrail_change" || remove.Classification != "agent_network_guardrail_delete" {
		t.Fatalf("unexpected classifications: create=%+v update=%+v delete=%+v", create, update, remove)
	}
	for _, report := range []ImpactReport{create, update, remove} {
		if report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
			t.Fatalf("unexpected report: %+v", report)
		}
	}
}

func TestAgentNetworkPolicyImpactsAreConservative(t *testing.T) {
	create, err := AgentNetworkPolicyCreateImpact([]byte(`{"name":"engineering","source_groups":["group-1"]}`))
	if err != nil {
		t.Fatal(err)
	}
	update, err := AgentNetworkPolicyUpdateImpact([]byte(`{"id":"policy-1","enabled":true}`), []byte(`{"id":"policy-1","enabled":false}`))
	if err != nil {
		t.Fatal(err)
	}
	remove, err := AgentNetworkPolicyDeleteImpact([]byte(`{"id":"policy-1","enabled":true}`))
	if err != nil {
		t.Fatal(err)
	}
	if create.Classification != "agent_network_policy_create" || update.Classification != "agent_network_policy_change" || remove.Classification != "agent_network_policy_delete" {
		t.Fatalf("unexpected classifications: create=%+v update=%+v delete=%+v", create, update, remove)
	}
	for _, report := range []ImpactReport{create, update, remove} {
		if report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
			t.Fatalf("unexpected report: %+v", report)
		}
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

func TestAgentNetworkProviderMutationsAreConservative(t *testing.T) {
	create, err := AgentNetworkProviderCreateImpact([]byte(`{"name":"OpenAI"}`))
	if err != nil {
		t.Fatal(err)
	}
	update, err := AgentNetworkProviderUpdateImpact([]byte(`{"id":"provider-1","name":"OpenAI"}`), []byte(`{"id":"provider-1","name":"OpenAI","enabled":false}`))
	if err != nil {
		t.Fatal(err)
	}
	remove, err := AgentNetworkProviderDeleteImpact([]byte(`{"id":"provider-1","name":"OpenAI"}`))
	if err != nil {
		t.Fatal(err)
	}
	for _, report := range []ImpactReport{create, update, remove} {
		if report.Reachability != "potentially_changed" || report.Confidence != "medium" || report.Completeness["state"] != "unknown" {
			t.Fatalf("unexpected provider report: %+v", report)
		}
	}
	if create.Classification != "agent_network_provider_create" || update.Classification != "agent_network_provider_change" || remove.Classification != "agent_network_provider_delete" {
		t.Fatalf("unexpected provider classifications: %q %q %q", create.Classification, update.Classification, remove.Classification)
	}
}

func TestUserLifecycleMutationsAreConservative(t *testing.T) {
	create, err := UserCreateImpact([]byte(`{"id":"user-1","email":"a@example.com"}`))
	if err != nil {
		t.Fatal(err)
	}
	update, err := UserUpdateImpact([]byte(`{"id":"user-1","role":"user"}`), []byte(`{"id":"user-1","role":"admin"}`))
	if err != nil {
		t.Fatal(err)
	}
	remove, err := UserDeleteImpact([]byte(`{"id":"user-1"}`))
	if err != nil {
		t.Fatal(err)
	}
	approve, err := UserApproveImpact([]byte(`{"id":"user-1","pending_approval":true}`))
	if err != nil {
		t.Fatal(err)
	}
	reject, err := UserRejectImpact([]byte(`{"id":"user-1","pending_approval":true}`))
	if err != nil {
		t.Fatal(err)
	}
	for _, report := range []ImpactReport{create, update, remove, approve, reject} {
		if report.Reachability != "potentially_changed" || report.Confidence != "medium" || report.Completeness["state"] != "unknown" {
			t.Fatalf("unexpected user report: %+v", report)
		}
	}
	if create.Classification != "user_create" || update.Classification != "user_change" || remove.Classification != "user_delete" || approve.Classification != "user_approve" || reject.Classification != "user_reject" {
		t.Fatalf("unexpected user classifications: %q %q %q %q %q", create.Classification, update.Classification, remove.Classification, approve.Classification, reject.Classification)
	}
}

func TestUserTokenDeleteImpactIsComplete(t *testing.T) {
	report, err := UserTokenDeleteImpact([]byte(`{"id":"token-1","name":"agent"}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "user_token_delete" || report.Confidence != "high" || report.Completeness["state"] != "complete" {
		t.Fatalf("unexpected token delete report: %+v", report)
	}
}

func TestUserTokenCreateImpactIsComplete(t *testing.T) {
	report, err := UserTokenCreateImpact([]byte(`{"id":"token-1","name":"agent"}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "user_token_create" || report.Reachability != "unchanged" || report.Confidence != "high" || report.Completeness["state"] != "complete" {
		t.Fatalf("unexpected token create report: %+v", report)
	}
}

func TestSetupKeyDeleteImpactIsConservative(t *testing.T) {
	report, err := SetupKeyDeleteImpact([]byte(`{"id":"key-1","name":"bootstrap","valid":true}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "setup_key_delete" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected setup key report: %+v", report)
	}
}

func TestSetupKeyCreateImpactIsConservative(t *testing.T) {
	report, err := SetupKeyCreateImpact([]byte(`{"id":"key-1","name":"bootstrap","valid":true}`))
	if err != nil {
		t.Fatal(err)
	}
	if report.Classification != "setup_key_create" || report.Reachability != "potentially_changed" || report.Completeness["state"] != "unknown" {
		t.Fatalf("unexpected setup key create report: %+v", report)
	}
}
