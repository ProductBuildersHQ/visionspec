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
	"github.com/ProductBuildersHQ/visionspec/pkg/profiles"
	"github.com/ProductBuildersHQ/visionspec/pkg/rubrics"
	"github.com/ProductBuildersHQ/visionspec/pkg/templates"
	"github.com/ProductBuildersHQ/visionspec/pkg/types"
	"github.com/spf13/cobra"
)

// Config allows customization of CLI behavior.
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

	// DefaultProfile is the profile to use when none is specified.
	// If empty, uses no profile (default visionspec behavior).
	DefaultProfile string

	// Version is the CLI version string.
	Version string
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		TemplateLoader: templates.DefaultLoader(),
		RubricLoader:   rubrics.DefaultLoader(),
		SpecConfig:     types.DefaultSpecConfig(),
		ProfileLoader:  profiles.DefaultLoader(),
		Version:        "0.3.0",
	}
}

// ConfigFromProfile creates a Config from a profile.
func ConfigFromProfile(profile *profiles.Profile) *Config {
	return &Config{
		TemplateLoader: profile.GetTemplateLoader(),
		RubricLoader:   profile.GetRubricLoader(),
		SpecConfig:     profile.GetSpecConfig(),
		ProfileLoader:  profiles.DefaultLoader(),
		Version:        "0.3.0",
	}
}

// GetSpecConfig returns the SpecConfig, falling back to defaults if nil.
func (c *Config) GetSpecConfig() *types.SpecConfig {
	if c == nil || c.SpecConfig == nil {
		return types.DefaultSpecConfig()
	}
	return c.SpecConfig
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
	)
}

// CommandSet contains all visionspec commands.
type CommandSet struct {
	Init       *cobra.Command
	Create     *cobra.Command
	Lint       *cobra.Command
	Status     *cobra.Command
	Eval       *cobra.Command
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
	}
}
