// Package align provides drift resolution workflow capabilities.
package align

import (
	"fmt"
	"strings"
	"time"
)

// ResolutionStrategy defines how to resolve a discrepancy.
type ResolutionStrategy string

const (
	// StrategyUpdateSpec updates the specification to match implementation.
	StrategyUpdateSpec ResolutionStrategy = "update_spec"
	// StrategyUpdateCode updates the implementation to match specification.
	StrategyUpdateCode ResolutionStrategy = "update_code"
	// StrategyAddSpec adds missing specification for undocumented code.
	StrategyAddSpec ResolutionStrategy = "add_spec"
	// StrategyRemoveCode removes undocumented/deprecated code.
	StrategyRemoveCode ResolutionStrategy = "remove_code"
	// StrategyDefer defers resolution to a later time.
	StrategyDefer ResolutionStrategy = "defer"
	// StrategyIgnore explicitly ignores the discrepancy.
	StrategyIgnore ResolutionStrategy = "ignore"
)

// Resolution represents a resolution decision for a discrepancy.
type Resolution struct {
	DiscrepancyID string             `json:"discrepancy_id"`
	Strategy      ResolutionStrategy `json:"strategy"`
	Description   string             `json:"description,omitempty"`
	AssignedTo    string             `json:"assigned_to,omitempty"`
	DueDate       *time.Time         `json:"due_date,omitempty"`
	Status        ResolutionStatus   `json:"status"`
	CreatedAt     time.Time          `json:"created_at"`
	ResolvedAt    *time.Time         `json:"resolved_at,omitempty"`
	Notes         string             `json:"notes,omitempty"`
	LinkedTasks   []string           `json:"linked_tasks,omitempty"`
}

// ResolutionStatus indicates the current state of a resolution.
type ResolutionStatus string

const (
	StatusPending    ResolutionStatus = "pending"
	StatusInProgress ResolutionStatus = "in_progress"
	StatusResolved   ResolutionStatus = "resolved"
	StatusDeferred   ResolutionStatus = "deferred"
	StatusIgnored    ResolutionStatus = "ignored"
)

// ResolutionPlan contains a complete resolution plan for an alignment result.
type ResolutionPlan struct {
	ProjectName string              `json:"project_name"`
	GeneratedAt time.Time           `json:"generated_at"`
	AlignmentID string              `json:"alignment_id,omitempty"`
	Resolutions []Resolution        `json:"resolutions"`
	Summary     ResolutionSummary   `json:"summary"`
	Priorities  []PrioritizedAction `json:"priorities,omitempty"`
}

// ResolutionSummary provides aggregate statistics.
type ResolutionSummary struct {
	TotalDiscrepancies int `json:"total_discrepancies"`
	UpdateSpec         int `json:"update_spec"`
	UpdateCode         int `json:"update_code"`
	AddSpec            int `json:"add_spec"`
	RemoveCode         int `json:"remove_code"`
	Deferred           int `json:"deferred"`
	Ignored            int `json:"ignored"`
	Pending            int `json:"pending"`
	Resolved           int `json:"resolved"`
}

// PrioritizedAction represents a prioritized resolution action.
type PrioritizedAction struct {
	Order         int                `json:"order"`
	DiscrepancyID string             `json:"discrepancy_id"`
	Strategy      ResolutionStrategy `json:"strategy"`
	Priority      string             `json:"priority"` // critical, high, medium, low
	Effort        string             `json:"effort"`   // small, medium, large
	Description   string             `json:"description"`
	Blockers      []string           `json:"blockers,omitempty"`
}

// ResolutionEngine generates resolution plans from alignment results.
type ResolutionEngine struct {
	defaultStrategies map[DiscrepancyType]ResolutionStrategy
}

