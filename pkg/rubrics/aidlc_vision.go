//nolint:dupl // Rubric definitions are intentionally similar in structure
package rubrics

import "github.com/ProductBuildersHQ/visionspec/pkg/types"

func init() {
	Register(NewAIDLCVisionRubricSet())
}

// NewAIDLCVisionRubricSet creates the rubric set for AIDLC Vision Documents.
func NewAIDLCVisionRubricSet() *RubricSet {
	return &RubricSet{
		SpecType:     types.SpecTypeAIDLCVision,
		Name:         "AIDLC Vision Document Evaluation",
		Description:  "Evaluates AIDLC Vision Documents for clarity, stakeholder alignment, and strategic direction",
		PassCriteria: DefaultPassCriteria(),
		Categories: []*Category{
			{
				ID:          "clarity",
				Name:        "Vision Clarity",
				Description: "Clarity and specificity of the vision statement",
				Weight:      0.25,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Vision is clear, inspiring, and articulates a specific future state with measurable outcomes",
					Partial: "Vision exists but lacks specificity or measurable outcomes",
					Fail:    "Vision is missing, vague, or generic without clear direction",
				},
			},
			{
				ID:          "completeness",
				Name:        "Completeness",
				Description: "Coverage of all required vision components",
				Weight:      0.25,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Document covers problem statement, target users, value proposition, and success criteria comprehensively",
					Partial: "Some required components present but gaps exist",
					Fail:    "Major components missing or incomplete",
				},
			},
			{
				ID:          "alignment",
				Name:        "Stakeholder Alignment",
				Description: "Evidence of stakeholder input and alignment",
				Weight:      0.25,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Clear stakeholder identification, documented input, and explicit alignment with organizational goals",
					Partial: "Stakeholders identified but limited evidence of alignment",
					Fail:    "No stakeholder alignment documented",
				},
			},
			{
				ID:          "feasibility",
				Name:        "Feasibility Assessment",
				Description: "Realistic assessment of technical and business feasibility",
				Weight:      0.25,
				Required:    false,
				Criteria: CategoricalCriteria{
					Pass:    "Feasibility addressed with clear constraints, risks, and mitigation strategies",
					Partial: "Basic feasibility mentioned but lacks depth",
					Fail:    "No feasibility assessment included",
				},
			},
		},
	}
}
