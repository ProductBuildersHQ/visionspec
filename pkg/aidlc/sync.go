package aidlc

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// SyncDirection indicates the direction of synchronization.
type SyncDirection string

const (
	// SyncToAIDLC exports VisionSpec documents to AIDLC format.
	SyncToAIDLC SyncDirection = "to_aidlc"
	// SyncFromAIDLC imports AIDLC documents to VisionSpec format.
	SyncFromAIDLC SyncDirection = "from_aidlc"
	// SyncBidirectional performs two-way sync based on timestamps.
	SyncBidirectional SyncDirection = "bidirectional"
)

// SyncAction represents a single sync operation.
type SyncAction struct {
	// Direction is the sync direction.
	Direction SyncDirection `json:"direction" yaml:"direction"`

	// DocType is the document type being synced.
	DocType DocType `json:"doc_type" yaml:"doc_type"`

	// SourcePath is the source file path.
	SourcePath string `json:"source_path" yaml:"source_path"`

	// DestPath is the destination file path.
	DestPath string `json:"dest_path" yaml:"dest_path"`

	// Action describes the operation (create, update, delete).
	Action string `json:"action" yaml:"action"`

	// Reason explains why the action is needed.
	Reason string `json:"reason,omitempty" yaml:"reason,omitempty"`
}

// SyncDiff represents the difference between VisionSpec and AIDLC directories.
type SyncDiff struct {
	// VisionSpecDir is the .visionspec directory path.
	VisionSpecDir string `json:"visionspec_dir" yaml:"visionspec_dir"`

	// AIDLCDocsDir is the aidlc-docs directory path.
	AIDLCDocsDir string `json:"aidlc_docs_dir" yaml:"aidlc_docs_dir"`

	// Actions lists required sync operations.
	Actions []SyncAction `json:"actions" yaml:"actions"`

	// Conflicts lists documents with conflicts requiring resolution.
	Conflicts []SyncConflict `json:"conflicts,omitempty" yaml:"conflicts,omitempty"`

	// ComputedAt is when the diff was computed.
	ComputedAt time.Time `json:"computed_at" yaml:"computed_at"`
}

// HasChanges returns whether there are any changes to sync.
func (d *SyncDiff) HasChanges() bool {
	return len(d.Actions) > 0 || len(d.Conflicts) > 0
}

// SyncConflict represents a sync conflict requiring resolution.
type SyncConflict struct {
	// DocType is the conflicting document type.
	DocType DocType `json:"doc_type" yaml:"doc_type"`

	// VisionSpecPath is the VisionSpec file path.
	VisionSpecPath string `json:"visionspec_path" yaml:"visionspec_path"`

	// AIDLCPath is the AIDLC file path.
	AIDLCPath string `json:"aidlc_path" yaml:"aidlc_path"`

	// VisionSpecModTime is when the VisionSpec file was modified.
	VisionSpecModTime time.Time `json:"visionspec_mod_time" yaml:"visionspec_mod_time"`

	// AIDLCModTime is when the AIDLC file was modified.
	AIDLCModTime time.Time `json:"aidlc_mod_time" yaml:"aidlc_mod_time"`

	// Reason explains the conflict.
	Reason string `json:"reason" yaml:"reason"`
}

// SyncResult contains the result of a sync operation.
type SyncResult struct {
	// Direction is the sync direction performed.
	Direction SyncDirection `json:"direction" yaml:"direction"`

	// Created lists newly created files.
	Created []string `json:"created,omitempty" yaml:"created,omitempty"`

	// Updated lists updated files.
	Updated []string `json:"updated,omitempty" yaml:"updated,omitempty"`

	// Skipped lists skipped files (conflicts or errors).
	Skipped []string `json:"skipped,omitempty" yaml:"skipped,omitempty"`

	// Errors lists any errors encountered.
	Errors []string `json:"errors,omitempty" yaml:"errors,omitempty"`

	// CompletedAt is when the sync completed.
	CompletedAt time.Time `json:"completed_at" yaml:"completed_at"`
}

// Success returns whether the sync completed without errors.
func (r *SyncResult) Success() bool {
	return len(r.Errors) == 0
}

