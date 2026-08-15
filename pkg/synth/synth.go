// Package synth provides spec synthesis (generation) capabilities.
package synth

import (
	"context"
	"fmt"
	"strings"

	"github.com/ProductBuildersHQ/visionspec/pkg/templates"
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

// SynthesisInput contains the input documents for synthesis.
type SynthesisInput struct {
	MRD          string // Market Requirements Document
	PRD          string // Product Requirements Document
	UXD          string // User Experience Design
	TRD          string // Technical Requirements Document (for IRD synthesis)
	Constitution string // Project/org constitution
	Press        string // Press Release (for FAQ and PRD synthesis)
	FAQ          string // FAQ document (for PRD synthesis)
	Context      string // Aggregated context summary (for grounding)
}

// SynthesisResult contains the generated document and metadata.
type SynthesisResult struct {
	SpecType types.SpecType
	Content  string
	Sources  []types.SpecType // Source specs used
}

// Synthesize generates a spec of the given type from input documents.
func (s *Synthesizer) Synthesize(ctx context.Context, targetType types.SpecType, input SynthesisInput) (*SynthesisResult, error) {
	// Get template for target type
	tmpl, err := templates.Get(targetType)
	if err != nil {
		return nil, fmt.Errorf("no template for %s: %w", targetType, err)
	}

	// Build prompt based on target type
	prompt, sources := s.buildPrompt(targetType, input, tmpl.Content)

	// Call LLM
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

// buildPrompt constructs the synthesis prompt for a target spec type.
func (s *Synthesizer) buildPrompt(targetType types.SpecType, input SynthesisInput, template string) (string, []types.SpecType) {
	var sb strings.Builder
	var sources []types.SpecType

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
		if input.MRD != "" {
			sb.WriteString("## Market Requirements Document (MRD)\n\n")
			sb.WriteString(input.MRD)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypeMRD)
		}
		if input.PRD != "" {
			sb.WriteString("## Product Requirements Document (PRD)\n\n")
			sb.WriteString(input.PRD)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypePRD)
		}
		if input.UXD != "" {
			sb.WriteString("## User Experience Design (UXD)\n\n")
			sb.WriteString(input.UXD)
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
		if input.TRD != "" {
			sb.WriteString("## Technical Requirements Document (TRD)\n\n")
			sb.WriteString(input.TRD)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypeTRD)
		}

	case types.SpecTypePress:
		sb.WriteString("Generate a Press Release following Amazon's Working Backwards methodology.\n")
		sb.WriteString("This is the vision document that defines the customer experience.\n\n")
		if input.MRD != "" {
			sb.WriteString("## Market Requirements Document (MRD)\n\n")
			sb.WriteString(input.MRD)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypeMRD)
		}
		// PRD is optional - only include if available for enrichment
		if input.PRD != "" {
			sb.WriteString("## Product Requirements Document (PRD) [Optional Context]\n\n")
			sb.WriteString(input.PRD)
			sb.WriteString("\n\n")
		}

	case types.SpecTypeFAQ:
		sb.WriteString("Generate an FAQ document that anticipates customer and stakeholder questions.\n")
		sb.WriteString("Challenge assumptions in the Press Release and surface potential gaps.\n\n")
		if input.MRD != "" {
			sb.WriteString("## Market Requirements Document (MRD)\n\n")
			sb.WriteString(input.MRD)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypeMRD)
		}
		if input.Press != "" {
			sb.WriteString("## Press Release (Vision)\n\n")
			sb.WriteString(input.Press)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypePress)
		}

	case types.SpecTypePRD:
		sb.WriteString("Generate a Product Requirements Document (PRD) based on the Working Backwards artifacts.\n")
		sb.WriteString("The Press Release defines the vision; the FAQ clarifies scope and concerns.\n")
		sb.WriteString("Translate these into detailed, testable product requirements.\n\n")
		if input.MRD != "" {
			sb.WriteString("## Market Requirements Document (MRD)\n\n")
			sb.WriteString(input.MRD)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypeMRD)
		}
		if input.Press != "" {
			sb.WriteString("## Press Release (Vision)\n\n")
			sb.WriteString(input.Press)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypePress)
		}
		if input.FAQ != "" {
			sb.WriteString("## FAQ (Scope Clarification)\n\n")
			sb.WriteString(input.FAQ)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypeFAQ)
		}

	case types.SpecTypeNarrative1P:
		sb.WriteString("Generate a 1-page executive narrative summarizing the product.\n\n")
		if input.MRD != "" {
			sb.WriteString("## Market Requirements Document (MRD)\n\n")
			sb.WriteString(input.MRD)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypeMRD)
		}
		if input.PRD != "" {
			sb.WriteString("## Product Requirements Document (PRD)\n\n")
			sb.WriteString(input.PRD)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypePRD)
		}

	case types.SpecTypeNarrative6P:
		sb.WriteString("Generate a 6-page narrative document following AWS format.\n\n")
		if input.MRD != "" {
			sb.WriteString("## Market Requirements Document (MRD)\n\n")
			sb.WriteString(input.MRD)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypeMRD)
		}
		if input.PRD != "" {
			sb.WriteString("## Product Requirements Document (PRD)\n\n")
			sb.WriteString(input.PRD)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypePRD)
		}
		if input.UXD != "" {
			sb.WriteString("## User Experience Design (UXD)\n\n")
			sb.WriteString(input.UXD)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypeUXD)
		}

	case types.SpecTypeTPD:
		sb.WriteString("Generate a Test Plan Document (TPD) based on the following source specs.\n")
		sb.WriteString("Derive test cases from PRD acceptance criteria, technical tests from TRD, and user journey tests from UXD.\n\n")
		if input.PRD != "" {
			sb.WriteString("## Product Requirements Document (PRD)\n\n")
			sb.WriteString("Use acceptance criteria to derive functional test cases:\n\n")
			sb.WriteString(input.PRD)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypePRD)
		}
		if input.TRD != "" {
			sb.WriteString("## Technical Requirements Document (TRD)\n\n")
			sb.WriteString("Use API design, data models, and NFRs to derive technical test cases:\n\n")
			sb.WriteString(input.TRD)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypeTRD)
		}
		if input.UXD != "" {
			sb.WriteString("## User Experience Design (UXD)\n\n")
			sb.WriteString("Use user journeys to derive E2E and UAT test scenarios:\n\n")
			sb.WriteString(input.UXD)
			sb.WriteString("\n\n")
			sources = append(sources, types.SpecTypeUXD)
		}
	}

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

	return sb.String(), sources
}

// RequiredSources returns the source spec types needed to synthesize a target type.
// This implements the Working Backwards flow where:
//
//	MRD → Press → FAQ → PRD → (UXD) → TRD → IRD
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

// CanSynthesize returns whether a spec type can be synthesized.
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
