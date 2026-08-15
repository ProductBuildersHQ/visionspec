package templates

import (
	"fmt"
	"sort"

	spectemplate "github.com/ProductBuildersHQ/visionspec/pkg/template"
	sws "github.com/ProductBuildersHQ/visionspec/pkg/workflows"
	prtemplates "github.com/grokify/prism-roadmap/templates"

	"github.com/ProductBuildersHQ/visionspec/pkg/types"
)

// mapLoader serves templates from an in-memory map keyed by spec type,
// as loaded from a specification-workflow-spec workflow.
type mapLoader struct {
	templates map[string]*spectemplate.Template
}

// NewMapLoader creates a loader backed by a map of spec type to template,
// as provided by specification-workflow-spec's LoadedWorkflow.Templates.
func NewMapLoader(templates map[string]*spectemplate.Template) Loader {
	return &mapLoader{templates: templates}
}

func (l *mapLoader) Load(specType types.SpecType) (*Template, error) {
	tmpl, ok := l.templates[string(specType)]
	if !ok {
		return nil, fmt.Errorf("template not found for spec type %q", specType)
	}
	return &Template{
		SpecType: specType,
		Content:  tmpl.Content,
	}, nil
}

func (l *mapLoader) Available() []types.SpecType {
	result := make([]types.SpecType, 0, len(l.templates))
	for name := range l.templates {
		result = append(result, types.SpecType(name))
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return result
}

// canvasLoader resolves canonical canvas templates (BMC, OpportunitySpec, and
// other canvases) from prism-roadmap, the single source of truth for them.
func canvasLoader() Loader {
	return NewSubFSLoader(prtemplates.FS())
}

// LoaderForWorkflow returns a template loader for a loaded workflow.
// The workflow's own templates take priority (falling back to the default
// embedded templates when the workflow has none), with prism-roadmap's
// canonical canvas templates chained last.
func LoaderForWorkflow(w *sws.LoadedWorkflow) Loader {
	var base Loader
	if w != nil && len(w.Templates) > 0 {
		base = NewMapLoader(w.Templates)
	} else {
		base = DefaultLoader()
	}
	return NewChainLoader(base, canvasLoader())
}
