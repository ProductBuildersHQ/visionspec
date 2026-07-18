// Package aidlc provides AIDLC document evaluation using structured-evaluation rubrics.
package aidlc

import (
	"context"
	"fmt"
	"time"

	"github.com/ProductBuildersHQ/visionspec/pkg/rubrics"
	"github.com/ProductBuildersHQ/visionspec/pkg/types"
	"github.com/plexusone/structured-evaluation/rubric"
)

// docTypeToSpecType maps AIDLC DocType to visionspec SpecType for rubric lookup.
var docTypeToSpecType = map[DocType]types.SpecType{
	// Inception phase
	DocVisionDocument:   types.SpecTypeAIDLCVision,
	DocRequirementsSpec: types.SpecTypeAIDLCRequirements,
	DocTechnicalSpec:    types.SpecTypeAIDLCTechnical,
	DocArchitectureSpec: types.SpecTypeAIDLCArchitecture,

	// Construction phase
	DocImplementationPlan: types.SpecTypeAIDLCImplementation,
	DocTestPlan:           types.SpecTypeAIDLCTestPlan,
	DocIntegrationPlan:    types.SpecTypeAIDLCIntegration,
	DocSecurityReview:     types.SpecTypeAIDLCSecurity,

	// Operations phase
	DocRunbook:        types.SpecTypeAIDLCRunbook,
	DocMonitoringPlan: types.SpecTypeAIDLCMonitoring,
	DocDisasterPlan:   types.SpecTypeAIDLCDisaster,
	DocSLODocument:    types.SpecTypeAIDLCSLO,
}

// GetSpecType returns the visionspec SpecType for an AIDLC DocType.
func GetSpecType(docType DocType) (types.SpecType, bool) {
	st, ok := docTypeToSpecType[docType]
	return st, ok
}

// Judge is the interface for LLM-based document evaluation.
// Implementations provide the actual LLM interaction for evaluating documents.
type Judge interface {
	// EvaluateCategory evaluates a document against a single category's criteria.
	// Returns the category result with score, reasoning, and any findings.
	EvaluateCategory(ctx context.Context, content string, category *rubric.Category) (*rubric.CategoryResult, error)
}

// Evaluator evaluates AIDLC documents using structured-evaluation rubrics.
type Evaluator struct {
	// Judge is the LLM judge for evaluation.
	// When nil, Evaluate returns a stub result.
	Judge Judge
}

// NewEvaluator creates a new AIDLC document evaluator.
func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

// NewEvaluatorWithJudge creates a new evaluator with the specified judge.
func NewEvaluatorWithJudge(judge Judge) *Evaluator {
	return &Evaluator{Judge: judge}
}

// EvaluationResult contains the result of evaluating an AIDLC document.
type EvaluationResult struct {
	// Document is the evaluated document.
	Document *Document

	// RubricSet is the rubric used for evaluation.
	RubricSet *rubric.RubricSet

	// Report is the full evaluation report.
	Report *rubric.Rubric

	// QualityScore is the aggregated quality score in AIDLC format.
	QualityScore *QualityScore

	// Decision is the evaluation decision.
	Decision *rubric.Decision

	// Error contains any evaluation error.
	Error error
}

// GetRubricSet returns the structured-evaluation rubric set for a document type.
func GetRubricSet(docType DocType) (*rubric.RubricSet, error) {
	specType, ok := GetSpecType(docType)
	if !ok {
		return nil, fmt.Errorf("no spec type mapping for AIDLC doc type %q", docType)
	}

	rs, err := rubrics.Get(specType)
	if err != nil {
		return nil, fmt.Errorf("failed to get rubric for %s: %w", specType, err)
	}

	return rs.ToEvaluationRubricSet(), nil
}

