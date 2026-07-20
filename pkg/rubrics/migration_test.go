package rubrics

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/plexusone/structured-evaluation/rubric"
	"gopkg.in/yaml.v3"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	// file = <root>/pkg/rubrics/migration_test.go
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

// TestAllRubricYAMLsAreStructuredEvaluationNative is the migration gate: every
// *.rubric.yaml shipped in the repo must parse directly as structured-evaluation
// native (the single canonical format) with at least one category. It guards
// against a legacy flat rubric silently reappearing after the flat parser was
// removed.
func TestAllRubricYAMLsAreStructuredEvaluationNative(t *testing.T) {
	root := repoRoot(t)
	var scanned int
	err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(p, ".rubric.yaml") {
			return nil
		}
		content, err := os.ReadFile(p) //nolint:gosec // G304/G122: reading repo-local rubric fixtures in a test
		if err != nil {
			return err
		}
		var rs rubric.RubricSet
		if err := yaml.Unmarshal(content, &rs); err != nil {
			t.Errorf("%s: does not parse as structured-evaluation native: %v", p, err)
			return nil
		}
		if len(rs.Categories) == 0 {
			t.Errorf("%s: no categories after structured-evaluation parse (likely still flat)", p)
		}
		scanned++
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if scanned < 80 {
		t.Errorf("scanned %d rubric files, expected at least 80", scanned)
	}
}
