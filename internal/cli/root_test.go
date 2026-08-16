package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/version"
)

func TestExecuteVersionJSON(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute(context.Background(), []string{"--json", "version"}, &stdout, &stderr, version.Info{Version: "1.2.3", Commit: "abc", Date: "now"})
	if code != 0 {
		t.Fatalf("exit code = %d, stderr = %q", code, stderr.String())
	}
	var result struct {
		Schema string `json:"schema"`
		OK     bool   `json:"ok"`
		Data   struct {
			Version string `json:"version"`
		} `json:"data"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode result: %v", err)
	}
	if result.Schema != "nb/v1/version-result" || !result.OK || result.Data.Version != "1.2.3" {
		t.Fatalf("unexpected result: %s", stdout.String())
	}
}

func TestExecuteSchemaList(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute(context.Background(), []string{"schema", "list"}, &stdout, &stderr, version.Current())
	if code != 0 || stderr.Len() != 0 {
		t.Fatalf("code=%d stderr=%q", code, stderr.String())
	}
	if stdout.String() == "" {
		t.Fatal("schema list was empty")
	}
}
