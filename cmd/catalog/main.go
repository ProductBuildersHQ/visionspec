// Command catalog emits the complete specification-workflow-spec catalog as
// JSON: spec types, resolved workflows, tool integrations, and pipelines.
//
// Consumers (the ProductBuildersHQ website, VisionStudio) use this export to
// generate pages and navigation deterministically from the library data.
//
//	go run ./cmd/catalog            # JSON to stdout
//	go run ./cmd/catalog -o out.json
//
// The output is deterministic: entries are sorted and no timestamps are
// included, so the same library version always produces the same bytes.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/ProductBuildersHQ/specification-workflow-spec/pkg/integration"
	"github.com/ProductBuildersHQ/specification-workflow-spec/pkg/integrations"
	"github.com/ProductBuildersHQ/specification-workflow-spec/pkg/loop"
	"github.com/ProductBuildersHQ/specification-workflow-spec/pkg/loops"
	"github.com/ProductBuildersHQ/specification-workflow-spec/pkg/pipeline"
	"github.com/ProductBuildersHQ/specification-workflow-spec/pkg/pipelines"
	"github.com/ProductBuildersHQ/specification-workflow-spec/pkg/spectype"
	"github.com/ProductBuildersHQ/specification-workflow-spec/pkg/workflow"
	"github.com/ProductBuildersHQ/specification-workflow-spec/pkg/workflows"
)

// Catalog is the exported shape.
type Catalog struct {
	// Version is the catalog format version.
	Version string `json:"version"`

	// SpecTypes is the full spec type registry.
	SpecTypes []spectype.SpecType `json:"spec_types"`

	// Workflows are the resolved definition-side workflows.
	Workflows []WorkflowEntry `json:"workflows"`

	// Integrations are the execution-side tool descriptors.
	Integrations []*integration.Integration `json:"integrations"`

	// Pipelines link definition workflows to execution integrations.
	Pipelines []*pipeline.Pipeline `json:"pipelines"`

	// LoopSystems are the modeled loop systems (e.g., the PBHQ two-loop
	// system and the AWS AI-DLC single-loop reference).
	LoopSystems []*loop.System `json:"loop_systems"`
}

// WorkflowEntry is a resolved workflow plus its template/rubric inventory.
type WorkflowEntry struct {
	*workflow.Workflow

	// TemplateSpecs lists spec type IDs that have templates in this workflow.
	TemplateSpecs []string `json:"template_specs,omitempty"`

	// RubricSpecs lists spec type IDs that have rubrics in this workflow.
	RubricSpecs []string `json:"rubric_specs,omitempty"`
}

func main() {
	out := flag.String("o", "", "output file (default stdout)")
	flag.Parse()

	catalog, err := buildCatalog()
	if err != nil {
		fmt.Fprintf(os.Stderr, "catalog: %v\n", err)
		os.Exit(1)
	}

	data, err := json.MarshalIndent(catalog, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "catalog: marshaling: %v\n", err)
		os.Exit(1)
	}
	data = append(data, '\n')

	if *out == "" {
		if _, err := os.Stdout.Write(data); err != nil {
			fmt.Fprintf(os.Stderr, "catalog: writing stdout: %v\n", err)
			os.Exit(1)
		}
		return
	}
	if err := os.WriteFile(*out, data, 0o600); err != nil {
		fmt.Fprintf(os.Stderr, "catalog: writing %s: %v\n", *out, err)
		os.Exit(1)
	}
}

func buildCatalog() (*Catalog, error) {
	loader := workflows.DefaultLoader()
	names := loader.Available()
	sort.Strings(names)

	entries := make([]WorkflowEntry, 0, len(names))
	for _, name := range names {
		loaded, err := loader.Load(name)
		if err != nil {
			return nil, fmt.Errorf("loading workflow %s: %w", name, err)
		}

		entries = append(entries, WorkflowEntry{
			Workflow:      loaded.Workflow,
			TemplateSpecs: sortedKeys(loaded.Templates),
			RubricSpecs:   sortedKeys(loaded.Rubrics),
		})
	}

	return &Catalog{
		Version:      "1",
		SpecTypes:    spectype.CoreSpecTypes(),
		Workflows:    entries,
		Integrations: integrations.All(),
		Pipelines:    pipelines.All(),
		LoopSystems:  loops.All(),
	}, nil
}

func sortedKeys[M ~map[string]V, V any](m M) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
