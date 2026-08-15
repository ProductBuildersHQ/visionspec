package layout

import (
	"testing"
)

func TestDefaultLayout(t *testing.T) {
	l := DefaultLayout()

	if l.SourceDir != "source" {
		t.Errorf("expected SourceDir=source, got %s", l.SourceDir)
	}
	if l.GTMDir != "gtm" {
		t.Errorf("expected GTMDir=gtm, got %s", l.GTMDir)
	}
	if l.EvalDir != "evals" {
		t.Errorf("expected EvalDir=evals, got %s", l.EvalDir)
	}
}

func TestDirForCategory(t *testing.T) {
	l := DefaultLayout()

	tests := []struct {
		category string
		want     string
	}{
		{"source", "source"},
		{"gtm", "gtm"},
		{"technical", "technical"},
		{"strategic", "strategic"},
		{"execution", "execution"},
		{"output", ""},
		{"unknown", ""},
		{"SOURCE", "source"}, // case insensitive
	}

	for _, tt := range tests {
		got := l.DirForCategory(tt.category)
		if got != tt.want {
			t.Errorf("DirForCategory(%q) = %q, want %q", tt.category, got, tt.want)
		}
	}
}

func TestSpecPath(t *testing.T) {
	l := DefaultLayout()

	tests := []struct {
		specType string
		category string
		want     string
	}{
		{"mrd", "source", "source/mrd.md"},
		{"press", "gtm", "gtm/press.md"},
		{"trd", "technical", "technical/trd.md"},
		{"spec", "output", "spec.md"},
	}

	for _, tt := range tests {
		got := l.SpecPath(tt.specType, tt.category)
		if got != tt.want {
			t.Errorf("SpecPath(%q, %q) = %q, want %q", tt.specType, tt.category, got, tt.want)
		}
	}
}

func TestEvalPath(t *testing.T) {
	l := DefaultLayout()

	got := l.EvalPath("mrd")
	want := "evals/mrd.eval.json"
	if got != want {
		t.Errorf("EvalPath(mrd) = %q, want %q", got, want)
	}
}

func TestCustomLayout(t *testing.T) {
	l := &Layout{
		SourceDir:     "specs/source",
		GTMDir:        "specs/gtm",
		EvalDir:       "results",
		SpecExtension: ".markdown",
		EvalExtension: ".result.json",
	}

	if got := l.SpecPath("mrd", "source"); got != "specs/source/mrd.markdown" {
		t.Errorf("custom SpecPath = %q", got)
	}

	if got := l.EvalPath("mrd"); got != "results/mrd.result.json" {
		t.Errorf("custom EvalPath = %q", got)
	}
}

func TestNilLayout(t *testing.T) {
	var l *Layout

	if got := l.DirForCategory("source"); got != "source" {
		t.Errorf("nil layout DirForCategory = %q, want source", got)
	}

	if got := l.SpecFilename("mrd"); got != "mrd.md" {
		t.Errorf("nil layout SpecFilename = %q, want mrd.md", got)
	}
}
