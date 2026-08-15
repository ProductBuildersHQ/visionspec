// Package main is the entry point for the visionspec CLI.
package main

import (
	"os"

	"github.com/ProductBuildersHQ/visionspec/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
