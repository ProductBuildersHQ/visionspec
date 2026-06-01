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

	"github.com/ProductBuildersHQ/visionspec/internal/mcp"
	"github.com/ProductBuildersHQ/visionspec/pkg/config"
	ctxpkg "github.com/ProductBuildersHQ/visionspec/pkg/context"
	"github.com/ProductBuildersHQ/visionspec/pkg/context/sources"
	"github.com/ProductBuildersHQ/visionspec/pkg/drift"
	"github.com/ProductBuildersHQ/visionspec/pkg/eval"
	"github.com/ProductBuildersHQ/visionspec/pkg/lint"
	"github.com/ProductBuildersHQ/visionspec/pkg/mkdocs"
	"github.com/ProductBuildersHQ/visionspec/pkg/profiles"
	"github.com/ProductBuildersHQ/visionspec/pkg/reconcile"
	"github.com/ProductBuildersHQ/visionspec/pkg/rubrics"
	"github.com/ProductBuildersHQ/visionspec/pkg/rules"
	"github.com/ProductBuildersHQ/visionspec/pkg/specgraph"
	"github.com/ProductBuildersHQ/visionspec/pkg/status"
	"github.com/ProductBuildersHQ/visionspec/pkg/synth"
	"github.com/ProductBuildersHQ/visionspec/pkg/target"
	"github.com/ProductBuildersHQ/visionspec/pkg/templates"
	"github.com/ProductBuildersHQ/visionspec/pkg/testgen"
	"github.com/ProductBuildersHQ/visionspec/pkg/types"
	"github.com/plexusone/structured-evaluation/claims"
	"github.com/plexusone/structured-evaluation/rubric"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// initCmd creates the init command.
func initCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init <project-name>",
		Short: "Initialize a new visionspec project",
		Long: `Initialize a new visionspec project with the canonical directory structure.

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
  └── visionspec.yaml   # Project configuration`,
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
	header := "# visionspec project configuration\n# See: https://github.com/ProductBuildersHQ/visionspec\n\n"
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
	fmt.Printf("\n✅ Created visionspec project: %s\n\n", projectName)
	fmt.Println("Directory structure:")
	fmt.Printf("  %s/\n", projectName)
	fmt.Println("  ├── source/        # Human-authored specs (mrd, prd, uxd)")
	fmt.Println("  ├── gtm/           # Synthesized GTM docs (press, faq, narrative)")
	fmt.Println("  ├── technical/     # Synthesized technical docs (trd, ird)")
	fmt.Println("  ├── eval/          # Evaluation results")
	fmt.Println("  └── visionspec.yaml # Project configuration")
	fmt.Println()
	fmt.Println("Working Backwards workflow:")
	fmt.Println("  1. Write MRD:           visionspec create mrd")
	fmt.Println("  2. Synthesize vision:   visionspec synthesize press")
	fmt.Println("  3. Challenge scope:     visionspec synthesize faq")
	fmt.Println("  4. Derive requirements: visionspec synthesize prd")
	fmt.Println("  5. Review narratives:   visionspec synthesize narrative-1p")
	fmt.Println("  6. Write UXD:           visionspec create uxd")
	fmt.Println("  7. Technical specs:     visionspec synthesize trd && visionspec synthesize ird")
	fmt.Println("  8. Reconcile:           visionspec reconcile")
	fmt.Println()
	fmt.Println("All synthesized docs are editable - refine them in git or with AI assistants.")

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

// createCmd creates the create command for scaffolding new specs.
func createCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <spec-type>",
		Short: "Create a new spec from template",
		Long: `Create a new spec file from a template.

Supported spec types:
  Source specs:   mrd, prd, uxd
  GTM specs:      press, faq, narrative-1p, narrative-6p
  Technical:      trd, ird

The command must be run from within a visionspec project directory.

Examples:
  visionspec create mrd          # Create MRD from template
  visionspec create press        # Create Press Release from template`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreate(cmd, args, cfg)
		},
	}

	cmd.Flags().Bool("force", false, "Overwrite existing file")

	return cmd
}

func runCreate(cmd *cobra.Command, args []string, cfg *Config) error {
	specTypeStr := strings.ToLower(args[0])
	force, _ := cmd.Flags().GetBool("force")

	// Find project root
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	projectPath, err := config.FindProjectRoot(cwd)
	if err != nil {
		return fmt.Errorf("not in a visionspec project (no visionspec.yaml found)")
	}

	// Load project config
	project, err := config.Load(projectPath)
	if err != nil {
		return fmt.Errorf("loading project config: %w", err)
	}

	// Parse spec type
	specType := types.SpecType(specTypeStr)
	if !specType.IsValid() {
		// List available spec types
		available := templates.Available()
		names := make([]string, len(available))
		for i, t := range available {
			names[i] = string(t)
		}
		return fmt.Errorf("invalid spec type %q (available: %s)", specTypeStr, strings.Join(names, ", "))
	}

	// Get template
	loader := cfg.TemplateLoader
	if loader == nil {
		loader = templates.DefaultLoader()
	}

	tmpl, err := loader.Load(specType)
	if err != nil {
		return fmt.Errorf("loading template for %s: %w", specType, err)
	}

	// Determine output path
	outputPath := config.SpecPath(projectPath, specType)

	// Check if file exists
	if _, err := os.Stat(outputPath); err == nil && !force {
		return fmt.Errorf("file %s already exists (use --force to overwrite)", outputPath)
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	// Render template
	opts := templates.DefaultRenderOptions()
	opts.ProjectName = project.Name
	content := tmpl.Render(opts)

	// Write file
	if err := os.WriteFile(outputPath, []byte(strings.TrimSpace(content)+"\n"), 0600); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	fmt.Printf("✓ Created %s\n", outputPath)
	fmt.Printf("\nNext step: Edit %s to add your content\n", filepath.Base(outputPath))

	return nil
}

// lintCmd creates the lint command.
func lintCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lint [project]",
		Short: "Validate directory structure and naming conventions",
		Long: `Validate that the project follows visionspec conventions:

  - Directory structure matches canonical layout
  - File naming follows conventions (lowercase specs, kebab-case projects)
  - Required specs are present
  - Config file is valid

Examples:
  visionspec lint                    # Lint all projects
  visionspec lint user-onboarding    # Lint specific project
  visionspec lint --format json      # Output as JSON`,
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

Output formats:
  - JSON eval result (default): prd.eval.json
  - Structured evaluation report: prd.evaluation.json
  - Claims report: prd.claims.json
  - Summary report: eval-summary.json (with --summary)
  - Markdown report: prd.eval.md (with --markdown)

Examples:
  visionspec eval prd              # Evaluate PRD
  visionspec eval --all            # Evaluate all specs
  visionspec eval --source         # Evaluate source specs only
  visionspec eval --gtm            # Evaluate GTM docs only
  visionspec eval prd --verbose    # Show detailed findings
  visionspec eval --all --summary  # Generate summary with embedded reports
  visionspec eval prd --markdown   # Generate markdown report`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runEval(cmd, args, cfg)
		},
	}

	cmd.Flags().Bool("all", false, "Evaluate all specs")
	cmd.Flags().Bool("source", false, "Evaluate source specs only")
	cmd.Flags().Bool("gtm", false, "Evaluate GTM docs only")
	cmd.Flags().Bool("technical", false, "Evaluate technical docs only")
	cmd.Flags().Bool("verbose", false, "Show detailed findings and explanations")
	cmd.Flags().Bool("markdown", false, "Generate markdown report")
	cmd.Flags().Bool("summary", false, "Generate summary report with embedded reports")
	cmd.Flags().Bool("claims", false, "Generate claims report from findings")

	return cmd
}

