package analysis

import (
	"encoding/json"
	"fmt"
	"sort"
)

// ImpactReport is the deterministic mutation-impact evidence stored with a
// stage. It is deliberately conservative: a group update is treated as
// reachability-neutral only when the normalized before/after documents differ
// in the group name alone.
type ImpactReport struct {
	Classification    string         `json:"classification"`
	Reachability      string         `json:"reachability"`
	AffectedPeers     []string       `json:"affected_peer_ids"`
	AffectedResources []string       `json:"affected_resource_ids"`
	Confidence        string         `json:"confidence"`
	Evidence          []string       `json:"evidence"`
	Completeness      map[string]any `json:"completeness"`
}

func PolicyUpdateImpact(before, intendedAfter []byte) (ImpactReport, error) {
	var beforeObject, afterObject map[string]any
	if err := json.Unmarshal(before, &beforeObject); err != nil {
		return ImpactReport{}, fmt.Errorf("decode policy impact preimage: %w", err)
	}
	if err := json.Unmarshal(intendedAfter, &afterObject); err != nil {
		return ImpactReport{}, fmt.Errorf("decode policy impact intended state: %w", err)
	}
	rulesChanged := !sameJSONValue(beforeObject["rules"], afterObject["rules"])
	if !rulesChanged {
		return ImpactReport{
			Classification:    "metadata_only",
			Reachability:      "unchanged",
			AffectedPeers:     []string{},
			AffectedResources: []string{},
			Confidence:        "high",
			Evidence:          []string{"policy metadata changed without changing policy rules"},
			Completeness:      map[string]any{"state": "complete", "reason": nil},
		}, nil
	}
	return ImpactReport{
		Classification:    "policy_rule_change",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"policy rules changed; affected peers and resources require live topology analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "policy_rule_change_requires_topology"},
	}, nil
}

func PolicyDeleteImpact(before []byte) (ImpactReport, error) {
	var object map[string]any
	if err := json.Unmarshal(before, &object); err != nil {
		return ImpactReport{}, fmt.Errorf("decode policy delete preimage: %w", err)
	}
	return ImpactReport{
		Classification:    "policy_delete",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"deleting a policy can remove access edges; affected peers and resources require live topology analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "policy_delete_requires_topology"},
	}, nil
}

func PolicyCreateImpact(intendedAfter []byte) (ImpactReport, error) {
	var object map[string]any
	if err := json.Unmarshal(intendedAfter, &object); err != nil {
		return ImpactReport{}, fmt.Errorf("decode policy create intent: %w", err)
	}
	return ImpactReport{
		Classification:    "policy_create",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"creating a policy can add access edges; affected peers and resources require live topology analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "policy_create_requires_topology"},
	}, nil
}

func DNSZoneCreateImpact(intendedAfter []byte) (ImpactReport, error) {
	var object map[string]any
	if err := json.Unmarshal(intendedAfter, &object); err != nil {
		return ImpactReport{}, fmt.Errorf("decode dns zone create intent: %w", err)
	}
	return ImpactReport{
		Classification:    "dns_zone_create",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"creating a DNS zone can change name-resolution behavior for distributed peers; affected peers and records require live analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "dns_zone_create_requires_dns_analysis"},
	}, nil
}

func DNSZoneDeleteImpact(before []byte) (ImpactReport, error) {
	var object map[string]any
	if err := json.Unmarshal(before, &object); err != nil {
		return ImpactReport{}, fmt.Errorf("decode dns zone delete preimage: %w", err)
	}
	return ImpactReport{
		Classification:    "dns_zone_delete",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"deleting a DNS zone can change name-resolution behavior for distributed peers; affected peers and records require live analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "dns_zone_delete_requires_dns_analysis"},
	}, nil
}

func DNSZoneUpdateImpact(before, intendedAfter []byte) (ImpactReport, error) {
	var beforeObject, afterObject map[string]any
	if err := json.Unmarshal(before, &beforeObject); err != nil {
		return ImpactReport{}, fmt.Errorf("decode dns zone impact preimage: %w", err)
	}
	if err := json.Unmarshal(intendedAfter, &afterObject); err != nil {
		return ImpactReport{}, fmt.Errorf("decode dns zone impact intended state: %w", err)
	}
	return ImpactReport{
		Classification:    "dns_zone_change",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"updating a DNS zone can change name-resolution behavior for distributed peers; affected peers and records require live analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "dns_zone_update_requires_dns_analysis"},
	}, nil
}

