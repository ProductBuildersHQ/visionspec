// Package git provides git repository context source.
package git

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bmatcuk/doublestar/v4"
	ctx "github.com/plexusone/multispec/pkg/context"
)

// Source analyzes a git repository for context.
type Source struct {
	config ctx.RepositoryConfig
	name   string
}

// NewSource creates a new git source from configuration.
func NewSource(cfg ctx.RepositoryConfig) (*Source, error) {
	// Resolve path
	path := cfg.Path

	// Clone from URL if no local path provided
	if path == "" && cfg.URL != "" {
		clonedPath, err := cloneRepo(cloneOptions{
			URL:     cfg.URL,
			Branch:  cfg.Branch,
			Sparse:  cfg.Sparse,
			Shallow: true,
			Depth:   1,
		})
		if err != nil {
			return nil, fmt.Errorf("cloning repository: %w", err)
		}
		path = clonedPath
		cfg.Path = path
	}

	if path == "" {
		return nil, fmt.Errorf("repository path or URL is required")
	}

	// Validate path exists
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("repository path: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("repository path is not a directory: %s", path)
	}

	// Set defaults
	if len(cfg.Analyze) == 0 {
		cfg.Analyze = []string{"structure", "dependencies", "api_schemas", "readme"}
	}
	if cfg.MaxDepth == 0 {
		cfg.MaxDepth = 5
	}

	// Extract name from URL or path
	var name string
	if cfg.URL != "" {
		name = extractRepoName(cfg.URL)
	} else {
		name = filepath.Base(path)
	}

	return &Source{
		config: cfg,
		name:   name,
	}, nil
}

// Name returns the source identifier.
func (s *Source) Name() string {
	return fmt.Sprintf("git:%s", s.name)
}

// Type returns the source type.
func (s *Source) Type() ctx.SourceType {
	return ctx.SourceTypeGit
}

// String returns a human-readable description.
func (s *Source) String() string {
	return fmt.Sprintf("Git repository: %s", s.config.Path)
}

// Fetch retrieves context from the git repository.
func (s *Source) Fetch(c context.Context) (*ctx.ContextData, error) {
	start := time.Now()

	code := &ctx.CodeContext{
		RepoPath:  s.config.Path,
		Languages: make(map[string]int),
	}

	// Get git info
	s.extractGitInfo(code)

	// Run configured analyzers
	var errors []string
	for _, analysis := range s.config.Analyze {
		select {
		case <-c.Done():
			return nil, c.Err()
		default:
		}

		switch analysis {
		case "structure":
			if err := s.analyzeStructure(code); err != nil {
				errors = append(errors, fmt.Sprintf("structure: %v", err))
			}
		case "dependencies":
			if err := s.analyzeDependencies(code); err != nil {
				errors = append(errors, fmt.Sprintf("dependencies: %v", err))
			}
		case "api_schemas":
			if err := s.analyzeAPIs(code); err != nil {
				errors = append(errors, fmt.Sprintf("api_schemas: %v", err))
			}
		case "readme":
			if err := s.extractREADME(code); err != nil {
				errors = append(errors, fmt.Sprintf("readme: %v", err))
			}
		}
	}

	return &ctx.ContextData{
		Source:    s.Name(),
		Type:      ctx.SourceTypeGit,
		FetchedAt: time.Now(),
		Duration:  time.Since(start),
		Code:      code,
		Summary:   ctx.GenerateCodeSummary(code),
		Errors:    errors,
	}, nil
}

// extractGitInfo extracts branch and commit from git.
func (s *Source) extractGitInfo(code *ctx.CodeContext) {
	// Read HEAD
	headPath := filepath.Join(s.config.Path, ".git", "HEAD")
	headContent, err := os.ReadFile(headPath)
	if err != nil {
		return
	}

	head := strings.TrimSpace(string(headContent))
	if strings.HasPrefix(head, "ref: refs/heads/") {
		code.Branch = strings.TrimPrefix(head, "ref: refs/heads/")

		// Read commit hash
		refPath := filepath.Join(s.config.Path, ".git", "refs", "heads", code.Branch)
		if commitContent, err := os.ReadFile(refPath); err == nil {
			code.Commit = strings.TrimSpace(string(commitContent))[:7]
		}
	} else if len(head) == 40 {
		// Detached HEAD
		code.Commit = head[:7]
	}
}

