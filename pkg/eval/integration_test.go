package eval_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/ProductBuildersHQ/visionspec/pkg/eval"
	"github.com/ProductBuildersHQ/visionspec/pkg/types"
)

// mockResult creates a sample Result for testing.
func mockResult() *eval.Result {
	return &eval.Result{
		SpecType:  types.SpecTypePRD,
		Timestamp: time.Now(),
		Score:     8.5,
		Passed:    true,
		Categories: []eval.CategoryResult{
			{
				ID:          "completeness",
				Name:        "Completeness",
				Score:       8.0,
				Weight:      0.3,
				Explanation: "Good coverage of requirements",
			},
			{
				ID:          "clarity",
				Name:        "Clarity",
				Score:       9.0,
				Weight:      0.3,
				Explanation: "Well-written and clear",
			},
		},
		Findings: []eval.Finding{
			{
				Severity:       "medium",
				Category:       "completeness",
				Title:          "Missing error handling",
				Description:    "Error handling scenarios not documented",
				Recommendation: "Add error handling section",
			},
		},
		Decision: "pass",
		Summary:  "Overall good PRD with minor improvements needed",
		Judge: eval.JudgeMetadata{
			Model:       "claude-opus-4-5-20251101",
			Provider:    "anthropic",
			Temperature: 0.0,
			Tokens:      1500,
		},
	}
}

func TestMarkdownRenderer(t *testing.T) {
	result := mockResult()
	renderer := eval.NewMarkdownRenderer()

	var buf bytes.Buffer
	err := renderer.Render(&buf, result)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	output := buf.String()

	// Check for expected content
	if !bytes.Contains([]byte(output), []byte("# PRD Evaluation Report")) {
		t.Error("Missing header")
	}
	if !bytes.Contains([]byte(output), []byte("✅ PASS")) {
		t.Error("Missing PASS decision")
	}
	if !bytes.Contains([]byte(output), []byte("Completeness")) {
		t.Error("Missing category name")
	}
	if !bytes.Contains([]byte(output), []byte("Missing error handling")) {
		t.Error("Missing finding title")
	}
}

func TestTerminalRenderer(t *testing.T) {
	result := mockResult()
	renderer := eval.NewTerminalRenderer(false)

	var buf bytes.Buffer
	err := renderer.Render(&buf, result)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	output := buf.String()

	// Check for expected content
	if !bytes.Contains([]byte(output), []byte("prd")) {
		t.Error("Missing spec type")
	}
	if !bytes.Contains([]byte(output), []byte("PASS")) {
		t.Error("Missing PASS decision")
	}
}

func TestTerminalRendererVerbose(t *testing.T) {
	result := mockResult()
	renderer := eval.NewTerminalRenderer(true)

	var buf bytes.Buffer
	err := renderer.Render(&buf, result)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	output := buf.String()

	// Verbose mode should include explanations
	if !bytes.Contains([]byte(output), []byte("Good coverage")) {
		t.Error("Missing category explanation in verbose mode")
	}
}

func TestToClaimsReport(t *testing.T) {
	result := mockResult()
	claimsReport := result.ToClaimsReport("prd.md")

	// Should have claims for each finding + summary
	expectedClaims := len(result.Findings) + 1 // +1 for summary
	if len(claimsReport.Claims) != expectedClaims {
		t.Errorf("Expected %d claims, got %d", expectedClaims, len(claimsReport.Claims))
	}

	// First claim should be the finding
	if claimsReport.Claims[0].ID != "finding-1" {
		t.Errorf("Expected first claim ID 'finding-1', got %s", claimsReport.Claims[0].ID)
	}

	// Should be finalized
	if claimsReport.Decision.Status == "" {
		t.Error("Claims report should be finalized")
	}
}

func TestEvalSummary(t *testing.T) {
	result := mockResult()
	summary := eval.NewEvalSummary("test-project", "1.0.0")
	summary.AddResult("prd", result, nil, nil)

	if !summary.IsAllPassing() {
		t.Error("Expected summary to be passing")
	}

	avgScore := summary.TotalScore()
	if avgScore != 8.5 {
		t.Errorf("Expected average score 8.5, got %.1f", avgScore)
	}
}

func TestToSummaryReport(t *testing.T) {
	result := mockResult()
	summary := eval.NewEvalSummary("test-project", "1.0.0")
	summary.AddResult("prd", result, nil, nil)

	report := summary.ToSummaryReport("SPEC EVALUATION")

	if report.Project != "test-project" {
		t.Errorf("Expected project 'test-project', got %s", report.Project)
	}
	if report.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got %s", report.Version)
	}
	if len(report.Teams) != 1 {
		t.Errorf("Expected 1 team section, got %d", len(report.Teams))
	}
}
