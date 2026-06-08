package workflows

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadRepo(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()

	// Create workflows directory structure
	workflowDir := filepath.Join(tmpDir, "workflows", "test-workflow", "product")
	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		t.Fatalf("failed to create workflow dir: %v", err)
	}

	// Create core-workflow.md
	workflowFile := filepath.Join(workflowDir, "core-workflow.md")
	content := []byte("# Test Workflow\n\n## Phase 1\n\nDo something.")
	if err := os.WriteFile(workflowFile, content, 0600); err != nil {
		t.Fatalf("failed to write workflow file: %v", err)
	}

	// Create rule-details directory
	ruleDetailsDir := filepath.Join(tmpDir, "rule-details")
	if err := os.MkdirAll(ruleDetailsDir, 0755); err != nil {
		t.Fatalf("failed to create rule-details dir: %v", err)
	}

	// Create templates directory
	templatesDir := filepath.Join(tmpDir, "templates", "default")
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		t.Fatalf("failed to create templates dir: %v", err)
	}

	// Create rubrics directory
	rubricsDir := filepath.Join(tmpDir, "rubrics", "default")
	if err := os.MkdirAll(rubricsDir, 0755); err != nil {
		t.Fatalf("failed to create rubrics dir: %v", err)
	}

	// Load the repo
	repo, err := LoadRepo(tmpDir)
	if err != nil {
		t.Fatalf("failed to load repo: %v", err)
	}

	// Verify workflows
	workflows := repo.ListWorkflows()
	if len(workflows) != 1 {
		t.Errorf("expected 1 workflow, got %d", len(workflows))
	}

	// Get the workflow
	wf, err := repo.GetWorkflow("test-workflow/product")
	if err != nil {
		t.Fatalf("failed to get workflow: %v", err)
	}

	if wf.Name != "test-workflow" {
		t.Errorf("expected name 'test-workflow', got '%s'", wf.Name)
	}
	if wf.Level != "product" {
		t.Errorf("expected level 'product', got '%s'", wf.Level)
	}
	if wf.ID() != "test-workflow/product" {
		t.Errorf("expected ID 'test-workflow/product', got '%s'", wf.ID())
	}
}

func TestLoadRepo_NotFound(t *testing.T) {
	tmpDir := t.TempDir()

	_, err := LoadRepo(filepath.Join(tmpDir, "nonexistent"))
	if err == nil {
		t.Error("expected error for nonexistent repo")
	}
}

func TestRepo_GetWorkflow_NotFound(t *testing.T) {
	repo := &Repo{
		Workflows: make(map[string]*Workflow),
	}

	_, err := repo.GetWorkflow("nonexistent/product")
	if err == nil {
		t.Error("expected error for nonexistent workflow")
	}
}

func TestRepo_TemplateLoader(t *testing.T) {
	tmpDir := t.TempDir()

	// Create templates directories
	defaultDir := filepath.Join(tmpDir, "templates", "default")
	if err := os.MkdirAll(defaultDir, 0755); err != nil {
		t.Fatalf("failed to create default templates dir: %v", err)
	}

	methodologyDir := filepath.Join(tmpDir, "templates", "aws-working-backwards")
	if err := os.MkdirAll(methodologyDir, 0755); err != nil {
		t.Fatalf("failed to create methodology templates dir: %v", err)
	}

	repo := &Repo{
		Path:          tmpDir,
		templatesPath: filepath.Join(tmpDir, "templates"),
		rubricsPath:   filepath.Join(tmpDir, "rubrics"),
	}

	// Test methodology-specific loader
	loader := repo.TemplateLoader("aws-working-backwards")
	if loader == nil {
		t.Error("expected non-nil template loader")
	}

	// Test fallback to default
	loader = repo.TemplateLoader("nonexistent")
	if loader == nil {
		t.Error("expected non-nil template loader for fallback")
	}
}

func TestRepo_HasExtension(t *testing.T) {
	tmpDir := t.TempDir()

	// Create extension directory
	extDir := filepath.Join(tmpDir, "extensions", "acme-corp")
	if err := os.MkdirAll(extDir, 0755); err != nil {
		t.Fatalf("failed to create extension dir: %v", err)
	}

	repo := &Repo{
		Path:           tmpDir,
		extensionsPath: filepath.Join(tmpDir, "extensions"),
	}

	if !repo.HasExtension("acme-corp") {
		t.Error("expected extension 'acme-corp' to exist")
	}

	if repo.HasExtension("nonexistent") {
		t.Error("expected extension 'nonexistent' to not exist")
	}
}

