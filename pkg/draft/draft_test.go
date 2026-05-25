package draft

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ProductBuildersHQ/visionspec/pkg/types"
)

func TestDraftPath(t *testing.T) {
	projectPath := filepath.Join("test", "project")

	tests := []struct {
		specType types.SpecType
		want     string
	}{
		{types.SpecTypeMRD, filepath.Join("test", "project", "source", "mrd.draft.md")},
		{types.SpecTypePRD, filepath.Join("test", "project", "source", "prd.draft.md")},
		{types.SpecTypeUXD, filepath.Join("test", "project", "source", "uxd.draft.md")},
	}

	for _, tt := range tests {
		t.Run(string(tt.specType), func(t *testing.T) {
			got := DraftPath(projectPath, tt.specType)
			if got != tt.want {
				t.Errorf("DraftPath() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMetadataPath(t *testing.T) {
	projectPath := filepath.Join("test", "project")

	got := MetadataPath(projectPath, types.SpecTypeMRD)
	want := filepath.Join("test", "project", "source", "mrd.draft.json")
	if got != want {
		t.Errorf("MetadataPath() = %q, want %q", got, want)
	}
}

func TestStartAndGet(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "test-project")
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatal(err)
	}

	// Start should succeed
	draft, err := Start(projectPath, types.SpecTypeMRD)
	if err != nil {
		t.Fatalf("Start() failed: %v", err)
	}

	if draft.Metadata.SpecType != types.SpecTypeMRD {
		t.Errorf("SpecType = %s, want MRD", draft.Metadata.SpecType)
	}
	if draft.Metadata.Version != 1 {
		t.Errorf("Version = %d, want 1", draft.Metadata.Version)
	}
	if draft.Content == "" {
		t.Error("Content is empty")
	}

	// Draft should exist
	if !Exists(projectPath, types.SpecTypeMRD) {
		t.Error("Exists() returned false after Start()")
	}

	// Get should return the same draft
	retrieved, err := Get(projectPath, types.SpecTypeMRD)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	if retrieved.Content != draft.Content {
		t.Error("Retrieved content doesn't match original")
	}
}

func TestStartDuplicateError(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "test-project")
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatal(err)
	}

	// Start first draft
	_, err := Start(projectPath, types.SpecTypeMRD)
	if err != nil {
		t.Fatal(err)
	}

	// Start should fail with ErrDraftExists
	_, err = Start(projectPath, types.SpecTypeMRD)
	if err != ErrDraftExists {
		t.Errorf("Start() error = %v, want ErrDraftExists", err)
	}
}

func TestUpdate(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "test-project")
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatal(err)
	}

	// Start draft
	draft, err := Start(projectPath, types.SpecTypePRD)
	if err != nil {
		t.Fatal(err)
	}

	// Update content
	newContent := "# Updated PRD\n\nNew content here."
	updated, err := Update(projectPath, types.SpecTypePRD, newContent)
	if err != nil {
		t.Fatalf("Update() failed: %v", err)
	}

	if updated.Content != newContent {
		t.Error("Updated content doesn't match")
	}
	if updated.Metadata.Version != draft.Metadata.Version+1 {
		t.Errorf("Version = %d, want %d", updated.Metadata.Version, draft.Metadata.Version+1)
	}
}

func TestDiscard(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "test-project")
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatal(err)
	}

	// Start draft
	_, err := Start(projectPath, types.SpecTypeUXD)
	if err != nil {
		t.Fatal(err)
	}

	// Discard
	err = Discard(projectPath, types.SpecTypeUXD)
	if err != nil {
		t.Fatalf("Discard() failed: %v", err)
	}

	// Should not exist
	if Exists(projectPath, types.SpecTypeUXD) {
		t.Error("Draft still exists after Discard()")
	}
}

func TestDiscardNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "test-project")
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatal(err)
	}

	err := Discard(projectPath, types.SpecTypeMRD)
	if err != ErrDraftNotFound {
		t.Errorf("Discard() error = %v, want ErrDraftNotFound", err)
	}
}

func TestFinalize(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "test-project")
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatal(err)
	}

	// Start and update draft
	_, err := Start(projectPath, types.SpecTypeMRD)
	if err != nil {
		t.Fatal(err)
	}

	content := "# Final MRD Content"
	_, err = Update(projectPath, types.SpecTypeMRD, content)
	if err != nil {
		t.Fatal(err)
	}

	// Finalize
	err = Finalize(projectPath, types.SpecTypeMRD)
	if err != nil {
		t.Fatalf("Finalize() failed: %v", err)
	}

	// Draft should be gone
	if Exists(projectPath, types.SpecTypeMRD) {
		t.Error("Draft still exists after Finalize()")
	}

	// Spec should exist at final location
	specPath := filepath.Join(projectPath, "source", "mrd.md")
	data, err := os.ReadFile(specPath)
	if err != nil {
		t.Fatalf("Failed to read finalized spec: %v", err)
	}

	if string(data) != content {
		t.Error("Finalized spec content doesn't match")
	}
}

func TestAddEvalResult(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "test-project")
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatal(err)
	}

	// Start draft
	_, err := Start(projectPath, types.SpecTypePRD)
	if err != nil {
		t.Fatal(err)
	}

	// Add eval result
	err = AddEvalResult(projectPath, types.SpecTypePRD, 7.5, true, 3)
	if err != nil {
		t.Fatalf("AddEvalResult() failed: %v", err)
	}

	// Verify
	draft, err := Get(projectPath, types.SpecTypePRD)
	if err != nil {
		t.Fatal(err)
	}

	if len(draft.Metadata.EvalHistory) != 1 {
		t.Fatalf("EvalHistory length = %d, want 1", len(draft.Metadata.EvalHistory))
	}

	eval := draft.Metadata.EvalHistory[0]
	if eval.Score != 7.5 {
		t.Errorf("Score = %f, want 7.5", eval.Score)
	}
	if !eval.Passed {
		t.Error("Passed = false, want true")
	}
	if eval.Findings != 3 {
		t.Errorf("Findings = %d, want 3", eval.Findings)
	}
}

func TestListDrafts(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "test-project")
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatal(err)
	}

	// Start multiple drafts
	_, err := Start(projectPath, types.SpecTypeMRD)
	if err != nil {
		t.Fatal(err)
	}
	_, err = Start(projectPath, types.SpecTypePRD)
	if err != nil {
		t.Fatal(err)
	}

	// List drafts
	drafts, err := ListDrafts(projectPath)
	if err != nil {
		t.Fatalf("ListDrafts() failed: %v", err)
	}

	if len(drafts) != 2 {
		t.Errorf("ListDrafts() returned %d drafts, want 2", len(drafts))
	}
}
