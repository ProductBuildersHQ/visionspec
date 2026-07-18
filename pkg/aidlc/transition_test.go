package aidlc

import (
	"testing"
)

func TestPhaseStatus(t *testing.T) {
	statuses := []PhaseStatus{
		PhaseStatusPending,
		PhaseStatusInProgress,
		PhaseStatusCompleted,
		PhaseStatusBlocked,
	}

	for _, status := range statuses {
		if status == "" {
			t.Error("PhaseStatus constant is empty")
		}
	}
}

func TestDefaultTransitionRules(t *testing.T) {
	rules := DefaultTransitionRules()

	// Should have rules for construction and operations phases
	if _, ok := rules[PhaseConstruction]; !ok {
		t.Error("Missing rule for PhaseConstruction")
	}
	if _, ok := rules[PhaseOperations]; !ok {
		t.Error("Missing rule for PhaseOperations")
	}

	// Check construction rule
	constructionRule := rules[PhaseConstruction]
	if len(constructionRule.RequiredDocs) == 0 {
		t.Error("Construction phase should have required docs")
	}
	if constructionRule.MinScore <= 0 {
		t.Error("Construction phase should have minimum score")
	}
}

func TestWorkflow_CurrentPhase(t *testing.T) {
	workflow := DefaultWorkflow()

	// Initially should be inception
	phase := workflow.CurrentPhase()
	if phase != PhaseInception {
		t.Errorf("CurrentPhase() = %s, want %s", phase, PhaseInception)
	}

	// Start a construction document - should move to construction
	nodeID := string(DocImplementationPlan)
	if node, ok := workflow.Nodes[nodeID]; ok {
		node.Status = NodeInProgress
	}

	phase = workflow.CurrentPhase()
	if phase != PhaseConstruction {
		t.Errorf("CurrentPhase() = %s, want %s when construction doc in progress", phase, PhaseConstruction)
	}

	// Start an operations document - should move to operations
	nodeID = string(DocRunbook)
	if node, ok := workflow.Nodes[nodeID]; ok {
		node.Status = NodeInProgress
	}

	phase = workflow.CurrentPhase()
	if phase != PhaseOperations {
		t.Errorf("CurrentPhase() = %s, want %s when operations doc in progress", phase, PhaseOperations)
	}
}

func TestWorkflow_GetPhaseStatus(t *testing.T) {
	workflow := DefaultWorkflow()

	// Initially pending
	status := workflow.GetPhaseStatus(PhaseInception)
	if status != PhaseStatusPending {
		t.Errorf("GetPhaseStatus(inception) = %s, want %s", status, PhaseStatusPending)
	}

	// Start a document
	nodeID := string(DocVisionDocument)
	if node, ok := workflow.Nodes[nodeID]; ok {
		node.Status = NodeInProgress
	}

	status = workflow.GetPhaseStatus(PhaseInception)
	if status != PhaseStatusInProgress {
		t.Errorf("GetPhaseStatus(inception) = %s, want %s after starting doc", status, PhaseStatusInProgress)
	}

	// Complete all required inception docs
	for _, docType := range InceptionDocTypes() {
		nodeID := string(docType)
		if node, ok := workflow.Nodes[nodeID]; ok {
			if node.Required {
				node.Status = NodeCompleted
			}
		}
	}

	status = workflow.GetPhaseStatus(PhaseInception)
	if status != PhaseStatusCompleted {
		t.Errorf("GetPhaseStatus(inception) = %s, want %s after completing", status, PhaseStatusCompleted)
	}
}

func TestWorkflow_GetPhaseStatus_Blocked(t *testing.T) {
	workflow := DefaultWorkflow()

	// Set a required node to failed
	nodeID := string(DocVisionDocument)
	if node, ok := workflow.Nodes[nodeID]; ok {
		node.Status = NodeFailed
	}

	status := workflow.GetPhaseStatus(PhaseInception)
	if status != PhaseStatusBlocked {
		t.Errorf("GetPhaseStatus(inception) = %s, want %s when doc failed", status, PhaseStatusBlocked)
	}
}

func TestWorkflow_CanTransitionTo(t *testing.T) {
	workflow := DefaultWorkflow()
	rules := DefaultTransitionRules()

	// Cannot transition forward without completing current phase
	result := workflow.CanTransitionTo(PhaseConstruction, rules)
	if result.Success {
		t.Error("Should not be able to transition without completing inception")
	}

	// Complete inception phase
	for _, docType := range InceptionDocTypes() {
		nodeID := string(docType)
		if node, ok := workflow.Nodes[nodeID]; ok {
			node.Status = NodeCompleted
			node.Score = &QualityScore{Score: 0.8}
		}
	}

	// Now should be able to transition
	result = workflow.CanTransitionTo(PhaseConstruction, rules)
	if !result.Success {
		t.Errorf("Should be able to transition after completing inception: %v", result.BlockingIssues)
	}
}

func TestWorkflow_CanTransitionTo_BackwardsNotAllowed(t *testing.T) {
	workflow := DefaultWorkflow()
	rules := DefaultTransitionRules()

	// Complete inception
	for _, docType := range InceptionDocTypes() {
		nodeID := string(docType)
		if node, ok := workflow.Nodes[nodeID]; ok {
			node.Status = NodeCompleted
		}
	}

	// Try to transition backwards
	result := workflow.CanTransitionTo(PhaseInception, rules)
	if result.Success {
		t.Error("Should not be able to transition backwards")
	}
	if len(result.BlockingIssues) == 0 {
		t.Error("Should have blocking issues for backwards transition")
	}
}

