package types

import (
	"time"
)

// Project represents a visionspec project.
type Project struct {
	// Name is the project identifier (kebab-case).
	Name string `json:"name" yaml:"name"`

	// Path is the absolute path to the project directory.
	Path string `json:"path" yaml:"path"`

	// Constitution is the path to the constitution file (relative or absolute).
	Constitution string `json:"constitution,omitempty" yaml:"constitution,omitempty"`

	// LLM configures the LLM provider for evaluations and synthesis.
	LLM *LLMConfig `json:"llm,omitempty" yaml:"llm,omitempty"`

	// Specs contains the status of each spec in the project.
	Specs map[SpecType]*Spec `json:"specs,omitempty" yaml:"specs,omitempty"`

	// Approvals tracks approval status for each spec.
	Approvals map[SpecType]*Approval `json:"approvals,omitempty" yaml:"approvals,omitempty"`

	// Targets configures export targets.
	Targets TargetConfig `json:"targets,omitempty" yaml:"targets,omitempty"`

	// SpecRequirements configures which specs are required and their settings.
	// This appears as "spec_config:" in visionspec.yaml.
	SpecRequirements map[string]*SpecRequirement `json:"spec_config,omitempty" yaml:"spec_config,omitempty"`

	// Context configures context sources for grounding.
	Context *ContextConfig `json:"context,omitempty" yaml:"context,omitempty"`

	// Rubrics configures custom rubric loading.
	Rubrics *RubricsConfig `json:"rubrics,omitempty" yaml:"rubrics,omitempty"`

	// Workflow specifies the workflow methodology and level (e.g., "aws-working-backwards/product").
	// Deprecated: Use RequirementsMethodology instead.
	Workflow string `json:"workflow,omitempty" yaml:"workflow,omitempty"`

	// RequirementsMethodology specifies the requirements methodology
	// (e.g., "aws-working-backwards/product", "big-tech-product", "lean-startup").
	// This defines WHAT to build.
	RequirementsMethodology string `json:"requirements_methodology,omitempty" yaml:"requirements_methodology,omitempty"`

	// ImplementationMethodology specifies the implementation methodology
	// (e.g., "aidlc", "speckit", "none").
	// This defines HOW to build.
	ImplementationMethodology ImplementationMethodology `json:"implementation_methodology,omitempty" yaml:"implementation_methodology,omitempty"`

	// Execution tracks the state of exported execution targets.
	Execution *ExecutionState `json:"execution,omitempty" yaml:"execution,omitempty"`

	// CreatedAt is when the project was initialized.
	CreatedAt time.Time `json:"created_at" yaml:"created_at"`

	// UpdatedAt is when the project was last modified.
	UpdatedAt time.Time `json:"updated_at" yaml:"updated_at"`
}

// ExecutionState tracks the state of tasks in an exported target.
type ExecutionState struct {
	Target   string           `json:"target" yaml:"target"`
	SyncedAt time.Time        `json:"synced_at" yaml:"synced_at"`
	Tasks    []ExecutionTask  `json:"tasks" yaml:"tasks"`
	Summary  ExecutionSummary `json:"summary" yaml:"summary"`
}

// ExecutionTask represents a task in the execution state.
type ExecutionTask struct {
	ID     string `json:"id" yaml:"id"`
	Title  string `json:"title" yaml:"title"`
	Status string `json:"status" yaml:"status"` // todo, in_progress, done, blocked
}

// ExecutionSummary provides aggregate statistics.
type ExecutionSummary struct {
	TotalTasks int `json:"total_tasks" yaml:"total_tasks"`
	TodoCount  int `json:"todo_count" yaml:"todo_count"`
	InProgress int `json:"in_progress" yaml:"in_progress"`
	DoneCount  int `json:"done_count" yaml:"done_count"`
}

// LLMConfig configures the LLM provider for a project.
type LLMConfig struct {
	// Provider is the LLM provider (anthropic, openai, gemini, etc.).
	Provider string `json:"provider,omitempty" yaml:"provider,omitempty"`

	// Model is the specific model to use (e.g., claude-sonnet-4-20250514, gpt-4o).
	Model string `json:"model,omitempty" yaml:"model,omitempty"`

	// Temperature controls randomness (0.0 = deterministic, higher = more random).
	Temperature *float64 `json:"temperature,omitempty" yaml:"temperature,omitempty"`

	// MaxTokens limits the response length.
	MaxTokens *int `json:"max_tokens,omitempty" yaml:"max_tokens,omitempty"`
}

