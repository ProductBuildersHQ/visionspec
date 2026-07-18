// Package aidlc provides workflow progress tracking for AIDLC workflows.
package aidlc

import (
	"fmt"
	"sort"
	"time"

	"github.com/grokify/pidl"
)

// StepExecutionStatus represents the execution status of a workflow step.
type StepExecutionStatus string

const (
	StepStatusPending    StepExecutionStatus = "pending"
	StepStatusReady      StepExecutionStatus = "ready"
	StepStatusBlocked    StepExecutionStatus = "blocked"
	StepStatusInProgress StepExecutionStatus = "in_progress"
	StepStatusCompleted  StepExecutionStatus = "completed"
	StepStatusFailed     StepExecutionStatus = "failed"
	StepStatusSkipped    StepExecutionStatus = "skipped"
)

// WorkflowExecutionContext tracks AIDLC workflow execution.
type WorkflowExecutionContext struct {
	// Workflow is the underlying AIDLC workflow.
	Workflow *Workflow

	// Protocol is the pidl protocol representation.
	Protocol *pidl.Protocol

	// StartTime is when execution started.
	StartTime time.Time

	// StepTimings tracks timing for each step.
	StepTimings map[string]*StepTiming

	// StepStatus tracks the status of each step.
	StepStatus map[string]StepExecutionStatus

	// ExecutionOrder is the topologically sorted execution order.
	ExecutionOrder []string

	// Dependencies maps step ID to its required predecessor step IDs.
	Dependencies map[string][]string

	// Dependents maps step ID to steps that depend on it.
	Dependents map[string][]string
}

// StepTiming tracks timing metrics for a workflow step.
type StepTiming struct {
	// StepID is the step identifier.
	StepID string `json:"step_id" yaml:"step_id"`

	// StartTime is when the step started.
	StartTime time.Time `json:"start_time,omitempty" yaml:"start_time,omitempty"`

	// EndTime is when the step completed.
	EndTime time.Time `json:"end_time,omitempty" yaml:"end_time,omitempty"`

	// Duration is the step duration.
	Duration time.Duration `json:"duration,omitempty" yaml:"duration,omitempty"`

	// Status is the current step status.
	Status StepExecutionStatus `json:"status" yaml:"status"`
}

// ExecutionMetrics contains overall execution metrics.
type ExecutionMetrics struct {
	// TotalSteps is the total number of steps.
	TotalSteps int `json:"total_steps" yaml:"total_steps"`

	// CompletedSteps is the number of completed steps.
	CompletedSteps int `json:"completed_steps" yaml:"completed_steps"`

	// InProgressSteps is the number of in-progress steps.
	InProgressSteps int `json:"in_progress_steps" yaml:"in_progress_steps"`

	// BlockedSteps is the number of blocked steps.
	BlockedSteps int `json:"blocked_steps" yaml:"blocked_steps"`

	// PendingSteps is the number of pending steps.
	PendingSteps int `json:"pending_steps" yaml:"pending_steps"`

	// FailedSteps is the number of failed steps.
	FailedSteps int `json:"failed_steps" yaml:"failed_steps"`

	// ProgressPercent is the completion percentage.
	ProgressPercent float64 `json:"progress_percent" yaml:"progress_percent"`

	// ElapsedTime is the total elapsed time.
	ElapsedTime time.Duration `json:"elapsed_time" yaml:"elapsed_time"`

	// EstimatedRemaining is the estimated remaining time (if available).
	EstimatedRemaining time.Duration `json:"estimated_remaining,omitempty" yaml:"estimated_remaining,omitempty"`

	// PhaseMetrics contains per-phase metrics.
	PhaseMetrics map[string]*PhaseMetrics `json:"phase_metrics,omitempty" yaml:"phase_metrics,omitempty"`
}

// PhaseMetrics contains metrics for a specific phase.
type PhaseMetrics struct {
	// PhaseID is the phase identifier.
	PhaseID string `json:"phase_id" yaml:"phase_id"`

	// TotalSteps is the total steps in this phase.
	TotalSteps int `json:"total_steps" yaml:"total_steps"`

	// CompletedSteps is the completed steps in this phase.
	CompletedSteps int `json:"completed_steps" yaml:"completed_steps"`

	// ProgressPercent is the phase completion percentage.
	ProgressPercent float64 `json:"progress_percent" yaml:"progress_percent"`

	// Status is the phase status.
	Status PhaseStatus `json:"status" yaml:"status"`
}

