package align

import (
	"strings"
	"testing"
	"time"
)

func TestNewResolutionEngine(t *testing.T) {
	engine := NewResolutionEngine()

	if engine == nil {
		t.Fatal("NewResolutionEngine should not return nil")
	}

	if len(engine.defaultStrategies) == 0 {
		t.Error("Should have default strategies")
	}
}

func TestResolutionEngine_GeneratePlan(t *testing.T) {
	engine := NewResolutionEngine()

	result := &AlignmentResult{
		Project: "test-project",
		Discrepancies: []Discrepancy{
			{
				ID:       "disc-1",
				Type:     DiscrepancyMissingFeature,
				Severity: SeverityHigh,
				SpecRef:  "FR-001",
			},
			{
				ID:       "disc-2",
				Type:     DiscrepancyUndocumentedCode,
				Severity: SeverityMedium,
				CodeRef:  "api/handler.go:50",
			},
			{
				ID:       "disc-3",
				Type:     DiscrepancyDiverged,
				Severity: SeverityCritical,
				SpecRef:  "FR-002",
				CodeRef:  "auth/login.go:100",
			},
		},
	}

	plan := engine.GeneratePlan(result)

	if plan.ProjectName != "test-project" {
		t.Errorf("ProjectName = %s, want test-project", plan.ProjectName)
	}

	if len(plan.Resolutions) != 3 {
		t.Errorf("Should have 3 resolutions, got %d", len(plan.Resolutions))
	}

	// Check strategies are assigned
	for _, res := range plan.Resolutions {
		if res.Strategy == "" {
			t.Errorf("Resolution %s has no strategy", res.DiscrepancyID)
		}
		if res.Status != StatusPending {
			t.Errorf("New resolution should have pending status")
		}
	}

	// Check summary
	if plan.Summary.TotalDiscrepancies != 3 {
		t.Errorf("Summary.TotalDiscrepancies = %d, want 3", plan.Summary.TotalDiscrepancies)
	}

	// Check priorities
	if len(plan.Priorities) != 3 {
		t.Errorf("Should have 3 prioritized actions")
	}

	// First action should be critical severity
	if plan.Priorities[0].Priority != "critical" {
		t.Error("First prioritized action should be critical")
	}
}

func TestResolutionEngine_SuggestResolution(t *testing.T) {
	engine := NewResolutionEngine()

	tests := []struct {
		discType DiscrepancyType
		expected ResolutionStrategy
	}{
		{DiscrepancyMissingFeature, StrategyUpdateCode},
		{DiscrepancyUndocumentedCode, StrategyAddSpec},
		{DiscrepancyDiverged, StrategyUpdateSpec},
		{DiscrepancyPartialImplementation, StrategyUpdateCode},
		{DiscrepancyBehaviorMismatch, StrategyUpdateCode},
	}

	for _, tt := range tests {
		disc := Discrepancy{
			ID:   "test",
			Type: tt.discType,
		}

		res := engine.suggestResolution(disc)

		if res.Strategy != tt.expected {
			t.Errorf("suggestResolution(%s) = %s, want %s", tt.discType, res.Strategy, tt.expected)
		}
	}
}

func TestResolutionEngine_CalculateSummary(t *testing.T) {
	engine := NewResolutionEngine()

	resolutions := []Resolution{
		{Strategy: StrategyUpdateSpec, Status: StatusPending},
		{Strategy: StrategyUpdateCode, Status: StatusResolved},
		{Strategy: StrategyAddSpec, Status: StatusPending},
		{Strategy: StrategyRemoveCode, Status: StatusPending},
		{Strategy: StrategyDefer, Status: StatusDeferred},
		{Strategy: StrategyIgnore, Status: StatusIgnored},
	}

	summary := engine.calculateSummary(resolutions)

	if summary.TotalDiscrepancies != 6 {
		t.Errorf("TotalDiscrepancies = %d, want 6", summary.TotalDiscrepancies)
	}

	if summary.UpdateSpec != 1 {
		t.Errorf("UpdateSpec = %d, want 1", summary.UpdateSpec)
	}

	if summary.UpdateCode != 1 {
		t.Errorf("UpdateCode = %d, want 1", summary.UpdateCode)
	}

	if summary.Pending != 3 {
		t.Errorf("Pending = %d, want 3", summary.Pending)
	}

	if summary.Resolved != 1 {
		t.Errorf("Resolved = %d, want 1", summary.Resolved)
	}
}

