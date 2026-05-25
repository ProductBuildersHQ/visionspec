// Package config handles configuration loading for visionspec.
package config

import (
	"os"
	"path/filepath"

	"github.com/ProductBuildersHQ/visionspec/pkg/types"
	"gopkg.in/yaml.v3"
)

const (
	// ConfigFileName is the canonical config file name.
	ConfigFileName = "visionspec.yaml"

	// LegacyConfigFileName supports migration from multispec.
	LegacyConfigFileName = "multispec.yaml"

	// SpecsDir is the canonical specs directory name.
	SpecsDir = "docs/specs"

	// ConstitutionFile is the repo-level constitution file name.
	ConstitutionFile = "CONSTITUTION.md"

	// RoadmapFile is the repo-level roadmap file name.
	RoadmapFile = "ROADMAP.md"
)

// Subdirectories within a project.
const (
	SourceDir    = "source"
	GTMDir       = "gtm"
	TechnicalDir = "technical"
	EvalDir      = "eval"
)

// Load loads the project configuration from visionspec.yaml (or legacy multispec.yaml).
func Load(projectPath string) (*types.Project, error) {
	configPath := filepath.Join(projectPath, ConfigFileName)

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Try legacy config file for backward compatibility
			legacyPath := filepath.Join(projectPath, LegacyConfigFileName)
			data, err = os.ReadFile(legacyPath)
			if err != nil {
				if os.IsNotExist(err) {
					// Return empty project if no config exists
					return &types.Project{
						Path: projectPath,
						Name: filepath.Base(projectPath),
					}, nil
				}
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	var project types.Project
	if err := yaml.Unmarshal(data, &project); err != nil {
		return nil, err
	}

	project.Path = projectPath
	return &project, nil
}

// Save saves the project configuration to visionspec.yaml.
func Save(project *types.Project) error {
	configPath := filepath.Join(project.Path, ConfigFileName)

	data, err := yaml.Marshal(project)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600)
}

// FindProjectRoot finds the project root by looking for visionspec.yaml (or legacy multispec.yaml).
func FindProjectRoot(startPath string) (string, error) {
	path := startPath
	for {
		configPath := filepath.Join(path, ConfigFileName)
		if _, err := os.Stat(configPath); err == nil {
			return path, nil
		}

		parent := filepath.Dir(path)
		if parent == path {
			return "", os.ErrNotExist
		}
		path = parent
	}
}

// FindSpecsDir finds the specs directory from the current path.
func FindSpecsDir(startPath string) (string, error) {
	// Look for docs/specs directory
	path := startPath
	for {
		specsPath := filepath.Join(path, SpecsDir)
		if info, err := os.Stat(specsPath); err == nil && info.IsDir() {
			return specsPath, nil
		}

		parent := filepath.Dir(path)
		if parent == path {
			return "", os.ErrNotExist
		}
		path = parent
	}
}

// ProjectPath returns the full path for a project within the specs directory.
func ProjectPath(specsDir, projectName string) string {
	return filepath.Join(specsDir, projectName)
}

// SpecPath returns the full path for a spec file within a project.
func SpecPath(projectPath string, specType types.SpecType) string {
	dir := specType.Dir()
	if dir == "" {
		return filepath.Join(projectPath, specType.Filename())
	}
	return filepath.Join(projectPath, dir, specType.Filename())
}

// EvalPath returns the full path for an eval file within a project.
func EvalPath(projectPath string, specType types.SpecType) string {
	return filepath.Join(projectPath, EvalDir, specType.EvalFilename())
}

// FindConstitution finds the constitution file from multiple locations.
// Search order (first found wins):
// 1. Repo-level: docs/specs/CONSTITUTION.md (from project path)
// 2. Org-level: ~/.config/visionspec/CONSTITUTION.md
// Returns the path if found, empty string otherwise.
func FindConstitution(projectPath string) string {
	// Try repo-level constitution
	specsDir := filepath.Dir(projectPath) // Go up from project to docs/specs
	repoConstitution := filepath.Join(specsDir, ConstitutionFile)
	if _, err := os.Stat(repoConstitution); err == nil {
		return repoConstitution
	}

	// Try org-level constitution
	homeDir, err := os.UserHomeDir()
	if err == nil {
		orgConstitution := filepath.Join(homeDir, ".config", "visionspec", ConstitutionFile)
		if _, err := os.Stat(orgConstitution); err == nil {
			return orgConstitution
		}
	}

	return ""
}

// LoadConstitution loads the constitution content from the first found location.
// Returns empty string if no constitution file exists.
func LoadConstitution(projectPath string) string {
	constitutionPath := FindConstitution(projectPath)
	if constitutionPath == "" {
		return ""
	}

	content, err := os.ReadFile(constitutionPath)
	if err != nil {
		return ""
	}

	return string(content)
}
