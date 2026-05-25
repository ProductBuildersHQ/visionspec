//nolint:dupl // Rubric definitions are intentionally similar in structure
package rubrics

import "github.com/ProductBuildersHQ/visionspec/pkg/types"

func init() {
	Register(NewNarrative1PRubricSet())
}

// NewNarrative1PRubricSet creates the rubric set for 1-Pager Executive Narratives.
func NewNarrative1PRubricSet() *RubricSet {
	return &RubricSet{
		SpecType:     types.SpecTypeNarrative1P,
		Name:         "Narrative 1-Pager Evaluation",
		Description:  "Evaluates Executive 1-Pager narratives for clarity, impact, and completeness",
		PassCriteria: DefaultPassCriteria(),
		Categories: []*Category{
			{
				ID:          "opportunity_clarity",
				Name:        "Opportunity Clarity",
				Description: "How clearly the market opportunity or problem is articulated",
				Weight:      0.25,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Compelling opportunity with clear market timing and urgency",
					Partial: "Opportunity stated but lacks compelling urgency",
					Fail:    "No clear opportunity articulated or poorly justified",
				},
			},
			{
				ID:          "solution_value",
				Name:        "Solution Value Proposition",
				Description: "Clarity of the solution and its value proposition",
				Weight:      0.25,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Crystal clear value proposition with direct link to opportunity",
					Partial: "Solution described but value unclear",
					Fail:    "No clear solution articulated or disconnected from problem",
				},
			},
			{
				ID:          "target_specificity",
				Name:        "Target Customer Specificity",
				Description: "How specifically the target customer is defined",
				Weight:      0.15,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Specific, well-defined customer segment with clear characteristics",
					Partial: "Customer mentioned but too broad",
					Fail:    "No target customer defined or definition is vague",
				},
			},
			{
				ID:          "benefits_concrete",
				Name:        "Concrete Benefits",
				Description: "Whether benefits are outcome-focused and concrete",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "Benefits are specific, measurable outcomes customers care about",
					Partial: "Benefits listed but feature-focused rather than outcome-focused",
					Fail:    "No clear benefits stated or benefits are vague",
				},
			},
			{
				ID:          "ask_clarity",
				Name:        "Clear Ask",
				Description: "Whether the ask (resources, decisions needed) is clear",
				Weight:      0.15,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Specific ask with clear next steps and timeline",
					Partial: "Ask present but vague on specifics",
					Fail:    "No ask or call to action",
				},
			},
		},
	}
}
