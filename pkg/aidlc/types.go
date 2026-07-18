// Package aidlc provides types and utilities for integrating with AWS AI DLC Workflows.
// It enables bidirectional sync between VisionSpec (.visionspec/) and AIDLC (aidlc-docs/)
// directories, visualization of AIDLC workflow phases, and LLM-as-judge evaluation.
package aidlc

import "time"

// Phase represents an AIDLC workflow phase.
type Phase string

const (
	// PhaseInception is the initial discovery and requirements phase.
	PhaseInception Phase = "inception"
	// PhaseConstruction is the implementation planning and testing phase.
	PhaseConstruction Phase = "construction"
	// PhaseOperations is the deployment and operational readiness phase.
	PhaseOperations Phase = "operations"
)

// String returns the phase name.
func (p Phase) String() string {
	return string(p)
}

// Order returns the phase order (0-indexed).
func (p Phase) Order() int {
	switch p {
	case PhaseInception:
		return 0
	case PhaseConstruction:
		return 1
	case PhaseOperations:
		return 2
	default:
		return -1
	}
}

// IsValid returns whether this is a known phase.
func (p Phase) IsValid() bool {
	return p.Order() >= 0
}

// AllPhases returns all AIDLC phases in order.
func AllPhases() []Phase {
	return []Phase{PhaseInception, PhaseConstruction, PhaseOperations}
}

// DocType represents an AIDLC document type.
type DocType string

const (
	// Inception phase documents
	DocVisionDocument   DocType = "vision_document"
	DocRequirementsSpec DocType = "requirements_spec"
	DocTechnicalSpec    DocType = "technical_spec"
	DocArchitectureSpec DocType = "architecture_spec"

	// Construction phase documents
	DocImplementationPlan DocType = "implementation_plan"
	DocTestPlan           DocType = "test_plan"
	DocIntegrationPlan    DocType = "integration_plan"
	DocSecurityReview     DocType = "security_review"

	// Operations phase documents
	DocRunbook        DocType = "runbook"
	DocMonitoringPlan DocType = "monitoring_plan"
	DocDisasterPlan   DocType = "disaster_recovery_plan"
	DocSLODocument    DocType = "slo_document"
)

// String returns the document type name.
func (d DocType) String() string {
	return string(d)
}

// DisplayName returns a human-readable name for the document type.
func (d DocType) DisplayName() string {
	switch d {
	case DocVisionDocument:
		return "Vision Document"
	case DocRequirementsSpec:
		return "Requirements Specification"
	case DocTechnicalSpec:
		return "Technical Specification"
	case DocArchitectureSpec:
		return "Architecture Specification"
	case DocImplementationPlan:
		return "Implementation Plan"
	case DocTestPlan:
		return "Test Plan"
	case DocIntegrationPlan:
		return "Integration Plan"
	case DocSecurityReview:
		return "Security Review"
	case DocRunbook:
		return "Runbook"
	case DocMonitoringPlan:
		return "Monitoring Plan"
	case DocDisasterPlan:
		return "Disaster Recovery Plan"
	case DocSLODocument:
		return "SLO Document"
	default:
		return string(d)
	}
}

// Phase returns the phase this document type belongs to.
func (d DocType) Phase() Phase {
	switch d {
	case DocVisionDocument, DocRequirementsSpec, DocTechnicalSpec, DocArchitectureSpec:
		return PhaseInception
	case DocImplementationPlan, DocTestPlan, DocIntegrationPlan, DocSecurityReview:
		return PhaseConstruction
	case DocRunbook, DocMonitoringPlan, DocDisasterPlan, DocSLODocument:
		return PhaseOperations
	default:
		return ""
	}
}

// IsValid returns whether this is a known document type.
func (d DocType) IsValid() bool {
	return d.Phase() != ""
}

// Filename returns the canonical filename for this document type.
func (d DocType) Filename() string {
	return string(d) + ".md"
}

// InceptionDocTypes returns document types in the inception phase.
func InceptionDocTypes() []DocType {
	return []DocType{
		DocVisionDocument,
		DocRequirementsSpec,
		DocTechnicalSpec,
		DocArchitectureSpec,
	}
}

// ConstructionDocTypes returns document types in the construction phase.
func ConstructionDocTypes() []DocType {
	return []DocType{
		DocImplementationPlan,
		DocTestPlan,
		DocIntegrationPlan,
		DocSecurityReview,
	}
}

// OperationsDocTypes returns document types in the operations phase.
func OperationsDocTypes() []DocType {
	return []DocType{
		DocRunbook,
		DocMonitoringPlan,
		DocDisasterPlan,
		DocSLODocument,
	}
}

// AllDocTypes returns all document types in workflow order.
func AllDocTypes() []DocType {
	var result []DocType
	result = append(result, InceptionDocTypes()...)
	result = append(result, ConstructionDocTypes()...)
	result = append(result, OperationsDocTypes()...)
	return result
}

// DocTypesByPhase returns document types for a given phase.
func DocTypesByPhase(phase Phase) []DocType {
	switch phase {
	case PhaseInception:
		return InceptionDocTypes()
	case PhaseConstruction:
		return ConstructionDocTypes()
	case PhaseOperations:
		return OperationsDocTypes()
	default:
		return nil
	}
}

