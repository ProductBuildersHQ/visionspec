// Package workflows provides integration with spec-workflows repositories.
//
// This package enables visionspec to load workflows, templates, and rubrics
// from external spec-workflows repositories, supporting organization
// customization through the fork pattern.
//
// Usage:
//
//	// Auto-discover spec-workflows repository
//	repo, err := workflows.DiscoverRepo("")
//
//	// Load from a specific path
//	repo, err := workflows.LoadRepo("/path/to/spec-workflows")
//
//	// Get a workflow
//	workflow, err := repo.GetWorkflow("aws-working-backwards/product")
//
//	// Get loaders for templates and rubrics
//	templateLoader := repo.TemplateLoader("aws-working-backwards")
//	rubricLoader := repo.RubricLoader("aws-working-backwards")
package workflows

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ProductBuildersHQ/visionspec/pkg/rubrics"
	"github.com/ProductBuildersHQ/visionspec/pkg/templates"
)

const (
	// EnvWorkflowsRepo is the environment variable for specifying the workflows repo path.
	EnvWorkflowsRepo = "VISIONSPEC_WORKFLOWS_REPO"

	// DirName is the standard directory name for spec-workflows.
	DirName = "spec-workflows"

	// HiddenDirName is the hidden directory name for spec-workflows.
	HiddenDirName = ".spec-workflows"

	// UserConfigDir is the user-level config directory name.
	UserConfigDir = ".config/visionspec"
)

// DiscoverRepo attempts to find a spec-workflows repository using the following
// search order:
//
//  1. explicitPath parameter (if non-empty)
//  2. VISIONSPEC_WORKFLOWS_REPO environment variable
//  3. Walk up from current directory looking for spec-workflows/ or .spec-workflows/
//  4. ~/.config/visionspec/spec-workflows/
//
// Returns the loaded Repo if found, or an error if no repository is discovered.
func DiscoverRepo(explicitPath string) (*Repo, error) {
	// 1. Explicit path takes highest priority
	if explicitPath != "" {
		return LoadRepo(explicitPath)
	}

	// 2. Environment variable
	if envPath := os.Getenv(EnvWorkflowsRepo); envPath != "" {
		return LoadRepo(envPath)
	}

	// 3. Walk up from current directory
	cwd, err := os.Getwd()
	if err == nil {
		if path := findRepoUpward(cwd); path != "" {
			return LoadRepo(path)
		}
	}

	// 4. User config directory
	if homeDir, err := os.UserHomeDir(); err == nil {
		userPath := filepath.Join(homeDir, UserConfigDir, DirName)
		if isValidRepo(userPath) {
			return LoadRepo(userPath)
		}
	}

	return nil, fmt.Errorf("no spec-workflows repository found. Use --workflows-repo flag, set %s, or clone to ~/.config/visionspec/spec-workflows/", EnvWorkflowsRepo)
}

// DiscoverRepoPath returns the path to a spec-workflows repository without loading it.
// Uses the same search order as DiscoverRepo. Returns empty string if not found.
func DiscoverRepoPath(explicitPath string) string {
	// 1. Explicit path
	if explicitPath != "" && isValidRepo(explicitPath) {
		return explicitPath
	}

	// 2. Environment variable
	if envPath := os.Getenv(EnvWorkflowsRepo); envPath != "" && isValidRepo(envPath) {
		return envPath
	}

	// 3. Walk up from current directory
	if cwd, err := os.Getwd(); err == nil {
		if path := findRepoUpward(cwd); path != "" {
			return path
		}
	}

	// 4. User config directory
	if homeDir, err := os.UserHomeDir(); err == nil {
		userPath := filepath.Join(homeDir, UserConfigDir, DirName)
		if isValidRepo(userPath) {
			return userPath
		}
	}

	return ""
}

// findRepoUpward walks up the directory tree looking for spec-workflows.
func findRepoUpward(startDir string) string {
	dir := startDir
	for {
		// Check for spec-workflows/
		if path := filepath.Join(dir, DirName); isValidRepo(path) {
			return path
		}
		// Check for .spec-workflows/
		if path := filepath.Join(dir, HiddenDirName); isValidRepo(path) {
			return path
		}

		// Move to parent directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root
			return ""
		}
		dir = parent
	}
}

// isValidRepo checks if a path contains a valid spec-workflows repository.
// A valid repo must have a workflows/ directory.
func isValidRepo(path string) bool {
	cleanPath := filepath.Clean(path)
	workflowsDir := filepath.Join(cleanPath, "workflows")
	info, err := os.Stat(workflowsDir) //nolint:gosec // Path is user-provided, traversal is expected
	return err == nil && info.IsDir()
}

// Workflow represents a workflow definition from spec-workflows.
type Workflow struct {
	// Name is the methodology name (e.g., "aws-working-backwards")
	Name string

	// Level is "product" or "feature"
	Level string

	// Description describes the workflow
	Description string

	// EntryPoint is the path to core-workflow.md
	EntryPoint string

	// RuleDetails is the path to rule-details directory
	RuleDetails string

	// Templates is the path to templates directory for this methodology
	Templates string

	// Rubrics is the path to rubrics directory for this methodology
	Rubrics string
}

