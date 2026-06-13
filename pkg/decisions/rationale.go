// Package decisions provides rationale graph generation from decision logs.
package decisions

import (
	"fmt"
	"sort"
	"strings"
)

// RationaleGraph represents the relationships between decisions.
type RationaleGraph struct {
	Nodes []RationaleNode `json:"nodes"`
	Edges []RationaleEdge `json:"edges"`
}

// RationaleNode represents a decision in the graph.
type RationaleNode struct {
	ID      string         `json:"id"`
	Title   string         `json:"title"`
	Status  DecisionStatus `json:"status"`
	Project string         `json:"project,omitempty"`
	Tags    []string       `json:"tags,omitempty"`
	Weight  int            `json:"weight"` // Number of connections
}

// RationaleEdge represents a relationship between decisions.
type RationaleEdge struct {
	From  string   `json:"from"`
	To    string   `json:"to"`
	Type  EdgeType `json:"type"`
	Label string   `json:"label,omitempty"`
}

// EdgeType categorizes the relationship between decisions.
type EdgeType string

const (
	EdgeTypeRelated    EdgeType = "related"
	EdgeTypeSupersedes EdgeType = "supersedes"
	EdgeTypeDependsOn  EdgeType = "depends_on"
	EdgeTypeConflicts  EdgeType = "conflicts"
	EdgeTypeEnables    EdgeType = "enables"
)

// BuildRationaleGraph creates a graph from a decision log.
func BuildRationaleGraph(log *DecisionLog) *RationaleGraph {
	graph := &RationaleGraph{
		Nodes: []RationaleNode{},
		Edges: []RationaleEdge{},
	}

	// Track connections for weight calculation
	connections := make(map[string]int)

	decisions := log.List()

	// Build edges first to count connections
	for _, d := range decisions {
		// Related decisions
		for _, relatedID := range d.Related {
			graph.Edges = append(graph.Edges, RationaleEdge{
				From:  d.ID,
				To:    relatedID,
				Type:  EdgeTypeRelated,
				Label: "related to",
			})
			connections[d.ID]++
			connections[relatedID]++
		}

		// Supersedes relationships
		if d.Supersedes != "" {
			graph.Edges = append(graph.Edges, RationaleEdge{
				From:  d.ID,
				To:    d.Supersedes,
				Type:  EdgeTypeSupersedes,
				Label: "supersedes",
			})
			connections[d.ID]++
			connections[d.Supersedes]++
		}
	}

	// Build nodes with weights
	for _, d := range decisions {
		node := RationaleNode{
			ID:      d.ID,
			Title:   d.Title,
			Status:  d.Status,
			Project: d.Project,
			Tags:    d.Tags,
			Weight:  connections[d.ID],
		}
		graph.Nodes = append(graph.Nodes, node)
	}

	return graph
}

// RenderMermaid exports the graph as a Mermaid diagram.
func (g *RationaleGraph) RenderMermaid() string {
	var sb strings.Builder

	sb.WriteString("graph TD\n")

	// Sort nodes for consistent output
	sortedNodes := make([]RationaleNode, len(g.Nodes))
	copy(sortedNodes, g.Nodes)
	sort.Slice(sortedNodes, func(i, j int) bool {
		return sortedNodes[i].ID < sortedNodes[j].ID
	})

	// Add nodes with styling based on status
	for _, node := range sortedNodes {
		nodeID := sanitizeMermaidID(node.ID)
		label := fmt.Sprintf("%s: %s", node.ID, truncateString(node.Title, 30))

		style := ""
		switch node.Status {
		case StatusAccepted:
			style = ":::accepted"
		case StatusProposed:
			style = ":::proposed"
		case StatusDeprecated:
			style = ":::deprecated"
		case StatusSuperseded:
			style = ":::superseded"
		case StatusRejected:
			style = ":::rejected"
		}

		sb.WriteString(fmt.Sprintf("    %s[\"%s\"]%s\n", nodeID, label, style))
	}

	sb.WriteString("\n")

	// Add edges
	for _, edge := range g.Edges {
		fromID := sanitizeMermaidID(edge.From)
		toID := sanitizeMermaidID(edge.To)

		arrow := "-->"
		switch edge.Type {
		case EdgeTypeSupersedes:
			arrow = "==>|supersedes|"
		case EdgeTypeRelated:
			arrow = "-.->|related|"
		case EdgeTypeDependsOn:
			arrow = "-->|depends on|"
		case EdgeTypeConflicts:
			arrow = "-.->|conflicts|"
		case EdgeTypeEnables:
			arrow = "-->|enables|"
		}

		sb.WriteString(fmt.Sprintf("    %s %s %s\n", fromID, arrow, toID))
	}

	// Add styles
	sb.WriteString("\n")
	sb.WriteString("    classDef accepted fill:#9f6,stroke:#333,stroke-width:2px\n")
	sb.WriteString("    classDef proposed fill:#ff9,stroke:#333,stroke-width:2px\n")
	sb.WriteString("    classDef deprecated fill:#f96,stroke:#333,stroke-width:2px\n")
	sb.WriteString("    classDef superseded fill:#999,stroke:#333,stroke-width:2px\n")
	sb.WriteString("    classDef rejected fill:#f66,stroke:#333,stroke-width:2px\n")

	return sb.String()
}

