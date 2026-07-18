package types

// MethodologyType categorizes methodologies into requirements vs implementation.
type MethodologyType string

const (
	// MethodologyTypeRequirements covers "WHAT to build" methodologies
	// (AWS Working Backwards, Big Tech, Lean Startup, JTBD, etc.)
	MethodologyTypeRequirements MethodologyType = "requirements"

	// MethodologyTypeImplementation covers "HOW to build" methodologies
	// (AIDLC, SpecKit, etc.)
	MethodologyTypeImplementation MethodologyType = "implementation"
)

// ImplementationMethodology defines available implementation workflows.
type ImplementationMethodology string

const (
	// ImplMethodologyNone indicates no implementation methodology is selected.
	ImplMethodologyNone ImplementationMethodology = "none"

	// ImplMethodologyAIDLC is the AWS AI DLC Workflows methodology
	// with 3 phases: Inception, Construction, Operations.
	ImplMethodologyAIDLC ImplementationMethodology = "aidlc"

	// ImplMethodologySpecKit is the GitHub SpecKit methodology.
	ImplMethodologySpecKit ImplementationMethodology = "speckit"
)

// ImplementationMethodologyInfo provides metadata about an implementation methodology.
type ImplementationMethodologyInfo struct {
	ID          ImplementationMethodology `json:"id" yaml:"id"`
	Name        string                    `json:"name" yaml:"name"`
	Description string                    `json:"description" yaml:"description"`
}

// AvailableImplementationMethodologies returns metadata for all available
// implementation methodologies.
func AvailableImplementationMethodologies() []ImplementationMethodologyInfo {
	return []ImplementationMethodologyInfo{
		{
			ID:          ImplMethodologyNone,
			Name:        "None",
			Description: "No implementation methodology - use requirements specs only",
		},
		{
			ID:          ImplMethodologyAIDLC,
			Name:        "AIDLC",
			Description: "AWS AI DLC Workflows with 3 phases: Inception, Construction, Operations",
		},
		{
			ID:          ImplMethodologySpecKit,
			Name:        "SpecKit",
			Description: "GitHub SpecKit for structured implementation specifications",
		},
	}
}

// IsValidImplementationMethodology checks if the given methodology is valid.
func IsValidImplementationMethodology(m ImplementationMethodology) bool {
	switch m {
	case ImplMethodologyNone, ImplMethodologyAIDLC, ImplMethodologySpecKit:
		return true
	default:
		return false
	}
}
