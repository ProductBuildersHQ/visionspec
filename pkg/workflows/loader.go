package workflows

import (
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/ProductBuildersHQ/visionspec/pkg/template"
	"github.com/ProductBuildersHQ/visionspec/pkg/workflow"
	oscompatfs "github.com/grokify/oscompat/fs"
	"github.com/plexusone/structured-evaluation/rubric"
	"gopkg.in/yaml.v3"
)

// Loader loads workflows by name.
type Loader interface {
	// Load returns a workflow by name.
	Load(name string) (*LoadedWorkflow, error)

	// Available returns all available workflow names.
	Available() []string
}

// DefaultLoader returns the default loader using embedded workflows.
func DefaultLoader() Loader {
	return NewResolvingLoader(&embeddedLoader{})
}

// embeddedLoader loads workflows from the embedded registry.
type embeddedLoader struct{}

func (l *embeddedLoader) Load(name string) (*LoadedWorkflow, error) {
	return Get(name)
}

func (l *embeddedLoader) Available() []string {
	return List()
}

// FileLoader loads workflows from a filesystem directory.
type FileLoader struct {
	dir string
}

// NewFileLoader creates a loader that reads workflows from a directory.
// Each workflow is a subdirectory containing profile.yaml, templates/, and rubrics/.
func NewFileLoader(dir string) *FileLoader {
	return &FileLoader{dir: dir}
}

func (l *FileLoader) Load(name string) (*LoadedWorkflow, error) {
	workflowDir := filepath.Join(l.dir, name)

	// Load profile.yaml
	profilePath := filepath.Join(workflowDir, "profile.yaml")
	data, err := os.ReadFile(profilePath)
	if err != nil {
		return nil, fmt.Errorf("workflow %q not found: %w", name, err)
	}

	w, err := workflow.ParseYAML(data)
	if err != nil {
		return nil, fmt.Errorf("parsing workflow %q: %w", name, err)
	}

	loaded := &LoadedWorkflow{
		Workflow:  w,
		Templates: make(map[string]*template.Template),
		Rubrics:   make(map[string]*rubric.RubricSet),
	}

	// Load templates
	templatesDir := filepath.Join(workflowDir, "templates")
	if err := l.loadTemplates(loaded, templatesDir); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("loading templates: %w", err)
	}

	// Load rubrics
	rubricsDir := filepath.Join(workflowDir, "rubrics")
	if err := l.loadRubrics(loaded, rubricsDir); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("loading rubrics: %w", err)
	}

	return loaded, nil
}

func (l *FileLoader) loadTemplates(loaded *LoadedWorkflow, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		specType := strings.TrimSuffix(entry.Name(), ".md")
		content, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return fmt.Errorf("reading template %s: %w", entry.Name(), err)
		}

		loaded.Templates[specType] = &template.Template{
			ID:       specType,
			SpecType: specType,
			Content:  oscompatfs.NormalizeLineEndings(string(content)),
		}
	}

	return nil
}

func (l *FileLoader) loadRubrics(loaded *LoadedWorkflow, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".rubric.yaml") {
			continue
		}

		specType := strings.TrimSuffix(entry.Name(), ".rubric.yaml")
		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return fmt.Errorf("reading rubric %s: %w", entry.Name(), err)
		}

		var rs rubric.RubricSet
		if err := yaml.Unmarshal(data, &rs); err != nil {
			return fmt.Errorf("parsing rubric %s: %w", entry.Name(), err)
		}

		loaded.Rubrics[specType] = &rs
	}

	return nil
}

func (l *FileLoader) Available() []string {
	entries, err := os.ReadDir(l.dir)
	if err != nil {
		return nil
	}

	var names []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		profilePath := filepath.Join(l.dir, entry.Name(), "profile.yaml")
		if _, err := os.Stat(profilePath); err == nil {
			names = append(names, entry.Name())
		}
	}
	return names
}

// FSLoader loads workflows from an fs.FS (for embed or testing).
type FSLoader struct {
	fsys fs.FS
	dir  string
}

// NewFSLoader creates a loader from an fs.FS.
func NewFSLoader(fsys fs.FS, dir string) *FSLoader {
	return &FSLoader{fsys: fsys, dir: dir}
}

func (l *FSLoader) Load(name string) (*LoadedWorkflow, error) {
	workflowDir := l.dir + "/" + name

	// Load profile.yaml
	profilePath := workflowDir + "/profile.yaml"
	data, err := fs.ReadFile(l.fsys, profilePath)
	if err != nil {
		return nil, fmt.Errorf("workflow %q not found: %w", name, err)
	}

	w, err := workflow.ParseYAML(data)
	if err != nil {
		return nil, fmt.Errorf("parsing workflow %q: %w", name, err)
	}

	loaded := &LoadedWorkflow{
		Workflow:  w,
		Templates: make(map[string]*template.Template),
		Rubrics:   make(map[string]*rubric.RubricSet),
	}

	// Load templates
	templatesDir := workflowDir + "/templates"
	if err := l.loadTemplates(loaded, templatesDir); err != nil {
		// Templates are optional
	}

	// Load rubrics
	rubricsDir := workflowDir + "/rubrics"
	if err := l.loadRubrics(loaded, rubricsDir); err != nil {
		// Rubrics are optional
	}

	return loaded, nil
}

