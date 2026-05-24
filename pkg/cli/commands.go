package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/plexusone/multispec/internal/mcp"
	"github.com/plexusone/multispec/pkg/config"
	"github.com/plexusone/multispec/pkg/eval"
	"github.com/plexusone/multispec/pkg/lint"
	"github.com/plexusone/multispec/pkg/profiles"
	"github.com/plexusone/multispec/pkg/reconcile"
	"github.com/plexusone/multispec/pkg/specgraph"
	"github.com/plexusone/multispec/pkg/status"
	"github.com/plexusone/multispec/pkg/synth"
	"github.com/plexusone/multispec/pkg/target"
	"github.com/plexusone/multispec/pkg/templates"
	"github.com/plexusone/multispec/pkg/types"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// initCmd creates the init command.
func initCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init <project-name>",
		Short: "Initialize a new multispec project",
		Long: `Initialize a new multispec project with the canonical directory structure.

The project name must be kebab-case (lowercase with hyphens).

Profiles:
  --profile 0-1         Minimal for idea validation
  --profile startup     Lightweight for pre-PMF startups
  --profile growth      Metrics-driven for 1-N scaling
  --profile enterprise  Comprehensive for post-PMF enterprises

Creates:
  docs/specs/<project>/
  ├── source/          # Human-authored specs (mrd, prd, uxd)
  ├── gtm/             # LLM-generated GTM docs (press, faq, narrative)
  ├── technical/       # LLM-generated technical docs (trd, ird)
  ├── eval/            # Evaluation results
  └── multispec.yaml   # Project configuration`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit(cmd, args, cfg)
		},
	}

	cmd.Flags().String("constitution", "", "Path to constitution file (relative or absolute)")
	cmd.Flags().Bool("with-templates", false, "Create template spec files")
	cmd.Flags().String("profile", "", "Configuration profile (0-1, startup, growth, enterprise)")

	return cmd
}

var kebabCaseRegex = regexp.MustCompile(`^[a-z][a-z0-9]*(-[a-z0-9]+)*$`)

