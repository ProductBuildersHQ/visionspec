package integration

import "testing"

func validIntegration() *Integration {
	return &Integration{
		ID:   "sample",
		Name: "Sample Tool",
		Kind: KindExecution,
		Detection: Detection{
			RootMarkers: []string{".sample/"},
		},
		Artifacts: []Artifact{
			{ID: "spec", Path: "specs/<feature>/spec.md", Role: RoleSpec, Phase: "specify", Produced: true},
			{ID: "tasks", Path: "specs/<feature>/tasks.md", Role: RoleTasks, Phase: "plan", Produced: true},
		},
		Lifecycle: []Phase{
			{ID: "specify", Name: "Specify", Produces: []string{"spec"}},
			{ID: "plan", Name: "Plan", Consumes: []string{"spec"}, Produces: []string{"tasks"}},
		},
		Status: &StatusModel{Method: StatusTaskCheckboxes, Source: "tasks"},
	}
}

func TestValidate_OK(t *testing.T) {
	if err := validIntegration().Validate(); err != nil {
		t.Fatalf("expected valid integration, got: %v", err)
	}
}

func TestValidate_Errors(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*Integration)
	}{
		{"missing id", func(in *Integration) { in.ID = "" }},
		{"missing name", func(in *Integration) { in.Name = "" }},
		{"invalid kind", func(in *Integration) { in.Kind = "bogus" }},
		{"no detection", func(in *Integration) { in.Detection = Detection{} }},
		{"duplicate artifact id", func(in *Integration) {
			in.Artifacts = append(in.Artifacts, Artifact{ID: "spec", Path: "x.md", Role: RoleSpec})
		}},
		{"artifact empty path", func(in *Integration) {
			in.Artifacts = append(in.Artifacts, Artifact{ID: "extra", Role: RoleSpec})
		}},
		{"phase references unknown artifact", func(in *Integration) {
			in.Lifecycle = append(in.Lifecycle, Phase{ID: "impl", Name: "Implement", Consumes: []string{"nope"}})
		}},
		{"status source not an artifact", func(in *Integration) {
			in.Status = &StatusModel{Method: StatusTaskCheckboxes, Source: "nonexistent"}
		}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			in := validIntegration()
			tc.mutate(in)
			if err := in.Validate(); err == nil {
				t.Fatalf("expected error for %q, got nil", tc.name)
			}
		})
	}
}

func TestValidate_StatusGlobSourceAllowed(t *testing.T) {
	in := validIntegration()
	in.Status = &StatusModel{Method: StatusPhaseFiles, Source: "specs/*/spec.md"}
	if err := in.Validate(); err != nil {
		t.Fatalf("glob status source should be allowed, got: %v", err)
	}
}

func TestParseYAML(t *testing.T) {
	data := []byte(`
id: openspec
name: OpenSpec
kind: execution
detection:
  root_markers:
    - openspec/
artifacts:
  - id: proposal
    path: openspec/changes/<id>/proposal.md
    role: proposal
    produced: true
`)
	in, err := ParseYAML(data)
	if err != nil {
		t.Fatalf("ParseYAML: %v", err)
	}
	if in.ID != "openspec" {
		t.Fatalf("expected id openspec, got %q", in.ID)
	}
	if in.Kind != KindExecution {
		t.Fatalf("expected execution kind, got %q", in.Kind)
	}
	if len(in.Artifacts) != 1 || in.Artifacts[0].Role != RoleProposal {
		t.Fatalf("unexpected artifacts: %+v", in.Artifacts)
	}
	if err := in.Validate(); err != nil {
		t.Fatalf("parsed integration should validate: %v", err)
	}
}
