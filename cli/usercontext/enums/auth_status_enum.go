package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// AuthStatusType defines the authorization state for GitHub CLI and provider tools.
type AuthStatusType string

const (
	StatusAuthorized      AuthStatusType = "authorized"
	StatusUnauthenticated AuthStatusType = "unauthenticated"
	StatusToolMissing     AuthStatusType = "tool_missing"
)

// IsValid reports whether the value matches a known variant.
func (a AuthStatusType) IsValid() bool {
	switch a {
	case StatusAuthorized, StatusUnauthenticated, StatusToolMissing:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (a AuthStatusType) String() string {
	return string(a)
}

// ParseAuthStatusType parses a raw string into a validated AuthStatusType.
func ParseAuthStatusType(s string) (AuthStatusType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllAuthStatusTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "statusauthorized":
		return StatusAuthorized, nil
	case "authorized":
		return StatusAuthorized, nil
	case "statusunauthenticated":
		return StatusUnauthenticated, nil
	case "unauthenticated":
		return StatusUnauthenticated, nil
	case "statustoolmissing":
		return StatusToolMissing, nil
	case "toolmissing":
		return StatusToolMissing, nil
	}
	return "", apperror.NewValidation("ParseAuthStatusType", "E1000", fmt.Sprintf("invalid authstatustype: %q", s)).WithContext("value", s)
}

// AllAuthStatusTypes returns all valid variants.
func AllAuthStatusTypes() []AuthStatusType {
	return []AuthStatusType{
		StatusAuthorized,
		StatusUnauthenticated,
		StatusToolMissing,
	}
}

// IsAuthorized reports whether the value is StatusAuthorized.
func (a AuthStatusType) IsAuthorized() bool {
	return a == StatusAuthorized
}

// IsUnauthenticated reports whether the value is StatusUnauthenticated.
func (a AuthStatusType) IsUnauthenticated() bool {
	return a == StatusUnauthenticated
}

// IsToolMissing reports whether the value is StatusToolMissing.
func (a AuthStatusType) IsToolMissing() bool {
	return a == StatusToolMissing
}

// HasStatus reports whether the value is non-empty and valid.
func (a AuthStatusType) HasStatus() bool {
	return a.IsValid()
}
