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
