// Package apptypes defines application type specifications for VisionSpec.
// Each app type (website, microservice, mobile, desktop) has specific
// constraints, required artifacts, and validation rules.
package apptypes

import (
	"fmt"

	"github.com/ProductBuildersHQ/visionspec/pkg/types"
)

// AppType represents the type of application being built.
type AppType string

const (
	AppTypeWebsite      AppType = "website"
	AppTypeMicroservice AppType = "microservice"
	AppTypeMobile       AppType = "mobile"
	AppTypeDesktop      AppType = "desktop"
	AppTypeCLI          AppType = "cli"
	AppTypeLibrary      AppType = "library"
)

// AppTypeSpec defines the specification for an application type.
type AppTypeSpec struct {
	APIVersion string           `json:"apiVersion" yaml:"apiVersion"` // visionspec/v1
	Kind       string           `json:"kind" yaml:"kind"`             // AppTypeSpec
	Metadata   AppMetadata      `json:"metadata" yaml:"metadata"`
	Artifacts  Artifacts        `json:"artifacts" yaml:"artifacts"`
	Defaults   AppDefaults      `json:"defaults" yaml:"defaults"`
	Specs      SpecRequirements `json:"specs" yaml:"specs"`
	Prompts    AppPrompts       `json:"prompts,omitempty" yaml:"prompts,omitempty"`
}

// AppMetadata contains app type identification.
type AppMetadata struct {
	Name        AppType `json:"name" yaml:"name"`
	Version     string  `json:"version" yaml:"version"`
	Description string  `json:"description,omitempty" yaml:"description,omitempty"`
}

// Artifacts defines what this app type produces.
type Artifacts struct {
	Required []ArtifactType `json:"required" yaml:"required"`
	Optional []ArtifactType `json:"optional,omitempty" yaml:"optional,omitempty"`
}

// ArtifactType represents a build artifact type.
type ArtifactType string

const (
	ArtifactBinary         ArtifactType = "binary"
	ArtifactContainerImage ArtifactType = "container-image"
	ArtifactOpenAPISpec    ArtifactType = "openapi-spec"
	ArtifactProtoSpec      ArtifactType = "proto-spec"
	ArtifactHelmChart      ArtifactType = "helm-chart"
	ArtifactPulumiModule   ArtifactType = "pulumi-module"
	ArtifactStaticSite     ArtifactType = "static-site"
	ArtifactMobileApp      ArtifactType = "mobile-app"
	ArtifactDesktopApp     ArtifactType = "desktop-app"
	ArtifactNPMPackage     ArtifactType = "npm-package"
	ArtifactGoModule       ArtifactType = "go-module"
	ArtifactPyPIPackage    ArtifactType = "pypi-package"
	ArtifactWASMModule     ArtifactType = "wasm-module"
)

// AppDefaults defines default configurations for this app type.
type AppDefaults struct {
	Technical      TechnicalDefaults      `json:"technical,omitempty" yaml:"technical,omitempty"`
	Infrastructure InfrastructureDefaults `json:"infrastructure,omitempty" yaml:"infrastructure,omitempty"`
}

// TechnicalDefaults defines TRD-level defaults for this app type.
type TechnicalDefaults struct {
	APIStyles       []string `json:"apiStyles,omitempty" yaml:"apiStyles,omitempty"`             // e.g., ["rest", "grpc"]
	EmbeddedDB      *bool    `json:"embeddedDb,omitempty" yaml:"embeddedDb,omitempty"`           // Allow SQLite?
	StatefulAllowed *bool    `json:"statefulAllowed,omitempty" yaml:"statefulAllowed,omitempty"` // Can store state locally?
}

// InfrastructureDefaults defines IRD-level defaults for this app type.
type InfrastructureDefaults struct {
	Containerized         *bool    `json:"containerized,omitempty" yaml:"containerized,omitempty"`
	Orchestration         []string `json:"orchestration,omitempty" yaml:"orchestration,omitempty"` // e.g., ["kubernetes", "ecs"]
	HorizontalScaling     *bool    `json:"horizontalScaling,omitempty" yaml:"horizontalScaling,omitempty"`
	CDNRequired           *bool    `json:"cdnRequired,omitempty" yaml:"cdnRequired,omitempty"`
	MinAvailabilityTarget string   `json:"minAvailabilityTarget,omitempty" yaml:"minAvailabilityTarget,omitempty"` // e.g., "99.9"
}

// SpecRequirements defines which VisionSpec documents are required/optional.
type SpecRequirements struct {
	Required []types.SpecType `json:"required" yaml:"required"`
	Optional []types.SpecType `json:"optional,omitempty" yaml:"optional,omitempty"`
}

// AppPrompts contains LLM guidance for this app type.
type AppPrompts struct {
	WhenToUse         string `json:"whenToUse,omitempty" yaml:"whenToUse,omitempty"`
	WhenNotToUse      string `json:"whenNotToUse,omitempty" yaml:"whenNotToUse,omitempty"`
	KeyConsiderations string `json:"keyConsiderations,omitempty" yaml:"keyConsiderations,omitempty"`
}

// ValidAppTypes returns all valid app types.
func ValidAppTypes() []AppType {
	return []AppType{
		AppTypeWebsite,
		AppTypeMicroservice,
		AppTypeMobile,
		AppTypeDesktop,
		AppTypeCLI,
		AppTypeLibrary,
	}
}

// IsValid returns whether this is a known app type.
func (a AppType) IsValid() bool {
	for _, valid := range ValidAppTypes() {
		if a == valid {
			return true
		}
	}
	return false
}

// String returns the string representation.
func (a AppType) String() string {
	return string(a)
}

// Validate validates an AppTypeSpec.
func (s *AppTypeSpec) Validate() error {
	if s.APIVersion == "" {
		return fmt.Errorf("apiVersion is required")
	}
	if s.Kind != "AppTypeSpec" {
		return fmt.Errorf("kind must be 'AppTypeSpec'")
	}
	if !s.Metadata.Name.IsValid() {
		return fmt.Errorf("invalid app type: %s", s.Metadata.Name)
	}
	if len(s.Artifacts.Required) == 0 {
		return fmt.Errorf("at least one required artifact must be specified")
	}
	if len(s.Specs.Required) == 0 {
		return fmt.Errorf("at least one required spec must be specified")
	}
	return nil
}

// RequiresSpec returns true if this app type requires the given spec.
func (s *AppTypeSpec) RequiresSpec(spec types.SpecType) bool {
	for _, r := range s.Specs.Required {
		if r == spec {
			return true
		}
	}
	return false
}

// AllowsSpec returns true if this app type allows the given spec (required or optional).
func (s *AppTypeSpec) AllowsSpec(spec types.SpecType) bool {
	for _, r := range s.Specs.Required {
		if r == spec {
			return true
		}
	}
	for _, o := range s.Specs.Optional {
		if o == spec {
			return true
		}
	}
	return false
}
