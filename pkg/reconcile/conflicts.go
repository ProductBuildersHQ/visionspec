// Package reconcile provides spec reconciliation capabilities.
package reconcile

import (
	"regexp"
	"strings"

	"github.com/plexusone/multispec/pkg/types"
)

// ConflictType represents the type of conflict between specs.
type ConflictType string

const (
	ConflictTypeRequirement ConflictType = "requirement"
	ConflictTypeConstraint  ConflictType = "constraint"
	ConflictTypeTradeoff    ConflictType = "tradeoff"
	ConflictTypeMissing     ConflictType = "missing"
)

// ConflictSeverity represents the severity of a conflict.
type ConflictSeverity string

const (
	SeverityHigh   ConflictSeverity = "high"
	SeverityMedium ConflictSeverity = "medium"
	SeverityLow    ConflictSeverity = "low"
)

// ConflictDetector analyzes specs for potential conflicts.
type ConflictDetector struct {
	// Patterns for detecting potential issues
	performancePatterns  *regexp.Regexp
	securityPatterns     *regexp.Regexp
	scalabilityPatterns  *regexp.Regexp
	timelinePatterns     *regexp.Regexp
	budgetPatterns       *regexp.Regexp
	priorityPatterns     *regexp.Regexp
	requirementIDPattern *regexp.Regexp
}

// NewConflictDetector creates a new conflict detector.
func NewConflictDetector() *ConflictDetector {
	return &ConflictDetector{
		performancePatterns:  regexp.MustCompile(`(?i)(latency|response time|throughput|performance|<\s*\d+\s*ms)`),
		securityPatterns:     regexp.MustCompile(`(?i)(security|authentication|authorization|encryption|GDPR|HIPAA|SOC|compliance)`),
		scalabilityPatterns:  regexp.MustCompile(`(?i)(scale|concurrent|users|traffic|load|horizontal|vertical)`),
		timelinePatterns:     regexp.MustCompile(`(?i)(deadline|timeline|phase \d|sprint|milestone|by Q\d|week \d)`),
		budgetPatterns:       regexp.MustCompile(`(?i)(budget|cost|resource|headcount|team size)`),
		priorityPatterns:     regexp.MustCompile(`(?i)(P0|P1|P2|must have|should have|could have|won't have|critical|high priority)`),
		requirementIDPattern: regexp.MustCompile(`(?i)(REQ-\d+|US-\d+|FR-\d+|NFR-\d+)`),
	}
}

// DetectedConflict represents a potential conflict found during analysis.
type DetectedConflict struct {
	Conflict
	Confidence float64 `json:"confidence"` // 0.0-1.0
}

// DetectConflicts analyzes input specs and returns potential conflicts.
func (d *ConflictDetector) DetectConflicts(input ReconcileInput) []DetectedConflict {
	var conflicts []DetectedConflict
	conflictID := 1

	// Check for performance vs. feature conflicts
	if input.PRD != "" && input.TRD != "" {
		if c := d.detectPerformanceConflicts(input, &conflictID); c != nil {
			conflicts = append(conflicts, *c)
		}
	}

	// Check for security vs. usability conflicts
	if input.UXD != "" && input.TRD != "" {
		if c := d.detectSecurityUsabilityConflicts(input, &conflictID); c != nil {
			conflicts = append(conflicts, *c)
		}
	}

	// Check for scope vs. timeline conflicts
	if input.PRD != "" && input.MRD != "" {
		if c := d.detectScopeTimelineConflicts(input, &conflictID); c != nil {
			conflicts = append(conflicts, *c)
		}
	}

	// Check for missing traceability
	missingTrace := d.detectMissingTraceability(input, &conflictID)
	conflicts = append(conflicts, missingTrace...)

	// Check for constraint conflicts
	constraintConflicts := d.detectConstraintConflicts(input, &conflictID)
	conflicts = append(conflicts, constraintConflicts...)

	return conflicts
}

