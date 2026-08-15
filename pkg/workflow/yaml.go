package workflow

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Principles and artifacts support a YAML shorthand form in addition to the
// canonical struct form:
//
//	principles:
//	  - customer_obsession: "Start with the customer and work backwards"
//
// is equivalent to:
//
//	principles:
//	  - id: customer_obsession
//	    name: Customer Obsession
//	    description: "Start with the customer and work backwards"
//
// The shorthand is a single-key mapping whose key is not a canonical field name.

// principleFields are the canonical field names of Principle.
var principleFields = map[string]bool{
	"id": true, "name": true, "description": true, "source": true,
}

// UnmarshalYAML supports both canonical and shorthand principle forms.
func (p *Principle) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("principle must be a mapping, got %v", node.Kind)
	}

	// Shorthand: single key-value pair with a non-canonical key.
	if len(node.Content) == 2 && !principleFields[node.Content[0].Value] {
		p.ID = node.Content[0].Value
		p.Name = titleFromID(p.ID)
		return node.Content[1].Decode(&p.Description)
	}

	type plain Principle
	var pp plain
	if err := node.Decode(&pp); err != nil {
		return err
	}
	*p = Principle(pp)
	return nil
}

// artifactFields are the canonical field names of Artifact.
var artifactFields = map[string]bool{
	"id": true, "description": true, "category": true,
}

// UnmarshalYAML supports both a flat artifact sequence and a mapping of
// category name to artifact sequence (e.g., primary/supporting groups).
func (as *Artifacts) UnmarshalYAML(node *yaml.Node) error {
	switch node.Kind {
	case yaml.SequenceNode:
		var items []Artifact
		if err := node.Decode(&items); err != nil {
			return err
		}
		*as = items
		return nil
	case yaml.MappingNode:
		var result []Artifact
		for i := 0; i+1 < len(node.Content); i += 2 {
			category := node.Content[i].Value
			var items []Artifact
			if err := node.Content[i+1].Decode(&items); err != nil {
				return fmt.Errorf("artifacts group %q: %w", category, err)
			}
			for j := range items {
				items[j].Category = category
			}
			result = append(result, items...)
		}
		*as = result
		return nil
	default:
		return fmt.Errorf("artifacts must be a sequence or mapping, got %v", node.Kind)
	}
}

// UnmarshalYAML supports canonical, shorthand, and bare-string artifact forms.
func (a *Artifact) UnmarshalYAML(node *yaml.Node) error {
	// Bare string: "- press_release"
	if node.Kind == yaml.ScalarNode {
		a.ID = node.Value
		return nil
	}

	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("artifact must be a string or mapping, got %v", node.Kind)
	}

	// Shorthand: single key-value pair with a non-canonical key.
	if len(node.Content) == 2 && !artifactFields[node.Content[0].Value] {
		a.ID = node.Content[0].Value
		return node.Content[1].Decode(&a.Description)
	}

	type plain Artifact
	var aa plain
	if err := node.Decode(&aa); err != nil {
		return err
	}
	*a = Artifact(aa)
	return nil
}

// UnmarshalYAML supports both the object form and a bare-string shorthand for a
// spec source:
//
//	template: {from: enterprise}
//	rubric: local
//
// The scalar form sets From directly.
func (s *SpecSource) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.ScalarNode {
		s.From = node.Value
		return nil
	}
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("spec source must be a string or mapping, got %v", node.Kind)
	}
	type plain SpecSource
	var ps plain
	if err := node.Decode(&ps); err != nil {
		return err
	}
	*s = SpecSource(ps)
	return nil
}

// titleFromID converts an identifier like "customer_obsession" to "Customer Obsession".
func titleFromID(id string) string {
	words := strings.FieldsFunc(id, func(r rune) bool {
		return r == '_' || r == '-'
	})
	for i, w := range words {
		if w != "" {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}
