// Package version provides spec version tracking with git-like history.
package version

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/ProductBuildersHQ/visionspec/pkg/config"
	"github.com/ProductBuildersHQ/visionspec/pkg/types"
)

// Version represents a specific version of a spec.
type Version struct {
	Number    int       `json:"number"`              // Sequential version number (1, 2, 3...)
	Hash      string    `json:"hash"`                // SHA256 of content
	Timestamp time.Time `json:"timestamp"`           // When version was created
	Author    string    `json:"author,omitempty"`    // Who created this version
	Message   string    `json:"message,omitempty"`   // Description of changes
	Size      int       `json:"size"`                // Content size in bytes
}

// History tracks all versions of a spec.
type History struct {
	SpecType types.SpecType `json:"spec_type"`
	Versions []Version      `json:"versions"`
}

// historyFilename returns the history filename for a spec type.
func historyFilename(specType types.SpecType) string {
	return string(specType) + ".history.json"
}

// historyPath returns the full path to the history file.
func historyPath(projectPath string, specType types.SpecType) string {
	return filepath.Join(projectPath, config.EvalDir, historyFilename(specType))
}

// versionPath returns the path to a specific version file.
func versionPath(projectPath string, specType types.SpecType, version int) string {
	return filepath.Join(projectPath, config.EvalDir, "versions",
		fmt.Sprintf("%s.v%d.md", specType, version))
}

// computeHash calculates SHA256 hash of content.
func computeHash(content string) string {
	h := sha256.New()
	h.Write([]byte(content))
	return hex.EncodeToString(h.Sum(nil))[:12] // Short hash like git
}

// GetHistory loads the version history for a spec.
func GetHistory(projectPath string, specType types.SpecType) (*History, error) {
	path := historyPath(projectPath, specType)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &History{
				SpecType: specType,
				Versions: []Version{},
			}, nil
		}
		return nil, fmt.Errorf("reading history: %w", err)
	}

	var history History
	if err := json.Unmarshal(data, &history); err != nil {
		return nil, fmt.Errorf("parsing history: %w", err)
	}

	return &history, nil
}

// SaveHistory persists the version history.
func SaveHistory(projectPath string, history *History) error {
	path := historyPath(projectPath, history.SpecType)

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling history: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing history: %w", err)
	}

	return nil
}

// CreateVersion creates a new version from the current spec content.
func CreateVersion(projectPath string, specType types.SpecType, opts CreateOptions) (*Version, error) {
	// Read current spec content
	specPath := config.SpecPath(projectPath, specType)
	content, err := os.ReadFile(specPath)
	if err != nil {
		return nil, fmt.Errorf("reading spec: %w", err)
	}

	// Load history
	history, err := GetHistory(projectPath, specType)
	if err != nil {
		return nil, err
	}

	// Compute hash and check for changes
	hash := computeHash(string(content))
	if len(history.Versions) > 0 {
		latest := history.Versions[len(history.Versions)-1]
		if latest.Hash == hash {
			return nil, ErrNoChanges
		}
	}

	// Create new version
	nextNum := 1
	if len(history.Versions) > 0 {
		nextNum = history.Versions[len(history.Versions)-1].Number + 1
	}

	version := Version{
		Number:    nextNum,
		Hash:      hash,
		Timestamp: time.Now(),
		Author:    opts.Author,
		Message:   opts.Message,
		Size:      len(content),
	}

	// Store version content
	vPath := versionPath(projectPath, specType, nextNum)
	if err := os.MkdirAll(filepath.Dir(vPath), 0755); err != nil {
		return nil, fmt.Errorf("creating versions directory: %w", err)
	}
	if err := os.WriteFile(vPath, content, 0600); err != nil {
		return nil, fmt.Errorf("storing version: %w", err)
	}

	// Update history
	history.Versions = append(history.Versions, version)
	if err := SaveHistory(projectPath, history); err != nil {
		return nil, err
	}

	return &version, nil
}

// CreateOptions configures version creation.
type CreateOptions struct {
	Author  string
	Message string
}

// GetVersion retrieves a specific version.
func GetVersion(projectPath string, specType types.SpecType, versionNum int) (*Version, string, error) {
	history, err := GetHistory(projectPath, specType)
	if err != nil {
		return nil, "", err
	}

	for _, v := range history.Versions {
		if v.Number == versionNum {
			// Load content
			vPath := versionPath(projectPath, specType, versionNum)
			content, err := os.ReadFile(vPath)
			if err != nil {
				return nil, "", fmt.Errorf("reading version file: %w", err)
			}
			return &v, string(content), nil
		}
	}

	return nil, "", ErrVersionNotFound
}

// GetLatestVersion returns the most recent version.
func GetLatestVersion(projectPath string, specType types.SpecType) (*Version, string, error) {
	history, err := GetHistory(projectPath, specType)
	if err != nil {
		return nil, "", err
	}

	if len(history.Versions) == 0 {
		return nil, "", ErrNoVersions
	}

	latest := history.Versions[len(history.Versions)-1]
	vPath := versionPath(projectPath, specType, latest.Number)
	content, err := os.ReadFile(vPath)
	if err != nil {
		return nil, "", fmt.Errorf("reading version file: %w", err)
	}

	return &latest, string(content), nil
}

// Revert restores a spec to a previous version.
func Revert(projectPath string, specType types.SpecType, versionNum int, message string) (*Version, error) {
	// Get target version content
	_, content, err := GetVersion(projectPath, specType, versionNum)
	if err != nil {
		return nil, err
	}

	// Write to spec file
	specPath := config.SpecPath(projectPath, specType)
	if err := os.WriteFile(specPath, []byte(content), 0600); err != nil {
		return nil, fmt.Errorf("writing spec: %w", err)
	}

	// Create new version for the revert
	if message == "" {
		message = fmt.Sprintf("Revert to version %d", versionNum)
	}
	return CreateVersion(projectPath, specType, CreateOptions{
		Message: message,
	})
}

// ListVersions returns all versions for a spec.
func ListVersions(projectPath string, specType types.SpecType) ([]Version, error) {
	history, err := GetHistory(projectPath, specType)
	if err != nil {
		return nil, err
	}

	// Return in reverse chronological order
	versions := make([]Version, len(history.Versions))
	copy(versions, history.Versions)
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].Number > versions[j].Number
	})

	return versions, nil
}

// Errors
var (
	ErrNoChanges       = fmt.Errorf("no changes since last version")
	ErrVersionNotFound = fmt.Errorf("version not found")
	ErrNoVersions      = fmt.Errorf("no versions exist")
)