func runInit(cmd *cobra.Command, args []string, cfg *Config) error {
	projectName := args[0]

	// Validate project name is kebab-case
	if !kebabCaseRegex.MatchString(projectName) {
		return fmt.Errorf("invalid project name %q: must be kebab-case (e.g., 'user-onboarding')", projectName)
	}

	// Load profile if specified
	profileName, _ := cmd.Flags().GetString("profile")
	if profileName != "" {
		loader := cfg.ProfileLoader
		if loader == nil {
			loader = profiles.DefaultLoader()
		}

		profile, err := loader.Load(profileName)
		if err != nil {
			return fmt.Errorf("loading profile %q: %w", profileName, err)
		}

		// Apply profile settings to config
		if profile.SpecConfig != nil {
			cfg.SpecConfig = profile.SpecConfig
		}
		if profile.TemplateLoader != nil {
			cfg.TemplateLoader = templates.NewChainLoader(profile.TemplateLoader, cfg.TemplateLoader)
		}
		if profile.RubricLoader != nil {
			cfg.RubricLoader = profile.RubricLoader
		}

		fmt.Printf("Using profile: %s\n", profile.Name)
		fmt.Printf("  %s\n\n", profile.Description)
	}

	// Find or create specs directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	specsDir, err := config.FindSpecsDir(cwd)
	if err != nil {
		// Create docs/specs if it doesn't exist
		specsDir = filepath.Join(cwd, config.SpecsDir)
		if err := os.MkdirAll(specsDir, 0755); err != nil {
			return fmt.Errorf("failed to create specs directory: %w", err)
		}
		fmt.Printf("Created specs directory: %s\n", specsDir)
	}

	projectPath := filepath.Join(specsDir, projectName)

	// Check if project already exists
	if _, err := os.Stat(projectPath); err == nil {
		return fmt.Errorf("project %q already exists at %s", projectName, projectPath)
	}

	// Create project directories
	dirs := []string{
		projectPath,
		filepath.Join(projectPath, config.SourceDir),
		filepath.Join(projectPath, config.GTMDir),
		filepath.Join(projectPath, config.TechnicalDir),
		filepath.Join(projectPath, config.EvalDir),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create project config
	constitution, _ := cmd.Flags().GetString("constitution")
	if constitution == "" {
		// Default to repo-level constitution
		constitution = fmt.Sprintf("../%s", config.ConstitutionFile)
	}

	project := &types.Project{
		Name:         projectName,
		Path:         projectPath,
		Constitution: constitution,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Targets: types.TargetConfig{
			Default: "speckit",
			SpecKit: &types.SpecKitConfig{
				Enabled:         true,
				BranchNumbering: "sequential",
			},
		},
	}

	// Save config
	configPath := filepath.Join(projectPath, config.ConfigFileName)
	data, err := yaml.Marshal(project)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Add header comment
	header := "# multispec project configuration\n# See: https://github.com/plexusone/multispec\n\n"
	if err := os.WriteFile(configPath, []byte(header+string(data)), 0600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	// Create template files if requested
	withTemplates, _ := cmd.Flags().GetBool("with-templates")
	if withTemplates {
		if err := createTemplateFiles(projectPath, cfg); err != nil {
			return fmt.Errorf("failed to create template files: %w", err)
		}
	}

	// Print summary
	fmt.Printf("\n✅ Created multispec project: %s\n\n", projectName)
	fmt.Println("Directory structure:")
	fmt.Printf("  %s/\n", projectName)
	fmt.Println("  ├── source/        # Human-authored specs")
	fmt.Println("  ├── gtm/           # LLM-generated GTM docs")
	fmt.Println("  ├── technical/     # LLM-generated technical docs")
	fmt.Println("  ├── eval/          # Evaluation results")
	fmt.Println("  └── multispec.yaml # Project configuration")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Create source specs: mrd.md, prd.md, uxd.md in source/")
	fmt.Println("  2. Run evaluations: multispec eval --all")
	fmt.Println("  3. Generate GTM docs: multispec synthesize press")
	fmt.Println("  4. Generate technical docs: multispec synthesize trd")
	fmt.Println("  5. Reconcile: multispec reconcile")

	return nil
}

func createTemplateFiles(projectPath string, cfg *Config) error {
	loader := cfg.TemplateLoader
	if loader == nil {
		loader = templates.DefaultLoader()
	}

	specConfig := cfg.GetSpecConfig()

	// Create templates for all required source specs
	for _, specName := range specConfig.RequiredSpecs() {
		category := specConfig.GetCategory(specName)
		if category != types.CategorySource {
			continue // Only create source templates
		}

		// Get the template name (may be aliased)
		templateName := specConfig.GetTemplate(specName)
		specType := types.SpecType(templateName)

		tmpl, err := loader.Load(specType)
		if err != nil {
			// Try loading by spec name if template name didn't work
			tmpl, err = loader.Load(types.SpecType(specName))
			if err != nil {
				fmt.Printf("⚠ Template not found for %s, skipping\n", specName)
				continue
			}
		}

		// Determine output directory based on category
		dir := config.SourceDir
		path := filepath.Join(projectPath, dir, specName+".md")

		if err := os.WriteFile(path, []byte(strings.TrimSpace(tmpl.Content)+"\n"), 0600); err != nil {
			return err
		}
		fmt.Printf("Created template: %s\n", filepath.Base(path))
	}

	return nil
}

// lintCmd creates the lint command.
func lintCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLint(cmd, args, cfg)
		},
	}

	cmd.Flags().String("format", "text", "Output format: text, json")
	cmd.Flags().Bool("ci", false, "Exit with non-zero code if lint fails")

	return cmd
}

