// Package synthesis defines rules for generating specs from other specs.
//
// Synthesis rules define the dependency graph between spec types, allowing
// automated generation of downstream documents from upstream sources.
package synthesis

// Rule defines how a spec can be synthesized from source specs.
type Rule struct {
	// Sources are the spec type IDs required to synthesize this spec.
	Sources []string `json:"sources" jsonschema:"required,description=Source spec type IDs"`

	// Guidance is the prompt context for LLM synthesis.
	Guidance string `json:"guidance,omitempty" jsonschema:"description=LLM prompt guidance for synthesis"`

	// PromptContext is additional context for the synthesis prompt.
	PromptContext string `json:"prompt_context,omitempty" jsonschema:"description=Additional synthesis prompt context"`

	// Required indicates all sources must be present (vs. best-effort).
	Required bool `json:"required,omitempty" jsonschema:"description=Whether all sources are required"`

	// Priority determines synthesis order when multiple rules exist.
	Priority int `json:"priority,omitempty" jsonschema:"description=Synthesis priority (higher = earlier)"`
}

// DAG represents the synthesis dependency graph.
type DAG struct {
	// Nodes are the spec type IDs in topological order.
	Nodes []string `json:"nodes" jsonschema:"required,description=Spec types in topological order"`

	// Edges are source->target dependencies.
	Edges []Edge `json:"edges" jsonschema:"required,description=Dependency edges"`
}

// Edge represents a synthesis dependency.
type Edge struct {
	// Source is the upstream spec type ID.
	Source string `json:"source" jsonschema:"required,description=Upstream spec type"`

	// Target is the downstream spec type ID (synthesized from source).
	Target string `json:"target" jsonschema:"required,description=Downstream spec type"`
}

// BuildDAG constructs a synthesis DAG from rules.
func BuildDAG(rules map[string]*Rule) *DAG {
	dag := &DAG{
		Nodes: make([]string, 0),
		Edges: make([]Edge, 0),
	}

	seen := make(map[string]bool)

	for target, rule := range rules {
		if !seen[target] {
			dag.Nodes = append(dag.Nodes, target)
			seen[target] = true
		}

		for _, source := range rule.Sources {
			if !seen[source] {
				dag.Nodes = append(dag.Nodes, source)
				seen[source] = true
			}

			dag.Edges = append(dag.Edges, Edge{
				Source: source,
				Target: target,
			})
		}
	}

	return dag
}