// SyncEngine handles bidirectional sync between VisionSpec and AIDLC directories.
type SyncEngine struct {
	// VisionSpecDir is the .visionspec directory path.
	VisionSpecDir string

	// AIDLCDocsDir is the aidlc-docs directory path.
	AIDLCDocsDir string

	// ConflictStrategy determines how conflicts are resolved.
	ConflictStrategy ConflictStrategy

	// DryRun prevents actual file modifications when true.
	DryRun bool
}

// ConflictStrategy determines how sync conflicts are resolved.
type ConflictStrategy string

const (
	// ConflictSkip skips conflicting files.
	ConflictSkip ConflictStrategy = "skip"
	// ConflictNewerWins uses the newer file.
	ConflictNewerWins ConflictStrategy = "newer_wins"
	// ConflictVisionSpecWins prefers VisionSpec files.
	ConflictVisionSpecWins ConflictStrategy = "visionspec_wins"
	// ConflictAIDLCWins prefers AIDLC files.
	ConflictAIDLCWins ConflictStrategy = "aidlc_wins"
)

// NewSyncEngine creates a new sync engine.
func NewSyncEngine(visionSpecDir, aidlcDocsDir string) *SyncEngine {
	return &SyncEngine{
		VisionSpecDir:    visionSpecDir,
		AIDLCDocsDir:     aidlcDocsDir,
		ConflictStrategy: ConflictNewerWins,
		DryRun:           false,
	}
}

// DiffState computes the difference between VisionSpec and AIDLC directories.
func (e *SyncEngine) DiffState(ctx context.Context) (*SyncDiff, error) {
	diff := &SyncDiff{
		VisionSpecDir: e.VisionSpecDir,
		AIDLCDocsDir:  e.AIDLCDocsDir,
		Actions:       make([]SyncAction, 0),
		Conflicts:     make([]SyncConflict, 0),
		ComputedAt:    time.Now(),
	}

	// Scan VisionSpec documents
	visionSpecDocs, err := e.scanVisionSpec(ctx)
	if err != nil {
		return nil, fmt.Errorf("scan visionspec: %w", err)
	}

	// Scan AIDLC documents
	aidlcDocs, err := e.scanAIDLC(ctx)
	if err != nil {
		return nil, fmt.Errorf("scan aidlc: %w", err)
	}

	// Compare documents
	for docType, vsDoc := range visionSpecDocs {
		aidlcDoc, exists := aidlcDocs[docType]
		if !exists {
			// Document only in VisionSpec - export to AIDLC
			diff.Actions = append(diff.Actions, SyncAction{
				Direction:  SyncToAIDLC,
				DocType:    docType,
				SourcePath: vsDoc.Path,
				DestPath:   e.aidlcPathForDoc(docType),
				Action:     "create",
				Reason:     "Document exists in VisionSpec but not in AIDLC",
			})
			continue
		}

		// Both exist - check for differences
		if vsDoc.Checksum != aidlcDoc.Checksum {
			// Content differs - check timestamps for conflict
			if vsDoc.UpdatedAt.After(aidlcDoc.UpdatedAt) {
				diff.Actions = append(diff.Actions, SyncAction{
					Direction:  SyncToAIDLC,
					DocType:    docType,
					SourcePath: vsDoc.Path,
					DestPath:   aidlcDoc.Path,
					Action:     "update",
					Reason:     "VisionSpec version is newer",
				})
			} else if aidlcDoc.UpdatedAt.After(vsDoc.UpdatedAt) {
				diff.Actions = append(diff.Actions, SyncAction{
					Direction:  SyncFromAIDLC,
					DocType:    docType,
					SourcePath: aidlcDoc.Path,
					DestPath:   vsDoc.Path,
					Action:     "update",
					Reason:     "AIDLC version is newer",
				})
			} else {
				// Same timestamp but different content - conflict
				diff.Conflicts = append(diff.Conflicts, SyncConflict{
					DocType:           docType,
					VisionSpecPath:    vsDoc.Path,
					AIDLCPath:         aidlcDoc.Path,
					VisionSpecModTime: vsDoc.UpdatedAt,
					AIDLCModTime:      aidlcDoc.UpdatedAt,
					Reason:            "Content differs with same modification time",
				})
			}
		}
	}

	// Check for documents only in AIDLC
	for docType, aidlcDoc := range aidlcDocs {
		if _, exists := visionSpecDocs[docType]; !exists {
			diff.Actions = append(diff.Actions, SyncAction{
				Direction:  SyncFromAIDLC,
				DocType:    docType,
				SourcePath: aidlcDoc.Path,
				DestPath:   e.visionSpecPathForDoc(docType),
				Action:     "create",
				Reason:     "Document exists in AIDLC but not in VisionSpec",
			})
		}
	}

	return diff, nil
}

