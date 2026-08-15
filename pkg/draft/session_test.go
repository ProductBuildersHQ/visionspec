package draft

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ProductBuildersHQ/visionspec/pkg/types"
)

func TestStartSession(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "test-project")
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatal(err)
	}

	// Start new session
	session, err := StartSession(projectPath, types.SpecTypeMRD)
	if err != nil {
		t.Fatalf("StartSession() failed: %v", err)
	}

	if !session.IsNew() {
		t.Error("IsNew() returned false for new session")
	}
	if session.Version() != 1 {
		t.Errorf("Version() = %d, want 1", session.Version())
	}
	if session.Status() != StatusNew {
		t.Errorf("Status() = %s, want %s", session.Status(), StatusNew)
	}
}

func TestStartSessionResumes(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "test-project")
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatal(err)
	}

	// Start first session
	session1, err := StartSession(projectPath, types.SpecTypePRD)
	if err != nil {
		t.Fatal(err)
	}

	// Update content
	err = session1.Update("# Updated content")
	if err != nil {
		t.Fatal(err)
	}

	// Start again should resume
	session2, err := StartSession(projectPath, types.SpecTypePRD)
	if err != nil {
		t.Fatalf("StartSession() resume failed: %v", err)
	}

	if session2.Version() != 2 {
		t.Errorf("Resumed session Version() = %d, want 2", session2.Version())
	}
	if session2.Content() != "# Updated content" {
		t.Error("Resumed session has wrong content")
	}
}

func TestResumeSession(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "test-project")
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatal(err)
	}

	// Start session
	_, err := StartSession(projectPath, types.SpecTypeUXD)
	if err != nil {
		t.Fatal(err)
	}

	// Resume
	session, err := ResumeSession(projectPath, types.SpecTypeUXD)
	if err != nil {
		t.Fatalf("ResumeSession() failed: %v", err)
	}

	if session.SpecType != types.SpecTypeUXD {
		t.Errorf("SpecType = %s, want UXD", session.SpecType)
	}
}

func TestResumeSessionNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "test-project")
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatal(err)
	}

	_, err := ResumeSession(projectPath, types.SpecTypeMRD)
	if err != ErrDraftNotFound {
		t.Errorf("ResumeSession() error = %v, want ErrDraftNotFound", err)
	}
}

func TestSessionUpdate(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "test-project")
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatal(err)
	}

	session, err := StartSession(projectPath, types.SpecTypeMRD)
	if err != nil {
		t.Fatal(err)
	}

	err = session.Update("# New content")
	if err != nil {
		t.Fatalf("Update() failed: %v", err)
	}

	if session.Content() != "# New content" {
		t.Error("Content not updated")
	}
	if session.Version() != 2 {
		t.Errorf("Version() = %d, want 2", session.Version())
	}
}

func TestSessionFinalize(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "test-project")
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatal(err)
	}

	session, err := StartSession(projectPath, types.SpecTypeMRD)
	if err != nil {
		t.Fatal(err)
	}

	err = session.Finalize()
	if err != nil {
		t.Fatalf("Finalize() failed: %v", err)
	}

	// Draft should be gone
	if Exists(projectPath, types.SpecTypeMRD) {
		t.Error("Draft still exists after Finalize()")
	}
}

func TestSessionDiscard(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "test-project")
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatal(err)
	}

	session, err := StartSession(projectPath, types.SpecTypePRD)
	if err != nil {
		t.Fatal(err)
	}

	err = session.Discard()
	if err != nil {
		t.Fatalf("Discard() failed: %v", err)
	}

	if Exists(projectPath, types.SpecTypePRD) {
		t.Error("Draft still exists after Discard()")
	}
}

