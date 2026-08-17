// Package synth provides spec synthesis (generation) capabilities.
package synth

import (
	"context"
	"fmt"
	"strings"

	"github.com/ProductBuildersHQ/visionspec/pkg/workflow"
	sws "github.com/ProductBuildersHQ/visionspec/pkg/workflows"

	"github.com/ProductBuildersHQ/visionspec/pkg/types"
)

// Synthesizer generates specs from source documents using LLM.
type Synthesizer struct {
	client LLMClient
}

// LLMClient defines the interface for LLM operations.
type LLMClient interface {
	Complete(ctx context.Context, prompt string) (string, error)
}

// NewSynthesizer creates a new synthesizer with the given LLM client.
func NewSynthesizer(client LLMClient) *Synthesizer {
	return &Synthesizer{client: client}
}

// SynthesisInput contains the input documents for synthesis. Sources is
// keyed by spec type so callers aren't limited to a fixed set of document
// types — a workflow's synthesis rule can name any spec type as a source
// (opportunity-spec, narrative-6p, v2mom-*, etc.), not just the handful
// the pre-merge Working Backwards flow used.
type SynthesisInput struct {
	Sources      map[types.SpecType]string
	Constitution string // Project/org constitution
	Context      string // Aggregated context summary (for grounding)
}

// Set records the content of a source spec, lazily initializing the
// underlying map.
func (in *SynthesisInput) Set(specType types.SpecType, content string) {
	if in.Sources == nil {
		in.Sources = make(map[types.SpecType]string)
	}
	in.Sources[specType] = content
}

// Get returns the recorded content for a source spec, or "" if absent.
func (in SynthesisInput) Get(specType types.SpecType) string {
	return in.Sources[specType]
}

// SynthesisResult contains the generated document and metadata.
type SynthesisResult struct {
	SpecType types.SpecType
	Content  string
	Sources  []types.SpecType // Source specs used
}

// Synthesize generates a spec of the given type from input documents.
//
// rule is the target's synthesis rule from the project's actual configured
// workflow (profile.yaml's synthesis: block), when one is resolved. When
// rule is nil — no workflow configured, or the workflow declares no rule
// for targetType — the prompt falls back to the original hardcoded
// pre-merge Working Backwards prompts for the eight types they cover.
func (s *Synthesizer) Synthesize(ctx context.Context, targetType types.SpecType, input SynthesisInput, rule *workflow.SynthesisRule, templateContent string) (*SynthesisResult, error) {
	prompt, sources := s.buildPrompt(targetType, input, templateContent, rule)

	content, err := s.client.Complete(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("LLM synthesis failed: %w", err)
	}

	return &SynthesisResult{
		SpecType: targetType,
		Content:  content,
		Sources:  sources,
	}, nil
}

// buildPrompt constructs the synthesis prompt for a target spec type,
// preferring a workflow-supplied rule over the legacy hardcoded prompts.
func (s *Synthesizer) buildPrompt(targetType types.SpecType, input SynthesisInput, template string, rule *workflow.SynthesisRule) (string, []types.SpecType) {
	if rule != nil {
		return s.buildPromptFromRule(targetType, input, template, rule)
	}
	return s.buildLegacyPrompt(targetType, input, template)
}

// buildPromptFromRule builds a prompt from a workflow's synthesis rule:
// its declared sources (in order, whichever are actually present), its
// guidance/prompt_context — the methodology-specific instructions each
// profile.yaml already authors — and the resolved template.
func (s *Synthesizer) buildPromptFromRule(targetType types.SpecType, input SynthesisInput, template string, rule *workflow.SynthesisRule) (string, []types.SpecType) {
	var sb strings.Builder
	var sources []types.SpecType

	sb.WriteString("You are a technical writer synthesizing specification documents.\n\n")
	fmt.Fprintf(&sb, "Generate a %s document based on the following source specs.\n\n", targetType)

	if input.Context != "" {
		sb.WriteString("## Codebase Context (for grounding)\n\n")
		sb.WriteString("The following context describes the existing codebase and systems. ")
		sb.WriteString("Use this to ground your decisions in reality:\n\n")
		sb.WriteString(input.Context)
		sb.WriteString("\n\n")
	}

	for _, srcName := range rule.Sources {
		srcType := types.SpecType(srcName)
		content := input.Get(srcType)
		if content == "" {
			continue
		}
		fmt.Fprintf(&sb, "## Source: %s\n\n", srcType)
		sb.WriteString(content)
		sb.WriteString("\n\n")
		sources = append(sources, srcType)
	}

	if rule.Guidance != "" {
		sb.WriteString("## Synthesis Guidance\n\n")
		sb.WriteString(rule.Guidance)
		sb.WriteString("\n\n")
	}
	if rule.PromptContext != "" {
		sb.WriteString("## Additional Context\n\n")
		sb.WriteString(rule.PromptContext)
		sb.WriteString("\n\n")
	}

	writeCommonPromptTail(&sb, input, template)

	return sb.String(), sources
}

