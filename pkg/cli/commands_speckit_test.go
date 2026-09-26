package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func runSpeckitExport(t *testing.T, args ...string) error {
	t.Helper()
	root := &cobra.Command{Use: "test"}
	AddCommandsTo(root, DefaultConfig())

	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs(append([]string{"speckit", "export"}, args...))
	return root.Execute()
}

func TestSpeckitExportAndCheck(t *testing.T) {
	dir := t.TempDir()

	if err := runSpeckitExport(t, "aws-one-way-door", "-o", dir); err != nil {
		t.Fatalf("export failed: %v", err)
	}

	manifest := filepath.Join(dir, "extensions", "one-way-door", "extension.yml")
	if _, err := os.Stat(manifest); err != nil {
		t.Fatalf("expected manifest not written: %v", err)
	}

	if err := runSpeckitExport(t, "aws-one-way-door", "-o", dir, "--check"); err != nil {
		t.Fatalf("--check reported drift immediately after export: %v", err)
	}

	if err := os.WriteFile(manifest, []byte("tampered\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := runSpeckitExport(t, "aws-one-way-door", "-o", dir, "--check"); err == nil {
		t.Fatal("--check did not report drift after tampering")
	}
}

func TestSpeckitExportUnknownWorkflow(t *testing.T) {
	if err := runSpeckitExport(t, "no-such-workflow", "-o", t.TempDir()); err == nil {
		t.Fatal("expected error for unknown workflow")
	}
}
