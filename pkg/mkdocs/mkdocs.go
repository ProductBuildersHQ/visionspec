// Package mkdocs generates MkDocs-compatible markdown files for multispec projects.
package mkdocs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/plexusone/multispec/pkg/config"
	"github.com/plexusone/multispec/pkg/status"
	"github.com/plexusone/multispec/pkg/types"
)

// ProjectIndexOptions configures project index generation.
type ProjectIndexOptions struct {
	// IncludeGraphLink adds a link to the graph visualization if available.
	IncludeGraphLink bool
	// GraphPath is the relative path to the graph HTML file.
	GraphPath string
}

// GenerateProjectIndex generates index.md for a single project.
func GenerateProjectIndex(w io.Writer, report *status.Report, opts ProjectIndexOptions) error {
	// Header with status badge
	statusEmoji := "🔴"
	statusText := "NOT READY"
	if report.Readiness.Ready {
		statusEmoji = "🟢"
		statusText = "READY"
	} else {
		// Check if in progress (some specs present)
		if report.Summary.PresentSpecs > 0 {
			statusEmoji = "🟡"
			statusText = "IN PROGRESS"
		}
	}

	fmt.Fprintf(w, "# Project: %s\n\n", report.Project)
	fmt.Fprintf(w, "**Status:** %s %s\n\n", statusEmoji, statusText)

	// Readiness gates
	fmt.Fprintf(w, "## Readiness Gates\n\n")
	for _, gate := range report.Readiness.Gates {
		icon := "❌"
		if gate.Passed {
			icon = "✅"
		}
		fmt.Fprintf(w, "- %s %s\n", icon, gate.Name)
	}
	fmt.Fprintf(w, "\n")

	// Group specs by category
	specsByCategory := make(map[types.SpecCategory][]status.SpecStatus)
	for _, spec := range report.Specs {
		specsByCategory[spec.Category] = append(specsByCategory[spec.Category], spec)
	}

	// Render specs table
	fmt.Fprintf(w, "## Specifications\n\n")
	fmt.Fprintf(w, "| Type | Category | Status | Eval | Approved |\n")
	fmt.Fprintf(w, "|------|----------|--------|------|----------|\n")

	for _, spec := range report.Specs {
		existsIcon := "❌"
		if spec.Exists {
			existsIcon = "✅"
		}

		evalText := "-"
		if spec.EvalStatus != nil && spec.EvalStatus.Exists {
			if spec.EvalStatus.Decision != "" {
				evalText = spec.EvalStatus.Decision
			} else {
				evalText = "✅"
			}
		}

		approvedIcon := "-"
		if spec.Approval != nil {
			approvedIcon = "✅"
		}

		typeText := string(spec.Type)
		if spec.Required {
			typeText = fmt.Sprintf("**%s***", spec.Type)
		}

		// Add link to spec file if it exists
		if spec.Exists {
			dir := spec.Type.Dir()
			if dir == "" {
				typeText = fmt.Sprintf("[%s](%s)", typeText, spec.Filename)
			} else {
				typeText = fmt.Sprintf("[%s](%s/%s)", typeText, dir, spec.Filename)
			}
		}

		fmt.Fprintf(w, "| %s | %s | %s | %s | %s |\n",
			typeText, spec.Category, existsIcon, evalText, approvedIcon)
	}
	fmt.Fprintf(w, "\n*\\* = required*\n\n")

	// Graph metrics if available
	if report.GraphMetrics != nil {
		fmt.Fprintf(w, "## Traceability Metrics\n\n")
		fmt.Fprintf(w, "- **Requirements:** %d extracted\n", report.GraphMetrics.TotalRequirements)
		fmt.Fprintf(w, "- **User Stories:** %d extracted\n", report.GraphMetrics.TotalUserStories)
		fmt.Fprintf(w, "- **Trace Coverage:** %.0f%%\n", report.GraphMetrics.TraceCoverage*100)
		fmt.Fprintf(w, "- **Conflicts:** %d detected\n", report.GraphMetrics.ConflictCount)

		if opts.IncludeGraphLink && opts.GraphPath != "" {
			fmt.Fprintf(w, "\n[View Graph →](%s)\n", opts.GraphPath)
		}
		fmt.Fprintf(w, "\n")
	}

	// Summary
	fmt.Fprintf(w, "## Summary\n\n")
	fmt.Fprintf(w, "- **Total specs:** %d\n", report.Summary.TotalSpecs)
	fmt.Fprintf(w, "- **Present:** %d\n", report.Summary.PresentSpecs)
	fmt.Fprintf(w, "- **Evaluated:** %d\n", report.Summary.EvaluatedSpecs)
	fmt.Fprintf(w, "- **Approved:** %d\n", report.Summary.ApprovedSpecs)

	fmt.Fprintf(w, "\n---\n")
	fmt.Fprintf(w, "*Generated at %s by MultiSpec*\n", report.GeneratedAt.Format(time.RFC3339))

	return nil
}

