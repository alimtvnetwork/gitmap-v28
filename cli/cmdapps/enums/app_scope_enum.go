package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// AppScopeType designates system-wide vs user-specific installations.
type AppScopeType string

const (
	// ScopeSystem indicates a system-wide installed application.
	ScopeSystem AppScopeType = "system"
	// ScopeUser indicates a user-scoped application installation.
	ScopeUser AppScopeType = "user"
)

// IsValid reports whether the value matches a known variant.
func (a AppScopeType) IsValid() bool {
	switch a {
	case ScopeSystem, ScopeUser:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (a AppScopeType) String() string {
	return string(a)
}

// ParseAppScopeType parses a raw string into a validated AppScopeType.
func ParseAppScopeType(s string) (AppScopeType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllAppScopeTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "scopesystem":
		return ScopeSystem, nil
	case "system":
		return ScopeSystem, nil
	case "scopeuser":
		return ScopeUser, nil
	case "user":
		return ScopeUser, nil
	}
	return "", apperror.NewValidation("ParseAppScopeType", "E1000", fmt.Sprintf("invalid appscopetype: %q", s)).WithContext("value", s)
}

// AllAppScopeTypes returns all valid variants.
func AllAppScopeTypes() []AppScopeType {
	return []AppScopeType{
		ScopeSystem,
		ScopeUser,
	}
}

// IsSystem reports whether the value is ScopeSystem.
func (a AppScopeType) IsSystem() bool {
	return a == ScopeSystem
}

// IsUser reports whether the value is ScopeUser.
func (a AppScopeType) IsUser() bool {
	return a == ScopeUser
}

// HasScope reports whether the value is non-empty and valid.
func (a AppScopeType) HasScope() bool {
	return a.IsValid()
}
