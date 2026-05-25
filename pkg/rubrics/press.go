//nolint:dupl // Rubric definitions are intentionally similar in structure
package rubrics

import "github.com/ProductBuildersHQ/visionspec/pkg/types"

func init() {
	Register(NewPressRubricSet())
}

// NewPressRubricSet creates the rubric set for Press Release documents.
func NewPressRubricSet() *RubricSet {
	return &RubricSet{
		SpecType:     types.SpecTypePress,
		Name:         "Press Release Evaluation",
		Description:  "Evaluates Press Release documents for clarity, customer focus, and completeness",
		PassCriteria: DefaultPassCriteria(),
		Categories: []*Category{
			{
				ID:          "headline_impact",
				Name:        "Headline Impact",
				Description: "Whether the headline clearly communicates the primary customer benefit",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Compelling headline with clear customer benefit and urgency",
					Partial: "Headline present but benefit unclear",
					Fail:    "No clear headline or irrelevant",
				},
			},
			{
				ID:          "customer_problem",
				Name:        "Customer Problem",
				Description: "How well the customer problem is articulated with authentic voice",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Vivid problem description with authentic customer quotes",
					Partial: "Problem stated but lacks customer voice",
					Fail:    "No clear customer problem",
				},
			},
			{
				ID:          "solution_clarity",
				Name:        "Solution Clarity",
				Description: "Whether the solution and benefits are clearly explained",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Crystal clear solution with concrete, measurable benefits",
					Partial: "Solution described but benefits vague",
					Fail:    "No clear solution articulated",
				},
			},
			{
				ID:          "customer_validation",
				Name:        "Customer Validation",
				Description: "Quality of customer quotes and validation",
				Weight:      0.15,
				Required:    false,
				Criteria: CategoricalCriteria{
					Pass:    "Authentic customer quotes with specific success stories",
					Partial: "Generic quotes without specificity",
					Fail:    "No customer validation",
				},
			},
			{
				ID:          "call_to_action",
				Name:        "Call to Action",
				Description: "Clarity of availability, pricing, and next steps",
				Weight:      0.15,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Clear availability, pricing, and specific action steps",
					Partial: "CTA present but vague",
					Fail:    "No call to action",
				},
			},
			{
				ID:          "readability",
				Name:        "Readability",
				Description: "Whether the press release is written in plain, accessible language",
				Weight:      0.10,
				Required:    false,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "Clear, jargon-free prose accessible to general audience",
					Partial: "Somewhat technical but understandable",
					Fail:    "Incomprehensible to target audience",
				},
			},
		},
	}
}
