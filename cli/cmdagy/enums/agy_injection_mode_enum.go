package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// AgyInjectionModeType indicates whether injection was routed via CLI, IDE, or skipped.
type AgyInjectionModeType string

const (
	AgyInjectionModeNone   AgyInjectionModeType = "none"
	AgyInjectionModeCLI    AgyInjectionModeType = "cli"
	AgyInjectionModeIDE    AgyInjectionModeType = "ide"
	AgyInjectionModeQueued AgyInjectionModeType = "queued"
)

// IsValid reports whether the injection mode matches a known variant.
func (m AgyInjectionModeType) IsValid() bool {
	switch m {
	case AgyInjectionModeNone, AgyInjectionModeCLI, AgyInjectionModeIDE, AgyInjectionModeQueued:
		return true
	default:
		return false
	}
}

// String returns the string representation of the injection mode.
func (m AgyInjectionModeType) String() string {
	return string(m)
}

// ParseAgyInjectionModeType parses a raw string into a validated AgyInjectionModeType.
func ParseAgyInjectionModeType(s string) (AgyInjectionModeType, *apperror.AppError) {
	val := AgyInjectionModeType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseAgyInjectionModeType", "E1000", fmt.Sprintf("invalid agy injection mode: %q", s)).WithContext("value", s)
}

// AllAgyInjectionModeTypes returns all valid AgyInjectionModeType variants.
func AllAgyInjectionModeTypes() []AgyInjectionModeType {
	return []AgyInjectionModeType{
		AgyInjectionModeNone,
		AgyInjectionModeCLI,
		AgyInjectionModeIDE,
		AgyInjectionModeQueued,
	}
}

// IsNone reports whether the injection mode is none.
func (m AgyInjectionModeType) IsNone() bool {
	return m == AgyInjectionModeNone
}

// IsCLI reports whether the injection mode is CLI.
func (m AgyInjectionModeType) IsCLI() bool {
	return m == AgyInjectionModeCLI
}

// IsIDE reports whether the injection mode is IDE.
func (m AgyInjectionModeType) IsIDE() bool {
	return m == AgyInjectionModeIDE
}

// IsQueued reports whether the injection mode is queued.
func (m AgyInjectionModeType) IsQueued() bool {
	return m == AgyInjectionModeQueued
}

// HasInjectionMode reports whether the injection mode is non-empty and valid.
func (m AgyInjectionModeType) HasInjectionMode() bool {
	return m.IsValid()
}