// NewResolutionEngine creates a new resolution engine with default strategies.
func NewResolutionEngine() *ResolutionEngine {
	return &ResolutionEngine{
		defaultStrategies: map[DiscrepancyType]ResolutionStrategy{
			DiscrepancyMissingFeature:        StrategyUpdateCode,
			DiscrepancyUndocumentedCode:      StrategyAddSpec,
			DiscrepancyDiverged:              StrategyUpdateSpec,
			DiscrepancyPartialImplementation: StrategyUpdateCode,
			DiscrepancyBehaviorMismatch:      StrategyUpdateCode,
		},
	}
}

// GeneratePlan creates a resolution plan from an alignment result.
func (e *ResolutionEngine) GeneratePlan(result *AlignmentResult) *ResolutionPlan {
	plan := &ResolutionPlan{
		ProjectName: result.Project,
		GeneratedAt: time.Now(),
		Resolutions: []Resolution{},
	}

	// Generate resolutions for each discrepancy
	for _, disc := range result.Discrepancies {
		resolution := e.suggestResolution(disc)
		plan.Resolutions = append(plan.Resolutions, resolution)
	}

	// Calculate summary
	plan.Summary = e.calculateSummary(plan.Resolutions)

	// Generate prioritized actions
	plan.Priorities = e.prioritize(result.Discrepancies, plan.Resolutions)

	return plan
}

// suggestResolution suggests a resolution for a discrepancy.
func (e *ResolutionEngine) suggestResolution(disc Discrepancy) Resolution {
	strategy := e.defaultStrategies[disc.Type]
	if strategy == "" {
		strategy = StrategyDefer
	}

	resolution := Resolution{
		DiscrepancyID: disc.ID,
		Strategy:      strategy,
		Description:   e.generateDescription(disc, strategy),
		Status:        StatusPending,
		CreatedAt:     time.Now(),
	}

	return resolution
}

// generateDescription creates a human-readable description of the resolution.
func (e *ResolutionEngine) generateDescription(disc Discrepancy, strategy ResolutionStrategy) string {
	var action string
	switch strategy {
	case StrategyUpdateSpec:
		action = "Update the specification"
	case StrategyUpdateCode:
		action = "Update the implementation"
	case StrategyAddSpec:
		action = "Add specification"
	case StrategyRemoveCode:
		action = "Remove or deprecate"
	case StrategyDefer:
		action = "Defer resolution"
	case StrategyIgnore:
		action = "Mark as intentional"
	}

	var target string
	switch disc.Type {
	case DiscrepancyMissingFeature:
		target = fmt.Sprintf("to implement %q", disc.SpecRef)
	case DiscrepancyUndocumentedCode:
		target = fmt.Sprintf("for %q", disc.CodeRef)
	case DiscrepancyDiverged:
		target = fmt.Sprintf("to reflect actual behavior of %q", disc.CodeRef)
	case DiscrepancyPartialImplementation:
		target = fmt.Sprintf("to complete %q", disc.SpecRef)
	case DiscrepancyBehaviorMismatch:
		target = fmt.Sprintf("to match specified behavior for %q", disc.SpecRef)
	default:
		target = disc.Description
	}

	return fmt.Sprintf("%s %s", action, target)
}

// calculateSummary computes aggregate statistics for resolutions.
func (e *ResolutionEngine) calculateSummary(resolutions []Resolution) ResolutionSummary {
	summary := ResolutionSummary{
		TotalDiscrepancies: len(resolutions),
	}

	for _, res := range resolutions {
		switch res.Strategy {
		case StrategyUpdateSpec:
			summary.UpdateSpec++
		case StrategyUpdateCode:
			summary.UpdateCode++
		case StrategyAddSpec:
			summary.AddSpec++
		case StrategyRemoveCode:
			summary.RemoveCode++
		case StrategyDefer:
			summary.Deferred++
		case StrategyIgnore:
			summary.Ignored++
		}

		switch res.Status {
		case StatusPending:
			summary.Pending++
		case StatusResolved:
			summary.Resolved++
		}
	}

	return summary
}

