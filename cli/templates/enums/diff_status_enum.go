package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strconv"
	"strings"
)

// DiffStatusType enumerates the high-level outcomes of a single Diff call.
type DiffStatusType int

const (
	// DiffNoChange means the template body matches the on-disk block
	// (or both are absent). Exit code 0 territory.
	DiffNoChange DiffStatusType = 0
	// DiffMissingBlock means the file exists but has no gitmap block
	// for tag — `add` would insert one.
	DiffMissingBlock DiffStatusType = 1
	// DiffMissingFile means the file itself is absent — `add` would
	// create it.
	DiffMissingFile DiffStatusType = 2
	// DiffBlockChanged means the file has a gitmap block but its body
	// differs from the template — `add` would update it.
	DiffBlockChanged DiffStatusType = 3
)

// IsValid reports whether the value matches a known variant.
func (d DiffStatusType) IsValid() bool {
	switch d {
	case DiffNoChange, DiffMissingBlock, DiffMissingFile, DiffBlockChanged:
		return true
	default:
		return false
	}
}

// String returns the numeric representation of the value.
func (d DiffStatusType) String() string {
	return strconv.Itoa(int(d))
}

// ParseDiffStatusType parses a raw string into a validated DiffStatusType.
func ParseDiffStatusType(s string) (DiffStatusType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	if n, err := strconv.Atoi(trimmed); err == nil {
		val := DiffStatusType(n)
		if val.IsValid() {
			return val, nil
		}
	}
	for _, val := range AllDiffStatusTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "diffnochange":
		return DiffNoChange, nil
	case "nochange":
		return DiffNoChange, nil
	case "diffmissingblock":
		return DiffMissingBlock, nil
	case "missingblock":
		return DiffMissingBlock, nil
	case "diffmissingfile":
		return DiffMissingFile, nil
	case "missingfile":
		return DiffMissingFile, nil
	case "diffblockchanged":
		return DiffBlockChanged, nil
	case "blockchanged":
		return DiffBlockChanged, nil
	}
	return 0, apperror.NewValidation("ParseDiffStatusType", "E1000", fmt.Sprintf("invalid diffstatustype: %q", s)).WithContext("value", s)
}

// AllDiffStatusTypes returns all valid variants.
func AllDiffStatusTypes() []DiffStatusType {
	return []DiffStatusType{
		DiffNoChange,
		DiffMissingBlock,
		DiffMissingFile,
		DiffBlockChanged,
	}
}

// IsNoChange reports whether the value is DiffNoChange.
func (d DiffStatusType) IsNoChange() bool {
	return d == DiffNoChange
}

// IsMissingBlock reports whether the value is DiffMissingBlock.
func (d DiffStatusType) IsMissingBlock() bool {
	return d == DiffMissingBlock
}

// IsMissingFile reports whether the value is DiffMissingFile.
func (d DiffStatusType) IsMissingFile() bool {
	return d == DiffMissingFile
}

// IsBlockChanged reports whether the value is DiffBlockChanged.
func (d DiffStatusType) IsBlockChanged() bool {
	return d == DiffBlockChanged
}

// HasStatus reports whether the value is non-empty and valid.
func (d DiffStatusType) HasStatus() bool {
	return d.IsValid()
}