func TestWorkflow_CanTransitionTo_SkipNotAllowed(t *testing.T) {
	workflow := DefaultWorkflow()
	rules := DefaultTransitionRules()

	// Try to skip from inception to operations
	result := workflow.CanTransitionTo(PhaseOperations, rules)
	if result.Success {
		t.Error("Should not be able to skip phases")
	}
}

func TestWorkflow_TransitionTo(t *testing.T) {
	workflow := DefaultWorkflow()
	rules := DefaultTransitionRules()

	// Complete inception phase
	for _, docType := range InceptionDocTypes() {
		nodeID := string(docType)
		if node, ok := workflow.Nodes[nodeID]; ok {
			node.Status = NodeCompleted
			node.Score = &QualityScore{Score: 0.8}
		}
	}

	// Transition
	result, err := workflow.TransitionTo(PhaseConstruction, rules)
	if err != nil {
		t.Fatalf("TransitionTo() error = %v", err)
	}
	if !result.Success {
		t.Errorf("TransitionTo() success = false, want true")
	}
}

func TestWorkflow_TransitionTo_Blocked(t *testing.T) {
	workflow := DefaultWorkflow()
	rules := DefaultTransitionRules()

	// Don't complete inception
	result, err := workflow.TransitionTo(PhaseConstruction, rules)
	if err == nil {
		t.Error("TransitionTo() should return error when blocked")
	}
	if result.Success {
		t.Error("TransitionTo() success should be false when blocked")
	}
}

func TestWorkflow_GetPhaseRequirements(t *testing.T) {
	workflow := DefaultWorkflow()

	reqs := workflow.GetPhaseRequirements(PhaseInception)

	if reqs.Phase != PhaseInception {
		t.Errorf("Phase = %s, want %s", reqs.Phase, PhaseInception)
	}
	if len(reqs.RequiredDocs) == 0 {
		t.Error("Should have required docs")
	}
	if reqs.ProgressPercent != 0 {
		t.Error("Progress should be 0 initially")
	}
	if reqs.CanAdvance {
		t.Error("Should not be able to advance without completing docs")
	}
}

func TestWorkflow_AllPhaseRequirements(t *testing.T) {
	workflow := DefaultWorkflow()

	allReqs := workflow.AllPhaseRequirements()

	if len(allReqs) != 3 {
		t.Errorf("Should have requirements for 3 phases, got %d", len(allReqs))
	}
}

func TestWorkflow_ValidatePhaseTransition(t *testing.T) {
	workflow := DefaultWorkflow()

	// Validate transition to construction without completing inception
	issues := workflow.ValidatePhaseTransition(PhaseConstruction)
	if len(issues) == 0 {
		t.Error("Should have issues when inception not complete")
	}

	// Complete inception
	for _, docType := range InceptionDocTypes() {
		nodeID := string(docType)
		if node, ok := workflow.Nodes[nodeID]; ok {
			node.Status = NodeCompleted
		}
	}

	issues = workflow.ValidatePhaseTransition(PhaseConstruction)
	if len(issues) != 0 {
		t.Errorf("Should have no issues after completing inception: %v", issues)
	}
}

func TestNewTransitionLog(t *testing.T) {
	log := NewTransitionLog()
	if log == nil {
		t.Fatal("NewTransitionLog() returned nil")
	}
	if log.Entries == nil {
		t.Error("Entries should be initialized")
	}
	if len(log.Entries) != 0 {
		t.Error("Entries should be empty initially")
	}
}

func TestTransitionLog_AddEntry(t *testing.T) {
	log := NewTransitionLog()

	log.AddEntry(PhaseInception, PhaseConstruction, "admin", "Approved for construction")

	if len(log.Entries) != 1 {
		t.Fatalf("Should have 1 entry, got %d", len(log.Entries))
	}

	entry := log.Entries[0]
	if entry.FromPhase != PhaseInception {
		t.Errorf("FromPhase = %s, want %s", entry.FromPhase, PhaseInception)
	}
	if entry.ToPhase != PhaseConstruction {
		t.Errorf("ToPhase = %s, want %s", entry.ToPhase, PhaseConstruction)
	}
	if entry.ApprovedBy != "admin" {
		t.Errorf("ApprovedBy = %s, want admin", entry.ApprovedBy)
	}
	if entry.Notes != "Approved for construction" {
		t.Errorf("Notes = %s, want 'Approved for construction'", entry.Notes)
	}
	if entry.Timestamp.IsZero() {
		t.Error("Timestamp should be set")
	}
}

func TestTransitionLog_LatestEntry(t *testing.T) {
	log := NewTransitionLog()

	// Empty log
	if log.LatestEntry() != nil {
		t.Error("LatestEntry should be nil for empty log")
	}

	// Add entries
	log.AddEntry(PhaseInception, PhaseConstruction, "user1", "")
	log.AddEntry(PhaseConstruction, PhaseOperations, "user2", "")

	latest := log.LatestEntry()
	if latest == nil {
		t.Fatal("LatestEntry should not be nil")
	}
	if latest.ToPhase != PhaseOperations {
		t.Errorf("Latest entry ToPhase = %s, want %s", latest.ToPhase, PhaseOperations)
	}
}

func TestTransitionLog_EntriesForPhase(t *testing.T) {
	log := NewTransitionLog()

	log.AddEntry(PhaseInception, PhaseConstruction, "user1", "")
	log.AddEntry(PhaseConstruction, PhaseOperations, "user2", "")

	// Entries involving inception
	entries := log.EntriesForPhase(PhaseInception)
	if len(entries) != 1 {
		t.Errorf("Should have 1 entry for inception, got %d", len(entries))
	}

	// Entries involving construction
	entries = log.EntriesForPhase(PhaseConstruction)
	if len(entries) != 2 {
		t.Errorf("Should have 2 entries for construction, got %d", len(entries))
	}
}
