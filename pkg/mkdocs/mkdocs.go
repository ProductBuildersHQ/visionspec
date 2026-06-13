// Package mkdocs generates MkDocs-compatible markdown files for visionspec projects.
package mkdocs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ProductBuildersHQ/visionspec/pkg/config"
	"github.com/ProductBuildersHQ/visionspec/pkg/status"
	"github.com/ProductBuildersHQ/visionspec/pkg/types"
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
	fmt.Fprintf(w, "*Generated at %s by VisionSpec*\n", report.GeneratedAt.Format(time.RFC3339))

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
	fmt.Fprintf(w, "*Generated at %s by VisionSpec*\n", time.Now().Format(time.RFC3339))

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

		// Check if it's a valid visionspec project
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

// NavItem represents a single item in the MkDocs navigation.
type NavItem struct {
	Title    string    `yaml:"title,omitempty"`
	Path     string    `yaml:"-"`
	Children []NavItem `yaml:"-"`
}

// MarshalYAML implements custom YAML marshaling for NavItem.
func (n NavItem) MarshalYAML() (interface{}, error) {
	if len(n.Children) > 0 {
		// Section with children
		children := make([]interface{}, len(n.Children))
		for i, child := range n.Children {
			m, err := child.MarshalYAML()
			if err != nil {
				return nil, err
			}
			children[i] = m
		}
		return map[string]interface{}{n.Title: children}, nil
	}
	// Leaf item
	return map[string]string{n.Title: n.Path}, nil
}

// GenerateNavigation generates the nav section for mkdocs.yml from a specs directory.
func GenerateNavigation(specsDir string) ([]NavItem, error) {
	var nav []NavItem

	// Add specs landing
	nav = append(nav, NavItem{Title: "Overview", Path: "specs/index.md"})

	// Scan projects
	projects, err := ScanProjects(specsDir)
	if err != nil {
		return nil, err
	}

	for _, p := range projects {
		projectNav := NavItem{Title: p.Name}

		// Add project index
		projectNav.Children = append(projectNav.Children, NavItem{
			Title: "Overview",
			Path:  fmt.Sprintf("specs/%s/index.md", p.Name),
		})

		// Scan for spec files
		projectPath := filepath.Join(specsDir, p.Name)
		for _, dir := range []string{"source", "gtm", "technical"} {
			dirPath := filepath.Join(projectPath, dir)
			if entries, err := os.ReadDir(dirPath); err == nil {
				for _, entry := range entries {
					if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
						continue
					}
					// Skip drafts
					if strings.Contains(entry.Name(), ".draft.") {
						continue
					}
					name := strings.TrimSuffix(entry.Name(), ".md")
					projectNav.Children = append(projectNav.Children, NavItem{
						Title: strings.ToUpper(name),
						Path:  fmt.Sprintf("specs/%s/%s/%s", p.Name, dir, entry.Name()),
					})
				}
			}
		}

		// Add eval directory if it has files
		evalDir := filepath.Join(projectPath, "eval")
		if entries, err := os.ReadDir(evalDir); err == nil && len(entries) > 0 {
			evalNav := NavItem{Title: "Evaluations"}
			for _, entry := range entries {
				if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
					continue
				}
				name := strings.TrimSuffix(entry.Name(), ".eval.md")
				evalNav.Children = append(evalNav.Children, NavItem{
					Title: strings.ToUpper(name) + " Eval",
					Path:  fmt.Sprintf("specs/%s/eval/%s", p.Name, entry.Name()),
				})
			}
			if len(evalNav.Children) > 0 {
				projectNav.Children = append(projectNav.Children, evalNav)
			}
		}

		nav = append(nav, projectNav)
	}

	// Add global docs if present
	constitutionPath := filepath.Join(specsDir, config.ConstitutionFile)
	if _, err := os.Stat(constitutionPath); err == nil {
		nav = append(nav, NavItem{
			Title: "Constitution",
			Path:  fmt.Sprintf("specs/%s", config.ConstitutionFile),
		})
	}

	roadmapPath := filepath.Join(specsDir, config.RoadmapFile)
	if _, err := os.Stat(roadmapPath); err == nil {
		nav = append(nav, NavItem{
			Title: "Roadmap",
			Path:  fmt.Sprintf("specs/%s", config.RoadmapFile),
		})
	}

	return nav, nil
}

