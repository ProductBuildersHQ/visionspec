// Package cli provides a composable CLI for visionspec.
//
// Organizations can import this package to build custom CLI tools
// that include visionspec commands alongside their own:
//
//	package main
//
//	import (
//		"github.com/spf13/cobra"
//		"github.com/ProductBuildersHQ/visionspec/pkg/cli"
//	)
//
//	func main() {
//		root := &cobra.Command{Use: "org-spec"}
//		cfg := cli.DefaultConfig()
//		cli.AddCommandsTo(root, cfg)
//		root.AddCommand(customCmd)
//		root.Execute()
//	}
package cli

import (
	"github.com/ProductBuildersHQ/visionspec/pkg/apptypes"
	"github.com/ProductBuildersHQ/visionspec/pkg/constitution"
	"github.com/ProductBuildersHQ/visionspec/pkg/profiles"
	"github.com/ProductBuildersHQ/visionspec/pkg/rubrics"
	"github.com/ProductBuildersHQ/visionspec/pkg/templates"
	"github.com/ProductBuildersHQ/visionspec/pkg/types"
	"github.com/ProductBuildersHQ/visionspec/pkg/workflows"
	"github.com/spf13/cobra"
)

// Config allows customization of CLI behavior.
// Organizations can provide custom loaders to override defaults:
//
//	cfg := cli.DefaultConfig()
//	cfg.TemplateLoader = templates.NewChainLoader(
//		orgTemplates,     // Organization-specific (prescriptive)
//		cfg.TemplateLoader, // Fall back to visionspec defaults
//	)
//	cfg.ConstitutionLoader = constitution.NewChainLoader(
//		orgConstitutions,
//		cfg.ConstitutionLoader,
//	)
type Config struct {
	// TemplateLoader loads spec templates.
	// If nil, uses embedded templates.
	TemplateLoader templates.Loader

	// RubricLoader loads evaluation rubrics.
	// If nil, uses embedded rubrics.
	RubricLoader rubrics.Loader

	// SpecConfig defines which specs are required and their settings.
	// If nil, uses default visionspec requirements.
	SpecConfig *types.SpecConfig

	// ProfileLoader loads configuration profiles.
	// If nil, uses default profiles.
	ProfileLoader profiles.Loader

	// ConstitutionLoader loads organization/team/project constitutions.
	// If nil, no built-in constitutions are available.
	// Organizations typically provide their own loader with prescriptive defaults.
	ConstitutionLoader constitution.Loader

	// AppTypeLoader loads app type specifications.
	// If nil, uses built-in app type specs.
	AppTypeLoader apptypes.Loader

	// DefaultProfile is the profile to use when none is specified.
	// If empty, uses no profile (default visionspec behavior).
	DefaultProfile string

	// WorkflowsRepoPath is the path to a spec-workflows repository.
	// If set, templates and rubrics are loaded from this repo.
	WorkflowsRepoPath string

	// Version is the CLI version string.
	Version string
}

// DefaultConfig returns the default configuration.
// This provides flexible, choice-based defaults suitable for open source use.
// Organizations should use NewOrgConfig() or customize loaders for prescriptive defaults.
func DefaultConfig() *Config {
	return &Config{
		TemplateLoader:     templates.DefaultLoader(),
		RubricLoader:       rubrics.DefaultLoader(),
		SpecConfig:         types.DefaultSpecConfig(),
		ProfileLoader:      profiles.DefaultLoader(),
		ConstitutionLoader: nil, // No built-in constitutions; orgs provide their own
		AppTypeLoader:      apptypes.DefaultLoader(),
		Version:            "0.3.0",
	}
}

// ConfigFromProfile creates a Config from a profile.
func ConfigFromProfile(profile *profiles.Profile) *Config {
	return &Config{
		TemplateLoader:     profile.GetTemplateLoader(),
		RubricLoader:       profile.GetRubricLoader(),
		SpecConfig:         profile.GetSpecConfig(),
		ProfileLoader:      profiles.DefaultLoader(),
		ConstitutionLoader: nil, // Profiles don't include constitutions yet
		AppTypeLoader:      apptypes.DefaultLoader(),
		Version:            "0.3.0",
	}
}

// GetSpecConfig returns the SpecConfig, falling back to defaults if nil.
func (c *Config) GetSpecConfig() *types.SpecConfig {
	if c == nil || c.SpecConfig == nil {
		return types.DefaultSpecConfig()
	}
	return c.SpecConfig
}

// GetConstitutionLoader returns the ConstitutionLoader.
// Returns nil if no loader is configured (organizations must provide their own).
func (c *Config) GetConstitutionLoader() constitution.Loader {
	if c == nil {
		return nil
	}
	return c.ConstitutionLoader
}

// GetAppTypeLoader returns the AppTypeLoader, falling back to built-in specs if nil.
func (c *Config) GetAppTypeLoader() apptypes.Loader {
	if c == nil || c.AppTypeLoader == nil {
		return apptypes.DefaultLoader()
	}
	return c.AppTypeLoader
}

// GetWorkflowsRepo loads and returns the spec-workflows repository.
// Uses auto-discovery if no explicit path is configured.
// Returns nil, nil if no repository is found (not an error for optional usage).
func (c *Config) GetWorkflowsRepo() (*workflows.Repo, error) {
	if c == nil {
		return nil, nil
	}
	// Use discovery with explicit path (empty string triggers auto-discovery)
	repo, err := workflows.DiscoverRepo(c.WorkflowsRepoPath)
	if err != nil {
		// If discovery fails and no explicit path was set, return nil (optional)
		if c.WorkflowsRepoPath == "" {
			return nil, nil
		}
		return nil, err
	}
	return repo, nil
}

