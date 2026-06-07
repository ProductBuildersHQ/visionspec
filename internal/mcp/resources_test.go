package mcp

import (
	"context"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func makeReadResourceRequest(uri string) *mcp.ReadResourceRequest {
	return &mcp.ReadResourceRequest{
		Params: &mcp.ReadResourceParams{
			URI: uri,
		},
	}
}

func TestHandleTemplateResource(t *testing.T) {
	s := &Server{}

	tests := []struct {
		name     string
		uri      string
		wantErr  bool
		contains string
	}{
		{
			name:     "valid MRD template",
			uri:      "template://mrd",
			wantErr:  false,
			contains: "Market Requirements Document",
		},
		{
			name:     "valid PRD template",
			uri:      "template://prd",
			wantErr:  false,
			contains: "Product Requirements Document",
		},
		{
			name:     "valid UXD template",
			uri:      "template://uxd",
			wantErr:  false,
			contains: "User Experience",
		},
		{
			name:    "invalid template",
			uri:     "template://nonexistent",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := makeReadResourceRequest(tt.uri)

			result, err := s.handleTemplateResource(context.Background(), req)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(result.Contents) == 0 {
				t.Error("expected at least one content block")
				return
			}

			content := result.Contents[0]
			if content.URI != tt.uri {
				t.Errorf("expected URI %s, got %s", tt.uri, content.URI)
			}

			if content.MIMEType != "text/markdown" {
				t.Errorf("expected MIME type text/markdown, got %s", content.MIMEType)
			}

			if tt.contains != "" && !strings.Contains(content.Text, tt.contains) {
				t.Errorf("expected content to contain %q", tt.contains)
			}
		})
	}
}

func TestHandleRubricResource(t *testing.T) {
	s := &Server{}

	tests := []struct {
		name     string
		uri      string
		wantErr  bool
		contains string
	}{
		{
			name:     "valid MRD rubric",
			uri:      "rubric://mrd",
			wantErr:  false,
			contains: "spec_type: mrd",
		},
		{
			name:     "valid PRD rubric",
			uri:      "rubric://prd",
			wantErr:  false,
			contains: "spec_type: prd",
		},
		{
			name:    "invalid rubric",
			uri:     "rubric://nonexistent",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := makeReadResourceRequest(tt.uri)

			result, err := s.handleRubricResource(context.Background(), req)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(result.Contents) == 0 {
				t.Error("expected at least one content block")
				return
			}

			content := result.Contents[0]
			if content.URI != tt.uri {
				t.Errorf("expected URI %s, got %s", tt.uri, content.URI)
			}

			if content.MIMEType != "application/yaml" {
				t.Errorf("expected MIME type application/yaml, got %s", content.MIMEType)
			}

			if tt.contains != "" && !strings.Contains(content.Text, tt.contains) {
				t.Errorf("expected content to contain %q", tt.contains)
			}
		})
	}
}

func TestHandleProfileResource(t *testing.T) {
	s := &Server{}

	tests := []struct {
		name     string
		uri      string
		wantErr  bool
		contains string
	}{
		{
			name:     "startup profile",
			uri:      "profile://startup",
			wantErr:  false,
			contains: "name: startup",
		},
		{
			name:     "growth profile",
			uri:      "profile://growth",
			wantErr:  false,
			contains: "name: growth",
		},
		{
			name:     "enterprise profile",
			uri:      "profile://enterprise",
			wantErr:  false,
			contains: "name: enterprise",
		},
		{
			name:    "invalid profile",
			uri:     "profile://nonexistent",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := makeReadResourceRequest(tt.uri)

			result, err := s.handleProfileResource(context.Background(), req)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(result.Contents) == 0 {
				t.Error("expected at least one content block")
				return
			}

			content := result.Contents[0]
			if content.URI != tt.uri {
				t.Errorf("expected URI %s, got %s", tt.uri, content.URI)
			}

			if content.MIMEType != "application/yaml" {
				t.Errorf("expected MIME type application/yaml, got %s", content.MIMEType)
			}

			if tt.contains != "" && !strings.Contains(content.Text, tt.contains) {
				t.Errorf("expected content to contain %q, got: %s", tt.contains, content.Text)
			}
		})
	}
}

func TestHandleListTemplates(t *testing.T) {
	s := &Server{}

	req := makeReadResourceRequest("visionspec://templates")

	result, err := s.handleListTemplates(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Contents) == 0 {
		t.Fatal("expected at least one content block")
	}

	content := result.Contents[0]
	if content.MIMEType != "application/json" {
		t.Errorf("expected MIME type application/json, got %s", content.MIMEType)
	}

	// Check JSON contains templates array
	if !strings.Contains(content.Text, `"templates"`) {
		t.Error("expected JSON to contain 'templates' key")
	}

	// Check it includes some known templates
	if !strings.Contains(content.Text, "mrd") {
		t.Error("expected templates list to include 'mrd'")
	}
	if !strings.Contains(content.Text, "prd") {
		t.Error("expected templates list to include 'prd'")
	}
}

func TestHandleListRubrics(t *testing.T) {
	s := &Server{}

	req := makeReadResourceRequest("visionspec://rubrics")

	result, err := s.handleListRubrics(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Contents) == 0 {
		t.Fatal("expected at least one content block")
	}

	content := result.Contents[0]
	if content.MIMEType != "application/json" {
		t.Errorf("expected MIME type application/json, got %s", content.MIMEType)
	}

	// Check JSON contains rubrics array
	if !strings.Contains(content.Text, `"rubrics"`) {
		t.Error("expected JSON to contain 'rubrics' key")
	}
}

func TestHandleListProfiles(t *testing.T) {
	s := &Server{}

	req := makeReadResourceRequest("visionspec://profiles")

	result, err := s.handleListProfiles(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Contents) == 0 {
		t.Fatal("expected at least one content block")
	}

	content := result.Contents[0]
	if content.MIMEType != "application/json" {
		t.Errorf("expected MIME type application/json, got %s", content.MIMEType)
	}

	// Check JSON contains profiles array
	if !strings.Contains(content.Text, `"profiles"`) {
		t.Error("expected JSON to contain 'profiles' key")
	}

	// Check it includes known profiles
	if !strings.Contains(content.Text, "startup") {
		t.Error("expected profiles list to include 'startup'")
	}
	if !strings.Contains(content.Text, "growth") {
		t.Error("expected profiles list to include 'growth'")
	}
	if !strings.Contains(content.Text, "enterprise") {
		t.Error("expected profiles list to include 'enterprise'")
	}
}