// ExportToAIDLC exports VisionSpec documents to AIDLC format.
func (e *SyncEngine) ExportToAIDLC(ctx context.Context) (*SyncResult, error) {
	result := &SyncResult{
		Direction:   SyncToAIDLC,
		Created:     make([]string, 0),
		Updated:     make([]string, 0),
		Skipped:     make([]string, 0),
		Errors:      make([]string, 0),
		CompletedAt: time.Now(),
	}

	// Ensure AIDLC directory exists
	if !e.DryRun {
		if err := os.MkdirAll(e.AIDLCDocsDir, 0755); err != nil {
			return nil, fmt.Errorf("create aidlc directory: %w", err)
		}
	}

	// Scan VisionSpec documents
	visionSpecDocs, err := e.scanVisionSpec(ctx)
	if err != nil {
		return nil, fmt.Errorf("scan visionspec: %w", err)
	}

	for docType, doc := range visionSpecDocs {
		destPath := e.aidlcPathForDoc(docType)

		// Check if AIDLC file exists
		_, err := os.Stat(destPath)
		isNew := os.IsNotExist(err)

		if e.DryRun {
			if isNew {
				result.Created = append(result.Created, destPath)
			} else {
				result.Updated = append(result.Updated, destPath)
			}
			continue
		}

		// Create phase subdirectory
		phaseDir := filepath.Dir(destPath)
		if err := os.MkdirAll(phaseDir, 0755); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("create directory %s: %v", phaseDir, err))
			continue
		}

		// Convert and write AIDLC format
		aidlcContent, err := e.convertToAIDLCFormat(doc)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("convert %s: %v", doc.Path, err))
			continue
		}

		if err := os.WriteFile(destPath, []byte(aidlcContent), 0644); err != nil { //nolint:gosec // G306: Document files need to be readable
			result.Errors = append(result.Errors, fmt.Sprintf("write %s: %v", destPath, err))
			continue
		}

		if isNew {
			result.Created = append(result.Created, destPath)
		} else {
			result.Updated = append(result.Updated, destPath)
		}
	}

	result.CompletedAt = time.Now()
	return result, nil
}

// ImportFromAIDLC imports AIDLC documents to VisionSpec format.
func (e *SyncEngine) ImportFromAIDLC(ctx context.Context) (*SyncResult, error) {
	result := &SyncResult{
		Direction:   SyncFromAIDLC,
		Created:     make([]string, 0),
		Updated:     make([]string, 0),
		Skipped:     make([]string, 0),
		Errors:      make([]string, 0),
		CompletedAt: time.Now(),
	}

	// Ensure VisionSpec directory exists
	if !e.DryRun {
		if err := os.MkdirAll(e.VisionSpecDir, 0755); err != nil {
			return nil, fmt.Errorf("create visionspec directory: %w", err)
		}
	}

	// Scan AIDLC documents
	aidlcDocs, err := e.scanAIDLC(ctx)
	if err != nil {
		return nil, fmt.Errorf("scan aidlc: %w", err)
	}

	for docType, doc := range aidlcDocs {
		destPath := e.visionSpecPathForDoc(docType)

		// Check if VisionSpec file exists
		_, err := os.Stat(destPath)
		isNew := os.IsNotExist(err)

		if e.DryRun {
			if isNew {
				result.Created = append(result.Created, destPath)
			} else {
				result.Updated = append(result.Updated, destPath)
			}
			continue
		}

		// Convert and write VisionSpec format
		vsContent, err := e.convertToVisionSpecFormat(doc)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("convert %s: %v", doc.Path, err))
			continue
		}

		if err := os.WriteFile(destPath, []byte(vsContent), 0644); err != nil { //nolint:gosec // G306: Document files need to be readable
			result.Errors = append(result.Errors, fmt.Sprintf("write %s: %v", destPath, err))
			continue
		}

		if isNew {
			result.Created = append(result.Created, destPath)
		} else {
			result.Updated = append(result.Updated, destPath)
		}
	}

	result.CompletedAt = time.Now()
	return result, nil
}

