// Package specgraph provides requirement graph extraction from visionspec specs.
// It integrates with graphize to build knowledge graphs from spec documents.
//
//nolint:dupl // Extraction functions are intentionally similar in structure
package specgraph

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ProductBuildersHQ/visionspec/pkg/config"
	"github.com/ProductBuildersHQ/visionspec/pkg/types"
	"github.com/plexusone/graphfs/pkg/graph"
	"github.com/plexusone/graphize/pkg/exporters/graphml"
	"github.com/plexusone/graphize/pkg/exporters/htmlsite"
)

// Node types for spec graphs.
const (
	NodeTypeRequirement    = "requirement"
	NodeTypeUserStory      = "user_story"
	NodeTypeConstraint     = "constraint"
	NodeTypeAcceptanceCrit = "acceptance_criteria"
	NodeTypeDecision       = "decision"
	NodeTypeTradeoff       = "tradeoff"
	NodeTypeCapability     = "capability"
	NodeTypeSection        = "section"
	NodeTypeSpec           = "spec"
)

// Edge types for spec relationships.
const (
	EdgeTypeTracesTo      = "traces_to"
	EdgeTypeDerivedFrom   = "derived_from"
	EdgeTypeConflictsWith = "conflicts_with"
	EdgeTypeSatisfies     = "satisfies"
	EdgeTypeDependsOn     = "depends_on"
	EdgeTypeBlocks        = "blocks"
	EdgeTypeSupersedes    = "supersedes"
	EdgeTypeContains      = "contains"
)

// SpecExtractor extracts requirement graphs from visionspec projects.
type SpecExtractor struct {
	projectPath string
}

// NewSpecExtractor creates a new spec extractor for a project.
func NewSpecExtractor(projectPath string) *SpecExtractor {
	return &SpecExtractor{projectPath: projectPath}
}

// Extract builds a requirement graph from all specs in the project.
func (e *SpecExtractor) Extract() (*graph.Graph, error) {
	g := graph.NewGraph()

	// Get all spec files
	specTypes := []types.SpecType{
		types.SpecTypeMRD,
		types.SpecTypePRD,
		types.SpecTypeUXD,
		types.SpecTypeTRD,
		types.SpecTypeIRD,
		types.SpecTypePress,
		types.SpecTypeFAQ,
		types.SpecTypeNarrative1P,
		types.SpecTypeNarrative6P,
	}

	for _, specType := range specTypes {
		specPath := config.SpecPath(e.projectPath, specType)
		if _, err := os.Stat(specPath); os.IsNotExist(err) {
			continue
		}

		nodes, edges, err := e.extractSpec(specPath, specType)
		if err != nil {
			return nil, fmt.Errorf("extracting %s: %w", specType, err)
		}

		for _, node := range nodes {
			g.AddNode(node)
		}
		for _, edge := range edges {
			g.AddEdge(edge)
		}
	}

	// Add cross-spec edges based on traceability
	e.addTraceabilityEdges(g)

	return g, nil
}

// extractSpec extracts nodes and edges from a single spec file.
func (e *SpecExtractor) extractSpec(path string, specType types.SpecType) ([]*graph.Node, []*graph.Edge, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	var nodes []*graph.Node
	var edges []*graph.Edge

	// Create spec root node
	relPath, _ := filepath.Rel(e.projectPath, path)
	specID := makeSpecID(specType)
	nodes = append(nodes, &graph.Node{
		ID:    specID,
		Type:  NodeTypeSpec,
		Label: string(specType),
		Attrs: map[string]string{
			"path":      relPath,
			"spec_type": string(specType),
		},
	})

	// Parse markdown content
	sections := parseMarkdownSections(string(content))

	for _, section := range sections {
		sectionID := makeID(specType, "section", section.Title)
		nodes = append(nodes, &graph.Node{
			ID:    sectionID,
			Type:  NodeTypeSection,
			Label: section.Title,
			Attrs: map[string]string{
				"level":     fmt.Sprintf("%d", section.Level),
				"spec_type": string(specType),
			},
		})

		edges = append(edges, &graph.Edge{
			From: specID,
			To:   sectionID,
			Type: EdgeTypeContains,
		})

		// Extract specific node types based on spec type and section
		extracted := e.extractFromSection(section, specType)
		for _, node := range extracted.Nodes {
			nodes = append(nodes, node)
			edges = append(edges, &graph.Edge{
				From: sectionID,
				To:   node.ID,
				Type: EdgeTypeContains,
			})
		}
		edges = append(edges, extracted.Edges...)
	}

	return nodes, edges, nil
}

