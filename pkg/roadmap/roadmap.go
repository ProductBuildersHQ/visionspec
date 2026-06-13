// Package roadmap provides ROADMAP.md generation from project specifications.
//
// This package analyzes project specs, dependencies, and priorities to generate
// a cohesive roadmap document that tracks milestones, phases, and deliverables.
package roadmap

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Phase represents a roadmap phase with milestones.
type Phase struct {
	ID          string      `json:"id" yaml:"id"`
	Name        string      `json:"name" yaml:"name"`
	Description string      `json:"description,omitempty" yaml:"description,omitempty"`
	StartDate   *time.Time  `json:"start_date,omitempty" yaml:"start_date,omitempty"`
	EndDate     *time.Time  `json:"end_date,omitempty" yaml:"end_date,omitempty"`
	Status      PhaseStatus `json:"status" yaml:"status"`
	Milestones  []Milestone `json:"milestones,omitempty" yaml:"milestones,omitempty"`
	Projects    []string    `json:"projects,omitempty" yaml:"projects,omitempty"`
	Order       int         `json:"order" yaml:"order"`
}

// PhaseStatus indicates the current state of a phase.
type PhaseStatus string

const (
	PhaseStatusPlanned    PhaseStatus = "planned"
	PhaseStatusInProgress PhaseStatus = "in_progress"
	PhaseStatusCompleted  PhaseStatus = "completed"
	PhaseStatusBlocked    PhaseStatus = "blocked"
)

