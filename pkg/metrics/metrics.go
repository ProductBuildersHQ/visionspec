// Package metrics provides evaluation and reconciliation metrics tracking.
//
// This package collects and aggregates metrics about spec quality,
// evaluation results, and reconciliation outcomes to provide insights
// into the specification process.
package metrics

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// MetricType categorizes different types of metrics.
type MetricType string

const (
	MetricTypeEval      MetricType = "evaluation"
	MetricTypeReconcile MetricType = "reconciliation"
	MetricTypeAlign     MetricType = "alignment"
	MetricTypeDrift     MetricType = "drift"
	MetricTypeCoverage  MetricType = "coverage"
)

// Metric represents a single metric data point.
type Metric struct {
	Type      MetricType        `json:"type"`
	Name      string            `json:"name"`
	Value     float64           `json:"value"`
	Unit      string            `json:"unit,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
	Labels    map[string]string `json:"labels,omitempty"`
}

// EvalMetrics contains evaluation-related metrics.
type EvalMetrics struct {
	TotalEvaluations int                `json:"total_evaluations"`
	PassCount        int                `json:"pass_count"`
	FailCount        int                `json:"fail_count"`
	AverageScore     float64            `json:"average_score"`
	ScoresBySpec     map[string]float64 `json:"scores_by_spec"`
	FindingsCount    int                `json:"findings_count"`
	BySeverity       map[string]int     `json:"by_severity"`
	Trend            []TrendPoint       `json:"trend,omitempty"`
}

// ReconcileMetrics contains reconciliation-related metrics.
type ReconcileMetrics struct {
	TotalReconciliations int        `json:"total_reconciliations"`
	SuccessCount         int        `json:"success_count"`
	ConflictCount        int        `json:"conflict_count"`
	AverageTime          float64    `json:"average_time_seconds"`
	SpecsIncluded        int        `json:"specs_included"`
	TasksGenerated       int        `json:"tasks_generated"`
	LastReconcile        *time.Time `json:"last_reconcile,omitempty"`
}

// AlignMetrics contains alignment-related metrics.
type AlignMetrics struct {
	AlignmentScore   float64 `json:"alignment_score"`
	CoveragePercent  float64 `json:"coverage_percent"`
	DiscrepancyCount int     `json:"discrepancy_count"`
	CriticalCount    int     `json:"critical_count"`
	HighCount        int     `json:"high_count"`
	MediumCount      int     `json:"medium_count"`
	LowCount         int     `json:"low_count"`
	MissingFeatures  int     `json:"missing_features"`
	UndocumentedCode int     `json:"undocumented_code"`
}

// DriftMetrics contains drift-related metrics.
type DriftMetrics struct {
	HasDrift       bool    `json:"has_drift"`
	DriftScore     float64 `json:"drift_score"`
	ItemCount      int     `json:"item_count"`
	CriticalCount  int     `json:"critical_count"`
	HighCount      int     `json:"high_count"`
	TrendDirection string  `json:"trend_direction"` // "improving", "stable", "degrading"
}

// TrendPoint represents a point in a trend over time.
type TrendPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
	Label     string    `json:"label,omitempty"`
}

// ProjectMetrics contains all metrics for a project.
type ProjectMetrics struct {
	Project      string            `json:"project"`
	GeneratedAt  time.Time         `json:"generated_at"`
	Eval         *EvalMetrics      `json:"evaluation,omitempty"`
	Reconcile    *ReconcileMetrics `json:"reconciliation,omitempty"`
	Align        *AlignMetrics     `json:"alignment,omitempty"`
	Drift        *DriftMetrics     `json:"drift,omitempty"`
	Summary      MetricsSummary    `json:"summary"`
	HealthScore  float64           `json:"health_score"`  // 0-100
	HealthStatus string            `json:"health_status"` // "healthy", "warning", "critical"
}

// MetricsSummary provides high-level project health indicators.
type MetricsSummary struct {
	TotalSpecs     int     `json:"total_specs"`
	ApprovedSpecs  int     `json:"approved_specs"`
	EvaluatedSpecs int     `json:"evaluated_specs"`
	PassingEvals   int     `json:"passing_evals"`
	OverallScore   float64 `json:"overall_score"`
	ReadinessScore float64 `json:"readiness_score"`
	QualityScore   float64 `json:"quality_score"`
}

// Collector gathers metrics from various sources.
type Collector struct {
	projectPath string
	history     *MetricsHistory
}

// NewCollector creates a new metrics collector.
func NewCollector(projectPath string) (*Collector, error) {
	history, err := NewMetricsHistory(filepath.Join(projectPath, ".visionspec", "metrics.json"))
	if err != nil {
		return nil, fmt.Errorf("loading metrics history: %w", err)
	}
	return &Collector{
		projectPath: projectPath,
		history:     history,
	}, nil
}

// History returns the metrics history tracker.
func (c *Collector) History() *MetricsHistory {
	return c.history
}

// Collect gathers all available metrics for the project.
func (c *Collector) Collect() (*ProjectMetrics, error) {
	metrics := &ProjectMetrics{
		Project:     filepath.Base(c.projectPath),
		GeneratedAt: time.Now(),
	}

	// Collect eval metrics
	evalMetrics, err := c.collectEvalMetrics()
	if err == nil && evalMetrics != nil {
		metrics.Eval = evalMetrics
	}

	// Collect reconcile metrics
	reconcileMetrics, err := c.collectReconcileMetrics()
	if err == nil && reconcileMetrics != nil {
		metrics.Reconcile = reconcileMetrics
	}

	// Calculate summary and health score
	metrics.Summary = c.calculateSummary(metrics)
	metrics.HealthScore, metrics.HealthStatus = c.calculateHealthScore(metrics)

	return metrics, nil
}

// collectEvalMetrics reads evaluation results from eval/ directory.
func (c *Collector) collectEvalMetrics() (*EvalMetrics, error) {
	evalDir := filepath.Join(c.projectPath, "eval")
	entries, err := os.ReadDir(evalDir)
	if err != nil {
		return nil, err
	}

	metrics := &EvalMetrics{
		ScoresBySpec: make(map[string]float64),
		BySeverity:   make(map[string]int),
	}

	var totalScore float64

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".eval.json") {
			continue
		}

		// Parse eval file
		evalPath := filepath.Join(evalDir, entry.Name())
		data, err := os.ReadFile(evalPath)
		if err != nil {
			continue
		}

		var eval struct {
			Score    float64 `json:"score"`
			Decision string  `json:"decision"`
			Findings []struct {
				Severity string `json:"severity"`
			} `json:"findings"`
		}

		if err := json.Unmarshal(data, &eval); err != nil {
			continue
		}

		metrics.TotalEvaluations++
		totalScore += eval.Score

		specName := entry.Name()[:len(entry.Name())-len(".eval.json")]
		metrics.ScoresBySpec[specName] = eval.Score

		if eval.Decision == "pass" || eval.Score >= 7.0 {
			metrics.PassCount++
		} else {
			metrics.FailCount++
		}

		for _, finding := range eval.Findings {
			metrics.FindingsCount++
			metrics.BySeverity[finding.Severity]++
		}
	}

	if metrics.TotalEvaluations > 0 {
		metrics.AverageScore = totalScore / float64(metrics.TotalEvaluations)
	}

	return metrics, nil
}

// collectReconcileMetrics reads reconciliation data.
func (c *Collector) collectReconcileMetrics() (*ReconcileMetrics, error) {
	metrics := &ReconcileMetrics{}

	// Check if spec.md exists
	specPath := filepath.Join(c.projectPath, "spec.md")
	info, err := os.Stat(specPath)
	if err != nil {
		return nil, err
	}

	modTime := info.ModTime()
	metrics.LastReconcile = &modTime
	metrics.SuccessCount = 1

	// Count specs included
	for _, dir := range []string{"source", "gtm", "technical"} {
		dirPath := filepath.Join(c.projectPath, dir)
		if entries, err := os.ReadDir(dirPath); err == nil {
			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
					metrics.SpecsIncluded++
				}
			}
		}
	}

	return metrics, nil
}

// calculateSummary computes summary metrics.
func (c *Collector) calculateSummary(metrics *ProjectMetrics) MetricsSummary {
	summary := MetricsSummary{}

	if metrics.Eval != nil {
		summary.EvaluatedSpecs = metrics.Eval.TotalEvaluations
		summary.PassingEvals = metrics.Eval.PassCount
		summary.QualityScore = metrics.Eval.AverageScore * 10 // Convert to 0-100
	}

	if metrics.Reconcile != nil {
		summary.TotalSpecs = metrics.Reconcile.SpecsIncluded
	}

	// Calculate overall score
	if summary.EvaluatedSpecs > 0 {
		summary.OverallScore = (float64(summary.PassingEvals) / float64(summary.EvaluatedSpecs)) * 100
	}

	return summary
}

// calculateHealthScore computes an overall health score.
func (c *Collector) calculateHealthScore(metrics *ProjectMetrics) (float64, string) {
	score := 100.0

	// Deduct for eval issues
	if metrics.Eval != nil {
		if metrics.Eval.TotalEvaluations > 0 {
			passRate := float64(metrics.Eval.PassCount) / float64(metrics.Eval.TotalEvaluations)
			score -= (1 - passRate) * 30
		}
		// Deduct for critical/high findings
		score -= float64(getMapValue(metrics.Eval.BySeverity, "critical")) * 5
		score -= float64(getMapValue(metrics.Eval.BySeverity, "high")) * 2
	}

	// Deduct for alignment issues
	if metrics.Align != nil {
		score -= float64(metrics.Align.CriticalCount) * 10
		score -= float64(metrics.Align.HighCount) * 5
		score -= float64(metrics.Align.MediumCount) * 1
	}

	// Deduct for drift
	if metrics.Drift != nil && metrics.Drift.HasDrift {
		score -= float64(metrics.Drift.CriticalCount) * 10
		score -= float64(metrics.Drift.HighCount) * 5
	}

	// Clamp score
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	// Determine status
	status := "healthy"
	if score < 50 {
		status = "critical"
	} else if score < 75 {
		status = "warning"
	}

	return score, status
}

// getMapValue helper to get count from a map.
func getMapValue(m map[string]int, key string) int {
	if v, ok := m[key]; ok {
		return v
	}
	return 0
}

// MetricsHistory tracks metrics over time.
type MetricsHistory struct {
	path    string
	entries []MetricsHistoryEntry
}

// MetricsHistoryEntry represents a single history entry.
type MetricsHistoryEntry struct {
	Timestamp   time.Time `json:"timestamp"`
	HealthScore float64   `json:"health_score"`
	EvalScore   float64   `json:"eval_score,omitempty"`
	AlignScore  float64   `json:"align_score,omitempty"`
}

// NewMetricsHistory creates a new metrics history tracker.
// Returns an error if the history file exists but cannot be parsed.
func NewMetricsHistory(path string) (*MetricsHistory, error) {
	h := &MetricsHistory{path: path}
	if err := h.load(); err != nil {
		return nil, err
	}
	return h, nil
}

// load reads history from file.
// Returns nil if file doesn't exist (empty history is valid).
// Returns error if file exists but cannot be parsed.
func (h *MetricsHistory) load() error {
	data, err := os.ReadFile(h.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No history file yet, start fresh
		}
		return fmt.Errorf("reading metrics history: %w", err)
	}
	if err := json.Unmarshal(data, &h.entries); err != nil {
		return fmt.Errorf("parsing metrics history: %w", err)
	}
	return nil
}

// Save persists history to file.
func (h *MetricsHistory) Save() error {
	dir := filepath.Dir(h.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(h.entries, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(h.path, data, 0600)
}

// Add adds a new entry to history.
func (h *MetricsHistory) Add(entry MetricsHistoryEntry) {
	h.entries = append(h.entries, entry)

	// Keep last 100 entries
	if len(h.entries) > 100 {
		h.entries = h.entries[len(h.entries)-100:]
	}
}

// Recent returns the most recent entries.
func (h *MetricsHistory) Recent(count int) []MetricsHistoryEntry {
	if count > len(h.entries) {
		count = len(h.entries)
	}
	return h.entries[len(h.entries)-count:]
}

// Trend returns the trend direction based on recent history.
func (h *MetricsHistory) Trend() string {
	if len(h.entries) < 2 {
		return "stable"
	}

	recent := h.Recent(5)
	if len(recent) < 2 {
		return "stable"
	}

	first := recent[0].HealthScore
	last := recent[len(recent)-1].HealthScore

	diff := last - first
	if diff > 5 {
		return "improving"
	} else if diff < -5 {
		return "degrading"
	}
	return "stable"
}
