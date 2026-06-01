// Package rubrics provides evaluation rubrics for spec types.
package rubrics

import (
	"fmt"

	"github.com/ProductBuildersHQ/visionspec/pkg/types"
	"github.com/plexusone/structured-evaluation/rubric"
)

// RubricSet contains all criteria for evaluating a spec type.
type RubricSet struct {
	SpecType    types.SpecType
	Name        string
	Description string
	Categories  []*Category
	PassCriteria
}

// Category represents a scoring category with categorical criteria.
type Category struct {
	ID          string
	Name        string
	Description string
	Weight      float64
	Required    bool
	Criteria    CategoricalCriteria
}

// CategoricalCriteria defines pass/partial/fail criteria for a category.
type CategoricalCriteria struct {
	Pass    string
	Partial string
	Fail    string
}

// PassCriteria defines what constitutes a passing evaluation.
type PassCriteria struct {
	RequireAllPass bool // All categories must pass
	MaxCritical    int
	MaxHigh        int
	MaxMedium      int // -1 = unlimited
}

// DefaultPassCriteria returns the default pass criteria.
func DefaultPassCriteria() PassCriteria {
	return PassCriteria{
		RequireAllPass: false,
		MaxCritical:    0,
		MaxHigh:        0,
		MaxMedium:      -1, // unlimited
	}
}

// StrictPassCriteria returns stricter pass criteria.
func StrictPassCriteria() PassCriteria {
	return PassCriteria{
		RequireAllPass: true,
		MaxCritical:    0,
		MaxHigh:        0,
		MaxMedium:      3,
	}
}

// registry maps spec types to their rubric sets.
var registry = make(map[types.SpecType]*RubricSet)

// Register adds a rubric set to the registry.
func Register(rs *RubricSet) {
	registry[rs.SpecType] = rs
}

// Get returns the rubric set for a spec type.
func Get(specType types.SpecType) (*RubricSet, error) {
	rs, ok := registry[specType]
	if !ok {
		return nil, fmt.Errorf("rubric not found for spec type %q", specType)
	}
	return rs, nil
}

// MustGet returns the rubric set for a spec type, panicking on error.
func MustGet(specType types.SpecType) *RubricSet {
	rs, err := Get(specType)
	if err != nil {
		panic(err)
	}
	return rs
}

// Available returns all spec types with registered rubrics.
func Available() []types.SpecType {
	var result []types.SpecType
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

// TotalWeight returns the sum of all category weights.
func (rs *RubricSet) TotalWeight() float64 {
	var total float64
	for _, cat := range rs.Categories {
		total += cat.Weight
	}
	return total
}

// CategoryByID returns a category by its ID.
func (rs *RubricSet) CategoryByID(id string) (*Category, bool) {
	for _, cat := range rs.Categories {
		if cat.ID == id {
			return cat, true
		}
	}
	return nil, false
}

// ToEvaluationRubricSet converts to structured-evaluation RubricSet.
func (rs *RubricSet) ToEvaluationRubricSet() *rubric.RubricSet {
	evalRS := rubric.NewRubricSet(
		string(rs.SpecType)+"-rubric",
		rs.Name,
		"1.0",
	)
	evalRS.Description = rs.Description

	for _, cat := range rs.Categories {
		evalCat := rubric.NewCategory(cat.ID, cat.Name, cat.Description).
			SetWeight(cat.Weight).
			SetRequired(cat.Required).
			WithPassPartialFail(
				[]string{cat.Criteria.Pass},
				[]string{cat.Criteria.Partial},
				[]string{cat.Criteria.Fail},
			)
		evalRS.AddCategory(*evalCat)
	}

	return evalRS
}
