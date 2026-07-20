package rubrics

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/plexusone/structured-evaluation/rubric"
)

// RubricYAML represents a rubric definition in visionspec's legacy flat YAML
// format (criteria as a pass/partial/fail map). It is retained only to read
// not-yet-migrated profile rubrics; parseRubricYAML converts them to the
// canonical structured-evaluation rubric.RubricSet. New rubrics are authored
// directly in the structured-evaluation format.
type RubricYAML struct {
	SpecType     string           `yaml:"spec_type"`
	Name         string           `yaml:"name"`
	Description  string           `yaml:"description"`
	Version      string           `yaml:"version"`
	Categories   []CategoryYAML   `yaml:"categories"`
	PassCriteria PassCriteriaYAML `yaml:"pass_criteria"`
}

// CategoryYAML represents a category in the legacy flat YAML format.
type CategoryYAML struct {
	ID          string       `yaml:"id"`
	Name        string       `yaml:"name"`
	Description string       `yaml:"description"`
	Weight      float64      `yaml:"weight"`
	Required    bool         `yaml:"required"`
	Criteria    CriteriaYAML `yaml:"criteria"`
}

// CriteriaYAML represents pass/partial/fail criteria in the legacy flat format.
type CriteriaYAML struct {
	Pass    string `yaml:"pass"`
	Partial string `yaml:"partial"`
	Fail    string `yaml:"fail"`
}

// PassCriteriaYAML represents pass criteria in the legacy flat format.
type PassCriteriaYAML struct {
	RequireAllPass bool `yaml:"require_all_pass"`
	MaxCritical    int  `yaml:"max_critical"`
	MaxHigh        int  `yaml:"max_high"`
	MaxMedium      int  `yaml:"max_medium"`
}

func sliceOrNil(s string) []string {
	if s == "" {
		return nil
	}
	return []string{s}
}

func minCategoriesPassing(requireAll bool) string {
	if requireAll {
		return "all"
	}
	return "all_required"
}

// ToRubricSet converts a legacy flat RubricYAML to the canonical
// structured-evaluation rubric.RubricSet (categorical scale per category).
func (r *RubricYAML) ToRubricSet() (*rubric.RubricSet, error) {
	if r.SpecType == "" {
		return nil, fmt.Errorf("spec_type is required")
	}
	if r.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if len(r.Categories) == 0 {
		return nil, fmt.Errorf("at least one category is required")
	}

	version := r.Version
	if version == "" {
		version = "1.0"
	}

	rs := rubric.NewRubricSet(r.SpecType+"-rubric", r.Name, version)
	rs.Description = r.Description
	rs.PassCriteria = rubric.RubricPassCriteria{
		MinCategoriesPassing: minCategoriesPassing(r.PassCriteria.RequireAllPass),
		MaxFindings: &rubric.FindingLimits{
			Critical: r.PassCriteria.MaxCritical,
			High:     r.PassCriteria.MaxHigh,
			Medium:   r.PassCriteria.MaxMedium,
			Low:      -1,
		},
	}

	for i, cat := range r.Categories {
		if cat.ID == "" {
			return nil, fmt.Errorf("category %d: id is required", i)
		}
		if cat.Name == "" {
			return nil, fmt.Errorf("category %d: name is required", i)
		}
		rs.AddCategory(*rubric.NewCategory(cat.ID, cat.Name, cat.Description).
			SetWeight(cat.Weight).SetRequired(cat.Required).
			WithPassPartialFail(
				sliceOrNil(cat.Criteria.Pass),
				sliceOrNil(cat.Criteria.Partial),
				sliceOrNil(cat.Criteria.Fail),
			))
	}

	return rs, nil
}

// parseRubricYAML parses rubric YAML content into a rubric.RubricSet. It first
// tries the canonical structured-evaluation format (rich weighted criteria or
// categorical scales); if that does not yield categories, it falls back to
// visionspec's legacy flat format. sourceName is used only for error context.
func parseRubricYAML(content []byte, sourceName string) (*rubric.RubricSet, error) {
	var rs rubric.RubricSet
	if err := yaml.Unmarshal(content, &rs); err == nil && len(rs.Categories) > 0 {
		return &rs, nil
	}

	var flat RubricYAML
	if err := yaml.Unmarshal(content, &flat); err != nil {
		return nil, fmt.Errorf("parsing rubric %s: %w", sourceName, err)
	}
	return flat.ToRubricSet()
}

// WriteRubricYAML writes a rubric.RubricSet to a YAML file in the canonical
// structured-evaluation format.
func WriteRubricYAML(path string, rs *rubric.RubricSet) error {
	data, err := yaml.Marshal(rs)
	if err != nil {
		return fmt.Errorf("marshaling rubric: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}
	return nil
}