// Evaluate evaluates an AIDLC document using its rubric.
func (e *Evaluator) Evaluate(ctx context.Context, doc *Document) (*EvaluationResult, error) {
	result := &EvaluationResult{
		Document: doc,
	}

	// Get the rubric for this document type
	rs, err := GetRubricSet(doc.Type)
	if err != nil {
		result.Error = err
		return result, err
	}
	result.RubricSet = rs

	// Create a new evaluation report
	report := rubric.NewRubric(string(doc.Type), doc.Path)
	report.Metadata.DocumentTitle = doc.Title
	report.RubricID = rs.ID
	report.RubricVersion = rs.Version
	result.Report = report

	// If no judge is configured, return a stub result
	if e.Judge == nil {
		result.QualityScore = createStubQualityScore(rs)
		return result, nil
	}

	// Evaluate each category using the judge
	var allFindings []rubric.Finding
	for _, cat := range rs.Categories {
		catResult, err := e.Judge.EvaluateCategory(ctx, doc.Content, &cat)
		if err != nil {
			// Record error as a failing category result
			report.Categories = append(report.Categories, rubric.CategoryResult{
				Category:  cat.ID,
				Score:     rubric.ScoreFail,
				Reasoning: fmt.Sprintf("Evaluation error: %v", err),
			})
			continue
		}
		report.Categories = append(report.Categories, *catResult)
		allFindings = append(allFindings, catResult.Findings...)
	}
	report.Findings = allFindings

	// Compute the decision
	decision := report.Evaluate(rs)
	result.Decision = &decision

	// Convert to QualityScore
	result.QualityScore = reportToQualityScore(report, rs)

	return result, nil
}

// createStubQualityScore creates a stub quality score when no judge is available.
func createStubQualityScore(rs *rubric.RubricSet) *QualityScore {
	qs := &QualityScore{
		Rating:      RatingNeedsImprovement,
		Score:       0.5,
		EvaluatedAt: time.Now(),
		Dimensions:  make(map[string]DimensionScore),
	}
	for _, cat := range rs.Categories {
		qs.Dimensions[cat.ID] = DimensionScore{
			ID:     cat.ID,
			Name:   cat.Name,
			Score:  0.5,
			Weight: cat.Weight,
		}
	}
	return qs
}

// reportToQualityScore converts a rubric report to QualityScore.
func reportToQualityScore(report *rubric.Rubric, rs *rubric.RubricSet) *QualityScore {
	qs := &QualityScore{
		EvaluatedAt: time.Now(),
		Dimensions:  make(map[string]DimensionScore),
		Issues:      []Issue{},
	}

	var totalWeight float64
	var weightedScore float64

	for _, result := range report.Categories {
		cat := rs.GetCategory(result.Category)
		weight := 1.0
		if cat != nil {
			weight = cat.Weight
		}

		// Convert score to numeric (0-1 range)
		score := scoreValueToNumeric(result.Score)

		weightedScore += score * weight
		totalWeight += weight

		// Add to dimensions
		dimName := result.Category
		if cat != nil {
			dimName = cat.Name
		}
		qs.Dimensions[result.Category] = DimensionScore{
			ID:       result.Category,
			Name:     dimName,
			Score:    score,
			Weight:   weight,
			Findings: convertFindings(result.Findings),
		}
	}

	// Convert all findings to issues
	for _, finding := range report.Findings {
		qs.Issues = append(qs.Issues, Issue{
			Severity:   convertSeverity(finding.Severity),
			Category:   finding.Category,
			Code:       string(finding.Code),
			Message:    finding.Description,
			Location:   finding.Location,
			Suggestion: finding.Recommendation,
		})
	}

	// Calculate overall score
	if totalWeight > 0 {
		qs.Score = weightedScore / totalWeight
	}

	// Determine rating
	qs.Rating = scoreToRating(qs.Score)

	return qs
}

// scoreValueToNumeric converts a categorical score to a numeric value (0-1).
func scoreValueToNumeric(score rubric.ScoreValue) float64 {
	switch score {
	case rubric.ScorePass:
		return 1.0
	case rubric.ScorePartial:
		return 0.6
	case rubric.ScoreFail:
		return 0.2
	default:
		return 0.5
	}
}

