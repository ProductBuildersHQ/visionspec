package profiles

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/ProductBuildersHQ/visionspec/pkg/rubrics"
	"github.com/ProductBuildersHQ/visionspec/pkg/templates"
)

// Loader loads profiles by name.
type Loader interface {
	// Load returns a profile by name.
	Load(name string) (*Profile, error)

	// Available returns all available profile names.
	Available() []string
}

// embedFSLoader loads profiles from an embedded filesystem.
type embedFSLoader struct {
	fs  embed.FS
	dir string
}

// NewEmbedFSLoader creates a loader from an embedded filesystem.
// The dir parameter specifies the base directory within the embed.FS.
func NewEmbedFSLoader(fsys embed.FS, dir string) Loader {
	return &embedFSLoader{fs: fsys, dir: dir}
}

func (l *embedFSLoader) Load(name string) (*Profile, error) {
	// Use path.Join (not filepath.Join) for embed.FS which always uses forward slashes
	profileDir := path.Join(l.dir, name)

	// Load profile.yaml
	profilePath := path.Join(profileDir, "profile.yaml")
	data, err := l.fs.ReadFile(profilePath)
	if err != nil {
		return nil, fmt.Errorf("profile %q not found: %w", name, err)
	}

	var py ProfileYAML
	if err := yaml.Unmarshal(data, &py); err != nil {
		return nil, fmt.Errorf("parsing profile %q: %w", name, err)
	}

	profile := py.ToProfile()
	if err := profile.Validate(); err != nil {
		return nil, fmt.Errorf("invalid profile %q: %w", name, err)
	}

	// Set up template loader for this profile
	templatesDir := path.Join(profileDir, "templates")
	if _, err := l.fs.ReadDir(templatesDir); err == nil {
		// Create a sub-filesystem for templates
		subFS, err := fs.Sub(l.fs, templatesDir)
		if err == nil {
			profile.TemplateLoader = templates.NewSubFSLoader(subFS)
		}
	}

	// Set up rubric loader for this profile
	rubricsDir := path.Join(profileDir, "rubrics")
	if _, err := l.fs.ReadDir(rubricsDir); err == nil {
		// Create a sub-filesystem for rubrics
		subFS, err := fs.Sub(l.fs, rubricsDir)
		if err == nil {
			profile.RubricLoader = rubrics.NewSubFSLoader(subFS)
		}
	}

	return profile, nil
}

func (l *embedFSLoader) Available() []string {
	entries, err := l.fs.ReadDir(l.dir)
	if err != nil {
		return nil
	}

	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			// Check if it has a profile.yaml
			// Use path.Join (not filepath.Join) for embed.FS which always uses forward slashes
			profilePath := path.Join(l.dir, entry.Name(), "profile.yaml")
			if _, err := l.fs.ReadFile(profilePath); err == nil {
				names = append(names, entry.Name())
			}
		}
	}
	return names
}

// fileLoader loads profiles from the filesystem.
type fileLoader struct {
	dir string
}

// NewFileLoader creates a loader from a filesystem directory.
func NewFileLoader(dir string) Loader {
	return &fileLoader{dir: dir}
}

func (l *fileLoader) Load(name string) (*Profile, error) {
	profileDir := filepath.Join(l.dir, name)

	// Load profile.yaml
	profilePath := filepath.Join(profileDir, "profile.yaml")
	data, err := os.ReadFile(profilePath)
	if err != nil {
		return nil, fmt.Errorf("profile %q not found: %w", name, err)
	}

	var py ProfileYAML
	if err := yaml.Unmarshal(data, &py); err != nil {
		return nil, fmt.Errorf("parsing profile %q: %w", name, err)
	}

	profile := py.ToProfile()
	profile.Path = profileDir

	if err := profile.Validate(); err != nil {
		return nil, fmt.Errorf("invalid profile %q: %w", name, err)
	}

	// Set up template loader for this profile
	templatesDir := filepath.Join(profileDir, "templates")
	if info, err := os.Stat(templatesDir); err == nil && info.IsDir() {
		profile.TemplateLoader = templates.NewFileLoader(templatesDir)
	}

	// Set up rubric loader for this profile
	rubricsDir := filepath.Join(profileDir, "rubrics")
	if info, err := os.Stat(rubricsDir); err == nil && info.IsDir() {
		profile.RubricLoader = rubrics.NewFileLoader(rubricsDir)
	}

	return profile, nil
}

func (l *fileLoader) Available() []string {
	entries, err := os.ReadDir(l.dir)
	if err != nil {
		return nil
	}

	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			profilePath := filepath.Join(l.dir, entry.Name(), "profile.yaml")
			if _, err := os.Stat(profilePath); err == nil {
				names = append(names, entry.Name())
			}
		}
	}
	return names
}

// chainLoader tries multiple loaders in order.
type chainLoader struct {
	loaders []Loader
}

// NewChainLoader creates a loader that tries multiple loaders in order.
func NewChainLoader(loaders ...Loader) Loader {
	return &chainLoader{loaders: loaders}
}

func (l *chainLoader) Load(name string) (*Profile, error) {
	var lastErr error
	for _, loader := range l.loaders {
		profile, err := loader.Load(name)
		if err == nil {
			return profile, nil
		}
		lastErr = err
	}
	return nil, lastErr
}

func (l *chainLoader) Available() []string {
	seen := make(map[string]bool)
	var names []string

	for _, loader := range l.loaders {
		for _, name := range loader.Available() {
			if !seen[name] {
				seen[name] = true
				names = append(names, name)
			}
		}
	}
	return names
}

// resolvingLoader resolves profile inheritance.
type resolvingLoader struct {
	base Loader
}

// NewResolvingLoader creates a loader that resolves profile inheritance.
func NewResolvingLoader(base Loader) Loader {
	return &resolvingLoader{base: base}
}

func (l *resolvingLoader) Load(name string) (*Profile, error) {
	return l.loadWithChain(name, nil)
}

func (l *resolvingLoader) loadWithChain(name string, chain []string) (*Profile, error) {
	// Check for circular inheritance
	for _, n := range chain {
		if n == name {
			return nil, fmt.Errorf("circular inheritance detected: %s", strings.Join(append(chain, name), " -> "))
		}
	}

	profile, err := l.base.Load(name)
	if err != nil {
		return nil, err
	}

	// If this profile extends another, load and merge
	if profile.Extends != "" {
		parent, err := l.loadWithChain(profile.Extends, append(chain, name))
		if err != nil {
			return nil, fmt.Errorf("loading parent profile %q: %w", profile.Extends, err)
		}
		profile = profile.Merge(parent)
	}

	return profile, nil
}

func (l *resolvingLoader) Available() []string {
	return l.base.Available()
}
