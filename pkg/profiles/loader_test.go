package profiles

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ProductBuildersHQ/visionspec/pkg/types"
)

func TestFileLoaderLoad(t *testing.T) {
	// Create temp profile directory
	tmpDir := t.TempDir()
	profileDir := filepath.Join(tmpDir, "test-profile")
	if err := os.MkdirAll(profileDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create profile.yaml
	profileYAML := `name: test-profile
description: Test profile for unit tests
spec_config:
  prd:
    required: true
    category: source
`
	if err := os.WriteFile(filepath.Join(profileDir, "profile.yaml"), []byte(profileYAML), 0600); err != nil {
		t.Fatal(err)
	}

	loader := NewFileLoader(tmpDir)

	t.Run("loads valid profile", func(t *testing.T) {
		profile, err := loader.Load("test-profile")
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}
		if profile.Name != "test-profile" {
			t.Errorf("Name = %v, want test-profile", profile.Name)
		}
		if profile.Description != "Test profile for unit tests" {
			t.Errorf("Description = %v, want Test profile for unit tests", profile.Description)
		}
	})

	t.Run("returns error for missing profile", func(t *testing.T) {
		_, err := loader.Load("nonexistent")
		if err == nil {
			t.Error("Load() should return error for missing profile")
		}
	})
}

func TestFileLoaderAvailable(t *testing.T) {
	// Create temp directory with multiple profiles
	tmpDir := t.TempDir()

	// Create two valid profiles
	for _, name := range []string{"profile-a", "profile-b"} {
		profileDir := filepath.Join(tmpDir, name)
		if err := os.MkdirAll(profileDir, 0755); err != nil {
			t.Fatal(err)
		}
		profileYAML := "name: " + name + "\ndescription: Test profile\n"
		if err := os.WriteFile(filepath.Join(profileDir, "profile.yaml"), []byte(profileYAML), 0600); err != nil {
			t.Fatal(err)
		}
	}

	// Create a directory without profile.yaml
	invalidDir := filepath.Join(tmpDir, "not-a-profile")
	if err := os.MkdirAll(invalidDir, 0755); err != nil {
		t.Fatal(err)
	}

	loader := NewFileLoader(tmpDir)
	available := loader.Available()

	if len(available) != 2 {
		t.Errorf("Available() returned %d profiles, want 2", len(available))
	}

	found := make(map[string]bool)
	for _, name := range available {
		found[name] = true
	}

	if !found["profile-a"] || !found["profile-b"] {
		t.Error("Available() missing expected profiles")
	}

	if found["not-a-profile"] {
		t.Error("Available() should not include directories without profile.yaml")
	}
}

func TestChainLoaderLoad(t *testing.T) {
	// Create two temp directories with different profiles
	tmpDir1 := t.TempDir()
	tmpDir2 := t.TempDir()

	// Profile A in first loader
	profileDirA := filepath.Join(tmpDir1, "profile-a")
	if err := os.MkdirAll(profileDirA, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(profileDirA, "profile.yaml"), []byte("name: profile-a\ndescription: From loader 1\n"), 0600); err != nil {
		t.Fatal(err)
	}

	// Profile B in second loader
	profileDirB := filepath.Join(tmpDir2, "profile-b")
	if err := os.MkdirAll(profileDirB, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(profileDirB, "profile.yaml"), []byte("name: profile-b\ndescription: From loader 2\n"), 0600); err != nil {
		t.Fatal(err)
	}

	loader := NewChainLoader(NewFileLoader(tmpDir1), NewFileLoader(tmpDir2))

	t.Run("loads from first loader", func(t *testing.T) {
		profile, err := loader.Load("profile-a")
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}
		if profile.Name != "profile-a" {
			t.Errorf("Name = %v, want profile-a", profile.Name)
		}
	})

	t.Run("loads from second loader", func(t *testing.T) {
		profile, err := loader.Load("profile-b")
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}
		if profile.Name != "profile-b" {
			t.Errorf("Name = %v, want profile-b", profile.Name)
		}
	})

	t.Run("returns error for missing profile", func(t *testing.T) {
		_, err := loader.Load("nonexistent")
		if err == nil {
			t.Error("Load() should return error for missing profile")
		}
	})
}

