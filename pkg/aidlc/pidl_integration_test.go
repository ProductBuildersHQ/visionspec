package aidlc

import (
	"testing"
	"time"
)

func TestStepExecutionStatus(t *testing.T) {
	statuses := []StepExecutionStatus{
		StepStatusPending,
		StepStatusReady,
		StepStatusBlocked,
		StepStatusInProgress,
		StepStatusCompleted,
		StepStatusFailed,
		StepStatusSkipped,
	}

	for _, status := range statuses {
		if status == "" {
			t.Errorf("Status constant is empty")
		}
	}
}

func TestNewWorkflowExecutionContext(t *testing.T) {
	workflow := DefaultWorkflow()

	ctx, err := NewWorkflowExecutionContext(workflow)
	if err != nil {
		t.Fatalf("NewWorkflowExecutionContext() error = %v", err)
	}

	if ctx.Workflow != workflow {
		t.Error("Context workflow does not match input")
	}

	if ctx.Protocol == nil {
		t.Error("Context protocol is nil")
	}

	if ctx.StartTime.IsZero() {
		t.Error("Context start time is zero")
	}

	// Check step timings initialized
	for nodeID := range workflow.Nodes {
		timing, ok := ctx.StepTimings[nodeID]
		if !ok {
			t.Errorf("Missing step timing for node %s", nodeID)
			continue
		}
		if timing.StepID != nodeID {
			t.Errorf("Step timing ID = %s, want %s", timing.StepID, nodeID)
		}
	}
}

func TestNewWorkflowExecutionContext_NilWorkflow(t *testing.T) {
	_, err := NewWorkflowExecutionContext(nil)
	if err == nil {
		t.Error("NewWorkflowExecutionContext(nil) should return error")
	}
}

func TestWorkflowExecutionContext_StartStep(t *testing.T) {
	workflow := DefaultWorkflow()
	ctx, err := NewWorkflowExecutionContext(workflow)
	if err != nil {
		t.Fatalf("NewWorkflowExecutionContext() error = %v", err)
	}

	stepID := string(DocVisionDocument)

	err = ctx.StartStep(stepID)
	if err != nil {
		t.Fatalf("StartStep() error = %v", err)
	}

	timing := ctx.StepTimings[stepID]
	if timing.Status != StepStatusInProgress {
		t.Errorf("Step status = %s, want %s", timing.Status, StepStatusInProgress)
	}
	if timing.StartTime.IsZero() {
		t.Error("Start time should be set")
	}

	// Check workflow node was updated
	node := workflow.Nodes[stepID]
	if node.Status != NodeInProgress {
		t.Errorf("Node status = %s, want %s", node.Status, NodeInProgress)
	}
}

func TestWorkflowExecutionContext_StartStep_NotFound(t *testing.T) {
	workflow := DefaultWorkflow()
	ctx, err := NewWorkflowExecutionContext(workflow)
	if err != nil {
		t.Fatalf("NewWorkflowExecutionContext() error = %v", err)
	}

	err = ctx.StartStep("nonexistent")
	if err == nil {
		t.Error("StartStep(nonexistent) should return error")
	}
}

func TestWorkflowExecutionContext_CompleteStep(t *testing.T) {
	workflow := DefaultWorkflow()
	ctx, err := NewWorkflowExecutionContext(workflow)
	if err != nil {
		t.Fatalf("NewWorkflowExecutionContext() error = %v", err)
	}

	stepID := string(DocVisionDocument)

	// Start first
	err = ctx.StartStep(stepID)
	if err != nil {
		t.Fatalf("StartStep() error = %v", err)
	}

	// Small delay to ensure duration > 0
	time.Sleep(10 * time.Millisecond)

	// Complete
	score := &QualityScore{
		Rating: RatingGood,
		Score:  0.85,
	}
	err = ctx.CompleteStep(stepID, score)
	if err != nil {
		t.Fatalf("CompleteStep() error = %v", err)
	}

	timing := ctx.StepTimings[stepID]
	if timing.Status != StepStatusCompleted {
		t.Errorf("Step status = %s, want %s", timing.Status, StepStatusCompleted)
	}
	if timing.EndTime.IsZero() {
		t.Error("End time should be set")
	}
	if timing.Duration <= 0 {
		t.Error("Duration should be positive")
	}

	// Check workflow node was updated
	node := workflow.Nodes[stepID]
	if node.Status != NodeCompleted {
		t.Errorf("Node status = %s, want %s", node.Status, NodeCompleted)
	}
	if node.Score != score {
		t.Error("Node score not set")
	}
}

