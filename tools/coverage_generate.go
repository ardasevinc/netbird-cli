//go:build ignore

package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type operation struct {
	ID             string `json:"id"`
	Method         string `json:"method"`
	Path           string `json:"path"`
	Implementation string `json:"implementation"`
	Verification   string `json:"verification"`
}

type manifest struct {
	Schema     string      `json:"schema"`
	Source     string      `json:"source_baseline"`
	Support    []string    `json:"support_window"`
	Operations []operation `json:"operations"`
	Summary    summary     `json:"summary"`
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

type inventory struct {
	Schema     string      `json:"schema"`
	Source     source      `json:"source"`
	Operations []operation `json:"operations"`
}

type source struct {
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
	Path       string `json:"path"`
	SHA256     string `json:"sha256"`
	PathCount  int    `json:"path_count"`
	OpCount    int    `json:"operation_count"`
}

var (
	pathPattern   = regexp.MustCompile(`^  (/[^:]+):\s*$`)
	methodPattern = regexp.MustCompile(`^    (get|post|put|patch|delete|head|options|trace):\s*$`)
)

func main() {
	sourcePath := flag.String("source", "", "path to a pinned NetBird OpenAPI YAML file")
	manifestPath := flag.String("manifest", "coverage/manifest.json", "coverage manifest to refresh")
	inventoryPath := flag.String("inventory", "coverage/sources/netbird-v0.77.0/openapi-inventory.json", "generated source inventory")
	tag := flag.String("tag", "v0.77.0", "source tag")
	flag.Parse()
	if *sourcePath == "" {
		fail("--source is required")
	}
	data, err := os.ReadFile(*sourcePath)
	if err != nil {
		fail("read source: %v", err)
	}
	ops := parseOperations(data)
	if len(ops) == 0 {
		fail("no OpenAPI operations found")
	}
	sha := sha256.Sum256(data)
	for i := range ops {
		ops[i].ID = ownedID(ops[i].Method, ops[i].Path)
	}
	applyReviewedOverrides(ops)
	if err := refreshManifest(*manifestPath, ops, *tag); err != nil {
		fail("refresh manifest: %v", err)
	}
	inv := inventory{
		Schema:     "nb/v1/coverage-source",
		Source:     source{Repository: "github.com/netbirdio/netbird", Tag: *tag, Path: "shared/management/http/api/openapi.yml", SHA256: hex.EncodeToString(sha[:]), PathCount: countPaths(data), OpCount: len(ops)},
		Operations: ops,
	}
	if err := writeJSON(*inventoryPath, inv); err != nil {
		fail("write inventory: %v", err)
	}
}

func applyReviewedOverrides(ops []operation) {
	reviewed := map[string]operation{
		"GET /api/accounts":            {ID: "accounts.list", Implementation: "implemented", Verification: "contract_verified"},
		"GET /api/peers":               {ID: "peers.list", Implementation: "implemented", Verification: "contract_verified"},
		"GET /api/peers/{peerId}":      {ID: "peers.get", Implementation: "implemented", Verification: "contract_verified"},
		"GET /api/policies":            {ID: "policies.list", Implementation: "implemented", Verification: "contract_verified"},
		"GET /api/policies/{policyId}": {ID: "policies.get", Implementation: "implemented", Verification: "contract_verified"},
		"GET /api/users":               {ID: "users.list", Implementation: "implemented", Verification: "contract_verified"},
		"GET /api/users/invites":       {ID: "users.invites", Implementation: "implemented", Verification: "contract_verified"},
	}
	for i, item := range ops {
		if override, ok := reviewed[item.Method+" "+item.Path]; ok {
			override.Method = item.Method
			override.Path = item.Path
			ops[i] = override
		}
	}
}

func parseOperations(data []byte) []operation {
	var result []operation
	var current string
	insidePaths := false
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "paths:" {
			insidePaths = true
			continue
		}
		if !insidePaths {
			continue
		}
		if match := pathPattern.FindStringSubmatch(line); match != nil {
			current = match[1]
			continue
		}
		if match := methodPattern.FindStringSubmatch(line); match != nil && current != "" {
			result = append(result, operation{Method: strings.ToUpper(match[1]), Path: current})
		}
	}
	return result
}

func countPaths(data []byte) int {
	count := 0
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		if pathPattern.MatchString(scanner.Text()) {
			count++
		}
	}
	return count
}

func ownedID(method, path string) string {
	clean := strings.NewReplacer("/", ".", "{", "", "}", "", "-", "_").Replace(path)
	clean = strings.Trim(clean, ".")
	return strings.ToLower(method) + "." + clean
}

func refreshManifest(path string, discovered []operation, tag string) error {
	var existing manifest
	if data, err := os.ReadFile(path); err == nil {
		if err := json.Unmarshal(data, &existing); err != nil {
			return err
		}
	}
	byKey := make(map[string]operation, len(existing.Operations))
	for _, item := range existing.Operations {
		item.Path = canonicalPath(item.Path)
		byKey[item.Method+" "+item.Path] = item
	}
	for _, item := range discovered {
		key := item.Method + " " + item.Path
		current, ok := byKey[key]
		if !ok {
			current = item
			current.Implementation = "discovered"
		}
		if item.Implementation != "" && item.Implementation != "discovered" && current.Implementation == "discovered" {
			current.ID = item.ID
			current.Implementation = item.Implementation
			current.Verification = item.Verification
		}
		if current.ID == "" {
			current.ID = item.ID
		}
		byKey[key] = current
	}
	merged := make([]operation, 0, len(byKey))
	for _, item := range byKey {
		merged = append(merged, item)
	}
	sort.Slice(merged, func(i, j int) bool {
		if merged[i].Path == merged[j].Path {
			return merged[i].Method < merged[j].Method
		}
		return merged[i].Path < merged[j].Path
	})
	existing.Schema = "nb/v1/coverage-manifest"
	existing.Source = "netbird-" + tag
	if len(existing.Support) == 0 {
		existing.Support = []string{"current", "previous-minor", "two-minors-back"}
	}
	existing.Operations = merged
	existing.Summary = summary{}
	for _, item := range merged {
		switch item.Implementation {
		case "implemented":
			existing.Summary.Implemented++
		case "classified":
			existing.Summary.Classified++
		default:
			existing.Summary.Discovered++
		}
		switch item.Verification {
		case "contract_verified":
			existing.Summary.ContractVerified++
		case "disposable_verified":
			existing.Summary.DisposableVerified++
		case "live_verified":
			existing.Summary.LiveVerified++
		case "unverified_live":
			existing.Summary.UnverifiedLive++
		}
	}
	return writeJSON(path, existing)
}

func canonicalPath(path string) string {
	if path == "/api/groups/{id}" {
		return "/api/groups/{groupId}"
	}
	return path
}

func writeJSON(path string, value any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
