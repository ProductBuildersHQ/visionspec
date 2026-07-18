//nolint:dupl // Rubric definitions are intentionally similar in structure
package rubrics

import "github.com/ProductBuildersHQ/visionspec/pkg/types"

func init() {
	Register(NewAIDLCTechnicalRubricSet())
}

// NewAIDLCTechnicalRubricSet creates the rubric set for AIDLC Technical Specs.
func NewAIDLCTechnicalRubricSet() *RubricSet {
	return &RubricSet{
		SpecType:     types.SpecTypeAIDLCTechnical,
		Name:         "AIDLC Technical Specification Evaluation",
		Description:  "Evaluates AIDLC Technical Specifications for architecture soundness, scalability, and implementation clarity",
		PassCriteria: DefaultPassCriteria(),
		Categories: []*Category{
			{
				ID:          "architecture",
				Name:        "Architecture Soundness",
				Description: "Technical architecture is well-designed and appropriate",
				Weight:      0.25,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Architecture follows established patterns, addresses scalability, and includes clear component interactions",
					Partial: "Architecture exists but has gaps in design or scalability considerations",
					Fail:    "Architecture is missing or fundamentally flawed",
				},
			},
			{
				ID:          "scalability",
				Name:        "Scalability",
				Description: "Design addresses growth and scale requirements",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Scalability explicitly addressed with capacity planning, bottleneck analysis, and growth strategies",
					Partial: "Some scalability considerations present but incomplete",
					Fail:    "No scalability considerations in design",
				},
			},
			{
				ID:          "implementation_clarity",
				Name:        "Implementation Clarity",
				Description: "Technical details sufficient for implementation",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "APIs, data models, and interfaces clearly defined with sufficient detail for implementation",
					Partial: "Some technical details present but gaps exist",
					Fail:    "Insufficient detail for implementation",
				},
			},
			{
				ID:          "technology_choices",
				Name:        "Technology Choices",
				Description: "Technology selections are justified and appropriate",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Technology choices documented with clear rationale, alternatives considered, and trade-offs explained",
					Partial: "Technologies mentioned but rationale incomplete",
					Fail:    "Technology choices not justified or inappropriate",
				},
			},
			{
				ID:          "dependencies",
				Name:        "Dependencies",
				Description: "External and internal dependencies identified",
				Weight:      0.15,
				Required:    false,
				Criteria: CategoricalCriteria{
					Pass:    "All dependencies documented with version requirements, fallback strategies, and integration points",
					Partial: "Dependencies listed but incomplete information",
					Fail:    "Dependencies not identified",
				},
			},
		},
	}
}