// WriteNavigation writes the nav section to a YAML file fragment.
func WriteNavigation(specsDir string, outputPath string) error {
	nav, err := GenerateNavigation(specsDir)
	if err != nil {
		return err
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("creating nav file: %w", err)
	}
	defer f.Close()

	fmt.Fprintf(f, "# Auto-generated by visionspec docs\n")
	fmt.Fprintf(f, "# Include this in your mkdocs.yml nav section\n\n")
	fmt.Fprintf(f, "nav:\n")
	fmt.Fprintf(f, "  - Specs:\n")

	for _, item := range nav {
		if len(item.Children) == 0 {
			fmt.Fprintf(f, "    - %s: %s\n", item.Title, item.Path)
		} else {
			fmt.Fprintf(f, "    - %s:\n", item.Title)
			for _, child := range item.Children {
				if len(child.Children) == 0 {
					fmt.Fprintf(f, "      - %s: %s\n", child.Title, child.Path)
				} else {
					fmt.Fprintf(f, "      - %s:\n", child.Title)
					for _, subchild := range child.Children {
						fmt.Fprintf(f, "        - %s: %s\n", subchild.Title, subchild.Path)
					}
				}
			}
		}
	}

	return nil
}

// GraphEmbedOptions configures graph embedding in documentation.
type GraphEmbedOptions struct {
	// Format is the graph output format: "mermaid", "iframe", "svg"
	Format string
	// Width is the iframe/svg width (e.g., "100%", "800px")
	Width string
	// Height is the iframe/svg height (e.g., "600px")
	Height string
	// GraphURL is the URL to the graph HTML file for iframe embedding
	GraphURL string
	// GraphPath is the local path to the graph file
	GraphPath string
	// Title is the optional title for the graph section
	Title string
	// Description is the optional description
	Description string
}

// DefaultGraphEmbedOptions returns default graph embed options.
func DefaultGraphEmbedOptions() GraphEmbedOptions {
	return GraphEmbedOptions{
		Format: "mermaid",
		Width:  "100%",
		Height: "600px",
		Title:  "Specification Graph",
	}
}

// GenerateGraphEmbed generates markdown for embedding a graph visualization.
func GenerateGraphEmbed(w io.Writer, nodes []GraphNode, edges []GraphEdge, opts GraphEmbedOptions) error {
	if opts.Title != "" {
		fmt.Fprintf(w, "## %s\n\n", opts.Title)
	}
	if opts.Description != "" {
		fmt.Fprintf(w, "%s\n\n", opts.Description)
	}

	switch opts.Format {
	case "iframe":
		return generateIframeEmbed(w, opts)
	case "svg":
		return generateSVGEmbed(w, opts)
	default:
		return generateMermaidEmbed(w, nodes, edges)
	}
}

// GraphNode represents a node in the spec graph.
type GraphNode struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Type     string `json:"type"` // requirement, feature, task, acceptance
	Status   string `json:"status,omitempty"`
	Category string `json:"category,omitempty"`
}

// GraphEdge represents an edge in the spec graph.
type GraphEdge struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Relation string `json:"relation"` // implements, depends_on, traces_to, validates
}

