package catalog

import (
	"encoding/json"
	"fmt"
)

type Operation struct {
	ID             string `json:"id"`
	Method         string `json:"method"`
	Path           string `json:"path"`
	Availability   string `json:"availability,omitempty"`
	Implementation string `json:"implementation"`
	Verification   string `json:"verification"`
}

type coverageDocument struct {
	Operations []Operation `json:"operations"`
}

func ReadOperation(id string) (Operation, error) {
	manifest, err := CoverageManifest()
	if err != nil {
		return Operation{}, err
	}
	var document coverageDocument
	if err := json.Unmarshal(manifest, &document); err != nil {
		return Operation{}, fmt.Errorf("decode coverage manifest: %w", err)
	}
	for _, operation := range document.Operations {
		if operation.ID == id {
			if operation.Method != "GET" {
				return Operation{}, fmt.Errorf("operation %q is %s, only GET operations are admitted by api get", id, operation.Method)
			}
			return operation, nil
		}
	}
	return Operation{}, fmt.Errorf("operation %q is not in the embedded coverage manifest", id)
}
