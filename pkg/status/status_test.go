package status

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ProductBuildersHQ/visionspec/pkg/types"
)

func TestGenerate(t *testing.T) {
	// Create temp directory with project structure
	tmpDir := t.TempDir()

	// Create subdirectories
	dirs := []string{"source", "gtm", "technical", "eval"}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755); err != nil {
			t.Fatalf("failed to create dir %s: %v", dir, err)
		}
	}

	project := &types.Project{
		Name: "test-project",
		Path: tmpDir,
	}

	report, err := Generate(project)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if report.Project != "test-project" {
		t.Errorf("Project = %q, want %q", report.Project, "test-project")
	}

	if report.Path != tmpDir {
		t.Errorf("Path = %q, want %q", report.Path, tmpDir)
	}

	if report.Summary.TotalSpecs != 11 {
		t.Errorf("TotalSpecs = %d, want 11", report.Summary.TotalSpecs)
	}

	// All specs should be missing
	if report.Summary.PresentSpecs != 0 {
		t.Errorf("PresentSpecs = %d, want 0", report.Summary.PresentSpecs)
	}

	// Should not be ready
	if report.Readiness.Ready {
		t.Error("expected project to not be ready")
	}
}

func TestGenerateWithSpecs(t *testing.T) {
	tmpDir := t.TempDir()

	// Create subdirectories
	dirs := []string{"source", "gtm", "technical", "eval"}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755); err != nil {
			t.Fatalf("failed to create dir %s: %v", dir, err)
		}
	}

	// Create some spec files
	specs := []string{
		filepath.Join(tmpDir, "source", "mrd.md"),
		filepath.Join(tmpDir, "source", "prd.md"),
	}
	for _, spec := range specs {
		if err := os.WriteFile(spec, []byte("# Test\n"), 0600); err != nil {
			t.Fatalf("failed to create spec %s: %v", spec, err)
		}
	}

	project := &types.Project{
		Name: "test-project",
		Path: tmpDir,
	}

	report, err := Generate(project)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if report.Summary.PresentSpecs != 2 {
		t.Errorf("PresentSpecs = %d, want 2", report.Summary.PresentSpecs)
	}
}

func TestRenderText(t *testing.T) {
	report := &Report{
		Project: "test-project",
		Path:    "/path/to/project",
		Readiness: types.ReadinessStatus{
			Ready:   false,
			Summary: "Not ready: 1 blocker",
			Gates: []types.ReadinessGate{
				{Name: "Test gate", Passed: true, Message: "Passed"},
				{Name: "Failed gate", Passed: false, Message: "Failed"},
			},
		},
		Specs: []SpecStatus{
			{Type: types.SpecTypeMRD, Category: types.CategorySource, Exists: true, Required: true, Status: types.StatusDraft},
		},
		Summary: Summary{TotalSpecs: 1, PresentSpecs: 1},
	}

	var buf bytes.Buffer
	if err := RenderText(&buf, report); err != nil {
		t.Fatalf("RenderText failed: %v", err)
	}

	output := buf.String()

	// Check key content
	if !strings.Contains(output, "test-project") {
		t.Error("output should contain project name")
	}
	if !strings.Contains(output, "NOT READY") {
		t.Error("output should contain NOT READY")
	}
	if !strings.Contains(output, "[+] Test gate") {
		t.Error("output should contain passed gate")
	}
	if !strings.Contains(output, "[X] Failed gate") {
		t.Error("output should contain failed gate")
	}
}