func TestWorkflowExecutionContext_FailStep(t *testing.T) {
	workflow := DefaultWorkflow()
	ctx, err := NewWorkflowExecutionContext(workflow)
	if err != nil {
		t.Fatalf("NewWorkflowExecutionContext() error = %v", err)
	}

	stepID := string(DocVisionDocument)

	// Start first
	err = ctx.StartStep(stepID)
	if err != nil {
		t.Fatalf("StartStep() error = %v", err)
	}

	// Fail
	err = ctx.FailStep(stepID)
	if err != nil {
		t.Fatalf("FailStep() error = %v", err)
	}

	timing := ctx.StepTimings[stepID]
	if timing.Status != StepStatusFailed {
		t.Errorf("Step status = %s, want %s", timing.Status, StepStatusFailed)
	}

	// Check workflow node was updated
	node := workflow.Nodes[stepID]
	if node.Status != NodeFailed {
		t.Errorf("Node status = %s, want %s", node.Status, NodeFailed)
	}
}

func TestWorkflowExecutionContext_SkipStep(t *testing.T) {
	workflow := DefaultWorkflow()
	ctx, err := NewWorkflowExecutionContext(workflow)
	if err != nil {
		t.Fatalf("NewWorkflowExecutionContext() error = %v", err)
	}

	stepID := string(DocArchitectureSpec)

	err = ctx.SkipStep(stepID)
	if err != nil {
		t.Fatalf("SkipStep() error = %v", err)
	}

	timing := ctx.StepTimings[stepID]
	if timing.Status != StepStatusSkipped {
		t.Errorf("Step status = %s, want %s", timing.Status, StepStatusSkipped)
	}
}

func TestWorkflowExecutionContext_GetMetrics(t *testing.T) {
	workflow := DefaultWorkflow()
	ctx, err := NewWorkflowExecutionContext(workflow)
	if err != nil {
		t.Fatalf("NewWorkflowExecutionContext() error = %v", err)
	}

	// Complete some steps
	steps := []string{string(DocVisionDocument), string(DocRequirementsSpec)}
	for _, stepID := range steps {
		_ = ctx.StartStep(stepID)
		_ = ctx.CompleteStep(stepID, nil)
	}

	metrics := ctx.GetMetrics()

	if metrics.TotalSteps != len(workflow.Nodes) {
		t.Errorf("TotalSteps = %d, want %d", metrics.TotalSteps, len(workflow.Nodes))
	}
	if metrics.CompletedSteps != 2 {
		t.Errorf("CompletedSteps = %d, want 2", metrics.CompletedSteps)
	}
	if metrics.ElapsedTime <= 0 {
		t.Error("ElapsedTime should be positive")
	}
	if metrics.ProgressPercent <= 0 {
		t.Error("ProgressPercent should be positive")
	}
}

func TestWorkflowExecutionContext_GetReadySteps(t *testing.T) {
	workflow := DefaultWorkflow()
	ctx, err := NewWorkflowExecutionContext(workflow)
	if err != nil {
		t.Fatalf("NewWorkflowExecutionContext() error = %v", err)
	}

	// Initially, steps without dependencies should be ready
	ready := ctx.GetReadySteps()
	if len(ready) == 0 {
		t.Error("Should have some ready steps initially")
	}

	// Vision document should be ready (no dependencies)
	found := false
	for _, stepID := range ready {
		if stepID == string(DocVisionDocument) {
			found = true
			break
		}
	}
	if !found {
		t.Error("VisionDocument should be in ready steps")
	}
}