// buildLegacyPrompt is the original hardcoded per-type prompt builder,
// used only when no workflow rule is available for targetType (e.g. no
// --profile was used at init, or the resolved workflow doesn't define a
// synthesis rule for this type).
func (s *Synthesizer) buildLegacyPrompt(targetType types.SpecType, input SynthesisInput, template string) (string, []types.SpecType) {
	var sb strings.Builder
	var sources []types.SpecType

	mrd := input.Get(types.SpecTypeMRD)
	prd := input.Get(types.SpecTypePRD)
	uxd := input.Get(types.SpecTypeUXD)
	trd := input.Get(types.SpecTypeTRD)
	press := input.Get(types.SpecTypePress)
	faq := input.Get(types.SpecTypeFAQ)

	sb.WriteString("You are a technical writer synthesizing specification documents.\n\n")

	switch targetType {
	case types.SpecTypeTRD:
		sb.WriteString("Generate a Technical Requirements Document (TRD) based on the following source specs.\n\n")
		if input.Context != "" {
			sb.WriteString("## Codebase Context (for grounding)\n\n")
			sb.WriteString("The following context describes the existing codebase and systems. ")
			sb.WriteString("Use this to ground your technical decisions in reality:\n\n")
			sb.WriteString(input.Context)
			sb.WriteString("\n\n")
		}
		if mrd != "" {
			sb.WriteString("## Market Requirements Document (MRD)\n\n")
			sb.WriteString(mrd)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypeMRD)
		}
		if prd != "" {
			sb.WriteString("## Product Requirements Document (PRD)\n\n")
			sb.WriteString(prd)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypePRD)
		}
		if uxd != "" {
			sb.WriteString("## User Experience Design (UXD)\n\n")
			sb.WriteString(uxd)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypeUXD)
		}

	case types.SpecTypeIRD:
		sb.WriteString("Generate an Infrastructure Requirements Document (IRD) based on the following source specs.\n\n")
		if input.Context != "" {
			sb.WriteString("## Codebase Context (for grounding)\n\n")
			sb.WriteString("The following context describes the existing codebase and infrastructure. ")
			sb.WriteString("Use this to ground your infrastructure decisions in reality:\n\n")
			sb.WriteString(input.Context)
			sb.WriteString("\n\n")
		}
		if trd != "" {
			sb.WriteString("## Technical Requirements Document (TRD)\n\n")
			sb.WriteString(trd)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypeTRD)
		}

	case types.SpecTypePress:
		sb.WriteString("Generate a Press Release following Amazon's Working Backwards methodology.\n")
		sb.WriteString("This is the vision document that defines the customer experience.\n\n")
		if mrd != "" {
			sb.WriteString("## Market Requirements Document (MRD)\n\n")
			sb.WriteString(mrd)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypeMRD)
		}
		// PRD is optional - only include if available for enrichment
		if prd != "" {
			sb.WriteString("## Product Requirements Document (PRD) [Optional Context]\n\n")
			sb.WriteString(prd)
			sb.WriteString("\n\n")
		}

	case types.SpecTypeFAQ:
		sb.WriteString("Generate an FAQ document that anticipates customer and stakeholder questions.\n")
		sb.WriteString("Challenge assumptions in the Press Release and surface potential gaps.\n\n")
		if mrd != "" {
			sb.WriteString("## Market Requirements Document (MRD)\n\n")
			sb.WriteString(mrd)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypeMRD)
		}
		if press != "" {
			sb.WriteString("## Press Release (Vision)\n\n")
			sb.WriteString(press)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypePress)
		}

	case types.SpecTypePRD:
		sb.WriteString("Generate a Product Requirements Document (PRD) based on the Working Backwards artifacts.\n")
		sb.WriteString("The Press Release defines the vision; the FAQ clarifies scope and concerns.\n")
		sb.WriteString("Translate these into detailed, testable product requirements.\n\n")
		if mrd != "" {
			sb.WriteString("## Market Requirements Document (MRD)\n\n")
			sb.WriteString(mrd)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypeMRD)
		}
		if press != "" {
			sb.WriteString("## Press Release (Vision)\n\n")
			sb.WriteString(press)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypePress)
		}
		if faq != "" {
			sb.WriteString("## FAQ (Scope Clarification)\n\n")
			sb.WriteString(faq)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypeFAQ)
		}

	case types.SpecTypeNarrative1P:
		sb.WriteString("Generate a 1-page executive narrative summarizing the product.\n\n")
		if mrd != "" {
			sb.WriteString("## Market Requirements Document (MRD)\n\n")
			sb.WriteString(mrd)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypeMRD)
		}
		if prd != "" {
			sb.WriteString("## Product Requirements Document (PRD)\n\n")
			sb.WriteString(prd)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypePRD)
		}

	case types.SpecTypeNarrative6P:
		sb.WriteString("Generate a 6-page narrative document following AWS format.\n\n")
		if mrd != "" {
			sb.WriteString("## Market Requirements Document (MRD)\n\n")
			sb.WriteString(mrd)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypeMRD)
		}
		if prd != "" {
			sb.WriteString("## Product Requirements Document (PRD)\n\n")
			sb.WriteString(prd)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypePRD)
		}
		if uxd != "" {
			sb.WriteString("## User Experience Design (UXD)\n\n")
			sb.WriteString(uxd)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypeUXD)
		}

	case types.SpecTypeTPD:
		sb.WriteString("Generate a Test Plan Document (TPD) based on the following source specs.\n")
		sb.WriteString("Derive test cases from PRD acceptance criteria, technical tests from TRD, and user journey tests from UXD.\n\n")
		if prd != "" {
			sb.WriteString("## Product Requirements Document (PRD)\n\n")
			sb.WriteString("Use acceptance criteria to derive functional test cases:\n\n")
			sb.WriteString(prd)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypePRD)
		}
		if trd != "" {
			sb.WriteString("## Technical Requirements Document (TRD)\n\n")
			sb.WriteString("Use API design, data models, and NFRs to derive technical test cases:\n\n")
			sb.WriteString(trd)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypeTRD)
		}
		if uxd != "" {
			sb.WriteString("## User Experience Design (UXD)\n\n")
			sb.WriteString("Use user journeys to derive E2E and UAT test scenarios:\n\n")
			sb.WriteString(uxd)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypeUXD)
		}
	}

	writeCommonPromptTail(&sb, input, template)

	return sb.String(), sources
}

