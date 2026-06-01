// Package eval provides evaluation orchestration for spec documents.
package eval

import (
	"fmt"
	"strings"
	"time"

	"github.com/plexusone/structured-evaluation/claims"
)

// ToClaimsReport extracts claims from evaluation findings.
// Each finding becomes a claim with internal validation based on the evaluation.
func (r *Result) ToClaimsReport(document string) *claims.ClaimsReport {
	report := claims.NewClaimsReport(document)
	report.Metadata.GeneratedBy = "visionspec eval"
	report.Metadata.GeneratedAt = r.Timestamp

	// Each finding becomes a claim
	for i, f := range r.Findings {
		claimID := fmt.Sprintf("finding-%d", i+1)

		// Map finding severity to claim category
		category := mapSeverityToCategory(f.Severity)

		claim := claims.NewClaim(
			claimID,
			f.Title+": "+f.Description,
			category,
			claims.Location{Section: f.Category},
		)

		// Set internal validation - findings are validated by the LLM evaluation
		validation := claims.NewInternalValidation(
			claims.MethodObservation,
			"",    // No specific evidence path for LLM findings
			false, // LLM evaluations are not fully reproducible
		)
		validation.Internal.ValidatedBy = fmt.Sprintf("%s via %s", r.Judge.Provider, r.Judge.Model)
		validation.Internal.ValidatedAt = r.Timestamp

		claim.SetValidation(validation)

		// Set verdict based on severity
		verdict, rationale := mapSeverityToVerdict(f.Severity)
		claim.SetVerdict(verdict, rationale)

		report.AddClaim(*claim)
	}

	// Add summary claim
	if r.Summary != "" {
		summaryClaim := claims.NewClaim(
			"summary",
			r.Summary,
			claims.ClaimGuidance,
			claims.Location{Section: "summary"},
		)
		summaryClaim.SetValidation(claims.NewInternalValidation(
			claims.MethodObservation,
			"",
			false,
		))
		summaryClaim.SetVerdict(claims.VerdictVerified, "Summary derived from LLM evaluation")
		report.AddClaim(*summaryClaim)
	}

	report.Finalize()
	return report
}

// mapSeverityToCategory maps finding severity to claim category.
func mapSeverityToCategory(severity string) claims.ClaimCategory {
	switch strings.ToLower(severity) {
	case "critical", "high":
		return claims.ClaimTechnicalFinding
	case "medium", "low":
		return claims.ClaimGuidance
	case "info":
		return claims.ClaimMetadata
	default:
		return claims.ClaimTechnicalFinding
	}
}

// mapSeverityToVerdict maps finding severity to claim verdict.
func mapSeverityToVerdict(severity string) (claims.Verdict, string) {
	switch strings.ToLower(severity) {
	case "critical":
		return claims.VerdictNeedsReview, "Critical finding requires immediate attention"
	case "high":
		return claims.VerdictNeedsReview, "High severity finding should be addressed"
	case "medium":
		return claims.VerdictVerified, "Medium finding identified in evaluation"
	case "low":
		return claims.VerdictVerified, "Low severity observation"
	case "info":
		return claims.VerdictVerified, "Informational observation"
	default:
		return claims.VerdictNeedsReview, "Unknown severity, review required"
	}
}

// ClaimsFromCategoryResults creates claims from category results.
func ClaimsFromCategoryResults(categories []CategoryResult, timestamp time.Time) []claims.Claim {
	result := make([]claims.Claim, 0, len(categories))

	for i, cat := range categories {
		claimID := fmt.Sprintf("category-%d", i+1)

		claim := claims.NewClaim(
			claimID,
			fmt.Sprintf("%s scored %.1f/10: %s", cat.Name, cat.Score, cat.Explanation),
			claims.ClaimTechnicalFinding,
			claims.Location{Section: cat.ID},
		)

		validation := claims.NewInternalValidation(
			claims.MethodObservation,
			"",
			false,
		)
		claim.SetValidation(validation)

		// Verdict based on score
		if cat.Score >= 7.0 {
			claim.SetVerdict(claims.VerdictVerified, "Category passed evaluation threshold")
		} else if cat.Score >= 5.0 {
			claim.SetVerdict(claims.VerdictNeedsReview, "Category partially met criteria")
		} else {
			claim.SetVerdict(claims.VerdictRejected, "Category failed to meet minimum threshold")
		}

		result = append(result, *claim)
	}

	return result
}
