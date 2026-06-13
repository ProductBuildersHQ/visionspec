// Package decisions provides decision log persistence and rationale tracking.
//
// This package captures architectural decisions, their rationale, and the
// relationships between decisions to create a comprehensive decision history.
package decisions

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Decision represents a single architectural or design decision.
type Decision struct {
	ID             string            `json:"id" yaml:"id"`
	Title          string            `json:"title" yaml:"title"`
	Status         DecisionStatus    `json:"status" yaml:"status"`
	Date           time.Time         `json:"date" yaml:"date"`
	DecisionMakers []string          `json:"decision_makers,omitempty" yaml:"decision_makers,omitempty"`
	Context        string            `json:"context" yaml:"context"`
	Decision       string            `json:"decision" yaml:"decision"`
	Rationale      string            `json:"rationale" yaml:"rationale"`
	Consequences   []string          `json:"consequences,omitempty" yaml:"consequences,omitempty"`
	Alternatives   []Alternative     `json:"alternatives,omitempty" yaml:"alternatives,omitempty"`
	Related        []string          `json:"related,omitempty" yaml:"related,omitempty"`       // Related decision IDs
	Supersedes     string            `json:"supersedes,omitempty" yaml:"supersedes,omitempty"` // ID of decision this replaces
	SupersededBy   string            `json:"superseded_by,omitempty" yaml:"superseded_by,omitempty"`
	Tags           []string          `json:"tags,omitempty" yaml:"tags,omitempty"`
	Project        string            `json:"project,omitempty" yaml:"project,omitempty"`
	Specs          []string          `json:"specs,omitempty" yaml:"specs,omitempty"` // Affected specs
	Metadata       map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// DecisionStatus indicates the current state of a decision.
type DecisionStatus string

const (
	StatusProposed   DecisionStatus = "proposed"
	StatusAccepted   DecisionStatus = "accepted"
	StatusDeprecated DecisionStatus = "deprecated"
	StatusSuperseded DecisionStatus = "superseded"
	StatusRejected   DecisionStatus = "rejected"
)

// Alternative represents an option that was considered but not chosen.
type Alternative struct {
	Title       string   `json:"title" yaml:"title"`
	Description string   `json:"description,omitempty" yaml:"description,omitempty"`
	Pros        []string `json:"pros,omitempty" yaml:"pros,omitempty"`
	Cons        []string `json:"cons,omitempty" yaml:"cons,omitempty"`
	RejectedFor string   `json:"rejected_for,omitempty" yaml:"rejected_for,omitempty"`
}

// DecisionLog manages a collection of decisions.
type DecisionLog struct {
	path      string
	decisions map[string]*Decision
	order     []string // Maintains insertion order
}

// NewDecisionLog creates a new decision log at the given path.
func NewDecisionLog(path string) (*DecisionLog, error) {
	log := &DecisionLog{
		path:      path,
		decisions: make(map[string]*Decision),
		order:     []string{},
	}

	if err := log.load(); err != nil {
		return nil, err
	}

	return log, nil
}

// load reads decisions from the log file.
func (l *DecisionLog) load() error {
	data, err := os.ReadFile(l.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Empty log is valid
		}
		return fmt.Errorf("reading decision log: %w", err)
	}

	var decisions []Decision
	if err := yaml.Unmarshal(data, &decisions); err != nil {
		return fmt.Errorf("parsing decision log: %w", err)
	}

	for i := range decisions {
		d := &decisions[i]
		l.decisions[d.ID] = d
		l.order = append(l.order, d.ID)
	}

	return nil
}

