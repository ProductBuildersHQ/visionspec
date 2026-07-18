//nolint:dupl // Rubric definitions are intentionally similar in structure
package rubrics

import "github.com/ProductBuildersHQ/visionspec/pkg/types"

func init() {
	Register(NewAIDLCRunbookRubricSet())
}

// NewAIDLCRunbookRubricSet creates the rubric set for AIDLC Runbooks.
func NewAIDLCRunbookRubricSet() *RubricSet {
	return &RubricSet{
		SpecType:     types.SpecTypeAIDLCRunbook,
		Name:         "AIDLC Runbook Evaluation",
		Description:  "Evaluates AIDLC Runbooks for operational clarity, rollback procedures, and troubleshooting guidance",
		PassCriteria: DefaultPassCriteria(),
		Categories: []*Category{
			{
				ID:          "step_clarity",
				Name:        "Step Clarity",
				Description: "Operational procedures are clear and unambiguous",
				Weight:      0.25,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Step-by-step procedures with clear commands, expected outputs, and decision points",
					Partial: "Procedures exist but some steps unclear",
					Fail:    "Procedures missing or too vague to follow",
				},
			},
			{
				ID:          "rollback_procedures",
				Name:        "Rollback Procedures",
				Description: "Ability to revert changes safely",
				Weight:      0.25,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Complete rollback procedures for each operation with verification steps",
					Partial: "Some rollback info present but incomplete",
					Fail:    "Rollback procedures not documented",
				},
			},
			{
				ID:          "troubleshooting",
				Name:        "Troubleshooting",
				Description: "Guidance for common issues and errors",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Troubleshooting guide with common errors, diagnostic steps, and resolution procedures",
					Partial: "Some troubleshooting info present",
					Fail:    "No troubleshooting guidance",
				},
			},
			{
				ID:          "prerequisites",
				Name:        "Prerequisites",
				Description: "Required access, tools, and preparation",
				Weight:      0.15,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Prerequisites clearly listed with access requirements, tool versions, and preparation checklist",
					Partial: "Some prerequisites mentioned",
					Fail:    "Prerequisites not documented",
				},
			},
			{
				ID:          "contacts",
				Name:        "Escalation Contacts",
				Description: "Contact information for escalation",
				Weight:      0.15,
				Required:    false,
				Criteria: CategoricalCriteria{
					Pass:    "Escalation path with contacts, on-call schedules, and severity guidelines",
					Partial: "Some contact info present",
					Fail:    "Escalation contacts not documented",
				},
			},
		},
	}
}
