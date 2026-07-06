# MCP Tools

The VisionSpec MCP server provides the following tools for project management, spec operations, and draft authoring.

## Project Tools

### list_projects

List all VisionSpec projects in the current repository.

**Arguments:** None

**Returns:**

```json
{
  "projects": ["user-onboarding", "payment-flow"],
  "specs_dir": "docs/specs",
  "count": 2
}
```

### get_project_status

Get the status and readiness of a specific project.

**Arguments:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `project` | string | Yes | Project name |

**Returns:**

```json
{
  "project": "user-onboarding",
  "path": "docs/specs/user-onboarding",
  "generated_at": "2024-01-15T10:30:00Z",
  "readiness": {
    "ready": false,
    "summary": "Not ready: 2 blockers",
    "gates": [
      {"name": "Required specs present", "passed": true, "message": "..."},
      {"name": "Evaluations passing", "passed": true, "message": "..."},
      {"name": "Approvals obtained", "passed": false, "message": "..."},
      {"name": "Execution spec generated", "passed": false, "message": "..."}
    ]
  },
  "specs": [...],
  "summary": {
    "total_specs": 10,
    "present_specs": 3,
    "evaluated_specs": 0,
    "approved_specs": 0
  }
}
```

### get_workflow

Get the workflow DAG for a project, showing spec dependencies, phases, and progress.

**Arguments:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `project` | string | Yes | Project name |
| `include_mermaid` | boolean | No | Include Mermaid diagram in response |
| `format` | string | No | Output format: `json` (default), `mermaid`, `dot` |

**Returns:**

```json
{
  "name": "big-tech-product",
  "description": "Big Tech methodology for products",
  "phases": [
    {"id": "discovery", "name": "Discovery", "order": 1, "nodes": ["mrd"]},
    {"id": "vision", "name": "Vision", "order": 2, "nodes": ["press", "faq"]}
  ],
  "nodes": {
    "mrd": {
      "id": "mrd",
      "name": "Market Requirements",
      "phase": "discovery",
      "status": "completed",
      "depends_on": [],
      "automated": false
    },
    "press": {
      "id": "press",
      "name": "Press Release",
      "phase": "vision",
      "status": "in_progress",
      "depends_on": ["mrd"],
      "automated": true
    }
  },
  "progress": {
    "completed": 3,
    "total": 10,
    "percent": 30.0
  },
  "mermaid": "graph TD\n    mrd[Market Requirements]\n    ..."
}
```

## Spec Tools

### get_spec

Get the content of a specification.

**Arguments:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `project` | string | Yes | Project name |
| `spec_type` | string | Yes | Spec type (mrd, prd, uxd, trd, ird, press, faq, narrative-1p, narrative-6p) |

**Returns:**

```json
{
  "project": "user-onboarding",
  "spec_type": "prd",
  "path": "docs/specs/user-onboarding/source/prd.md",
  "content": "# Product Requirements Document\n..."
}
```

### get_eval

Get evaluation results for a specification.

**Arguments:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `project` | string | Yes | Project name |
| `spec_type` | string | Yes | Spec type |

**Returns:**

```json
{
  "project": "user-onboarding",
  "spec_type": "prd",
  "evaluation": {
    "decision": "pass",
    "score": 0.85,
    "findings": [...]
  }
}
```

### run_eval

Run evaluation on a specification using LLM judges.

**Arguments:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `project` | string | Yes | Project name |
| `spec_type` | string | Yes | Spec type |

**Returns:**

```json
{
  "project": "user-onboarding",
  "spec_type": "prd",
  "decision": "pass",
  "findings_count": 3,
  "message": "Evaluation completed successfully"
}
```

### synthesize

Generate a spec from source documents using LLM synthesis.

**Arguments:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `project` | string | Yes | Project name |
| `spec_type` | string | Yes | Spec type to generate (trd, ird, press, faq, narrative-1p, narrative-6p) |

**Returns:**

```json
{
  "project": "user-onboarding",
  "spec_type": "trd",
  "path": "docs/specs/user-onboarding/technical/trd.md",
  "sources": ["mrd", "prd", "uxd"],
  "message": "Synthesis completed successfully"
}
```

### reconcile

