package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// LevelType enumerates the supported severities.
type LevelType string

const (
	LevelDebug LevelType = "debug"
	LevelInfo  LevelType = "info"
	LevelWarn  LevelType = "warn"
	LevelError LevelType = "error"
)

// IsValid reports whether the value matches a known variant.
func (l LevelType) IsValid() bool {
	switch l {
	case LevelDebug, LevelInfo, LevelWarn, LevelError:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (l LevelType) String() string {
	return string(l)
}

// ParseLevelType parses a raw string into a validated LevelType.
func ParseLevelType(s string) (LevelType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllLevelTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "leveldebug":
		return LevelDebug, nil
	case "debug":
		return LevelDebug, nil
	case "levelinfo":
		return LevelInfo, nil
	case "info":
		return LevelInfo, nil
	case "levelwarn":
		return LevelWarn, nil
	case "warn":
		return LevelWarn, nil
	case "levelerror":
		return LevelError, nil
	case "error":
		return LevelError, nil
	}
	return "", apperror.NewValidation("ParseLevelType", "E1000", fmt.Sprintf("invalid leveltype: %q", s)).WithContext("value", s)
}

// AllLevelTypes returns all valid variants.
func AllLevelTypes() []LevelType {
	return []LevelType{
		LevelDebug,
		LevelInfo,
		LevelWarn,
		LevelError,
	}
}

// IsDebug reports whether the value is LevelDebug.
func (l LevelType) IsDebug() bool {
	return l == LevelDebug
}

// IsInfo reports whether the value is LevelInfo.
func (l LevelType) IsInfo() bool {
	return l == LevelInfo
}

// IsWarn reports whether the value is LevelWarn.
func (l LevelType) IsWarn() bool {
	return l == LevelWarn
}

// IsError reports whether the value is LevelError.
func (l LevelType) IsError() bool {
	return l == LevelError
}

// HasLevel reports whether the value is non-empty and valid.
func (l LevelType) HasLevel() bool {
	return l.IsValid()
}