// convertFindings converts structured-evaluation findings to AIDLC issues.
func convertFindings(findings []rubric.Finding) []Issue {
	var issues []Issue
	for _, f := range findings {
		issues = append(issues, Issue{
			Severity:   convertSeverity(f.Severity),
			Category:   f.Category,
			Code:       string(f.Code),
			Message:    f.Description,
			Location:   f.Location,
			Suggestion: f.Recommendation,
		})
	}
	return issues
}

// convertSeverity converts structured-evaluation severity to AIDLC severity.
func convertSeverity(s rubric.Severity) IssueSeverity {
	switch s {
	case rubric.SeverityCritical:
		return SeverityCritical
	case rubric.SeverityHigh:
		return SeverityHigh
	case rubric.SeverityMedium:
		return SeverityMedium
	case rubric.SeverityLow:
		return SeverityLow
	case rubric.SeverityInfo:
		return SeverityInfo
	default:
		return SeverityMedium
	}
}

// scoreToRating converts a numeric score to a quality rating.
func scoreToRating(score float64) QualityRating {
	switch {
	case score >= 0.9:
		return RatingExcellent
	case score >= 0.7:
		return RatingGood
	case score >= 0.5:
		return RatingNeedsImprovement
	default:
		return RatingPoor
	}
}

// EvaluateDocument is a convenience function to evaluate a document with a judge.
func EvaluateDocument(ctx context.Context, doc *Document, judge Judge) (*EvaluationResult, error) {
	eval := NewEvaluatorWithJudge(judge)
	return eval.Evaluate(ctx, doc)
}

// EvaluateContent evaluates raw content for a specific document type.
func EvaluateContent(ctx context.Context, docType DocType, content string, judge Judge) (*EvaluationResult, error) {
	doc := &Document{
		Type:    docType,
		Phase:   docType.Phase(),
		Title:   docType.DisplayName(),
		Content: content,
		Status:  StatusReview,
	}
	return EvaluateDocument(ctx, doc, judge)
}

// CheckRubricAvailable verifies that a rubric exists for the document type.
func CheckRubricAvailable(docType DocType) error {
	_, err := GetRubricSet(docType)
	return err
}

// AllDocTypesWithRubrics returns all document types that have rubrics defined.
func AllDocTypesWithRubrics() []DocType {
	var result []DocType
	for docType := range docTypeToSpecType {
		if CheckRubricAvailable(docType) == nil {
			result = append(result, docType)
		}
	}
	return result
}

// BuildCategoryPrompt builds the evaluation prompt for a category.
// This can be used by Judge implementations to construct prompts.
func BuildCategoryPrompt(docType DocType, content string, cat *rubric.Category) string {
	displayName := docType.DisplayName()

	prompt := fmt.Sprintf(`Evaluate the following %s document for the category "%s".

Category: %s
Description: %s

Criteria:
`, displayName, cat.Name, cat.Name, cat.Description)

	for _, opt := range cat.Scale.Options {
		prompt += fmt.Sprintf("\n%s:\n", opt.Label)
		for _, criterion := range opt.Criteria {
			prompt += fmt.Sprintf("  - %s\n", criterion)
		}
	}

	prompt += fmt.Sprintf("\n\nDocument Content:\n---\n%s\n---\n", content)

	if cat.EvaluationPrompt != "" {
		prompt += fmt.Sprintf("\nAdditional guidance: %s\n", cat.EvaluationPrompt)
	}

	prompt += `
Respond with a JSON object containing:
{
  "score": "pass" | "partial" | "fail",
  "confidence": 0.0-1.0,
  "reasoning": "Explanation of the score",
  "evidence": ["Specific quotes or observations"],
  "findings": [
    {
      "severity": "critical" | "high" | "medium" | "low" | "info",
      "title": "Brief issue title",
      "description": "Detailed explanation",
      "recommendation": "How to fix"
    }
  ]
}
`

	return prompt
}
