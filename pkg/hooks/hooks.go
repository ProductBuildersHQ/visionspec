// Package hooks provides Git hooks management for visionspec.
//
// This package enables automatic validation of specs on git operations:
// - pre-commit: Lint changed spec files
// - pre-push: Evaluate specs and check for blockers
// - commit-msg: Validate commit message references spec
package hooks

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/grokify/oscompat/fs"
)

// HookType represents a Git hook type.
type HookType string

const (
	HookPreCommit  HookType = "pre-commit"
	HookPrePush    HookType = "pre-push"
	HookCommitMsg  HookType = "commit-msg"
	HookPostCommit HookType = "post-commit"
)

// Hook represents a Git hook configuration.
type Hook struct {
	Type        HookType `json:"type" yaml:"type"`
	Enabled     bool     `json:"enabled" yaml:"enabled"`
	Description string   `json:"description" yaml:"description"`
	Script      string   `json:"script" yaml:"script"`
}

// Config holds hooks configuration.
type Config struct {
	Enabled  bool              `json:"enabled" yaml:"enabled"`
	HooksDir string            `json:"hooks_dir,omitempty" yaml:"hooks_dir,omitempty"`
	Hooks    map[HookType]Hook `json:"hooks,omitempty" yaml:"hooks,omitempty"`
}

// DefaultConfig returns default hooks configuration.
func DefaultConfig() Config {
	return Config{
		Enabled: true,
		Hooks: map[HookType]Hook{
			HookPreCommit: {
				Type:        HookPreCommit,
				Enabled:     true,
				Description: "Lint spec files before commit",
			},
			HookPrePush: {
				Type:        HookPrePush,
				Enabled:     true,
				Description: "Validate specs before push",
			},
		},
	}
}

// Manager handles Git hooks operations.
type Manager struct {
	repoRoot string
	config   Config
}

// NewManager creates a new hooks manager.
func NewManager(repoRoot string, config Config) *Manager {
	return &Manager{
		repoRoot: repoRoot,
		config:   config,
	}
}

// GitHooksDir returns the path to the .git/hooks directory.
func (m *Manager) GitHooksDir() (string, error) {
	// Check if we're in a git repo
	gitDir := filepath.Join(m.repoRoot, ".git")
	info, err := os.Stat(gitDir)
	if err != nil {
		return "", fmt.Errorf("not a git repository (no .git found)")
	}

	// Handle worktrees where .git is a file
	if !info.IsDir() {
		// Read the gitdir path from the file
		content, err := os.ReadFile(gitDir)
		if err != nil {
			return "", fmt.Errorf("reading .git file: %w", err)
		}
		gitDir = strings.TrimSpace(strings.TrimPrefix(string(content), "gitdir: "))
	}

	hooksDir := filepath.Join(gitDir, "hooks")
	return hooksDir, nil
}

// Install installs visionspec Git hooks.
func (m *Manager) Install(hookTypes []HookType) (*InstallResult, error) {
	hooksDir, err := m.GitHooksDir()
	if err != nil {
		return nil, err
	}

	// Create hooks directory if it doesn't exist
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		return nil, fmt.Errorf("creating hooks directory: %w", err)
	}

	result := &InstallResult{
		HooksDir: hooksDir,
	}

	for _, hookType := range hookTypes {
		hookPath := filepath.Join(hooksDir, string(hookType))
		template := GetTemplate(hookType)

		// Check if hook already exists
		if _, err := os.Stat(hookPath); err == nil {
			// Backup existing hook
			backupPath := hookPath + ".backup"
			if err := os.Rename(hookPath, backupPath); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("failed to backup %s: %v", hookType, err))
				continue
			}
			result.BackedUp = append(result.BackedUp, string(hookType))
		}

		// Write hook (must be executable)
		if err := os.WriteFile(hookPath, []byte(template), 0755); err != nil { //nolint:gosec // G306: Git hooks must be executable
			result.Errors = append(result.Errors, fmt.Sprintf("failed to install %s: %v", hookType, err))
			continue
		}

		result.Installed = append(result.Installed, string(hookType))
	}

	return result, nil
}

