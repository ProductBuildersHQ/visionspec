// Package rules provides embedded workflow rules for AI assistant orchestration.
package rules

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed embedded/*
var embeddedFS embed.FS

// Export copies workflow rules to the specified directory.
// If destDir is empty, it defaults to ".visionspec-rules".
func Export(destDir string) ([]string, error) {
	if destDir == "" {
		destDir = ".visionspec-rules"
	}

	// Create the output directory
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return nil, fmt.Errorf("creating output directory: %w", err)
	}

	var files []string

	err := fs.WalkDir(embeddedFS, "embedded", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root "embedded" directory
		if path == "embedded" {
			return nil
		}

		// Calculate destination path (strip "embedded/" prefix)
		relPath, err := filepath.Rel("embedded", path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(destDir, relPath)

		if d.IsDir() {
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return fmt.Errorf("creating directory %s: %w", destPath, err)
			}
			return nil
		}

		// Read file from embedded FS
		content, err := embeddedFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading embedded file %s: %w", path, err)
		}

		// Write file to destination
		if err := os.WriteFile(destPath, content, 0600); err != nil {
			return fmt.Errorf("writing file %s: %w", destPath, err)
		}

		files = append(files, destPath)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("exporting rules: %w", err)
	}

	return files, nil
}

// List returns the paths of all embedded rule files.
func List() ([]string, error) {
	var files []string

	err := fs.WalkDir(embeddedFS, "embedded", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == "embedded" || d.IsDir() {
			return nil
		}

		// Strip "embedded/" prefix
		relPath, err := filepath.Rel("embedded", path)
		if err != nil {
			return err
		}
		files = append(files, relPath)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("listing rules: %w", err)
	}

	return files, nil
}

// Get returns the content of a specific rule file.
func Get(path string) ([]byte, error) {
	fullPath := filepath.Join("embedded", path)
	content, err := embeddedFS.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("reading rule %s: %w", path, err)
	}
	return content, nil
}
