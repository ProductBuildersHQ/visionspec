//nolint:dupl // Rubric definitions are intentionally similar in structure
package rubrics

import "github.com/ProductBuildersHQ/visionspec/pkg/types"

func init() {
	Register(NewAIDLCIntegrationRubricSet())
}

// NewAIDLCIntegrationRubricSet creates the rubric set for AIDLC Integration Plans.
func NewAIDLCIntegrationRubricSet() *RubricSet {
	return &RubricSet{
		SpecType:     types.SpecTypeAIDLCIntegration,
		Name:         "AIDLC Integration Plan Evaluation",
		Description:  "Evaluates AIDLC Integration Plans for service integration, data synchronization, and rollout strategy",
		PassCriteria: DefaultPassCriteria(),
		Categories: []*Category{
			{
				ID:          "integration_points",
				Name:        "Integration Points",
				Description: "Identification and documentation of all integration points",
				Weight:      0.25,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "All integration points documented with protocols, data formats, and error handling",
					Partial: "Integration points listed but incomplete details",
					Fail:    "Integration points not identified",
				},
			},
			{
				ID:          "data_sync",
				Name:        "Data Synchronization",
				Description: "Data consistency and synchronization approach",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "Data sync strategy defined with consistency requirements, conflict resolution, and rollback procedures",
					Partial: "Some data sync considerations present",
					Fail:    "Data synchronization not addressed",
				},
			},
			{
				ID:          "rollout_strategy",
				Name:        "Rollout Strategy",
				Description: "Phased rollout and migration approach",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "Clear rollout phases with feature flags, canary deployment, and rollback triggers",
					Partial: "Basic rollout plan exists but lacks detail",
					Fail:    "No rollout strategy defined",
				},
			},
			{
				ID:          "backwards_compatibility",
				Name:        "Backwards Compatibility",
				Description: "API versioning and compatibility handling",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "API versioning strategy with deprecation timeline and migration guides",
					Partial: "Compatibility mentioned but incomplete strategy",
					Fail:    "Backwards compatibility not addressed",
				},
			},
			{
				ID:          "monitoring",
				Name:        "Integration Monitoring",
				Description: "Monitoring for integration health",
				Weight:      0.15,
				Required:    false,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "Integration monitoring defined with health checks, alerts, and dashboards",
					Partial: "Some monitoring considerations present",
					Fail:    "Integration monitoring not addressed",
				},
			},
		},
	}
}
