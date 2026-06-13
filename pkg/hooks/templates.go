package hooks

import "strings"

// GetTemplate returns the script template for a hook type.
func GetTemplate(hookType HookType) string {
	switch hookType {
	case HookPreCommit:
		return preCommitTemplate
	case HookPrePush:
		return prePushTemplate
	case HookCommitMsg:
		return commitMsgTemplate
	case HookPostCommit:
		return postCommitTemplate
	default:
		return ""
	}
}

// GetTemplateDescription returns a description of what a hook does.
func GetTemplateDescription(hookType HookType) string {
	switch hookType {
	case HookPreCommit:
		return "Lints changed spec files before commit"
	case HookPrePush:
		return "Validates specs and checks for blockers before push"
	case HookCommitMsg:
		return "Validates commit message format"
	case HookPostCommit:
		return "Updates spec status after commit"
	default:
		return ""
	}
}

const preCommitTemplate = `#!/bin/bash
# visionspec pre-commit hook
# Lints changed spec files before commit

set -e

# Find visionspec executable
VISIONSPEC="visionspec"
if ! command -v "$VISIONSPEC" &> /dev/null; then
    # Try common paths
    if [ -x "./visionspec" ]; then
        VISIONSPEC="./visionspec"
    elif [ -x "$HOME/go/bin/visionspec" ]; then
        VISIONSPEC="$HOME/go/bin/visionspec"
    else
        echo "visionspec not found, skipping pre-commit hook"
        exit 0
    fi
fi

# Get changed spec files
CHANGED_SPECS=$(git diff --cached --name-only --diff-filter=ACM | grep -E '^docs/specs/.*\.md$' || true)

if [ -z "$CHANGED_SPECS" ]; then
    # No spec files changed, nothing to lint
    exit 0
fi

echo "Linting changed spec files..."

# Run lint on each changed file's project
PROJECTS=$(echo "$CHANGED_SPECS" | sed -E 's|docs/specs/([^/]+)/.*|\1|' | sort -u)

FAILED=0
for PROJECT in $PROJECTS; do
    echo "Checking project: $PROJECT"
    if ! $VISIONSPEC lint "$PROJECT" 2>/dev/null; then
        echo "Lint failed for $PROJECT"
        FAILED=1
    fi
done

if [ $FAILED -eq 1 ]; then
    echo ""
    echo "Spec lint failed. Please fix the issues before committing."
    echo "Run 'visionspec lint' for details."
    exit 1
fi

echo "Spec lint passed."
exit 0
`

const prePushTemplate = `#!/bin/bash
# visionspec pre-push hook
# Validates specs before push

set -e

# Find visionspec executable
VISIONSPEC="visionspec"
if ! command -v "$VISIONSPEC" &> /dev/null; then
    if [ -x "./visionspec" ]; then
        VISIONSPEC="./visionspec"
    elif [ -x "$HOME/go/bin/visionspec" ]; then
        VISIONSPEC="$HOME/go/bin/visionspec"
    else
        echo "visionspec not found, skipping pre-push hook"
        exit 0
    fi
fi

# Get the remote and branch being pushed to
REMOTE="$1"
URL="$2"

echo "Validating specs before push to $REMOTE..."

# Run status check
if ! $VISIONSPEC status --format=json 2>/dev/null | grep -q '"has_blockers":false'; then
    echo ""
    echo "WARNING: Specs have blockers. Consider resolving before push."
    echo "Run 'visionspec status' for details."
    # Note: This is a warning, not a failure. Remove 'exit 0' to make it blocking.
    exit 0
fi

# Run drift check if spec.md exists
if [ -f "spec.md" ] || [ -f "docs/specs/*/spec.md" ]; then
    echo "Checking for spec drift..."
    if ! $VISIONSPEC drift --ci 2>/dev/null; then
        echo ""
        echo "WARNING: Spec drift detected."
        echo "Run 'visionspec drift' for details."
        # Note: This is a warning, not a failure. Remove 'exit 0' to make it blocking.
        exit 0
    fi
fi

echo "Spec validation passed."
exit 0
`

const commitMsgTemplate = `#!/bin/bash
# visionspec commit-msg hook
# Validates commit message format

COMMIT_MSG_FILE="$1"
COMMIT_MSG=$(cat "$COMMIT_MSG_FILE")

# Check for conventional commit format
if ! echo "$COMMIT_MSG" | head -1 | grep -qE '^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert)(\(.+\))?: .+'; then
    echo ""
    echo "Commit message should follow conventional commits format:"
    echo "  <type>(<scope>): <description>"
    echo ""
    echo "Types: feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert"
    echo ""
    echo "Example: feat(auth): add OAuth2 login support"
    echo ""
    # Note: This is a warning. Remove 'exit 0' to make it blocking.
    exit 0
fi

exit 0
`

const postCommitTemplate = `#!/bin/bash
# visionspec post-commit hook
# Updates spec status after commit (informational only)

# Find visionspec executable
VISIONSPEC="visionspec"
if ! command -v "$VISIONSPEC" &> /dev/null; then
    exit 0
fi

# Get changed files in the last commit
CHANGED_SPECS=$(git diff --name-only HEAD~1 HEAD 2>/dev/null | grep -E '^docs/specs/.*\.md$' || true)

if [ -n "$CHANGED_SPECS" ]; then
    echo ""
    echo "Spec files updated. Consider running:"
    echo "  visionspec status    # Check overall spec status"
    echo "  visionspec eval      # Evaluate spec quality"
fi

exit 0
`

// CustomTemplate allows creating a custom hook template.
type CustomTemplate struct {
	HookType    HookType
	Description string
	Script      string
}

// GenerateCustomHook generates a hook script from a custom template.
func GenerateCustomHook(custom CustomTemplate) string {
	var sb strings.Builder

	sb.WriteString("#!/bin/bash\n")
	sb.WriteString("# visionspec " + string(custom.HookType) + " hook (custom)\n")
	if custom.Description != "" {
		sb.WriteString("# " + custom.Description + "\n")
	}
	sb.WriteString("\n")
	sb.WriteString(custom.Script)
	sb.WriteString("\n")

	return sb.String()
}