// Save persists the decision log to disk.
func (l *DecisionLog) Save() error {
	dir := filepath.Dir(l.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	// Preserve order
	decisions := make([]Decision, 0, len(l.decisions))
	for _, id := range l.order {
		if d, ok := l.decisions[id]; ok {
			decisions = append(decisions, *d)
		}
	}

	data, err := yaml.Marshal(decisions)
	if err != nil {
		return fmt.Errorf("marshaling decisions: %w", err)
	}

	return os.WriteFile(l.path, data, 0600)
}

// Add adds a new decision to the log.
func (l *DecisionLog) Add(d *Decision) error {
	if d.ID == "" {
		d.ID = l.generateID()
	}

	if _, exists := l.decisions[d.ID]; exists {
		return fmt.Errorf("decision %q already exists", d.ID)
	}

	if d.Date.IsZero() {
		d.Date = time.Now()
	}

	l.decisions[d.ID] = d
	l.order = append(l.order, d.ID)

	return nil
}

// Update updates an existing decision.
func (l *DecisionLog) Update(d *Decision) error {
	if _, exists := l.decisions[d.ID]; !exists {
		return fmt.Errorf("decision %q not found", d.ID)
	}

	l.decisions[d.ID] = d
	return nil
}

// Get retrieves a decision by ID.
func (l *DecisionLog) Get(id string) (*Decision, error) {
	d, ok := l.decisions[id]
	if !ok {
		return nil, fmt.Errorf("decision %q not found", id)
	}
	return d, nil
}

// List returns all decisions in order.
func (l *DecisionLog) List() []*Decision {
	result := make([]*Decision, 0, len(l.order))
	for _, id := range l.order {
		if d, ok := l.decisions[id]; ok {
			result = append(result, d)
		}
	}
	return result
}

// ListByStatus returns decisions filtered by status.
func (l *DecisionLog) ListByStatus(status DecisionStatus) []*Decision {
	var result []*Decision
	for _, id := range l.order {
		if d, ok := l.decisions[id]; ok && d.Status == status {
			result = append(result, d)
		}
	}
	return result
}

// ListByTag returns decisions filtered by tag.
func (l *DecisionLog) ListByTag(tag string) []*Decision {
	var result []*Decision
	for _, id := range l.order {
		d, ok := l.decisions[id]
		if !ok {
			continue
		}
		for _, t := range d.Tags {
			if t == tag {
				result = append(result, d)
				break
			}
		}
	}
	return result
}

// ListByProject returns decisions filtered by project.
func (l *DecisionLog) ListByProject(project string) []*Decision {
	var result []*Decision
	for _, id := range l.order {
		if d, ok := l.decisions[id]; ok && d.Project == project {
			result = append(result, d)
		}
	}
	return result
}

// Supersede marks a decision as superseded by another.
func (l *DecisionLog) Supersede(oldID, newID string) error {
	old, ok := l.decisions[oldID]
	if !ok {
		return fmt.Errorf("decision %q not found", oldID)
	}

	newD, ok := l.decisions[newID]
	if !ok {
		return fmt.Errorf("decision %q not found", newID)
	}

	old.Status = StatusSuperseded
	old.SupersededBy = newID
	newD.Supersedes = oldID

	return nil
}

// generateID creates a new unique decision ID.
func (l *DecisionLog) generateID() string {
	// Format: ADR-NNNN where NNNN is sequential
	maxNum := 0
	for id := range l.decisions {
		if strings.HasPrefix(id, "ADR-") {
			var num int
			if _, err := fmt.Sscanf(id, "ADR-%d", &num); err == nil {
				if num > maxNum {
					maxNum = num
				}
			}
		}
	}
	return fmt.Sprintf("ADR-%04d", maxNum+1)
}

// RenderMarkdown renders a single decision as Markdown (ADR format).
func (d *Decision) RenderMarkdown() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %s: %s\n\n", d.ID, d.Title))
	sb.WriteString(fmt.Sprintf("**Status:** %s\n\n", formatStatus(d.Status)))
	sb.WriteString(fmt.Sprintf("**Date:** %s\n\n", d.Date.Format("2006-01-02")))

	if len(d.DecisionMakers) > 0 {
		sb.WriteString(fmt.Sprintf("**Decision Makers:** %s\n\n", strings.Join(d.DecisionMakers, ", ")))
	}

	if d.Project != "" {
		sb.WriteString(fmt.Sprintf("**Project:** %s\n\n", d.Project))
	}

	if len(d.Tags) > 0 {
		sb.WriteString(fmt.Sprintf("**Tags:** %s\n\n", strings.Join(d.Tags, ", ")))
	}

	sb.WriteString("## Context\n\n")
	sb.WriteString(d.Context)
	sb.WriteString("\n\n")

	sb.WriteString("## Decision\n\n")
	sb.WriteString(d.Decision)
	sb.WriteString("\n\n")

	sb.WriteString("## Rationale\n\n")
	sb.WriteString(d.Rationale)
	sb.WriteString("\n\n")

	if len(d.Consequences) > 0 {
		sb.WriteString("## Consequences\n\n")
		for _, c := range d.Consequences {
			sb.WriteString(fmt.Sprintf("- %s\n", c))
		}
		sb.WriteString("\n")
	}

	if len(d.Alternatives) > 0 {
		sb.WriteString("## Alternatives Considered\n\n")
		for _, alt := range d.Alternatives {
			sb.WriteString(fmt.Sprintf("### %s\n\n", alt.Title))
			if alt.Description != "" {
				sb.WriteString(alt.Description + "\n\n")
			}
			if len(alt.Pros) > 0 {
				sb.WriteString("**Pros:**\n")
				for _, p := range alt.Pros {
					sb.WriteString(fmt.Sprintf("- %s\n", p))
				}
				sb.WriteString("\n")
			}
			if len(alt.Cons) > 0 {
				sb.WriteString("**Cons:**\n")
				for _, c := range alt.Cons {
					sb.WriteString(fmt.Sprintf("- %s\n", c))
				}
				sb.WriteString("\n")
			}
			if alt.RejectedFor != "" {
				sb.WriteString(fmt.Sprintf("**Rejected because:** %s\n\n", alt.RejectedFor))
			}
		}
	}

	if len(d.Related) > 0 {
		sb.WriteString("## Related Decisions\n\n")
		for _, r := range d.Related {
			sb.WriteString(fmt.Sprintf("- %s\n", r))
		}
		sb.WriteString("\n")
	}

	if d.Supersedes != "" {
		sb.WriteString(fmt.Sprintf("**Supersedes:** %s\n\n", d.Supersedes))
	}

	if d.SupersededBy != "" {
		sb.WriteString(fmt.Sprintf("**Superseded by:** %s\n\n", d.SupersededBy))
	}

	return sb.String()
}

