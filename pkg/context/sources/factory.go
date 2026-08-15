// Package sources provides factory functions to create context sources.
package sources

import (
	"fmt"
	"os"
	"path/filepath"

	ctx "github.com/ProductBuildersHQ/visionspec/pkg/context"
	"github.com/ProductBuildersHQ/visionspec/pkg/context/file"
	"github.com/ProductBuildersHQ/visionspec/pkg/context/git"
	"github.com/ProductBuildersHQ/visionspec/pkg/context/graphize"
	"github.com/ProductBuildersHQ/visionspec/pkg/context/mcp"
)

// CreateFromConfig creates Source instances from configuration.
func CreateFromConfig(cfg *ctx.Config) ([]ctx.Source, error) {
	var sources []ctx.Source
	var errors []string

	// Create git repository sources
	for _, repoCfg := range cfg.Repositories {
		src, err := git.NewSource(repoCfg)
		if err != nil {
			errors = append(errors, fmt.Sprintf("git source %s: %v", repoCfg.Path, err))
			continue
		}
		sources = append(sources, src)

		// Auto-detect graphize in the repo
		if repoCfg.Graphize == "auto" || repoCfg.Graphize == "" {
			graphizePath := filepath.Join(repoCfg.Path, ".graphize")
			if info, err := os.Stat(graphizePath); err == nil && info.IsDir() {
				gSrc, err := graphize.NewSource(ctx.GraphizeConfig{
					Path: repoCfg.Path,
					Name: fmt.Sprintf("%s-graphize", filepath.Base(repoCfg.Path)),
				})
				if err == nil {
					sources = append(sources, gSrc)
				}
			}
		}
	}

	// Create standalone graphize sources
	for _, graphCfg := range cfg.Graphize {
		src, err := graphize.NewSource(graphCfg)
		if err != nil {
			errors = append(errors, fmt.Sprintf("graphize source %s: %v", graphCfg.Path, err))
			continue
		}
		sources = append(sources, src)
	}

	// Create file sources
	for _, fileCfg := range cfg.Files {
		src, err := file.NewSource(fileCfg)
		if err != nil {
			errors = append(errors, fmt.Sprintf("file source %s: %v", fileCfg.Path, err))
			continue
		}
		sources = append(sources, src)
	}

	// Create MCP server sources
	for name, mcpCfg := range cfg.MCPServers {
		src, err := mcp.NewSource(name, mcpCfg)
		if err != nil {
			errors = append(errors, fmt.Sprintf("mcp source %s: %v", name, err))
			continue
		}
		sources = append(sources, src)
	}

	if len(errors) > 0 && len(sources) == 0 {
		return nil, fmt.Errorf("failed to create any sources: %v", errors)
	}

	return sources, nil
}

// BuildAggregator creates an Aggregator with sources from configuration.
func BuildAggregator(projectName string, cfg *ctx.Config) (*ctx.Aggregator, error) {
	if cfg == nil {
		cfg = ctx.DefaultConfig()
	}

	agg := ctx.NewAggregator(projectName, cfg)

	sources, err := CreateFromConfig(cfg)
	if err != nil {
		return nil, err
	}

	for _, src := range sources {
		agg.AddSource(src)
	}

	return agg, nil
}
