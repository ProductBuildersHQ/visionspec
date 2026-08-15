// Package gate defines phase gates and approval checkpoints.
//
// Gates are decision points in the specification workflow where
// human approval or automated checks must pass before proceeding.
package gate

// Gate defines an approval checkpoint in the workflow.
type Gate struct {
	// ID is the gate identifier.
	ID string `json:"id" jsonschema:"required,description=Gate identifier"`

	// Name is the human-readable gate name.
	Name string `json:"name" jsonschema:"required,description=Gate name"`

	// Description explains the gate purpose.
	Description string `json:"description,omitempty" jsonschema:"description=Gate purpose"`

	// Type indicates the gate type.
	Type GateType `json:"type" jsonschema:"required,enum=approval,enum=evaluation,enum=automated"`

	// After is the spec type ID after which this gate applies.
	After string `json:"after,omitempty" jsonschema:"description=Spec type ID after which gate applies"`

	// Phase is the phase ID after which this gate applies.
	Phase string `json:"phase,omitempty" jsonschema:"description=Phase ID after which gate applies"`

	// Conditions are the requirements to pass this gate.
	Conditions []Condition `json:"conditions,omitempty" jsonschema:"description=Gate conditions"`

	// Approvers defines who can approve this gate.
	Approvers *ApproverConfig `json:"approvers,omitempty" jsonschema:"description=Approver configuration"`

	// Required indicates whether this gate is mandatory.
	Required bool `json:"required,omitempty" jsonschema:"description=Whether gate is mandatory"`

	// Blocking indicates whether failure blocks progress.
	Blocking bool `json:"blocking,omitempty" jsonschema:"description=Whether failure blocks progress"`
}

// GateType indicates the type of gate.
type GateType string

const (
	// GateApproval requires human approval.
	GateApproval GateType = "approval"

	// GateEvaluation requires evaluation score threshold.
	GateEvaluation GateType = "evaluation"

	// GateAutomated requires automated checks to pass.
	GateAutomated GateType = "automated"
)

// Condition defines a gate requirement.
type Condition struct {
	// Type is the condition type.
	Type ConditionType `json:"type" jsonschema:"required,enum=spec_exists,enum=spec_approved,enum=eval_score,enum=no_blocking_findings,enum=custom"`

	// SpecType is the spec type ID (for spec-related conditions).
	SpecType string `json:"specType,omitempty" jsonschema:"description=Spec type ID"`

	// MinScore is the minimum evaluation score (for eval_score).
	MinScore float64 `json:"minScore,omitempty" jsonschema:"minimum=0,maximum=100,description=Minimum score"`

	// Custom is a custom condition expression.
	Custom string `json:"custom,omitempty" jsonschema:"description=Custom condition expression"`
}

// ConditionType indicates the type of condition.
type ConditionType string

const (
	// ConditionSpecExists checks that a spec file exists.
	ConditionSpecExists ConditionType = "spec_exists"

	// ConditionSpecApproved checks that a spec has been approved.
	ConditionSpecApproved ConditionType = "spec_approved"

	// ConditionEvalScore checks evaluation score meets threshold.
	ConditionEvalScore ConditionType = "eval_score"

	// ConditionNoBlockingFindings checks no critical/high findings.
	ConditionNoBlockingFindings ConditionType = "no_blocking_findings"

	// ConditionCustom uses a custom expression.
	ConditionCustom ConditionType = "custom"
)

// ApproverConfig defines who can approve a gate.
type ApproverConfig struct {
	// Roles are the roles that can approve (e.g., "tech_lead", "stakeholder").
	Roles []string `json:"roles,omitempty" jsonschema:"description=Approver roles"`

	// Users are specific user IDs that can approve.
	Users []string `json:"users,omitempty" jsonschema:"description=Specific approver user IDs"`

	// MinApprovals is the minimum number of approvals required.
	MinApprovals int `json:"minApprovals,omitempty" jsonschema:"minimum=1,description=Minimum approvals required"`

	// RequireAll indicates all listed approvers must approve.
	RequireAll bool `json:"requireAll,omitempty" jsonschema:"description=Whether all approvers must approve"`
}

// GateResult represents the outcome of a gate check.
type GateResult struct {
	// GateID is the gate that was checked.
	GateID string `json:"gateId" jsonschema:"required,description=Gate identifier"`

	// Passed indicates whether the gate passed.
	Passed bool `json:"passed" jsonschema:"required,description=Whether gate passed"`

	// ConditionResults are the individual condition outcomes.
	ConditionResults []ConditionResult `json:"conditionResults,omitempty" jsonschema:"description=Individual condition outcomes"`

	// Approvals are the recorded approvals.
	Approvals []Approval `json:"approvals,omitempty" jsonschema:"description=Recorded approvals"`

	// BlockingReason explains why the gate failed.
	BlockingReason string `json:"blockingReason,omitempty" jsonschema:"description=Failure reason"`

	// CheckedAt is when the gate was checked.
	CheckedAt string `json:"checkedAt,omitempty" jsonschema:"format=date-time,description=Check timestamp"`
}

// ConditionResult is the outcome of a single condition check.
type ConditionResult struct {
	// Type is the condition type.
	Type ConditionType `json:"type" jsonschema:"required,description=Condition type"`

	// Passed indicates whether the condition passed.
	Passed bool `json:"passed" jsonschema:"required,description=Whether condition passed"`

	// Detail provides additional context.
	Detail string `json:"detail,omitempty" jsonschema:"description=Additional context"`
}

// Approval records a human approval.
type Approval struct {
	// ApproverID is the approver's user ID.
	ApproverID string `json:"approverId" jsonschema:"required,description=Approver user ID"`

	// ApproverRole is the approver's role.
	ApproverRole string `json:"approverRole,omitempty" jsonschema:"description=Approver role"`

	// ApprovedAt is when approval was given.
	ApprovedAt string `json:"approvedAt" jsonschema:"required,format=date-time,description=Approval timestamp"`

	// Comment is an optional approval comment.
	Comment string `json:"comment,omitempty" jsonschema:"description=Approval comment"`
}