// RenderIndex renders the decision log as a Markdown index.
func (l *DecisionLog) RenderIndex() string {
	var sb strings.Builder

	sb.WriteString("# Architecture Decision Records\n\n")
	sb.WriteString("This document tracks architectural decisions made for this project.\n\n")

	// Group by status
	byStatus := make(map[DecisionStatus][]*Decision)
	for _, d := range l.List() {
		byStatus[d.Status] = append(byStatus[d.Status], d)
	}

	// Show accepted first
	statusOrder := []DecisionStatus{StatusAccepted, StatusProposed, StatusDeprecated, StatusSuperseded, StatusRejected}

	for _, status := range statusOrder {
		decisions := byStatus[status]
		if len(decisions) == 0 {
			continue
		}

		sb.WriteString(fmt.Sprintf("## %s\n\n", formatStatus(status)))
		sb.WriteString("| ID | Title | Date | Tags |\n")
		sb.WriteString("|-----|-------|------|------|\n")

		for _, d := range decisions {
			tags := strings.Join(d.Tags, ", ")
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
				d.ID, d.Title, d.Date.Format("2006-01-02"), tags))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// ExportJSON exports the decision log as JSON.
func (l *DecisionLog) ExportJSON() ([]byte, error) {
	decisions := l.List()
	return json.MarshalIndent(decisions, "", "  ")
}

// Statistics returns decision log statistics.
type Statistics struct {
	Total       int            `json:"total"`
	ByStatus    map[string]int `json:"by_status"`
	ByTag       map[string]int `json:"by_tag"`
	ByProject   map[string]int `json:"by_project"`
	RecentCount int            `json:"recent_count"` // Last 30 days
}

// GetStatistics returns statistics about the decision log.
func (l *DecisionLog) GetStatistics() *Statistics {
	stats := &Statistics{
		Total:     len(l.decisions),
		ByStatus:  make(map[string]int),
		ByTag:     make(map[string]int),
		ByProject: make(map[string]int),
	}

	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	for _, d := range l.decisions {
		stats.ByStatus[string(d.Status)]++

		for _, tag := range d.Tags {
			stats.ByTag[tag]++
		}

		if d.Project != "" {
			stats.ByProject[d.Project]++
		}

		if d.Date.After(thirtyDaysAgo) {
			stats.RecentCount++
		}
	}

	return stats
}

// Helper functions

func formatStatus(status DecisionStatus) string {
	switch status {
	case StatusAccepted:
		return "✅ Accepted"
	case StatusProposed:
		return "📋 Proposed"
	case StatusDeprecated:
		return "⚠️ Deprecated"
	case StatusSuperseded:
		return "🔄 Superseded"
	case StatusRejected:
		return "❌ Rejected"
	default:
		return string(status)
	}
}

// SaveDecisionFiles saves each decision as an individual Markdown file.
func (l *DecisionLog) SaveDecisionFiles(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	for _, d := range l.List() {
		filename := fmt.Sprintf("%s.md", strings.ToLower(d.ID))
		path := filepath.Join(dir, filename)

		content := d.RenderMarkdown()
		if err := os.WriteFile(path, []byte(content), 0600); err != nil {
			return fmt.Errorf("writing %s: %w", filename, err)
		}
	}

	// Write index
	indexPath := filepath.Join(dir, "README.md")
	indexContent := l.RenderIndex()
	if err := os.WriteFile(indexPath, []byte(indexContent), 0600); err != nil {
		return fmt.Errorf("writing index: %w", err)
	}

	return nil
}

// GetAllTags returns all unique tags used in decisions.
func (l *DecisionLog) GetAllTags() []string {
	tagSet := make(map[string]bool)
	for _, d := range l.decisions {
		for _, tag := range d.Tags {
			tagSet[tag] = true
		}
	}

	tags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
}
