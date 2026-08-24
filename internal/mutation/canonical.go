package mutation

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

func CanonicalJSON(data []byte) ([]byte, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, fmt.Errorf("decode JSON: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("JSON contains more than one value")
		}
		return nil, fmt.Errorf("decode trailing JSON: %w", err)
	}
	encoded, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("encode canonical JSON: %w", err)
	}
	return encoded, nil
}

func Digest(parts ...[]byte) (string, error) {
	hash := sha256.New()
	for _, part := range parts {
		canonical, err := CanonicalJSON(part)
		if err != nil {
			return "", err
		}
		_, _ = hash.Write([]byte{0})
		_, _ = hash.Write(canonical)
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func Equivalent(left, right []byte) (bool, error) {
	leftValue, err := decodeCanonicalValue(left)
	if err != nil {
		return false, err
	}
	rightValue, err := decodeCanonicalValue(right)
	if err != nil {
		return false, err
	}
	return equivalentValue(leftValue, rightValue), nil
}

// EquivalentUnorderedCollection compares top-level JSON collections without
// assigning meaning to server response order. Objects and nested arrays keep
// the normal canonical comparison rules, and duplicate items are preserved.
func EquivalentUnorderedCollection(left, right []byte) (bool, error) {
	leftValue, err := decodeCanonicalValue(left)
	if err != nil {
		return false, err
	}
	rightValue, err := decodeCanonicalValue(right)
	if err != nil {
		return false, err
	}

	leftCollection, leftOK := topLevelCollection(leftValue)
	rightCollection, rightOK := topLevelCollection(rightValue)
	if !leftOK || !rightOK {
		return Equivalent(left, right)
	}
	if len(leftCollection) != len(rightCollection) {
		return false, nil
	}

	leftItems, err := canonicalCollectionItems(leftCollection)
	if err != nil {
		return false, err
	}
	rightItems, err := canonicalCollectionItems(rightCollection)
	if err != nil {
		return false, err
	}
	for i := range leftItems {
		if !bytes.Equal(leftItems[i], rightItems[i]) {
			return false, nil
		}
	}
	return true, nil
}

// NormalizedResultSchema reports the stable nb result schema when value is a
// wrapped command result rather than an exact upstream payload. Such results
// are useful for display, but cannot serve as mutation preimages because the
// normalized model may omit upstream fields.
func NormalizedResultSchema(value []byte) (string, bool, error) {
	var object map[string]json.RawMessage
	if err := json.Unmarshal(value, &object); err != nil {
		if strings.HasPrefix(strings.TrimSpace(string(value)), "[") || strings.TrimSpace(string(value)) == "null" {
			return "", false, nil
		}
		return "", false, err
	}
	var schema string
	rawSchema, hasSchema := object["schema"]
	if !hasSchema {
		return "", false, nil
	}
	if err := json.Unmarshal(rawSchema, &schema); err != nil {
		return "", false, fmt.Errorf("decode result schema: %w", err)
	}
	if !strings.HasPrefix(schema, "nb/v1/") || len(object["data"]) == 0 {
		return "", false, nil
	}
	return schema, true, nil
}

func topLevelCollection(value any) ([]any, bool) {
	if value == nil {
		return []any{}, true
	}
	collection, ok := value.([]any)
	return collection, ok
}

func canonicalCollectionItems(collection []any) ([][]byte, error) {
	items := make([][]byte, 0, len(collection))
	for _, item := range collection {
		encoded, err := json.Marshal(item)
		if err != nil {
			return nil, fmt.Errorf("encode collection item: %w", err)
		}
		items = append(items, encoded)
	}
	sort.Slice(items, func(i, j int) bool {
		return bytes.Compare(items[i], items[j]) < 0
	})
	return items, nil
}

func decodeCanonicalValue(data []byte) (any, error) {
	canonical, err := CanonicalJSON(data)
	if err != nil {
		return nil, err
	}
	var value any
	decoder := json.NewDecoder(bytes.NewReader(canonical))
	decoder.UseNumber()
	if err := decoder.Decode(&value); err != nil {
		return nil, fmt.Errorf("decode canonical JSON: %w", err)
	}
	return value, nil
}

func equivalentValue(left, right any) bool {
	if left == nil {
		if values, ok := right.([]any); ok {
			return len(values) == 0
		}
		return right == nil
	}
	if right == nil {
		if values, ok := left.([]any); ok {
			return len(values) == 0
		}
		return false
	}

	leftObject, leftIsObject := left.(map[string]any)
	rightObject, rightIsObject := right.(map[string]any)
	if leftIsObject || rightIsObject {
		if !leftIsObject || !rightIsObject || len(leftObject) != len(rightObject) {
			return false
		}
		for key, leftValue := range leftObject {
			rightValue, ok := rightObject[key]
			if !ok || !equivalentValue(leftValue, rightValue) {
				return false
			}
		}
		return true
	}

	leftArray, leftIsArray := left.([]any)
	rightArray, rightIsArray := right.([]any)
	if leftIsArray || rightIsArray {
		if !leftIsArray || !rightIsArray || len(leftArray) != len(rightArray) {
			return false
		}
		for index := range leftArray {
			if !equivalentValue(leftArray[index], rightArray[index]) {
				return false
			}
		}
		return true
	}

	return bytes.Equal(mustCanonicalScalar(left), mustCanonicalScalar(right))
}

func mustCanonicalScalar(value any) []byte {
	encoded, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	return encoded
}