func runEval(cmd *cobra.Command, args []string, cfg *Config) error {
	allFlag, _ := cmd.Flags().GetBool("all")
	sourceFlag, _ := cmd.Flags().GetBool("source")
	gtmFlag, _ := cmd.Flags().GetBool("gtm")
	technicalFlag, _ := cmd.Flags().GetBool("technical")
	verboseFlag, _ := cmd.Flags().GetBool("verbose")
	markdownFlag, _ := cmd.Flags().GetBool("markdown")
	summaryFlag, _ := cmd.Flags().GetBool("summary")
	claimsFlag, _ := cmd.Flags().GetBool("claims")

	// Find project root
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	projectPath, err := config.FindProjectRoot(cwd)
	if err != nil {
		return fmt.Errorf("not in a visionspec project (no visionspec.yaml found)")
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
	rubricLoader := cfg.RubricLoader
	if rubricLoader == nil {
		rubricLoader = rubrics.DefaultLoader()
	}
	evaluator.SetRubricLoader(rubricLoader)

	ctx := context.Background()

	// Collect results for summary report
	evalSummary := eval.NewEvalSummary(project.Name, "")

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

		// Generate structured evaluation report
		rubricSet, _ := rubricLoader.Load(specType)
		var evalReport *rubric.Rubric
		if rubricSet != nil {
			evalReport = result.ToEvaluationReport(rubricSet)
			evalReportPath := filepath.Clean(strings.TrimSuffix(evalPath, ".eval.json") + ".evaluation.json")
			if evalReportData, err := json.MarshalIndent(evalReport, "", "  "); err == nil {
				_ = os.WriteFile(evalReportPath, evalReportData, 0600) //nolint:gosec // G703: path derived from validated config.EvalPath, cleaned with filepath.Clean
			}
		}

		// Generate claims report if requested or for summary
		var claimsReport *claims.ClaimsReport
		if claimsFlag || summaryFlag {
			claimsReport = result.ToClaimsReport(string(specType) + ".md")
			if claimsFlag {
				claimsPath := filepath.Clean(strings.TrimSuffix(evalPath, ".eval.json") + ".claims.json")
				if claimsData, err := json.MarshalIndent(claimsReport, "", "  "); err == nil {
					_ = os.WriteFile(claimsPath, claimsData, 0600) //nolint:gosec // G703: path derived from validated config.EvalPath, cleaned with filepath.Clean
				}
			}
		}

		// Generate markdown report if requested
		if markdownFlag {
			mdPath := strings.TrimSuffix(evalPath, ".json") + ".md"
			mdFile, err := os.Create(mdPath)
			if err == nil {
				renderer := eval.NewMarkdownRenderer()
				_ = renderer.Render(mdFile, result)
				mdFile.Close()
			}
		}

		// Add to summary
		evalSummary.AddResult(string(specType), result, evalReport, claimsReport)

		// Print summary
		if verboseFlag {
			renderer := eval.NewTerminalRenderer(true)
			_ = renderer.Render(os.Stdout, result)
		} else if result.Passed {
			fmt.Printf("✓ %s: %.1f/10 PASS (%d findings)\n", specType, result.Score, len(result.Findings))
		} else {
			fmt.Printf("✗ %s: %.1f/10 FAIL (%d findings)\n", specType, result.Score, len(result.Findings))
		}
	}

	// Generate summary report if requested
	if summaryFlag && len(specTypes) > 0 {
		summaryReport := evalSummary.ToSummaryReport("SPEC EVALUATION")
		summaryPath := filepath.Join(projectPath, "eval", "eval-summary.json")
		if summaryData, err := json.MarshalIndent(summaryReport, "", "  "); err == nil {
			if err := os.WriteFile(summaryPath, summaryData, 0600); err == nil {
				fmt.Printf("✓ Generated summary report: %s\n", summaryPath)
			}
		}
	}

	return nil
}

// synthesizeCmd creates the synthesize command.
func synthesizeCmd(cfg *Config) *cobra.Command { //nolint:unparam // cfg reserved for future use
	cmd := &cobra.Command{
		Use:   "synthesize <type>",
		Short: "Generate specs using Working Backwards methodology",
		Long: `Generate specification documents from source specs using LLM synthesis.

Working Backwards Flow:
  visionspec synthesize press        # MRD → press.md (vision document)
  visionspec synthesize faq          # MRD + Press → faq.md (scope clarification)
  visionspec synthesize prd          # MRD + Press + FAQ → prd.md (detailed requirements)

Technical Synthesis:
  visionspec synthesize trd          # MRD + PRD + UXD + CONSTITUTION + CONTEXT → trd.md
  visionspec synthesize tpd          # PRD + TRD + UXD → tpd.md (test plan)
  visionspec synthesize ird          # TRD + CONSTITUTION + CONTEXT → ird.md

Narrative Documents:
  visionspec synthesize narrative-1p # MRD + PRD → narrative-1p.md
  visionspec synthesize narrative-6p # MRD + PRD + UXD → narrative-6p.md

Context grounding:
  For TRD, TPD, and IRD, if context sources are configured, the synthesizer
  will gather codebase context to ground technical decisions in reality.`,
		Args: cobra.ExactArgs(1),
		RunE: runSynthesize,
	}

	cmd.Flags().Bool("eval", false, "Run evaluation after synthesis")
	cmd.Flags().Bool("no-context", false, "Skip context gathering for technical synthesis")

	return cmd
}

