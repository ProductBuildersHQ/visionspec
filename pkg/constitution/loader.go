package constitution

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Loader loads constitutions from various sources.
type Loader interface {
	// Load returns the constitution for a given name.
	// Names follow the pattern: "org/<name>", "team/<name>", "project/<name>"
	Load(name string) (*Constitution, error)

	// LoadByLevel returns all constitutions at a given level.
	LoadByLevel(level Level) ([]*Constitution, error)

	// Available returns all available constitution names.
	Available() []string
}

// embeddedLoader loads constitutions from embedded files.
type embeddedLoader struct {
	fs  embed.FS
	dir string
}

// NewEmbeddedLoader creates a loader from embedded filesystem.
// Constitutions should be organized as: {dir}/org/*.yaml, {dir}/team/*.yaml, {dir}/project/*.yaml
//
// Usage:
//
//	//go:embed constitutions/**/*.yaml
//	var orgConstitutions embed.FS
//
//	loader := constitution.NewEmbeddedLoader(orgConstitutions, "constitutions")
func NewEmbeddedLoader(fsys embed.FS, dir string) Loader {
	return &embeddedLoader{fs: fsys, dir: dir}
}

func (l *embeddedLoader) Load(name string) (*Constitution, error) {
	// Parse name: "org/example" -> dir="org", file="example.yaml"
	parts := strings.SplitN(name, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid constitution name %q: expected format 'level/name'", name)
	}

	path := fmt.Sprintf("%s/%s/%s.yaml", l.dir, parts[0], parts[1])
	content, err := l.fs.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("constitution not found: %s: %w", name, err)
	}

	var c Constitution
	if err := yaml.Unmarshal(content, &c); err != nil {
		return nil, fmt.Errorf("parsing constitution %s: %w", name, err)
	}

	return &c, nil
}

func (l *embeddedLoader) LoadByLevel(level Level) ([]*Constitution, error) {
	var result []*Constitution

	levelDir := fmt.Sprintf("%s/%s", l.dir, level)
	entries, err := fs.ReadDir(l.fs, levelDir)
	if err != nil {
		// Level directory doesn't exist - not an error
		return result, nil
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}

		name := fmt.Sprintf("%s/%s", level, strings.TrimSuffix(entry.Name(), ".yaml"))
		c, err := l.Load(name)
		if err != nil {
			continue // Skip invalid files
		}
		result = append(result, c)
	}

	return result, nil
}

func (l *embeddedLoader) Available() []string {
	var result []string

	for _, level := range []Level{LevelOrganization, LevelTeam, LevelProject} {
		levelDir := fmt.Sprintf("%s/%s", l.dir, level)
		entries, err := fs.ReadDir(l.fs, levelDir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
				continue
			}
			name := fmt.Sprintf("%s/%s", level, strings.TrimSuffix(entry.Name(), ".yaml"))
			result = append(result, name)
		}
	}

	return result
}

// fileLoader loads constitutions from a directory.
type fileLoader struct {
	dir string
}

// NewFileLoader creates a loader that reads constitutions from a directory.
// Constitutions should be organized as: {dir}/org/*.yaml, {dir}/team/*.yaml, {dir}/project/*.yaml
func NewFileLoader(dir string) Loader {
	return &fileLoader{dir: dir}
}

func (l *fileLoader) Load(name string) (*Constitution, error) {
	parts := strings.SplitN(name, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid constitution name %q: expected format 'level/name'", name)
	}

	path := filepath.Join(l.dir, parts[0], parts[1]+".yaml")
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("constitution not found: %s", name)
		}
		return nil, fmt.Errorf("reading constitution %s: %w", name, err)
	}

	var c Constitution
	if err := yaml.Unmarshal(content, &c); err != nil {
		return nil, fmt.Errorf("parsing constitution %s: %w", name, err)
	}

	return &c, nil
}

func (l *fileLoader) LoadByLevel(level Level) ([]*Constitution, error) {
	var result []*Constitution

	levelDir := filepath.Join(l.dir, string(level))
	entries, err := os.ReadDir(levelDir)
	if err != nil {
		return result, nil
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}

		name := fmt.Sprintf("%s/%s", level, strings.TrimSuffix(entry.Name(), ".yaml"))
		c, err := l.Load(name)
		if err != nil {
			continue
		}
		result = append(result, c)
	}

	return result, nil
}