func TestSessionRecordEval(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "test-project")
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatal(err)
	}

	session, err := StartSession(projectPath, types.SpecTypeUXD)
	if err != nil {
		t.Fatal(err)
	}

	// Record failing eval
	err = session.RecordEval(6.5, false, 5)
	if err != nil {
		t.Fatalf("RecordEval() failed: %v", err)
	}

	if session.EvalCount() != 1 {
		t.Errorf("EvalCount() = %d, want 1", session.EvalCount())
	}
	if session.Status() != StatusFailing {
		t.Errorf("Status() = %s, want %s", session.Status(), StatusFailing)
	}
	if session.HasPassingEval() {
		t.Error("HasPassingEval() returned true")
	}

	// Record passing eval
	err = session.RecordEval(8.0, true, 2)
	if err != nil {
		t.Fatal(err)
	}

	if session.EvalCount() != 2 {
		t.Errorf("EvalCount() = %d, want 2", session.EvalCount())
	}
	if session.Status() != StatusPassing {
		t.Errorf("Status() = %s, want %s", session.Status(), StatusPassing)
	}
	if !session.HasPassingEval() {
		t.Error("HasPassingEval() returned false")
	}
}

func TestSessionLastEval(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "test-project")
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatal(err)
	}

	session, err := StartSession(projectPath, types.SpecTypeMRD)
	if err != nil {
		t.Fatal(err)
	}

	// No evals yet
	if session.LastEval() != nil {
		t.Error("LastEval() returned non-nil for no evals")
	}

	// Add eval
	if err := session.RecordEval(7.0, true, 1); err != nil {
		t.Fatal(err)
	}

	lastEval := session.LastEval()
	if lastEval == nil {
		t.Fatal("LastEval() returned nil after recording eval")
	}
	if lastEval.Score != 7.0 {
		t.Errorf("LastEval().Score = %f, want 7.0", lastEval.Score)
	}
}

func TestSessionStatus(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "test-project")
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatal(err)
	}

	session, err := StartSession(projectPath, types.SpecTypePRD)
	if err != nil {
		t.Fatal(err)
	}

	// New
	if session.Status() != StatusNew {
		t.Errorf("Status() = %s, want %s", session.Status(), StatusNew)
	}

	// Editing
	if err := session.Update("# Updated"); err != nil {
		t.Fatal(err)
	}
	if session.Status() != StatusEditing {
		t.Errorf("Status() = %s, want %s", session.Status(), StatusEditing)
	}

	// Failing
	if err := session.RecordEval(5.0, false, 10); err != nil {
		t.Fatal(err)
	}
	if session.Status() != StatusFailing {
		t.Errorf("Status() = %s, want %s", session.Status(), StatusFailing)
	}

	// Passing
	if err := session.RecordEval(8.5, true, 1); err != nil {
		t.Fatal(err)
	}
	if session.Status() != StatusPassing {
		t.Errorf("Status() = %s, want %s", session.Status(), StatusPassing)
	}
}

func TestSessionSummary(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "test-project")
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatal(err)
	}

	session, err := StartSession(projectPath, types.SpecTypeMRD)
	if err != nil {
		t.Fatal(err)
	}

	summary := session.Summary()

	if summary.SpecType != types.SpecTypeMRD {
		t.Errorf("SpecType = %s, want MRD", summary.SpecType)
	}
	if summary.SpecName != "Market Requirements Document" {
		t.Errorf("SpecName = %s, want 'Market Requirements Document'", summary.SpecName)
	}
	if summary.Version != 1 {
		t.Errorf("Version = %d, want 1", summary.Version)
	}
	if summary.Status != StatusNew {
		t.Errorf("Status = %s, want %s", summary.Status, StatusNew)
	}
	if !summary.HasTemplate {
		t.Error("HasTemplate = false, want true")
	}
	if summary.ContentLength == 0 {
		t.Error("ContentLength = 0, want > 0")
	}
}

func TestSessionInstructions(t *testing.T) {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "test-project")
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatal(err)
	}

	session, err := StartSession(projectPath, types.SpecTypeMRD)
	if err != nil {
		t.Fatal(err)
	}

	// Should have instructions for new status
	instructions := session.Instructions()
	if instructions == "" {
		t.Error("Instructions() returned empty string")
	}
}
