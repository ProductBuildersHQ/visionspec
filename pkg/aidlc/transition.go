// Package aidlc provides phase transition logic for AIDLC workflows.
package aidlc

import (
	"fmt"
	"time"
)

// PhaseStatus represents the status of a workflow phase.
type PhaseStatus string

const (
	// PhaseStatusPending indicates the phase has not started.
	PhaseStatusPending PhaseStatus = "pending"
	// PhaseStatusInProgress indicates the phase is currently active.
	PhaseStatusInProgress PhaseStatus = "in_progress"
	// PhaseStatusCompleted indicates all required documents are complete.
	PhaseStatusCompleted PhaseStatus = "completed"
	// PhaseStatusBlocked indicates the phase cannot proceed due to issues.
	PhaseStatusBlocked PhaseStatus = "blocked"
)

// TransitionRule defines requirements for transitioning to a phase.
type TransitionRule struct {
	// TargetPhase is the phase to transition to.
	TargetPhase Phase `json:"target_phase" yaml:"target_phase"`

	// RequiredDocs are document types that must be completed.
	RequiredDocs []DocType `json:"required_docs" yaml:"required_docs"`

	// MinScore is the minimum quality score required (0-1).
	MinScore float64 `json:"min_score,omitempty" yaml:"min_score,omitempty"`

	// AllowPartial allows transition with partial scores.
	AllowPartial bool `json:"allow_partial,omitempty" yaml:"allow_partial,omitempty"`

	// RequireApproval requires explicit approval for transition.
	RequireApproval bool `json:"require_approval,omitempty" yaml:"require_approval,omitempty"`
}

// TransitionResult contains the result of a phase transition attempt.
type TransitionResult struct {
	// Success indicates if the transition succeeded.
	Success bool `json:"success" yaml:"success"`

	// FromPhase is the source phase.
	FromPhase Phase `json:"from_phase" yaml:"from_phase"`

	// ToPhase is the target phase.
	ToPhase Phase `json:"to_phase" yaml:"to_phase"`

	// BlockingDocs are documents preventing transition.
	BlockingDocs []DocType `json:"blocking_docs,omitempty" yaml:"blocking_docs,omitempty"`

	// BlockingIssues are issues preventing transition.
	BlockingIssues []string `json:"blocking_issues,omitempty" yaml:"blocking_issues,omitempty"`

	// Timestamp is when the transition was attempted.
	Timestamp time.Time `json:"timestamp" yaml:"timestamp"`
}

// TransitionLog tracks phase transitions over time.
type TransitionLog struct {
	// Entries are the transition history.
	Entries []TransitionEntry `json:"entries" yaml:"entries"`
}

// TransitionEntry records a single transition event.
type TransitionEntry struct {
	// FromPhase is the source phase.
	FromPhase Phase `json:"from_phase" yaml:"from_phase"`

	// ToPhase is the target phase.
	ToPhase Phase `json:"to_phase" yaml:"to_phase"`

	// Timestamp is when the transition occurred.
	Timestamp time.Time `json:"timestamp" yaml:"timestamp"`

	// ApprovedBy is who approved the transition (if required).
	ApprovedBy string `json:"approved_by,omitempty" yaml:"approved_by,omitempty"`

	// Notes are optional notes about the transition.
	Notes string `json:"notes,omitempty" yaml:"notes,omitempty"`
}

// DefaultTransitionRules returns the default transition rules.
func DefaultTransitionRules() map[Phase]TransitionRule {
	return map[Phase]TransitionRule{
		PhaseConstruction: {
			TargetPhase: PhaseConstruction,
			RequiredDocs: []DocType{
				DocVisionDocument,
				DocRequirementsSpec,
				DocTechnicalSpec,
			},
			MinScore:     0.7,
			AllowPartial: true,
		},
		PhaseOperations: {
			TargetPhase: PhaseOperations,
			RequiredDocs: []DocType{
				DocImplementationPlan,
				DocTestPlan,
				DocSecurityReview,
			},
			MinScore:     0.7,
			AllowPartial: false,
		},
	}
}

// CurrentPhase determines the current workflow phase based on node statuses.
func (w *Workflow) CurrentPhase() Phase {
	// Start from the last phase and work backwards to find the first incomplete phase
	for i := len(w.Phases) - 1; i >= 0; i-- {
		phase := Phase(w.Phases[i].ID)
		status := w.GetPhaseStatus(phase)
		if status == PhaseStatusCompleted {
			// If we're on the last phase and it's complete, we're done
			if i == len(w.Phases)-1 {
				return phase
			}
			// Otherwise, we're in the next phase
			continue
		}
		if status == PhaseStatusInProgress {
			return phase
		}
	}

	// Default to inception if nothing is started
	return PhaseInception
}

