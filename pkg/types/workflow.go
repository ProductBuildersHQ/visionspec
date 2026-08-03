package types

import (
	swf "github.com/ProductBuildersHQ/specification-workflow-spec/pkg/workflow"
)

// SpecConfigFromWorkflow converts a specification-workflow-spec Workflow's
// spec configuration into a visionspec SpecConfig.
// Returns nil if the workflow has no spec configuration.
func SpecConfigFromWorkflow(w *swf.Workflow) *SpecConfig {
	if w == nil || w.SpecConfig == nil {
		return nil
	}

	sc := &SpecConfig{
		Specs: make(map[string]*SpecRequirement, len(w.SpecConfig)),
	}
	for specType, req := range w.SpecConfig {
		if req == nil {
			continue
		}
		sc.Specs[specType] = &SpecRequirement{
			Required: req.Required,
			Category: SpecCategory(req.Category),
			Template: req.Template,
			Rubric:   req.Rubric,
		}
	}
	return sc
}
