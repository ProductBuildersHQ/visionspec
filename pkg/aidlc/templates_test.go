package aidlc

import (
	"strings"
	"testing"
)

func TestGetTemplate(t *testing.T) {
	// Test all known document types have templates
	for _, docType := range AllDocTypes() {
		t.Run(string(docType), func(t *testing.T) {
			tmpl, ok := GetTemplate(docType)
			if !ok {
				t.Errorf("GetTemplate(%s) returned false, want true", docType)
				return
			}
			if tmpl == nil {
				t.Errorf("GetTemplate(%s) returned nil template", docType)
				return
			}
			if tmpl.DocType != docType {
				t.Errorf("Template.DocType = %s, want %s", tmpl.DocType, docType)
			}
			if tmpl.Name == "" {
				t.Errorf("Template.Name is empty for %s", docType)
			}
			if tmpl.Content == "" {
				t.Errorf("Template.Content is empty for %s", docType)
			}
		})
	}

	// Test unknown document type
	t.Run("unknown", func(t *testing.T) {
		_, ok := GetTemplate(DocType("unknown"))
		if ok {
			t.Error("GetTemplate(unknown) returned true, want false")
		}
	})
}

func TestAllTemplates(t *testing.T) {
	templates := AllTemplates()
	if len(templates) != 12 {
		t.Errorf("AllTemplates() returned %d templates, want 12", len(templates))
	}

	// Verify each document type has a template
	for _, docType := range AllDocTypes() {
		if _, ok := templates[docType]; !ok {
			t.Errorf("AllTemplates() missing template for %s", docType)
		}
	}
}

func TestRenderTemplate(t *testing.T) {
	data := TemplateData{
		ProjectName: "TestProject",
		Title:       "Test Title",
		Author:      "Test Author",
		Date:        "2024-01-15",
		Version:     "1.0",
		Description: "Test description",
	}

	// Test rendering each template type
	for _, docType := range AllDocTypes() {
		t.Run(string(docType), func(t *testing.T) {
			content, err := RenderTemplate(docType, data)
			if err != nil {
				t.Errorf("RenderTemplate(%s) error = %v", docType, err)
				return
			}
			if content == "" {
				t.Errorf("RenderTemplate(%s) returned empty content", docType)
				return
			}

			// Check that template variables are substituted
			if !strings.Contains(content, "Test Title") {
				t.Errorf("RenderTemplate(%s) did not substitute Title", docType)
			}
			if !strings.Contains(content, "Test Author") {
				t.Errorf("RenderTemplate(%s) did not substitute Author", docType)
			}
		})
	}

	// Test unknown document type
	t.Run("unknown", func(t *testing.T) {
		_, err := RenderTemplate(DocType("unknown"), data)
		if err == nil {
			t.Error("RenderTemplate(unknown) should return error")
		}
	})
}

func TestDefaultTemplateData(t *testing.T) {
	data := DefaultTemplateData("MyProject")

	if data.ProjectName != "MyProject" {
		t.Errorf("ProjectName = %q, want %q", data.ProjectName, "MyProject")
	}
	if data.Title != "MyProject" {
		t.Errorf("Title = %q, want %q", data.Title, "MyProject")
	}
	if data.Version != "1.0" {
		t.Errorf("Version = %q, want %q", data.Version, "1.0")
	}
	if data.Date == "" {
		t.Error("Date should not be empty")
	}
	if data.Custom == nil {
		t.Error("Custom map should not be nil")
	}
}

func TestTemplateHasSections(t *testing.T) {
	// Verify each template has defined sections
	for _, docType := range AllDocTypes() {
		t.Run(string(docType), func(t *testing.T) {
			tmpl, ok := GetTemplate(docType)
			if !ok {
				t.Skipf("No template for %s", docType)
			}

			if len(tmpl.Sections) == 0 {
				t.Errorf("Template for %s has no sections defined", docType)
			}

			// Check each section has required fields
			for _, section := range tmpl.Sections {
				if section.ID == "" {
					t.Errorf("Section has empty ID in template %s", docType)
				}
				if section.Title == "" {
					t.Errorf("Section %s has empty Title in template %s", section.ID, docType)
				}
			}
		})
	}
}

func TestTemplateContentHasFrontmatter(t *testing.T) {
	// Verify templates have YAML frontmatter
	for _, docType := range AllDocTypes() {
		t.Run(string(docType), func(t *testing.T) {
			tmpl, ok := GetTemplate(docType)
			if !ok {
				t.Skipf("No template for %s", docType)
			}

			if !strings.HasPrefix(tmpl.Content, "---") {
				t.Errorf("Template for %s missing YAML frontmatter", docType)
			}

			// Check frontmatter ends properly
			parts := strings.SplitN(tmpl.Content, "---", 3)
			if len(parts) < 3 {
				t.Errorf("Template for %s has malformed frontmatter", docType)
			}
		})
	}
}

func TestTemplatePhaseGrouping(t *testing.T) {
	// Verify templates are grouped correctly by phase
	phases := map[Phase][]DocType{
		PhaseInception:    {DocVisionDocument, DocRequirementsSpec, DocTechnicalSpec, DocArchitectureSpec},
		PhaseConstruction: {DocImplementationPlan, DocTestPlan, DocIntegrationPlan, DocSecurityReview},
		PhaseOperations:   {DocRunbook, DocMonitoringPlan, DocDisasterPlan, DocSLODocument},
	}

	for phase, docTypes := range phases {
		t.Run(string(phase), func(t *testing.T) {
			for _, docType := range docTypes {
				tmpl, ok := GetTemplate(docType)
				if !ok {
					t.Errorf("Missing template for %s in %s phase", docType, phase)
					continue
				}

				if tmpl.DocType.Phase() != phase {
					t.Errorf("Template %s has phase %s, expected %s",
						docType, tmpl.DocType.Phase(), phase)
				}
			}
		})
	}
}