func runLint(cmd *cobra.Command, args []string, cfg *Config) error {
	format, _ := cmd.Flags().GetString("format")
	ci, _ := cmd.Flags().GetBool("ci")

	// Find specs directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	specsDir := filepath.Join(cwd, config.SpecsDir)

	// Get SpecConfig from CLI config
	specConfig := cfg.GetSpecConfig()

	linter := lint.NewWithConfig(specsDir, specConfig)

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

// statusCmd creates the status command.
func statusCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show project status and readiness",
		Long: `Show the status of all specs, evaluations, and approvals for a project.

Displays readiness gates and indicates whether the project is ready
for AI-assisted development.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatus(cmd, args, cfg)
		},
	}

	cmd.Flags().String("format", "text", "Output format: text, json, html, markdown")
	cmd.Flags().Bool("ci", false, "CI mode: exit non-zero if not ready")

	return cmd
}

func runStatus(cmd *cobra.Command, _ []string, cfg *Config) error {
	projectPath, err := resolveProjectPath(cmd)
	if err != nil {
		return err
	}

	project, err := config.Load(projectPath)
	if err != nil {
		return fmt.Errorf("failed to load project: %w", err)
	}

	// Get SpecConfig: CLI config takes precedence, then project config
	specConfig := cfg.GetSpecConfig()
	if cfg.SpecConfig == nil {
		specConfig = project.GetSpecConfig()
	}

	report, err := status.GenerateWithConfig(project, specConfig)
	if err != nil {
		return fmt.Errorf("failed to generate status: %w", err)
	}

	format, _ := cmd.Flags().GetString("format")

	switch format {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(report)
	case "html":
		return status.RenderHTML(os.Stdout, report)
	case "markdown":
		return status.RenderMarkdown(os.Stdout, report)
	default:
		return status.RenderText(os.Stdout, report)
	}
}

func resolveProjectPath(cmd *cobra.Command) (string, error) {
	projectFlag, _ := cmd.Flags().GetString("project")

	if projectFlag != "" {
		// Check if it's an absolute path
		if _, err := os.Stat(projectFlag); err == nil {
			return projectFlag, nil
		}

		// Try as project name under specs directory
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}

		specsDir, err := config.FindSpecsDir(cwd)
		if err != nil {
			return "", fmt.Errorf("specs directory not found")
		}

		projectPath := config.ProjectPath(specsDir, projectFlag)
		if _, err := os.Stat(projectPath); err != nil {
			return "", fmt.Errorf("project %q not found", projectFlag)
		}

		return projectPath, nil
	}

	// Try to find project root from current directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	projectPath, err := config.FindProjectRoot(cwd)
	if err != nil {
		return "", fmt.Errorf("no project found (use --project flag or run from project directory)")
	}

	return projectPath, nil
}

// evalCmd creates the eval command.
func evalCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "eval [spec-type]",
		Short: "Evaluate specs using LLM judges",
		Long: `Evaluate specification documents using LLM-as-a-Judge.

Examples:
  multispec eval prd          # Evaluate PRD
  multispec eval --all        # Evaluate all specs
  multispec eval --source     # Evaluate source specs only
  multispec eval --gtm        # Evaluate GTM docs only`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runEval(cmd, args, cfg)
		},
	}

	cmd.Flags().Bool("all", false, "Evaluate all specs")
	cmd.Flags().Bool("source", false, "Evaluate source specs only")
	cmd.Flags().Bool("gtm", false, "Evaluate GTM docs only")
	cmd.Flags().Bool("technical", false, "Evaluate technical docs only")

	return cmd
}

