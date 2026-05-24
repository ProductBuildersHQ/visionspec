// Package types defines the core data structures for multispec.
package types

import "time"

// SpecType represents the type of specification document.
type SpecType string

const (
	// Source specs (human-authored)
	SpecTypeMRD SpecType = "mrd" // Market Requirements Document
	SpecTypePRD SpecType = "prd" // Product Requirements Document
	SpecTypeUXD SpecType = "uxd" // User Experience Design

	// GTM specs (LLM-generated)
	SpecTypePress       SpecType = "press"        // Press Release
	SpecTypeFAQ         SpecType = "faq"          // FAQ
	SpecTypeNarrative1P SpecType = "narrative-1p" // 1-Pager Executive Narrative
	SpecTypeNarrative6P SpecType = "narrative-6p" // 6-Pager AWS-Style Narrative

	// Technical specs (LLM-generated)
	SpecTypeTRD SpecType = "trd" // Technical Requirements Document
	SpecTypeIRD SpecType = "ird" // Infrastructure Requirements Document

	// Output specs
	SpecTypeSpec         SpecType = "spec"          // Reconciled execution spec
	SpecTypeCurrentTruth SpecType = "current-truth" // Post-ship state
)

// SpecCategory groups spec types by their category.
type SpecCategory string

const (
	CategorySource    SpecCategory = "source"
	CategoryGTM       SpecCategory = "gtm"
	CategoryTechnical SpecCategory = "technical"
	CategoryOutput    SpecCategory = "output"
)

// Category returns the category for a spec type.
func (s SpecType) Category() SpecCategory {
	switch s {
	case SpecTypeMRD, SpecTypePRD, SpecTypeUXD:
		return CategorySource
	case SpecTypePress, SpecTypeFAQ, SpecTypeNarrative1P, SpecTypeNarrative6P:
		return CategoryGTM
	case SpecTypeTRD, SpecTypeIRD:
		return CategoryTechnical
	case SpecTypeSpec, SpecTypeCurrentTruth:
		return CategoryOutput
	default:
		return ""
	}
}

// IsRequired returns whether this spec type is required for reconciliation.
func (s SpecType) IsRequired() bool {
	switch s {
	case SpecTypeMRD, SpecTypePRD, SpecTypeUXD, SpecTypeTRD:
		return true
	default:
		return false
	}
}

// Filename returns the canonical filename for this spec type.
func (s SpecType) Filename() string {
	return string(s) + ".md"
}

// EvalFilename returns the canonical eval filename for this spec type.
func (s SpecType) EvalFilename() string {
	return string(s) + ".eval.json"
}

// Dir returns the directory where this spec type should be stored.
func (s SpecType) Dir() string {
	switch s.Category() {
	case CategorySource:
		return "source"
	case CategoryGTM:
		return "gtm"
	case CategoryTechnical:
		return "technical"
	case CategoryOutput:
		return "" // root of project
	default:
		return ""
	}
}

// IsValid returns whether this is a known spec type.
func (s SpecType) IsValid() bool {
	return s.Category() != ""
}

// SpecStatus represents the status of a specification.
type SpecStatus string

const (
	StatusMissing   SpecStatus = "missing"   // File does not exist
	StatusDraft     SpecStatus = "draft"     // File exists, not evaluated
	StatusEvaluated SpecStatus = "evaluated" // Evaluation complete
	StatusApproved  SpecStatus = "approved"  // Approved by reviewer
	StatusRejected  SpecStatus = "rejected"  // Rejected, needs revision
	StatusOutdated  SpecStatus = "outdated"  // Source changed since eval
)

// Spec represents a specification document.
type Spec struct {
	Type      SpecType   `json:"type" yaml:"type"`
	Path      string     `json:"path" yaml:"path"`
	Status    SpecStatus `json:"status" yaml:"status"`
	UpdatedAt time.Time  `json:"updated_at" yaml:"updated_at"`
	Checksum  string     `json:"checksum,omitempty" yaml:"checksum,omitempty"`
}

// AllSpecTypes returns all spec types in workflow order.
func AllSpecTypes() []SpecType {
	return []SpecType{
		// Source
		SpecTypeMRD,
		SpecTypePRD,
		SpecTypeUXD,
		// GTM
		SpecTypePress,
		SpecTypeFAQ,
		SpecTypeNarrative1P,
		SpecTypeNarrative6P,
		// Technical
		SpecTypeTRD,
		SpecTypeIRD,
		// Output
		SpecTypeSpec,
		SpecTypeCurrentTruth,
	}
}

// SourceSpecTypes returns source spec types.
func SourceSpecTypes() []SpecType {
	return []SpecType{SpecTypeMRD, SpecTypePRD, SpecTypeUXD}
}

// GTMSpecTypes returns GTM spec types.
func GTMSpecTypes() []SpecType {
	return []SpecType{SpecTypePress, SpecTypeFAQ, SpecTypeNarrative1P, SpecTypeNarrative6P}
}

// TechnicalSpecTypes returns technical spec types.
func TechnicalSpecTypes() []SpecType {
	return []SpecType{SpecTypeTRD, SpecTypeIRD}
}
