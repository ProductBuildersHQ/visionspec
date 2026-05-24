package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/plexusone/multispec/pkg/types"
)

func TestLoadMissingConfig(t *testing.T) {
	tmpDir := t.TempDir()

	project, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if project.Path != tmpDir {
		t.Errorf("Path = %v, want %v", project.Path, tmpDir)
	}

	if project.Name != filepath.Base(tmpDir) {
		t.Errorf("Name = %v, want %v", project.Name, filepath.Base(tmpDir))
	}
}

func TestLoadExistingConfig(t *testing.T) {
	tmpDir := t.TempDir()

	configContent := `name: test-project
description: A test project
`
	err := os.WriteFile(filepath.Join(tmpDir, ConfigFileName), []byte(configContent), 0600)
	if err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	project, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if project.Name != "test-project" {
		t.Errorf("Name = %v, want test-project", project.Name)
	}

	if project.Path != tmpDir {
		t.Errorf("Path = %v, want %v", project.Path, tmpDir)
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()

	project := &types.Project{
		Name: "my-project",
		Path: tmpDir,
		LLM: &types.LLMConfig{
			Provider: "anthropic",
			Model:    "claude-sonnet-4-20250514",
		},
	}

	err := Save(project)
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	loaded, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if loaded.Name != project.Name {
		t.Errorf("Name = %v, want %v", loaded.Name, project.Name)
	}

	if loaded.LLM == nil {
		t.Fatal("LLM config not loaded")
	}

	if loaded.LLM.Provider != "anthropic" {
		t.Errorf("LLM.Provider = %v, want anthropic", loaded.LLM.Provider)
	}
}

func TestFindProjectRoot(t *testing.T) {
	tmpDir := t.TempDir()

	// Create nested directory structure
	projectDir := filepath.Join(tmpDir, "myproject")
	subDir := filepath.Join(projectDir, "source")
	err := os.MkdirAll(subDir, 0750)
	if err != nil {
		t.Fatalf("Failed to create dirs: %v", err)
	}

	// Create config file in project dir
	err = os.WriteFile(filepath.Join(projectDir, ConfigFileName), []byte("name: myproject\n"), 0600)
	if err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Find from subdirectory
	root, err := FindProjectRoot(subDir)
	if err != nil {
		t.Fatalf("FindProjectRoot() failed: %v", err)
	}

	if root != projectDir {
		t.Errorf("FindProjectRoot() = %v, want %v", root, projectDir)
	}
}

func TestFindProjectRootNotFound(t *testing.T) {
	tmpDir := t.TempDir()

	_, err := FindProjectRoot(tmpDir)
	if !os.IsNotExist(err) {
		t.Errorf("FindProjectRoot() error = %v, want os.ErrNotExist", err)
	}
}

func TestFindSpecsDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Create docs/specs structure
	specsDir := filepath.Join(tmpDir, SpecsDir)
	err := os.MkdirAll(specsDir, 0750)
	if err != nil {
		t.Fatalf("Failed to create dirs: %v", err)
	}

	// Create a project subdirectory
	projectDir := filepath.Join(specsDir, "my-project", "source")
	err = os.MkdirAll(projectDir, 0750)
	if err != nil {
		t.Fatalf("Failed to create project dir: %v", err)
	}

	// Find from project subdirectory
	found, err := FindSpecsDir(projectDir)
	if err != nil {
		t.Fatalf("FindSpecsDir() failed: %v", err)
	}

	if found != specsDir {
		t.Errorf("FindSpecsDir() = %v, want %v", found, specsDir)
	}
}

func TestFindSpecsDirNotFound(t *testing.T) {
	tmpDir := t.TempDir()

	_, err := FindSpecsDir(tmpDir)
	if !os.IsNotExist(err) {
		t.Errorf("FindSpecsDir() error = %v, want os.ErrNotExist", err)
	}
}

func TestProjectPath(t *testing.T) {
	specsDir := filepath.Join("docs", "specs")
	got := ProjectPath(specsDir, "my-project")
	want := filepath.Join("docs", "specs", "my-project")
	if got != want {
		t.Errorf("ProjectPath() = %v, want %v", got, want)
	}
}

func TestSpecPath(t *testing.T) {
	projectDir := "project"
	tests := []struct {
		specType types.SpecType
		expected string
	}{
		{types.SpecTypeMRD, filepath.Join(projectDir, "source", "mrd.md")},
		{types.SpecTypePRD, filepath.Join(projectDir, "source", "prd.md")},
		{types.SpecTypeUXD, filepath.Join(projectDir, "source", "uxd.md")},
		{types.SpecTypePress, filepath.Join(projectDir, "gtm", "press.md")},
		{types.SpecTypeFAQ, filepath.Join(projectDir, "gtm", "faq.md")},
		{types.SpecTypeTRD, filepath.Join(projectDir, "technical", "trd.md")},
		{types.SpecTypeIRD, filepath.Join(projectDir, "technical", "ird.md")},
		{types.SpecTypeSpec, filepath.Join(projectDir, "spec.md")},
	}

	for _, tt := range tests {
		t.Run(string(tt.specType), func(t *testing.T) {
			got := SpecPath(projectDir, tt.specType)
			if got != tt.expected {
				t.Errorf("SpecPath() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestEvalPath(t *testing.T) {
	projectDir := "project"
	tests := []struct {
		specType types.SpecType
		expected string
	}{
		{types.SpecTypeMRD, filepath.Join(projectDir, "eval", "mrd.eval.json")},
		{types.SpecTypePRD, filepath.Join(projectDir, "eval", "prd.eval.json")},
		{types.SpecTypeTRD, filepath.Join(projectDir, "eval", "trd.eval.json")},
	}

	for _, tt := range tests {
		t.Run(string(tt.specType), func(t *testing.T) {
			got := EvalPath(projectDir, tt.specType)
			if got != tt.expected {
				t.Errorf("EvalPath() = %v, want %v", got, tt.expected)
			}
		})
	}
}
