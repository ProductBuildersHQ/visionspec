package pipelines

import (
	"testing"

	"github.com/ProductBuildersHQ/specification-workflow-spec/pkg/integrations"
	"github.com/ProductBuildersHQ/specification-workflow-spec/pkg/loops"
	"github.com/ProductBuildersHQ/specification-workflow-spec/pkg/spectype"
	"github.com/ProductBuildersHQ/specification-workflow-spec/pkg/workflows"
)

// TestCatalogCrossReferencesResolve guards the whole embedded catalog graph.
//
// Each package's own Validate() only checks internal integrity; nothing checks
// references that cross package boundaries — a pipeline that names a workflow or
// spec type, an integration artifact mapped to a spec type ID, a loop station
// mapped to a workflow/integration/spec type. A typo in any of those parses and
// validates cleanly, then silently breaks downstream consumers (the website,
// VisionStudio) with no build signal. This test is that missing signal.
func TestCatalogCrossReferencesResolve(t *testing.T) {
	specIDs := knownSpecIDs()
	workflowNames := knownWorkflowNames()
	integrationIDs := knownIntegrationIDs()

	checkSpec := func(id, context string) {
		if id != "" && !specIDs[id] {
			t.Errorf("%s: spec type %q not found in registry", context, id)
		}
	}
	checkWorkflow := func(name, context string) {
		if name != "" && !workflowNames[name] {
			t.Errorf("%s: workflow %q not found", context, name)
		}
	}
	checkIntegration := func(id, context string) {
		if id != "" && !integrationIDs[id] {
			t.Errorf("%s: integration %q not found", context, id)
		}
	}

	// Pipelines -> workflows, integrations, spec types.
	for _, p := range All() {
		checkWorkflow(p.Definition.Workflow, "pipeline "+p.ID+" definition")
		checkIntegration(p.Execution.Integration, "pipeline "+p.ID+" execution")
		for _, out := range p.Definition.Outputs {
			checkSpec(out, "pipeline "+p.ID+" output")
		}
		for _, h := range p.Handoffs {
			checkSpec(h.SpecType, "pipeline "+p.ID+" handoff")
		}
	}

	// Integrations -> spec types (artifact mappings and input ports).
	for _, in := range integrations.All() {
		for _, a := range in.Artifacts {
			checkSpec(a.SpecType, "integration "+in.ID+" artifact "+a.ID)
		}
		for _, port := range in.Inputs {
			for _, st := range port.SpecTypes {
				checkSpec(st, "integration "+in.ID+" input port "+port.ID)
			}
		}
	}

	// Loop stations -> spec types, workflows, integrations.
	for _, sys := range loops.All() {
		for _, l := range sys.Loops {
			for _, st := range l.Stations {
				ctx := "loop " + sys.ID + " station " + l.ID + "." + st.ID
				for _, s := range st.SpecTypes {
					checkSpec(s, ctx)
				}
				for _, w := range st.Workflows {
					checkWorkflow(w, ctx)
				}
				for _, i := range st.Integrations {
					checkIntegration(i, ctx)
				}
			}
		}
	}
}

// knownSpecIDs returns the set of valid spec type IDs, including aliases, since
// cross-references may use either the canonical ID or a documented alias.
func knownSpecIDs() map[string]bool {
	ids := map[string]bool{}
	for _, s := range spectype.CoreSpecTypes() {
		ids[s.ID] = true
		for _, alias := range s.Aliases {
			ids[alias] = true
		}
	}
	return ids
}

func knownWorkflowNames() map[string]bool {
	names := map[string]bool{}
	for _, n := range workflows.DefaultLoader().Available() {
		names[n] = true
	}
	return names
}

func knownIntegrationIDs() map[string]bool {
	ids := map[string]bool{}
	for _, id := range integrations.List() {
		ids[id] = true
	}
	return ids
}
