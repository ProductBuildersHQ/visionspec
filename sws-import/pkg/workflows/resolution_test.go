package workflows

import (
	"sort"
	"testing"
)

// TestRequiredSpecsResolve guards the whole embedded workflow library: every
// spec declared in a workflow's spec_config — required or optional — must
// resolve to both a template and a rubric, whether from the workflow's own
// directory, an extends chain, or an explicit spec_config template/rubric
// source. "Optional" in spec_config means optional to use when running the
// workflow (a team may choose not to produce that document), not optional for
// the library to actually support — a spec a workflow offers but can't
// render or evaluate if chosen is a gap, not a legitimate omission.
//
// This is the failure mode the pbhq-lite gap represented: a workflow whose
// spec_config named a spec, but whose templates/rubrics never resolved a
// file for it, with no test to catch the drift. There is no allowlist —
// every gap this test finds should be fixed, not suppressed.
func TestRequiredSpecsResolve(t *testing.T) {
	loader := DefaultLoader()
	names := loader.Available()
	sort.Strings(names)

	for _, name := range names {
		w, err := loader.Load(name)
		if err != nil {
			t.Errorf("workflow %q: load error: %v", name, err)
			continue
		}

		specs := make([]string, 0, len(w.Workflow.SpecConfig))
		for id, req := range w.Workflow.SpecConfig {
			if req != nil {
				specs = append(specs, id)
			}
		}
		sort.Strings(specs)

		for _, spec := range specs {
			if _, ok := w.Templates[spec]; !ok {
				t.Errorf("workflow %q: spec %q has no resolved template", name, spec)
			}
			if _, ok := w.Rubrics[spec]; !ok {
				t.Errorf("workflow %q: spec %q has no resolved rubric", name, spec)
			}
		}
	}
}
