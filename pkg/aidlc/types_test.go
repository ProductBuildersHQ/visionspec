package aidlc

import (
	"testing"
)

func TestPhaseString(t *testing.T) {
	tests := []struct {
		phase Phase
		want  string
	}{
		{PhaseInception, "inception"},
		{PhaseConstruction, "construction"},
		{PhaseOperations, "operations"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.phase.String(); got != tt.want {
				t.Errorf("Phase.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPhaseOrder(t *testing.T) {
	tests := []struct {
		phase Phase
		want  int
	}{
		{PhaseInception, 0},
		{PhaseConstruction, 1},
		{PhaseOperations, 2},
		{Phase("unknown"), -1},
	}

	for _, tt := range tests {
		t.Run(string(tt.phase), func(t *testing.T) {
			if got := tt.phase.Order(); got != tt.want {
				t.Errorf("Phase.Order() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestPhaseIsValid(t *testing.T) {
	tests := []struct {
		phase Phase
		want  bool
	}{
		{PhaseInception, true},
		{PhaseConstruction, true},
		{PhaseOperations, true},
		{Phase("unknown"), false},
		{Phase(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.phase), func(t *testing.T) {
			if got := tt.phase.IsValid(); got != tt.want {
				t.Errorf("Phase.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAllPhases(t *testing.T) {
	phases := AllPhases()
	if len(phases) != 3 {
		t.Errorf("AllPhases() returned %d phases, want 3", len(phases))
	}

	// Check order
	expected := []Phase{PhaseInception, PhaseConstruction, PhaseOperations}
	for i, phase := range phases {
		if phase != expected[i] {
			t.Errorf("AllPhases()[%d] = %s, want %s", i, phase, expected[i])
		}
	}
}

func TestDocTypeDisplayName(t *testing.T) {
	tests := []struct {
		docType DocType
		want    string
	}{
		{DocVisionDocument, "Vision Document"},
		{DocRequirementsSpec, "Requirements Specification"},
		{DocTechnicalSpec, "Technical Specification"},
		{DocArchitectureSpec, "Architecture Specification"},
		{DocImplementationPlan, "Implementation Plan"},
		{DocTestPlan, "Test Plan"},
		{DocIntegrationPlan, "Integration Plan"},
		{DocSecurityReview, "Security Review"},
		{DocRunbook, "Runbook"},
		{DocMonitoringPlan, "Monitoring Plan"},
		{DocDisasterPlan, "Disaster Recovery Plan"},
		{DocSLODocument, "SLO Document"},
		{DocType("unknown"), "unknown"},
	}

	for _, tt := range tests {
		t.Run(string(tt.docType), func(t *testing.T) {
			if got := tt.docType.DisplayName(); got != tt.want {
				t.Errorf("DocType.DisplayName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDocTypePhase(t *testing.T) {
	tests := []struct {
		docType DocType
		want    Phase
	}{
		// Inception phase
		{DocVisionDocument, PhaseInception},
		{DocRequirementsSpec, PhaseInception},
		{DocTechnicalSpec, PhaseInception},
		{DocArchitectureSpec, PhaseInception},
		// Construction phase
		{DocImplementationPlan, PhaseConstruction},
		{DocTestPlan, PhaseConstruction},
		{DocIntegrationPlan, PhaseConstruction},
		{DocSecurityReview, PhaseConstruction},
		// Operations phase
		{DocRunbook, PhaseOperations},
		{DocMonitoringPlan, PhaseOperations},
		{DocDisasterPlan, PhaseOperations},
		{DocSLODocument, PhaseOperations},
		// Unknown
		{DocType("unknown"), Phase("")},
	}

	for _, tt := range tests {
		t.Run(string(tt.docType), func(t *testing.T) {
			if got := tt.docType.Phase(); got != tt.want {
				t.Errorf("DocType.Phase() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestDocTypeIsValid(t *testing.T) {
	tests := []struct {
		docType DocType
		want    bool
	}{
		{DocVisionDocument, true},
		{DocRequirementsSpec, true},
		{DocRunbook, true},
		{DocType("unknown"), false},
		{DocType(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.docType), func(t *testing.T) {
			if got := tt.docType.IsValid(); got != tt.want {
				t.Errorf("DocType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocTypeFilename(t *testing.T) {
	tests := []struct {
		docType DocType
		want    string
	}{
		{DocVisionDocument, "vision_document.md"},
		{DocRequirementsSpec, "requirements_spec.md"},
		{DocRunbook, "runbook.md"},
	}

	for _, tt := range tests {
		t.Run(string(tt.docType), func(t *testing.T) {
			if got := tt.docType.Filename(); got != tt.want {
				t.Errorf("DocType.Filename() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestInceptionDocTypes(t *testing.T) {
	docs := InceptionDocTypes()
	if len(docs) != 4 {
		t.Errorf("InceptionDocTypes() returned %d docs, want 4", len(docs))
	}

	for _, doc := range docs {
		if doc.Phase() != PhaseInception {
			t.Errorf("InceptionDocTypes() contains %s with phase %s", doc, doc.Phase())
		}
	}
}

func TestConstructionDocTypes(t *testing.T) {
	docs := ConstructionDocTypes()
	if len(docs) != 4 {
		t.Errorf("ConstructionDocTypes() returned %d docs, want 4", len(docs))
	}

	for _, doc := range docs {
		if doc.Phase() != PhaseConstruction {
			t.Errorf("ConstructionDocTypes() contains %s with phase %s", doc, doc.Phase())
		}
	}
}

func TestOperationsDocTypes(t *testing.T) {
	docs := OperationsDocTypes()
	if len(docs) != 4 {
		t.Errorf("OperationsDocTypes() returned %d docs, want 4", len(docs))
	}

	for _, doc := range docs {
		if doc.Phase() != PhaseOperations {
			t.Errorf("OperationsDocTypes() contains %s with phase %s", doc, doc.Phase())
		}
	}
}

func TestAllDocTypes(t *testing.T) {
	docs := AllDocTypes()
	if len(docs) != 12 {
		t.Errorf("AllDocTypes() returned %d docs, want 12", len(docs))
	}

	// Verify order: inception -> construction -> operations
	var lastPhaseOrder int = -1
	for _, doc := range docs {
		phaseOrder := doc.Phase().Order()
		if phaseOrder < lastPhaseOrder {
			t.Errorf("AllDocTypes() not in phase order at %s", doc)
		}
		lastPhaseOrder = phaseOrder
	}
}

func TestDocTypesByPhase(t *testing.T) {
	tests := []struct {
		phase Phase
		count int
	}{
		{PhaseInception, 4},
		{PhaseConstruction, 4},
		{PhaseOperations, 4},
		{Phase("unknown"), 0},
	}

	for _, tt := range tests {
		t.Run(string(tt.phase), func(t *testing.T) {
			docs := DocTypesByPhase(tt.phase)
			if len(docs) != tt.count {
				t.Errorf("DocTypesByPhase(%s) returned %d docs, want %d", tt.phase, len(docs), tt.count)
			}
		})
	}
}

func TestQualityRatingScore(t *testing.T) {
	tests := []struct {
		rating QualityRating
		want   float64
	}{
		{RatingExcellent, 1.0},
		{RatingGood, 0.75},
		{RatingNeedsImprovement, 0.5},
		{RatingPoor, 0.25},
		{QualityRating("unknown"), 0.0},
	}

	for _, tt := range tests {
		t.Run(string(tt.rating), func(t *testing.T) {
			if got := tt.rating.Score(); got != tt.want {
				t.Errorf("QualityRating.Score() = %f, want %f", got, tt.want)
			}
		})
	}
}

func TestIssueSeverityWeight(t *testing.T) {
	tests := []struct {
		severity IssueSeverity
		want     float64
	}{
		{SeverityCritical, 4.0},
		{SeverityHigh, 3.0},
		{SeverityMedium, 2.0},
		{SeverityLow, 1.0},
		{SeverityInfo, 0.5},
		{IssueSeverity("unknown"), 0.0},
	}

	for _, tt := range tests {
		t.Run(string(tt.severity), func(t *testing.T) {
			if got := tt.severity.Weight(); got != tt.want {
				t.Errorf("IssueSeverity.Weight() = %f, want %f", got, tt.want)
			}
		})
	}
}
