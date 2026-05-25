// Example: 1-N Growth CLI using profiles
//
// This example demonstrates using the built-in "growth" profile
// for the 1-N scaling phase.
//
// Build with:
//
//	go build -o growth-spec ./examples/1-n-growth
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
		Use:   "growth-spec",
		Short: "Spec management for 1-N growth phase",
		Long: `growth-spec uses the "growth" profile for metrics-driven development.

PRD, UXD, and FAQ are required. Focus on experiments and growth metrics.`,
		Version: "1.0.0",
	}

	// Load the growth profile
	profile, err := profiles.DefaultLoader().Load("growth")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading profile: %v\n", err)
		os.Exit(1)
	}

	// Create config from profile
	cfg := cli.ConfigFromProfile(profile)
	cfg.Version = "1.0.0-growth"

	// Add visionspec commands
	cli.AddCommandsTo(root, cfg)

	// Add growth-specific commands
	root.AddCommand(metricsCmd())
	root.AddCommand(experimentCmd())

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func metricsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "metrics",
		Short: "Show AARRR metrics framework",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("AARRR Pirate Metrics:")
			fmt.Println("  1. Acquisition")
			fmt.Println("  2. Activation")
			fmt.Println("  3. Retention")
			fmt.Println("  4. Revenue")
			fmt.Println("  5. Referral")
			return nil
		},
	}
}

func experimentCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "experiment",
		Short: "Experimentation framework",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Experiment Framework:")
			fmt.Println("  1. Hypothesis")
			fmt.Println("  2. Design (control/variant)")
			fmt.Println("  3. Implement (feature flags)")
			fmt.Println("  4. Analyze")
			fmt.Println("  5. Decide")
			return nil
		},
	}
}