// generateMermaidEmbed generates a Mermaid diagram.
func generateMermaidEmbed(w io.Writer, nodes []GraphNode, edges []GraphEdge) error {
	fmt.Fprintf(w, "```mermaid\ngraph TD\n")

	// Group nodes by type for styling
	nodesByType := make(map[string][]GraphNode)
	for _, node := range nodes {
		nodesByType[node.Type] = append(nodesByType[node.Type], node)
	}

	// Write nodes
	for _, node := range nodes {
		shape := getNodeShape(node.Type)
		label := escapeLabel(node.Label)
		fmt.Fprintf(w, "    %s%s%s%s\n", node.ID, shape[0], label, shape[1])
	}

	fmt.Fprintf(w, "\n")

	// Write edges
	for _, edge := range edges {
		arrow := getArrowStyle(edge.Relation)
		fmt.Fprintf(w, "    %s %s %s\n", edge.From, arrow, edge.To)
	}

	// Add styles
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "    classDef requirement fill:#e1f5fe,stroke:#01579b\n")
	fmt.Fprintf(w, "    classDef feature fill:#e8f5e9,stroke:#2e7d32\n")
	fmt.Fprintf(w, "    classDef task fill:#fff3e0,stroke:#ef6c00\n")
	fmt.Fprintf(w, "    classDef acceptance fill:#f3e5f5,stroke:#7b1fa2\n")

	// Apply classes
	for nodeType, typeNodes := range nodesByType {
		if len(typeNodes) > 0 {
			ids := make([]string, len(typeNodes))
			for i, n := range typeNodes {
				ids[i] = n.ID
			}
			fmt.Fprintf(w, "    class %s %s\n", strings.Join(ids, ","), nodeType)
		}
	}

	fmt.Fprintf(w, "```\n\n")
	return nil
}

// getNodeShape returns Mermaid shape delimiters based on node type.
func getNodeShape(nodeType string) [2]string {
	switch nodeType {
	case "requirement":
		return [2]string{"[[", "]]"} // Stadium-shaped
	case "feature":
		return [2]string{"[", "]"} // Rectangle
	case "task":
		return [2]string{"((", "))"} // Circle
	case "acceptance":
		return [2]string{"{", "}"} // Diamond
	default:
		return [2]string{"[", "]"}
	}
}

// getArrowStyle returns Mermaid arrow style based on relation type.
func getArrowStyle(relation string) string {
	switch relation {
	case "implements":
		return "==>"
	case "depends_on":
		return "-->"
	case "traces_to":
		return "-.->"
	case "validates":
		return "--o"
	default:
		return "-->"
	}
}

// escapeLabel escapes special characters for Mermaid labels.
func escapeLabel(label string) string {
	// Escape quotes and special characters
	label = strings.ReplaceAll(label, "\"", "'")
	label = strings.ReplaceAll(label, "[", "(")
	label = strings.ReplaceAll(label, "]", ")")
	label = strings.ReplaceAll(label, "{", "(")
	label = strings.ReplaceAll(label, "}", ")")
	// Truncate long labels
	if len(label) > 40 {
		label = label[:37] + "..."
	}
	return "\"" + label + "\""
}

// generateIframeEmbed generates an iframe embed for external graph.
func generateIframeEmbed(w io.Writer, opts GraphEmbedOptions) error {
	if opts.GraphURL == "" && opts.GraphPath != "" {
		opts.GraphURL = opts.GraphPath
	}
	if opts.GraphURL == "" {
		return fmt.Errorf("GraphURL or GraphPath required for iframe embed")
	}

	fmt.Fprintf(w, "<iframe src=\"%s\" width=\"%s\" height=\"%s\" frameborder=\"0\"></iframe>\n\n",
		opts.GraphURL, opts.Width, opts.Height)
	return nil
}

// generateSVGEmbed generates an embedded SVG or link to SVG.
func generateSVGEmbed(w io.Writer, opts GraphEmbedOptions) error {
	if opts.GraphPath == "" {
		return fmt.Errorf("GraphPath required for SVG embed")
	}

	// Check if file exists and read content
	content, err := os.ReadFile(opts.GraphPath)
	if err != nil {
		// Fall back to image link
		fmt.Fprintf(w, "![%s](%s)\n\n", opts.Title, opts.GraphPath)
		return nil
	}

	// Embed SVG directly
	fmt.Fprintf(w, "%s\n\n", string(content))
	return nil
}