func (l *fileLoader) Available() []string {
	var result []string

	for _, level := range []Level{LevelOrganization, LevelTeam, LevelProject} {
		levelDir := filepath.Join(l.dir, string(level))
		entries, err := os.ReadDir(levelDir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
				continue
			}
			name := fmt.Sprintf("%s/%s", level, strings.TrimSuffix(entry.Name(), ".yaml"))
			result = append(result, name)
		}
	}

	return result
}

// chainLoader tries multiple loaders in order.
type chainLoader struct {
	loaders []Loader
}

// NewChainLoader creates a loader that tries multiple loaders in order.
// The first loader that can load a constitution wins.
// This is useful for organization-specific overrides:
//
//	loader := constitution.NewChainLoader(
//		orgLoader,        // Try org-specific first
//		defaultLoader,    // Fall back to visionspec defaults
//	)
func NewChainLoader(loaders ...Loader) Loader {
	return &chainLoader{loaders: loaders}
}

func (l *chainLoader) Load(name string) (*Constitution, error) {
	var lastErr error

	for _, loader := range l.loaders {
		c, err := loader.Load(name)
		if err == nil {
			return c, nil
		}
		lastErr = err
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, fmt.Errorf("no loader could load constitution %q", name)
}

func (l *chainLoader) LoadByLevel(level Level) ([]*Constitution, error) {
	seen := make(map[string]bool)
	var result []*Constitution

	for _, loader := range l.loaders {
		constitutions, err := loader.LoadByLevel(level)
		if err != nil {
			continue
		}
		for _, c := range constitutions {
			key := fmt.Sprintf("%s/%s", c.Metadata.Level, c.Metadata.Name)
			if !seen[key] {
				seen[key] = true
				result = append(result, c)
			}
		}
	}

	return result, nil
}

func (l *chainLoader) Available() []string {
	seen := make(map[string]bool)
	var result []string

	for _, loader := range l.loaders {
		for _, name := range loader.Available() {
			if !seen[name] {
				seen[name] = true
				result = append(result, name)
			}
		}
	}

	return result
}

// memoryLoader holds constitutions in memory.
// Useful for testing and programmatic configuration.
type memoryLoader struct {
	constitutions map[string]*Constitution
}

// NewMemoryLoader creates a loader from in-memory constitutions.
func NewMemoryLoader(constitutions ...*Constitution) Loader {
	l := &memoryLoader{
		constitutions: make(map[string]*Constitution),
	}
	for _, c := range constitutions {
		name := fmt.Sprintf("%s/%s", c.Metadata.Level, c.Metadata.Name)
		l.constitutions[name] = c
	}
	return l
}

func (l *memoryLoader) Load(name string) (*Constitution, error) {
	c, ok := l.constitutions[name]
	if !ok {
		return nil, fmt.Errorf("constitution not found: %s", name)
	}
	return c, nil
}

func (l *memoryLoader) LoadByLevel(level Level) ([]*Constitution, error) {
	var result []*Constitution
	prefix := string(level) + "/"

	for name, c := range l.constitutions {
		if strings.HasPrefix(name, prefix) {
			result = append(result, c)
		}
	}

	return result, nil
}

func (l *memoryLoader) Available() []string {
	result := make([]string, 0, len(l.constitutions))
	for name := range l.constitutions {
		result = append(result, name)
	}
	return result
}

// Resolver resolves constitution chains with inheritance.
type Resolver struct {
	loader Loader
}

// NewResolver creates a constitution resolver.
func NewResolver(loader Loader) *Resolver {
	return &Resolver{loader: loader}
}

// ResolveChain resolves a constitution by loading and merging its inheritance chain.
// Given "project/myproject", it loads:
//   - The project constitution
//   - Its inherited team constitution (if any)
//   - The team's inherited org constitution (if any)
//
// Then merges them in order: org → team → project
func (r *Resolver) ResolveChain(name string) (*Constitution, error) {
	c, err := r.loader.Load(name)
	if err != nil {
		return nil, err
	}

	// Build chain from leaf to root
	chain := []*Constitution{c}
	current := c

	for current.Metadata.Inherits != "" {
		parent, err := r.loader.Load(current.Metadata.Inherits)
		if err != nil {
			return nil, fmt.Errorf("loading inherited constitution %q: %w", current.Metadata.Inherits, err)
		}
		chain = append(chain, parent)
		current = parent
	}

	// Reverse to get root-to-leaf order
	for i, j := 0, len(chain)-1; i < j; i, j = i+1, j-1 {
		chain[i], chain[j] = chain[j], chain[i]
	}

	// Merge chain
	return Resolve(chain...)
}
