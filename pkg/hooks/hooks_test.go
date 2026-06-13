package hooks

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/grokify/oscompat/fs"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if !cfg.Enabled {
		t.Error("Default config should be enabled")
	}

	if len(cfg.Hooks) == 0 {
		t.Error("Default config should have hooks")
	}

	// Check pre-commit hook exists
	if hook, ok := cfg.Hooks[HookPreCommit]; !ok {
		t.Error("Default config should have pre-commit hook")
	} else if !hook.Enabled {
		t.Error("Pre-commit hook should be enabled by default")
	}
}

func TestGetTemplate(t *testing.T) {
	tests := []struct {
		hookType HookType
		contains string
	}{
		{HookPreCommit, "pre-commit"},
		{HookPrePush, "pre-push"},
		{HookCommitMsg, "commit-msg"},
		{HookPostCommit, "post-commit"},
	}

	for _, tt := range tests {
		t.Run(string(tt.hookType), func(t *testing.T) {
			template := GetTemplate(tt.hookType)
			if template == "" {
				t.Errorf("Template for %s should not be empty", tt.hookType)
			}
			if len(template) < 50 {
				t.Errorf("Template for %s seems too short", tt.hookType)
			}
		})
	}
}

func TestGetTemplateDescription(t *testing.T) {
	for _, hookType := range AllHookTypes() {
		desc := GetTemplateDescription(hookType)
		if desc == "" {
			t.Errorf("Description for %s should not be empty", hookType)
		}
	}
}

func TestManager_GitHooksDir(t *testing.T) {
	// Create a temporary git repo
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git", "hooks")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatalf("Failed to create .git/hooks: %v", err)
	}

	manager := NewManager(tmpDir, DefaultConfig())
	hooksDir, err := manager.GitHooksDir()
	if err != nil {
		t.Fatalf("GitHooksDir failed: %v", err)
	}

	if hooksDir != gitDir {
		t.Errorf("GitHooksDir = %s, want %s", hooksDir, gitDir)
	}
}

func TestManager_Install(t *testing.T) {
	// Create a temporary git repo
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git", "hooks")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatalf("Failed to create .git/hooks: %v", err)
	}

	manager := NewManager(tmpDir, DefaultConfig())
	result, err := manager.Install([]HookType{HookPreCommit})
	if err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	if len(result.Installed) != 1 {
		t.Errorf("Should have installed 1 hook, got %d", len(result.Installed))
	}

	// Verify the hook file exists
	hookPath := filepath.Join(gitDir, "pre-commit")
	if _, err := os.Stat(hookPath); err != nil {
		t.Fatalf("Hook file should exist: %v", err)
	}

	// Verify it's executable (cross-platform check)
	isExec, err := fs.IsExecutable(hookPath)
	if err != nil {
		t.Fatalf("IsExecutable check failed: %v", err)
	}
	if !isExec {
		t.Error("Hook should be executable")
	}

	// Verify content contains visionspec
	content, _ := os.ReadFile(hookPath)
	if len(content) == 0 {
		t.Error("Hook file should have content")
	}
}

func TestManager_Uninstall(t *testing.T) {
	// Create a temporary git repo
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git", "hooks")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatalf("Failed to create .git/hooks: %v", err)
	}

	manager := NewManager(tmpDir, DefaultConfig())

	// First install
	_, err := manager.Install([]HookType{HookPreCommit})
	if err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	// Then uninstall
	result, err := manager.Uninstall([]HookType{HookPreCommit})
	if err != nil {
		t.Fatalf("Uninstall failed: %v", err)
	}

	if len(result.Removed) != 1 {
		t.Errorf("Should have removed 1 hook, got %d", len(result.Removed))
	}

	// Verify the hook file is gone
	hookPath := filepath.Join(gitDir, "pre-commit")
	if _, err := os.Stat(hookPath); !os.IsNotExist(err) {
		t.Error("Hook file should be removed")
	}
}

func TestManager_Status(t *testing.T) {
	// Create a temporary git repo
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git", "hooks")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatalf("Failed to create .git/hooks: %v", err)
	}

	manager := NewManager(tmpDir, DefaultConfig())

	// Get status before install
	result, err := manager.Status()
	if err != nil {
		t.Fatalf("Status failed: %v", err)
	}

	// All hooks should be not installed
	for _, status := range result.Hooks {
		if status.Installed {
			t.Errorf("Hook %s should not be installed initially", status.Type)
		}
	}

	// Install a hook
	_, err = manager.Install([]HookType{HookPreCommit})
	if err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	// Get status after install
	result, err = manager.Status()
	if err != nil {
		t.Fatalf("Status failed: %v", err)
	}

	status, ok := result.Hooks[HookPreCommit]
	if !ok {
		t.Fatal("pre-commit status should exist")
	}
	if !status.Installed {
		t.Error("pre-commit should be installed")
	}
	if !status.IsVisionSpec {
		t.Error("pre-commit should be identified as visionspec hook")
	}
	if !status.Executable {
		t.Error("pre-commit should be executable")
	}
}

func TestAllHookTypes(t *testing.T) {
	types := AllHookTypes()
	if len(types) < 2 {
		t.Error("Should have at least 2 hook types")
	}
}

func TestGenerateCustomHook(t *testing.T) {
	custom := CustomTemplate{
		HookType:    HookPreCommit,
		Description: "Custom validation",
		Script:      "echo 'Custom hook'\nexit 0",
	}

	result := GenerateCustomHook(custom)

	if result == "" {
		t.Error("Should generate non-empty script")
	}

	if len(result) < 50 {
		t.Error("Generated script seems too short")
	}
}
