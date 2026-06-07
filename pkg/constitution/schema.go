// Package constitution defines the hierarchical constitution system for VisionSpec.
// Constitutions define organizational defaults that flow down: org → team → project.
package constitution

import "time"

// Level represents the hierarchy level of a constitution.
type Level string

const (
	LevelOrganization Level = "organization"
	LevelTeam         Level = "team"
	LevelProject      Level = "project"
)

// Constitution represents a configuration document that defines defaults
// for VisionSpec artifacts at a given hierarchy level.
type Constitution struct {
	APIVersion     string         `json:"apiVersion" yaml:"apiVersion"` // visionspec/v1
	Kind           string         `json:"kind" yaml:"kind"`             // Constitution
	Metadata       Metadata       `json:"metadata" yaml:"metadata"`
	Technical      Technical      `json:"technical,omitempty" yaml:"technical,omitempty"`
	Infrastructure Infrastructure `json:"infrastructure,omitempty" yaml:"infrastructure,omitempty"`
	Security       Security       `json:"security,omitempty" yaml:"security,omitempty"`
	Prompts        Prompts        `json:"prompts,omitempty" yaml:"prompts,omitempty"`
}

// Metadata contains constitution identification and hierarchy info.
type Metadata struct {
	Name        string    `json:"name" yaml:"name"`
	Level       Level     `json:"level" yaml:"level"`
	Inherits    string    `json:"inherits,omitempty" yaml:"inherits,omitempty"` // e.g., "org/plexusone"
	Version     string    `json:"version,omitempty" yaml:"version,omitempty"`
	Description string    `json:"description,omitempty" yaml:"description,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty" yaml:"updatedAt,omitempty"`
}

// Technical defines TRD-scoped defaults.
type Technical struct {
	Languages     Languages         `json:"languages,omitempty" yaml:"languages,omitempty"`
	APIs          APIs              `json:"apis,omitempty" yaml:"apis,omitempty"`
	Database      Database          `json:"database,omitempty" yaml:"database,omitempty"`
	Tenancy       Tenancy           `json:"tenancy,omitempty" yaml:"tenancy,omitempty"`
	Observability TechObservability `json:"observability,omitempty" yaml:"observability,omitempty"`
}

// Languages defines programming language choices.
type Languages struct {
	Backend  LanguageChoice `json:"backend,omitempty" yaml:"backend,omitempty"`
	Frontend LanguageChoice `json:"frontend,omitempty" yaml:"frontend,omitempty"`
	WASM     string         `json:"wasm,omitempty" yaml:"wasm,omitempty"` // e.g., "rust"
}

// LanguageChoice defines a language selection with allowed alternatives.
type LanguageChoice struct {
	Primary           string   `json:"primary" yaml:"primary"`                                         // e.g., "go"
	Allowed           []string `json:"allowed,omitempty" yaml:"allowed,omitempty"`                     // e.g., ["go", "rust"]
	ExceptionsRequire string   `json:"exceptionsRequire,omitempty" yaml:"exceptionsRequire,omitempty"` // e.g., "approval"
}

// APIs defines API style and framework choices.
type APIs struct {
	REST RESTConfig `json:"rest,omitempty" yaml:"rest,omitempty"`
	GRPC GRPCConfig `json:"grpc,omitempty" yaml:"grpc,omitempty"`
}

// RESTConfig defines REST API configuration.
type RESTConfig struct {
	Framework  string `json:"framework,omitempty" yaml:"framework,omitempty"`   // e.g., "huma-chi"
	SpecFormat string `json:"specFormat,omitempty" yaml:"specFormat,omitempty"` // e.g., "openapi-3.1"
	StyleGuide string `json:"styleGuide,omitempty" yaml:"styleGuide,omitempty"` // e.g., "google-api-design-guide"
}

// GRPCConfig defines gRPC configuration.
type GRPCConfig struct {
	Framework string `json:"framework,omitempty" yaml:"framework,omitempty"` // e.g., "connect-go"
}

// Database defines database choices.
type Database struct {
	Relational   string `json:"relational,omitempty" yaml:"relational,omitempty"`     // e.g., "postgresql"
	Document     string `json:"document,omitempty" yaml:"document,omitempty"`         // e.g., "mongodb"
	Cache        string `json:"cache,omitempty" yaml:"cache,omitempty"`               // e.g., "redis"
	Search       string `json:"search,omitempty" yaml:"search,omitempty"`             // e.g., "elasticsearch"
	MultiTenancy string `json:"multiTenancy,omitempty" yaml:"multiTenancy,omitempty"` // e.g., "rls", "schema", "database"
}

