// Package workflow defines specification workflow configurations.
//
// A Workflow bundles spec requirements, synthesis rules, and evaluation
// criteria into a cohesive workflow configuration (e.g., "aws-one-way-door",
// "big-tech-feature", "pbhq-lite").
package workflow

// Workflow represents a complete specification workflow configuration.
type Workflow struct {
	// Name is the workflow identifier (e.g., "aws-one-way-door", "pbhq-lite").
	Name string `json:"name" yaml:"name" jsonschema:"required,description=Workflow identifier"`

	// Description explains the workflow's purpose and use case.
	Description string `json:"description,omitempty" yaml:"description,omitempty" jsonschema:"description=Workflow purpose and target audience"`

	// Extends is the name of a parent workflow to inherit from.
	Extends string `json:"extends,omitempty" yaml:"extends,omitempty" jsonschema:"description=Parent workflow to inherit settings from"`

	// Abstract indicates this workflow is a base for other workflows (not directly usable).
	Abstract bool `json:"abstract,omitempty" yaml:"abstract,omitempty" jsonschema:"description=True if this workflow cannot be used directly"`

	// Methodology documents the underlying product methodology.
	Methodology *Methodology `json:"methodology,omitempty" yaml:"methodology,omitempty" jsonschema:"description=Underlying product methodology documentation"`

	// SpecConfig defines which specs are required/optional.
	SpecConfig map[string]*SpecRequirement `json:"spec_config,omitempty" yaml:"spec_config,omitempty" jsonschema:"description=Spec requirements by spec type ID"`

	// Synthesis defines how specs are generated from other specs.
	Synthesis map[string]*SynthesisRule `json:"synthesis,omitempty" yaml:"synthesis,omitempty" jsonschema:"description=Synthesis rules by target spec type"`

	// Execution defines the ordered phases and gates.
	Execution *Execution `json:"execution,omitempty" yaml:"execution,omitempty" jsonschema:"description=Phase ordering and gates"`

	// Evaluation defines pass/fail thresholds.
	Evaluation *EvaluationConfig `json:"evaluation,omitempty" yaml:"evaluation,omitempty" jsonschema:"description=Evaluation thresholds"`
}

// SynthesisRule defines how a spec can be synthesized from source specs.
type SynthesisRule struct {
	// Sources are the spec type IDs required to synthesize this spec.
	Sources []string `json:"sources" yaml:"sources" jsonschema:"required,description=Source spec type IDs"`

	// Guidance is the prompt context for LLM synthesis.
	Guidance string `json:"guidance,omitempty" yaml:"guidance,omitempty" jsonschema:"description=LLM prompt guidance for synthesis"`

	// PromptContext is additional context for the synthesis prompt.
	PromptContext string `json:"prompt_context,omitempty" yaml:"prompt_context,omitempty" jsonschema:"description=Additional synthesis prompt context"`

	// Required indicates all sources must be present (vs. best-effort).
	Required bool `json:"required,omitempty" yaml:"required,omitempty" jsonschema:"description=Whether all sources are required"`

	// Priority determines synthesis order when multiple rules exist.
	Priority int `json:"priority,omitempty" yaml:"priority,omitempty" jsonschema:"description=Synthesis priority (higher = earlier)"`
}

// Methodology documents the underlying product development methodology.
type Methodology struct {
	// Name is the methodology name (e.g., "Amazon Working Backwards").
	Name string `json:"name" yaml:"name" jsonschema:"required,description=Methodology name"`

	// Description explains the methodology.
	Description string `json:"description,omitempty" yaml:"description,omitempty" jsonschema:"description=Methodology overview"`

	// Creator is the person/company who created the methodology.
	Creator string `json:"creator,omitempty" yaml:"creator,omitempty" jsonschema:"description=Methodology creator"`

	// Source is the origin company or publication of the methodology.
	Source string `json:"source,omitempty" yaml:"source,omitempty" jsonschema:"description=Origin company or publication"`

	// Reference is a URL to the canonical methodology documentation.
	Reference string `json:"reference,omitempty" yaml:"reference,omitempty" jsonschema:"format=uri,description=URL to methodology documentation"`

	// Principles are the core principles of the methodology.
	Principles []Principle `json:"principles,omitempty" yaml:"principles,omitempty" jsonschema:"description=Core methodology principles"`

	// Artifacts are the key artifacts produced by the methodology.
	Artifacts Artifacts `json:"artifacts,omitempty" yaml:"artifacts,omitempty" jsonschema:"description=Key artifacts"`
}

