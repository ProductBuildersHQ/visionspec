// Package mcp implements the Model Context Protocol server for visionspec.
package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ProductBuildersHQ/visionspec/pkg/align"
	"github.com/ProductBuildersHQ/visionspec/pkg/config"
	ctxpkg "github.com/ProductBuildersHQ/visionspec/pkg/context"
	"github.com/ProductBuildersHQ/visionspec/pkg/draft"
	"github.com/ProductBuildersHQ/visionspec/pkg/eval"
	"github.com/ProductBuildersHQ/visionspec/pkg/reconcile"
	"github.com/ProductBuildersHQ/visionspec/pkg/status"
	"github.com/ProductBuildersHQ/visionspec/pkg/synth"
	"github.com/ProductBuildersHQ/visionspec/pkg/target"
	"github.com/ProductBuildersHQ/visionspec/pkg/templates"
	"github.com/ProductBuildersHQ/visionspec/pkg/testmap"
	"github.com/ProductBuildersHQ/visionspec/pkg/types"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Server implements the MCP server for visionspec.
type Server struct {
	server *mcp.Server
}

// NewServer creates a new MCP server.
func NewServer() *Server {
	impl := &mcp.Implementation{
		Name:    "visionspec",
		Version: "0.2.0",
	}

	server := mcp.NewServer(impl, nil)

	s := &Server{server: server}
	s.registerTools()
	s.registerResources()

	return s
}

// Serve starts the MCP server on stdio transport.
func (s *Server) Serve(ctx context.Context) error {
	transport := &mcp.StdioTransport{}
	return s.server.Run(ctx, transport)
}

// Tool argument types

// ProjectArgs contains a project name argument.
type ProjectArgs struct {
	Project string `json:"project" jsonschema:"description=Project name"`
}

// ProjectSpecArgs contains project and spec type arguments.
type ProjectSpecArgs struct {
	Project  string `json:"project" jsonschema:"description=Project name"`
	SpecType string `json:"spec_type" jsonschema:"description=Spec type (mrd, prd, uxd, trd, etc.)"`
}

// ApproveArgs contains approve command arguments.
type ApproveArgs struct {
	Project  string `json:"project" jsonschema:"description=Project name"`
	SpecType string `json:"spec_type" jsonschema:"description=Spec type"`
	Approver string `json:"approver,omitempty" jsonschema:"description=Approver identifier"`
}

// ExportArgs contains export command arguments.
type ExportArgs struct {
	Project string `json:"project" jsonschema:"description=Project name"`
	Target  string `json:"target" jsonschema:"description=Target (speckit, gsd, gastown, gascity)"`
}

// StartDraftArgs contains start_draft arguments.
type StartDraftArgs struct {
	Project  string `json:"project" jsonschema:"description=Project name"`
	SpecType string `json:"spec_type" jsonschema:"description=Spec type (mrd, prd, uxd)"`
}

// UpdateDraftArgs contains update_draft arguments.
type UpdateDraftArgs struct {
	Project  string `json:"project" jsonschema:"description=Project name"`
	SpecType string `json:"spec_type" jsonschema:"description=Spec type (mrd, prd, uxd)"`
	Content  string `json:"content" jsonschema:"description=Updated draft content (full markdown document)"`
}

// DraftArgs contains draft operation arguments.
type DraftArgs struct {
	Project  string `json:"project" jsonschema:"description=Project name"`
	SpecType string `json:"spec_type" jsonschema:"description=Spec type (mrd, prd, uxd)"`
}

func (s *Server) registerTools() {
	// Tool: list_projects
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "list_projects",
		Description: "List all visionspec projects",
	}, s.handleListProjects)

	// Tool: get_project_status
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "get_project_status",
		Description: "Get status and readiness for a project",
	}, s.handleGetProjectStatus)

	// Tool: get_spec
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "get_spec",
		Description: "Get the content of a specification",
	}, s.handleGetSpec)

	// Tool: get_eval
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "get_eval",
		Description: "Get evaluation results for a specification",
	}, s.handleGetEval)

	// Tool: run_eval
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "run_eval",
		Description: "Run evaluation on a specification",
	}, s.handleRunEval)

	// Tool: synthesize
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "synthesize",
		Description: "Generate a spec from source documents",
	}, s.handleSynthesize)

	// Tool: reconcile
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "reconcile",
		Description: "Generate unified execution spec from approved specs",
	}, s.handleReconcile)

	// Tool: approve
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "approve",
		Description: "Approve a specification",
	}, s.handleApprove)

	// Tool: export
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "export",
		Description: "Export specs to target execution system",
	}, s.handleExport)

	// Draft authoring tools

	// Tool: start_draft
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "start_draft",
		Description: "Initialize a new draft for a source spec (MRD, PRD, UXD). Returns template content and instructions.",
	}, s.handleStartDraft)

	// Tool: get_draft
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "get_draft",
		Description: "Get the current content and metadata of a draft",
	}, s.handleGetDraft)

	// Tool: update_draft
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "update_draft",
		Description: "Save updated content to a draft. Increments version number.",
	}, s.handleUpdateDraft)

	// Tool: eval_draft
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "eval_draft",
		Description: "Evaluate a draft against its rubric. Returns findings, score, and pass/fail decision.",
	}, s.handleEvalDraft)

	// Tool: finalize_draft
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "finalize_draft",
		Description: "Promote a draft to final spec. Moves content from draft to spec location.",
	}, s.handleFinalizeDraft)

	// Tool: discard_draft
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "discard_draft",
		Description: "Delete a draft and its metadata",
	}, s.handleDiscardDraft)

	// Tool: list_drafts
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "list_drafts",
		Description: "List all drafts in a project",
	}, s.handleListDrafts)

	// Tool: align
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "align",
		Description: "Check alignment between spec.md and implementation. Returns discrepancies, alignment score, and coverage metrics.",
	}, s.handleAlign)

	// Tool: get_resolution_plan
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "get_resolution_plan",
		Description: "Generate a resolution plan for alignment discrepancies. Returns prioritized actions and strategies.",
	}, s.handleGetResolutionPlan)

	// Tool: get_test_coverage
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "get_test_coverage",
		Description: "Get test coverage mapping between requirements and tests.",
	}, s.handleGetTestCoverage)

	// Tool: get_execution_context
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "get_execution_context",
		Description: "Get execution context for AI agents including spec summary, requirements, and implementation guidance.",
	}, s.handleGetExecutionContext)

	// Tool: get_execution_status
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "get_execution_status",
		Description: "Get execution status tracking for a project - which requirements are implemented, in progress, or pending.",
	}, s.handleGetExecutionStatus)

	// Tool: track_requirement
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "track_requirement",
		Description: "Track implementation progress for a specific requirement. Updates status and records evidence.",
	}, s.handleTrackRequirement)

	// Register execution prompts
	s.registerPrompts()
}

// Tool handlers

type emptyArgs struct{}

