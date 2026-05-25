//nolint:dupl // Rubric definitions are intentionally similar in structure
package rubrics

import "github.com/ProductBuildersHQ/visionspec/pkg/types"

func init() {
	Register(NewMRDRubricSet())
}

// NewMRDRubricSet creates the rubric set for Market Requirements Documents.
func NewMRDRubricSet() *RubricSet {
	return &RubricSet{
		SpecType:     types.SpecTypeMRD,
		Name:         "MRD Evaluation",
		Description:  "Evaluates Market Requirements Documents for completeness, clarity, and market insight",
		PassCriteria: DefaultPassCriteria(),
		Categories: []*Category{
			{
				ID:          "problem_statement",
				Name:        "Problem Statement",
				Description: "Clarity and relevance of the market problem being addressed",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Problem is clearly articulated with quantified impact, timing justification, and cost of inaction",
					Partial: "Problem is stated but lacks depth or supporting evidence",
					Fail:    "Problem statement is missing or unclear",
				},
			},
			{
				ID:          "target_market",
				Name:        "Target Market",
				Description: "Definition and sizing of target market segments",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "Clear TAM/SAM/SOM with detailed segment profiles and sizing methodology",
					Partial: "Basic market definition present but incomplete",
					Fail:    "Market definition missing or completely inadequate",
				},
			},
			{
				ID:          "competitive_analysis",
				Name:        "Competitive Analysis",
				Description: "Understanding of competitive landscape and differentiation",
				Weight:      0.15,
				Required:    false,
				Criteria: CategoricalCriteria{
					Pass:    "Comprehensive analysis of direct/indirect competitors with clear differentiation strategy",
					Partial: "Basic competitor list but limited analysis",
					Fail:    "Competitive analysis missing",
				},
			},
			{
				ID:          "requirements_clarity",
				Name:        "Requirements Clarity",
				Description: "Clarity and prioritization of market requirements",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Requirements are prioritized (MoSCoW), traceable, and well-rationalized",
					Partial: "Requirements listed but prioritization is weak",
					Fail:    "Requirements section is inadequate",
				},
			},
			{
				ID:          "business_alignment",
				Name:        "Business Alignment",
				Description: "Alignment with business goals and measurable success metrics",
				Weight:      0.15,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Clear strategic alignment with quantified success metrics and timelines",
					Partial: "Business alignment stated but metrics are weak",
					Fail:    "No clear business alignment",
				},
			},
			{
				ID:          "risk_assessment",
				Name:        "Risk Assessment",
				Description: "Identification of constraints, assumptions, and risks",
				Weight:      0.10,
				Required:    false,
				Criteria: CategoricalCriteria{
					Pass:    "Comprehensive risk matrix with impact, probability, and mitigation strategies",
					Partial: "Risks listed but mitigation is incomplete",
					Fail:    "Risk assessment missing",
				},
			},
		},
	}
}
