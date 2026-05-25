package target

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func init() {
	Register(&GasCityTarget{})
}

// GasCityTarget exports to GasCity multi-agent format.
// GasCity orchestrates multiple agents with city.toml, agent definitions, and orders.
type GasCityTarget struct{}

// Name returns the target name.
func (t *GasCityTarget) Name() string {
	return "gascity"
}

// Description returns a description of the target.
func (t *GasCityTarget) Description() string {
	return "GasCity multi-agent orchestration (city.toml, agents, orders)"
}

// Capabilities returns what this target supports.
func (t *GasCityTarget) Capabilities() Capabilities {
	return Capabilities{
		SequentialTasks:   true,
		ParallelExecution: true,
		MultiAgent:        true,
		Verification:      true,
		DependencyGraph:   true,
	}
}

// Validate checks if the spec can be exported to this target.
func (t *GasCityTarget) Validate(spec string) error {
	if spec == "" {
		return fmt.Errorf("spec content is empty")
	}
	return nil
}

// Export exports the spec to GasCity format.
func (t *GasCityTarget) Export(spec string, config ExportConfig) (*ExportResult, error) {
	if err := t.Validate(spec); err != nil {
		return nil, err
	}

	// Determine output directory
	outputDir := config.OutputDir
	if outputDir == "" {
		outputDir = ".gascity"
	}

	// Create directory structure
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("creating output directory: %w", err)
	}

	agentsDir := filepath.Join(outputDir, "agents")
	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		return nil, fmt.Errorf("creating agents directory: %w", err)
	}

	ordersDir := filepath.Join(outputDir, "orders")
	if err := os.MkdirAll(ordersDir, 0755); err != nil {
		return nil, fmt.Errorf("creating orders directory: %w", err)
	}

	var files []string

	// Extract agents and orders from spec
	agents := t.extractAgents(spec)
	orders := t.extractOrders(spec, agents)

	// Generate city.toml
	city := t.generateCity(config.ProjectName, agents, orders)
	cityPath := filepath.Join(outputDir, "city.toml")
	if err := os.WriteFile(cityPath, []byte(city), 0600); err != nil {
		return nil, fmt.Errorf("writing city.toml: %w", err)
	}
	files = append(files, cityPath)

	// Generate agent files
	for _, agent := range agents {
		agentContent := t.generateAgent(agent)
		agentPath := filepath.Join(agentsDir, agent.ID+".toml")
		if err := os.WriteFile(agentPath, []byte(agentContent), 0600); err != nil {
			return nil, fmt.Errorf("writing agent %s: %w", agent.ID, err)
		}
		files = append(files, agentPath)
	}

	// Generate order files
	for _, order := range orders {
		orderContent := t.generateOrder(order)
		orderPath := filepath.Join(ordersDir, order.ID+".toml")
		if err := os.WriteFile(orderPath, []byte(orderContent), 0600); err != nil {
			return nil, fmt.Errorf("writing order %s: %w", order.ID, err)
		}
		files = append(files, orderPath)
	}

	return &ExportResult{
		Target:    t.Name(),
		OutputDir: outputDir,
		Files:     files,
		Success:   true,
		Message:   fmt.Sprintf("Exported to %s with %d agents and %d orders", outputDir, len(agents), len(orders)),
	}, nil
}

// Agent represents a GasCity agent.
type Agent struct {
	ID           string
	Name         string
	Role         string
	Capabilities []string
	Priority     int
}

// Order represents a GasCity order (task assignment).
type Order struct {
	ID          string
	Description string
	AssignedTo  string   // Agent ID
	DependsOn   []string // Order IDs
	Status      string   // pending, in_progress, complete
	Priority    int
}

// extractAgents determines agents from the spec content.
func (t *GasCityTarget) extractAgents(spec string) []Agent {
	agents := []Agent{}
	agentMap := make(map[string]bool)

	// Look for role-based keywords to determine needed agents
	specLower := strings.ToLower(spec)

	// Common agent types
	agentTypes := []struct {
		keywords []string
		agent    Agent
	}{
		{
			keywords: []string{"implement", "code", "develop", "build", "create"},
			agent:    Agent{ID: "dev", Name: "Developer", Role: "development", Capabilities: []string{"code", "test", "debug"}, Priority: 1},
		},
		{
			keywords: []string{"test", "qa", "quality", "verify", "validate"},
			agent:    Agent{ID: "qa", Name: "QA Engineer", Role: "testing", Capabilities: []string{"test", "report", "verify"}, Priority: 2},
		},
		{
			keywords: []string{"deploy", "release", "infrastructure", "ops", "devops"},
			agent:    Agent{ID: "ops", Name: "Operations", Role: "operations", Capabilities: []string{"deploy", "monitor", "scale"}, Priority: 3},
		},
		{
			keywords: []string{"document", "documentation", "readme", "api doc"},
			agent:    Agent{ID: "docs", Name: "Documentation", Role: "documentation", Capabilities: []string{"write", "review", "publish"}, Priority: 4},
		},
		{
			keywords: []string{"review", "approve", "sign-off", "audit"},
			agent:    Agent{ID: "reviewer", Name: "Reviewer", Role: "review", Capabilities: []string{"review", "approve", "feedback"}, Priority: 5},
		},
		{
			keywords: []string{"design", "architecture", "system design", "api design"},
			agent:    Agent{ID: "architect", Name: "Architect", Role: "design", Capabilities: []string{"design", "review", "guide"}, Priority: 0},
		},
	}

	for _, at := range agentTypes {
		for _, keyword := range at.keywords {
			if strings.Contains(specLower, keyword) && !agentMap[at.agent.ID] {
				agents = append(agents, at.agent)
				agentMap[at.agent.ID] = true
				break
			}
		}
	}

	// Ensure at least one agent
	if len(agents) == 0 {
		agents = append(agents, Agent{
			ID:           "default",
			Name:         "Default Agent",
			Role:         "general",
			Capabilities: []string{"execute", "report"},
			Priority:     1,
		})
	}

	return agents
}