Generate unified execution spec from approved specs.

**Arguments:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `project` | string | Yes | Project name |

**Returns:**

```json
{
  "project": "user-onboarding",
  "path": "docs/specs/user-onboarding/spec.md",
  "sources": ["mrd", "prd", "uxd", "trd"],
  "message": "Reconciliation completed successfully"
}
```

### approve

Approve a specification for reconciliation.

**Arguments:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `project` | string | Yes | Project name |
| `spec_type` | string | Yes | Spec type |
| `approver` | string | No | Approver identifier |

**Returns:**

```json
{
  "project": "user-onboarding",
  "spec_type": "prd",
  "approver": "pm@example.com",
  "approved_at": "2024-01-15T10:30:00Z",
  "message": "Spec approved successfully"
}
```

### export

Export specs to a target execution system.

**Arguments:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `project` | string | Yes | Project name |
| `target` | string | Yes | Target system (speckit, gsd, gastown, gascity) |

**Returns:**

```json
{
  "project": "user-onboarding",
  "target": "speckit",
  "output_dir": "docs/specs/user-onboarding/.speckit",
  "files": ["spec.md", "plan.md", "tasks.md"],
  "message": "Export completed successfully"
}
```

## Draft Authoring Tools

These tools support the AI-assisted authoring workflow for creating and refining specs.

### start_draft

Initialize a new draft for a spec type.

**Arguments:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `project` | string | Yes | Project name |
| `spec_type` | string | Yes | Spec type (mrd, prd, uxd) |

**Returns:**

```json
{
  "project": "user-onboarding",
  "spec_type": "prd",
  "template": "# Product Requirements Document\n...",
  "draft_path": "docs/specs/user-onboarding/source/prd.draft.md",
  "message": "Draft initialized"
}
```

### get_draft

Get the current draft content.

**Arguments:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `project` | string | Yes | Project name |
| `spec_type` | string | Yes | Spec type |

**Returns:**

```json
{
  "project": "user-onboarding",
  "spec_type": "prd",
  "content": "# Product Requirements Document\n...",
  "metadata": {
    "version": 3,
    "started_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "eval_history": [...]
  }
}
```

### update_draft

Save updated draft content.

**Arguments:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `project` | string | Yes | Project name |
| `spec_type` | string | Yes | Spec type |
| `content` | string | Yes | Draft content |

**Returns:**

```json
{
  "project": "user-onboarding",
  "spec_type": "prd",
  "version": 4,
  "message": "Draft updated"
}
```

### eval_draft

Evaluate the current draft against the rubric.

**Arguments:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `project` | string | Yes | Project name |
| `spec_type` | string | Yes | Spec type |

**Returns:**

```json
{
  "project": "user-onboarding",
  "spec_type": "prd",
  "decision": "partial",
  "findings": [
    {
      "category": "Problem Clarity",
      "rating": "pass",
      "rationale": "Problem statement is clear and well-defined"
    },
    {
      "category": "User Stories",
      "rating": "partial",
      "rationale": "Some user stories lack acceptance criteria"
    }
  ],
  "suggestions": ["Add acceptance criteria to user stories in section 3"]
}
```

### finalize_draft

Promote draft to final spec.

**Arguments:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `project` | string | Yes | Project name |
| `spec_type` | string | Yes | Spec type |

**Returns:**

```json
{
  "project": "user-onboarding",
  "spec_type": "prd",
  "path": "docs/specs/user-onboarding/source/prd.md",
  "message": "Draft finalized"
}
```

### discard_draft

Delete a draft without saving.

**Arguments:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `project` | string | Yes | Project name |
| `spec_type` | string | Yes | Spec type |

**Returns:**

```json
{
  "project": "user-onboarding",
  "spec_type": "prd",
  "message": "Draft discarded"
}
```

### list_drafts

List all drafts in a project.

**Arguments:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `project` | string | Yes | Project name |

**Returns:**

```json
{
  "project": "user-onboarding",
  "drafts": [
    {
      "spec_type": "prd",
      "version": 3,
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

## Error Handling

All tools return errors in a consistent format:

```json
{
  "error": "error message here"
}
```

The response will have `isError: true` set in the MCP CallToolResult.
