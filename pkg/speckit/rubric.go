package speckit

import (
	"fmt"
	"strings"

	"github.com/plexusone/structured-evaluation/rubric"
)

// renderRubricGuardrails renders a rubric as a markdown self-evaluation
// checklist for embedding in a command body. Categories become checklist
// items; the pass-band criteria (or composite criteria) become the evidence
// an agent must satisfy before reporting the document done.
func renderRubricGuardrails(rs *rubric.RubricSet) string {
	var b strings.Builder
	for i := range rs.Categories {
		c := &rs.Categories[i]
		marker := "advisory"
		if c.Required {
			marker = "required"
		}
		if c.Blocking {
			marker += ", blocking"
		}
		fmt.Fprintf(&b, "- [ ] **%s** (%s): %s\n", c.Name, marker, c.Description)
		for _, line := range passIndicators(c) {
			fmt.Fprintf(&b, "  - %s\n", line)
		}
		if p := strings.TrimSpace(c.EvaluationPrompt); p != "" {
			fmt.Fprintf(&b, "  - Evaluation focus: %s\n", p)
		}
	}
	if line := passCriteriaLine(rs.PassCriteria); line != "" {
		fmt.Fprintf(&b, "\nOverall pass: %s\n", line)
	}
	return b.String()
}

// passIndicators returns the concrete pass-band signals for a category:
// composite categories list each criterion's pass description; simple
// categories list the criteria of the categorical "pass" option.
func passIndicators(c *rubric.Category) []string {
	if c.IsComposite() {
		lines := make([]string, 0, len(c.Criteria))
		for _, cr := range c.Criteria {
			desc := cr.Pass.Description
			if desc == "" && len(cr.Pass.Indicators) > 0 {
				desc = strings.Join(cr.Pass.Indicators, "; ")
			}
			if desc == "" {
				lines = append(lines, cr.Name)
				continue
			}
			lines = append(lines, fmt.Sprintf("%s — %s", cr.Name, desc))
		}
		return lines
	}
	for _, opt := range c.Scale.Options {
		if opt.Value == "pass" {
			return opt.Criteria
		}
	}
	return nil
}

// passCriteriaLine summarizes a rubric's overall pass criteria in one line.
func passCriteriaLine(pc rubric.RubricPassCriteria) string {
	var parts []string
	switch pc.MinCategoriesPassing {
	case "":
	case "all":
		parts = append(parts, "every category passes")
	case "all_required":
		parts = append(parts, "every required category passes")
	default:
		parts = append(parts, "at least "+pc.MinCategoriesPassing+" categories pass")
	}
	if mf := pc.MaxFindings; mf != nil {
		parts = append(parts, fmt.Sprintf("findings allowed: %s critical, %s high, %s medium, %s low",
			findingLimit(mf.Critical), findingLimit(mf.High), findingLimit(mf.Medium), findingLimit(mf.Low)))
	}
	if st := pc.ScoreThresholds; st != nil {
		parts = append(parts, fmt.Sprintf("weighted score >= %d (partial >= %d)", st.Pass, st.Partial))
	}
	return strings.Join(parts, "; ")
}

func findingLimit(n int) string {
	if n < 0 {
		return "unlimited"
	}
	return fmt.Sprintf("%d", n)
}
