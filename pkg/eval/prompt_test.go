package eval

import (
	"strings"
	"testing"

	"github.com/plexusone/structured-evaluation/rubric"
)

func TestBuildEvalPromptRendersFlatCriteria(t *testing.T) {
	rs := rubric.NewRubricSet("prd-rubric", "PRD", "1.0")
	rs.AddCategory(*rubric.NewCategory("problem", "Problem", "Problem clarity").
		SetWeight(1).
		WithPassPartialFail(
			[]string{"Clear measurable problem"},
			[]string{"Vague problem"},
			[]string{"Missing problem"},
		))

	p := buildEvalPrompt(rs, "doc")
	for _, want := range []string{"Scoring guidance", "Clear measurable problem", "Missing problem"} {
		if !strings.Contains(p, want) {
			t.Errorf("flat prompt missing %q", want)
		}
	}
}

func TestBuildEvalPromptRendersRichCriteria(t *testing.T) {
	rs := rubric.NewRubricSet("disc-rubric", "Discovery", "1.0")
	rs.AddCategory(rubric.Category{
		ID:     "cov",
		Name:   "Coverage",
		Weight: 2,
		Criteria: []rubric.Criterion{{
			ID:     "desir",
			Name:   "Desirability",
			Weight: 25,
			Pass: rubric.CriterionLevel{
				Description: "Desirability assumptions identified",
				Indicators:  []string{"customer demand cited"},
			},
		}},
	})

	p := buildEvalPrompt(rs, "doc")
	for _, want := range []string{"Weighted criteria", "Desirability", "Desirability assumptions identified", "customer demand cited"} {
		if !strings.Contains(p, want) {
			t.Errorf("rich prompt missing %q", want)
		}
	}
}
