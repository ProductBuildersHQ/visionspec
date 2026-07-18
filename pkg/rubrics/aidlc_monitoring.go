//nolint:dupl // Rubric definitions are intentionally similar in structure
package rubrics

import "github.com/ProductBuildersHQ/visionspec/pkg/types"

func init() {
	Register(NewAIDLCMonitoringRubricSet())
}

// NewAIDLCMonitoringRubricSet creates the rubric set for AIDLC Monitoring Plans.
func NewAIDLCMonitoringRubricSet() *RubricSet {
	return &RubricSet{
		SpecType:     types.SpecTypeAIDLCMonitoring,
		Name:         "AIDLC Monitoring Plan Evaluation",
		Description:  "Evaluates AIDLC Monitoring Plans for metrics, alerting, and observability",
		PassCriteria: DefaultPassCriteria(),
		Categories: []*Category{
			{
				ID:          "metrics",
				Name:        "Metrics Definition",
				Description: "Key metrics and their collection approach",
				Weight:      0.25,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Key metrics defined (RED/USE) with collection method, retention policy, and dashboards",
					Partial: "Some metrics defined but incomplete coverage",
					Fail:    "Metrics not defined",
				},
			},
			{
				ID:          "alerts",
				Name:        "Alerting Strategy",
				Description: "Alert rules and notification channels",
				Weight:      0.25,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Alert rules with thresholds, severity levels, and notification channels defined",
					Partial: "Some alerts defined but incomplete strategy",
					Fail:    "Alerting not addressed",
				},
			},
			{
				ID:          "slos",
				Name:        "Service Level Objectives",
				Description: "SLO definitions and tracking",
				Weight:      0.20,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "SLOs defined with targets, error budgets, and burn rate alerting",
					Partial: "Some SLOs mentioned but incomplete",
					Fail:    "SLOs not defined",
				},
			},
			{
				ID:          "logging",
				Name:        "Logging Strategy",
				Description: "Log collection and analysis approach",
				Weight:      0.15,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Logging strategy with log levels, structured logging, and centralized aggregation",
					Partial: "Basic logging mentioned",
					Fail:    "Logging strategy not addressed",
				},
			},
			{
				ID:          "tracing",
				Name:        "Distributed Tracing",
				Description: "Request tracing and correlation",
				Weight:      0.15,
				Required:    false,
				Criteria: CategoricalCriteria{
					Pass:    "Distributed tracing implemented with correlation IDs and trace visualization",
					Partial: "Tracing mentioned but incomplete",
					Fail:    "Distributed tracing not addressed",
				},
			},
		},
	}
}
