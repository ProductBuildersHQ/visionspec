// Package eval provides evaluation orchestration for spec documents.
package eval

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/plexusone/multispec/pkg/rubrics"
	"github.com/plexusone/multispec/pkg/types"
	"github.com/plexusone/structured-evaluation/evaluation"
)

// Result represents the outcome of an evaluation.
type Result struct {
	SpecType   types.SpecType      `json:"spec_type"`
	Timestamp  time.Time           `json:"timestamp"`
	Score      float64             `json:"score"`
	Passed     bool                `json:"passed"`
	Categories []CategoryResult    `json:"categories"`
	Findings   []Finding           `json:"findings"`
	Decision   string              `json:"decision"`
	Summary    string              `json:"summary"`
	Judge      JudgeMetadata       `json:"judge"`
}

// CategoryResult contains the evaluation result for a category.
type CategoryResult struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Score       float64 `json:"score"`
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
	llm *LLMClient
}

// NewEvaluator creates a new evaluator with the given LLM client.
func NewEvaluator(llm *LLMClient) *Evaluator {
	return &Evaluator{llm: llm}
}

// Evaluate runs evaluation on content against the rubric for the given spec type.
func (e *Evaluator) Evaluate(ctx context.Context, specType types.SpecType, content string) (*Result, error) {
	// Get rubric
	rubricSet, err := rubrics.Get(specType)
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
func buildEvalPrompt(rubricSet *rubrics.RubricSet, content string) string {
	prompt := fmt.Sprintf(`You are an expert document evaluator. Evaluate the following %s against the provided rubric.

## Rubric

Evaluate each category on a scale of 0-10:

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
      "score": 7.5,
      "explanation": "Brief explanation of the score"
    }
  ],
  "findings": [
    {
      "severity": "critical|high|medium|low|info",
      "category": "category_id",
      "title": "Short title",
      "description": "Detailed description",
      "recommendation": "How to fix"
    }
  ],
  "summary": "Overall assessment in 2-3 sentences"
}

Severity levels:
- critical: Fundamental issues that block approval
- high: Significant issues that should be fixed
- medium: Notable issues worth addressing
- low: Minor improvements
- info: Informational observations

Respond with ONLY the JSON, no additional text.`, content)

	return prompt
}

// evalResponse is the expected JSON structure from the LLM.
type evalResponse struct {
	Categories []struct {
		ID          string  `json:"id"`
		Score       float64 `json:"score"`
		Explanation string  `json:"explanation"`
	} `json:"categories"`
	Findings []Finding `json:"findings"`
	Summary  string    `json:"summary"`
}

// parseEvalResponse parses the LLM response into a Result.
func parseEvalResponse(specType types.SpecType, rubricSet *rubrics.RubricSet, response string, metadata JudgeMetadata) (*Result, error) {
	var resp evalResponse
	if err := json.Unmarshal([]byte(response), &resp); err != nil {
		return nil, fmt.Errorf("invalid JSON response: %w", err)
	}

	// Build category results and compute weighted score
	var categories []CategoryResult
	var totalScore float64
	var totalWeight float64

	for _, cat := range rubricSet.Categories {
		// Find matching category in response
		score := 5.0 // default
		var explanation string
		for _, respCat := range resp.Categories {
			if respCat.ID == cat.ID {
				score = respCat.Score
				explanation = respCat.Explanation
				break
			}
		}

		categories = append(categories, CategoryResult{
			ID:          cat.ID,
			Name:        cat.Name,
			Score:       score,
			Weight:      cat.Weight,
			Explanation: explanation,
		})

		totalScore += score * cat.Weight
		totalWeight += cat.Weight
	}

	// Compute final score
	finalScore := totalScore / totalWeight

	// Determine pass/fail based on rubric criteria
	passed := evaluatePassCriteria(finalScore, resp.Findings, rubricSet.PassCriteria)

	// Determine decision
	decision := "fail"
	if passed {
		decision = "pass"
	}

	return &Result{
		SpecType:   specType,
		Timestamp:  time.Now(),
		Score:      finalScore,
		Passed:     passed,
		Categories: categories,
		Findings:   resp.Findings,
		Decision:   decision,
		Summary:    resp.Summary,
		Judge:      metadata,
	}, nil
}

// evaluatePassCriteria checks if the evaluation passes based on criteria.
func evaluatePassCriteria(score float64, findings []Finding, criteria rubrics.PassCriteria) bool {
	if score < criteria.MinScore {
		return false
	}

	// Count findings by severity
	var critical, high, medium int
	for _, f := range findings {
		switch f.Severity {
		case "critical":
			critical++
		case "high":
			high++
		case "medium":
			medium++
		}
	}

	if critical > criteria.MaxCritical {
		return false
	}
	if high > criteria.MaxHigh {
		return false
	}
	if criteria.MaxMedium >= 0 && medium > criteria.MaxMedium {
		return false
	}

	return true
}

// ToEvaluationReport converts the result to a structured-evaluation report.
func (r *Result) ToEvaluationReport() *evaluation.EvaluationReport {
	report := evaluation.NewEvaluationReport(string(r.SpecType), "")

	// Add category scores
	for _, cat := range r.Categories {
		report.AddCategory(evaluation.NewCategoryScore(cat.ID, cat.Weight, cat.Score, cat.Explanation))
	}

	// Add findings
	for _, f := range r.Findings {
		severity := evaluation.SeverityMedium
		switch f.Severity {
		case "critical":
			severity = evaluation.SeverityCritical
		case "high":
			severity = evaluation.SeverityHigh
		case "medium":
			severity = evaluation.SeverityMedium
		case "low":
			severity = evaluation.SeverityLow
		case "info":
			severity = evaluation.SeverityInfo
		}

		report.AddFinding(evaluation.Finding{
			Severity:       severity,
			Category:       f.Category,
			Title:          f.Title,
			Description:    f.Description,
			Recommendation: f.Recommendation,
			Evidence:       f.Evidence,
		})
	}

	// Finalize
	report.Finalize("multispec eval")

	return report
}
