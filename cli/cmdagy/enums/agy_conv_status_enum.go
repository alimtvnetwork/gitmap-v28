package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// AgyConvStatusType represents the execution status of an active conversation.
type AgyConvStatusType string

const (
	AgyConvStatusIdle    AgyConvStatusType = "idle"
	AgyConvStatusRunning AgyConvStatusType = "running"
	AgyConvStatusUnknown AgyConvStatusType = "unknown"
)

// IsValid reports whether the conversation status matches a known variant.
func (s AgyConvStatusType) IsValid() bool {
	switch s {
	case AgyConvStatusIdle, AgyConvStatusRunning, AgyConvStatusUnknown:
		return true
	default:
		return false
	}
}

// String returns the string representation of the conversation status.
func (s AgyConvStatusType) String() string {
	return string(s)
}

// ParseAgyConvStatusType parses a raw string into a validated AgyConvStatusType.
func ParseAgyConvStatusType(s string) (AgyConvStatusType, *apperror.AppError) {
	val := AgyConvStatusType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseAgyConvStatusType", "E1000", fmt.Sprintf("invalid agy conversation status: %q", s)).WithContext("value", s)
}

// AllAgyConvStatusTypes returns all valid AgyConvStatusType variants.
func AllAgyConvStatusTypes() []AgyConvStatusType {
	return []AgyConvStatusType{
		AgyConvStatusIdle,
		AgyConvStatusRunning,
		AgyConvStatusUnknown,
	}
}

// IsIdle reports whether the conversation status is idle.
func (s AgyConvStatusType) IsIdle() bool {
	return s == AgyConvStatusIdle
}

// IsRunning reports whether the conversation status is running.
func (s AgyConvStatusType) IsRunning() bool {
	return s == AgyConvStatusRunning
}

// IsUnknown reports whether the conversation status is unknown.
func (s AgyConvStatusType) IsUnknown() bool {
	return s == AgyConvStatusUnknown
}

// HasConvStatus reports whether the conversation status is non-empty and valid.
func (s AgyConvStatusType) HasConvStatus() bool {
	return s.IsValid()
}
