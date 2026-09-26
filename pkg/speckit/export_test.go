package speckit

import (
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/ProductBuildersHQ/visionspec/pkg/workflows"
)

const testFamily = "aws-one-way-door"

func buildTestFiles(t *testing.T) []File {
	t.Helper()
	// The resolving loader merges the extends chain (aws-one-way-door extends
	// enterprise), which is where the technical categories come from.
	lw, err := workflows.DefaultLoader().Load(testFamily)
	if err != nil {
		t.Fatal(err)
	}
	files, err := BuildExtensions(lw, Options{})
	if err != nil {
		t.Fatal(err)
	}
	return files
}

func fileByPath(t *testing.T, files []File, p string) File {
	t.Helper()
	for _, f := range files {
		if f.Path == p {
			return f
		}
	}
	t.Fatalf("file %s not in generated set", p)
	return File{}
}

func TestBuildExtensionsSplitsByCategory(t *testing.T) {
	files := buildTestFiles(t)

	wantCommands := map[string]bool{
		"extensions/one-way-door/commands/speckit.one-way-door.press.md":                       true,
		"extensions/one-way-door/commands/speckit.one-way-door.faq.md":                         true,
		"extensions/one-way-door/commands/speckit.one-way-door.mrd.md":                         true,
		"extensions/one-way-door/commands/speckit.one-way-door.prd.md":                         true,
		"extensions/one-way-door/commands/speckit.one-way-door.narrative-6p.md":                true,
		"extensions/one-way-door/commands/speckit.one-way-door.uxd.md":                         true,
		"extensions/one-way-door-engineering/commands/speckit.one-way-door-engineering.trd.md": true,
		"extensions/one-way-door-engineering/commands/speckit.one-way-door-engineering.tpd.md": true,
		"extensions/one-way-door-engineering/commands/speckit.one-way-door-engineering.ird.md": true,
	}
	got := map[string]bool{}
	for _, f := range files {
		if strings.Contains(f.Path, "/commands/") {
			got[f.Path] = true
		}
	}
	for p := range wantCommands {
		if !got[p] {
			t.Errorf("missing command file %s", p)
		}
	}
	for p := range got {
		if !wantCommands[p] {
			t.Errorf("unexpected command file %s", p)
		}
	}
}

func TestManifestsSatisfySpecKitConstraints(t *testing.T) {
	files := buildTestFiles(t)
	commandNamePattern := regexp.MustCompile(`^speckit\.[a-z0-9-]+\.[a-z0-9-]+$`)

	for _, extID := range []string{"one-way-door", "one-way-door-engineering"} {
		raw := fileByPath(t, files, "extensions/"+extID+"/extension.yml")

		var m Manifest
		if err := yaml.Unmarshal(raw.Content, &m); err != nil {
			t.Fatalf("%s: extension.yml does not parse: %v", extID, err)
		}
		if m.SchemaVersion != "1.0" {
			t.Errorf("%s: schema_version = %q, want 1.0", extID, m.SchemaVersion)
		}
		if m.Extension.ID != extID {
			t.Errorf("%s: id = %q", extID, m.Extension.ID)
		}
		if len(m.Extension.Description) >= 200 {
			t.Errorf("%s: description is %d chars, spec-kit requires < 200", extID, len(m.Extension.Description))
		}
		if m.Requires.SpecKitVersion == "" {
			t.Errorf("%s: requires.speckit_version is empty", extID)
		}
		if len(m.Provides.Commands) == 0 {
			t.Errorf("%s: no commands provided", extID)
		}

		// Manifest ↔ fileset completeness, both directions.
		for _, c := range m.Provides.Commands {
			if !commandNamePattern.MatchString(c.Name) {
				t.Errorf("%s: command name %q violates spec-kit pattern", extID, c.Name)
			}
			if c.Description == "" {
				t.Errorf("%s: command %s has no description", extID, c.Name)
			}
			fileByPath(t, files, "extensions/"+extID+"/"+c.File)
		}
		for _, tm := range m.Provides.Templates {
			fileByPath(t, files, "extensions/"+extID+"/"+tm.File)
		}
		declared := map[string]bool{}
		for _, c := range m.Provides.Commands {
			declared[filepath.ToSlash(c.File)] = true
		}
		for _, tm := range m.Provides.Templates {
			declared[filepath.ToSlash(tm.File)] = true
		}
		prefix := "extensions/" + extID + "/"
		for _, f := range files {
			rel := strings.TrimPrefix(f.Path, prefix)
			if rel == f.Path || rel == "extension.yml" || rel == "README.md" {
				continue
			}
			if !declared[rel] {
				t.Errorf("%s: file %s not declared in manifest provides", extID, rel)
			}
		}
	}
}

