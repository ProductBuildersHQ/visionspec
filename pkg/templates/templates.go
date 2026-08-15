// Package templates provides embedded templates for spec types.
package templates

import (
	"embed"
	"fmt"
	"strings"
	"time"

	"github.com/ProductBuildersHQ/visionspec/pkg/types"
)

//go:embed mrd.md prd.md uxd.md narrative-1p.md narrative-6p.md press.md faq.md trd.md tpd.md ird.md
var templateFS embed.FS

// Template represents a loaded template with metadata.
type Template struct {
	SpecType types.SpecType
	Content  string
}

// Get returns the template for a spec type.
func Get(specType types.SpecType) (*Template, error) {
	filename := string(specType) + ".md"
	content, err := templateFS.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("template not found for spec type %q: %w", specType, err)
	}

	return &Template{
		SpecType: specType,
		Content:  string(content),
	}, nil
}

// MustGet returns the template for a spec type, panicking on error.
func MustGet(specType types.SpecType) *Template {
	t, err := Get(specType)
	if err != nil {
		panic(err)
	}
	return t
}

// Available returns all available template spec types.
func Available() []types.SpecType {
	return []types.SpecType{
		// Source specs
		types.SpecTypeMRD,
		types.SpecTypePRD,
		types.SpecTypeUXD,
		// GTM specs
		types.SpecTypePress,
		types.SpecTypeFAQ,
		types.SpecTypeNarrative1P,
		types.SpecTypeNarrative6P,
		// Technical specs
		types.SpecTypeTRD,
		types.SpecTypeTPD,
		types.SpecTypeIRD,
	}
}

// HasTemplate returns true if a template exists for the spec type.
func HasTemplate(specType types.SpecType) bool {
	_, err := Get(specType)
	return err == nil
}

// RenderOptions contains options for rendering a template.
type RenderOptions struct {
	ProjectName string
	Author      string
	Date        string
}

// DefaultRenderOptions returns default options for rendering.
func DefaultRenderOptions() RenderOptions {
	return RenderOptions{
		ProjectName: "Untitled Project",
		Author:      "Unknown",
		Date:        time.Now().Format("2006-01-02"),
	}
}

// Render renders the template with the given options.
func (t *Template) Render(opts RenderOptions) string {
	content := t.Content

	// Replace placeholders
	replacements := map[string]string{
		"{project_name}": opts.ProjectName,
		"{author}":       opts.Author,
		"{date}":         opts.Date,
	}

	for placeholder, value := range replacements {
		content = strings.ReplaceAll(content, placeholder, value)
	}

	return content
}

// RenderForProject is a convenience function to render a template for a project.
func RenderForProject(specType types.SpecType, projectName string) (string, error) {
	t, err := Get(specType)
	if err != nil {
		return "", err
	}

	opts := DefaultRenderOptions()
	opts.ProjectName = projectName

	return t.Render(opts), nil
}
