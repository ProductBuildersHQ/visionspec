package rubrics

import (
	"embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/plexusone/multispec/pkg/types"
)

//go:embed testdata/*.rubric.yaml
var testRubrics embed.FS

func TestEmbeddedRubricLoader(t *testing.T) {
	loader := EmbeddedLoader()

	// Test loading a known rubric
	rs, err := loader.Load(types.SpecTypePRD)
	if err != nil {
		t.Fatalf("Load(prd) failed: %v", err)
	}

	if rs.SpecType != types.SpecTypePRD {
		t.Errorf("SpecType = %v, want %v", rs.SpecType, types.SpecTypePRD)
	}

	if rs.Name == "" {
		t.Error("Name is empty")
	}

	// Test available rubrics
	available := loader.Available()
	if len(available) == 0 {
		t.Error("Available() returned empty list")
	}
}

func TestFileRubricLoader(t *testing.T) {
	// Create temp directory with test rubric
	tmpDir := t.TempDir()

	// Write a custom rubric
	rubricContent := `
spec_type: security
name: "Security Rubric"
description: "Custom security rubric"
version: "1.0"

categories:
  - id: threat-modeling
    name: "Threat Modeling"
    description: "Are threats identified?"
    weight: 2.0
    required: true
    criteria:
      pass: "Comprehensive threat model"
      partial: "Basic threats identified"
      fail: "No threat analysis"

pass_criteria:
  require_all_pass: false
  max_critical: 0
  max_high: 0
  max_medium: 3
`
	err := os.WriteFile(filepath.Join(tmpDir, "security.rubric.yaml"), []byte(rubricContent), 0600)
	if err != nil {
		t.Fatalf("Failed to write rubric: %v", err)
	}

	loader := NewFileLoader(tmpDir)

	// Test loading custom rubric
	rs, err := loader.Load(types.SpecType("security"))
	if err != nil {
		t.Fatalf("Load(security) failed: %v", err)
	}

	if rs.Name != "Security Rubric" {
		t.Errorf("Name = %q, want %q", rs.Name, "Security Rubric")
	}

	if len(rs.Categories) != 1 {
		t.Errorf("Categories count = %d, want 1", len(rs.Categories))
	}

	if rs.Categories[0].ID != "threat-modeling" {
		t.Errorf("Category ID = %q, want %q", rs.Categories[0].ID, "threat-modeling")
	}

	// Test loading non-existent rubric
	_, err = loader.Load(types.SpecType("nonexistent"))
	if err == nil {
		t.Error("Expected error for non-existent rubric")
	}

	// Test available rubrics
	available := loader.Available()
	if len(available) != 1 {
		t.Errorf("Available() returned %d rubrics, want 1", len(available))
	}
}

func TestChainRubricLoader(t *testing.T) {
	// Create temp directory with override rubric
	tmpDir := t.TempDir()

	// Write a custom PRD rubric that overrides the embedded one
	rubricContent := `
spec_type: prd
name: "Custom PRD Rubric"
description: "Overrides the default"
version: "2.0"

categories:
  - id: custom-category
    name: "Custom Category"
    description: "A custom category"
    weight: 1.0
    required: true
    criteria:
      pass: "Custom pass"
      partial: "Custom partial"
      fail: "Custom fail"

pass_criteria:
  require_all_pass: false
  max_critical: 0
  max_high: 0
  max_medium: 5
`
	err := os.WriteFile(filepath.Join(tmpDir, "prd.rubric.yaml"), []byte(rubricContent), 0600)
	if err != nil {
		t.Fatalf("Failed to write rubric: %v", err)
	}

	// Chain: file loader first (override), then embedded (fallback)
	chain := NewChainLoader(
		NewFileLoader(tmpDir),
		EmbeddedLoader(),
	)

	// PRD should come from file loader (override)
	rs, err := chain.Load(types.SpecTypePRD)
	if err != nil {
		t.Fatalf("Load(prd) failed: %v", err)
	}

	if rs.Name != "Custom PRD Rubric" {
		t.Error("Chain loader did not use file loader for override")
	}

	// MRD should come from embedded loader (fallback)
	rs, err = chain.Load(types.SpecTypeMRD)
	if err != nil {
		t.Fatalf("Load(mrd) failed: %v", err)
	}

	if rs.Name == "" {
		t.Error("Chain loader did not fall back to embedded loader")
	}

	// Available should include both
	available := chain.Available()
	if len(available) < 2 {
		t.Errorf("Available() returned %d rubrics, expected at least 2", len(available))
	}
}