func runSynthesize(cmd *cobra.Command, args []string) error {
	specTypeArg := args[0]
	evalFlag, _ := cmd.Flags().GetBool("eval")
	noContext, _ := cmd.Flags().GetBool("no-context")

	// Parse spec type
	specType := types.SpecType(specTypeArg)
	if !synth.CanSynthesize(specType) {
		return fmt.Errorf("cannot synthesize %s (valid: press, faq, prd, trd, tpd, ird, narrative-1p, narrative-6p)", specTypeArg)
	}

	// Find project root
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	projectPath, err := config.FindProjectRoot(cwd)
	if err != nil {
		return fmt.Errorf("not in a visionspec project (no visionspec.yaml found)")
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
	if content, err := os.ReadFile(config.SpecPath(projectPath, types.SpecTypeFAQ)); err == nil {
		input.FAQ = string(content)
	}

	// Load constitution from repo-level or org-level
	input.Constitution = config.LoadConstitution(projectPath)

	// Gather context for TRD/TPD/IRD synthesis (grounding)
	if !noContext && (specType == types.SpecTypeTRD || specType == types.SpecTypeTPD || specType == types.SpecTypeIRD) {
		ctxCfg := getContextConfig(project, projectPath)
		if ctxCfg.HasSources() {
			fmt.Println("⋯ Gathering codebase context for grounding...")
			agg, err := sources.BuildAggregator(project.Name, ctxCfg)
			if err == nil && agg.SourceCount() > 0 {
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
				ac, err := agg.Gather(ctx)
				cancel()
				if err == nil {
					input.Context = ac.Summary
					fmt.Printf("  Gathered context from %d sources\n", len(ac.Sources))
				} else {
					fmt.Printf("  Warning: context gathering failed: %v\n", err)
				}
			}
		}
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
		return fmt.Errorf("not in a visionspec project (no visionspec.yaml found)")
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
		fmt.Println("\nApprove specs with: visionspec approve <spec-type>")
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

	// Load constitution from repo-level or org-level
	input.Constitution = config.LoadConstitution(projectPath)

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

	// Write spec.md output
	outputPath := config.SpecPath(projectPath, types.SpecTypeSpec)
	if err := os.WriteFile(outputPath, []byte(result.Content), 0600); err != nil {
		return fmt.Errorf("writing spec.md: %w", err)
	}

	fmt.Printf("✓ Generated %s\n", outputPath)
	fmt.Printf("  Sources: %v\n", result.Sources)

	if len(result.Conflicts) > 0 {
		fmt.Printf("  Conflicts detected: %d\n", len(result.Conflicts))
		for _, c := range result.Conflicts {
			status := "⚠"
			if c.Resolution != "" {
				status = "✓"
			}
			fmt.Printf("    %s %s: %s\n", status, c.ID, c.Description)
		}
	}

	// Write spec.eval.json with reconciliation metadata
	evalOutput := map[string]any{
		"spec_type":    "spec",
		"generated_at": result.GeneratedAt.Format(time.RFC3339),
		"sources":      result.Sources,
		"conflicts":    result.Conflicts,
		"decision_log": result.DecisionLog,
		"status":       getReconcileStatus(result.Conflicts),
	}

	evalJSON, err := json.MarshalIndent(evalOutput, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling eval output: %w", err)
	}

	evalPath := filepath.Join(projectPath, "eval", "spec.eval.json")
	if err := os.WriteFile(evalPath, evalJSON, 0600); err != nil {
		return fmt.Errorf("writing spec.eval.json: %w", err)
	}

	fmt.Printf("✓ Generated %s\n", evalPath)

	return nil
}

// getReconcileStatus determines overall reconciliation status.
func getReconcileStatus(conflicts []reconcile.Conflict) string {
	unresolvedHigh := 0
	for _, c := range conflicts {
		if c.Resolution == "" && c.Severity == "high" {
			unresolvedHigh++
		}
	}

	if unresolvedHigh > 0 {
		return "needs_review"
	}
	if len(conflicts) > 0 {
		return "reconciled_with_tradeoffs"
	}
	return "reconciled"
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
  visionspec approve prd                    # Approve PRD
  visionspec approve trd --approver=eng@co  # Approve with approver`,
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
		return fmt.Errorf("not in a visionspec project (no visionspec.yaml found)")
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
  visionspec export speckit`,
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
		return fmt.Errorf("not in a visionspec project (no visionspec.yaml found)")
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
			return fmt.Errorf("spec.md not found - run 'visionspec reconcile' first")
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

	// Pass constitution to SpecKit if found (repo-level or org-level)
	if constitutionPath := config.FindConstitution(projectPath); constitutionPath != "" {
		if exportConfig.Options == nil {
			exportConfig.Options = make(map[string]any)
		}
		exportConfig.Options["constitution_path"] = constitutionPath
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
  visionspec graph extract                    # Extract graph from current project
  visionspec graph export --format html       # Export graph as HTML
  visionspec graph export --format graphml    # Export graph as GraphML
  visionspec graph query --type requirement   # List all requirements`,
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
		return fmt.Errorf("not in a visionspec project (no visionspec.yaml found)")
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
		return fmt.Errorf("not in a visionspec project (no visionspec.yaml found)")
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
		return fmt.Errorf("not in a visionspec project (no visionspec.yaml found)")
	}

	// Load graph
	graphPath := filepath.Join(projectPath, ".graphize", "spec-graph.json")
	g, err := specgraph.LoadJSON(graphPath)
	if err != nil {
		return fmt.Errorf("loading graph (run 'visionspec graph extract' first): %w", err)
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
      "visionspec": {
        "command": "visionspec",
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
  visionspec profiles list              # List available profiles
  visionspec profiles show startup      # Show profile details
  visionspec profiles export startup ./my-profile  # Export for customization
  visionspec init my-project --profile startup`,
	}

	cmd.AddCommand(profilesListCmd(cfg))
	cmd.AddCommand(profilesShowCmd(cfg))
	cmd.AddCommand(profilesExportCmd(cfg))

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
			fmt.Println("Use with: visionspec init <project> --profile <name>")

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

func profilesExportCmd(cfg *Config) *cobra.Command {
	return &cobra.Command{
		Use:   "export <profile-name> <output-dir>",
		Short: "Export a profile to a directory for customization",
		Long: `Export a built-in profile to a directory so you can customize it.

This creates a complete profile directory with:
  - profile.yaml     Configuration file
  - templates/       Template files (.md)
  - rubrics/         Rubric files (.rubric.yaml)

You can then modify these files and use them as a custom profile.

Examples:
  # Export enterprise profile to customize
  visionspec profiles export enterprise ./my-profile

  # Use the exported profile
  visionspec init my-project --profile-dir ./my-profile`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			profileName := args[0]
			outputDir := args[1]

			loader := cfg.ProfileLoader
			if loader == nil {
				loader = profiles.DefaultLoader()
			}

			// Verify profile exists
			profile, err := loader.Load(profileName)
			if err != nil {
				return fmt.Errorf("profile %q not found: %w", profileName, err)
			}

			// Create output directory
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				return fmt.Errorf("creating output directory: %w", err)
			}

			// Export profile.yaml
			profileYAML := profiles.ProfileToYAML(profile)
			profilePath := filepath.Join(outputDir, "profile.yaml")
			if err := profiles.WriteProfileYAML(profilePath, profileYAML); err != nil {
				return fmt.Errorf("writing profile.yaml: %w", err)
			}
			fmt.Printf("Created %s\n", profilePath)

			// Export templates
			if profile.TemplateLoader != nil {
				templatesDir := filepath.Join(outputDir, "templates")
				if err := os.MkdirAll(templatesDir, 0755); err != nil {
					return fmt.Errorf("creating templates directory: %w", err)
				}

				for _, specType := range profile.TemplateLoader.Available() {
					tmpl, err := profile.TemplateLoader.Load(specType)
					if err != nil {
						continue
					}
					filename := string(specType) + ".md"
					path := filepath.Join(templatesDir, filename)
					if err := os.WriteFile(path, []byte(tmpl.Content), 0600); err != nil {
						return fmt.Errorf("writing template %s: %w", filename, err)
					}
					fmt.Printf("Created %s\n", path)
				}
			}

			// Export rubrics
			if profile.RubricLoader != nil {
				rubricsDir := filepath.Join(outputDir, "rubrics")
				if err := os.MkdirAll(rubricsDir, 0755); err != nil {
					return fmt.Errorf("creating rubrics directory: %w", err)
				}

				for _, specType := range profile.RubricLoader.Available() {
					rubric, err := profile.RubricLoader.Load(specType)
					if err != nil {
						continue
					}
					filename := string(specType) + ".rubric.yaml"
					path := filepath.Join(rubricsDir, filename)
					if err := rubrics.WriteRubricYAML(path, rubric); err != nil {
						return fmt.Errorf("writing rubric %s: %w", filename, err)
					}
					fmt.Printf("Created %s\n", path)
				}
			}

			fmt.Println()
			fmt.Printf("Profile exported to %s\n", outputDir)
			fmt.Println()
			fmt.Println("To use this profile:")
			fmt.Printf("  visionspec init my-project --profile-dir %s\n", outputDir)

			return nil
		},
	}
}

// contextCmd creates the context command with subcommands.
func contextCmd(cfg *Config) *cobra.Command { //nolint:unparam // cfg reserved for future use
	cmd := &cobra.Command{
		Use:   "context <subcommand>",
		Short: "Gather and manage codebase context for grounding",
		Long: `Gather context from git repositories, graphize graphs, and external sources.

Context is used to ground spec synthesis in the reality of existing codebases,
requirement traceability, and external tool state.

Subcommands:
  gather    Collect context from all configured sources
  show      Display current context summary
  save      Save context snapshot to file
  load      Load context snapshot from file
  sources   List configured context sources

Examples:
  visionspec context gather                  # Gather context from all sources
  visionspec context show                    # Show context summary
  visionspec context save --output ctx.json  # Save snapshot
  visionspec context sources                 # List configured sources`,
	}

	cmd.AddCommand(contextGatherCmd(cfg))
	cmd.AddCommand(contextShowCmd(cfg))
	cmd.AddCommand(contextSaveCmd(cfg))
	cmd.AddCommand(contextSourcesCmd(cfg))

	return cmd
}

func contextGatherCmd(cfg *Config) *cobra.Command { //nolint:unparam // cfg reserved for future use
	cmd := &cobra.Command{
		Use:   "gather",
		Short: "Gather context from all configured sources",
		RunE:  runContextGather,
	}

	cmd.Flags().Duration("timeout", 2*time.Minute, "Timeout for gathering context")
	cmd.Flags().String("format", "text", "Output format: text, json")
	cmd.Flags().Bool("refresh", false, "Refresh cache before gathering")

	return cmd
}

func runContextGather(cmd *cobra.Command, args []string) error {
	timeout, _ := cmd.Flags().GetDuration("timeout")
	format, _ := cmd.Flags().GetString("format")
	refresh, _ := cmd.Flags().GetBool("refresh")

	// Find project root
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	projectPath, err := config.FindProjectRoot(cwd)
	if err != nil {
		return fmt.Errorf("not in a visionspec project (no visionspec.yaml found)")
	}

	// Load project config
	project, err := config.Load(projectPath)
	if err != nil {
		return fmt.Errorf("loading project config: %w", err)
	}

	// Get context configuration from project
	ctxCfg := getContextConfig(project, projectPath)

	// Create aggregator
	agg, err := sources.BuildAggregator(project.Name, ctxCfg)
	if err != nil {
		return fmt.Errorf("building aggregator: %w", err)
	}

	if agg.SourceCount() == 0 {
		fmt.Println("No context sources configured.")
		fmt.Println("\nAdd sources to visionspec.yaml:")
		fmt.Println("  context:")
		fmt.Println("    repositories:")
		fmt.Println("      - path: /path/to/repo")
		fmt.Println("    files:")
		fmt.Println("      - path: architecture.md")
		return nil
	}

	fmt.Printf("Gathering context from %d sources...\n", agg.SourceCount())

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var ac *ctxpkg.AggregatedContext
	if refresh {
		ac, err = agg.Refresh(ctx)
	} else {
		ac, err = agg.Gather(ctx)
	}
	if err != nil {
		return fmt.Errorf("gathering context: %w", err)
	}

	switch format {
	case "json":
		data, err := ac.ToJSON()
		if err != nil {
			return fmt.Errorf("marshaling context: %w", err)
		}
		fmt.Println(string(data))
	default:
		fmt.Println(ac.Summary)
	}

	return nil
}

func contextShowCmd(cfg *Config) *cobra.Command { //nolint:unparam // cfg reserved for future use
	return &cobra.Command{
		Use:   "show",
		Short: "Show current context summary",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Find project root
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("getting working directory: %w", err)
			}

			projectPath, err := config.FindProjectRoot(cwd)
			if err != nil {
				return fmt.Errorf("not in a visionspec project (no visionspec.yaml found)")
			}

			// Try to load existing snapshot
			snapshotPath := filepath.Join(projectPath, ".context-snapshot.json")
			ac, err := ctxpkg.LoadSnapshot(snapshotPath)
			if err != nil {
				fmt.Println("No context snapshot found. Run 'visionspec context gather' first.")
				return nil
			}

			fmt.Println(ac.Summary)
			return nil
		},
	}
}

