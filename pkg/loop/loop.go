// Package loop defines the intermediate representation for AI development
// loop systems: cycles of human-gated AI work such as the ProductBuildersHQ
// two-loop system (Product Loop + Builder Loop) and single-loop reference
// systems like AWS AI-DLC.
//
// A System contains one or more Loops. Each Loop is an ordered cycle of
// Stations (the last station feeds back to the first). Stations may be gates —
// approval checkpoints whose authority (human vs. policy) shifts with a team's
// autonomy level — and may map to spec types, definition workflows, or
// execution integrations from the rest of this library, letting consumers
// render a loop as the organizing geography of the catalog.
//
// Seams connect stations across loops (e.g., the Product Loop's Approve
// station hands the Product Baseline to the Builder Loop's Accept station).
// This package holds declarative data only; consumers (the website,
// VisionStudio) render and act on it.
package loop

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Actor identifies who performs a station's work.
type Actor string

const (
	// ActorHuman means the station is performed by people.
	ActorHuman Actor = "human"

	// ActorAI means the station is performed by AI (agents, models).
	ActorAI Actor = "ai"

	// ActorHumanAI means the station is collaborative.
	ActorHumanAI Actor = "human+ai"
)

// GateAuthority identifies who holds a gate's approval authority by default.
type GateAuthority string

const (
	// AuthorityHuman means a person must approve.
	AuthorityHuman GateAuthority = "human"

	// AuthorityPolicy means policy-governed automation may approve.
	AuthorityPolicy GateAuthority = "policy"
)

// System is a complete loop system.
type System struct {
	// ID is the canonical identifier (e.g., "pbhq-two-loop", "aws-ai-dlc").
	ID string `json:"id" yaml:"id" jsonschema:"required,description=Canonical identifier for the loop system"`

	// Name is the human-readable name.
	Name string `json:"name" yaml:"name" jsonschema:"required,description=Human-readable name"`

	// Vendor is the defining organization (e.g., "ProductBuildersHQ", "AWS").
	Vendor string `json:"vendor,omitempty" yaml:"vendor,omitempty" jsonschema:"description=Defining organization"`

	// Description explains the system.
	Description string `json:"description,omitempty" yaml:"description,omitempty" jsonschema:"description=System purpose and scope"`

	// Reference is a URL to canonical documentation.
	Reference string `json:"reference,omitempty" yaml:"reference,omitempty" jsonschema:"format=uri,description=Documentation URL"`

	// Loops are the cycles in this system, outermost first.
	Loops []Loop `json:"loops" yaml:"loops" jsonschema:"required,description=Cycles in this system (outermost first)"`

	// Seams connect stations across loops.
	Seams []Seam `json:"seams,omitempty" yaml:"seams,omitempty" jsonschema:"description=Cross-loop connections"`
}

// Loop is one cycle of stations. The final station feeds back to the first.
type Loop struct {
	// ID is the loop identifier, unique within the system (e.g., "product").
	ID string `json:"id" yaml:"id" jsonschema:"required,description=Loop identifier (unique within the system)"`

	// Name is the human-readable name (e.g., "Product Loop").
	Name string `json:"name" yaml:"name" jsonschema:"required,description=Loop name"`

	// Owner is the accountable role (e.g., "Product person", "Builder").
	Owner string `json:"owner,omitempty" yaml:"owner,omitempty" jsonschema:"description=Accountable role"`

	// Question is the loop's verification question (e.g., "Right thing?").
	Question string `json:"question,omitempty" yaml:"question,omitempty" jsonschema:"description=The loop's verification question"`

	// Description explains the loop.
	Description string `json:"description,omitempty" yaml:"description,omitempty" jsonschema:"description=Loop purpose"`

	// Stations are the ordered steps; the last cycles back to the first.
	Stations []Station `json:"stations" yaml:"stations" jsonschema:"required,description=Ordered stations (last cycles back to first)"`
}

