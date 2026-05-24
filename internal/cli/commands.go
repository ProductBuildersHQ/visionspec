package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/plexusone/multispec/internal/mcp"
	"github.com/plexusone/multispec/pkg/config"
	"github.com/plexusone/multispec/pkg/eval"
	"github.com/plexusone/multispec/pkg/lint"
	"github.com/plexusone/multispec/pkg/reconcile"
	"github.com/plexusone/multispec/pkg/synth"
	"github.com/plexusone/multispec/pkg/target"
	"github.com/plexusone/multispec/pkg/types"
	"github.com/spf13/cobra"
)

var lintCmd = &cobra.Command{
	Use:   "lint [project]",
	Short: "Validate directory structure and naming conventions",
	Long: `Validate that the project follows multispec conventions:

  - Directory structure matches canonical layout
  - File naming follows conventions (lowercase specs, kebab-case projects)
  - Required specs are present
  - Config file is valid

Examples:
  multispec lint                    # Lint all projects
  multispec lint user-onboarding    # Lint specific project
  multispec lint --format json      # Output as JSON`,
	RunE: runLint,
}

func runLint(cmd *cobra.Command, args []string) error {
	format, _ := cmd.Flags().GetString("format")
	ci, _ := cmd.Flags().GetBool("ci")

	// Find specs directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	specsDir := filepath.Join(cwd, config.SpecsDir)

	linter := lint.New(specsDir)

	var result *lint.Result

	if len(args) > 0 {
		// Lint specific project
		projectName := args[0]
		projectPath := filepath.Join(specsDir, projectName)

		if _, err := os.Stat(projectPath); os.IsNotExist(err) {
			return fmt.Errorf("project not found: %s", projectName)
		}

		result, err = linter.LintProject(projectName, projectPath)
	} else {
		// Lint all projects
		result, err = linter.LintAll()
	}

	if err != nil {
		return fmt.Errorf("linting: %w", err)
	}

	// Output result
	switch format {
	case "json":
		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling result: %w", err)
		}
		fmt.Println(string(data))
	default:
		fmt.Print(result.FormatText())
	}

	// Exit with error code for CI
	if ci && !result.Passed {
		os.Exit(1)
	}

	return nil
}

var evalCmd = &cobra.Command{
	Use:   "eval [spec-type]",
	Short: "Evaluate specs using LLM judges",
	Long: `Evaluate specification documents using LLM-as-a-Judge.

Examples:
  multispec eval prd          # Evaluate PRD
  multispec eval --all        # Evaluate all specs
  multispec eval --source     # Evaluate source specs only
  multispec eval --gtm        # Evaluate GTM docs only`,
	RunE: runEval,
}

func runEval(cmd *cobra.Command, args []string) error {
	allFlag, _ := cmd.Flags().GetBool("all")
	sourceFlag, _ := cmd.Flags().GetBool("source")
	gtmFlag, _ := cmd.Flags().GetBool("gtm")
	technicalFlag, _ := cmd.Flags().GetBool("technical")

	// Find project root
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	projectPath, err := config.FindProjectRoot(cwd)
	if err != nil {
		return fmt.Errorf("not in a multispec project (no multispec.yaml found)")
	}

	// Load project config for LLM settings
	project, err := config.Load(projectPath)
	if err != nil {
		return fmt.Errorf("loading project config: %w", err)
	}

	// Determine which specs to evaluate
	var specTypes []types.SpecType

	if len(args) > 0 {
		// Evaluate specific spec type
		specType := types.SpecType(args[0])
		if !specType.IsValid() {
			return fmt.Errorf("invalid spec type: %s", args[0])
		}
		specTypes = append(specTypes, specType)
	} else if allFlag {
		specTypes = types.AllSpecTypes()
	} else if sourceFlag {
		specTypes = types.SourceSpecTypes()
	} else if gtmFlag {
		specTypes = types.GTMSpecTypes()
	} else if technicalFlag {
		specTypes = types.TechnicalSpecTypes()
	} else {
		return fmt.Errorf("specify a spec type or use --all, --source, --gtm, or --technical")
	}

	// Create LLM client
	llmClient, err := eval.NewLLMClientFromProject(project.LLM)
	if err != nil {
		return fmt.Errorf("initializing LLM: %w", err)
	}
	defer func() { _ = llmClient.Close() }()

	evaluator := eval.NewEvaluator(llmClient)
	ctx := context.Background()

	// Evaluate each spec
	for _, specType := range specTypes {
		specPath := config.SpecPath(projectPath, specType)
		content, err := os.ReadFile(specPath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("⊘ %s: not found, skipping\n", specType)
				continue
			}
			return fmt.Errorf("reading %s: %w", specType, err)
		}

		fmt.Printf("⋯ Evaluating %s...\n", specType)

		result, err := evaluator.Evaluate(ctx, specType, string(content))
		if err != nil {
			fmt.Printf("✗ %s: evaluation failed: %v\n", specType, err)
			continue
		}

		// Write eval result to file
		evalPath := config.EvalPath(projectPath, specType)
		evalDir := filepath.Dir(evalPath)
		if err := os.MkdirAll(evalDir, 0755); err != nil {
			return fmt.Errorf("creating eval directory: %w", err)
		}

		evalData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling eval result: %w", err)
		}

		if err := os.WriteFile(evalPath, evalData, 0600); err != nil { //nolint:gosec // G703: evalPath from config.EvalPath with validated project path
			return fmt.Errorf("writing eval file: %w", err)
		}

		// Print summary
		if result.Passed {
			fmt.Printf("✓ %s: %.1f/10 PASS (%d findings)\n", specType, result.Score, len(result.Findings))
		} else {
			fmt.Printf("✗ %s: %.1f/10 FAIL (%d findings)\n", specType, result.Score, len(result.Findings))
		}
	}

	return nil
}