// Approval represents an approval record for a spec.
type Approval struct {
	Approver   string    `json:"approver" yaml:"approver"`
	ApprovedAt time.Time `json:"approved_at" yaml:"approved_at"`
	Comment    string    `json:"comment,omitempty" yaml:"comment,omitempty"`
}

// TargetConfig configures export targets for a project.
type TargetConfig struct {
	Default string `json:"default,omitempty" yaml:"default,omitempty"`

	SpecKit  *SpecKitConfig  `json:"speckit,omitempty" yaml:"speckit,omitempty"`
	GSD      *GSDConfig      `json:"gsd,omitempty" yaml:"gsd,omitempty"`
	GasTown  *GasTownConfig  `json:"gastown,omitempty" yaml:"gastown,omitempty"`
	GasCity  *GasCityConfig  `json:"gascity,omitempty" yaml:"gascity,omitempty"`
	AIDLC    *AIDLCConfig    `json:"aidlc,omitempty" yaml:"aidlc,omitempty"`
	OpenSpec *OpenSpecConfig `json:"openspec,omitempty" yaml:"openspec,omitempty"`
}

// SpecKitConfig configures the SpecKit export target.
type SpecKitConfig struct {
	Enabled         bool   `json:"enabled" yaml:"enabled"`
	OutputDir       string `json:"output_dir,omitempty" yaml:"output_dir,omitempty"`
	BranchNumbering string `json:"branch_numbering,omitempty" yaml:"branch_numbering,omitempty"` // "sequential" or "timestamp"
}

// GSDConfig configures the GSD export target.
type GSDConfig struct {
	Enabled      bool   `json:"enabled" yaml:"enabled"`
	OutputDir    string `json:"output_dir,omitempty" yaml:"output_dir,omitempty"`
	ModelProfile string `json:"model_profile,omitempty" yaml:"model_profile,omitempty"` // "balanced", "quality", "budget"
}

// GasTownConfig configures the GasTown export target.
type GasTownConfig struct {
	Enabled     bool   `json:"enabled" yaml:"enabled"`
	FormulaType string `json:"formula_type,omitempty" yaml:"formula_type,omitempty"` // "convoy", "workflow", "expansion"
	Rig         string `json:"rig,omitempty" yaml:"rig,omitempty"`
}

// GasCityConfig configures the GasCity export target.
type GasCityConfig struct {
	Enabled bool   `json:"enabled" yaml:"enabled"`
	CityDir string `json:"city_dir,omitempty" yaml:"city_dir,omitempty"`
}

// OpenSpecConfig configures the OpenSpec export target.
type OpenSpecConfig struct {
	Enabled bool `json:"enabled" yaml:"enabled"`
}

// AIDLCConfig configures the AWS AI-DLC Workflows export target.
type AIDLCConfig struct {
	Enabled   bool   `json:"enabled" yaml:"enabled"`
	OutputDir string `json:"output_dir,omitempty" yaml:"output_dir,omitempty"` // default: ".aidlc"
}

// ContextConfig configures context sources for grounding specs.
type ContextConfig struct {
	// Repositories are git repositories to analyze.
	Repositories []RepositoryContextConfig `json:"repositories,omitempty" yaml:"repositories,omitempty"`

	// Graphize are standalone graphize graph paths.
	Graphize []GraphizeContextConfig `json:"graphize,omitempty" yaml:"graphize,omitempty"`

	// Files are local files to include as context.
	Files []FileContextConfig `json:"files,omitempty" yaml:"files,omitempty"`

	// MCPServers are MCP servers for external context.
	MCPServers map[string]MCPServerContextConfig `json:"mcp_servers,omitempty" yaml:"mcp_servers,omitempty"`

	// CacheTTL is how long to cache context data.
	CacheTTL time.Duration `json:"cache_ttl,omitempty" yaml:"cache_ttl,omitempty"`
}