func TestCommandChainIntegrity(t *testing.T) {
	files := buildTestFiles(t)
	prefix := "extensions/one-way-door/commands/speckit.one-way-door."

	content := func(spec string) string {
		return string(fileByPath(t, files, prefix+spec+".md").Content)
	}

	faq := content("faq")
	if !strings.Contains(faq, "__SPECKIT_COMMAND_ONE-WAY-DOOR_PRD__") {
		t.Error("faq does not chain to prd")
	}
	if !strings.Contains(faq, "__SPECKIT_COMMAND_ONE-WAY-DOOR_MRD__") {
		t.Error("faq does not offer the optional mrd deepening")
	}
	if !strings.Contains(faq, "## Review Gate: PR/FAQ Review — STOP") {
		t.Error("faq is missing the PR/FAQ review gate")
	}

	sixPager := content("narrative-6p")
	if !strings.Contains(sixPager, "## Review Gate: Decision Meeting — STOP") {
		t.Error("narrative-6p is missing the decision-meeting gate")
	}
	if !strings.Contains(sixPager, "mandatory") {
		t.Error("narrative-6p gate is not marked mandatory")
	}

	uxd := content("uxd")
	if !strings.Contains(uxd, "__SPECKIT_COMMAND_ONE-WAY-DOOR-ENGINEERING_TRD__") {
		t.Error("uxd does not hand off to the engineering extension")
	}
	// The fallback must target speckit.specify, not speckit.plan directly:
	// speckit.plan requires a FEATURE_SPEC that only speckit.specify's
	// create-new-feature script produces.
	if !strings.Contains(uxd, "__SPECKIT_COMMAND_SPECIFY__") {
		t.Error("uxd has no core-specify fallback handoff")
	}
	if strings.Contains(uxd, "hand off to `__SPECKIT_COMMAND_PLAN__`") {
		t.Error("uxd hands off directly to plan, which has no FEATURE_SPEC without an intervening specify")
	}

	ird := string(fileByPath(t, files, "extensions/one-way-door-engineering/commands/speckit.one-way-door-engineering.ird.md").Content)
	if !strings.Contains(ird, "__SPECKIT_COMMAND_SPECIFY__") {
		t.Error("ird does not hand off to core specify")
	}
}

func TestGeneratedFileHygiene(t *testing.T) {
	files := buildTestFiles(t)
	tokenPattern := regexp.MustCompile(`__SPECKIT_COMMAND_[A-Z0-9_-]+__`)
	anyToken := regexp.MustCompile(`__SPECKIT[A-Za-z0-9_-]*__`)

	for _, f := range files {
		content := string(f.Content)
		if !strings.Contains(content, "Code generated by visionspec speckit export; DO NOT EDIT.") {
			t.Errorf("%s: missing provenance header", f.Path)
		}
		// Every token present must be well-formed so spec-kit's install-time
		// rewrite recognizes it.
		for _, tok := range anyToken.FindAllString(content, -1) {
			if !tokenPattern.MatchString(tok) {
				t.Errorf("%s: malformed command token %q", f.Path, tok)
			}
		}
		if strings.Contains(f.Path, "/commands/") && !strings.HasPrefix(content, "---\ndescription: ") {
			t.Errorf("%s: command file does not start with description frontmatter", f.Path)
		}
	}
}

func TestBuildExtensionsDeterministic(t *testing.T) {
	a := buildTestFiles(t)
	b := buildTestFiles(t)
	if len(a) != len(b) {
		t.Fatalf("file count differs between runs: %d vs %d", len(a), len(b))
	}
	for i := range a {
		if a[i].Path != b[i].Path || string(a[i].Content) != string(b[i].Content) {
			t.Errorf("output not deterministic at %s", a[i].Path)
		}
	}
}

func TestOptionsValidation(t *testing.T) {
	lw, err := workflows.DefaultLoader().Load(testFamily)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := BuildExtensions(lw, Options{ProductID: "Bad_ID"}); err == nil {
		t.Error("invalid product id accepted")
	}
	if _, err := BuildExtensions(lw, Options{Version: "1.0"}); err == nil {
		t.Error("invalid version accepted")
	}
}
