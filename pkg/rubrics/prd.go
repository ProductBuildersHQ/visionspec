//nolint:dupl // Rubric definitions are intentionally similar in structure
package rubrics

import "github.com/ProductBuildersHQ/visionspec/pkg/types"

func init() {
	Register(NewPRDRubricSet())
}

// NewPRDRubricSet creates the rubric set for Product Requirements Documents.
func NewPRDRubricSet() *RubricSet {
	return &RubricSet{
		SpecType:     types.SpecTypePRD,
		Name:         "PRD Evaluation",
		Description:  "Evaluates Product Requirements Documents for completeness, testability, and clarity",
		PassCriteria: DefaultPassCriteria(),
		Categories: []*Category{
			{
				ID:          "problem_definition",
				Name:        "Problem Definition",
				Description: "Clarity of the problem statement from user perspective",
				Weight:      0.15,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Problem is specific, measurable, and clearly impacts users",
					Partial: "Problem exists but lacks specificity",
					Fail:    "No clear problem statement",
				},
			},
			{
				ID:          "goals_scope",
				Name:        "Goals and Scope",
				Description: "Clarity of goals and explicit non-goals",
				Weight:      0.15,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "Goals are SMART, non-goals explicitly stated, scope is bounded",
					Partial: "Goals present but scope could be clearer",
					Fail:    "Goals and scope missing",
				},
			},
			{
				ID:          "user_stories",
				Name:        "User Stories",
				Description: "Quality and completeness of user stories",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "All stories follow standard format, are prioritized, with clear acceptance criteria",
					Partial: "Stories exist but some lack acceptance criteria",
					Fail:    "User stories missing or inadequate",
				},
			},
			{
				ID:          "functional_requirements",
				Name:        "Functional Requirements",
				Description: "Completeness and testability of functional requirements",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Requirements are complete, traceable to stories, and testable",
					Partial: "Requirements exist but some are vague",
					Fail:    "Functional requirements missing",
				},
			},
			{
				ID:          "nonfunctional_requirements",
				Name:        "Non-Functional Requirements",
				Description: "Coverage of performance, security, reliability requirements",
				Weight:      0.15,
				Required:    false,
				Criteria: CategoricalCriteria{
					Pass:    "Comprehensive NFRs with measurable targets for performance, security, accessibility",
					Partial: "Some NFRs present but incomplete coverage",
					Fail:    "NFRs missing",
				},
			},
			{
				ID:          "success_metrics",
				Name:        "Success Metrics",
				Description: "Definition of measurable success criteria",
				Weight:      0.10,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "SMART metrics with baselines, targets, and measurement methods",
					Partial: "Some metrics but lacking specificity",
					Fail:    "Success metrics missing",
				},
			},
			{
				ID:          "dependencies",
				Name:        "Dependencies and Risks",
				Description: "Identification of dependencies and open questions",
				Weight:      0.05,
				Required:    false,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "Dependencies mapped with owners, status, and risk mitigation",
					Partial: "Some dependencies noted",
					Fail:    "Dependencies not addressed",
				},
			},
		},
	}
}
