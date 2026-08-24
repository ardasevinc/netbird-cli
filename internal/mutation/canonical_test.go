package mutation

import "testing"

func TestCanonicalJSONIgnoresObjectFormatting(t *testing.T) {
	equivalent, err := Equivalent([]byte(`{"b":2,"a":1}`), []byte(`{"a":1,"b":2}`))
	if err != nil || !equivalent {
		t.Fatalf("equivalent=%t err=%v", equivalent, err)
	}
}

func TestCanonicalJSONPreservesArrayOrder(t *testing.T) {
	equivalent, err := Equivalent([]byte(`{"items":[1,2]}`), []byte(`{"items":[2,1]}`))
	if err != nil || equivalent {
		t.Fatalf("equivalent=%t err=%v", equivalent, err)
	}
}

func TestEquivalentTreatsNullAndEmptyArrayAsSameCollection(t *testing.T) {
	equivalent, err := Equivalent([]byte(`{"resources":null,"routers":[]}`), []byte(`{"resources":[],"routers":null}`))
	if err != nil || !equivalent {
		t.Fatalf("equivalent=%t err=%v", equivalent, err)
	}
}

func TestEquivalentDoesNotTreatNullAndNonEmptyArrayAsSame(t *testing.T) {
	equivalent, err := Equivalent([]byte(`{"resources":null}`), []byte(`{"resources":[{"id":"r1"}]}`))
	if err != nil || equivalent {
		t.Fatalf("equivalent=%t err=%v", equivalent, err)
	}
}

func TestEquivalentUnorderedCollectionIgnoresTopLevelOrder(t *testing.T) {
	equivalent, err := EquivalentUnorderedCollection(
		[]byte(`[{"id":"user-1","groups":["a","b"]},{"id":"user-2"}]`),
		[]byte(`[{"id":"user-2"},{"groups":["a","b"],"id":"user-1"}]`),
	)
	if err != nil || !equivalent {
		t.Fatalf("equivalent=%t err=%v", equivalent, err)
	}
}

func TestEquivalentUnorderedCollectionPreservesNestedArrayOrder(t *testing.T) {
	equivalent, err := EquivalentUnorderedCollection(
		[]byte(`[{"id":"user-1","groups":["a","b"]}]`),
		[]byte(`[{"id":"user-1","groups":["b","a"]}]`),
	)
	if err != nil || equivalent {
		t.Fatalf("equivalent=%t err=%v", equivalent, err)
	}
}

func TestEquivalentUnorderedCollectionPreservesDuplicateItems(t *testing.T) {
	equivalent, err := EquivalentUnorderedCollection(
		[]byte(`[{"id":"user-1"},{"id":"user-1"}]`),
		[]byte(`[{"id":"user-1"}]`),
	)
	if err != nil || equivalent {
		t.Fatalf("equivalent=%t err=%v", equivalent, err)
	}
}

func TestEquivalentUnorderedCollectionRejectsChangedItem(t *testing.T) {
	equivalent, err := EquivalentUnorderedCollection(
		[]byte(`[{"id":"user-1","role":"admin"}]`),
		[]byte(`[{"id":"user-1","role":"auditor"}]`),
	)
	if err != nil || equivalent {
		t.Fatalf("equivalent=%t err=%v", equivalent, err)
	}
}

func TestNormalizedResultSchemaRecognizesWrappedStableResult(t *testing.T) {
	schema, normalized, err := NormalizedResultSchema([]byte(`{"schema":"nb/v1/groups-list-result","data":{"groups":[]}}`))
	if err != nil || !normalized || schema != "nb/v1/groups-list-result" {
		t.Fatalf("schema=%q normalized=%t err=%v", schema, normalized, err)
	}
}
