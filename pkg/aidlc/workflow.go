package aidlc

import (
	"fmt"
	"strings"
)

// Workflow represents the full AIDLC workflow DAG.
type Workflow struct {
	// Name is the workflow name.
	Name string `json:"name" yaml:"name"`

	// Description is a brief description of the workflow.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// Phases are the workflow phases in order.
	Phases []WorkflowPhase `json:"phases" yaml:"phases"`

	// Nodes are all workflow nodes indexed by ID.
	Nodes map[string]*WorkflowNode `json:"nodes" yaml:"nodes"`

	// Edges define dependencies between nodes.
	Edges []WorkflowEdge `json:"edges" yaml:"edges"`
}

// WorkflowPhase represents a phase in the workflow.
type WorkflowPhase struct {
	// ID is the phase identifier.
	ID string `json:"id" yaml:"id"`

	// Name is the display name.
	Name string `json:"name" yaml:"name"`

	// Description describes the phase purpose.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// Order is the phase order (0-indexed).
	Order int `json:"order" yaml:"order"`

	// NodeIDs are the node IDs in this phase.
	NodeIDs []string `json:"node_ids" yaml:"node_ids"`
}

// WorkflowNode represents a document node in the workflow.
type WorkflowNode struct {
	// ID is the unique node identifier.
	ID string `json:"id" yaml:"id"`

	// DocType is the AIDLC document type.
	DocType DocType `json:"doc_type" yaml:"doc_type"`

	// Phase is the workflow phase.
	Phase Phase `json:"phase" yaml:"phase"`

	// Name is the display name.
	Name string `json:"name" yaml:"name"`

	// Description describes the node purpose.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// Status is the node status.
	Status NodeStatus `json:"status" yaml:"status"`

	// Score is the quality score (if evaluated).
	Score *QualityScore `json:"score,omitempty" yaml:"score,omitempty"`

	// DependsOn lists node IDs this node depends on.
	DependsOn []string `json:"depends_on,omitempty" yaml:"depends_on,omitempty"`

	// Blocks lists node IDs blocked by this node.
	Blocks []string `json:"blocks,omitempty" yaml:"blocks,omitempty"`

	// Required indicates if this node is required.
	Required bool `json:"required" yaml:"required"`

	// Automated indicates if this node is LLM-generated.
	Automated bool `json:"automated,omitempty" yaml:"automated,omitempty"`

	// Metadata contains additional node data.
	Metadata map[string]any `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// NodeStatus represents the status of a workflow node.
type NodeStatus string

const (
	NodePending    NodeStatus = "pending"
	NodeReady      NodeStatus = "ready"
	NodeInProgress NodeStatus = "in_progress"
	NodeCompleted  NodeStatus = "completed"
	NodeBlocked    NodeStatus = "blocked"
	NodeSkipped    NodeStatus = "skipped"
	NodeFailed     NodeStatus = "failed"
)

// WorkflowEdge represents a dependency edge between nodes.
type WorkflowEdge struct {
	// From is the source node ID.
	From string `json:"from" yaml:"from"`

	// To is the target node ID.
	To string `json:"to" yaml:"to"`

	// Type is the edge type (dependency, blocks, suggests).
	Type EdgeType `json:"type" yaml:"type"`

	// Label is an optional edge label.
	Label string `json:"label,omitempty" yaml:"label,omitempty"`
}

// EdgeType represents the type of workflow edge.
type EdgeType string

const (
	// EdgeDependency indicates the target depends on the source.
	EdgeDependency EdgeType = "dependency"
	// EdgeBlocks indicates the source blocks the target.
	EdgeBlocks EdgeType = "blocks"
	// EdgeSuggests indicates the source suggests the target.
	EdgeSuggests EdgeType = "suggests"
)

// WorkflowProgress tracks completion progress.
type WorkflowProgress struct {
	// Completed is the number of completed nodes.
	Completed int `json:"completed" yaml:"completed"`

	// Total is the total number of nodes.
	Total int `json:"total" yaml:"total"`

	// Percent is the completion percentage (0-100).
	Percent float64 `json:"percent" yaml:"percent"`

	// PhaseProgress maps phase IDs to progress.
	PhaseProgress map[string]float64 `json:"phase_progress,omitempty" yaml:"phase_progress,omitempty"`
}

// NewWorkflow creates a new AIDLC workflow.
func NewWorkflow() *Workflow {
	return &Workflow{
		Name:        "AIDLC Workflow",
		Description: "AWS AI DLC Development Lifecycle workflow",
		Phases:      make([]WorkflowPhase, 0),
		Nodes:       make(map[string]*WorkflowNode),
		Edges:       make([]WorkflowEdge, 0),
	}
}

// DefaultWorkflow creates the default AIDLC workflow with all document types.
func DefaultWorkflow() *Workflow {
	w := NewWorkflow()

	// Create phases
	for _, phase := range AllPhases() {
		wp := WorkflowPhase{
			ID:      string(phase),
			Name:    phaseDisplayName(phase),
			Order:   phase.Order(),
			NodeIDs: make([]string, 0),
		}

		// Add nodes for this phase
		for _, docType := range DocTypesByPhase(phase) {
			nodeID := string(docType)
			node := &WorkflowNode{
				ID:        nodeID,
				DocType:   docType,
				Phase:     phase,
				Name:      docType.DisplayName(),
				Status:    NodePending,
				Required:  isRequiredDoc(docType),
				Automated: isAutomatedDoc(docType),
				DependsOn: make([]string, 0),
				Blocks:    make([]string, 0),
				Metadata:  make(map[string]any),
			}
			w.Nodes[nodeID] = node
			wp.NodeIDs = append(wp.NodeIDs, nodeID)
		}

		w.Phases = append(w.Phases, wp)
	}

	// Create default dependencies
	w.addDefaultDependencies()

	return w
}

// addDefaultDependencies adds the standard AIDLC dependencies.
func (w *Workflow) addDefaultDependencies() {
	// Inception phase dependencies
	w.AddEdge(string(DocVisionDocument), string(DocRequirementsSpec), EdgeDependency)
	w.AddEdge(string(DocRequirementsSpec), string(DocTechnicalSpec), EdgeDependency)
	w.AddEdge(string(DocRequirementsSpec), string(DocArchitectureSpec), EdgeDependency)

	// Inception -> Construction transition
	w.AddEdge(string(DocTechnicalSpec), string(DocImplementationPlan), EdgeDependency)
	w.AddEdge(string(DocArchitectureSpec), string(DocImplementationPlan), EdgeDependency)

	// Construction phase dependencies
	w.AddEdge(string(DocImplementationPlan), string(DocTestPlan), EdgeDependency)
	w.AddEdge(string(DocImplementationPlan), string(DocIntegrationPlan), EdgeDependency)
	w.AddEdge(string(DocTechnicalSpec), string(DocSecurityReview), EdgeDependency)

	// Construction -> Operations transition
	w.AddEdge(string(DocImplementationPlan), string(DocRunbook), EdgeDependency)
	w.AddEdge(string(DocTestPlan), string(DocMonitoringPlan), EdgeDependency)

	// Operations phase dependencies
	w.AddEdge(string(DocRunbook), string(DocDisasterPlan), EdgeDependency)
	w.AddEdge(string(DocMonitoringPlan), string(DocSLODocument), EdgeDependency)

	// Update node dependencies and blocks
	for _, edge := range w.Edges {
		if edge.Type == EdgeDependency {
			if toNode, ok := w.Nodes[edge.To]; ok {
				toNode.DependsOn = append(toNode.DependsOn, edge.From)
			}
			if fromNode, ok := w.Nodes[edge.From]; ok {
				fromNode.Blocks = append(fromNode.Blocks, edge.To)
			}
		}
	}
}

// AddEdge adds an edge to the workflow.
func (w *Workflow) AddEdge(from, to string, edgeType EdgeType) {
	w.Edges = append(w.Edges, WorkflowEdge{
		From: from,
		To:   to,
		Type: edgeType,
	})
}

// GetNode returns a node by ID.
func (w *Workflow) GetNode(id string) (*WorkflowNode, bool) {
	node, ok := w.Nodes[id]
	return node, ok
}

// GetPhase returns a phase by ID.
func (w *Workflow) GetPhase(id string) (*WorkflowPhase, bool) {
	for i := range w.Phases {
		if w.Phases[i].ID == id {
			return &w.Phases[i], true
		}
	}
	return nil, false
}

// Progress computes the current workflow progress.
func (w *Workflow) Progress() WorkflowProgress {
	progress := WorkflowProgress{
		Total:         len(w.Nodes),
		PhaseProgress: make(map[string]float64),
	}

	for _, node := range w.Nodes {
		if node.Status == NodeCompleted {
			progress.Completed++
		}
	}

	if progress.Total > 0 {
		progress.Percent = float64(progress.Completed) / float64(progress.Total) * 100
	}

	// Compute per-phase progress
	for _, phase := range w.Phases {
		var completed, total int
		for _, nodeID := range phase.NodeIDs {
			if node, ok := w.Nodes[nodeID]; ok {
				total++
				if node.Status == NodeCompleted {
					completed++
				}
			}
		}
		if total > 0 {
			progress.PhaseProgress[phase.ID] = float64(completed) / float64(total) * 100
		}
	}

	return progress
}

// ReadyNodes returns nodes that are ready to be worked on.
func (w *Workflow) ReadyNodes() []*WorkflowNode {
	var ready []*WorkflowNode
	for _, node := range w.Nodes {
		if node.Status == NodePending || node.Status == NodeReady {
			if w.canStart(node) {
				ready = append(ready, node)
			}
		}
	}
	return ready
}

// canStart checks if a node's dependencies are satisfied.
func (w *Workflow) canStart(node *WorkflowNode) bool {
	for _, depID := range node.DependsOn {
		if depNode, ok := w.Nodes[depID]; ok {
			if depNode.Status != NodeCompleted && depNode.Status != NodeSkipped {
				return false
			}
		}
	}
	return true
}

// UpdateNodeStatus updates a node's status and recalculates dependent statuses.
func (w *Workflow) UpdateNodeStatus(nodeID string, status NodeStatus, score *QualityScore) error {
	node, ok := w.Nodes[nodeID]
	if !ok {
		return fmt.Errorf("node not found: %s", nodeID)
	}

	node.Status = status
	node.Score = score

	// Update blocked nodes
	if status == NodeCompleted || status == NodeSkipped {
		for _, blockedID := range node.Blocks {
			if blockedNode, ok := w.Nodes[blockedID]; ok {
				if blockedNode.Status == NodeBlocked && w.canStart(blockedNode) {
					blockedNode.Status = NodeReady
				}
			}
		}
	} else if status == NodeFailed || status == NodeBlocked {
		for _, blockedID := range node.Blocks {
			if blockedNode, ok := w.Nodes[blockedID]; ok {
				if blockedNode.Status == NodePending || blockedNode.Status == NodeReady {
					blockedNode.Status = NodeBlocked
				}
			}
		}
	}

	return nil
}

// ToMermaid generates a Mermaid flowchart representation.
func (w *Workflow) ToMermaid() string {
	var sb strings.Builder
	sb.WriteString("flowchart TD\n")

	// Add subgraphs for phases
	for _, phase := range w.Phases {
		sb.WriteString(fmt.Sprintf("    subgraph %s[\"%s\"]\n", phase.ID, phase.Name))
		for _, nodeID := range phase.NodeIDs {
			if node, ok := w.Nodes[nodeID]; ok {
				shape := "([%s])"
				if node.Status == NodeCompleted {
					shape = "[[%s]]"
				} else if node.Status == NodeInProgress {
					shape = "((%s))"
				}
				sb.WriteString(fmt.Sprintf("        %s%s\n", nodeID, fmt.Sprintf(shape, node.Name)))
			}
		}
		sb.WriteString("    end\n")
	}

	// Add edges
	for _, edge := range w.Edges {
		arrow := "-->"
		if edge.Type == EdgeSuggests {
			arrow = "-.->|suggests|"
		}
		sb.WriteString(fmt.Sprintf("    %s %s %s\n", edge.From, arrow, edge.To))
	}

	// Add styling
	sb.WriteString("\n    classDef completed fill:#22c55e,stroke:#16a34a,color:#fff\n")
	sb.WriteString("    classDef inProgress fill:#3b82f6,stroke:#2563eb,color:#fff\n")
	sb.WriteString("    classDef blocked fill:#ef4444,stroke:#dc2626,color:#fff\n")
	sb.WriteString("    classDef pending fill:#6b7280,stroke:#4b5563,color:#fff\n")

	for nodeID, node := range w.Nodes {
		switch node.Status {
		case NodeCompleted:
			sb.WriteString(fmt.Sprintf("    class %s completed\n", nodeID))
		case NodeInProgress:
			sb.WriteString(fmt.Sprintf("    class %s inProgress\n", nodeID))
		case NodeBlocked, NodeFailed:
			sb.WriteString(fmt.Sprintf("    class %s blocked\n", nodeID))
		default:
			sb.WriteString(fmt.Sprintf("    class %s pending\n", nodeID))
		}
	}

	return sb.String()
}

// UpdateFromState updates the workflow from an AIDLC state.
func (w *Workflow) UpdateFromState(state *State) {
	// Update completed nodes
	for _, docType := range state.CompletedDocs {
		if node, ok := w.Nodes[string(docType)]; ok {
			node.Status = NodeCompleted
			if score, ok := state.DocumentScores[docType]; ok {
				node.Score = score
			}
		}
	}

	// Update in-progress nodes
	for _, docType := range state.InProgressDocs {
		if node, ok := w.Nodes[string(docType)]; ok {
			node.Status = NodeInProgress
		}
	}

	// Update pending nodes - mark as ready if dependencies are met
	for _, docType := range state.PendingDocs {
		if node, ok := w.Nodes[string(docType)]; ok {
			if w.canStart(node) {
				node.Status = NodeReady
			} else {
				node.Status = NodeBlocked
			}
		}
	}
}

// phaseDisplayName returns the display name for a phase.
func phaseDisplayName(phase Phase) string {
	switch phase {
	case PhaseInception:
		return "Inception"
	case PhaseConstruction:
		return "Construction"
	case PhaseOperations:
		return "Operations"
	default:
		return string(phase)
	}
}

// isRequiredDoc returns whether a document type is required.
func isRequiredDoc(docType DocType) bool {
	switch docType {
	case DocVisionDocument, DocRequirementsSpec, DocTechnicalSpec,
		DocImplementationPlan, DocTestPlan, DocRunbook:
		return true
	default:
		return false
	}
}

// isAutomatedDoc returns whether a document type is typically LLM-generated.
func isAutomatedDoc(docType DocType) bool {
	switch docType {
	case DocSecurityReview, DocMonitoringPlan, DocDisasterPlan:
		return true
	default:
		return false
	}
}
