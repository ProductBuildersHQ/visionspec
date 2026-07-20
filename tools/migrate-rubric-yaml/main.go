// Command migrate-rubric-yaml is a one-time migration that rewrites visionspec's
// legacy flat rubric YAML files (criteria as a pass/partial/fail map) into the
// canonical structured-evaluation format (categorical Scale options). Files that
// already parse as structured-evaluation native (rich weighted criteria) are
// left untouched.
//
// It is self-contained (it does not import visionspec's rubrics package, which
// no longer defines the legacy types) so it remains buildable after the flat
// parser is removed. Every rewrite is verified by re-parsing the emitted YAML
// and comparing it to the intended RubricSet before the file is written.
//
// Usage: go run ./tools/migrate-rubric-yaml [root]   (root defaults to ".")
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	serubric "github.com/plexusone/structured-evaluation/rubric"
	"gopkg.in/yaml.v3"
)

// flat mirrors visionspec's legacy flat rubric YAML.
type flat struct {
	SpecType     string       `yaml:"spec_type"`
	Name         string       `yaml:"name"`
	Description  string       `yaml:"description"`
	Version      string       `yaml:"version"`
	Categories   []flatCat    `yaml:"categories"`
	PassCriteria flatPassCrit `yaml:"pass_criteria"`
}

type flatCat struct {
	ID          string       `yaml:"id"`
	Name        string       `yaml:"name"`
	Description string       `yaml:"description"`
	Weight      float64      `yaml:"weight"`
	Required    bool         `yaml:"required"`
	Criteria    flatCriteria `yaml:"criteria"`
}

type flatCriteria struct {
	Pass    string `yaml:"pass"`
	Partial string `yaml:"partial"`
	Fail    string `yaml:"fail"`
}

type flatPassCrit struct {
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

// flatToSE converts a parsed flat rubric to the canonical structured-evaluation
// RubricSet. specType is taken from the filename when the file omits spec_type.
func flatToSE(f flat, specType string) (*serubric.RubricSet, error) {
	if f.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if len(f.Categories) == 0 {
		return nil, fmt.Errorf("at least one category is required")
	}
	version := f.Version
	if version == "" {
		version = "1.0"
	}
	rs := serubric.NewRubricSet(specType+"-rubric", f.Name, version)
	rs.Description = f.Description
	rs.PassCriteria = serubric.RubricPassCriteria{
		MinCategoriesPassing: minCategoriesPassing(f.PassCriteria.RequireAllPass),
		MaxFindings: &serubric.FindingLimits{
			Critical: f.PassCriteria.MaxCritical,
			High:     f.PassCriteria.MaxHigh,
			Medium:   f.PassCriteria.MaxMedium,
			Low:      -1,
		},
	}
	for i, c := range f.Categories {
		if c.ID == "" {
			return nil, fmt.Errorf("category %d: id is required", i)
		}
		if c.Name == "" {
			return nil, fmt.Errorf("category %d: name is required", i)
		}
		rs.AddCategory(*serubric.NewCategory(c.ID, c.Name, c.Description).
			SetWeight(c.Weight).SetRequired(c.Required).
			WithPassPartialFail(
				sliceOrNil(c.Criteria.Pass),
				sliceOrNil(c.Criteria.Partial),
				sliceOrNil(c.Criteria.Fail),
			))
	}
	return rs, nil
}

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	var files []string
	//nolint:gosec // G703: one-time migration over repo-local rubric files
	err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(p, ".rubric.yaml") && !strings.Contains(p, "/.git/") {
			files = append(files, p)
		}
		return nil
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "walk:", err)
		os.Exit(1)
	}
	sort.Strings(files)

	migrated, skipped, failed := 0, 0, 0
	for _, f := range files {
		content, err := os.ReadFile(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "FAIL read %s: %v\n", f, err)
			failed++
			continue
		}

		// Skip files that already parse as structured-evaluation native.
		var probe serubric.RubricSet
		if err := yaml.Unmarshal(content, &probe); err == nil && len(probe.Categories) > 0 {
			skipped++
			continue
		}

		var ff flat
		if err := yaml.Unmarshal(content, &ff); err != nil {
			fmt.Fprintf(os.Stderr, "FAIL parse %s: %v\n", f, err)
			failed++
			continue
		}
		specType := strings.TrimSuffix(filepath.Base(f), ".rubric.yaml")
		if ff.SpecType != "" {
			specType = ff.SpecType
		}

		rs, err := flatToSE(ff, specType)
		if err != nil {
			fmt.Fprintf(os.Stderr, "FAIL convert %s: %v\n", f, err)
			failed++
			continue
		}

		out, err := yaml.Marshal(rs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "FAIL marshal %s: %v\n", f, err)
			failed++
			continue
		}

		// Round-trip verify: re-parse the emitted YAML and compare.
		var back serubric.RubricSet
		if err := yaml.Unmarshal(out, &back); err != nil {
			fmt.Fprintf(os.Stderr, "FAIL reparse %s: %v\n", f, err)
			failed++
			continue
		}
		if !reflect.DeepEqual(*rs, back) {
			fmt.Fprintf(os.Stderr, "FAIL roundtrip mismatch %s\n", f)
			failed++
			continue
		}

		if err := os.WriteFile(f, out, 0600); err != nil {
			fmt.Fprintf(os.Stderr, "FAIL write %s: %v\n", f, err)
			failed++
			continue
		}
		migrated++
	}

	fmt.Printf("migrated=%d skipped(se-native)=%d failed=%d total=%d\n", migrated, skipped, failed, len(files))
	if failed > 0 {
		os.Exit(1)
	}
}
