package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// ContextKeyType is a context key for macro recursion state.
type ContextKeyType string

const (
	CallStackKey    ContextKeyType = "macroCallStack"
	CurrentMacroKey ContextKeyType = "currentMacroName"
)

// IsValid reports whether the context key matches a known key.
func (k ContextKeyType) IsValid() bool {
	switch k {
	case CallStackKey, CurrentMacroKey:
		return true
	default:
		return false
	}
}

// String returns the string representation of the context key.
func (k ContextKeyType) String() string {
	return string(k)
}

// ParseContextKeyType parses a raw string into a validated ContextKeyType.
func ParseContextKeyType(s string) (ContextKeyType, *apperror.AppError) {
	val := ContextKeyType(strings.TrimSpace(s))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseContextKeyType", "E1000", fmt.Sprintf("invalid macro context key: %q", s)).WithContext("value", s)
}

// AllContextKeyTypes returns all valid ContextKeyType variants.
func AllContextKeyTypes() []ContextKeyType {
	return []ContextKeyType{CallStackKey, CurrentMacroKey}
}

// IsCallStack reports whether the key is the call stack key.
func (k ContextKeyType) IsCallStack() bool { return k == CallStackKey }

// IsCurrentMacro reports whether the key is the current macro key.
func (k ContextKeyType) IsCurrentMacro() bool { return k == CurrentMacroKey }
