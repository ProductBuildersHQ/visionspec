# Contributing

Guidelines for contributing to VisionSpec.

## Development Setup

### Prerequisites

- Go 1.21 or later
- Make
- golangci-lint (for linting)

### Clone and Build

```bash
git clone https://github.com/ProductBuildersHQ/visionspec.git
cd visionspec
make build
```

### Run Tests

```bash
make test
```

### Run Linter

```bash
make lint
```

## Project Structure

```
visionspec/
├── cmd/
│   ├── visionspec/      # Main CLI entry point
│   └── mcp-server/     # MCP server entry point
├── internal/
│   ├── cli/            # CLI commands (Cobra)
│   └── mcp/            # MCP server implementation
├── pkg/
│   ├── config/         # Configuration loading
│   ├── lint/           # Linting logic
│   ├── status/         # Status generation
│   └── types/          # Core types
├── docs/               # MkDocs documentation
└── Makefile
```

## Code Style

- Use `gofmt` for formatting
- Follow standard Go conventions
- Run `golangci-lint run` before committing
- Keep functions focused and small

## Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add new feature
fix: resolve bug
docs: update documentation
refactor: restructure code
test: add tests
chore: maintenance tasks
```

## Pull Requests

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and linter
5. Submit a pull request

## Adding a New Command

1. Create a new file in `internal/cli/`
2. Define the Cobra command
3. Register in `internal/cli/root.go`
4. Add tests
5. Update documentation in `docs/cli/`

## Adding an MCP Tool

1. Define argument types in `internal/mcp/server.go`
2. Add the tool registration in `registerTools()`
3. Implement the handler function
4. Add documentation in `docs/mcp/tools.md`

## Documentation

Documentation is built with MkDocs. To preview locally:

```bash
pip install mkdocs-material
mkdocs serve
```

## Roadmap

See [ROADMAP.md](../specs/ROADMAP.md) for planned features and implementation status.
