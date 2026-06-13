package metrics

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewCollector(t *testing.T) {
	tmpDir := t.TempDir()
	collector, err := NewCollector(tmpDir)
	if err != nil {
		t.Fatalf("NewCollector failed: %v", err)
	}

	if collector == nil {
		t.Fatal("NewCollector returned nil")
	}

	if collector.projectPath != tmpDir {
		t.Errorf("projectPath = %s, want %s", collector.projectPath, tmpDir)
	}
}

func TestCollector_Collect(t *testing.T) {
	tmpDir := t.TempDir()

	// Create an eval directory with an eval file
	evalDir := filepath.Join(tmpDir, "eval")
	if err := os.MkdirAll(evalDir, 0755); err != nil {
		t.Fatalf("Failed to create eval dir: %v", err)
	}

	evalData := map[string]interface{}{
		"score":    8.5,
		"decision": "pass",
		"findings": []map[string]interface{}{
			{"severity": "low"},
			{"severity": "medium"},
		},
	}
	evalBytes, _ := json.Marshal(evalData)
	if err := os.WriteFile(filepath.Join(evalDir, "test.eval.json"), evalBytes, 0600); err != nil {
		t.Fatalf("Failed to write eval file: %v", err)
	}

	// Create spec.md and source directory
	if err := os.WriteFile(filepath.Join(tmpDir, "spec.md"), []byte("# Test Spec"), 0600); err != nil {
		t.Fatalf("Failed to write spec.md: %v", err)
	}
	sourceDir := filepath.Join(tmpDir, "source")
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "prd.md"), []byte("# PRD"), 0600); err != nil {
		t.Fatalf("Failed to write prd.md: %v", err)
	}

	collector, err := NewCollector(tmpDir)
	if err != nil {
		t.Fatalf("NewCollector failed: %v", err)
	}
	metrics, err := collector.Collect()
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	if metrics == nil {
		t.Fatal("Collect returned nil metrics")
	}

	// Check project name
	if metrics.Project != filepath.Base(tmpDir) {
		t.Errorf("Project = %s, want %s", metrics.Project, filepath.Base(tmpDir))
	}

	// Check eval metrics
	if metrics.Eval == nil {
		t.Error("Eval metrics should not be nil")
	} else {
		if metrics.Eval.TotalEvaluations != 1 {
			t.Errorf("TotalEvaluations = %d, want 1", metrics.Eval.TotalEvaluations)
		}
		if metrics.Eval.PassCount != 1 {
			t.Errorf("PassCount = %d, want 1", metrics.Eval.PassCount)
		}
		if metrics.Eval.AverageScore != 8.5 {
			t.Errorf("AverageScore = %f, want 8.5", metrics.Eval.AverageScore)
		}
		if metrics.Eval.FindingsCount != 2 {
			t.Errorf("FindingsCount = %d, want 2", metrics.Eval.FindingsCount)
		}
	}

	// Check reconcile metrics
	if metrics.Reconcile == nil {
		t.Error("Reconcile metrics should not be nil")
	} else {
		if metrics.Reconcile.SpecsIncluded < 1 {
			t.Errorf("SpecsIncluded = %d, want >= 1", metrics.Reconcile.SpecsIncluded)
		}
	}

	// Check health score
	if metrics.HealthScore < 0 || metrics.HealthScore > 100 {
		t.Errorf("HealthScore = %f, want 0-100", metrics.HealthScore)
	}
}

func TestCollector_History(t *testing.T) {
	tmpDir := t.TempDir()
	collector, err := NewCollector(tmpDir)
	if err != nil {
		t.Fatalf("NewCollector failed: %v", err)
	}

	history := collector.History()
	if history == nil {
		t.Fatal("History returned nil")
	}
}

func TestMetricsHistory_AddAndRecent(t *testing.T) {
	tmpDir := t.TempDir()
	historyPath := filepath.Join(tmpDir, "metrics.json")
	history, err := NewMetricsHistory(historyPath)
	if err != nil {
		t.Fatalf("NewMetricsHistory failed: %v", err)
	}

	// Add some entries
	now := time.Now()
	for i := 0; i < 5; i++ {
		history.Add(MetricsHistoryEntry{
			Timestamp:   now.Add(time.Duration(i) * time.Hour),
			HealthScore: float64(50 + i*10),
			EvalScore:   float64(7 + float64(i)*0.5),
		})
	}

	// Check recent entries
	recent := history.Recent(3)
	if len(recent) != 3 {
		t.Errorf("Recent(3) returned %d entries, want 3", len(recent))
	}

	// Most recent should have highest score
	if recent[len(recent)-1].HealthScore != 90 {
		t.Errorf("Most recent HealthScore = %f, want 90", recent[len(recent)-1].HealthScore)
	}
}