func TestChainLoaderAvailable(t *testing.T) {
	// Create two temp directories with different profiles
	tmpDir1 := t.TempDir()
	tmpDir2 := t.TempDir()

	// Profile A in first loader
	profileDirA := filepath.Join(tmpDir1, "profile-a")
	if err := os.MkdirAll(profileDirA, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(profileDirA, "profile.yaml"), []byte("name: profile-a\n"), 0600); err != nil {
		t.Fatal(err)
	}

	// Profile B in second loader
	profileDirB := filepath.Join(tmpDir2, "profile-b")
	if err := os.MkdirAll(profileDirB, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(profileDirB, "profile.yaml"), []byte("name: profile-b\n"), 0600); err != nil {
		t.Fatal(err)
	}

	// Profile A also in second loader (should be deduplicated)
	profileDirA2 := filepath.Join(tmpDir2, "profile-a")
	if err := os.MkdirAll(profileDirA2, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(profileDirA2, "profile.yaml"), []byte("name: profile-a\n"), 0600); err != nil {
		t.Fatal(err)
	}

	loader := NewChainLoader(NewFileLoader(tmpDir1), NewFileLoader(tmpDir2))
	available := loader.Available()

	// Should have 2 unique profiles (profile-a should not be duplicated)
	if len(available) != 2 {
		t.Errorf("Available() returned %d profiles, want 2", len(available))
	}
}

func TestResolvingLoaderLoad(t *testing.T) {
	tmpDir := t.TempDir()

	// Create parent profile
	parentDir := filepath.Join(tmpDir, "parent")
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		t.Fatal(err)
	}
	parentYAML := `name: parent
description: Parent profile
spec_config:
  mrd:
    required: true
    category: source
`
	if err := os.WriteFile(filepath.Join(parentDir, "profile.yaml"), []byte(parentYAML), 0600); err != nil {
		t.Fatal(err)
	}

	// Create child profile that extends parent
	childDir := filepath.Join(tmpDir, "child")
	if err := os.MkdirAll(childDir, 0755); err != nil {
		t.Fatal(err)
	}
	childYAML := `name: child
description: Child profile
extends: parent
spec_config:
  prd:
    required: true
    category: source
`
	if err := os.WriteFile(filepath.Join(childDir, "profile.yaml"), []byte(childYAML), 0600); err != nil {
		t.Fatal(err)
	}

	loader := NewResolvingLoader(NewFileLoader(tmpDir))

	t.Run("resolves inheritance", func(t *testing.T) {
		profile, err := loader.Load("child")
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}
		if profile.Name != "child" {
			t.Errorf("Name = %v, want child", profile.Name)
		}

		// Should have both parent's mrd and child's prd
		if profile.SpecConfig.Specs["mrd"] == nil {
			t.Error("Profile should inherit mrd from parent")
		}
		if profile.SpecConfig.Specs["prd"] == nil {
			t.Error("Profile should have prd from child")
		}
	})
}

func TestResolvingLoaderCircularDetection(t *testing.T) {
	tmpDir := t.TempDir()

	// Create profile A that extends B
	profileDirA := filepath.Join(tmpDir, "profile-a")
	if err := os.MkdirAll(profileDirA, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(profileDirA, "profile.yaml"), []byte("name: profile-a\nextends: profile-b\n"), 0600); err != nil {
		t.Fatal(err)
	}

	// Create profile B that extends A (circular)
	profileDirB := filepath.Join(tmpDir, "profile-b")
	if err := os.MkdirAll(profileDirB, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(profileDirB, "profile.yaml"), []byte("name: profile-b\nextends: profile-a\n"), 0600); err != nil {
		t.Fatal(err)
	}

	loader := NewResolvingLoader(NewFileLoader(tmpDir))

	_, err := loader.Load("profile-a")
	if err == nil {
		t.Error("Load() should return error for circular inheritance")
	}
}