// RenderDOT exports the graph in Graphviz DOT format.
func (g *RationaleGraph) RenderDOT() string {
	var sb strings.Builder

	sb.WriteString("digraph RationaleGraph {\n")
	sb.WriteString("    rankdir=TB;\n")
	sb.WriteString("    node [shape=box, style=rounded];\n")
	sb.WriteString("\n")

	// Add nodes
	for _, node := range g.Nodes {
		nodeID := sanitizeDOTID(node.ID)
		label := fmt.Sprintf("%s\\n%s", node.ID, truncateString(node.Title, 25))

		color := "white"
		switch node.Status {
		case StatusAccepted:
			color = "lightgreen"
		case StatusProposed:
			color = "lightyellow"
		case StatusDeprecated:
			color = "lightsalmon"
		case StatusSuperseded:
			color = "lightgray"
		case StatusRejected:
			color = "lightcoral"
		}

		sb.WriteString(fmt.Sprintf("    %s [label=\"%s\", fillcolor=%s, style=filled];\n",
			nodeID, label, color))
	}

	sb.WriteString("\n")

	// Add edges
	for _, edge := range g.Edges {
		fromID := sanitizeDOTID(edge.From)
		toID := sanitizeDOTID(edge.To)

		style := ""
		switch edge.Type {
		case EdgeTypeSupersedes:
			style = "[style=bold, color=blue]"
		case EdgeTypeRelated:
			style = "[style=dashed]"
		case EdgeTypeConflicts:
			style = "[style=dotted, color=red]"
		}

		sb.WriteString(fmt.Sprintf("    %s -> %s %s;\n", fromID, toID, style))
	}

	sb.WriteString("}\n")
	return sb.String()
}

// GetClusters groups nodes by project.
func (g *RationaleGraph) GetClusters() map[string][]RationaleNode {
	clusters := make(map[string][]RationaleNode)
	for _, node := range g.Nodes {
		project := node.Project
		if project == "" {
			project = "unassigned"
		}
		clusters[project] = append(clusters[project], node)
	}
	return clusters
}

// GetConnectedComponents finds disconnected subgraphs.
func (g *RationaleGraph) GetConnectedComponents() [][]string {
	if len(g.Nodes) == 0 {
		return nil
	}

	// Build adjacency list
	adj := make(map[string][]string)
	nodeSet := make(map[string]bool)

	for _, node := range g.Nodes {
		nodeSet[node.ID] = true
		if _, ok := adj[node.ID]; !ok {
			adj[node.ID] = []string{}
		}
	}

	for _, edge := range g.Edges {
		adj[edge.From] = append(adj[edge.From], edge.To)
		adj[edge.To] = append(adj[edge.To], edge.From) // Undirected for components
	}

	// BFS to find components
	visited := make(map[string]bool)
	var components [][]string

	for nodeID := range nodeSet {
		if visited[nodeID] {
			continue
		}

		component := []string{}
		queue := []string{nodeID}

		for len(queue) > 0 {
			current := queue[0]
			queue = queue[1:]

			if visited[current] {
				continue
			}
			visited[current] = true
			component = append(component, current)

			for _, neighbor := range adj[current] {
				if !visited[neighbor] {
					queue = append(queue, neighbor)
				}
			}
		}

		sort.Strings(component)
		components = append(components, component)
	}

	return components
}

