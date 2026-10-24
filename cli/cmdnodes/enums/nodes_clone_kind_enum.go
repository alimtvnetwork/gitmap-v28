package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// NodesCloneKindType defines the clone operation mode.
type NodesCloneKindType string

const (
	// CloneKindClone executes standard multi-repo or single-repo clone.
	CloneKindClone NodesCloneKindType = "clone"
	// CloneKindCFR executes clone-fix-repo (remediation & auto-setup).
	CloneKindCFR NodesCloneKindType = "cfr"
	// CloneKindCFRP executes clone-fix-repo with public visibility promotion.
	CloneKindCFRP NodesCloneKindType = "cfrp"
)

// IsValid reports whether the clone kind matches a known kind.
func (k NodesCloneKindType) IsValid() bool {
	switch k {
	case CloneKindClone, CloneKindCFR, CloneKindCFRP:
		return true
	default:
		return false
	}
}

// String returns the string representation of the clone kind.
func (k NodesCloneKindType) String() string {
	return string(k)
}

// ParseNodesCloneKindType parses a raw string into a validated NodesCloneKindType.
func ParseNodesCloneKindType(s string) (NodesCloneKindType, *apperror.AppError) {
	val := NodesCloneKindType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseNodesCloneKindType", "E1000", fmt.Sprintf("invalid nodes clone kind: %q", s)).WithContext("value", s)
}

// AllNodesCloneKindTypes returns all valid NodesCloneKindType variants.
func AllNodesCloneKindTypes() []NodesCloneKindType {
	return []NodesCloneKindType{CloneKindClone, CloneKindCFR, CloneKindCFRP}
}

// IsClone reports whether the kind is clone.
func (k NodesCloneKindType) IsClone() bool { return k == CloneKindClone }

// IsCFR reports whether the kind is cfr.
func (k NodesCloneKindType) IsCFR() bool { return k == CloneKindCFR }

// IsCFRP reports whether the kind is cfrp.
func (k NodesCloneKindType) IsCFRP() bool { return k == CloneKindCFRP }
