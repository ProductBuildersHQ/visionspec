// Package file provides local file context source.
package file

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	ctx "github.com/ProductBuildersHQ/visionspec/pkg/context"
)

// Source reads context from local files.
type Source struct {
	config ctx.FileConfig
}

// NewSource creates a new file source from configuration.
func NewSource(cfg ctx.FileConfig) (*Source, error) {
	if cfg.Path == "" {
		return nil, fmt.Errorf("file path is required")
	}

	// Validate file exists
	info, err := os.Stat(cfg.Path)
	if err != nil {
		return nil, fmt.Errorf("file: %w", err)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("path is a directory, not a file: %s", cfg.Path)
	}

	return &Source{config: cfg}, nil
}

// Name returns the source identifier.
func (s *Source) Name() string {
	return fmt.Sprintf("file:%s", filepath.Base(s.config.Path))
}

// Type returns the source type.
func (s *Source) Type() ctx.SourceType {
	return ctx.SourceTypeFile
}

// String returns a human-readable description.
func (s *Source) String() string {
	return fmt.Sprintf("File: %s", s.config.Path)
}

// Fetch reads the file and returns its context.
func (s *Source) Fetch(c context.Context) (*ctx.ContextData, error) {
	start := time.Now()

	data, err := os.ReadFile(s.config.Path)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	content := string(data)

	// Truncate if too large
	maxSize := s.config.MaxSize
	if maxSize == 0 {
		maxSize = 50000 // 50KB default
	}
	if int64(len(content)) > maxSize {
		content = content[:maxSize] + "\n...(truncated)"
	}

	// Detect file type and format
	fileType := s.config.Type
	if fileType == "" {
		fileType = detectFileType(s.config.Path)
	}

	format := detectFormat(s.config.Path)

	file := &ctx.FileContext{
		Path:    s.config.Path,
		Type:    fileType,
		Content: content,
		Format:  format,
	}

	return &ctx.ContextData{
		Source:    s.Name(),
		Type:      ctx.SourceTypeFile,
		FetchedAt: time.Now(),
		Duration:  time.Since(start),
		File:      file,
		Summary:   fmt.Sprintf("File: %s (%s, %d bytes)", filepath.Base(s.config.Path), format, len(data)),
	}, nil
}

// detectFileType guesses the semantic type of a file.
func detectFileType(path string) string {
	name := strings.ToLower(filepath.Base(path))

	// Architecture docs
	if strings.Contains(name, "architecture") ||
		strings.Contains(name, "design") ||
		strings.HasPrefix(name, "adr") {
		return "architecture"
	}

	// API specs
	if strings.Contains(name, "openapi") ||
		strings.Contains(name, "swagger") ||
		strings.HasSuffix(name, ".graphql") ||
		strings.HasSuffix(name, ".proto") {
		return "api_spec"
	}

	// README
	if strings.HasPrefix(name, "readme") {
		return "readme"
	}

	// Diagrams
	if strings.HasSuffix(name, ".puml") ||
		strings.HasSuffix(name, ".mmd") ||
		strings.HasSuffix(name, ".dot") {
		return "diagram"
	}

	// Requirements
	if strings.Contains(name, "requirements") ||
		strings.Contains(name, "spec") {
		return "requirements"
	}

	// Config
	if strings.Contains(name, "config") ||
		name == "settings.yaml" ||
		name == "settings.json" {
		return "config"
	}

	return "document"
}

// detectFormat returns the file format based on extension.
func detectFormat(path string) string {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".md", ".markdown":
		return "markdown"
	case ".yaml", ".yml":
		return "yaml"
	case ".json":
		return "json"
	case ".toml":
		return "toml"
	case ".xml":
		return "xml"
	case ".txt":
		return "text"
	case ".graphql", ".gql":
		return "graphql"
	case ".proto":
		return "proto"
	case ".puml":
		return "plantuml"
	case ".mmd":
		return "mermaid"
	default:
		return "text"
	}
}