func DNSRecordCreateImpact(intendedAfter []byte) (ImpactReport, error) {
	var object map[string]any
	if err := json.Unmarshal(intendedAfter, &object); err != nil {
		return ImpactReport{}, fmt.Errorf("decode dns record create intent: %w", err)
	}
	return ImpactReport{
		Classification:    "dns_record_create",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"creating a DNS record can change name-resolution behavior for distributed peers; affected peers and records require live analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "dns_record_create_requires_dns_analysis"},
	}, nil
}

func DNSRecordUpdateImpact(before, intendedAfter []byte) (ImpactReport, error) {
	var beforeObject, afterObject map[string]any
	if err := json.Unmarshal(before, &beforeObject); err != nil {
		return ImpactReport{}, fmt.Errorf("decode dns record impact preimage: %w", err)
	}
	if err := json.Unmarshal(intendedAfter, &afterObject); err != nil {
		return ImpactReport{}, fmt.Errorf("decode dns record impact intended state: %w", err)
	}
	return ImpactReport{
		Classification:    "dns_record_change",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"updating a DNS record can change name-resolution behavior for distributed peers; affected peers and records require live analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "dns_record_update_requires_dns_analysis"},
	}, nil
}

func DNSRecordDeleteImpact(before []byte) (ImpactReport, error) {
	var object map[string]any
	if err := json.Unmarshal(before, &object); err != nil {
		return ImpactReport{}, fmt.Errorf("decode dns record delete preimage: %w", err)
	}
	return ImpactReport{
		Classification:    "dns_record_delete",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"deleting a DNS record can change name-resolution behavior for distributed peers; affected peers and records require live analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "dns_record_delete_requires_dns_analysis"},
	}, nil
}

func DNSNameserverCreateImpact(intendedAfter []byte) (ImpactReport, error) {
	var object map[string]any
	if err := json.Unmarshal(intendedAfter, &object); err != nil {
		return ImpactReport{}, fmt.Errorf("decode dns nameserver create intent: %w", err)
	}
	return ImpactReport{
		Classification:    "dns_nameserver_create",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"creating a nameserver group can change resolver behavior for distributed peers; affected peers and domains require live analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "dns_nameserver_create_requires_dns_analysis"},
	}, nil
}

func DNSNameserverUpdateImpact(before, intendedAfter []byte) (ImpactReport, error) {
	var beforeObject, afterObject map[string]any
	if err := json.Unmarshal(before, &beforeObject); err != nil {
		return ImpactReport{}, fmt.Errorf("decode dns nameserver impact preimage: %w", err)
	}
	if err := json.Unmarshal(intendedAfter, &afterObject); err != nil {
		return ImpactReport{}, fmt.Errorf("decode dns nameserver impact intended state: %w", err)
	}
	return ImpactReport{
		Classification:    "dns_nameserver_change",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"updating a nameserver group can change resolver behavior for distributed peers; affected peers and domains require live analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "dns_nameserver_update_requires_dns_analysis"},
	}, nil
}

func DNSNameserverDeleteImpact(before []byte) (ImpactReport, error) {
	var object map[string]any
	if err := json.Unmarshal(before, &object); err != nil {
		return ImpactReport{}, fmt.Errorf("decode dns nameserver delete preimage: %w", err)
	}
	return ImpactReport{
		Classification:    "dns_nameserver_delete",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"deleting a nameserver group can change resolver behavior for distributed peers; affected peers and domains require live analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "dns_nameserver_delete_requires_dns_analysis"},
	}, nil
}

