.PHONY: build test lint clean install docs docs-serve

# Build variables
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -ldflags "-X github.com/plexusone/multispec/internal/cli.Version=$(VERSION) -X github.com/plexusone/multispec/internal/cli.Commit=$(COMMIT)"

# Build targets
build: build-cli build-mcp

build-cli:
	go build $(LDFLAGS) -o bin/multispec ./cmd/multispec

build-mcp:
	go build $(LDFLAGS) -o bin/multispec-mcp ./cmd/mcp-server

# Install to GOPATH/bin
install:
	go install $(LDFLAGS) ./cmd/multispec
	go install $(LDFLAGS) ./cmd/mcp-server

# Testing
test:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Linting
lint:
	golangci-lint run

# Clean
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Development
dev: build
	./bin/multispec --help

# Dependencies
deps:
	go mod tidy
	go mod verify

# Generate (for future code generation)
generate:
	go generate ./...

# Documentation
docs:
	mkdocs build --strict

docs-serve:
	mkdocs serve

docs-deploy:
	mkdocs gh-deploy
