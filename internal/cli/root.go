// Package cli implements the visionspec command-line interface.
//
// This package provides backward compatibility with the existing CLI.
// For composable CLI building, use pkg/cli instead.
package cli

import (
	"fmt"

	pkgcli "github.com/ProductBuildersHQ/visionspec/pkg/cli"
	"github.com/spf13/cobra"
)

var (
	// Version is set at build time.
	Version = "dev"

	// Commit is set at build time.
	Commit = "unknown"
)

var rootCmd *cobra.Command

func init() {
	rootCmd = &cobra.Command{
		Use:   "visionspec",
		Short: "Vision-first specification orchestration using Working Backwards",
		Long: `VisionSpec implements Amazon's Working Backwards methodology for
specification development. Start with the customer experience (Press Release),
derive requirements (PRD), then build technical specs (TRD, IRD).

It provides:
  - Working Backwards flow (MRD → Press → FAQ → PRD)
  - Technical synthesis (TRD, IRD)
  - Structured evaluation with LLM judges
  - Reconciliation into unified execution specs
  - Export to SpecKit, GSD, GasTown, GasCity`,
		Version: fmt.Sprintf("%s (commit: %s)", Version, Commit),
	}

	// Global flags
	rootCmd.PersistentFlags().StringP("project", "p", "", "Project name or path")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().String("workflows-repo", "", "Path to spec-workflows repository for templates and rubrics")

	// Use pkg/cli for composable commands
	cfg := pkgcli.DefaultConfig()
	cfg.Version = Version

	// Apply workflows-repo flag to config if provided
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		workflowsRepo, _ := cmd.Flags().GetString("workflows-repo")
		if workflowsRepo != "" {
			cfg.WorkflowsRepoPath = workflowsRepo
		}
		return nil
	}

	pkgcli.AddCommandsTo(rootCmd, cfg)
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
