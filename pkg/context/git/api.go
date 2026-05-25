package git

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	ctx "github.com/ProductBuildersHQ/visionspec/pkg/context"
	"gopkg.in/yaml.v3"
)

// analyzeAPIs detects and parses API schema files.
// This is a best-effort operation; failures are silently ignored.
func (s *Source) analyzeAPIs(code *ctx.CodeContext) {
	// Find API schema files
	patterns := []struct {
		glob   string
		format string
	}{
		{"**/openapi.yaml", "openapi"},
		{"**/openapi.yml", "openapi"},
		{"**/openapi.json", "openapi"},
		{"**/swagger.yaml", "openapi"},
		{"**/swagger.yml", "openapi"},
		{"**/swagger.json", "openapi"},
		{"**/*.graphql", "graphql"},
		{"**/*.gql", "graphql"},
		{"**/*.proto", "proto"},
		{"**/api/**/*.yaml", "openapi"},
		{"**/api/**/*.json", "openapi"},
	}

	seen := make(map[string]bool)

	for _, p := range patterns {
		matches, err := filepath.Glob(filepath.Join(s.config.Path, strings.Replace(p.glob, "**", "*", -1)))
		if err != nil {
			continue
		}

		// Try deeper matching
		_ = filepath.Walk(s.config.Path, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			name := info.Name()
			if matched, _ := filepath.Match(filepath.Base(p.glob), name); matched {
				matches = append(matches, path)
			}
			return nil
		})

		for _, match := range matches {
			if seen[match] {
				continue
			}
			seen[match] = true

			relPath, _ := filepath.Rel(s.config.Path, match)

			switch p.format {
			case "openapi":
				if api := s.parseOpenAPI(match, relPath); api != nil {
					code.APIs = append(code.APIs, *api)
				}
			case "graphql":
				if api := s.parseGraphQL(match, relPath); api != nil {
					code.APIs = append(code.APIs, *api)
				}
			case "proto":
				if api := s.parseProto(match, relPath); api != nil {
					code.APIs = append(code.APIs, *api)
				}
			}
		}
	}
}

// parseOpenAPI extracts info from OpenAPI/Swagger specs.
func (s *Source) parseOpenAPI(path, relPath string) *ctx.APISchema {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	// Try to detect if it's actually an OpenAPI spec
	var spec struct {
		OpenAPI string `json:"openapi" yaml:"openapi"`
		Swagger string `json:"swagger" yaml:"swagger"`
		Info    struct {
			Title   string `json:"title" yaml:"title"`
			Version string `json:"version" yaml:"version"`
		} `json:"info" yaml:"info"`
		Paths map[string]map[string]struct {
			Summary     string `json:"summary" yaml:"summary"`
			OperationID string `json:"operationId" yaml:"operationId"`
		} `json:"paths" yaml:"paths"`
	}

	// Try JSON first, then YAML
	if err := json.Unmarshal(data, &spec); err != nil {
		if err := yaml.Unmarshal(data, &spec); err != nil {
			return nil
		}
	}

	// Verify it's actually an OpenAPI spec
	if spec.OpenAPI == "" && spec.Swagger == "" {
		return nil
	}

	api := &ctx.APISchema{
		Path:   relPath,
		Format: "openapi",
		Title:  spec.Info.Title,
	}

	// Extract routes
	for path, methods := range spec.Paths {
		for method, op := range methods {
			method = strings.ToUpper(method)
			if method == "GET" || method == "POST" || method == "PUT" ||
				method == "DELETE" || method == "PATCH" {
				api.Routes = append(api.Routes, ctx.APIRoute{
					Method:      method,
					Path:        path,
					Summary:     op.Summary,
					OperationID: op.OperationID,
				})
			}
		}
	}

	return api
}

// parseGraphQL extracts info from GraphQL schema files.
func (s *Source) parseGraphQL(path, relPath string) *ctx.APISchema {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	content := string(data)

	api := &ctx.APISchema{
		Path:   relPath,
		Format: "graphql",
	}

	// Extract queries and mutations
	queryRE := regexp.MustCompile(`(?m)^\s*(type\s+Query\s*\{[\s\S]*?\})`)
	mutationRE := regexp.MustCompile(`(?m)^\s*(type\s+Mutation\s*\{[\s\S]*?\})`)
	fieldRE := regexp.MustCompile(`^\s*(\w+)\s*(?:\([^)]*\))?\s*:\s*(\S+)`)

	// Find Query type
	if matches := queryRE.FindStringSubmatch(content); len(matches) > 1 {
		queryBlock := matches[1]
		lines := strings.Split(queryBlock, "\n")
		for _, line := range lines {
			if fieldMatches := fieldRE.FindStringSubmatch(line); len(fieldMatches) > 2 {
				api.Routes = append(api.Routes, ctx.APIRoute{
					Method:  "QUERY",
					Path:    fieldMatches[1],
					Summary: fieldMatches[2],
				})
			}
		}
	}

	// Find Mutation type
	if matches := mutationRE.FindStringSubmatch(content); len(matches) > 1 {
		mutBlock := matches[1]
		lines := strings.Split(mutBlock, "\n")
		for _, line := range lines {
			if fieldMatches := fieldRE.FindStringSubmatch(line); len(fieldMatches) > 2 {
				api.Routes = append(api.Routes, ctx.APIRoute{
					Method:  "MUTATION",
					Path:    fieldMatches[1],
					Summary: fieldMatches[2],
				})
			}
		}
	}

	// Only return if we found something
	if len(api.Routes) == 0 {
		return nil
	}

	return api
}

// parseProto extracts info from Protocol Buffer files.
func (s *Source) parseProto(path, relPath string) *ctx.APISchema {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	content := string(data)

	api := &ctx.APISchema{
		Path:   relPath,
		Format: "proto",
	}

	// Extract service definitions
	serviceRE := regexp.MustCompile(`service\s+(\w+)\s*\{([^}]*)\}`)
	rpcRE := regexp.MustCompile(`rpc\s+(\w+)\s*\(\s*(\w+)\s*\)\s*returns\s*\(\s*(\w+)\s*\)`)

	for _, serviceMatch := range serviceRE.FindAllStringSubmatch(content, -1) {
		if len(serviceMatch) < 3 {
			continue
		}
		serviceName := serviceMatch[1]
		serviceBody := serviceMatch[2]

		for _, rpcMatch := range rpcRE.FindAllStringSubmatch(serviceBody, -1) {
			if len(rpcMatch) < 4 {
				continue
			}
			api.Routes = append(api.Routes, ctx.APIRoute{
				Method:  "RPC",
				Path:    serviceName + "/" + rpcMatch[1],
				Summary: rpcMatch[2] + " -> " + rpcMatch[3],
			})
		}
	}

	// Extract package name for title
	packageRE := regexp.MustCompile(`package\s+([^;]+);`)
	if matches := packageRE.FindStringSubmatch(content); len(matches) > 1 {
		api.Title = matches[1]
	}

	// Only return if we found services
	if len(api.Routes) == 0 {
		return nil
	}

	return api
}

// extractREADME reads the README file.
// This is a best-effort operation; failures are silently ignored.
func (s *Source) extractREADME(code *ctx.CodeContext) {
	readmeNames := []string{"README.md", "README", "README.txt", "readme.md"}

	for _, name := range readmeNames {
		path := filepath.Join(s.config.Path, name)
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		content := string(data)

		// Truncate if too large
		maxSize := 5000
		if len(content) > maxSize {
			content = content[:maxSize] + "\n...(truncated)"
		}

		code.README = content
		return
	}
}
