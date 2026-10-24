package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// CommitActionKindType defines semantic action types.
type CommitActionKindType string

const (
	CommitActionStandard CommitActionKindType = "commits"
	CommitActionFeature  CommitActionKindType = "cpf"
	CommitActionBug      CommitActionKindType = "cpb"
	CommitActionRelease  CommitActionKindType = "cpr"
	CommitActionFix      CommitActionKindType = "commit-fix"
)

// IsValid reports whether the commit action kind matches a known kind.
func (k CommitActionKindType) IsValid() bool {
	switch k {
	case CommitActionStandard, CommitActionFeature, CommitActionBug, CommitActionRelease, CommitActionFix:
		return true
	default:
		return false
	}
}

// String returns the string representation of the commit action kind.
func (k CommitActionKindType) String() string {
	return string(k)
}

// ParseCommitActionKindType parses a raw string into a validated CommitActionKindType.
func ParseCommitActionKindType(s string) (CommitActionKindType, *apperror.AppError) {
	val := CommitActionKindType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseCommitActionKindType", "E1000", fmt.Sprintf("invalid commit action kind: %q", s)).WithContext("value", s)
}

// AllCommitActionKindTypes returns all valid CommitActionKindType variants.
func AllCommitActionKindTypes() []CommitActionKindType {
	return []CommitActionKindType{CommitActionStandard, CommitActionFeature, CommitActionBug, CommitActionRelease, CommitActionFix}
}

// IsStandard reports whether the kind is commits.
func (k CommitActionKindType) IsStandard() bool { return k == CommitActionStandard }

// IsFeature reports whether the kind is cpf.
func (k CommitActionKindType) IsFeature() bool { return k == CommitActionFeature }

// IsBug reports whether the kind is cpb.
func (k CommitActionKindType) IsBug() bool { return k == CommitActionBug }

// IsRelease reports whether the kind is cpr.
func (k CommitActionKindType) IsRelease() bool { return k == CommitActionRelease }

// IsFix reports whether the kind is commit-fix.
func (k CommitActionKindType) IsFix() bool { return k == CommitActionFix }