// Tenancy defines multi-tenancy configuration.
type Tenancy struct {
	Model          TenancyModel `json:"model" yaml:"model"`
	Isolation      string       `json:"isolation,omitempty" yaml:"isolation,omitempty"`           // e.g., "rls", "schema", "database"
	TenantIDHeader string       `json:"tenantIdHeader,omitempty" yaml:"tenantIdHeader,omitempty"` // e.g., "X-Tenant-ID"
	TenantIDClaim  string       `json:"tenantIdClaim,omitempty" yaml:"tenantIdClaim,omitempty"`   // e.g., "tenant_id" (JWT claim)
}

// TenancyModel represents the tenancy architecture.
type TenancyModel string

const (
	TenancySingleTenant TenancyModel = "single-tenant"
	TenancyMultiTenant  TenancyModel = "multi-tenant"
)

// TechObservability defines instrumentation choices (TRD scope).
type TechObservability struct {
	Library        string `json:"library,omitempty" yaml:"library,omitempty"` // e.g., "github.com/plexusone/omniobserve"
	AutoInstrument bool   `json:"autoInstrument,omitempty" yaml:"autoInstrument,omitempty"`
	MetricsSDK     string `json:"metricsSDK,omitempty" yaml:"metricsSDK,omitempty"` // e.g., "prometheus/client_golang"
	TracesSDK      string `json:"tracesSDK,omitempty" yaml:"tracesSDK,omitempty"`   // e.g., "opentelemetry-go"
	LoggingSDK     string `json:"loggingSDK,omitempty" yaml:"loggingSDK,omitempty"` // e.g., "slog"
}

// Infrastructure defines IRD-scoped defaults.
type Infrastructure struct {
	IaC           IaC                `json:"iac,omitempty" yaml:"iac,omitempty"`
	Observability InfraObservability `json:"observability,omitempty" yaml:"observability,omitempty"`
	LocalDev      LocalDev           `json:"localDev,omitempty" yaml:"localDev,omitempty"`
	Cloud         Cloud              `json:"cloud,omitempty" yaml:"cloud,omitempty"`
	Availability  Availability       `json:"availability,omitempty" yaml:"availability,omitempty"`
}

// IaC defines Infrastructure as Code choices.
type IaC struct {
	Tool     IaCTool `json:"tool" yaml:"tool"`
	Language string  `json:"language,omitempty" yaml:"language,omitempty"` // e.g., "go", "typescript"
	RepoPath string  `json:"repoPath,omitempty" yaml:"repoPath,omitempty"` // e.g., "infra/"
}

// IaCTool represents the IaC tool choice.
type IaCTool string

const (
	IaCPulumi         IaCTool = "pulumi"
	IaCCDK            IaCTool = "cdk"
	IaCTerraform      IaCTool = "terraform"
	IaCCloudFormation IaCTool = "cloudformation"
	IaCNone           IaCTool = "none"
)

// InfraObservability defines observability infrastructure (IRD scope).
type InfraObservability struct {
	Metrics ObservabilityPillar `json:"metrics,omitempty" yaml:"metrics,omitempty"`
	Traces  ObservabilityPillar `json:"traces,omitempty" yaml:"traces,omitempty"`
	Logging ObservabilityPillar `json:"logging,omitempty" yaml:"logging,omitempty"`
}

// ObservabilityPillar defines a single observability pillar configuration.
type ObservabilityPillar struct {
	Platform      string `json:"platform,omitempty" yaml:"platform,omitempty"`           // e.g., "prometheus", "loki"
	Visualization string `json:"visualization,omitempty" yaml:"visualization,omitempty"` // e.g., "grafana"
	Backend       string `json:"backend,omitempty" yaml:"backend,omitempty"`             // e.g., "langfuse" (for traces)
	Collector     string `json:"collector,omitempty" yaml:"collector,omitempty"`         // e.g., "opentelemetry-collector"
}

// LocalDev defines local development environment preferences.
type LocalDev struct {
	Priority []LocalDevTarget `json:"priority,omitempty" yaml:"priority,omitempty"`
}

// LocalDevTarget represents a local development target.
type LocalDevTarget string

const (
	LocalDevBinaries   LocalDevTarget = "binaries"
	LocalDevPodman     LocalDevTarget = "podman"
	LocalDevDocker     LocalDevTarget = "docker"
	LocalDevLocalStack LocalDevTarget = "localstack"
	LocalDevMinikube   LocalDevTarget = "minikube"
	LocalDevKind       LocalDevTarget = "kind"
)

// Cloud defines cloud provider configuration.
type Cloud struct {
	Provider string   `json:"provider,omitempty" yaml:"provider,omitempty"` // e.g., "aws", "gcp", "azure"
	Regions  []string `json:"regions,omitempty" yaml:"regions,omitempty"`
}