// analyzeStructure builds the directory tree.
func (s *Source) analyzeStructure(code *ctx.CodeContext) error {
	root := &ctx.TreeNode{
		Name: filepath.Base(s.config.Path),
		Type: "dir",
	}

	err := s.walkDir(s.config.Path, root, 0, code)
	if err != nil {
		return err
	}

	code.Structure = root
	return nil
}

func (s *Source) walkDir(path string, node *ctx.TreeNode, depth int, code *ctx.CodeContext) error {
	if depth >= s.config.MaxDepth {
		return nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		name := entry.Name()

		// Skip hidden files and common excludes
		if strings.HasPrefix(name, ".") || s.isExcluded(name) {
			continue
		}

		child := &ctx.TreeNode{
			Name: name,
		}

		if entry.IsDir() {
			child.Type = "dir"
			if err := s.walkDir(filepath.Join(path, name), child, depth+1, code); err != nil {
				continue // Skip directories we can't read
			}
			// Only add non-empty directories
			if len(child.Children) > 0 || depth < 2 {
				node.Children = append(node.Children, child)
			}
		} else {
			child.Type = "file"
			if info, err := entry.Info(); err == nil {
				child.Size = info.Size()
			}

			// Count language LOC
			s.countLanguage(name, child.Size, code)

			// Include file if it matches include patterns or is important
			if s.shouldIncludeFile(name) {
				node.Children = append(node.Children, child)
			}
		}
	}

	return nil
}

func (s *Source) isExcluded(name string) bool {
	// Default exclusions
	excludes := []string{
		"node_modules", "vendor", "dist", "build", "__pycache__",
		".git", ".svn", ".hg", "target", "bin", "obj",
	}

	for _, ex := range excludes {
		if name == ex {
			return true
		}
	}

	for _, pattern := range s.config.Exclude {
		if matched, _ := doublestar.Match(pattern, name); matched {
			return true
		}
	}

	return false
}

func (s *Source) shouldIncludeFile(name string) bool {
	// Always include certain files
	important := []string{
		"go.mod", "go.sum", "package.json", "package-lock.json",
		"Cargo.toml", "requirements.txt", "pyproject.toml",
		"Dockerfile", "docker-compose.yml", "Makefile",
		"README.md", "README", "LICENSE",
	}

	for _, imp := range important {
		if name == imp {
			return true
		}
	}

	// Check include patterns
	if len(s.config.Include) > 0 {
		for _, pattern := range s.config.Include {
			if matched, _ := doublestar.Match(pattern, name); matched {
				return true
			}
		}
		return false
	}

	// Default: include source files
	exts := []string{".go", ".py", ".js", ".ts", ".java", ".rs", ".rb", ".c", ".cpp", ".h"}
	for _, ext := range exts {
		if strings.HasSuffix(name, ext) {
			return true
		}
	}

	return false
}

func (s *Source) countLanguage(name string, size int64, code *ctx.CodeContext) {
	// Rough LOC estimate: ~40 chars per line average
	loc := int(size / 40)
	if loc == 0 {
		loc = 1
	}

	ext := filepath.Ext(name)
	lang := extToLanguage(ext)
	if lang != "" {
		code.Languages[lang] += loc
	}
}

func extToLanguage(ext string) string {
	switch ext {
	case ".go":
		return "Go"
	case ".py":
		return "Python"
	case ".js":
		return "JavaScript"
	case ".ts", ".tsx":
		return "TypeScript"
	case ".java":
		return "Java"
	case ".rs":
		return "Rust"
	case ".rb":
		return "Ruby"
	case ".c", ".h":
		return "C"
	case ".cpp", ".cc", ".cxx", ".hpp":
		return "C++"
	case ".cs":
		return "C#"
	case ".swift":
		return "Swift"
	case ".kt":
		return "Kotlin"
	case ".scala":
		return "Scala"
	case ".php":
		return "PHP"
	case ".sh", ".bash":
		return "Shell"
	case ".yaml", ".yml":
		return "YAML"
	case ".json":
		return "JSON"
	case ".md":
		return "Markdown"
	default:
		return ""
	}
}
