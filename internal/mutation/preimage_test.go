package mutation

import "testing"

func TestClassifyPreimage(t *testing.T) {
	state, err := ClassifyPreimage([]byte(`{"name":"old"}`), []byte(`{"name":"old"}`), []byte(`{"name":"new"}`))
	if err != nil || state != PreimageMatches {
		t.Fatalf("state=%s err=%v", state, err)
	}
	state, err = ClassifyPreimage([]byte(`{"name":"old"}`), []byte(`{"name":"other"}`), []byte(`{"name":"new"}`))
	if err != nil || state != PreimageDrifted {
		t.Fatalf("state=%s err=%v", state, err)
	}
	state, err = ClassifyPreimage([]byte(`{"name":"old"}`), []byte(`{"name":"new"}`), []byte(`{"name":"new"}`))
	if err != nil || state != PreimageAlreadySatisfied {
		t.Fatalf("state=%s err=%v", state, err)
	}
}
