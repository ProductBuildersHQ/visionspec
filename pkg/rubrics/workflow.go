package rubrics

import (
	"fmt"
	"sort"

	sws "github.com/ProductBuildersHQ/visionspec/pkg/workflows"
	prrubrics "github.com/grokify/prism-roadmap/rubrics"
	"github.com/plexusone/structured-evaluation/rubric"

	"github.com/ProductBuildersHQ/visionspec/pkg/types"
)

// mapLoader serves rubrics from an in-memory map keyed by spec type,
// as loaded from a specification-workflow-spec workflow.
type mapLoader struct {
	rubrics map[string]*rubric.RubricSet
}

// NewMapLoader creates a loader backed by a map of spec type to rubric set,
// as provided by specification-workflow-spec's LoadedWorkflow.Rubrics.
func NewMapLoader(rubrics map[string]*rubric.RubricSet) Loader {
	return &mapLoader{rubrics: rubrics}
}

func (l *mapLoader) Load(specType types.SpecType) (*rubric.RubricSet, error) {
	rs, ok := l.rubrics[string(specType)]
	if !ok {
		return nil, fmt.Errorf("rubric not found for spec type %q", specType)
	}
	return rs, nil
}

func (l *mapLoader) Available() []types.SpecType {
	result := make([]types.SpecType, 0, len(l.rubrics))
	for name := range l.rubrics {
		result = append(result, types.SpecType(name))
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return result
}

// canvasLoader resolves canonical canvas rubrics from prism-roadmap,
// the single source of truth for them.
func canvasLoader() Loader {
	return NewSubFSLoader(prrubrics.FS())
}

// LoaderForWorkflow returns a rubric loader for a loaded workflow.
// The workflow's own rubrics take priority (falling back to the default
// embedded rubrics when the workflow has none), with prism-roadmap's
// canonical canvas rubrics chained last.
func LoaderForWorkflow(w *sws.LoadedWorkflow) Loader {
	var base Loader
	if w != nil && len(w.Rubrics) > 0 {
		base = NewMapLoader(w.Rubrics)
	} else {
		base = DefaultLoader()
	}
	return NewChainLoader(base, canvasLoader())
}