// FindCentralNodes returns nodes with the most connections.
func (g *RationaleGraph) FindCentralNodes(topN int) []RationaleNode {
	// Copy and sort by weight
	nodes := make([]RationaleNode, len(g.Nodes))
	copy(nodes, g.Nodes)

	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Weight > nodes[j].Weight
	})

	if topN > len(nodes) {
		topN = len(nodes)
	}

	return nodes[:topN]
}

// FindPath finds a path between two decisions (BFS).
func (g *RationaleGraph) FindPath(fromID, toID string) []string {
	if fromID == toID {
		return []string{fromID}
	}

	// Build adjacency list
	adj := make(map[string][]string)
	for _, edge := range g.Edges {
		adj[edge.From] = append(adj[edge.From], edge.To)
	}

	// BFS
	visited := make(map[string]bool)
	parent := make(map[string]string)
	queue := []string{fromID}
	visited[fromID] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == toID {
			// Reconstruct path
			path := []string{}
			for node := toID; node != ""; node = parent[node] {
				path = append([]string{node}, path...)
			}
			return path
		}

		for _, neighbor := range adj[current] {
			if !visited[neighbor] {
				visited[neighbor] = true
				parent[neighbor] = current
				queue = append(queue, neighbor)
			}
		}
	}

	return nil // No path found
}

// Helper functions

func sanitizeMermaidID(s string) string {
	// Mermaid IDs can't have hyphens at the start
	s = strings.ReplaceAll(s, "-", "_")
	return strings.ToLower(s)
}

func sanitizeDOTID(s string) string {
	// DOT IDs should be simple identifiers
	s = strings.ReplaceAll(s, "-", "_")
	return strings.ToLower(s)
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// ImpactAnalysis represents the impact of changing a decision.
type ImpactAnalysis struct {
	DecisionID         string   `json:"decision_id"`
	DirectlyAffected   []string `json:"directly_affected"`
	IndirectlyAffected []string `json:"indirectly_affected"`
	TotalImpact        int      `json:"total_impact"`
}

// AnalyzeImpact determines what would be affected by changing a decision.
func (g *RationaleGraph) AnalyzeImpact(decisionID string) *ImpactAnalysis {
	analysis := &ImpactAnalysis{
		DecisionID:         decisionID,
		DirectlyAffected:   []string{},
		IndirectlyAffected: []string{},
	}

	// Build reverse adjacency (what depends on this)
	dependents := make(map[string][]string)
	for _, edge := range g.Edges {
		if edge.Type == EdgeTypeSupersedes || edge.Type == EdgeTypeDependsOn {
			dependents[edge.To] = append(dependents[edge.To], edge.From)
		}
	}

	// Find directly affected
	analysis.DirectlyAffected = dependents[decisionID]

	// Find indirectly affected (transitive closure)
	visited := make(map[string]bool)
	visited[decisionID] = true
	for _, d := range analysis.DirectlyAffected {
		visited[d] = true
	}

	var findIndirect func(id string)
	findIndirect = func(id string) {
		for _, dep := range dependents[id] {
			if !visited[dep] {
				visited[dep] = true
				analysis.IndirectlyAffected = append(analysis.IndirectlyAffected, dep)
				findIndirect(dep)
			}
		}
	}

	for _, d := range analysis.DirectlyAffected {
		findIndirect(d)
	}

	analysis.TotalImpact = len(analysis.DirectlyAffected) + len(analysis.IndirectlyAffected)

	sort.Strings(analysis.DirectlyAffected)
	sort.Strings(analysis.IndirectlyAffected)

	return analysis
}