// ID returns the workflow identifier (e.g., "aws-working-backwards/product")
func (w *Workflow) ID() string {
	return w.Name + "/" + w.Level
}

// Repo represents a spec-workflows repository.
type Repo struct {
	// Path is the local filesystem path to the repo
	Path string

	// Workflows maps workflow IDs to Workflow definitions
	Workflows map[string]*Workflow

	// ruleDetailsPath is the path to the rule-details directory
	ruleDetailsPath string

	// templatesPath is the path to the templates directory
	templatesPath string

	// rubricsPath is the path to the rubrics directory
	rubricsPath string

	// extensionsPath is the path to the extensions directory
	extensionsPath string
}

// LoadRepo loads a spec-workflows repository from the given path.
func LoadRepo(path string) (*Repo, error) {
	// Validate path exists
	cleanPath := filepath.Clean(path)
	if _, err := os.Stat(cleanPath); err != nil { //nolint:gosec // Path is user-provided
		return nil, fmt.Errorf("spec-workflows repo not found at %s: %w", cleanPath, err)
	}

	repo := &Repo{
		Path:            path,
		Workflows:       make(map[string]*Workflow),
		ruleDetailsPath: filepath.Join(path, "rule-details"),
		templatesPath:   filepath.Join(path, "templates"),
		rubricsPath:     filepath.Join(path, "rubrics"),
		extensionsPath:  filepath.Join(path, "extensions"),
	}

	// Scan workflows directory
	workflowsDir := filepath.Join(path, "workflows")
	entries, err := os.ReadDir(workflowsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read workflows directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		methodologyName := entry.Name()
		methodologyDir := filepath.Join(workflowsDir, methodologyName)

		// Check for product and feature subdirectories
		for _, level := range []string{"product", "feature"} {
			levelDir := filepath.Join(methodologyDir, level)
			entryPoint := filepath.Join(levelDir, "core-workflow.md")

			if _, err := os.Stat(entryPoint); err == nil { //nolint:gosec // Path constructed from trusted repo
				workflowID := methodologyName + "/" + level
				repo.Workflows[workflowID] = &Workflow{
					Name:        methodologyName,
					Level:       level,
					EntryPoint:  entryPoint,
					RuleDetails: repo.ruleDetailsPath,
					Templates:   filepath.Join(repo.templatesPath, methodologyName),
					Rubrics:     filepath.Join(repo.rubricsPath, methodologyName),
				}
			}
		}
	}

	return repo, nil
}

// GetWorkflow returns the workflow with the given ID.
func (r *Repo) GetWorkflow(id string) (*Workflow, error) {
	wf, ok := r.Workflows[id]
	if !ok {
		return nil, fmt.Errorf("workflow not found: %s", id)
	}
	return wf, nil
}

// ListWorkflows returns all available workflow IDs.
func (r *Repo) ListWorkflows() []string {
	var ids []string
	for id := range r.Workflows {
		ids = append(ids, id)
	}
	return ids
}

// TemplateLoader returns a template loader for the given methodology.
// Falls back to default templates if methodology-specific templates don't exist.
func (r *Repo) TemplateLoader(methodology string) templates.Loader {
	// Try methodology-specific templates first
	methodologyPath := filepath.Join(r.templatesPath, methodology)
	if _, err := os.Stat(methodologyPath); err == nil {
		return templates.NewChainLoader(
			templates.NewFileLoader(methodologyPath),
			templates.NewFileLoader(filepath.Join(r.templatesPath, "default")),
		)
	}

	// Fall back to default templates
	return templates.NewFileLoader(filepath.Join(r.templatesPath, "default"))
}

// RubricLoader returns a rubric loader for the given methodology.
// Falls back to default rubrics if methodology-specific rubrics don't exist.
func (r *Repo) RubricLoader(methodology string) rubrics.Loader {
	// Try methodology-specific rubrics first
	methodologyPath := filepath.Join(r.rubricsPath, methodology)
	if _, err := os.Stat(methodologyPath); err == nil {
		return rubrics.NewChainLoader(
			rubrics.NewFileLoader(methodologyPath),
			rubrics.NewFileLoader(filepath.Join(r.rubricsPath, "default")),
		)
	}

	// Fall back to default rubrics
	return rubrics.NewFileLoader(filepath.Join(r.rubricsPath, "default"))
}

// RuleDetailsPath returns the path to the rule-details directory.
func (r *Repo) RuleDetailsPath() string {
	return r.ruleDetailsPath
}

// HasExtension checks if an extension exists in the repo.
func (r *Repo) HasExtension(name string) bool {
	extPath := filepath.Join(r.extensionsPath, name)
	_, err := os.Stat(extPath)
	return err == nil
}

// ExtensionPath returns the path to an extension directory.
func (r *Repo) ExtensionPath(name string) string {
	return filepath.Join(r.extensionsPath, name)
}