// errorResult creates an error response for MCP tools.
func errorResult(message string) (*mcp.CallToolResult, any, error) {
	result := map[string]any{
		"error": message,
	}
	data, _ := json.Marshal(result)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
		IsError: true,
	}, nil, nil
}

func (s *Server) handleListProjects(ctx context.Context, req *mcp.CallToolRequest, args emptyArgs) (*mcp.CallToolResult, any, error) {
	// Find specs directory from current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return errorResult("failed to get working directory: " + err.Error())
	}

	specsDir, err := config.FindSpecsDir(cwd)
	if err != nil {
		return errorResult("specs directory not found")
	}

	// List all project directories
	entries, err := os.ReadDir(specsDir)
	if err != nil {
		return errorResult("failed to read specs directory: " + err.Error())
	}

	var projects []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		// Skip hidden directories
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		// Check if it has a visionspec.yaml
		configPath := filepath.Join(specsDir, entry.Name(), config.ConfigFileName)
		if _, err := os.Stat(configPath); err == nil {
			projects = append(projects, entry.Name())
		}
	}

	result := map[string]any{
		"projects":  projects,
		"specs_dir": specsDir,
		"count":     len(projects),
	}
	data, _ := json.Marshal(result)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

func (s *Server) handleGetProjectStatus(ctx context.Context, req *mcp.CallToolRequest, args ProjectArgs) (*mcp.CallToolResult, any, error) {
	// Find specs directory
	cwd, err := os.Getwd()
	if err != nil {
		return errorResult("failed to get working directory: " + err.Error())
	}

	specsDir, err := config.FindSpecsDir(cwd)
	if err != nil {
		return errorResult("specs directory not found")
	}

	// Load project
	projectPath := config.ProjectPath(specsDir, args.Project)
	project, err := config.Load(projectPath)
	if err != nil {
		return errorResult("failed to load project: " + err.Error())
	}

	// Generate status report
	report, err := status.Generate(project)
	if err != nil {
		return errorResult("failed to generate status: " + err.Error())
	}

	data, _ := json.Marshal(report)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

func (s *Server) handleGetSpec(ctx context.Context, req *mcp.CallToolRequest, args ProjectSpecArgs) (*mcp.CallToolResult, any, error) {
	// Parse spec type
	specType := types.SpecType(strings.ToLower(args.SpecType))

	// Get project path
	projectPath, err := getProjectPath(args.Project)
	if err != nil {
		return errorResult("failed to find project: " + err.Error())
	}

	// Read spec content
	specPath := config.SpecPath(projectPath, specType)
	content, err := os.ReadFile(specPath)
	if err != nil {
		if os.IsNotExist(err) {
			return errorResult("spec not found: " + args.SpecType)
		}
		return errorResult("failed to read spec: " + err.Error())
	}

	result := map[string]any{
		"project":   args.Project,
		"spec_type": string(specType),
		"path":      specPath,
		"content":   string(content),
	}
	data, _ := json.Marshal(result)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

func (s *Server) handleGetEval(ctx context.Context, req *mcp.CallToolRequest, args ProjectSpecArgs) (*mcp.CallToolResult, any, error) {
	// Parse spec type
	specType := types.SpecType(strings.ToLower(args.SpecType))

	// Get project path
	projectPath, err := getProjectPath(args.Project)
	if err != nil {
		return errorResult("failed to find project: " + err.Error())
	}

	// Read eval file
	evalPath := config.EvalPath(projectPath, specType)
	content, err := os.ReadFile(evalPath)
	if err != nil {
		if os.IsNotExist(err) {
			return errorResult("evaluation not found for: " + args.SpecType)
		}
		return errorResult("failed to read evaluation: " + err.Error())
	}

	// Parse eval JSON to return structured data
	var evalData map[string]any
	if err := json.Unmarshal(content, &evalData); err != nil {
		return errorResult("failed to parse evaluation: " + err.Error())
	}

	result := map[string]any{
		"project":    args.Project,
		"spec_type":  string(specType),
		"path":       evalPath,
		"evaluation": evalData,
	}
	data, _ := json.Marshal(result)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

func (s *Server) handleRunEval(ctx context.Context, req *mcp.CallToolRequest, args ProjectSpecArgs) (*mcp.CallToolResult, any, error) {
	// Parse spec type
	specType := types.SpecType(strings.ToLower(args.SpecType))

	// Get project path
	projectPath, err := getProjectPath(args.Project)
	if err != nil {
		return errorResult("failed to find project: " + err.Error())
	}

	// Load project config for LLM settings
	project, err := config.Load(projectPath)
	if err != nil {
		return errorResult("failed to load project config: " + err.Error())
	}

	// Read spec content
	specPath := config.SpecPath(projectPath, specType)
	content, err := os.ReadFile(specPath)
	if err != nil {
		if os.IsNotExist(err) {
			return errorResult("spec not found: " + args.SpecType)
		}
		return errorResult("failed to read spec: " + err.Error())
	}

	// Create LLM client from project config (falls back to environment)
	llmClient, err := eval.NewLLMClientFromProject(project.LLM)
	if err != nil {
		return errorResult("failed to initialize LLM: " + err.Error())
	}
	defer func() { _ = llmClient.Close() }()

	// Run evaluation
	evaluator := eval.NewEvaluator(llmClient)
	evalResult, err := evaluator.Evaluate(ctx, specType, string(content))
	if err != nil {
		return errorResult("evaluation failed: " + err.Error())
	}

	result := map[string]any{
		"project":    args.Project,
		"spec_type":  string(specType),
		"score":      evalResult.Score,
		"passed":     evalResult.Passed,
		"decision":   evalResult.Decision,
		"categories": evalResult.Categories,
		"findings":   evalResult.Findings,
		"summary":    evalResult.Summary,
		"judge":      evalResult.Judge,
	}
	data, _ := json.Marshal(result)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

func (s *Server) handleSynthesize(ctx context.Context, req *mcp.CallToolRequest, args ProjectSpecArgs) (*mcp.CallToolResult, any, error) {
	// Parse spec type
	specType := types.SpecType(strings.ToLower(args.SpecType))

	// Validate that this spec type can be synthesized
	if !synth.CanSynthesize(specType) {
		return errorResult("spec type cannot be synthesized: " + args.SpecType + " (only trd, ird, press, faq, narrative-1p, narrative-6p)")
	}

	// Get project path
	projectPath, err := getProjectPath(args.Project)
	if err != nil {
		return errorResult("failed to find project: " + err.Error())
	}

	// Load project config for LLM settings
	project, err := config.Load(projectPath)
	if err != nil {
		return errorResult("failed to load project config: " + err.Error())
	}

	// Check required source specs exist
	requiredSources := synth.RequiredSources(specType)
	for _, srcType := range requiredSources {
		srcPath := config.SpecPath(projectPath, srcType)
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			return errorResult("missing required source spec: " + string(srcType))
		}
	}

	// Load source specs into input
	input := synth.SynthesisInput{}
	if content, err := os.ReadFile(config.SpecPath(projectPath, types.SpecTypeMRD)); err == nil {
		input.MRD = string(content)
	}
	if content, err := os.ReadFile(config.SpecPath(projectPath, types.SpecTypePRD)); err == nil {
		input.PRD = string(content)
	}
	if content, err := os.ReadFile(config.SpecPath(projectPath, types.SpecTypeUXD)); err == nil {
		input.UXD = string(content)
	}
	if content, err := os.ReadFile(config.SpecPath(projectPath, types.SpecTypeTRD)); err == nil {
		input.TRD = string(content)
	}
	if content, err := os.ReadFile(config.SpecPath(projectPath, types.SpecTypePress)); err == nil {
		input.Press = string(content)
	}

	// Load constitution if exists
	constitutionPath := filepath.Join(projectPath, "..", "CONSTITUTION.md")
	if content, err := os.ReadFile(constitutionPath); err == nil {
		input.Constitution = string(content)
	}

	// Create LLM client
	llmClient, err := eval.NewLLMClientFromProject(project.LLM)
	if err != nil {
		return errorResult("failed to initialize LLM: " + err.Error())
	}
	defer func() { _ = llmClient.Close() }()

	// Create synthesizer with adapter
	synthesizer := synth.NewSynthesizer(&synthLLMAdapter{client: llmClient})

	// Run synthesis
	synthResult, err := synthesizer.Synthesize(ctx, specType, input)
	if err != nil {
		return errorResult("synthesis failed: " + err.Error())
	}

	// Write output to file
	outputPath := config.SpecPath(projectPath, specType)
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return errorResult("failed to create output directory: " + err.Error())
	}
	if err := os.WriteFile(outputPath, []byte(synthResult.Content), 0600); err != nil {
		return errorResult("failed to write output: " + err.Error())
	}

	result := map[string]any{
		"project":   args.Project,
		"spec_type": string(specType),
		"path":      outputPath,
		"sources":   synthResult.Sources,
		"message":   "Synthesis completed successfully",
	}
	data, _ := json.Marshal(result)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

// synthLLMAdapter adapts eval.LLMClient to synth.LLMClient interface.
type synthLLMAdapter struct {
	client *eval.LLMClient
}

func (a *synthLLMAdapter) Complete(ctx context.Context, prompt string) (string, error) {
	content, _, err := a.client.Complete(ctx, prompt)
	return content, err
}

func (s *Server) handleReconcile(ctx context.Context, req *mcp.CallToolRequest, args ProjectArgs) (*mcp.CallToolResult, any, error) {
	// Get project path
	projectPath, err := getProjectPath(args.Project)
	if err != nil {
		return errorResult("failed to find project: " + err.Error())
	}

	// Load project config
	project, err := config.Load(projectPath)
	if err != nil {
		return errorResult("failed to load project config: " + err.Error())
	}

	// Check approvals
	approved, missing := reconcile.CheckApprovals(project.Approvals)
	if len(missing) > 0 {
		var missingStr []string
		for _, m := range missing {
			missingStr = append(missingStr, string(m))
		}
		return errorResult("missing approvals for: " + strings.Join(missingStr, ", "))
	}

	// Load all approved specs
	input := reconcile.ReconcileInput{
		ProjectName: project.Name,
	}
	if content, err := os.ReadFile(config.SpecPath(projectPath, types.SpecTypeMRD)); err == nil {
		input.MRD = string(content)
	}
	if content, err := os.ReadFile(config.SpecPath(projectPath, types.SpecTypePRD)); err == nil {
		input.PRD = string(content)
	}
	if content, err := os.ReadFile(config.SpecPath(projectPath, types.SpecTypeUXD)); err == nil {
		input.UXD = string(content)
	}
	if content, err := os.ReadFile(config.SpecPath(projectPath, types.SpecTypeTRD)); err == nil {
		input.TRD = string(content)
	}
	if content, err := os.ReadFile(config.SpecPath(projectPath, types.SpecTypeIRD)); err == nil {
		input.IRD = string(content)
	}

	// Load constitution if exists
	constitutionPath := filepath.Join(projectPath, "..", "CONSTITUTION.md")
	if content, err := os.ReadFile(constitutionPath); err == nil {
		input.Constitution = string(content)
	}

	// Create LLM client
	llmClient, err := eval.NewLLMClientFromProject(project.LLM)
	if err != nil {
		return errorResult("failed to initialize LLM: " + err.Error())
	}
	defer func() { _ = llmClient.Close() }()

	// Create reconciler
	reconciler := reconcile.NewReconciler(&reconcileLLMAdapter{client: llmClient})

	// Run reconciliation
	reconcileResult, err := reconciler.Reconcile(ctx, input)
	if err != nil {
		return errorResult("reconciliation failed: " + err.Error())
	}

	// Write output to spec.md
	outputPath := config.SpecPath(projectPath, types.SpecTypeSpec)
	if err := os.WriteFile(outputPath, []byte(reconcileResult.Content), 0600); err != nil {
		return errorResult("failed to write spec.md: " + err.Error())
	}

	result := map[string]any{
		"project":      args.Project,
		"path":         outputPath,
		"sources":      reconcileResult.Sources,
		"approved":     approved,
		"conflicts":    reconcileResult.Conflicts,
		"generated_at": reconcileResult.GeneratedAt.Format(time.RFC3339),
		"message":      "Reconciliation completed successfully",
	}
	data, _ := json.Marshal(result)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

// reconcileLLMAdapter adapts eval.LLMClient to reconcile.LLMClient interface.
type reconcileLLMAdapter struct {
	client *eval.LLMClient
}

func (a *reconcileLLMAdapter) Complete(ctx context.Context, prompt string) (string, error) {
	content, _, err := a.client.Complete(ctx, prompt)
	return content, err
}

func (s *Server) handleApprove(ctx context.Context, req *mcp.CallToolRequest, args ApproveArgs) (*mcp.CallToolResult, any, error) {
	// Parse spec type
	specType := types.SpecType(strings.ToLower(args.SpecType))

	// Validate spec type is approvable
	if !specType.IsValid() {
		return errorResult("invalid spec type: " + args.SpecType)
	}

	// Get project path
	projectPath, err := getProjectPath(args.Project)
	if err != nil {
		return errorResult("failed to find project: " + err.Error())
	}

	// Load project config
	project, err := config.Load(projectPath)
	if err != nil {
		return errorResult("failed to load project config: " + err.Error())
	}

	// Check that spec exists
	specPath := config.SpecPath(projectPath, specType)
	if _, err := os.Stat(specPath); os.IsNotExist(err) {
		return errorResult("spec not found: " + args.SpecType + " - cannot approve non-existent spec")
	}

	// Initialize approvals map if needed
	if project.Approvals == nil {
		project.Approvals = make(map[types.SpecType]*types.Approval)
	}

	// Determine approver
	approver := args.Approver
	if approver == "" {
		approver = "unknown"
	}

	// Record approval
	project.Approvals[specType] = &types.Approval{
		Approver:   approver,
		ApprovedAt: time.Now(),
	}
	project.UpdatedAt = time.Now()

	// Save project config
	if err := config.Save(project); err != nil {
		return errorResult("failed to save approval: " + err.Error())
	}

	result := map[string]any{
		"project":     args.Project,
		"spec_type":   string(specType),
		"approver":    approver,
		"approved_at": project.Approvals[specType].ApprovedAt.Format(time.RFC3339),
		"message":     "Spec approved successfully",
	}
	data, _ := json.Marshal(result)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

func (s *Server) handleExport(ctx context.Context, req *mcp.CallToolRequest, args ExportArgs) (*mcp.CallToolResult, any, error) {
	// Get project path
	projectPath, err := getProjectPath(args.Project)
	if err != nil {
		return errorResult("failed to find project: " + err.Error())
	}

	// Load project config
	project, err := config.Load(projectPath)
	if err != nil {
		return errorResult("failed to load project config: " + err.Error())
	}

	// Read spec.md
	specPath := config.SpecPath(projectPath, types.SpecTypeSpec)
	specContent, err := os.ReadFile(specPath)
	if err != nil {
		if os.IsNotExist(err) {
			return errorResult("spec.md not found - run reconcile first")
		}
		return errorResult("failed to read spec.md: " + err.Error())
	}

	// Get target adapter
	t, err := target.Get(args.Target)
	if err != nil {
		available := target.Available()
		return errorResult(fmt.Sprintf("unknown target: %s (available: %v)", args.Target, available))
	}

	// Get export config from project
	exportConfig := target.ProjectTargetConfig(project, args.Target)
	if exportConfig.OutputDir == "" {
		exportConfig.OutputDir = filepath.Join(projectPath, "export", args.Target)
	}

	// Run export
	exportResult, err := t.Export(string(specContent), *exportConfig)
	if err != nil {
		return errorResult("export failed: " + err.Error())
	}

	result := map[string]any{
		"project":    args.Project,
		"target":     args.Target,
		"output_dir": exportResult.OutputDir,
		"files":      exportResult.Files,
		"message":    exportResult.Message,
	}
	data, _ := json.Marshal(result)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

// Draft authoring handlers

// getProjectPath resolves the project path from the project name.
func getProjectPath(projectName string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	specsDir, err := config.FindSpecsDir(cwd)
	if err != nil {
		return "", err
	}

	return config.ProjectPath(specsDir, projectName), nil
}

// parseSpecType parses and validates a spec type string.
func parseSpecType(s string) (types.SpecType, error) {
	st := types.SpecType(strings.ToLower(s))

	// Validate it's a source spec type with a template
	if !templates.HasTemplate(st) {
		return "", errors.New("spec type must be mrd, prd, or uxd")
	}

	return st, nil
}

func (s *Server) handleStartDraft(ctx context.Context, req *mcp.CallToolRequest, args StartDraftArgs) (*mcp.CallToolResult, any, error) {
	// Validate spec type
	specType, err := parseSpecType(args.SpecType)
	if err != nil {
		return errorResult(err.Error())
	}

	// Get project path
	projectPath, err := getProjectPath(args.Project)
	if err != nil {
		return errorResult("failed to find project: " + err.Error())
	}

	// Start or resume session
	session, err := draft.StartSession(projectPath, specType)
	if err != nil {
		return errorResult("failed to start draft: " + err.Error())
	}

	result := map[string]any{
		"project":      args.Project,
		"spec_type":    string(specType),
		"is_new":       session.IsNew(),
		"version":      session.Version(),
		"status":       string(session.Status()),
		"instructions": session.Instructions(),
		"content":      session.Content(),
	}
	data, _ := json.Marshal(result)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

func (s *Server) handleGetDraft(ctx context.Context, req *mcp.CallToolRequest, args DraftArgs) (*mcp.CallToolResult, any, error) {
	// Validate spec type
	specType, err := parseSpecType(args.SpecType)
	if err != nil {
		return errorResult(err.Error())
	}

	// Get project path
	projectPath, err := getProjectPath(args.Project)
	if err != nil {
		return errorResult("failed to find project: " + err.Error())
	}

	// Get draft
	d, err := draft.Get(projectPath, specType)
	if err != nil {
		if errors.Is(err, draft.ErrDraftNotFound) {
			return errorResult("no draft found for " + args.SpecType)
		}
		return errorResult("failed to get draft: " + err.Error())
	}

	result := map[string]any{
		"project":      args.Project,
		"spec_type":    string(specType),
		"version":      d.Metadata.Version,
		"started_at":   d.Metadata.StartedAt,
		"updated_at":   d.Metadata.UpdatedAt,
		"eval_count":   len(d.Metadata.EvalHistory),
		"eval_history": d.Metadata.EvalHistory,
		"content":      d.Content,
	}
	data, _ := json.Marshal(result)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

func (s *Server) handleUpdateDraft(ctx context.Context, req *mcp.CallToolRequest, args UpdateDraftArgs) (*mcp.CallToolResult, any, error) {
	// Validate spec type
	specType, err := parseSpecType(args.SpecType)
	if err != nil {
		return errorResult(err.Error())
	}

	// Get project path
	projectPath, err := getProjectPath(args.Project)
	if err != nil {
		return errorResult("failed to find project: " + err.Error())
	}

	// Update draft
	d, err := draft.Update(projectPath, specType, args.Content)
	if err != nil {
		if errors.Is(err, draft.ErrDraftNotFound) {
			return errorResult("no draft found - use start_draft first")
		}
		return errorResult("failed to update draft: " + err.Error())
	}

	result := map[string]any{
		"project":    args.Project,
		"spec_type":  string(specType),
		"version":    d.Metadata.Version,
		"updated_at": d.Metadata.UpdatedAt,
		"message":    "Draft updated successfully",
	}
	data, _ := json.Marshal(result)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

func (s *Server) handleEvalDraft(ctx context.Context, req *mcp.CallToolRequest, args DraftArgs) (*mcp.CallToolResult, any, error) {
	// Validate spec type
	specType, err := parseSpecType(args.SpecType)
	if err != nil {
		return errorResult(err.Error())
	}

	// Get project path
	projectPath, err := getProjectPath(args.Project)
	if err != nil {
		return errorResult("failed to find project: " + err.Error())
	}

	// Load project config for LLM settings
	project, err := config.Load(projectPath)
	if err != nil {
		return errorResult("failed to load project config: " + err.Error())
	}

	// Get draft content
	d, err := draft.Get(projectPath, specType)
	if err != nil {
		if errors.Is(err, draft.ErrDraftNotFound) {
			return errorResult("no draft found - use start_draft first")
		}
		return errorResult("failed to get draft: " + err.Error())
	}

	// Create LLM client from project config (falls back to environment)
	llmClient, err := eval.NewLLMClientFromProject(project.LLM)
	if err != nil {
		return errorResult("failed to initialize LLM: " + err.Error())
	}
	defer func() { _ = llmClient.Close() }()

	// Run evaluation
	evaluator := eval.NewEvaluator(llmClient)
	evalResult, err := evaluator.Evaluate(ctx, specType, d.Content)
	if err != nil {
		return errorResult("evaluation failed: " + err.Error())
	}

	// Record eval result in draft metadata
	if err := draft.AddEvalResult(projectPath, specType, evalResult.Score, evalResult.Passed, len(evalResult.Findings)); err != nil {
		// Log but don't fail - evaluation was successful
		_ = err
	}

	result := map[string]any{
		"project":    args.Project,
		"spec_type":  string(specType),
		"score":      evalResult.Score,
		"passed":     evalResult.Passed,
		"decision":   evalResult.Decision,
		"categories": evalResult.Categories,
		"findings":   evalResult.Findings,
		"summary":    evalResult.Summary,
		"judge":      evalResult.Judge,
	}
	data, _ := json.Marshal(result)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

func (s *Server) handleFinalizeDraft(ctx context.Context, req *mcp.CallToolRequest, args DraftArgs) (*mcp.CallToolResult, any, error) {
	// Validate spec type
	specType, err := parseSpecType(args.SpecType)
	if err != nil {
		return errorResult(err.Error())
	}

	// Get project path
	projectPath, err := getProjectPath(args.Project)
	if err != nil {
		return errorResult("failed to find project: " + err.Error())
	}

	// Check draft exists
	if !draft.Exists(projectPath, specType) {
		return errorResult("no draft found - use start_draft first")
	}

	// Finalize
	if err := draft.Finalize(projectPath, specType); err != nil {
		return errorResult("failed to finalize draft: " + err.Error())
	}

	specPath := config.SpecPath(projectPath, specType)
	result := map[string]any{
		"project":   args.Project,
		"spec_type": string(specType),
		"spec_path": specPath,
		"message":   "Draft finalized successfully",
	}
	data, _ := json.Marshal(result)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

func (s *Server) handleDiscardDraft(ctx context.Context, req *mcp.CallToolRequest, args DraftArgs) (*mcp.CallToolResult, any, error) {
	// Validate spec type
	specType, err := parseSpecType(args.SpecType)
	if err != nil {
		return errorResult(err.Error())
	}

	// Get project path
	projectPath, err := getProjectPath(args.Project)
	if err != nil {
		return errorResult("failed to find project: " + err.Error())
	}

	// Discard draft
	if err := draft.Discard(projectPath, specType); err != nil {
		if errors.Is(err, draft.ErrDraftNotFound) {
			return errorResult("no draft found for " + args.SpecType)
		}
		return errorResult("failed to discard draft: " + err.Error())
	}

	result := map[string]any{
		"project":   args.Project,
		"spec_type": string(specType),
		"message":   "Draft discarded successfully",
	}
	data, _ := json.Marshal(result)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

func (s *Server) handleListDrafts(ctx context.Context, req *mcp.CallToolRequest, args ProjectArgs) (*mcp.CallToolResult, any, error) {
	// Get project path
	projectPath, err := getProjectPath(args.Project)
	if err != nil {
		return errorResult("failed to find project: " + err.Error())
	}

	// List drafts
	drafts, err := draft.ListDrafts(projectPath)
	if err != nil {
		return errorResult("failed to list drafts: " + err.Error())
	}

	// Build response
	var draftList []map[string]any
	for _, d := range drafts {
		draftInfo := map[string]any{
			"spec_type":   string(d.Metadata.SpecType),
			"version":     d.Metadata.Version,
			"started_at":  d.Metadata.StartedAt.Format(time.RFC3339),
			"updated_at":  d.Metadata.UpdatedAt.Format(time.RFC3339),
			"eval_count":  len(d.Metadata.EvalHistory),
			"content_len": len(d.Content),
		}
		// Include latest eval result if available
		if len(d.Metadata.EvalHistory) > 0 {
			latest := d.Metadata.EvalHistory[len(d.Metadata.EvalHistory)-1]
			draftInfo["latest_eval"] = map[string]any{
				"score":  latest.Score,
				"passed": latest.Passed,
			}
		}
		draftList = append(draftList, draftInfo)
	}

	result := map[string]any{
		"project": args.Project,
		"drafts":  draftList,
		"count":   len(draftList),
	}
	data, _ := json.Marshal(result)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

// AlignArgs contains align command arguments.
type AlignArgs struct {
	Project       string `json:"project" jsonschema:"description=Project name"`
	GenerateTruth bool   `json:"generate_truth,omitempty" jsonschema:"description=Generate current-truth.md document"`
}

func (s *Server) handleAlign(ctx context.Context, req *mcp.CallToolRequest, args AlignArgs) (*mcp.CallToolResult, any, error) {
	// Get project path
	projectPath, err := getProjectPath(args.Project)
	if err != nil {
		return errorResult("failed to find project: " + err.Error())
	}

	// Load project config
	project, err := config.Load(projectPath)
	if err != nil {
		return errorResult("failed to load project config: " + err.Error())
	}

	// Read spec.md
	specPath := filepath.Join(projectPath, "spec.md")
	specContent, err := os.ReadFile(specPath)
	if err != nil {
		if os.IsNotExist(err) {
			return errorResult("spec.md not found - run reconcile first")
		}
		return errorResult("failed to read spec.md: " + err.Error())
	}

	// Load context from cache if available
	var aggregatedCtx *ctxpkg.AggregatedContext
	contextCachePath := filepath.Join(projectPath, ".visionspec", "context-cache.json")
	if contextData, err := os.ReadFile(contextCachePath); err == nil {
		var ctxData ctxpkg.AggregatedContext
		if json.Unmarshal(contextData, &ctxData) == nil {
			aggregatedCtx = &ctxData
		}
	}

	if aggregatedCtx == nil {
		aggregatedCtx = &ctxpkg.AggregatedContext{
			Project: project.Name,
		}
	}

	// Run alignment check
	aligner := align.NewAligner()
	opts := align.AlignOptions{
		MinSeverity:     align.SeverityLow,
		IncludeEvidence: true,
	}

	alignResult, err := aligner.Align(string(specContent), aggregatedCtx, opts)
	if err != nil {
		return errorResult("alignment check failed: " + err.Error())
	}

	// Build result
	result := map[string]any{
		"project":         args.Project,
		"generated_at":    alignResult.GeneratedAt.Format(time.RFC3339),
		"alignment_score": alignResult.Summary.AlignmentScore,
		"is_aligned":      alignResult.Summary.IsAligned,
		"summary": map[string]any{
			"total_discrepancies": alignResult.Summary.TotalDiscrepancies,
			"critical_count":      alignResult.Summary.CriticalCount,
			"high_count":          alignResult.Summary.HighCount,
		},
		"coverage": map[string]any{
			"total_requirements":  alignResult.Coverage.TotalRequirements,
			"implemented":         alignResult.Coverage.ImplementedCount,
			"partial":             alignResult.Coverage.PartialCount,
			"missing":             alignResult.Coverage.MissingCount,
			"coverage_percentage": alignResult.Coverage.CoveragePercentage,
		},
		"discrepancies": alignResult.Discrepancies,
	}

	// Generate current-truth.md if requested
	if args.GenerateTruth {
		truth := align.GenerateCurrentTruth(alignResult)
		truthContent, err := truth.RenderMarkdown()
		if err != nil {
			return errorResult("failed to render current-truth: " + err.Error())
		}

		truthPath := filepath.Join(projectPath, "current-truth.md")
		if err := os.WriteFile(truthPath, []byte(truthContent), 0600); err != nil { //nolint:gosec // G703: projectPath is from config, not user input
			return errorResult("failed to write current-truth.md: " + err.Error())
		}
		result["truth_path"] = truthPath
	}

	data, _ := json.Marshal(result)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

func (s *Server) handleGetResolutionPlan(ctx context.Context, req *mcp.CallToolRequest, args AlignArgs) (*mcp.CallToolResult, any, error) {
	// Get project path
	projectPath, err := getProjectPath(args.Project)
	if err != nil {
		return errorResult("failed to find project: " + err.Error())
	}

	// Load project config
	project, err := config.Load(projectPath)
	if err != nil {
		return errorResult("failed to load project config: " + err.Error())
	}

	// Read spec.md
	specPath := filepath.Join(projectPath, "spec.md")
	specContent, err := os.ReadFile(specPath)
	if err != nil {
		if os.IsNotExist(err) {
			return errorResult("spec.md not found - run reconcile first")
		}
		return errorResult("failed to read spec.md: " + err.Error())
	}

	// Load context from cache if available
	var aggregatedCtx *ctxpkg.AggregatedContext
	contextCachePath := filepath.Join(projectPath, ".visionspec", "context-cache.json")
	if contextData, err := os.ReadFile(contextCachePath); err == nil {
		var ctxData ctxpkg.AggregatedContext
		if json.Unmarshal(contextData, &ctxData) == nil {
			aggregatedCtx = &ctxData
		}
	}

	if aggregatedCtx == nil {
		aggregatedCtx = &ctxpkg.AggregatedContext{
			Project: project.Name,
		}
	}

	// Run alignment check
	aligner := align.NewAligner()
	opts := align.AlignOptions{
		MinSeverity:     align.SeverityLow,
		IncludeEvidence: true,
	}

	alignResult, err := aligner.Align(string(specContent), aggregatedCtx, opts)
	if err != nil {
		return errorResult("alignment check failed: " + err.Error())
	}

	// Generate resolution plan
	engine := align.NewResolutionEngine()
	plan := engine.GeneratePlan(alignResult)

	result := map[string]any{
		"project":      args.Project,
		"generated_at": plan.GeneratedAt.Format(time.RFC3339),
		"summary": map[string]any{
			"total":       plan.Summary.TotalDiscrepancies,
			"update_spec": plan.Summary.UpdateSpec,
			"update_code": plan.Summary.UpdateCode,
			"add_spec":    plan.Summary.AddSpec,
			"deferred":    plan.Summary.Deferred,
		},
		"priorities":  plan.Priorities,
		"resolutions": plan.Resolutions,
		"progress":    plan.GetProgress(),
	}

	data, _ := json.Marshal(result)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

func (s *Server) handleGetTestCoverage(ctx context.Context, req *mcp.CallToolRequest, args ProjectArgs) (*mcp.CallToolResult, any, error) {
	// Get project path
	projectPath, err := getProjectPath(args.Project)
	if err != nil {
		return errorResult("failed to find project: " + err.Error())
	}

	// Find the repository root (one level up from specs dir usually)
	repoRoot := filepath.Dir(filepath.Dir(projectPath))

	// Create test mapper
	mapper := testmap.NewMapper(repoRoot)
	mapping, err := mapper.Map()
	if err != nil {
		return errorResult("failed to generate test coverage: " + err.Error())
	}

	result := map[string]any{
		"project":      args.Project,
		"generated_at": mapping.GeneratedAt.Format(time.RFC3339),
		"summary": map[string]any{
			"total_requirements":   mapping.Summary.TotalRequirements,
			"covered":              mapping.Summary.CoveredRequirements,
			"partial":              mapping.Summary.PartialRequirements,
			"uncovered":            mapping.Summary.UncoveredRequirements,
			"overall_coverage_pct": mapping.Summary.OverallCoverage,
			"total_tests":          mapping.Summary.TotalTests,
			"mapped_tests":         mapping.Summary.MappedTests,
		},
		"requirements": mapping.Requirements,
		"unmapped": map[string]any{
			"requirements": mapping.Unmapped.Requirements,
			"tests":        mapping.Unmapped.Tests,
		},
	}

	data, _ := json.Marshal(result)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

// ExecutionContextArgs contains arguments for get_execution_context.
type ExecutionContextArgs struct {
	Project     string `json:"project" jsonschema:"description=Project name"`
	IncludeSpec bool   `json:"include_spec,omitempty" jsonschema:"description=Include full spec content"`
	IncludeTRD  bool   `json:"include_trd,omitempty" jsonschema:"description=Include technical requirements"`
}

func (s *Server) handleGetExecutionContext(ctx context.Context, req *mcp.CallToolRequest, args ExecutionContextArgs) (*mcp.CallToolResult, any, error) {
	// Get project path
	projectPath, err := getProjectPath(args.Project)
	if err != nil {
		return errorResult("failed to find project: " + err.Error())
	}

	// Load project config
	project, err := config.Load(projectPath)
	if err != nil {
		return errorResult("failed to load project config: " + err.Error())
	}

	// Build execution context
	execContext := map[string]any{
		"project": project.Name,
	}

	// Get project status
	projectStatus, err := status.Generate(project)
	if err == nil {
		execContext["status"] = map[string]any{
			"readiness": projectStatus.Readiness,
			"summary": map[string]any{
				"total":     projectStatus.Summary.TotalSpecs,
				"present":   projectStatus.Summary.PresentSpecs,
				"approved":  projectStatus.Summary.ApprovedSpecs,
				"evaluated": projectStatus.Summary.EvaluatedSpecs,
			},
		}
	}

	// Load spec.md if requested
	if args.IncludeSpec {
		specPath := filepath.Join(projectPath, "spec.md")
		if specContent, err := os.ReadFile(specPath); err == nil {
			execContext["spec"] = string(specContent)
		}
	}

	// Load TRD if requested
	if args.IncludeTRD {
		trdPath := config.SpecPath(projectPath, types.SpecTypeTRD)
		if trdContent, err := os.ReadFile(trdPath); err == nil {
			execContext["trd"] = string(trdContent)
		}
	}

	// Load cached context snapshot if available
	contextCachePath := filepath.Join(projectPath, ".visionspec", "context-cache.json")
	if contextData, err := os.ReadFile(contextCachePath); err == nil {
		var ctxData ctxpkg.AggregatedContext
		if json.Unmarshal(contextData, &ctxData) == nil {
			execContext["codebase"] = map[string]any{
				"has_code":  ctxData.HasCode,
				"has_graph": ctxData.HasGraph,
				"summary":   ctxData.Summary,
			}
		}
	}

	// Generate execution guidance
	guidance := []string{
		"1. Review the spec.md for complete requirements",
		"2. Check TRD for technical implementation details",
		"3. Use 'align' tool to check implementation status",
		"4. Use 'get_resolution_plan' for prioritized actions",
	}
	execContext["guidance"] = guidance

	data, _ := json.Marshal(execContext)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

// ExecutionStatusArgs contains arguments for get_execution_status.
type ExecutionStatusArgs struct {
	Project string `json:"project" jsonschema:"description=Project name"`
}

func (s *Server) handleGetExecutionStatus(ctx context.Context, req *mcp.CallToolRequest, args ExecutionStatusArgs) (*mcp.CallToolResult, any, error) {
	// Get project path
	projectPath, err := getProjectPath(args.Project)
	if err != nil {
		return errorResult("failed to find project: " + err.Error())
	}

	// Load execution status file
	statusPath := filepath.Join(projectPath, ".visionspec", "execution-status.json")
	var execStatus ExecutionStatus
	if statusData, err := os.ReadFile(statusPath); err == nil {
		_ = json.Unmarshal(statusData, &execStatus)
	} else {
		// Initialize with defaults
		execStatus = ExecutionStatus{
			Project:      args.Project,
			Requirements: make(map[string]RequirementStatus),
		}
	}

	// Load spec.md to extract requirements if not tracked
	specPath := filepath.Join(projectPath, "spec.md")
	if specContent, err := os.ReadFile(specPath); err == nil {
		reqs := extractRequirementsFromSpec(string(specContent))
		for _, reqID := range reqs {
			if _, exists := execStatus.Requirements[reqID]; !exists {
				execStatus.Requirements[reqID] = RequirementStatus{
					ID:     reqID,
					Status: "pending",
				}
			}
		}
	}

	// Calculate summary
	summary := ExecutionSummary{}
	for _, req := range execStatus.Requirements {
		switch req.Status {
		case "implemented":
			summary.Implemented++
		case "in_progress":
			summary.InProgress++
		case "blocked":
			summary.Blocked++
		default:
			summary.Pending++
		}
	}
	summary.Total = len(execStatus.Requirements)
	if summary.Total > 0 {
		summary.Progress = float64(summary.Implemented) / float64(summary.Total) * 100
	}

	result := map[string]any{
		"project":      args.Project,
		"updated_at":   execStatus.UpdatedAt.Format(time.RFC3339),
		"summary":      summary,
		"requirements": execStatus.Requirements,
	}

	data, _ := json.Marshal(result)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

// TrackRequirementArgs contains arguments for track_requirement.
type TrackRequirementArgs struct {
	Project       string `json:"project" jsonschema:"description=Project name"`
	RequirementID string `json:"requirement_id" jsonschema:"description=Requirement identifier (e.g., REQ-001)"`
	Status        string `json:"status" jsonschema:"description=New status: pending, in_progress, implemented, blocked"`
	Evidence      string `json:"evidence,omitempty" jsonschema:"description=Implementation evidence (file paths, commit refs)"`
	Notes         string `json:"notes,omitempty" jsonschema:"description=Implementation notes"`
}

func (s *Server) handleTrackRequirement(ctx context.Context, req *mcp.CallToolRequest, args TrackRequirementArgs) (*mcp.CallToolResult, any, error) {
	// Validate status
	validStatuses := map[string]bool{"pending": true, "in_progress": true, "implemented": true, "blocked": true}
	if !validStatuses[args.Status] {
		return errorResult("invalid status: must be pending, in_progress, implemented, or blocked")
	}

	// Get project path
	projectPath, err := getProjectPath(args.Project)
	if err != nil {
		return errorResult("failed to find project: " + err.Error())
	}

	// Load or create execution status
	statusPath := filepath.Join(projectPath, ".visionspec", "execution-status.json")
	var execStatus ExecutionStatus
	if statusData, err := os.ReadFile(statusPath); err == nil {
		_ = json.Unmarshal(statusData, &execStatus)
	} else {
		execStatus = ExecutionStatus{
			Project:      args.Project,
			Requirements: make(map[string]RequirementStatus),
		}
	}

	// Update requirement status
	reqStatus := execStatus.Requirements[args.RequirementID]
	reqStatus.ID = args.RequirementID
	reqStatus.Status = args.Status
	reqStatus.UpdatedAt = time.Now()
	if args.Evidence != "" {
		reqStatus.Evidence = append(reqStatus.Evidence, args.Evidence)
	}
	if args.Notes != "" {
		reqStatus.Notes = args.Notes
	}
	execStatus.Requirements[args.RequirementID] = reqStatus
	execStatus.UpdatedAt = time.Now()

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(statusPath), 0755); err != nil {
		return errorResult("failed to create status directory: " + err.Error())
	}

	// Save status
	statusData, _ := json.MarshalIndent(execStatus, "", "  ")
	if err := os.WriteFile(statusPath, statusData, 0644); err != nil {
		return errorResult("failed to save execution status: " + err.Error())
	}

	result := map[string]any{
		"project":        args.Project,
		"requirement_id": args.RequirementID,
		"status":         args.Status,
		"updated_at":     reqStatus.UpdatedAt.Format(time.RFC3339),
		"message":        "Requirement status updated successfully",
	}

	data, _ := json.Marshal(result)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

// ExecutionStatus tracks implementation progress.
type ExecutionStatus struct {
	Project      string                       `json:"project"`
	UpdatedAt    time.Time                    `json:"updated_at"`
	Requirements map[string]RequirementStatus `json:"requirements"`
}

// RequirementStatus tracks a single requirement's implementation status.
type RequirementStatus struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"` // pending, in_progress, implemented, blocked
	UpdatedAt time.Time `json:"updated_at"`
	Evidence  []string  `json:"evidence,omitempty"`
	Notes     string    `json:"notes,omitempty"`
}

// ExecutionSummary provides aggregate execution statistics.
type ExecutionSummary struct {
	Total       int     `json:"total"`
	Implemented int     `json:"implemented"`
	InProgress  int     `json:"in_progress"`
	Blocked     int     `json:"blocked"`
	Pending     int     `json:"pending"`
	Progress    float64 `json:"progress_pct"`
}

// extractRequirementsFromSpec extracts requirement IDs from spec content.
func extractRequirementsFromSpec(content string) []string {
	var reqs []string
	seen := make(map[string]bool)

	// Look for REQ-XXX patterns
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		// Extract REQ-XXX identifiers
		for i := 0; i < len(line); i++ {
			if i+4 < len(line) && line[i:i+4] == "REQ-" {
				// Find the end of the ID
				end := i + 4
				for end < len(line) && (line[end] >= '0' && line[end] <= '9' || line[end] == '-' || (line[end] >= 'A' && line[end] <= 'Z')) {
					end++
				}
				if end > i+4 {
					reqID := line[i:end]
					if !seen[reqID] {
						seen[reqID] = true
						reqs = append(reqs, reqID)
					}
				}
			}
		}
	}
	return reqs
}

// registerPrompts registers MCP prompts for guided implementation.
func (s *Server) registerPrompts() {
	// Prompt: implement_requirement
	s.server.AddPrompt(&mcp.Prompt{
		Name:        "implement_requirement",
		Description: "Guide implementation of a specific requirement from the spec",
		Arguments: []*mcp.PromptArgument{
			{Name: "project", Description: "Project name", Required: true},
			{Name: "requirement_id", Description: "Requirement ID (e.g., REQ-001)", Required: true},
		},
	}, s.handleImplementRequirementPrompt)

	// Prompt: verify_acceptance
	s.server.AddPrompt(&mcp.Prompt{
		Name:        "verify_acceptance",
		Description: "Guide verification of acceptance criteria for a requirement",
		Arguments: []*mcp.PromptArgument{
			{Name: "project", Description: "Project name", Required: true},
			{Name: "requirement_id", Description: "Requirement ID to verify", Required: true},
		},
	}, s.handleVerifyAcceptancePrompt)

	// Prompt: resolve_drift
	s.server.AddPrompt(&mcp.Prompt{
		Name:        "resolve_drift",
		Description: "Guide resolution of detected drift between spec and implementation",
		Arguments: []*mcp.PromptArgument{
			{Name: "project", Description: "Project name", Required: true},
			{Name: "category", Description: "Drift category: missing_feature, undocumented_code, or diverged", Required: false},
		},
	}, s.handleResolveDriftPrompt)
}

func (s *Server) handleImplementRequirementPrompt(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	project := req.Params.Arguments["project"]
	reqID := req.Params.Arguments["requirement_id"]

	// Get project path
	projectPath, err := getProjectPath(project)
	if err != nil {
		return nil, fmt.Errorf("failed to find project: %w", err)
	}

	// Load spec.md
	specPath := filepath.Join(projectPath, "spec.md")
	specContent, err := os.ReadFile(specPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read spec: %w", err)
	}

	// Load TRD for technical details
	trdPath := config.SpecPath(projectPath, types.SpecTypeTRD)
	trdContent, _ := os.ReadFile(trdPath)

	promptContent := fmt.Sprintf(`# Implement Requirement: %s

## Context
You are implementing requirement %s from the visionspec project "%s".

## Specification
%s

## Technical Requirements
%s

## Implementation Guidance

1. **Understand the Requirement**: Review the specification section for %s
2. **Check Technical Details**: Review the TRD for implementation approach
3. **Implement Incrementally**: Break down into small, testable changes
4. **Track Progress**: Use 'track_requirement' tool to update status
5. **Verify Acceptance**: Ensure all acceptance criteria are met

## Next Steps
- Mark requirement as 'in_progress' when you start
- Implement the feature following TRD architecture
- Write tests to verify acceptance criteria
- Mark requirement as 'implemented' when done

Please proceed with implementing %s.`,
		reqID, reqID, project,
		string(specContent),
		string(trdContent),
		reqID, reqID)

	return &mcp.GetPromptResult{
		Description: fmt.Sprintf("Implementation guide for %s in %s", reqID, project),
		Messages: []*mcp.PromptMessage{
			{
				Role:    "user",
				Content: &mcp.TextContent{Text: promptContent},
			},
		},
	}, nil
}

func (s *Server) handleVerifyAcceptancePrompt(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	project := req.Params.Arguments["project"]
	reqID := req.Params.Arguments["requirement_id"]

	// Get project path
	projectPath, err := getProjectPath(project)
	if err != nil {
		return nil, fmt.Errorf("failed to find project: %w", err)
	}

	// Load spec.md
	specPath := filepath.Join(projectPath, "spec.md")
	specContent, err := os.ReadFile(specPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read spec: %w", err)
	}

	promptContent := fmt.Sprintf(`# Verify Acceptance Criteria: %s

## Context
You are verifying that requirement %s has been correctly implemented in project "%s".

## Specification
%s

## Verification Steps

1. **Locate Acceptance Criteria**: Find the acceptance criteria for %s in the spec
2. **Check Implementation**: Verify each criterion is implemented
3. **Run Tests**: Ensure tests pass for this requirement
4. **Document Evidence**: Record implementation evidence

## Verification Checklist
- [ ] All acceptance criteria identified
- [ ] Implementation matches each criterion
- [ ] Tests exist and pass
- [ ] Edge cases handled
- [ ] Documentation updated

## Recording Results
Use 'track_requirement' to update status to 'implemented' with evidence of completion.

Please verify the implementation of %s.`,
		reqID, reqID, project,
		string(specContent),
		reqID, reqID)

	return &mcp.GetPromptResult{
		Description: fmt.Sprintf("Acceptance verification guide for %s in %s", reqID, project),
		Messages: []*mcp.PromptMessage{
			{
				Role:    "user",
				Content: &mcp.TextContent{Text: promptContent},
			},
		},
	}, nil
}

func (s *Server) handleResolveDriftPrompt(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	project := req.Params.Arguments["project"]
	category := req.Params.Arguments["category"]
	if category == "" {
		category = "all"
	}

	// Get project path
	projectPath, err := getProjectPath(project)
	if err != nil {
		return nil, fmt.Errorf("failed to find project: %w", err)
	}

	// Load spec.md
	specPath := filepath.Join(projectPath, "spec.md")
	specContent, err := os.ReadFile(specPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read spec: %w", err)
	}

	promptContent := fmt.Sprintf(`# Resolve Specification Drift

## Context
You are resolving drift between the specification and implementation in project "%s".
Focus category: %s

## Current Specification
%s

## Drift Resolution Strategy

### For Missing Features (spec exists, implementation doesn't)
1. Evaluate if feature is still needed
2. If yes: implement the feature
3. If no: update spec to remove requirement

### For Undocumented Code (implementation exists, not in spec)
1. Evaluate if code is needed
2. If yes: add to specification
3. If no: remove the code

### For Diverged (both exist but differ)
1. Determine correct behavior
2. Update implementation OR spec to match
3. Document the decision

## Tools Available
- 'align' - Check current alignment status
- 'get_resolution_plan' - Get prioritized resolution actions
- 'track_requirement' - Update requirement status

## Next Steps
1. Run 'align' to get current discrepancies
2. Run 'get_resolution_plan' for prioritized actions
3. Resolve each discrepancy following the strategy above
4. Re-run 'align' to verify resolution

Please analyze and resolve drift in project "%s".`,
		project, category,
		string(specContent),
		project)

	return &mcp.GetPromptResult{
		Description: fmt.Sprintf("Drift resolution guide for %s", project),
		Messages: []*mcp.PromptMessage{
			{
				Role:    "user",
				Content: &mcp.TextContent{Text: promptContent},
			},
		},
	}, nil
}
