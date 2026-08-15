# Release v0.1.0

**Release Date:** 2026-08-02

Initial release of the specification-workflow-spec library, providing Go types and utilities for spec-driven product development workflows.

## Overview

This library defines the core primitives for VisionSpec-style specification workflows, including spec types, synthesis rules, evaluation rubrics, and filesystem layout conventions.

## Features

### Spec Type Registry (`pkg/spectype`)

Comprehensive registry of specification types organized by category:

- **Source specs:** MRD, PRD, UXD, TRD, IRD
- **GTM specs:** Press Release, FAQ, Narrative-6P
- **Strategic specs:** V2MOM variants (Vision, Values, Methods, Obstacles, Measures)
- **Execution specs:** TPD (Test Plan)
- **Output specs:** Reconciled Spec, Current Truth

### Synthesis Engine (`pkg/synthesis`)

Defines spec generation dependency rules for AI-assisted synthesis:

- Maps which specs can be generated from which inputs
- Supports multi-source synthesis (e.g., TRD from PRD + UXD)
- Enables workflow automation

### Profile Configuration (`pkg/profile`)

YAML-based workflow configuration:

- Define custom spec workflows per project
- Configure evaluation rubrics and gates
- Enable/disable specific spec types

### Diagram Generation (`pkg/diagram`)

Visual workflow representation:

- D2 diagram generation for spec DAGs
- Mermaid diagram generation
- Supports category grouping and styling

### Layout Conventions (`pkg/layout`)

Cross-platform filesystem layout:

- Standard directory structure for spec projects
- Configurable paths for source, GTM, technical, and execution specs
- Platform-agnostic path handling

### Supporting Packages

- **`pkg/gate`**: Quality gate definitions
- **`pkg/rubricdef`**: Evaluation rubric definitions
- **`pkg/template`**: Spec template primitives
- **`schema/`**: JSON Schema definitions for validation

## Installation

```bash
go get github.com/ProductBuildersHQ/specification-workflow-spec@v0.1.0
```

## Requirements

- Go 1.24 or later

## Links

- [Repository](https://github.com/ProductBuildersHQ/specification-workflow-spec)
- [Changelog](https://github.com/ProductBuildersHQ/specification-workflow-spec/blob/main/CHANGELOG.md)
