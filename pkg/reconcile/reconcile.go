// Package reconcile provides spec reconciliation capabilities.
package reconcile

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ProductBuildersHQ/visionspec/pkg/types"
)

// Reconciler generates unified specs from approved source specs.
type Reconciler struct {
	client   LLMClient
	detector *ConflictDetector
}

// LLMClient defines the interface for LLM operations.
type LLMClient interface {
	Complete(ctx context.Context, prompt string) (string, error)
}

// NewReconciler creates a new reconciler with the given LLM client.
func NewReconciler(client LLMClient) *Reconciler {
	return &Reconciler{
		client:   client,
		detector: NewConflictDetector(),
	}
}

// ReconcileInput contains all approved specs for reconciliation.
type ReconcileInput struct {
	ProjectName  string
	MRD          string
	PRD          string
	UXD          string
	TRD          string
	IRD          string
	Press        string
	FAQ          string
	Constitution string
}

// Conflict represents a detected conflict between specs.
type Conflict struct {
	ID          string           `json:"id"`
	Type        string           `json:"type"` // requirement, constraint, tradeoff
	Description string           `json:"description"`
	Sources     []types.SpecType `json:"sources"`
	Severity    string           `json:"severity"` // high, medium, low
	Resolution  string           `json:"resolution,omitempty"`
}

// ReconcileResult contains the unified spec and metadata.
type ReconcileResult struct {
	Content     string           `json:"content"`
	Sources     []types.SpecType `json:"sources"`
	Conflicts   []Conflict       `json:"conflicts,omitempty"`
	GeneratedAt time.Time        `json:"generated_at"`
	DecisionLog []string         `json:"decision_log,omitempty"`
}

// Reconcile generates a unified spec.md from approved specs.
func (r *Reconciler) Reconcile(ctx context.Context, input ReconcileInput) (*ReconcileResult, error) {
	// Detect potential conflicts before reconciliation
	detectedConflicts := r.detector.DetectConflicts(input)

	// Build prompt including detected conflicts
	prompt, sources := r.buildPrompt(input, detectedConflicts)

	// Call LLM
	content, err := r.client.Complete(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("LLM reconciliation failed: %w", err)
	}

	// Parse resolved conflicts from output
	resolvedConflicts := ParseConflictsFromOutput(content)

	// Merge detected conflicts with resolutions from output
	finalConflicts := mergeConflicts(detectedConflicts, resolvedConflicts)

	// Build decision log from conflicts
	decisionLog := buildDecisionLog(finalConflicts)

	return &ReconcileResult{
		Content:     content,
		Sources:     sources,
		Conflicts:   finalConflicts,
		GeneratedAt: time.Now(),
		DecisionLog: decisionLog,
	}, nil
}

// mergeConflicts combines detected conflicts with their resolutions.
func mergeConflicts(detected []DetectedConflict, resolved []Conflict) []Conflict {
	var merged []Conflict

	// Add detected conflicts
	for _, d := range detected {
		c := d.Conflict
		// Try to find a matching resolution
		for _, r := range resolved {
			if strings.Contains(strings.ToLower(r.Description), strings.ToLower(c.Description)[:min(20, len(c.Description))]) {
				c.Resolution = r.Resolution
				break
			}
		}
		merged = append(merged, c)
	}

	// Add any resolved conflicts not in detected list
	for _, r := range resolved {
		found := false
		for _, m := range merged {
			if m.ID == r.ID {
				found = true
				break
			}
		}
		if !found {
			merged = append(merged, r)
		}
	}

	return merged
}

// buildDecisionLog creates a summary of decisions made.
func buildDecisionLog(conflicts []Conflict) []string {
	var log []string
	for _, c := range conflicts {
		entry := c.Description
		if c.Resolution != "" {
			entry += " → " + c.Resolution
		}
		log = append(log, entry)
	}
	return log
}

