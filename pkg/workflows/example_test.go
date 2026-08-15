package workflows_test

import (
	"fmt"

	"github.com/ProductBuildersHQ/specification-workflow-spec/pkg/workflows"
)

func Example() {
	// Load a workflow by name - no filesystem access needed
	w, err := workflows.Get("aws-two-way-door")
	if err != nil {
		panic(err)
	}

	// Access workflow metadata
	fmt.Printf("Workflow: %s\n", w.Workflow.Name)
	fmt.Printf("Extends: %s\n", w.Workflow.Extends)
	fmt.Printf("Methodology: %s\n", w.Workflow.Methodology.Name)

	// Access template content directly
	press := w.Templates["press"]
	fmt.Printf("Press template length: %d chars\n", len(press.Content))

	// Access rubric categories
	pressRubric := w.Rubrics["press"]
	fmt.Printf("Press rubric categories: %d\n", len(pressRubric.Categories))

	// Check required specs
	required := w.Workflow.RequiredSpecs()
	fmt.Printf("Required specs count: %d\n", len(required))

	// Output:
	// Workflow: aws-two-way-door
	// Extends: enterprise
	// Methodology: Amazon Working Backwards (Feature)
	// Press template length: 6338 chars
	// Press rubric categories: 9
	// Required specs count: 2
}

func ExampleList() {
	// List all available workflows
	names := workflows.List()
	fmt.Printf("Available workflows: %d\n", len(names))

	// Output:
	// Available workflows: 25
}