func TestIsValidRepo(t *testing.T) {
	tmpDir := t.TempDir()

	// Not valid without workflows directory
	if isValidRepo(tmpDir) {
		t.Error("expected invalid repo without workflows dir")
	}

	// Valid with workflows directory
	workflowsDir := filepath.Join(tmpDir, "workflows")
	if err := os.MkdirAll(workflowsDir, 0755); err != nil {
		t.Fatalf("failed to create workflows dir: %v", err)
	}

	if !isValidRepo(tmpDir) {
		t.Error("expected valid repo with workflows dir")
	}
}

func TestFindRepoUpward(t *testing.T) {
	tmpDir := t.TempDir()

	// Create nested directory structure
	nestedDir := filepath.Join(tmpDir, "project", "src", "pkg")
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatalf("failed to create nested dir: %v", err)
	}

	// Create spec-workflows at project level
	specWorkflowsDir := filepath.Join(tmpDir, "project", "spec-workflows", "workflows")
	if err := os.MkdirAll(specWorkflowsDir, 0755); err != nil {
		t.Fatalf("failed to create spec-workflows dir: %v", err)
	}

	// Find from nested directory should find parent spec-workflows
	found := findRepoUpward(nestedDir)
	expected := filepath.Join(tmpDir, "project", "spec-workflows")
	if found != expected {
		t.Errorf("expected %s, got %s", expected, found)
	}

	// Find from directory without spec-workflows should return empty
	emptyDir := filepath.Join(tmpDir, "empty")
	if err := os.MkdirAll(emptyDir, 0755); err != nil {
		t.Fatalf("failed to create empty dir: %v", err)
	}

	found = findRepoUpward(emptyDir)
	if found != "" {
		t.Errorf("expected empty string, got %s", found)
	}
}

func TestFindRepoUpward_HiddenDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Create nested directory structure
	nestedDir := filepath.Join(tmpDir, "project", "src")
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatalf("failed to create nested dir: %v", err)
	}

	// Create .spec-workflows (hidden) at project level
	hiddenDir := filepath.Join(tmpDir, "project", ".spec-workflows", "workflows")
	if err := os.MkdirAll(hiddenDir, 0755); err != nil {
		t.Fatalf("failed to create hidden spec-workflows dir: %v", err)
	}

	// Find from nested directory should find parent .spec-workflows
	found := findRepoUpward(nestedDir)
	expected := filepath.Join(tmpDir, "project", ".spec-workflows")
	if found != expected {
		t.Errorf("expected %s, got %s", expected, found)
	}
}

func TestDiscoverRepoPath_ExplicitPath(t *testing.T) {
	tmpDir := t.TempDir()

	// Create valid repo
	workflowsDir := filepath.Join(tmpDir, "workflows")
	if err := os.MkdirAll(workflowsDir, 0755); err != nil {
		t.Fatalf("failed to create workflows dir: %v", err)
	}

	// Explicit path should be returned when valid
	path := DiscoverRepoPath(tmpDir)
	if path != tmpDir {
		t.Errorf("expected %s, got %s", tmpDir, path)
	}

	// Note: When explicit path is invalid, DiscoverRepoPath falls back to
	// auto-discovery (env var, upward search, user config). This is correct
	// behavior - the function finds a repo, just not from the explicit path.
}

func TestDiscoverRepo_ExplicitPath(t *testing.T) {
	tmpDir := t.TempDir()

	// Create valid repo structure
	workflowDir := filepath.Join(tmpDir, "workflows", "test", "product")
	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		t.Fatalf("failed to create workflow dir: %v", err)
	}

	workflowFile := filepath.Join(workflowDir, "core-workflow.md")
	if err := os.WriteFile(workflowFile, []byte("# Test"), 0600); err != nil {
		t.Fatalf("failed to write workflow file: %v", err)
	}

	// Discover with explicit path
	repo, err := DiscoverRepo(tmpDir)
	if err != nil {
		t.Fatalf("failed to discover repo: %v", err)
	}

	if repo.Path != tmpDir {
		t.Errorf("expected path %s, got %s", tmpDir, repo.Path)
	}
}
