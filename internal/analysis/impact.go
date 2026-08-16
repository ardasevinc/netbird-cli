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
