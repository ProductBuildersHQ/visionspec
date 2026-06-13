package target

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

func init() {
	Register(&OpenSpecTarget{})
}

// OpenSpecTarget exports to OpenSpec portable format.
// OpenSpec is a standardized JSON/YAML schema for requirements that can be
// imported by various tools and systems.
type OpenSpecTarget struct{}

// Name returns the target name.
func (t *OpenSpecTarget) Name() string {
	return "openspec"
}

// Description returns a description of the target.
func (t *OpenSpecTarget) Description() string {
	return "OpenSpec portable format (JSON/YAML) for cross-tool compatibility"
}

// Capabilities returns what this target supports.
func (t *OpenSpecTarget) Capabilities() Capabilities {
	return Capabilities{
		SequentialTasks:   true,
		ParallelExecution: true,
		MultiAgent:        false,
		Verification:      true,
		DependencyGraph:   true,
	}
}

// Validate checks if the spec can be exported to this target.
func (t *OpenSpecTarget) Validate(spec string) error {
	if spec == "" {
		return fmt.Errorf("spec content is empty")
	}
	return nil
}

// OpenSpecDocument is the root document structure.
type OpenSpecDocument struct {
	Version     string               `json:"version" yaml:"version"`
	Metadata    OpenSpecMetadata     `json:"metadata" yaml:"metadata"`
	Overview    OpenSpecOverview     `json:"overview" yaml:"overview"`
	Features    []OpenSpecFeature    `json:"features" yaml:"features"`
	Tasks       []OpenSpecTask       `json:"tasks" yaml:"tasks"`
	Acceptance  []OpenSpecCriteria   `json:"acceptance" yaml:"acceptance"`
	Constraints []OpenSpecConstraint `json:"constraints,omitempty" yaml:"constraints,omitempty"`
	Appendix    *OpenSpecAppendix    `json:"appendix,omitempty" yaml:"appendix,omitempty"`
}

// OpenSpecMetadata contains document metadata.
type OpenSpecMetadata struct {
	Project     string            `json:"project" yaml:"project"`
	Title       string            `json:"title,omitempty" yaml:"title,omitempty"`
	Description string            `json:"description,omitempty" yaml:"description,omitempty"`
	Author      string            `json:"author,omitempty" yaml:"author,omitempty"`
	CreatedAt   time.Time         `json:"created_at" yaml:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" yaml:"updated_at"`
	Status      string            `json:"status" yaml:"status"` // draft, approved, implemented
	Tags        []string          `json:"tags,omitempty" yaml:"tags,omitempty"`
	Labels      map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	SourceSpec  string            `json:"source_spec,omitempty" yaml:"source_spec,omitempty"`
}

// OpenSpecOverview provides project overview.
type OpenSpecOverview struct {
	Problem    string   `json:"problem" yaml:"problem"`
	Solution   string   `json:"solution" yaml:"solution"`
	Goals      []string `json:"goals" yaml:"goals"`
	NonGoals   []string `json:"non_goals,omitempty" yaml:"non_goals,omitempty"`
	Audience   []string `json:"audience,omitempty" yaml:"audience,omitempty"`
	Milestones []string `json:"milestones,omitempty" yaml:"milestones,omitempty"`
}

