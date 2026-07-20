// Package mcp implements the Model Context Protocol server for visionspec.
package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/ProductBuildersHQ/visionspec/pkg/profiles"
	"github.com/ProductBuildersHQ/visionspec/pkg/rubrics"
	"github.com/ProductBuildersHQ/visionspec/pkg/templates"
	"github.com/ProductBuildersHQ/visionspec/pkg/types"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerResources registers MCP resources for templates, rubrics, and profiles.
func (s *Server) registerResources() {
	// Template resources
	s.server.AddResourceTemplate(&mcp.ResourceTemplate{
		URITemplate: "template://{spec_type}",
		Name:        "Spec Template",
		Description: "Get the template for a spec type (mrd, prd, uxd, etc.)",
		MIMEType:    "text/markdown",
	}, s.handleTemplateResource)

	// Rubric resources
	s.server.AddResourceTemplate(&mcp.ResourceTemplate{
		URITemplate: "rubric://{spec_type}",
		Name:        "Evaluation Rubric",
		Description: "Get the evaluation rubric for a spec type",
		MIMEType:    "application/yaml",
	}, s.handleRubricResource)

	// Profile resources
	s.server.AddResourceTemplate(&mcp.ResourceTemplate{
		URITemplate: "profile://{profile_name}",
		Name:        "Configuration Profile",
		Description: "Get a configuration profile (startup, growth, enterprise, 0-1)",
		MIMEType:    "application/yaml",
	}, s.handleProfileResource)

	// List resources
	s.server.AddResource(&mcp.Resource{
		URI:         "visionspec://templates",
		Name:        "Available Templates",
		Description: "List all available spec templates",
		MIMEType:    "application/json",
	}, s.handleListTemplates)

	s.server.AddResource(&mcp.Resource{
		URI:         "visionspec://rubrics",
		Name:        "Available Rubrics",
		Description: "List all available evaluation rubrics",
		MIMEType:    "application/json",
	}, s.handleListRubrics)

	s.server.AddResource(&mcp.Resource{
		URI:         "visionspec://profiles",
		Name:        "Available Profiles",
		Description: "List all available configuration profiles",
		MIMEType:    "application/json",
	}, s.handleListProfiles)
}

// handleTemplateResource returns the template content for a spec type.
func (s *Server) handleTemplateResource(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	// Extract spec type from URI
	uri := req.Params.URI
	specType := strings.TrimPrefix(uri, "template://")

	loader := templates.DefaultLoader()
	tmpl, err := loader.Load(types.SpecType(specType))
	if err != nil {
		return nil, fmt.Errorf("template not found: %s", specType)
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      uri,
				MIMEType: "text/markdown",
				Text:     tmpl.Content,
			},
		},
	}, nil
}

// handleRubricResource returns the rubric content for a spec type.
func (s *Server) handleRubricResource(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	// Extract spec type from URI
	uri := req.Params.URI
	specType := strings.TrimPrefix(uri, "rubric://")

	loader := rubrics.DefaultLoader()
	rubric, err := loader.Load(types.SpecType(specType))
	if err != nil {
		return nil, fmt.Errorf("rubric not found: %s", specType)
	}

	// Format rubric as YAML-like structure
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s Evaluation Rubric\n\n", rubric.Name))
	sb.WriteString(fmt.Sprintf("spec_type: %s\n", specType))
	sb.WriteString(fmt.Sprintf("name: %s\n", rubric.Name))
	sb.WriteString(fmt.Sprintf("description: %s\n\n", rubric.Description))

	sb.WriteString("categories:\n")
	for _, cat := range rubric.Categories {
		sb.WriteString(fmt.Sprintf("  - id: %s\n", cat.ID))
		sb.WriteString(fmt.Sprintf("    name: %s\n", cat.Name))
		sb.WriteString(fmt.Sprintf("    weight: %.2f\n", cat.Weight))
		sb.WriteString(fmt.Sprintf("    description: %s\n", cat.Description))
	}

	sb.WriteString("\npass_criteria:\n")
	sb.WriteString(fmt.Sprintf("  min_categories_passing: %s\n", rubric.PassCriteria.MinCategoriesPassing))
	if lim := rubric.PassCriteria.MaxFindings; lim != nil {
		sb.WriteString(fmt.Sprintf("  max_critical: %d\n", lim.Critical))
		sb.WriteString(fmt.Sprintf("  max_high: %d\n", lim.High))
		sb.WriteString(fmt.Sprintf("  max_medium: %d\n", lim.Medium))
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      uri,
				MIMEType: "application/yaml",
				Text:     sb.String(),
			},
		},
	}, nil
}