// writeCommonPromptTail appends the constitution, template, and standard
// instructions shared by both the rule-driven and legacy prompt builders.
func writeCommonPromptTail(sb *strings.Builder, input SynthesisInput, template string) {
	if input.Constitution != "" {
		sb.WriteString("## Constitution (Guiding Principles)\n\n")
		sb.WriteString(input.Constitution)
		sb.WriteString("\n\n")
	}

	sb.WriteString("## Template\n\n")
	sb.WriteString("Use the following template structure for your output:\n\n")
	sb.WriteString(template)
	sb.WriteString("\n\n")

	sb.WriteString("## Instructions\n\n")
	sb.WriteString("1. Fill in all sections of the template based on the source documents\n")
	sb.WriteString("2. Replace placeholder text with actual content\n")
	sb.WriteString("3. Ensure traceability to source requirements\n")
	sb.WriteString("4. Be specific and concrete, avoid vague statements\n")
	sb.WriteString("5. Output ONLY the completed document, no explanations\n")
}

// ResolveRule looks up a workflow's synthesis rule for a spec type,
// including inheritance — LoadedWorkflow.Workflow.Synthesis already
// reflects extends-resolved, per-key-merged rules, so this is a plain map
// lookup. ok is true whenever the workflow declares a rule at all,
// including an explicit empty-Sources rule (e.g. aws-one-way-door's press:
// sources: [] deliberately overrides an inherited press<-mrd rule to mark
// press as human-authored, not synthesizable). Callers use ok=false to
// mean "fall back to the legacy hardcoded tables", not ok=true-with-empty.
func ResolveRule(w *workflow.Workflow, specType types.SpecType) (rule *workflow.SynthesisRule, ok bool) {
	if w == nil || w.Synthesis == nil {
		return nil, false
	}
	rule, ok = w.Synthesis[string(specType)]
	return rule, ok
}

