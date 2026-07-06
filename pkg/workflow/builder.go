package workflow

// Builder provides a fluent API for constructing workflows.
type Builder struct {
	workflow     *Workflow
	currentPhase string
	err          error
}

// NewBuilder creates a new workflow builder.
func NewBuilder(name string) *Builder {
	return &Builder{
		workflow: New(name),
	}
}

// Description sets the workflow description.
func (b *Builder) Description(desc string) *Builder {
	if b.err != nil {
		return b
	}
	b.workflow.Description = desc
	return b
}

// Phase adds a phase and sets it as current for subsequent nodes.
func (b *Builder) Phase(id, name string, order int) *Builder {
	if b.err != nil {
		return b
	}
	b.workflow.AddPhase(id, name, order)
	b.currentPhase = id
	return b
}

// PhaseWithDescription adds a phase with description.
func (b *Builder) PhaseWithDescription(id, name, desc string, order int) *Builder {
	if b.err != nil {
		return b
	}
	phase := b.workflow.AddPhase(id, name, order)
	phase.Description = desc
	b.currentPhase = id
	return b
}

// Node adds a node to the current phase.
func (b *Builder) Node(id, name string) *NodeBuilder {
	return &NodeBuilder{
		builder: b,
		node: &Node{
			ID:     id,
			Name:   name,
			Phase:  b.currentPhase,
			Status: StatusPending,
		},
	}
}

// Build finalizes and validates the workflow.
func (b *Builder) Build() (*Workflow, error) {
	if b.err != nil {
		return nil, b.err
	}
	if err := b.workflow.Validate(); err != nil {
		return nil, err
	}
	return b.workflow, nil
}

// MustBuild finalizes the workflow and panics on error.
func (b *Builder) MustBuild() *Workflow {
	w, err := b.Build()
	if err != nil {
		panic(err)
	}
	return w
}

// NodeBuilder provides a fluent API for constructing nodes.
type NodeBuilder struct {
	builder *Builder
	node    *Node
}

// Description sets the node description.
func (nb *NodeBuilder) Description(desc string) *NodeBuilder {
	nb.node.Description = desc
	return nb
}

// Type sets the node type.
func (nb *NodeBuilder) Type(t NodeType) *NodeBuilder {
	nb.node.Type = t
	return nb
}

// DependsOn sets the node dependencies.
func (nb *NodeBuilder) DependsOn(deps ...string) *NodeBuilder {
	nb.node.DependsOn = deps
	return nb
}

// Automated marks the node as machine-generated.
func (nb *NodeBuilder) Automated() *NodeBuilder {
	nb.node.Automated = true
	return nb
}

// Status sets the initial status.
func (nb *NodeBuilder) Status(s Status) *NodeBuilder {
	nb.node.Status = s
	return nb
}

// Metadata sets a metadata key-value pair.
func (nb *NodeBuilder) Metadata(key string, value any) *NodeBuilder {
	if nb.node.Metadata == nil {
		nb.node.Metadata = make(map[string]any)
	}
	nb.node.Metadata[key] = value
	return nb
}

// Add finalizes the node and adds it to the workflow.
func (nb *NodeBuilder) Add() *Builder {
	if nb.builder.err != nil {
		return nb.builder
	}
	if err := nb.builder.workflow.AddNode(nb.node); err != nil {
		nb.builder.err = err
	}
	return nb.builder
}