// Station is one step in a loop.
type Station struct {
	// ID is the station identifier, unique within the loop (e.g., "define").
	ID string `json:"id" yaml:"id" jsonschema:"required,description=Station identifier (unique within the loop)"`

	// Name is the human-readable name (e.g., "Define").
	Name string `json:"name" yaml:"name" jsonschema:"required,description=Station name"`

	// Actor performs the station's work.
	Actor Actor `json:"actor" yaml:"actor" jsonschema:"required,enum=human,enum=ai,enum=human+ai"`

	// Description explains the station.
	Description string `json:"description,omitempty" yaml:"description,omitempty" jsonschema:"description=Station purpose"`

	// Gate marks this station as an approval checkpoint.
	Gate bool `json:"gate,omitempty" yaml:"gate,omitempty" jsonschema:"description=True if this station is an approval checkpoint"`

	// GateAuthority is the default approval authority when Gate is true.
	GateAuthority GateAuthority `json:"gate_authority,omitempty" yaml:"gate_authority,omitempty" jsonschema:"enum=human,enum=policy,description=Default approval authority when gate is true"`

	// AutonomyNote describes how the station changes as autonomy (ASDM level)
	// rises (e.g., "human at L5 and below; policy-gated at L6+").
	AutonomyNote string `json:"autonomy_note,omitempty" yaml:"autonomy_note,omitempty" jsonschema:"description=How the station changes as autonomy rises"`

	// SpecTypes lists spec type IDs produced or consumed at this station.
	SpecTypes []string `json:"spec_types,omitempty" yaml:"spec_types,omitempty" jsonschema:"description=Spec type IDs produced or consumed here"`

	// Workflows lists definition workflow names that operate at this station.
	Workflows []string `json:"workflows,omitempty" yaml:"workflows,omitempty" jsonschema:"description=Definition workflow names mapped to this station"`

	// Integrations lists execution integration IDs that operate at this station.
	Integrations []string `json:"integrations,omitempty" yaml:"integrations,omitempty" jsonschema:"description=Execution integration IDs mapped to this station"`
}

// Seam connects a station in one loop to a station in another.
type Seam struct {
	// From is the source station reference as "loopID.stationID".
	From string `json:"from" yaml:"from" jsonschema:"required,description=Source station as loopID.stationID"`

	// To is the destination station reference as "loopID.stationID".
	To string `json:"to" yaml:"to" jsonschema:"required,description=Destination station as loopID.stationID"`

	// Artifact names what crosses the seam (e.g., "Product Baseline").
	Artifact string `json:"artifact,omitempty" yaml:"artifact,omitempty" jsonschema:"description=What crosses the seam"`

	// Description explains the seam.
	Description string `json:"description,omitempty" yaml:"description,omitempty" jsonschema:"description=Seam explanation"`
}

// ParseYAML parses a System from YAML bytes.
func ParseYAML(data []byte) (*System, error) {
	var s System
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("unmarshaling loop system: %w", err)
	}
	return &s, nil
}

// Validate checks required fields and referential integrity.
func (s *System) Validate() error {
	if s.ID == "" {
		return fmt.Errorf("loop system: id is required")
	}
	if s.Name == "" {
		return fmt.Errorf("loop system %q: name is required", s.ID)
	}
	if len(s.Loops) == 0 {
		return fmt.Errorf("loop system %q: at least one loop is required", s.ID)
	}

	stationRefs := make(map[string]bool)
	loopIDs := make(map[string]bool, len(s.Loops))
	for _, l := range s.Loops {
		if l.ID == "" {
			return fmt.Errorf("loop system %q: loop with empty id", s.ID)
		}
		if loopIDs[l.ID] {
			return fmt.Errorf("loop system %q: duplicate loop id %q", s.ID, l.ID)
		}
		loopIDs[l.ID] = true
		if len(l.Stations) == 0 {
			return fmt.Errorf("loop system %q: loop %q has no stations", s.ID, l.ID)
		}

		stationIDs := make(map[string]bool, len(l.Stations))
		for _, st := range l.Stations {
			if st.ID == "" {
				return fmt.Errorf("loop system %q: loop %q has station with empty id", s.ID, l.ID)
			}
			if stationIDs[st.ID] {
				return fmt.Errorf("loop system %q: loop %q duplicate station id %q", s.ID, l.ID, st.ID)
			}
			stationIDs[st.ID] = true
			switch st.Actor {
			case ActorHuman, ActorAI, ActorHumanAI:
			default:
				return fmt.Errorf("loop system %q: station %s.%s has invalid actor %q", s.ID, l.ID, st.ID, st.Actor)
			}
			if st.GateAuthority != "" && !st.Gate {
				return fmt.Errorf("loop system %q: station %s.%s sets gate_authority without gate", s.ID, l.ID, st.ID)
			}
			stationRefs[l.ID+"."+st.ID] = true
		}
	}

	for i, seam := range s.Seams {
		for _, ref := range []string{seam.From, seam.To} {
			if !strings.Contains(ref, ".") || !stationRefs[ref] {
				return fmt.Errorf("loop system %q: seam %d references unknown station %q", s.ID, i, ref)
			}
		}
	}

	return nil
}
