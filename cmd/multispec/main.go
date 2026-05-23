// Package main is the entry point for the multispec CLI.
package main

import (
	"os"

	"github.com/plexusone/multispec/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