func DNSSettingsUpdateImpact(before, intendedAfter []byte) (ImpactReport, error) {
	var beforeObject, afterObject map[string]any
	if err := json.Unmarshal(before, &beforeObject); err != nil {
		return ImpactReport{}, fmt.Errorf("decode dns settings impact preimage: %w", err)
	}
	if err := json.Unmarshal(intendedAfter, &afterObject); err != nil {
		return ImpactReport{}, fmt.Errorf("decode dns settings impact intended state: %w", err)
	}
	return ImpactReport{
		Classification:    "dns_settings_change",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"updating DNS settings can change resolver behavior for distributed peers; affected peers and domains require live analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "dns_settings_update_requires_dns_analysis"},
	}, nil
}

func AccountUpdateImpact(before, intendedAfter []byte) (ImpactReport, error) {
	var beforeObject, afterObject map[string]any
	if err := json.Unmarshal(before, &beforeObject); err != nil {
		return ImpactReport{}, fmt.Errorf("decode account impact preimage: %w", err)
	}
	if err := json.Unmarshal(intendedAfter, &afterObject); err != nil {
		return ImpactReport{}, fmt.Errorf("decode account impact intended state: %w", err)
	}
	return ImpactReport{
		Classification:    "account_change",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"updating account settings can change authentication, posture, DNS, or network behavior; affected peers and resources require live analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "account_update_requires_management_analysis"},
	}, nil
}

func AccountDeleteImpact(before []byte) (ImpactReport, error) {
	var object map[string]any
	if err := json.Unmarshal(before, &object); err != nil {
		return ImpactReport{}, fmt.Errorf("decode account delete preimage: %w", err)
	}
	return ImpactReport{
		Classification:    "account_delete",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"deleting an account can remove management access and account resources; affected peers and resources require live analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "account_delete_requires_management_analysis"},
	}, nil
}

func PostureCheckCreateImpact(intendedAfter []byte) (ImpactReport, error) {
	var object map[string]any
	if err := json.Unmarshal(intendedAfter, &object); err != nil {
		return ImpactReport{}, fmt.Errorf("decode posture check create intent: %w", err)
	}
	return ImpactReport{
		Classification:    "posture_check_create",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"creating a posture check can change policy admission for peers; affected peers and policies require live analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "posture_check_create_requires_policy_analysis"},
	}, nil
}

func PostureCheckUpdateImpact(before, intendedAfter []byte) (ImpactReport, error) {
	var beforeObject, afterObject map[string]any
	if err := json.Unmarshal(before, &beforeObject); err != nil {
		return ImpactReport{}, fmt.Errorf("decode posture check impact preimage: %w", err)
	}
	if err := json.Unmarshal(intendedAfter, &afterObject); err != nil {
		return ImpactReport{}, fmt.Errorf("decode posture check impact intended state: %w", err)
	}
	return ImpactReport{
		Classification:    "posture_check_change",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"updating a posture check can change policy admission for peers; affected peers and policies require live analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "posture_check_update_requires_policy_analysis"},
	}, nil
}

func PostureCheckDeleteImpact(before []byte) (ImpactReport, error) {
	var object map[string]any
	if err := json.Unmarshal(before, &object); err != nil {
		return ImpactReport{}, fmt.Errorf("decode posture check delete preimage: %w", err)
	}
	return ImpactReport{
		Classification:    "posture_check_delete",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"deleting a posture check can change policy admission for peers; affected peers and policies require live analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "posture_check_delete_requires_policy_analysis"},
	}, nil
}

func IngressPeerCreateImpact(intendedAfter []byte) (ImpactReport, error) {
	var object map[string]any
	if err := json.Unmarshal(intendedAfter, &object); err != nil {
		return ImpactReport{}, fmt.Errorf("decode ingress peer create intent: %w", err)
	}
	return ImpactReport{
		Classification:    "ingress_peer_create",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"creating an ingress peer can add an external reachability path; affected peers and allocations require live topology analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "ingress_peer_create_requires_topology_analysis"},
	}, nil
}

