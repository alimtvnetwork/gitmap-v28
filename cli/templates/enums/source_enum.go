package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strconv"
	"strings"
)

// SourceType describes where a resolved template came from.
type SourceType int

const (
	SourceNone  SourceType = 0
	SourceUser  SourceType = 1
	SourceEmbed SourceType = 2
)

// IsValid reports whether the value matches a known variant.
func (s SourceType) IsValid() bool {
	switch s {
	case SourceNone, SourceUser, SourceEmbed:
		return true
	default:
		return false
	}
}

// String returns the numeric representation of the value.
func (s SourceType) String() string {
	return strconv.Itoa(int(s))
}

// ParseSourceType parses a raw string into a validated SourceType.
func ParseSourceType(s string) (SourceType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	if n, err := strconv.Atoi(trimmed); err == nil {
		val := SourceType(n)
		if val.IsValid() {
			return val, nil
		}
	}
	for _, val := range AllSourceTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "sourcenone":
		return SourceNone, nil
	case "none":
		return SourceNone, nil
	case "sourceuser":
		return SourceUser, nil
	case "user":
		return SourceUser, nil
	case "sourceembed":
		return SourceEmbed, nil
	case "embed":
		return SourceEmbed, nil
	}
	return 0, apperror.NewValidation("ParseSourceType", "E1000", fmt.Sprintf("invalid sourcetype: %q", s)).WithContext("value", s)
}

// AllSourceTypes returns all valid variants.
func AllSourceTypes() []SourceType {
	return []SourceType{
		SourceNone,
		SourceUser,
		SourceEmbed,
	}
}

// IsNone reports whether the value is SourceNone.
func (s SourceType) IsNone() bool {
	return s == SourceNone
}

// IsUser reports whether the value is SourceUser.
func (s SourceType) IsUser() bool {
	return s == SourceUser
}

// IsEmbed reports whether the value is SourceEmbed.
func (s SourceType) IsEmbed() bool {
	return s == SourceEmbed
}

// HasSource reports whether the value is non-empty and valid.
func (s SourceType) HasSource() bool {
	return s.IsValid()
}