// extractOrders extracts orders from the spec.
func (t *GasCityTarget) extractOrders(spec string, agents []Agent) []Order {
	var orders []Order

	// Look for task lists
	taskRE := regexp.MustCompile(`(?m)^[-*]\s*\[[ x]\]\s*(.+)$`)
	matches := taskRE.FindAllStringSubmatch(spec, -1)

	for i, match := range matches {
		if len(match) > 1 {
			desc := strings.TrimSpace(match[1])
			status := "pending"
			if strings.Contains(match[0], "[x]") || strings.Contains(match[0], "[X]") {
				status = "complete"
			}

			// Assign to appropriate agent based on keywords
			assignedTo := t.matchAgentForTask(desc, agents)

			order := Order{
				ID:          fmt.Sprintf("order-%03d", i+1),
				Description: desc,
				AssignedTo:  assignedTo,
				Status:      status,
				Priority:    i + 1,
			}

			// Simple linear dependency
			if i > 0 {
				order.DependsOn = []string{fmt.Sprintf("order-%03d", i)}
			}

			orders = append(orders, order)
		}
	}

	// Look for numbered tasks
	numberedRE := regexp.MustCompile(`(?m)^\d+\.\s+(.+)$`)
	numberedMatches := numberedRE.FindAllStringSubmatch(spec, -1)
	startIdx := len(orders)
	for i, match := range numberedMatches {
		if len(match) > 1 {
			desc := strings.TrimSpace(match[1])
			if strings.HasPrefix(desc, "#") || len(desc) < 5 {
				continue
			}

			assignedTo := t.matchAgentForTask(desc, agents)

			order := Order{
				ID:          fmt.Sprintf("order-%03d", startIdx+i+1),
				Description: desc,
				AssignedTo:  assignedTo,
				Status:      "pending",
				Priority:    startIdx + i + 1,
			}

			if len(orders) > 0 {
				order.DependsOn = []string{orders[len(orders)-1].ID}
			}

			orders = append(orders, order)
		}
	}

	// Default orders if none found
	if len(orders) == 0 {
		devAgent := "default"
		qaAgent := "default"
		opsAgent := "default"
		for _, a := range agents {
			if a.Role == "development" {
				devAgent = a.ID
			} else if a.Role == "testing" {
				qaAgent = a.ID
			} else if a.Role == "operations" {
				opsAgent = a.ID
			}
		}

		orders = []Order{
			{ID: "order-001", Description: "Review specification", AssignedTo: devAgent, Status: "pending", Priority: 1},
			{ID: "order-002", Description: "Implement core functionality", AssignedTo: devAgent, Status: "pending", Priority: 2, DependsOn: []string{"order-001"}},
			{ID: "order-003", Description: "Write tests", AssignedTo: qaAgent, Status: "pending", Priority: 3, DependsOn: []string{"order-002"}},
			{ID: "order-004", Description: "Deploy", AssignedTo: opsAgent, Status: "pending", Priority: 4, DependsOn: []string{"order-003"}},
		}
	}

	return orders
}

// matchAgentForTask finds the best agent for a task.
func (t *GasCityTarget) matchAgentForTask(desc string, agents []Agent) string {
	descLower := strings.ToLower(desc)

	// Keywords to agent role mapping
	roleKeywords := map[string][]string{
		"development":   {"implement", "code", "build", "create", "develop", "fix", "refactor"},
		"testing":       {"test", "qa", "verify", "validate", "check"},
		"operations":    {"deploy", "release", "infrastructure", "monitor", "scale"},
		"documentation": {"document", "write doc", "readme", "api doc"},
		"review":        {"review", "approve", "sign-off"},
		"design":        {"design", "architect", "plan"},
	}

	for _, agent := range agents {
		if keywords, ok := roleKeywords[agent.Role]; ok {
			for _, keyword := range keywords {
				if strings.Contains(descLower, keyword) {
					return agent.ID
				}
			}
		}
	}

	// Default to first agent
	if len(agents) > 0 {
		return agents[0].ID
	}
	return "default"
}

