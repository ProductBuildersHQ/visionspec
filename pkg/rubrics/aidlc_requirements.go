//nolint:dupl // Rubric definitions are intentionally similar in structure
package rubrics

import "github.com/ProductBuildersHQ/visionspec/pkg/types"

func init() {
	Register(NewAIDLCRequirementsRubricSet())
}

// NewAIDLCRequirementsRubricSet creates the rubric set for AIDLC Requirements Specs.
func NewAIDLCRequirementsRubricSet() *RubricSet {
	return &RubricSet{
		SpecType:     types.SpecTypeAIDLCRequirements,
		Name:         "AIDLC Requirements Specification Evaluation",
		Description:  "Evaluates AIDLC Requirements Specifications for testability, traceability, and completeness",
		PassCriteria: DefaultPassCriteria(),
		Categories: []*Category{
			{
				ID:          "testability",
				Name:        "Testability",
				Description: "Requirements are verifiable and testable",
				Weight:      0.25,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "All requirements have clear acceptance criteria with specific, measurable test conditions",
					Partial: "Most requirements testable but some lack clear acceptance criteria",
					Fail:    "Requirements are vague and not testable",
				},
			},
			{
				ID:          "traceability",
				Name:        "Traceability",
				Description: "Requirements link to vision and can be traced to implementation",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "Each requirement traces to vision goals with unique identifiers for tracking",
					Partial: "Some traceability exists but inconsistent",
					Fail:    "No traceability to vision or other documents",
				},
			},
			{
				ID:          "no_ambiguity",
				Name:        "Clarity",
				Description: "Requirements are unambiguous and clearly stated",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "Requirements use precise language, avoid ambiguous terms, and have single interpretation",
					Partial: "Some requirements have minor ambiguities",
					Fail:    "Requirements contain significant ambiguities or multiple interpretations",
				},
			},
			{
				ID:          "completeness",
				Name:        "Completeness",
				Description: "Coverage of functional and non-functional requirements",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "Comprehensive coverage of functional, non-functional, and edge case requirements",
					Partial: "Core requirements present but gaps in edge cases or NFRs",
					Fail:    "Major requirement categories missing",
				},
			},
			{
				ID:          "prioritization",
				Name:        "Prioritization",
				Description: "Requirements are prioritized with clear rationale",
				Weight:      0.15,
				Required:    false,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "All requirements prioritized (MoSCoW or similar) with documented rationale",
					Partial: "Some prioritization exists but incomplete",
					Fail:    "No prioritization of requirements",
				},
			},
		},
	}
}
