package context

import (
	"fmt"
	"strings"
)

// GenerateSummary creates an LLM-friendly summary from aggregated context.
func (ac *AggregatedContext) GenerateSummary() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# Context Summary for %s\n\n", ac.Project))
	sb.WriteString(fmt.Sprintf("Gathered: %s (took %s)\n", ac.GatheredAt.Format("2006-01-02 15:04:05"), ac.Duration))
	sb.WriteString(fmt.Sprintf("Sources: %d\n\n", len(ac.Sources)))

	// Code context summary
	if ac.HasCode {
		sb.WriteString("## Codebase\n\n")
		for _, code := range ac.CodeContexts() {
			sb.WriteString(fmt.Sprintf("### %s", code.RepoPath))
			if code.Branch != "" {
				sb.WriteString(fmt.Sprintf(" (%s@%s)", code.Branch, code.Commit))
			}
			sb.WriteString("\n\n")

			if len(code.Languages) > 0 {
				sb.WriteString("**Languages:** ")
				langs := make([]string, 0, len(code.Languages))
				for lang, loc := range code.Languages {
					langs = append(langs, fmt.Sprintf("%s (%d LOC)", lang, loc))
				}
				sb.WriteString(strings.Join(langs, ", "))
				sb.WriteString("\n\n")
			}

			if len(code.Dependencies) > 0 {
				sb.WriteString("**Key Dependencies:**\n")
				for i, dep := range code.Dependencies {
					if i >= 10 {
						sb.WriteString(fmt.Sprintf("- ... and %d more\n", len(code.Dependencies)-10))
						break
					}
					sb.WriteString(fmt.Sprintf("- %s@%s (%s)\n", dep.Name, dep.Version, dep.Source))
				}
				sb.WriteString("\n")
			}

			if len(code.APIs) > 0 {
				sb.WriteString("**APIs:**\n")
				for _, api := range code.APIs {
					sb.WriteString(fmt.Sprintf("- %s (%s)", api.Path, api.Format))
					if api.Title != "" {
						sb.WriteString(fmt.Sprintf(": %s", api.Title))
					}
					sb.WriteString("\n")
					for _, route := range api.Routes {
						sb.WriteString(fmt.Sprintf("  - %s %s\n", route.Method, route.Path))
					}
				}
				sb.WriteString("\n")
			}
		}
	}

	// Graph context summary
	if ac.HasGraph {
		sb.WriteString("## Requirement Graphs\n\n")
		for _, graph := range ac.GraphContexts() {
			sb.WriteString(fmt.Sprintf("### %s\n\n", graph.GraphPath))
			sb.WriteString(fmt.Sprintf("Nodes: %d, Edges: %d\n\n", graph.NodeCount, graph.EdgeCount))

			if len(graph.Requirements) > 0 {
				sb.WriteString("**Requirements:**\n")
				for i, req := range graph.Requirements {
					if i >= 10 {
						sb.WriteString(fmt.Sprintf("- ... and %d more\n", len(graph.Requirements)-10))
						break
					}
					sb.WriteString(fmt.Sprintf("- %s: %s\n", req.ID, req.Title))
				}
				sb.WriteString("\n")
			}

			if len(graph.Decisions) > 0 {
				sb.WriteString("**Architectural Decisions:**\n")
				for i, dec := range graph.Decisions {
					if i >= 5 {
						sb.WriteString(fmt.Sprintf("- ... and %d more\n", len(graph.Decisions)-5))
						break
					}
					sb.WriteString(fmt.Sprintf("- %s: %s\n", dec.ID, dec.Title))
				}
				sb.WriteString("\n")
			}

			if len(graph.Traceability) > 0 {
				sb.WriteString("**Traceability Links:**\n")
				for i, link := range graph.Traceability {
					if i >= 10 {
						sb.WriteString(fmt.Sprintf("- ... and %d more\n", len(graph.Traceability)-10))
						break
					}
					sb.WriteString(fmt.Sprintf("- %s → %s (%s)\n", link.FromID, link.ToID, link.Relation))
				}
				sb.WriteString("\n")
			}
		}
	}

	// External context summary
	if ac.HasExternal {
		sb.WriteString("## External Tools\n\n")
		for _, ext := range ac.ExternalContexts() {
			sb.WriteString(fmt.Sprintf("### %s (%s)\n\n", ext.ServerName, ext.ServerType))

			if len(ext.Issues) > 0 {
				sb.WriteString("**Issues:**\n")
				for i, issue := range ext.Issues {
					if i >= 10 {
						sb.WriteString(fmt.Sprintf("- ... and %d more\n", len(ext.Issues)-10))
						break
					}
					sb.WriteString(fmt.Sprintf("- %s [%s]: %s\n", issue.Key, issue.Status, issue.Summary))
				}
				sb.WriteString("\n")
			}

			if len(ext.Pages) > 0 {
				sb.WriteString("**Pages:**\n")
				for i, page := range ext.Pages {
					if i >= 5 {
						sb.WriteString(fmt.Sprintf("- ... and %d more\n", len(ext.Pages)-5))
						break
					}
					sb.WriteString(fmt.Sprintf("- %s\n", page.Title))
				}
				sb.WriteString("\n")
			}
		}
	}

	// File context summary
	if ac.HasFiles {
		sb.WriteString("## Local Files\n\n")
		for _, file := range ac.FileContexts() {
			sb.WriteString(fmt.Sprintf("- %s (%s)\n", file.Path, file.Type))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// GenerateCodeSummary creates a summary for a single code context.
func GenerateCodeSummary(code *CodeContext) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Repository: %s\n", code.RepoPath))
	if code.Branch != "" {
		sb.WriteString(fmt.Sprintf("Branch: %s (commit: %s)\n", code.Branch, code.Commit))
	}

	if len(code.Languages) > 0 {
		sb.WriteString("Languages: ")
		langs := make([]string, 0, len(code.Languages))
		for lang, loc := range code.Languages {
			langs = append(langs, fmt.Sprintf("%s=%d", lang, loc))
		}
		sb.WriteString(strings.Join(langs, ", "))
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("Dependencies: %d\n", len(code.Dependencies)))
	sb.WriteString(fmt.Sprintf("APIs: %d\n", len(code.APIs)))

	return sb.String()
}

// GenerateGraphSummary creates a summary for a single graph context.
func GenerateGraphSummary(graph *GraphContext) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Graph: %s\n", graph.GraphPath))
	if graph.Version != "" {
		sb.WriteString(fmt.Sprintf("Version: %s", graph.Version))
		if graph.Tool != "" {
			sb.WriteString(fmt.Sprintf(" (tool: %s)", graph.Tool))
		}
		sb.WriteString("\n")
	}
	sb.WriteString(fmt.Sprintf("Nodes: %d, Edges: %d\n", graph.NodeCount, graph.EdgeCount))

	// Node breakdown
	if graph.RequirementCount > 0 || graph.CodeCount > 0 || graph.TestCount > 0 {
		sb.WriteString(fmt.Sprintf("  Requirements: %d, Code: %d, Tests: %d\n",
			graph.RequirementCount, graph.CodeCount, graph.TestCount))
	}

	// Typed nodes
	if len(graph.Requirements) > 0 || len(graph.Decisions) > 0 || len(graph.Constraints) > 0 {
		sb.WriteString(fmt.Sprintf("Typed: %d requirements, %d decisions, %d constraints\n",
			len(graph.Requirements), len(graph.Decisions), len(graph.Constraints)))
	}

	// Coverage metrics
	if graph.RequirementCount > 0 {
		sb.WriteString(fmt.Sprintf("Coverage: %.1f%% code, %.1f%% test\n",
			graph.CodeCoverage, graph.TestCoverage))
	}

	// Traceability
	if graph.Stats.TraceabilityPct > 0 {
		sb.WriteString(fmt.Sprintf("Traceability: %.1f%%\n", graph.Stats.TraceabilityPct))
	}

	return sb.String()
}

// RenderTreeToString renders a tree node to a string representation.
func RenderTreeToString(node *TreeNode, prefix string, isLast bool) string {
	if node == nil {
		return ""
	}

	var sb strings.Builder

	// Add connector
	connector := "├── "
	if isLast {
		connector = "└── "
	}
	if prefix == "" {
		connector = ""
	}

	sb.WriteString(prefix)
	sb.WriteString(connector)
	sb.WriteString(node.Name)
	if node.Type == "dir" {
		sb.WriteString("/")
	}
	sb.WriteString("\n")

	// Update prefix for children
	childPrefix := prefix
	if prefix != "" {
		if isLast {
			childPrefix += "    "
		} else {
			childPrefix += "│   "
		}
	}

	// Render children
	for i, child := range node.Children {
		isLastChild := i == len(node.Children)-1
		sb.WriteString(RenderTreeToString(child, childPrefix, isLastChild))
	}

	return sb.String()
}