func contextSaveCmd(cfg *Config) *cobra.Command { //nolint:unparam // cfg reserved for future use
	cmd := &cobra.Command{
		Use:   "save",
		Short: "Save context snapshot to file",
		RunE:  runContextSave,
	}

	cmd.Flags().StringP("output", "o", "", "Output file path (default: .context-snapshot.json)")

	return cmd
}

func runContextSave(cmd *cobra.Command, args []string) error {
	output, _ := cmd.Flags().GetString("output")

	// Find project root
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	projectPath, err := config.FindProjectRoot(cwd)
	if err != nil {
		return fmt.Errorf("not in a visionspec project (no visionspec.yaml found)")
	}

	// Load project config
	project, err := config.Load(projectPath)
	if err != nil {
		return fmt.Errorf("loading project config: %w", err)
	}

	// Get context configuration
	ctxCfg := getContextConfig(project, projectPath)

	// Create aggregator and gather
	agg, err := sources.BuildAggregator(project.Name, ctxCfg)
	if err != nil {
		return fmt.Errorf("building aggregator: %w", err)
	}

	if agg.SourceCount() == 0 {
		return fmt.Errorf("no context sources configured")
	}

	fmt.Printf("Gathering context from %d sources...\n", agg.SourceCount())

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	ac, err := agg.Gather(ctx)
	if err != nil {
		return fmt.Errorf("gathering context: %w", err)
	}

	// Determine output path
	if output == "" {
		output = filepath.Join(projectPath, ".context-snapshot.json")
	}

	// Save snapshot
	if err := ctxpkg.SaveSnapshot(ac, output); err != nil {
		return fmt.Errorf("saving snapshot: %w", err)
	}

	fmt.Printf("Saved context snapshot to: %s\n", output)
	return nil
}

