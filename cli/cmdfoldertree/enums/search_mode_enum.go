package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// SearchModeType defines the pattern matching strategy for tree searching.
type SearchModeType string

// SearchMode is an alias for SearchModeType.
type SearchMode = SearchModeType

const (
	ModeWildcard   SearchModeType = "wildcard"
	ModeStartsWith SearchModeType = "startsWith"
	ModeContains   SearchModeType = "contains"
	ModeEndsWith   SearchModeType = "endsWith"
	ModeGrep       SearchModeType = "grep"
)

// IsValid reports whether the search mode matches a known mode.
func (m SearchModeType) IsValid() bool {
	switch m {
	case ModeWildcard, ModeStartsWith, ModeContains, ModeEndsWith, ModeGrep:
		return true
	default:
		return false
	}
}

// String returns the string representation of the search mode.
func (m SearchModeType) String() string {
	return string(m)
}

// ParseSearchModeType parses a raw string into a validated SearchModeType.
func ParseSearchModeType(s string) (SearchModeType, *apperror.AppError) {
	cleaned := strings.TrimSpace(s)
	for _, val := range AllSearchModeTypes() {
		if val.String() == cleaned {
			return val, nil
		}
	}

	return "", apperror.NewValidation("ParseSearchModeType", "E1000", fmt.Sprintf("invalid tree search mode: %q", s)).WithContext("value", s)
}

// AllSearchModeTypes returns all valid SearchModeType variants.
func AllSearchModeTypes() []SearchModeType {
	return []SearchModeType{ModeWildcard, ModeStartsWith, ModeContains, ModeEndsWith, ModeGrep}
}

// IsWildcard reports whether the mode is wildcard.
func (m SearchModeType) IsWildcard() bool { return m == ModeWildcard }

// IsStartsWith reports whether the mode is startsWith.
func (m SearchModeType) IsStartsWith() bool { return m == ModeStartsWith }

// IsContains reports whether the mode is contains.
func (m SearchModeType) IsContains() bool { return m == ModeContains }

// IsEndsWith reports whether the mode is endsWith.
func (m SearchModeType) IsEndsWith() bool { return m == ModeEndsWith }

// IsGrep reports whether the mode is grep.
func (m SearchModeType) IsGrep() bool { return m == ModeGrep }
