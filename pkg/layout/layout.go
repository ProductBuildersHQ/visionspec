// Package layout defines filesystem conventions for specification projects.
//
// A typical visionspec project layout:
//
//	project/
//	├── source/           # Source specs (human-authored)
//	│   ├── mrd.md
//	│   ├── prd.md
//	│   └── uxd.md
//	├── gtm/              # GTM specs (synthesized)
//	│   ├── press.md
//	│   ├── faq.md
//	│   └── narrative-6p.md
//	├── technical/        # Technical specs
//	│   ├── trd.md
//	│   └── ird.md
//	├── strategic/        # Strategic planning (V2MOM)
//	│   └── v2mom-vision.md
//	├── execution/        # Execution phase specs
//	│   └── tpd.md
//	├── evals/            # Evaluation results
//	│   ├── mrd.eval.json
//	│   └── prd.eval.json
//	├── spec.md           # Reconciled output spec
//	└── current-truth.md  # Post-ship state
package layout

import (
	"path"
	"path/filepath"
	"strings"
)

// Layout defines the filesystem layout for a specification project.
type Layout struct {
	// SourceDir is the directory for source specs (default: "source").
	SourceDir string `json:"source_dir,omitempty" yaml:"source_dir,omitempty"`

	// GTMDir is the directory for GTM specs (default: "gtm").
	GTMDir string `json:"gtm_dir,omitempty" yaml:"gtm_dir,omitempty"`

	// TechnicalDir is the directory for technical specs (default: "technical").
	TechnicalDir string `json:"technical_dir,omitempty" yaml:"technical_dir,omitempty"`

	// StrategicDir is the directory for strategic specs (default: "strategic").
	StrategicDir string `json:"strategic_dir,omitempty" yaml:"strategic_dir,omitempty"`

	// ExecutionDir is the directory for execution specs (default: "execution").
	ExecutionDir string `json:"execution_dir,omitempty" yaml:"execution_dir,omitempty"`

	// EvalDir is the directory for evaluation results (default: "evals").
	EvalDir string `json:"eval_dir,omitempty" yaml:"eval_dir,omitempty"`

	// SpecExtension is the file extension for specs (default: ".md").
	SpecExtension string `json:"spec_extension,omitempty" yaml:"spec_extension,omitempty"`

	// EvalExtension is the file extension for evals (default: ".eval.json").
	EvalExtension string `json:"eval_extension,omitempty" yaml:"eval_extension,omitempty"`
}

// DefaultLayout returns the default filesystem layout.
func DefaultLayout() *Layout {
	return &Layout{
		SourceDir:     "source",
		GTMDir:        "gtm",
		TechnicalDir:  "technical",
		StrategicDir:  "strategic",
		ExecutionDir:  "execution",
		EvalDir:       "evals",
		SpecExtension: ".md",
		EvalExtension: ".eval.json",
	}
}

// DirForCategory returns the directory for a given category.
func (l *Layout) DirForCategory(category string) string {
	if l == nil {
		l = DefaultLayout()
	}

	switch strings.ToLower(category) {
	case "source":
		if l.SourceDir != "" {
			return l.SourceDir
		}
		return "source"
	case "gtm":
		if l.GTMDir != "" {
			return l.GTMDir
		}
		return "gtm"
	case "technical":
		if l.TechnicalDir != "" {
			return l.TechnicalDir
		}
		return "technical"
	case "strategic":
		if l.StrategicDir != "" {
			return l.StrategicDir
		}
		return "strategic"
	case "execution":
		if l.ExecutionDir != "" {
			return l.ExecutionDir
		}
		return "execution"
	case "output":
		return "" // Output specs go in project root
	default:
		return ""
	}
}

// SpecFilename returns the filename for a spec type.
func (l *Layout) SpecFilename(specType string) string {
	ext := ".md"
	if l != nil && l.SpecExtension != "" {
		ext = l.SpecExtension
	}
	return specType + ext
}

// EvalFilename returns the evaluation filename for a spec type.
func (l *Layout) EvalFilename(specType string) string {
	ext := ".eval.json"
	if l != nil && l.EvalExtension != "" {
		ext = l.EvalExtension
	}
	return specType + ext
}

// SpecPath returns the full path for a spec, relative to the project root.
// Always uses forward slashes for cross-platform consistency.
func (l *Layout) SpecPath(specType, category string) string {
	dir := l.DirForCategory(category)
	filename := l.SpecFilename(specType)
	if dir == "" {
		return filename
	}
	return path.Join(dir, filename)
}

// EvalPath returns the full path for an evaluation result, relative to project root.
// Always uses forward slashes for cross-platform consistency.
func (l *Layout) EvalPath(specType string) string {
	evalDir := "evals"
	if l != nil && l.EvalDir != "" {
		evalDir = l.EvalDir
	}
	return path.Join(evalDir, l.EvalFilename(specType))
}

// SpecPathAbs returns the absolute path for a spec.
func (l *Layout) SpecPathAbs(projectRoot, specType, category string) string {
	return filepath.Join(projectRoot, l.SpecPath(specType, category))
}

// EvalPathAbs returns the absolute path for an evaluation result.
func (l *Layout) EvalPathAbs(projectRoot, specType string) string {
	return filepath.Join(projectRoot, l.EvalPath(specType))
}

// AllDirs returns all spec directories (excluding root/output).
func (l *Layout) AllDirs() []string {
	if l == nil {
		l = DefaultLayout()
	}
	return []string{
		l.SourceDir,
		l.GTMDir,
		l.TechnicalDir,
		l.StrategicDir,
		l.ExecutionDir,
		l.EvalDir,
	}
}

// CategoryDirs returns a map of category to directory.
func (l *Layout) CategoryDirs() map[string]string {
	return map[string]string{
		"source":    l.DirForCategory("source"),
		"gtm":       l.DirForCategory("gtm"),
		"technical": l.DirForCategory("technical"),
		"strategic": l.DirForCategory("strategic"),
		"execution": l.DirForCategory("execution"),
		"output":    "", // root
	}
}