// Sync performs bidirectional sync based on the diff.
func (e *SyncEngine) Sync(ctx context.Context) (*SyncResult, error) {
	diff, err := e.DiffState(ctx)
	if err != nil {
		return nil, fmt.Errorf("compute diff: %w", err)
	}

	result := &SyncResult{
		Direction:   SyncBidirectional,
		Created:     make([]string, 0),
		Updated:     make([]string, 0),
		Skipped:     make([]string, 0),
		Errors:      make([]string, 0),
		CompletedAt: time.Now(),
	}

	// Handle conflicts based on strategy
	for _, conflict := range diff.Conflicts {
		switch e.ConflictStrategy {
		case ConflictSkip:
			result.Skipped = append(result.Skipped, conflict.VisionSpecPath)
		case ConflictNewerWins:
			if conflict.VisionSpecModTime.After(conflict.AIDLCModTime) {
				diff.Actions = append(diff.Actions, SyncAction{
					Direction:  SyncToAIDLC,
					DocType:    conflict.DocType,
					SourcePath: conflict.VisionSpecPath,
					DestPath:   conflict.AIDLCPath,
					Action:     "update",
					Reason:     "Conflict resolved: VisionSpec is newer",
				})
			} else {
				diff.Actions = append(diff.Actions, SyncAction{
					Direction:  SyncFromAIDLC,
					DocType:    conflict.DocType,
					SourcePath: conflict.AIDLCPath,
					DestPath:   conflict.VisionSpecPath,
					Action:     "update",
					Reason:     "Conflict resolved: AIDLC is newer",
				})
			}
		case ConflictVisionSpecWins:
			diff.Actions = append(diff.Actions, SyncAction{
				Direction:  SyncToAIDLC,
				DocType:    conflict.DocType,
				SourcePath: conflict.VisionSpecPath,
				DestPath:   conflict.AIDLCPath,
				Action:     "update",
				Reason:     "Conflict resolved: VisionSpec wins",
			})
		case ConflictAIDLCWins:
			diff.Actions = append(diff.Actions, SyncAction{
				Direction:  SyncFromAIDLC,
				DocType:    conflict.DocType,
				SourcePath: conflict.AIDLCPath,
				DestPath:   conflict.VisionSpecPath,
				Action:     "update",
				Reason:     "Conflict resolved: AIDLC wins",
			})
		}
	}

	// Execute actions
	for _, action := range diff.Actions {
		if e.DryRun {
			if action.Action == "create" {
				result.Created = append(result.Created, action.DestPath)
			} else {
				result.Updated = append(result.Updated, action.DestPath)
			}
			continue
		}

		// Ensure destination directory exists
		destDir := filepath.Dir(action.DestPath)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("create directory %s: %v", destDir, err))
			continue
		}

		// Copy file
		if err := copyFile(action.SourcePath, action.DestPath); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("copy %s to %s: %v", action.SourcePath, action.DestPath, err))
			continue
		}

		if action.Action == "create" {
			result.Created = append(result.Created, action.DestPath)
		} else {
			result.Updated = append(result.Updated, action.DestPath)
		}
	}

	result.CompletedAt = time.Now()
	return result, nil
}

