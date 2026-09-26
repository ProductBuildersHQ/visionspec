package speckit

import (
	"fmt"
	"path"
	"strings"

	"github.com/ProductBuildersHQ/visionspec/pkg/workflow"
	"github.com/ProductBuildersHQ/visionspec/pkg/workflows"
)

// BuildExtensions exports a loaded workflow family as spec-kit extension
// packages: a product extension (source and gtm specs) and, when the family
// has technical specs, an engineering add-on extension. File paths are
// relative to the export root (extensions/<id>/...). Output is deterministic
// for a given workflow and options.
func BuildExtensions(lw *workflows.LoadedWorkflow, opts Options) ([]File, error) {
	if lw == nil || lw.Workflow == nil {
		return nil, fmt.Errorf("nil workflow")
	}
	w := lw.Workflow
	if w.Execution == nil || len(w.Execution.Sequence) == 0 {
		return nil, fmt.Errorf("workflow %q has no execution sequence to export", w.Name)
	}

	opts, err := opts.withDefaults(w.Name)
	if err != nil {
		return nil, err
	}

	product, engineering := partitionSequence(w)
	if len(product) == 0 {
		return nil, fmt.Errorf("workflow %q has no product (source/gtm) specs in its sequence", w.Name)
	}

	// The product chain ends at the engineering add-on when one exists,
	// falling back to core spec-kit either way. The fallback hands off to
	// speckit.specify, not speckit.plan directly: speckit.plan hard-requires
	// a FEATURE_SPEC that only speckit.specify's create-new-feature script
	// produces, so a chain that skipped straight to plan would find no
	// current feature. speckit.specify's $ARGUMENTS is free-form natural
	// language, so reading the chain's own documents back into it is a
	// synthesis step, not a re-authoring one — it mirrors how VisionSpec's
	// own speckit export target already treats spec.md as downstream of the
	// PRD/TRD, not hand-authored independently.
	productFinal := "Hand off to " + specifyHandoff([]string{"prd", "uxd"}) + "\n"
	if len(engineering) > 0 {
		productFinal = fmt.Sprintf(
			"If the %s extension is installed, continue with `%s` (%s). Otherwise hand off to %s\n",
			opts.EngineeringID, commandToken(opts.EngineeringID, engineering[0]), specTitle(engineering[0]),
			specifyHandoff([]string{"prd", "uxd"}))
	}
	engineeringFinal := "The specification chain is complete. Hand off to " + specifyHandoff([]string{"prd", "uxd", "trd"}) + "\n"

	files, err := buildExtension(lw, opts, opts.ProductID, product, productFinal, productMeta)
	if err != nil {
		return nil, err
	}
	if len(engineering) > 0 {
		engFiles, err := buildExtension(lw, opts, opts.EngineeringID, engineering, engineeringFinal, engineeringMeta)
		if err != nil {
			return nil, err
		}
		files = append(files, engFiles...)
	}
	return files, nil
}

// partitionSequence splits the execution sequence into the product chain
// (source, gtm, and any other non-technical categories) and the engineering
// chain (technical), preserving sequence order.
func partitionSequence(w *workflow.Workflow) (product, engineering []string) {
	for _, spec := range w.Execution.Sequence {
		category := ""
		if req, ok := w.SpecConfig[spec]; ok {
			category = req.Category
		}
		if category == "technical" {
			engineering = append(engineering, spec)
		} else {
			product = append(product, spec)
		}
	}
	return product, engineering
}

// extensionMeta derives the manifest name, description, and tags for one of
// the two extensions of a family.
type extensionMeta func(w *workflow.Workflow, opts Options, specs []string) (name, description string, tags []string)

func productMeta(w *workflow.Workflow, opts Options, specs []string) (string, string, []string) {
	name := specTitle(opts.ProductID) + " Product Chain"
	description := chainDescription(w, specs, "requirements chain with review gates and quality rubrics")
	return name, description, extensionTags(w, "product")
}

func engineeringMeta(w *workflow.Workflow, opts Options, specs []string) (string, string, []string) {
	name := specTitle(opts.ProductID) + " Engineering Add-On"
	description := chainDescription(w, specs, "technical spec chain extending the product extension")
	return name, description, extensionTags(w, "engineering")
}

// chainDescription builds a manifest description (< 200 chars) from the
// methodology name, the spec chain, and a suffix.
func chainDescription(w *workflow.Workflow, specs []string, suffix string) string {
	method := "Spec-driven"
	if w.Methodology != nil && w.Methodology.Name != "" {
		method = w.Methodology.Name
	}
	desc := fmt.Sprintf("%s %s: %s", method, suffix, strings.Join(specs, " → "))
	return truncate(desc, maxDescription-1)
}

