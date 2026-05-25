// Example: 0-1 Product CLI using profiles
//
// This example demonstrates using the built-in "0-1" profile
// with custom commands for the idea validation phase.
//
// Build with:
//
//	go build -o zero-to-one ./examples/0-1-product
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
		Use:   "zero-to-one",
		Short: "Minimal spec management for 0-1 product development",
		Long: `zero-to-one is a minimalist visionspec CLI for the 0-1 phase.

Going from nothing to something? You need:
- A clear hypothesis to test
- Just enough structure to communicate
- Maximum speed, minimum overhead

This CLI uses the "0-1" profile with a single "hypothesis" document.`,
		Version: "1.0.0",
	}

	// Load the 0-1 profile
	profile, err := profiles.DefaultLoader().Load("0-1")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading profile: %v\n", err)
		os.Exit(1)
	}

	// Create config from profile
	cfg := cli.ConfigFromProfile(profile)
	cfg.Version = "1.0.0-zero-to-one"

	// Add visionspec commands
	cli.AddCommandsTo(root, cfg)

	// Add 0-1 specific commands
	root.AddCommand(hypothesisCmd())
	root.AddCommand(validateCmd())
	root.AddCommand(pivotCmd())

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func hypothesisCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "hypothesis",
		Short: "Show hypothesis template",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Hypothesis Template:")
			fmt.Println()
			fmt.Println("  We believe that [target users]")
			fmt.Println("  have a problem with [problem]")
			fmt.Println()
			fmt.Println("  We believe that [solution]")
			fmt.Println("  will solve this problem")
			fmt.Println()
			fmt.Println("  We will know we're right when [success metric]")
			return nil
		},
	}
}

func validateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Check if your hypothesis is testable",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Hypothesis Validation Checklist:")
			fmt.Println()
			fmt.Println("  [ ] Target users are specific")
			fmt.Println("  [ ] Problem is concrete and observable")
			fmt.Println("  [ ] Solution is buildable in days")
			fmt.Println("  [ ] Success metric is measurable")
			fmt.Println("  [ ] You can test with <100 users")
			return nil
		},
	}
}

func pivotCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pivot",
		Short: "Guidance on when and how to pivot",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Pivot when:")
			fmt.Println("  - Metric not improving after 3 iterations")
			fmt.Println("  - Users won't pay/use daily")
			fmt.Println("  - Market too small")
			return nil
		},
	}
}