// CanSynthesizeWithRule reports whether specType can be synthesized given
// an optionally-resolved workflow rule from ResolveRule. A resolved rule
// with zero sources (deliberately human-authored, e.g. press under
// aws-one-way-door) is never synthesizable, regardless of the legacy
// table. With no resolved rule, falls back to CanSynthesize.
func CanSynthesizeWithRule(rule *workflow.SynthesisRule, ok bool, specType types.SpecType) bool {
	if ok {
		return rule != nil && len(rule.Sources) > 0
	}
	return CanSynthesize(specType)
}

// SourcesForRule returns the required source types given an optionally-
// resolved workflow rule from ResolveRule, falling back to
// RequiredSources when the workflow declares no rule for specType.
func SourcesForRule(rule *workflow.SynthesisRule, ok bool, specType types.SpecType) []types.SpecType {
	if ok {
		if rule == nil {
			return nil
		}
		sources := make([]types.SpecType, 0, len(rule.Sources))
		for _, s := range rule.Sources {
			sources = append(sources, types.SpecType(s))
		}
		return sources
	}
	return RequiredSources(specType)
}

// LoadWorkflowRule loads profileName via loader (falling back to
// workflows.DefaultLoader() when loader is nil) and resolves specType's
// synthesis rule from it. Returns a nil LoadedWorkflow and ok=false when
// profileName is empty or the profile fails to load, so callers can fall
// back to the legacy hardcoded tables via CanSynthesize/RequiredSources.
func LoadWorkflowRule(loader sws.Loader, profileName string, specType types.SpecType) (loaded *sws.LoadedWorkflow, rule *workflow.SynthesisRule, ok bool) {
	if profileName == "" {
		return nil, nil, false
	}
	if loader == nil {
		loader = sws.DefaultLoader()
	}
	w, err := loader.Load(profileName)
	if err != nil || w.Workflow == nil {
		return nil, nil, false
	}
	rule, ok = ResolveRule(w.Workflow, specType)
	return w, rule, ok
}

// RequiredSources returns the source spec types needed to synthesize a
// target type under the legacy pre-merge Working Backwards flow. Used only
// as a fallback when no workflow is configured, or the configured workflow
// declares no synthesis rule for targetType — callers with a resolved
// workflow should prefer SourcesForRule instead.
func RequiredSources(targetType types.SpecType) []types.SpecType {
	switch targetType {
	case types.SpecTypePress:
		// Working Backwards: Press comes first, from MRD only
		return []types.SpecType{types.SpecTypeMRD}
	case types.SpecTypeFAQ:
		// Working Backwards: FAQ clarifies scope using MRD + Press
		return []types.SpecType{types.SpecTypeMRD, types.SpecTypePress}
	case types.SpecTypePRD:
		// Working Backwards: PRD derived from MRD + Press + FAQ
		return []types.SpecType{types.SpecTypeMRD, types.SpecTypePress, types.SpecTypeFAQ}
	case types.SpecTypeTRD:
		return []types.SpecType{types.SpecTypeMRD, types.SpecTypePRD}
	case types.SpecTypeTPD:
		return []types.SpecType{types.SpecTypePRD, types.SpecTypeTRD, types.SpecTypeUXD}
	case types.SpecTypeIRD:
		return []types.SpecType{types.SpecTypeTRD}
	case types.SpecTypeNarrative1P:
		return []types.SpecType{types.SpecTypeMRD, types.SpecTypePRD}
	case types.SpecTypeNarrative6P:
		return []types.SpecType{types.SpecTypeMRD, types.SpecTypePRD}
	default:
		return nil
	}
}

// CanSynthesize returns whether a spec type can be synthesized under the
// legacy pre-merge Working Backwards flow. Used only as a fallback when no
// workflow is configured — callers with a resolved workflow should check
// whether its Synthesis map declares a rule for the type instead, since a
// workflow can make types synthesizable (e.g. enterprise's bmc) that this
// hardcoded list doesn't know about.
func CanSynthesize(specType types.SpecType) bool {
	switch specType {
	case types.SpecTypePRD, // PRD is synthesizable via Working Backwards flow
		types.SpecTypeTRD, types.SpecTypeTPD, types.SpecTypeIRD,
		types.SpecTypePress, types.SpecTypeFAQ,
		types.SpecTypeNarrative1P, types.SpecTypeNarrative6P:
		return true
	default:
		return false
	}
}
