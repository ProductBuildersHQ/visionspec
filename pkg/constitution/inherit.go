package constitution

import (
	"fmt"
	"reflect"
)

// Merge merges a child constitution into a parent, with child values taking precedence.
// Zero values in child are ignored (parent value preserved).
// This implements the inheritance model: org → team → project.
func Merge(parent, child *Constitution) *Constitution {
	if parent == nil {
		return child
	}
	if child == nil {
		return parent
	}

	result := &Constitution{
		APIVersion:     coalesce(child.APIVersion, parent.APIVersion),
		Kind:           coalesce(child.Kind, parent.Kind),
		Metadata:       mergeMetadata(parent.Metadata, child.Metadata),
		Technical:      mergeTechnical(parent.Technical, child.Technical),
		Infrastructure: mergeInfrastructure(parent.Infrastructure, child.Infrastructure),
		Security:       mergeSecurity(parent.Security, child.Security),
		Prompts:        mergePrompts(parent.Prompts, child.Prompts),
	}

	return result
}

// Resolve resolves a constitution chain, applying inheritance from root to leaf.
// constitutions should be ordered from root (org) to leaf (project).
func Resolve(constitutions ...*Constitution) (*Constitution, error) {
	if len(constitutions) == 0 {
		return nil, fmt.Errorf("no constitutions provided")
	}

	result := constitutions[0]
	for i := 1; i < len(constitutions); i++ {
		result = Merge(result, constitutions[i])
	}

	return result, nil
}

// ValidateInheritance validates that a child constitution's inherits field
// references a valid parent.
func ValidateInheritance(child *Constitution, availableParents map[string]*Constitution) error {
	if child.Metadata.Inherits == "" {
		// No inheritance, valid for org-level
		if child.Metadata.Level != LevelOrganization {
			return fmt.Errorf("non-organization constitution must specify 'inherits'")
		}
		return nil
	}

	if _, ok := availableParents[child.Metadata.Inherits]; !ok {
		return fmt.Errorf("inherited constitution not found: %s", child.Metadata.Inherits)
	}

	return nil
}