func TestMetricsHistory_Trend(t *testing.T) {
	tests := []struct {
		name   string
		scores []float64
		want   string
	}{
		{
			name:   "improving",
			scores: []float64{50, 55, 60, 65, 70},
			want:   "improving",
		},
		{
			name:   "degrading",
			scores: []float64{90, 85, 80, 75, 70},
			want:   "degrading",
		},
		{
			name:   "stable",
			scores: []float64{75, 76, 74, 75, 76},
			want:   "stable",
		},
		{
			name:   "insufficient_data",
			scores: []float64{75},
			want:   "stable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			history, err := NewMetricsHistory(filepath.Join(tmpDir, "metrics.json"))
			if err != nil {
				t.Fatalf("NewMetricsHistory failed: %v", err)
			}

			now := time.Now()
			for i, score := range tt.scores {
				history.Add(MetricsHistoryEntry{
					Timestamp:   now.Add(time.Duration(i) * time.Hour),
					HealthScore: score,
				})
			}

			trend := history.Trend()
			if trend != tt.want {
				t.Errorf("Trend() = %s, want %s", trend, tt.want)
			}
		})
	}
}

func TestMetricsHistory_SaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	historyPath := filepath.Join(tmpDir, ".visionspec", "metrics.json")
	history, err := NewMetricsHistory(historyPath)
	if err != nil {
		t.Fatalf("NewMetricsHistory failed: %v", err)
	}

	// Add entries
	now := time.Now()
	history.Add(MetricsHistoryEntry{
		Timestamp:   now,
		HealthScore: 85.5,
		EvalScore:   8.0,
		AlignScore:  90.0,
	})

	// Save
	if err := history.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(historyPath); err != nil {
		t.Fatalf("History file not created: %v", err)
	}

	// Load in new instance
	history2, err := NewMetricsHistory(historyPath)
	if err != nil {
		t.Fatalf("NewMetricsHistory (reload) failed: %v", err)
	}
	recent := history2.Recent(1)
	if len(recent) != 1 {
		t.Fatalf("Loaded history has %d entries, want 1", len(recent))
	}

	if recent[0].HealthScore != 85.5 {
		t.Errorf("Loaded HealthScore = %f, want 85.5", recent[0].HealthScore)
	}
}

func TestCalculateHealthScore(t *testing.T) {
	tests := []struct {
		name       string
		metrics    *ProjectMetrics
		wantMin    float64
		wantMax    float64
		wantStatus string
	}{
		{
			name: "perfect_health",
			metrics: &ProjectMetrics{
				Eval: &EvalMetrics{
					TotalEvaluations: 10,
					PassCount:        10,
					BySeverity:       map[string]int{},
				},
			},
			wantMin:    95,
			wantMax:    100,
			wantStatus: "healthy",
		},
		{
			name: "some_failures",
			metrics: &ProjectMetrics{
				Eval: &EvalMetrics{
					TotalEvaluations: 10,
					PassCount:        5,
					BySeverity:       map[string]int{},
				},
			},
			// 50% pass rate = -15 points (0.5 * 30), score = 85, status = "healthy"
			wantMin:    80,
			wantMax:    90,
			wantStatus: "healthy",
		},
		{
			name: "many_critical_findings",
			metrics: &ProjectMetrics{
				Eval: &EvalMetrics{
					TotalEvaluations: 10,
					PassCount:        0, // -30 points
					BySeverity: map[string]int{
						"critical": 10, // -50 points
					},
				},
			},
			// Total -80 points, score = 20, status = "critical"
			wantMin:    0,
			wantMax:    30,
			wantStatus: "critical",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			collector, err := NewCollector(tmpDir)
			if err != nil {
				t.Fatalf("NewCollector failed: %v", err)
			}
			score, status := collector.calculateHealthScore(tt.metrics)

			if score < tt.wantMin || score > tt.wantMax {
				t.Errorf("HealthScore = %f, want %f-%f", score, tt.wantMin, tt.wantMax)
			}
			if status != tt.wantStatus {
				t.Errorf("Status = %s, want %s", status, tt.wantStatus)
			}
		})
	}
}

