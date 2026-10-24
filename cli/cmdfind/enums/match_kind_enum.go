package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strconv"
	"strings"
)

// MatchKindType is a strongly-typed enum.
type MatchKindType int

const (
	MatchExact      MatchKindType = 0
	MatchContains   MatchKindType = 1
	MatchStartsWith MatchKindType = 2
	MatchEndsWith   MatchKindType = 3
	MatchWildcard   MatchKindType = 4
)

// IsValid reports whether the value matches a known variant.
func (m MatchKindType) IsValid() bool {
	switch m {
	case MatchExact, MatchContains, MatchStartsWith, MatchEndsWith, MatchWildcard:
		return true
	default:
		return false
	}
}

// String returns the numeric representation of the value.
func (m MatchKindType) String() string {
	return strconv.Itoa(int(m))
}

// ParseMatchKindType parses a raw string into a validated MatchKindType.
func ParseMatchKindType(s string) (MatchKindType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	if n, err := strconv.Atoi(trimmed); err == nil {
		val := MatchKindType(n)
		if val.IsValid() {
			return val, nil
		}
	}
	for _, val := range AllMatchKindTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "matchexact":
		return MatchExact, nil
	case "exact":
		return MatchExact, nil
	case "matchcontains":
		return MatchContains, nil
	case "contains":
		return MatchContains, nil
	case "matchstartswith":
		return MatchStartsWith, nil
	case "startswith":
		return MatchStartsWith, nil
	case "matchendswith":
		return MatchEndsWith, nil
	case "endswith":
		return MatchEndsWith, nil
	case "matchwildcard":
		return MatchWildcard, nil
	case "wildcard":
		return MatchWildcard, nil
	}
	return 0, apperror.NewValidation("ParseMatchKindType", "E1000", fmt.Sprintf("invalid matchkindtype: %q", s)).WithContext("value", s)
}

// AllMatchKindTypes returns all valid variants.
func AllMatchKindTypes() []MatchKindType {
	return []MatchKindType{
		MatchExact,
		MatchContains,
		MatchStartsWith,
		MatchEndsWith,
		MatchWildcard,
	}
}

// IsExact reports whether the value is MatchExact.
func (m MatchKindType) IsExact() bool {
	return m == MatchExact
}

// IsContains reports whether the value is MatchContains.
func (m MatchKindType) IsContains() bool {
	return m == MatchContains
}

// IsStartsWith reports whether the value is MatchStartsWith.
func (m MatchKindType) IsStartsWith() bool {
	return m == MatchStartsWith
}

// IsEndsWith reports whether the value is MatchEndsWith.
func (m MatchKindType) IsEndsWith() bool {
	return m == MatchEndsWith
}

// IsWildcard reports whether the value is MatchWildcard.
func (m MatchKindType) IsWildcard() bool {
	return m == MatchWildcard
}

// HasKind reports whether the value is non-empty and valid.
func (m MatchKindType) HasKind() bool {
	return m.IsValid()
}
