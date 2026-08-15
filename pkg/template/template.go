// Package template defines the structure of spec templates.
//
// Templates provide the scaffolding for spec documents, including
// required sections, optional sections, and guidance for authors.
package template

// Template defines the structure of a spec template.
type Template struct {
	// ID is the template identifier (usually matches spec type ID).
	ID string `json:"id" jsonschema:"required,description=Template identifier"`

	// SpecType is the spec type this template produces.
	SpecType string `json:"specType" jsonschema:"required,description=Spec type ID this template produces"`

	// Version is the template version.
	Version string `json:"version,omitempty" jsonschema:"description=Template version"`

	// Description explains the template purpose.
	Description string `json:"description,omitempty" jsonschema:"description=Template purpose"`

	// Sections define the document structure.
	Sections []Section `json:"sections,omitempty" jsonschema:"description=Template sections"`

	// Frontmatter defines YAML frontmatter schema.
	Frontmatter *Frontmatter `json:"frontmatter,omitempty" jsonschema:"description=YAML frontmatter schema"`

	// Metadata about the template.
	Metadata *TemplateMetadata `json:"metadata,omitempty" jsonschema:"description=Template metadata"`

	// Content is the raw markdown template with placeholders.
	// Populated when loading templates via embed for filesystem-free access.
	Content string `json:"content,omitempty" jsonschema:"description=Raw markdown template content with placeholders"`
}

// Section defines a section within a template.
type Section struct {
	// ID is the section identifier.
	ID string `json:"id" jsonschema:"required,description=Section identifier"`

	// Title is the section heading.
	Title string `json:"title" jsonschema:"required,description=Section heading"`

	// Description explains what belongs in this section.
	Description string `json:"description,omitempty" jsonschema:"description=Section guidance"`

	// Required indicates whether this section must be present.
	Required bool `json:"required,omitempty" jsonschema:"description=Whether section is required"`

	// Level is the heading level (1-6).
	Level int `json:"level,omitempty" jsonschema:"minimum=1,maximum=6,description=Heading level"`

	// Subsections are nested sections.
	Subsections []Section `json:"subsections,omitempty" jsonschema:"description=Nested sections"`

	// Placeholder is example content for the section.
	Placeholder string `json:"placeholder,omitempty" jsonschema:"description=Example content"`

	// MinLength is the minimum content length (characters).
	MinLength int `json:"minLength,omitempty" jsonschema:"minimum=0,description=Minimum content length"`

	// MaxLength is the maximum content length (characters).
	MaxLength int `json:"maxLength,omitempty" jsonschema:"minimum=0,description=Maximum content length"`
}

// Frontmatter defines the YAML frontmatter schema.
type Frontmatter struct {
	// Required are required frontmatter fields.
	Required []FrontmatterField `json:"required,omitempty" jsonschema:"description=Required frontmatter fields"`

	// Optional are optional frontmatter fields.
	Optional []FrontmatterField `json:"optional,omitempty" jsonschema:"description=Optional frontmatter fields"`
}

// FrontmatterField defines a frontmatter field.
type FrontmatterField struct {
	// Name is the field name.
	Name string `json:"name" jsonschema:"required,description=Field name"`

	// Type is the field type (string, number, boolean, array, object).
	Type string `json:"type" jsonschema:"required,enum=string,enum=number,enum=boolean,enum=array,enum=object"`

	// Description explains the field purpose.
	Description string `json:"description,omitempty" jsonschema:"description=Field purpose"`

	// Default is the default value.
	Default any `json:"default,omitempty" jsonschema:"description=Default value"`

	// Enum lists allowed values.
	Enum []string `json:"enum,omitempty" jsonschema:"description=Allowed values"`
}

// TemplateMetadata contains template provenance information.
type TemplateMetadata struct {
	// Author is the template author.
	Author string `json:"author,omitempty" jsonschema:"description=Template author"`

	// Source is the methodology or company this template derives from.
	Source string `json:"source,omitempty" jsonschema:"description=Template origin"`

	// LastUpdated is the last modification date.
	LastUpdated string `json:"lastUpdated,omitempty" jsonschema:"format=date,description=Last modification date"`

	// Tags are categorization tags.
	Tags []string `json:"tags,omitempty" jsonschema:"description=Categorization tags"`
}
