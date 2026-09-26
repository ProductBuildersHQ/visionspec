package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ProductBuildersHQ/visionspec/pkg/speckit"
	sws "github.com/ProductBuildersHQ/visionspec/pkg/workflows"
)

// speckitCmd creates the speckit command group.
func speckitCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "speckit",
		Short: "Export workflows as GitHub Spec Kit plugins",
		Long: `Export visionspec workflows as GitHub Spec Kit plugins.

Spec Kit (github.com/github/spec-kit) supports third-party extensions,
workflows, and bundles. These commands package visionspec workflow families
in those official formats.`,
	}
	cmd.AddCommand(speckitExportCmd(cfg))
	return cmd
}

// speckitExportCmd creates the speckit export command.
func speckitExportCmd(cfg *Config) *cobra.Command {
	var (
		outputDir     string
		check         bool
		productID     string
		engineeringID string
		extVersion    string
	)

	cmd := &cobra.Command{
		Use:   "export [workflow]",
		Short: "Export a workflow family as Spec Kit extensions",
		Long: `Export a workflow family as GitHub Spec Kit extension packages.

The family is split into two extensions by spec category: a product
extension (source and gtm specs — the requirements chain) and an
engineering add-on (technical specs). Each spec type becomes a Spec Kit
command that authors the document from the family's template and
self-evaluates against the family's rubric.

The output under <output-dir>/extensions/ is generated content: regenerate
it with this command rather than editing it, and use --check in CI to catch
drift between the committed copy and the generator.

Examples:
  # Export the default family into ./speckit
  visionspec speckit export

  # Verify committed artifacts are in sync (CI drift gate)
  visionspec speckit export --check

  # Export another family with a custom extension id
  visionspec speckit export big-tech -o ./speckit --product-id big-tech`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workflowName := "aws-one-way-door"
			if len(args) == 1 {
				workflowName = args[0]
			}

			loader := cfg.WorkflowLoader
			if loader == nil {
				loader = sws.DefaultLoader()
			}
			lw, err := loader.Load(workflowName)
			if err != nil {
				return fmt.Errorf("workflow %q not found: %w", workflowName, err)
			}

			opts := speckit.Options{
				ProductID:     productID,
				EngineeringID: engineeringID,
				Version:       extVersion,
			}

			files, err := speckit.BuildExtensions(lw, opts)
			if err != nil {
				return fmt.Errorf("building spec-kit extensions for %q: %w", workflowName, err)
			}

			wf, err := speckit.BuildWorkflow(lw, opts)
			if err != nil {
				return fmt.Errorf("building spec-kit workflow for %q: %w", workflowName, err)
			}
			files = append(files, wf)

			if check {
				drift, err := speckit.Check(outputDir, files)
				if err != nil {
					return fmt.Errorf("checking spec-kit artifacts: %w", err)
				}
				if len(drift) > 0 {
					for _, p := range drift {
						fmt.Fprintf(cmd.ErrOrStderr(), "drift: %s\n", p)
					}
					return fmt.Errorf("%d spec-kit artifact(s) out of sync; run: visionspec speckit export %s -o %s",
						len(drift), workflowName, outputDir)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Spec Kit artifacts in %s are in sync (%d files)\n", outputDir, len(files))
				return nil
			}

			if err := speckit.WriteFiles(outputDir, files); err != nil {
				return fmt.Errorf("writing spec-kit artifacts: %w", err)
			}
			for _, f := range files {
				fmt.Fprintf(cmd.OutOrStdout(), "Created %s/%s\n", outputDir, f.Path)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&outputDir, "output", "o", "speckit", "Output directory for generated artifacts")
	cmd.Flags().BoolVar(&check, "check", false, "Verify committed artifacts match the generator instead of writing")
	cmd.Flags().StringVar(&productID, "product-id", "", "Extension id for the product chain (default: family name without aws- prefix)")
	cmd.Flags().StringVar(&engineeringID, "engineering-id", "", "Extension id for the engineering add-on (default: <product-id>-engineering)")
	cmd.Flags().StringVar(&extVersion, "ext-version", "", "Extension version X.Y.Z (default: 0.1.0)")

	return cmd
}
