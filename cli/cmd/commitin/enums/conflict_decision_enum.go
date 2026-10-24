package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// ConflictDecisionType is the action a higher-level replay loop takes for
// a single conflicting source commit.
type ConflictDecisionType int

const (
	ConflictDecisionTakeTheirs ConflictDecisionType = iota
	ConflictDecisionAbort
)

// IsValid reports whether the conflict decision matches a known decision.
func (d ConflictDecisionType) IsValid() bool {
	switch d {
	case ConflictDecisionTakeTheirs, ConflictDecisionAbort:
		return true
	default:
		return false
	}
}

// String returns the name of the conflict decision.
func (d ConflictDecisionType) String() string {
	switch d {
	case ConflictDecisionTakeTheirs:
		return "TakeTheirs"
	case ConflictDecisionAbort:
		return "Abort"
	default:
		return ""
	}
}

// ParseConflictDecisionType parses a canonical name into a validated ConflictDecisionType.
func ParseConflictDecisionType(s string) (ConflictDecisionType, *apperror.AppError) {
	for _, val := range AllConflictDecisions() {
		if strings.EqualFold(val.String(), strings.TrimSpace(s)) {
			return val, nil
		}
	}

	return 0, apperror.NewValidation("ParseConflictDecisionType", "E1000", fmt.Sprintf("invalid conflict decision: %q", s)).WithContext("value", s)
}

// AllConflictDecisions returns all valid ConflictDecisionType variants.
func AllConflictDecisions() []ConflictDecisionType {
	return []ConflictDecisionType{ConflictDecisionTakeTheirs, ConflictDecisionAbort}
}

// IsTakeTheirs reports whether the decision is take-theirs.
func (d ConflictDecisionType) IsTakeTheirs() bool { return d == ConflictDecisionTakeTheirs }

// IsAbort reports whether the decision is abort.
func (d ConflictDecisionType) IsAbort() bool { return d == ConflictDecisionAbort }
