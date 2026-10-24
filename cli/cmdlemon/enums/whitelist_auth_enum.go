package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// WhitelistAuthType is a strongly-typed enum.
type WhitelistAuthType string

const (
	AuthAllowedType WhitelistAuthType = "ALLOWED"
	AuthDeniedType  WhitelistAuthType = "DENIED"
)

// IsValid reports whether the value matches a known variant.
func (w WhitelistAuthType) IsValid() bool {
	switch w {
	case AuthAllowedType, AuthDeniedType:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (w WhitelistAuthType) String() string {
	return string(w)
}

// ParseWhitelistAuthType parses a raw string into a validated WhitelistAuthType.
func ParseWhitelistAuthType(s string) (WhitelistAuthType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllWhitelistAuthTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "authallowedtype":
		return AuthAllowedType, nil
	case "allowedtype":
		return AuthAllowedType, nil
	case "authdeniedtype":
		return AuthDeniedType, nil
	case "deniedtype":
		return AuthDeniedType, nil
	}
	return "", apperror.NewValidation("ParseWhitelistAuthType", "E1000", fmt.Sprintf("invalid whitelistauthtype: %q", s)).WithContext("value", s)
}

// AllWhitelistAuthTypes returns all valid variants.
func AllWhitelistAuthTypes() []WhitelistAuthType {
	return []WhitelistAuthType{
		AuthAllowedType,
		AuthDeniedType,
	}
}

// IsAllowedType reports whether the value is AuthAllowedType.
func (w WhitelistAuthType) IsAllowedType() bool {
	return w == AuthAllowedType
}

// IsDeniedType reports whether the value is AuthDeniedType.
func (w WhitelistAuthType) IsDeniedType() bool {
	return w == AuthDeniedType
}

// HasAuth reports whether the value is non-empty and valid.
func (w WhitelistAuthType) HasAuth() bool {
	return w.IsValid()
}
