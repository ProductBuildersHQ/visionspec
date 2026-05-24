// Package config handles configuration loading for multispec.
package config

import (
	"os"
	"path/filepath"

	"github.com/plexusone/multispec/pkg/types"
	"gopkg.in/yaml.v3"
)

const (
	// ConfigFileName is the canonical config file name.
	ConfigFileName = "multispec.yaml"

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

// Load loads the project configuration from multispec.yaml.
func Load(projectPath string) (*types.Project, error) {
	configPath := filepath.Join(projectPath, ConfigFileName)

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return empty project if config doesn't exist
			return &types.Project{
				Path: projectPath,
				Name: filepath.Base(projectPath),
			}, nil
		}
		return nil, err
	}

	var project types.Project
	if err := yaml.Unmarshal(data, &project); err != nil {
		return nil, err
	}

	project.Path = projectPath
	return &project, nil
}

// Save saves the project configuration to multispec.yaml.
func Save(project *types.Project) error {
	configPath := filepath.Join(project.Path, ConfigFileName)

	data, err := yaml.Marshal(project)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600)
}

// FindProjectRoot finds the project root by looking for multispec.yaml.
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
