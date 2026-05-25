//nolint:dupl // Rubric definitions are intentionally similar in structure
package rubrics

import "github.com/ProductBuildersHQ/visionspec/pkg/types"

func init() {
	Register(NewIRDRubricSet())
}

// NewIRDRubricSet creates the rubric set for Infrastructure Requirements Documents.
func NewIRDRubricSet() *RubricSet {
	return &RubricSet{
		SpecType:     types.SpecTypeIRD,
		Name:         "IRD Evaluation",
		Description:  "Evaluates Infrastructure Requirements Documents for completeness, security, and operability",
		PassCriteria: DefaultPassCriteria(),
		Categories: []*Category{
			{
				ID:          "architecture_completeness",
				Name:        "Architecture Completeness",
				Description: "Whether all infrastructure components are documented",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Complete infrastructure diagram with all components and connections",
					Partial: "Basic infrastructure documented but gaps exist",
					Fail:    "No clear infrastructure architecture",
				},
			},
			{
				ID:          "security_design",
				Name:        "Security Design",
				Description: "Coverage of security controls, IAM, encryption, and compliance",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Comprehensive security design with IAM, encryption, secrets, and compliance",
					Partial: "Basic security present but incomplete",
					Fail:    "No security design",
				},
			},
			{
				ID:          "availability_dr",
				Name:        "Availability and DR",
				Description: "Coverage of availability targets, failover, and disaster recovery",
				Weight:      0.15,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "Clear SLAs, multi-region strategy, and DR procedures",
					Partial: "Basic availability documented",
					Fail:    "No availability or DR planning",
				},
			},
			{
				ID:          "observability",
				Name:        "Observability",
				Description: "Coverage of logging, metrics, tracing, and alerting",
				Weight:      0.15,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "Complete observability stack with dashboards and alert definitions",
					Partial: "Basic observability present",
					Fail:    "No observability planning",
				},
			},
			{
				ID:          "capacity_cost",
				Name:        "Capacity and Cost",
				Description: "Capacity planning and cost estimation",
				Weight:      0.15,
				Required:    false,
				Criteria: CategoricalCriteria{
					Pass:    "Detailed capacity projections with cost breakdown and scaling strategy",
					Partial: "Basic capacity and cost present",
					Fail:    "No capacity or cost planning",
				},
			},
			{
				ID:          "operability",
				Name:        "Operability",
				Description: "Whether infrastructure can be operated and maintained",
				Weight:      0.15,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "Clear IaC approach, CI/CD infrastructure, and runbooks",
					Partial: "Basic operational procedures",
					Fail:    "Cannot be operated from this document",
				},
			},
		},
	}
}
