//nolint:dupl // Rubric definitions are intentionally similar in structure
package rubrics

import "github.com/ProductBuildersHQ/visionspec/pkg/types"

func init() {
	Register(NewTRDRubricSet())
}

// NewTRDRubricSet creates the rubric set for Technical Requirements Documents.
func NewTRDRubricSet() *RubricSet {
	return &RubricSet{
		SpecType:     types.SpecTypeTRD,
		Name:         "TRD Evaluation",
		Description:  "Evaluates Technical Requirements Documents for completeness, clarity, and implementability",
		PassCriteria: DefaultPassCriteria(),
		Categories: []*Category{
			{
				ID:          "architecture_clarity",
				Name:        "Architecture Clarity",
				Description: "How clearly the system architecture is described",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Clear architecture with diagrams, component descriptions, and data flow",
					Partial: "Basic architecture present but lacks detail",
					Fail:    "No clear architecture described",
				},
			},
			{
				ID:          "prd_traceability",
				Name:        "PRD Traceability",
				Description: "Whether technical design traces back to PRD requirements",
				Weight:      0.15,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "Every PRD requirement has corresponding technical design",
					Partial: "Some traceability but gaps exist",
					Fail:    "No traceability to requirements",
				},
			},
			{
				ID:          "api_design",
				Name:        "API Design",
				Description: "Quality of API/interface definitions",
				Weight:      0.15,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "Complete API spec with request/response schemas and error handling",
					Partial: "Basic API definition present",
					Fail:    "No API design",
				},
			},
			{
				ID:          "data_design",
				Name:        "Data Design",
				Description: "Quality of data models and storage design",
				Weight:      0.15,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "Complete data models with relationships, storage rationale, and migration plan",
					Partial: "Basic data design present",
					Fail:    "No data design",
				},
			},
			{
				ID:          "nfr_coverage",
				Name:        "NFR Coverage",
				Description: "Coverage of non-functional requirements (performance, security, etc.)",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Comprehensive NFRs with specific targets and measurement plans",
					Partial: "Some NFRs but gaps in coverage",
					Fail:    "No NFR coverage",
				},
			},
			{
				ID:          "implementability",
				Name:        "Implementability",
				Description: "Whether the design is implementable by an engineer",
				Weight:      0.15,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Design is actionable with clear path to implementation",
					Partial: "Design needs some clarification before implementation",
					Fail:    "Cannot be implemented from this document",
				},
			},
		},
	}
}