func IngressPeerUpdateImpact(before, intendedAfter []byte) (ImpactReport, error) {
	var beforeObject, afterObject map[string]any
	if err := json.Unmarshal(before, &beforeObject); err != nil {
		return ImpactReport{}, fmt.Errorf("decode ingress peer impact preimage: %w", err)
	}
	if err := json.Unmarshal(intendedAfter, &afterObject); err != nil {
		return ImpactReport{}, fmt.Errorf("decode ingress peer impact intended state: %w", err)
	}
	return ImpactReport{
		Classification:    "ingress_peer_change",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"updating an ingress peer can change an external reachability path; affected peers and allocations require live topology analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "ingress_peer_update_requires_topology_analysis"},
	}, nil
}

func IngressPeerDeleteImpact(before []byte) (ImpactReport, error) {
	var object map[string]any
	if err := json.Unmarshal(before, &object); err != nil {
		return ImpactReport{}, fmt.Errorf("decode ingress peer delete preimage: %w", err)
	}
	return ImpactReport{
		Classification:    "ingress_peer_delete",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"deleting an ingress peer can remove an external reachability path; affected peers and allocations require live topology analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "ingress_peer_delete_requires_topology_analysis"},
	}, nil
}

func AgentNetworkSettingsUpdateImpact(before, intendedAfter []byte) (ImpactReport, error) {
	var beforeObject, afterObject map[string]any
	if err := json.Unmarshal(before, &beforeObject); err != nil {
		return ImpactReport{}, fmt.Errorf("decode agent-network settings impact preimage: %w", err)
	}
	if err := json.Unmarshal(intendedAfter, &afterObject); err != nil {
		return ImpactReport{}, fmt.Errorf("decode agent-network settings impact intended state: %w", err)
	}
	return ImpactReport{
		Classification:    "agent_network_settings_change",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"updating agent-network settings can change Cloud agent routing and provider behavior; affected peers and account resources require capability-aware live analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "agent_network_settings_update_requires_capability_analysis"},
	}, nil
}

func AgentNetworkSettingsCreateImpact(intendedAfter []byte) (ImpactReport, error) {
	var object map[string]any
	if err := json.Unmarshal(intendedAfter, &object); err != nil {
		return ImpactReport{}, fmt.Errorf("decode agent-network settings create intent: %w", err)
	}
	return ImpactReport{
		Classification:    "agent_network_settings_create",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"creating agent-network settings can enable Cloud agent routing and provider behavior; affected peers and account resources require capability-aware live analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "agent_network_settings_create_requires_capability_analysis"},
	}, nil
}

func AgentNetworkSettingsDeleteImpact(before []byte) (ImpactReport, error) {
	var object map[string]any
	if err := json.Unmarshal(before, &object); err != nil {
		return ImpactReport{}, fmt.Errorf("decode agent-network settings delete preimage: %w", err)
	}
	return ImpactReport{
		Classification:    "agent_network_settings_delete",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"deleting agent-network settings can disable Cloud agent routing and provider behavior; affected peers and account resources require capability-aware live analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "agent_network_settings_delete_requires_capability_analysis"},
	}, nil
}

func AgentNetworkBudgetRuleCreateImpact(intendedAfter []byte) (ImpactReport, error) {
	var object map[string]any
	if err := json.Unmarshal(intendedAfter, &object); err != nil {
		return ImpactReport{}, fmt.Errorf("decode agent-network budget rule create intent: %w", err)
	}
	return ImpactReport{
		Classification:    "agent_network_budget_rule_create",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"creating an agent-network budget rule can change spend and token admission for targeted callers; affected peers and account resources require capability-aware live analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "agent_network_budget_rule_create_requires_capability_analysis"},
	}, nil
}

func AgentNetworkBudgetRuleUpdateImpact(before, intendedAfter []byte) (ImpactReport, error) {
	var beforeObject, afterObject map[string]any
	if err := json.Unmarshal(before, &beforeObject); err != nil {
		return ImpactReport{}, fmt.Errorf("decode agent-network budget rule impact preimage: %w", err)
	}
	if err := json.Unmarshal(intendedAfter, &afterObject); err != nil {
		return ImpactReport{}, fmt.Errorf("decode agent-network budget rule impact intended state: %w", err)
	}
	return ImpactReport{
		Classification:    "agent_network_budget_rule_change",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"updating an agent-network budget rule can change spend and token admission for targeted callers; affected peers and account resources require capability-aware live analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "agent_network_budget_rule_update_requires_capability_analysis"},
	}, nil
}