// Availability defines availability and reliability targets.
type Availability struct {
	Target      AvailabilityTarget `json:"target" yaml:"target"`
	RTO         string             `json:"rto,omitempty" yaml:"rto,omitempty"` // e.g., "1h", "15m"
	RPO         string             `json:"rpo,omitempty" yaml:"rpo,omitempty"` // e.g., "15m", "1h"
	MultiRegion bool               `json:"multiRegion,omitempty" yaml:"multiRegion,omitempty"`
	MultiAZ     bool               `json:"multiAZ,omitempty" yaml:"multiAZ,omitempty"`
	DRStrategy  string             `json:"drStrategy,omitempty" yaml:"drStrategy,omitempty"` // e.g., "active-passive", "active-active"
}

// AvailabilityTarget represents SLA availability percentages.
type AvailabilityTarget string

const (
	Availability99    AvailabilityTarget = "99"     // 99% - 3.65 days/year downtime
	Availability999   AvailabilityTarget = "99.9"   // 99.9% - 8.76 hours/year downtime
	Availability9999  AvailabilityTarget = "99.99"  // 99.99% - 52.6 minutes/year downtime
	Availability99999 AvailabilityTarget = "99.999" // 99.999% - 5.26 minutes/year downtime
)

// DowntimePerYear returns the approximate downtime per year for the target.
func (a AvailabilityTarget) DowntimePerYear() string {
	switch a {
	case Availability99:
		return "3.65 days"
	case Availability999:
		return "8.76 hours"
	case Availability9999:
		return "52.6 minutes"
	case Availability99999:
		return "5.26 minutes"
	default:
		return "unknown"
	}
}

// Security defines security configuration.
type Security struct {
	Secrets    SecretsConfig    `json:"secrets,omitempty" yaml:"secrets,omitempty"`
	Encryption EncryptionConfig `json:"encryption,omitempty" yaml:"encryption,omitempty"`
	Auth       AuthConfig       `json:"auth,omitempty" yaml:"auth,omitempty"`
}

// SecretsConfig defines secrets management.
type SecretsConfig struct {
	Provider string `json:"provider,omitempty" yaml:"provider,omitempty"` // e.g., "aws-secrets-manager", "vault"
}

// EncryptionConfig defines encryption requirements.
type EncryptionConfig struct {
	AtRest    string `json:"atRest,omitempty" yaml:"atRest,omitempty"`       // e.g., "aes-256"
	InTransit string `json:"inTransit,omitempty" yaml:"inTransit,omitempty"` // e.g., "tls-1.3"
}

// AuthConfig defines authentication configuration.
type AuthConfig struct {
	Provider string   `json:"provider,omitempty" yaml:"provider,omitempty"` // e.g., "oidc", "saml"
	MFA      string   `json:"mfa,omitempty" yaml:"mfa,omitempty"`           // e.g., "required", "optional"
	Methods  []string `json:"methods,omitempty" yaml:"methods,omitempty"`   // e.g., ["jwt", "api-key"]
}

// Prompts contains LLM guidance prompts for decision-making.
type Prompts struct {
	LanguageChoice     string `json:"languageChoice,omitempty" yaml:"languageChoice,omitempty"`
	APIDesign          string `json:"apiDesign,omitempty" yaml:"apiDesign,omitempty"`
	DatabaseChoice     string `json:"databaseChoice,omitempty" yaml:"databaseChoice,omitempty"`
	TenancyChoice      string `json:"tenancyChoice,omitempty" yaml:"tenancyChoice,omitempty"`
	AvailabilityChoice string `json:"availabilityChoice,omitempty" yaml:"availabilityChoice,omitempty"`
}

// ValidTenancyModels returns all valid tenancy models.
func ValidTenancyModels() []TenancyModel {
	return []TenancyModel{TenancySingleTenant, TenancyMultiTenant}
}

// ValidAvailabilityTargets returns all valid availability targets.
func ValidAvailabilityTargets() []AvailabilityTarget {
	return []AvailabilityTarget{
		Availability99,
		Availability999,
		Availability9999,
		Availability99999,
	}
}

// ValidIaCTools returns all valid IaC tools.
func ValidIaCTools() []IaCTool {
	return []IaCTool{
		IaCPulumi,
		IaCCDK,
		IaCTerraform,
		IaCCloudFormation,
		IaCNone,
	}
}

// ValidLocalDevTargets returns all valid local dev targets.
func ValidLocalDevTargets() []LocalDevTarget {
	return []LocalDevTarget{
		LocalDevBinaries,
		LocalDevPodman,
		LocalDevDocker,
		LocalDevLocalStack,
		LocalDevMinikube,
		LocalDevKind,
	}
}