// coalesce returns the first non-zero string.
func coalesce(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

// mergeMetadata merges metadata, child takes precedence.
func mergeMetadata(parent, child Metadata) Metadata {
	return Metadata{
		Name:        coalesce(child.Name, parent.Name),
		Level:       child.Level,    // Always use child's level
		Inherits:    child.Inherits, // Always use child's inherits
		Version:     coalesce(child.Version, parent.Version),
		Description: coalesce(child.Description, parent.Description),
		UpdatedAt:   child.UpdatedAt, // Always use child's timestamp
	}
}

// mergeTechnical merges technical configurations.
func mergeTechnical(parent, child Technical) Technical {
	return Technical{
		Languages:     mergeLanguages(parent.Languages, child.Languages),
		APIs:          mergeAPIs(parent.APIs, child.APIs),
		Database:      mergeDatabase(parent.Database, child.Database),
		Tenancy:       mergeTenancy(parent.Tenancy, child.Tenancy),
		Observability: mergeTechObservability(parent.Observability, child.Observability),
	}
}

func mergeLanguages(parent, child Languages) Languages {
	return Languages{
		Backend:  mergeLanguageChoice(parent.Backend, child.Backend),
		Frontend: mergeLanguageChoice(parent.Frontend, child.Frontend),
		WASM:     coalesce(child.WASM, parent.WASM),
	}
}

func mergeLanguageChoice(parent, child LanguageChoice) LanguageChoice {
	result := LanguageChoice{
		Primary:           coalesce(child.Primary, parent.Primary),
		ExceptionsRequire: coalesce(child.ExceptionsRequire, parent.ExceptionsRequire),
	}
	if len(child.Allowed) > 0 {
		result.Allowed = child.Allowed
	} else {
		result.Allowed = parent.Allowed
	}
	return result
}

func mergeAPIs(parent, child APIs) APIs {
	return APIs{
		REST: RESTConfig{
			Framework:  coalesce(child.REST.Framework, parent.REST.Framework),
			SpecFormat: coalesce(child.REST.SpecFormat, parent.REST.SpecFormat),
			StyleGuide: coalesce(child.REST.StyleGuide, parent.REST.StyleGuide),
		},
		GRPC: GRPCConfig{
			Framework: coalesce(child.GRPC.Framework, parent.GRPC.Framework),
		},
	}
}

func mergeDatabase(parent, child Database) Database {
	return Database{
		Relational:   coalesce(child.Relational, parent.Relational),
		Document:     coalesce(child.Document, parent.Document),
		Cache:        coalesce(child.Cache, parent.Cache),
		Search:       coalesce(child.Search, parent.Search),
		MultiTenancy: coalesce(child.MultiTenancy, parent.MultiTenancy),
	}
}

func mergeTenancy(parent, child Tenancy) Tenancy {
	model := child.Model
	if model == "" {
		model = parent.Model
	}
	return Tenancy{
		Model:          model,
		Isolation:      coalesce(child.Isolation, parent.Isolation),
		TenantIDHeader: coalesce(child.TenantIDHeader, parent.TenantIDHeader),
		TenantIDClaim:  coalesce(child.TenantIDClaim, parent.TenantIDClaim),
	}
}

func mergeTechObservability(parent, child TechObservability) TechObservability {
	return TechObservability{
		Library:        coalesce(child.Library, parent.Library),
		AutoInstrument: child.AutoInstrument || parent.AutoInstrument,
		MetricsSDK:     coalesce(child.MetricsSDK, parent.MetricsSDK),
		TracesSDK:      coalesce(child.TracesSDK, parent.TracesSDK),
		LoggingSDK:     coalesce(child.LoggingSDK, parent.LoggingSDK),
	}
}

// mergeInfrastructure merges infrastructure configurations.
func mergeInfrastructure(parent, child Infrastructure) Infrastructure {
	return Infrastructure{
		IaC:           mergeIaC(parent.IaC, child.IaC),
		Observability: mergeInfraObservability(parent.Observability, child.Observability),
		LocalDev:      mergeLocalDev(parent.LocalDev, child.LocalDev),
		Cloud:         mergeCloud(parent.Cloud, child.Cloud),
		Availability:  mergeAvailability(parent.Availability, child.Availability),
	}
}

func mergeIaC(parent, child IaC) IaC {
	tool := child.Tool
	if tool == "" {
		tool = parent.Tool
	}
	return IaC{
		Tool:     tool,
		Language: coalesce(child.Language, parent.Language),
		RepoPath: coalesce(child.RepoPath, parent.RepoPath),
	}
}

func mergeInfraObservability(parent, child InfraObservability) InfraObservability {
	return InfraObservability{
		Metrics: mergeObservabilityPillar(parent.Metrics, child.Metrics),
		Traces:  mergeObservabilityPillar(parent.Traces, child.Traces),
		Logging: mergeObservabilityPillar(parent.Logging, child.Logging),
	}
}

func mergeObservabilityPillar(parent, child ObservabilityPillar) ObservabilityPillar {
	return ObservabilityPillar{
		Platform:      coalesce(child.Platform, parent.Platform),
		Visualization: coalesce(child.Visualization, parent.Visualization),
		Backend:       coalesce(child.Backend, parent.Backend),
		Collector:     coalesce(child.Collector, parent.Collector),
	}
}

func mergeLocalDev(parent, child LocalDev) LocalDev {
	if len(child.Priority) > 0 {
		return child
	}
	return parent
}

func mergeCloud(parent, child Cloud) Cloud {
	result := Cloud{
		Provider: coalesce(child.Provider, parent.Provider),
	}
	if len(child.Regions) > 0 {
		result.Regions = child.Regions
	} else {
		result.Regions = parent.Regions
	}
	return result
}

func mergeAvailability(parent, child Availability) Availability {
	target := child.Target
	if target == "" {
		target = parent.Target
	}
	return Availability{
		Target:      target,
		RTO:         coalesce(child.RTO, parent.RTO),
		RPO:         coalesce(child.RPO, parent.RPO),
		MultiRegion: child.MultiRegion || parent.MultiRegion,
		MultiAZ:     child.MultiAZ || parent.MultiAZ,
		DRStrategy:  coalesce(child.DRStrategy, parent.DRStrategy),
	}
}

// mergeSecurity merges security configurations.
func mergeSecurity(parent, child Security) Security {
	return Security{
		Secrets: SecretsConfig{
			Provider: coalesce(child.Secrets.Provider, parent.Secrets.Provider),
		},
		Encryption: EncryptionConfig{
			AtRest:    coalesce(child.Encryption.AtRest, parent.Encryption.AtRest),
			InTransit: coalesce(child.Encryption.InTransit, parent.Encryption.InTransit),
		},
		Auth: mergeAuth(parent.Auth, child.Auth),
	}
}

func mergeAuth(parent, child AuthConfig) AuthConfig {
	result := AuthConfig{
		Provider: coalesce(child.Provider, parent.Provider),
		MFA:      coalesce(child.MFA, parent.MFA),
	}
	if len(child.Methods) > 0 {
		result.Methods = child.Methods
	} else {
		result.Methods = parent.Methods
	}
	return result
}

// mergePrompts merges prompt configurations.
func mergePrompts(parent, child Prompts) Prompts {
	return Prompts{
		LanguageChoice:     coalesce(child.LanguageChoice, parent.LanguageChoice),
		APIDesign:          coalesce(child.APIDesign, parent.APIDesign),
		DatabaseChoice:     coalesce(child.DatabaseChoice, parent.DatabaseChoice),
		TenancyChoice:      coalesce(child.TenancyChoice, parent.TenancyChoice),
		AvailabilityChoice: coalesce(child.AvailabilityChoice, parent.AvailabilityChoice),
	}
}

// IsZero returns true if the value is the zero value for its type.
func IsZero(v interface{}) bool {
	return reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface())
}