// RepositoryContextConfig configures a git repository context source.
type RepositoryContextConfig struct {
	Path     string   `json:"path,omitempty" yaml:"path,omitempty"`
	URL      string   `json:"url,omitempty" yaml:"url,omitempty"`
	Branch   string   `json:"branch,omitempty" yaml:"branch,omitempty"`
	Include  []string `json:"include,omitempty" yaml:"include,omitempty"`
	Exclude  []string `json:"exclude,omitempty" yaml:"exclude,omitempty"`
	Analyze  []string `json:"analyze,omitempty" yaml:"analyze,omitempty"`
	Graphize string   `json:"graphize,omitempty" yaml:"graphize,omitempty"`
	MaxDepth int      `json:"max_depth,omitempty" yaml:"max_depth,omitempty"`
}

// GraphizeContextConfig configures a graphize context source.
type GraphizeContextConfig struct {
	Path         string   `json:"path" yaml:"path"`
	Name         string   `json:"name,omitempty" yaml:"name,omitempty"`
	IncludeNodes []string `json:"include_nodes,omitempty" yaml:"include_nodes,omitempty"`
	IncludeEdges []string `json:"include_edges,omitempty" yaml:"include_edges,omitempty"`
}

// FileContextConfig configures a local file context source.
type FileContextConfig struct {
	Path    string `json:"path" yaml:"path"`
	Type    string `json:"type,omitempty" yaml:"type,omitempty"`
	MaxSize int64  `json:"max_size,omitempty" yaml:"max_size,omitempty"`
}

// MCPServerContextConfig configures an MCP server context source.
type MCPServerContextConfig struct {
	Command string            `json:"command" yaml:"command"`
	Args    []string          `json:"args,omitempty" yaml:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty" yaml:"env,omitempty"`
	Config  map[string]any    `json:"config,omitempty" yaml:"config,omitempty"`
	Timeout time.Duration     `json:"timeout,omitempty" yaml:"timeout,omitempty"`
}

// RubricsConfig configures custom rubric loading for a project.
type RubricsConfig struct {
	// Directory is a path to a directory containing .rubric.yaml files.
	// Rubrics are named: {spec-type}.rubric.yaml (e.g., prd.rubric.yaml).
	Directory string `json:"directory,omitempty" yaml:"directory,omitempty"`

	// Overrides maps spec types to specific rubric file paths.
	// This allows using different rubric files for specific specs.
	Overrides map[SpecType]string `json:"overrides,omitempty" yaml:"overrides,omitempty"`

	// StrictMode requires all categories to pass (no partial scores).
	StrictMode bool `json:"strict_mode,omitempty" yaml:"strict_mode,omitempty"`

	// MaxCritical is the maximum number of critical findings allowed (default: 0).
	MaxCritical int `json:"max_critical,omitempty" yaml:"max_critical,omitempty"`

	// MaxHigh is the maximum number of high findings allowed (default: 0).
	MaxHigh int `json:"max_high,omitempty" yaml:"max_high,omitempty"`

	// MaxMedium is the maximum number of medium findings allowed (-1 = unlimited).
	MaxMedium int `json:"max_medium,omitempty" yaml:"max_medium,omitempty"`
}

// ReadinessGate represents a readiness check for a project.
type ReadinessGate struct {
	Name    string `json:"name" yaml:"name"`
	Passed  bool   `json:"passed" yaml:"passed"`
	Message string `json:"message,omitempty" yaml:"message,omitempty"`
}

// ReadinessStatus represents the overall readiness of a project.
type ReadinessStatus struct {
	Ready   bool            `json:"ready" yaml:"ready"`
	Gates   []ReadinessGate `json:"gates" yaml:"gates"`
	Summary string          `json:"summary" yaml:"summary"`
}

// GetSpecConfig returns a SpecConfig wrapper for the project's spec requirements.
// This provides helper methods like IsRequired(), GetCategory(), etc.
func (p *Project) GetSpecConfig() *SpecConfig {
	if p == nil {
		return DefaultSpecConfig()
	}
	if p.SpecRequirements == nil {
		return DefaultSpecConfig()
	}
	// Merge project requirements with defaults
	config := DefaultSpecConfig()
	config.Merge(&SpecConfig{Specs: p.SpecRequirements})
	return config
}