// Section represents a markdown section.
type Section struct {
	Title   string
	Level   int
	Content string
}

// parseMarkdownSections parses markdown into sections.
func parseMarkdownSections(content string) []Section {
	var sections []Section
	lines := strings.Split(content, "\n")

	headingPattern := regexp.MustCompile(`^(#{1,6})\s+(.+)$`)

	var current *Section
	var contentLines []string

	for _, line := range lines {
		if matches := headingPattern.FindStringSubmatch(line); len(matches) > 0 {
			// Save previous section
			if current != nil {
				current.Content = strings.Join(contentLines, "\n")
				sections = append(sections, *current)
			}

			current = &Section{
				Title: strings.TrimSpace(matches[2]),
				Level: len(matches[1]),
			}
			contentLines = nil
		} else {
			contentLines = append(contentLines, line)
		}
	}

	// Save last section
	if current != nil {
		current.Content = strings.Join(contentLines, "\n")
		sections = append(sections, *current)
	}

	return sections
}

// ExtractedContent holds nodes and edges extracted from a section.
type ExtractedContent struct {
	Nodes []*graph.Node
	Edges []*graph.Edge
}

// extractFromSection extracts specific nodes based on spec type and section content.
func (e *SpecExtractor) extractFromSection(section Section, specType types.SpecType) ExtractedContent {
	var result ExtractedContent

	// Extract requirements from PRD sections
	if specType == types.SpecTypePRD {
		result.Nodes = append(result.Nodes, e.extractRequirements(section, specType)...)
		result.Nodes = append(result.Nodes, e.extractUserStories(section, specType)...)
		result.Nodes = append(result.Nodes, e.extractAcceptanceCriteria(section, specType)...)
	}

	// Extract constraints from TRD
	if specType == types.SpecTypeTRD {
		result.Nodes = append(result.Nodes, e.extractConstraints(section, specType)...)
		result.Nodes = append(result.Nodes, e.extractDecisions(section, specType)...)
	}

	// Extract constraints from MRD
	if specType == types.SpecTypeMRD {
		result.Nodes = append(result.Nodes, e.extractConstraints(section, specType)...)
	}

	return result
}

// extractRequirements extracts requirement nodes from section content.
func (e *SpecExtractor) extractRequirements(section Section, specType types.SpecType) []*graph.Node {
	var nodes []*graph.Node

	// Look for requirement patterns: "REQ-XXX:", "FR-XXX:", numbered lists with "shall/must"
	reqPattern := regexp.MustCompile(`(?m)^[\-\*\d\.]+\s*(REQ-\d+|FR-\d+)?:?\s*(.+(?:shall|must|should).+)$`)
	matches := reqPattern.FindAllStringSubmatch(section.Content, -1)

	for i, match := range matches {
		id := ""
		text := ""
		if len(match) > 2 && match[1] != "" {
			id = match[1]
			text = match[2]
		} else if len(match) > 2 {
			id = fmt.Sprintf("REQ-%s-%d", strings.ToUpper(string(specType)), i+1)
			text = match[2]
		} else {
			continue
		}

		nodeID := makeID(specType, "req", id)
		nodes = append(nodes, &graph.Node{
			ID:    nodeID,
			Type:  NodeTypeRequirement,
			Label: truncate(text, 60),
			Attrs: map[string]string{
				"req_id":    id,
				"full_text": text,
				"spec_type": string(specType),
				"section":   section.Title,
			},
		})
	}

	return nodes
}

// extractUserStories extracts user story nodes.
func (e *SpecExtractor) extractUserStories(section Section, specType types.SpecType) []*graph.Node {
	var nodes []*graph.Node

	// Pattern: "As a <role>, I want <goal> so that <benefit>"
	storyPattern := regexp.MustCompile(`(?i)As a\s+([^,]+),\s*I want\s+([^,]+)(?:,?\s*so that\s+(.+))?`)
	matches := storyPattern.FindAllStringSubmatch(section.Content, -1)

	for i, match := range matches {
		if len(match) < 3 {
			continue
		}
		role := strings.TrimSpace(match[1])
		goal := strings.TrimSpace(match[2])
		benefit := ""
		if len(match) > 3 {
			benefit = strings.TrimSpace(match[3])
		}

		id := fmt.Sprintf("US-%s-%d", strings.ToUpper(string(specType)), i+1)
		nodeID := makeID(specType, "story", id)

		label := fmt.Sprintf("As %s, %s", role, truncate(goal, 40))
		nodes = append(nodes, &graph.Node{
			ID:    nodeID,
			Type:  NodeTypeUserStory,
			Label: label,
			Attrs: map[string]string{
				"story_id":  id,
				"role":      role,
				"goal":      goal,
				"benefit":   benefit,
				"spec_type": string(specType),
				"section":   section.Title,
			},
		})
	}

	return nodes
}