// scanVisionSpec scans the VisionSpec directory for documents.
func (e *SyncEngine) scanVisionSpec(ctx context.Context) (map[DocType]*Document, error) {
	docs := make(map[DocType]*Document)

	if _, err := os.Stat(e.VisionSpecDir); os.IsNotExist(err) {
		return docs, nil
	}

	err := filepath.Walk(e.VisionSpecDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".md") && !strings.HasSuffix(path, ".yaml") && !strings.HasSuffix(path, ".json") {
			return nil
		}

		doc, err := e.parseVisionSpecFile(path, info)
		if err != nil {
			return nil // Skip unparseable files
		}
		if doc.Type != "" {
			docs[doc.Type] = doc
		}
		return nil
	})

	return docs, err
}

// scanAIDLC scans the AIDLC directory for documents.
func (e *SyncEngine) scanAIDLC(ctx context.Context) (map[DocType]*Document, error) {
	docs := make(map[DocType]*Document)

	if _, err := os.Stat(e.AIDLCDocsDir); os.IsNotExist(err) {
		return docs, nil
	}

	err := filepath.Walk(e.AIDLCDocsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		doc, err := e.parseAIDLCFile(path, info)
		if err != nil {
			return nil // Skip unparseable files
		}
		if doc.Type != "" {
			docs[doc.Type] = doc
		}
		return nil
	})

	return docs, err
}

// parseVisionSpecFile parses a VisionSpec file.
func (e *SyncEngine) parseVisionSpecFile(path string, info os.FileInfo) (*Document, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	doc := &Document{
		Path:      path,
		Content:   string(content),
		UpdatedAt: info.ModTime(),
		Checksum:  computeChecksum(content),
	}

	// Infer document type from filename
	baseName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	doc.Type = inferDocTypeFromName(baseName)

	// Parse frontmatter if present
	if strings.HasPrefix(doc.Content, "---") {
		parts := strings.SplitN(doc.Content, "---", 3)
		if len(parts) >= 3 {
			var meta map[string]any
			if err := yaml.Unmarshal([]byte(parts[1]), &meta); err == nil {
				doc.Metadata = meta
				if t, ok := meta["type"].(string); ok {
					doc.Type = DocType(t)
				}
				if title, ok := meta["title"].(string); ok {
					doc.Title = title
				}
			}
		}
	}

	if doc.Type != "" {
		doc.Phase = doc.Type.Phase()
	}

	return doc, nil
}

// parseAIDLCFile parses an AIDLC file.
func (e *SyncEngine) parseAIDLCFile(path string, info os.FileInfo) (*Document, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	doc := &Document{
		Path:      path,
		Content:   string(content),
		UpdatedAt: info.ModTime(),
		Checksum:  computeChecksum(content),
	}

	// Infer type from path structure (aidlc-docs/{phase}/{doc_type}.md)
	relPath, _ := filepath.Rel(e.AIDLCDocsDir, path)
	parts := strings.Split(relPath, string(filepath.Separator))
	if len(parts) >= 2 {
		doc.Phase = Phase(parts[0])
		baseName := strings.TrimSuffix(parts[len(parts)-1], ".md")
		doc.Type = DocType(baseName)
	} else {
		baseName := strings.TrimSuffix(filepath.Base(path), ".md")
		doc.Type = inferDocTypeFromName(baseName)
		if doc.Type != "" {
			doc.Phase = doc.Type.Phase()
		}
	}

	// Parse frontmatter
	if strings.HasPrefix(doc.Content, "---") {
		parts := strings.SplitN(doc.Content, "---", 3)
		if len(parts) >= 3 {
			var meta map[string]any
			if err := yaml.Unmarshal([]byte(parts[1]), &meta); err == nil {
				doc.Metadata = meta
				if title, ok := meta["title"].(string); ok {
					doc.Title = title
				}
			}
		}
	}

	return doc, nil
}

// aidlcPathForDoc returns the AIDLC path for a document type.
func (e *SyncEngine) aidlcPathForDoc(docType DocType) string {
	return filepath.Join(e.AIDLCDocsDir, string(docType.Phase()), docType.Filename())
}

// visionSpecPathForDoc returns the VisionSpec path for a document type.
func (e *SyncEngine) visionSpecPathForDoc(docType DocType) string {
	return filepath.Join(e.VisionSpecDir, docType.Filename())
}