// detectPerformanceConflicts checks for conflicts between feature requirements and performance constraints.
func (d *ConflictDetector) detectPerformanceConflicts(input ReconcileInput, id *int) *DetectedConflict {
	prdHasPerf := d.performancePatterns.MatchString(input.PRD)
	trdHasPerf := d.performancePatterns.MatchString(input.TRD)

	// Look for specific performance numbers that might conflict
	prdNumbers := extractNumbers(input.PRD)
	trdNumbers := extractNumbers(input.TRD)

	if prdHasPerf && trdHasPerf && len(prdNumbers) > 0 && len(trdNumbers) > 0 {
		// Check if there are different performance targets
		if !numbersMatch(prdNumbers, trdNumbers) {
			c := DetectedConflict{
				Conflict: Conflict{
					ID:          generateConflictID(id),
					Type:        string(ConflictTypeConstraint),
					Description: "Performance requirements in PRD may conflict with technical constraints in TRD. Different latency/throughput targets specified.",
					Sources:     []types.SpecType{types.SpecTypePRD, types.SpecTypeTRD},
					Severity:    string(SeverityMedium),
				},
				Confidence: 0.7,
			}
			return &c
		}
	}

	return nil
}

// detectSecurityUsabilityConflicts checks for conflicts between security and UX.
func (d *ConflictDetector) detectSecurityUsabilityConflicts(input ReconcileInput, id *int) *DetectedConflict {
	trdHasSecurity := d.securityPatterns.MatchString(input.TRD)
	uxdMentionsFriction := strings.Contains(strings.ToLower(input.UXD), "frictionless") ||
		strings.Contains(strings.ToLower(input.UXD), "seamless") ||
		strings.Contains(strings.ToLower(input.UXD), "one-click")

	if trdHasSecurity && uxdMentionsFriction {
		c := DetectedConflict{
			Conflict: Conflict{
				ID:          generateConflictID(id),
				Type:        string(ConflictTypeTradeoff),
				Description: "Security requirements in TRD may conflict with frictionless UX goals. Authentication/authorization adds user friction.",
				Sources:     []types.SpecType{types.SpecTypeUXD, types.SpecTypeTRD},
				Severity:    string(SeverityMedium),
			},
			Confidence: 0.6,
		}
		return &c
	}

	return nil
}

// detectScopeTimelineConflicts checks for scope vs timeline conflicts.
func (d *ConflictDetector) detectScopeTimelineConflicts(input ReconcileInput, id *int) *DetectedConflict {
	// Count requirements in PRD
	reqCount := len(d.requirementIDPattern.FindAllString(input.PRD, -1))
	if reqCount == 0 {
		// Count bullet points as proxy for requirements
		reqCount = strings.Count(input.PRD, "- ") + strings.Count(input.PRD, "* ")
	}

	// Check if MRD has aggressive timeline
	hasAggressiveTimeline := strings.Contains(strings.ToLower(input.MRD), "aggressive") ||
		strings.Contains(strings.ToLower(input.MRD), "tight deadline") ||
		strings.Contains(strings.ToLower(input.MRD), "asap") ||
		d.timelinePatterns.MatchString(input.MRD)

	// Large scope + aggressive timeline = potential conflict
	if reqCount > 20 && hasAggressiveTimeline {
		c := DetectedConflict{
			Conflict: Conflict{
				ID:          generateConflictID(id),
				Type:        string(ConflictTypeTradeoff),
				Description: "Large feature scope in PRD may conflict with timeline expectations in MRD. Consider phased delivery.",
				Sources:     []types.SpecType{types.SpecTypeMRD, types.SpecTypePRD},
				Severity:    string(SeverityHigh),
			},
			Confidence: 0.65,
		}
		return &c
	}

	return nil
}

// detectMissingTraceability checks for requirements without technical coverage.
func (d *ConflictDetector) detectMissingTraceability(input ReconcileInput, id *int) []DetectedConflict {
	var conflicts []DetectedConflict

	// Extract requirement IDs from PRD
	prdReqs := d.requirementIDPattern.FindAllString(input.PRD, -1)

	// Check if TRD references those requirements
	if len(prdReqs) > 0 && input.TRD != "" {
		missingInTRD := []string{}
		for _, req := range prdReqs {
			if !strings.Contains(input.TRD, req) {
				missingInTRD = append(missingInTRD, req)
			}
		}

		if len(missingInTRD) > 0 {
			c := DetectedConflict{
				Conflict: Conflict{
					ID:          generateConflictID(id),
					Type:        string(ConflictTypeMissing),
					Description: "Requirements not traced to TRD: " + strings.Join(missingInTRD[:min(5, len(missingInTRD))], ", "),
					Sources:     []types.SpecType{types.SpecTypePRD, types.SpecTypeTRD},
					Severity:    string(SeverityMedium),
				},
				Confidence: 0.9,
			}
			conflicts = append(conflicts, c)
		}
	}

	return conflicts
}

