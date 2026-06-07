package profiles

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ProductBuildersHQ/visionspec/pkg/types"
)

func TestProfileValidate(t *testing.T) {
	tests := []struct {
		name    string
		profile *Profile
		wantErr bool
	}{
		{
			name:    "valid profile",
			profile: &Profile{Name: "test"},
			wantErr: false,
		},
		{
			name:    "empty name",
			profile: &Profile{Name: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.profile.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Profile.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProfileGetSpecConfig(t *testing.T) {
	t.Run("nil spec config returns default", func(t *testing.T) {
		p := &Profile{Name: "test"}
		config := p.GetSpecConfig()
		if config == nil {
			t.Error("GetSpecConfig() returned nil, want default config")
		}
	})

	t.Run("returns configured spec config", func(t *testing.T) {
		specConfig := &types.SpecConfig{
			Specs: map[string]*types.SpecRequirement{
				"custom": {Required: true, Category: types.CategorySource},
			},
		}
		p := &Profile{Name: "test", SpecConfig: specConfig}
		config := p.GetSpecConfig()
		if config != specConfig {
			t.Error("GetSpecConfig() did not return configured spec config")
		}
	})
}

func TestProfileMerge(t *testing.T) {
	parent := &Profile{
		Name:        "parent",
		Description: "Parent profile",
		SpecConfig: &types.SpecConfig{
			Specs: map[string]*types.SpecRequirement{
				"prd": {Required: true, Category: types.CategorySource},
				"mrd": {Required: true, Category: types.CategorySource},
			},
		},
	}

	child := &Profile{
		Name:        "child",
		Description: "Child profile",
		SpecConfig: &types.SpecConfig{
			Specs: map[string]*types.SpecRequirement{
				"prd": {Required: false, Category: types.CategorySource}, // Override
				"trd": {Required: true, Category: types.CategoryTechnical},
			},
		},
	}

	merged := child.Merge(parent)

	t.Run("merged has child name", func(t *testing.T) {
		if merged.Name != "child" {
			t.Errorf("Name = %v, want child", merged.Name)
		}
	})

	t.Run("merged inherits parent specs", func(t *testing.T) {
		if merged.SpecConfig.Specs["mrd"] == nil {
			t.Error("Merged profile missing mrd from parent")
		}
	})

	t.Run("child overrides parent spec", func(t *testing.T) {
		if merged.SpecConfig.Specs["prd"].Required != false {
			t.Error("Child did not override parent prd.Required")
		}
	})

	t.Run("merged has child new specs", func(t *testing.T) {
		if merged.SpecConfig.Specs["trd"] == nil {
			t.Error("Merged profile missing trd from child")
		}
	})
}

func TestProfileYAMLToProfile(t *testing.T) {
	py := &ProfileYAML{
		Name:        "test",
		Description: "Test profile",
		Extends:     "enterprise",
		SpecConfig: map[string]*types.SpecRequirement{
			"prd": {Required: true, Category: types.CategorySource},
		},
	}

	p := py.ToProfile()

	if p.Name != py.Name {
		t.Errorf("Name = %v, want %v", p.Name, py.Name)
	}
	if p.Description != py.Description {
		t.Errorf("Description = %v, want %v", p.Description, py.Description)
	}
	if p.Extends != py.Extends {
		t.Errorf("Extends = %v, want %v", p.Extends, py.Extends)
	}
	if p.SpecConfig == nil {
		t.Error("SpecConfig is nil")
	}
	if p.SpecConfig.Specs["prd"] == nil {
		t.Error("SpecConfig missing prd")
	}
}

func TestProfileToYAML(t *testing.T) {
	p := &Profile{
		Name:        "test",
		Description: "Test profile",
		Extends:     "enterprise",
		SpecConfig: &types.SpecConfig{
			Specs: map[string]*types.SpecRequirement{
				"prd": {Required: true, Category: types.CategorySource},
			},
		},
	}

	py := ProfileToYAML(p)

	if py.Name != p.Name {
		t.Errorf("Name = %v, want %v", py.Name, p.Name)
	}
	if py.SpecConfig["prd"] == nil {
		t.Error("SpecConfig missing prd")
	}
}

func TestWriteProfileYAML(t *testing.T) {
	py := &ProfileYAML{
		Name:        "test",
		Description: "Test profile",
	}

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "profile.yaml")

	err := WriteProfileYAML(path, py)
	if err != nil {
		t.Fatalf("WriteProfileYAML() error = %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("WriteProfileYAML() did not create file")
	}
}

func TestProfileSummary(t *testing.T) {
	t.Run("returns description if set", func(t *testing.T) {
		p := &Profile{Name: "test", Description: "Test description"}
		if p.Summary() != "Test description" {
			t.Errorf("Summary() = %v, want Test description", p.Summary())
		}
	})

	t.Run("returns name if no description", func(t *testing.T) {
		p := &Profile{Name: "test"}
		if p.Summary() != "test" {
			t.Errorf("Summary() = %v, want test", p.Summary())
		}
	})
}

func TestProfileRequiredSpecs(t *testing.T) {
	t.Run("nil spec config returns nil", func(t *testing.T) {
		p := &Profile{Name: "test"}
		if p.RequiredSpecs() != nil {
			t.Error("RequiredSpecs() should return nil for nil spec config")
		}
	})

	t.Run("returns required specs", func(t *testing.T) {
		p := &Profile{
			Name: "test",
			SpecConfig: &types.SpecConfig{
				Specs: map[string]*types.SpecRequirement{
					"prd": {Required: true, Category: types.CategorySource},
					"mrd": {Required: false, Category: types.CategorySource},
				},
			},
		}
		required := p.RequiredSpecs()
		found := false
		for _, s := range required {
			if s == "prd" {
				found = true
				break
			}
		}
		if !found {
			t.Error("RequiredSpecs() should include prd")
		}
	})
}