// GetPhaseStatus returns the status of a specific phase.
func (w *Workflow) GetPhaseStatus(phase Phase) PhaseStatus {
	wp, ok := w.GetPhase(string(phase))
	if !ok {
		return PhaseStatusPending
	}

	var hasCompleted, hasInProgress, hasPending bool

	for _, nodeID := range wp.NodeIDs {
		node, ok := w.Nodes[nodeID]
		if !ok {
			continue
		}

		switch node.Status {
		case NodeCompleted:
			hasCompleted = true
		case NodeInProgress:
			hasInProgress = true
		case NodePending, NodeReady:
			if node.Required {
				hasPending = true
			}
		case NodeBlocked, NodeFailed:
			return PhaseStatusBlocked
		}
	}

	// Determine overall phase status
	if hasPending {
		if hasInProgress || hasCompleted {
			return PhaseStatusInProgress
		}
		return PhaseStatusPending
	}
	if hasInProgress {
		return PhaseStatusInProgress
	}
	if hasCompleted {
		return PhaseStatusCompleted
	}
	return PhaseStatusPending
}

// CanTransitionTo checks if the workflow can transition to a target phase.
func (w *Workflow) CanTransitionTo(targetPhase Phase, rules map[Phase]TransitionRule) *TransitionResult {
	result := &TransitionResult{
		FromPhase: w.CurrentPhase(),
		ToPhase:   targetPhase,
		Timestamp: time.Now(),
	}

	// Cannot transition to a phase that comes before current phase
	if targetPhase.Order() <= w.CurrentPhase().Order() {
		result.BlockingIssues = append(result.BlockingIssues,
			fmt.Sprintf("cannot transition backwards from %s to %s", w.CurrentPhase(), targetPhase))
		return result
	}

	// Cannot skip phases
	if targetPhase.Order() > w.CurrentPhase().Order()+1 {
		result.BlockingIssues = append(result.BlockingIssues,
			fmt.Sprintf("cannot skip phases: must complete %s first", w.getNextPhase()))
		return result
	}

	// Check transition rules for the target phase
	rule, ok := rules[targetPhase]
	if !ok {
		// No rules defined - allow transition
		result.Success = true
		return result
	}

	// Check required documents
	for _, docType := range rule.RequiredDocs {
		node, ok := w.Nodes[string(docType)]
		if !ok {
			result.BlockingDocs = append(result.BlockingDocs, docType)
			result.BlockingIssues = append(result.BlockingIssues,
				fmt.Sprintf("required document %s not found in workflow", docType.DisplayName()))
			continue
		}

		if node.Status != NodeCompleted && node.Status != NodeSkipped {
			result.BlockingDocs = append(result.BlockingDocs, docType)
			result.BlockingIssues = append(result.BlockingIssues,
				fmt.Sprintf("required document %s is not complete (status: %s)", docType.DisplayName(), node.Status))
			continue
		}

		// Check quality score if required
		if rule.MinScore > 0 && node.Score != nil {
			if node.Score.Score < rule.MinScore {
				result.BlockingDocs = append(result.BlockingDocs, docType)
				result.BlockingIssues = append(result.BlockingIssues,
					fmt.Sprintf("document %s score %.2f is below minimum %.2f",
						docType.DisplayName(), node.Score.Score, rule.MinScore))
			}
		}
	}

	result.Success = len(result.BlockingDocs) == 0 && len(result.BlockingIssues) == 0
	return result
}

// TransitionTo attempts to transition the workflow to a new phase.
func (w *Workflow) TransitionTo(targetPhase Phase, rules map[Phase]TransitionRule) (*TransitionResult, error) {
	result := w.CanTransitionTo(targetPhase, rules)
	if !result.Success {
		return result, fmt.Errorf("transition blocked: %v", result.BlockingIssues)
	}

	// Update node statuses for the new phase
	newPhaseNodes := DocTypesByPhase(targetPhase)
	for _, docType := range newPhaseNodes {
		nodeID := string(docType)
		if node, ok := w.Nodes[nodeID]; ok {
			if node.Status == NodePending || node.Status == NodeBlocked {
				if w.canStart(node) {
					node.Status = NodeReady
				}
			}
		}
	}

	return result, nil
}

// getNextPhase returns the next phase in the workflow.
func (w *Workflow) getNextPhase() Phase {
	current := w.CurrentPhase()
	switch current {
	case PhaseInception:
		return PhaseConstruction
	case PhaseConstruction:
		return PhaseOperations
	default:
		return current
	}
}