func TestEstimateEffort(t *testing.T) {
	tests := []struct {
		strategy ResolutionStrategy
		severity Severity
		expected string
	}{
		{StrategyIgnore, SeverityMedium, "small"},
		{StrategyAddSpec, SeverityHigh, "small"},
		{StrategyUpdateSpec, SeverityMedium, "medium"},
		{StrategyRemoveCode, SeverityLow, "medium"},
		{StrategyUpdateCode, SeverityCritical, "large"},
		{StrategyUpdateCode, SeverityMedium, "medium"},
		{StrategyDefer, SeverityLow, "small"},
	}

	for _, tt := range tests {
		disc := Discrepancy{Severity: tt.severity}
		result := estimateEffort(disc, tt.strategy)

		if result != tt.expected {
			t.Errorf("estimateEffort(%s, %s) = %s, want %s",
				tt.strategy, tt.severity, result, tt.expected)
		}
	}
}

func TestSortPrioritizedActions(t *testing.T) {
	actions := []PrioritizedAction{
		{Priority: "low", Effort: "small"},
		{Priority: "critical", Effort: "large"},
		{Priority: "high", Effort: "small"},
		{Priority: "critical", Effort: "small"},
	}

	sortPrioritizedActions(actions)

	// Critical should come first
	if actions[0].Priority != "critical" {
		t.Error("First action should be critical")
	}

	// Among criticals, smaller effort should come first
	if actions[0].Effort != "small" {
		t.Error("First critical action should have small effort")
	}

	// High should come before low
	highIdx := -1
	lowIdx := -1
	for i, a := range actions {
		if a.Priority == "high" && highIdx == -1 {
			highIdx = i
		}
		if a.Priority == "low" && lowIdx == -1 {
			lowIdx = i
		}
	}

	if highIdx > lowIdx {
		t.Error("High priority should come before low")
	}
}

func TestResolutionPlan_RenderMarkdown(t *testing.T) {
	now := time.Now()

	plan := &ResolutionPlan{
		ProjectName: "test-project",
		GeneratedAt: now,
		Summary: ResolutionSummary{
			TotalDiscrepancies: 3,
			UpdateSpec:         1,
			UpdateCode:         1,
			AddSpec:            1,
		},
		Resolutions: []Resolution{
			{
				DiscrepancyID: "disc-1",
				Strategy:      StrategyUpdateCode,
				Description:   "Update the implementation to implement FR-001",
				Status:        StatusPending,
			},
		},
		Priorities: []PrioritizedAction{
			{
				Order:         1,
				DiscrepancyID: "disc-1",
				Strategy:      StrategyUpdateCode,
				Priority:      "high",
				Effort:        "medium",
				Description:   "Update the implementation to implement FR-001",
			},
		},
	}

	md := plan.RenderMarkdown()

	if md == "" {
		t.Error("Should generate Markdown")
	}

	if !strings.Contains(md, "# Drift Resolution Plan") {
		t.Error("Should have title")
	}

	if !strings.Contains(md, "test-project") {
		t.Error("Should contain project name")
	}

	if !strings.Contains(md, "Summary") {
		t.Error("Should have summary section")
	}

	if !strings.Contains(md, "Prioritized Actions") {
		t.Error("Should have priorities section")
	}

	if !strings.Contains(md, "Resolution Details") {
		t.Error("Should have resolution details")
	}
}

