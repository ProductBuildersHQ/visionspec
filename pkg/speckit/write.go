package speckit

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

// WriteFiles persists a generated file set under dir, creating directories
// as needed.
func WriteFiles(dir string, files []File) error {
	for _, f := range files {
		dest := filepath.Join(dir, filepath.FromSlash(f.Path))
		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return fmt.Errorf("creating directory for %s: %w", f.Path, err)
		}
		if err := os.WriteFile(dest, f.Content, 0600); err != nil {
			return fmt.Errorf("writing %s: %w", f.Path, err)
		}
	}
	return nil
}

// Check compares a generated file set against the copy under dir and returns
// the relative paths that are missing or differ. An empty result means the
// committed artifacts are in sync with the generator.
func Check(dir string, files []File) ([]string, error) {
	var drift []string
	for _, f := range files {
		dest := filepath.Join(dir, filepath.FromSlash(f.Path))
		existing, err := os.ReadFile(dest)
		if os.IsNotExist(err) {
			drift = append(drift, f.Path)
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", f.Path, err)
		}
		if !bytes.Equal(existing, f.Content) {
			drift = append(drift, f.Path)
		}
	}
	return drift, nil
}
