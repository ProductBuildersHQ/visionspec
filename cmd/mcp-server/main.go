// Package main is the entry point for the multispec MCP server.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/plexusone/multispec/internal/mcp"
)

func main() {
	server := mcp.NewServer()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := server.Serve(ctx); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
