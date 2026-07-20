package rubrics

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/plexusone/structured-evaluation/rubric"
)

// parseRubricYAML parses canonical structured-evaluation rubric YAML into a
// rubric.RubricSet. It is the single rubric format across the ecosystem: flat
// rubrics use categorical Scale options, rich rubrics use weighted Criteria.
// sourceName is used only for error context.
func parseRubricYAML(content []byte, sourceName string) (*rubric.RubricSet, error) {
	var rs rubric.RubricSet
	if err := yaml.Unmarshal(content, &rs); err != nil {
		return nil, fmt.Errorf("parsing rubric %s: %w", sourceName, err)
	}
	if len(rs.Categories) == 0 {
		return nil, fmt.Errorf("rubric %s: at least one category is required", sourceName)
	}
	return &rs, nil
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