// prioritize orders resolutions by importance and effort.
func (e *ResolutionEngine) prioritize(discrepancies []Discrepancy, resolutions []Resolution) []PrioritizedAction {
	var actions []PrioritizedAction

	// Create a map for quick lookup
	discMap := make(map[string]Discrepancy)
	for _, d := range discrepancies {
		discMap[d.ID] = d
	}

	for _, res := range resolutions {
		disc, ok := discMap[res.DiscrepancyID]
		if !ok {
			continue
		}

		action := PrioritizedAction{
			DiscrepancyID: res.DiscrepancyID,
			Strategy:      res.Strategy,
			Description:   res.Description,
			Priority:      string(disc.Severity),
			Effort:        estimateEffort(disc, res.Strategy),
		}

		actions = append(actions, action)
	}

	// Sort by priority (critical first) then by effort (small first)
	sortPrioritizedActions(actions)

	// Assign order numbers
	for i := range actions {
		actions[i].Order = i + 1
	}

	return actions
}

// estimateEffort estimates the effort required for a resolution.
func estimateEffort(disc Discrepancy, strategy ResolutionStrategy) string {
	// Heuristics based on discrepancy type and strategy
	switch strategy {
	case StrategyIgnore:
		return "small"
	case StrategyAddSpec:
		return "small"
	case StrategyUpdateSpec:
		return "medium"
	case StrategyRemoveCode:
		return "medium"
	case StrategyUpdateCode:
		if disc.Severity == SeverityCritical {
			return "large"
		}
		return "medium"
	case StrategyDefer:
		return "small"
	}
	return "medium"
}

// sortPrioritizedActions sorts actions by priority and effort.
func sortPrioritizedActions(actions []PrioritizedAction) {
	// Custom sort: critical > high > medium > low, then small > medium > large
	priorityOrder := map[string]int{
		"critical": 0,
		"high":     1,
		"medium":   2,
		"low":      3,
	}
	effortOrder := map[string]int{
		"small":  0,
		"medium": 1,
		"large":  2,
	}

	// Simple bubble sort (good enough for typical sizes)
	for i := 0; i < len(actions); i++ {
		for j := i + 1; j < len(actions); j++ {
			pi := priorityOrder[actions[i].Priority]
			pj := priorityOrder[actions[j].Priority]
			if pj < pi || (pj == pi && effortOrder[actions[j].Effort] < effortOrder[actions[i].Effort]) {
				actions[i], actions[j] = actions[j], actions[i]
			}
		}
	}
}