func AgentNetworkBudgetRuleDeleteImpact(before []byte) (ImpactReport, error) {
	var object map[string]any
	if err := json.Unmarshal(before, &object); err != nil {
		return ImpactReport{}, fmt.Errorf("decode agent-network budget rule delete preimage: %w", err)
	}
	return ImpactReport{
		Classification:    "agent_network_budget_rule_delete",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"deleting an agent-network budget rule can remove spend and token admission limits for targeted callers; affected peers and account resources require capability-aware live analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "agent_network_budget_rule_delete_requires_capability_analysis"},
	}, nil
}

func AgentNetworkGuardrailCreateImpact(intendedAfter []byte) (ImpactReport, error) {
	var object map[string]any
	if err := json.Unmarshal(intendedAfter, &object); err != nil {
		return ImpactReport{}, fmt.Errorf("decode agent-network guardrail create intent: %w", err)
	}
	return ImpactReport{
		Classification:    "agent_network_guardrail_create",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"creating an agent-network guardrail can alter model and prompt admission for attached policies; affected peers and account resources require capability-aware live analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "agent_network_guardrail_create_requires_capability_analysis"},
	}, nil
}

func AgentNetworkGuardrailUpdateImpact(before, intendedAfter []byte) (ImpactReport, error) {
	var beforeObject, afterObject map[string]any
	if err := json.Unmarshal(before, &beforeObject); err != nil {
		return ImpactReport{}, fmt.Errorf("decode agent-network guardrail impact preimage: %w", err)
	}
	if err := json.Unmarshal(intendedAfter, &afterObject); err != nil {
		return ImpactReport{}, fmt.Errorf("decode agent-network guardrail impact intended state: %w", err)
	}
	return ImpactReport{
		Classification:    "agent_network_guardrail_change",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"updating an agent-network guardrail can alter model and prompt admission for attached policies; affected peers and account resources require capability-aware live analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "agent_network_guardrail_update_requires_capability_analysis"},
	}, nil
}

func AgentNetworkGuardrailDeleteImpact(before []byte) (ImpactReport, error) {
	var object map[string]any
	if err := json.Unmarshal(before, &object); err != nil {
		return ImpactReport{}, fmt.Errorf("decode agent-network guardrail delete preimage: %w", err)
	}
	return ImpactReport{
		Classification:    "agent_network_guardrail_delete",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"deleting an agent-network guardrail can remove model and prompt admission controls from attached policies; affected peers and account resources require capability-aware live analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "agent_network_guardrail_delete_requires_capability_analysis"},
	}, nil
}

func RouteUpdateImpact(before, intendedAfter []byte) (ImpactReport, error) {
	var beforeObject, afterObject map[string]any
	if err := json.Unmarshal(before, &beforeObject); err != nil {
		return ImpactReport{}, fmt.Errorf("decode route impact preimage: %w", err)
	}
	if err := json.Unmarshal(intendedAfter, &afterObject); err != nil {
		return ImpactReport{}, fmt.Errorf("decode route impact intended state: %w", err)
	}
	changed := changedKeys(beforeObject, afterObject)
	if len(changed) == 0 || (len(changed) == 1 && changed[0] == "description") {
		return ImpactReport{
			Classification:    "metadata_only",
			Reachability:      "unchanged",
			AffectedPeers:     []string{},
			AffectedResources: []string{},
			Confidence:        "high",
			Evidence:          []string{"route description changed without changing routing behavior"},
			Completeness:      map[string]any{"state": "complete", "reason": nil},
		}, nil
	}
	return ImpactReport{
		Classification:    "route_change",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{fmt.Sprintf("route update changes routing fields: %v; affected peers and resources require live topology analysis", changed)},
		Completeness:      map[string]any{"state": "unknown", "reason": "route_change_requires_topology"},
	}, nil
}

func RouteDeleteImpact(before []byte) (ImpactReport, error) {
	var object map[string]any
	if err := json.Unmarshal(before, &object); err != nil {
		return ImpactReport{}, fmt.Errorf("decode route delete preimage: %w", err)
	}
	return ImpactReport{
		Classification:    "route_delete",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"deleting a route can change network reachability; affected peers and resources require live topology analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "route_delete_requires_topology"},
	}, nil
}

