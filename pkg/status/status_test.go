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

	if report.Summary.TotalSpecs != 12 {
		t.Errorf("TotalSpecs = %d, want 12", report.Summary.TotalSpecs)
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

func TestNewRichReport(t *testing.T) {
	tmpDir := t.TempDir()

	report := &Report{
		Project: "test-project",
		Path:    tmpDir,
		Specs: []SpecStatus{
			{
				Type:     types.SpecTypeMRD,
				Category: types.CategorySource,
				Exists:   true,
				Required: true,
				EvalStatus: &EvalStatus{
					Exists:   true,
					Decision: "pass",
					Categories: &CategoryBreakdown{
						Pass:  6,
						Total: 6,
					},
				},
			},
			{
				Type:     types.SpecTypePRD,
				Category: types.CategorySource,
				Exists:   true,
				Required: true,
				EvalStatus: &EvalStatus{
					Exists:   true,
					Decision: "pass",
					Categories: &CategoryBreakdown{
						Pass:    5,
						Partial: 1,
						Total:   6,
					},
				},
			},
			{
				Type:     types.SpecTypeIRD,
				Category: types.CategoryTechnical,
				Exists:   false,
				Required: true,
			},
		},
		Summary: Summary{TotalSpecs: 3, PresentSpecs: 2, EvaluatedSpecs: 2},
	}

	rr := NewRichReport(report)

	// Check pipeline was built
	if len(rr.Pipeline) == 0 {
		t.Error("expected pipeline to be built")
	}

	// Check aggregate categories
	if rr.AggregateCategories.Pass != 11 {
		t.Errorf("AggregateCategories.Pass = %d, want 11", rr.AggregateCategories.Pass)
	}
	if rr.AggregateCategories.Partial != 1 {
		t.Errorf("AggregateCategories.Partial = %d, want 1", rr.AggregateCategories.Partial)
	}
	if rr.AggregateCategories.Total != 12 {
		t.Errorf("AggregateCategories.Total = %d, want 12", rr.AggregateCategories.Total)
	}

	// Check completion percent (2 complete out of 9 pipeline stages = 22%)
	// MRD and PRD are complete with passing evals
	if rr.CompletionPercent < 20 || rr.CompletionPercent > 25 {
		t.Errorf("CompletionPercent = %d, want ~22", rr.CompletionPercent)
	}
}

func TestRichReportSetters(t *testing.T) {
	report := &Report{
		Project: "test-project",
		Path:    t.TempDir(),
		Specs:   []SpecStatus{},
	}

	rr := NewRichReport(report)

	// Test SetHighlight
	rr.SetHighlight(types.SpecTypeMRD, "Market analysis")
	if rr.Highlights[types.SpecTypeMRD] != "Market analysis" {
		t.Error("SetHighlight failed")
	}

	// Test SetMediumFindings
	rr.SetMediumFindings([]Finding{
		{Spec: types.SpecTypeTPD, Severity: "medium", Description: "Missing rollback testing"},
	})
	if len(rr.MediumFindings) != 1 {
		t.Error("SetMediumFindings failed")
	}

	// Test SetNextSteps
	rr.SetNextSteps([]string{"Create IRD", "Create spec.md"})
	if len(rr.NextSteps) != 2 {
		t.Error("SetNextSteps failed")
	}

	// Test SetKeyDecisions
	rr.SetKeyDecisions([]KeyDecision{
		{Area: "Cloud", Choice: "AWS with Pulumi"},
	})
	if len(rr.KeyDecisions) != 1 {
		t.Error("SetKeyDecisions failed")
	}

	// Test SetReadySummary
	rr.SetReadySummary("Reconciliation and implementation")
	if rr.ReadySummary != "Reconciliation and implementation" {
		t.Error("SetReadySummary failed")
	}
}

func TestRenderRichText(t *testing.T) {
	tmpDir := t.TempDir()

	report := &Report{
		Project: "test-project",
		Path:    tmpDir,
		Specs: []SpecStatus{
			{
				Type:     types.SpecTypeMRD,
				Category: types.CategorySource,
				Exists:   true,
				Required: true,
				EvalStatus: &EvalStatus{
					Exists:   true,
					Decision: "pass",
					Categories: &CategoryBreakdown{
						Pass:  6,
						Total: 6,
					},
					Findings: struct {
						Critical int `json:"critical"`
						High     int `json:"high"`
						Medium   int `json:"medium"`
						Low      int `json:"low"`
						Info     int `json:"info"`
					}{Info: 2},
				},
			},
		},
		Summary: Summary{TotalSpecs: 1, PresentSpecs: 1, EvaluatedSpecs: 1},
	}

	rr := NewRichReport(report)
	rr.SetHighlight(types.SpecTypeMRD, "Market analysis, competitive positioning")
	rr.SetMediumFindings([]Finding{
		{Spec: types.SpecTypeTPD, Severity: "medium", Description: "Missing rollback testing"},
	})
	rr.SetNextSteps([]string{"Create IRD", "Create spec.md"})
	rr.SetReadySummary("Reconciliation and implementation")

	var buf bytes.Buffer
	if err := RenderRichText(&buf, rr); err != nil {
		t.Fatalf("RenderRichText failed: %v", err)
	}

	output := buf.String()

	// Check key content
	if !strings.Contains(output, "VisionSpec Status") {
		t.Error("output should contain VisionSpec Status")
	}
	if !strings.Contains(output, "Pipeline Progress") {
		t.Error("output should contain Pipeline Progress")
	}
	if !strings.Contains(output, "Summary") {
		t.Error("output should contain Summary")
	}
	if !strings.Contains(output, "Medium Findings") {
		t.Error("output should contain Medium Findings")
	}
	if !strings.Contains(output, "Next Steps") {
		t.Error("output should contain Next Steps")
	}
	if !strings.Contains(output, "Ready for:") {
		t.Error("output should contain Ready for:")
	}
	// Check box-drawing characters
	if !strings.Contains(output, "┌") {
		t.Error("output should contain box-drawing characters")
	}
}

func TestRenderRichMarkdown(t *testing.T) {
	tmpDir := t.TempDir()

	report := &Report{
		Project: "test-project",
		Path:    tmpDir,
		Specs: []SpecStatus{
			{
				Type:     types.SpecTypeMRD,
				Category: types.CategorySource,
				Exists:   true,
				Required: true,
				EvalStatus: &EvalStatus{
					Exists:   true,
					Decision: "pass",
					Categories: &CategoryBreakdown{
						Pass:  6,
						Total: 6,
					},
				},
			},
		},
		Summary: Summary{TotalSpecs: 1, PresentSpecs: 1, EvaluatedSpecs: 1},
	}

	rr := NewRichReport(report)
	rr.SetKeyDecisions([]KeyDecision{
		{Area: "Cloud", Choice: "AWS with Pulumi"},
		{Area: "IaC", Choice: "Pulumi Go SDK"},
	})

	var buf bytes.Buffer
	if err := RenderRichMarkdown(&buf, rr); err != nil {
		t.Fatalf("RenderRichMarkdown failed: %v", err)
	}

	output := buf.String()

	// Check markdown structure
	if !strings.Contains(output, "# VisionSpec Status:") {
		t.Error("output should contain markdown header")
	}
	if !strings.Contains(output, "## Pipeline Progress") {
		t.Error("output should contain Pipeline Progress section")
	}
	if !strings.Contains(output, "## Summary") {
		t.Error("output should contain Summary section")
	}
	if !strings.Contains(output, "## Key Decisions") {
		t.Error("output should contain Key Decisions section")
	}
	if !strings.Contains(output, "| Decision | Choice |") {
		t.Error("output should contain key decisions table")
	}
}

func TestBuildPipeline(t *testing.T) {
	tmpDir := t.TempDir()

	report := &Report{
		Path: tmpDir,
		Specs: []SpecStatus{
			{Type: types.SpecTypeMRD, Exists: true, EvalStatus: &EvalStatus{Decision: "pass"}},
			{Type: types.SpecTypePRD, Exists: true, Approval: &types.Approval{}},
			{Type: types.SpecTypeUXD, Exists: true}, // Draft
			{Type: types.SpecTypeTRD, Exists: false},
		},
	}

	pipeline := buildPipeline(report)

	// Check pipeline length (should include all standard stages)
	if len(pipeline) != 9 {
		t.Errorf("len(pipeline) = %d, want 9", len(pipeline))
	}

	// Find MRD stage - should be complete (has passing eval)
	for _, stage := range pipeline {
		if stage.Type == types.SpecTypeMRD {
			if stage.Status != StageComplete {
				t.Errorf("MRD status = %s, want %s", stage.Status, StageComplete)
			}
		}
		if stage.Type == types.SpecTypePRD {
			if stage.Status != StageComplete {
				t.Errorf("PRD status = %s, want %s (has approval)", stage.Status, StageComplete)
			}
		}
		if stage.Type == types.SpecTypeUXD {
			if stage.Status != StagePending {
				t.Errorf("UXD status = %s, want %s (draft)", stage.Status, StagePending)
			}
		}
		if stage.Type == types.SpecTypeTRD {
			if stage.Status != StageMissing {
				t.Errorf("TRD status = %s, want %s", stage.Status, StageMissing)
			}
		}
	}
}

func TestFormatFindings(t *testing.T) {
	tests := []struct {
		name string
		eval *EvalStatus
		want string
	}{
		{
			name: "nil eval",
			eval: nil,
			want: "-",
		},
		{
			name: "no findings",
			eval: &EvalStatus{},
			want: "-",
		},
		{
			name: "single finding type",
			eval: &EvalStatus{
				Findings: struct {
					Critical int `json:"critical"`
					High     int `json:"high"`
					Medium   int `json:"medium"`
					Low      int `json:"low"`
					Info     int `json:"info"`
				}{Info: 2},
			},
			want: "2 info",
		},
		{
			name: "multiple finding types",
			eval: &EvalStatus{
				Findings: struct {
					Critical int `json:"critical"`
					High     int `json:"high"`
					Medium   int `json:"medium"`
					Low      int `json:"low"`
					Info     int `json:"info"`
				}{Medium: 1, Low: 2, Info: 1},
			},
			want: "1 medium, 2 low, 1 info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatFindings(tt.eval)
			if got != tt.want {
				t.Errorf("formatFindings() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		s      string
		maxLen int
		want   string
	}{
		{"short", 10, "short"},
		{"exactly10!", 10, "exactly10!"},
		{"this is a long string", 10, "this is..."},
	}

	for _, tt := range tests {
		got := truncate(tt.s, tt.maxLen)
		if got != tt.want {
			t.Errorf("truncate(%q, %d) = %q, want %q", tt.s, tt.maxLen, got, tt.want)
		}
	}
}