// extractAcceptanceCriteria extracts acceptance criteria nodes.
func (e *SpecExtractor) extractAcceptanceCriteria(section Section, specType types.SpecType) []*graph.Node {
	var nodes []*graph.Node

	// Pattern: "Given/When/Then" or "AC-XXX:"
	acPattern := regexp.MustCompile(`(?mi)^[\-\*]\s*(AC-\d+:?\s*)?(?:Given|When|Then)\s+(.+)$`)
	matches := acPattern.FindAllStringSubmatch(section.Content, -1)

	for i, match := range matches {
		if len(match) < 3 {
			continue
		}
		id := strings.TrimSpace(match[1])
		if id == "" {
			id = fmt.Sprintf("AC-%s-%d", strings.ToUpper(string(specType)), i+1)
		}
		text := strings.TrimSpace(match[2])

		nodeID := makeID(specType, "ac", id)
		nodes = append(nodes, &graph.Node{
			ID:    nodeID,
			Type:  NodeTypeAcceptanceCrit,
			Label: truncate(text, 60),
			Attrs: map[string]string{
				"ac_id":     id,
				"full_text": text,
				"spec_type": string(specType),
				"section":   section.Title,
			},
		})
	}

	return nodes
}

// extractConstraints extracts constraint nodes.
//
//nolint:dupl // Similar structure to extractDecisions is intentional for clarity
func (e *SpecExtractor) extractConstraints(section Section, specType types.SpecType) []*graph.Node {
	var nodes []*graph.Node

	// Look for constraint keywords
	constraintKeywords := []string{"constraint", "limitation", "must not", "cannot", "prohibited"}
	sectionLower := strings.ToLower(section.Title)

	isConstraintSection := false
	for _, keyword := range constraintKeywords {
		if strings.Contains(sectionLower, keyword) {
			isConstraintSection = true
			break
		}
	}

	if !isConstraintSection {
		return nodes
	}

	// Extract list items as constraints
	listPattern := regexp.MustCompile(`(?m)^[\-\*\d\.]+\s+(.+)$`)
	matches := listPattern.FindAllStringSubmatch(section.Content, -1)

	for i, match := range matches {
		if len(match) < 2 {
			continue
		}
		text := strings.TrimSpace(match[1])
		if len(text) < 10 {
			continue
		}

		id := fmt.Sprintf("CON-%s-%d", strings.ToUpper(string(specType)), i+1)
		nodeID := makeID(specType, "constraint", id)
		nodes = append(nodes, &graph.Node{
			ID:    nodeID,
			Type:  NodeTypeConstraint,
			Label: truncate(text, 60),
			Attrs: map[string]string{
				"constraint_id": id,
				"full_text":     text,
				"spec_type":     string(specType),
				"section":       section.Title,
			},
		})
	}

	return nodes
}

// extractDecisions extracts architectural decision nodes.
func (e *SpecExtractor) extractDecisions(section Section, specType types.SpecType) []*graph.Node {
	var nodes []*graph.Node

	// Look for decision keywords in section title
	decisionKeywords := []string{"decision", "architecture", "design", "approach", "solution"}
	sectionLower := strings.ToLower(section.Title)

	isDecisionSection := false
	for _, keyword := range decisionKeywords {
		if strings.Contains(sectionLower, keyword) {
			isDecisionSection = true
			break
		}
	}

	if !isDecisionSection {
		return nodes
	}

	// Extract ADR-style decisions or list items
	listPattern := regexp.MustCompile(`(?m)^[\-\*\d\.]+\s+(.+)$`)
	matches := listPattern.FindAllStringSubmatch(section.Content, -1)

	for i, match := range matches {
		if len(match) < 2 {
			continue
		}
		text := strings.TrimSpace(match[1])
		if len(text) < 10 {
			continue
		}

		id := fmt.Sprintf("DEC-%s-%d", strings.ToUpper(string(specType)), i+1)
		nodeID := makeID(specType, "decision", id)
		nodes = append(nodes, &graph.Node{
			ID:    nodeID,
			Type:  NodeTypeDecision,
			Label: truncate(text, 60),
			Attrs: map[string]string{
				"decision_id": id,
				"full_text":   text,
				"spec_type":   string(specType),
				"section":     section.Title,
			},
		})
	}

	return nodes
}

