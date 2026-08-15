// Package pipelines provides the embedded default library of pipelines linking
// definition-side workflows to execution-side integrations (see pkg/pipeline).
//
// Pipelines load from embedded YAML with no filesystem access, mirroring the
// pkg/workflows and pkg/integrations patterns. This package holds data only.
package pipelines

import (
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"
	"sync"

	"github.com/ProductBuildersHQ/visionspec/pkg/pipeline"
)

//go:embed default/*.yaml
var defaultFS embed.FS

var (
	mu       sync.RWMutex
	registry = map[string]*pipeline.Pipeline{}
)

func init() {
	if err := loadAll(); err != nil {
		panic(fmt.Sprintf("failed to load embedded pipelines: %v", err))
	}
}

// loadAll discovers and loads every pipeline from the embedded filesystem.
func loadAll() error {
	entries, err := fs.ReadDir(defaultFS, "default")
	if err != nil {
		return fmt.Errorf("reading default pipelines directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}

		data, err := fs.ReadFile(defaultFS, "default/"+entry.Name())
		if err != nil {
			return fmt.Errorf("reading pipeline %s: %w", entry.Name(), err)
		}

		p, err := pipeline.ParseYAML(data)
		if err != nil {
			return fmt.Errorf("parsing pipeline %s: %w", entry.Name(), err)
		}
		if err := p.Validate(); err != nil {
			return fmt.Errorf("invalid pipeline %s: %w", entry.Name(), err)
		}

		Register(p)
	}

	return nil
}

// Register adds or replaces a pipeline in the registry.
func Register(p *pipeline.Pipeline) {
	mu.Lock()
	defer mu.Unlock()
	registry[p.ID] = p
}

// Get returns the pipeline with the given ID.
func Get(id string) (*pipeline.Pipeline, error) {
	mu.RLock()
	defer mu.RUnlock()
	p, ok := registry[id]
	if !ok {
		return nil, fmt.Errorf("pipeline %q not found", id)
	}
	return p, nil
}

// List returns all registered pipeline IDs, sorted.
func List() []string {
	mu.RLock()
	defer mu.RUnlock()
	ids := make([]string, 0, len(registry))
	for id := range registry {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

// All returns every registered pipeline, sorted by ID.
func All() []*pipeline.Pipeline {
	ids := List()
	out := make([]*pipeline.Pipeline, 0, len(ids))
	for _, id := range ids {
		// id came from List() over the same registry, which nothing deletes
		// from; a miss is an invariant violation, not a runtime condition.
		p, err := Get(id)
		if err != nil {
			panic(fmt.Sprintf("pipelines: registry inconsistency: %v", err))
		}
		out = append(out, p)
	}
	return out
}
