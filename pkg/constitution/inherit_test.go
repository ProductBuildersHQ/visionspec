package constitution

import (
	"testing"
)

func TestMerge(t *testing.T) {
	t.Run("child overrides parent", func(t *testing.T) {
		parent := &Constitution{
			APIVersion: "visionspec/v1",
			Technical: Technical{
				Languages: Languages{
					Backend: LanguageChoice{
						Primary: "go",
						Allowed: []string{"go"},
					},
				},
			},
			Infrastructure: Infrastructure{
				Availability: Availability{
					Target: Availability999,
				},
			},
		}

		child := &Constitution{
			Technical: Technical{
				Languages: Languages{
					Backend: LanguageChoice{
						Allowed: []string{"go", "rust"}, // Override allowed
					},
				},
			},
			Infrastructure: Infrastructure{
				Availability: Availability{
					Target: Availability9999, // Override target
				},
			},
		}

		result := Merge(parent, child)

		// Child values should win
		if result.Infrastructure.Availability.Target != Availability9999 {
			t.Errorf("expected child target 99.99, got %q", result.Infrastructure.Availability.Target)
		}

		// Child allowed should override
		if len(result.Technical.Languages.Backend.Allowed) != 2 {
			t.Errorf("expected 2 allowed languages, got %d", len(result.Technical.Languages.Backend.Allowed))
		}

		// Parent primary should be preserved (child is zero)
		if result.Technical.Languages.Backend.Primary != "go" {
			t.Errorf("expected primary 'go', got %q", result.Technical.Languages.Backend.Primary)
		}

		// Parent APIVersion should be preserved
		if result.APIVersion != "visionspec/v1" {
			t.Errorf("expected APIVersion 'visionspec/v1', got %q", result.APIVersion)
		}
	})

	t.Run("nil parent returns child", func(t *testing.T) {
		child := &Constitution{
			Metadata: Metadata{Name: "child"},
		}
		result := Merge(nil, child)
		if result.Metadata.Name != "child" {
			t.Errorf("expected 'child', got %q", result.Metadata.Name)
		}
	})

	t.Run("nil child returns parent", func(t *testing.T) {
		parent := &Constitution{
			Metadata: Metadata{Name: "parent"},
		}
		result := Merge(parent, nil)
		if result.Metadata.Name != "parent" {
			t.Errorf("expected 'parent', got %q", result.Metadata.Name)
		}
	})
}

func TestResolve(t *testing.T) {
	t.Run("three-level hierarchy", func(t *testing.T) {
		org := &Constitution{
			APIVersion: "visionspec/v1",
			Metadata:   Metadata{Name: "org", Level: LevelOrganization},
			Technical: Technical{
				Languages: Languages{
					Backend: LanguageChoice{Primary: "go"},
				},
				Tenancy: Tenancy{Model: TenancyMultiTenant},
			},
			Infrastructure: Infrastructure{
				IaC:          IaC{Tool: IaCPulumi},
				Availability: Availability{Target: Availability999},
			},
		}

		team := &Constitution{
			Metadata: Metadata{Name: "team", Level: LevelTeam, Inherits: "org/org"},
			Infrastructure: Infrastructure{
				Availability: Availability{Target: Availability9999}, // Team wants higher
			},
		}

		project := &Constitution{
			Metadata: Metadata{Name: "project", Level: LevelProject, Inherits: "team/team"},
			Infrastructure: Infrastructure{
				Availability: Availability{
					RTO: "30m", // Project-specific
				},
			},
		}

		result, err := Resolve(org, team, project)
		if err != nil {
			t.Fatalf("Resolve failed: %v", err)
		}

		// Should have org's language
		if result.Technical.Languages.Backend.Primary != "go" {
			t.Errorf("expected 'go', got %q", result.Technical.Languages.Backend.Primary)
		}

		// Should have org's IaC
		if result.Infrastructure.IaC.Tool != IaCPulumi {
			t.Errorf("expected 'pulumi', got %q", result.Infrastructure.IaC.Tool)
		}

		// Should have team's availability target (overrode org)
		if result.Infrastructure.Availability.Target != Availability9999 {
			t.Errorf("expected '99.99', got %q", result.Infrastructure.Availability.Target)
		}

		// Should have project's RTO
		if result.Infrastructure.Availability.RTO != "30m" {
			t.Errorf("expected '30m', got %q", result.Infrastructure.Availability.RTO)
		}

		// Should have project's metadata
		if result.Metadata.Level != LevelProject {
			t.Errorf("expected 'project' level, got %q", result.Metadata.Level)
		}
	})

	t.Run("empty returns error", func(t *testing.T) {
		_, err := Resolve()
		if err == nil {
			t.Error("expected error for empty constitutions")
		}
	})
}

func TestValidateInheritance(t *testing.T) {
	parents := map[string]*Constitution{
		"org/example": {Metadata: Metadata{Name: "example", Level: LevelOrganization}},
	}

	t.Run("org without inherits is valid", func(t *testing.T) {
		c := &Constitution{
			Metadata: Metadata{Level: LevelOrganization},
		}
		if err := ValidateInheritance(c, parents); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("team without inherits is invalid", func(t *testing.T) {
		c := &Constitution{
			Metadata: Metadata{Level: LevelTeam},
		}
		if err := ValidateInheritance(c, parents); err == nil {
			t.Error("expected error for team without inherits")
		}
	})

	t.Run("valid inherits reference", func(t *testing.T) {
		c := &Constitution{
			Metadata: Metadata{Level: LevelTeam, Inherits: "org/example"},
		}
		if err := ValidateInheritance(c, parents); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("invalid inherits reference", func(t *testing.T) {
		c := &Constitution{
			Metadata: Metadata{Level: LevelTeam, Inherits: "org/nonexistent"},
		}
		if err := ValidateInheritance(c, parents); err == nil {
			t.Error("expected error for invalid inherits reference")
		}
	})
}
