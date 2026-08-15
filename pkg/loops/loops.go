// Package loops provides the embedded default library of loop systems (see
// pkg/loop): the ProductBuildersHQ two-loop system and reference systems such
// as AWS AI-DLC.
//
// Systems load from embedded YAML with no filesystem access, mirroring the
// pkg/workflows, pkg/integrations, and pkg/pipelines patterns. This package
// holds data only.
package loops

import (
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"
	"sync"

	"github.com/ProductBuildersHQ/visionspec/pkg/loop"
)

//go:embed default/*.yaml
var defaultFS embed.FS

var (
	mu       sync.RWMutex
	registry = map[string]*loop.System{}
)

func init() {
	if err := loadAll(); err != nil {
		panic(fmt.Sprintf("failed to load embedded loop systems: %v", err))
	}
}

// loadAll discovers and loads every loop system from the embedded filesystem.
func loadAll() error {
	entries, err := fs.ReadDir(defaultFS, "default")
	if err != nil {
		return fmt.Errorf("reading default loop systems directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}

		data, err := fs.ReadFile(defaultFS, "default/"+entry.Name())
		if err != nil {
			return fmt.Errorf("reading loop system %s: %w", entry.Name(), err)
		}

		s, err := loop.ParseYAML(data)
		if err != nil {
			return fmt.Errorf("parsing loop system %s: %w", entry.Name(), err)
		}
		if err := s.Validate(); err != nil {
			return fmt.Errorf("invalid loop system %s: %w", entry.Name(), err)
		}

		Register(s)
	}

	return nil
}

// Register adds or replaces a loop system in the registry.
func Register(s *loop.System) {
	mu.Lock()
	defer mu.Unlock()
	registry[s.ID] = s
}

// Get returns the loop system with the given ID.
func Get(id string) (*loop.System, error) {
	mu.RLock()
	defer mu.RUnlock()
	s, ok := registry[id]
	if !ok {
		return nil, fmt.Errorf("loop system %q not found", id)
	}
	return s, nil
}

// List returns all registered loop system IDs, sorted.
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

// All returns every registered loop system, sorted by ID.
func All() []*loop.System {
	ids := List()
	out := make([]*loop.System, 0, len(ids))
	for _, id := range ids {
		// id came from List() over the same registry, which nothing deletes
		// from; a miss is an invariant violation, not a runtime condition.
		s, err := Get(id)
		if err != nil {
			panic(fmt.Sprintf("loops: registry inconsistency: %v", err))
		}
		out = append(out, s)
	}
	return out
}