// convertToAIDLCFormat converts a VisionSpec document to AIDLC format.
func (e *SyncEngine) convertToAIDLCFormat(doc *Document) (string, error) {
	// AIDLC uses markdown with YAML frontmatter
	var sb strings.Builder

	// Write frontmatter
	sb.WriteString("---\n")
	meta := map[string]any{
		"type":       string(doc.Type),
		"phase":      string(doc.Phase),
		"title":      doc.Title,
		"synced_at":  time.Now().Format(time.RFC3339),
		"source":     "visionspec",
		"source_dir": e.VisionSpecDir,
	}
	if doc.Metadata != nil {
		for k, v := range doc.Metadata {
			if _, exists := meta[k]; !exists {
				meta[k] = v
			}
		}
	}
	metaBytes, err := yaml.Marshal(meta)
	if err != nil {
		return "", fmt.Errorf("marshal metadata: %w", err)
	}
	sb.Write(metaBytes)
	sb.WriteString("---\n\n")

	// Write content (strip existing frontmatter if present)
	content := doc.Content
	if strings.HasPrefix(content, "---") {
		parts := strings.SplitN(content, "---", 3)
		if len(parts) >= 3 {
			content = strings.TrimSpace(parts[2])
		}
	}
	sb.WriteString(content)

	return sb.String(), nil
}

// convertToVisionSpecFormat converts an AIDLC document to VisionSpec format.
func (e *SyncEngine) convertToVisionSpecFormat(doc *Document) (string, error) {
	// VisionSpec uses markdown with YAML frontmatter
	var sb strings.Builder

	// Write frontmatter
	sb.WriteString("---\n")
	meta := map[string]any{
		"type":       string(doc.Type),
		"title":      doc.Title,
		"synced_at":  time.Now().Format(time.RFC3339),
		"source":     "aidlc",
		"source_dir": e.AIDLCDocsDir,
	}
	if doc.Metadata != nil {
		for k, v := range doc.Metadata {
			if k != "phase" && k != "synced_at" && k != "source" && k != "source_dir" {
				meta[k] = v
			}
		}
	}
	metaBytes, err := yaml.Marshal(meta)
	if err != nil {
		return "", fmt.Errorf("marshal metadata: %w", err)
	}
	sb.Write(metaBytes)
	sb.WriteString("---\n\n")

	// Write content (strip existing frontmatter if present)
	content := doc.Content
	if strings.HasPrefix(content, "---") {
		parts := strings.SplitN(content, "---", 3)
		if len(parts) >= 3 {
			content = strings.TrimSpace(parts[2])
		}
	}
	sb.WriteString(content)

	return sb.String(), nil
}

// inferDocTypeFromName infers document type from a filename.
func inferDocTypeFromName(name string) DocType {
	nameLower := strings.ToLower(name)
	nameLower = strings.ReplaceAll(nameLower, "-", "_")
	nameLower = strings.ReplaceAll(nameLower, " ", "_")

	// Direct match
	for _, dt := range AllDocTypes() {
		if nameLower == string(dt) {
			return dt
		}
	}

	// Partial match
	mappings := map[string]DocType{
		"vision":         DocVisionDocument,
		"requirements":   DocRequirementsSpec,
		"technical":      DocTechnicalSpec,
		"architecture":   DocArchitectureSpec,
		"implementation": DocImplementationPlan,
		"test":           DocTestPlan,
		"integration":    DocIntegrationPlan,
		"security":       DocSecurityReview,
		"runbook":        DocRunbook,
		"monitoring":     DocMonitoringPlan,
		"disaster":       DocDisasterPlan,
		"slo":            DocSLODocument,
	}

	for key, docType := range mappings {
		if strings.Contains(nameLower, key) {
			return docType
		}
	}

	return ""
}

// computeChecksum computes a SHA-256 checksum for content.
func computeChecksum(content []byte) string {
	hash := sha256.Sum256(content)
	return hex.EncodeToString(hash[:])
}

// copyFile copies a file from src to dst.
func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	dest, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, source)
	return err
}
