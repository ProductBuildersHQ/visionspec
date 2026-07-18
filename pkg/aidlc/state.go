package aidlc

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// State represents the parsed aidlc-state.md file content.
type State struct {
	// CurrentPhase is the current workflow phase.
	CurrentPhase Phase `json:"current_phase" yaml:"current_phase"`

	// CurrentDocument is the document currently being worked on.
	CurrentDocument DocType `json:"current_document,omitempty" yaml:"current_document,omitempty"`

	// CompletedDocs lists completed document types.
	CompletedDocs []DocType `json:"completed_docs" yaml:"completed_docs"`

	// PendingDocs lists pending document types.
	PendingDocs []DocType `json:"pending_docs" yaml:"pending_docs"`

	// InProgressDocs lists in-progress document types.
	InProgressDocs []DocType `json:"in_progress_docs,omitempty" yaml:"in_progress_docs,omitempty"`

	// DocumentScores maps document types to their quality scores.
	DocumentScores map[DocType]*QualityScore `json:"document_scores,omitempty" yaml:"document_scores,omitempty"`

	// PhaseProgress tracks progress per phase (0.0-1.0).
	PhaseProgress map[Phase]float64 `json:"phase_progress,omitempty" yaml:"phase_progress,omitempty"`

	// LastUpdated is when the state was last modified.
	LastUpdated time.Time `json:"last_updated" yaml:"last_updated"`

	// WorkflowStarted is when the workflow began.
	WorkflowStarted time.Time `json:"workflow_started,omitempty" yaml:"workflow_started,omitempty"`

	// Metadata contains additional parsed state data.
	Metadata map[string]any `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// OverallProgress returns the overall workflow progress (0.0-1.0).
func (s *State) OverallProgress() float64 {
	if len(s.CompletedDocs) == 0 && len(s.PendingDocs) == 0 {
		return 0.0
	}
	total := len(s.CompletedDocs) + len(s.PendingDocs) + len(s.InProgressDocs)
	if total == 0 {
		return 0.0
	}
	return float64(len(s.CompletedDocs)) / float64(total)
}

// IsPhaseComplete returns whether a phase is fully complete.
func (s *State) IsPhaseComplete(phase Phase) bool {
	phaseDocs := DocTypesByPhase(phase)
	for _, doc := range phaseDocs {
		found := false
		for _, completed := range s.CompletedDocs {
			if completed == doc {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// NextDocument returns the next document to work on, or empty if none.
func (s *State) NextDocument() DocType {
	if len(s.InProgressDocs) > 0 {
		return s.InProgressDocs[0]
	}
	if len(s.PendingDocs) > 0 {
		return s.PendingDocs[0]
	}
	return ""
}

// StateParser parses aidlc-state.md files.
type StateParser struct {
	// patterns for parsing markdown state files
	phasePattern    *regexp.Regexp
	docPattern      *regexp.Regexp
	scorePattern    *regexp.Regexp
	statusPattern   *regexp.Regexp
	progressPattern *regexp.Regexp
}

// NewStateParser creates a new state parser.
func NewStateParser() *StateParser {
	return &StateParser{
		phasePattern:    regexp.MustCompile(`(?i)^##?\s*(Inception|Construction|Operations)\s*Phase`),
		docPattern:      regexp.MustCompile(`(?i)^\s*[-*]\s*\[([xX ])\]\s*(.+?)(?:\s*[-:]\s*(.+))?$`),
		scorePattern:    regexp.MustCompile(`(?i)(EXCELLENT|GOOD|NEEDS_IMPROVEMENT|POOR)\s*(?:\(([0-9.]+)\))?`),
		statusPattern:   regexp.MustCompile(`(?i)Status:\s*(.+)`),
		progressPattern: regexp.MustCompile(`(?i)Progress:\s*([0-9.]+)%?`),
	}
}

// newState creates a new empty state with the given last updated time.
func newState(lastUpdated time.Time) *State {
	return &State{
		CurrentPhase:   PhaseInception,
		CompletedDocs:  make([]DocType, 0),
		PendingDocs:    make([]DocType, 0),
		InProgressDocs: make([]DocType, 0),
		DocumentScores: make(map[DocType]*QualityScore),
		PhaseProgress:  make(map[Phase]float64),
		LastUpdated:    lastUpdated,
		Metadata:       make(map[string]any),
	}
}

