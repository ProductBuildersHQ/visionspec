// Package integration defines declarative descriptors for external
// spec-driven-development tools (e.g., GitHub Spec-Kit, AWS AI-DLC Workflows,
// OpenSpec, AWS Kiro).
//
// An Integration describes how to recognize a tool's project on disk, the
// artifacts and lifecycle it uses, how a viewer computes progress from those
// artifacts, and where its inputs and outputs connect to the specification
// workflow. It holds NO execution or filesystem-scanning logic: consumers
// (visionspec, visionstudio) implement detection and status against this
// contract. This keeps the specification-workflow-spec layer purely
// declarative, consistent with the rest of the library.
//
// Integrations complement Workflows. A Workflow (see pkg/workflow) is a
// definition-side methodology that produces what-to-build specs. An Integration
// most often describes an execution-side tool that consumes those specs and
// produces code — the "Definition" and "Execution" halves surfaced in the
// VisionStudio viewer. A Pipeline (see pkg/pipeline) links the two.
package integration

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Kind classifies an integration by which side of the spec-driven-development
// divide it serves.
type Kind string

const (
	// KindDefinition produces what-to-build specs (requirements, design).
	KindDefinition Kind = "definition"

	// KindExecution consumes specs and produces code, tests, or infrastructure.
	KindExecution Kind = "execution"

	// KindHybrid spans both definition and execution.
	KindHybrid Kind = "hybrid"
)

// Invocation is a mode by which a tool is driven.
type Invocation string

const (
	// InvocationCLI is a command-line entrypoint.
	InvocationCLI Invocation = "cli"

	// InvocationSlashCommand is an agent slash command (e.g., "/specify").
	InvocationSlashCommand Invocation = "slash-command"

	// InvocationMCP is a Model Context Protocol server.
	InvocationMCP Invocation = "mcp"

	// InvocationIDE is an integrated development environment feature.
	InvocationIDE Invocation = "ide"

	// InvocationLibrary is an importable SDK.
	InvocationLibrary Invocation = "library"
)

// SourceKind indicates how the descriptor was modeled, and therefore how much
// confidence a consumer should place in its on-disk detail.
type SourceKind string

const (
	// SourceRepo means the descriptor was modeled from an inspectable
	// open-source repository.
	SourceRepo SourceKind = "repo"

	// SourceDocs means the descriptor was modeled from published documentation
	// or prose only (no inspectable repository).
	SourceDocs SourceKind = "docs"
)

// ArtifactRole classifies an on-disk artifact by the part it plays in the
// tool's lifecycle.
type ArtifactRole string

const (
	// RoleConstitution is a governing set of principles or rules.
	RoleConstitution ArtifactRole = "constitution"

	// RoleRequirement captures what must be built.
	RoleRequirement ArtifactRole = "requirement"

	// RoleDesign captures how it will be built (architecture, approach).
	RoleDesign ArtifactRole = "design"

	// RolePlan sequences work into steps or phases.
	RolePlan ArtifactRole = "plan"

	// RoleTasks is a checklist of executable units of work.
	RoleTasks ArtifactRole = "tasks"

	// RoleProposal is a delta-based change proposal.
	RoleProposal ArtifactRole = "proposal"

	// RoleSpec is a canonical capability specification.
	RoleSpec ArtifactRole = "spec"

	// RoleResearch captures investigation or discovery notes.
	RoleResearch ArtifactRole = "research"

	// RoleContract is an interface or API contract.
	RoleContract ArtifactRole = "contract"

	// RoleConfig is tool configuration.
	RoleConfig ArtifactRole = "config"

	// RoleOutput is generated code, tests, or infrastructure.
	RoleOutput ArtifactRole = "output"
)

// StatusMethod names a strategy for computing progress from on-disk artifacts.
type StatusMethod string

const (
	// StatusTaskCheckboxes counts completed vs. total checklist items in a
	// tasks artifact (e.g., "- [x]" vs. "- [ ]").
	StatusTaskCheckboxes StatusMethod = "task-checkboxes"

	// StatusPhaseFiles infers progress from the presence or absence of the
	// artifacts each lifecycle phase produces.
	StatusPhaseFiles StatusMethod = "phase-files"

	// StatusDirState infers progress from which directory an item lives in
	// (e.g., active vs. archived).
	StatusDirState StatusMethod = "dir-state"

	// StatusFrontmatter reads a status field from artifact frontmatter.
	StatusFrontmatter StatusMethod = "frontmatter"
)