func (l *FSLoader) loadTemplates(loaded *LoadedWorkflow, dir string) error {
	entries, err := fs.ReadDir(l.fsys, dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		specType := strings.TrimSuffix(entry.Name(), ".md")
		content, err := fs.ReadFile(l.fsys, dir+"/"+entry.Name())
		if err != nil {
			return fmt.Errorf("reading template %s: %w", entry.Name(), err)
		}

		loaded.Templates[specType] = &template.Template{
			ID:       specType,
			SpecType: specType,
			Content:  oscompatfs.NormalizeLineEndings(string(content)),
		}
	}

	return nil
}

func (l *FSLoader) loadRubrics(loaded *LoadedWorkflow, dir string) error {
	entries, err := fs.ReadDir(l.fsys, dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".rubric.yaml") {
			continue
		}

		specType := strings.TrimSuffix(entry.Name(), ".rubric.yaml")
		data, err := fs.ReadFile(l.fsys, dir+"/"+entry.Name())
		if err != nil {
			return fmt.Errorf("reading rubric %s: %w", entry.Name(), err)
		}

		var rs rubric.RubricSet
		if err := yaml.Unmarshal(data, &rs); err != nil {
			return fmt.Errorf("parsing rubric %s: %w", entry.Name(), err)
		}

		loaded.Rubrics[specType] = &rs
	}

	return nil
}

func (l *FSLoader) Available() []string {
	entries, err := fs.ReadDir(l.fsys, l.dir)
	if err != nil {
		return nil
	}

	var names []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		profilePath := l.dir + "/" + entry.Name() + "/profile.yaml"
		if _, err := fs.Stat(l.fsys, profilePath); err == nil {
			names = append(names, entry.Name())
		}
	}
	return names
}

// ChainLoader tries multiple loaders in order.
type ChainLoader struct {
	loaders []Loader
}

// NewChainLoader creates a loader that tries multiple loaders in order.
// The first loader to successfully load a workflow wins.
func NewChainLoader(loaders ...Loader) *ChainLoader {
	return &ChainLoader{loaders: loaders}
}

func (l *ChainLoader) Load(name string) (*LoadedWorkflow, error) {
	var lastErr error
	for _, loader := range l.loaders {
		w, err := loader.Load(name)
		if err == nil {
			return w, nil
		}
		lastErr = err
	}
	return nil, lastErr
}

func (l *ChainLoader) Available() []string {
	seen := make(map[string]bool)
	var names []string

	for _, loader := range l.loaders {
		for _, name := range loader.Available() {
			if !seen[name] {
				seen[name] = true
				names = append(names, name)
			}
		}
	}
	return names
}

// ResolvingLoader resolves workflow inheritance (extends).
type ResolvingLoader struct {
	base Loader
}

// NewResolvingLoader creates a loader that resolves workflow inheritance.
func NewResolvingLoader(base Loader) *ResolvingLoader {
	return &ResolvingLoader{base: base}
}

func (l *ResolvingLoader) Load(name string) (*LoadedWorkflow, error) {
	return l.loadWithChain(name, nil)
}

// applyExplicitSources honors per-spec template/rubric provenance declared in
// spec_config (SpecRequirement.Template/.Rubric). A declared source overrides
// whatever the extends chain resolved, and is enforced: if the named source
// does not actually provide the file, loading fails. This makes provenance
// authoritative rather than implicit. The chain carries every workflow on the
// current resolution path (extends links and spec sources alike) so a source
// that leads back to a workflow already being resolved fails as a cycle
// instead of recursing forever.
func (l *ResolvingLoader) applyExplicitSources(loaded *LoadedWorkflow, name string, chain []string) error {
	for spec, req := range loaded.Workflow.SpecConfig {
		if req == nil {
			continue
		}
		if req.Template != nil {
			t, err := l.sourceTemplate(name, req.Template.From, spec, chain)
			if err != nil {
				return fmt.Errorf("workflow %q spec %q template: %w", name, spec, err)
			}
			loaded.Templates[spec] = t
		}
		if req.Rubric != nil {
			r, err := l.sourceRubric(name, req.Rubric.From, spec, chain)
			if err != nil {
				return fmt.Errorf("workflow %q spec %q rubric: %w", name, spec, err)
			}
			loaded.Rubrics[spec] = r
		}
	}
	return nil
}