var synthesizeCmd = &cobra.Command{
	Use:   "synthesize <type>",
	Short: "Generate specs from source documents",
	Long: `Generate specification documents from source specs.

GTM synthesis (Working Backwards):
  multispec synthesize press        # MRD + PRD → press.md
  multispec synthesize faq          # press.md → faq.md
  multispec synthesize narrative-1p # MRD + PRD → narrative-1p.md
  multispec synthesize narrative-6p # MRD + PRD + UXD → narrative-6p.md

Technical synthesis:
  multispec synthesize trd          # MRD + PRD + UXD + CONSTITUTION → trd.md
  multispec synthesize ird          # TRD + CONSTITUTION → ird.md`,
	Args: cobra.ExactArgs(1),
	RunE: runSynthesize,
}

func runSynthesize(cmd *cobra.Command, args []string) error {
	specTypeArg := args[0]
	evalFlag, _ := cmd.Flags().GetBool("eval")

	// Parse spec type
	specType := types.SpecType(specTypeArg)
	if !synth.CanSynthesize(specType) {
		return fmt.Errorf("cannot synthesize %s (valid: trd, ird, press, faq, narrative-1p, narrative-6p)", specTypeArg)
	}

	// Find project root
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	projectPath, err := config.FindProjectRoot(cwd)
	if err != nil {
		return fmt.Errorf("not in a multispec project (no multispec.yaml found)")
	}

	// Load project config
	project, err := config.Load(projectPath)
	if err != nil {
		return fmt.Errorf("loading project config: %w", err)
	}

	// Check required sources exist
	requiredSources := synth.RequiredSources(specType)
	for _, srcType := range requiredSources {
		srcPath := config.SpecPath(projectPath, srcType)
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			return fmt.Errorf("missing required source spec: %s", srcType)
		}
	}

	// Load source specs
	input := synth.SynthesisInput{}
	if content, err := os.ReadFile(config.SpecPath(projectPath, types.SpecTypeMRD)); err == nil {
		input.MRD = string(content)
	}
	if content, err := os.ReadFile(config.SpecPath(projectPath, types.SpecTypePRD)); err == nil {
		input.PRD = string(content)
	}
	if content, err := os.ReadFile(config.SpecPath(projectPath, types.SpecTypeUXD)); err == nil {
		input.UXD = string(content)
	}
	if content, err := os.ReadFile(config.SpecPath(projectPath, types.SpecTypeTRD)); err == nil {
		input.TRD = string(content)
	}
	if content, err := os.ReadFile(config.SpecPath(projectPath, types.SpecTypePress)); err == nil {
		input.Press = string(content)
	}

	// Load constitution if exists
	constitutionPath := filepath.Join(projectPath, "..", "CONSTITUTION.md")
	if content, err := os.ReadFile(constitutionPath); err == nil {
		input.Constitution = string(content)
	}

	// Create LLM client
	llmClient, err := eval.NewLLMClientFromProject(project.LLM)
	if err != nil {
		return fmt.Errorf("initializing LLM: %w", err)
	}
	defer func() { _ = llmClient.Close() }()

	// Create synthesizer
	synthesizer := synth.NewSynthesizer(&cliSynthLLMAdapter{client: llmClient})

	fmt.Printf("⋯ Synthesizing %s from %v...\n", specType, requiredSources)

	ctx := context.Background()
	result, err := synthesizer.Synthesize(ctx, specType, input)
	if err != nil {
		return fmt.Errorf("synthesis failed: %w", err)
	}

	// Write output
	outputPath := config.SpecPath(projectPath, specType)
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}
	if err := os.WriteFile(outputPath, []byte(result.Content), 0600); err != nil {
		return fmt.Errorf("writing output: %w", err)
	}

	fmt.Printf("✓ Generated %s\n", outputPath)

	// Run evaluation if requested
	if evalFlag {
		fmt.Printf("⋯ Evaluating %s...\n", specType)
		evaluator := eval.NewEvaluator(llmClient)
		evalResult, err := evaluator.Evaluate(ctx, specType, result.Content)
		if err != nil {
			fmt.Printf("✗ Evaluation failed: %v\n", err)
		} else {
			evalPath := config.EvalPath(projectPath, specType)
			evalDir := filepath.Dir(evalPath)
			if err := os.MkdirAll(evalDir, 0755); err != nil {
				fmt.Printf("⚠ Failed to create eval directory: %v\n", err)
			} else if evalData, err := json.MarshalIndent(evalResult, "", "  "); err != nil {
				fmt.Printf("⚠ Failed to marshal eval result: %v\n", err)
			} else if err := os.WriteFile(evalPath, evalData, 0600); err != nil {
				fmt.Printf("⚠ Failed to write eval file: %v\n", err)
			}
			if evalResult.Passed {
				fmt.Printf("✓ %s: %.1f/10 PASS\n", specType, evalResult.Score)
			} else {
				fmt.Printf("✗ %s: %.1f/10 FAIL\n", specType, evalResult.Score)
			}
		}
	}

	return nil
}