func RouteCreateImpact(intendedAfter []byte) (ImpactReport, error) {
	var object map[string]any
	if err := json.Unmarshal(intendedAfter, &object); err != nil {
		return ImpactReport{}, fmt.Errorf("decode route create intent: %w", err)
	}
	return ImpactReport{
		Classification:    "route_create",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"creating a route can add a reachable destination or alter path selection; affected peers and resources require live topology analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "route_create_requires_topology"},
	}, nil
}

func PeerUpdateImpact(before, intendedAfter []byte) (ImpactReport, error) {
	var beforeObject, afterObject map[string]any
	if err := json.Unmarshal(before, &beforeObject); err != nil {
		return ImpactReport{}, fmt.Errorf("decode peer impact preimage: %w", err)
	}
	if err := json.Unmarshal(intendedAfter, &afterObject); err != nil {
		return ImpactReport{}, fmt.Errorf("decode peer impact intended state: %w", err)
	}
	changed := changedKeys(beforeObject, afterObject)
	if len(changed) == 0 || (len(changed) == 1 && changed[0] == "name") {
		return ImpactReport{
			Classification:    "metadata_only",
			Reachability:      "unchanged",
			AffectedPeers:     []string{},
			AffectedResources: []string{},
			Confidence:        "high",
			Evidence:          []string{"peer name changed without changing peer access or connectivity state"},
			Completeness:      map[string]any{"state": "complete", "reason": nil},
		}, nil
	}
	return ImpactReport{
		Classification:    "peer_change",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{fmt.Sprintf("peer update changes access or connectivity fields: %v; affected peers and resources require live topology analysis", changed)},
		Completeness:      map[string]any{"state": "unknown", "reason": "peer_change_requires_topology"},
	}, nil
}

func PeerDeleteImpact(before []byte) (ImpactReport, error) {
	var object map[string]any
	if err := json.Unmarshal(before, &object); err != nil {
		return ImpactReport{}, fmt.Errorf("decode peer delete preimage: %w", err)
	}
	return ImpactReport{
		Classification:    "peer_delete",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"deleting a peer removes an access and connectivity principal; affected peers and resources require live topology analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "peer_delete_requires_topology"},
	}, nil
}

func NetworkUpdateImpact(before, intendedAfter []byte) (ImpactReport, error) {
	var beforeObject, afterObject map[string]any
	if err := json.Unmarshal(before, &beforeObject); err != nil {
		return ImpactReport{}, fmt.Errorf("decode network impact preimage: %w", err)
	}
	if err := json.Unmarshal(intendedAfter, &afterObject); err != nil {
		return ImpactReport{}, fmt.Errorf("decode network impact intended state: %w", err)
	}
	changed := changedKeys(beforeObject, afterObject)
	metadataOnly := true
	for _, key := range changed {
		if key != "name" && key != "description" {
			metadataOnly = false
			break
		}
	}
	if metadataOnly {
		return ImpactReport{
			Classification:    "metadata_only",
			Reachability:      "unchanged",
			AffectedPeers:     []string{},
			AffectedResources: []string{},
			Confidence:        "high",
			Evidence:          []string{"network metadata changed without changing attached policies, resources, or routers"},
			Completeness:      map[string]any{"state": "complete", "reason": nil},
		}, nil
	}
	return ImpactReport{
		Classification:    "network_change",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{fmt.Sprintf("network update changes topology fields: %v; affected peers and resources require live topology analysis", changed)},
		Completeness:      map[string]any{"state": "unknown", "reason": "network_change_requires_topology"},
	}, nil
}

func NetworkDeleteImpact(before []byte) (ImpactReport, error) {
	var object map[string]any
	if err := json.Unmarshal(before, &object); err != nil {
		return ImpactReport{}, fmt.Errorf("decode network delete preimage: %w", err)
	}
	return ImpactReport{
		Classification:    "network_delete",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"deleting a network can remove attached policies, resources, and routers; affected peers and resources require live topology analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "network_delete_requires_topology"},
	}, nil
}

