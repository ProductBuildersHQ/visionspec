// Package types defines the core data structures for visionspec.
package types

import "sort"

// SpecRequirement defines requirements for a spec type.
type SpecRequirement struct {
	// Required indicates whether this spec is mandatory for the project.
	Required bool `json:"required" yaml:"required"`

	// Category specifies the spec category (source, gtm, technical).
	// Only needed for custom spec types; built-in types have implicit categories.
	Category SpecCategory `json:"category,omitempty" yaml:"category,omitempty"`

	// Template specifies which template to use (defaults to spec type name).
	Template string `json:"template,omitempty" yaml:"template,omitempty"`

	// Rubric specifies which rubric to use (defaults to spec type name).
	Rubric string `json:"rubric,omitempty" yaml:"rubric,omitempty"`
}

// SpecConfig holds per-spec configuration for a project.
type SpecConfig struct {
	// Specs maps spec type names to their requirements.
	Specs map[string]*SpecRequirement `json:"specs,omitempty" yaml:"specs,omitempty"`
}

// NewSpecConfig creates a new SpecConfig with default values.
func NewSpecConfig() *SpecConfig {
	return &SpecConfig{
		Specs: make(map[string]*SpecRequirement),
	}
}

// DefaultSpecConfig returns the default spec configuration.
func DefaultSpecConfig() *SpecConfig {
	return &SpecConfig{
		Specs: map[string]*SpecRequirement{
			// Source specs
			string(SpecTypeMRD): {Required: true, Category: CategorySource},
			string(SpecTypePRD): {Required: true, Category: CategorySource},
			string(SpecTypeUXD): {Required: true, Category: CategorySource},
			// GTM specs (optional by default)
			string(SpecTypePress):       {Required: false, Category: CategoryGTM},
			string(SpecTypeFAQ):         {Required: false, Category: CategoryGTM},
			string(SpecTypeNarrative1P): {Required: false, Category: CategoryGTM},
			string(SpecTypeNarrative6P): {Required: false, Category: CategoryGTM},
			// Technical specs
			string(SpecTypeTRD): {Required: true, Category: CategoryTechnical},
			string(SpecTypeIRD): {Required: false, Category: CategoryTechnical},
			// Output specs
			string(SpecTypeSpec):         {Required: false, Category: CategoryOutput},
			string(SpecTypeCurrentTruth): {Required: false, Category: CategoryOutput},
		},
	}
}

// IsRequired returns whether a spec type is required.
// It checks the config first, then falls back to built-in defaults.
func (sc *SpecConfig) IsRequired(specType string) bool {
	if sc == nil || sc.Specs == nil {
		// Fall back to built-in default
		return SpecType(specType).IsRequired()
	}

	if req, ok := sc.Specs[specType]; ok {
		return req.Required
	}

	// Fall back to built-in default for known types
	return SpecType(specType).IsRequired()
}

// GetCategory returns the category for a spec type.
func (sc *SpecConfig) GetCategory(specType string) SpecCategory {
	if sc != nil && sc.Specs != nil {
		if req, ok := sc.Specs[specType]; ok && req.Category != "" {
			return req.Category
		}
	}

	// Fall back to built-in category for known types
	return SpecType(specType).Category()
}

// GetTemplate returns the template name for a spec type.
// Returns the spec type name if no custom template is configured.
func (sc *SpecConfig) GetTemplate(specType string) string {
	if sc != nil && sc.Specs != nil {
		if req, ok := sc.Specs[specType]; ok && req.Template != "" {
			return req.Template
		}
	}
	return specType
}

// GetRubric returns the rubric name for a spec type.
// Returns the spec type name if no custom rubric is configured.
func (sc *SpecConfig) GetRubric(specType string) string {
	if sc != nil && sc.Specs != nil {
		if req, ok := sc.Specs[specType]; ok && req.Rubric != "" {
			return req.Rubric
		}
	}
	return specType
}

// GetRequirement returns the requirement for a spec type.
// Returns nil if not configured.
func (sc *SpecConfig) GetRequirement(specType string) *SpecRequirement {
	if sc == nil || sc.Specs == nil {
		return nil
	}
	return sc.Specs[specType]
}

// SetRequirement sets the requirement for a spec type.
func (sc *SpecConfig) SetRequirement(specType string, req *SpecRequirement) {
	if sc.Specs == nil {
		sc.Specs = make(map[string]*SpecRequirement)
	}
	sc.Specs[specType] = req
}

// CustomSpecs returns all custom (non-built-in) spec types.
func (sc *SpecConfig) CustomSpecs() []string {
	if sc == nil || sc.Specs == nil {
		return nil
	}

	var custom []string
	for name := range sc.Specs {
		if !SpecType(name).IsValid() {
			custom = append(custom, name)
		}
	}
	return custom
}

// AllSpecs returns all configured spec types (built-in + custom).
func (sc *SpecConfig) AllSpecs() []string {
	specs := make(map[string]bool)

	// Add built-in types
	for _, st := range AllSpecTypes() {
		specs[string(st)] = true
	}

	// Add custom types from config
	if sc != nil && sc.Specs != nil {
		for name := range sc.Specs {
			specs[name] = true
		}
	}

	result := make([]string, 0, len(specs))
	for name := range specs {
		result = append(result, name)
	}

	// Sort by true workflow order (AWS Working Backwards flow)
	// Discovery → Vision (PR/FAQ/Narrative) → Product (PRD/UXD) → Technical → Output
	// The 6-Pager is the leadership approval gate BEFORE detailed PRD/UXD work
	workflowOrder := map[string]int{
		// Discovery phase
		"opportunity-spec": 0,
		"mrd":              1,
		// Vision phase (Working Backwards - PR/FAQ then 6-Pager for approval)
		"press":        2,
		"faq":          3,
		"narrative-6p": 4, // Leadership approval gate
		"narrative-1p": 5, // Executive summary (optional)
		"bmc":          6, // Business model (optional)
		// Product phase (detailed work after approval)
		"prd": 7,
		"uxd": 8,
		// Technical phase
		"trd": 9,
		"ird": 10,
		"tpd": 11,
		// Output phase
		"spec":          12,
		"current-truth": 13,
	}

	sort.Slice(result, func(i, j int) bool {
		wfOrderI, okI := workflowOrder[result[i]]
		wfOrderJ, okJ := workflowOrder[result[j]]
		if okI && okJ {
			return wfOrderI < wfOrderJ
		}
		// Known specs come before unknown
		if okI != okJ {
			return okI
		}
		// Fall back to alphabetical for unknown specs
		return result[i] < result[j]
	})

	return result
}

// RequiredSpecs returns all required spec types.
func (sc *SpecConfig) RequiredSpecs() []string {
	var required []string
	for _, name := range sc.AllSpecs() {
		if sc.IsRequired(name) {
			required = append(required, name)
		}
	}
	return required
}

// Merge merges another SpecConfig into this one.
// Values from other override values in this config.
func (sc *SpecConfig) Merge(other *SpecConfig) {
	if other == nil || other.Specs == nil {
		return
	}

	if sc.Specs == nil {
		sc.Specs = make(map[string]*SpecRequirement)
	}

	for name, req := range other.Specs {
		sc.Specs[name] = req
	}
}

// SpecsByCategory returns all spec types in the given category.
func (sc *SpecConfig) SpecsByCategory(category SpecCategory) []string {
	var specs []string
	for _, name := range sc.AllSpecs() {
		if sc.GetCategory(name) == category {
			specs = append(specs, name)
		}
	}
	return specs
}