// handleProfileResource returns the profile content.
func (s *Server) handleProfileResource(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	// Extract profile name from URI
	uri := req.Params.URI
	profileName := strings.TrimPrefix(uri, "profile://")

	loader := profiles.DefaultLoader()
	profile, err := loader.Load(profileName)
	if err != nil {
		return nil, fmt.Errorf("profile not found: %s", profileName)
	}

	// Format profile as YAML-like structure
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s Profile\n\n", profile.Name))
	sb.WriteString(fmt.Sprintf("name: %s\n", profile.Name))
	sb.WriteString(fmt.Sprintf("description: %s\n", profile.Description))
	if profile.Extends != "" {
		sb.WriteString(fmt.Sprintf("extends: %s\n", profile.Extends))
	}

	sb.WriteString("\nspec_config:\n")
	if profile.SpecConfig != nil {
		for specType, req := range profile.SpecConfig.Specs {
			sb.WriteString(fmt.Sprintf("  %s:\n", specType))
			sb.WriteString(fmt.Sprintf("    required: %v\n", req.Required))
			sb.WriteString(fmt.Sprintf("    category: %s\n", req.Category))
			if req.Template != "" {
				sb.WriteString(fmt.Sprintf("    template: %s\n", req.Template))
			}
			if req.Rubric != "" {
				sb.WriteString(fmt.Sprintf("    rubric: %s\n", req.Rubric))
			}
		}
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      uri,
				MIMEType: "application/yaml",
				Text:     sb.String(),
			},
		},
	}, nil
}

// resourceListItem represents an item in a resource list.
type resourceListItem struct {
	key   string // e.g., "spec_type" or "name"
	value string
	uri   string
}

// buildResourceListJSON builds a JSON list for resources.
func buildResourceListJSON(listKey string, items []resourceListItem) string {
	var sb strings.Builder
	sb.WriteString("{\n")
	sb.WriteString(fmt.Sprintf("  \"%s\": [\n", listKey))
	for i, item := range items {
		comma := ","
		if i == len(items)-1 {
			comma = ""
		}
		sb.WriteString(fmt.Sprintf("    {\"%s\": \"%s\", \"uri\": \"%s\"}%s\n", item.key, item.value, item.uri, comma))
	}
	sb.WriteString("  ],\n")
	sb.WriteString(fmt.Sprintf("  \"count\": %d\n", len(items)))
	sb.WriteString("}\n")
	return sb.String()
}

// handleListTemplates returns a list of available templates.
func (s *Server) handleListTemplates(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	loader := templates.DefaultLoader()
	available := loader.Available()

	items := make([]resourceListItem, len(available))
	for i, specType := range available {
		items[i] = resourceListItem{key: "spec_type", value: string(specType), uri: "template://" + string(specType)}
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     buildResourceListJSON("templates", items),
			},
		},
	}, nil
}

// handleListRubrics returns a list of available rubrics.
func (s *Server) handleListRubrics(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	loader := rubrics.DefaultLoader()
	available := loader.Available()

	items := make([]resourceListItem, len(available))
	for i, specType := range available {
		items[i] = resourceListItem{key: "spec_type", value: string(specType), uri: "rubric://" + string(specType)}
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     buildResourceListJSON("rubrics", items),
			},
		},
	}, nil
}

// handleListProfiles returns a list of available profiles.
func (s *Server) handleListProfiles(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	loader := profiles.DefaultLoader()
	available := loader.Available()

	items := make([]resourceListItem, len(available))
	for i, name := range available {
		items[i] = resourceListItem{key: "name", value: name, uri: "profile://" + name}
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     buildResourceListJSON("profiles", items),
			},
		},
	}, nil
}
