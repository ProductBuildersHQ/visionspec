package templates

import (
	"strings"
	"testing"

	"github.com/ProductBuildersHQ/visionspec/pkg/types"
)

func TestGet(t *testing.T) {
	tests := []struct {
		specType    types.SpecType
		wantErr     bool
		wantContain string
	}{
		{types.SpecTypeMRD, false, "Market Requirements Document"},
		{types.SpecTypePRD, false, "Product Requirements Document"},
		{types.SpecTypeUXD, false, "User Experience Design"},
		{types.SpecTypeNarrative1P, false, "Executive Summary"},
		{types.SpecTypeNarrative6P, false, "Customer Problem"},
		{types.SpecTypePress, false, "Press Release"},
		{types.SpecTypeFAQ, false, "FAQ"},
		{types.SpecTypeTRD, false, "Technical Requirements Document"},
		{types.SpecTypeIRD, false, "Infrastructure Requirements Document"},
		{"unknown", true, ""}, // No template for unknown type
	}

	for _, tt := range tests {
		t.Run(string(tt.specType), func(t *testing.T) {
			tmpl, err := Get(tt.specType)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Get(%s) expected error, got nil", tt.specType)
				}
				return
			}
			if err != nil {
				t.Errorf("Get(%s) unexpected error: %v", tt.specType, err)
				return
			}
			if !strings.Contains(tmpl.Content, tt.wantContain) {
				t.Errorf("Get(%s) content does not contain %q", tt.specType, tt.wantContain)
			}
		})
	}
}

func TestMustGet(t *testing.T) {
	// Should not panic for valid types
	tmpl := MustGet(types.SpecTypeMRD)
	if tmpl == nil {
		t.Error("MustGet(MRD) returned nil")
	}

	// Should panic for invalid types
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustGet(unknown) did not panic")
		}
	}()
	MustGet("unknown")
}

func TestAvailable(t *testing.T) {
	available := Available()
	if len(available) != 9 {
		t.Errorf("Available() returned %d types, expected 9", len(available))
	}

	expected := map[types.SpecType]bool{
		types.SpecTypeMRD:         true,
		types.SpecTypePRD:         true,
		types.SpecTypeUXD:         true,
		types.SpecTypePress:       true,
		types.SpecTypeFAQ:         true,
		types.SpecTypeNarrative1P: true,
		types.SpecTypeNarrative6P: true,
		types.SpecTypeTRD:         true,
		types.SpecTypeIRD:         true,
	}

	for _, st := range available {
		if !expected[st] {
			t.Errorf("Available() returned unexpected type: %s", st)
		}
	}
}

func TestHasTemplate(t *testing.T) {
	hasTemplateTypes := []types.SpecType{
		types.SpecTypeMRD,
		types.SpecTypePRD,
		types.SpecTypeUXD,
		types.SpecTypePress,
		types.SpecTypeFAQ,
		types.SpecTypeNarrative1P,
		types.SpecTypeNarrative6P,
		types.SpecTypeTRD,
		types.SpecTypeIRD,
	}
	for _, st := range hasTemplateTypes {
		if !HasTemplate(st) {
			t.Errorf("HasTemplate(%s) returned false", st)
		}
	}
	if HasTemplate("unknown") {
		t.Error("HasTemplate(unknown) returned true")
	}
}

func TestRender(t *testing.T) {
	tmpl, err := Get(types.SpecTypeMRD)
	if err != nil {
		t.Fatalf("Get(MRD) failed: %v", err)
	}

	opts := RenderOptions{
		ProjectName: "Test Project",
		Author:      "Test Author",
		Date:        "2024-01-15",
	}

	rendered := tmpl.Render(opts)

	if !strings.Contains(rendered, "Test Project") {
		t.Error("Rendered content does not contain project name")
	}
	if !strings.Contains(rendered, "Test Author") {
		t.Error("Rendered content does not contain author")
	}
	if !strings.Contains(rendered, "2024-01-15") {
		t.Error("Rendered content does not contain date")
	}
}

func TestRenderForProject(t *testing.T) {
	content, err := RenderForProject(types.SpecTypePRD, "My Project")
	if err != nil {
		t.Fatalf("RenderForProject failed: %v", err)
	}

	if !strings.Contains(content, "My Project") {
		t.Error("Rendered content does not contain project name")
	}
}

func TestDefaultRenderOptions(t *testing.T) {
	opts := DefaultRenderOptions()
	if opts.ProjectName == "" {
		t.Error("DefaultRenderOptions().ProjectName is empty")
	}
	if opts.Author == "" {
		t.Error("DefaultRenderOptions().Author is empty")
	}
	if opts.Date == "" {
		t.Error("DefaultRenderOptions().Date is empty")
	}
}
