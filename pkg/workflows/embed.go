package workflows

import (
	"embed"
	"fmt"
	"io/fs"
	"path"
	"strings"

	"github.com/ProductBuildersHQ/specification-workflow-spec/pkg/template"
	"github.com/ProductBuildersHQ/specification-workflow-spec/pkg/workflow"
	oscompatfs "github.com/grokify/oscompat/fs"
	"github.com/plexusone/structured-evaluation/rubric"
	"gopkg.in/yaml.v3"
)

//go:embed default/*
var defaultFS embed.FS

func init() {
	if err := loadAllWorkflows(); err != nil {
		panic(fmt.Sprintf("failed to load embedded workflows: %v", err))
	}
}

// loadAllWorkflows discovers and loads all workflows from the embedded filesystem.
func loadAllWorkflows() error {
	entries, err := fs.ReadDir(defaultFS, "default")
	if err != nil {
		return fmt.Errorf("reading default workflows directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		workflowName := entry.Name()
		workflowDir := path.Join("default", workflowName)

		loaded, err := loadWorkflow(workflowDir)
		if err != nil {
			return fmt.Errorf("loading workflow %s: %w", workflowName, err)
		}

		Register(workflowName, loaded)
	}

	return nil
}

// loadWorkflow loads a single workflow from the embedded filesystem.
func loadWorkflow(dir string) (*LoadedWorkflow, error) {
	// Load profile.yaml (still named profile.yaml in the embedded files)
	workflowPath := path.Join(dir, "profile.yaml")
	workflowData, err := fs.ReadFile(defaultFS, workflowPath)
	if err != nil {
		return nil, fmt.Errorf("reading profile.yaml: %w", err)
	}

	w, err := workflow.ParseYAML(workflowData)
	if err != nil {
		return nil, fmt.Errorf("parsing profile.yaml: %w", err)
	}

	loaded := &LoadedWorkflow{
		Workflow:  w,
		Templates: make(map[string]*template.Template),
		Rubrics:   make(map[string]*rubric.RubricSet),
	}

	// Load templates
	templatesDir := path.Join(dir, "templates")
	if err := loadTemplates(loaded, templatesDir); err != nil {
		// Templates directory is optional
		if !isNotExist(err) {
			return nil, fmt.Errorf("loading templates: %w", err)
		}
	}

	// Load rubrics
	rubricsDir := path.Join(dir, "rubrics")
	if err := loadRubrics(loaded, rubricsDir); err != nil {
		// Rubrics directory is optional
		if !isNotExist(err) {
			return nil, fmt.Errorf("loading rubrics: %w", err)
		}
	}

	return loaded, nil
}

// loadTemplates loads all templates from a templates directory.
func loadTemplates(loaded *LoadedWorkflow, dir string) error {
	entries, err := fs.ReadDir(defaultFS, dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".md") {
			continue
		}

		// Extract spec type from filename (e.g., "press.md" -> "press")
		specType := strings.TrimSuffix(name, ".md")

		content, err := fs.ReadFile(defaultFS, path.Join(dir, name))
		if err != nil {
			return fmt.Errorf("reading template %s: %w", name, err)
		}

		// Normalize CRLF so embedded content is identical across OS checkouts.
		tmpl := &template.Template{
			ID:       specType,
			SpecType: specType,
			Content:  oscompatfs.NormalizeLineEndings(string(content)),
		}

		loaded.Templates[specType] = tmpl
	}

	return nil
}

// loadRubrics loads all rubrics from a rubrics directory.
func loadRubrics(loaded *LoadedWorkflow, dir string) error {
	entries, err := fs.ReadDir(defaultFS, dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".rubric.yaml") {
			continue
		}

		// Extract spec type from filename (e.g., "press.rubric.yaml" -> "press")
		specType := strings.TrimSuffix(name, ".rubric.yaml")

		data, err := fs.ReadFile(defaultFS, path.Join(dir, name))
		if err != nil {
			return fmt.Errorf("reading rubric %s: %w", name, err)
		}

		var rs rubric.RubricSet
		if err := yaml.Unmarshal(data, &rs); err != nil {
			return fmt.Errorf("parsing rubric %s: %w", name, err)
		}

		loaded.Rubrics[specType] = &rs
	}

	return nil
}

// isNotExist checks if an error is a "not exist" error.
func isNotExist(err error) bool {
	return err != nil && strings.Contains(err.Error(), "file does not exist")
}