// Milestone represents a specific deliverable within a phase.
type Milestone struct {
	ID           string          `json:"id" yaml:"id"`
	Name         string          `json:"name" yaml:"name"`
	Description  string          `json:"description,omitempty" yaml:"description,omitempty"`
	DueDate      *time.Time      `json:"due_date,omitempty" yaml:"due_date,omitempty"`
	Status       MilestoneStatus `json:"status" yaml:"status"`
	Priority     Priority        `json:"priority" yaml:"priority"`
	Deliverables []string        `json:"deliverables,omitempty" yaml:"deliverables,omitempty"`
	Dependencies []string        `json:"dependencies,omitempty" yaml:"dependencies,omitempty"` // Other milestone IDs
	Projects     []string        `json:"projects,omitempty" yaml:"projects,omitempty"`
	RMI          string          `json:"rmi,omitempty" yaml:"rmi,omitempty"` // Roadmap Item ID
	Tags         []string        `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// MilestoneStatus indicates milestone completion state.
type MilestoneStatus string

const (
	MilestoneStatusPending    MilestoneStatus = "pending"
	MilestoneStatusInProgress MilestoneStatus = "in_progress"
	MilestoneStatusCompleted  MilestoneStatus = "completed"
	MilestoneStatusBlocked    MilestoneStatus = "blocked"
	MilestoneStatusDeferred   MilestoneStatus = "deferred"
)

// Priority indicates milestone importance.
type Priority string

const (
	PriorityCritical Priority = "critical"
	PriorityHigh     Priority = "high"
	PriorityMedium   Priority = "medium"
	PriorityLow      Priority = "low"
)

// Roadmap represents the complete project roadmap.
type Roadmap struct {
	Title       string            `json:"title" yaml:"title"`
	Description string            `json:"description,omitempty" yaml:"description,omitempty"`
	Version     string            `json:"version,omitempty" yaml:"version,omitempty"`
	UpdatedAt   time.Time         `json:"updated_at" yaml:"updated_at"`
	Phases      []Phase           `json:"phases" yaml:"phases"`
	Metadata    map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// Generator creates roadmaps from project data.
type Generator struct {
	specsRoot string
}

// NewGenerator creates a new roadmap generator.
func NewGenerator(specsRoot string) *Generator {
	return &Generator{specsRoot: specsRoot}
}

// ProjectInfo contains project metadata for roadmap generation.
type ProjectInfo struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Status      string   `yaml:"status"`
	Priority    string   `yaml:"priority"`
	Phase       string   `yaml:"phase"`
	Tags        []string `yaml:"tags"`
	DueDate     string   `yaml:"due_date"`
}

// Generate creates a roadmap from project specs.
func (g *Generator) Generate() (*Roadmap, error) {
	roadmap := &Roadmap{
		Title:     "Project Roadmap",
		UpdatedAt: time.Now(),
		Phases:    []Phase{},
	}

	// Read existing ROADMAP.yaml if it exists
	roadmapPath := filepath.Join(g.specsRoot, "ROADMAP.yaml")
	if data, err := os.ReadFile(roadmapPath); err == nil {
		if err := yaml.Unmarshal(data, roadmap); err != nil {
			return nil, fmt.Errorf("parsing ROADMAP.yaml: %w", err)
		}
	}

	// Scan projects and merge into roadmap
	projects, err := g.scanProjects()
	if err != nil {
		return nil, err
	}

	// Group projects by phase
	phaseProjects := make(map[string][]ProjectInfo)
	for _, proj := range projects {
		phase := proj.Phase
		if phase == "" {
			phase = "backlog"
		}
		phaseProjects[phase] = append(phaseProjects[phase], proj)
	}

	// Update or create phases
	existingPhases := make(map[string]*Phase)
	for i := range roadmap.Phases {
		existingPhases[roadmap.Phases[i].ID] = &roadmap.Phases[i]
	}

	for phaseName, projs := range phaseProjects {
		if phase, ok := existingPhases[phaseName]; ok {
			// Update existing phase
			for _, proj := range projs {
				if !contains(phase.Projects, proj.Name) {
					phase.Projects = append(phase.Projects, proj.Name)
				}
			}
		} else {
			// Create new phase
			newPhase := Phase{
				ID:       phaseName,
				Name:     formatPhaseName(phaseName),
				Status:   PhaseStatusPlanned,
				Projects: []string{},
				Order:    len(roadmap.Phases),
			}
			for _, proj := range projs {
				newPhase.Projects = append(newPhase.Projects, proj.Name)
			}
			roadmap.Phases = append(roadmap.Phases, newPhase)
		}
	}

	// Sort phases by order
	sort.Slice(roadmap.Phases, func(i, j int) bool {
		return roadmap.Phases[i].Order < roadmap.Phases[j].Order
	})

	return roadmap, nil
}

// scanProjects reads all project configs.
func (g *Generator) scanProjects() ([]ProjectInfo, error) {
	var projects []ProjectInfo

	entries, err := os.ReadDir(g.specsRoot)
	if err != nil {
		return nil, fmt.Errorf("reading specs root: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		configPath := filepath.Join(g.specsRoot, entry.Name(), "visionspec.yaml")
		data, err := os.ReadFile(configPath)
		if err != nil {
			continue // Skip projects without config
		}

		var config struct {
			Project ProjectInfo `yaml:"project"`
		}
		if err := yaml.Unmarshal(data, &config); err != nil {
			continue
		}

		if config.Project.Name == "" {
			config.Project.Name = entry.Name()
		}

		projects = append(projects, config.Project)
	}

	return projects, nil
}

// RenderMarkdown renders the roadmap as Markdown.
func (r *Roadmap) RenderMarkdown(w io.Writer) error {
	fmt.Fprintf(w, "# %s\n\n", r.Title)

	if r.Description != "" {
		fmt.Fprintf(w, "%s\n\n", r.Description)
	}

	fmt.Fprintf(w, "_Last updated: %s_\n\n", r.UpdatedAt.Format("2006-01-02"))

	if r.Version != "" {
		fmt.Fprintf(w, "**Version:** %s\n\n", r.Version)
	}

	fmt.Fprintf(w, "---\n\n")

	// Table of contents
	fmt.Fprintf(w, "## Contents\n\n")
	for i, phase := range r.Phases {
		fmt.Fprintf(w, "%d. [%s](#%s)\n", i+1, phase.Name, slugify(phase.Name))
	}
	fmt.Fprintf(w, "\n---\n\n")

	// Render each phase
	for _, phase := range r.Phases {
		fmt.Fprintf(w, "## %s\n\n", phase.Name)

		// Status badge
		statusBadge := getStatusBadge(string(phase.Status))
		fmt.Fprintf(w, "**Status:** %s\n\n", statusBadge)

		if phase.Description != "" {
			fmt.Fprintf(w, "%s\n\n", phase.Description)
		}

		// Dates
		if phase.StartDate != nil || phase.EndDate != nil {
			fmt.Fprintf(w, "**Timeline:** ")
			if phase.StartDate != nil {
				fmt.Fprintf(w, "%s", phase.StartDate.Format("Jan 2006"))
			}
			if phase.EndDate != nil {
				if phase.StartDate != nil {
					fmt.Fprintf(w, " - ")
				}
				fmt.Fprintf(w, "%s", phase.EndDate.Format("Jan 2006"))
			}
			fmt.Fprintf(w, "\n\n")
		}

		// Projects
		if len(phase.Projects) > 0 {
			fmt.Fprintf(w, "### Projects\n\n")
			for _, proj := range phase.Projects {
				fmt.Fprintf(w, "- [ ] %s\n", proj)
			}
			fmt.Fprintf(w, "\n")
		}

		// Milestones
		if len(phase.Milestones) > 0 {
			fmt.Fprintf(w, "### Milestones\n\n")
			for _, ms := range phase.Milestones {
				checkbox := "[ ]"
				if ms.Status == MilestoneStatusCompleted {
					checkbox = "[x]"
				}

				priorityBadge := ""
				switch ms.Priority {
				case PriorityCritical:
					priorityBadge = " 🔴"
				case PriorityHigh:
					priorityBadge = " 🟠"
				case PriorityMedium:
					priorityBadge = " 🟡"
				}

				rmiRef := ""
				if ms.RMI != "" {
					rmiRef = fmt.Sprintf(" `%s`", ms.RMI)
				}

				fmt.Fprintf(w, "- %s **%s**%s%s\n", checkbox, ms.Name, priorityBadge, rmiRef)

				if ms.Description != "" {
					fmt.Fprintf(w, "  - %s\n", ms.Description)
				}

				if ms.DueDate != nil {
					fmt.Fprintf(w, "  - Due: %s\n", ms.DueDate.Format("2006-01-02"))
				}

				if len(ms.Deliverables) > 0 {
					fmt.Fprintf(w, "  - Deliverables:\n")
					for _, d := range ms.Deliverables {
						fmt.Fprintf(w, "    - %s\n", d)
					}
				}
			}
			fmt.Fprintf(w, "\n")
		}

		fmt.Fprintf(w, "---\n\n")
	}

	// Summary statistics
	fmt.Fprintf(w, "## Summary\n\n")
	r.renderSummary(w)

	return nil
}

// renderSummary outputs roadmap statistics.
func (r *Roadmap) renderSummary(w io.Writer) {
	totalMilestones := 0
	completedMilestones := 0
	totalProjects := 0

	projectSet := make(map[string]bool)
	for _, phase := range r.Phases {
		totalMilestones += len(phase.Milestones)
		for _, ms := range phase.Milestones {
			if ms.Status == MilestoneStatusCompleted {
				completedMilestones++
			}
		}
		for _, p := range phase.Projects {
			projectSet[p] = true
		}
	}
	totalProjects = len(projectSet)

	fmt.Fprintf(w, "| Metric | Value |\n")
	fmt.Fprintf(w, "|--------|-------|\n")
	fmt.Fprintf(w, "| Total Phases | %d |\n", len(r.Phases))
	fmt.Fprintf(w, "| Total Projects | %d |\n", totalProjects)
	fmt.Fprintf(w, "| Total Milestones | %d |\n", totalMilestones)
	if totalMilestones > 0 {
		progress := float64(completedMilestones) / float64(totalMilestones) * 100
		fmt.Fprintf(w, "| Progress | %.0f%% (%d/%d) |\n", progress, completedMilestones, totalMilestones)
	}
	fmt.Fprintf(w, "\n")
}

// Save writes the roadmap to YAML and Markdown files.
func (r *Roadmap) Save(specsRoot string) error {
	// Save YAML
	yamlPath := filepath.Join(specsRoot, "ROADMAP.yaml")
	yamlData, err := yaml.Marshal(r)
	if err != nil {
		return fmt.Errorf("marshaling YAML: %w", err)
	}
	if err := os.WriteFile(yamlPath, yamlData, 0600); err != nil {
		return fmt.Errorf("writing ROADMAP.yaml: %w", err)
	}

	// Save Markdown
	mdPath := filepath.Join(specsRoot, "ROADMAP.md")
	f, err := os.Create(mdPath)
	if err != nil {
		return fmt.Errorf("creating ROADMAP.md: %w", err)
	}
	defer f.Close()

	if err := r.RenderMarkdown(f); err != nil {
		return fmt.Errorf("rendering markdown: %w", err)
	}

	return nil
}

// Load reads a roadmap from YAML.
func Load(path string) (*Roadmap, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading roadmap: %w", err)
	}

	var roadmap Roadmap
	if err := yaml.Unmarshal(data, &roadmap); err != nil {
		return nil, fmt.Errorf("parsing roadmap: %w", err)
	}

	return &roadmap, nil
}

// Helper functions

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func formatPhaseName(id string) string {
	// Convert kebab-case to Title Case
	words := strings.Split(id, "-")
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}

func getStatusBadge(status string) string {
	switch status {
	case "completed":
		return "✅ Completed"
	case "in_progress":
		return "🚧 In Progress"
	case "blocked":
		return "🚫 Blocked"
	case "planned":
		return "📋 Planned"
	default:
		return status
	}
}

// AddMilestone adds a milestone to a phase.
func (r *Roadmap) AddMilestone(phaseID string, milestone Milestone) error {
	for i := range r.Phases {
		if r.Phases[i].ID == phaseID {
			r.Phases[i].Milestones = append(r.Phases[i].Milestones, milestone)
			return nil
		}
	}
	return fmt.Errorf("phase %q not found", phaseID)
}

// UpdateMilestoneStatus updates the status of a milestone.
func (r *Roadmap) UpdateMilestoneStatus(milestoneID string, status MilestoneStatus) error {
	for i := range r.Phases {
		for j := range r.Phases[i].Milestones {
			if r.Phases[i].Milestones[j].ID == milestoneID {
				r.Phases[i].Milestones[j].Status = status
				return nil
			}
		}
	}
	return fmt.Errorf("milestone %q not found", milestoneID)
}

// GetProgress returns overall roadmap progress (0-100).
func (r *Roadmap) GetProgress() float64 {
	total := 0
	completed := 0

	for _, phase := range r.Phases {
		for _, ms := range phase.Milestones {
			total++
			if ms.Status == MilestoneStatusCompleted {
				completed++
			}
		}
	}

	if total == 0 {
		return 0
	}
	return float64(completed) / float64(total) * 100
}