func NetworkCreateImpact(intendedAfter []byte) (ImpactReport, error) {
	var object map[string]any
	if err := json.Unmarshal(intendedAfter, &object); err != nil {
		return ImpactReport{}, fmt.Errorf("decode network create intent: %w", err)
	}
	return ImpactReport{
		Classification:    "network_create",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"creating a network can add a topology scope for policies, resources, and routers; affected peers and resources require live topology analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "network_create_requires_topology"},
	}, nil
}

func NetworkResourceDeleteImpact(before []byte) (ImpactReport, error) {
	var object map[string]any
	if err := json.Unmarshal(before, &object); err != nil {
		return ImpactReport{}, fmt.Errorf("decode network resource delete preimage: %w", err)
	}
	return ImpactReport{
		Classification:    "network_resource_delete",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"deleting a network resource can remove a reachable destination; affected peers and routes require live topology analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "network_resource_delete_requires_topology"},
	}, nil
}

func NetworkResourceUpdateImpact(before, intendedAfter []byte) (ImpactReport, error) {
	var beforeObject, afterObject map[string]any
	if err := json.Unmarshal(before, &beforeObject); err != nil {
		return ImpactReport{}, fmt.Errorf("decode network resource impact preimage: %w", err)
	}
	if err := json.Unmarshal(intendedAfter, &afterObject); err != nil {
		return ImpactReport{}, fmt.Errorf("decode network resource impact intended state: %w", err)
	}
	changed := changedKeys(beforeObject, afterObject)
	metadataOnly := true
	for _, key := range changed {
		if key != "name" && key != "description" {
			metadataOnly = false
			break
		}
	}
	if metadataOnly {
		return ImpactReport{
			Classification:    "metadata_only",
			Reachability:      "unchanged",
			AffectedPeers:     []string{},
			AffectedResources: []string{},
			Confidence:        "high",
			Evidence:          []string{"network resource metadata changed without changing its address, enablement, or group assignment"},
			Completeness:      map[string]any{"state": "complete", "reason": nil},
		}, nil
	}
	return ImpactReport{
		Classification:    "network_resource_change",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{fmt.Sprintf("network resource update changes topology fields: %v; affected peers and routes require live topology analysis", changed)},
		Completeness:      map[string]any{"state": "unknown", "reason": "network_resource_change_requires_topology"},
	}, nil
}

func NetworkResourceCreateImpact(intendedAfter []byte) (ImpactReport, error) {
	var object map[string]any
	if err := json.Unmarshal(intendedAfter, &object); err != nil {
		return ImpactReport{}, fmt.Errorf("decode network resource create intent: %w", err)
	}
	return ImpactReport{
		Classification:    "network_resource_create",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"creating a network resource can add a reachable destination; affected peers and routes require live topology analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "network_resource_create_requires_topology"},
	}, nil
}

func NetworkRouterDeleteImpact(before []byte) (ImpactReport, error) {
	var object map[string]any
	if err := json.Unmarshal(before, &object); err != nil {
		return ImpactReport{}, fmt.Errorf("decode network router delete preimage: %w", err)
	}
	return ImpactReport{
		Classification:    "network_router_delete",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"deleting a network router can change route availability and peer reachability; affected peers and resources require live topology analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "network_router_delete_requires_topology"},
	}, nil
}

func NetworkRouterUpdateImpact(before, intendedAfter []byte) (ImpactReport, error) {
	var beforeObject, afterObject map[string]any
	if err := json.Unmarshal(before, &beforeObject); err != nil {
		return ImpactReport{}, fmt.Errorf("decode network router impact preimage: %w", err)
	}
	if err := json.Unmarshal(intendedAfter, &afterObject); err != nil {
		return ImpactReport{}, fmt.Errorf("decode network router impact intended state: %w", err)
	}
	changed := changedKeys(beforeObject, afterObject)
	if len(changed) == 0 {
		return ImpactReport{
			Classification:    "metadata_only",
			Reachability:      "unchanged",
			AffectedPeers:     []string{},
			AffectedResources: []string{},
			Confidence:        "high",
			Evidence:          []string{"network router state is unchanged"},
			Completeness:      map[string]any{"state": "complete", "reason": nil},
		}, nil
	}
	return ImpactReport{
		Classification:    "network_router_change",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{fmt.Sprintf("network router update changes reachability fields: %v; affected peers and resources require live topology analysis", changed)},
		Completeness:      map[string]any{"state": "unknown", "reason": "network_router_change_requires_topology"},
	}, nil
}