func contextSourcesCmd(cfg *Config) *cobra.Command { //nolint:unparam // cfg reserved for future use
	return &cobra.Command{
		Use:   "sources",
		Short: "List configured context sources",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Find project root
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("getting working directory: %w", err)
			}

			projectPath, err := config.FindProjectRoot(cwd)
			if err != nil {
				return fmt.Errorf("not in a visionspec project (no visionspec.yaml found)")
			}

			// Load project config
			project, err := config.Load(projectPath)
			if err != nil {
				return fmt.Errorf("loading project config: %w", err)
			}

			// Get context configuration
			ctxCfg := getContextConfig(project, projectPath)

			fmt.Println("Configured context sources:")
			fmt.Println()

			if len(ctxCfg.Repositories) > 0 {
				fmt.Println("Git Repositories:")
				for _, repo := range ctxCfg.Repositories {
					fmt.Printf("  - %s\n", repo.Path)
					if repo.Graphize == "auto" || repo.Graphize == "true" {
						fmt.Printf("    (graphize: %s)\n", repo.Graphize)
					}
				}
				fmt.Println()
			}

			if len(ctxCfg.Graphize) > 0 {
				fmt.Println("Graphize Graphs:")
				for _, g := range ctxCfg.Graphize {
					name := g.Name
					if name == "" {
						name = filepath.Base(g.Path)
					}
					fmt.Printf("  - %s (%s)\n", name, g.Path)
				}
				fmt.Println()
			}

			if len(ctxCfg.Files) > 0 {
				fmt.Println("Local Files:")
				for _, f := range ctxCfg.Files {
					fileType := f.Type
					if fileType == "" {
						fileType = "document"
					}
					fmt.Printf("  - %s (%s)\n", f.Path, fileType)
				}
				fmt.Println()
			}

			if len(ctxCfg.MCPServers) > 0 {
				fmt.Println("MCP Servers:")
				for name := range ctxCfg.MCPServers {
					fmt.Printf("  - %s (not yet implemented)\n", name)
				}
				fmt.Println()
			}

			if !ctxCfg.HasSources() {
				fmt.Println("  (no sources configured)")
				fmt.Println()
				fmt.Println("Add sources to visionspec.yaml:")
				fmt.Println("  context:")
				fmt.Println("    repositories:")
				fmt.Println("      - path: /path/to/repo")
			}

			return nil
		},
	}
}

// getContextConfig extracts context configuration from project.
func getContextConfig(project *types.Project, projectPath string) *ctxpkg.Config {
	cfg := ctxpkg.DefaultConfig()
	cfg.ProjectName = project.Name

	// Check if project has context configuration
	if project.Context != nil {
		// Map project context config to ctxpkg.Config
		if project.Context.Repositories != nil {
			for _, repo := range project.Context.Repositories {
				cfg.Repositories = append(cfg.Repositories, ctxpkg.RepositoryConfig{
					Path:     repo.Path,
					URL:      repo.URL,
					Branch:   repo.Branch,
					Include:  repo.Include,
					Exclude:  repo.Exclude,
					Analyze:  repo.Analyze,
					Graphize: repo.Graphize,
					MaxDepth: repo.MaxDepth,
				})
			}
		}
		if project.Context.Graphize != nil {
			for _, g := range project.Context.Graphize {
				cfg.Graphize = append(cfg.Graphize, ctxpkg.GraphizeConfig{
					Path:         g.Path,
					Name:         g.Name,
					IncludeNodes: g.IncludeNodes,
					IncludeEdges: g.IncludeEdges,
				})
			}
		}
		if project.Context.Files != nil {
			for _, f := range project.Context.Files {
				cfg.Files = append(cfg.Files, ctxpkg.FileConfig{
					Path:    f.Path,
					Type:    f.Type,
					MaxSize: f.MaxSize,
				})
			}
		}
		if project.Context.MCPServers != nil {
			cfg.MCPServers = make(map[string]ctxpkg.MCPServerConfig)
			for name, srv := range project.Context.MCPServers {
				cfg.MCPServers[name] = ctxpkg.MCPServerConfig{
					Command: srv.Command,
					Args:    srv.Args,
					Env:     srv.Env,
					Config:  srv.Config,
					Timeout: srv.Timeout,
				}
			}
		}
		if project.Context.CacheTTL > 0 {
			cfg.CacheTTL = project.Context.CacheTTL
		}
	}

	// Auto-detect: if no repos configured, use current project path
	if len(cfg.Repositories) == 0 {
		// Check if project path has a .git directory
		gitPath := filepath.Join(projectPath, "..", "..", "..", ".git")
		if info, err := os.Stat(gitPath); err == nil && info.IsDir() {
			repoPath := filepath.Join(projectPath, "..", "..", "..")
			cfg.Repositories = append(cfg.Repositories, ctxpkg.RepositoryConfig{
				Path:     repoPath,
				Graphize: "auto",
			})
		}
	}

	return cfg
}

// docsCmd creates the docs command for MkDocs generation.
func docsCmd(cfg *Config) *cobra.Command { //nolint:unparam // cfg reserved for future use
	cmd := &cobra.Command{
		Use:   "docs <subcommand>",
		Short: "Generate MkDocs-compatible documentation",
		Long: `Generate markdown files for MkDocs integration.

Subcommands:
  generate    Generate all index.md files for projects and specs landing page
  project     Generate index.md for a specific project

Examples:
  visionspec docs generate             # Generate all docs
  visionspec docs project my-project   # Generate docs for specific project`,
	}

	cmd.AddCommand(docsGenerateCmd(cfg))
	cmd.AddCommand(docsProjectCmd(cfg))

	return cmd
}

func docsGenerateCmd(cfg *Config) *cobra.Command { //nolint:unparam // cfg reserved for future use
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate all MkDocs index files",
		Long: `Generate index.md files for all projects and the specs landing page.

Creates:
  - docs/specs/index.md (specs landing page)
  - docs/specs/{project}/index.md (for each project)`,
		RunE: runDocsGenerate,
	}

	cmd.Flags().Bool("with-graph", false, "Include graph metrics in reports")

	return cmd
}

