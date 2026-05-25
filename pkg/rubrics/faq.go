//nolint:dupl // Rubric definitions are intentionally similar in structure
package rubrics

import "github.com/ProductBuildersHQ/visionspec/pkg/types"

func init() {
	Register(NewFAQRubricSet())
}

// NewFAQRubricSet creates the rubric set for FAQ documents.
func NewFAQRubricSet() *RubricSet {
	return &RubricSet{
		SpecType:     types.SpecTypeFAQ,
		Name:         "FAQ Evaluation",
		Description:  "Evaluates FAQ documents for completeness, clarity, and anticipation of customer questions",
		PassCriteria: DefaultPassCriteria(),
		Categories: []*Category{
			{
				ID:          "question_coverage",
				Name:        "Question Coverage",
				Description: "Whether the FAQ anticipates the questions customers will actually ask",
				Weight:      0.25,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Comprehensive coverage of likely questions across all topics",
					Partial: "Basic questions covered but gaps exist",
					Fail:    "Inadequate question coverage",
				},
			},
			{
				ID:          "answer_clarity",
				Name:        "Answer Clarity",
				Description: "How clearly and directly questions are answered",
				Weight:      0.25,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Direct, clear answers with no ambiguity",
					Partial: "Answers present but could be clearer",
					Fail:    "Answers don't address questions",
				},
			},
			{
				ID:          "customer_language",
				Name:        "Customer Language",
				Description: "Whether questions and answers use customer language, not internal jargon",
				Weight:      0.15,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Questions phrased as customers would ask, answers in plain language",
					Partial: "Some jargon but generally accessible",
					Fail:    "Written for internal audience, not customers",
				},
			},
			{
				ID:          "pricing_transparency",
				Name:        "Pricing Transparency",
				Description: "Clarity of pricing and cost-related information",
				Weight:      0.15,
				Required:    false,
				Criteria: CategoricalCriteria{
					Pass:    "Clear pricing with examples and cost estimation guidance",
					Partial: "Basic pricing present but lacks detail",
					Fail:    "No pricing information",
				},
			},
			{
				ID:          "getting_started",
				Name:        "Getting Started Clarity",
				Description: "How well the path to adoption is explained",
				Weight:      0.10,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "Clear step-by-step path from interest to success",
					Partial: "Basic getting started info present",
					Fail:    "No getting started guidance",
				},
			},
			{
				ID:          "objection_handling",
				Name:        "Objection Handling",
				Description: "Whether the FAQ addresses likely concerns and objections",
				Weight:      0.10,
				Required:    false,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "Proactively addresses concerns with honest, confident responses",
					Partial: "Some objections covered",
					Fail:    "No objection handling",
				},
			},
		},
	}
}
