// Package eval provides evaluation orchestration for spec documents.
package eval

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/plexusone/structured-evaluation/rubric"

	"github.com/ProductBuildersHQ/visionspec/pkg/rubrics"
	"github.com/ProductBuildersHQ/visionspec/pkg/types"
)

// Result represents the outcome of an evaluation.
type Result struct {
	// Schema version for backwards compatibility
	SchemaVersion string `json:"schemaVersion,omitempty"`

	SpecType  types.SpecType `json:"spec_type"`
	Timestamp time.Time      `json:"timestamp"`

	// V1 score (0-10, deprecated but kept for compatibility)
	Score float64 `json:"score"`

	// V2 integer score (1-5)
	IntScore rubric.IntegerScore `json:"intScore,omitempty"`

	// V2 confidence (0.0-1.0)
	Confidence float64 `json:"confidence,omitempty"`

	// Pass/fail gate
	Passed bool `json:"passed"`

	// V2 blocking reason codes
	Blocking []rubric.ReasonCode `json:"blocking,omitempty"`

	Categories []CategoryResult `json:"categories"`
	Findings   []Finding        `json:"findings"`
	Decision   string           `json:"decision"`
	Summary    string           `json:"summary"`
	Judge      JudgeMetadata    `json:"judge"`
}

// CategoryResult contains the evaluation result for a category.
type CategoryResult struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	// V1 score (0-10, deprecated)
	Score float64 `json:"score"`

	// V2 integer score (1-5)
	IntScore rubric.IntegerScore `json:"intScore,omitempty"`

	// V2 confidence (0.0-1.0)
	Confidence float64 `json:"confidence,omitempty"`

	// V2 reason codes
	ReasonCodes []rubric.ReasonCode `json:"reasonCodes,omitempty"`

	Weight      float64 `json:"weight"`
	Explanation string  `json:"explanation"`
}