// OpenSpecFeature represents a product feature.
type OpenSpecFeature struct {
	ID           string                 `json:"id" yaml:"id"`
	Name         string                 `json:"name" yaml:"name"`
	Description  string                 `json:"description" yaml:"description"`
	Priority     string                 `json:"priority" yaml:"priority"` // must, should, could, wont
	Status       string                 `json:"status" yaml:"status"`     // proposed, approved, implemented, verified
	Category     string                 `json:"category,omitempty" yaml:"category,omitempty"`
	Requirements []OpenSpecRequirement  `json:"requirements" yaml:"requirements"`
	Dependencies []string               `json:"dependencies,omitempty" yaml:"dependencies,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// OpenSpecRequirement represents a specific requirement.
type OpenSpecRequirement struct {
	ID           string `json:"id" yaml:"id"`
	Type         string `json:"type" yaml:"type"` // functional, nonfunctional, constraint
	Description  string `json:"description" yaml:"description"`
	Rationale    string `json:"rationale,omitempty" yaml:"rationale,omitempty"`
	Priority     string `json:"priority" yaml:"priority"`
	Status       string `json:"status" yaml:"status"`
	Verification string `json:"verification,omitempty" yaml:"verification,omitempty"`
}

// OpenSpecTask represents an implementation task.
type OpenSpecTask struct {
	ID          string   `json:"id" yaml:"id"`
	Title       string   `json:"title" yaml:"title"`
	Description string   `json:"description,omitempty" yaml:"description,omitempty"`
	Type        string   `json:"type" yaml:"type"` // feature, bugfix, refactor, docs, test
	Priority    string   `json:"priority" yaml:"priority"`
	Status      string   `json:"status" yaml:"status"` // todo, in_progress, done, blocked
	Assignee    string   `json:"assignee,omitempty" yaml:"assignee,omitempty"`
	FeatureID   string   `json:"feature_id,omitempty" yaml:"feature_id,omitempty"`
	DependsOn   []string `json:"depends_on,omitempty" yaml:"depends_on,omitempty"`
	Blocks      []string `json:"blocks,omitempty" yaml:"blocks,omitempty"`
	Estimate    string   `json:"estimate,omitempty" yaml:"estimate,omitempty"`
	Labels      []string `json:"labels,omitempty" yaml:"labels,omitempty"`
}

// OpenSpecCriteria represents acceptance criteria.
type OpenSpecCriteria struct {
	ID          string `json:"id" yaml:"id"`
	FeatureID   string `json:"feature_id,omitempty" yaml:"feature_id,omitempty"`
	Description string `json:"description" yaml:"description"`
	Type        string `json:"type" yaml:"type"` // functional, performance, security, usability
	Given       string `json:"given,omitempty" yaml:"given,omitempty"`
	When        string `json:"when,omitempty" yaml:"when,omitempty"`
	Then        string `json:"then,omitempty" yaml:"then,omitempty"`
	Status      string `json:"status" yaml:"status"` // pending, passed, failed
}

// OpenSpecConstraint represents a project constraint.
type OpenSpecConstraint struct {
	ID          string `json:"id" yaml:"id"`
	Type        string `json:"type" yaml:"type"` // technical, business, regulatory, timeline
	Description string `json:"description" yaml:"description"`
	Impact      string `json:"impact,omitempty" yaml:"impact,omitempty"`
	Mitigation  string `json:"mitigation,omitempty" yaml:"mitigation,omitempty"`
}

// OpenSpecAppendix contains additional reference material.
type OpenSpecAppendix struct {
	Glossary    map[string]string   `json:"glossary,omitempty" yaml:"glossary,omitempty"`
	References  []OpenSpecReference `json:"references,omitempty" yaml:"references,omitempty"`
	Diagrams    []string            `json:"diagrams,omitempty" yaml:"diagrams,omitempty"`
	RawSections map[string]string   `json:"raw_sections,omitempty" yaml:"raw_sections,omitempty"`
}

// OpenSpecReference represents a reference document.
type OpenSpecReference struct {
	Title string `json:"title" yaml:"title"`
	URL   string `json:"url,omitempty" yaml:"url,omitempty"`
	Type  string `json:"type,omitempty" yaml:"type,omitempty"`
}

// Export exports the spec to OpenSpec format.
func (t *OpenSpecTarget) Export(spec string, config ExportConfig) (*ExportResult, error) {
	if err := t.Validate(spec); err != nil {
		return nil, err
	}

	// Determine output format
	format := "json"
	if f, ok := config.Options["format"].(string); ok && f != "" {
		format = strings.ToLower(f)
	}

	// Determine output directory
	outputDir := config.OutputDir
	if outputDir == "" {
		outputDir = "openspec"
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("creating output directory: %w", err)
	}

	// Parse spec into OpenSpec document
	doc := t.parseSpec(spec, config.ProjectName)

	var files []string

	// Write main document
	var data []byte
	var ext string
	var err error

	switch format {
	case "yaml", "yml":
		data, err = yaml.Marshal(doc)
		ext = "yaml"
	default:
		data, err = json.MarshalIndent(doc, "", "  ")
		ext = "json"
	}
	if err != nil {
		return nil, fmt.Errorf("marshaling document: %w", err)
	}

	mainFile := filepath.Join(outputDir, fmt.Sprintf("openspec.%s", ext))
	if err := os.WriteFile(mainFile, data, 0600); err != nil {
		return nil, fmt.Errorf("writing openspec file: %w", err)
	}
	files = append(files, mainFile)

	// Also write separate files for features and tasks if requested
	if separate, ok := config.Options["separate_files"].(bool); ok && separate {
		separateFiles, err := t.writeSeparateFiles(doc, outputDir, ext)
		if err != nil {
			return nil, err
		}
		files = append(files, separateFiles...)
	}

	return &ExportResult{
		Target:    t.Name(),
		OutputDir: outputDir,
		Files:     files,
		Success:   true,
		Message:   fmt.Sprintf("Exported to %s format in %s", format, outputDir),
	}, nil
}

// writeSeparateFiles writes individual feature and task files.
func (t *OpenSpecTarget) writeSeparateFiles(doc OpenSpecDocument, outputDir, ext string) ([]string, error) {
	var files []string

	// Write features
	if len(doc.Features) > 0 {
		featuresDir := filepath.Join(outputDir, "features")
		if err := os.MkdirAll(featuresDir, 0755); err != nil {
			return nil, fmt.Errorf("creating features directory: %w", err)
		}

		for _, feature := range doc.Features {
			var data []byte
			var err error
			if ext == "yaml" {
				data, err = yaml.Marshal(feature)
			} else {
				data, err = json.MarshalIndent(feature, "", "  ")
			}
			if err != nil {
				continue
			}

			filename := filepath.Join(featuresDir, fmt.Sprintf("%s.%s", feature.ID, ext))
			if err := os.WriteFile(filename, data, 0600); err != nil {
				continue
			}
			files = append(files, filename)
		}
	}

	// Write tasks
	if len(doc.Tasks) > 0 {
		tasksDir := filepath.Join(outputDir, "tasks")
		if err := os.MkdirAll(tasksDir, 0755); err != nil {
			return nil, fmt.Errorf("creating tasks directory: %w", err)
		}

		for _, task := range doc.Tasks {
			var data []byte
			var err error
			if ext == "yaml" {
				data, err = yaml.Marshal(task)
			} else {
				data, err = json.MarshalIndent(task, "", "  ")
			}
			if err != nil {
				continue
			}

			filename := filepath.Join(tasksDir, fmt.Sprintf("%s.%s", task.ID, ext))
			if err := os.WriteFile(filename, data, 0600); err != nil {
				continue
			}
			files = append(files, filename)
		}
	}

	return files, nil
}

// parseSpec parses markdown spec into OpenSpec document.
func (t *OpenSpecTarget) parseSpec(spec, projectName string) OpenSpecDocument {
	now := time.Now()

	doc := OpenSpecDocument{
		Version: "1.0.0",
		Metadata: OpenSpecMetadata{
			Project:    projectName,
			CreatedAt:  now,
			UpdatedAt:  now,
			Status:     "approved",
			SourceSpec: "spec.md",
		},
		Overview: OpenSpecOverview{
			Problem:  "See source spec for details",
			Solution: "See source spec for details",
		},
	}

	// Parse sections
	lines := strings.Split(spec, "\n")
	currentSection := ""
	var currentContent strings.Builder
	sectionContents := make(map[string]string)

	for _, line := range lines {
		if strings.HasPrefix(line, "# ") {
			// Save previous section
			if currentSection != "" {
				sectionContents[currentSection] = currentContent.String()
			}
			currentSection = strings.TrimSpace(strings.TrimPrefix(line, "# "))
			currentContent.Reset()
		} else if strings.HasPrefix(line, "## ") {
			// Save previous section
			if currentSection != "" {
				sectionContents[currentSection] = currentContent.String()
			}
			currentSection = strings.TrimSpace(strings.TrimPrefix(line, "## "))
			currentContent.Reset()
		} else {
			currentContent.WriteString(line)
			currentContent.WriteString("\n")
		}
	}
	// Save last section
	if currentSection != "" {
		sectionContents[currentSection] = currentContent.String()
	}

	// Extract title from first H1
	for _, line := range lines {
		if strings.HasPrefix(line, "# ") {
			doc.Metadata.Title = strings.TrimSpace(strings.TrimPrefix(line, "# "))
			break
		}
	}

	// Extract overview
	for section, content := range sectionContents {
		sectionLower := strings.ToLower(section)
		if strings.Contains(sectionLower, "problem") || strings.Contains(sectionLower, "background") {
			doc.Overview.Problem = strings.TrimSpace(content)
		} else if strings.Contains(sectionLower, "solution") || strings.Contains(sectionLower, "approach") {
			doc.Overview.Solution = strings.TrimSpace(content)
		} else if strings.Contains(sectionLower, "goal") {
			doc.Overview.Goals = t.extractBulletPoints(content)
		} else if strings.Contains(sectionLower, "non-goal") || strings.Contains(sectionLower, "out of scope") {
			doc.Overview.NonGoals = t.extractBulletPoints(content)
		} else if strings.Contains(sectionLower, "audience") || strings.Contains(sectionLower, "user") {
			doc.Overview.Audience = t.extractBulletPoints(content)
		}
	}

	// Extract features
	doc.Features = t.extractFeatures(spec)

	// Extract tasks
	doc.Tasks = t.extractTasks(spec)

	// Extract acceptance criteria
	doc.Acceptance = t.extractAcceptanceCriteria(spec)

	// Store raw sections in appendix
	doc.Appendix = &OpenSpecAppendix{
		RawSections: sectionContents,
	}

	return doc
}

// extractBulletPoints extracts bullet points from content.
func (t *OpenSpecTarget) extractBulletPoints(content string) []string {
	var points []string
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "- ") {
			points = append(points, strings.TrimPrefix(line, "- "))
		} else if strings.HasPrefix(line, "* ") {
			points = append(points, strings.TrimPrefix(line, "* "))
		} else if matched, _ := regexp.MatchString(`^\d+\.`, line); matched {
			// Numbered list
			parts := strings.SplitN(line, ".", 2)
			if len(parts) > 1 {
				points = append(points, strings.TrimSpace(parts[1]))
			}
		}
	}
	return points
}

// extractFeatures extracts features from spec.
func (t *OpenSpecTarget) extractFeatures(spec string) []OpenSpecFeature {
	var features []OpenSpecFeature
	featureNum := 0

	lines := strings.Split(spec, "\n")
	inFeatureSection := false
	var currentFeature *OpenSpecFeature

	for i, line := range lines {
		lineLower := strings.ToLower(line)

		// Detect feature sections
		if strings.HasPrefix(line, "## ") || strings.HasPrefix(line, "### ") {
			title := strings.TrimSpace(strings.TrimLeft(line, "# "))

			// Check if it's a feature-like section
			if strings.Contains(lineLower, "feature") ||
				strings.Contains(lineLower, "requirement") ||
				strings.Contains(lineLower, "capability") ||
				strings.Contains(lineLower, "function") {
				inFeatureSection = true
			}

			if inFeatureSection && !strings.HasPrefix(line, "## Features") {
				// Save previous feature
				if currentFeature != nil {
					features = append(features, *currentFeature)
				}

				featureNum++
				currentFeature = &OpenSpecFeature{
					ID:          fmt.Sprintf("F-%03d", featureNum),
					Name:        title,
					Description: "",
					Priority:    "should",
					Status:      "proposed",
				}
			}
		}

		// Extract requirements from bullet points under features
		if currentFeature != nil {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
				reqText := strings.TrimPrefix(strings.TrimPrefix(line, "- "), "* ")
				priority := "should"
				reqType := "functional"

				// Detect priority keywords
				if strings.Contains(strings.ToUpper(reqText), "MUST") || strings.Contains(strings.ToUpper(reqText), "REQUIRED") {
					priority = "must"
				} else if strings.Contains(strings.ToUpper(reqText), "COULD") || strings.Contains(strings.ToUpper(reqText), "OPTIONAL") {
					priority = "could"
				}

				// Detect requirement type
				if strings.Contains(strings.ToLower(reqText), "performance") || strings.Contains(strings.ToLower(reqText), "latency") {
					reqType = "nonfunctional"
				} else if strings.Contains(strings.ToLower(reqText), "security") || strings.Contains(strings.ToLower(reqText), "auth") {
					reqType = "nonfunctional"
				}

				currentFeature.Requirements = append(currentFeature.Requirements, OpenSpecRequirement{
					ID:          fmt.Sprintf("%s-R%03d", currentFeature.ID, len(currentFeature.Requirements)+1),
					Type:        reqType,
					Description: reqText,
					Priority:    priority,
					Status:      "proposed",
				})
			} else if line != "" && !strings.HasPrefix(line, "#") {
				// Add to feature description
				if currentFeature.Description != "" {
					currentFeature.Description += " "
				}
				currentFeature.Description += line
				_ = i // use i to avoid unused variable warning
			}
		}
	}

	// Save last feature
	if currentFeature != nil {
		features = append(features, *currentFeature)
	}

	return features
}

// extractTasks extracts tasks from spec.
func (t *OpenSpecTarget) extractTasks(spec string) []OpenSpecTask {
	var tasks []OpenSpecTask
	taskNum := 0

	lines := strings.Split(spec, "\n")
	inTaskSection := false

	for _, line := range lines {
		lineLower := strings.ToLower(line)

		// Detect task sections
		if strings.HasPrefix(line, "## ") {
			inTaskSection = strings.Contains(lineLower, "task") ||
				strings.Contains(lineLower, "implementation") ||
				strings.Contains(lineLower, "work") ||
				strings.Contains(lineLower, "action")
		}

		// Extract tasks from checklists and bullet points
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "- [ ]") {
			taskNum++
			tasks = append(tasks, OpenSpecTask{
				ID:       fmt.Sprintf("T-%03d", taskNum),
				Title:    strings.TrimSpace(strings.TrimPrefix(line, "- [ ]")),
				Type:     "feature",
				Priority: "medium",
				Status:   "todo",
			})
		} else if strings.HasPrefix(line, "- [x]") || strings.HasPrefix(line, "- [X]") {
			taskNum++
			title := strings.TrimSpace(strings.TrimPrefix(line, "- [x]"))
			title = strings.TrimSpace(strings.TrimPrefix(title, "- [X]"))
			tasks = append(tasks, OpenSpecTask{
				ID:       fmt.Sprintf("T-%03d", taskNum),
				Title:    title,
				Type:     "feature",
				Priority: "medium",
				Status:   "done",
			})
		} else if inTaskSection && strings.HasPrefix(line, "- ") {
			taskNum++
			tasks = append(tasks, OpenSpecTask{
				ID:       fmt.Sprintf("T-%03d", taskNum),
				Title:    strings.TrimSpace(strings.TrimPrefix(line, "- ")),
				Type:     "feature",
				Priority: "medium",
				Status:   "todo",
			})
		}
	}

	return tasks
}

// extractAcceptanceCriteria extracts acceptance criteria from spec.
func (t *OpenSpecTarget) extractAcceptanceCriteria(spec string) []OpenSpecCriteria {
	var criteria []OpenSpecCriteria
	criteriaNum := 0

	lines := strings.Split(spec, "\n")
	inAcceptanceSection := false
	var currentCriteria *OpenSpecCriteria

	for _, line := range lines {
		lineLower := strings.ToLower(line)

		// Detect acceptance criteria sections
		if strings.HasPrefix(line, "## ") || strings.HasPrefix(line, "### ") {
			inAcceptanceSection = strings.Contains(lineLower, "acceptance") ||
				strings.Contains(lineLower, "criteria") ||
				strings.Contains(lineLower, "verification") ||
				strings.Contains(lineLower, "test")

			// Save previous criteria
			if currentCriteria != nil {
				criteria = append(criteria, *currentCriteria)
				currentCriteria = nil
			}
		}

		if !inAcceptanceSection {
			continue
		}

		line = strings.TrimSpace(line)

		// Parse Given/When/Then format
		if strings.HasPrefix(strings.ToLower(line), "given") {
			criteriaNum++
			currentCriteria = &OpenSpecCriteria{
				ID:     fmt.Sprintf("AC-%03d", criteriaNum),
				Type:   "functional",
				Status: "pending",
				Given:  strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(line, "Given"), "given")),
			}
		} else if currentCriteria != nil && strings.HasPrefix(strings.ToLower(line), "when") {
			currentCriteria.When = strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(line, "When"), "when"))
		} else if currentCriteria != nil && strings.HasPrefix(strings.ToLower(line), "then") {
			currentCriteria.Then = strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(line, "Then"), "then"))
			currentCriteria.Description = fmt.Sprintf("Given %s, when %s, then %s",
				currentCriteria.Given, currentCriteria.When, currentCriteria.Then)
		}

		// Also extract bullet point criteria
		if strings.HasPrefix(line, "- ") && currentCriteria == nil {
			criteriaNum++
			criteria = append(criteria, OpenSpecCriteria{
				ID:          fmt.Sprintf("AC-%03d", criteriaNum),
				Description: strings.TrimSpace(strings.TrimPrefix(line, "- ")),
				Type:        "functional",
				Status:      "pending",
			})
		}
	}

	// Save last criteria
	if currentCriteria != nil {
		criteria = append(criteria, *currentCriteria)
	}

	return criteria
}