// cliSynthLLMAdapter adapts eval.LLMClient to synth.LLMClient interface.
type cliSynthLLMAdapter struct {
	client *eval.LLMClient
}

func (a *cliSynthLLMAdapter) Complete(ctx context.Context, prompt string) (string, error) {
	content, _, err := a.client.Complete(ctx, prompt)
	return content, err
}

var reconcileCmd = &cobra.Command{
	Use:   "reconcile",
	Short: "Generate unified execution spec from approved specs",
	Long: `Reconcile all approved specifications into a unified execution spec.

This command:
  1. Loads all approved source, GTM, and technical specs
  2. Detects conflicts and missing traceability
  3. Generates spec.md (unified execution spec)
  4. Generates spec.eval.json (reconciliation evaluation)`,
	RunE: runReconcile,
}

func runReconcile(cmd *cobra.Command, args []string) error {
	// Find project root
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	projectPath, err := config.FindProjectRoot(cwd)
	if err != nil {
		return fmt.Errorf("not in a multispec project (no multispec.yaml found)")
	}

	// Load project config
	project, err := config.Load(projectPath)
	if err != nil {
		return fmt.Errorf("loading project config: %w", err)
	}

	// Check approvals
	approved, missing := reconcile.CheckApprovals(project.Approvals)
	if len(missing) > 0 {
		fmt.Println("Missing approvals:")
		for _, m := range missing {
			fmt.Printf("  ✗ %s\n", m)
		}
		fmt.Println("\nApprove specs with: multispec approve <spec-type>")
		return fmt.Errorf("cannot reconcile without required approvals")
	}

	fmt.Printf("✓ All required specs approved: %v\n", approved)

	// Load all approved specs
	input := reconcile.ReconcileInput{
		ProjectName: project.Name,
	}
	if content, err := os.ReadFile(config.SpecPath(projectPath, types.SpecTypeMRD)); err == nil {
		input.MRD = string(content)
	}
	if content, err := os.ReadFile(config.SpecPath(projectPath, types.SpecTypePRD)); err == nil {
		input.PRD = string(content)
	}
	if content, err := os.ReadFile(config.SpecPath(projectPath, types.SpecTypeUXD)); err == nil {
		input.UXD = string(content)
	}
	if content, err := os.ReadFile(config.SpecPath(projectPath, types.SpecTypeTRD)); err == nil {
		input.TRD = string(content)
	}
	if content, err := os.ReadFile(config.SpecPath(projectPath, types.SpecTypeIRD)); err == nil {
		input.IRD = string(content)
	}

	// Load constitution if exists
	constitutionPath := filepath.Join(projectPath, "..", "CONSTITUTION.md")
	if content, err := os.ReadFile(constitutionPath); err == nil {
		input.Constitution = string(content)
	}

	// Create LLM client
	llmClient, err := eval.NewLLMClientFromProject(project.LLM)
	if err != nil {
		return fmt.Errorf("initializing LLM: %w", err)
	}
	defer func() { _ = llmClient.Close() }()

	// Create reconciler
	reconciler := reconcile.NewReconciler(&cliReconcileLLMAdapter{client: llmClient})

	fmt.Println("⋯ Reconciling specs...")

	ctx := context.Background()
	result, err := reconciler.Reconcile(ctx, input)
	if err != nil {
		return fmt.Errorf("reconciliation failed: %w", err)
	}

	// Write output
	outputPath := config.SpecPath(projectPath, types.SpecTypeSpec)
	if err := os.WriteFile(outputPath, []byte(result.Content), 0600); err != nil {
		return fmt.Errorf("writing spec.md: %w", err)
	}

	fmt.Printf("✓ Generated %s\n", outputPath)
	fmt.Printf("  Sources: %v\n", result.Sources)

	if len(result.Conflicts) > 0 {
		fmt.Printf("  Conflicts: %d (see spec.md Decision Log)\n", len(result.Conflicts))
	}

	return nil
}

