package rubrics

import (
	"fmt"
	"os"

	"github.com/plexusone/multispec/pkg/types"
	"gopkg.in/yaml.v3"
)

// RubricYAML represents a rubric definition in YAML format.
type RubricYAML struct {
	SpecType     string           `yaml:"spec_type"`
	Name         string           `yaml:"name"`
	Description  string           `yaml:"description"`
	Version      string           `yaml:"version"`
	Categories   []CategoryYAML   `yaml:"categories"`
	PassCriteria PassCriteriaYAML `yaml:"pass_criteria"`
}

// CategoryYAML represents a category in YAML format.
type CategoryYAML struct {
	ID          string       `yaml:"id"`
	Name        string       `yaml:"name"`
	Description string       `yaml:"description"`
	Weight      float64      `yaml:"weight"`
	Required    bool         `yaml:"required"`
	Criteria    CriteriaYAML `yaml:"criteria"`
}

// CriteriaYAML represents pass/partial/fail criteria in YAML format.
type CriteriaYAML struct {
	Pass    string `yaml:"pass"`
	Partial string `yaml:"partial"`
	Fail    string `yaml:"fail"`
}

// PassCriteriaYAML represents pass criteria in YAML format.
type PassCriteriaYAML struct {
	RequireAllPass bool `yaml:"require_all_pass"`
	MaxCritical    int  `yaml:"max_critical"`
	MaxHigh        int  `yaml:"max_high"`
	MaxMedium      int  `yaml:"max_medium"`
}

// ToRubricSet converts a RubricYAML to a RubricSet.
func (r *RubricYAML) ToRubricSet() (*RubricSet, error) {
	if r.SpecType == "" {
		return nil, fmt.Errorf("spec_type is required")
	}
	if r.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if len(r.Categories) == 0 {
		return nil, fmt.Errorf("at least one category is required")
	}

	rs := &RubricSet{
		SpecType:    types.SpecType(r.SpecType),
		Name:        r.Name,
		Description: r.Description,
		Categories:  make([]*Category, len(r.Categories)),
		PassCriteria: PassCriteria{
			RequireAllPass: r.PassCriteria.RequireAllPass,
			MaxCritical:    r.PassCriteria.MaxCritical,
			MaxHigh:        r.PassCriteria.MaxHigh,
			MaxMedium:      r.PassCriteria.MaxMedium,
		},
	}

	for i, cat := range r.Categories {
		if cat.ID == "" {
			return nil, fmt.Errorf("category %d: id is required", i)
		}
		if cat.Name == "" {
			return nil, fmt.Errorf("category %d: name is required", i)
		}

		rs.Categories[i] = &Category{
			ID:          cat.ID,
			Name:        cat.Name,
			Description: cat.Description,
			Weight:      cat.Weight,
			Required:    cat.Required,
			Criteria: CategoricalCriteria{
				Pass:    cat.Criteria.Pass,
				Partial: cat.Criteria.Partial,
				Fail:    cat.Criteria.Fail,
			},
		}
	}

	return rs, nil
}

// ToYAML converts a RubricSet to RubricYAML.
func (rs *RubricSet) ToYAML() *RubricYAML {
	yaml := &RubricYAML{
		SpecType:    string(rs.SpecType),
		Name:        rs.Name,
		Description: rs.Description,
		Categories:  make([]CategoryYAML, len(rs.Categories)),
		PassCriteria: PassCriteriaYAML{
			RequireAllPass: rs.PassCriteria.RequireAllPass,
			MaxCritical:    rs.PassCriteria.MaxCritical,
			MaxHigh:        rs.PassCriteria.MaxHigh,
			MaxMedium:      rs.PassCriteria.MaxMedium,
		},
	}

	for i, cat := range rs.Categories {
		yaml.Categories[i] = CategoryYAML{
			ID:          cat.ID,
			Name:        cat.Name,
			Description: cat.Description,
			Weight:      cat.Weight,
			Required:    cat.Required,
			Criteria: CriteriaYAML{
				Pass:    cat.Criteria.Pass,
				Partial: cat.Criteria.Partial,
				Fail:    cat.Criteria.Fail,
			},
		}
	}

	return yaml
}

// WriteRubricYAML writes a RubricSet to a YAML file.
func WriteRubricYAML(path string, rs *RubricSet) error {
	yamlData := rs.ToYAML()
	data, err := yaml.Marshal(yamlData)
	if err != nil {
		return fmt.Errorf("marshaling rubric: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}
