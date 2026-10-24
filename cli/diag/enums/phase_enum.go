package enums

import (
	"fmt"
	"strings"
)

// PhaseType tags which subsystem produced the failure. Stored as the
// JSON top-level array key (`scan` / `clone`) so consumers can
// branch on it without parsing the inner Entry first.
type PhaseType string

const (
	// PhaseScan covers everything that runs under `gitmap scan`:
	// scanner ReadDir failures, background ls-remote / shallow-clone
	// probe errors, project-detection failures we choose to surface.
	PhaseScan PhaseType = "scan"
	// PhaseClone covers everything under `gitmap clone-next` batch
	// mode (--all / --csv): per-repo clone failures, post-clone
	// hooks that errored, and any other non-skip cn outcome.
	PhaseClone PhaseType = "clone"
)

// IsValid reports whether the value matches a known variant.
func (p PhaseType) IsValid() bool {
	switch p {
	case PhaseScan, PhaseClone:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (p PhaseType) String() string {
	return string(p)
}

// ParsePhaseType parses a raw string into a validated PhaseType.
// ParsePhaseType returns a plain error (not *apperror.AppError) by
// recorded exemption: apperror imports cli/diag, so an AppError Parse
// here would close an import cycle (diag -> diag/enums -> apperror -> diag).
func ParsePhaseType(s string) (PhaseType, error) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllPhaseTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "phasescan":
		return PhaseScan, nil
	case "scan":
		return PhaseScan, nil
	case "phaseclone":
		return PhaseClone, nil
	case "clone":
		return PhaseClone, nil
	}
	return "", fmt.Errorf("invalid phasetype: %q", s)
}

// AllPhaseTypes returns all valid variants.
func AllPhaseTypes() []PhaseType {
	return []PhaseType{
		PhaseScan,
		PhaseClone,
	}
}

// IsScan reports whether the value is PhaseScan.
func (p PhaseType) IsScan() bool {
	return p == PhaseScan
}

// IsClone reports whether the value is PhaseClone.
func (p PhaseType) IsClone() bool {
	return p == PhaseClone
}

// HasPhase reports whether the value is non-empty and valid.
func (p PhaseType) HasPhase() bool {
	return p.IsValid()
}
