// Example: Building a custom CLI with visionspec
//
// This example demonstrates how organizations can build their own
// CLI tools that include visionspec commands alongside custom commands.
//
// Templates and rubrics are compiled into the binary using //go:embed,
// so the CLI can be distributed as a single executable.
//
// Build with:
//
//	go build -o org-spec ./examples/org-cli
//
// Usage:
//
//	org-spec init my-project         # Uses visionspec init with org templates
//	org-spec lint                    # Uses visionspec lint
//	org-spec eval prd                # Uses org rubrics for evaluation
//	org-spec policy list             # Custom organization command
package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/ProductBuildersHQ/visionspec/pkg/cli"
	"github.com/ProductBuildersHQ/visionspec/pkg/rubrics"
	"github.com/ProductBuildersHQ/visionspec/pkg/templates"
	"github.com/spf13/cobra"
)

// Embed organization templates into the binary.
// These override or extend the default visionspec templates.
//
//go:embed templates/*.md
var orgTemplates embed.FS

// Embed organization rubrics into the binary.
// These override or extend the default visionspec rubrics.
//
//go:embed rubrics/*.rubric.yaml
var orgRubrics embed.FS

func main() {
	root := &cobra.Command{
		Use:   "org-spec",
		Short: "Organization specification management",
		Long: `org-spec is a customized visionspec CLI for Acme Corp.

It includes all standard visionspec commands plus organization-specific
features like policy management and custom templates/rubrics compiled
into a single binary.`,
		Version: "1.0.0",
	}

	// Configure visionspec with organization-specific loaders
	cfg := cli.DefaultConfig()
	cfg.Version = "1.0.0-acme"

	// Use embedded org templates, falling back to visionspec defaults.
	// Templates are compiled into the binary - no external files needed.
	cfg.TemplateLoader = templates.NewChainLoader(
		templates.NewEmbedFSLoader(orgTemplates, "templates"),
		templates.EmbeddedLoader(), // Fallback to visionspec defaults
	)

	// Use embedded org rubrics, falling back to visionspec defaults.
	// Rubrics are compiled into the binary - no external files needed.
	cfg.RubricLoader = rubrics.NewChainLoader(
		rubrics.NewEmbedFSLoader(orgRubrics, "rubrics"),
		rubrics.EmbeddedLoader(), // Fallback to visionspec defaults
	)

	// Add all visionspec commands
	cli.AddCommandsTo(root, cfg)

	// Add organization-specific commands
	root.AddCommand(policyCmd())

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

// policyCmd creates organization-specific policy commands.
func policyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "policy",
		Short: "Manage organization policies",
		Long:  `Commands for managing organization-specific policies and compliance.`,
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List active policies",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Active policies:")
			fmt.Println("  - security: SAST required for all PRDs")
			fmt.Println("  - compliance: SOC2 controls in TRD")
			fmt.Println("  - accessibility: WCAG 2.1 AA in UXD")
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "apply",
		Short: "Apply policies to current project",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Applying organization policies...")
			fmt.Println("✓ Added security requirements to prd.md")
			fmt.Println("✓ Added compliance requirements to trd.md")
			fmt.Println("✓ Added accessibility requirements to uxd.md")
			return nil
		},
	})

	return cmd
}
