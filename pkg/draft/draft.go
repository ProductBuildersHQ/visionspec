// Package draft provides draft management for spec authoring.
package draft

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ProductBuildersHQ/visionspec/pkg/config"
	"github.com/ProductBuildersHQ/visionspec/pkg/templates"
	"github.com/ProductBuildersHQ/visionspec/pkg/types"
)

// ErrDraftNotFound indicates no draft exists.
var ErrDraftNotFound = errors.New("draft not found")

// ErrDraftExists indicates a draft already exists.
var ErrDraftExists = errors.New("draft already exists")

// ErrSpecExists indicates the final spec already exists.
var ErrSpecExists = errors.New("spec already exists")

// Draft represents a draft document with its metadata.
type Draft struct {
	Metadata Metadata
	Content  string
}

// Metadata contains draft metadata.
type Metadata struct {
	SpecType    types.SpecType `json:"spec_type"`
	ProjectName string         `json:"project_name"`
	StartedAt   time.Time      `json:"started_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Version     int            `json:"version"`
	EvalHistory []EvalSummary  `json:"eval_history,omitempty"`
}

// EvalSummary summarizes a previous evaluation.
type EvalSummary struct {
	Timestamp time.Time `json:"timestamp"`
	Score     float64   `json:"score"`
	Passed    bool      `json:"passed"`
	Findings  int       `json:"findings"`
}

// draftFilename returns the draft content filename.
func draftFilename(specType types.SpecType) string {
	return string(specType) + ".draft.md"
}

// metadataFilename returns the draft metadata filename.
func metadataFilename(specType types.SpecType) string {
	return string(specType) + ".draft.json"
}

// DraftPath returns the full path for a draft file.
func DraftPath(projectPath string, specType types.SpecType) string {
	dir := specType.Dir()
	if dir == "" {
		return filepath.Join(projectPath, draftFilename(specType))
	}
	return filepath.Join(projectPath, dir, draftFilename(specType))
}

// MetadataPath returns the full path for a draft metadata file.
func MetadataPath(projectPath string, specType types.SpecType) string {
	dir := specType.Dir()
	if dir == "" {
		return filepath.Join(projectPath, metadataFilename(specType))
	}
	return filepath.Join(projectPath, dir, metadataFilename(specType))
}

// Exists checks if a draft exists for the given project and spec type.
func Exists(projectPath string, specType types.SpecType) bool {
	draftPath := DraftPath(projectPath, specType)
	_, err := os.Stat(draftPath)
	return err == nil
}

// SpecExists checks if the final spec already exists.
func SpecExists(projectPath string, specType types.SpecType) bool {
	specPath := config.SpecPath(projectPath, specType)
	_, err := os.Stat(specPath)
	return err == nil
}

// Start initializes a new draft for a spec type.
// Returns the template content and metadata, or error if draft/spec already exists.
func Start(projectPath string, specType types.SpecType) (*Draft, error) {
	// Check if draft already exists
	if Exists(projectPath, specType) {
		return nil, ErrDraftExists
	}

	// Get the template
	tmpl, err := templates.Get(specType)
	if err != nil {
		return nil, fmt.Errorf("no template for spec type %s: %w", specType, err)
	}

	// Get project name from path
	projectName := filepath.Base(projectPath)

	// Render template
	opts := templates.RenderOptions{
		ProjectName: projectName,
		Author:      "",
		Date:        time.Now().Format("2006-01-02"),
	}
	content := tmpl.Render(opts)

	// Create metadata
	now := time.Now()
	metadata := Metadata{
		SpecType:    specType,
		ProjectName: projectName,
		StartedAt:   now,
		UpdatedAt:   now,
		Version:     1,
	}

	// Ensure directory exists
	dir := filepath.Dir(DraftPath(projectPath, specType))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Write draft content
	if err := os.WriteFile(DraftPath(projectPath, specType), []byte(content), 0600); err != nil {
		return nil, fmt.Errorf("failed to write draft: %w", err)
	}

	// Write metadata
	metaData, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %w", err)
	}
	if err := os.WriteFile(MetadataPath(projectPath, specType), metaData, 0600); err != nil {
		return nil, fmt.Errorf("failed to write metadata: %w", err)
	}

	return &Draft{
		Metadata: metadata,
		Content:  content,
	}, nil
}

// Get retrieves an existing draft.
func Get(projectPath string, specType types.SpecType) (*Draft, error) {
	if !Exists(projectPath, specType) {
		return nil, ErrDraftNotFound
	}

	// Read content
	content, err := os.ReadFile(DraftPath(projectPath, specType))
	if err != nil {
		return nil, fmt.Errorf("failed to read draft: %w", err)
	}

	// Read metadata
	metaBytes, err := os.ReadFile(MetadataPath(projectPath, specType))
	if err != nil {
		// Create default metadata if missing
		return &Draft{
			Metadata: Metadata{
				SpecType:    specType,
				ProjectName: filepath.Base(projectPath),
				UpdatedAt:   time.Now(),
				Version:     1,
			},
			Content: string(content),
		}, nil
	}

	var metadata Metadata
	if err := json.Unmarshal(metaBytes, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	return &Draft{
		Metadata: metadata,
		Content:  string(content),
	}, nil
}

// Update saves updated content to the draft.
func Update(projectPath string, specType types.SpecType, content string) (*Draft, error) {
	draft, err := Get(projectPath, specType)
	if err != nil {
		return nil, err
	}

	// Update content
	draft.Content = content
	draft.Metadata.UpdatedAt = time.Now()
	draft.Metadata.Version++

	// Write content
	if err := os.WriteFile(DraftPath(projectPath, specType), []byte(content), 0600); err != nil {
		return nil, fmt.Errorf("failed to write draft: %w", err)
	}

	// Write metadata
	metaData, err := json.MarshalIndent(draft.Metadata, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %w", err)
	}
	if err := os.WriteFile(MetadataPath(projectPath, specType), metaData, 0600); err != nil {
		return nil, fmt.Errorf("failed to write metadata: %w", err)
	}

	return draft, nil
}

// Discard removes a draft and its metadata.
func Discard(projectPath string, specType types.SpecType) error {
	if !Exists(projectPath, specType) {
		return ErrDraftNotFound
	}

	// Remove draft file
	if err := os.Remove(DraftPath(projectPath, specType)); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove draft: %w", err)
	}

	// Remove metadata file
	if err := os.Remove(MetadataPath(projectPath, specType)); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove metadata: %w", err)
	}

	return nil
}

// Finalize promotes a draft to the final spec location.
func Finalize(projectPath string, specType types.SpecType) error {
	draft, err := Get(projectPath, specType)
	if err != nil {
		return err
	}

	specPath := config.SpecPath(projectPath, specType)

	// Ensure directory exists
	dir := filepath.Dir(specPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write to final location
	if err := os.WriteFile(specPath, []byte(draft.Content), 0600); err != nil { //nolint:gosec // G703: specPath from config.SpecPath with validated project path
		return fmt.Errorf("failed to write spec: %w", err)
	}

	// Remove draft files
	if err := Discard(projectPath, specType); err != nil {
		// Log but don't fail - spec was written successfully
		return nil
	}

	return nil
}

// AddEvalResult adds an evaluation result to the draft metadata.
func AddEvalResult(projectPath string, specType types.SpecType, score float64, passed bool, findings int) error {
	draft, err := Get(projectPath, specType)
	if err != nil {
		return err
	}

	summary := EvalSummary{
		Timestamp: time.Now(),
		Score:     score,
		Passed:    passed,
		Findings:  findings,
	}

	draft.Metadata.EvalHistory = append(draft.Metadata.EvalHistory, summary)
	draft.Metadata.UpdatedAt = time.Now()

	// Write metadata
	metaData, err := json.MarshalIndent(draft.Metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}
	if err := os.WriteFile(MetadataPath(projectPath, specType), metaData, 0600); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	return nil
}

// ListDrafts returns all drafts in a project.
func ListDrafts(projectPath string) ([]*Draft, error) {
	var drafts []*Draft

	for _, specType := range templates.Available() {
		if Exists(projectPath, specType) {
			draft, err := Get(projectPath, specType)
			if err != nil {
				continue
			}
			drafts = append(drafts, draft)
		}
	}

	return drafts, nil
}
