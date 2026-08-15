// Package specworkflow adapts the generic workflow package for visionspec.
//
// It generates workflows from specification workflows, mapping spec types to
// workflow nodes and synthesis dependencies to edges.
package specworkflow

import (
	sws "github.com/ProductBuildersHQ/specification-workflow-spec/pkg/workflows"
	"github.com/ProductBuildersHQ/specification-workflow-spec/pkg/workflow"
	"github.com/ProductBuildersHQ/specification-workflow-spec/pkg/types"
	"github.com/ProductBuildersHQ/specification-workflow-spec/pkg/execgraph"
)

// Node types for visionspec
const (
	TypeSource    execgraph.NodeType = "source"
	TypeGTM       execgraph.NodeType = "gtm"
	TypeTechnical execgraph.NodeType = "technical"
	TypeOutput    execgraph.NodeType = "output"
)

// Spec metadata for display
var specMeta = map[string]struct {
	Name        string
	Description string
}{
	"mrd":              {Name: "Market Requirements", Description: "Define the market opportunity and business case"},
	"opportunity-spec": {Name: "Opportunity Spec", Description: "Define the opportunity using Patton + Cagan framework"},
	"press":            {Name: "Press Release", Description: "Working backwards: announce the finished product"},
	"faq":              {Name: "FAQ", Description: "Challenge claims and surface gaps"},
	"narrative-1p":     {Name: "1-Pager", Description: "Executive summary narrative"},
	"narrative-6p":     {Name: "6-Pager", Description: "Amazon-style detailed narrative"},
	"prd":              {Name: "Product Requirements", Description: "Define what to build"},
	"uxd":              {Name: "User Experience", Description: "Design the user experience"},
	"trd":              {Name: "Technical Design", Description: "Define how to build it"},
	"tpd":              {Name: "Test Plan", Description: "Define how to test it"},
	"ird":              {Name: "Infrastructure", Description: "Define infrastructure requirements"},
	"spec":             {Name: "Execution Spec", Description: "Reconciled spec for execution"},
	"current-truth":    {Name: "Current Truth", Description: "Post-ship state documentation"},
	"hypothesis":       {Name: "Hypothesis", Description: "Define and test hypotheses"},
	"bmc":              {Name: "Business Model Canvas", Description: "Business model analysis"},
}

// Phase definitions
var phaseDefs = []struct {
	ID          string
	Name        string
	Description string
	Order       int
	Categories  []types.SpecCategory
}{
	{"discovery", "Discovery", "Define the market and opportunity", 1, []types.SpecCategory{types.CategorySource}},
	{"vision", "Vision", "Working backwards from the customer announcement", 2, []types.SpecCategory{types.CategoryGTM}},
	{"product", "Product", "Define what to build and the experience", 3, []types.SpecCategory{types.CategorySource}},
	{"technical", "Technical", "Define how to build and test", 4, []types.SpecCategory{types.CategoryTechnical}},
	{"reconcile", "Reconciliation", "Generate unified execution spec", 5, []types.SpecCategory{types.CategoryOutput}},
}

// FromWorkflow generates a visionspec workflow from a loaded specification workflow.
func FromWorkflow(w *sws.LoadedWorkflow) (*execgraph.Workflow, error) {
	b := execgraph.NewBuilder(w.Workflow.Name).
		Description(w.Workflow.Description)

	// Add phases
	for _, pd := range phaseDefs {
		b.PhaseWithDescription(pd.ID, pd.Name, pd.Description, pd.Order)
	}

	specConfig := types.SpecConfigFromWorkflow(w.Workflow)
	if specConfig == nil {
		specConfig = types.DefaultSpecConfig()
	}

	// Add nodes for each spec in the config
	for specName, req := range specConfig.Specs {
		if req == nil {
			continue
		}

		meta := specMeta[specName]
		if meta.Name == "" {
			meta.Name = specName
		}

		// Determine phase based on category
		phase := categoryToPhase(req.Category)

		// Get dependencies from the workflow's own synthesis configuration
		deps := synthesisSources(w.Workflow, specName)

		nodeType := categoryToNodeType(req.Category)
		automated := isAutomated(req.Category)

		b.Phase(phase, "", 0) // Set current phase (order doesn't matter here)
		b.Node(specName, meta.Name).
			Description(meta.Description).
			Type(nodeType).
			DependsOn(deps...).
			Automated().
			Metadata("category", string(req.Category)).
			Metadata("required", req.Required)

		if !automated {
			// Human-authored specs aren't marked automated
			b.Node(specName, meta.Name).
				Description(meta.Description).
				Type(nodeType).
				DependsOn(deps...).
				Metadata("category", string(req.Category)).
				Metadata("required", req.Required).
				Add()
		} else {
			b.Node(specName, meta.Name).
				Description(meta.Description).
				Type(nodeType).
				DependsOn(deps...).
				Automated().
				Metadata("category", string(req.Category)).
				Metadata("required", req.Required).
				Add()
		}
	}

	return b.Build()
}

// synthesisSources returns the workflow's own declared synthesis sources for
// a spec type, or nil if the workflow declares no synthesis rule for it
// (e.g. a human-authored document like the press release).
func synthesisSources(w *workflow.Workflow, specName string) []string {
	rule, ok := w.Synthesis[specName]
	if !ok || rule == nil {
		return nil
	}
	return rule.Sources
}

func categoryToPhase(cat types.SpecCategory) string {
	switch cat {
	case types.CategorySource:
		return "discovery" // or "product" depending on spec type
	case types.CategoryGTM:
		return "vision"
	case types.CategoryTechnical:
		return "technical"
	case types.CategoryOutput:
		return "reconcile"
	default:
		return "discovery"
	}
}

func categoryToNodeType(cat types.SpecCategory) execgraph.NodeType {
	switch cat {
	case types.CategorySource:
		return TypeSource
	case types.CategoryGTM:
		return TypeGTM
	case types.CategoryTechnical:
		return TypeTechnical
	case types.CategoryOutput:
		return TypeOutput
	default:
		return TypeSource
	}
}

func isAutomated(cat types.SpecCategory) bool {
	// GTM and Technical specs are typically LLM-synthesized
	return cat == types.CategoryGTM || cat == types.CategoryTechnical || cat == types.CategoryOutput
}

// UpdateFromProject updates workflow node statuses from project state.
func UpdateFromProject(w *execgraph.Workflow, project *types.Project) {
	if project == nil || project.Specs == nil {
		return
	}

	for specType, spec := range project.Specs {
		node, ok := w.GetNode(string(specType))
		if !ok {
			continue
		}

		// Map spec status to workflow status
		switch spec.Status {
		case types.StatusDraft:
			node.Status = execgraph.StatusInProgress
		case types.StatusEvaluated:
			// Evaluated specs are in progress until approved
			node.Status = execgraph.StatusInProgress
		case types.StatusApproved:
			node.Status = execgraph.StatusCompleted
		case types.StatusRejected:
			node.Status = execgraph.StatusBlocked
		case types.StatusMissing:
			if w.IsReady(node.ID) {
				node.Status = execgraph.StatusReady
			} else {
				node.Status = execgraph.StatusPending
			}
		default:
			if w.IsReady(node.ID) {
				node.Status = execgraph.StatusReady
			} else {
				node.Status = execgraph.StatusPending
			}
		}
	}

	// Update ready status for pending nodes
	for _, node := range w.Nodes {
		if node.Status == execgraph.StatusPending && w.IsReady(node.ID) {
			node.Status = execgraph.StatusReady
		}
	}
}
