module github.com/ProductBuildersHQ/specification-workflow-spec

go 1.25.5

require (
	github.com/grokify/oscompat v0.5.0
	github.com/invopop/jsonschema v0.14.0
	github.com/plexusone/structured-evaluation v0.14.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/buger/jsonparser v1.6.1 // indirect
	github.com/pb33f/ordered-map/v2 v2.3.1 // indirect
	go.yaml.in/yaml/v4 v4.0.0-rc.6 // indirect
)

// TODO(release): drop this replace once structured-evaluation's layered-rubric
// schema (Class/Blocking/Evaluation, JudgeInstructions) is tagged v0.14.0, then
// bump the require above to that version. This directive must not survive a
// push (see Pre-Push Checklist) — see INIT-SPECWORKFLOWSPEC-001 TRD-012/§5.
// replace github.com/plexusone/structured-evaluation => ../../plexusone/structured-evaluation
