// Package lint validates visionspec directory structure and naming conventions.
package lint

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ProductBuildersHQ/visionspec/pkg/config"
	"github.com/ProductBuildersHQ/visionspec/pkg/types"
)

// Severity levels for lint findings.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo    Severity = "info"
)

// Finding represents a single lint issue.
type Finding struct {
	Path     string   `json:"path"`
	Rule     string   `json:"rule"`
	Message  string   `json:"message"`
	Severity Severity `json:"severity"`
}

// Result contains the complete lint output.
type Result struct {
	Project  string    `json:"project"`
	Findings []Finding `json:"findings"`
	Errors   int       `json:"errors"`
	Warnings int       `json:"warnings"`
	Passed   bool      `json:"passed"`
}

// Rules
const (
	RuleProjectNaming      = "project-naming"
	RuleDirectoryStructure = "directory-structure"
	RuleSpecNaming         = "spec-naming"
	RuleConfigExists       = "config-exists"
	RuleRepoLevelNaming    = "repo-level-naming"
	RuleEvalNaming         = "eval-naming"
	RuleUnknownFile        = "unknown-file"
)

var (
	kebabCaseRegex = regexp.MustCompile(`^[a-z][a-z0-9]*(-[a-z0-9]+)*$`)
	specNameRegex  = regexp.MustCompile(`^[a-z]+\.md$`)
	evalNameRegex  = regexp.MustCompile(`^[a-z]+\.eval\.json$`)
)

// Linter validates visionspec projects.
type Linter struct {
	specsDir   string
	specConfig *types.SpecConfig
	findings   []Finding
}

// New creates a new Linter for the given specs directory.
func New(specsDir string) *Linter {
	return &Linter{
		specsDir:   specsDir,
		specConfig: types.DefaultSpecConfig(),
		findings:   make([]Finding, 0),
	}
}

// NewWithConfig creates a new Linter with a custom spec configuration.
func NewWithConfig(specsDir string, specConfig *types.SpecConfig) *Linter {
	if specConfig == nil {
		specConfig = types.DefaultSpecConfig()
	}
	return &Linter{
		specsDir:   specsDir,
		specConfig: specConfig,
		findings:   make([]Finding, 0),
	}
}

// LintAll validates all projects in the specs directory.
func (l *Linter) LintAll() (*Result, error) {
	result := &Result{
		Project:  "all",
		Findings: make([]Finding, 0),
	}

	// Check if specs directory exists
	if _, err := os.Stat(l.specsDir); os.IsNotExist(err) {
		result.Findings = append(result.Findings, Finding{
			Path:     l.specsDir,
			Rule:     RuleDirectoryStructure,
			Message:  fmt.Sprintf("specs directory does not exist: %s", l.specsDir),
			Severity: SeverityError,
		})
		result.Errors = 1
		return result, nil
	}

	// Check repo-level files
	l.lintRepoLevel()

	// Find all projects
	entries, err := os.ReadDir(l.specsDir)
	if err != nil {
		return nil, fmt.Errorf("reading specs directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			// Check repo-level file naming (CAPS.md)
			continue
		}

		projectName := entry.Name()
		projectPath := filepath.Join(l.specsDir, projectName)

		// Skip hidden directories
		if strings.HasPrefix(projectName, ".") {
			continue
		}

		projectResult, err := l.LintProject(projectName, projectPath)
		if err != nil {
			return nil, fmt.Errorf("linting project %s: %w", projectName, err)
		}

		result.Findings = append(result.Findings, projectResult.Findings...)
	}

	// Count errors and warnings
	for _, f := range result.Findings {
		switch f.Severity {
		case SeverityError:
			result.Errors++
		case SeverityWarning:
			result.Warnings++
		}
	}

	result.Passed = result.Errors == 0
	return result, nil
}

// LintProject validates a single project.
func (l *Linter) LintProject(name, path string) (*Result, error) {
	l.findings = make([]Finding, 0)

	// Check project name is kebab-case
	l.checkProjectNaming(name, path)

	// Check directory structure
	l.checkDirectoryStructure(path)

	// Check config exists
	l.checkConfigExists(path)

	// Check spec files
	l.checkSpecFiles(path)

	// Check eval files
	l.checkEvalFiles(path)

	result := &Result{
		Project:  name,
		Findings: l.findings,
	}

	for _, f := range l.findings {
		switch f.Severity {
		case SeverityError:
			result.Errors++
		case SeverityWarning:
			result.Warnings++
		}
	}

	result.Passed = result.Errors == 0
	return result, nil
}

func (l *Linter) addFinding(path, rule, message string, severity Severity) {
	l.findings = append(l.findings, Finding{
		Path:     path,
		Rule:     rule,
		Message:  message,
		Severity: severity,
	})
}

func (l *Linter) lintRepoLevel() {
	entries, err := os.ReadDir(l.specsDir)
	if err != nil {
		return
	}

	allowedRepoFiles := map[string]bool{
		"CONSTITUTION.md": true,
		"ROADMAP.md":      true,
		"README.md":       true,
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		path := filepath.Join(l.specsDir, name)

		// Check if it's a known repo-level file
		if allowedRepoFiles[name] {
			continue
		}

		// Check if it follows CAPS.md pattern for unknown files
		if strings.HasSuffix(name, ".md") {
			if name != strings.ToUpper(name) {
				l.findings = append(l.findings, Finding{
					Path:     path,
					Rule:     RuleRepoLevelNaming,
					Message:  fmt.Sprintf("repo-level markdown files should be CAPS.md: %s", name),
					Severity: SeverityWarning,
				})
			}
		}
	}
}