// addTraceabilityEdges adds cross-spec traceability edges.
func (e *SpecExtractor) addTraceabilityEdges(g *graph.Graph) {
	// Collect nodes by type
	var prdReqs []*graph.Node
	var trdDecisions []*graph.Node
	for _, node := range g.Nodes {
		if node.Type == NodeTypeRequirement {
			prdReqs = append(prdReqs, node)
		}
		if node.Type == NodeTypeDecision {
			trdDecisions = append(trdDecisions, node)
		}
	}

	// PRD requirements trace to TRD decisions
	for _, req := range prdReqs {
		// Simple text matching for traceability
		reqText := strings.ToLower(req.Attrs["full_text"])
		for _, dec := range trdDecisions {
			decText := strings.ToLower(dec.Attrs["full_text"])
			if textSimilarity(reqText, decText) > 0.3 {
				g.AddEdge(&graph.Edge{
					From:            req.ID,
					To:              dec.ID,
					Type:            EdgeTypeTracesTo,
					Confidence:      graph.ConfidenceInferred,
					ConfidenceScore: textSimilarity(reqText, decText),
					Attrs: map[string]string{
						"reason": "Text similarity",
					},
				})
			}
		}
	}

	// TRD is derived from PRD and MRD
	trdSpec := g.GetNode(makeSpecID(types.SpecTypeTRD))
	prdSpec := g.GetNode(makeSpecID(types.SpecTypePRD))
	mrdSpec := g.GetNode(makeSpecID(types.SpecTypeMRD))

	if trdSpec != nil && prdSpec != nil {
		g.AddEdge(&graph.Edge{
			From:       trdSpec.ID,
			To:         prdSpec.ID,
			Type:       EdgeTypeDerivedFrom,
			Confidence: graph.ConfidenceExtracted,
		})
	}
	if trdSpec != nil && mrdSpec != nil {
		g.AddEdge(&graph.Edge{
			From:       trdSpec.ID,
			To:         mrdSpec.ID,
			Type:       EdgeTypeDerivedFrom,
			Confidence: graph.ConfidenceExtracted,
		})
	}
}

// Helper functions

func makeSpecID(specType types.SpecType) string {
	return "spec_" + strings.ToLower(string(specType))
}

