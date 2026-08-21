package coveragecheck

import (
	"encoding/json"
	"fmt"
)

type operation struct {
	Method         string `json:"method"`
	Path           string `json:"path"`
	Implementation string `json:"implementation"`
	Verification   string `json:"verification"`
}

type summary struct {
	Implemented        int `json:"implemented"`
	Classified         int `json:"classified"`
	Discovered         int `json:"discovered"`
	ContractVerified   int `json:"contract_verified"`
	DisposableVerified int `json:"disposable_verified"`
	LiveVerified       int `json:"live_verified"`
	UnverifiedLive     int `json:"unverified_live"`
}

type declaredSummary struct {
	Implemented        *int `json:"implemented"`
	Classified         *int `json:"classified"`
	Discovered         *int `json:"discovered"`
	ContractVerified   *int `json:"contract_verified"`
	DisposableVerified *int `json:"disposable_verified"`
	LiveVerified       *int `json:"live_verified"`
	UnverifiedLive     *int `json:"unverified_live"`
}

type inventory struct {
	Operations []operation `json:"operations"`
}

type manifest struct {
	Operations []operation     `json:"operations"`
	Summary    declaredSummary `json:"summary"`
}

func Validate(inventoryJSON, manifestJSON []byte) (int, error) {
	var source inventory
	if err := json.Unmarshal(inventoryJSON, &source); err != nil {
		return 0, fmt.Errorf("decode source inventory: %w", err)
	}
	var declared manifest
	if err := json.Unmarshal(manifestJSON, &declared); err != nil {
		return 0, fmt.Errorf("decode coverage manifest: %w", err)
	}
	declaredCounts, err := declared.Summary.value()
	if err != nil {
		return 0, err
	}

	want := make(map[string]struct{}, len(source.Operations))
	for _, item := range source.Operations {
		key := operationKey(item)
		if _, exists := want[key]; exists {
			return 0, fmt.Errorf("duplicate source operation %s", key)
		}
		want[key] = struct{}{}
	}
	got := make(map[string]operation, len(declared.Operations))
	for _, item := range declared.Operations {
		key := operationKey(item)
		if _, exists := got[key]; exists {
			return 0, fmt.Errorf("duplicate declared operation %s", key)
		}
		got[key] = item
	}
	for key := range want {
		item, ok := got[key]
		if !ok {
			return 0, fmt.Errorf("source operation missing from coverage manifest: %s", key)
		}
		if item.Implementation == "" {
			return 0, fmt.Errorf("source operation has no implementation state: %s", key)
		}
		if item.Verification == "" {
			return 0, fmt.Errorf("source operation has no verification state: %s", key)
		}
	}
	for key := range got {
		if _, ok := want[key]; !ok {
			return 0, fmt.Errorf("coverage manifest contains operation absent from pinned source: %s", key)
		}
	}

	actual := summary{}
	for _, item := range declared.Operations {
		switch item.Implementation {
		case "implemented":
			actual.Implemented++
		case "classified":
			actual.Classified++
		case "discovered":
			actual.Discovered++
		default:
			return 0, fmt.Errorf("operation %s has unknown implementation state %q", operationKey(item), item.Implementation)
		}
		switch item.Verification {
		case "contract_verified":
			actual.ContractVerified++
		case "disposable_verified":
			actual.DisposableVerified++
		case "live_verified":
			actual.LiveVerified++
		case "unverified_live":
			actual.UnverifiedLive++
		default:
			return 0, fmt.Errorf("operation %s has unknown verification state %q", operationKey(item), item.Verification)
		}
	}
	if actual != declaredCounts {
		return 0, fmt.Errorf("coverage summary does not match operation rows: summary=%+v rows=%+v", declaredCounts, actual)
	}
	return len(want), nil
}

func operationKey(item operation) string {
	return item.Method + " " + item.Path
}

func (declared declaredSummary) value() (summary, error) {
	fields := map[string]*int{
		"implemented": declared.Implemented, "classified": declared.Classified, "discovered": declared.Discovered,
		"contract_verified": declared.ContractVerified, "disposable_verified": declared.DisposableVerified,
		"live_verified": declared.LiveVerified, "unverified_live": declared.UnverifiedLive,
	}
	for name, value := range fields {
		if value == nil {
			return summary{}, fmt.Errorf("coverage summary is missing %s", name)
		}
	}
	return summary{
		Implemented: *declared.Implemented, Classified: *declared.Classified, Discovered: *declared.Discovered,
		ContractVerified: *declared.ContractVerified, DisposableVerified: *declared.DisposableVerified,
		LiveVerified: *declared.LiveVerified, UnverifiedLive: *declared.UnverifiedLive,
	}, nil
}
