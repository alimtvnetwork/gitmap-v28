package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strconv"
	"strings"
)

// SourceKindType enumerates spec §2.3 cases. PascalCase per Core memory
// rules; literal strings live in this file (single use site).
type SourceKindType uint8

const (
	// SourceKindCloned is spec §2.3 case 1: URL → git clone.
	SourceKindCloned SourceKindType = 1
	// SourceKindExistingRepo is spec §2.3 case 2: dir w/ .git → reuse.
	SourceKindExistingRepo SourceKindType = 2
	// SourceKindInitInPlace is spec §2.3 case 3: dir w/o .git → git init.
	SourceKindInitInPlace SourceKindType = 3
	// SourceKindCreatedAndInit is spec §2.3 case 4: missing → mkdir + init.
	SourceKindCreatedAndInit SourceKindType = 4
)

// IsValid reports whether the value matches a known variant.
func (s SourceKindType) IsValid() bool {
	switch s {
	case SourceKindCloned, SourceKindExistingRepo, SourceKindInitInPlace, SourceKindCreatedAndInit:
		return true
	default:
		return false
	}
}

// String returns the numeric representation of the value.
func (s SourceKindType) String() string {
	return strconv.Itoa(int(s))
}

// ParseSourceKindType parses a raw string into a validated SourceKindType.
func ParseSourceKindType(s string) (SourceKindType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	if n, err := strconv.Atoi(trimmed); err == nil {
		val := SourceKindType(n)
		if val.IsValid() {
			return val, nil
		}
	}
	for _, val := range AllSourceKindTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "sourcekindcloned":
		return SourceKindCloned, nil
	case "cloned":
		return SourceKindCloned, nil
	case "sourcekindexistingrepo":
		return SourceKindExistingRepo, nil
	case "existingrepo":
		return SourceKindExistingRepo, nil
	case "sourcekindinitinplace":
		return SourceKindInitInPlace, nil
	case "initinplace":
		return SourceKindInitInPlace, nil
	case "sourcekindcreatedandinit":
		return SourceKindCreatedAndInit, nil
	case "createdandinit":
		return SourceKindCreatedAndInit, nil
	}
	return 0, apperror.NewValidation("ParseSourceKindType", "E1000", fmt.Sprintf("invalid sourcekindtype: %q", s)).WithContext("value", s)
}

// AllSourceKindTypes returns all valid variants.
func AllSourceKindTypes() []SourceKindType {
	return []SourceKindType{
		SourceKindCloned,
		SourceKindExistingRepo,
		SourceKindInitInPlace,
		SourceKindCreatedAndInit,
	}
}

// IsCloned reports whether the value is SourceKindCloned.
func (s SourceKindType) IsCloned() bool {
	return s == SourceKindCloned
}

// IsExistingRepo reports whether the value is SourceKindExistingRepo.
func (s SourceKindType) IsExistingRepo() bool {
	return s == SourceKindExistingRepo
}

// IsInitInPlace reports whether the value is SourceKindInitInPlace.
func (s SourceKindType) IsInitInPlace() bool {
	return s == SourceKindInitInPlace
}

// IsCreatedAndInit reports whether the value is SourceKindCreatedAndInit.
func (s SourceKindType) IsCreatedAndInit() bool {
	return s == SourceKindCreatedAndInit
}

// HasKind reports whether the value is non-empty and valid.
func (s SourceKindType) HasKind() bool {
	return s.IsValid()
}