// GetWorkflowsRepoPath returns the path to the spec-workflows repository.
// Uses auto-discovery if no explicit path is configured.
// Returns empty string if no repository is found.
func (c *Config) GetWorkflowsRepoPath() string {
	if c == nil {
		return ""
	}
	return workflows.DiscoverRepoPath(c.WorkflowsRepoPath)
}

// GetTemplateLoaderForWorkflow returns a template loader for a specific workflow.
// If a workflows repo is found (explicit or auto-discovered), it creates a chain
// loader that prioritizes the workflow's methodology-specific templates, then
// falls back to defaults.
func (c *Config) GetTemplateLoaderForWorkflow(methodology string) templates.Loader {
	if c == nil {
		return templates.DefaultLoader()
	}

	// Try to get workflows repo (uses auto-discovery)
	repo, err := c.GetWorkflowsRepo()
	if err == nil && repo != nil {
		workflowLoader := repo.TemplateLoader(methodology)
		if c.TemplateLoader != nil {
			return templates.NewChainLoader(workflowLoader, c.TemplateLoader)
		}
		return workflowLoader
	}

	// Fall back to configured loader or defaults
	if c.TemplateLoader != nil {
		return c.TemplateLoader
	}
	return templates.DefaultLoader()
}

// GetRubricLoaderForWorkflow returns a rubric loader for a specific workflow.
// If a workflows repo is found (explicit or auto-discovered), it creates a chain
// loader that prioritizes the workflow's methodology-specific rubrics, then
// falls back to defaults.
func (c *Config) GetRubricLoaderForWorkflow(methodology string) rubrics.Loader {
	if c == nil {
		return rubrics.DefaultLoader()
	}

	// Try to get workflows repo (uses auto-discovery)
	repo, err := c.GetWorkflowsRepo()
	if err == nil && repo != nil {
		workflowLoader := repo.RubricLoader(methodology)
		if c.RubricLoader != nil {
			return rubrics.NewChainLoader(workflowLoader, c.RubricLoader)
		}
		return workflowLoader
	}

	// Fall back to configured loader or defaults
	if c.RubricLoader != nil {
		return c.RubricLoader
	}
	return rubrics.DefaultLoader()
}

// ListAvailableWorkflows returns all available workflow IDs from the workflows repo.
// Uses auto-discovery if no explicit path is configured.
// Returns nil if no workflows repo is found.
func (c *Config) ListAvailableWorkflows() []string {
	if c == nil {
		return nil
	}
	repo, err := c.GetWorkflowsRepo()
	if err != nil || repo == nil {
		return nil
	}
	return repo.ListWorkflows()
}

// AddCommandsTo adds all visionspec commands to a root command.
func AddCommandsTo(root *cobra.Command, cfg *Config) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	// Store config in root for subcommands to access
	root.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		cmd.SetContext(WithConfig(cmd.Context(), cfg))
	}

	// Add all commands
	cmds := Commands(cfg)
	root.AddCommand(
		cmds.Init,
		cmds.Create,
		cmds.Lint,
		cmds.Status,
		cmds.Eval,
		cmds.Render,
		cmds.Synthesize,
		cmds.Reconcile,
		cmds.Approve,
		cmds.Export,
		cmds.Targets,
		cmds.Graph,
		cmds.Docs,
		cmds.Serve,
		cmds.Profiles,
		cmds.Context,
		cmds.Rules,
		cmds.Generate,
		cmds.Sync,
		cmds.Drift,
		cmds.Align,
		cmds.Hooks,
		cmds.Watch,
		cmds.Workflows,
		cmds.Metrics,
		cmds.Version,
	)
}

// CommandSet contains all visionspec commands.
type CommandSet struct {
	Init       *cobra.Command
	Create     *cobra.Command
	Lint       *cobra.Command
	Status     *cobra.Command
	Eval       *cobra.Command
	Render     *cobra.Command
	Synthesize *cobra.Command
	Reconcile  *cobra.Command
	Approve    *cobra.Command
	Export     *cobra.Command
	Targets    *cobra.Command
	Graph      *cobra.Command
	Docs       *cobra.Command
	Serve      *cobra.Command
	Profiles   *cobra.Command
	Context    *cobra.Command
	Rules      *cobra.Command
	Generate   *cobra.Command
	Sync       *cobra.Command
	Drift      *cobra.Command
	Align      *cobra.Command
	Hooks      *cobra.Command
	Watch      *cobra.Command
	Workflows  *cobra.Command
	Metrics    *cobra.Command
	Version    *cobra.Command
}

// Commands returns all visionspec commands.
// Use this for selective command inclusion.
func Commands(cfg *Config) *CommandSet {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	return &CommandSet{
		Init:       initCmd(cfg),
		Create:     createCmd(cfg),
		Lint:       lintCmd(cfg),
		Status:     statusCmd(cfg),
		Eval:       evalCmd(cfg),
		Render:     renderCmd(cfg),
		Synthesize: synthesizeCmd(cfg),
		Reconcile:  reconcileCmd(cfg),
		Approve:    approveCmd(cfg),
		Export:     exportCmd(cfg),
		Targets:    targetsCmd(cfg),
		Graph:      graphCmd(cfg),
		Docs:       docsCmd(cfg),
		Serve:      serveCmd(cfg),
		Profiles:   profilesCmd(cfg),
		Context:    contextCmd(cfg),
		Rules:      rulesCmd(cfg),
		Generate:   generateCmd(cfg),
		Sync:       syncCmd(cfg),
		Drift:      driftCmd(cfg),
		Align:      alignCmd(cfg),
		Hooks:      hooksCmd(cfg),
		Watch:      watchCmd(cfg),
		Workflows:  workflowsCmd(cfg),
		Metrics:    metricsCmd(cfg),
		Version:    versionCmd(cfg),
	}
}
