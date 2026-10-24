package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// ArchiveInstallStrategyType represents detected deployment mechanism.
type ArchiveInstallStrategyType string

const (
	StrategyBinaryApp ArchiveInstallStrategyType = "binary_app"
	StrategyScript    ArchiveInstallStrategyType = "install_script"
	StrategySource    ArchiveInstallStrategyType = "source_build"
	StrategySingleGz  ArchiveInstallStrategyType = "single_gzip"
	StrategyUnknown   ArchiveInstallStrategyType = "unknown"
)

// IsValid reports whether the strategy matches a known archive install strategy.
func (s ArchiveInstallStrategyType) IsValid() bool {
	switch s {
	case StrategyBinaryApp, StrategyScript, StrategySource,
		StrategySingleGz, StrategyUnknown:
		return true
	default:
		return false
	}
}

// String returns the string representation of the strategy.
func (s ArchiveInstallStrategyType) String() string {
	return string(s)
}

// ParseArchiveInstallStrategyType parses a raw string into a validated ArchiveInstallStrategyType.
func ParseArchiveInstallStrategyType(s string) (ArchiveInstallStrategyType, *apperror.AppError) {
	val := ArchiveInstallStrategyType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseArchiveInstallStrategyType", "E1000", fmt.Sprintf("invalid archive install strategy: %q", s)).WithContext("value", s)
}

// AllArchiveInstallStrategyTypes returns all valid ArchiveInstallStrategyType variants.
func AllArchiveInstallStrategyTypes() []ArchiveInstallStrategyType {
	return []ArchiveInstallStrategyType{
		StrategyBinaryApp,
		StrategyScript,
		StrategySource,
		StrategySingleGz,
		StrategyUnknown,
	}
}

// IsBinaryApp reports whether the strategy is a binary application.
func (s ArchiveInstallStrategyType) IsBinaryApp() bool {
	return s == StrategyBinaryApp
}

// IsScript reports whether the strategy is an install script.
func (s ArchiveInstallStrategyType) IsScript() bool {
	return s == StrategyScript
}

// IsSource reports whether the strategy is a source build.
func (s ArchiveInstallStrategyType) IsSource() bool {
	return s == StrategySource
}

// IsSingleGz reports whether the strategy is a single gzip archive.
func (s ArchiveInstallStrategyType) IsSingleGz() bool {
	return s == StrategySingleGz
}

// IsUnknown reports whether the strategy is unknown.
func (s ArchiveInstallStrategyType) IsUnknown() bool {
	return s == StrategyUnknown
}
