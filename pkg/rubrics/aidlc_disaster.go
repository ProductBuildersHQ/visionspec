//nolint:dupl // Rubric definitions are intentionally similar in structure
package rubrics

import "github.com/ProductBuildersHQ/visionspec/pkg/types"

func init() {
	Register(NewAIDLCDisasterRubricSet())
}

// NewAIDLCDisasterRubricSet creates the rubric set for AIDLC Disaster Recovery Plans.
func NewAIDLCDisasterRubricSet() *RubricSet {
	return &RubricSet{
		SpecType:     types.SpecTypeAIDLCDisaster,
		Name:         "AIDLC Disaster Recovery Plan Evaluation",
		Description:  "Evaluates AIDLC Disaster Recovery Plans for RTO/RPO, backup strategy, and recovery procedures",
		PassCriteria: StrictPassCriteria(), // DR plans require stricter criteria
		Categories: []*Category{
			{
				ID:          "rto_rpo",
				Name:        "RTO/RPO Definition",
				Description: "Recovery time and point objectives",
				Weight:      0.25,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "RTO and RPO clearly defined per service tier with business justification",
					Partial: "RTO/RPO mentioned but not per service or lacking justification",
					Fail:    "RTO/RPO not defined",
				},
			},
			{
				ID:          "backup_strategy",
				Name:        "Backup Strategy",
				Description: "Data backup and replication approach",
				Weight:      0.25,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Backup schedule, retention policy, geographic distribution, and encryption documented",
					Partial: "Basic backup info present but incomplete",
					Fail:    "Backup strategy not documented",
				},
			},
			{
				ID:          "recovery_procedures",
				Name:        "Recovery Procedures",
				Description: "Step-by-step recovery process",
				Weight:      0.25,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Detailed recovery procedures with verification steps, dependencies, and communication plan",
					Partial: "Recovery procedures exist but lack detail",
					Fail:    "Recovery procedures not documented",
				},
			},
			{
				ID:          "testing",
				Name:        "DR Testing",
				Description: "Disaster recovery testing plan",
				Weight:      0.15,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Regular DR testing schedule with test scenarios, success criteria, and improvement process",
					Partial: "DR testing mentioned but no schedule or criteria",
					Fail:    "DR testing not addressed",
				},
			},
			{
				ID:          "failover",
				Name:        "Failover Architecture",
				Description: "Multi-region and failover design",
				Weight:      0.10,
				Required:    false,
				Criteria: CategoricalCriteria{
					Pass:    "Failover architecture documented with automatic/manual triggers and health checks",
					Partial: "Basic failover info present",
					Fail:    "Failover architecture not addressed",
				},
			},
		},
	}
}