// PortMedium classifies what flows across an input or output seam.
type PortMedium string

const (
	// MediumSpecFile is a specification document.
	MediumSpecFile PortMedium = "spec-file"

	// MediumCode is source code.
	MediumCode PortMedium = "code"

	// MediumTests is test code or test plans.
	MediumTests PortMedium = "tests"

	// MediumPrompt is a natural-language intent or description.
	MediumPrompt PortMedium = "prompt"

	// MediumExport is a structured export handed to another tool.
	MediumExport PortMedium = "export"
)

// Integration is a declarative descriptor of an external spec-driven-development
// tool.
type Integration struct {
	// ID is the canonical identifier (e.g., "spec-kit", "ai-dlc", "openspec").
	ID string `json:"id" yaml:"id" jsonschema:"required,description=Canonical identifier for the integration"`

	// Name is the human-readable name (e.g., "GitHub Spec-Kit").
	Name string `json:"name" yaml:"name" jsonschema:"required,description=Human-readable name"`

	// ShortName is an abbreviated name for compact UI display.
	ShortName string `json:"short_name,omitempty" yaml:"short_name,omitempty" jsonschema:"description=Abbreviated name for compact display"`

	// Vendor is the owning organization or project (e.g., "GitHub", "AWS Labs").
	Vendor string `json:"vendor,omitempty" yaml:"vendor,omitempty" jsonschema:"description=Owning organization or project"`

	// Description explains the tool's purpose in one or two sentences.
	Description string `json:"description,omitempty" yaml:"description,omitempty" jsonschema:"description=Purpose of the tool"`

	// Kind classifies the tool as definition, execution, or hybrid.
	Kind Kind `json:"kind" yaml:"kind" jsonschema:"required,enum=definition,enum=execution,enum=hybrid"`

	// Source indicates whether the descriptor was modeled from a repository or
	// from documentation only.
	Source SourceKind `json:"source,omitempty" yaml:"source,omitempty" jsonschema:"enum=repo,enum=docs,description=How this descriptor was modeled"`

	// Homepage is the tool's homepage URL.
	Homepage string `json:"homepage,omitempty" yaml:"homepage,omitempty" jsonschema:"format=uri,description=Homepage URL"`

	// Repo is the tool's source repository URL.
	Repo string `json:"repo,omitempty" yaml:"repo,omitempty" jsonschema:"format=uri,description=Source repository URL"`

	// Reference is a URL to the tool's canonical documentation.
	Reference string `json:"reference,omitempty" yaml:"reference,omitempty" jsonschema:"format=uri,description=Canonical documentation URL"`

	// License is the SPDX license identifier or name (e.g., "MIT").
	License string `json:"license,omitempty" yaml:"license,omitempty" jsonschema:"description=License identifier"`

	// Invocations lists the modes by which the tool is driven.
	Invocations []Invocation `json:"invocations,omitempty" yaml:"invocations,omitempty" jsonschema:"description=How the tool is driven"`

	// CLICommand is the primary CLI entrypoint, if any (e.g., "specify", "openspec").
	CLICommand string `json:"cli_command,omitempty" yaml:"cli_command,omitempty" jsonschema:"description=Primary CLI command name"`

	// Agents lists supported AI agents or runtimes (e.g., "Claude Code", "Cursor").
	Agents []string `json:"agents,omitempty" yaml:"agents,omitempty" jsonschema:"description=Supported AI agents or runtimes"`

	// Detection describes how to recognize a project using this tool on disk.
	Detection Detection `json:"detection" yaml:"detection" jsonschema:"required,description=How to detect the tool on disk"`

	// Artifacts are the files the tool reads or writes.
	Artifacts []Artifact `json:"artifacts,omitempty" yaml:"artifacts,omitempty" jsonschema:"description=Files the tool reads or writes"`

	// Lifecycle is the ordered list of phases the tool moves through.
	Lifecycle []Phase `json:"lifecycle,omitempty" yaml:"lifecycle,omitempty" jsonschema:"description=Ordered lifecycle phases"`

	// Status declares how a viewer computes progress from artifacts.
	Status *StatusModel `json:"status,omitempty" yaml:"status,omitempty" jsonschema:"description=How to compute progress from artifacts"`

	// Inputs are the seams where upstream specs feed into the tool.
	Inputs []Port `json:"inputs,omitempty" yaml:"inputs,omitempty" jsonschema:"description=Input seams (what feeds in)"`

	// Outputs are the seams where the tool emits its results.
	Outputs []Port `json:"outputs,omitempty" yaml:"outputs,omitempty" jsonschema:"description=Output seams (what it emits)"`
}

