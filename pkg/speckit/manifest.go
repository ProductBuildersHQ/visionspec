package speckit

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Manifest mirrors spec-kit's extension.yml (plugin schema 1.0).
type Manifest struct {
	SchemaVersion string           `yaml:"schema_version"`
	Extension     ManifestMeta     `yaml:"extension"`
	Requires      ManifestRequires `yaml:"requires"`
	Provides      ManifestProvides `yaml:"provides"`
	Tags          []string         `yaml:"tags,omitempty"`
}

// ManifestMeta is the extension identity block.
type ManifestMeta struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
	Author      string `yaml:"author"`
	Repository  string `yaml:"repository"`
	License     string `yaml:"license"`
	Category    string `yaml:"category,omitempty"`
	Effect      string `yaml:"effect,omitempty"`
}

// ManifestRequires declares compatibility constraints.
type ManifestRequires struct {
	SpecKitVersion string `yaml:"speckit_version"`
}

// ManifestProvides lists the components the extension ships.
type ManifestProvides struct {
	Commands  []ManifestCommand  `yaml:"commands,omitempty"`
	Templates []ManifestTemplate `yaml:"templates,omitempty"`
}

// ManifestCommand is one provided slash command.
type ManifestCommand struct {
	Name        string `yaml:"name"`
	File        string `yaml:"file"`
	Description string `yaml:"description"`
}

// ManifestTemplate is one provided document template asset.
type ManifestTemplate struct {
	Name string `yaml:"name"`
	File string `yaml:"file"`
}

// renderManifest serializes a Manifest with a provenance header, validating
// the spec-kit schema constraints that matter for install-time acceptance.
func renderManifest(m Manifest, family string) ([]byte, error) {
	if !extensionIDRe.MatchString(m.Extension.ID) {
		return nil, fmt.Errorf("extension id %q: must match %s", m.Extension.ID, extensionIDRe)
	}
	if len(m.Extension.Description) >= maxDescription {
		return nil, fmt.Errorf("extension %s: description is %d chars; spec-kit requires < %d",
			m.Extension.ID, len(m.Extension.Description), maxDescription)
	}
	for _, c := range m.Provides.Commands {
		if !commandNameRe.MatchString(c.Name) {
			return nil, fmt.Errorf("extension %s: command name %q: must match %s", m.Extension.ID, c.Name, commandNameRe)
		}
		if c.Description == "" {
			return nil, fmt.Errorf("extension %s: command %s: description is required", m.Extension.ID, c.Name)
		}
	}
	for _, t := range m.Provides.Templates {
		if !extensionIDRe.MatchString(t.Name) {
			return nil, fmt.Errorf("extension %s: template name %q: must match %s", m.Extension.ID, t.Name, extensionIDRe)
		}
	}

	data, err := yaml.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("marshaling extension.yml for %s: %w", m.Extension.ID, err)
	}
	return append([]byte(yamlHeader(family)), data...), nil
}