func extensionTags(w *workflow.Workflow, audience string) []string {
	tags := []string{audience, "process"}
	if w.Methodology != nil && w.Methodology.Name != "" {
		tags = append(tags, slugify(w.Methodology.Name))
	}
	tags = append(tags, slugify(w.Name))
	return tags
}

// slugify lowercases and reduces a string to [a-z0-9-].
func slugify(s string) string {
	var b strings.Builder
	lastDash := true
	for _, r := range strings.ToLower(s) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			lastDash = false
		default:
			if !lastDash {
				b.WriteByte('-')
				lastDash = true
			}
		}
	}
	return strings.TrimSuffix(b.String(), "-")
}

// buildExtension assembles one extension package: manifest, README,
// commands, and template assets.
func buildExtension(lw *workflows.LoadedWorkflow, opts Options, extID string, specs []string, finalNext string, meta extensionMeta) ([]File, error) {
	w := lw.Workflow
	root := path.Join("extensions", extID)

	gates := map[string]*workflow.ReviewGate{}
	if w.Execution != nil {
		for i := range w.Execution.ReviewGates {
			g := &w.Execution.ReviewGates[i]
			gates[g.After] = g
		}
	}

	var files []File
	m := Manifest{
		SchemaVersion: "1.0",
		Requires:      ManifestRequires{SpecKitVersion: opts.SpecKitVersion},
	}
	name, description, tags := meta(w, opts, specs)
	m.Extension = ManifestMeta{
		ID:          extID,
		Name:        name,
		Version:     opts.Version,
		Description: description,
		Author:      opts.Author,
		Repository:  opts.Repository,
		License:     opts.License,
		Category:    "process",
		Effect:      "read-write",
	}
	m.Tags = tags

	for i, spec := range specs {
		cmdFile := "commands/" + commandName(extID, spec) + ".md"
		m.Provides.Commands = append(m.Provides.Commands, ManifestCommand{
			Name:        commandName(extID, spec),
			File:        cmdFile,
			Description: commandDescription(spec, w.SpecConfig[spec]),
		})

		content := renderCommand(commandInput{
			lw:        lw,
			opts:      opts,
			extID:     extID,
			specs:     specs,
			index:     i,
			gate:      gates[spec],
			finalNext: finalNext,
		})
		files = append(files, File{Path: path.Join(root, cmdFile), Content: content})

		if tmpl, ok := lw.Templates[spec]; ok {
			m.Provides.Templates = append(m.Provides.Templates, ManifestTemplate{
				Name: spec,
				File: "templates/" + spec + ".md",
			})
			templateContent := markdownHeader(opts.Family) + tmpl.Content
			files = append(files, File{Path: path.Join(root, "templates", spec+".md"), Content: []byte(templateContent)})
		}
	}

	manifest, err := renderManifest(m, opts.Family)
	if err != nil {
		return nil, err
	}
	readme := renderExtensionReadme(m, opts, specs)

	// Manifest and README first, then commands and templates in chain order.
	out := make([]File, 0, len(files)+2)
	out = append(out,
		File{Path: path.Join(root, "extension.yml"), Content: manifest},
		File{Path: path.Join(root, "README.md"), Content: readme},
	)
	return append(out, files...), nil
}

// renderExtensionReadme renders a minimal README for a packaged extension.
func renderExtensionReadme(m Manifest, opts Options, specs []string) []byte {
	var b strings.Builder
	b.WriteString(markdownHeader(opts.Family))
	fmt.Fprintf(&b, "# %s\n\n%s\n\n", m.Extension.Name, m.Extension.Description)
	fmt.Fprintf(&b, "Generated from the visionspec `%s` workflow family. Regenerate with `visionspec speckit export %s`; do not edit by hand.\n\n",
		opts.Family, opts.Family)
	fmt.Fprintf(&b, "## Install\n\n```bash\nspecify extension add %s --dev <path-to-this-directory>\n```\n\n", m.Extension.ID)
	b.WriteString("## Commands\n\n")
	for _, c := range m.Provides.Commands {
		fmt.Fprintf(&b, "- `/%s` — %s\n", c.Name, c.Description)
	}
	fmt.Fprintf(&b, "\nDocuments are written to `.specify/initiatives/<slug>/` in chain order: %s.\n", strings.Join(specs, " → "))
	return []byte(b.String())
}