// cliReconcileLLMAdapter adapts eval.LLMClient to reconcile.LLMClient interface.
type cliReconcileLLMAdapter struct {
	client *eval.LLMClient
}

func (a *cliReconcileLLMAdapter) Complete(ctx context.Context, prompt string) (string, error) {
	content, _, err := a.client.Complete(ctx, prompt)
	return content, err
}

var approveCmd = &cobra.Command{
	Use:   "approve <spec-type>",
	Short: "Approve a spec for reconciliation",
	Long: `Mark a specification as approved.

Examples:
  multispec approve prd                    # Approve PRD
  multispec approve trd --approver=eng@co  # Approve with approver`,
	Args: cobra.ExactArgs(1),
	RunE: runApprove,
}

func runApprove(cmd *cobra.Command, args []string) error {
	specTypeArg := args[0]
	approver, _ := cmd.Flags().GetString("approver")
	comment, _ := cmd.Flags().GetString("comment")

	// Parse spec type
	specType := types.SpecType(specTypeArg)
	if !specType.IsValid() {
		return fmt.Errorf("invalid spec type: %s", specTypeArg)
	}

	// Find project root
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	projectPath, err := config.FindProjectRoot(cwd)
	if err != nil {
		return fmt.Errorf("not in a multispec project (no multispec.yaml found)")
	}

	// Load project config
	project, err := config.Load(projectPath)
	if err != nil {
		return fmt.Errorf("loading project config: %w", err)
	}

	// Check that spec exists
	specPath := config.SpecPath(projectPath, specType)
	if _, err := os.Stat(specPath); os.IsNotExist(err) {
		return fmt.Errorf("spec not found: %s - cannot approve non-existent spec", specType)
	}

	// Initialize approvals map if needed
	if project.Approvals == nil {
		project.Approvals = make(map[types.SpecType]*types.Approval)
	}

	// Determine approver
	if approver == "" {
		// Try to get from git config or environment
		approver = os.Getenv("USER")
		if approver == "" {
			approver = "unknown"
		}
	}

	// Record approval
	project.Approvals[specType] = &types.Approval{
		Approver:   approver,
		ApprovedAt: time.Now(),
		Comment:    comment,
	}
	project.UpdatedAt = time.Now()

	// Save project config
	if err := config.Save(project); err != nil {
		return fmt.Errorf("saving approval: %w", err)
	}

	fmt.Printf("✓ Approved %s by %s\n", specType, approver)
	if comment != "" {
		fmt.Printf("  Comment: %s\n", comment)
	}

	return nil
}

var exportCmd = &cobra.Command{
	Use:   "export <target>",
	Short: "Export specs to target execution systems",
	Long: `Export the reconciled spec to downstream execution systems.

Targets:
  speckit   - GitHub Spec-Kit format
  gsd       - Get Shit Done format (not yet implemented)
  gastown   - GasTown formula/beads (not yet implemented)
  gascity   - GasCity city.toml (not yet implemented)
  openspec  - OpenSpec portable format (not yet implemented)

Examples:
  multispec export speckit`,
	Args: cobra.ExactArgs(1),
	RunE: runExport,
}

