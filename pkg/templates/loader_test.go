package templates

import (
	"embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/ProductBuildersHQ/visionspec/pkg/types"
)

//go:embed testdata/*.md
var testTemplates embed.FS

func TestEmbeddedLoader(t *testing.T) {
	loader := EmbeddedLoader()

	// Test loading a known template
	tmpl, err := loader.Load(types.SpecTypePRD)
	if err != nil {
		t.Fatalf("Load(prd) failed: %v", err)
	}

	if tmpl.SpecType != types.SpecTypePRD {
		t.Errorf("SpecType = %v, want %v", tmpl.SpecType, types.SpecTypePRD)
	}

	if tmpl.Content == "" {
		t.Error("Content is empty")
	}

	// Test available templates
	available := loader.Available()
	if len(available) == 0 {
		t.Error("Available() returned empty list")
	}
}

func TestFileLoader(t *testing.T) {
	// Create temp directory with test templates
	tmpDir := t.TempDir()

	// Write a custom template
	customContent := "# Custom Security Spec\n\nThis is a security template."
	err := os.WriteFile(filepath.Join(tmpDir, "security.md"), []byte(customContent), 0600)
	if err != nil {
		t.Fatalf("Failed to write template: %v", err)
	}

	loader := NewFileLoader(tmpDir)

	// Test loading custom template
	tmpl, err := loader.Load(types.SpecType("security"))
	if err != nil {
		t.Fatalf("Load(security) failed: %v", err)
	}

	if tmpl.Content != customContent {
		t.Errorf("Content = %q, want %q", tmpl.Content, customContent)
	}

	// Test loading non-existent template
	_, err = loader.Load(types.SpecType("nonexistent"))
	if err == nil {
		t.Error("Expected error for non-existent template")
	}

	// Test available templates
	available := loader.Available()
	if len(available) != 1 {
		t.Errorf("Available() returned %d templates, want 1", len(available))
	}
}

func TestFileLoaderRegisterCustomType(t *testing.T) {
	tmpDir := t.TempDir()
	loader := NewFileLoader(tmpDir)

	loader.RegisterCustomType("compliance", types.CategorySource)

	// The custom type is registered but we need to verify it's tracked
	if _, ok := loader.customTypes["compliance"]; !ok {
		t.Error("Custom type not registered")
	}
}

func TestChainLoader(t *testing.T) {
	// Create temp directory with override template
	tmpDir := t.TempDir()

	// Write a custom PRD template that overrides the embedded one
	customContent := "# Custom PRD\n\nThis overrides the default."
	err := os.WriteFile(filepath.Join(tmpDir, "prd.md"), []byte(customContent), 0600)
	if err != nil {
		t.Fatalf("Failed to write template: %v", err)
	}

	// Chain: file loader first (override), then embedded (fallback)
	chain := NewChainLoader(
		NewFileLoader(tmpDir),
		EmbeddedLoader(),
	)

	// PRD should come from file loader (override)
	tmpl, err := chain.Load(types.SpecTypePRD)
	if err != nil {
		t.Fatalf("Load(prd) failed: %v", err)
	}

	if tmpl.Content != customContent {
		t.Error("Chain loader did not use file loader for override")
	}

	// MRD should come from embedded loader (fallback)
	tmpl, err = chain.Load(types.SpecTypeMRD)
	if err != nil {
		t.Fatalf("Load(mrd) failed: %v", err)
	}

	if tmpl.Content == "" {
		t.Error("Chain loader did not fall back to embedded loader")
	}

	// Available should include both
	available := chain.Available()
	if len(available) < 2 {
		t.Errorf("Available() returned %d templates, expected at least 2", len(available))
	}
}

func TestChainLoaderNoLoaders(t *testing.T) {
	chain := NewChainLoader()

	_, err := chain.Load(types.SpecTypePRD)
	if err == nil {
		t.Error("Expected error with no loaders")
	}
}

func TestDefaultLoader(t *testing.T) {
	loader := DefaultLoader()

	// Should be able to load embedded templates
	tmpl, err := loader.Load(types.SpecTypePRD)
	if err != nil {
		t.Fatalf("Load(prd) failed: %v", err)
	}

	if tmpl.Content == "" {
		t.Error("Content is empty")
	}
}

func TestLoadWithLoader(t *testing.T) {
	// Test with nil loader (should use default)
	tmpl, err := LoadWithLoader(nil, types.SpecTypePRD)
	if err != nil {
		t.Fatalf("LoadWithLoader(nil, prd) failed: %v", err)
	}

	if tmpl.Content == "" {
		t.Error("Content is empty")
	}

	// Test with explicit loader
	tmpl, err = LoadWithLoader(EmbeddedLoader(), types.SpecTypeMRD)
	if err != nil {
		t.Fatalf("LoadWithLoader(embedded, mrd) failed: %v", err)
	}

	if tmpl.Content == "" {
		t.Error("Content is empty")
	}
}

func TestEmbedFSLoader(t *testing.T) {
	loader := NewEmbedFSLoader(testTemplates, "testdata")

	// Test loading a template from embedded FS
	tmpl, err := loader.Load(types.SpecType("custom"))
	if err != nil {
		t.Fatalf("Load(custom) failed: %v", err)
	}

	if tmpl.SpecType != types.SpecType("custom") {
		t.Errorf("SpecType = %v, want custom", tmpl.SpecType)
	}

	if tmpl.Content == "" {
		t.Error("Content is empty")
	}

	if !contains(tmpl.Content, "Custom Template") {
		t.Errorf("Content does not contain expected text: %s", tmpl.Content)
	}

	// Test loading non-existent template
	_, err = loader.Load(types.SpecType("nonexistent"))
	if err == nil {
		t.Error("Expected error for non-existent template")
	}

	// Test available templates
	available := loader.Available()
	if len(available) != 1 {
		t.Errorf("Available() returned %d templates, want 1", len(available))
	}
	if available[0] != types.SpecType("custom") {
		t.Errorf("Available()[0] = %v, want custom", available[0])
	}
}

func TestEmbedFSLoaderInChain(t *testing.T) {
	// Test that EmbedFSLoader works in a chain with fallback
	chain := NewChainLoader(
		NewEmbedFSLoader(testTemplates, "testdata"),
		EmbeddedLoader(),
	)

	// Custom should come from EmbedFSLoader
	tmpl, err := chain.Load(types.SpecType("custom"))
	if err != nil {
		t.Fatalf("Load(custom) failed: %v", err)
	}
	if !contains(tmpl.Content, "Custom Template") {
		t.Error("Did not load from EmbedFSLoader")
	}

	// PRD should fall back to EmbeddedLoader
	tmpl, err = chain.Load(types.SpecTypePRD)
	if err != nil {
		t.Fatalf("Load(prd) failed: %v", err)
	}
	if tmpl.Content == "" {
		t.Error("Failed to fall back to EmbeddedLoader")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
