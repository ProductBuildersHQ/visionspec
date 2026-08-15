package apptypes

import (
	"testing"

	"github.com/ProductBuildersHQ/visionspec/pkg/types"
)

func TestAppTypeIsValid(t *testing.T) {
	tests := []struct {
		appType AppType
		valid   bool
	}{
		{AppTypeWebsite, true},
		{AppTypeMicroservice, true},
		{AppTypeMobile, true},
		{AppTypeDesktop, true},
		{AppTypeCLI, true},
		{AppTypeLibrary, true},
		{AppType("invalid"), false},
		{AppType(""), false},
	}

	for _, tc := range tests {
		t.Run(string(tc.appType), func(t *testing.T) {
			if got := tc.appType.IsValid(); got != tc.valid {
				t.Errorf("IsValid() = %v, want %v", got, tc.valid)
			}
		})
	}
}

func TestValidAppTypes(t *testing.T) {
	appTypes := ValidAppTypes()
	if len(appTypes) != 6 {
		t.Errorf("expected 6 app types, got %d", len(appTypes))
	}
}

func TestAppTypeSpecValidate(t *testing.T) {
	t.Run("valid spec", func(t *testing.T) {
		spec := MicroserviceSpec()
		if err := spec.Validate(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("missing apiVersion", func(t *testing.T) {
		spec := &AppTypeSpec{
			Kind:     "AppTypeSpec",
			Metadata: AppMetadata{Name: AppTypeMicroservice},
			Artifacts: Artifacts{
				Required: []ArtifactType{ArtifactBinary},
			},
			Specs: SpecRequirements{
				Required: []types.SpecType{types.SpecTypeMRD},
			},
		}
		if err := spec.Validate(); err == nil {
			t.Error("expected error for missing apiVersion")
		}
	})

	t.Run("wrong kind", func(t *testing.T) {
		spec := &AppTypeSpec{
			APIVersion: "visionspec/v1",
			Kind:       "WrongKind",
			Metadata:   AppMetadata{Name: AppTypeMicroservice},
			Artifacts: Artifacts{
				Required: []ArtifactType{ArtifactBinary},
			},
			Specs: SpecRequirements{
				Required: []types.SpecType{types.SpecTypeMRD},
			},
		}
		if err := spec.Validate(); err == nil {
			t.Error("expected error for wrong kind")
		}
	})

	t.Run("invalid app type", func(t *testing.T) {
		spec := &AppTypeSpec{
			APIVersion: "visionspec/v1",
			Kind:       "AppTypeSpec",
			Metadata:   AppMetadata{Name: AppType("invalid")},
			Artifacts: Artifacts{
				Required: []ArtifactType{ArtifactBinary},
			},
			Specs: SpecRequirements{
				Required: []types.SpecType{types.SpecTypeMRD},
			},
		}
		if err := spec.Validate(); err == nil {
			t.Error("expected error for invalid app type")
		}
	})

	t.Run("no required artifacts", func(t *testing.T) {
		spec := &AppTypeSpec{
			APIVersion: "visionspec/v1",
			Kind:       "AppTypeSpec",
			Metadata:   AppMetadata{Name: AppTypeMicroservice},
			Artifacts:  Artifacts{},
			Specs: SpecRequirements{
				Required: []types.SpecType{types.SpecTypeMRD},
			},
		}
		if err := spec.Validate(); err == nil {
			t.Error("expected error for no required artifacts")
		}
	})

	t.Run("no required specs", func(t *testing.T) {
		spec := &AppTypeSpec{
			APIVersion: "visionspec/v1",
			Kind:       "AppTypeSpec",
			Metadata:   AppMetadata{Name: AppTypeMicroservice},
			Artifacts: Artifacts{
				Required: []ArtifactType{ArtifactBinary},
			},
			Specs: SpecRequirements{},
		}
		if err := spec.Validate(); err == nil {
			t.Error("expected error for no required specs")
		}
	})
}

func TestRequiresSpec(t *testing.T) {
	spec := MicroserviceSpec()

	t.Run("required spec", func(t *testing.T) {
		if !spec.RequiresSpec(types.SpecTypeMRD) {
			t.Error("expected MRD to be required")
		}
		if !spec.RequiresSpec(types.SpecTypeTRD) {
			t.Error("expected TRD to be required")
		}
	})

	t.Run("optional spec", func(t *testing.T) {
		if spec.RequiresSpec(types.SpecTypeUXD) {
			t.Error("expected UXD to NOT be required")
		}
	})
}

func TestAllowsSpec(t *testing.T) {
	spec := MicroserviceSpec()

	t.Run("required spec is allowed", func(t *testing.T) {
		if !spec.AllowsSpec(types.SpecTypeMRD) {
			t.Error("expected MRD to be allowed")
		}
	})

	t.Run("optional spec is allowed", func(t *testing.T) {
		if !spec.AllowsSpec(types.SpecTypeUXD) {
			t.Error("expected UXD to be allowed")
		}
	})
}
