package git

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	ctx "github.com/plexusone/multispec/pkg/context"
)

// analyzeDependencies extracts dependencies from package manager files.
func (s *Source) analyzeDependencies(code *ctx.CodeContext) error {
	// Go modules
	if deps, err := s.parseGoMod(); err == nil {
		code.Dependencies = append(code.Dependencies, deps...)
	}

	// Node.js
	if deps, err := s.parsePackageJSON(); err == nil {
		code.Dependencies = append(code.Dependencies, deps...)
	}

	// Python
	if deps, err := s.parseRequirementsTxt(); err == nil {
		code.Dependencies = append(code.Dependencies, deps...)
	}

	// Rust
	if deps, err := s.parseCargoToml(); err == nil {
		code.Dependencies = append(code.Dependencies, deps...)
	}

	return nil
}

// parseGoMod extracts dependencies from go.mod.
func (s *Source) parseGoMod() ([]ctx.Dependency, error) {
	path := filepath.Join(s.config.Path, "go.mod")
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var deps []ctx.Dependency
	scanner := bufio.NewScanner(file)
	inRequire := false
	requireRE := regexp.MustCompile(`^\s*(\S+)\s+(\S+)`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "require (") {
			inRequire = true
			continue
		}
		if line == ")" {
			inRequire = false
			continue
		}

		// Single-line require
		if strings.HasPrefix(line, "require ") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				deps = append(deps, ctx.Dependency{
					Name:    parts[1],
					Version: parts[2],
					Type:    "direct",
					Source:  "go.mod",
				})
			}
			continue
		}

		// Multi-line require block
		if inRequire {
			matches := requireRE.FindStringSubmatch(line)
			if len(matches) >= 3 {
				depType := "direct"
				if strings.Contains(line, "// indirect") {
					depType = "indirect"
				}
				deps = append(deps, ctx.Dependency{
					Name:    matches[1],
					Version: matches[2],
					Type:    depType,
					Source:  "go.mod",
				})
			}
		}
	}

	return deps, scanner.Err()
}

// parsePackageJSON extracts dependencies from package.json.
func (s *Source) parsePackageJSON() ([]ctx.Dependency, error) {
	path := filepath.Join(s.config.Path, "package.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}

	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}

	var deps []ctx.Dependency

	for name, version := range pkg.Dependencies {
		deps = append(deps, ctx.Dependency{
			Name:    name,
			Version: version,
			Type:    "direct",
			Source:  "package.json",
		})
	}

	for name, version := range pkg.DevDependencies {
		deps = append(deps, ctx.Dependency{
			Name:    name,
			Version: version,
			Type:    "dev",
			Source:  "package.json",
		})
	}

	return deps, nil
}

// parseRequirementsTxt extracts dependencies from requirements.txt.
func (s *Source) parseRequirementsTxt() ([]ctx.Dependency, error) {
	path := filepath.Join(s.config.Path, "requirements.txt")
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var deps []ctx.Dependency
	scanner := bufio.NewScanner(file)
	// Matches: package==1.0.0, package>=1.0.0, package~=1.0.0, etc.
	depRE := regexp.MustCompile(`^([a-zA-Z0-9_-]+)([>=<~!]+)?(.*)$`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Skip -r, -e, etc.
		if strings.HasPrefix(line, "-") {
			continue
		}

		matches := depRE.FindStringSubmatch(line)
		if len(matches) >= 2 {
			version := ""
			if len(matches) >= 4 {
				version = matches[2] + matches[3]
			}
			deps = append(deps, ctx.Dependency{
				Name:    matches[1],
				Version: version,
				Type:    "direct",
				Source:  "requirements.txt",
			})
		}
	}

	return deps, scanner.Err()
}

// parseCargoToml extracts dependencies from Cargo.toml.
func (s *Source) parseCargoToml() ([]ctx.Dependency, error) {
	path := filepath.Join(s.config.Path, "Cargo.toml")
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var deps []ctx.Dependency
	scanner := bufio.NewScanner(file)
	inDeps := false
	inDevDeps := false
	// Simple TOML parsing - matches: name = "version" or name = { version = "x" }
	simpleRE := regexp.MustCompile(`^([a-zA-Z0-9_-]+)\s*=\s*"([^"]+)"`)
	tableRE := regexp.MustCompile(`^([a-zA-Z0-9_-]+)\s*=\s*\{.*version\s*=\s*"([^"]+)"`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "[dependencies]" {
			inDeps = true
			inDevDeps = false
			continue
		}
		if line == "[dev-dependencies]" {
			inDeps = false
			inDevDeps = true
			continue
		}
		if strings.HasPrefix(line, "[") {
			inDeps = false
			inDevDeps = false
			continue
		}

		if inDeps || inDevDeps {
			depType := "direct"
			if inDevDeps {
				depType = "dev"
			}

			if matches := simpleRE.FindStringSubmatch(line); len(matches) >= 3 {
				deps = append(deps, ctx.Dependency{
					Name:    matches[1],
					Version: matches[2],
					Type:    depType,
					Source:  "Cargo.toml",
				})
			} else if matches := tableRE.FindStringSubmatch(line); len(matches) >= 3 {
				deps = append(deps, ctx.Dependency{
					Name:    matches[1],
					Version: matches[2],
					Type:    depType,
					Source:  "Cargo.toml",
				})
			}
		}
	}

	return deps, scanner.Err()
}
