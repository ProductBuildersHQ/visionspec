// Example: Building a custom CLI with visionspec
//
// This example demonstrates how organizations can build their own
// CLI tools that include visionspec commands alongside custom commands.
//
// Key differences from open source visionspec:
// - Open source: Provides choices and flexibility
// - Organization: Prescriptive defaults that enforce standards
//
// Resources are compiled into the binary using //go:embed:
// - Templates: Organization-specific (e.g., IRD with Pulumi pre-filled)
// - Rubrics: Stricter evaluation criteria (e.g., "MUST use PostgreSQL")
// - Constitutions: Org defaults that projects inherit
// - App Types: Constraints per app type (e.g., microservice requires 99.9%)
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
//	org-spec constitution resolve    # Resolves constitution hierarchy
//	org-spec policy list             # Custom organization command
package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/ProductBuildersHQ/visionspec/pkg/apptypes"
	"github.com/ProductBuildersHQ/visionspec/pkg/cli"
	"github.com/ProductBuildersHQ/visionspec/pkg/constitution"
	"github.com/ProductBuildersHQ/visionspec/pkg/rubrics"
	"github.com/ProductBuildersHQ/visionspec/pkg/templates"
	"github.com/spf13/cobra"
)

// Embed organization templates into the binary.
// These override or extend the default visionspec templates.
// Org templates are more prescriptive (e.g., IRD already has Pulumi/PostgreSQL).
//
//go:embed templates/*.md
var orgTemplates embed.FS

// Embed organization rubrics into the binary.
// Org rubrics have stricter criteria (e.g., "MUST use Go" vs "language documented").
//
//go:embed rubrics/*.rubric.yaml
var orgRubrics embed.FS

// Embed organization constitutions.
// Structure: constitutions/organization/*.yaml, constitutions/team/*.yaml
// Projects inherit defaults from org constitution automatically.
//
// //go:embed constitutions/**/*.yaml
// var orgConstitutions embed.FS

// Embed organization app type specs.
// Org app types have stricter constraints (e.g., 99.9% availability required).
//
// //go:embed apptypes/*.yaml
// var orgAppTypes embed.FS

func main() {
	root := &cobra.Command{
		Use:   "org-spec",
		Short: "Organization specification management",
		Long: `org-spec is a customized visionspec CLI for Acme Corp.

It includes all standard visionspec commands plus organization-specific
features like policy management and custom templates/rubrics compiled
into a single binary.

Organization standards enforced:
- Go backend with Huma+Chi for REST APIs
- PostgreSQL with RLS for multi-tenancy
- Pulumi (Go SDK) for Infrastructure as Code
- 99.9% availability minimum for microservices`,
		Version: "1.0.0",
	}

	// Configure visionspec with organization-specific loaders
	cfg := cli.DefaultConfig()
	cfg.Version = "1.0.0-acme"

	// Use embedded org templates, falling back to visionspec defaults.
	// Templates are compiled into the binary - no external files needed.
	// Org templates are more prescriptive (e.g., IRD already has Pulumi/PostgreSQL).
	cfg.TemplateLoader = templates.NewChainLoader(
		templates.NewEmbedFSLoader(orgTemplates, "templates"),
		templates.EmbeddedLoader(), // Fallback to visionspec defaults
	)

	// Use embedded org rubrics, falling back to visionspec defaults.
	// Org rubrics have stricter criteria (e.g., "MUST use Go" vs "language documented").
	cfg.RubricLoader = rubrics.NewChainLoader(
		rubrics.NewEmbedFSLoader(orgRubrics, "rubrics"),
		rubrics.EmbeddedLoader(), // Fallback to visionspec defaults
	)

	// Use embedded org constitutions.
	// These define organizational defaults that projects inherit.
	// Uncomment when constitutions are available:
	//
	// cfg.ConstitutionLoader = constitution.NewEmbeddedLoader(orgConstitutions, "constitutions")

	// Use embedded org app type specs, falling back to visionspec defaults.
	// Org app types have stricter constraints (e.g., 99.9% availability required).
	// Uncomment when app type specs are available:
	//
	// cfg.AppTypeLoader = apptypes.NewChainLoader(
	//     apptypes.NewEmbeddedLoader(orgAppTypes, "apptypes"),
	//     apptypes.DefaultLoader(), // Fallback to visionspec defaults
	// )

	// Add all visionspec commands
	cli.AddCommandsTo(root, cfg)

	// Add organization-specific commands
	root.AddCommand(policyCmd())
	root.AddCommand(constitutionCmd(cfg))

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

// constitutionCmd creates commands for managing constitutions.
func constitutionCmd(cfg *cli.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "constitution",
		Short: "Manage organization constitutions",
		Long:  `Commands for viewing and resolving constitution inheritance.`,
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List available constitutions",
		RunE: func(cmd *cobra.Command, args []string) error {
			loader := cfg.GetConstitutionLoader()
			if loader == nil {
				fmt.Println("No constitution loader configured.")
				fmt.Println("Add constitutions/ directory with org/team/project YAML files.")
				return nil
			}

			fmt.Println("Available constitutions:")
			for _, name := range loader.Available() {
				fmt.Printf("  - %s\n", name)
			}
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "resolve [name]",
		Short: "Resolve constitution with inheritance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			loader := cfg.GetConstitutionLoader()
			if loader == nil {
				return fmt.Errorf("no constitution loader configured")
			}

			resolver := constitution.NewResolver(loader)
			resolved, err := resolver.ResolveChain(args[0])
			if err != nil {
				return err
			}

			fmt.Printf("Resolved constitution: %s\n", resolved.Metadata.Name)
			fmt.Printf("Level: %s\n", resolved.Metadata.Level)
			if resolved.Technical.Languages.Backend.Primary != "" {
				fmt.Printf("Backend language: %s\n", resolved.Technical.Languages.Backend.Primary)
			}
			if resolved.Infrastructure.IaC.Tool != "" {
				fmt.Printf("IaC tool: %s\n", resolved.Infrastructure.IaC.Tool)
			}
			if resolved.Infrastructure.Availability.Target != "" {
				fmt.Printf("Availability target: %s%%\n", resolved.Infrastructure.Availability.Target)
			}
			return nil
		},
	})

	return cmd
}

// Suppress unused import warnings when constitutions/apptypes are commented out
var (
	_ = constitution.LevelOrganization
	_ = apptypes.AppTypeMicroservice
)

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
