package context

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Aggregator fetches and combines context from multiple sources.
type Aggregator struct {
	sources     []Source
	cache       *Cache
	config      *Config
	projectName string
}

// NewAggregator creates an aggregator from configuration.
// Sources are created based on the config but not yet fetched.
func NewAggregator(projectName string, cfg *Config) *Aggregator {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	cacheTTL := cfg.CacheTTL
	if cacheTTL == 0 {
		cacheTTL = 5 * time.Minute
	}

	return &Aggregator{
		sources:     make([]Source, 0),
		cache:       NewCache(cacheTTL),
		config:      cfg,
		projectName: projectName,
	}
}

// AddSource adds a source to the aggregator.
func (a *Aggregator) AddSource(src Source) {
	a.sources = append(a.sources, src)
}

// Sources returns all registered sources.
func (a *Aggregator) Sources() []Source {
	return a.sources
}

// SourceCount returns the number of sources.
func (a *Aggregator) SourceCount() int {
	return len(a.sources)
}

// Gather fetches context from all sources concurrently.
func (a *Aggregator) Gather(ctx context.Context) (*AggregatedContext, error) {
	start := time.Now()

	if len(a.sources) == 0 {
		return &AggregatedContext{
			Project:    a.projectName,
			GatheredAt: time.Now(),
			Duration:   time.Since(start),
			Sources:    []*ContextData{},
		}, nil
	}

	results := make([]*ContextData, len(a.sources))
	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error

	for i, src := range a.sources {
		wg.Add(1)
		go func(idx int, s Source) {
			defer wg.Done()

			// Check cache first
			if cached := a.cache.Get(s.Name()); cached != nil {
				mu.Lock()
				results[idx] = cached
				mu.Unlock()
				return
			}

			// Fetch from source
			fetchStart := time.Now()
			data, err := s.Fetch(ctx)
			fetchDuration := time.Since(fetchStart)

			if err != nil {
				data = &ContextData{
					Source:    s.Name(),
					Type:      s.Type(),
					FetchedAt: time.Now(),
					Duration:  fetchDuration,
					Errors:    []string{err.Error()},
				}
				mu.Lock()
				if firstErr == nil {
					firstErr = fmt.Errorf("source %s: %w", s.Name(), err)
				}
				mu.Unlock()
			} else {
				data.Duration = fetchDuration
			}

			// Cache successful results
			if len(data.Errors) == 0 {
				a.cache.Set(s.Name(), data)
			}

			mu.Lock()
			results[idx] = data
			mu.Unlock()
		}(i, src)
	}

	wg.Wait()

	// Build aggregated context
	ac := &AggregatedContext{
		Project:    a.projectName,
		GatheredAt: time.Now(),
		Duration:   time.Since(start),
		Sources:    results,
	}

	// Set helper flags and count errors
	for _, src := range results {
		if src == nil {
			continue
		}
		switch src.Type {
		case SourceTypeGit:
			if src.Code != nil {
				ac.HasCode = true
			}
		case SourceTypeGraphize:
			if src.Graph != nil {
				ac.HasGraph = true
			}
		case SourceTypeMCP:
			if src.External != nil {
				ac.HasExternal = true
			}
		case SourceTypeFile:
			if src.File != nil {
				ac.HasFiles = true
			}
		}
		ac.ErrorCount += len(src.Errors)
	}

	// Generate combined summary
	ac.Summary = ac.GenerateSummary()

	return ac, nil
}

// GatherWithTimeout fetches context with a timeout.
func (a *Aggregator) GatherWithTimeout(timeout time.Duration) (*AggregatedContext, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return a.Gather(ctx)
}

// Refresh clears the cache and re-fetches all sources.
func (a *Aggregator) Refresh(ctx context.Context) (*AggregatedContext, error) {
	a.cache.Clear()
	return a.Gather(ctx)
}

// ClearCache clears the context cache.
func (a *Aggregator) ClearCache() {
	a.cache.Clear()
}

// Cache provides a simple in-memory cache with TTL.
type Cache struct {
	entries map[string]*cacheEntry
	ttl     time.Duration
	mu      sync.RWMutex
}

type cacheEntry struct {
	data      *ContextData
	expiresAt time.Time
}

// NewCache creates a new cache with the given TTL.
func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		entries: make(map[string]*cacheEntry),
		ttl:     ttl,
	}
}

// Get retrieves a cached entry if it exists and hasn't expired.
func (c *Cache) Get(key string) *ContextData {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.entries[key]
	if !ok {
		return nil
	}

	if time.Now().After(entry.expiresAt) {
		return nil
	}

	return entry.data
}

// Set stores a value in the cache.
func (c *Cache) Set(key string, data *ContextData) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = &cacheEntry{
		data:      data,
		expiresAt: time.Now().Add(c.ttl),
	}
}

// Clear removes all entries from the cache.
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries = make(map[string]*cacheEntry)
}

// Size returns the number of cached entries.
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.entries)
}

// Cleanup removes expired entries.
func (c *Cache) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, entry := range c.entries {
		if now.After(entry.expiresAt) {
			delete(c.entries, key)
		}
	}
}
