// Package graphize provides graphize requirement traceability context source.
package graphize

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	ctx "github.com/plexusone/multispec/pkg/context"
)

// Source analyzes graphize requirement traceability data.
type Source struct {
	config ctx.GraphizeConfig
	name   string
}

// NewSource creates a new graphize source from configuration.
func NewSource(cfg ctx.GraphizeConfig) (*Source, error) {
	path := cfg.Path
	if path == "" {
		return nil, fmt.Errorf("graphize path is required")
	}

	// Check if .graphize directory exists
	graphizePath := filepath.Join(path, ".graphize")
	info, err := os.Stat(graphizePath)
	if err != nil {
		return nil, fmt.Errorf("graphize directory: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf(".graphize is not a directory: %s", graphizePath)
	}

	name := filepath.Base(path)

	return &Source{
		config: cfg,
		name:   name,
	}, nil
}

// Name returns the source identifier.
func (s *Source) Name() string {
	return fmt.Sprintf("graphize:%s", s.name)
}

// Type returns the source type.
func (s *Source) Type() ctx.SourceType {
	return ctx.SourceTypeGraphize
}

// String returns a human-readable description.
func (s *Source) String() string {
	return fmt.Sprintf("Graphize traceability: %s", s.config.Path)
}

// Fetch retrieves context from graphize data.
func (s *Source) Fetch(c context.Context) (*ctx.ContextData, error) {
	start := time.Now()

	graph := &ctx.GraphContext{}

	graphizePath := filepath.Join(s.config.Path, ".graphize")

	// Load nodes
	if err := s.loadNodes(graphizePath, graph); err != nil {
		return nil, fmt.Errorf("loading nodes: %w", err)
	}

	// Load edges
	if err := s.loadEdges(graphizePath, graph); err != nil {
		return nil, fmt.Errorf("loading edges: %w", err)
	}

	// Load metadata
	s.loadMetadata(graphizePath, graph)

	// Compute coverage metrics
	s.computeCoverage(graph)

	return &ctx.ContextData{
		Source:    s.Name(),
		Type:      ctx.SourceTypeGraphize,
		FetchedAt: time.Now(),
		Duration:  time.Since(start),
		Graph:     graph,
		Summary:   ctx.GenerateGraphSummary(graph),
	}, nil
}

// graphizeNode represents a node in the graphize format.
type graphizeNode struct {
	ID       string            `json:"id"`
	Type     string            `json:"type"`
	Label    string            `json:"label"`
	Path     string            `json:"path,omitempty"`
	Line     int               `json:"line,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// graphizeEdge represents an edge in the graphize format.
type graphizeEdge struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Type   string `json:"type"`
}

// graphizeGraph represents the complete graph.
type graphizeGraph struct {
	Nodes []graphizeNode `json:"nodes"`
	Edges []graphizeEdge `json:"edges"`
}

func (s *Source) loadNodes(graphizePath string, graph *ctx.GraphContext) error {
	// Try graph.json first (complete format)
	graphFile := filepath.Join(graphizePath, "graph.json")
	if data, err := os.ReadFile(graphFile); err == nil {
		var g graphizeGraph
		if err := json.Unmarshal(data, &g); err != nil {
			return fmt.Errorf("parsing graph.json: %w", err)
		}

		for _, node := range g.Nodes {
			gn := ctx.GraphNode{
				ID:       node.ID,
				Type:     node.Type,
				Label:    node.Label,
				Path:     node.Path,
				Line:     node.Line,
				Metadata: node.Metadata,
			}
			graph.Nodes = append(graph.Nodes, gn)

			// Count by type
			switch node.Type {
			case "requirement":
				graph.RequirementCount++
			case "code":
				graph.CodeCount++
			case "test":
				graph.TestCount++
			}
		}
		return nil
	}

	// Try nodes.json (separate format)
	nodesFile := filepath.Join(graphizePath, "nodes.json")
	data, err := os.ReadFile(nodesFile)
	if err != nil {
		return fmt.Errorf("reading nodes: %w", err)
	}

	var nodes []graphizeNode
	if err := json.Unmarshal(data, &nodes); err != nil {
		return fmt.Errorf("parsing nodes.json: %w", err)
	}

	for _, node := range nodes {
		gn := ctx.GraphNode{
			ID:       node.ID,
			Type:     node.Type,
			Label:    node.Label,
			Path:     node.Path,
			Line:     node.Line,
			Metadata: node.Metadata,
		}
		graph.Nodes = append(graph.Nodes, gn)

		switch node.Type {
		case "requirement":
			graph.RequirementCount++
		case "code":
			graph.CodeCount++
		case "test":
			graph.TestCount++
		}
	}

	return nil
}

func (s *Source) loadEdges(graphizePath string, graph *ctx.GraphContext) error {
	// Try graph.json first
	graphFile := filepath.Join(graphizePath, "graph.json")
	if data, err := os.ReadFile(graphFile); err == nil {
		var g graphizeGraph
		if err := json.Unmarshal(data, &g); err != nil {
			return nil // Already parsed in loadNodes
		}

		for _, edge := range g.Edges {
			graph.Edges = append(graph.Edges, ctx.GraphEdge{
				Source: edge.Source,
				Target: edge.Target,
				Type:   edge.Type,
			})
		}
		return nil
	}

	// Try edges.json
	edgesFile := filepath.Join(graphizePath, "edges.json")
	data, err := os.ReadFile(edgesFile)
	if err != nil {
		return nil // Edges are optional
	}

	var edges []graphizeEdge
	if err := json.Unmarshal(data, &edges); err != nil {
		return fmt.Errorf("parsing edges.json: %w", err)
	}

	for _, edge := range edges {
		graph.Edges = append(graph.Edges, ctx.GraphEdge{
			Source: edge.Source,
			Target: edge.Target,
			Type:   edge.Type,
		})
	}

	return nil
}

func (s *Source) loadMetadata(graphizePath string, graph *ctx.GraphContext) {
	metaFile := filepath.Join(graphizePath, "metadata.json")
	data, err := os.ReadFile(metaFile)
	if err != nil {
		return // Metadata is optional
	}

	var meta struct {
		Version   string `json:"version"`
		CreatedAt string `json:"created_at"`
		Tool      string `json:"tool"`
	}

	if err := json.Unmarshal(data, &meta); err != nil {
		return
	}

	graph.Version = meta.Version
	graph.Tool = meta.Tool
}

func (s *Source) computeCoverage(graph *ctx.GraphContext) {
	if graph.RequirementCount == 0 {
		return
	}

	// Count requirements with at least one code link
	linkedReqs := make(map[string]bool)
	testedReqs := make(map[string]bool)

	for _, edge := range graph.Edges {
		switch edge.Type {
		case "implements", "satisfies":
			// Code implements requirement
			linkedReqs[edge.Target] = true
		case "tests", "verifies":
			// Test verifies requirement
			testedReqs[edge.Target] = true
		}
	}

	graph.LinkedRequirements = len(linkedReqs)
	graph.TestedRequirements = len(testedReqs)

	if graph.RequirementCount > 0 {
		graph.CodeCoverage = float64(graph.LinkedRequirements) / float64(graph.RequirementCount) * 100
		graph.TestCoverage = float64(graph.TestedRequirements) / float64(graph.RequirementCount) * 100
	}
}
