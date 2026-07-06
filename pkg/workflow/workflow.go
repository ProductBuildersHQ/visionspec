// Package workflow provides a generic workflow intermediate representation (IR)
// for modeling directed acyclic graphs (DAGs) of work items with dependencies.
//
// This package is designed to be domain-agnostic and could be extracted
// as a standalone library for general workflow orchestration.
//
// Key concepts:
//   - Workflow: A named DAG with phases and nodes
//   - Phase: A logical grouping/stage in the workflow
//   - Node: A work item with dependencies, status, and metadata
//   - Status: The state of a node (pending, in_progress, completed, blocked)
package workflow

import (
	"encoding/json"
	"fmt"
	"sort"
)

// Status represents the state of a workflow node.
type Status string

const (
	StatusPending    Status = "pending"     // Not started, dependencies not met
	StatusReady      Status = "ready"       // Dependencies met, ready to start
	StatusInProgress Status = "in_progress" // Work in progress
	StatusCompleted  Status = "completed"   // Successfully completed
	StatusBlocked    Status = "blocked"     // Blocked by failed dependency
	StatusSkipped    Status = "skipped"     // Intentionally skipped
)

// NodeType categorizes nodes for grouping and visualization.
type NodeType string

// Workflow represents a directed acyclic graph of work items.
type Workflow struct {
	// Name identifies the workflow (e.g., "aws-working-backwards", "big-tech-product")
	Name string `json:"name" yaml:"name"`

	// Description provides human-readable context
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// Phases define logical groupings/stages in order
	Phases []Phase `json:"phases" yaml:"phases"`

	// Nodes are the work items in the workflow
	Nodes map[string]*Node `json:"nodes" yaml:"nodes"`
}

// Phase represents a logical stage in the workflow.
type Phase struct {
	// ID uniquely identifies the phase
	ID string `json:"id" yaml:"id"`

	// Name is the display name
	Name string `json:"name" yaml:"name"`

	// Description provides context
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// Order determines phase sequence (lower = earlier)
	Order int `json:"order" yaml:"order"`

	// Nodes lists the node IDs in this phase
	Nodes []string `json:"nodes" yaml:"nodes"`
}

