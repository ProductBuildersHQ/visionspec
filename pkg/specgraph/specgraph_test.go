package specgraph

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/plexusone/graphfs/pkg/graph"
)

func TestSaveAndLoadJSON(t *testing.T) {
	// Create a test graph
	g := graph.NewGraph()
	g.AddNode(&graph.Node{
		ID:    "test_node_1",
		Type:  NodeTypeRequirement,
		Label: "Test Requirement",
		Attrs: map[string]string{
			"spec_type": "prd",
			"full_text": "The system shall do something",
		},
	})
	g.AddNode(&graph.Node{
		ID:    "test_node_2",
		Type:  NodeTypeDecision,
		Label: "Test Decision",
		Attrs: map[string]string{
			"spec_type": "trd",
		},
	})
	g.AddEdge(&graph.Edge{
		From: "test_node_1",
		To:   "test_node_2",
		Type: EdgeTypeTracesTo,
	})

	// Create temp directory
	tmpDir := t.TempDir()
	graphPath := filepath.Join(tmpDir, "test-graph.json")

	// Save
	err := SaveJSON(g, graphPath)
	if err != nil {
		t.Fatalf("SaveJSON failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(graphPath); os.IsNotExist(err) {
		t.Fatal("Graph file was not created")
	}

	// Load
	loaded, err := LoadJSON(graphPath)
	if err != nil {
		t.Fatalf("LoadJSON failed: %v", err)
	}

	// Verify
	if len(loaded.Nodes) != 2 {
		t.Errorf("Expected 2 nodes, got %d", len(loaded.Nodes))
	}
	if len(loaded.Edges) != 1 {
		t.Errorf("Expected 1 edge, got %d", len(loaded.Edges))
	}

	node1 := loaded.GetNode("test_node_1")
	if node1 == nil {
		t.Fatal("Node test_node_1 not found")
	}
	if node1.Type != NodeTypeRequirement {
		t.Errorf("Expected type %s, got %s", NodeTypeRequirement, node1.Type)
	}
}

func TestQuery(t *testing.T) {
	// Create a test graph
	g := graph.NewGraph()
	g.AddNode(&graph.Node{
		ID:    "prd_req_1",
		Type:  NodeTypeRequirement,
		Label: "PRD Requirement 1",
		Attrs: map[string]string{"spec_type": "prd"},
	})
	g.AddNode(&graph.Node{
		ID:    "prd_req_2",
		Type:  NodeTypeRequirement,
		Label: "PRD Requirement 2",
		Attrs: map[string]string{"spec_type": "prd"},
	})
	g.AddNode(&graph.Node{
		ID:    "trd_dec_1",
		Type:  NodeTypeDecision,
		Label: "TRD Decision 1",
		Attrs: map[string]string{"spec_type": "trd"},
	})
	g.AddNode(&graph.Node{
		ID:    "prd_story_1",
		Type:  NodeTypeUserStory,
		Label: "PRD User Story",
		Attrs: map[string]string{"spec_type": "prd"},
	})

	tests := []struct {
		name     string
		filter   QueryFilter
		expected int
	}{
		{
			name:     "no filter",
			filter:   QueryFilter{},
			expected: 4,
		},
		{
			name:     "filter by requirement type",
			filter:   QueryFilter{NodeType: NodeTypeRequirement},
			expected: 2,
		},
		{
			name:     "filter by prd spec",
			filter:   QueryFilter{SpecType: "prd"},
			expected: 3,
		},
		{
			name:     "filter by trd spec",
			filter:   QueryFilter{SpecType: "trd"},
			expected: 1,
		},
		{
			name:     "filter by type and spec",
			filter:   QueryFilter{NodeType: NodeTypeRequirement, SpecType: "prd"},
			expected: 2,
		},
		{
			name:     "no matches",
			filter:   QueryFilter{NodeType: NodeTypeConstraint},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Query(g, tt.filter)
			if result.Count != tt.expected {
				t.Errorf("Expected %d nodes, got %d", tt.expected, result.Count)
			}
			if len(result.Nodes) != tt.expected {
				t.Errorf("Expected %d nodes in slice, got %d", tt.expected, len(result.Nodes))
			}
		})
	}
}

func TestExport(t *testing.T) {
	// Create a test graph
	g := graph.NewGraph()
	g.AddNode(&graph.Node{
		ID:    "test_node",
		Type:  NodeTypeRequirement,
		Label: "Test",
	})

	tmpDir := t.TempDir()

	tests := []struct {
		name       string
		format     ExportFormat
		expectFile string
	}{
		{
			name:       "export json",
			format:     FormatJSON,
			expectFile: "spec-graph.json",
		},
		{
			name:       "export graphml",
			format:     FormatGraphML,
			expectFile: "spec-graph.graphml",
		},
		{
			name:       "export html",
			format:     FormatHTML,
			expectFile: "graph.html",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputDir := filepath.Join(tmpDir, string(tt.format))

			result, err := Export(g, ExportOptions{
				Format:    tt.format,
				OutputDir: outputDir,
				Title:     "Test Graph",
			})
			if err != nil {
				t.Fatalf("Export failed: %v", err)
			}

			expectedPath := filepath.Join(outputDir, tt.expectFile)
			if result.OutputPath != expectedPath {
				t.Errorf("Expected output path %s, got %s", expectedPath, result.OutputPath)
			}

			if _, err := os.Stat(result.OutputPath); os.IsNotExist(err) {
				t.Errorf("Output file was not created: %s", result.OutputPath)
			}

			if result.NodeCount != 1 {
				t.Errorf("Expected 1 node, got %d", result.NodeCount)
			}
		})
	}
}

func TestExportInvalidFormat(t *testing.T) {
	g := graph.NewGraph()
	tmpDir := t.TempDir()

	_, err := Export(g, ExportOptions{
		Format:    "invalid",
		OutputDir: tmpDir,
	})
	if err == nil {
		t.Error("Expected error for invalid format")
	}
}

func TestTextSimilarity(t *testing.T) {
	tests := []struct {
		a, b     string
		minScore float64
		maxScore float64
	}{
		{"hello world", "hello world", 0.9, 1.0},
		{"the quick brown fox", "the lazy brown dog", 0.2, 0.5},
		{"completely different", "nothing similar here", 0.0, 0.1},
		{"", "something", 0.0, 0.0},
		{"something", "", 0.0, 0.0},
	}

	for _, tt := range tests {
		score := textSimilarity(tt.a, tt.b)
		if score < tt.minScore || score > tt.maxScore {
			t.Errorf("textSimilarity(%q, %q) = %f, expected between %f and %f",
				tt.a, tt.b, score, tt.minScore, tt.maxScore)
		}
	}
}