// ProjectSummary contains summary information for the specs landing page.
type ProjectSummary struct {
	Name          string
	Path          string
	Status        string // ready, in_progress, not_ready
	StatusEmoji   string
	Progress      int // Percentage
	LastUpdated   time.Time
	SpecCount     int
	EvalCount     int
	ApprovalCount int
}

// SpecsLandingOptions configures the specs landing page generation.
type SpecsLandingOptions struct {
	// IncludeConstitution adds a link to CONSTITUTION.md if it exists.
	IncludeConstitution bool
	// IncludeRoadmap adds a link to ROADMAP.md if it exists.
	IncludeRoadmap bool
	// ConstitutionPath is the relative path to CONSTITUTION.md.
	ConstitutionPath string
	// RoadmapPath is the relative path to ROADMAP.md.
	RoadmapPath string
}

// GenerateSpecsLanding generates the main docs/specs/index.md file.
func GenerateSpecsLanding(w io.Writer, projects []ProjectSummary, opts SpecsLandingOptions) error {
	fmt.Fprintf(w, "# Specifications\n\n")

	// Projects table
	if len(projects) > 0 {
		fmt.Fprintf(w, "## Projects\n\n")
		fmt.Fprintf(w, "| Project | Status | Progress | Last Updated |\n")
		fmt.Fprintf(w, "|---------|--------|----------|---------------|\n")

		for _, p := range projects {
			fmt.Fprintf(w, "| [%s](%s/) | %s %s | %d%% | %s |\n",
				p.Name, p.Name, p.StatusEmoji, p.Status,
				p.Progress, p.LastUpdated.Format("2006-01-02"))
		}
		fmt.Fprintf(w, "\n")
	}

	// Global resources
	if opts.IncludeConstitution || opts.IncludeRoadmap {
		fmt.Fprintf(w, "## Global Resources\n\n")
		if opts.IncludeConstitution && opts.ConstitutionPath != "" {
			fmt.Fprintf(w, "- [CONSTITUTION.md](%s) - Repository governance\n", opts.ConstitutionPath)
		}
		if opts.IncludeRoadmap && opts.RoadmapPath != "" {
			fmt.Fprintf(w, "- [ROADMAP.md](%s) - Development roadmap\n", opts.RoadmapPath)
		}
		fmt.Fprintf(w, "\n")
	}

	// Metrics summary
	totalSpecs := 0
	totalEvals := 0
	totalApprovals := 0
	readyCount := 0
	inProgressCount := 0

	for _, p := range projects {
		totalSpecs += p.SpecCount
		totalEvals += p.EvalCount
		totalApprovals += p.ApprovalCount
		if p.Status == "Ready" {
			readyCount++
		} else if p.Status == "In Progress" {
			inProgressCount++
		}
	}

	fmt.Fprintf(w, "## Metrics Summary\n\n")
	fmt.Fprintf(w, "- **Total Projects:** %d\n", len(projects))
	fmt.Fprintf(w, "- **Ready:** %d\n", readyCount)
	fmt.Fprintf(w, "- **In Progress:** %d\n", inProgressCount)
	fmt.Fprintf(w, "- **Total Specs:** %d\n", totalSpecs)
	fmt.Fprintf(w, "- **Evaluated:** %d\n", totalEvals)
	fmt.Fprintf(w, "- **Approved:** %d\n", totalApprovals)

	fmt.Fprintf(w, "\n---\n")
	fmt.Fprintf(w, "*Generated at %s by MultiSpec*\n", time.Now().Format(time.RFC3339))

	return nil
}