// WorkflowToPIDL converts an AIDLC workflow to a pidl Protocol.
// This provides a basic protocol representation for visualization.
func WorkflowToPIDL(w *Workflow) *pidl.Protocol {
	protocol := &pidl.Protocol{
		ProtocolMeta: pidl.ProtocolMeta{
			ID:          w.Name,
			Name:        w.Name,
			Description: w.Description,
		},
		Entities: make([]pidl.Entity, 0),
		Phases:   make([]pidl.Phase, 0),
		Flows:    make([]pidl.Flow, 0),
	}

	// Convert phases
	for _, phase := range w.Phases {
		protocol.Phases = append(protocol.Phases, pidl.Phase{
			ID:          phase.ID,
			Name:        phase.Name,
			Description: phase.Description,
		})
	}

	// Convert nodes to entities
	for _, node := range w.Nodes {
		entity := pidl.Entity{
			ID:          node.ID,
			Name:        node.Name,
			Type:        pidl.EntityTypeOther,
			Description: node.Description,
		}
		protocol.Entities = append(protocol.Entities, entity)
	}

	// Convert edges to flows
	for _, edge := range w.Edges {
		protocol.Flows = append(protocol.Flows, pidl.Flow{
			From:   edge.From,
			To:     edge.To,
			Action: string(edge.Type),
			Mode:   pidl.FlowModeRequest,
		})
	}

	return protocol
}

// NewWorkflowExecutionContext creates a new execution context for an AIDLC workflow.
func NewWorkflowExecutionContext(w *Workflow) (*WorkflowExecutionContext, error) {
	if w == nil {
		return nil, fmt.Errorf("workflow cannot be nil")
	}

	protocol := WorkflowToPIDL(w)

	ctx := &WorkflowExecutionContext{
		Workflow:     w,
		Protocol:     protocol,
		StartTime:    time.Now(),
		StepTimings:  make(map[string]*StepTiming),
		StepStatus:   make(map[string]StepExecutionStatus),
		Dependencies: make(map[string][]string),
		Dependents:   make(map[string][]string),
	}

	// Build dependency graph
	for _, edge := range w.Edges {
		ctx.Dependencies[edge.To] = appendUniqueString(ctx.Dependencies[edge.To], edge.From)
		ctx.Dependents[edge.From] = appendUniqueString(ctx.Dependents[edge.From], edge.To)
	}

	// Also use node.DependsOn
	for nodeID, node := range w.Nodes {
		for _, depID := range node.DependsOn {
			ctx.Dependencies[nodeID] = appendUniqueString(ctx.Dependencies[nodeID], depID)
			ctx.Dependents[depID] = appendUniqueString(ctx.Dependents[depID], nodeID)
		}
	}

	// Calculate execution order
	ctx.ExecutionOrder = ctx.topologicalSort()

	// Initialize step timings and status
	for nodeID := range w.Nodes {
		ctx.StepTimings[nodeID] = &StepTiming{
			StepID: nodeID,
			Status: StepStatusPending,
		}
		ctx.StepStatus[nodeID] = StepStatusPending
	}

	// Sync initial status from workflow
	ctx.SyncFromWorkflow()

	return ctx, nil
}

// topologicalSort returns steps in topological execution order.
func (ctx *WorkflowExecutionContext) topologicalSort() []string {
	// Kahn's algorithm for topological sort
	inDegree := make(map[string]int)
	nodes := make(map[string]bool)

	// Initialize in-degrees
	for nodeID := range ctx.Workflow.Nodes {
		nodes[nodeID] = true
		inDegree[nodeID] = 0
	}

	// Calculate in-degrees from dependencies
	for nodeID := range nodes {
		for _, dep := range ctx.Dependencies[nodeID] {
			if nodes[dep] {
				inDegree[nodeID]++
			}
		}
	}

	// Start with nodes that have no dependencies
	var queue []string
	for id, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, id)
		}
	}
	sort.Strings(queue)

	var result []string
	for len(queue) > 0 {
		// Take first node
		node := queue[0]
		queue = queue[1:]
		result = append(result, node)

		// Reduce in-degree of dependents
		dependents := ctx.Dependents[node]
		sort.Strings(dependents)
		for _, dependent := range dependents {
			if !nodes[dependent] {
				continue
			}
			inDegree[dependent]--
			if inDegree[dependent] == 0 {
				queue = append(queue, dependent)
				sort.Strings(queue)
			}
		}
	}

	return result
}

