package eval

import (
	"testing"

	"github.com/plexusone/structured-evaluation/rubric"
)

func TestEvaluatePassCriteriaStrictMode(t *testing.T) {
	pc := rubric.RubricPassCriteria{MaxFindings: &rubric.FindingLimits{Critical: 0, High: 0, Medium: -1, Low: -1}}
	good := []CategoryResult{{ID: "a", IntScore: 4}, {ID: "b", IntScore: 3}}
	weak := []CategoryResult{{ID: "a", IntScore: 4}, {ID: "b", IntScore: 2}}

	if passed, _ := evaluatePassCriteriaV2(rubric.ScoreGood, weak, nil, pc, true); passed {
		t.Error("strict mode should fail when a category is below Acceptable")
	}
	if passed, _ := evaluatePassCriteriaV2(rubric.ScoreGood, good, nil, pc, true); !passed {
		t.Error("strict mode should pass when all categories >= Acceptable")
	}
	if passed, _ := evaluatePassCriteriaV2(rubric.ScoreGood, weak, nil, pc, false); !passed {
		t.Error("non-strict should pass despite a weak category")
	}
}

func TestEvaluatePassCriteriaFindingLimits(t *testing.T) {
	cats := []CategoryResult{{ID: "a", IntScore: 4}}
	pc := rubric.RubricPassCriteria{MaxFindings: &rubric.FindingLimits{Critical: 0, High: 0, Medium: 2, Low: -1}}
	findings := []Finding{{Severity: "medium"}, {Severity: "medium"}, {Severity: "medium"}}
	if passed, _ := evaluatePassCriteriaV2(rubric.ScoreGood, cats, findings, pc, false); passed {
		t.Error("should fail when medium findings exceed the configured limit")
	}
}
