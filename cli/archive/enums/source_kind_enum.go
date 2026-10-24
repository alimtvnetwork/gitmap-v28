package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strconv"
	"strings"
)

// SourceKindType classifies one entry on the command line.
type SourceKindType int

const (
	SourceLocal SourceKindType = 0
	SourceHTTP  SourceKindType = 1
	SourceGit   SourceKindType = 2
)

// IsValid reports whether the value matches a known variant.
func (s SourceKindType) IsValid() bool {
	switch s {
	case SourceLocal, SourceHTTP, SourceGit:
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
	case "sourcelocal":
		return SourceLocal, nil
	case "local":
		return SourceLocal, nil
	case "sourcehttp":
		return SourceHTTP, nil
	case "http":
		return SourceHTTP, nil
	case "sourcegit":
		return SourceGit, nil
	case "git":
		return SourceGit, nil
	}
	return 0, apperror.NewValidation("ParseSourceKindType", "E1000", fmt.Sprintf("invalid sourcekindtype: %q", s)).WithContext("value", s)
}

// AllSourceKindTypes returns all valid variants.
func AllSourceKindTypes() []SourceKindType {
	return []SourceKindType{
		SourceLocal,
		SourceHTTP,
		SourceGit,
	}
}

// IsLocal reports whether the value is SourceLocal.
func (s SourceKindType) IsLocal() bool {
	return s == SourceLocal
}

// IsHTTP reports whether the value is SourceHTTP.
func (s SourceKindType) IsHTTP() bool {
	return s == SourceHTTP
}

// IsGit reports whether the value is SourceGit.
func (s SourceKindType) IsGit() bool {
	return s == SourceGit
}

// HasKind reports whether the value is non-empty and valid.
func (s SourceKindType) HasKind() bool {
	return s.IsValid()
}
