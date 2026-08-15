// Package integrations provides the embedded default library of external
// spec-driven-development tool descriptors (see pkg/integration).
//
// Descriptors load from embedded YAML with no filesystem access, mirroring the
// pkg/workflows pattern. Consumers (visionspec, visionstudio) read these to
// detect and represent tool status; this package holds data only.
package integrations

import (
	"embed"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"sync"

	"github.com/ProductBuildersHQ/visionspec/pkg/integration"
)

//go:embed default/*
var defaultFS embed.FS

var (
	mu       sync.RWMutex
	registry = map[string]*integration.Integration{}
)

func init() {
	if err := loadAll(); err != nil {
		panic(fmt.Sprintf("failed to load embedded integrations: %v", err))
	}
}

// loadAll discovers and loads every integration from the embedded filesystem.
func loadAll() error {
	entries, err := fs.ReadDir(defaultFS, "default")
	if err != nil {
		return fmt.Errorf("reading default integrations directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		id := entry.Name()
		data, err := fs.ReadFile(defaultFS, path.Join("default", id, "integration.yaml"))
		if err != nil {
			return fmt.Errorf("reading integration %s: %w", id, err)
		}

		in, err := integration.ParseYAML(data)
		if err != nil {
			return fmt.Errorf("parsing integration %s: %w", id, err)
		}
		if err := in.Validate(); err != nil {
			return fmt.Errorf("invalid integration %s: %w", id, err)
		}
		if in.ID != id {
			return fmt.Errorf("integration %s: id %q does not match directory name", id, in.ID)
		}

		Register(in)
	}

	return nil
}

// Register adds or replaces an integration in the registry.
func Register(in *integration.Integration) {
	mu.Lock()
	defer mu.Unlock()
	registry[in.ID] = in
}

// Get returns the integration with the given ID.
func Get(id string) (*integration.Integration, error) {
	mu.RLock()
	defer mu.RUnlock()
	in, ok := registry[id]
	if !ok {
		return nil, fmt.Errorf("integration %q not found", id)
	}
	return in, nil
}

// List returns all registered integration IDs, sorted.
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

// All returns every registered integration, sorted by ID.
func All() []*integration.Integration {
	ids := List()
	out := make([]*integration.Integration, 0, len(ids))
	for _, id := range ids {
		// id came from List() over the same registry, which nothing deletes
		// from; a miss is an invariant violation, not a runtime condition.
		in, err := Get(id)
		if err != nil {
			panic(fmt.Sprintf("integrations: registry inconsistency: %v", err))
		}
		out = append(out, in)
	}
	return out
}