// Finding represents an issue found during evaluation.
type Finding struct {
	Severity       string `json:"severity"`
	Category       string `json:"category"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Recommendation string `json:"recommendation"`
	Evidence       string `json:"evidence,omitempty"`

	// V2 fields
	Code     rubric.ReasonCode `json:"code,omitempty"`
	Location string            `json:"location,omitempty"`
}

// JudgeMetadata records information about the LLM judge.
type JudgeMetadata struct {
	Model       string  `json:"model"`
	Provider    string  `json:"provider"`
	Temperature float64 `json:"temperature"`
	Tokens      int     `json:"tokens"`
}

// Evaluator performs evaluations using an LLM judge.
type Evaluator struct {
	llm          *LLMClient
	rubricLoader rubrics.Loader
}

// NewEvaluator creates a new evaluator with the given LLM client.
func NewEvaluator(llm *LLMClient) *Evaluator {
	return &Evaluator{
		llm:          llm,
		rubricLoader: rubrics.DefaultLoader(),
	}
}

// SetRubricLoader sets a custom rubric loader for evaluation.
func (e *Evaluator) SetRubricLoader(loader rubrics.Loader) {
	if loader != nil {
		e.rubricLoader = loader
	}
}

// Evaluate runs evaluation on content against the rubric for the given spec type.
func (e *Evaluator) Evaluate(ctx context.Context, specType types.SpecType, content string) (*Result, error) {
	// Get rubric from loader
	rubricSet, err := e.rubricLoader.Load(specType)
	if err != nil {
		return nil, fmt.Errorf("no rubric for spec type %s: %w", specType, err)
	}

	// Build evaluation prompt
	prompt := buildEvalPrompt(rubricSet, content)

	// Call LLM
	response, metadata, err := e.llm.Complete(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("LLM evaluation failed: %w", err)
	}

	// Parse response
	result, err := parseEvalResponse(specType, rubricSet, response, metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to parse evaluation response: %w", err)
	}

	return result, nil
}

// buildEvalPrompt constructs the evaluation prompt for the LLM.
func buildEvalPrompt(rubricSet *rubric.RubricSet, content string) string {
	prompt := fmt.Sprintf(`You are an expert document evaluator. Evaluate the following %s against the provided rubric.

## Rubric

Evaluate each category on a scale of 1-5:
- 5 (Excellent): Exceeds expectations
- 4 (Good): Meets expectations well
- 3 (Acceptable): Meets minimum requirements
- 2 (Major Revisions): Significant work needed
- 1 (Unacceptable): Does not meet requirements

`, rubricSet.Name)

	for _, cat := range rubricSet.Categories {
		prompt += fmt.Sprintf("### %s (Weight: %.0f%%)\n%s\n\n", cat.Name, cat.Weight*100, cat.Description)
	}

	prompt += fmt.Sprintf(`## Document to Evaluate

%s

## Instructions

Provide your evaluation in the following JSON format:
{
  "categories": [
    {
      "id": "category_id",
      "intScore": 4,
      "confidence": 0.85,
      "reasonCodes": ["AMBIGUOUS_REQUIREMENT"],
      "explanation": "Brief explanation of the score"
    }
  ],
  "findings": [
    {
      "severity": "critical|high|medium|low|info",
      "category": "category_id",
      "code": "AMBIGUOUS_REQUIREMENT",
      "title": "Short title",
      "description": "Detailed description",
      "recommendation": "How to fix",
      "location": "Section 2.1"
    }
  ],
  "summary": "Overall assessment in 2-3 sentences",
  "overallConfidence": 0.85
}

### Score Scale (1-5)
- 5: Excellent - Exceeds expectations
- 4: Good - Meets expectations well
- 3: Acceptable - Meets minimum requirements
- 2: Major Revisions - Significant work needed
- 1: Unacceptable - Does not meet requirements

### Severity Levels
- critical: Fundamental issues that block approval
- high: Significant issues that should be fixed
- medium: Notable issues worth addressing
- low: Minor improvements
- info: Informational observations

### Common Reason Codes
- AMBIGUOUS_REQUIREMENT: Requirement lacks specificity
- MISSING_ACCEPTANCE_CRITERIA: No acceptance criteria
- UNMEASURABLE_SUCCESS_METRIC: Metric cannot be measured
- MISSING_USER_PERSONA: No user persona defined
- SECURITY_GAP: Security concern not addressed
- INCOMPLETE_ERROR_HANDLING: Error handling incomplete
- MISSING_API_CONTRACT: API not specified
- SCALABILITY_CONCERN: Scalability not addressed

### Confidence Values (0.0-1.0)
- 0.9+: Very confident in assessment
- 0.7-0.9: Confident
- 0.5-0.7: Somewhat confident
- <0.5: Low confidence, may need human review

Respond with ONLY the JSON, no additional text.`, content)

	return prompt
}

// evalResponse is the expected JSON structure from the LLM (v2).
type evalResponse struct {
	Categories []struct {
		ID          string   `json:"id"`
		IntScore    int      `json:"intScore"`
		Score       float64  `json:"score"` // Legacy support
		Confidence  float64  `json:"confidence"`
		ReasonCodes []string `json:"reasonCodes"`
		Explanation string   `json:"explanation"`
	} `json:"categories"`
	Findings []struct {
		Severity       string `json:"severity"`
		Category       string `json:"category"`
		Code           string `json:"code"`
		Title          string `json:"title"`
		Description    string `json:"description"`
		Recommendation string `json:"recommendation"`
		Evidence       string `json:"evidence,omitempty"`
		Location       string `json:"location,omitempty"`
	} `json:"findings"`
	Summary           string  `json:"summary"`
	OverallConfidence float64 `json:"overallConfidence"`
}

// parseEvalResponse parses the LLM response into a Result.
func parseEvalResponse(specType types.SpecType, rubricSet *rubric.RubricSet, response string, metadata JudgeMetadata) (*Result, error) {
	var resp evalResponse
	if err := json.Unmarshal([]byte(response), &resp); err != nil {
		return nil, fmt.Errorf("invalid JSON response: %w", err)
	}

	// Build category results and compute weighted score
	var categories []CategoryResult
	var totalIntScore float64
	var totalWeight float64
	var minConfidence float64 = 1.0

	for _, cat := range rubricSet.Categories {
		// Find matching category in response
		intScore := 3 // default to Acceptable
		var confidence float64 = 0.8
		var reasonCodes []rubric.ReasonCode
		var explanation string

		for _, respCat := range resp.Categories {
			if respCat.ID == cat.ID {
				intScore = respCat.IntScore
				// Fallback to legacy score if intScore not provided
				if intScore == 0 && respCat.Score > 0 {
					intScore = legacyScoreToIntScore(respCat.Score)
				}
				confidence = respCat.Confidence
				if confidence == 0 {
					confidence = 0.8 // default confidence
				}
				for _, code := range respCat.ReasonCodes {
					reasonCodes = append(reasonCodes, rubric.ReasonCode(code))
				}
				explanation = respCat.Explanation
				break
			}
		}

		// Clamp intScore to valid range
		if intScore < 1 {
			intScore = 1
		}
		if intScore > 5 {
			intScore = 5
		}

		// Track minimum confidence
		if confidence < minConfidence {
			minConfidence = confidence
		}

		categories = append(categories, CategoryResult{
			ID:          cat.ID,
			Name:        cat.Name,
			Score:       intScoreToLegacy(rubric.IntegerScore(intScore)),
			IntScore:    rubric.IntegerScore(intScore),
			Confidence:  confidence,
			ReasonCodes: reasonCodes,
			Weight:      cat.Weight,
			Explanation: explanation,
		})

		totalIntScore += float64(intScore) * cat.Weight
		totalWeight += cat.Weight
	}

	// Compute final integer score (weighted average, rounded)
	finalIntScore := rubric.ParseIntegerScore(int(totalIntScore/totalWeight + 0.5))

	// Convert findings
	var findings []Finding
	for _, f := range resp.Findings {
		findings = append(findings, Finding{
			Severity:       f.Severity,
			Category:       f.Category,
			Code:           rubric.ReasonCode(f.Code),
			Title:          f.Title,
			Description:    f.Description,
			Recommendation: f.Recommendation,
			Evidence:       f.Evidence,
			Location:       f.Location,
		})
	}

	// Use overall confidence from response or compute from categories
	confidence := resp.OverallConfidence
	if confidence == 0 {
		confidence = minConfidence
	}

	// Determine pass/fail using v2 criteria
	passed, blocking := evaluatePassCriteriaV2(finalIntScore, findings, rubricSet.PassCriteria)

	// Determine decision
	decision := "fail"
	if passed {
		decision = "pass"
	}

	return &Result{
		SchemaVersion: rubric.SchemaVersionV2,
		SpecType:      specType,
		Timestamp:     time.Now(),
		Score:         intScoreToLegacy(finalIntScore),
		IntScore:      finalIntScore,
		Confidence:    confidence,
		Passed:        passed,
		Blocking:      blocking,
		Categories:    categories,
		Findings:      findings,
		Decision:      decision,
		Summary:       resp.Summary,
		Judge:         metadata,
	}, nil
}

// legacyScoreToIntScore converts a 0-10 score to 1-5.
func legacyScoreToIntScore(score float64) int {
	switch {
	case score >= 9.0:
		return 5
	case score >= 7.0:
		return 4
	case score >= 5.0:
		return 3
	case score >= 3.0:
		return 2
	default:
		return 1
	}
}

// intScoreToLegacy converts a 1-5 score to legacy 0-10 scale.
func intScoreToLegacy(score rubric.IntegerScore) float64 {
	switch score {
	case rubric.ScoreExcellent:
		return 9.5
	case rubric.ScoreGood:
		return 8.0
	case rubric.ScoreAcceptable:
		return 6.5
	case rubric.ScoreMajorRevisions:
		return 4.0
	default:
		return 2.0
	}
}

// evaluatePassCriteriaV2 checks pass criteria using integer scores and returns blocking codes.
func evaluatePassCriteriaV2(score rubric.IntegerScore, findings []Finding, criteria rubric.RubricPassCriteria) (bool, []rubric.ReasonCode) {
	var blocking []rubric.ReasonCode

	// Score must be at least 3 (Acceptable) to pass
	if score < rubric.ScoreAcceptable {
		return false, blocking
	}

	// Finding limits by severity; a nil limit set blocks any critical/high.
	limits := criteria.MaxFindings
	if limits == nil {
		limits = &rubric.FindingLimits{Critical: 0, High: 0, Medium: -1, Low: -1}
	}

	// Count findings by severity and collect blocking codes
	var critical, high, medium int
	for _, f := range findings {
		switch f.Severity {
		case "critical":
			critical++
			if f.Code != "" {
				blocking = append(blocking, f.Code)
			}
		case "high":
			high++
			if f.Code != "" {
				blocking = append(blocking, f.Code)
			}
		case "medium":
			medium++
		}
	}

	// Check blocking thresholds (a negative limit means unlimited)
	if limits.Critical >= 0 && critical > limits.Critical {
		return false, blocking
	}
	if limits.High >= 0 && high > limits.High {
		return false, blocking
	}
	if limits.Medium >= 0 && medium > limits.Medium {
		return false, blocking
	}

	// If score is exactly Acceptable (3), check if there are any high-severity findings
	// that should block even though count thresholds aren't exceeded
	if score == rubric.ScoreAcceptable && len(blocking) > 0 {
		return false, blocking
	}

	return true, nil
}

// ToEvaluationReport converts the result to a structured-evaluation report.
// The rubricSet parameter is required for finalization.
func (r *Result) ToEvaluationReport(rubricSet *rubric.RubricSet) *rubric.Rubric {
	report := rubric.NewRubric(string(r.SpecType), "")

	// Set v2 fields
	report.SetIntScore(r.IntScore)
	report.SetConfidence(r.Confidence)
	report.SetPass(r.Passed)
	report.SetBlocking(r.Blocking)

	// Add category results with v2 fields
	for _, cat := range r.Categories {
		cr := rubric.NewCategoryResultWithIntScore(
			cat.ID,
			cat.IntScore,
			cat.Confidence,
			cat.Explanation,
		)
		cr.AddReasonCodes(cat.ReasonCodes...)
		report.AddCategoryResult(*cr)
	}

	// Add findings with v2 fields
	for _, f := range r.Findings {
		severity := rubric.SeverityMedium
		switch f.Severity {
		case "critical":
			severity = rubric.SeverityCritical
		case "high":
			severity = rubric.SeverityHigh
		case "medium":
			severity = rubric.SeverityMedium
		case "low":
			severity = rubric.SeverityLow
		case "info":
			severity = rubric.SeverityInfo
		}

		finding := rubric.Finding{
			Severity:       severity,
			Category:       f.Category,
			Code:           f.Code,
			Title:          f.Title,
			Description:    f.Description,
			Recommendation: f.Recommendation,
			Evidence:       f.Evidence,
			Location:       f.Location,
		}
		report.AddFinding(finding)
	}

	// Finalize with rubric
	report.Finalize(rubricSet, "visionspec eval")

	return report
}
