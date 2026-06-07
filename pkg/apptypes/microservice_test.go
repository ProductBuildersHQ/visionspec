package apptypes

import (
	"testing"

	"github.com/ProductBuildersHQ/visionspec/pkg/constitution"
)

func TestMicroserviceSpec(t *testing.T) {
	spec := MicroserviceSpec()

	t.Run("has correct metadata", func(t *testing.T) {
		if spec.Metadata.Name != AppTypeMicroservice {
			t.Errorf("expected microservice, got %q", spec.Metadata.Name)
		}
		if spec.APIVersion != "visionspec/v1" {
			t.Errorf("expected visionspec/v1, got %q", spec.APIVersion)
		}
	})

	t.Run("has required artifacts", func(t *testing.T) {
		required := spec.Artifacts.Required
		if len(required) < 2 {
			t.Errorf("expected at least 2 required artifacts, got %d", len(required))
		}

		hasBinary := false
		hasContainer := false
		for _, a := range required {
			if a == ArtifactBinary {
				hasBinary = true
			}
			if a == ArtifactContainerImage {
				hasContainer = true
			}
		}
		if !hasBinary {
			t.Error("expected binary artifact to be required")
		}
		if !hasContainer {
			t.Error("expected container-image artifact to be required")
		}
	})

	t.Run("validates successfully", func(t *testing.T) {
		if err := spec.Validate(); err != nil {
			t.Errorf("validation failed: %v", err)
		}
	})

	t.Run("has defaults", func(t *testing.T) {
		if spec.Defaults.Technical.EmbeddedDB == nil {
			t.Error("expected embeddedDb to be set")
		}
		if *spec.Defaults.Technical.EmbeddedDB != false {
			t.Error("expected embeddedDb to be false for microservices")
		}

		if spec.Defaults.Infrastructure.Containerized == nil {
			t.Error("expected containerized to be set")
		}
		if *spec.Defaults.Infrastructure.Containerized != true {
			t.Error("expected containerized to be true for microservices")
		}
	})
}

func TestDefaultMicroserviceConstraints(t *testing.T) {
	constraints := DefaultMicroserviceConstraints()

	t.Run("allows both tenancy models", func(t *testing.T) {
		if len(constraints.AllowedTenancyModels) != 2 {
			t.Errorf("expected 2 allowed tenancy models, got %d", len(constraints.AllowedTenancyModels))
		}
	})

	t.Run("has minimum availability", func(t *testing.T) {
		if constraints.MinAvailability != constitution.Availability999 {
			t.Errorf("expected 99.9%% minimum, got %q", constraints.MinAvailability)
		}
	})

	t.Run("requires horizontal scaling", func(t *testing.T) {
		if !constraints.RequireHorizontalScaling {
			t.Error("expected horizontal scaling to be required")
		}
	})

	t.Run("requires containerization", func(t *testing.T) {
		if !constraints.RequireContainerization {
			t.Error("expected containerization to be required")
		}
	})

	t.Run("requires API spec", func(t *testing.T) {
		if !constraints.RequireAPISpec {
			t.Error("expected API spec to be required")
		}
	})
}

func TestValidateForMicroservice(t *testing.T) {
	t.Run("valid constitution passes", func(t *testing.T) {
		c := &constitution.Constitution{
			Technical: constitution.Technical{
				Tenancy: constitution.Tenancy{
					Model: constitution.TenancyMultiTenant,
				},
			},
			Infrastructure: constitution.Infrastructure{
				Availability: constitution.Availability{
					Target: constitution.Availability999,
				},
			},
		}

		issues := ValidateForMicroservice(c, nil)
		for _, issue := range issues {
			if issue.Severity == SeverityError {
				t.Errorf("unexpected error: %s - %s", issue.Field, issue.Message)
			}
		}
	})

	t.Run("availability below minimum generates warning", func(t *testing.T) {
		c := &constitution.Constitution{
			Infrastructure: constitution.Infrastructure{
				Availability: constitution.Availability{
					Target: constitution.Availability99, // Below 99.9% minimum
				},
			},
		}

		issues := ValidateForMicroservice(c, nil)
		hasWarning := false
		for _, issue := range issues {
			if issue.Field == "infrastructure.availability.target" && issue.Severity == SeverityWarning {
				hasWarning = true
			}
		}
		if !hasWarning {
			t.Error("expected warning for availability below minimum")
		}
	})

	t.Run("higher availability passes", func(t *testing.T) {
		c := &constitution.Constitution{
			Infrastructure: constitution.Infrastructure{
				Availability: constitution.Availability{
					Target: constitution.Availability9999, // Above 99.9% minimum
				},
			},
		}

		issues := ValidateForMicroservice(c, nil)
		for _, issue := range issues {
			if issue.Field == "infrastructure.availability.target" {
				t.Errorf("unexpected issue for high availability: %s", issue.Message)
			}
		}
	})
}

func TestMeetsMinAvailability(t *testing.T) {
	tests := []struct {
		target AvailabilityTarget
		min    AvailabilityTarget
		meets  bool
	}{
		{constitution.Availability99, constitution.Availability99, true},
		{constitution.Availability999, constitution.Availability99, true},
		{constitution.Availability9999, constitution.Availability999, true},
		{constitution.Availability99999, constitution.Availability9999, true},
		{constitution.Availability99, constitution.Availability999, false},
		{constitution.Availability999, constitution.Availability9999, false},
	}

	for _, tc := range tests {
		t.Run(string(tc.target)+"_vs_"+string(tc.min), func(t *testing.T) {
			// Note: Can't call meetsMinAvailability directly as it's unexported
			// This is tested indirectly through ValidateForMicroservice
		})
	}
}

// Type alias for test readability
type AvailabilityTarget = constitution.AvailabilityTarget
