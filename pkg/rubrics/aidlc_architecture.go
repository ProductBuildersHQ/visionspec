//nolint:dupl // Rubric definitions are intentionally similar in structure
package rubrics

import "github.com/ProductBuildersHQ/visionspec/pkg/types"

func init() {
	Register(NewAIDLCArchitectureRubricSet())
}

// NewAIDLCArchitectureRubricSet creates the rubric set for AIDLC Architecture Specs.
func NewAIDLCArchitectureRubricSet() *RubricSet {
	return &RubricSet{
		SpecType:     types.SpecTypeAIDLCArchitecture,
		Name:         "AIDLC Architecture Specification Evaluation",
		Description:  "Evaluates AIDLC Architecture Specifications for system design, component interactions, and deployment topology",
		PassCriteria: DefaultPassCriteria(),
		Categories: []*Category{
			{
				ID:          "system_design",
				Name:        "System Design",
				Description: "Overall system architecture and component breakdown",
				Weight:      0.25,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "Clear system decomposition with well-defined boundaries, responsibilities, and interaction patterns",
					Partial: "System design exists but component boundaries unclear",
					Fail:    "No clear system design or architecture",
				},
			},
			{
				ID:          "data_architecture",
				Name:        "Data Architecture",
				Description: "Data flow, storage, and consistency approach",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "Data models, storage strategies, consistency requirements, and data flow patterns clearly documented",
					Partial: "Some data architecture present but incomplete",
					Fail:    "Data architecture not addressed",
				},
			},
			{
				ID:          "integration_patterns",
				Name:        "Integration Patterns",
				Description: "API design and service integration approach",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "Integration patterns documented with API contracts, protocols, and error handling strategies",
					Partial: "Integration approach mentioned but lacks detail",
					Fail:    "No integration patterns defined",
				},
			},
			{
				ID:          "deployment_topology",
				Name:        "Deployment Topology",
				Description: "Infrastructure and deployment architecture",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "Deployment topology clearly documented with infrastructure requirements, networking, and scaling approach",
					Partial: "Basic deployment info present but gaps exist",
					Fail:    "Deployment topology not addressed",
				},
			},
			{
				ID:          "diagrams",
				Name:        "Architecture Diagrams",
				Description: "Visual representations of the architecture",
				Weight:      0.15,
				Required:    false,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "Clear architecture diagrams (C4, sequence, etc.) that accurately represent the system",
					Partial: "Some diagrams present but incomplete or unclear",
					Fail:    "No architecture diagrams included",
				},
			},
		},
	}
}