// Artifacts is a list of methodology artifacts.
// In YAML it accepts either a flat sequence or a mapping of category name
// to sequence (e.g., primary/supporting groups).
type Artifacts []Artifact

// Artifact is a named artifact produced by a methodology.
type Artifact struct {
	// ID is the artifact identifier (e.g., "press_release").
	ID string `json:"id" yaml:"id" jsonschema:"required,description=Artifact identifier"`

	// Description explains the artifact.
	Description string `json:"description,omitempty" yaml:"description,omitempty" jsonschema:"description=Artifact explanation"`

	// Category optionally groups artifacts (e.g., "primary", "supporting").
	Category string `json:"category,omitempty" yaml:"category,omitempty" jsonschema:"description=Artifact grouping"`
}

// Principle is a named principle with description.
type Principle struct {
	// ID is the principle identifier (e.g., "customer_obsession").
	ID string `json:"id" yaml:"id" jsonschema:"required,description=Principle identifier"`

	// Name is the human-readable name.
	Name string `json:"name" yaml:"name" jsonschema:"required,description=Principle name"`

	// Description explains the principle.
	Description string `json:"description,omitempty" yaml:"description,omitempty" jsonschema:"description=Principle explanation"`

	// Source is the origin (e.g., "Amazon", "Google").
	Source string `json:"source,omitempty" yaml:"source,omitempty" jsonschema:"description=Origin company or methodology"`
}

// SpecRequirement defines whether a spec type is required and its configuration.
type SpecRequirement struct {
	// Required indicates whether this spec must be present.
	Required bool `json:"required" yaml:"required" jsonschema:"description=Whether this spec is required"`

	// Category overrides the default category for this spec type.
	Category string `json:"category,omitempty" yaml:"category,omitempty" jsonschema:"enum=source,enum=gtm,enum=technical,enum=execution,enum=output,enum=strategic"`

	// Description provides workflow-specific context for this spec.
	Description string `json:"description,omitempty" yaml:"description,omitempty" jsonschema:"description=Workflow-specific description"`

	// Template declares where this spec's document template comes from. Omit to
	// use default resolution (the workflow's own templates/ dir, then the
	// extends chain).
	Template *SpecSource `json:"template,omitempty" yaml:"template,omitempty" jsonschema:"description=Provenance of this spec's template"`

	// Rubric declares where this spec's evaluation rubric comes from. Omit to
	// use default resolution (the workflow's own rubrics/ dir, then the extends
	// chain).
	Rubric *SpecSource `json:"rubric,omitempty" yaml:"rubric,omitempty" jsonschema:"description=Provenance of this spec's rubric"`
}

// SourceLocal is the SpecSource.From sentinel meaning this workflow's own
// templates/ or rubrics/ directory (as opposed to another workflow).
const SourceLocal = "local"

// SpecSource declares the provenance of a spec's template or rubric: the
// workflow that owns the file. The sentinel "local" (SourceLocal) means this
// workflow's own directory; any other value names the workflow to resolve the
// file from. Declaring a source makes provenance explicit and loader-enforced —
// a source that does not actually provide the file is a load-time error.
//
// A source must not lead back to the declaring workflow: it is resolved with
// full inheritance, so naming a descendant (whose extends chain necessarily
// returns to the declaring workflow) or any other workflow whose resolution
// path re-enters it is a circular reference and a load-time error. Point
// sources at ancestors or unrelated workflows; content owned by a descendant
// must be physically copied, not referenced.
//
// In YAML it accepts the object form or a bare-string shorthand:
//
//	template: {from: enterprise}
//	rubric: local
type SpecSource struct {
	// From is the owning workflow name, or "local" for this workflow's own dir.
	From string `json:"from" yaml:"from" jsonschema:"required,description=Owning workflow name, or \"local\" for this workflow's own directory"`
}

