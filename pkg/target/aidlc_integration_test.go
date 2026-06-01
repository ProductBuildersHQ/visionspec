package target

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestAIDLCIntegration tests the full AIDLC export workflow.
func TestAIDLCIntegration(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a mock project structure
	projectDir := filepath.Join(tmpDir, "test-project")
	sourceDir := filepath.Join(projectDir, "source")
	gtmDir := filepath.Join(projectDir, "gtm")
	technicalDir := filepath.Join(projectDir, "technical")

	for _, dir := range []string{sourceDir, gtmDir, technicalDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Create mock spec files
	mrdContent := `# Market Requirements Document

## Problem Statement

Users struggle to manage their tasks efficiently across multiple devices.

## Target Audience

- Busy professionals
- Students
- Project managers

## Business Goals

- Achieve 100K MAU in first year
- 70% retention rate
`

	prdContent := `# Product Requirements Document

## User Stories

### US-001: Create Task
As a user, I want to create tasks quickly so I can capture ideas before I forget them.

### US-002: Sync Tasks
As a user, I want my tasks synced across devices so I can access them anywhere.

## Acceptance Criteria

- Task creation takes < 2 seconds
- Sync happens within 5 seconds of change
`

	pressContent := `# Press Release

## Summary

TaskSync launches revolutionary cross-device task management.

## Lead

Today we announce TaskSync, the task manager that keeps you organized everywhere.
`

	trdContent := `# Technical Requirements Document

## Architecture

Microservices-based architecture with:
- API Gateway
- Task Service
- Sync Service
- User Service

## APIs

### POST /tasks
Create a new task.

### GET /tasks
List all tasks for user.

## Non-Functional Requirements

- 99.9% uptime
- < 100ms API response time
`

	irdContent := `# Infrastructure Requirements Document

## Infrastructure

- AWS EKS for container orchestration
- RDS PostgreSQL for data storage
- ElastiCache Redis for caching

## Deployment

- Blue-green deployment strategy
- Automated rollback on failure
`

	specContent := `# Execution Spec: TaskSync

## Requirements

### Functional Requirements

1. Create tasks with title and description
2. Sync tasks across devices in real-time
3. Mark tasks as complete

### Non-Functional Requirements

1. Response time < 100ms
2. 99.9% availability

## Acceptance Criteria

- All user stories pass verification
- Performance benchmarks met
`

	// Write spec files
	specFiles := map[string]string{
		filepath.Join(sourceDir, "mrd.md"):    mrdContent,
		filepath.Join(sourceDir, "prd.md"):    prdContent,
		filepath.Join(gtmDir, "press.md"):     pressContent,
		filepath.Join(technicalDir, "trd.md"): trdContent,
		filepath.Join(technicalDir, "ird.md"): irdContent,
		filepath.Join(projectDir, "spec.md"):  specContent,
	}

	for path, content := range specFiles {
		if err := os.WriteFile(path, []byte(content), 0600); err != nil {
			t.Fatalf("Failed to write %s: %v", path, err)
		}
	}

	// Create AIDLC export target
	target := &AIDLCTarget{}
	outputDir := filepath.Join(tmpDir, "aidlc-output")

	// Configure export with all spec content
	config := ExportConfig{
		ProjectName: "tasksync",
		OutputDir:   outputDir,
		Options: map[string]any{
			"mrd":     mrdContent,
			"prd":     prdContent,
			"press":   pressContent,
			"trd":     trdContent,
			"ird":     irdContent,
			"context": "Existing Go microservices with gRPC APIs",
		},
	}

	// Export
	result, err := target.Export(specContent, config)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// Verify result
	if !result.Success {
		t.Errorf("Export should succeed")
	}

	if result.Target != "aidlc" {
		t.Errorf("Target = %q, want %q", result.Target, "aidlc")
	}

	if len(result.Files) != 3 {
		t.Errorf("Expected 3 files, got %d", len(result.Files))
	}

	// Verify vision-document.md
	visionPath := filepath.Join(outputDir, "vision-document.md")
	visionContent, err := os.ReadFile(visionPath)
	if err != nil {
		t.Fatalf("Failed to read vision-document.md: %v", err)
	}

	visionStr := string(visionContent)

	// Check structure
	if !strings.Contains(visionStr, "# Vision Document") {
		t.Error("Vision document should have title")
	}
	if !strings.Contains(visionStr, "## Problem/Opportunity") {
		t.Error("Vision document should have Problem/Opportunity section")
	}
	if !strings.Contains(visionStr, "## Vision Statement") {
		t.Error("Vision document should have Vision Statement section")
	}
	if !strings.Contains(visionStr, "## Goals and Success Metrics") {
		t.Error("Vision document should have Goals section")
	}

	// Check content mapping from MRD
	if !strings.Contains(visionStr, "tasks efficiently") {
		t.Error("Vision document should contain MRD problem content")
	}

	// Check content mapping from Press
	if !strings.Contains(visionStr, "TaskSync") {
		t.Error("Vision document should contain Press content")
	}

	// Verify technical-environment.md
	techPath := filepath.Join(outputDir, "technical-environment.md")
	techContent, err := os.ReadFile(techPath)
	if err != nil {
		t.Fatalf("Failed to read technical-environment.md: %v", err)
	}

	techStr := string(techContent)

	if !strings.Contains(techStr, "# Technical Environment") {
		t.Error("Tech environment should have title")
	}
	if !strings.Contains(techStr, "## Existing Systems") {
		t.Error("Tech environment should have Existing Systems section")
	}
	if !strings.Contains(techStr, "gRPC") {
		t.Error("Tech environment should contain context about existing systems")
	}
	if !strings.Contains(techStr, "Microservices") {
		t.Error("Tech environment should contain TRD architecture content")
	}
	if !strings.Contains(techStr, "AWS EKS") {
		t.Error("Tech environment should contain IRD infrastructure content")
	}

	// Verify imported-requirements.md
	reqPath := filepath.Join(outputDir, "imported-requirements.md")
	reqContent, err := os.ReadFile(reqPath)
	if err != nil {
		t.Fatalf("Failed to read imported-requirements.md: %v", err)
	}

	reqStr := string(reqContent)

	if !strings.Contains(reqStr, "# Imported Requirements: tasksync") {
		t.Error("Requirements should have project name in title")
	}
	if !strings.Contains(reqStr, "Create tasks with title") {
		t.Error("Requirements should contain spec.md content")
	}
	if !strings.Contains(reqStr, "AI-DLC Integration") {
		t.Error("Requirements should have AI-DLC integration instructions")
	}
}

// TestAIDLCExportMinimal tests export with minimal input.
func TestAIDLCExportMinimal(t *testing.T) {
	tmpDir := t.TempDir()

	target := &AIDLCTarget{}
	outputDir := filepath.Join(tmpDir, "minimal")

	config := ExportConfig{
		ProjectName: "minimal-project",
		OutputDir:   outputDir,
		Options:     map[string]any{},
	}

	specContent := "# Minimal Spec\n\nJust a simple spec."

	result, err := target.Export(specContent, config)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	if !result.Success {
		t.Error("Export should succeed even with minimal input")
	}

	// All three files should be created
	for _, fname := range []string{"vision-document.md", "technical-environment.md", "imported-requirements.md"} {
		fpath := filepath.Join(outputDir, fname)
		if _, err := os.Stat(fpath); os.IsNotExist(err) {
			t.Errorf("Expected file %s to exist", fname)
		}
	}
}
