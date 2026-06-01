package rules

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestList(t *testing.T) {
	files, err := List()
	if err != nil {
		t.Fatalf("List() failed: %v", err)
	}

	if len(files) == 0 {
		t.Error("List() returned no files")
	}

	// Should contain core workflow
	found := false
	for _, f := range files {
		if strings.Contains(f, "core-workflow.md") {
			found = true
			break
		}
	}
	if !found {
		t.Error("List() should include core-workflow.md")
	}
}

func TestGet(t *testing.T) {
	content, err := Get("core-workflow.md")
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	if len(content) == 0 {
		t.Error("Get() returned empty content")
	}

	if !strings.Contains(string(content), "VisionSpec Core Workflow") {
		t.Error("core-workflow.md should contain expected header")
	}
}

func TestGetNotFound(t *testing.T) {
	_, err := Get("nonexistent.md")
	if err == nil {
		t.Error("Get() should fail for nonexistent file")
	}
}

func TestExport(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "rules")

	files, err := Export(outputDir)
	if err != nil {
		t.Fatalf("Export() failed: %v", err)
	}

	if len(files) == 0 {
		t.Error("Export() should create files")
	}

	// Verify core-workflow.md was created
	coreWorkflow := filepath.Join(outputDir, "core-workflow.md")
	if _, err := os.Stat(coreWorkflow); os.IsNotExist(err) {
		t.Error("Export() should create core-workflow.md")
	}

	// Verify phases directory was created
	phasesDir := filepath.Join(outputDir, "phases")
	if _, err := os.Stat(phasesDir); os.IsNotExist(err) {
		t.Error("Export() should create phases directory")
	}

	// Verify gates directory was created
	gatesDir := filepath.Join(outputDir, "gates")
	if _, err := os.Stat(gatesDir); os.IsNotExist(err) {
		t.Error("Export() should create gates directory")
	}

	// Verify frameworks directory was created
	frameworksDir := filepath.Join(outputDir, "frameworks")
	if _, err := os.Stat(frameworksDir); os.IsNotExist(err) {
		t.Error("Export() should create frameworks directory")
	}
}

func TestExportDefaultDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Change to temp directory
	originalWd, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalWd) }()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Export with empty dir (should use default)
	files, err := Export("")
	if err != nil {
		t.Fatalf("Export() failed: %v", err)
	}

	if len(files) == 0 {
		t.Error("Export() should create files")
	}

	// Verify default directory was created
	defaultDir := filepath.Join(tmpDir, ".visionspec-rules")
	if _, err := os.Stat(defaultDir); os.IsNotExist(err) {
		t.Error("Export() should create .visionspec-rules directory by default")
	}
}