// ExtractGraphFromSpec extracts graph nodes and edges from spec content.
// This is a simplified extraction - a full implementation would use
// more sophisticated NLP/parsing.
func ExtractGraphFromSpec(specContent string, projectName string) ([]GraphNode, []GraphEdge) {
	var nodes []GraphNode
	var edges []GraphEdge

	lines := strings.Split(specContent, "\n")
	currentSection := ""
	nodeNum := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Track section
		if strings.HasPrefix(trimmed, "## ") {
			currentSection = strings.TrimPrefix(trimmed, "## ")
			continue
		}

		// Extract requirements from bullet points
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			content := strings.TrimPrefix(strings.TrimPrefix(trimmed, "- "), "* ")

			nodeNum++
			nodeType := inferNodeType(currentSection, content)
			node := GraphNode{
				ID:       fmt.Sprintf("%s_%03d", strings.ToLower(nodeType[:3]), nodeNum),
				Label:    content,
				Type:     nodeType,
				Category: currentSection,
			}
			nodes = append(nodes, node)
		}
	}

	// Generate some basic edges (in a real implementation, this would
	// use dependency analysis)
	for i := 1; i < len(nodes); i++ {
		if nodes[i-1].Type == "requirement" && nodes[i].Type == "feature" {
			edges = append(edges, GraphEdge{
				From:     nodes[i-1].ID,
				To:       nodes[i].ID,
				Relation: "implements",
			})
		} else if nodes[i].Type == "task" && i > 0 {
			edges = append(edges, GraphEdge{
				From:     nodes[i-1].ID,
				To:       nodes[i].ID,
				Relation: "depends_on",
			})
		}
	}

	return nodes, edges
}

// inferNodeType infers the node type from section and content.
func inferNodeType(section, content string) string {
	sectionLower := strings.ToLower(section)
	contentLower := strings.ToLower(content)

	switch {
	case strings.Contains(sectionLower, "requirement") || strings.Contains(contentLower, "must") || strings.Contains(contentLower, "shall"):
		return "requirement"
	case strings.Contains(sectionLower, "feature") || strings.Contains(sectionLower, "capability"):
		return "feature"
	case strings.Contains(sectionLower, "task") || strings.Contains(sectionLower, "work") || strings.HasPrefix(content, "[ ]"):
		return "task"
	case strings.Contains(sectionLower, "acceptance") || strings.Contains(sectionLower, "criteria") || strings.Contains(contentLower, "given"):
		return "acceptance"
	default:
		return "feature"
	}
}

// WriteGraphPage generates a dedicated graph visualization page.
func WriteGraphPage(projectPath string, specContent string, opts GraphEmbedOptions) error {
	projectName := filepath.Base(projectPath)
	nodes, edges := ExtractGraphFromSpec(specContent, projectName)

	graphPath := filepath.Join(projectPath, "graph.md")
	f, err := os.Create(graphPath)
	if err != nil {
		return fmt.Errorf("creating graph.md: %w", err)
	}
	defer f.Close()

	fmt.Fprintf(f, "# %s Specification Graph\n\n", projectName)
	fmt.Fprintf(f, "This graph shows the relationships between requirements, features, tasks, and acceptance criteria.\n\n")

	// Legend
	fmt.Fprintf(f, "## Legend\n\n")
	fmt.Fprintf(f, "- **Stadium shapes**: Requirements\n")
	fmt.Fprintf(f, "- **Rectangles**: Features\n")
	fmt.Fprintf(f, "- **Circles**: Tasks\n")
	fmt.Fprintf(f, "- **Diamonds**: Acceptance Criteria\n\n")

	// Generate graph
	opts.Title = "Traceability Graph"
	opts.Description = ""
	if err := GenerateGraphEmbed(f, nodes, edges, opts); err != nil {
		return err
	}

	// Summary
	fmt.Fprintf(f, "## Summary\n\n")
	fmt.Fprintf(f, "- **Total Nodes:** %d\n", len(nodes))
	fmt.Fprintf(f, "- **Total Edges:** %d\n", len(edges))

	// Count by type
	byType := make(map[string]int)
	for _, node := range nodes {
		byType[node.Type]++
	}
	for nodeType, count := range byType {
		// Capitalize first letter
		title := nodeType
		if len(title) > 0 {
			title = strings.ToUpper(string(title[0])) + title[1:]
		}
		fmt.Fprintf(f, "- **%s:** %d\n", title, count)
	}

	fmt.Fprintf(f, "\n---\n")
	fmt.Fprintf(f, "*Generated by VisionSpec*\n")

	return nil
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