// Node represents a single work item in the workflow.
type Node struct {
	// ID uniquely identifies the node
	ID string `json:"id" yaml:"id"`

	// Name is the display name
	Name string `json:"name" yaml:"name"`

	// Description provides context
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// Type categorizes the node (domain-specific, e.g., "source", "gtm", "technical")
	Type NodeType `json:"type,omitempty" yaml:"type,omitempty"`

	// Phase is the phase ID this node belongs to
	Phase string `json:"phase" yaml:"phase"`

	// DependsOn lists node IDs that must complete before this node can start
	DependsOn []string `json:"depends_on,omitempty" yaml:"depends_on,omitempty"`

	// Status is the current state of the node
	Status Status `json:"status" yaml:"status"`

	// Automated indicates if this node is machine-generated vs human-authored
	Automated bool `json:"automated,omitempty" yaml:"automated,omitempty"`

	// Metadata holds domain-specific data
	Metadata map[string]any `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// New creates a new empty workflow.
func New(name string) *Workflow {
	return &Workflow{
		Name:   name,
		Phases: make([]Phase, 0),
		Nodes:  make(map[string]*Node),
	}
}

// AddPhase adds a phase to the workflow.
func (w *Workflow) AddPhase(id, name string, order int) *Phase {
	phase := Phase{
		ID:    id,
		Name:  name,
		Order: order,
		Nodes: make([]string, 0),
	}
	w.Phases = append(w.Phases, phase)
	// Keep phases sorted by order
	sort.Slice(w.Phases, func(i, j int) bool {
		return w.Phases[i].Order < w.Phases[j].Order
	})
	return &w.Phases[len(w.Phases)-1]
}

// AddNode adds a node to the workflow.
func (w *Workflow) AddNode(node *Node) error {
	if node.ID == "" {
		return fmt.Errorf("node ID is required")
	}
	if _, exists := w.Nodes[node.ID]; exists {
		return fmt.Errorf("node %q already exists", node.ID)
	}
	if node.Status == "" {
		node.Status = StatusPending
	}
	w.Nodes[node.ID] = node

	// Add to phase's node list
	for i := range w.Phases {
		if w.Phases[i].ID == node.Phase {
			w.Phases[i].Nodes = append(w.Phases[i].Nodes, node.ID)
			break
		}
	}

	return nil
}

// GetNode returns a node by ID.
func (w *Workflow) GetNode(id string) (*Node, bool) {
	node, ok := w.Nodes[id]
	return node, ok
}

// Dependencies returns the direct dependencies of a node.
func (w *Workflow) Dependencies(nodeID string) []*Node {
	node, ok := w.Nodes[nodeID]
	if !ok {
		return nil
	}

	deps := make([]*Node, 0, len(node.DependsOn))
	for _, depID := range node.DependsOn {
		if dep, ok := w.Nodes[depID]; ok {
			deps = append(deps, dep)
		}
	}
	return deps
}

// Dependents returns nodes that depend on the given node.
func (w *Workflow) Dependents(nodeID string) []*Node {
	var dependents []*Node
	for _, node := range w.Nodes {
		for _, depID := range node.DependsOn {
			if depID == nodeID {
				dependents = append(dependents, node)
				break
			}
		}
	}
	return dependents
}

// IsReady returns true if all dependencies of a node are completed.
func (w *Workflow) IsReady(nodeID string) bool {
	node, ok := w.Nodes[nodeID]
	if !ok {
		return false
	}

	for _, depID := range node.DependsOn {
		dep, ok := w.Nodes[depID]
		if !ok {
			return false
		}
		if dep.Status != StatusCompleted && dep.Status != StatusSkipped {
			return false
		}
	}
	return true
}

// ReadyNodes returns all nodes that are ready to start.
func (w *Workflow) ReadyNodes() []*Node {
	var ready []*Node
	for _, node := range w.Nodes {
		if node.Status == StatusPending && w.IsReady(node.ID) {
			ready = append(ready, node)
		}
	}
	return ready
}

// UpdateStatus updates a node's status and returns affected nodes.
func (w *Workflow) UpdateStatus(nodeID string, status Status) error {
	node, ok := w.Nodes[nodeID]
	if !ok {
		return fmt.Errorf("node %q not found", nodeID)
	}
	node.Status = status
	return nil
}

// TopologicalSort returns nodes in dependency order.
func (w *Workflow) TopologicalSort() ([]*Node, error) {
	visited := make(map[string]bool)
	inStack := make(map[string]bool)
	var result []*Node

	var visit func(id string) error
	visit = func(id string) error {
		if inStack[id] {
			return fmt.Errorf("cycle detected at node %q", id)
		}
		if visited[id] {
			return nil
		}

		inStack[id] = true
		node := w.Nodes[id]
		for _, depID := range node.DependsOn {
			if err := visit(depID); err != nil {
				return err
			}
		}
		inStack[id] = false
		visited[id] = true
		result = append(result, node)
		return nil
	}

	for id := range w.Nodes {
		if err := visit(id); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// Progress returns completion statistics.
func (w *Workflow) Progress() (completed, total int, percent float64) {
	total = len(w.Nodes)
	for _, node := range w.Nodes {
		if node.Status == StatusCompleted || node.Status == StatusSkipped {
			completed++
		}
	}
	if total > 0 {
		percent = float64(completed) / float64(total) * 100
	}
	return
}

// Validate checks the workflow for errors.
func (w *Workflow) Validate() error {
	// Check for missing dependencies
	for id, node := range w.Nodes {
		for _, depID := range node.DependsOn {
			if _, ok := w.Nodes[depID]; !ok {
				return fmt.Errorf("node %q depends on non-existent node %q", id, depID)
			}
		}
	}

	// Check for cycles
	_, err := w.TopologicalSort()
	if err != nil {
		return err
	}

	return nil
}

// Clone creates a deep copy of the workflow.
func (w *Workflow) Clone() *Workflow {
	data, err := json.Marshal(w)
	if err != nil {
		panic(fmt.Sprintf("workflow Clone: marshal failed: %v", err))
	}
	var clone Workflow
	if err := json.Unmarshal(data, &clone); err != nil {
		panic(fmt.Sprintf("workflow Clone: unmarshal failed: %v", err))
	}
	return &clone
}

// NodesByPhase returns nodes grouped by phase in order.
func (w *Workflow) NodesByPhase() []struct {
	Phase Phase
	Nodes []*Node
} {
	var result []struct {
		Phase Phase
		Nodes []*Node
	}

	for _, phase := range w.Phases {
		nodes := make([]*Node, 0, len(phase.Nodes))
		for _, nodeID := range phase.Nodes {
			if node, ok := w.Nodes[nodeID]; ok {
				nodes = append(nodes, node)
			}
		}
		result = append(result, struct {
			Phase Phase
			Nodes []*Node
		}{Phase: phase, Nodes: nodes})
	}

	return result
}
