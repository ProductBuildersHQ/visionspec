package types

import "time"

// Project represents a multispec project.
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

	// CreatedAt is when the project was initialized.
	CreatedAt time.Time `json:"created_at" yaml:"created_at"`

	// UpdatedAt is when the project was last modified.
	UpdatedAt time.Time `json:"updated_at" yaml:"updated_at"`
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
