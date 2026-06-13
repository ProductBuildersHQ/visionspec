// Package deps provides cross-project dependency tracking for VisionSpec.
//
// This package enables projects to declare dependencies on other projects,
// track dependency versions, and detect conflicts or circular dependencies.
package deps

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

// Dependency represents a dependency on another project.
type Dependency struct {
	Project     string            `json:"project" yaml:"project"`
	Version     string            `json:"version,omitempty" yaml:"version,omitempty"`
	Type        DependencyType    `json:"type" yaml:"type"`
	Required    bool              `json:"required" yaml:"required"`
	Description string            `json:"description,omitempty" yaml:"description,omitempty"`
	Specs       []string          `json:"specs,omitempty" yaml:"specs,omitempty"` // Which specs are affected
	Metadata    map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// DependencyType categorizes the nature of the dependency.
type DependencyType string

const (
	// DepTypeBlocks indicates this project blocks the dependent project.
	DepTypeBlocks DependencyType = "blocks"
	// DepTypeRequires indicates this project requires the dependency.
	DepTypeRequires DependencyType = "requires"
	// DepTypeExtends indicates this project extends the dependency.
	DepTypeExtends DependencyType = "extends"
	// DepTypeRelated indicates projects are related but not strictly dependent.
	DepTypeRelated DependencyType = "related"
)

// DependencyGraph represents the full dependency graph across projects.
type DependencyGraph struct {
	Projects map[string]*ProjectNode `json:"projects"`
	Edges    []DependencyEdge        `json:"edges"`
}

// ProjectNode represents a project in the dependency graph.
type ProjectNode struct {
	Name         string       `json:"name"`
	Path         string       `json:"path"`
	Version      string       `json:"version,omitempty"`
	Status       string       `json:"status,omitempty"` // draft, approved, completed
	Dependencies []Dependency `json:"dependencies,omitempty"`
	Dependents   []string     `json:"dependents,omitempty"` // Projects that depend on this
}

// DependencyEdge represents a directed edge in the dependency graph.
type DependencyEdge struct {
	From     string         `json:"from"`
	To       string         `json:"to"`
	Type     DependencyType `json:"type"`
	Required bool           `json:"required"`
}

// Manager handles dependency operations for a project.
type Manager struct {
	projectPath string
	specsRoot   string
}

// NewManager creates a new dependency manager.
func NewManager(projectPath string) *Manager {
	// Find specs root (docs/specs or just the project path)
	specsRoot := projectPath
	if _, err := os.Stat(filepath.Join(projectPath, "docs", "specs")); err == nil {
		specsRoot = filepath.Join(projectPath, "docs", "specs")
	}

	return &Manager{
		projectPath: projectPath,
		specsRoot:   specsRoot,
	}
}

// GetDependencies returns the dependencies declared in visionspec.yaml.
func (m *Manager) GetDependencies() ([]Dependency, error) {
	configPath := filepath.Join(m.projectPath, "visionspec.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No dependencies declared
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var config struct {
		Dependencies []Dependency `yaml:"dependencies"`
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	return config.Dependencies, nil
}

// AddDependency adds a dependency to the project's visionspec.yaml.
func (m *Manager) AddDependency(dep Dependency) error {
	configPath := filepath.Join(m.projectPath, "visionspec.yaml")

	// Read existing config
	var config map[string]interface{}
	data, err := os.ReadFile(configPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("reading config: %w", err)
		}
		config = make(map[string]interface{})
	} else {
		if err := yaml.Unmarshal(data, &config); err != nil {
			return fmt.Errorf("parsing config: %w", err)
		}
	}

	// Get or create dependencies list
	deps, ok := config["dependencies"].([]interface{})
	if !ok {
		deps = []interface{}{}
	}

	// Add new dependency
	depMap := map[string]interface{}{
		"project":  dep.Project,
		"type":     string(dep.Type),
		"required": dep.Required,
	}
	if dep.Version != "" {
		depMap["version"] = dep.Version
	}
	if dep.Description != "" {
		depMap["description"] = dep.Description
	}
	if len(dep.Specs) > 0 {
		depMap["specs"] = dep.Specs
	}

	deps = append(deps, depMap)
	config["dependencies"] = deps

	// Write back
	output, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	return os.WriteFile(configPath, output, 0600)
}

// RemoveDependency removes a dependency from the project's visionspec.yaml.
func (m *Manager) RemoveDependency(projectName string) error {
	configPath := filepath.Join(m.projectPath, "visionspec.yaml")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("reading config: %w", err)
	}

	var config map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("parsing config: %w", err)
	}

	deps, ok := config["dependencies"].([]interface{})
	if !ok {
		return nil // No dependencies to remove
	}

	// Filter out the dependency
	filtered := []interface{}{}
	for _, d := range deps {
		depMap, ok := d.(map[string]interface{})
		if !ok {
			continue
		}
		if depMap["project"] != projectName {
			filtered = append(filtered, d)
		}
	}

	config["dependencies"] = filtered

	output, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	return os.WriteFile(configPath, output, 0600)
}

