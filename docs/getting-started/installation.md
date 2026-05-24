# Installation

## From Source

### Prerequisites

- Go 1.21 or later

### Install via `go install`

```bash
go install github.com/plexusone/multispec/cmd/multispec@latest
```

This installs the `multispec` CLI to your `$GOPATH/bin` directory.

### Build from Source

```bash
git clone https://github.com/plexusone/multispec.git
cd multispec
make build
```

The binaries will be placed in the `bin/` directory:

- `bin/multispec` - Main CLI
- `bin/multispec-mcp` - MCP server

### Install Locally

```bash
make install
```

This copies the binaries to `$GOPATH/bin`.

## Verify Installation

```bash
multispec --version
multispec --help
```

## MCP Server

The MCP server is a separate binary for integration with AI assistants:

```bash
# Install both binaries
go install github.com/plexusone/multispec/cmd/multispec@latest
go install github.com/plexusone/multispec/cmd/multispec-mcp@latest
```

See [MCP Server](../mcp/index.md) for configuration details.
