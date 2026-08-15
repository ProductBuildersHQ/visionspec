package rubrics

import (
	"embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/ProductBuildersHQ/visionspec/pkg/types"
	"github.com/plexusone/structured-evaluation/rubric"
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

	if rs.ID != "prd-rubric" {
		t.Errorf("ID = %v, want prd-rubric", rs.ID)
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
id: security-rubric
name: "Security Rubric"
description: "Custom security rubric"
version: "1.0"
passCriteria:
  minCategoriesPassing: all_required
  maxFindingsSeverity: {critical: 0, high: 0, medium: 3, low: -1}
categories:
  - id: threat-modeling
    name: "Threat Modeling"
    description: "Are threats identified?"
    weight: 2.0
    required: true
    scale:
      type: categorical
      options:
        - {value: pass, label: Pass, criteria: ["Comprehensive threat model"]}
        - {value: partial, label: Partial, criteria: ["Basic threats identified"]}
        - {value: fail, label: Fail, criteria: ["No threat analysis"]}
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
id: prd-rubric
name: "Custom PRD Rubric"
description: "Overrides the default"
version: "2.0"
passCriteria:
  minCategoriesPassing: all_required
  maxFindingsSeverity: {critical: 0, high: 0, medium: 5, low: -1}
categories:
  - id: custom-category
    name: "Custom Category"
    description: "A custom category"
    weight: 1.0
    required: true
    scale:
      type: categorical
      options:
        - {value: pass, label: Pass, criteria: ["Custom pass"]}
        - {value: partial, label: Partial, criteria: ["Custom partial"]}
        - {value: fail, label: Fail, criteria: ["Custom fail"]}
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

func TestParseRubricYAMLValidation(t *testing.T) {
	tests := []struct {
		name        string
		yaml        string
		expectError bool
	}{
		{
			name: "valid rubric",
			yaml: `
id: test-rubric
name: Test Rubric
categories:
  - id: cat1
    name: Category 1
    scale:
      type: categorical
      options:
        - {value: pass, criteria: ["ok"]}
`,
			expectError: false,
		},
		{
			name: "no categories",
			yaml: `
id: test-rubric
name: Test Rubric
`,
			expectError: true,
		},
		{
			name:        "not yaml",
			yaml:        "\t: : not: valid",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseRubricYAML([]byte(tt.yaml), tt.name)
			if tt.expectError && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestWriteRubricYAMLRoundTrip(t *testing.T) {
	rs := rubric.NewRubricSet("prd-rubric", "Test Rubric", "1.0")
	rs.Description = "A test rubric"
	rs.PassCriteria = StrictPassCriteria()
	rs.AddCategory(*rubric.NewCategory("cat1", "Category 1", "First category").
		SetWeight(1.5).SetRequired(true).
		WithPassPartialFail([]string{"Pass criteria"}, []string{"Partial criteria"}, []string{"Fail criteria"}))

	dir := t.TempDir()
	path := filepath.Join(dir, "prd.rubric.yaml")
	if err := WriteRubricYAML(path, rs); err != nil {
		t.Fatalf("WriteRubricYAML: %v", err)
	}

	got, err := NewFileLoader(dir).Load(types.SpecTypePRD)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Name != "Test Rubric" {
		t.Errorf("Name = %q, want Test Rubric", got.Name)
	}
	if len(got.Categories) != 1 || got.Categories[0].ID != "cat1" {
		t.Fatalf("categories = %+v, want one cat1", got.Categories)
	}
	if len(passOptionCriteria(&got.Categories[0])) == 0 {
		t.Error("round-tripped category lost pass criteria")
	}
}

func TestEmbedFSRubricLoader(t *testing.T) {
	loader := NewEmbedFSLoader(testRubrics, "testdata")

	// Test loading a rubric from embedded FS
	rs, err := loader.Load(types.SpecType("custom"))
	if err != nil {
		t.Fatalf("Load(custom) failed: %v", err)
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
