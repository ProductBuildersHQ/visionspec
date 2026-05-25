package context

import (
	"encoding/json"
	"fmt"
	"os"
)

// Snapshot represents a saved context state.
type Snapshot struct {
	// Version of the snapshot format
	Version string `json:"version"`

	// The aggregated context
	Context *AggregatedContext `json:"context"`
}

const snapshotVersion = "1.0"

// SaveSnapshot saves an aggregated context to a file.
func SaveSnapshot(ac *AggregatedContext, path string) error {
	snapshot := &Snapshot{
		Version: snapshotVersion,
		Context: ac,
	}

	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling snapshot: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing snapshot: %w", err)
	}

	return nil
}

// LoadSnapshot loads an aggregated context from a file.
func LoadSnapshot(path string) (*AggregatedContext, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading snapshot: %w", err)
	}

	var snapshot Snapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return nil, fmt.Errorf("parsing snapshot: %w", err)
	}

	// Version compatibility check
	if snapshot.Version != snapshotVersion {
		// For now, just warn - future versions may need migration
		fmt.Fprintf(os.Stderr, "Warning: snapshot version %s, current %s\n", snapshot.Version, snapshotVersion)
	}

	return snapshot.Context, nil
}

// ToJSON serializes the aggregated context to JSON.
func (ac *AggregatedContext) ToJSON() ([]byte, error) {
	return json.MarshalIndent(ac, "", "  ")
}

// FromJSON deserializes an aggregated context from JSON.
func FromJSON(data []byte) (*AggregatedContext, error) {
	var ac AggregatedContext
	if err := json.Unmarshal(data, &ac); err != nil {
		return nil, err
	}
	return &ac, nil
}
