package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// PromptStatusType defines running or enqueued prompt status.
type PromptStatusType string

const (
	// PromptStatusRunning indicates an active executing prompt.
	PromptStatusRunning PromptStatusType = "running"
	// PromptStatusEnqueued indicates a pending queued prompt.
	PromptStatusEnqueued PromptStatusType = "enqueued"
)

// IsValid reports whether the prompt status is a recognized enumeration value.
func (s PromptStatusType) IsValid() bool {
	switch s {
	case PromptStatusRunning, PromptStatusEnqueued:
		return true
	default:
		return false
	}
}

// String returns the string representation of the prompt status.
func (s PromptStatusType) String() string {
	return string(s)
}

// ParsePromptStatusType parses a raw string into a validated PromptStatusType.
func ParsePromptStatusType(s string) (PromptStatusType, *apperror.AppError) {
	val := PromptStatusType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParsePromptStatusType", "E1000", fmt.Sprintf("invalid prompt status: %q", s)).WithContext("value", s)
}

// AllPromptStatusTypes returns all valid PromptStatusType variants.
func AllPromptStatusTypes() []PromptStatusType {
	return []PromptStatusType{
		PromptStatusRunning,
		PromptStatusEnqueued,
	}
}

// IsRunning reports whether the prompt status is running.
func (s PromptStatusType) IsRunning() bool {
	return s == PromptStatusRunning
}

// IsEnqueued reports whether the prompt status is enqueued.
func (s PromptStatusType) IsEnqueued() bool {
	return s == PromptStatusEnqueued
}
