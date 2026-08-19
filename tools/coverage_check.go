//go:build ignore

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type operation struct {
	Method         string `json:"method"`
	Path           string `json:"path"`
	Implementation string `json:"implementation"`
}

type summary struct {
	Implemented int `json:"implemented"`
	Classified  int `json:"classified"`
	Discovered  int `json:"discovered"`
}

type inventory struct {
	Operations []operation `json:"operations"`
}

type manifest struct {
	Operations []operation `json:"operations"`
	Summary    summary     `json:"summary"`
}

func main() {
	inventoryPath := flag.String("inventory", "coverage/sources/netbird-v0.77.0/openapi-inventory.json", "source inventory")
	manifestPath := flag.String("manifest", "coverage/manifest.json", "coverage manifest")
	flag.Parse()
	var source inventory
	var declared manifest
	readJSON(*inventoryPath, &source)
	readJSON(*manifestPath, &declared)
	want := make(map[string]struct{}, len(source.Operations))
	for _, item := range source.Operations {
		want[item.Method+" "+item.Path] = struct{}{}
	}
	got := make(map[string]operation, len(declared.Operations))
	for _, item := range declared.Operations {
		key := item.Method + " " + item.Path
		if _, exists := got[key]; exists {
			fail("duplicate declared operation %s", key)
		}
		got[key] = item
	}
	for key := range want {
		item, ok := got[key]
		if !ok {
			fail("source operation missing from coverage manifest: %s", key)
		}
		if item.Implementation == "" {
			fail("source operation has no implementation state: %s", key)
		}
	}
	for key := range got {
		if _, ok := want[key]; !ok {
			fail("coverage manifest contains operation absent from pinned source: %s", key)
		}
	}
	actual := summary{}
	for _, item := range declared.Operations {
		switch item.Implementation {
		case "implemented":
			actual.Implemented++
		case "classified":
			actual.Classified++
		default:
			actual.Discovered++
		}
	}
	if actual != declared.Summary {
		fail("coverage summary does not match operation rows: summary=%+v rows=%+v", declared.Summary, actual)
	}
	fmt.Printf("coverage source reconciled: %d operations\n", len(want))
}

func readJSON(path string, value any) {
	data, err := os.ReadFile(path)
	if err != nil {
		fail("read %s: %v", path, err)
	}
	if err := json.Unmarshal(data, value); err != nil {
		fail("decode %s: %v", path, err)
	}
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