// Detection describes how to recognize a project using this tool on disk.
type Detection struct {
	// RootMarkers are paths (relative to a project root) whose presence
	// indicates the tool. Any match is sufficient (e.g., "openspec/", ".specify/").
	RootMarkers []string `json:"root_markers,omitempty" yaml:"root_markers,omitempty" jsonschema:"description=Paths whose presence indicates the tool (any match)"`

	// Globs are additional file globs that corroborate detection.
	Globs []string `json:"globs,omitempty" yaml:"globs,omitempty" jsonschema:"description=Additional corroborating file globs"`
}

// Artifact is a file (or file glob) the tool reads or writes.
type Artifact struct {
	// ID is the artifact identifier, unique within the integration.
	ID string `json:"id" yaml:"id" jsonschema:"required,description=Artifact identifier"`

	// Path is the path or glob relative to the project root
	// (e.g., "specs/<feature>/spec.md"). Angle-bracket segments denote variables.
	Path string `json:"path" yaml:"path" jsonschema:"required,description=Path or glob relative to project root"`

	// Role classifies the artifact.
	Role ArtifactRole `json:"role" yaml:"role" jsonschema:"required,enum=constitution,enum=requirement,enum=design,enum=plan,enum=tasks,enum=proposal,enum=spec,enum=research,enum=contract,enum=config,enum=output"`

	// Description explains the artifact's purpose.
	Description string `json:"description,omitempty" yaml:"description,omitempty" jsonschema:"description=Artifact purpose"`

	// Phase is the lifecycle phase ID that produces this artifact.
	Phase string `json:"phase,omitempty" yaml:"phase,omitempty" jsonschema:"description=Lifecycle phase ID that produces this artifact"`

	// SpecType optionally maps this artifact to a specification-workflow-spec
	// spec type ID (e.g., "trd", "prd"), enabling cross-tool alignment.
	SpecType string `json:"spec_type,omitempty" yaml:"spec_type,omitempty" jsonschema:"description=Mapped specification-workflow-spec spec type ID"`

	// Produced is true if the tool writes this artifact, false if it only
	// consumes it.
	Produced bool `json:"produced,omitempty" yaml:"produced,omitempty" jsonschema:"description=True if the tool writes this artifact"`
}

// Phase is a step in the tool's lifecycle.
type Phase struct {
	// ID is the phase identifier.
	ID string `json:"id" yaml:"id" jsonschema:"required,description=Phase identifier"`

	// Name is the human-readable phase name.
	Name string `json:"name" yaml:"name" jsonschema:"required,description=Phase name"`

	// Description explains the phase.
	Description string `json:"description,omitempty" yaml:"description,omitempty" jsonschema:"description=Phase purpose"`

	// Command is the CLI command or slash command that drives the phase
	// (e.g., "/specify", "openspec archive").
	Command string `json:"command,omitempty" yaml:"command,omitempty" jsonschema:"description=Command or slash command that drives the phase"`

	// Produces lists artifact IDs emitted by this phase.
	Produces []string `json:"produces,omitempty" yaml:"produces,omitempty" jsonschema:"description=Artifact IDs emitted by this phase"`

	// Consumes lists artifact IDs required by this phase.
	Consumes []string `json:"consumes,omitempty" yaml:"consumes,omitempty" jsonschema:"description=Artifact IDs required by this phase"`
}

