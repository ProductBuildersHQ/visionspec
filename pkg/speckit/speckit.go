// Package speckit exports visionspec workflow families as GitHub Spec Kit
// extension packages (github.com/github/spec-kit, plugin schema 1.0).
//
// A family is split into two extensions by spec category: a product extension
// (source and gtm specs — the requirements chain, ending at the decision
// gate and UXD) and an engineering add-on extension (technical specs). Each
// spec type in the family's execution sequence becomes a spec-kit command
// that authors the document from the family's template and self-evaluates
// against the family's rubric.
//
// The output is a deterministic in-memory file set; WriteFiles persists it
// and Check detects drift against a previously written copy (the prism-sync
// pattern for committed generated artifacts).
package speckit

import (
	"fmt"
	"regexp"
	"strings"
)

// File is one generated file, with a path relative to the export root.
type File struct {
	Path    string
	Content []byte
}

// Options configures an export. Zero-value fields are filled by defaults
// derived from the workflow family.
type Options struct {
	// Family is the visionspec workflow name, used for provenance headers
	// and id defaults. Defaults to the loaded workflow's name.
	Family string

	// ProductID is the spec-kit extension id for the product chain.
	// Defaults to the family name with any "aws-" prefix stripped.
	ProductID string

	// EngineeringID is the spec-kit extension id for the technical add-on.
	// Defaults to ProductID + "-engineering".
	EngineeringID string

	// WorkflowID is the spec-kit workflow id for the generated end-to-end
	// chain (BuildWorkflow). Defaults to ProductID + "-chain".
	WorkflowID string

	// Version is the extension version (X.Y.Z). Defaults to "0.1.0".
	Version string

	// Author, Repository, and License populate the extension manifests.
	Author     string
	Repository string
	License    string

	// SpecKitVersion is the requires.speckit_version constraint.
	// Defaults to ">=1.0.0".
	SpecKitVersion string
}

var (
	extensionIDRe  = regexp.MustCompile(`^[a-z0-9-]+$`)
	versionRe      = regexp.MustCompile(`^\d+\.\d+\.\d+$`)
	commandNameRe  = regexp.MustCompile(`^speckit\.[a-z0-9-]+\.[a-z0-9-]+$`)
	maxDescription = 200
)

func (o Options) withDefaults(family string) (Options, error) {
	if o.Family == "" {
		o.Family = family
	}
	if o.ProductID == "" {
		o.ProductID = strings.TrimPrefix(o.Family, "aws-")
	}
	if o.EngineeringID == "" {
		o.EngineeringID = o.ProductID + "-engineering"
	}
	if o.WorkflowID == "" {
		o.WorkflowID = o.ProductID + "-chain"
	}
	if o.Version == "" {
		o.Version = "0.1.0"
	}
	if o.Author == "" {
		o.Author = "ProductBuildersHQ"
	}
	if o.Repository == "" {
		o.Repository = "https://github.com/ProductBuildersHQ/visionspec"
	}
	if o.License == "" {
		o.License = "MIT"
	}
	if o.SpecKitVersion == "" {
		o.SpecKitVersion = ">=1.0.0"
	}

	if !extensionIDRe.MatchString(o.ProductID) {
		return o, fmt.Errorf("invalid product extension id %q: must match %s", o.ProductID, extensionIDRe)
	}
	if !extensionIDRe.MatchString(o.EngineeringID) {
		return o, fmt.Errorf("invalid engineering extension id %q: must match %s", o.EngineeringID, extensionIDRe)
	}
	if !extensionIDRe.MatchString(o.WorkflowID) {
		return o, fmt.Errorf("invalid workflow id %q: must match %s", o.WorkflowID, extensionIDRe)
	}
	if !versionRe.MatchString(o.Version) {
		return o, fmt.Errorf("invalid extension version %q: must be X.Y.Z", o.Version)
	}
	return o, nil
}

// commandName returns the spec-kit command name for a spec type in an
// extension: speckit.<ext-id>.<spec>.
func commandName(extID, spec string) string {
	return "speckit." + extID + "." + spec
}

// commandToken returns the __SPECKIT_COMMAND_*__ cross-reference token for a
// command. Spec-kit rewrites the token at install time by lowercasing and
// replacing "_" with "." (hyphens survive), so the extension id and spec type
// are joined with "_" and everything else keeps its hyphens.
func commandToken(extID, spec string) string {
	return "__SPECKIT_COMMAND_" + strings.ToUpper(extID) + "_" + strings.ToUpper(spec) + "__"
}

// coreToken returns the token for a core spec-kit command (e.g. "PLAN").
func coreToken(cmd string) string {
	return "__SPECKIT_COMMAND_" + strings.ToUpper(cmd) + "__"
}

// specifyHandoff renders the instruction for handing a finished document
// chain off to spec-kit's core flow. It always targets speckit.specify, not
// speckit.plan: speckit.plan resolves its input from a FEATURE_SPEC that
// only speckit.specify's create-new-feature script produces, so a chain
// that skipped straight to plan would find no current feature. Passing the
// chain's own documents as speckit.specify's free-form $ARGUMENTS is a
// synthesis step, not a re-authoring one.
func specifyHandoff(docs []string) string {
	return fmt.Sprintf(
		"`%s` — read %s, and pass their scope, requirements, and user stories as its description. "+
			"This synthesizes spec-kit's own spec.md from what you just wrote, ready for `%s`.",
		coreToken("SPECIFY"), joinInitDocs(docs), coreToken("PLAN"))
}

// joinInitDocs renders a list of spec type ids as a prose-joined list of
// INIT_DIR-relative document references, e.g. "`INIT_DIR/prd.md` and
// `INIT_DIR/uxd.md`".
func joinInitDocs(specs []string) string {
	parts := make([]string, len(specs))
	for i, s := range specs {
		parts[i] = fmt.Sprintf("`INIT_DIR/%s.md`", s)
	}
	switch len(parts) {
	case 0:
		return ""
	case 1:
		return parts[0]
	case 2:
		return parts[0] + " and " + parts[1]
	default:
		return strings.Join(parts[:len(parts)-1], ", ") + ", and " + parts[len(parts)-1]
	}
}

// truncate shortens s to at most n runes, cutting at a word boundary and
// appending an ellipsis when truncation happens.
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	cut := s[:n-1]
	if i := strings.LastIndexByte(cut, ' '); i > 0 {
		cut = cut[:i]
	}
	return cut + "…"
}

// firstSentence returns the first sentence of s (up to the first ". "), or s
// unchanged if no sentence boundary is found.
func firstSentence(s string) string {
	if i := strings.Index(s, ". "); i > 0 {
		return s[:i+1]
	}
	return s
}
