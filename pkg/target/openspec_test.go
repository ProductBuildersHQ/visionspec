package target

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestOpenSpecTarget_Interface(t *testing.T) {
	target := &OpenSpecTarget{}

	// Test Name
	if target.Name() != "openspec" {
		t.Errorf("Name() = %s, want openspec", target.Name())
	}

	// Test Description
	if target.Description() == "" {
		t.Error("Description() should not be empty")
	}

	// Test Capabilities
	caps := target.Capabilities()
	if !caps.SequentialTasks {
		t.Error("Capabilities().SequentialTasks should be true")
	}
	if !caps.DependencyGraph {
		t.Error("Capabilities().DependencyGraph should be true")
	}
}

func TestOpenSpecTarget_Validate(t *testing.T) {
	target := &OpenSpecTarget{}

	// Empty spec should fail
	if err := target.Validate(""); err == nil {
		t.Error("Validate should fail for empty spec")
	}

	// Non-empty spec should pass
	if err := target.Validate("# Test Spec"); err != nil {
		t.Errorf("Validate should pass for non-empty spec: %v", err)
	}
}

func TestOpenSpecTarget_Export_JSON(t *testing.T) {
	target := &OpenSpecTarget{}
	tmpDir := t.TempDir()

	spec := `# Test Project Spec

## Problem

Users need a better way to manage tasks.

## Solution

Build a task management system.

## Goals

- Easy to use
- Fast performance

## Features

### Task Management

- Create tasks MUST be supported
- Delete tasks SHOULD be available
- Archive tasks COULD be optional

## Tasks

- [ ] Design API
- [ ] Implement backend
- [x] Write tests

## Acceptance Criteria

- Given a user is logged in
- When they create a task
- Then the task should appear in their list
`

	config := ExportConfig{
		ProjectName: "test-project",
		OutputDir:   tmpDir,
		Options:     map[string]any{"format": "json"},
	}

	result, err := target.Export(spec, config)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	if !result.Success {
		t.Error("Export should succeed")
	}

	if len(result.Files) == 0 {
		t.Error("Export should create at least one file")
	}

	// Verify the output file exists and is valid JSON
	outputFile := filepath.Join(tmpDir, "openspec.json")
	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	var doc OpenSpecDocument
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	// Verify content
	if doc.Metadata.Project != "test-project" {
		t.Errorf("Project = %s, want test-project", doc.Metadata.Project)
	}

	if doc.Version != "1.0.0" {
		t.Errorf("Version = %s, want 1.0.0", doc.Version)
	}

	if len(doc.Overview.Goals) != 2 {
		t.Errorf("Goals count = %d, want 2", len(doc.Overview.Goals))
	}

	// Should have extracted features
	if len(doc.Features) == 0 {
		t.Error("Should have extracted features")
	}

	// Should have extracted tasks
	if len(doc.Tasks) < 3 {
		t.Errorf("Tasks count = %d, want at least 3", len(doc.Tasks))
	}

	// Should have at least one done task
	hasDone := false
	for _, task := range doc.Tasks {
		if task.Status == "done" {
			hasDone = true
			break
		}
	}
	if !hasDone {
		t.Error("Should have at least one done task")
	}
}

func TestOpenSpecTarget_Export_YAML(t *testing.T) {
	target := &OpenSpecTarget{}
	tmpDir := t.TempDir()

	spec := "# Simple Spec\n\nSome content."

	config := ExportConfig{
		ProjectName: "yaml-project",
		OutputDir:   tmpDir,
		Options:     map[string]any{"format": "yaml"},
	}

	result, err := target.Export(spec, config)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	if !result.Success {
		t.Error("Export should succeed")
	}

	// Verify YAML file exists
	outputFile := filepath.Join(tmpDir, "openspec.yaml")
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("YAML output file should exist")
	}
}

func TestOpenSpecTarget_Export_SeparateFiles(t *testing.T) {
	target := &OpenSpecTarget{}
	tmpDir := t.TempDir()

	spec := `# Project

## Features

### Feature One

- Requirement A

## Tasks

- [ ] Task 1
- [ ] Task 2
`

	config := ExportConfig{
		ProjectName: "separate-files-project",
		OutputDir:   tmpDir,
		Options: map[string]any{
			"format":         "json",
			"separate_files": true,
		},
	}

	result, err := target.Export(spec, config)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	if len(result.Files) < 2 {
		t.Errorf("Should have created multiple files, got %d", len(result.Files))
	}

	// Check features directory exists
	featuresDir := filepath.Join(tmpDir, "features")
	if _, err := os.Stat(featuresDir); os.IsNotExist(err) {
		t.Error("Features directory should exist")
	}

	// Check tasks directory exists
	tasksDir := filepath.Join(tmpDir, "tasks")
	if _, err := os.Stat(tasksDir); os.IsNotExist(err) {
		t.Error("Tasks directory should exist")
	}
}

func TestOpenSpecTarget_ExtractBulletPoints(t *testing.T) {
	target := &OpenSpecTarget{}

	content := `
- Point one
- Point two
* Point three
1. Numbered one
2. Numbered two
`

	points := target.extractBulletPoints(content)

	if len(points) != 5 {
		t.Errorf("Expected 5 points, got %d", len(points))
	}
}

func TestOpenSpecTarget_Registration(t *testing.T) {
	// Verify the target is registered
	target, err := Get("openspec")
	if err != nil {
		t.Fatalf("OpenSpec target should be registered: %v", err)
	}

	if target.Name() != "openspec" {
		t.Errorf("Got wrong target: %s", target.Name())
	}
}
