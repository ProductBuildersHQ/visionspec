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
				ID:          "required_declarations",
				Name:        "Required Declarations",
				Description: "Explicit declarations for IaC approach and observability pillars (VisionSpec provides no defaults)",
				Weight:      0.15,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "IaC approach explicitly declared (Pulumi/CDK/Terraform/CloudFormation/Other/None with justification). All three observability pillars (metrics, traces, logging) explicitly declared as Implementing or None with justification.",
					Partial: "Some declarations present but incomplete (e.g., IaC stated but observability pillars missing, or vice versa)",
					Fail:    "Required declarations missing. IaC approach not stated, or observability pillars not explicitly declared.",
				},
			},
			{
				ID:          "architecture_completeness",
				Name:        "Architecture Completeness",
				Description: "Whether all infrastructure components are documented",
				Weight:      0.15,
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
				Weight:      0.15,
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
				ID:          "observability_implementation",
				Name:        "Observability Implementation",
				Description: "Implementation details for declared observability pillars (must align with Section 2.2 declaration)",
				Weight:      0.10,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "Each pillar declared as 'Implementing' has complete implementation details (platform, config, dashboards/alerts). Pillars declared as 'None' correctly marked N/A.",
					Partial: "Implementation details present but incomplete or inconsistent with declarations",
					Fail:    "No implementation details for declared pillars, or declarations contradict implementation",
				},
			},
			{
				ID:          "capacity_cost",
				Name:        "Capacity and Cost",
				Description: "Capacity planning and cost estimation",
				Weight:      0.10,
				Required:    false,
				Criteria: CategoricalCriteria{
					Pass:    "Detailed capacity projections with cost breakdown and scaling strategy",
					Partial: "Basic capacity and cost present",
					Fail:    "No capacity or cost planning",
				},
			},
			{
				ID:          "iac_implementation",
				Name:        "IaC Implementation",
				Description: "Implementation details for declared IaC approach (must align with Section 2.1 declaration)",
				Weight:      0.10,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "IaC declaration in Section 2.1 has corresponding implementation details (repo, modules, CI/CD). If 'No IaC' declared, manual procedures documented.",
					Partial: "IaC mentioned but implementation incomplete or inconsistent with declaration",
					Fail:    "No IaC implementation details, or implementation contradicts Section 2.1 declaration",
				},
			},
			{
				ID:          "operability",
				Name:        "Operability",
				Description: "Whether infrastructure can be operated and maintained",
				Weight:      0.10,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "Clear CI/CD infrastructure, runbooks, and operational procedures aligned with IaC choice",
					Partial: "Basic operational procedures",
					Fail:    "Cannot be operated from this document",
				},
			},
		},
	}
}