func TestRubricYAMLValidation(t *testing.T) {
	tests := []struct {
		name        string
		yaml        RubricYAML
		expectError bool
	}{
		{
			name: "valid rubric",
			yaml: RubricYAML{
				SpecType: "test",
				Name:     "Test Rubric",
				Categories: []CategoryYAML{
					{ID: "cat1", Name: "Category 1"},
				},
			},
			expectError: false,
		},
		{
			name: "missing spec_type",
			yaml: RubricYAML{
				Name: "Test Rubric",
				Categories: []CategoryYAML{
					{ID: "cat1", Name: "Category 1"},
				},
			},
			expectError: true,
		},
		{
			name: "missing name",
			yaml: RubricYAML{
				SpecType: "test",
				Categories: []CategoryYAML{
					{ID: "cat1", Name: "Category 1"},
				},
			},
			expectError: true,
		},
		{
			name: "no categories",
			yaml: RubricYAML{
				SpecType:   "test",
				Name:       "Test Rubric",
				Categories: []CategoryYAML{},
			},
			expectError: true,
		},
		{
			name: "category missing id",
			yaml: RubricYAML{
				SpecType: "test",
				Name:     "Test Rubric",
				Categories: []CategoryYAML{
					{Name: "Category 1"},
				},
			},
			expectError: true,
		},
		{
			name: "category missing name",
			yaml: RubricYAML{
				SpecType: "test",
				Name:     "Test Rubric",
				Categories: []CategoryYAML{
					{ID: "cat1"},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.yaml.ToRubricSet()
			if tt.expectError && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestRubricSetToYAML(t *testing.T) {
	rs := &RubricSet{
		SpecType:    types.SpecTypePRD,
		Name:        "Test Rubric",
		Description: "A test rubric",
		Categories: []*Category{
			{
				ID:          "cat1",
				Name:        "Category 1",
				Description: "First category",
				Weight:      1.5,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Pass criteria",
					Partial: "Partial criteria",
					Fail:    "Fail criteria",
				},
			},
		},
		PassCriteria: PassCriteria{
			RequireAllPass: true,
			MaxCritical:    0,
			MaxHigh:        1,
			MaxMedium:      3,
		},
	}

	yaml := rs.ToYAML()

	if yaml.SpecType != "prd" {
		t.Errorf("SpecType = %q, want %q", yaml.SpecType, "prd")
	}

	if yaml.Name != "Test Rubric" {
		t.Errorf("Name = %q, want %q", yaml.Name, "Test Rubric")
	}

	if len(yaml.Categories) != 1 {
		t.Fatalf("Categories count = %d, want 1", len(yaml.Categories))
	}

	if yaml.Categories[0].ID != "cat1" {
		t.Errorf("Category ID = %q, want %q", yaml.Categories[0].ID, "cat1")
	}

	if yaml.PassCriteria.MaxHigh != 1 {
		t.Errorf("MaxHigh = %d, want 1", yaml.PassCriteria.MaxHigh)
	}
}

func TestEmbedFSRubricLoader(t *testing.T) {
	loader := NewEmbedFSLoader(testRubrics, "testdata")

	// Test loading a rubric from embedded FS
	rs, err := loader.Load(types.SpecType("custom"))
	if err != nil {
		t.Fatalf("Load(custom) failed: %v", err)
	}

	if rs.SpecType != types.SpecType("custom") {
		t.Errorf("SpecType = %v, want custom", rs.SpecType)
	}

	if rs.Name != "Custom Test Rubric" {
		t.Errorf("Name = %q, want %q", rs.Name, "Custom Test Rubric")
	}

	if len(rs.Categories) != 1 {
		t.Errorf("Categories count = %d, want 1", len(rs.Categories))
	}

	// Test loading non-existent rubric
	_, err = loader.Load(types.SpecType("nonexistent"))
	if err == nil {
		t.Error("Expected error for non-existent rubric")
	}

	// Test available rubrics
	available := loader.Available()
	if len(available) != 1 {
		t.Errorf("Available() returned %d rubrics, want 1", len(available))
	}
	if available[0] != types.SpecType("custom") {
		t.Errorf("Available()[0] = %v, want custom", available[0])
	}
}

func TestEmbedFSRubricLoaderInChain(t *testing.T) {
	// Test that EmbedFSLoader works in a chain with fallback
	chain := NewChainLoader(
		NewEmbedFSLoader(testRubrics, "testdata"),
		EmbeddedLoader(),
	)

	// Custom should come from EmbedFSLoader
	rs, err := chain.Load(types.SpecType("custom"))
	if err != nil {
		t.Fatalf("Load(custom) failed: %v", err)
	}
	if rs.Name != "Custom Test Rubric" {
		t.Error("Did not load from EmbedFSLoader")
	}

	// PRD should fall back to EmbeddedLoader
	rs, err = chain.Load(types.SpecTypePRD)
	if err != nil {
		t.Fatalf("Load(prd) failed: %v", err)
	}
	if rs.Name == "" {
		t.Error("Failed to fall back to EmbeddedLoader")
	}
}
