package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strconv"
	"strings"
)

// MergeOutcomeType describes what Merge did to the target file.
type MergeOutcomeType int

const (
	// MergeCreated means the target file did not exist and was created.
	MergeCreated MergeOutcomeType = 0
	// MergeInserted means the file existed but had no prior gitmap block;
	// the new block was appended at the end.
	MergeInserted MergeOutcomeType = 1
	// MergeUpdated means a prior gitmap block was found and its body was
	// replaced (or kept identical — see Changed).
	MergeUpdated MergeOutcomeType = 2
)

// IsValid reports whether the value matches a known variant.
func (m MergeOutcomeType) IsValid() bool {
	switch m {
	case MergeCreated, MergeInserted, MergeUpdated:
		return true
	default:
		return false
	}
}

// String returns the numeric representation of the value.
func (m MergeOutcomeType) String() string {
	return strconv.Itoa(int(m))
}

// ParseMergeOutcomeType parses a raw string into a validated MergeOutcomeType.
func ParseMergeOutcomeType(s string) (MergeOutcomeType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	if n, err := strconv.Atoi(trimmed); err == nil {
		val := MergeOutcomeType(n)
		if val.IsValid() {
			return val, nil
		}
	}
	for _, val := range AllMergeOutcomeTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "mergecreated":
		return MergeCreated, nil
	case "created":
		return MergeCreated, nil
	case "mergeinserted":
		return MergeInserted, nil
	case "inserted":
		return MergeInserted, nil
	case "mergeupdated":
		return MergeUpdated, nil
	case "updated":
		return MergeUpdated, nil
	}
	return 0, apperror.NewValidation("ParseMergeOutcomeType", "E1000", fmt.Sprintf("invalid mergeoutcometype: %q", s)).WithContext("value", s)
}

// AllMergeOutcomeTypes returns all valid variants.
func AllMergeOutcomeTypes() []MergeOutcomeType {
	return []MergeOutcomeType{
		MergeCreated,
		MergeInserted,
		MergeUpdated,
	}
}

// IsCreated reports whether the value is MergeCreated.
func (m MergeOutcomeType) IsCreated() bool {
	return m == MergeCreated
}

// IsInserted reports whether the value is MergeInserted.
func (m MergeOutcomeType) IsInserted() bool {
	return m == MergeInserted
}

// IsUpdated reports whether the value is MergeUpdated.
func (m MergeOutcomeType) IsUpdated() bool {
	return m == MergeUpdated
}

// HasOutcome reports whether the value is non-empty and valid.
func (m MergeOutcomeType) HasOutcome() bool {
	return m.IsValid()
}
