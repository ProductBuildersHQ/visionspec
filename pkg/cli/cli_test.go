package cli

import (
	"testing"

	"github.com/ProductBuildersHQ/visionspec/pkg/rubrics"
	"github.com/ProductBuildersHQ/visionspec/pkg/templates"
	"github.com/ProductBuildersHQ/visionspec/pkg/types"
	"github.com/spf13/cobra"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg == nil {
		t.Fatal("DefaultConfig() returned nil")
	}

	if cfg.TemplateLoader == nil {
		t.Error("TemplateLoader is nil")
	}

	if cfg.RubricLoader == nil {
		t.Error("RubricLoader is nil")
	}

	if cfg.Version == "" {
		t.Error("Version is empty")
	}
}

func TestCommands(t *testing.T) {
	cfg := DefaultConfig()
	cmds := Commands(cfg)

	if cmds == nil {
		t.Fatal("Commands() returned nil")
	}

	// Check that all commands are present
	if cmds.Init == nil {
		t.Error("Init command is nil")
	}
	if cmds.Lint == nil {
		t.Error("Lint command is nil")
	}
	if cmds.Status == nil {
		t.Error("Status command is nil")
	}
	if cmds.Eval == nil {
		t.Error("Eval command is nil")
	}
	if cmds.Synthesize == nil {
		t.Error("Synthesize command is nil")
	}
	if cmds.Reconcile == nil {
		t.Error("Reconcile command is nil")
	}
	if cmds.Approve == nil {
		t.Error("Approve command is nil")
	}
	if cmds.Export == nil {
		t.Error("Export command is nil")
	}
	if cmds.Targets == nil {
		t.Error("Targets command is nil")
	}
	if cmds.Graph == nil {
		t.Error("Graph command is nil")
	}
	if cmds.Serve == nil {
		t.Error("Serve command is nil")
	}
}

func TestAddCommandsTo(t *testing.T) {
	root := &cobra.Command{Use: "test"}
	cfg := DefaultConfig()

	AddCommandsTo(root, cfg)

	// Check that commands were added
	if len(root.Commands()) == 0 {
		t.Error("No commands were added to root")
	}

	// Check specific commands exist
	names := make(map[string]bool)
	for _, cmd := range root.Commands() {
		names[cmd.Name()] = true
	}

	expected := []string{"init", "lint", "status", "eval", "synthesize", "reconcile", "approve", "export", "targets", "graph", "serve"}
	for _, name := range expected {
		if !names[name] {
			t.Errorf("Command %q not found", name)
		}
	}
}

func TestAddCommandsToWithNilConfig(t *testing.T) {
	root := &cobra.Command{Use: "test"}

	// Should not panic with nil config
	AddCommandsTo(root, nil)

	// Commands should still be added with default config
	if len(root.Commands()) == 0 {
		t.Error("No commands were added to root with nil config")
	}
}

func TestCommandsWithNilConfig(t *testing.T) {
	// Should not panic with nil config
	cmds := Commands(nil)

	if cmds == nil {
		t.Fatal("Commands(nil) returned nil")
	}

	// All commands should be present
	if cmds.Init == nil {
		t.Error("Init command is nil")
	}
}

func TestSelectiveCommandInclusion(t *testing.T) {
	// Demonstrate selective command inclusion
	root := &cobra.Command{Use: "org-spec"}
	cfg := DefaultConfig()

	// Get all commands
	cmds := Commands(cfg)

	// Only add subset of commands
	root.AddCommand(cmds.Init)
	root.AddCommand(cmds.Lint)
	root.AddCommand(cmds.Status)

	// Verify only selected commands are present
	if len(root.Commands()) != 3 {
		t.Errorf("Expected 3 commands, got %d", len(root.Commands()))
	}
}

func TestCustomConfig(t *testing.T) {
	// Test with custom loaders
	tmpDir := t.TempDir()

	customTemplateLoader := templates.NewChainLoader(
		templates.NewFileLoader(tmpDir),
		templates.EmbeddedLoader(),
	)

	customRubricLoader := rubrics.NewChainLoader(
		rubrics.NewFileLoader(tmpDir),
		rubrics.EmbeddedLoader(),
	)

	cfg := &Config{
		TemplateLoader: customTemplateLoader,
		RubricLoader:   customRubricLoader,
		Version:        "custom-1.0.0",
	}

	// Verify custom loaders are used
	if cfg.TemplateLoader == nil {
		t.Error("Custom TemplateLoader not set")
	}

	if cfg.RubricLoader == nil {
		t.Error("Custom RubricLoader not set")
	}

	// Test that Available() works with custom loaders
	available := cfg.TemplateLoader.Available()
	if len(available) == 0 {
		t.Error("No templates available from custom loader")
	}

	// Verify built-in templates are still accessible via chain
	_, err := cfg.TemplateLoader.Load(types.SpecTypeMRD)
	if err != nil {
		t.Errorf("Failed to load MRD template from custom loader: %v", err)
	}
}