// DocumentStatus represents the status of an AIDLC document.
type DocumentStatus string

const (
	StatusPending    DocumentStatus = "pending"
	StatusDraft      DocumentStatus = "draft"
	StatusInProgress DocumentStatus = "in_progress"
	StatusReview     DocumentStatus = "review"
	StatusApproved   DocumentStatus = "approved"
	StatusRejected   DocumentStatus = "rejected"
)

// Document represents a parsed AIDLC document.
type Document struct {
	// Type is the document type.
	Type DocType `json:"type" yaml:"type"`

	// Phase is the workflow phase this document belongs to.
	Phase Phase `json:"phase" yaml:"phase"`

	// Path is the file path relative to the project root.
	Path string `json:"path" yaml:"path"`

	// Title is the document title.
	Title string `json:"title" yaml:"title"`

	// Status is the current document status.
	Status DocumentStatus `json:"status" yaml:"status"`

	// Content is the raw markdown content.
	Content string `json:"content,omitempty" yaml:"content,omitempty"`

	// Metadata contains frontmatter and other parsed metadata.
	Metadata map[string]any `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	// Score is the quality evaluation score (if evaluated).
	Score *QualityScore `json:"score,omitempty" yaml:"score,omitempty"`

	// UpdatedAt is when the document was last modified.
	UpdatedAt time.Time `json:"updated_at" yaml:"updated_at"`

	// Checksum is the content hash for change detection.
	Checksum string `json:"checksum,omitempty" yaml:"checksum,omitempty"`
}

// QualityRating represents the overall quality assessment.
type QualityRating string

const (
	RatingExcellent        QualityRating = "EXCELLENT"
	RatingGood             QualityRating = "GOOD"
	RatingNeedsImprovement QualityRating = "NEEDS_IMPROVEMENT"
	RatingPoor             QualityRating = "POOR"
)

// Score returns a numeric score (0.0-1.0) for the rating.
func (r QualityRating) Score() float64 {
	switch r {
	case RatingExcellent:
		return 1.0
	case RatingGood:
		return 0.75
	case RatingNeedsImprovement:
		return 0.5
	case RatingPoor:
		return 0.25
	default:
		return 0.0
	}
}

// QualityScore represents the evaluation score for a document.
type QualityScore struct {
	// Rating is the overall quality rating.
	Rating QualityRating `json:"rating" yaml:"rating"`

	// Score is a numeric score (0.0-1.0).
	Score float64 `json:"score" yaml:"score"`

	// Issues are identified problems.
	Issues []Issue `json:"issues,omitempty" yaml:"issues,omitempty"`

	// Dimensions contains per-dimension scores (if available).
	Dimensions map[string]DimensionScore `json:"dimensions,omitempty" yaml:"dimensions,omitempty"`

	// EvaluatedAt is when the evaluation was performed.
	EvaluatedAt time.Time `json:"evaluated_at" yaml:"evaluated_at"`
}

// DimensionScore represents a score for a single evaluation dimension.
type DimensionScore struct {
	// ID is the dimension identifier.
	ID string `json:"id" yaml:"id"`

	// Name is the dimension display name.
	Name string `json:"name" yaml:"name"`

	// Score is the dimension score (0.0-1.0).
	Score float64 `json:"score" yaml:"score"`

	// Weight is the dimension weight in the overall score.
	Weight float64 `json:"weight" yaml:"weight"`

	// Findings are dimension-specific issues.
	Findings []Issue `json:"findings,omitempty" yaml:"findings,omitempty"`
}

// IssueSeverity represents the severity of an issue.
type IssueSeverity string

const (
	SeverityCritical IssueSeverity = "critical"
	SeverityHigh     IssueSeverity = "high"
	SeverityMedium   IssueSeverity = "medium"
	SeverityLow      IssueSeverity = "low"
	SeverityInfo     IssueSeverity = "info"
)

// Weight returns a numeric weight for prioritization (higher = more severe).
func (s IssueSeverity) Weight() float64 {
	switch s {
	case SeverityCritical:
		return 4.0
	case SeverityHigh:
		return 3.0
	case SeverityMedium:
		return 2.0
	case SeverityLow:
		return 1.0
	case SeverityInfo:
		return 0.5
	default:
		return 0.0
	}
}

// Issue represents a quality issue found during evaluation.
type Issue struct {
	// Severity is the issue severity.
	Severity IssueSeverity `json:"severity" yaml:"severity"`

	// Category is the issue category (e.g., "clarity", "completeness").
	Category string `json:"category" yaml:"category"`

	// Code is a machine-readable issue code (e.g., "MISSING_ACCEPTANCE_CRITERIA").
	Code string `json:"code,omitempty" yaml:"code,omitempty"`

	// Message is a human-readable description.
	Message string `json:"message" yaml:"message"`

	// Location references where the issue was found (e.g., section or line).
	Location string `json:"location,omitempty" yaml:"location,omitempty"`

	// Suggestion provides remediation guidance.
	Suggestion string `json:"suggestion,omitempty" yaml:"suggestion,omitempty"`
}