// StatusModel declares how a viewer computes progress from on-disk artifacts.
type StatusModel struct {
	// Method is the progress-computation strategy.
	Method StatusMethod `json:"method" yaml:"method" jsonschema:"required,enum=task-checkboxes,enum=phase-files,enum=dir-state,enum=frontmatter"`

	// Source is the artifact ID or glob the method reads.
	Source string `json:"source,omitempty" yaml:"source,omitempty" jsonschema:"description=Artifact ID or glob the method reads"`

	// Description explains how status is derived.
	Description string `json:"description,omitempty" yaml:"description,omitempty" jsonschema:"description=How status is derived"`
}

// Port describes a seam where specs flow in or artifacts flow out.
type Port struct {
	// ID is the port identifier, unique within the integration.
	ID string `json:"id" yaml:"id" jsonschema:"required,description=Port identifier"`

	// Medium classifies what flows across the seam.
	Medium PortMedium `json:"medium" yaml:"medium" jsonschema:"required,enum=spec-file,enum=code,enum=tests,enum=prompt,enum=export"`

	// Description explains the seam.
	Description string `json:"description,omitempty" yaml:"description,omitempty" jsonschema:"description=Seam explanation"`

	// SpecTypes lists specification-workflow-spec spec type IDs that map to this
	// port, used primarily for inputs (e.g., an execution tool that accepts a
	// "trd" and a "prd").
	SpecTypes []string `json:"spec_types,omitempty" yaml:"spec_types,omitempty" jsonschema:"description=Mapped spec type IDs (for inputs)"`

	// Artifacts lists artifact IDs associated with this port, used primarily for
	// outputs.
	Artifacts []string `json:"artifacts,omitempty" yaml:"artifacts,omitempty" jsonschema:"description=Associated artifact IDs (for outputs)"`
}

// ParseYAML parses an Integration from YAML bytes.
func ParseYAML(data []byte) (*Integration, error) {
	var in Integration
	if err := yaml.Unmarshal(data, &in); err != nil {
		return nil, fmt.Errorf("unmarshaling integration: %w", err)
	}
	return &in, nil
}

// Validate checks that required fields are present and internally consistent.
func (in *Integration) Validate() error {
	if in.ID == "" {
		return fmt.Errorf("integration: id is required")
	}
	if in.Name == "" {
		return fmt.Errorf("integration %q: name is required", in.ID)
	}
	switch in.Kind {
	case KindDefinition, KindExecution, KindHybrid:
	default:
		return fmt.Errorf("integration %q: invalid kind %q", in.ID, in.Kind)
	}
	if len(in.Detection.RootMarkers) == 0 && len(in.Detection.Globs) == 0 {
		return fmt.Errorf("integration %q: detection requires at least one root marker or glob", in.ID)
	}

	artifactIDs := make(map[string]bool, len(in.Artifacts))
	for _, a := range in.Artifacts {
		if a.ID == "" {
			return fmt.Errorf("integration %q: artifact with empty id", in.ID)
		}
		if artifactIDs[a.ID] {
			return fmt.Errorf("integration %q: duplicate artifact id %q", in.ID, a.ID)
		}
		artifactIDs[a.ID] = true
		if a.Path == "" {
			return fmt.Errorf("integration %q: artifact %q has empty path", in.ID, a.ID)
		}
	}

	// Phase artifact references must resolve.
	for _, p := range in.Lifecycle {
		if p.ID == "" {
			return fmt.Errorf("integration %q: phase with empty id", in.ID)
		}
		for _, ref := range append(append([]string{}, p.Produces...), p.Consumes...) {
			if !artifactIDs[ref] {
				return fmt.Errorf("integration %q: phase %q references unknown artifact %q", in.ID, p.ID, ref)
			}
		}
	}

	// Status source, if given, must reference a known artifact (unless it is a glob).
	if in.Status != nil && in.Status.Source != "" && !artifactIDs[in.Status.Source] {
		// A glob source (containing a path separator or wildcard) is allowed.
		if !looksLikePath(in.Status.Source) {
			return fmt.Errorf("integration %q: status source %q is not a known artifact id", in.ID, in.Status.Source)
		}
	}

	return nil
}

// looksLikePath reports whether s appears to be a path or glob rather than an
// artifact ID.
func looksLikePath(s string) bool {
	for _, r := range s {
		if r == '/' || r == '*' || r == '<' {
			return true
		}
	}
	return false
}
