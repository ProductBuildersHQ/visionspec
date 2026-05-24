package templates

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/plexusone/multispec/pkg/types"
)

// Loader loads templates from various sources.
type Loader interface {
	// Load returns the template for a spec type.
	Load(specType types.SpecType) (*Template, error)

	// Available returns all available spec types.
	Available() []types.SpecType
}

// embeddedLoader loads templates from embedded files.
type embeddedLoader struct{}

// EmbeddedLoader returns a loader that uses embedded templates.
func EmbeddedLoader() Loader {
	return &embeddedLoader{}
}

func (l *embeddedLoader) Load(specType types.SpecType) (*Template, error) {
	return Get(specType)
}

func (l *embeddedLoader) Available() []types.SpecType {
	return Available()
}

// fileLoader loads templates from a directory.
type fileLoader struct {
	dir         string
	customTypes map[string]types.SpecCategory
}

// NewFileLoader creates a loader that reads templates from a directory.
// Templates are named: {spec-type}.md (e.g., prd.md, security.md)
func NewFileLoader(dir string) *fileLoader {
	return &fileLoader{
		dir:         dir,
		customTypes: make(map[string]types.SpecCategory),
	}
}

// RegisterCustomType registers a custom spec type with its category.
// This is needed for spec types that aren't built-in.
func (l *fileLoader) RegisterCustomType(name string, category types.SpecCategory) {
	l.customTypes[name] = category
}

func (l *fileLoader) Load(specType types.SpecType) (*Template, error) {
	filename := string(specType) + ".md"
	path := filepath.Join(l.dir, filename)

	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("template not found for spec type %q in %s", specType, l.dir)
		}
		return nil, fmt.Errorf("reading template %s: %w", path, err)
	}

	return &Template{
		SpecType: specType,
		Content:  string(content),
	}, nil
}

func (l *fileLoader) Available() []types.SpecType {
	var result []types.SpecType

	entries, err := os.ReadDir(l.dir)
	if err != nil {
		return result
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".md") {
			continue
		}

		// Extract spec type from filename
		specType := types.SpecType(strings.TrimSuffix(name, ".md"))
		result = append(result, specType)
	}

	return result
}

// embedFSLoader loads templates from an embedded filesystem.
// This allows organizations to compile their templates into a single binary.
type embedFSLoader struct {
	fs  embed.FS
	dir string
}

// NewEmbedFSLoader creates a loader that reads templates from an embedded filesystem.
// This is useful for organizations that want to compile custom templates into their CLI.
//
// Usage:
//
//	//go:embed templates/*.md
//	var orgTemplates embed.FS
//
//	loader := templates.NewEmbedFSLoader(orgTemplates, "templates")
func NewEmbedFSLoader(fsys embed.FS, dir string) Loader {
	return &embedFSLoader{fs: fsys, dir: dir}
}

func (l *embedFSLoader) Load(specType types.SpecType) (*Template, error) {
	filename := string(specType) + ".md"
	path := l.dir + "/" + filename // embed.FS uses forward slashes

	content, err := l.fs.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("template not found for spec type %q: %w", specType, err)
	}

	return &Template{
		SpecType: specType,
		Content:  string(content),
	}, nil
}

func (l *embedFSLoader) Available() []types.SpecType {
	var result []types.SpecType

	entries, err := fs.ReadDir(l.fs, l.dir)
	if err != nil {
		return result
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".md") {
			continue
		}

		specType := types.SpecType(strings.TrimSuffix(name, ".md"))
		result = append(result, specType)
	}

	return result
}

// chainLoader tries multiple loaders in order.
type chainLoader struct {
	loaders []Loader
}

// NewChainLoader creates a loader that tries multiple loaders in order.
// The first loader that can load a template wins.
func NewChainLoader(loaders ...Loader) Loader {
	return &chainLoader{loaders: loaders}
}

func (l *chainLoader) Load(specType types.SpecType) (*Template, error) {
	var lastErr error

	for _, loader := range l.loaders {
		t, err := loader.Load(specType)
		if err == nil {
			return t, nil
		}
		lastErr = err
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, fmt.Errorf("no loader could load template for spec type %q", specType)
}

func (l *chainLoader) Available() []types.SpecType {
	seen := make(map[types.SpecType]bool)
	var result []types.SpecType

	for _, loader := range l.loaders {
		for _, specType := range loader.Available() {
			if !seen[specType] {
				seen[specType] = true
				result = append(result, specType)
			}
		}
	}

	return result
}

// subFSLoader loads templates from an fs.FS (sub-filesystem).
type subFSLoader struct {
	fsys fs.FS
}

// NewSubFSLoader creates a loader from an fs.FS interface.
// This is useful when working with fs.Sub() results.
func NewSubFSLoader(fsys fs.FS) Loader {
	return &subFSLoader{fsys: fsys}
}

func (l *subFSLoader) Load(specType types.SpecType) (*Template, error) {
	filename := string(specType) + ".md"

	content, err := fs.ReadFile(l.fsys, filename)
	if err != nil {
		return nil, fmt.Errorf("template not found for spec type %q: %w", specType, err)
	}

	return &Template{
		SpecType: specType,
		Content:  string(content),
	}, nil
}

func (l *subFSLoader) Available() []types.SpecType {
	var result []types.SpecType

	entries, err := fs.ReadDir(l.fsys, ".")
	if err != nil {
		return result
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".md") {
			continue
		}

		specType := types.SpecType(strings.TrimSuffix(name, ".md"))
		result = append(result, specType)
	}

	return result
}

// DefaultLoader returns the default template loader (embedded templates).
func DefaultLoader() Loader {
	return EmbeddedLoader()
}

// LoadWithLoader loads a template using a specific loader.
func LoadWithLoader(loader Loader, specType types.SpecType) (*Template, error) {
	if loader == nil {
		loader = DefaultLoader()
	}
	return loader.Load(specType)
}
