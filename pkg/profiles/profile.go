// Package profiles provides reusable configuration profiles for visionspec.
//
// Profiles bundle together spec requirements, templates, and rubrics into
// a cohesive configuration that can be selected at runtime or compiled
// into custom CLI tools.
//
// Default profiles are provided for common product stages:
//   - 0-1: Minimal configuration for idea validation
//   - startup: Lightweight configuration for pre-PMF startups
//   - growth: Metrics-driven configuration for 1-N scaling
//   - enterprise: Comprehensive configuration for post-PMF enterprises
//
// Organizations can create custom profiles that extend default profiles
// or define entirely new configurations.
package profiles

import (
	"fmt"
	"os"

	prrubrics "github.com/grokify/prism-roadmap/rubrics"
	prtemplates "github.com/grokify/prism-roadmap/templates"

	"github.com/ProductBuildersHQ/visionspec/pkg/rubrics"
	"github.com/ProductBuildersHQ/visionspec/pkg/templates"
	"github.com/ProductBuildersHQ/visionspec/pkg/types"
	"gopkg.in/yaml.v3"
)

// canvasTemplateLoader resolves canonical canvas templates (BMC, OpportunitySpec,
// and other canvases) from prism-roadmap, the single source of truth for them.
// VisionSpec no longer carries duplicate copies in profile directories.
func canvasTemplateLoader() templates.Loader {
	return templates.NewSubFSLoader(prtemplates.FS())
}

// canvasRubricLoader resolves canonical canvas rubrics from prism-roadmap.
func canvasRubricLoader() rubrics.Loader {
	return rubrics.NewSubFSLoader(prrubrics.FS())
}

// Profile represents a complete visionspec configuration profile.
type Profile struct {
	// Name is the profile identifier (e.g., "startup", "enterprise").
	Name string `yaml:"name" json:"name"`

	// Description explains the profile's purpose and use case.
	Description string `yaml:"description" json:"description"`

	// Extends is the name of a parent profile to inherit from.
	// Settings from this profile override the parent.
	Extends string `yaml:"extends,omitempty" json:"extends,omitempty"`

	// SpecConfig defines which specs are required and their settings.
	SpecConfig *types.SpecConfig `yaml:"spec_config,omitempty" json:"spec_config,omitempty"`

	// TemplateLoader provides templates for this profile.
	// Set by the loader, not from YAML.
	TemplateLoader templates.Loader `yaml:"-" json:"-"`

	// RubricLoader provides rubrics for this profile.
	// Set by the loader, not from YAML.
	RubricLoader rubrics.Loader `yaml:"-" json:"-"`

	// Path is the filesystem path where the profile was loaded from.
	// Empty for embedded profiles.
	Path string `yaml:"-" json:"-"`
}

// ProfileYAML is the YAML representation of a profile.
type ProfileYAML struct {
	Name        string                            `yaml:"name"`
	Description string                            `yaml:"description"`
	Extends     string                            `yaml:"extends,omitempty"`
	SpecConfig  map[string]*types.SpecRequirement `yaml:"spec_config,omitempty"`
}

// ToProfile converts ProfileYAML to a Profile.
func (py *ProfileYAML) ToProfile() *Profile {
	var specConfig *types.SpecConfig
	if py.SpecConfig != nil {
		specConfig = &types.SpecConfig{Specs: py.SpecConfig}
	}

	return &Profile{
		Name:        py.Name,
		Description: py.Description,
		Extends:     py.Extends,
		SpecConfig:  specConfig,
	}
}

// Validate checks if the profile is valid.
func (p *Profile) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("profile name is required")
	}
	return nil
}

// GetSpecConfig returns the effective SpecConfig, merging with parent if extends is set.
func (p *Profile) GetSpecConfig() *types.SpecConfig {
	if p.SpecConfig != nil {
		return p.SpecConfig
	}
	return types.DefaultSpecConfig()
}

// GetTemplateLoader returns the template loader, falling back to the default
// loader and then to prism-roadmap's canonical canvas templates.
func (p *Profile) GetTemplateLoader() templates.Loader {
	base := p.TemplateLoader
	if base == nil {
		base = templates.DefaultLoader()
	}
	return templates.NewChainLoader(base, canvasTemplateLoader())
}

// GetRubricLoader returns the rubric loader, falling back to the default loader
// and then to prism-roadmap's canonical canvas rubrics.
func (p *Profile) GetRubricLoader() rubrics.Loader {
	base := p.RubricLoader
	if base == nil {
		base = rubrics.DefaultLoader()
	}
	return rubrics.NewChainLoader(base, canvasRubricLoader())
}

// Merge combines this profile with a parent profile.
// Settings from this profile override the parent.
func (p *Profile) Merge(parent *Profile) *Profile {
	merged := &Profile{
		Name:        p.Name,
		Description: p.Description,
		Path:        p.Path,
	}

	// Merge spec config
	if parent.SpecConfig != nil {
		merged.SpecConfig = &types.SpecConfig{
			Specs: make(map[string]*types.SpecRequirement),
		}
		// Copy parent specs
		for k, v := range parent.SpecConfig.Specs {
			merged.SpecConfig.Specs[k] = v
		}
		// Override with this profile's specs
		if p.SpecConfig != nil {
			for k, v := range p.SpecConfig.Specs {
				merged.SpecConfig.Specs[k] = v
			}
		}
	} else if p.SpecConfig != nil {
		merged.SpecConfig = p.SpecConfig
	}

	// Chain template loaders (this profile first, then parent)
	if p.TemplateLoader != nil && parent.TemplateLoader != nil {
		merged.TemplateLoader = templates.NewChainLoader(p.TemplateLoader, parent.TemplateLoader)
	} else if p.TemplateLoader != nil {
		merged.TemplateLoader = p.TemplateLoader
	} else {
		merged.TemplateLoader = parent.TemplateLoader
	}

	// Chain rubric loaders (this profile first, then parent)
	if p.RubricLoader != nil && parent.RubricLoader != nil {
		merged.RubricLoader = rubrics.NewChainLoader(p.RubricLoader, parent.RubricLoader)
	} else if p.RubricLoader != nil {
		merged.RubricLoader = p.RubricLoader
	} else {
		merged.RubricLoader = parent.RubricLoader
	}

	return merged
}

// Summary returns a short description for display.
func (p *Profile) Summary() string {
	if p.Description != "" {
		return p.Description
	}
	return p.Name
}

// RequiredSpecs returns the list of required spec names.
func (p *Profile) RequiredSpecs() []string {
	if p.SpecConfig == nil {
		return nil
	}
	return p.SpecConfig.RequiredSpecs()
}

// ProfileToYAML converts a Profile back to ProfileYAML for serialization.
func ProfileToYAML(p *Profile) *ProfileYAML {
	py := &ProfileYAML{
		Name:        p.Name,
		Description: p.Description,
		Extends:     p.Extends,
	}

	if p.SpecConfig != nil && p.SpecConfig.Specs != nil {
		py.SpecConfig = p.SpecConfig.Specs
	}

	return py
}

// WriteProfileYAML writes a ProfileYAML to a file.
func WriteProfileYAML(path string, py *ProfileYAML) error {
	data, err := yaml.Marshal(py)
	if err != nil {
		return fmt.Errorf("marshaling profile: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}
