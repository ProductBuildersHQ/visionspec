module github.com/ProductBuildersHQ/visionspec/tools/prism-sync

go 1.25.5

require (
	github.com/grokify/prism-roadmap v0.16.0
	gopkg.in/yaml.v3 v3.0.1
)

require github.com/plexusone/structured-evaluation v0.13.0 // indirect

// TODO(release): drop this replace once the v2mom/ost templates and rubrics
// land in a tagged prism-roadmap release, then require that version above.
// This directive must not survive a push (see Pre-Push Checklist).
replace github.com/grokify/prism-roadmap => ../../../../grokify/prism-roadmap
