//nolint:dupl // Rubric definitions are intentionally similar in structure
package rubrics

import "github.com/ProductBuildersHQ/visionspec/pkg/types"

func init() {
	Register(NewAIDLCImplementationRubricSet())
}

// NewAIDLCImplementationRubricSet creates the rubric set for AIDLC Implementation Plans.
func NewAIDLCImplementationRubricSet() *RubricSet {
	return &RubricSet{
		SpecType:     types.SpecTypeAIDLCImplementation,
		Name:         "AIDLC Implementation Plan Evaluation",
		Description:  "Evaluates AIDLC Implementation Plans for task breakdown, dependencies, and risk mitigation",
		PassCriteria: DefaultPassCriteria(),
		Categories: []*Category{
			{
				ID:          "task_breakdown",
				Name:        "Task Breakdown",
				Description: "Granularity and clarity of implementation tasks",
				Weight:      0.25,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Tasks broken down to actionable units with clear scope, deliverables, and completion criteria",
					Partial: "Tasks exist but some are too large or unclear",
					Fail:    "Task breakdown missing or too high-level",
				},
			},
			{
				ID:          "dependencies",
				Name:        "Dependencies",
				Description: "Task dependencies and critical path identification",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Dependencies clearly mapped with critical path identified and blocking relationships documented",
					Partial: "Some dependencies noted but incomplete mapping",
					Fail:    "Dependencies not identified",
				},
			},
			{
				ID:          "risk_mitigation",
				Name:        "Risk Mitigation",
				Description: "Identification and mitigation of implementation risks",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Risks identified with probability, impact assessment, and concrete mitigation strategies",
					Partial: "Some risks noted but mitigation incomplete",
					Fail:    "No risk assessment included",
				},
			},
			{
				ID:          "resource_allocation",
				Name:        "Resource Allocation",
				Description: "Team and resource assignment to tasks",
				Weight:      0.15,
				Required:    false,
				Criteria: CategoricalCriteria{
					Pass:    "Resources assigned with skills mapping and capacity planning",
					Partial: "Some resource allocation present",
					Fail:    "No resource allocation defined",
				},
			},
			{
				ID:          "milestones",
				Name:        "Milestones",
				Description: "Key milestones and checkpoints defined",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Clear milestones with measurable success criteria and review gates",
					Partial: "Milestones exist but lack clear criteria",
					Fail:    "No milestones defined",
				},
			},
		},
	}
}
