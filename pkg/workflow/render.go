package workflow

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Renderer generates visual representations of workflows.
type Renderer interface {
	Render(w *Workflow) string
}

// MermaidRenderer renders workflows as Mermaid flowchart diagrams.
type MermaidRenderer struct {
	// Direction: TB (top-bottom), LR (left-right), BT, RL
	Direction string
	// ShowStatus includes status indicators
	ShowStatus bool
	// ShowPhases groups nodes by phase
	ShowPhases bool
}

// NewMermaidRenderer creates a Mermaid renderer with defaults.
func NewMermaidRenderer() *MermaidRenderer {
	return &MermaidRenderer{
		Direction:  "TB",
		ShowStatus: true,
		ShowPhases: true,
	}
}

// Render generates a Mermaid flowchart.
func (r *MermaidRenderer) Render(w *Workflow) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("flowchart %s\n", r.Direction))

	// Render phases as subgraphs
	if r.ShowPhases {
		for _, phase := range w.Phases {
			sb.WriteString(fmt.Sprintf("    subgraph %s[\"%s\"]\n", phase.ID, phase.Name))
			for _, nodeID := range phase.Nodes {
				if node, ok := w.Nodes[nodeID]; ok {
					sb.WriteString(fmt.Sprintf("        %s\n", r.renderNode(node)))
				}
			}
			sb.WriteString("    end\n")
		}
	} else {
		// Render nodes without phase grouping
		for _, node := range w.Nodes {
			sb.WriteString(fmt.Sprintf("    %s\n", r.renderNode(node)))
		}
	}

	// Render edges (dependencies)
	sb.WriteString("\n")
	for _, node := range w.Nodes {
		for _, depID := range node.DependsOn {
			sb.WriteString(fmt.Sprintf("    %s --> %s\n", depID, node.ID))
		}
	}

	// Add styling based on status
	if r.ShowStatus {
		sb.WriteString("\n")
		for _, node := range w.Nodes {
			style := r.statusStyle(node.Status)
			if style != "" {
				sb.WriteString(fmt.Sprintf("    style %s %s\n", node.ID, style))
			}
		}
	}

	return sb.String()
}

func (r *MermaidRenderer) renderNode(node *Node) string {
	// Use different shapes based on type
	shape := r.nodeShape(node)
	label := node.Name
	if r.ShowStatus {
		label = fmt.Sprintf("%s %s", r.statusIcon(node.Status), node.Name)
	}
	return fmt.Sprintf("%s%s", node.ID, fmt.Sprintf(shape, label))
}

func (r *MermaidRenderer) nodeShape(node *Node) string {
	if node.Automated {
		return "([%s])" // Stadium shape for automated
	}
	return "[\"%s\"]" // Rectangle for manual
}

func (r *MermaidRenderer) statusIcon(status Status) string {
	switch status {
	case StatusCompleted:
		return "✓"
	case StatusInProgress:
		return "◐"
	case StatusReady:
		return "○"
	case StatusBlocked:
		return "✗"
	case StatusSkipped:
		return "⊘"
	default:
		return "○"
	}
}

func (r *MermaidRenderer) statusStyle(status Status) string {
	switch status {
	case StatusCompleted:
		return "fill:#10b981,color:#fff"
	case StatusInProgress:
		return "fill:#3b82f6,color:#fff"
	case StatusReady:
		return "fill:#f59e0b,color:#fff"
	case StatusBlocked:
		return "fill:#ef4444,color:#fff"
	case StatusSkipped:
		return "fill:#6b7280,color:#fff"
	default:
		return "fill:#374151,color:#9ca3af"
	}
}

// DOTRenderer renders workflows as Graphviz DOT diagrams.
type DOTRenderer struct {
	// RankDir: TB, LR, BT, RL
	RankDir string
	// ShowStatus includes status indicators
	ShowStatus bool
	// Cluster groups nodes by phase
	Cluster bool
}

