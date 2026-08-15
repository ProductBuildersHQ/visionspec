package version

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ProductBuildersHQ/visionspec/pkg/types"
)

func setupTestProject(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "test-project")

	// Create directories
	dirs := []string{
		filepath.Join(projectPath, "source"),
		filepath.Join(projectPath, "eval"),
		filepath.Join(projectPath, "eval", "versions"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
	}

	return projectPath
}

func writeSpec(t *testing.T, projectPath string, specType types.SpecType, content string) {
	t.Helper()
	specPath := filepath.Join(projectPath, "source", string(specType)+".md")
	if err := os.WriteFile(specPath, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestCreateVersion(t *testing.T) {
	projectPath := setupTestProject(t)

	// Write initial spec
	writeSpec(t, projectPath, types.SpecTypeMRD, "# MRD v1\n\nInitial content.")

	// Create first version
	v1, err := CreateVersion(projectPath, types.SpecTypeMRD, CreateOptions{
		Author:  "test",
		Message: "Initial version",
	})
	if err != nil {
		t.Fatalf("CreateVersion() failed: %v", err)
	}

	if v1.Number != 1 {
		t.Errorf("Version.Number = %d, want 1", v1.Number)
	}
	if v1.Author != "test" {
		t.Errorf("Version.Author = %s, want 'test'", v1.Author)
	}
	if v1.Hash == "" {
		t.Error("Version.Hash is empty")
	}

	// Update spec and create second version
	writeSpec(t, projectPath, types.SpecTypeMRD, "# MRD v2\n\nUpdated content.")

	v2, err := CreateVersion(projectPath, types.SpecTypeMRD, CreateOptions{
		Message: "Second version",
	})
	if err != nil {
		t.Fatalf("CreateVersion() second failed: %v", err)
	}

	if v2.Number != 2 {
		t.Errorf("Version.Number = %d, want 2", v2.Number)
	}
	if v2.Hash == v1.Hash {
		t.Error("Hash should be different for different content")
	}
}

func TestCreateVersionNoChanges(t *testing.T) {
	projectPath := setupTestProject(t)

	content := "# Unchanged MRD\n\nSame content."
	writeSpec(t, projectPath, types.SpecTypeMRD, content)

	// Create first version
	_, err := CreateVersion(projectPath, types.SpecTypeMRD, CreateOptions{})
	if err != nil {
		t.Fatal(err)
	}

	// Try to create second version without changes
	_, err = CreateVersion(projectPath, types.SpecTypeMRD, CreateOptions{})
	if err != ErrNoChanges {
		t.Errorf("CreateVersion() error = %v, want ErrNoChanges", err)
	}
}

func TestGetVersion(t *testing.T) {
	projectPath := setupTestProject(t)

	content := "# Test MRD\n\nVersion content."
	writeSpec(t, projectPath, types.SpecTypeMRD, content)

	// Create version
	_, err := CreateVersion(projectPath, types.SpecTypeMRD, CreateOptions{
		Message: "Test version",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Get version
	v, retrievedContent, err := GetVersion(projectPath, types.SpecTypeMRD, 1)
	if err != nil {
		t.Fatalf("GetVersion() failed: %v", err)
	}

	if v.Number != 1 {
		t.Errorf("Version.Number = %d, want 1", v.Number)
	}
	if retrievedContent != content {
		t.Errorf("Content mismatch: got %q, want %q", retrievedContent, content)
	}
}

func TestGetVersionNotFound(t *testing.T) {
	projectPath := setupTestProject(t)

	_, _, err := GetVersion(projectPath, types.SpecTypeMRD, 999)
	if err != ErrVersionNotFound {
		t.Errorf("GetVersion() error = %v, want ErrVersionNotFound", err)
	}
}

func TestListVersions(t *testing.T) {
	projectPath := setupTestProject(t)

	// Create multiple versions
	for i := 1; i <= 3; i++ {
		writeSpec(t, projectPath, types.SpecTypePRD, "# PRD\n\nVersion "+string(rune('0'+i)))
		_, err := CreateVersion(projectPath, types.SpecTypePRD, CreateOptions{})
		if err != nil {
			t.Fatal(err)
		}
	}

	versions, err := ListVersions(projectPath, types.SpecTypePRD)
	if err != nil {
		t.Fatalf("ListVersions() failed: %v", err)
	}

	if len(versions) != 3 {
		t.Errorf("ListVersions() returned %d versions, want 3", len(versions))
	}

	// Should be in reverse chronological order
	if versions[0].Number != 3 {
		t.Errorf("First version should be 3, got %d", versions[0].Number)
	}
}

func TestRevert(t *testing.T) {
	projectPath := setupTestProject(t)

	// Create v1
	v1Content := "# MRD v1\n\nOriginal content."
	writeSpec(t, projectPath, types.SpecTypeMRD, v1Content)
	_, err := CreateVersion(projectPath, types.SpecTypeMRD, CreateOptions{Message: "v1"})
	if err != nil {
		t.Fatal(err)
	}

	// Create v2
	v2Content := "# MRD v2\n\nModified content."
	writeSpec(t, projectPath, types.SpecTypeMRD, v2Content)
	_, err = CreateVersion(projectPath, types.SpecTypeMRD, CreateOptions{Message: "v2"})
	if err != nil {
		t.Fatal(err)
	}

	// Revert to v1
	revertV, err := Revert(projectPath, types.SpecTypeMRD, 1, "")
	if err != nil {
		t.Fatalf("Revert() failed: %v", err)
	}

	// Should create v3 as the revert
	if revertV.Number != 3 {
		t.Errorf("Revert created version %d, want 3", revertV.Number)
	}

	// Current spec should have v1 content
	specPath := filepath.Join(projectPath, "source", "mrd.md")
	data, err := os.ReadFile(specPath)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != v1Content {
		t.Errorf("Spec content = %q, want %q", string(data), v1Content)
	}
}

func TestGetHistory(t *testing.T) {
	projectPath := setupTestProject(t)

	// No history initially
	history, err := GetHistory(projectPath, types.SpecTypeUXD)
	if err != nil {
		t.Fatalf("GetHistory() failed: %v", err)
	}

	if len(history.Versions) != 0 {
		t.Errorf("Initial history should be empty, got %d versions", len(history.Versions))
	}

	// Create versions
	writeSpec(t, projectPath, types.SpecTypeUXD, "# UXD v1")
	_, err = CreateVersion(projectPath, types.SpecTypeUXD, CreateOptions{})
	if err != nil {
		t.Fatal(err)
	}

	history, err = GetHistory(projectPath, types.SpecTypeUXD)
	if err != nil {
		t.Fatal(err)
	}

	if len(history.Versions) != 1 {
		t.Errorf("History should have 1 version, got %d", len(history.Versions))
	}
}
