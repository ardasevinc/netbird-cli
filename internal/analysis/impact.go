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
