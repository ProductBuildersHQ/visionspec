// Package rules provides embedded workflow rules for AI assistant orchestration.
package rules

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
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

	err := fs.WalkDir(embeddedFS, "embedded", func(embeddedPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root "embedded" directory
		if embeddedPath == "embedded" {
			return nil
		}

		// Calculate destination path (strip "embedded/" prefix)
		// embed.FS paths always use forward slashes, so use strings.TrimPrefix
		relPath := strings.TrimPrefix(embeddedPath, "embedded/")
		// Convert to OS-specific path for local filesystem
		destPath := filepath.Join(destDir, filepath.FromSlash(relPath))

		if d.IsDir() {
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return fmt.Errorf("creating directory %s: %w", destPath, err)
			}
			return nil
		}

		// Read file from embedded FS (embeddedPath uses forward slashes)
		content, err := embeddedFS.ReadFile(embeddedPath)
		if err != nil {
			return fmt.Errorf("reading embedded file %s: %w", embeddedPath, err)
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
// Paths are returned with forward slashes for cross-platform consistency.
func List() ([]string, error) {
	var files []string

	err := fs.WalkDir(embeddedFS, "embedded", func(embeddedPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if embeddedPath == "embedded" || d.IsDir() {
			return nil
		}

		// Strip "embedded/" prefix (embed.FS always uses forward slashes)
		relPath := strings.TrimPrefix(embeddedPath, "embedded/")
		files = append(files, relPath)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("listing rules: %w", err)
	}

	return files, nil
}

// Get returns the content of a specific rule file.
func Get(rulePath string) ([]byte, error) {
	// Normalize path separators to forward slashes for embed.FS
	normalizedPath := strings.ReplaceAll(rulePath, "\\", "/")
	// Use path.Join (not filepath.Join) because embed.FS always uses forward slashes
	fullPath := path.Join("embedded", normalizedPath)
	content, err := embeddedFS.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("reading rule %s: %w", rulePath, err)
	}
	return content, nil
}
