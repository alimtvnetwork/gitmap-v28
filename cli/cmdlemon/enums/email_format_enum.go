package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// EmailFormatType represents the encoding of an incoming email command.
type EmailFormatType string

const (
	EmailFormatTypeText    EmailFormatType = "text"
	EmailFormatTypeJSON    EmailFormatType = "json"
	EmailFormatTypeYAML    EmailFormatType = "yaml"
	EmailFormatTypeUnknown EmailFormatType = "unknown"
)

// IsValid reports whether the value matches a known variant.
func (e EmailFormatType) IsValid() bool {
	switch e {
	case EmailFormatTypeText, EmailFormatTypeJSON, EmailFormatTypeYAML, EmailFormatTypeUnknown:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (e EmailFormatType) String() string {
	return string(e)
}

// ParseEmailFormatType parses a raw string into a validated EmailFormatType.
func ParseEmailFormatType(s string) (EmailFormatType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllEmailFormatTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "emailformattypetext":
		return EmailFormatTypeText, nil
	case "text":
		return EmailFormatTypeText, nil
	case "emailformattypejson":
		return EmailFormatTypeJSON, nil
	case "json":
		return EmailFormatTypeJSON, nil
	case "emailformattypeyaml":
		return EmailFormatTypeYAML, nil
	case "yaml":
		return EmailFormatTypeYAML, nil
	case "emailformattypeunknown":
		return EmailFormatTypeUnknown, nil
	case "unknown":
		return EmailFormatTypeUnknown, nil
	}
	return "", apperror.NewValidation("ParseEmailFormatType", "E1000", fmt.Sprintf("invalid emailformattype: %q", s)).WithContext("value", s)
}

// AllEmailFormatTypes returns all valid variants.
func AllEmailFormatTypes() []EmailFormatType {
	return []EmailFormatType{
		EmailFormatTypeText,
		EmailFormatTypeJSON,
		EmailFormatTypeYAML,
		EmailFormatTypeUnknown,
	}
}

// IsText reports whether the value is EmailFormatTypeText.
func (e EmailFormatType) IsText() bool {
	return e == EmailFormatTypeText
}

// IsJSON reports whether the value is EmailFormatTypeJSON.
func (e EmailFormatType) IsJSON() bool {
	return e == EmailFormatTypeJSON
}

// IsYAML reports whether the value is EmailFormatTypeYAML.
func (e EmailFormatType) IsYAML() bool {
	return e == EmailFormatTypeYAML
}

// IsUnknown reports whether the value is EmailFormatTypeUnknown.
func (e EmailFormatType) IsUnknown() bool {
	return e == EmailFormatTypeUnknown
}

// HasFormat reports whether the value is non-empty and valid.
func (e EmailFormatType) HasFormat() bool {
	return e.IsValid()
}
