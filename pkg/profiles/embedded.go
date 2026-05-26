package profiles

import "embed"

// Embed all default profiles
//
// Stage-based profiles (by company maturity)
//
// Methodology profiles (by product development approach)
//
//go:embed default/0-1/profile.yaml default/0-1/templates/*.md default/0-1/rubrics/*.rubric.yaml
//go:embed default/startup/profile.yaml default/startup/templates/*.md default/startup/rubrics/*.rubric.yaml
//go:embed default/growth/profile.yaml default/growth/templates/*.md default/growth/rubrics/*.rubric.yaml
//go:embed default/enterprise/profile.yaml default/enterprise/templates/*.md default/enterprise/rubrics/*.rubric.yaml
//go:embed default/aws/profile.yaml default/aws/templates/*.md default/aws/rubrics/*.rubric.yaml
//go:embed default/google/profile.yaml default/google/templates/*.md default/google/rubrics/*.rubric.yaml
//go:embed default/stripe/profile.yaml default/stripe/templates/*.md default/stripe/rubrics/*.rubric.yaml
//go:embed default/lean-startup/profile.yaml default/lean-startup/templates/*.md default/lean-startup/rubrics/*.rubric.yaml
//go:embed default/design-thinking/profile.yaml default/design-thinking/templates/*.md default/design-thinking/rubrics/*.rubric.yaml
//go:embed default/jtbd/profile.yaml default/jtbd/templates/*.md default/jtbd/rubrics/*.rubric.yaml
var defaultProfiles embed.FS

// DefaultLoader returns a loader for built-in default profiles.
func DefaultLoader() Loader {
	return NewResolvingLoader(NewEmbedFSLoader(defaultProfiles, "default"))
}

// DefaultProfileNames returns the names of all default profiles.
//
// Stage-based profiles (by company maturity):
//   - 0-1: Pre-product-market-fit exploration
//   - startup: Early product development
//   - growth: Scaling product and team
//   - enterprise: Mature organization with compliance needs
//
// Methodology profiles (by product development approach):
//   - aws: Amazon Working Backwards with Leadership Principles
//   - google: Google Design Docs + RFC culture with OKRs
//   - lean-startup: Eric Ries' Build-Measure-Learn with validated learning
//   - design-thinking: Stanford d.school human-centered design
//   - jtbd: Clayton Christensen's Jobs-to-be-Done framework
var DefaultProfileNames = []string{
	// Stage-based
	"0-1", "startup", "growth", "enterprise",
	// Methodology-based
	"aws", "google", "stripe", "lean-startup", "design-thinking", "jtbd",
}

// IsDefaultProfile returns true if the name is a default profile.
func IsDefaultProfile(name string) bool {
	for _, n := range DefaultProfileNames {
		if n == name {
			return true
		}
	}
	return false
}
