package speckit

import (
	"path/filepath"
	"testing"

	"github.com/ProductBuildersHQ/visionspec/pkg/workflows"
)

// TestCommittedArtifactsInSync is the drift gate for the generated artifacts
// committed under the repo's speckit/ directory: they must match what the
// generator produces from the embedded workflow. Regenerate with
// `visionspec speckit export` after changing the family or the exporter.
func TestCommittedArtifactsInSync(t *testing.T) {
	lw, err := workflows.DefaultLoader().Load(testFamily)
	if err != nil {
		t.Fatal(err)
	}
	files, err := BuildExtensions(lw, Options{})
	if err != nil {
		t.Fatal(err)
	}
	wf, err := BuildWorkflow(lw, Options{})
	if err != nil {
		t.Fatal(err)
	}
	files = append(files, wf)

	dir := filepath.Join("..", "..", "speckit")
	drift, err := Check(dir, files)
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range drift {
		t.Errorf("committed artifact out of sync: speckit/%s", p)
	}
	if len(drift) > 0 {
		t.Log("regenerate with: go run ./cmd/visionspec speckit export aws-one-way-door -o speckit")
	}
}