// SyncFromWorkflow syncs the execution context from the workflow state.
func (ctx *WorkflowExecutionContext) SyncFromWorkflow() {
	for nodeID, node := range ctx.Workflow.Nodes {
		timing := ctx.StepTimings[nodeID]
		if timing == nil {
			timing = &StepTiming{StepID: nodeID}
			ctx.StepTimings[nodeID] = timing
		}

		var status StepExecutionStatus
		switch node.Status {
		case NodePending:
			status = StepStatusPending
		case NodeReady:
			status = StepStatusReady
		case NodeInProgress:
			status = StepStatusInProgress
			if timing.StartTime.IsZero() {
				timing.StartTime = time.Now()
			}
		case NodeCompleted:
			status = StepStatusCompleted
			if timing.EndTime.IsZero() {
				timing.EndTime = time.Now()
				if !timing.StartTime.IsZero() {
					timing.Duration = timing.EndTime.Sub(timing.StartTime)
				}
			}
		case NodeBlocked:
			status = StepStatusBlocked
		case NodeFailed:
			status = StepStatusFailed
		case NodeSkipped:
			status = StepStatusSkipped
		default:
			status = StepStatusPending
		}

		timing.Status = status
		ctx.StepStatus[nodeID] = status
	}
}

// StartStep marks a step as in progress.
func (ctx *WorkflowExecutionContext) StartStep(stepID string) error {
	timing, ok := ctx.StepTimings[stepID]
	if !ok {
		return fmt.Errorf("step not found: %s", stepID)
	}

	timing.StartTime = time.Now()
	timing.Status = StepStatusInProgress
	ctx.StepStatus[stepID] = StepStatusInProgress

	// Update the workflow
	if err := ctx.Workflow.UpdateNodeStatus(stepID, NodeInProgress, nil); err != nil {
		return err
	}

	return nil
}

// CompleteStep marks a step as completed.
func (ctx *WorkflowExecutionContext) CompleteStep(stepID string, score *QualityScore) error {
	timing, ok := ctx.StepTimings[stepID]
	if !ok {
		return fmt.Errorf("step not found: %s", stepID)
	}

	timing.EndTime = time.Now()
	timing.Duration = timing.EndTime.Sub(timing.StartTime)
	timing.Status = StepStatusCompleted
	ctx.StepStatus[stepID] = StepStatusCompleted

	// Update the workflow
	if err := ctx.Workflow.UpdateNodeStatus(stepID, NodeCompleted, score); err != nil {
		return err
	}

	return nil
}

// FailStep marks a step as failed.
func (ctx *WorkflowExecutionContext) FailStep(stepID string) error {
	timing, ok := ctx.StepTimings[stepID]
	if !ok {
		return fmt.Errorf("step not found: %s", stepID)
	}

	timing.EndTime = time.Now()
	timing.Duration = timing.EndTime.Sub(timing.StartTime)
	timing.Status = StepStatusFailed
	ctx.StepStatus[stepID] = StepStatusFailed

	// Update the workflow
	if err := ctx.Workflow.UpdateNodeStatus(stepID, NodeFailed, nil); err != nil {
		return err
	}

	return nil
}

// SkipStep marks a step as skipped.
func (ctx *WorkflowExecutionContext) SkipStep(stepID string) error {
	timing, ok := ctx.StepTimings[stepID]
	if !ok {
		return fmt.Errorf("step not found: %s", stepID)
	}

	timing.Status = StepStatusSkipped
	ctx.StepStatus[stepID] = StepStatusSkipped

	// Update the workflow
	if err := ctx.Workflow.UpdateNodeStatus(stepID, NodeSkipped, nil); err != nil {
		return err
	}

	return nil
}

