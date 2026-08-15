//go:build ignore

// This program generates JSON Schema files from Go types.
// Run with: go generate ./schema/...
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/invopop/jsonschema"

	"github.com/ProductBuildersHQ/visionspec/pkg/gate"
	"github.com/ProductBuildersHQ/visionspec/pkg/integration"
	"github.com/ProductBuildersHQ/visionspec/pkg/layout"
	"github.com/ProductBuildersHQ/visionspec/pkg/loop"
	"github.com/ProductBuildersHQ/visionspec/pkg/pipeline"
	"github.com/ProductBuildersHQ/visionspec/pkg/spectype"
	"github.com/ProductBuildersHQ/visionspec/pkg/synthesis"
	"github.com/ProductBuildersHQ/visionspec/pkg/template"
	"github.com/ProductBuildersHQ/visionspec/pkg/workflow"
)

func main() {
	schemas := []struct {
		name string
		typ  any
	}{
		{"spectype", spectype.SpecTypeRegistry{}},
		{"workflow", workflow.Workflow{}},
		{"template", template.Template{}},
		{"synthesis", synthesis.DAG{}},
		{"gate", gate.Gate{}},
		{"layout", layout.Layout{}},
		{"integration", integration.Integration{}},
		{"pipeline", pipeline.Pipeline{}},
		{"loop", loop.System{}},
	}

	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	for _, s := range schemas {
		if err := generateSchema(dir, s.name, s.typ); err != nil {
			fmt.Fprintf(os.Stderr, "error generating %s: %v\n", s.name, err)
			os.Exit(1)
		}
		fmt.Printf("generated %s.schema.json\n", s.name)
	}
}

func generateSchema(dir, name string, v any) error {
	// Use references for recursive types (like Section.Subsections)
	r := jsonschema.Reflector{}

	schema := r.Reflect(v)
	schema.ID = jsonschema.ID(fmt.Sprintf(
		"https://github.com/ProductBuildersHQ/specification-workflow-spec/schema/%s.schema.json",
		name,
	))
	schema.Title = toTitle(name)
	schema.Description = fmt.Sprintf("Schema for %s definitions", name)

	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling schema: %w", err)
	}

	path := filepath.Join(dir, name+".schema.json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}

func toTitle(name string) string {
	titles := map[string]string{
		"spectype":    "Specification Type Registry",
		"workflow":    "Specification Workflow",
		"template":    "Spec Template",
		"synthesis":   "Synthesis DAG",
		"gate":        "Phase Gate",
		"layout":      "Project Layout",
		"integration": "External Tool Integration",
		"pipeline":    "Definition-to-Execution Pipeline",
		"loop":        "Loop System",
	}
	if t, ok := titles[name]; ok {
		return t
	}
	return name
}

// Ensure types are used to avoid import errors
var _ = reflect.TypeOf(spectype.SpecType{})