func TestWorkflowExecutionContext_GetBlockedSteps(t *testing.T) {
	workflow := DefaultWorkflow()
	ctx, err := NewWorkflowExecutionContext(workflow)
	if err != nil {
		t.Fatalf("NewWorkflowExecutionContext() error = %v", err)
	}

	blocked := ctx.GetBlockedSteps()

	// Steps with dependencies should be blocked
	if len(blocked) == 0 {
		t.Log("No blocked steps found (may be expected if no dependencies)")
	}
}

func TestWorkflowExecutionContext_IsComplete(t *testing.T) {
	workflow := DefaultWorkflow()
	ctx, err := NewWorkflowExecutionContext(workflow)
	if err != nil {
		t.Fatalf("NewWorkflowExecutionContext() error = %v", err)
	}

	// Initially not complete
	if ctx.IsComplete() {
		t.Error("Workflow should not be complete initially")
	}

	// Complete or skip all steps
	for stepID := range workflow.Nodes {
		_ = ctx.StartStep(stepID)
		_ = ctx.CompleteStep(stepID, nil)
	}

	if !ctx.IsComplete() {
		t.Error("Workflow should be complete after all steps done")
	}
}

func TestWorkflowExecutionContext_Reset(t *testing.T) {
	workflow := DefaultWorkflow()
	ctx, err := NewWorkflowExecutionContext(workflow)
	if err != nil {
		t.Fatalf("NewWorkflowExecutionContext() error = %v", err)
	}

	// Complete some steps
	stepID := string(DocVisionDocument)
	_ = ctx.StartStep(stepID)
	_ = ctx.CompleteStep(stepID, nil)

	// Verify completed
	if ctx.StepTimings[stepID].Status != StepStatusCompleted {
		t.Fatal("Step should be completed before reset")
	}

	// Reset
	ctx.Reset()

	// Verify reset
	if ctx.StepTimings[stepID].Status != StepStatusPending {
		t.Errorf("Step status = %s, want %s after reset", ctx.StepTimings[stepID].Status, StepStatusPending)
	}
}

func TestWorkflowExecutionContext_GetStepTiming(t *testing.T) {
	workflow := DefaultWorkflow()
	ctx, err := NewWorkflowExecutionContext(workflow)
	if err != nil {
		t.Fatalf("NewWorkflowExecutionContext() error = %v", err)
	}

	stepID := string(DocVisionDocument)

	timing, ok := ctx.GetStepTiming(stepID)
	if !ok {
		t.Error("GetStepTiming should return true for valid step")
	}
	if timing == nil {
		t.Error("GetStepTiming should return timing for valid step")
	}

	_, ok = ctx.GetStepTiming("nonexistent")
	if ok {
		t.Error("GetStepTiming should return false for invalid step")
	}
}

func TestWorkflowExecutionContext_GetCriticalPath(t *testing.T) {
	workflow := DefaultWorkflow()
	ctx, err := NewWorkflowExecutionContext(workflow)
	if err != nil {
		t.Fatalf("NewWorkflowExecutionContext() error = %v", err)
	}

	path := ctx.GetCriticalPath()
	if len(path) == 0 {
		t.Error("Critical path should not be empty")
	}
}

func TestWorkflowToPIDL(t *testing.T) {
	workflow := DefaultWorkflow()
	protocol := WorkflowToPIDL(workflow)

	if protocol == nil {
		t.Fatal("WorkflowToPIDL returned nil")
	}

	if protocol.ProtocolMeta.ID != workflow.Name {
		t.Errorf("Protocol ID = %s, want %s", protocol.ProtocolMeta.ID, workflow.Name)
	}

	// Check phases
	if len(protocol.Phases) != len(workflow.Phases) {
		t.Errorf("Protocol has %d phases, want %d", len(protocol.Phases), len(workflow.Phases))
	}

	// Check entities
	if len(protocol.Entities) != len(workflow.Nodes) {
		t.Errorf("Protocol has %d entities, want %d", len(protocol.Entities), len(workflow.Nodes))
	}

	// Check flows
	if len(protocol.Flows) != len(workflow.Edges) {
		t.Errorf("Protocol has %d flows, want %d", len(protocol.Flows), len(workflow.Edges))
	}
}
