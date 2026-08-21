//go:build ignore

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ardasevinc/netbird-cli/internal/coveragecheck"
)

func main() {
	inventoryPath := flag.String("inventory", "coverage/sources/netbird-v0.77.0/openapi-inventory.json", "source inventory")
	manifestPath := flag.String("manifest", "coverage/manifest.json", "coverage manifest")
	flag.Parse()
	inventoryJSON, err := os.ReadFile(*inventoryPath)
	if err != nil {
		fail("read %s: %v", *inventoryPath, err)
	}
	manifestJSON, err := os.ReadFile(*manifestPath)
	if err != nil {
		fail("read %s: %v", *manifestPath, err)
	}
	count, err := coveragecheck.Validate(inventoryJSON, manifestJSON)
	if err != nil {
		fail("%v", err)
	}
	fmt.Printf("coverage source reconciled: %d operations\n", count)
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
