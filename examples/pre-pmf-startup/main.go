// Example: Pre-PMF Startup CLI using profiles
//
// This example demonstrates using the built-in "startup" profile
// for pre-product-market-fit startups.
//
// Build with:
//
//	go build -o startup-spec ./examples/pre-pmf-startup
package main

import (
	"fmt"
	"os"

	"github.com/ProductBuildersHQ/visionspec/pkg/cli"
	"github.com/ProductBuildersHQ/visionspec/pkg/profiles"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "startup-spec",
		Short: "Lightweight spec management for pre-PMF startups",
		Long: `startup-spec uses the "startup" profile for early-stage startups.

Only PRD is required. Add more specs as you grow.`,
		Version: "1.0.0",
	}

	// Load the startup profile
	profile, err := profiles.DefaultLoader().Load("startup")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading profile: %v\n", err)
		os.Exit(1)
	}

	// Create config from profile
	cfg := cli.ConfigFromProfile(profile)
	cfg.Version = "1.0.0-startup"

	// Add visionspec commands
	cli.AddCommandsTo(root, cfg)

	// Add startup-specific commands
	root.AddCommand(velocityCmd())

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func velocityCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "velocity",
		Short: "Check documentation velocity",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Recommended: PRD only")
			fmt.Println("Add when needed: MRD, UXD, TRD")
			fmt.Println("Wait until post-PMF: GTM docs")
			return nil
		},
	}
}
