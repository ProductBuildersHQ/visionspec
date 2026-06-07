package version

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ProductBuildersHQ/visionspec/pkg/config"
	"github.com/ProductBuildersHQ/visionspec/pkg/types"
)

// DiffLine represents a single line in a diff.
type DiffLine struct {
	Type    DiffType
	Content string
	LineNum int
}

// DiffType indicates whether a line was added, removed, or unchanged.
type DiffType string

const (
	DiffTypeUnchanged DiffType = " "
	DiffTypeAdded     DiffType = "+"
	DiffTypeRemoved   DiffType = "-"
)

// DiffResult contains the diff between two versions.
type DiffResult struct {
	SpecType   types.SpecType
	OldVersion int
	NewVersion int
	Lines      []DiffLine
	Additions  int
	Deletions  int
	Unchanged  int
}

// Diff compares two versions of a spec.
func Diff(projectPath string, specType types.SpecType, oldVersion, newVersion int) (*DiffResult, error) {
	// Get old version content
	var oldContent string
	if oldVersion == 0 {
		// Compare with nothing (show all as additions)
		oldContent = ""
	} else {
		_, content, err := GetVersion(projectPath, specType, oldVersion)
		if err != nil {
			return nil, fmt.Errorf("getting old version: %w", err)
		}
		oldContent = content
	}

	// Get new version content
	var newContent string
	if newVersion == 0 {
		// Compare with current file
		specPath := config.SpecPath(projectPath, specType)
		data, err := os.ReadFile(specPath)
		if err != nil {
			return nil, fmt.Errorf("reading current spec: %w", err)
		}
		newContent = string(data)
	} else {
		_, content, err := GetVersion(projectPath, specType, newVersion)
		if err != nil {
			return nil, fmt.Errorf("getting new version: %w", err)
		}
		newContent = content
	}

	// Perform diff
	return diffStrings(specType, oldVersion, newVersion, oldContent, newContent), nil
}

// DiffWithCurrent compares a version with the current spec file.
func DiffWithCurrent(projectPath string, specType types.SpecType, versionNum int) (*DiffResult, error) {
	return Diff(projectPath, specType, versionNum, 0)
}

// diffStrings performs a simple line-by-line diff.
// This is a simplified diff algorithm; for production, consider using a proper LCS algorithm.
func diffStrings(specType types.SpecType, oldVersion, newVersion int, oldContent, newContent string) *DiffResult {
	oldLines := splitLines(oldContent)
	newLines := splitLines(newContent)

	result := &DiffResult{
		SpecType:   specType,
		OldVersion: oldVersion,
		NewVersion: newVersion,
		Lines:      []DiffLine{},
	}

	// Simple diff: use longest common subsequence
	lcs := computeLCS(oldLines, newLines)
	lcsSet := make(map[string]bool)
	for _, line := range lcs {
		lcsSet[line] = true
	}

	oldIdx, newIdx := 0, 0
	lineNum := 0

	for oldIdx < len(oldLines) || newIdx < len(newLines) {
		lineNum++

		if oldIdx < len(oldLines) && newIdx < len(newLines) && oldLines[oldIdx] == newLines[newIdx] {
			// Lines match - unchanged
			result.Lines = append(result.Lines, DiffLine{
				Type:    DiffTypeUnchanged,
				Content: oldLines[oldIdx],
				LineNum: lineNum,
			})
			result.Unchanged++
			oldIdx++
			newIdx++
		} else if newIdx < len(newLines) && (oldIdx >= len(oldLines) || !lcsSet[newLines[newIdx]] || (oldIdx < len(oldLines) && lcsSet[oldLines[oldIdx]])) {
			// Line in new but not in old - addition
			if oldIdx >= len(oldLines) || !contains(oldLines[oldIdx:], newLines[newIdx]) {
				result.Lines = append(result.Lines, DiffLine{
					Type:    DiffTypeAdded,
					Content: newLines[newIdx],
					LineNum: lineNum,
				})
				result.Additions++
				newIdx++
			} else {
				// Line removed from old
				result.Lines = append(result.Lines, DiffLine{
					Type:    DiffTypeRemoved,
					Content: oldLines[oldIdx],
					LineNum: lineNum,
				})
				result.Deletions++
				oldIdx++
			}
		} else if oldIdx < len(oldLines) {
			// Line in old but not in new - removal
			result.Lines = append(result.Lines, DiffLine{
				Type:    DiffTypeRemoved,
				Content: oldLines[oldIdx],
				LineNum: lineNum,
			})
			result.Deletions++
			oldIdx++
		}
	}

	return result
}

// splitLines splits content into lines.
func splitLines(content string) []string {
	if content == "" {
		return []string{}
	}
	scanner := bufio.NewScanner(strings.NewReader(content))
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

// contains checks if a slice contains a string.
func contains(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// computeLCS computes the longest common subsequence of two string slices.
func computeLCS(a, b []string) []string {
	m, n := len(a), len(b)
	if m == 0 || n == 0 {
		return []string{}
	}

	// Build LCS table
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}

	// Backtrack to find LCS
	var lcs []string
	i, j := m, n
	for i > 0 && j > 0 {
		if a[i-1] == b[j-1] {
			lcs = append([]string{a[i-1]}, lcs...)
			i--
			j--
		} else if dp[i-1][j] > dp[i][j-1] {
			i--
		} else {
			j--
		}
	}

	return lcs
}

// FormatDiff formats a diff result for terminal output.
func (d *DiffResult) FormatDiff() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Diff: %s v%d → v%d\n", d.SpecType, d.OldVersion, d.NewVersion))
	sb.WriteString(fmt.Sprintf("+%d additions, -%d deletions, %d unchanged\n\n", d.Additions, d.Deletions, d.Unchanged))

	for _, line := range d.Lines {
		sb.WriteString(fmt.Sprintf("%s %s\n", line.Type, line.Content))
	}

	return sb.String()
}

// FormatCompact formats a diff in compact form (changes only).
func (d *DiffResult) FormatCompact() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Changes: %s v%d → v%d (+%d/-%d)\n", d.SpecType, d.OldVersion, d.NewVersion, d.Additions, d.Deletions))

	for _, line := range d.Lines {
		if line.Type != DiffTypeUnchanged {
			sb.WriteString(fmt.Sprintf("%s %s\n", line.Type, line.Content))
		}
	}

	return sb.String()
}