// Uninstall removes visionspec Git hooks.
func (m *Manager) Uninstall(hookTypes []HookType) (*UninstallResult, error) {
	hooksDir, err := m.GitHooksDir()
	if err != nil {
		return nil, err
	}

	result := &UninstallResult{
		HooksDir: hooksDir,
	}

	for _, hookType := range hookTypes {
		hookPath := filepath.Join(hooksDir, string(hookType))

		// Check if hook exists and is a visionspec hook
		content, err := os.ReadFile(hookPath)
		if err != nil {
			if os.IsNotExist(err) {
				result.Skipped = append(result.Skipped, string(hookType))
				continue
			}
			result.Errors = append(result.Errors, fmt.Sprintf("failed to read %s: %v", hookType, err))
			continue
		}

		if !strings.Contains(string(content), "visionspec") {
			result.Skipped = append(result.Skipped, string(hookType))
			continue
		}

		// Remove hook
		if err := os.Remove(hookPath); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("failed to remove %s: %v", hookType, err))
			continue
		}

		// Restore backup if it exists
		backupPath := hookPath + ".backup"
		if _, err := os.Stat(backupPath); err == nil {
			if err := os.Rename(backupPath, hookPath); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("failed to restore backup for %s: %v", hookType, err))
			} else {
				result.Restored = append(result.Restored, string(hookType))
			}
		}

		result.Removed = append(result.Removed, string(hookType))
	}

	return result, nil
}

// Status checks the status of installed hooks.
func (m *Manager) Status() (*StatusResult, error) {
	hooksDir, err := m.GitHooksDir()
	if err != nil {
		return nil, err
	}

	result := &StatusResult{
		HooksDir: hooksDir,
		Hooks:    make(map[HookType]HookStatus),
	}

	allHooks := []HookType{HookPreCommit, HookPrePush, HookCommitMsg, HookPostCommit}

	for _, hookType := range allHooks {
		hookPath := filepath.Join(hooksDir, string(hookType))
		status := HookStatus{
			Type: hookType,
		}

		content, err := os.ReadFile(hookPath)
		if err != nil {
			if os.IsNotExist(err) {
				status.Installed = false
			} else {
				status.Error = err.Error()
			}
		} else {
			status.Installed = true
			status.IsVisionSpec = strings.Contains(string(content), "visionspec")

			// Check if executable
			isExec, err := fs.IsExecutable(hookPath)
			if err == nil {
				status.Executable = isExec
			}
		}

		result.Hooks[hookType] = status
	}

	return result, nil
}

// InstallResult contains the result of installing hooks.
type InstallResult struct {
	HooksDir  string   `json:"hooks_dir"`
	Installed []string `json:"installed"`
	BackedUp  []string `json:"backed_up,omitempty"`
	Errors    []string `json:"errors,omitempty"`
}

// UninstallResult contains the result of uninstalling hooks.
type UninstallResult struct {
	HooksDir string   `json:"hooks_dir"`
	Removed  []string `json:"removed"`
	Restored []string `json:"restored,omitempty"`
	Skipped  []string `json:"skipped,omitempty"`
	Errors   []string `json:"errors,omitempty"`
}

// StatusResult contains the status of hooks.
type StatusResult struct {
	HooksDir string                  `json:"hooks_dir"`
	Hooks    map[HookType]HookStatus `json:"hooks"`
}

// HookStatus contains the status of a single hook.
type HookStatus struct {
	Type         HookType `json:"type"`
	Installed    bool     `json:"installed"`
	IsVisionSpec bool     `json:"is_visionspec"`
	Executable   bool     `json:"executable"`
	Error        string   `json:"error,omitempty"`
}

// FindRepoRoot finds the root of the git repository.
func FindRepoRoot(startPath string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = startPath
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not in a git repository")
	}
	return strings.TrimSpace(string(output)), nil
}

// AllHookTypes returns all supported hook types.
func AllHookTypes() []HookType {
	return []HookType{HookPreCommit, HookPrePush}
}
