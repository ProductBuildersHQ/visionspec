package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/plexusone/multispec/pkg/config"
	"github.com/plexusone/multispec/pkg/types"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var initCmd = &cobra.Command{
	Use:   "init <project-name>",
	Short: "Initialize a new multispec project",
	Long: `Initialize a new multispec project with the canonical directory structure.

The project name must be kebab-case (lowercase with hyphens).

Creates:
  docs/specs/<project>/
  ├── source/          # Human-authored specs (mrd, prd, uxd)
  ├── gtm/             # LLM-generated GTM docs (press, faq, narrative)
  ├── technical/       # LLM-generated technical docs (trd, ird)
  ├── eval/            # Evaluation results
  └── multispec.yaml   # Project configuration`,
	Args: cobra.ExactArgs(1),
	RunE: runInit,
}

var kebabCaseRegex = regexp.MustCompile(`^[a-z][a-z0-9]*(-[a-z0-9]+)*$`)

func init() {
	initCmd.Flags().String("constitution", "", "Path to constitution file (relative or absolute)")
	initCmd.Flags().Bool("with-templates", false, "Create template spec files")
}

func runInit(cmd *cobra.Command, args []string) error {
	projectName := args[0]

	// Validate project name is kebab-case
	if !kebabCaseRegex.MatchString(projectName) {
		return fmt.Errorf("invalid project name %q: must be kebab-case (e.g., 'user-onboarding')", projectName)
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
		if err := createTemplateFiles(projectPath); err != nil {
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

func createTemplateFiles(projectPath string) error {
	templates := map[string]string{
		filepath.Join(projectPath, config.SourceDir, "mrd.md"): mrdTemplate,
		filepath.Join(projectPath, config.SourceDir, "prd.md"): prdTemplate,
		filepath.Join(projectPath, config.SourceDir, "uxd.md"): uxdTemplate,
	}

	for path, content := range templates {
		if err := os.WriteFile(path, []byte(strings.TrimSpace(content)+"\n"), 0600); err != nil {
			return err
		}
		fmt.Printf("Created template: %s\n", filepath.Base(path))
	}

	return nil
}

const mrdTemplate = `
# Market Requirements Document

## Overview

<!-- Brief description of the market opportunity -->

## Target Market

### Primary Audience

<!-- Who is the primary customer? -->

### Secondary Audience

<!-- Who else might use this? -->

## Market Problem

<!-- What problem exists in the market? -->

## Competitive Landscape

<!-- Who are the competitors? What are their strengths/weaknesses? -->

## Business Metrics

### Success Criteria

<!-- How will we measure success? -->

### Key Performance Indicators

| Metric | Baseline | Target |
|--------|----------|--------|
| | | |

## Assumptions

<!-- What are we assuming to be true? -->

## Risks

<!-- What could go wrong? -->
`

const prdTemplate = `
# Product Requirements Document

## Overview

<!-- Brief description of the product/feature -->

## User Stories

### US-001: [Story Title]

**As a** [user type]
**I want** [goal]
**So that** [benefit]

**Acceptance Criteria:**

- [ ] Criterion 1
- [ ] Criterion 2

## Functional Requirements

### FR-001: [Requirement Title]

<!-- Detailed requirement description -->

**Priority:** High | Medium | Low
**Dependencies:** None

## Non-Functional Requirements

### NFR-001: Performance

<!-- Performance requirements -->

### NFR-002: Accessibility

<!-- Accessibility requirements -->

## Out of Scope

<!-- What is explicitly NOT included -->

## Open Questions

<!-- Questions that need to be resolved -->
`

const uxdTemplate = `
# User Experience Design

## Overview

<!-- Brief description of the user experience goals -->

## User Journeys

### Journey 1: [Journey Name]

1. **Entry Point:** How does the user arrive?
2. **Steps:**
   - Step 1
   - Step 2
3. **Exit Point:** What is the successful outcome?

## Interaction Flows

### Flow 1: [Flow Name]

<!-- Describe the interaction flow -->

## Accessibility Requirements

### WCAG Compliance

- Level: AA
- Requirements:
  - [ ] Keyboard navigation
  - [ ] Screen reader support
  - [ ] Color contrast

## Design Principles

1. **Principle 1:** Description
2. **Principle 2:** Description

## Open Questions

<!-- UX questions that need to be resolved -->
`
