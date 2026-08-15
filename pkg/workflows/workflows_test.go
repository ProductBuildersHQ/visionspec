package workflows

import (
	"testing"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"enterprise", false},
		{"aws-two-way-door", false},
		{"aws-one-way-door", false},
		{"nonexistent", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, err := Get(tt.name)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Get(%q) expected error, got nil", tt.name)
				}
				return
			}
			if err != nil {
				t.Fatalf("Get(%q) unexpected error: %v", tt.name, err)
			}
			if w.Workflow == nil {
				t.Errorf("Get(%q).Workflow is nil", tt.name)
			}
			if w.Workflow.Name != tt.name {
				t.Errorf("Get(%q).Workflow.Name = %q, want %q", tt.name, w.Workflow.Name, tt.name)
			}
		})
	}
}

func TestList(t *testing.T) {
	names := List()
	if len(names) < 3 {
		t.Errorf("List() returned %d workflows, want at least 3", len(names))
	}

	// Check expected workflows are present
	expected := map[string]bool{
		"enterprise":       false,
		"aws-two-way-door": false,
		"aws-one-way-door": false,
	}
	for _, name := range names {
		expected[name] = true
	}
	for name, found := range expected {
		if !found {
			t.Errorf("List() missing expected workflow %q", name)
		}
	}
}

func TestTemplatesLoaded(t *testing.T) {
	w, err := Get("aws-two-way-door")
	if err != nil {
		t.Fatalf("Get(aws-two-way-door) error: %v", err)
	}

	if len(w.Templates) == 0 {
		t.Error("aws-two-way-door has no templates loaded")
	}

	// Check press template has content
	press, ok := w.Templates["press"]
	if !ok {
		t.Fatal("aws-two-way-door missing press template")
	}
	if press.Content == "" {
		t.Error("press template has no content")
	}
	if press.SpecType != "press" {
		t.Errorf("press.SpecType = %q, want %q", press.SpecType, "press")
	}
}

func TestRubricsLoaded(t *testing.T) {
	w, err := Get("aws-two-way-door")
	if err != nil {
		t.Fatalf("Get(aws-two-way-door) error: %v", err)
	}

	if len(w.Rubrics) == 0 {
		t.Error("aws-two-way-door has no rubrics loaded")
	}

	// Check press rubric is loaded
	press, ok := w.Rubrics["press"]
	if !ok {
		t.Fatal("aws-two-way-door missing press rubric")
	}
	if press.ID == "" {
		t.Error("press rubric has no ID")
	}
	if len(press.Categories) == 0 {
		t.Error("press rubric has no categories")
	}
}

func TestMustGetPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustGet(nonexistent) did not panic")
		}
	}()
	MustGet("nonexistent")
}

func TestWorkflowInheritance(t *testing.T) {
	awsFeature, err := Get("aws-two-way-door")
	if err != nil {
		t.Fatalf("Get(aws-two-way-door) error: %v", err)
	}

	// aws-two-way-door extends enterprise
	if awsFeature.Workflow.Extends != "enterprise" {
		t.Errorf("aws-two-way-door.Extends = %q, want %q", awsFeature.Workflow.Extends, "enterprise")
	}

	// Check methodology is set
	if awsFeature.Workflow.Methodology == nil {
		t.Error("aws-two-way-door.Methodology is nil")
	}
	if awsFeature.Workflow.Methodology.Name == "" {
		t.Error("aws-two-way-door.Methodology.Name is empty")
	}
}
