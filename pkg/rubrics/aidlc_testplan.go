//nolint:dupl // Rubric definitions are intentionally similar in structure
package rubrics

import "github.com/ProductBuildersHQ/visionspec/pkg/types"

func init() {
	Register(NewAIDLCTestPlanRubricSet())
}

// NewAIDLCTestPlanRubricSet creates the rubric set for AIDLC Test Plans.
func NewAIDLCTestPlanRubricSet() *RubricSet {
	return &RubricSet{
		SpecType:     types.SpecTypeAIDLCTestPlan,
		Name:         "AIDLC Test Plan Evaluation",
		Description:  "Evaluates AIDLC Test Plans for coverage, edge cases, and automation readiness",
		PassCriteria: DefaultPassCriteria(),
		Categories: []*Category{
			{
				ID:          "coverage",
				Name:        "Test Coverage",
				Description: "Completeness of test coverage across requirements",
				Weight:      0.25,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "Test cases map to all requirements with clear traceability and coverage targets defined",
					Partial: "Some coverage exists but gaps in requirement mapping",
					Fail:    "Test coverage incomplete or no traceability",
				},
			},
			{
				ID:          "edge_cases",
				Name:        "Edge Cases",
				Description: "Coverage of boundary conditions and error scenarios",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "Edge cases, boundary conditions, and error scenarios thoroughly documented with test cases",
					Partial: "Some edge cases covered but not comprehensive",
					Fail:    "Edge cases not addressed",
				},
			},
			{
				ID:          "automation_readiness",
				Name:        "Automation Readiness",
				Description: "Plan for test automation and CI/CD integration",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "Clear automation strategy with tools, frameworks, and CI/CD integration plan",
					Partial: "Automation mentioned but strategy incomplete",
					Fail:    "No automation plan included",
				},
			},
			{
				ID:          "test_data",
				Name:        "Test Data",
				Description: "Test data strategy and management",
				Weight:      0.15,
				Required:    false,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "Test data requirements defined with generation strategy and data management plan",
					Partial: "Some test data considerations present",
					Fail:    "Test data not addressed",
				},
			},
			{
				ID:          "environments",
				Name:        "Test Environments",
				Description: "Test environment requirements and setup",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "Test environments clearly defined with setup procedures and environment parity requirements",
					Partial: "Environments mentioned but incomplete details",
					Fail:    "Test environments not addressed",
				},
			},
		},
	}
}