func TestDashboard_RenderJSON(t *testing.T) {
	metrics := &ProjectMetrics{
		Project:      "test-project",
		GeneratedAt:  time.Now(),
		HealthScore:  85.0,
		HealthStatus: "healthy",
		Summary: MetricsSummary{
			TotalSpecs:     5,
			EvaluatedSpecs: 3,
			PassingEvals:   2,
			OverallScore:   66.7,
		},
	}

	dashboard := NewDashboard(metrics)
	var buf bytes.Buffer

	err := dashboard.Render(&buf, FormatJSON)
	if err != nil {
		t.Fatalf("Render JSON failed: %v", err)
	}

	// Verify it's valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	if result["project"] != "test-project" {
		t.Errorf("project = %v, want test-project", result["project"])
	}
}

func TestDashboard_RenderTerminal(t *testing.T) {
	metrics := &ProjectMetrics{
		Project:      "test-project",
		GeneratedAt:  time.Now(),
		HealthScore:  85.0,
		HealthStatus: "healthy",
		Summary: MetricsSummary{
			TotalSpecs:     5,
			EvaluatedSpecs: 3,
			PassingEvals:   2,
		},
		Eval: &EvalMetrics{
			TotalEvaluations: 3,
			PassCount:        2,
			FailCount:        1,
			AverageScore:     7.5,
			BySeverity: map[string]int{
				"low":    3,
				"medium": 1,
			},
		},
	}

	dashboard := NewDashboard(metrics)
	var buf bytes.Buffer

	err := dashboard.Render(&buf, FormatTerminal)
	if err != nil {
		t.Fatalf("Render Terminal failed: %v", err)
	}

	output := buf.String()

	// Check for expected content
	if !bytes.Contains(buf.Bytes(), []byte("test-project")) {
		t.Error("Output should contain project name")
	}
	if !bytes.Contains(buf.Bytes(), []byte("Health Score")) {
		t.Error("Output should contain Health Score")
	}
	if len(output) < 100 {
		t.Error("Terminal output seems too short")
	}
}

func TestDashboard_RenderHTML(t *testing.T) {
	metrics := &ProjectMetrics{
		Project:      "test-project",
		GeneratedAt:  time.Now(),
		HealthScore:  75.0,
		HealthStatus: "warning",
		Summary: MetricsSummary{
			TotalSpecs: 5,
		},
	}

	dashboard := NewDashboard(metrics)
	var buf bytes.Buffer

	err := dashboard.Render(&buf, FormatHTML)
	if err != nil {
		t.Fatalf("Render HTML failed: %v", err)
	}

	output := buf.String()

	// Check for HTML structure
	if !bytes.Contains(buf.Bytes(), []byte("<!DOCTYPE html>")) {
		t.Error("Output should be HTML")
	}
	if !bytes.Contains(buf.Bytes(), []byte("test-project")) {
		t.Error("Output should contain project name")
	}
	if len(output) < 500 {
		t.Error("HTML output seems too short")
	}
}

func TestDashboard_RenderMarkdown(t *testing.T) {
	metrics := &ProjectMetrics{
		Project:      "test-project",
		GeneratedAt:  time.Now(),
		HealthScore:  90.0,
		HealthStatus: "healthy",
		Summary: MetricsSummary{
			TotalSpecs:   5,
			OverallScore: 80.0,
		},
		Eval: &EvalMetrics{
			TotalEvaluations: 5,
			PassCount:        4,
			FailCount:        1,
			AverageScore:     8.0,
		},
	}

	dashboard := NewDashboard(metrics)
	var buf bytes.Buffer

	err := dashboard.Render(&buf, FormatMarkdown)
	if err != nil {
		t.Fatalf("Render Markdown failed: %v", err)
	}

	output := buf.String()

	// Check for Markdown structure
	if !bytes.Contains(buf.Bytes(), []byte("# VisionSpec Project Metrics")) {
		t.Error("Output should have Markdown heading")
	}
	if !bytes.Contains(buf.Bytes(), []byte("|")) {
		t.Error("Output should have Markdown tables")
	}
	if len(output) < 200 {
		t.Error("Markdown output seems too short")
	}
}

func TestGetMapValue(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
	}

	if v := getMapValue(m, "a"); v != 1 {
		t.Errorf("getMapValue(m, 'a') = %d, want 1", v)
	}

	if v := getMapValue(m, "c"); v != 0 {
		t.Errorf("getMapValue(m, 'c') = %d, want 0", v)
	}

	if v := getMapValue(nil, "a"); v != 0 {
		t.Errorf("getMapValue(nil, 'a') = %d, want 0", v)
	}
}
