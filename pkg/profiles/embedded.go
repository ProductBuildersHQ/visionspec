package profiles

import "embed"

// Embed all default profiles
//
// Stage-based profiles (by company maturity)
//
// Methodology profiles (by product development approach)
//
// Stage-based profiles
//
//go:embed default/0-1/profile.yaml default/0-1/templates/*.md default/0-1/rubrics/*.yaml
//go:embed default/startup/profile.yaml default/startup/templates/*.md default/startup/rubrics/*.yaml
//go:embed default/growth/profile.yaml default/growth/templates/*.md default/growth/rubrics/*.yaml
//go:embed default/enterprise/profile.yaml default/enterprise/templates/*.md default/enterprise/rubrics/*.yaml
//
// AWS methodology profiles
//
//go:embed default/aws-product/profile.yaml default/aws-product/templates/*.md default/aws-product/rubrics/*.yaml
//go:embed default/aws-feature/profile.yaml default/aws-feature/templates/*.md default/aws-feature/rubrics/*.yaml
//
// Other company methodology profiles
//
//go:embed default/google/profile.yaml default/google/templates/*.md default/google/rubrics/*.yaml
//go:embed default/stripe/profile.yaml default/stripe/templates/*.md default/stripe/rubrics/*.yaml
//go:embed default/lean-startup/profile.yaml default/lean-startup/templates/*.md default/lean-startup/rubrics/*.yaml
//go:embed default/design-thinking/profile.yaml default/design-thinking/templates/*.md default/design-thinking/rubrics/*.yaml
//go:embed default/jtbd/profile.yaml default/jtbd/templates/*.md default/jtbd/rubrics/*.yaml
//go:embed default/shapeup/profile.yaml default/shapeup/templates/*.md default/shapeup/rubrics/*.yaml
//go:embed default/continuous-discovery/profile.yaml default/continuous-discovery/templates/*.md default/continuous-discovery/rubrics/*.yaml
//
// Big Tech profiles (combines 10 methodologies)
//
//go:embed default/big-tech/profile.yaml
//go:embed default/big-tech-product/profile.yaml default/big-tech-product/templates/*.md default/big-tech-product/rubrics/*.yaml
//go:embed default/big-tech-feature/profile.yaml default/big-tech-feature/templates/*.md default/big-tech-feature/rubrics/*.yaml
//
// Big Tech Essentials profiles (simplified 3-company synthesis)
//
//go:embed default/big-tech-essentials/profile.yaml
//go:embed default/big-tech-essentials-product/profile.yaml default/big-tech-essentials-product/rubrics/*.yaml
//go:embed default/big-tech-essentials-feature/profile.yaml default/big-tech-essentials-feature/rubrics/*.yaml
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
//   - aws-product: Amazon Working Backwards for new products
//   - aws-feature: Amazon Working Backwards for features
//   - google: Google Design Docs + RFC culture with OKRs
//   - stripe: Stripe API-first development
//   - lean-startup: Eric Ries' Build-Measure-Learn with validated learning
//   - design-thinking: Stanford d.school human-centered design
//   - jtbd: Clayton Christensen's Jobs-to-be-Done framework
//   - shapeup: Basecamp Shape Up methodology
//   - continuous-discovery: Teresa Torres Continuous Discovery Habits
//
// Big Tech profiles (combines 10 methodologies):
//   - big-tech: Abstract base for big-tech-product and big-tech-feature
//   - big-tech-product: Full 10-company synthesis for new products
//   - big-tech-feature: Full 10-company synthesis for features
//
// Big Tech Essentials profiles (simplified 3-company synthesis):
//   - big-tech-essentials: Abstract base (Amazon + Google + Stripe only)
//   - big-tech-essentials-product: Essentials for new products
//   - big-tech-essentials-feature: Essentials for features
var DefaultProfileNames = []string{
	// Stage-based
	"0-1", "startup", "growth", "enterprise",
	// AWS methodology
	"aws-product", "aws-feature",
	// Other company methodology
	"google", "stripe", "lean-startup", "design-thinking", "jtbd", "shapeup", "continuous-discovery",
	// Big Tech (full synthesis)
	"big-tech", "big-tech-product", "big-tech-feature",
	// Big Tech Essentials (simplified)
	"big-tech-essentials", "big-tech-essentials-product", "big-tech-essentials-feature",
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
