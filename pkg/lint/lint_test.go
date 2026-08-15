package lint

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ProductBuildersHQ/visionspec/pkg/config"
)

func TestKebabCaseRegex(t *testing.T) {
	valid := []string{
		"user-onboarding",
		"auth",
		"api-v2",
		"my-cool-project",
		"a",
		"a1",
		"test-123",
	}

	invalid := []string{
		"UserOnboarding",
		"user_onboarding",
		"User-Onboarding",
		"123-test",
		"-test",
		"test-",
		"test--double",
		"CAPS",
	}

	for _, name := range valid {
		if !kebabCaseRegex.MatchString(name) {
			t.Errorf("expected %q to be valid kebab-case", name)
		}
	}

	for _, name := range invalid {
		if kebabCaseRegex.MatchString(name) {
			t.Errorf("expected %q to be invalid kebab-case", name)
		}
	}
}

func TestSpecNameRegex(t *testing.T) {
	valid := []string{
		"mrd.md",
		"prd.md",
		"spec.md",
		"narrative.md",
	}

	invalid := []string{
		"MRD.md",
		"mrd.MD",
		"mrd.txt",
		"my-spec.md",
		"spec_file.md",
	}

	for _, name := range valid {
		if !specNameRegex.MatchString(name) {
			t.Errorf("expected %q to be valid spec name", name)
		}
	}

	for _, name := range invalid {
		if specNameRegex.MatchString(name) {
			t.Errorf("expected %q to be invalid spec name", name)
		}
	}
}

func TestEvalNameRegex(t *testing.T) {
	valid := []string{
		"mrd.eval.json",
		"prd.eval.json",
		"spec.eval.json",
	}

	invalid := []string{
		"mrd.json",
		"MRD.eval.json",
		"mrd.eval.yaml",
		"mrd-eval.json",
	}

	for _, name := range valid {
		if !evalNameRegex.MatchString(name) {
			t.Errorf("expected %q to be valid eval name", name)
		}
	}

	for _, name := range invalid {
		if evalNameRegex.MatchString(name) {
			t.Errorf("expected %q to be invalid eval name", name)
		}
	}
}

func TestLintProject(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "visionspec-lint-test")
	if err != nil {
		t.Fatalf("creating temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a valid project structure
	projectPath := filepath.Join(tmpDir, "user-onboarding")
	dirs := []string{"source", "gtm", "technical", "eval"}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(projectPath, dir), 0755); err != nil {
			t.Fatalf("creating dir %s: %v", dir, err)
		}
	}

	// Create config file
	if err := os.WriteFile(filepath.Join(projectPath, config.ConfigFileName), []byte("name: user-onboarding\n"), 0600); err != nil {
		t.Fatalf("creating config: %v", err)
	}

	// Lint the project
	linter := New(tmpDir)
	result, err := linter.LintProject("user-onboarding", projectPath)
	if err != nil {
		t.Fatalf("linting project: %v", err)
	}

	if !result.Passed {
		t.Errorf("expected lint to pass, got %d errors: %v", result.Errors, result.Findings)
	}
}

func TestLintProjectBadName(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "visionspec-lint-test")
	if err != nil {
		t.Fatalf("creating temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	projectPath := filepath.Join(tmpDir, "BadName")
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatalf("creating project dir: %v", err)
	}

	linter := New(tmpDir)
	result, err := linter.LintProject("BadName", projectPath)
	if err != nil {
		t.Fatalf("linting project: %v", err)
	}

	if result.Passed {
		t.Error("expected lint to fail for bad project name")
	}

	hasNamingError := false
	for _, f := range result.Findings {
		if f.Rule == RuleProjectNaming {
			hasNamingError = true
			break
		}
	}

	if !hasNamingError {
		t.Error("expected project-naming error")
	}
}

func TestLintProjectMissingDirs(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "visionspec-lint-test")
	if err != nil {
		t.Fatalf("creating temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	projectPath := filepath.Join(tmpDir, "test-project")
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatalf("creating project dir: %v", err)
	}

	linter := New(tmpDir)
	result, err := linter.LintProject("test-project", projectPath)
	if err != nil {
		t.Fatalf("linting project: %v", err)
	}

	if result.Passed {
		t.Error("expected lint to fail for missing directories")
	}

	// Should have 4 directory structure errors
	dirErrors := 0
	for _, f := range result.Findings {
		if f.Rule == RuleDirectoryStructure {
			dirErrors++
		}
	}

	if dirErrors != 4 {
		t.Errorf("expected 4 directory structure errors, got %d", dirErrors)
	}
}

func TestResultFormatText(t *testing.T) {
	result := &Result{
		Project: "test",
		Findings: []Finding{
			{
				Path:     "/path/to/file",
				Rule:     RuleProjectNaming,
				Message:  "test message",
				Severity: SeverityError,
			},
		},
		Errors:   1,
		Warnings: 0,
		Passed:   false,
	}

	output := result.FormatText()
	if output == "" {
		t.Error("expected non-empty output")
	}

	// Check for key content
	if !contains(output, "Lint failed") {
		t.Error("expected 'Lint failed' in output")
	}
	if !contains(output, "test message") {
		t.Error("expected error message in output")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr, 0))
}

func containsAt(s, substr string, start int) bool {
	for i := start; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
