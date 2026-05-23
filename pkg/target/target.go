// Package target provides export adapters for different execution systems.
package target

import (
	"fmt"

	"github.com/plexusone/multispec/pkg/types"
)

// Target defines the interface for export adapters.
type Target interface {
	// Name returns the target name.
	Name() string

	// Description returns a description of the target.
	Description() string

	// Capabilities returns what this target supports.
	Capabilities() Capabilities

	// Validate checks if the spec can be exported to this target.
	Validate(spec string) error

	// Export exports the spec to this target.
	Export(spec string, config ExportConfig) (*ExportResult, error)
}

// Capabilities describes what a target supports.
type Capabilities struct {
	SequentialTasks   bool `json:"sequential_tasks"`
	ParallelExecution bool `json:"parallel_execution"`
	MultiAgent        bool `json:"multi_agent"`
	Verification      bool `json:"verification"`
	DependencyGraph   bool `json:"dependency_graph"`
}

// ExportConfig contains configuration for export.
type ExportConfig struct {
	ProjectName string `json:"project_name"`
	OutputDir   string `json:"output_dir"`
	Options     map[string]any `json:"options,omitempty"`
}

// ExportResult contains the result of an export.
type ExportResult struct {
	Target      string   `json:"target"`
	OutputDir   string   `json:"output_dir"`
	Files       []string `json:"files"`
	Success     bool     `json:"success"`
	Message     string   `json:"message,omitempty"`
}

// Registry holds registered targets.
var registry = make(map[string]Target)

// Register adds a target to the registry.
func Register(target Target) {
	registry[target.Name()] = target
}

// Get retrieves a target by name.
func Get(name string) (Target, error) {
	target, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown target: %s", name)
	}
	return target, nil
}

// Available returns all registered target names.
func Available() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}

// ListTargets returns information about all registered targets.
func ListTargets() []TargetInfo {
	var infos []TargetInfo
	for _, target := range registry {
		infos = append(infos, TargetInfo{
			Name:         target.Name(),
			Description:  target.Description(),
			Capabilities: target.Capabilities(),
		})
	}
	return infos
}

// TargetInfo contains information about a target.
type TargetInfo struct {
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	Capabilities Capabilities `json:"capabilities"`
}

// ProjectTargetConfig returns the target config from a project.
func ProjectTargetConfig(project *types.Project, targetName string) *ExportConfig {
	cfg := &ExportConfig{
		ProjectName: project.Name,
		Options:     make(map[string]any),
	}

	switch targetName {
	case "speckit":
		if project.Targets.SpecKit != nil {
			cfg.OutputDir = project.Targets.SpecKit.OutputDir
			cfg.Options["branch_numbering"] = project.Targets.SpecKit.BranchNumbering
		}
	case "gsd":
		if project.Targets.GSD != nil {
			cfg.OutputDir = project.Targets.GSD.OutputDir
			cfg.Options["model_profile"] = project.Targets.GSD.ModelProfile
		}
	case "gastown":
		if project.Targets.GasTown != nil {
			cfg.Options["formula_type"] = project.Targets.GasTown.FormulaType
			cfg.Options["rig"] = project.Targets.GasTown.Rig
		}
	case "gascity":
		if project.Targets.GasCity != nil {
			cfg.OutputDir = project.Targets.GasCity.CityDir
		}
	}

	return cfg
}