// BuildGraph builds the full dependency graph across all projects.
func (m *Manager) BuildGraph() (*DependencyGraph, error) {
	graph := &DependencyGraph{
		Projects: make(map[string]*ProjectNode),
		Edges:    []DependencyEdge{},
	}

	// Find all projects
	entries, err := os.ReadDir(m.specsRoot)
	if err != nil {
		return nil, fmt.Errorf("reading specs root: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		// Skip non-project directories
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		projectPath := filepath.Join(m.specsRoot, entry.Name())
		configPath := filepath.Join(projectPath, "visionspec.yaml")

		// Check if this is a valid project
		if _, err := os.Stat(configPath); err != nil {
			continue
		}

		node := &ProjectNode{
			Name: entry.Name(),
			Path: projectPath,
		}

		// Read project config
		data, err := os.ReadFile(configPath)
		if err != nil {
			continue
		}

		var config struct {
			Version      string       `yaml:"version"`
			Status       string       `yaml:"status"`
			Dependencies []Dependency `yaml:"dependencies"`
		}
		if err := yaml.Unmarshal(data, &config); err != nil {
			continue
		}

		node.Version = config.Version
		node.Status = config.Status
		node.Dependencies = config.Dependencies

		graph.Projects[entry.Name()] = node

		// Add edges for dependencies
		for _, dep := range config.Dependencies {
			graph.Edges = append(graph.Edges, DependencyEdge{
				From:     entry.Name(),
				To:       dep.Project,
				Type:     dep.Type,
				Required: dep.Required,
			})
		}
	}

	// Compute dependents (reverse edges)
	for _, edge := range graph.Edges {
		if toNode, ok := graph.Projects[edge.To]; ok {
			toNode.Dependents = append(toNode.Dependents, edge.From)
		}
	}

	return graph, nil
}

// DetectCycles detects circular dependencies in the graph.
func (g *DependencyGraph) DetectCycles() [][]string {
	var cycles [][]string
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	var dfs func(node string, path []string) bool
	dfs = func(node string, path []string) bool {
		visited[node] = true
		recStack[node] = true
		path = append(path, node)

		project, ok := g.Projects[node]
		if !ok {
			return false
		}

		for _, dep := range project.Dependencies {
			if !visited[dep.Project] {
				if dfs(dep.Project, path) {
					return true
				}
			} else if recStack[dep.Project] {
				// Found a cycle
				cycleStart := -1
				for i, p := range path {
					if p == dep.Project {
						cycleStart = i
						break
					}
				}
				if cycleStart >= 0 {
					cycle := append([]string{}, path[cycleStart:]...)
					cycle = append(cycle, dep.Project)
					cycles = append(cycles, cycle)
				}
				return true
			}
		}

		recStack[node] = false
		return false
	}

	for name := range g.Projects {
		if !visited[name] {
			dfs(name, []string{})
		}
	}

	return cycles
}

// TopologicalSort returns projects in dependency order.
func (g *DependencyGraph) TopologicalSort() ([]string, error) {
	cycles := g.DetectCycles()
	if len(cycles) > 0 {
		return nil, fmt.Errorf("circular dependencies detected: %v", cycles[0])
	}

	var sorted []string
	visited := make(map[string]bool)

	var visit func(name string)
	visit = func(name string) {
		if visited[name] {
			return
		}
		visited[name] = true

		if project, ok := g.Projects[name]; ok {
			for _, dep := range project.Dependencies {
				visit(dep.Project)
			}
		}

		sorted = append(sorted, name)
	}

	// Visit all projects
	names := make([]string, 0, len(g.Projects))
	for name := range g.Projects {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		visit(name)
	}

	return sorted, nil
}

// GetAffectedProjects returns all projects affected by changes to a given project.
func (g *DependencyGraph) GetAffectedProjects(projectName string) []string {
	affected := make(map[string]bool)
	var traverse func(name string)
	traverse = func(name string) {
		if project, ok := g.Projects[name]; ok {
			for _, dep := range project.Dependents {
				if !affected[dep] {
					affected[dep] = true
					traverse(dep)
				}
			}
		}
	}
	traverse(projectName)

	result := make([]string, 0, len(affected))
	for name := range affected {
		result = append(result, name)
	}
	sort.Strings(result)
	return result
}

// Validate checks for dependency issues.
type ValidationResult struct {
	Valid    bool              `json:"valid"`
	Errors   []ValidationError `json:"errors,omitempty"`
	Warnings []ValidationError `json:"warnings,omitempty"`
}

// ValidationError represents a dependency validation issue.
type ValidationError struct {
	Type    string `json:"type"`
	Project string `json:"project"`
	Message string `json:"message"`
}

// Validate checks the dependency graph for issues.
func (g *DependencyGraph) Validate() *ValidationResult {
	result := &ValidationResult{Valid: true}

	// Check for cycles
	cycles := g.DetectCycles()
	for _, cycle := range cycles {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Type:    "circular_dependency",
			Project: cycle[0],
			Message: fmt.Sprintf("Circular dependency: %s", strings.Join(cycle, " -> ")),
		})
	}

	// Check for missing dependencies
	for name, project := range g.Projects {
		for _, dep := range project.Dependencies {
			if _, ok := g.Projects[dep.Project]; !ok {
				if dep.Required {
					result.Valid = false
					result.Errors = append(result.Errors, ValidationError{
						Type:    "missing_dependency",
						Project: name,
						Message: fmt.Sprintf("Required dependency '%s' not found", dep.Project),
					})
				} else {
					result.Warnings = append(result.Warnings, ValidationError{
						Type:    "missing_dependency",
						Project: name,
						Message: fmt.Sprintf("Optional dependency '%s' not found", dep.Project),
					})
				}
			}
		}
	}

	return result
}