// sourceTemplate resolves a spec's template from the named source workflow. The
// "local" sentinel (or the workflow's own name) resolves from the workflow's
// own directory only; any other name resolves from that workflow's fully
// resolved template set.
func (l *ResolvingLoader) sourceTemplate(self, from, spec string, chain []string) (*template.Template, error) {
	if from == "" {
		return nil, fmt.Errorf("empty source (use %q or a workflow name)", workflow.SourceLocal)
	}
	src, err := l.sourceWorkflow(self, from, chain)
	if err != nil {
		return nil, err
	}
	t, ok := src.Templates[spec]
	if !ok {
		return nil, fmt.Errorf("source %q has no template for %q", from, spec)
	}
	return t, nil
}

// sourceRubric resolves a spec's rubric from the named source workflow, with the
// same semantics as sourceTemplate.
func (l *ResolvingLoader) sourceRubric(self, from, spec string, chain []string) (*rubric.RubricSet, error) {
	if from == "" {
		return nil, fmt.Errorf("empty source (use %q or a workflow name)", workflow.SourceLocal)
	}
	src, err := l.sourceWorkflow(self, from, chain)
	if err != nil {
		return nil, err
	}
	r, ok := src.Rubrics[spec]
	if !ok {
		return nil, fmt.Errorf("source %q has no rubric for %q", from, spec)
	}
	return r, nil
}

// sourceWorkflow returns the assets to resolve a spec source against. For the
// "local" sentinel (or the workflow's own name) it returns the workflow's own
// directory assets only (no inheritance); otherwise it returns the named
// workflow's fully resolved assets, resolved on the same chain as the caller —
// a source naming a workflow whose resolution leads back to the declaring
// workflow (e.g. a descendant, whose extends chain necessarily returns here)
// is a cycle error, not infinite recursion.
func (l *ResolvingLoader) sourceWorkflow(self, from string, chain []string) (*LoadedWorkflow, error) {
	if from == workflow.SourceLocal || from == self {
		local, err := l.base.Load(self)
		if err != nil {
			return nil, fmt.Errorf("loading local assets: %w", err)
		}
		return local, nil
	}
	src, err := l.loadWithChain(from, chain)
	if err != nil {
		return nil, fmt.Errorf("loading source workflow %q: %w", from, err)
	}
	return src, nil
}

func (l *ResolvingLoader) loadWithChain(name string, chain []string) (*LoadedWorkflow, error) {
	// The chain holds every workflow currently being resolved, whether reached
	// via extends or via an explicit spec_config template/rubric source; any
	// path that returns to one of them is a cycle.
	if slices.Contains(chain, name) {
		return nil, fmt.Errorf("circular workflow reference detected (via extends or spec_config sources): %s", strings.Join(append(chain, name), " -> "))
	}

	loaded, err := l.base.Load(name)
	if err != nil {
		return nil, err
	}

	// Clone before appending: the chain is shared with sibling resolution paths
	// (each spec source and the extends link recurse from the same slice), so
	// growing it in place could leak this level's entry across paths.
	childChain := append(slices.Clone(chain), name)

	// Resolve this workflow's own declared template/rubric sources against real
	// workflows before merging upward. Resolving per level keeps the "local"
	// sentinel pinned to the workflow that declared it: once baked into this
	// level's Templates/Rubrics maps, a child that inherits the entry reuses the
	// resolved file rather than re-resolving "local" against its own directory.
	if err := l.applyExplicitSources(loaded, name, childChain); err != nil {
		return nil, err
	}

	// If this workflow extends another, load and merge
	if loaded.Workflow.Extends != "" {
		parent, err := l.loadWithChain(loaded.Workflow.Extends, childChain)
		if err != nil {
			return nil, fmt.Errorf("loading parent workflow %q: %w", loaded.Workflow.Extends, err)
		}
		loaded = mergeWorkflows(loaded, parent)
	}

	return loaded, nil
}

func (l *ResolvingLoader) Available() []string {
	return l.base.Available()
}

// mergeWorkflows merges a child workflow with its parent.
// Child settings override parent settings.
func mergeWorkflows(child, parent *LoadedWorkflow) *LoadedWorkflow {
	merged := &LoadedWorkflow{
		Workflow:  child.Workflow.Merge(parent.Workflow),
		Templates: make(map[string]*template.Template),
		Rubrics:   make(map[string]*rubric.RubricSet),
	}

	// Copy parent templates, then override with child
	maps.Copy(merged.Templates, parent.Templates)
	maps.Copy(merged.Templates, child.Templates)

	// Copy parent rubrics, then override with child
	maps.Copy(merged.Rubrics, parent.Rubrics)
	maps.Copy(merged.Rubrics, child.Rubrics)

	return merged
}