// clone returns a deep copy, or nil if the receiver is nil.
func (s *SpecSource) clone() *SpecSource {
	if s == nil {
		return nil
	}
	c := *s
	return &c
}

// Execution defines the ordered execution of specs.
type Execution struct {
	// Sequence is the ordered list of spec types to produce.
	Sequence []string `json:"sequence,omitempty" yaml:"sequence,omitempty" jsonschema:"description=Ordered spec type IDs"`

	// Phases groups specs into named phases.
	Phases []Phase `json:"phases,omitempty" yaml:"phases,omitempty" jsonschema:"description=Named workflow phases"`

	// IterationTrigger is the spec type that triggers iteration.
	IterationTrigger string `json:"iteration_trigger,omitempty" yaml:"iteration_trigger,omitempty" jsonschema:"description=Spec type that triggers workflow iteration"`

	// ReviewGates are approval checkpoints.
	ReviewGates []ReviewGate `json:"review_gates,omitempty" yaml:"review_gates,omitempty" jsonschema:"description=Approval checkpoints"`
}

// Phase is a named group of specs in the workflow.
type Phase struct {
	// ID is the phase identifier.
	ID string `json:"id" yaml:"id" jsonschema:"required,description=Phase identifier"`

	// Name is the human-readable phase name.
	Name string `json:"name" yaml:"name" jsonschema:"required,description=Phase name"`

	// Description explains the phase purpose.
	Description string `json:"description,omitempty" yaml:"description,omitempty" jsonschema:"description=Phase purpose"`

	// Specs are the spec types in this phase.
	Specs []string `json:"specs" yaml:"specs" jsonschema:"required,description=Spec type IDs in this phase"`
}

// ReviewGate is an approval checkpoint after a spec.
type ReviewGate struct {
	// After is the spec type after which this gate applies.
	After string `json:"after" yaml:"after" jsonschema:"required,description=Spec type ID after which gate applies"`

	// Action is the required action (e.g., "stakeholder_review", "tech_lead_review").
	Action string `json:"action" yaml:"action" jsonschema:"required,description=Required approval action"`

	// Required indicates whether passing this gate is mandatory.
	Required bool `json:"required,omitempty" yaml:"required,omitempty" jsonschema:"description=Whether gate is mandatory"`
}

// EvaluationConfig defines pass/fail thresholds.
type EvaluationConfig struct {
	// PassThreshold is the minimum score (0-100) to pass.
	PassThreshold int `json:"pass_threshold,omitempty" yaml:"pass_threshold,omitempty" jsonschema:"minimum=0,maximum=100,description=Minimum score to pass"`

	// PartialThreshold is the minimum score (0-100) for partial pass.
	PartialThreshold int `json:"partial_threshold,omitempty" yaml:"partial_threshold,omitempty" jsonschema:"minimum=0,maximum=100,description=Minimum score for partial pass"`

	// MaxFindingsSeverity defines maximum allowed findings by severity.
	MaxFindingsSeverity *FindingSeverityLimits `json:"max_findings_severity,omitempty" yaml:"max_findings_severity,omitempty" jsonschema:"description=Maximum findings allowed by severity"`
}

// FindingSeverityLimits defines maximum findings by severity level.
type FindingSeverityLimits struct {
	Critical int `json:"critical,omitempty" yaml:"critical,omitempty" jsonschema:"description=Max critical findings (-1 = unlimited)"`
	High     int `json:"high,omitempty" yaml:"high,omitempty" jsonschema:"description=Max high findings (-1 = unlimited)"`
	Medium   int `json:"medium,omitempty" yaml:"medium,omitempty" jsonschema:"description=Max medium findings (-1 = unlimited)"`
	Low      int `json:"low,omitempty" yaml:"low,omitempty" jsonschema:"description=Max low findings (-1 = unlimited)"`
}