func TestDefaultLoader(t *testing.T) {
	loader := DefaultLoader()
	if loader == nil {
		t.Fatal("DefaultLoader() returned nil")
	}

	available := loader.Available()
	if len(available) == 0 {
		t.Error("DefaultLoader() should have at least one profile")
	}

	// Verify we can load some expected profiles
	expectedProfiles := []string{"startup", "growth", "enterprise", "0-1"}
	for _, name := range expectedProfiles {
		profile, err := loader.Load(name)
		if err != nil {
			t.Errorf("Load(%q) error = %v", name, err)
			continue
		}
		if profile.Name != name {
			t.Errorf("Load(%q).Name = %v, want %v", name, profile.Name, name)
		}
	}
}

func TestIsDefaultProfile(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"0-1", true},
		{"startup", true},
		{"growth", true},
		{"enterprise", true},
		{"custom-profile", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if IsDefaultProfile(tt.name) != tt.expected {
				t.Errorf("IsDefaultProfile(%q) = %v, want %v", tt.name, !tt.expected, tt.expected)
			}
		})
	}
}

func TestFileLoaderWithTemplatesAndRubrics(t *testing.T) {
	tmpDir := t.TempDir()
	profileDir := filepath.Join(tmpDir, "full-profile")

	// Create directories
	if err := os.MkdirAll(filepath.Join(profileDir, "templates"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(profileDir, "rubrics"), 0755); err != nil {
		t.Fatal(err)
	}

	// Create profile.yaml
	profileYAML := `name: full-profile
description: Profile with templates and rubrics
spec_config:
  prd:
    required: true
    category: source
`
	if err := os.WriteFile(filepath.Join(profileDir, "profile.yaml"), []byte(profileYAML), 0600); err != nil {
		t.Fatal(err)
	}

	// Create a template
	templateContent := "# PRD Template\n\nThis is a test template.\n"
	if err := os.WriteFile(filepath.Join(profileDir, "templates", "prd.md"), []byte(templateContent), 0600); err != nil {
		t.Fatal(err)
	}

	// Create a rubric
	rubricContent := `spec_type: prd
name: PRD Rubric
categories:
  - id: completeness
    name: Completeness
    weight: 0.5
`
	if err := os.WriteFile(filepath.Join(profileDir, "rubrics", "prd.rubric.yaml"), []byte(rubricContent), 0600); err != nil {
		t.Fatal(err)
	}

	loader := NewFileLoader(tmpDir)
	profile, err := loader.Load("full-profile")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	t.Run("has template loader", func(t *testing.T) {
		if profile.TemplateLoader == nil {
			t.Error("Profile should have TemplateLoader")
		}
	})

	t.Run("has rubric loader", func(t *testing.T) {
		if profile.RubricLoader == nil {
			t.Error("Profile should have RubricLoader")
		}
	})

	t.Run("can load template", func(t *testing.T) {
		if profile.TemplateLoader == nil {
			t.Skip("No template loader")
		}
		tmpl, err := profile.TemplateLoader.Load(types.SpecTypePRD)
		if err != nil {
			t.Errorf("TemplateLoader.Load() error = %v", err)
		}
		if tmpl == nil {
			t.Error("Template should not be nil")
		}
	})

	t.Run("can load rubric", func(t *testing.T) {
		if profile.RubricLoader == nil {
			t.Skip("No rubric loader")
		}
		rubric, err := profile.RubricLoader.Load(types.SpecTypePRD)
		if err != nil {
			t.Errorf("RubricLoader.Load() error = %v", err)
		}
		if rubric == nil {
			t.Error("Rubric should not be nil")
		}
	})
}
