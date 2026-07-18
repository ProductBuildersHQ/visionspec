//nolint:dupl // Rubric definitions are intentionally similar in structure
package rubrics

import "github.com/ProductBuildersHQ/visionspec/pkg/types"

func init() {
	Register(NewAIDLCSecurityRubricSet())
}

// NewAIDLCSecurityRubricSet creates the rubric set for AIDLC Security Reviews.
func NewAIDLCSecurityRubricSet() *RubricSet {
	return &RubricSet{
		SpecType:     types.SpecTypeAIDLCSecurity,
		Name:         "AIDLC Security Review Evaluation",
		Description:  "Evaluates AIDLC Security Reviews for threat modeling, access control, and compliance",
		PassCriteria: StrictPassCriteria(), // Security reviews require stricter criteria
		Categories: []*Category{
			{
				ID:          "threat_modeling",
				Name:        "Threat Modeling",
				Description: "Identification and analysis of security threats",
				Weight:      0.25,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Comprehensive threat model (STRIDE or similar) with attack vectors, risk ratings, and mitigations",
					Partial: "Some threats identified but analysis incomplete",
					Fail:    "No threat modeling performed",
				},
			},
			{
				ID:          "access_control",
				Name:        "Access Control",
				Description: "Authentication and authorization design",
				Weight:      0.25,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Authentication/authorization mechanisms documented with least privilege, role-based access, and audit trails",
					Partial: "Access control mentioned but incomplete design",
					Fail:    "Access control not addressed",
				},
			},
			{
				ID:          "data_protection",
				Name:        "Data Protection",
				Description: "Encryption and sensitive data handling",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Data classification, encryption at rest/transit, and key management documented",
					Partial: "Some data protection measures mentioned",
					Fail:    "Data protection not addressed",
				},
			},
			{
				ID:          "compliance",
				Name:        "Compliance",
				Description: "Regulatory and compliance requirements",
				Weight:      0.15,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Applicable compliance requirements identified (GDPR, SOC2, etc.) with controls mapped",
					Partial: "Compliance mentioned but controls not mapped",
					Fail:    "Compliance requirements not addressed",
				},
			},
			{
				ID:          "incident_response",
				Name:        "Incident Response",
				Description: "Security incident handling procedures",
				Weight:      0.15,
				Required:    false,
				Criteria: CategoricalCriteria{
					Pass:    "Incident response plan with detection, escalation, and remediation procedures",
					Partial: "Basic incident response considerations",
					Fail:    "Incident response not addressed",
				},
			},
		},
	}
}