func (l *Linter) checkProjectNaming(name, path string) {
	if !kebabCaseRegex.MatchString(name) {
		l.addFinding(path, RuleProjectNaming,
			fmt.Sprintf("project name must be kebab-case: %s", name),
			SeverityError)
	}
}

func (l *Linter) checkDirectoryStructure(projectPath string) {
	requiredDirs := []string{
		config.SourceDir,
		config.GTMDir,
		config.TechnicalDir,
		config.EvalDir,
	}

	for _, dir := range requiredDirs {
		dirPath := filepath.Join(projectPath, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			l.addFinding(dirPath, RuleDirectoryStructure,
				fmt.Sprintf("missing required directory: %s", dir),
				SeverityError)
		}
	}
}

func (l *Linter) checkConfigExists(projectPath string) {
	configPath := filepath.Join(projectPath, config.ConfigFileName)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		l.addFinding(configPath, RuleConfigExists,
			"missing visionspec.yaml configuration file",
			SeverityWarning)
	}
}

func (l *Linter) checkSpecFiles(projectPath string) {
	// Check source specs
	sourceSpecs := l.specConfig.SpecsByCategory(types.CategorySource)
	l.checkSpecsInDirByName(filepath.Join(projectPath, config.SourceDir), sourceSpecs)

	// Check GTM specs
	gtmSpecs := l.specConfig.SpecsByCategory(types.CategoryGTM)
	l.checkSpecsInDirByName(filepath.Join(projectPath, config.GTMDir), gtmSpecs)

	// Check technical specs
	techSpecs := l.specConfig.SpecsByCategory(types.CategoryTechnical)
	l.checkSpecsInDirByName(filepath.Join(projectPath, config.TechnicalDir), techSpecs)
}

func (l *Linter) checkSpecsInDirByName(dirPath string, validSpecs []string) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return // Directory doesn't exist, already reported
	}

	validNames := make(map[string]bool)
	for _, specName := range validSpecs {
		// Convert spec name to filename (e.g., "prd" -> "prd.md")
		validNames[types.SpecType(specName).Filename()] = true
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		path := filepath.Join(dirPath, name)

		// Skip non-markdown files
		if !strings.HasSuffix(name, ".md") {
			continue
		}

		// Check naming convention (lowercase.md)
		if !specNameRegex.MatchString(name) {
			l.addFinding(path, RuleSpecNaming,
				fmt.Sprintf("spec file must be lowercase.md: %s", name),
				SeverityError)
			continue
		}

		// Check if it's a valid spec type for this directory
		if !validNames[name] {
			l.addFinding(path, RuleUnknownFile,
				fmt.Sprintf("unknown spec file in directory: %s", name),
				SeverityWarning)
		}
	}
}

func (l *Linter) checkEvalFiles(projectPath string) {
	evalDir := filepath.Join(projectPath, config.EvalDir)
	entries, err := os.ReadDir(evalDir)
	if err != nil {
		return // Directory doesn't exist, already reported
	}

	validEvalNames := make(map[string]bool)
	for _, specName := range l.specConfig.AllSpecs() {
		validEvalNames[types.SpecType(specName).EvalFilename()] = true
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		path := filepath.Join(evalDir, name)

		// Skip non-JSON files
		if !strings.HasSuffix(name, ".json") {
			continue
		}

		// Check eval naming convention (spec.eval.json)
		if !evalNameRegex.MatchString(name) {
			l.addFinding(path, RuleEvalNaming,
				fmt.Sprintf("eval file must be {spec}.eval.json: %s", name),
				SeverityError)
			continue
		}

		// Check if it's a valid eval file
		if !validEvalNames[name] {
			l.addFinding(path, RuleUnknownFile,
				fmt.Sprintf("unknown eval file: %s", name),
				SeverityWarning)
		}
	}
}

// FormatText renders the result as terminal text.
func (r *Result) FormatText() string {
	var sb strings.Builder

	if r.Passed {
		sb.WriteString(fmt.Sprintf("Lint passed for %s\n", r.Project))
	} else {
		sb.WriteString(fmt.Sprintf("Lint failed for %s\n", r.Project))
	}

	if len(r.Findings) == 0 {
		sb.WriteString("  No issues found\n")
		return sb.String()
	}

	sb.WriteString(fmt.Sprintf("  %d error(s), %d warning(s)\n\n", r.Errors, r.Warnings))

	for _, f := range r.Findings {
		icon := "  "
		switch f.Severity {
		case SeverityError:
			icon = "E "
		case SeverityWarning:
			icon = "W "
		case SeverityInfo:
			icon = "I "
		}

		sb.WriteString(fmt.Sprintf("%s[%s] %s\n", icon, f.Rule, f.Message))
		sb.WriteString(fmt.Sprintf("    %s\n", f.Path))
	}

	return sb.String()
}