// detectConstraintConflicts checks for explicit constraint conflicts.
func (d *ConflictDetector) detectConstraintConflicts(input ReconcileInput, id *int) []DetectedConflict {
	var conflicts []DetectedConflict

	// Check scalability constraints
	if input.TRD != "" && input.IRD != "" {
		trdScale := d.scalabilityPatterns.FindAllString(input.TRD, -1)
		irdScale := d.scalabilityPatterns.FindAllString(input.IRD, -1)

		if len(trdScale) > 0 && len(irdScale) > 0 {
			// Look for mismatched scale expectations
			trdText := strings.ToLower(input.TRD)
			irdText := strings.ToLower(input.IRD)

			// Check for horizontal vs vertical scaling mismatch
			trdHorizontal := strings.Contains(trdText, "horizontal") || strings.Contains(trdText, "auto-scale")
			irdVertical := strings.Contains(irdText, "vertical") && !strings.Contains(irdText, "horizontal")

			if trdHorizontal && irdVertical {
				c := DetectedConflict{
					Conflict: Conflict{
						ID:          generateConflictID(id),
						Type:        string(ConflictTypeConstraint),
						Description: "TRD assumes horizontal scaling but IRD describes vertical scaling approach.",
						Sources:     []types.SpecType{types.SpecTypeTRD, types.SpecTypeIRD},
						Severity:    string(SeverityHigh),
					},
					Confidence: 0.8,
				}
				conflicts = append(conflicts, c)
			}
		}
	}

	return conflicts
}

// generateConflictID generates a unique conflict ID.
func generateConflictID(counter *int) string {
	id := *counter
	*counter++
	return "CONFLICT-" + padNumber(id, 3)
}

// padNumber pads a number with leading zeros.
func padNumber(n, width int) string {
	s := ""
	for i := 0; i < width; i++ {
		s = "0" + s
	}
	ns := s + string(rune('0'+n%10))
	for n >= 10 {
		n /= 10
		ns = string(rune('0'+n%10)) + ns
	}
	// Simple padding
	result := make([]byte, width)
	for i := range result {
		result[i] = '0'
	}
	numStr := []byte(ns)
	copy(result[width-len(numStr):], numStr)
	return string(result)
}

// extractNumbers extracts numeric values from performance-related text.
func extractNumbers(text string) []int {
	numberPattern := regexp.MustCompile(`(\d+)\s*(ms|seconds?|users?|req/s|rps|qps)`)
	matches := numberPattern.FindAllStringSubmatch(text, -1)
	var numbers []int
	for _, m := range matches {
		if len(m) > 1 {
			var n int
			for _, c := range m[1] {
				n = n*10 + int(c-'0')
			}
			numbers = append(numbers, n)
		}
	}
	return numbers
}

// numbersMatch checks if two number sets have overlap.
func numbersMatch(a, b []int) bool {
	for _, na := range a {
		for _, nb := range b {
			if na == nb {
				return true
			}
		}
	}
	return false
}

// ParseConflictsFromOutput extracts resolved conflicts from LLM output.
func ParseConflictsFromOutput(content string) []Conflict {
	var conflicts []Conflict

	// Look for Decision Log section
	decisionLogPattern := regexp.MustCompile(`(?s)## Decision Log\s*\n(.*?)(?:\n##|\z)`)
	match := decisionLogPattern.FindStringSubmatch(content)

	if len(match) > 1 {
		decisionLog := match[1]

		// Parse individual decisions/conflicts
		decisionPattern := regexp.MustCompile(`(?m)^[-*]\s*\*\*([^*]+)\*\*[:\s]*(.+)$`)
		decisions := decisionPattern.FindAllStringSubmatch(decisionLog, -1)

		for i, d := range decisions {
			if len(d) > 2 {
				conflicts = append(conflicts, Conflict{
					ID:          "RESOLVED-" + padNumber(i+1, 3),
					Type:        string(ConflictTypeTradeoff),
					Description: strings.TrimSpace(d[1]),
					Resolution:  strings.TrimSpace(d[2]),
					Severity:    string(SeverityMedium), // Resolved conflicts are medium by default
				})
			}
		}
	}

	return conflicts
}
