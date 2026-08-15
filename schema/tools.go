//go:build tools

// This file pins build-only dependencies used by the schema generator
// (generate.go, which carries a //go:build ignore tag and is therefore invisible
// to `go mod tidy`). Without this blank import, `go mod tidy` strips
// github.com/invopop/jsonschema from go.mod/go.sum and `go generate ./schema/...`
// fails with a missing go.sum entry. The file is never compiled into the
// library: the "tools" build tag excludes it from normal builds.
package schema

import (
	_ "github.com/invopop/jsonschema"
)
