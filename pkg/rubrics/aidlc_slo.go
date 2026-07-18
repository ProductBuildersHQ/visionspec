//nolint:dupl // Rubric definitions are intentionally similar in structure
package rubrics

import "github.com/ProductBuildersHQ/visionspec/pkg/types"

func init() {
	Register(NewAIDLCSLORubricSet())
}

// NewAIDLCSLORubricSet creates the rubric set for AIDLC SLO Documents.
func NewAIDLCSLORubricSet() *RubricSet {
	return &RubricSet{
		SpecType:     types.SpecTypeAIDLCSLO,
		Name:         "AIDLC SLO Document Evaluation",
		Description:  "Evaluates AIDLC SLO Documents for objective clarity, measurement approach, and error budgets",
		PassCriteria: DefaultPassCriteria(),
		Categories: []*Category{
			{
				ID:          "sli_definition",
				Name:        "SLI Definition",
				Description: "Service Level Indicator specifications",
				Weight:      0.25,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "SLIs clearly defined with measurement method, data source, and calculation formula",
					Partial: "SLIs defined but measurement approach unclear",
					Fail:    "SLIs not defined",
				},
			},
			{
				ID:          "slo_targets",
				Name:        "SLO Targets",
				Description: "Service Level Objective targets",
				Weight:      0.25,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "SLO targets defined with specific percentages, time windows, and business rationale",
					Partial: "SLO targets exist but lack rationale or specificity",
					Fail:    "SLO targets not defined",
				},
			},
			{
				ID:          "error_budgets",
				Name:        "Error Budgets",
				Description: "Error budget definition and policies",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Error budgets defined with burn rate alerts, exhaustion policies, and reset procedures",
					Partial: "Error budgets mentioned but policies incomplete",
					Fail:    "Error budgets not defined",
				},
			},
			{
				ID:          "stakeholder_alignment",
				Name:        "Stakeholder Alignment",
				Description: "Agreement between teams on objectives",
				Weight:      0.15,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Stakeholder sign-off documented with escalation procedures and review cadence",
					Partial: "Some stakeholder info present",
					Fail:    "Stakeholder alignment not documented",
				},
			},
			{
				ID:          "reporting",
				Name:        "SLO Reporting",
				Description: "Regular SLO status reporting",
				Weight:      0.15,
				Required:    false,
				Criteria: CategoricalCriteria{
					Pass:    "Reporting cadence, dashboard links, and review meeting schedule defined",
					Partial: "Basic reporting mentioned",
					Fail:    "SLO reporting not addressed",
				},
			},
		},
	}
}
