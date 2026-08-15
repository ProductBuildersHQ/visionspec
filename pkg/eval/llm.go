package eval

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/ProductBuildersHQ/visionspec/pkg/types"
	"github.com/plexusone/omnillm-core"
)

// LLMConfig configures the LLM client.
type LLMConfig struct {
	Provider    string  // Provider name (openai, anthropic, gemini, etc.)
	Model       string  // Model name
	APIKey      string  // API key (optional if env var is set)
	Temperature float64 // Temperature for generation (default 0.0 for deterministic)
	MaxTokens   int     // Max tokens for response (default 4096)
}

// DefaultLLMConfig returns default configuration for evaluation.
func DefaultLLMConfig() LLMConfig {
	return LLMConfig{
		Provider:    "anthropic",
		Model:       "claude-sonnet-4-20250514",
		Temperature: 0.0, // Deterministic for evaluation
		MaxTokens:   4096,
	}
}

// LLMClient wraps omnillm for evaluation requests.
type LLMClient struct {
	client *omnillm.ChatClient
	config LLMConfig
}

// NewLLMClient creates a new LLM client with the given configuration.
func NewLLMClient(cfg LLMConfig) (*LLMClient, error) {
	// Validate config
	if cfg.Provider == "" {
		cfg.Provider = "anthropic"
	}
	if cfg.Model == "" {
		cfg.Model = "claude-sonnet-4-20250514"
	}
	if cfg.MaxTokens == 0 {
		cfg.MaxTokens = 4096
	}

	// Get API key from config or environment
	apiKey := cfg.APIKey
	if apiKey == "" {
		apiKey = getAPIKeyFromEnv(cfg.Provider)
	}
	if apiKey == "" {
		return nil, fmt.Errorf("no API key found for provider %s", cfg.Provider)
	}

	// Map provider name to omnillm provider
	providerName, err := mapProviderName(cfg.Provider)
	if err != nil {
		return nil, err
	}

	// Create omnillm client
	client, err := omnillm.NewClient(omnillm.ClientConfig{
		Providers: []omnillm.ProviderConfig{
			{
				Provider: providerName,
				APIKey:   apiKey,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create LLM client: %w", err)
	}

	return &LLMClient{
		client: client,
		config: cfg,
	}, nil
}

// Complete sends a prompt to the LLM and returns the response.
func (c *LLMClient) Complete(ctx context.Context, prompt string) (string, JudgeMetadata, error) {
	temperature := c.config.Temperature
	maxTokens := c.config.MaxTokens

	req := &omnillm.ChatCompletionRequest{
		Model: c.config.Model,
		Messages: []omnillm.Message{
			{
				Role:    omnillm.RoleUser,
				Content: prompt,
			},
		},
		Temperature: &temperature,
		MaxTokens:   &maxTokens,
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", JudgeMetadata{}, fmt.Errorf("chat completion failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", JudgeMetadata{}, errors.New("no response choices returned")
	}

	metadata := JudgeMetadata{
		Model:       c.config.Model,
		Provider:    c.config.Provider,
		Temperature: c.config.Temperature,
		Tokens:      resp.Usage.TotalTokens,
	}

	return resp.Choices[0].Message.Content, metadata, nil
}

// Close releases resources.
func (c *LLMClient) Close() error {
	return c.client.Close()
}

// mapProviderName maps a string provider name to omnillm.ProviderName.
func mapProviderName(name string) (omnillm.ProviderName, error) {
	switch name {
	case "openai":
		return omnillm.ProviderNameOpenAI, nil
	case "anthropic":
		return omnillm.ProviderNameAnthropic, nil
	case "gemini", "google":
		return omnillm.ProviderNameGemini, nil
	case "xai":
		return omnillm.ProviderNameXAI, nil
	case "glm":
		return omnillm.ProviderNameGLM, nil
	case "kimi":
		return omnillm.ProviderNameKimi, nil
	case "qwen":
		return omnillm.ProviderNameQwen, nil
	case "ollama":
		return omnillm.ProviderNameOllama, nil
	default:
		return "", fmt.Errorf("unknown provider: %s", name)
	}
}

// getAPIKeyFromEnv gets the API key from environment variables.
func getAPIKeyFromEnv(provider string) string {
	switch provider {
	case "openai":
		return os.Getenv("OPENAI_API_KEY")
	case "anthropic":
		return os.Getenv("ANTHROPIC_API_KEY")
	case "gemini", "google":
		return os.Getenv("GEMINI_API_KEY")
	case "xai":
		return os.Getenv("XAI_API_KEY")
	case "glm":
		return os.Getenv("GLM_API_KEY")
	case "kimi":
		return os.Getenv("KIMI_API_KEY")
	case "qwen":
		return os.Getenv("QWEN_API_KEY")
	default:
		return ""
	}
}

// NewLLMClientFromEnv creates an LLM client using environment configuration.
// It tries providers in order: ANTHROPIC, OPENAI, GEMINI.
func NewLLMClientFromEnv() (*LLMClient, error) {
	// Try providers in preference order
	providers := []struct {
		name   string
		envVar string
		model  string
	}{
		{"anthropic", "ANTHROPIC_API_KEY", "claude-sonnet-4-20250514"},
		{"openai", "OPENAI_API_KEY", "gpt-4o"},
		{"gemini", "GEMINI_API_KEY", "gemini-2.5-pro"},
	}

	for _, p := range providers {
		if apiKey := os.Getenv(p.envVar); apiKey != "" {
			return NewLLMClient(LLMConfig{
				Provider:    p.name,
				Model:       p.model,
				APIKey:      apiKey,
				Temperature: 0.0,
				MaxTokens:   4096,
			})
		}
	}

	return nil, errors.New("no LLM API key found in environment (set ANTHROPIC_API_KEY, OPENAI_API_KEY, or GEMINI_API_KEY)")
}

// NewLLMClientFromProject creates an LLM client using project configuration.
// Project config values take precedence; missing values fall back to environment defaults.
func NewLLMClientFromProject(projectCfg *types.LLMConfig) (*LLMClient, error) {
	// Start with defaults
	cfg := DefaultLLMConfig()

	// If no project config, use environment
	if projectCfg == nil {
		return NewLLMClientFromEnv()
	}

	// Override with project config values
	if projectCfg.Provider != "" {
		cfg.Provider = projectCfg.Provider
		// Update model to match provider default if model not specified
		if projectCfg.Model == "" {
			cfg.Model = defaultModelForProvider(projectCfg.Provider)
		}
	}
	if projectCfg.Model != "" {
		cfg.Model = projectCfg.Model
	}
	if projectCfg.Temperature != nil {
		cfg.Temperature = *projectCfg.Temperature
	}
	if projectCfg.MaxTokens != nil {
		cfg.MaxTokens = *projectCfg.MaxTokens
	}

	return NewLLMClient(cfg)
}

// defaultModelForProvider returns the default model for a given provider.
func defaultModelForProvider(provider string) string {
	switch provider {
	case "anthropic":
		return "claude-sonnet-4-20250514"
	case "openai":
		return "gpt-4o"
	case "gemini", "google":
		return "gemini-2.5-pro"
	case "xai":
		return "grok-3"
	case "glm":
		return "glm-4"
	case "kimi":
		return "moonshot-v1-8k"
	case "qwen":
		return "qwen-max"
	case "ollama":
		return "llama3.3"
	default:
		return ""
	}
}
