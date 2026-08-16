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