// ScanProjects scans a specs directory and returns summaries for all projects.
func ScanProjects(specsDir string) ([]ProjectSummary, error) {
	var summaries []ProjectSummary

	entries, err := os.ReadDir(specsDir)
	if err != nil {
		return nil, fmt.Errorf("reading specs directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		// Skip hidden directories and special files
		if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_") {
			continue
		}

		projectPath := filepath.Join(specsDir, name)

		// Check if it's a valid multispec project
		configPath := filepath.Join(projectPath, config.ConfigFileName)
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			continue
		}

		// Load project and generate report
		project, err := config.Load(projectPath)
		if err != nil {
			continue
		}

		report, err := status.Generate(project)
		if err != nil {
			continue
		}

		// Build summary
		summary := ProjectSummary{
			Name:          name,
			Path:          projectPath,
			LastUpdated:   report.GeneratedAt,
			SpecCount:     report.Summary.TotalSpecs,
			EvalCount:     report.Summary.EvaluatedSpecs,
			ApprovalCount: report.Summary.ApprovedSpecs,
		}

		// Calculate progress and status
		if report.Summary.TotalSpecs > 0 {
			summary.Progress = (report.Summary.ApprovedSpecs * 100) / report.Summary.TotalSpecs
		}

		if report.Readiness.Ready {
			summary.Status = "Ready"
			summary.StatusEmoji = "🟢"
		} else if report.Summary.PresentSpecs > 0 {
			summary.Status = "In Progress"
			summary.StatusEmoji = "🟡"
		} else {
			summary.Status = "Not Started"
			summary.StatusEmoji = "🔴"
		}

		summaries = append(summaries, summary)
	}

	// Sort by name
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].Name < summaries[j].Name
	})

	return summaries, nil
}

// WriteProjectIndex writes index.md for a project.
func WriteProjectIndex(projectPath string, report *status.Report, opts ProjectIndexOptions) error {
	indexPath := filepath.Join(projectPath, "index.md")

	f, err := os.Create(indexPath)
	if err != nil {
		return fmt.Errorf("creating index.md: %w", err)
	}
	defer f.Close()

	return GenerateProjectIndex(f, report, opts)
}

// WriteSpecsLanding writes the main specs/index.md file.
func WriteSpecsLanding(specsDir string, opts SpecsLandingOptions) error {
	// Check for global resources
	constitutionPath := filepath.Join(specsDir, config.ConstitutionFile)
	if _, err := os.Stat(constitutionPath); err == nil {
		opts.IncludeConstitution = true
		opts.ConstitutionPath = config.ConstitutionFile
	}

	roadmapPath := filepath.Join(specsDir, config.RoadmapFile)
	if _, err := os.Stat(roadmapPath); err == nil {
		opts.IncludeRoadmap = true
		opts.RoadmapPath = config.RoadmapFile
	}

	// Scan projects
	projects, err := ScanProjects(specsDir)
	if err != nil {
		return err
	}

	// Write index.md
	indexPath := filepath.Join(specsDir, "index.md")
	f, err := os.Create(indexPath)
	if err != nil {
		return fmt.Errorf("creating specs index.md: %w", err)
	}
	defer f.Close()

	return GenerateSpecsLanding(f, projects, opts)
}
