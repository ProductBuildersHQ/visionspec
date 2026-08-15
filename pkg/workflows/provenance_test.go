package workflows

import (
	"testing"

	"github.com/ProductBuildersHQ/specification-workflow-spec/pkg/template"
	"github.com/ProductBuildersHQ/specification-workflow-spec/pkg/workflow"
	"github.com/plexusone/structured-evaluation/rubric"
)

// TestPBHQLiteProvenanceResolves verifies that pbhq-lite's declared per-spec
// template/rubric sources resolve: PRD is local, and TRD/PLAN/ROADMAP templates
// are sourced from enterprise (filling the former gap) while their rubrics stay
// local.
func TestPBHQLiteProvenanceResolves(t *testing.T) {
	w, err := DefaultLoader().Load("pbhq-lite")
	if err != nil {
		t.Fatalf("Load(pbhq-lite) error: %v", err)
	}

	for _, spec := range []string{"prd", "trd", "plan", "roadmap"} {
		if w.Templates[spec] == nil {
			t.Errorf("pbhq-lite: template for %q did not resolve", spec)
		}
		if w.Rubrics[spec] == nil {
			t.Errorf("pbhq-lite: rubric for %q did not resolve", spec)
		}
	}

	// The borrowed TRD template must be enterprise's content, not empty.
	if got := w.Templates["trd"]; got != nil && got.Content == "" {
		t.Error("pbhq-lite: sourced trd template has empty content")
	}
}

// TestQuickFixInheritsResolvedProvenance guards the inheritance interaction:
// quick-fix extends pbhq-lite, so it inherits the "local" declarations. Those
// must stay pinned to pbhq-lite (already resolved) rather than re-resolving
// against quick-fix's own directory.
func TestQuickFixInheritsResolvedProvenance(t *testing.T) {
	w, err := DefaultLoader().Load("quick-fix")
	if err != nil {
		t.Fatalf("Load(quick-fix) error: %v", err)
	}
	for _, spec := range []string{"prd", "trd", "plan", "roadmap"} {
		if w.Templates[spec] == nil {
			t.Errorf("quick-fix: inherited template for %q did not resolve", spec)
		}
	}
}

// TestExplicitSourceEnforced verifies that a declared source which does not
// provide the file is a load-time error — the property that makes provenance
// authoritative rather than advisory.
func TestExplicitSourceEnforced(t *testing.T) {
	base := &stubLoader{workflows: map[string]*LoadedWorkflow{
		"child": {
			Workflow: &workflow.Workflow{
				Name: "child",
				SpecConfig: map[string]*workflow.SpecRequirement{
					// enterprise-style source, but the stub source lacks it.
					"trd": {Required: true, Template: &workflow.SpecSource{From: "empty-src"}},
				},
			},
			Templates: map[string]*template.Template{},
			Rubrics:   map[string]*rubric.RubricSet{},
		},
		"empty-src": {
			Workflow:  &workflow.Workflow{Name: "empty-src"},
			Templates: map[string]*template.Template{},
			Rubrics:   map[string]*rubric.RubricSet{},
		},
	}}

	_, err := NewResolvingLoader(base).Load("child")
	if err == nil {
		t.Fatal("expected error when declared template source has no such template")
	}
}

// stubLoader is an in-memory Loader for testing source resolution.
type stubLoader struct {
	workflows map[string]*LoadedWorkflow
}

func (s *stubLoader) Load(name string) (*LoadedWorkflow, error) {
	w, ok := s.workflows[name]
	if !ok {
		return nil, &NotFoundError{Name: name}
	}
	// Return a shallow copy so resolution does not mutate the fixture.
	return &LoadedWorkflow{
		Workflow:  w.Workflow,
		Templates: cloneTemplateMap(w.Templates),
		Rubrics:   cloneRubricMap(w.Rubrics),
	}, nil
}

func (s *stubLoader) Available() []string {
	names := make([]string, 0, len(s.workflows))
	for n := range s.workflows {
		names = append(names, n)
	}
	return names
}

func cloneTemplateMap(m map[string]*template.Template) map[string]*template.Template {
	out := make(map[string]*template.Template, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func cloneRubricMap(m map[string]*rubric.RubricSet) map[string]*rubric.RubricSet {
	out := make(map[string]*rubric.RubricSet, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
