package git

import (
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// cloneOptions configures repository cloning.
type cloneOptions struct {
	URL     string
	Branch  string
	Sparse  []string
	Shallow bool
	Depth   int
}

// cloneRepo clones a git repository to a temporary directory.
// Returns the path to the cloned repository.
func cloneRepo(opts cloneOptions) (string, error) {
	// Create cache directory
	cacheDir, err := getCacheDir()
	if err != nil {
		return "", fmt.Errorf("getting cache dir: %w", err)
	}

	// Generate unique directory name based on URL and branch
	hash := sha256.Sum256([]byte(opts.URL + opts.Branch))
	repoDir := filepath.Join(cacheDir, fmt.Sprintf("%x", hash[:8]))

	// Check if already cloned
	if info, err := os.Stat(filepath.Join(repoDir, ".git")); err == nil && info.IsDir() {
		// Update existing clone
		if err := updateRepo(repoDir, opts.Branch); err != nil {
			// If update fails, remove and re-clone
			os.RemoveAll(repoDir)
		} else {
			return repoDir, nil
		}
	}

	// Build git clone command
	args := []string{"clone"}

	// Shallow clone
	if opts.Shallow || opts.Depth > 0 {
		depth := opts.Depth
		if depth == 0 {
			depth = 1
		}
		args = append(args, "--depth", fmt.Sprintf("%d", depth))
	}

	// Single branch
	if opts.Branch != "" {
		args = append(args, "--branch", opts.Branch, "--single-branch")
	}

	// Sparse checkout
	if len(opts.Sparse) > 0 {
		args = append(args, "--sparse", "--filter=blob:none")
	}

	args = append(args, opts.URL, repoDir)

	// Run git clone
	cmd := exec.Command("git", args...)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git clone: %w", err)
	}

	// Configure sparse checkout paths
	if len(opts.Sparse) > 0 {
		if err := configureSparseCheckout(repoDir, opts.Sparse); err != nil {
			return "", fmt.Errorf("configuring sparse checkout: %w", err)
		}
	}

	return repoDir, nil
}

// updateRepo pulls the latest changes.
func updateRepo(repoDir, branch string) error {
	cmd := exec.Command("git", "-C", repoDir, "pull", "--ff-only")
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// configureSparseCheckout sets up sparse checkout for specific paths.
func configureSparseCheckout(repoDir string, paths []string) error {
	// Initialize sparse checkout
	cmd := exec.Command("git", "-C", repoDir, "sparse-checkout", "init", "--cone")
	if err := cmd.Run(); err != nil {
		return err
	}

	// Set sparse checkout paths
	args := append([]string{"-C", repoDir, "sparse-checkout", "set"}, paths...)
	cmd = exec.Command("git", args...)
	return cmd.Run()
}

// getCacheDir returns the cache directory for cloned repos.
func getCacheDir() (string, error) {
	// Use system temp directory with multispec subdirectory
	cacheDir := filepath.Join(os.TempDir(), "multispec-repos")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", err
	}
	return cacheDir, nil
}

// extractRepoName extracts repository name from URL.
func extractRepoName(url string) string {
	// Handle various URL formats:
	// https://github.com/owner/repo.git
	// git@github.com:owner/repo.git
	// /path/to/repo

	url = strings.TrimSuffix(url, ".git")

	// Handle SSH format
	if strings.Contains(url, ":") && !strings.Contains(url, "://") {
		parts := strings.Split(url, ":")
		if len(parts) == 2 {
			url = parts[1]
		}
	}

	// Extract last path component
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return "repo"
}