// PhaseRequirements returns the requirements for completing a phase.
type PhaseRequirements struct {
	// Phase is the phase being checked.
	Phase Phase `json:"phase" yaml:"phase"`

	// RequiredDocs are the required documents.
	RequiredDocs []DocType `json:"required_docs" yaml:"required_docs"`

	// CompletedDocs are the completed required documents.
	CompletedDocs []DocType `json:"completed_docs" yaml:"completed_docs"`

	// PendingDocs are the pending required documents.
	PendingDocs []DocType `json:"pending_docs" yaml:"pending_docs"`

	// ProgressPercent is the completion percentage.
	ProgressPercent float64 `json:"progress_percent" yaml:"progress_percent"`

	// CanAdvance indicates if the phase can advance.
	CanAdvance bool `json:"can_advance" yaml:"can_advance"`
}

// GetPhaseRequirements returns the requirements for a phase.
func (w *Workflow) GetPhaseRequirements(phase Phase) PhaseRequirements {
	reqs := PhaseRequirements{
		Phase:        phase,
		RequiredDocs: make([]DocType, 0),
	}

	wp, ok := w.GetPhase(string(phase))
	if !ok {
		return reqs
	}

	for _, nodeID := range wp.NodeIDs {
		node, ok := w.Nodes[nodeID]
		if !ok {
			continue
		}

		if node.Required {
			reqs.RequiredDocs = append(reqs.RequiredDocs, node.DocType)
			if node.Status == NodeCompleted {
				reqs.CompletedDocs = append(reqs.CompletedDocs, node.DocType)
			} else {
				reqs.PendingDocs = append(reqs.PendingDocs, node.DocType)
			}
		}
	}

	if len(reqs.RequiredDocs) > 0 {
		reqs.ProgressPercent = float64(len(reqs.CompletedDocs)) / float64(len(reqs.RequiredDocs)) * 100
	}

	reqs.CanAdvance = len(reqs.PendingDocs) == 0

	return reqs
}

// AllPhaseRequirements returns requirements for all phases.
func (w *Workflow) AllPhaseRequirements() []PhaseRequirements {
	var reqs []PhaseRequirements
	for _, phase := range AllPhases() {
		reqs = append(reqs, w.GetPhaseRequirements(phase))
	}
	return reqs
}

// ValidatePhaseTransition performs comprehensive validation for a phase transition.
func (w *Workflow) ValidatePhaseTransition(targetPhase Phase) []string {
	var issues []string

	// Check phase ordering
	currentPhase := w.CurrentPhase()
	if targetPhase.Order() < currentPhase.Order() {
		issues = append(issues, "cannot transition to a previous phase")
		return issues
	}

	// Check all intermediate phases
	for i := currentPhase.Order(); i < targetPhase.Order(); i++ {
		intermediatePhase := AllPhases()[i]
		reqs := w.GetPhaseRequirements(intermediatePhase)
		if !reqs.CanAdvance {
			issues = append(issues,
				fmt.Sprintf("phase %s has incomplete required documents: %v",
					intermediatePhase, reqs.PendingDocs))
		}
	}

	// Check for blocked or failed nodes
	for _, nodeID := range w.Phases[currentPhase.Order()].NodeIDs {
		if node, ok := w.Nodes[nodeID]; ok {
			if node.Required && (node.Status == NodeBlocked || node.Status == NodeFailed) {
				issues = append(issues,
					fmt.Sprintf("required document %s is %s", node.Name, node.Status))
			}
		}
	}

	return issues
}

// NewTransitionLog creates a new transition log.
func NewTransitionLog() *TransitionLog {
	return &TransitionLog{
		Entries: make([]TransitionEntry, 0),
	}
}

// AddEntry adds a transition entry to the log.
func (l *TransitionLog) AddEntry(from, to Phase, approvedBy, notes string) {
	l.Entries = append(l.Entries, TransitionEntry{
		FromPhase:  from,
		ToPhase:    to,
		Timestamp:  time.Now(),
		ApprovedBy: approvedBy,
		Notes:      notes,
	})
}

// LatestEntry returns the most recent transition entry.
func (l *TransitionLog) LatestEntry() *TransitionEntry {
	if len(l.Entries) == 0 {
		return nil
	}
	return &l.Entries[len(l.Entries)-1]
}

// EntriesForPhase returns all entries related to a phase.
func (l *TransitionLog) EntriesForPhase(phase Phase) []TransitionEntry {
	var entries []TransitionEntry
	for _, e := range l.Entries {
		if e.FromPhase == phase || e.ToPhase == phase {
			entries = append(entries, e)
		}
	}
	return entries
}