// NewDOTRenderer creates a DOT renderer with defaults.
func NewDOTRenderer() *DOTRenderer {
	return &DOTRenderer{
		RankDir:    "TB",
		ShowStatus: true,
		Cluster:    true,
	}
}

// Render generates a Graphviz DOT diagram.
func (r *DOTRenderer) Render(w *Workflow) string {
	var sb strings.Builder

	sb.WriteString("digraph workflow {\n")
	sb.WriteString(fmt.Sprintf("    rankdir=%s;\n", r.RankDir))
	sb.WriteString("    node [shape=box, style=rounded];\n")
	sb.WriteString("    edge [color=\"#6b7280\"];\n\n")

	// Render phases as clusters
	if r.Cluster {
		for i, phase := range w.Phases {
			sb.WriteString(fmt.Sprintf("    subgraph cluster_%d {\n", i))
			sb.WriteString(fmt.Sprintf("        label=\"%s\";\n", phase.Name))
			sb.WriteString("        style=dashed;\n")
			sb.WriteString("        color=\"#4b5563\";\n")
			for _, nodeID := range phase.Nodes {
				if node, ok := w.Nodes[nodeID]; ok {
					sb.WriteString(fmt.Sprintf("        %s;\n", r.renderNode(node)))
				}
			}
			sb.WriteString("    }\n\n")
		}
	} else {
		for _, node := range w.Nodes {
			sb.WriteString(fmt.Sprintf("    %s;\n", r.renderNode(node)))
		}
	}

	// Render edges
	sb.WriteString("\n")
	for _, node := range w.Nodes {
		for _, depID := range node.DependsOn {
			sb.WriteString(fmt.Sprintf("    %s -> %s;\n", depID, node.ID))
		}
	}

	sb.WriteString("}\n")
	return sb.String()
}

func (r *DOTRenderer) renderNode(node *Node) string {
	label := node.Name
	if r.ShowStatus {
		label = fmt.Sprintf("%s %s", r.statusIcon(node.Status), node.Name)
	}

	attrs := []string{
		fmt.Sprintf("label=\"%s\"", label),
	}

	// Style based on status
	if r.ShowStatus {
		fillColor, fontColor := r.statusColors(node.Status)
		attrs = append(attrs, fmt.Sprintf("fillcolor=\"%s\"", fillColor))
		attrs = append(attrs, fmt.Sprintf("fontcolor=\"%s\"", fontColor))
		attrs = append(attrs, "style=\"filled,rounded\"")
	}

	// Shape based on automated
	if node.Automated {
		attrs = append(attrs, "shape=ellipse")
	}

	return fmt.Sprintf("%s [%s]", node.ID, strings.Join(attrs, ", "))
}

func (r *DOTRenderer) statusIcon(status Status) string {
	switch status {
	case StatusCompleted:
		return "✓"
	case StatusInProgress:
		return "◐"
	case StatusReady:
		return "○"
	case StatusBlocked:
		return "✗"
	case StatusSkipped:
		return "⊘"
	default:
		return "○"
	}
}

func (r *DOTRenderer) statusColors(status Status) (fill, font string) {
	switch status {
	case StatusCompleted:
		return "#10b981", "#ffffff"
	case StatusInProgress:
		return "#3b82f6", "#ffffff"
	case StatusReady:
		return "#f59e0b", "#ffffff"
	case StatusBlocked:
		return "#ef4444", "#ffffff"
	case StatusSkipped:
		return "#6b7280", "#ffffff"
	default:
		return "#374151", "#9ca3af"
	}
}

// JSONRenderer renders workflows as JSON.
type JSONRenderer struct {
	Indent bool
}

// Render generates JSON representation.
func (r *JSONRenderer) Render(w *Workflow) string {
	var data []byte
	var err error
	if r.Indent {
		data, err = json.MarshalIndent(w, "", "  ")
	} else {
		data, err = json.Marshal(w)
	}
	if err != nil {
		return fmt.Sprintf(`{"error": %q}`, err.Error())
	}
	return string(data)
}
