package rubrics

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ProductBuildersHQ/visionspec/pkg/types"
	"github.com/plexusone/structured-evaluation/rubric"
)

// Loader loads rubrics from various sources.
type Loader interface {
	// Load returns the rubric for a spec type.
	Load(specType types.SpecType) (*rubric.RubricSet, error)

	// Available returns all available spec types with rubrics.
	Available() []types.SpecType
}

// embeddedLoader loads rubrics from the Go registry.
type embeddedLoader struct{}

// EmbeddedLoader returns a loader that uses embedded Go-defined rubrics.
func EmbeddedLoader() Loader {
	return &embeddedLoader{}
}

func (l *embeddedLoader) Load(specType types.SpecType) (*rubric.RubricSet, error) {
	return Get(specType)
}

func (l *embeddedLoader) Available() []types.SpecType {
	return Available()
}

// fileLoader loads rubrics from YAML files in a directory.
type fileLoader struct {
	dir string
}

// NewFileLoader creates a loader that reads rubrics from YAML files.
// Rubrics are named: {spec-type}.rubric.yaml (e.g., prd.rubric.yaml)
func NewFileLoader(dir string) *fileLoader {
	return &fileLoader{dir: dir}
}

func (l *fileLoader) Load(specType types.SpecType) (*rubric.RubricSet, error) {
	filename := string(specType) + ".rubric.yaml"
	path := filepath.Join(l.dir, filename)

	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("rubric not found for spec type %q in %s", specType, l.dir)
		}
		return nil, fmt.Errorf("reading rubric %s: %w", path, err)
	}

	return parseRubricYAML(content, path)
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
		if !strings.HasSuffix(name, ".rubric.yaml") {
			continue
		}

		// Extract spec type from filename
		specType := types.SpecType(strings.TrimSuffix(name, ".rubric.yaml"))
		result = append(result, specType)
	}

	return result
}

// embedFSLoader loads rubrics from an embedded filesystem.
// This allows organizations to compile their rubrics into a single binary.
type embedFSLoader struct {
	fs  embed.FS
	dir string
}

// NewEmbedFSLoader creates a loader that reads rubrics from an embedded filesystem.
// This is useful for organizations that want to compile custom rubrics into their CLI.
//
// Usage:
//
//	//go:embed rubrics/*.rubric.yaml
//	var orgRubrics embed.FS
//
//	loader := rubrics.NewEmbedFSLoader(orgRubrics, "rubrics")
func NewEmbedFSLoader(fsys embed.FS, dir string) Loader {
	return &embedFSLoader{fs: fsys, dir: dir}
}

func (l *embedFSLoader) Load(specType types.SpecType) (*rubric.RubricSet, error) {
	filename := string(specType) + ".rubric.yaml"
	path := l.dir + "/" + filename // embed.FS uses forward slashes

	content, err := l.fs.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("rubric not found for spec type %q: %w", specType, err)
	}

	return parseRubricYAML(content, path)
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
		if !strings.HasSuffix(name, ".rubric.yaml") {
			continue
		}

		specType := types.SpecType(strings.TrimSuffix(name, ".rubric.yaml"))
		result = append(result, specType)
	}

	return result
}

// chainLoader tries multiple loaders in order.
type chainLoader struct {
	loaders []Loader
}

// NewChainLoader creates a loader that tries multiple loaders in order.
// The first loader that can load a rubric wins.
func NewChainLoader(loaders ...Loader) Loader {
	return &chainLoader{loaders: loaders}
}

func (l *chainLoader) Load(specType types.SpecType) (*rubric.RubricSet, error) {
	var lastErr error

	for _, loader := range l.loaders {
		rs, err := loader.Load(specType)
		if err == nil {
			return rs, nil
		}
		lastErr = err
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, fmt.Errorf("no loader could load rubric for spec type %q", specType)
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

// subFSLoader loads rubrics from an fs.FS (sub-filesystem).
type subFSLoader struct {
	fsys fs.FS
}

// NewSubFSLoader creates a loader from an fs.FS interface.
// This is useful when working with fs.Sub() results.
func NewSubFSLoader(fsys fs.FS) Loader {
	return &subFSLoader{fsys: fsys}
}

func (l *subFSLoader) Load(specType types.SpecType) (*rubric.RubricSet, error) {
	filename := string(specType) + ".rubric.yaml"

	content, err := fs.ReadFile(l.fsys, filename)
	if err != nil {
		return nil, fmt.Errorf("rubric not found for spec type %q: %w", specType, err)
	}

	return parseRubricYAML(content, filename)
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
		if !strings.HasSuffix(name, ".rubric.yaml") {
			continue
		}

		specType := types.SpecType(strings.TrimSuffix(name, ".rubric.yaml"))
		result = append(result, specType)
	}

	return result
}

// DefaultLoader returns the default rubric loader (embedded rubrics).
func DefaultLoader() Loader {
	return EmbeddedLoader()
}

// LoadWithLoader loads a rubric using a specific loader.
func LoadWithLoader(loader Loader, specType types.SpecType) (*rubric.RubricSet, error) {
	if loader == nil {
		loader = DefaultLoader()
	}
	return loader.Load(specType)
}
