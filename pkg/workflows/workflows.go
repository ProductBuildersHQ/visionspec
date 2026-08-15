// Package workflows provides embedded default workflows with templates and rubrics.
//
// Workflows are loaded at init time from embedded filesystem and accessible
// via the registry without any filesystem I/O at runtime.
//
// Usage:
//
//	w, err := workflows.Get("aws-two-way-door")
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println(w.Workflow.Name)
//	fmt.Println(w.Templates["press"].Content)
//	fmt.Println(w.Rubrics["press"].Categories)
package workflows

import (
	"maps"

	"github.com/ProductBuildersHQ/visionspec/pkg/template"
	"github.com/ProductBuildersHQ/visionspec/pkg/workflow"
	"github.com/plexusone/structured-evaluation/rubric"
)

// LoadedWorkflow contains a fully loaded workflow with all associated templates and rubrics.
type LoadedWorkflow struct {
	// Workflow is the parsed workflow configuration.
	Workflow *workflow.Workflow

	// Templates maps spec type ID to its template (e.g., "press" -> press template).
	Templates map[string]*template.Template

	// Rubrics maps spec type ID to its rubric definition (e.g., "press" -> press rubric).
	// Rubrics use structured-evaluation's canonical RubricSet type.
	Rubrics map[string]*rubric.RubricSet
}

// registry holds all loaded workflows, keyed by workflow name.
var registry = make(map[string]*LoadedWorkflow)

// Get returns a loaded workflow by name.
// Returns an error if the workflow is not found.
func Get(name string) (*LoadedWorkflow, error) {
	if w, ok := registry[name]; ok {
		return w, nil
	}
	return nil, &NotFoundError{Name: name}
}

// MustGet returns a loaded workflow by name, panicking if not found.
func MustGet(name string) *LoadedWorkflow {
	w, err := Get(name)
	if err != nil {
		panic(err)
	}
	return w
}

// List returns the names of all available workflows.
func List() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}

// All returns all loaded workflows.
func All() map[string]*LoadedWorkflow {
	result := make(map[string]*LoadedWorkflow, len(registry))
	maps.Copy(result, registry)
	return result
}

// Register adds a workflow to the registry.
// This is typically called from init() in embed.go.
func Register(name string, w *LoadedWorkflow) {
	registry[name] = w
}

// NotFoundError is returned when a workflow is not found.
type NotFoundError struct {
	Name string
}

func (e *NotFoundError) Error() string {
	return "workflow not found: " + e.Name
}