func makeID(specType types.SpecType, nodeType, name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")
	name = regexp.MustCompile(`[^a-z0-9_]`).ReplaceAllString(name, "")
	return fmt.Sprintf("%s_%s_%s", strings.ToLower(string(specType)), nodeType, name)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// textSimilarity computes a simple word overlap similarity score.
func textSimilarity(a, b string) float64 {
	wordsA := strings.Fields(a)
	wordsB := strings.Fields(b)

	if len(wordsA) == 0 || len(wordsB) == 0 {
		return 0
	}

	wordSet := make(map[string]bool)
	for _, w := range wordsA {
		if len(w) > 3 { // Skip short words
			wordSet[w] = true
		}
	}

	overlap := 0
	for _, w := range wordsB {
		if len(w) > 3 && wordSet[w] {
			overlap++
		}
	}

	return float64(overlap) / float64(max(len(wordsA), len(wordsB)))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// SaveJSON saves a graph to a JSON file.
func SaveJSON(g *graph.Graph, path string) error {
	data, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadJSON loads a graph from a JSON file.
func LoadJSON(path string) (*graph.Graph, error) {
	data, err := os.ReadFile(path) //nolint:gosec // G304: path is from trusted source
	if err != nil {
		return nil, err
	}
	var g graph.Graph
	if err := json.Unmarshal(data, &g); err != nil {
		return nil, err
	}
	return &g, nil
}

// QueryFilter specifies filtering criteria for graph queries.
type QueryFilter struct {
	NodeType string // Filter by node type (e.g., "requirement", "user_story")
	SpecType string // Filter by spec type (e.g., "prd", "trd")
}

// QueryResult contains filtered nodes from a query.
type QueryResult struct {
	Nodes []*graph.Node
	Count int
}

// Query filters nodes in a graph based on the provided filter criteria.
func Query(g *graph.Graph, filter QueryFilter) *QueryResult {
	var nodes []*graph.Node
	for _, node := range g.Nodes {
		if filter.NodeType != "" && node.Type != filter.NodeType {
			continue
		}
		if filter.SpecType != "" && node.Attrs["spec_type"] != filter.SpecType {
			continue
		}
		nodes = append(nodes, node)
	}
	return &QueryResult{
		Nodes: nodes,
		Count: len(nodes),
	}
}

// ExportFormat specifies the output format for graph export.
type ExportFormat string

const (
	FormatHTML    ExportFormat = "html"
	FormatGraphML ExportFormat = "graphml"
	FormatJSON    ExportFormat = "json"
)

// ExportOptions configures graph export behavior.
type ExportOptions struct {
	Format    ExportFormat
	OutputDir string
	Title     string
}

// ExportResult contains information about the export operation.
type ExportResult struct {
	OutputPath string
	Format     ExportFormat
	NodeCount  int
	EdgeCount  int
}

// Export exports a graph to the specified format.
func Export(g *graph.Graph, opts ExportOptions) (*ExportResult, error) {
	if err := os.MkdirAll(opts.OutputDir, 0755); err != nil {
		return nil, fmt.Errorf("creating output directory: %w", err)
	}

	// Convert nodes map to slice for exporters
	var nodes []*graph.Node
	for _, node := range g.Nodes {
		nodes = append(nodes, node)
	}

	result := &ExportResult{
		Format:    opts.Format,
		NodeCount: len(nodes),
		EdgeCount: len(g.Edges),
	}

	switch opts.Format {
	case FormatHTML:
		gen := htmlsite.NewGenerator()
		if opts.Title != "" {
			gen.Title = opts.Title
		} else {
			gen.Title = "Spec Graph"
		}
		content, err := gen.Generate(nodes, g.Edges, nil)
		if err != nil {
			return nil, fmt.Errorf("generating HTML: %w", err)
		}
		result.OutputPath = filepath.Join(opts.OutputDir, "graph.html")
		if err := os.WriteFile(result.OutputPath, content.Index, 0600); err != nil {
			return nil, fmt.Errorf("writing HTML: %w", err)
		}

	case FormatGraphML:
		gen := graphml.NewGenerator()
		gen.GraphID = "spec-graph"
		gen.Description = "VisionSpec requirement graph"
		genResult, err := gen.Generate(nodes, g.Edges)
		if err != nil {
			return nil, fmt.Errorf("generating GraphML: %w", err)
		}
		result.OutputPath = filepath.Join(opts.OutputDir, "spec-graph.graphml")
		if err := os.WriteFile(result.OutputPath, genResult.Data, 0600); err != nil {
			return nil, fmt.Errorf("writing GraphML: %w", err)
		}

	case FormatJSON:
		result.OutputPath = filepath.Join(opts.OutputDir, "spec-graph.json")
		if err := SaveJSON(g, result.OutputPath); err != nil {
			return nil, fmt.Errorf("saving JSON: %w", err)
		}

	default:
		return nil, fmt.Errorf("unsupported format: %s", opts.Format)
	}

	return result, nil
}

// Metrics contains computed statistics from a spec graph.
type Metrics struct {
	TotalRequirements int
	TotalUserStories  int
	TotalConstraints  int
	TotalDecisions    int
	TraceCoverage     float64 // Percentage of requirements traced to TRD
	ConflictCount     int
}

// ComputeMetrics computes metrics from a graph.
func ComputeMetrics(g *graph.Graph) *Metrics {
	m := &Metrics{}

	// Count nodes by type
	reqsTraced := 0
	for _, node := range g.Nodes {
		switch node.Type {
		case NodeTypeRequirement:
			m.TotalRequirements++
		case NodeTypeUserStory:
			m.TotalUserStories++
		case NodeTypeConstraint:
			m.TotalConstraints++
		case NodeTypeDecision:
			m.TotalDecisions++
		}
	}

	// Count traced requirements (requirements with traces_to edges)
	for _, edge := range g.Edges {
		if edge.Type == EdgeTypeTracesTo {
			// Check if from node is a requirement
			fromNode := g.GetNode(edge.From)
			if fromNode != nil && fromNode.Type == NodeTypeRequirement {
				reqsTraced++
			}
		}
		if edge.Type == EdgeTypeConflictsWith {
			m.ConflictCount++
		}
	}

	// Calculate trace coverage
	if m.TotalRequirements > 0 {
		m.TraceCoverage = float64(reqsTraced) / float64(m.TotalRequirements)
	}

	return m
}
