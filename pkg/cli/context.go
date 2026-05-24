package cli

import "context"

// configKey is the context key for Config.
type configKey struct{}

// WithConfig returns a context with the config attached.
func WithConfig(ctx context.Context, cfg *Config) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, configKey{}, cfg)
}

// ConfigFromContext returns the Config from a context.
func ConfigFromContext(ctx context.Context) *Config {
	if ctx == nil {
		return nil
	}
	cfg, _ := ctx.Value(configKey{}).(*Config)
	return cfg
}
