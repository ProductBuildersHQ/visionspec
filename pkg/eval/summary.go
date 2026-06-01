// Package eval provides evaluation orchestration for spec documents.
package eval

import (
	"time"

	"github.com/plexusone/structured-evaluation/claims"
	"github.com/plexusone/structured-evaluation/rubric"
	"github.com/plexusone/structured-evaluation/summary"
)

// EvalSummary aggregates multiple evaluation results with embedded reports.
type EvalSummary struct {
	// Project is the project name.
	Project string

	// Version is the project version.
	Version string

	// Results are the individual evaluation results keyed by spec type.
	Results map[string]*Result

	// EvaluationReports are the structured evaluation reports.
	EvaluationReports map[string]*rubric.Rubric

	// ClaimsReports are the claims extracted from findings.
	ClaimsReports map[string]*claims.ClaimsReport
}

// NewEvalSummary creates a new evaluation summary.
func NewEvalSummary(project, version string) *EvalSummary {
	return &EvalSummary{
		Project:           project,
		Version:           version,
		Results:           make(map[string]*Result),
		EvaluationReports: make(map[string]*rubric.Rubric),
		ClaimsReports:     make(map[string]*claims.ClaimsReport),
	}
}

// AddResult adds an evaluation result to the summary.
func (s *EvalSummary) AddResult(specType string, result *Result, evalReport *rubric.Rubric, claimsReport *claims.ClaimsReport) {
	s.Results[specType] = result
	if evalReport != nil {
		s.EvaluationReports[specType] = evalReport
	}
	if claimsReport != nil {
		s.ClaimsReports[specType] = claimsReport
	}
}

// IsAllPassing returns true if all evaluations passed.
func (s *EvalSummary) IsAllPassing() bool {
	for _, r := range s.Results {
		if !r.Passed {
			return false
		}
	}
	return len(s.Results) > 0
}

// TotalScore returns the average score across all evaluations.
func (s *EvalSummary) TotalScore() float64 {
	if len(s.Results) == 0 {
		return 0
	}
	var total float64
	for _, r := range s.Results {
		total += r.Score
	}
	return total / float64(len(s.Results))
}

// ToSummaryReport converts to a structured-evaluation SummaryReport.
// The report embeds full-fidelity EvaluationReport and ClaimsReport.
func (s *EvalSummary) ToSummaryReport(phase string) *summary.SummaryReport {
	report := summary.NewSummaryReport(s.Project, s.Version, phase)
	report.GeneratedBy = "visionspec eval"
	report.GeneratedAt = time.Now().UTC()

	// Create team sections from results
	for specType, result := range s.Results {
		team := summary.TeamSection{
			ID:   specType,
			Name: string(result.SpecType) + " Evaluation",
		}

		// Add tasks for each category
		for _, cat := range result.Categories {
			status := summary.StatusNoGo
			if cat.Score >= 7.0 {
				status = summary.StatusGo
			} else if cat.Score >= 5.0 {
				status = summary.StatusWarn
			}

			team.Tasks = append(team.Tasks, summary.TaskResult{
				ID:     cat.ID,
				Status: status,
				Detail: cat.Explanation,
			})
		}

		// Compute section status
		team.ComputeStatus()
		report.AddTeam(team)
	}

	// Compute overall status
	report.ComputeOverallStatus()

	// Embed full-fidelity reports
	for specType, evalReport := range s.EvaluationReports {
		_ = report.EmbedEvaluationReport(specType, evalReport)
	}

	for specType, claimsReport := range s.ClaimsReports {
		_ = report.EmbedClaimsReport(specType, claimsReport)
	}

	return report
}

// CreateSingleEvalSummary creates a summary with a single evaluation.
func CreateSingleEvalSummary(
	project string,
	specType string,
	result *Result,
	evalReport *rubric.Rubric,
	claimsReport *claims.ClaimsReport,
) *summary.SummaryReport {
	s := NewEvalSummary(project, "")
	s.AddResult(specType, result, evalReport, claimsReport)
	return s.ToSummaryReport("SPEC EVALUATION")
}

// CreateMultiEvalSummary creates a summary with multiple evaluations.
func CreateMultiEvalSummary(
	project string,
	version string,
	results map[string]*Result,
	evalReports map[string]*rubric.Rubric,
	claimsReports map[string]*claims.ClaimsReport,
) *summary.SummaryReport {
	s := NewEvalSummary(project, version)
	for specType, result := range results {
		var evalReport *rubric.Rubric
		var claimsReport *claims.ClaimsReport
		if evalReports != nil {
			evalReport = evalReports[specType]
		}
		if claimsReports != nil {
			claimsReport = claimsReports[specType]
		}
		s.AddResult(specType, result, evalReport, claimsReport)
	}
	return s.ToSummaryReport("MULTI-SPEC EVALUATION")
}
