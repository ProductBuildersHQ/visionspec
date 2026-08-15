package apptypes

import (
	"github.com/ProductBuildersHQ/visionspec/pkg/constitution"
	"github.com/ProductBuildersHQ/visionspec/pkg/types"
)

// MicroserviceSpec returns the default AppTypeSpec for microservices.
func MicroserviceSpec() *AppTypeSpec {
	embedded := false
	stateful := false
	containerized := true
	horizontal := true
	cdn := false

	return &AppTypeSpec{
		APIVersion: "visionspec/v1",
		Kind:       "AppTypeSpec",
		Metadata: AppMetadata{
			Name:        AppTypeMicroservice,
			Version:     "1.0",
			Description: "Backend service with single bounded context, independent deployment lifecycle",
		},
		Artifacts: Artifacts{
			Required: []ArtifactType{
				ArtifactBinary,
				ArtifactContainerImage,
			},
			Optional: []ArtifactType{
				ArtifactOpenAPISpec,
				ArtifactProtoSpec,
				ArtifactHelmChart,
				ArtifactPulumiModule,
			},
		},
		Defaults: AppDefaults{
			Technical: TechnicalDefaults{
				APIStyles:       []string{"rest", "grpc"},
				EmbeddedDB:      &embedded,
				StatefulAllowed: &stateful,
			},
			Infrastructure: InfrastructureDefaults{
				Containerized:         &containerized,
				Orchestration:         []string{"kubernetes", "ecs", "standalone"},
				HorizontalScaling:     &horizontal,
				CDNRequired:           &cdn,
				MinAvailabilityTarget: "99.9",
			},
		},
		Specs: SpecRequirements{
			Required: []types.SpecType{
				types.SpecTypeMRD,
				types.SpecTypePRD,
				types.SpecTypeTRD,
				types.SpecTypeIRD,
			},
			Optional: []types.SpecType{
				types.SpecTypeUXD,   // Only if service has UI
				types.SpecTypePress, // For major features
				types.SpecTypeFAQ,
				types.SpecTypeTPD,
			},
		},
		Prompts: AppPrompts{
			WhenToUse: `Use microservice app type when:
- Service has single bounded context (one domain/responsibility)
- Independent deployment is required
- Team owns full lifecycle from code to production
- Service needs to scale independently of others
- Clear API contract with other services
- Separate data store from other services`,
			WhenNotToUse: `Consider alternatives when:
- Early stage, domain boundaries unclear → start with modular monolith
- Small team, low complexity → monolith may be simpler
- Tight coupling with other services → may indicate wrong boundary
- Shared database required → consider bounded context within monolith
- High inter-service communication → may indicate distributed monolith anti-pattern`,
			KeyConsiderations: `Key microservice considerations:
1. API Design: REST (Huma+Chi) or gRPC (Connect-Go), not both unless necessary
2. Data: Own data store, no shared databases, eventual consistency accepted
3. Observability: Distributed tracing critical (OpenTelemetry), correlation IDs required
4. Resilience: Circuit breakers, retries with backoff, graceful degradation
5. Testing: Contract tests (Pact) for API compatibility
6. Deployment: Blue/green or canary releases, feature flags for gradual rollout`,
		},
	}
}

// MicroserviceConstraints defines validation constraints for microservices.
type MicroserviceConstraints struct {
	// Tenancy constraints
	AllowedTenancyModels []constitution.TenancyModel

	// Availability constraints
	MinAvailability constitution.AvailabilityTarget
	MaxRTO          string // Maximum acceptable RTO
	MaxRPO          string // Maximum acceptable RPO

	// Scaling constraints
	RequireHorizontalScaling bool
	RequireContainerization  bool

	// API constraints
	RequireAPISpec bool   // Must have OpenAPI or Proto spec
	APISpecFormat  string // Preferred format: "openapi" or "proto"
}

// DefaultMicroserviceConstraints returns the default constraints for microservices.
func DefaultMicroserviceConstraints() *MicroserviceConstraints {
	return &MicroserviceConstraints{
		AllowedTenancyModels: []constitution.TenancyModel{
			constitution.TenancyMultiTenant,
			constitution.TenancySingleTenant,
		},
		MinAvailability:          constitution.Availability999,
		MaxRTO:                   "1h",
		MaxRPO:                   "15m",
		RequireHorizontalScaling: true,
		RequireContainerization:  true,
		RequireAPISpec:           true,
		APISpecFormat:            "openapi",
	}
}

// ValidateForMicroservice validates a constitution against microservice constraints.
func ValidateForMicroservice(c *constitution.Constitution, constraints *MicroserviceConstraints) []ValidationIssue {
	if constraints == nil {
		constraints = DefaultMicroserviceConstraints()
	}

	var issues []ValidationIssue

	// Validate tenancy model
	if c.Technical.Tenancy.Model != "" {
		valid := false
		for _, allowed := range constraints.AllowedTenancyModels {
			if c.Technical.Tenancy.Model == allowed {
				valid = true
				break
			}
		}
		if !valid {
			issues = append(issues, ValidationIssue{
				Field:    "technical.tenancy.model",
				Message:  "tenancy model not allowed for microservice",
				Severity: SeverityError,
			})
		}
	}

	// Validate availability target
	if c.Infrastructure.Availability.Target != "" {
		if !meetsMinAvailability(c.Infrastructure.Availability.Target, constraints.MinAvailability) {
			issues = append(issues, ValidationIssue{
				Field:    "infrastructure.availability.target",
				Message:  "availability target below minimum for microservice",
				Severity: SeverityWarning,
			})
		}
	}

	// Validate containerization requirement
	// Note: This would be validated against project spec, not constitution
	// Constitution defines defaults, project spec declares actual choices

	return issues
}

// ValidationIssue represents a validation problem.
type ValidationIssue struct {
	Field    string
	Message  string
	Severity ValidationSeverity
}

// ValidationSeverity indicates the severity of a validation issue.
type ValidationSeverity string

const (
	SeverityError   ValidationSeverity = "error"
	SeverityWarning ValidationSeverity = "warning"
	SeverityInfo    ValidationSeverity = "info"
)

// meetsMinAvailability checks if target meets or exceeds minimum.
func meetsMinAvailability(target, min constitution.AvailabilityTarget) bool {
	order := map[constitution.AvailabilityTarget]int{
		constitution.Availability99:    1,
		constitution.Availability999:   2,
		constitution.Availability9999:  3,
		constitution.Availability99999: 4,
	}
	return order[target] >= order[min]
}
