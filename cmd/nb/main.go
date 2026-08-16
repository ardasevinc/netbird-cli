package main

import (
	"context"
	"os"

	"github.com/ardasevinc/netbird-cli/internal/cli"
	"github.com/ardasevinc/netbird-cli/internal/version"
)

func main() {
	os.Exit(cli.Execute(context.Background(), os.Args[1:], os.Stdout, os.Stderr, version.Current()))
}