func runDocsGenerate(cmd *cobra.Command, args []string) error {
	withGraph, _ := cmd.Flags().GetBool("with-graph")

	// Find specs directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	specsDir, err := config.FindSpecsDir(cwd)
	if err != nil {
		return fmt.Errorf("not in a visionspec workspace (no docs/specs found)")
	}

	fmt.Println("⋯ Generating MkDocs files...")

	// Generate landing page
	if err := mkdocs.WriteSpecsLanding(specsDir, mkdocs.SpecsLandingOptions{}); err != nil {
		return fmt.Errorf("generating specs landing: %w", err)
	}
	fmt.Printf("  ✓ Generated %s/index.md\n", filepath.Base(specsDir))

	// Generate project indexes
	entries, err := os.ReadDir(specsDir)
	if err != nil {
		return fmt.Errorf("reading specs directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		projectPath := filepath.Join(specsDir, entry.Name())
		configPath := filepath.Join(projectPath, config.ConfigFileName)
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			continue
		}

		// Load project and generate report
		project, err := config.Load(projectPath)
		if err != nil {
			fmt.Printf("  ⚠ Skipping %s: %v\n", entry.Name(), err)
			continue
		}

		report, err := status.Generate(project)
		if err != nil {
			fmt.Printf("  ⚠ Skipping %s: %v\n", entry.Name(), err)
			continue
		}

		// Add graph metrics if requested
		if withGraph {
			extractor := specgraph.NewSpecExtractor(projectPath)
			graph, err := extractor.Extract()
			if err == nil {
				metrics := specgraph.ComputeMetrics(graph)
				report.GraphMetrics = &status.GraphMetrics{
					TotalRequirements: metrics.TotalRequirements,
					TotalUserStories:  metrics.TotalUserStories,
					TotalConstraints:  metrics.TotalConstraints,
					TotalDecisions:    metrics.TotalDecisions,
					TraceCoverage:     metrics.TraceCoverage,
					ConflictCount:     metrics.ConflictCount,
				}
			}
		}

		opts := mkdocs.ProjectIndexOptions{
			IncludeGraphLink: withGraph,
			GraphPath:        "graph/graph.html",
		}

		if err := mkdocs.WriteProjectIndex(projectPath, report, opts); err != nil {
			fmt.Printf("  ⚠ Error writing %s/index.md: %v\n", entry.Name(), err)
			continue
		}
		fmt.Printf("  ✓ Generated %s/index.md\n", entry.Name())
	}

	fmt.Println("\n✓ MkDocs documentation generated")
	return nil
}

func docsProjectCmd(cfg *Config) *cobra.Command { //nolint:unparam // cfg reserved for future use
	cmd := &cobra.Command{
		Use:   "project [project-name]",
		Short: "Generate index.md for a specific project",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runDocsProject,
	}

	cmd.Flags().Bool("with-graph", false, "Include graph metrics in report")

	return cmd
}

func runDocsProject(cmd *cobra.Command, args []string) error {
	withGraph, _ := cmd.Flags().GetBool("with-graph")

	// Find project path
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	var projectPath string
	if len(args) > 0 {
		specsDir, err := config.FindSpecsDir(cwd)
		if err != nil {
			return fmt.Errorf("not in a visionspec workspace (no docs/specs found)")
		}
		projectPath = filepath.Join(specsDir, args[0])
	} else {
		projectPath, err = config.FindProjectRoot(cwd)
		if err != nil {
			return fmt.Errorf("not in a visionspec project (no visionspec.yaml found)")
		}
	}

	// Load project and generate report
	project, err := config.Load(projectPath)
	if err != nil {
		return fmt.Errorf("loading project config: %w", err)
	}

	report, err := status.Generate(project)
	if err != nil {
		return fmt.Errorf("generating status report: %w", err)
	}

	// Add graph metrics if requested
	if withGraph {
		extractor := specgraph.NewSpecExtractor(projectPath)
		graph, err := extractor.Extract()
		if err == nil {
			metrics := specgraph.ComputeMetrics(graph)
			report.GraphMetrics = &status.GraphMetrics{
				TotalRequirements: metrics.TotalRequirements,
				TotalUserStories:  metrics.TotalUserStories,
				TotalConstraints:  metrics.TotalConstraints,
				TotalDecisions:    metrics.TotalDecisions,
				TraceCoverage:     metrics.TraceCoverage,
				ConflictCount:     metrics.ConflictCount,
			}
		}
	}

	opts := mkdocs.ProjectIndexOptions{
		IncludeGraphLink: withGraph,
		GraphPath:        "graph/graph.html",
	}

	if err := mkdocs.WriteProjectIndex(projectPath, report, opts); err != nil {
		return fmt.Errorf("writing index.md: %w", err)
	}

	fmt.Printf("✓ Generated %s/index.md\n", project.Name)
	return nil
}

// rulesCmd creates the rules command.
func rulesCmd(_ *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rules",
		Short: "Manage workflow rules for AI assistant orchestration",
		Long: `Workflow rules guide AI assistants (Claude Code, Kiro, Cursor) through
the VisionSpec specification workflow.

Rules provide:
  - Trigger patterns for activating VisionSpec workflows
  - Phase-by-phase guidance for spec creation
  - Framework-specific flows (AWS, Lean Startup, Design Thinking, etc.)
  - Evaluation and approval gates

Usage:
  visionspec rules list              # List available rules
  visionspec rules export [dir]      # Export rules to project directory`,
	}

	cmd.AddCommand(rulesListCmd())
	cmd.AddCommand(rulesExportCmd())

	return cmd
}

func rulesListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available workflow rules",
		RunE: func(cmd *cobra.Command, args []string) error {
			ruleFiles, err := rules.List()
			if err != nil {
				return fmt.Errorf("listing rules: %w", err)
			}

			fmt.Println("Available workflow rules:")
			fmt.Println()

			for _, f := range ruleFiles {
				fmt.Printf("  %s\n", f)
			}

			fmt.Println()
			fmt.Println("Export rules to your project with: visionspec rules export")
			return nil
		},
	}
}

func rulesExportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "export [output-dir]",
		Short: "Export workflow rules to a directory",
		Long: `Export VisionSpec workflow rules to your project.

This copies the workflow rules that guide AI assistants through
spec creation and review. The rules work with:
  - Claude Code (via CLAUDE.md reference)
  - AWS Kiro (via .kiro/steering/)
  - Cursor (via .cursor/rules/)

Examples:
  # Export to default location (.visionspec-rules)
  visionspec rules export

  # Export to custom directory
  visionspec rules export ./my-rules

After export, reference in your CLAUDE.md:
  See .visionspec-rules/ for VisionSpec workflow guidance.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			outputDir := ".visionspec-rules"
			if len(args) > 0 {
				outputDir = args[0]
			}

			files, err := rules.Export(outputDir)
			if err != nil {
				return fmt.Errorf("exporting rules: %w", err)
			}

			fmt.Printf("✓ Exported %d rule files to %s\n", len(files), outputDir)
			fmt.Println()
			fmt.Println("Contents:")
			fmt.Println("  core-workflow.md         - Main orchestration rules")
			fmt.Println("  phases/                  - Phase-by-phase guidance")
			fmt.Println("  gates/                   - Evaluation and approval gates")
			fmt.Println("  frameworks/              - Framework-specific flows")
			fmt.Println()
			fmt.Println("To use with Claude Code, add to your CLAUDE.md:")
			fmt.Println("  See .visionspec-rules/ for VisionSpec workflow guidance.")
			return nil
		},
	}
}

// generateCmd creates the generate command with subcommands.
func generateCmd(cfg *Config) *cobra.Command { //nolint:unparam // cfg reserved for future use
	cmd := &cobra.Command{
		Use:   "generate <subcommand>",
		Short: "Generate artifacts from specs",
		Long: `Generate various artifacts from specification documents.

Subcommands:
  tests     Generate test stubs from TPD (Test Plan Document)

Examples:
  visionspec generate tests --lang go --output ./tests
  visionspec generate tests --lang ts --framework jest
  visionspec generate tests --lang py`,
	}

	cmd.AddCommand(generateTestsCmd(cfg))

	return cmd
}

func generateTestsCmd(cfg *Config) *cobra.Command { //nolint:unparam // cfg reserved for future use
	cmd := &cobra.Command{
		Use:   "tests",
		Short: "Generate test stubs from TPD",
		Long: `Generate executable test stubs from the Test Plan Document (TPD).

Supported languages:
  go    Go (testing or testify framework)
  ts    TypeScript (Jest or Vitest)
  py    Python (pytest)

The command parses test tables from the TPD and generates corresponding
test stub files with TODOs for implementation.

Examples:
  visionspec generate tests --lang go                    # Go with standard testing
  visionspec generate tests --lang go --framework testify # Go with testify
  visionspec generate tests --lang ts --output ./tests   # TypeScript/Jest
  visionspec generate tests --lang py --group-by priority # Python grouped by priority`,
		RunE: runGenerateTests,
	}

	cmd.Flags().String("lang", "go", "Target language: go, ts, py")
	cmd.Flags().String("framework", "", "Test framework (go: testing|testify, ts: jest|vitest, py: pytest)")
	cmd.Flags().String("output", "", "Output directory (default: ./generated-tests)")
	cmd.Flags().String("package", "", "Package/module name")
	cmd.Flags().String("group-by", "type", "Group tests by: type, priority, file")

	return cmd
}

func runGenerateTests(cmd *cobra.Command, args []string) error {
	lang, _ := cmd.Flags().GetString("lang")
	framework, _ := cmd.Flags().GetString("framework")
	outputDir, _ := cmd.Flags().GetString("output")
	packageName, _ := cmd.Flags().GetString("package")
	groupBy, _ := cmd.Flags().GetString("group-by")

	// Find project root
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	projectPath, err := config.FindProjectRoot(cwd)
	if err != nil {
		return fmt.Errorf("not in a visionspec project (no visionspec.yaml found)")
	}

	// Load TPD
	tpdPath := config.SpecPath(projectPath, types.SpecTypeTPD)
	tpdContent, err := os.ReadFile(tpdPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("TPD not found at %s - run 'visionspec synthesize tpd' first", tpdPath)
		}
		return fmt.Errorf("reading TPD: %w", err)
	}

	// Parse TPD
	parser := testgen.NewParser()
	parsed, err := parser.Parse(string(tpdContent))
	if err != nil {
		return fmt.Errorf("parsing TPD: %w", err)
	}

	allTests := parsed.AllTestCases()
	if len(allTests) == 0 {
		fmt.Println("No test cases found in TPD")
		return nil
	}

	fmt.Printf("Found %d test cases in TPD\n", len(allTests))

	// Get generator
	gen, err := testgen.Get(lang)
	if err != nil {
		available := testgen.Available()
		return fmt.Errorf("unknown language %q (available: %v)", lang, available)
	}

	// Set defaults
	if outputDir == "" {
		outputDir = filepath.Join(projectPath, "generated-tests", lang)
	}
	if packageName == "" {
		packageName = "tests"
	}

	// Generate
	opts := testgen.GenerateOptions{
		OutputDir:     outputDir,
		PackageName:   packageName,
		TestFramework: framework,
		GroupBy:       groupBy,
	}

	fmt.Printf("Generating %s tests to %s...\n", lang, outputDir)

	result, err := gen.Generate(allTests, opts)
	if err != nil {
		return fmt.Errorf("generating tests: %w", err)
	}

	// Print summary
	fmt.Printf("\n✓ Generated %d test stubs\n", result.TotalTests)
	fmt.Printf("  Language:  %s\n", result.Language)
	fmt.Printf("  Framework: %s\n", result.Framework)
	fmt.Printf("  Output:    %s\n", result.OutputDir)
	fmt.Println("  Files:")
	for _, f := range result.Files {
		fmt.Printf("    - %s (%d tests)\n", filepath.Base(f.Path), f.TestCount)
	}

	fmt.Println("\nNext steps:")
	fmt.Println("  1. Review generated test stubs")
	fmt.Println("  2. Implement TODO sections")
	fmt.Println("  3. Remove @pytest.mark.skip / t.Skip() markers")

	return nil
}

// driftCmd creates the drift detection command.
func driftCmd(cfg *Config) *cobra.Command { //nolint:unparam // cfg reserved for future use
	cmd := &cobra.Command{
		Use:   "drift",
		Short: "Detect drift between spec and code",
		Long: `Detect drift between the reconciled spec.md and the actual codebase.

Drift detection identifies:
  - Unimplemented: Requirements in spec not found in code
  - Undocumented: Code that exists without spec coverage
  - Mismatch: Spec and code exist but differ

The command uses context sources configured in visionspec.yaml to analyze
the codebase. Run 'visionspec context gather' first to populate context.

Examples:
  visionspec drift                    # Full drift report
  visionspec drift --severity high    # Only high/critical items
  visionspec drift --format json      # JSON output
  visionspec drift --ci               # Exit non-zero if drift found`,
		RunE: runDrift,
	}

	cmd.Flags().String("format", "text", "Output format: text, json, markdown")
	cmd.Flags().String("severity", "low", "Minimum severity: low, medium, high, critical")
	cmd.Flags().Bool("ci", false, "CI mode: exit non-zero if drift detected")

	return cmd
}