// GetMetrics returns current execution metrics.
func (ctx *WorkflowExecutionContext) GetMetrics() *ExecutionMetrics {
	metrics := &ExecutionMetrics{
		TotalSteps:   len(ctx.StepTimings),
		ElapsedTime:  time.Since(ctx.StartTime),
		PhaseMetrics: make(map[string]*PhaseMetrics),
	}

	// Count step statuses
	for _, timing := range ctx.StepTimings {
		switch timing.Status {
		case StepStatusCompleted:
			metrics.CompletedSteps++
		case StepStatusInProgress:
			metrics.InProgressSteps++
		case StepStatusBlocked:
			metrics.BlockedSteps++
		case StepStatusPending, StepStatusReady:
			metrics.PendingSteps++
		case StepStatusFailed:
			metrics.FailedSteps++
		}
	}

	if metrics.TotalSteps > 0 {
		metrics.ProgressPercent = float64(metrics.CompletedSteps) / float64(metrics.TotalSteps) * 100
	}

	// Calculate per-phase metrics
	for _, phase := range ctx.Workflow.Phases {
		pm := &PhaseMetrics{
			PhaseID: phase.ID,
		}

		for _, nodeID := range phase.NodeIDs {
			pm.TotalSteps++
			if timing, ok := ctx.StepTimings[nodeID]; ok {
				if timing.Status == StepStatusCompleted {
					pm.CompletedSteps++
				}
			}
		}

		if pm.TotalSteps > 0 {
			pm.ProgressPercent = float64(pm.CompletedSteps) / float64(pm.TotalSteps) * 100
		}

		// Determine phase status
		if pm.CompletedSteps == pm.TotalSteps && pm.TotalSteps > 0 {
			pm.Status = PhaseStatusCompleted
		} else if pm.CompletedSteps > 0 {
			pm.Status = PhaseStatusInProgress
		} else {
			pm.Status = PhaseStatusPending
		}

		metrics.PhaseMetrics[phase.ID] = pm
	}

	// Estimate remaining time based on average step duration
	var totalDuration time.Duration
	var completedCount int
	for _, timing := range ctx.StepTimings {
		if timing.Duration > 0 {
			totalDuration += timing.Duration
			completedCount++
		}
	}

	if completedCount > 0 {
		avgDuration := totalDuration / time.Duration(completedCount)
		remainingSteps := metrics.TotalSteps - metrics.CompletedSteps - metrics.FailedSteps
		metrics.EstimatedRemaining = avgDuration * time.Duration(remainingSteps)
	}

	return metrics
}

// GetReadySteps returns steps that are ready to be executed.
func (ctx *WorkflowExecutionContext) GetReadySteps() []string {
	var ready []string
	for stepID, timing := range ctx.StepTimings {
		if timing.Status == StepStatusReady || timing.Status == StepStatusPending {
			// Check if all dependencies are satisfied
			allSatisfied := true
			for _, depID := range ctx.Dependencies[stepID] {
				if depStatus, ok := ctx.StepStatus[depID]; ok {
					if depStatus != StepStatusCompleted && depStatus != StepStatusSkipped {
						allSatisfied = false
						break
					}
				}
			}
			if allSatisfied {
				ready = append(ready, stepID)
			}
		}
	}
	sort.Strings(ready)
	return ready
}

// GetBlockedSteps returns steps that are blocked and their blocking reasons.
func (ctx *WorkflowExecutionContext) GetBlockedSteps() map[string][]string {
	blocked := make(map[string][]string)
	for stepID, timing := range ctx.StepTimings {
		if timing.Status == StepStatusPending || timing.Status == StepStatusBlocked {
			var reasons []string
			for _, depID := range ctx.Dependencies[stepID] {
				if depStatus, ok := ctx.StepStatus[depID]; ok {
					if depStatus != StepStatusCompleted && depStatus != StepStatusSkipped {
						reasons = append(reasons, fmt.Sprintf("waiting for %s (status: %s)", depID, depStatus))
					}
				}
			}
			if len(reasons) > 0 {
				blocked[stepID] = reasons
			}
		}
	}
	return blocked
}

// GetCriticalPath returns the steps on the critical path (longest execution path).
func (ctx *WorkflowExecutionContext) GetCriticalPath() []string {
	return ctx.ExecutionOrder
}

// IsComplete returns whether all steps are completed or failed/skipped.
func (ctx *WorkflowExecutionContext) IsComplete() bool {
	for _, timing := range ctx.StepTimings {
		if timing.Status == StepStatusPending ||
			timing.Status == StepStatusReady ||
			timing.Status == StepStatusInProgress ||
			timing.Status == StepStatusBlocked {
			return false
		}
	}
	return true
}

// GetStepTiming returns timing information for a specific step.
func (ctx *WorkflowExecutionContext) GetStepTiming(stepID string) (*StepTiming, bool) {
	timing, ok := ctx.StepTimings[stepID]
	return timing, ok
}

// Reset resets the execution context for a new run.
func (ctx *WorkflowExecutionContext) Reset() {
	ctx.StartTime = time.Now()
	ctx.ExecutionOrder = ctx.topologicalSort()

	for stepID := range ctx.StepTimings {
		ctx.StepTimings[stepID] = &StepTiming{
			StepID: stepID,
			Status: StepStatusPending,
		}
		ctx.StepStatus[stepID] = StepStatusPending
	}
}

// appendUniqueString appends an item to a slice if it doesn't already exist.
func appendUniqueString(slice []string, item string) []string {
	for _, s := range slice {
		if s == item {
			return slice
		}
	}
	return append(slice, item)
}
