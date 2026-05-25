// Example: Post-PMF Enterprise CLI using profiles
//
// This example demonstrates using the built-in "enterprise" profile
// for post-product-market-fit enterprises building SaaS products.
//
// Build with:
//
//	go build -o enterprise-spec ./examples/post-pmf-enterprise
package main

import (
	"fmt"
	"os"

	"github.com/ProductBuildersHQ/visionspec/pkg/cli"
	"github.com/ProductBuildersHQ/visionspec/pkg/profiles"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "enterprise-spec",
		Short: "Enterprise-grade spec management for SaaS products",
		Long: `enterprise-spec uses the "enterprise" profile for post-PMF companies.

All source specs (MRD, PRD, UXD), TRD, and GTM docs are required.
Includes security and WCAG 2.1 AA accessibility requirements.`,
		Version: "1.0.0",
	}

	// Load the enterprise profile
	profile, err := profiles.DefaultLoader().Load("enterprise")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading profile: %v\n", err)
		os.Exit(1)
	}

	// Create config from profile
	cfg := cli.ConfigFromProfile(profile)
	cfg.Version = "1.0.0-enterprise"

	// Add visionspec commands
	cli.AddCommandsTo(root, cfg)

	// Add enterprise-specific commands
	root.AddCommand(complianceCmd())
	root.AddCommand(securityCmd())

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

// complianceCmd provides compliance checking commands.
func complianceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compliance",
		Short: "Compliance management",
		Long:  `Commands for managing compliance requirements.`,
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "check",
		Short: "Check compliance requirements",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Enterprise Compliance Checklist:")
			fmt.Println()
			fmt.Println("Security Requirements:")
			fmt.Println("  [ ] Authentication requirements defined in PRD")
			fmt.Println("  [ ] Authorization model documented in TRD")
			fmt.Println("  [ ] Data encryption requirements specified")
			fmt.Println("  [ ] Audit logging requirements documented")
			fmt.Println()
			fmt.Println("Accessibility Requirements:")
			fmt.Println("  [ ] WCAG 2.1 AA compliance documented in UXD")
			fmt.Println("  [ ] Keyboard navigation specified")
			fmt.Println("  [ ] Screen reader support defined")
			fmt.Println("  [ ] Color contrast requirements specified")
			fmt.Println()
			fmt.Println("Data Protection:")
			fmt.Println("  [ ] PII handling documented")
			fmt.Println("  [ ] Data retention policy specified")
			fmt.Println("  [ ] GDPR/CCPA considerations addressed")
			fmt.Println()
			fmt.Println("Platform Requirements:")
			fmt.Println("  [ ] Web browser support matrix in UXD")
			fmt.Println("  [ ] Mobile OS version requirements defined")
			fmt.Println("  [ ] API versioning strategy in TRD")
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "frameworks",
		Short: "List supported compliance frameworks",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Supported Compliance Frameworks:")
			fmt.Println()
			fmt.Println("  SOC 2 Type II")
			fmt.Println("    - Security, availability, confidentiality controls")
			fmt.Println("    - Integrated into TRD templates")
			fmt.Println()
			fmt.Println("  GDPR")
			fmt.Println("    - Data protection requirements in PRD")
			fmt.Println("    - Privacy-by-design in UXD")
			fmt.Println()
			fmt.Println("  WCAG 2.1 AA")
			fmt.Println("    - Accessibility requirements in UXD")
			fmt.Println("    - Evaluated in UXD rubric")
			fmt.Println()
			fmt.Println("  OWASP Top 10")
			fmt.Println("    - Security requirements in PRD")
			fmt.Println("    - Security architecture in TRD")
			return nil
		},
	})

	return cmd
}

// securityCmd provides security-specific commands.
func securityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "security",
		Short: "Security requirements management",
		Long:  `Commands for managing security requirements across specs.`,
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "checklist",
		Short: "Show security checklist",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Security Requirements Checklist:")
			fmt.Println()
			fmt.Println("Authentication:")
			fmt.Println("  [ ] MFA support required")
			fmt.Println("  [ ] Session management defined")
			fmt.Println("  [ ] Password policy specified")
			fmt.Println()
			fmt.Println("Authorization:")
			fmt.Println("  [ ] RBAC/ABAC model documented")
			fmt.Println("  [ ] Resource-level permissions defined")
			fmt.Println("  [ ] API authorization strategy")
			fmt.Println()
			fmt.Println("Data Security:")
			fmt.Println("  [ ] Encryption at rest requirements")
			fmt.Println("  [ ] Encryption in transit (TLS 1.3)")
			fmt.Println("  [ ] Key management strategy")
			fmt.Println()
			fmt.Println("API Security:")
			fmt.Println("  [ ] Rate limiting requirements")
			fmt.Println("  [ ] Input validation strategy")
			fmt.Println("  [ ] CORS policy defined")
			fmt.Println()
			fmt.Println("Audit & Monitoring:")
			fmt.Println("  [ ] Audit log requirements")
			fmt.Println("  [ ] Security event alerting")
			fmt.Println("  [ ] Incident response plan")
			return nil
		},
	})

	return cmd
}
