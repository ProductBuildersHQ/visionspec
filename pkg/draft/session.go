package draft

import (
	"fmt"

	"github.com/plexusone/multispec/pkg/templates"
	"github.com/plexusone/multispec/pkg/types"
)

// Session manages the state of a draft authoring session.
type Session struct {
	ProjectPath string
	SpecType    types.SpecType
	Draft       *Draft
}

// StartSession begins a new draft session.
// If a draft already exists, it resumes that draft instead of creating a new one.
func StartSession(projectPath string, specType types.SpecType) (*Session, error) {
	var draft *Draft
	var err error

	// Check if draft exists
	if Exists(projectPath, specType) {
		// Resume existing draft
		draft, err = Get(projectPath, specType)
		if err != nil {
			return nil, fmt.Errorf("failed to resume draft: %w", err)
		}
	} else {
		// Start new draft
		draft, err = Start(projectPath, specType)
		if err != nil {
			return nil, fmt.Errorf("failed to start draft: %w", err)
		}
	}

	return &Session{
		ProjectPath: projectPath,
		SpecType:    specType,
		Draft:       draft,
	}, nil
}

// ResumeSession resumes an existing draft session.
func ResumeSession(projectPath string, specType types.SpecType) (*Session, error) {
	draft, err := Get(projectPath, specType)
	if err != nil {
		return nil, err
	}

	return &Session{
		ProjectPath: projectPath,
		SpecType:    specType,
		Draft:       draft,
	}, nil
}

// Update saves new content to the session's draft.
func (s *Session) Update(content string) error {
	draft, err := Update(s.ProjectPath, s.SpecType, content)
	if err != nil {
		return err
	}
	s.Draft = draft
	return nil
}

// Finalize promotes the draft to final spec.
func (s *Session) Finalize() error {
	return Finalize(s.ProjectPath, s.SpecType)
}

// Discard removes the draft.
func (s *Session) Discard() error {
	return Discard(s.ProjectPath, s.SpecType)
}

// IsNew returns true if this is a newly created draft (version 1).
func (s *Session) IsNew() bool {
	return s.Draft.Metadata.Version == 1
}

// Version returns the current draft version.
func (s *Session) Version() int {
	return s.Draft.Metadata.Version
}

// Content returns the current draft content.
func (s *Session) Content() string {
	return s.Draft.Content
}

// RecordEval records an evaluation result.
func (s *Session) RecordEval(score float64, passed bool, findings int) error {
	if err := AddEvalResult(s.ProjectPath, s.SpecType, score, passed, findings); err != nil {
		return err
	}

	// Refresh draft
	draft, err := Get(s.ProjectPath, s.SpecType)
	if err != nil {
		return err
	}
	s.Draft = draft
	return nil
}

// EvalCount returns the number of evaluations performed.
func (s *Session) EvalCount() int {
	return len(s.Draft.Metadata.EvalHistory)
}

// LastEval returns the most recent evaluation summary, or nil if none.
func (s *Session) LastEval() *EvalSummary {
	if len(s.Draft.Metadata.EvalHistory) == 0 {
		return nil
	}
	return &s.Draft.Metadata.EvalHistory[len(s.Draft.Metadata.EvalHistory)-1]
}

// HasPassingEval returns true if the draft has passed evaluation at least once.
func (s *Session) HasPassingEval() bool {
	for _, eval := range s.Draft.Metadata.EvalHistory {
		if eval.Passed {
			return true
		}
	}
	return false
}

// SessionStatus represents the status of a draft session.
type SessionStatus string

const (
	StatusNew       SessionStatus = "new"       // Just started, no edits
	StatusEditing   SessionStatus = "editing"   // Has edits, not evaluated
	StatusEvaluated SessionStatus = "evaluated" // Has been evaluated
	StatusPassing   SessionStatus = "passing"   // Passed evaluation
	StatusFailing   SessionStatus = "failing"   // Failed evaluation
	StatusReady     SessionStatus = "ready"     // Ready to finalize
)

// Status returns the current session status.
func (s *Session) Status() SessionStatus {
	if s.EvalCount() == 0 {
		if s.Draft.Metadata.Version == 1 {
			return StatusNew
		}
		return StatusEditing
	}

	lastEval := s.LastEval()
	if lastEval.Passed {
		return StatusPassing
	}
	return StatusFailing
}

// Instructions returns authoring instructions based on current status.
func (s *Session) Instructions() string {
	switch s.Status() {
	case StatusNew:
		return "This is a new draft. Fill in each section following the template guidance."
	case StatusEditing:
		return "Continue editing the draft. Run evaluation when ready."
	case StatusEvaluated, StatusFailing:
		lastEval := s.LastEval()
		return fmt.Sprintf("Last evaluation: %.1f (findings: %d). Address feedback and run evaluation again.",
			lastEval.Score, lastEval.Findings)
	case StatusPassing:
		return "Draft has passed evaluation. Ready to finalize or continue improving."
	default:
		return ""
	}
}

// Summary returns a brief summary of the session state.
type SessionSummary struct {
	SpecType      types.SpecType `json:"spec_type"`
	SpecName      string         `json:"spec_name"`
	Version       int            `json:"version"`
	Status        SessionStatus  `json:"status"`
	EvalCount     int            `json:"eval_count"`
	LastScore     *float64       `json:"last_score,omitempty"`
	LastPassed    *bool          `json:"last_passed,omitempty"`
	HasTemplate   bool           `json:"has_template"`
	ContentLength int            `json:"content_length"`
}

// Summary returns a summary of the session.
func (s *Session) Summary() SessionSummary {
	summary := SessionSummary{
		SpecType:      s.SpecType,
		SpecName:      specTypeName(s.SpecType),
		Version:       s.Version(),
		Status:        s.Status(),
		EvalCount:     s.EvalCount(),
		HasTemplate:   templates.HasTemplate(s.SpecType),
		ContentLength: len(s.Draft.Content),
	}

	if lastEval := s.LastEval(); lastEval != nil {
		summary.LastScore = &lastEval.Score
		summary.LastPassed = &lastEval.Passed
	}

	return summary
}

// specTypeName returns a human-readable name for the spec type.
func specTypeName(st types.SpecType) string {
	switch st {
	case types.SpecTypeMRD:
		return "Market Requirements Document"
	case types.SpecTypePRD:
		return "Product Requirements Document"
	case types.SpecTypeUXD:
		return "User Experience Design"
	default:
		return string(st)
	}
}