func runEval(cmd *cobra.Command, args []string, cfg *Config) error {
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

	// Get SpecConfig: CLI config takes precedence, then project config, then defaults
	specConfig := cfg.GetSpecConfig()
	if cfg.SpecConfig == nil {
		// Use project's spec config if CLI doesn't override
		specConfig = project.GetSpecConfig()
	}

	// Determine which specs to evaluate
	var specTypes []types.SpecType

	if len(args) > 0 {
		// Evaluate specific spec type - allow custom types
		specType := types.SpecType(args[0])
		// Check if it's in our config (either built-in or custom)
		if specConfig.GetRequirement(args[0]) == nil && !specType.IsValid() {
			return fmt.Errorf("invalid spec type: %s", args[0])
		}
		specTypes = append(specTypes, specType)
	} else if allFlag {
		for _, name := range specConfig.AllSpecs() {
			specTypes = append(specTypes, types.SpecType(name))
		}
	} else if sourceFlag {
		for _, name := range specConfig.SpecsByCategory(types.CategorySource) {
			specTypes = append(specTypes, types.SpecType(name))
		}
	} else if gtmFlag {
		for _, name := range specConfig.SpecsByCategory(types.CategoryGTM) {
			specTypes = append(specTypes, types.SpecType(name))
		}
	} else if technicalFlag {
		for _, name := range specConfig.SpecsByCategory(types.CategoryTechnical) {
			specTypes = append(specTypes, types.SpecType(name))
		}
	} else {
		return fmt.Errorf("specify a spec type or use --all, --source, --gtm, or --technical")
	}

	// Create LLM client
	llmClient, err := eval.NewLLMClientFromProject(project.LLM)
	if err != nil {
		return fmt.Errorf("initializing LLM: %w", err)
	}
	defer func() { _ = llmClient.Close() }()

	// Create evaluator with optional custom rubric loader
	evaluator := eval.NewEvaluator(llmClient)
	if cfg.RubricLoader != nil {
		evaluator.SetRubricLoader(cfg.RubricLoader)
	}

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

// synthesizeCmd creates the synthesize command.
func synthesizeCmd(cfg *Config) *cobra.Command { //nolint:unparam // cfg reserved for future use
	cmd := &cobra.Command{
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

	cmd.Flags().Bool("eval", false, "Run evaluation after synthesis")

	return cmd
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

// reconcileCmd creates the reconcile command.
func reconcileCmd(cfg *Config) *cobra.Command { //nolint:unparam // cfg reserved for future use
	cmd := &cobra.Command{
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

	return cmd
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

// approveCmd creates the approve command.
func approveCmd(cfg *Config) *cobra.Command { //nolint:unparam // cfg reserved for future use
	cmd := &cobra.Command{
		Use:   "approve <spec-type>",
		Short: "Approve a spec for reconciliation",
		Long: `Mark a specification as approved.

Examples:
  multispec approve prd                    # Approve PRD
  multispec approve trd --approver=eng@co  # Approve with approver`,
		Args: cobra.ExactArgs(1),
		RunE: runApprove,
	}

	cmd.Flags().String("approver", "", "Approver email or identifier")
	cmd.Flags().String("comment", "", "Approval comment")

	return cmd
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

// exportCmd creates the export command.
func exportCmd(cfg *Config) *cobra.Command { //nolint:unparam // cfg reserved for future use
	cmd := &cobra.Command{
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

	cmd.Flags().Bool("dry-run", false, "Show what would be exported without writing")
	cmd.Flags().String("output", "", "Output directory (default: target-specific)")

	return cmd
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

// targetsCmd creates the targets command.
func targetsCmd(cfg *Config) *cobra.Command { //nolint:unparam // cfg reserved for future use
	return &cobra.Command{
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
}

// graphCmd creates the graph command with subcommands.
func graphCmd(cfg *Config) *cobra.Command { //nolint:unparam // cfg reserved for future use
	cmd := &cobra.Command{
		Use:   "graph <subcommand>",
		Short: "Manage requirement graphs via graphize",
		Long: `Manage requirement graphs using graphize integration.

Subcommands:
  extract   Build graph from specs
  query     Query graph relationships
  export    Export to HTML/JSON/GraphML

Examples:
  multispec graph extract                    # Extract graph from current project
  multispec graph export --format html       # Export graph as HTML
  multispec graph export --format graphml    # Export graph as GraphML
  multispec graph query --type requirement   # List all requirements`,
	}

	// Add subcommands
	extractCmd := &cobra.Command{
		Use:   "extract",
		Short: "Extract requirement graph from specs",
		RunE:  runGraphExtract,
	}

	exportSubCmd := &cobra.Command{
		Use:   "export",
		Short: "Export graph to HTML/JSON/GraphML",
		RunE:  runGraphExport,
	}
	exportSubCmd.Flags().String("format", "html", "Export format: html, graphml, json")
	exportSubCmd.Flags().String("output", "", "Output directory (default: .graphize)")

	queryCmd := &cobra.Command{
		Use:   "query",
		Short: "Query graph nodes and relationships",
		RunE:  runGraphQuery,
	}
	queryCmd.Flags().String("type", "", "Filter by node type (requirement, user_story, constraint, decision)")
	queryCmd.Flags().String("spec", "", "Filter by spec type (mrd, prd, uxd, trd)")

	cmd.AddCommand(extractCmd)
	cmd.AddCommand(exportSubCmd)
	cmd.AddCommand(queryCmd)

	return cmd
}

func runGraphExtract(cmd *cobra.Command, args []string) error {
	// Find project root
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	projectPath, err := config.FindProjectRoot(cwd)
	if err != nil {
		return fmt.Errorf("not in a multispec project (no multispec.yaml found)")
	}

	// Extract graph
	extractor := specgraph.NewSpecExtractor(projectPath)
	g, err := extractor.Extract()
	if err != nil {
		return fmt.Errorf("extracting graph: %w", err)
	}

	// Save graph to .graphize directory
	graphDir := filepath.Join(projectPath, ".graphize")
	if err := os.MkdirAll(graphDir, 0755); err != nil {
		return fmt.Errorf("creating .graphize directory: %w", err)
	}

	graphPath := filepath.Join(graphDir, "spec-graph.json")
	if err := specgraph.SaveJSON(g, graphPath); err != nil {
		return fmt.Errorf("saving graph: %w", err)
	}

	fmt.Printf("Extracted graph with %d nodes and %d edges\n", len(g.Nodes), len(g.Edges))
	fmt.Printf("Saved to: %s\n", graphPath)

	// Print summary by node type
	typeCounts := make(map[string]int)
	for _, node := range g.Nodes {
		typeCounts[node.Type]++
	}
	fmt.Println("\nNode types:")
	for nodeType, count := range typeCounts {
		fmt.Printf("  %s: %d\n", nodeType, count)
	}

	return nil
}

func runGraphExport(cmd *cobra.Command, args []string) error {
	format, _ := cmd.Flags().GetString("format")
	output, _ := cmd.Flags().GetString("output")

	// Find project root
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	projectPath, err := config.FindProjectRoot(cwd)
	if err != nil {
		return fmt.Errorf("not in a multispec project (no multispec.yaml found)")
	}

	// Load graph
	graphPath := filepath.Join(projectPath, ".graphize", "spec-graph.json")
	g, err := specgraph.LoadJSON(graphPath)
	if err != nil {
		// Try extracting first
		fmt.Println("Graph not found, extracting...")
		extractor := specgraph.NewSpecExtractor(projectPath)
		g, err = extractor.Extract()
		if err != nil {
			return fmt.Errorf("extracting graph: %w", err)
		}
	}

	// Determine output path
	if output == "" {
		output = filepath.Join(projectPath, ".graphize")
	}

	// Export using library
	result, err := specgraph.Export(g, specgraph.ExportOptions{
		Format:    specgraph.ExportFormat(format),
		OutputDir: output,
		Title:     "Spec Graph",
	})
	if err != nil {
		return err
	}

	fmt.Printf("Exported %s to: %s\n", result.Format, result.OutputPath)
	return nil
}

func runGraphQuery(cmd *cobra.Command, args []string) error {
	nodeType, _ := cmd.Flags().GetString("type")
	specType, _ := cmd.Flags().GetString("spec")

	// Find project root
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	projectPath, err := config.FindProjectRoot(cwd)
	if err != nil {
		return fmt.Errorf("not in a multispec project (no multispec.yaml found)")
	}

	// Load graph
	graphPath := filepath.Join(projectPath, ".graphize", "spec-graph.json")
	g, err := specgraph.LoadJSON(graphPath)
	if err != nil {
		return fmt.Errorf("loading graph (run 'multispec graph extract' first): %w", err)
	}

	// Query using library
	result := specgraph.Query(g, specgraph.QueryFilter{
		NodeType: nodeType,
		SpecType: specType,
	})

	// Print results
	fmt.Printf("Found %d nodes\n\n", result.Count)
	for _, node := range result.Nodes {
		fmt.Printf("[%s] %s\n", node.Type, node.Label)
		fmt.Printf("  ID: %s\n", node.ID)
		if node.Attrs["spec_type"] != "" {
			fmt.Printf("  Spec: %s\n", node.Attrs["spec_type"])
		}
		if node.Attrs["full_text"] != "" && len(node.Attrs["full_text"]) > 100 {
			fmt.Printf("  Text: %s...\n", node.Attrs["full_text"][:100])
		} else if node.Attrs["full_text"] != "" {
			fmt.Printf("  Text: %s\n", node.Attrs["full_text"])
		}
		fmt.Println()
	}

	return nil
}

// serveCmd creates the serve command.
func serveCmd(cfg *Config) *cobra.Command { //nolint:unparam // cfg reserved for future use
	cmd := &cobra.Command{
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

	cmd.Flags().Int("port", 0, "HTTP port (0 for stdio transport)")
	cmd.Flags().String("transport", "stdio", "Transport: stdio, http, sse")

	return cmd
}

func runServe(cmd *cobra.Command, args []string) error {
	server := mcp.NewServer()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	return server.Serve(ctx)
}

// profilesCmd creates the profiles command.
func profilesCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profiles",
		Short: "Manage configuration profiles",
		Long: `Configuration profiles bundle spec requirements, templates, and rubrics.

Default profiles:
  0-1         Minimal for idea validation (hypothesis only)
  startup     Lightweight for pre-PMF startups (PRD only)
  growth      Metrics-driven for 1-N scaling (PRD, UXD, FAQ)
  enterprise  Comprehensive for post-PMF (all specs + security)

Usage:
  multispec profiles list              # List available profiles
  multispec profiles show startup      # Show profile details
  multispec init my-project --profile startup`,
	}

	cmd.AddCommand(profilesListCmd(cfg))
	cmd.AddCommand(profilesShowCmd(cfg))

	return cmd
}

func profilesListCmd(cfg *Config) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			loader := cfg.ProfileLoader
			if loader == nil {
				loader = profiles.DefaultLoader()
			}

			available := loader.Available()

			fmt.Println("Available profiles:")
			fmt.Println()

			for _, name := range available {
				profile, err := loader.Load(name)
				if err != nil {
					fmt.Printf("  %-12s (error loading)\n", name)
					continue
				}

				marker := ""
				if profiles.IsDefaultProfile(name) {
					marker = " [default]"
				}

				fmt.Printf("  %-12s %s%s\n", name, profile.Description, marker)
			}

			fmt.Println()
			fmt.Println("Use with: multispec init <project> --profile <name>")

			return nil
		},
	}
}

func profilesShowCmd(cfg *Config) *cobra.Command {
	return &cobra.Command{
		Use:   "show <profile-name>",
		Short: "Show profile details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			profileName := args[0]

			loader := cfg.ProfileLoader
			if loader == nil {
				loader = profiles.DefaultLoader()
			}

			profile, err := loader.Load(profileName)
			if err != nil {
				return fmt.Errorf("profile %q not found: %w", profileName, err)
			}

			fmt.Printf("Profile: %s\n", profile.Name)
			fmt.Printf("Description: %s\n", profile.Description)
			if profile.Extends != "" {
				fmt.Printf("Extends: %s\n", profile.Extends)
			}
			fmt.Println()

			// Show required specs
			fmt.Println("Required specs:")
			if profile.SpecConfig != nil {
				required := profile.RequiredSpecs()
				if len(required) == 0 {
					fmt.Println("  (none)")
				} else {
					for _, name := range required {
						category := profile.SpecConfig.GetCategory(name)
						fmt.Printf("  - %s (%s)\n", name, category)
					}
				}
			} else {
				fmt.Println("  (uses defaults)")
			}
			fmt.Println()

			// Show available templates
			fmt.Println("Custom templates:")
			if profile.TemplateLoader != nil {
				available := profile.TemplateLoader.Available()
				if len(available) == 0 {
					fmt.Println("  (none)")
				} else {
					for _, t := range available {
						fmt.Printf("  - %s\n", t)
					}
				}
			} else {
				fmt.Println("  (uses defaults)")
			}
			fmt.Println()

			// Show available rubrics
			fmt.Println("Custom rubrics:")
			if profile.RubricLoader != nil {
				available := profile.RubricLoader.Available()
				if len(available) == 0 {
					fmt.Println("  (none)")
				} else {
					for _, r := range available {
						fmt.Printf("  - %s\n", r)
					}
				}
			} else {
				fmt.Println("  (uses defaults)")
			}

			return nil
		},
	}
}