// generateCity creates the city.toml file.
func (t *GasCityTarget) generateCity(projectName string, agents []Agent, orders []Order) string {
	var sb strings.Builder

	sb.WriteString("# GasCity Configuration\n")
	sb.WriteString(fmt.Sprintf("# Generated by visionspec on %s\n\n", time.Now().Format("2006-01-02")))

	sb.WriteString("[city]\n")
	sb.WriteString(fmt.Sprintf("name = %q\n", projectName))
	sb.WriteString("version = \"1.0\"\n")
	sb.WriteString(fmt.Sprintf("created = %q\n", time.Now().Format(time.RFC3339)))
	sb.WriteString("source = \"multispec\"\n\n")

	sb.WriteString("[city.coordination]\n")
	sb.WriteString("mode = \"orchestrated\"  # orchestrated, autonomous, hybrid\n")
	sb.WriteString("checkpoint_interval = \"1h\"\n")
	sb.WriteString("timeout = \"24h\"\n\n")

	// Agents section
	sb.WriteString("[agents]\n")
	sb.WriteString("directory = \"agents\"\n")
	sb.WriteString(fmt.Sprintf("count = %d\n\n", len(agents)))

	for _, agent := range agents {
		sb.WriteString(fmt.Sprintf("[[agents.list]]\n"))
		sb.WriteString(fmt.Sprintf("id = %q\n", agent.ID))
		sb.WriteString(fmt.Sprintf("role = %q\n", agent.Role))
		sb.WriteString(fmt.Sprintf("priority = %d\n\n", agent.Priority))
	}

	// Orders section
	sb.WriteString("[orders]\n")
	sb.WriteString("directory = \"orders\"\n")
	sb.WriteString(fmt.Sprintf("count = %d\n\n", len(orders)))

	// Pending orders
	pending := []string{}
	for _, order := range orders {
		if order.Status == "pending" && len(order.DependsOn) == 0 {
			pending = append(pending, order.ID)
		}
	}
	sb.WriteString(fmt.Sprintf("entry_points = [%s]\n\n", quoteStrings(pending)))

	// Execution rules
	sb.WriteString("[execution]\n")
	sb.WriteString("parallel_agents = true\n")
	sb.WriteString("order_timeout = \"2h\"\n")
	sb.WriteString("on_failure = \"pause\"  # pause, retry, abort\n")
	sb.WriteString("max_retries = 3\n")

	return sb.String()
}

// generateAgent creates an agent TOML file.
func (t *GasCityTarget) generateAgent(agent Agent) string {
	var sb strings.Builder

	sb.WriteString("# GasCity Agent Definition\n\n")

	sb.WriteString("[agent]\n")
	sb.WriteString(fmt.Sprintf("id = %q\n", agent.ID))
	sb.WriteString(fmt.Sprintf("name = %q\n", agent.Name))
	sb.WriteString(fmt.Sprintf("role = %q\n", agent.Role))
	sb.WriteString(fmt.Sprintf("priority = %d\n\n", agent.Priority))

	sb.WriteString("[agent.capabilities]\n")
	sb.WriteString(fmt.Sprintf("skills = [%s]\n", quoteStrings(agent.Capabilities)))
	sb.WriteString("max_concurrent_orders = 3\n\n")

	sb.WriteString("[agent.communication]\n")
	sb.WriteString("protocol = \"message\"  # message, event, rpc\n")
	sb.WriteString("timeout = \"5m\"\n\n")

	sb.WriteString("[agent.resources]\n")
	sb.WriteString("memory = \"standard\"\n")
	sb.WriteString("compute = \"standard\"\n")

	return sb.String()
}

// generateOrder creates an order TOML file.
func (t *GasCityTarget) generateOrder(order Order) string {
	var sb strings.Builder

	sb.WriteString("# GasCity Order\n\n")

	sb.WriteString("[order]\n")
	sb.WriteString(fmt.Sprintf("id = %q\n", order.ID))
	sb.WriteString(fmt.Sprintf("description = %q\n", order.Description))
	sb.WriteString(fmt.Sprintf("assigned_to = %q\n", order.AssignedTo))
	sb.WriteString(fmt.Sprintf("status = %q\n", order.Status))
	sb.WriteString(fmt.Sprintf("priority = %d\n\n", order.Priority))

	if len(order.DependsOn) > 0 {
		sb.WriteString("[order.dependencies]\n")
		sb.WriteString(fmt.Sprintf("depends_on = [%s]\n", quoteStrings(order.DependsOn)))
		sb.WriteString("wait_strategy = \"all\"  # all, any\n\n")
	}

	sb.WriteString("[order.execution]\n")
	sb.WriteString("timeout = \"2h\"\n")
	sb.WriteString("retries = 3\n")
	sb.WriteString("on_complete = \"notify\"  # notify, chain, done\n")

	return sb.String()
}
