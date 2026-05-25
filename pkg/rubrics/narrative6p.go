//nolint:dupl // Rubric definitions are intentionally similar in structure
package rubrics

import "github.com/ProductBuildersHQ/visionspec/pkg/types"

func init() {
	Register(NewNarrative6PRubricSet())
}

// NewNarrative6PRubricSet creates the rubric set for AWS-Style 6-Pager Narratives.
func NewNarrative6PRubricSet() *RubricSet {
	return &RubricSet{
		SpecType:     types.SpecTypeNarrative6P,
		Name:         "Narrative 6-Pager Evaluation",
		Description:  "Evaluates AWS-style 6-pager narratives for depth, rigor, and completeness",
		PassCriteria: DefaultPassCriteria(),
		Categories: []*Category{
			{
				ID:          "tenets_clarity",
				Name:        "Tenets and Introduction",
				Description: "Quality of guiding principles and context setting",
				Weight:      0.10,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Clear, actionable tenets that guide tradeoff decisions",
					Partial: "Tenets present but generic or not actionable",
					Fail:    "No tenets or introduction",
				},
			},
			{
				ID:          "customer_depth",
				Name:        "Customer Problem Depth",
				Description: "How deeply the customer and their problem are understood",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Specific customer with vivid problem description, quotes, and anecdotes",
					Partial: "Customer and problem described but lacks depth",
					Fail:    "No meaningful customer or problem description",
				},
			},
			{
				ID:          "solution_narrative",
				Name:        "Solution Narrative",
				Description: "Quality of the customer-centric solution description",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Compelling customer journey with clear capabilities and scope",
					Partial: "Solution described but too technical or feature-focused",
					Fail:    "No coherent solution narrative",
				},
			},
			{
				ID:          "timing_justification",
				Name:        "Why Now Justification",
				Description: "Strength of the timing and urgency argument",
				Weight:      0.15,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Compelling market, competitive, and internal readiness arguments",
					Partial: "Timing argument present but weak",
					Fail:    "No timing justification or convincing reason for urgency",
				},
			},
			{
				ID:          "business_case",
				Name:        "Business Case Rigor",
				Description: "Quality of customer impact, business metrics, and ROI analysis",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "Quantified customer and business impact with clear assumptions",
					Partial: "Some metrics but lacking rigor or assumptions",
					Fail:    "No business case or weak unsupported claims",
				},
			},
			{
				ID:          "risks_dependencies",
				Name:        "Risks and Dependencies",
				Description: "Thoroughness of risk identification and mitigation",
				Weight:      0.15,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Comprehensive risks with likelihood, impact, and mitigations",
					Partial: "Some risks listed but incomplete analysis",
					Fail:    "No risk assessment",
				},
			},
		},
	}
}
