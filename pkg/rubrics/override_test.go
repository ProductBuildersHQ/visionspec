package rubrics

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ProductBuildersHQ/visionspec/pkg/types"
)

func TestOverrideLoader(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "custom-prd.yaml")
	content := `id: prd-rubric
name: Override PRD
categories:
  - id: c1
    name: C1
    scale:
      type: categorical
      options:
        - {value: pass, criteria: ["ok"]}
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	loader := NewOverrideLoader(map[types.SpecType]string{types.SpecTypePRD: path})

	rs, err := loader.Load(types.SpecTypePRD)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if rs.Name != "Override PRD" {
		t.Errorf("Name = %q, want Override PRD", rs.Name)
	}
	if _, err := loader.Load(types.SpecTypeMRD); err == nil {
		t.Error("expected error for unmapped spec type")
	}
}
