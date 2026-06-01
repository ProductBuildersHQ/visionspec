package target

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAIDLCTarget_Name(t *testing.T) {
	target := &AIDLCTarget{}
	if got := target.Name(); got != "aidlc" {
		t.Errorf("Name() = %q, want %q", got, "aidlc")
	}
}

func TestAIDLCTarget_Description(t *testing.T) {
	target := &AIDLCTarget{}
	desc := target.Description()
	if !strings.Contains(desc, "AI-DLC") {
		t.Error("Description should mention AI-DLC")
	}
	if !strings.Contains(desc, "vision-document.md") {
		t.Error("Description should mention vision-document.md")
	}
}

func TestAIDLCTarget_Capabilities(t *testing.T) {
	target := &AIDLCTarget{}
	caps := target.Capabilities()

	if !caps.SequentialTasks {
		t.Error("SequentialTasks should be true")
	}
	if !caps.MultiAgent {
		t.Error("MultiAgent should be true")
	}
	if !caps.Verification {
		t.Error("Verification should be true")
	}
}

func TestAIDLCTarget_Validate(t *testing.T) {
	target := &AIDLCTarget{}

	// Empty spec should fail
	if err := target.Validate(""); err == nil {
		t.Error("Validate() should fail for empty spec")
	}

	// Non-empty spec should pass
	if err := target.Validate("# Spec\n\nContent"); err != nil {
		t.Errorf("Validate() failed unexpectedly: %v", err)
	}
}

func TestAIDLCTarget_Export(t *testing.T) {
	target := &AIDLCTarget{}
	tmpDir := t.TempDir()

	spec := `# Reconciled Spec

## Features

- Feature 1: User authentication
- Feature 2: Dashboard display

## Acceptance Criteria

- Users can log in
- Dashboard loads in < 2s
`

	config := ExportConfig{
		ProjectName: "test-project",
		OutputDir:   tmpDir,
		Options: map[string]any{
			"mrd":   "## Problem\n\nUsers need better authentication.",
			"press": "## Summary\n\nAnnouncing secure login.",
			"prd":   "## Goals\n\n- Improve security\n- Better UX",
			"trd":   "## Architecture\n\nMicroservices-based.",
			"ird":   "## Infrastructure\n\nAWS-based deployment.",
		},
	}

	result, err := target.Export(spec, config)
	if err != nil {
		t.Fatalf("Export() failed: %v", err)
	}

	if !result.Success {
		t.Error("Export should succeed")
	}

	if result.Target != "aidlc" {
		t.Errorf("Target = %q, want %q", result.Target, "aidlc")
	}

	// Check files were created
	expectedFiles := []string{
		"vision-document.md",
		"technical-environment.md",
		"imported-requirements.md",
	}

	for _, fname := range expectedFiles {
		fpath := filepath.Join(tmpDir, fname)
		if _, err := os.Stat(fpath); os.IsNotExist(err) {
			t.Errorf("Expected file %s not created", fname)
		}
	}

	// Verify vision-document.md content
	visionContent, err := os.ReadFile(filepath.Join(tmpDir, "vision-document.md"))
	if err != nil {
		t.Fatalf("Failed to read vision-document.md: %v", err)
	}
	if !strings.Contains(string(visionContent), "Vision Document") {
		t.Error("vision-document.md should contain 'Vision Document' heading")
	}
	if !strings.Contains(string(visionContent), "AI-DLC") {
		t.Error("vision-document.md should reference AI-DLC")
	}

	// Verify technical-environment.md content
	techContent, err := os.ReadFile(filepath.Join(tmpDir, "technical-environment.md"))
	if err != nil {
		t.Fatalf("Failed to read technical-environment.md: %v", err)
	}
	if !strings.Contains(string(techContent), "Technical Environment") {
		t.Error("technical-environment.md should contain 'Technical Environment' heading")
	}

	// Verify imported-requirements.md content
	reqContent, err := os.ReadFile(filepath.Join(tmpDir, "imported-requirements.md"))
	if err != nil {
		t.Fatalf("Failed to read imported-requirements.md: %v", err)
	}
	if !strings.Contains(string(reqContent), "test-project") {
		t.Error("imported-requirements.md should contain project name")
	}
	if !strings.Contains(string(reqContent), "Feature 1") {
		t.Error("imported-requirements.md should contain spec content")
	}
}

func TestAIDLCTarget_ExportDefaultDir(t *testing.T) {
	target := &AIDLCTarget{}

	// Use a temp dir as the working directory
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalWd) }()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}

	config := ExportConfig{
		ProjectName: "test-project",
		// OutputDir empty - should use default
		Options: map[string]any{},
	}

	result, err := target.Export("# Spec content", config)
	if err != nil {
		t.Fatalf("Export() failed: %v", err)
	}

	if result.OutputDir != ".aidlc" {
		t.Errorf("OutputDir = %q, want %q", result.OutputDir, ".aidlc")
	}
}

func TestExtractSection(t *testing.T) {
	content := `# Document

## Problem Statement

This is the problem description.
It spans multiple lines.

## Solution

This is the solution.
`

	// Should extract Problem section
	result := extractSection(content, "Problem")
	if !strings.Contains(result, "problem description") {
		t.Errorf("extractSection should find Problem section, got: %s", result)
	}
	if strings.Contains(result, "solution") {
		t.Error("extractSection should not include next section")
	}

	// Should extract Solution section
	result = extractSection(content, "Solution")
	if !strings.Contains(result, "solution") {
		t.Errorf("extractSection should find Solution section, got: %s", result)
	}

	// Should handle missing section
	result = extractSection(content, "NonExistent")
	if !strings.Contains(result, "not found") {
		t.Error("extractSection should indicate content not found")
	}
}

func TestAIDLCTarget_Registered(t *testing.T) {
	// Verify AIDLC target is registered
	target, err := Get("aidlc")
	if err != nil {
		t.Fatalf("Get(aidlc) failed: %v", err)
	}
	if target == nil {
		t.Fatal("Get(aidlc) returned nil")
	}
	if target.Name() != "aidlc" {
		t.Errorf("Got wrong target: %s", target.Name())
	}
}