// buildPrompt constructs the reconciliation prompt.
func (r *Reconciler) buildPrompt(input ReconcileInput, conflicts []DetectedConflict) (string, []types.SpecType) {
	var sb strings.Builder
	var sources []types.SpecType

	sb.WriteString("You are a technical architect reconciling multiple specification documents into a unified execution spec.\n\n")
	sb.WriteString("## Task\n\n")
	sb.WriteString("Generate a unified SPEC.md that:\n")
	sb.WriteString("1. Consolidates all requirements from source specs\n")
	sb.WriteString("2. Resolves any conflicts between specs\n")
	sb.WriteString("3. Creates a clear task breakdown for implementation\n")
	sb.WriteString("4. Maintains traceability to source requirements\n\n")

	// Include detected conflicts if any
	if len(conflicts) > 0 {
		sb.WriteString("## Detected Conflicts (Require Resolution)\n\n")
		sb.WriteString("The following potential conflicts were detected during analysis. Address each in the Decision Log section:\n\n")
		for _, c := range conflicts {
			sb.WriteString(fmt.Sprintf("- **%s** [%s, %s]: %s\n", c.ID, c.Severity, c.Type, c.Description))
			sb.WriteString(fmt.Sprintf("  Sources: %v\n", c.Sources))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## Input Specifications\n\n")

	if input.MRD != "" {
		sb.WriteString("### Market Requirements Document (MRD)\n\n")
		sb.WriteString(input.MRD)
		sb.WriteString("\n\n")
		sources = append(sources, types.SpecTypeMRD)
	}

	if input.PRD != "" {
		sb.WriteString("### Product Requirements Document (PRD)\n\n")
		sb.WriteString(input.PRD)
		sb.WriteString("\n\n")
		sources = append(sources, types.SpecTypePRD)
	}

	if input.UXD != "" {
		sb.WriteString("### User Experience Design (UXD)\n\n")
		sb.WriteString(input.UXD)
		sb.WriteString("\n\n")
		sources = append(sources, types.SpecTypeUXD)
	}

	if input.TRD != "" {
		sb.WriteString("### Technical Requirements Document (TRD)\n\n")
		sb.WriteString(input.TRD)
		sb.WriteString("\n\n")
		sources = append(sources, types.SpecTypeTRD)
	}

	if input.IRD != "" {
		sb.WriteString("### Infrastructure Requirements Document (IRD)\n\n")
		sb.WriteString(input.IRD)
		sb.WriteString("\n\n")
		sources = append(sources, types.SpecTypeIRD)
	}

	if input.Constitution != "" {
		sb.WriteString("### Constitution (Guiding Principles)\n\n")
		sb.WriteString(input.Constitution)
		sb.WriteString("\n\n")
	}

	sb.WriteString("## Output Format\n\n")
	sb.WriteString("Generate a SPEC.md with the following structure:\n\n")
	sb.WriteString("```markdown\n")
	sb.WriteString("# Execution Specification: " + input.ProjectName + "\n\n")
	sb.WriteString("## Overview\n")
	sb.WriteString("[Brief summary of what this spec covers]\n\n")
	sb.WriteString("## Requirements Summary\n")
	sb.WriteString("[Consolidated requirements from all source specs]\n\n")
	sb.WriteString("## Technical Architecture\n")
	sb.WriteString("[Key architectural decisions from TRD]\n\n")
	sb.WriteString("## Implementation Tasks\n")
	sb.WriteString("[Ordered task list with dependencies]\n\n")
	sb.WriteString("## Non-Functional Requirements\n")
	sb.WriteString("[Performance, security, etc.]\n\n")
	sb.WriteString("## Acceptance Criteria\n")
	sb.WriteString("[How to verify implementation is complete]\n\n")
	sb.WriteString("## Traceability Matrix\n")
	sb.WriteString("[Map tasks to source requirements]\n\n")
	sb.WriteString("## Decision Log\n")
	sb.WriteString("[Any conflicts resolved and tradeoffs made]\n")
	sb.WriteString("```\n\n")

	sb.WriteString("## Instructions\n\n")
	sb.WriteString("1. Output ONLY the completed SPEC.md document\n")
	sb.WriteString("2. Ensure every requirement from source specs is addressed\n")
	sb.WriteString("3. Document any conflicts found and how they were resolved\n")
	sb.WriteString("4. Create actionable implementation tasks\n")
	sb.WriteString("5. Include a traceability matrix linking tasks to requirements\n")

	return sb.String(), sources
}

// CheckApprovals verifies that required specs are approved.
func CheckApprovals(approvals map[types.SpecType]*types.Approval) ([]types.SpecType, []types.SpecType) {
	required := []types.SpecType{
		types.SpecTypeMRD,
		types.SpecTypePRD,
		types.SpecTypeTRD,
	}

	var approved, missing []types.SpecType

	for _, specType := range required {
		if approval, ok := approvals[specType]; ok && approval != nil {
			approved = append(approved, specType)
		} else {
			missing = append(missing, specType)
		}
	}

	return approved, missing
}
