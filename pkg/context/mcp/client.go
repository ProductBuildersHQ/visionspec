// Package mcp provides MCP client for external context sources.
package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
	"time"
)

// Client manages communication with an MCP server.
type Client struct {
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	stdout  *bufio.Reader
	mu      sync.Mutex
	nextID  atomic.Int64
	pending map[int64]chan *Response
	pendMu  sync.Mutex
	closed  bool
}

// Request represents a JSON-RPC request.
type Request struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int64  `json:"id"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
}

// Response represents a JSON-RPC response.
type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int64           `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
}

// RPCError represents a JSON-RPC error.
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (e *RPCError) Error() string {
	return fmt.Sprintf("RPC error %d: %s", e.Code, e.Message)
}

// NewClient creates a new MCP client by spawning the server process.
func NewClient(command string, args []string, env map[string]string) (*Client, error) {
	cmd := exec.Command(command, args...)

	// Set environment
	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("creating stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("creating stdout pipe: %w", err)
	}

	// Capture stderr for debugging
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("starting MCP server: %w", err)
	}

	client := &Client{
		cmd:     cmd,
		stdin:   stdin,
		stdout:  bufio.NewReader(stdout),
		pending: make(map[int64]chan *Response),
	}

	// Start reading responses
	go client.readResponses()

	return client, nil
}

// Initialize sends the initialize request to the MCP server.
func (c *Client) Initialize(ctx context.Context) error {
	params := map[string]any{
		"protocolVersion": "2024-11-05",
		"capabilities":    map[string]any{},
		"clientInfo": map[string]string{
			"name":    "multispec",
			"version": "0.4.0",
		},
	}

	var result json.RawMessage
	if err := c.Call(ctx, "initialize", params, &result); err != nil {
		return fmt.Errorf("initialize: %w", err)
	}

	// Send initialized notification
	return c.Notify("notifications/initialized", nil)
}

// ListTools returns the available tools from the MCP server.
func (c *Client) ListTools(ctx context.Context) ([]Tool, error) {
	var result struct {
		Tools []Tool `json:"tools"`
	}

	if err := c.Call(ctx, "tools/list", nil, &result); err != nil {
		return nil, err
	}

	return result.Tools, nil
}

// Tool represents an MCP tool.
type Tool struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"inputSchema"`
}

// CallTool invokes a tool on the MCP server.
func (c *Client) CallTool(ctx context.Context, name string, args map[string]any) ([]ToolContent, error) {
	params := map[string]any{
		"name":      name,
		"arguments": args,
	}

	var result struct {
		Content []ToolContent `json:"content"`
		IsError bool          `json:"isError"`
	}

	if err := c.Call(ctx, "tools/call", params, &result); err != nil {
		return nil, err
	}

	if result.IsError {
		// Extract error message from content
		for _, c := range result.Content {
			if c.Type == "text" {
				return nil, fmt.Errorf("tool error: %s", c.Text)
			}
		}
		return nil, fmt.Errorf("tool returned error")
	}

	return result.Content, nil
}

// ToolContent represents content returned by a tool.
type ToolContent struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// Call makes a synchronous RPC call.
func (c *Client) Call(ctx context.Context, method string, params any, result any) error {
	id := c.nextID.Add(1)

	req := &Request{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  params,
	}

	// Create response channel
	respCh := make(chan *Response, 1)
	c.pendMu.Lock()
	c.pending[id] = respCh
	c.pendMu.Unlock()

	defer func() {
		c.pendMu.Lock()
		delete(c.pending, id)
		c.pendMu.Unlock()
	}()

	// Send request
	if err := c.send(req); err != nil {
		return err
	}

	// Wait for response
	select {
	case <-ctx.Done():
		return ctx.Err()
	case resp := <-respCh:
		if resp.Error != nil {
			return resp.Error
		}
		if result != nil && len(resp.Result) > 0 {
			return json.Unmarshal(resp.Result, result)
		}
		return nil
	}
}

// Notify sends a notification (no response expected).
func (c *Client) Notify(method string, params any) error {
	req := &Request{
		JSONRPC: "2.0",
		ID:      0, // Notifications have no ID
		Method:  method,
		Params:  params,
	}

	return c.send(req)
}

func (c *Client) send(req *Request) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return fmt.Errorf("client closed")
	}

	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	// Write with newline delimiter
	_, err = fmt.Fprintf(c.stdin, "%s\n", data)
	return err
}

func (c *Client) readResponses() {
	for {
		line, err := c.stdout.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Fprintf(os.Stderr, "MCP read error: %v\n", err)
			}
			return
		}

		var resp Response
		if err := json.Unmarshal(line, &resp); err != nil {
			continue // Skip malformed responses
		}

		// Route response to waiting caller
		c.pendMu.Lock()
		if ch, ok := c.pending[resp.ID]; ok {
			ch <- &resp
		}
		c.pendMu.Unlock()
	}
}

// Close shuts down the MCP client.
func (c *Client) Close() error {
	c.mu.Lock()
	c.closed = true
	c.mu.Unlock()

	c.stdin.Close()

	// Give server time to exit gracefully
	done := make(chan error, 1)
	go func() {
		done <- c.cmd.Wait()
	}()

	select {
	case <-done:
		return nil
	case <-time.After(5 * time.Second):
		_ = c.cmd.Process.Kill() // Best effort; we're already returning an error
		return fmt.Errorf("server did not exit gracefully")
	}
}
