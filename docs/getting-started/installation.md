# Installation

## From Source

### Prerequisites

- Go 1.26.1 or later

### Install via `go install`

```bash
go install github.com/ProductBuildersHQ/visionspec/cmd/visionspec@latest
```

This installs the `visionspec` CLI to your `$GOPATH/bin` directory.

### Build from Source

```bash
git clone https://github.com/ProductBuildersHQ/visionspec.git
cd visionspec
make build
```

The binaries will be placed in the `bin/` directory:

- `bin/visionspec` - Main CLI
- `bin/visionspec-mcp` - MCP server

### Install Locally

```bash
make install
```

This copies the binaries to `$GOPATH/bin`.

## Verify Installation

```bash
visionspec --version
visionspec --help
```

## MCP Server

The MCP server is a separate binary for integration with AI assistants:

```bash
# Install both binaries
go install github.com/ProductBuildersHQ/visionspec/cmd/visionspec@latest
go install github.com/ProductBuildersHQ/visionspec/cmd/visionspec-mcp@latest
```

See [MCP Server](../mcp/index.md) for configuration details.