func runDrift(cmd *cobra.Command, args []string) error {
	formatStr, _ := cmd.Flags().GetString("format")
	severityStr, _ := cmd.Flags().GetString("severity")
	ciMode, _ := cmd.Flags().GetBool("ci")

	// Parse format
	var format drift.RenderFormat
	switch formatStr {
	case "json":
		format = drift.FormatJSON
	case "markdown", "md":
		format = drift.FormatMarkdown
	default:
		format = drift.FormatText
	}

	// Parse severity
	var minSeverity drift.Severity
	switch severityStr {
	case "critical":
		minSeverity = drift.SeverityCritical
	case "high":
		minSeverity = drift.SeverityHigh
	case "medium":
		minSeverity = drift.SeverityMedium
	default:
		minSeverity = drift.SeverityLow
	}

	// Find project root
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	projectPath, err := config.FindProjectRoot(cwd)
	if err != nil {
		return fmt.Errorf("not in a visionspec project (no visionspec.yaml found)")
	}

	// Load project
	project, err := config.Load(projectPath)
	if err != nil {
		return fmt.Errorf("loading project: %w", err)
	}

	// Load spec.md
	specPath := filepath.Join(projectPath, "spec.md")
	specContent, err := os.ReadFile(specPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("spec.md not found - run 'visionspec reconcile' first")
		}
		return fmt.Errorf("reading spec.md: %w", err)
	}

	// Load context (if available)
	var aggregatedCtx *ctxpkg.AggregatedContext
	contextCachePath := filepath.Join(projectPath, ".visionspec", "context-cache.json")
	if contextData, err := os.ReadFile(contextCachePath); err == nil {
		var ctx ctxpkg.AggregatedContext
		if json.Unmarshal(contextData, &ctx) == nil {
			aggregatedCtx = &ctx
		}
	}

	if aggregatedCtx == nil {
		// Create minimal context
		aggregatedCtx = &ctxpkg.AggregatedContext{
			Project: project.Name,
		}
	}

	// Run drift detection
	detector := drift.NewDetector()
	opts := drift.DetectOptions{
		MinSeverity: minSeverity,
	}

	report, err := detector.Detect(string(specContent), aggregatedCtx, opts)
	if err != nil {
		return fmt.Errorf("drift detection failed: %w", err)
	}

	// Render report
	renderer := drift.NewRenderer(format)
	if err := renderer.Render(os.Stdout, report); err != nil {
		return fmt.Errorf("rendering report: %w", err)
	}

	// CI mode: exit with error if drift found
	if ciMode && report.HasDrift() {
		if report.HasBlockers() {
			return fmt.Errorf("drift detected with blockers (critical/high severity)")
		}
		return fmt.Errorf("drift detected")
	}

	return nil
}

// syncCmd creates the sync command.
func syncCmd(cfg *Config) *cobra.Command { //nolint:unparam // cfg reserved for future use
	cmd := &cobra.Command{
		Use:   "sync [target]",
		Short: "Sync task state from exported targets",
		Long: `Sync task state from exported execution targets.

This command retrieves the current state of tasks from an exported target
(SpecKit, GSD, GasTown) and updates the project's execution state.

Supported targets:
  speckit   GitHub Spec-Kit (tasks.md)
  gsd       GSD (STATE.md)
  gastown   GasTown (beads/*.toml)

Examples:
  visionspec sync speckit
  visionspec sync gsd
  visionspec sync gastown`,
		Args: cobra.MaximumNArgs(1),
		RunE: runSync,
	}

	cmd.Flags().BoolP("verbose", "v", false, "Show detailed task information")

	return cmd
}

func runSync(cmd *cobra.Command, args []string) error {
	verbose, _ := cmd.Flags().GetBool("verbose")

	// Find project root
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	projectPath, err := config.FindProjectRoot(cwd)
	if err != nil {
		return fmt.Errorf("not in a visionspec project (no visionspec.yaml found)")
	}

	// Load project
	project, err := config.Load(projectPath)
	if err != nil {
		return fmt.Errorf("loading project: %w", err)
	}

	// Determine target
	var targetName string
	if len(args) > 0 {
		targetName = args[0]
	} else if project.Targets.Default != "" {
		targetName = project.Targets.Default
	} else {
		// List available syncable targets
		syncable := target.SyncableTargets()
		if len(syncable) == 0 {
			return fmt.Errorf("no syncable targets available")
		}
		return fmt.Errorf("specify a target: visionspec sync <%s>", strings.Join(syncable, "|"))
	}

	// Get syncer
	syncer, err := target.GetSyncer(targetName)
	if err != nil {
		return err
	}

	// Get target config
	exportCfg := target.ProjectTargetConfig(project, targetName)

	// Check if sync is possible
	if !syncer.CanSync(*exportCfg) {
		return fmt.Errorf("cannot sync %s - export first with 'visionspec export %s'", targetName, targetName)
	}

	fmt.Printf("Syncing task state from %s...\n", targetName)

	// Perform sync
	result, err := syncer.Sync(*exportCfg)
	if err != nil {
		return fmt.Errorf("sync failed: %w", err)
	}

	// Update project execution state
	project.Execution = &types.ExecutionState{
		Target:   result.Target,
		SyncedAt: result.SyncedAt,
		Tasks:    make([]types.ExecutionTask, len(result.Tasks)),
		Summary: types.ExecutionSummary{
			TotalTasks: result.Summary.TotalTasks,
			TodoCount:  result.Summary.TodoCount,
			InProgress: result.Summary.InProgress,
			DoneCount:  result.Summary.DoneCount,
		},
	}
	for i, t := range result.Tasks {
		project.Execution.Tasks[i] = types.ExecutionTask{
			ID:     t.ID,
			Title:  t.Title,
			Status: t.Status,
		}
	}

	// Save project
	project.UpdatedAt = result.SyncedAt
	if err := config.Save(project); err != nil {
		return fmt.Errorf("saving project: %w", err)
	}

	// Print summary
	fmt.Printf("\n✓ Synced %d tasks from %s\n", result.Summary.TotalTasks, result.Target)
	fmt.Printf("  Todo:        %d\n", result.Summary.TodoCount)
	fmt.Printf("  In Progress: %d\n", result.Summary.InProgress)
	fmt.Printf("  Done:        %d\n", result.Summary.DoneCount)

	// Progress bar
	if result.Summary.TotalTasks > 0 {
		pct := float64(result.Summary.DoneCount) / float64(result.Summary.TotalTasks) * 100
		fmt.Printf("  Progress:    %.0f%%\n", pct)
	}

	// Show tasks if verbose
	if verbose && len(result.Tasks) > 0 {
		fmt.Println("\nTasks:")
		for _, t := range result.Tasks {
			statusIcon := "[ ]"
			switch t.Status {
			case "done":
				statusIcon = "[x]"
			case "in_progress":
				statusIcon = "[~]"
			}
			fmt.Printf("  %s %s: %s\n", statusIcon, t.ID, t.Title)
		}
	}

	return nil
}