func TestResolutionPlan_UpdateResolution(t *testing.T) {
	plan := &ResolutionPlan{
		Resolutions: []Resolution{
			{DiscrepancyID: "disc-1", Status: StatusPending},
			{DiscrepancyID: "disc-2", Status: StatusPending},
		},
	}

	// Update existing resolution
	err := plan.UpdateResolution("disc-1", StatusResolved, "Fixed the issue")
	if err != nil {
		t.Fatalf("UpdateResolution failed: %v", err)
	}

	// Verify update
	for _, res := range plan.Resolutions {
		if res.DiscrepancyID == "disc-1" {
			if res.Status != StatusResolved {
				t.Error("Status should be updated to resolved")
			}
			if res.Notes != "Fixed the issue" {
				t.Error("Notes should be updated")
			}
			if res.ResolvedAt == nil {
				t.Error("ResolvedAt should be set")
			}
		}
	}

	// Try to update non-existent resolution
	err = plan.UpdateResolution("disc-999", StatusResolved, "")
	if err == nil {
		t.Error("Should fail for non-existent resolution")
	}
}

func TestResolutionPlan_GetPendingResolutions(t *testing.T) {
	plan := &ResolutionPlan{
		Resolutions: []Resolution{
			{DiscrepancyID: "disc-1", Status: StatusPending},
			{DiscrepancyID: "disc-2", Status: StatusInProgress},
			{DiscrepancyID: "disc-3", Status: StatusResolved},
			{DiscrepancyID: "disc-4", Status: StatusIgnored},
		},
	}

	pending := plan.GetPendingResolutions()

	if len(pending) != 2 {
		t.Errorf("Should have 2 pending/in-progress resolutions, got %d", len(pending))
	}
}

func TestResolutionPlan_GetProgress(t *testing.T) {
	tests := []struct {
		resolutions []Resolution
		expected    float64
	}{
		{
			resolutions: []Resolution{},
			expected:    100,
		},
		{
			resolutions: []Resolution{
				{Status: StatusResolved},
				{Status: StatusPending},
			},
			expected: 50,
		},
		{
			resolutions: []Resolution{
				{Status: StatusResolved},
				{Status: StatusIgnored},
			},
			expected: 100,
		},
		{
			resolutions: []Resolution{
				{Status: StatusPending},
				{Status: StatusPending},
				{Status: StatusPending},
				{Status: StatusPending},
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		plan := &ResolutionPlan{Resolutions: tt.resolutions}
		progress := plan.GetProgress()

		if progress != tt.expected {
			t.Errorf("GetProgress() = %.0f, want %.0f", progress, tt.expected)
		}
	}
}

func TestCapitalize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"high", "High"},
		{"critical", "Critical"},
		{"", ""},
		{"A", "A"},
	}

	for _, tt := range tests {
		result := capitalize(tt.input)
		if result != tt.expected {
			t.Errorf("capitalize(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestStrategyIcon(t *testing.T) {
	strategies := []ResolutionStrategy{
		StrategyUpdateSpec,
		StrategyUpdateCode,
		StrategyAddSpec,
		StrategyRemoveCode,
		StrategyDefer,
		StrategyIgnore,
	}

	for _, s := range strategies {
		icon := strategyIcon(s)
		if icon == "" || icon == "❓" {
			// Only fail for known strategies returning unknown
			if s != "" {
				t.Errorf("strategyIcon(%s) should not return empty or unknown", s)
			}
		}
	}
}

func TestEffortBadge(t *testing.T) {
	tests := []struct {
		effort   string
		expected string
	}{
		{"small", "`S`"},
		{"medium", "`M`"},
		{"large", "`L`"},
		{"unknown", ""},
	}

	for _, tt := range tests {
		result := effortBadge(tt.effort)
		if result != tt.expected {
			t.Errorf("effortBadge(%q) = %q, want %q", tt.effort, result, tt.expected)
		}
	}
}