func runExport(cmd *cobra.Command, args []string) error {
	targetName := args[0]

	// Find project root
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	projectPath, err := config.FindProjectRoot(cwd)
	if err != nil {
		return fmt.Errorf("not in a multispec project (no multispec.yaml found)")
	}

	// Load project config
	project, err := config.Load(projectPath)
	if err != nil {
		return fmt.Errorf("loading project config: %w", err)
	}

	// Read spec.md
	specPath := config.SpecPath(projectPath, types.SpecTypeSpec)
	specContent, err := os.ReadFile(specPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("spec.md not found - run 'multispec reconcile' first")
		}
		return fmt.Errorf("reading spec.md: %w", err)
	}

	// Get target adapter
	t, err := target.Get(targetName)
	if err != nil {
		available := target.Available()
		return fmt.Errorf("unknown target: %s (available: %v)", targetName, available)
	}

	// Get export config
	exportConfig := target.ProjectTargetConfig(project, targetName)
	if exportConfig.OutputDir == "" {
		exportConfig.OutputDir = filepath.Join(projectPath, "export", targetName)
	}

	fmt.Printf("⋯ Exporting to %s...\n", targetName)

	// Run export
	result, err := t.Export(string(specContent), *exportConfig)
	if err != nil {
		return fmt.Errorf("export failed: %w", err)
	}

	fmt.Printf("✓ %s\n", result.Message)
	fmt.Printf("  Output: %s\n", result.OutputDir)
	fmt.Println("  Files:")
	for _, f := range result.Files {
		fmt.Printf("    - %s\n", f)
	}

	return nil
}

var targetsCmd = &cobra.Command{
	Use:   "targets",
	Short: "List available export targets",
	Long:  `List all available export targets and their capabilities.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Available targets:")
		fmt.Println()
		fmt.Println("  speckit   GitHub Spec-Kit format (spec.md, plan.md, tasks.md)")
		fmt.Println("  gsd       Get Shit Done format (PLAN.md, STATE.md)")
		fmt.Println("  gastown   GasTown formulas and beads")
		fmt.Println("  gascity   GasCity city.toml configuration")
		fmt.Println("  openspec  OpenSpec portable format (future)")
		return nil
	},
}

var graphCmd = &cobra.Command{
	Use:   "graph <subcommand>",
	Short: "Manage requirement graphs via graphize",
	Long: `Manage requirement graphs using graphize integration.

Subcommands:
  extract   Build graph from specs
  query     Query graph relationships
  export    Export to HTML/JSON/GraphML`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("graph: not yet implemented")
		return nil
	},
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start MCP server for AI assistant integration",
	Long: `Start a Model Context Protocol (MCP) server for integration
with AI coding assistants like Claude Code and Kiro CLI.

The MCP server provides tools for:
  - Listing projects and their status
  - Querying spec content and evaluations
  - Running synthesis and reconciliation
  - Exporting to targets

Configuration for Claude Code (~/.claude/claude_desktop_config.json):
  {
    "mcpServers": {
      "multispec": {
        "command": "multispec",
        "args": ["serve"]
      }
    }
  }`,
	RunE: runServe,
}

func runServe(cmd *cobra.Command, args []string) error {
	server := mcp.NewServer()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	return server.Serve(ctx)
}

func init() {
	// lint flags
	lintCmd.Flags().String("format", "text", "Output format: text, json")
	lintCmd.Flags().Bool("ci", false, "Exit with non-zero code if lint fails")

	// eval flags
	evalCmd.Flags().Bool("all", false, "Evaluate all specs")
	evalCmd.Flags().Bool("source", false, "Evaluate source specs only")
	evalCmd.Flags().Bool("gtm", false, "Evaluate GTM docs only")
	evalCmd.Flags().Bool("technical", false, "Evaluate technical docs only")

	// synthesize flags
	synthesizeCmd.Flags().Bool("eval", false, "Run evaluation after synthesis")

	// approve flags
	approveCmd.Flags().String("approver", "", "Approver email or identifier")
	approveCmd.Flags().String("comment", "", "Approval comment")

	// export flags
	exportCmd.Flags().Bool("dry-run", false, "Show what would be exported without writing")
	exportCmd.Flags().String("output", "", "Output directory (default: target-specific)")

	// serve flags
	serveCmd.Flags().Int("port", 0, "HTTP port (0 for stdio transport)")
	serveCmd.Flags().String("transport", "stdio", "Transport: stdio, http, sse")
}