// parseScanner parses state from a scanner, populating the provided state.
func (p *StateParser) parseScanner(scanner *bufio.Scanner, state *State) error {
	var currentPhase Phase
	var inProgressSection bool

	for scanner.Scan() {
		line := scanner.Text()

		// Check for phase headers
		if matches := p.phasePattern.FindStringSubmatch(line); matches != nil {
			phaseName := strings.ToLower(matches[1])
			switch phaseName {
			case "inception":
				currentPhase = PhaseInception
			case "construction":
				currentPhase = PhaseConstruction
			case "operations":
				currentPhase = PhaseOperations
			}
			continue
		}

		// Check for "In Progress" section
		if strings.Contains(strings.ToLower(line), "in progress") {
			inProgressSection = true
			continue
		}
		if strings.HasPrefix(line, "##") || strings.HasPrefix(line, "---") {
			inProgressSection = false
		}

		// Parse document checklist items
		if matches := p.docPattern.FindStringSubmatch(line); matches != nil {
			isChecked := strings.ToLower(matches[1]) == "x"
			docName := strings.TrimSpace(matches[2])
			extra := strings.TrimSpace(matches[3])

			docType := p.inferDocType(docName, currentPhase)
			if docType == "" {
				continue
			}

			if inProgressSection {
				state.InProgressDocs = append(state.InProgressDocs, docType)
				if state.CurrentDocument == "" {
					state.CurrentDocument = docType
				}
			} else if isChecked {
				state.CompletedDocs = append(state.CompletedDocs, docType)
				// Parse score if present
				if scoreMatches := p.scorePattern.FindStringSubmatch(extra); scoreMatches != nil {
					rating := QualityRating(strings.ToUpper(scoreMatches[1]))
					var score float64
					if scoreMatches[2] != "" {
						score, _ = strconv.ParseFloat(scoreMatches[2], 64)
					} else {
						score = rating.Score()
					}
					state.DocumentScores[docType] = &QualityScore{
						Rating:      rating,
						Score:       score,
						EvaluatedAt: state.LastUpdated,
					}
				}
			} else {
				state.PendingDocs = append(state.PendingDocs, docType)
			}
		}

		// Parse overall status
		if matches := p.statusPattern.FindStringSubmatch(line); matches != nil {
			status := strings.TrimSpace(matches[1])
			state.Metadata["status"] = status
		}

		// Parse progress
		if matches := p.progressPattern.FindStringSubmatch(line); matches != nil {
			progress, _ := strconv.ParseFloat(matches[1], 64)
			if progress > 1.0 {
				progress /= 100.0 // Convert percentage to decimal
			}
			if currentPhase != "" {
				state.PhaseProgress[currentPhase] = progress
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Set current phase based on progress
	state.CurrentPhase = p.determineCurrentPhase(state)

	return nil
}

// ParseFile parses an aidlc-state.md file.
func (p *StateParser) ParseFile(path string) (*State, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open state file: %w", err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat state file: %w", err)
	}

	state := newState(info.ModTime())

	scanner := bufio.NewScanner(f)
	if err := p.parseScanner(scanner, state); err != nil {
		return nil, fmt.Errorf("read state file: %w", err)
	}

	return state, nil
}

// inferDocType infers the DocType from a document name.
func (p *StateParser) inferDocType(name string, phase Phase) DocType {
	nameLower := strings.ToLower(name)

	// Direct mappings
	docMappings := map[string]DocType{
		"vision":          DocVisionDocument,
		"vision document": DocVisionDocument,
		"requirements":    DocRequirementsSpec,
		"technical spec":  DocTechnicalSpec,
		"architecture":    DocArchitectureSpec,
		"implementation":  DocImplementationPlan,
		"test plan":       DocTestPlan,
		"integration":     DocIntegrationPlan,
		"security":        DocSecurityReview,
		"runbook":         DocRunbook,
		"monitoring":      DocMonitoringPlan,
		"disaster":        DocDisasterPlan,
		"slo":             DocSLODocument,
	}

	for key, docType := range docMappings {
		if strings.Contains(nameLower, key) {
			return docType
		}
	}

	// Fall back to phase-based inference
	if phase != "" {
		phaseDocs := DocTypesByPhase(phase)
		if len(phaseDocs) > 0 {
			// Return first doc type if name contains "document" or "spec"
			if strings.Contains(nameLower, "document") || strings.Contains(nameLower, "spec") {
				return phaseDocs[0]
			}
		}
	}

	return ""
}

// determineCurrentPhase determines the current phase based on state.
func (p *StateParser) determineCurrentPhase(state *State) Phase {
	// If we have in-progress docs, use their phase
	if len(state.InProgressDocs) > 0 {
		return state.InProgressDocs[0].Phase()
	}

	// Check phases in order for incomplete docs
	for _, phase := range AllPhases() {
		if !state.IsPhaseComplete(phase) {
			return phase
		}
	}

	// All complete - return operations
	return PhaseOperations
}

// ParseString parses state from a markdown string.
func (p *StateParser) ParseString(content string) (*State, error) {
	state := newState(time.Now())

	scanner := bufio.NewScanner(strings.NewReader(content))
	if err := p.parseScanner(scanner, state); err != nil {
		return nil, fmt.Errorf("parse state string: %w", err)
	}

	return state, nil
}

// DefaultStateParser is the default state parser instance.
var DefaultStateParser = NewStateParser()

// ParseStateFile parses an aidlc-state.md file using the default parser.
func ParseStateFile(path string) (*State, error) {
	return DefaultStateParser.ParseFile(path)
}

// ParseStateString parses state from a markdown string using the default parser.
func ParseStateString(content string) (*State, error) {
	return DefaultStateParser.ParseString(content)
}
