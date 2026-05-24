// Package cli implements the multispec command-line interface.
//
// This package provides backward compatibility with the existing CLI.
// For composable CLI building, use pkg/cli instead.
package cli

import (
	"fmt"

	pkgcli "github.com/plexusone/multispec/pkg/cli"
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
		Use:   "multispec",
		Short: "Multi-domain specification orchestration for humans and AI agents",
		Long: `MultiSpec bridges the gap between organizational intent (MRD, PRD, UXD)
and executable specifications for AI coding agents.

It provides:
  - Domain-specific authoring (source specs)
  - GTM synthesis (press releases, FAQs, narratives)
  - Technical synthesis (TRD, IRD)
  - Structured evaluation with LLM judges
  - Reconciliation into unified execution specs
  - Export to SpecKit, GSD, GasTown, GasCity, OpenSpec`,
		Version: fmt.Sprintf("%s (commit: %s)", Version, Commit),
	}

	// Global flags
	rootCmd.PersistentFlags().StringP("project", "p", "", "Project name or path")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")

	// Use pkg/cli for composable commands
	cfg := pkgcli.DefaultConfig()
	cfg.Version = Version
	pkgcli.AddCommandsTo(rootCmd, cfg)
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