// ExportMermaid exports the dependency graph as Mermaid diagram.
func (g *DependencyGraph) ExportMermaid() string {
	var sb strings.Builder
	sb.WriteString("graph TD\n")

	// Sort for consistent output
	names := make([]string, 0, len(g.Projects))
	for name := range g.Projects {
		names = append(names, name)
	}
	sort.Strings(names)

	// Add nodes
	for _, name := range names {
		project := g.Projects[name]
		style := ""
		switch project.Status {
		case "completed":
			style = ":::completed"
		case "approved":
			style = ":::approved"
		case "draft":
			style = ":::draft"
		}
		sb.WriteString(fmt.Sprintf("    %s[\"%s\"]%s\n", sanitizeID(name), name, style))
	}

	sb.WriteString("\n")

	// Add edges
	for _, edge := range g.Edges {
		arrow := "-->"
		switch edge.Type {
		case DepTypeBlocks:
			arrow = "-.->|blocks|"
		case DepTypeRequires:
			arrow = "-->|requires|"
		case DepTypeExtends:
			arrow = "-->|extends|"
		case DepTypeRelated:
			arrow = "-.-|related|"
		}
		sb.WriteString(fmt.Sprintf("    %s %s %s\n",
			sanitizeID(edge.From), arrow, sanitizeID(edge.To)))
	}

	// Add styles
	sb.WriteString("\n")
	sb.WriteString("    classDef completed fill:#9f6,stroke:#333\n")
	sb.WriteString("    classDef approved fill:#69f,stroke:#333\n")
	sb.WriteString("    classDef draft fill:#ff9,stroke:#333\n")

	return sb.String()
}

// sanitizeID makes a string safe for use as a Mermaid node ID.
func sanitizeID(s string) string {
	return strings.ReplaceAll(s, "-", "_")
}

// DependencyReport contains a summary of project dependencies.
type DependencyReport struct {
	GeneratedAt      time.Time         `json:"generated_at"`
	TotalProjects    int               `json:"total_projects"`
	TotalEdges       int               `json:"total_edges"`
	Cycles           [][]string        `json:"cycles,omitempty"`
	CriticalPath     []string          `json:"critical_path,omitempty"`
	OrphanProjects   []string          `json:"orphan_projects,omitempty"` // No deps and no dependents
	LeafProjects     []string          `json:"leaf_projects,omitempty"`   // No dependents
	RootProjects     []string          `json:"root_projects,omitempty"`   // No dependencies
	ValidationResult *ValidationResult `json:"validation,omitempty"`
}

// GenerateReport creates a dependency report from the graph.
func (g *DependencyGraph) GenerateReport() *DependencyReport {
	report := &DependencyReport{
		GeneratedAt:   time.Now(),
		TotalProjects: len(g.Projects),
		TotalEdges:    len(g.Edges),
		Cycles:        g.DetectCycles(),
	}

	// Find orphans, leaves, and roots
	for name, project := range g.Projects {
		hasDeps := len(project.Dependencies) > 0
		hasDependents := len(project.Dependents) > 0

		if !hasDeps && !hasDependents {
			report.OrphanProjects = append(report.OrphanProjects, name)
		} else if !hasDependents {
			report.LeafProjects = append(report.LeafProjects, name)
		} else if !hasDeps {
			report.RootProjects = append(report.RootProjects, name)
		}
	}

	sort.Strings(report.OrphanProjects)
	sort.Strings(report.LeafProjects)
	sort.Strings(report.RootProjects)

	// Compute critical path (longest path)
	report.CriticalPath = g.findCriticalPath()

	// Validate
	report.ValidationResult = g.Validate()

	return report
}

// findCriticalPath finds the longest dependency path.
func (g *DependencyGraph) findCriticalPath() []string {
	var longest []string
	visited := make(map[string]bool)

	var dfs func(name string, path []string)
	dfs = func(name string, path []string) {
		if visited[name] {
			return
		}
		visited[name] = true
		path = append(path, name)

		if len(path) > len(longest) {
			longest = make([]string, len(path))
			copy(longest, path)
		}

		if project, ok := g.Projects[name]; ok {
			for _, dep := range project.Dependencies {
				dfs(dep.Project, path)
			}
		}

		visited[name] = false
	}

	for name := range g.Projects {
		dfs(name, []string{})
	}

	return longest
}

// ExportJSON exports the dependency graph as JSON.
func (g *DependencyGraph) ExportJSON() ([]byte, error) {
	return json.MarshalIndent(g, "", "  ")
}
