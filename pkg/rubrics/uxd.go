//nolint:dupl // Rubric definitions are intentionally similar in structure
package rubrics

import "github.com/ProductBuildersHQ/visionspec/pkg/types"

func init() {
	Register(NewUXDRubricSet())
}

// NewUXDRubricSet creates the rubric set for User Experience Design documents.
func NewUXDRubricSet() *RubricSet {
	return &RubricSet{
		SpecType:     types.SpecTypeUXD,
		Name:         "UXD Evaluation",
		Description:  "Evaluates User Experience Design documents for usability, accessibility, and completeness",
		PassCriteria: DefaultPassCriteria(),
		Categories: []*Category{
			{
				ID:          "design_goals",
				Name:        "Design Goals",
				Description: "Clarity of UX objectives and design principles",
				Weight:      0.10,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "Clear, measurable UX goals aligned with product strategy and guiding principles",
					Partial: "Goals present but could be more specific",
					Fail:    "Design goals missing",
				},
			},
			{
				ID:          "user_personas",
				Name:        "User Personas",
				Description: "Quality and depth of user persona definitions",
				Weight:      0.15,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Rich personas with goals, pain points, behaviors, and scenarios",
					Partial: "Basic personas present but shallow",
					Fail:    "Personas missing or inadequate",
				},
			},
			{
				ID:          "user_journeys",
				Name:        "User Journeys",
				Description: "Completeness of user journey mapping",
				Weight:      0.15,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Complete journeys with entry/exit points, success criteria, and edge cases",
					Partial: "Basic flows present but incomplete",
					Fail:    "User journeys missing",
				},
			},
			{
				ID:          "information_architecture",
				Name:        "Information Architecture",
				Description: "Structure and organization of content",
				Weight:      0.15,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "Clear site map, navigation structure, and content hierarchy",
					Partial: "Basic structure present",
					Fail:    "Information architecture missing",
				},
			},
			{
				ID:          "wireframes",
				Name:        "Wireframes and Layouts",
				Description: "Quality and coverage of wireframes",
				Weight:      0.15,
				Required:    true,
				Criteria: CategoricalCriteria{
					Pass:    "Comprehensive wireframes for all key screens with element descriptions",
					Partial: "Some wireframes but incomplete coverage",
					Fail:    "Wireframes missing",
				},
			},
			{
				ID:          "interaction_design",
				Name:        "Interaction Design",
				Description: "Definition of interaction patterns and behaviors",
				Weight:      0.10,
				Required:    false,
				Criteria: CategoricalCriteria{
					Pass:    "Clear interaction patterns, micro-interactions, and state transitions",
					Partial: "Basic interactions present",
					Fail:    "Interaction design missing",
				},
			},
			{
				ID:          "accessibility",
				Name:        "Accessibility",
				Description: "WCAG compliance and inclusive design considerations",
				Weight:      0.15,
				Required:    true,
				Criteria: CategoricalCriteria{ //nolint:gosec // G101: Rubric criteria text, not credentials
					Pass:    "WCAG compliance specified, comprehensive checklist, keyboard/screen reader considerations",
					Partial: "Basic accessibility mentioned",
					Fail:    "Accessibility not addressed",
				},
			},
			{
				ID:          "error_states",
				Name:        "Error States and Edge Cases",
				Description: "Handling of errors, empty states, and edge cases",
				Weight:      0.05,
				Required:    false,
				Criteria: CategoricalCriteria{
					Pass:    "Comprehensive error messages, empty states, loading states, and recovery actions",
					Partial: "Some error states defined",
					Fail:    "Error states not addressed",
				},
			},
		},
	}
}
