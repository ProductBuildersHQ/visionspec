package apptypes

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Loader loads app type specifications from various sources.
type Loader interface {
	// Load returns the app type spec for a given app type.
	Load(appType AppType) (*AppTypeSpec, error)

	// Available returns all available app types.
	Available() []AppType
}

// builtinLoader provides built-in app type specs.
type builtinLoader struct{}

// BuiltinLoader returns a loader with built-in app type specs.
func BuiltinLoader() Loader {
	return &builtinLoader{}
}

func (l *builtinLoader) Load(appType AppType) (*AppTypeSpec, error) {
	switch appType {
	case AppTypeMicroservice:
		return MicroserviceSpec(), nil
	// TODO: Add other built-in specs
	// case AppTypeWebsite:
	//     return WebsiteSpec(), nil
	default:
		return nil, fmt.Errorf("no built-in spec for app type %q", appType)
	}
}

func (l *builtinLoader) Available() []AppType {
	return []AppType{
		AppTypeMicroservice,
		// TODO: Add as implemented
	}
}

// embeddedLoader loads app type specs from embedded files.
type embeddedLoader struct {
	fs  embed.FS
	dir string
}

// NewEmbeddedLoader creates a loader from embedded filesystem.
// App type specs should be named: {dir}/{apptype}.yaml (e.g., microservice.yaml)
//
// Usage:
//
//	//go:embed apptypes/*.yaml
//	var orgAppTypes embed.FS
//
//	loader := apptypes.NewEmbeddedLoader(orgAppTypes, "apptypes")
func NewEmbeddedLoader(fsys embed.FS, dir string) Loader {
	return &embeddedLoader{fs: fsys, dir: dir}
}

func (l *embeddedLoader) Load(appType AppType) (*AppTypeSpec, error) {
	path := fmt.Sprintf("%s/%s.yaml", l.dir, appType)
	content, err := l.fs.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("app type spec not found: %s: %w", appType, err)
	}

	var spec AppTypeSpec
	if err := yaml.Unmarshal(content, &spec); err != nil {
		return nil, fmt.Errorf("parsing app type spec %s: %w", appType, err)
	}

	return &spec, nil
}

func (l *embeddedLoader) Available() []AppType {
	var result []AppType

	entries, err := fs.ReadDir(l.fs, l.dir)
	if err != nil {
		return result
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}

		appType := AppType(strings.TrimSuffix(entry.Name(), ".yaml"))
		if appType.IsValid() {
			result = append(result, appType)
		}
	}

	return result
}

// fileLoader loads app type specs from a directory.
type fileLoader struct {
	dir string
}

// NewFileLoader creates a loader that reads app type specs from a directory.
func NewFileLoader(dir string) Loader {
	return &fileLoader{dir: dir}
}

func (l *fileLoader) Load(appType AppType) (*AppTypeSpec, error) {
	path := filepath.Join(l.dir, string(appType)+".yaml")
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("app type spec not found: %s", appType)
		}
		return nil, fmt.Errorf("reading app type spec %s: %w", appType, err)
	}

	var spec AppTypeSpec
	if err := yaml.Unmarshal(content, &spec); err != nil {
		return nil, fmt.Errorf("parsing app type spec %s: %w", appType, err)
	}

	return &spec, nil
}

func (l *fileLoader) Available() []AppType {
	var result []AppType

	entries, err := os.ReadDir(l.dir)
	if err != nil {
		return result
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}

		appType := AppType(strings.TrimSuffix(entry.Name(), ".yaml"))
		if appType.IsValid() {
			result = append(result, appType)
		}
	}

	return result
}

// chainLoader tries multiple loaders in order.
type chainLoader struct {
	loaders []Loader
}

// NewChainLoader creates a loader that tries multiple loaders in order.
// The first loader that can load an app type spec wins.
// This is useful for organization-specific overrides:
//
//	loader := apptypes.NewChainLoader(
//		orgLoader,        // Try org-specific first (more prescriptive)
//		builtinLoader,    // Fall back to visionspec defaults (more flexible)
//	)
func NewChainLoader(loaders ...Loader) Loader {
	return &chainLoader{loaders: loaders}
}

func (l *chainLoader) Load(appType AppType) (*AppTypeSpec, error) {
	var lastErr error

	for _, loader := range l.loaders {
		spec, err := loader.Load(appType)
		if err == nil {
			return spec, nil
		}
		lastErr = err
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, fmt.Errorf("no loader could load app type spec %q", appType)
}

func (l *chainLoader) Available() []AppType {
	seen := make(map[AppType]bool)
	var result []AppType

	for _, loader := range l.loaders {
		for _, appType := range loader.Available() {
			if !seen[appType] {
				seen[appType] = true
				result = append(result, appType)
			}
		}
	}

	return result
}

// memoryLoader holds app type specs in memory.
type memoryLoader struct {
	specs map[AppType]*AppTypeSpec
}

// NewMemoryLoader creates a loader from in-memory app type specs.
func NewMemoryLoader(specs ...*AppTypeSpec) Loader {
	l := &memoryLoader{
		specs: make(map[AppType]*AppTypeSpec),
	}
	for _, spec := range specs {
		l.specs[spec.Metadata.Name] = spec
	}
	return l
}

func (l *memoryLoader) Load(appType AppType) (*AppTypeSpec, error) {
	spec, ok := l.specs[appType]
	if !ok {
		return nil, fmt.Errorf("app type spec not found: %s", appType)
	}
	return spec, nil
}

func (l *memoryLoader) Available() []AppType {
	result := make([]AppType, 0, len(l.specs))
	for appType := range l.specs {
		result = append(result, appType)
	}
	return result
}

// DefaultLoader returns the default app type loader (built-in specs).
func DefaultLoader() Loader {
	return BuiltinLoader()
}
