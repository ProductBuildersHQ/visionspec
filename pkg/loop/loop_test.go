package loop

import "testing"

func validSystem() *System {
	return &System{
		ID:   "sample",
		Name: "Sample System",
		Loops: []Loop{
			{
				ID:   "outer",
				Name: "Outer",
				Stations: []Station{
					{ID: "a", Name: "A", Actor: ActorHuman},
					{ID: "b", Name: "B", Actor: ActorAI, Gate: true, GateAuthority: AuthorityPolicy},
				},
			},
			{
				ID:   "inner",
				Name: "Inner",
				Stations: []Station{
					{ID: "c", Name: "C", Actor: ActorHumanAI},
				},
			},
		},
		Seams: []Seam{{From: "outer.b", To: "inner.c", Artifact: "handoff"}},
	}
}

func TestValidate_OK(t *testing.T) {
	if err := validSystem().Validate(); err != nil {
		t.Fatalf("expected valid system, got: %v", err)
	}
}

func TestValidate_Errors(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*System)
	}{
		{"missing id", func(s *System) { s.ID = "" }},
		{"missing name", func(s *System) { s.Name = "" }},
		{"no loops", func(s *System) { s.Loops = nil }},
		{"duplicate loop id", func(s *System) { s.Loops = append(s.Loops, s.Loops[0]) }},
		{"loop without stations", func(s *System) { s.Loops[1].Stations = nil }},
		{"duplicate station id", func(s *System) {
			s.Loops[0].Stations = append(s.Loops[0].Stations, s.Loops[0].Stations[0])
		}},
		{"invalid actor", func(s *System) { s.Loops[0].Stations[0].Actor = "robot" }},
		{"gate_authority without gate", func(s *System) {
			s.Loops[0].Stations[0].GateAuthority = AuthorityHuman
		}},
		{"seam to unknown station", func(s *System) {
			s.Seams = append(s.Seams, Seam{From: "outer.a", To: "inner.nope"})
		}},
		{"seam without dot", func(s *System) {
			s.Seams = append(s.Seams, Seam{From: "outera", To: "inner.c"})
		}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := validSystem()
			tc.mutate(s)
			if err := s.Validate(); err == nil {
				t.Fatalf("expected error for %q, got nil", tc.name)
			}
		})
	}
}

func TestParseYAML(t *testing.T) {
	data := []byte(`
id: demo
name: Demo
loops:
  - id: only
    name: Only Loop
    stations:
      - id: s1
        name: Step One
        actor: ai
`)
	s, err := ParseYAML(data)
	if err != nil {
		t.Fatalf("ParseYAML: %v", err)
	}
	if err := s.Validate(); err != nil {
		t.Fatalf("parsed system should validate: %v", err)
	}
	if s.Loops[0].Stations[0].Actor != ActorAI {
		t.Fatalf("unexpected actor: %q", s.Loops[0].Stations[0].Actor)
	}
}
