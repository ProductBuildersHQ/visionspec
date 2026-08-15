// Package schema provides embedded JSON Schema files generated from Go types.
//
//go:generate go run generate.go
package schema

import "embed"

// FS embeds the generated JSON Schema files.
//
//go:embed *.schema.json
var FS embed.FS

// SchemaFiles lists the available schema files.
var SchemaFiles = []string{
	"spectype.schema.json",
	"workflow.schema.json",
	"template.schema.json",
	"synthesis.schema.json",
	"gate.schema.json",
	"layout.schema.json",
	"integration.schema.json",
	"pipeline.schema.json",
	"loop.schema.json",
}
