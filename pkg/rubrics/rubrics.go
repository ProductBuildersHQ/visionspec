// Package rubrics provides evaluation rubrics for spec types. Rubrics use the
// shared structured-evaluation rubric definition (rubric.RubricSet) as the one
// canonical format across the ecosystem — there is no visionspec-local rubric
// type. Flat rubrics use categorical scales; rich rubrics use weighted criteria
// with indicators. Both are the same rubric.RubricSet.
package rubrics

import (
	"fmt"

	"github.com/ProductBuildersHQ/visionspec/pkg/types"
	"github.com/plexusone/structured-evaluation/rubric"
)

// DefaultPassCriteria returns the default pass criteria: no critical or high
// findings, unlimited medium/low, and all required categories must pass.
func DefaultPassCriteria() rubric.RubricPassCriteria {
	return rubric.RubricPassCriteria{
		MinCategoriesPassing: "all_required",
		MaxFindings:          &rubric.FindingLimits{Critical: 0, High: 0, Medium: -1, Low: -1},
	}
}

// StrictPassCriteria returns stricter pass criteria: no critical or high
// findings, at most 3 medium, and every category must pass.
func StrictPassCriteria() rubric.RubricPassCriteria {
	return rubric.RubricPassCriteria{
		MinCategoriesPassing: "all",
		MaxFindings:          &rubric.FindingLimits{Critical: 0, High: 0, Medium: 3, Low: -1},
	}
}

// registry maps spec types to their rubric sets.
var registry = make(map[types.SpecType]*rubric.RubricSet)

// Register adds a rubric set for a spec type to the registry.
func Register(specType types.SpecType, rs *rubric.RubricSet) {
	registry[specType] = rs
}

// Get returns the rubric set for a spec type.
func Get(specType types.SpecType) (*rubric.RubricSet, error) {
	rs, ok := registry[specType]
	if !ok {
		return nil, fmt.Errorf("rubric not found for spec type %q", specType)
	}
	return rs, nil
}

// MustGet returns the rubric set for a spec type, panicking on error.
func MustGet(specType types.SpecType) *rubric.RubricSet {
	rs, err := Get(specType)
	if err != nil {
		panic(err)
	}
	return rs
}

// Available returns all spec types with registered rubrics.
func Available() []types.SpecType {
	result := make([]types.SpecType, 0, len(registry))
	for st := range registry {
		result = append(result, st)
	}
	return result
}

// HasRubric returns true if a rubric exists for the spec type.
func HasRubric(specType types.SpecType) bool {
	_, ok := registry[specType]
	return ok
}