func NetworkRouterCreateImpact(intendedAfter []byte) (ImpactReport, error) {
	var object map[string]any
	if err := json.Unmarshal(intendedAfter, &object); err != nil {
		return ImpactReport{}, fmt.Errorf("decode network router create intent: %w", err)
	}
	return ImpactReport{
		Classification:    "network_router_create",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"creating a network router can add route availability and peer reachability; affected peers and resources require live topology analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "network_router_create_requires_topology"},
	}, nil
}

func GroupUpdateImpact(before, intendedAfter []byte) (ImpactReport, error) {
	var beforeObject, afterObject map[string]any
	if err := json.Unmarshal(before, &beforeObject); err != nil {
		return ImpactReport{}, fmt.Errorf("decode group impact preimage: %w", err)
	}
	if err := json.Unmarshal(intendedAfter, &afterObject); err != nil {
		return ImpactReport{}, fmt.Errorf("decode group impact intended state: %w", err)
	}

	changed := changedKeys(beforeObject, afterObject)
	if len(changed) == 0 || (len(changed) == 1 && changed[0] == "name") {
		return ImpactReport{
			Classification:    "metadata_only",
			Reachability:      "unchanged",
			AffectedPeers:     []string{},
			AffectedResources: []string{},
			Confidence:        "high",
			Evidence:          []string{"groups.update changes group metadata only; no membership or policy edge changed"},
			Completeness:      map[string]any{"state": "complete", "reason": nil},
		}, nil
	}
	return ImpactReport{
		Classification:    "unknown",
		Reachability:      "unknown",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "low",
		Evidence:          []string{fmt.Sprintf("group update changes unsupported fields: %v", changed)},
		Completeness:      map[string]any{"state": "unknown", "reason": "unsupported_group_fields"},
	}, nil
}

func GroupDeleteImpact(before []byte) (ImpactReport, error) {
	var object map[string]any
	if err := json.Unmarshal(before, &object); err != nil {
		return ImpactReport{}, fmt.Errorf("decode group delete preimage: %w", err)
	}
	return ImpactReport{
		Classification:    "group_delete",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"deleting a group can change policy membership and peer access; affected peers and resources require live topology analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "group_delete_requires_topology"},
	}, nil
}

func GroupCreateImpact(intendedAfter []byte) (ImpactReport, error) {
	var object map[string]any
	if err := json.Unmarshal(intendedAfter, &object); err != nil {
		return ImpactReport{}, fmt.Errorf("decode group create intent: %w", err)
	}
	return ImpactReport{
		Classification:    "group_create",
		Reachability:      "potentially_changed",
		AffectedPeers:     []string{},
		AffectedResources: []string{},
		Confidence:        "medium",
		Evidence:          []string{"creating a group can create a new policy or resource membership principal; affected peers and resources require live topology analysis"},
		Completeness:      map[string]any{"state": "unknown", "reason": "group_create_requires_topology"},
	}, nil
}

func changedKeys(before, after map[string]any) []string {
	keys := make(map[string]struct{}, len(before)+len(after))
	for key := range before {
		keys[key] = struct{}{}
	}
	for key := range after {
		keys[key] = struct{}{}
	}
	changed := make([]string, 0)
	for key := range keys {
		left, leftOK := before[key]
		right, rightOK := after[key]
		if !leftOK || !rightOK || !sameJSONValue(left, right) {
			changed = append(changed, key)
		}
	}
	sort.Strings(changed)
	return changed
}

func sameJSONValue(left, right any) bool {
	leftJSON, _ := json.Marshal(left)
	rightJSON, _ := json.Marshal(right)
	return string(leftJSON) == string(rightJSON)
}
