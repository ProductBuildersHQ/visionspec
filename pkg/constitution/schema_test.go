package constitution

import (
	"testing"
)

func TestAvailabilityTargetDowntime(t *testing.T) {
	tests := []struct {
		target   AvailabilityTarget
		expected string
	}{
		{Availability99, "3.65 days"},
		{Availability999, "8.76 hours"},
		{Availability9999, "52.6 minutes"},
		{Availability99999, "5.26 minutes"},
		{AvailabilityTarget("invalid"), "unknown"},
	}

	for _, tc := range tests {
		t.Run(string(tc.target), func(t *testing.T) {
			if got := tc.target.DowntimePerYear(); got != tc.expected {
				t.Errorf("DowntimePerYear() = %q, want %q", got, tc.expected)
			}
		})
	}
}

func TestValidEnums(t *testing.T) {
	t.Run("TenancyModels", func(t *testing.T) {
		models := ValidTenancyModels()
		if len(models) != 2 {
			t.Errorf("expected 2 tenancy models, got %d", len(models))
		}
	})

	t.Run("AvailabilityTargets", func(t *testing.T) {
		targets := ValidAvailabilityTargets()
		if len(targets) != 4 {
			t.Errorf("expected 4 availability targets, got %d", len(targets))
		}
	})

	t.Run("IaCTools", func(t *testing.T) {
		tools := ValidIaCTools()
		if len(tools) != 5 {
			t.Errorf("expected 5 IaC tools, got %d", len(tools))
		}
	})

	t.Run("LocalDevTargets", func(t *testing.T) {
		targets := ValidLocalDevTargets()
		if len(targets) != 6 {
			t.Errorf("expected 6 local dev targets, got %d", len(targets))
		}
	})
}

func TestConstitutionStructure(t *testing.T) {
	c := &Constitution{
		APIVersion: "visionspec/v1",
		Kind:       "Constitution",
		Metadata: Metadata{
			Name:  "test-org",
			Level: LevelOrganization,
		},
		Technical: Technical{
			Languages: Languages{
				Backend: LanguageChoice{
					Primary: "go",
					Allowed: []string{"go", "rust"},
				},
			},
			Tenancy: Tenancy{
				Model:     TenancyMultiTenant,
				Isolation: "rls",
			},
		},
		Infrastructure: Infrastructure{
			IaC: IaC{
				Tool:     IaCPulumi,
				Language: "go",
			},
			Availability: Availability{
				Target:  Availability999,
				MultiAZ: true,
			},
		},
	}

	if c.Metadata.Name != "test-org" {
		t.Errorf("expected name 'test-org', got %q", c.Metadata.Name)
	}
	if c.Technical.Tenancy.Model != TenancyMultiTenant {
		t.Errorf("expected multi-tenant, got %q", c.Technical.Tenancy.Model)
	}
	if c.Infrastructure.IaC.Tool != IaCPulumi {
		t.Errorf("expected pulumi, got %q", c.Infrastructure.IaC.Tool)
	}
}