func TestRenderHTML(t *testing.T) {
	report := &Report{
		Project: "test-project",
		Path:    "/path/to/project",
		Readiness: types.ReadinessStatus{
			Ready:   true,
			Summary: "Ready for AI-assisted development",
			Gates: []types.ReadinessGate{
				{Name: "Test gate", Passed: true, Message: "Passed"},
			},
		},
		Specs:   []SpecStatus{},
		Summary: Summary{TotalSpecs: 0},
	}

	var buf bytes.Buffer
	if err := RenderHTML(&buf, report); err != nil {
		t.Fatalf("RenderHTML failed: %v", err)
	}

	output := buf.String()

	// Check key content
	if !strings.Contains(output, "<!DOCTYPE html>") {
		t.Error("output should be valid HTML")
	}
	if !strings.Contains(output, "test-project") {
		t.Error("output should contain project name")
	}
	if !strings.Contains(output, "#28a745") {
		t.Error("output should contain green color for READY status")
	}
}

func TestRenderMarkdown(t *testing.T) {
	report := &Report{
		Project: "test-project",
		Path:    "/path/to/project",
		Readiness: types.ReadinessStatus{
			Ready:   false,
			Summary: "Not ready: 2 blockers",
			Gates: []types.ReadinessGate{
				{Name: "Gate 1", Passed: true, Message: "OK"},
				{Name: "Gate 2", Passed: false, Message: "Failed"},
			},
		},
		Specs: []SpecStatus{
			{Type: types.SpecTypeMRD, Category: types.CategorySource, Exists: true, Required: true},
			{Type: types.SpecTypePRD, Category: types.CategorySource, Exists: false, Required: true},
		},
		Summary: Summary{TotalSpecs: 2, PresentSpecs: 1},
	}

	var buf bytes.Buffer
	if err := RenderMarkdown(&buf, report); err != nil {
		t.Fatalf("RenderMarkdown failed: %v", err)
	}

	output := buf.String()

	// Check markdown structure
	if !strings.Contains(output, "# Project Status:") {
		t.Error("output should contain markdown header")
	}
	if !strings.Contains(output, "| Type | Category |") {
		t.Error("output should contain markdown table")
	}
	if !strings.Contains(output, ":white_check_mark:") {
		t.Error("output should contain checkmark emoji")
	}
	if !strings.Contains(output, ":x:") {
		t.Error("output should contain x emoji")
	}
}

func TestCalculateReadiness(t *testing.T) {
	// Test with all gates passing
	report := &Report{
		Path: t.TempDir(),
		Specs: []SpecStatus{
			{Type: types.SpecTypeMRD, Required: true, Exists: true, Approval: &types.Approval{}},
			{Type: types.SpecTypePRD, Required: true, Exists: true, Approval: &types.Approval{}},
			{Type: types.SpecTypeUXD, Required: true, Exists: true, Approval: &types.Approval{}},
			{Type: types.SpecTypeTRD, Required: true, Exists: true, Approval: &types.Approval{}},
		},
	}

	// Create spec.md to pass the last gate
	specPath := filepath.Join(report.Path, "spec.md")
	if err := os.WriteFile(specPath, []byte("# Spec\n"), 0600); err != nil {
		t.Fatalf("failed to create spec.md: %v", err)
	}

	status := calculateReadiness(report)

	if !status.Ready {
		t.Errorf("expected Ready to be true, got false. Summary: %s", status.Summary)
	}

	if len(status.Gates) != 4 {
		t.Errorf("expected 4 gates, got %d", len(status.Gates))
	}

	for _, gate := range status.Gates {
		if !gate.Passed {
			t.Errorf("expected gate %q to pass", gate.Name)
		}
	}
}

func TestPluralize(t *testing.T) {
	tests := []struct {
		n        int
		singular string
		plural   string
		want     string
	}{
		{0, "item", "items", "0 items"},
		{1, "item", "items", "1 item"},
		{2, "item", "items", "2 items"},
		{10, "blocker", "blockers", "10 blockers"},
	}

	for _, tt := range tests {
		got := pluralize(tt.n, tt.singular, tt.plural)
		if got != tt.want {
			t.Errorf("pluralize(%d, %q, %q) = %q, want %q", tt.n, tt.singular, tt.plural, got, tt.want)
		}
	}
}