// RenderMarkdown outputs the resolution plan as Markdown.
func (p *ResolutionPlan) RenderMarkdown() string {
	var sb strings.Builder

	sb.WriteString("# Drift Resolution Plan\n\n")
	sb.WriteString(fmt.Sprintf("**Project:** %s\n\n", p.ProjectName))
	sb.WriteString(fmt.Sprintf("**Generated:** %s\n\n", p.GeneratedAt.Format("2006-01-02 15:04:05")))

	// Summary
	sb.WriteString("## Summary\n\n")
	sb.WriteString("| Strategy | Count |\n")
	sb.WriteString("|----------|-------|\n")
	sb.WriteString(fmt.Sprintf("| Update Spec | %d |\n", p.Summary.UpdateSpec))
	sb.WriteString(fmt.Sprintf("| Update Code | %d |\n", p.Summary.UpdateCode))
	sb.WriteString(fmt.Sprintf("| Add Spec | %d |\n", p.Summary.AddSpec))
	sb.WriteString(fmt.Sprintf("| Remove Code | %d |\n", p.Summary.RemoveCode))
	sb.WriteString(fmt.Sprintf("| Deferred | %d |\n", p.Summary.Deferred))
	sb.WriteString(fmt.Sprintf("| Ignored | %d |\n", p.Summary.Ignored))
	sb.WriteString(fmt.Sprintf("| **Total** | **%d** |\n", p.Summary.TotalDiscrepancies))
	sb.WriteString("\n")

	// Prioritized actions
	if len(p.Priorities) > 0 {
		sb.WriteString("## Prioritized Actions\n\n")

		// Group by priority
		currentPriority := ""
		for _, action := range p.Priorities {
			if action.Priority != currentPriority {
				currentPriority = action.Priority
				sb.WriteString(fmt.Sprintf("### %s Priority\n\n", capitalize(currentPriority)))
			}

			strategyIcon := strategyIcon(action.Strategy)
			effortBadge := effortBadge(action.Effort)

			sb.WriteString(fmt.Sprintf("%d. %s %s %s\n",
				action.Order, strategyIcon, action.Description, effortBadge))
		}
		sb.WriteString("\n")
	}

	// Detailed resolutions
	sb.WriteString("## Resolution Details\n\n")
	for _, res := range p.Resolutions {
		sb.WriteString(fmt.Sprintf("### %s\n\n", res.DiscrepancyID))
		sb.WriteString(fmt.Sprintf("**Strategy:** %s\n\n", res.Strategy))
		sb.WriteString(fmt.Sprintf("**Status:** %s\n\n", res.Status))
		sb.WriteString(fmt.Sprintf("%s\n\n", res.Description))

		if res.AssignedTo != "" {
			sb.WriteString(fmt.Sprintf("**Assigned:** %s\n\n", res.AssignedTo))
		}
		if res.DueDate != nil {
			sb.WriteString(fmt.Sprintf("**Due:** %s\n\n", res.DueDate.Format("2006-01-02")))
		}
		if res.Notes != "" {
			sb.WriteString(fmt.Sprintf("**Notes:** %s\n\n", res.Notes))
		}
	}

	return sb.String()
}

// Helper functions

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func strategyIcon(strategy ResolutionStrategy) string {
	switch strategy {
	case StrategyUpdateSpec:
		return "📝"
	case StrategyUpdateCode:
		return "💻"
	case StrategyAddSpec:
		return "📄"
	case StrategyRemoveCode:
		return "🗑️"
	case StrategyDefer:
		return "⏳"
	case StrategyIgnore:
		return "🔇"
	default:
		return "❓"
	}
}

func effortBadge(effort string) string {
	switch effort {
	case "small":
		return "`S`"
	case "medium":
		return "`M`"
	case "large":
		return "`L`"
	default:
		return ""
	}
}

// UpdateResolution updates a resolution in the plan.
func (p *ResolutionPlan) UpdateResolution(discID string, status ResolutionStatus, notes string) error {
	for i := range p.Resolutions {
		if p.Resolutions[i].DiscrepancyID == discID {
			p.Resolutions[i].Status = status
			if notes != "" {
				p.Resolutions[i].Notes = notes
			}
			if status == StatusResolved {
				now := time.Now()
				p.Resolutions[i].ResolvedAt = &now
			}
			// Recalculate summary
			p.Summary = NewResolutionEngine().calculateSummary(p.Resolutions)
			return nil
		}
	}
	return fmt.Errorf("resolution for %q not found", discID)
}

// GetPendingResolutions returns all pending resolutions.
func (p *ResolutionPlan) GetPendingResolutions() []Resolution {
	var pending []Resolution
	for _, res := range p.Resolutions {
		if res.Status == StatusPending || res.Status == StatusInProgress {
			pending = append(pending, res)
		}
	}
	return pending
}

// GetProgress returns the completion percentage.
func (p *ResolutionPlan) GetProgress() float64 {
	if len(p.Resolutions) == 0 {
		return 100
	}

	resolved := 0
	for _, res := range p.Resolutions {
		if res.Status == StatusResolved || res.Status == StatusIgnored {
			resolved++
		}
	}

	return float64(resolved) / float64(len(p.Resolutions)) * 100
}
