package enums

import (
	"fmt"
	"strings"
)

// StaleOutputFormatType specifies the output formatting mode for stale repo reporting.
type StaleOutputFormatType string

const (
	StaleOutputFormatTable StaleOutputFormatType = "table"
	StaleOutputFormatJSON  StaleOutputFormatType = "json"
	StaleOutputFormatCSV   StaleOutputFormatType = "csv"
)

// IsValid reports whether the output format is a recognized enum value.
func (f StaleOutputFormatType) IsValid() bool {
	switch f {
	case StaleOutputFormatTable, StaleOutputFormatJSON, StaleOutputFormatCSV:
		return true
	default:
		return false
	}
}

// String returns the string representation of the output format.
func (f StaleOutputFormatType) String() string {
	return string(f)
}

// ParseStaleOutputFormatType parses a raw string into a StaleOutputFormatType.
func ParseStaleOutputFormatType(s string) (StaleOutputFormatType, error) {
	norm := strings.ToLower(strings.TrimSpace(s))
	if norm == "" {
		return StaleOutputFormatTable, nil
	}

	val := StaleOutputFormatType(norm)
	if val.IsValid() {
		return val, nil
	}

	return "", fmt.Errorf("invalid stale output format: %q", s)
}

// AllStaleOutputFormatTypes returns all supported StaleOutputFormatType values.
func AllStaleOutputFormatTypes() []StaleOutputFormatType {
	return []StaleOutputFormatType{
		StaleOutputFormatTable,
		StaleOutputFormatJSON,
		StaleOutputFormatCSV,
	}
}

// IsTable reports whether the format is tabular text.
func (f StaleOutputFormatType) IsTable() bool {
	return f == StaleOutputFormatTable
}

// IsJSON reports whether the format is JSON.
func (f StaleOutputFormatType) IsJSON() bool {
	return f == StaleOutputFormatJSON
}

// IsCSV reports whether the format is CSV.
func (f StaleOutputFormatType) IsCSV() bool {
	return f == StaleOutputFormatCSV
}

// HasFormat reports whether the format is valid and recognized.
func (f StaleOutputFormatType) HasFormat() bool {
	return f.IsValid()
}

// StaleActionType specifies the action mode for cmdstale execution.
type StaleActionType string

const (
	StaleActionInspect StaleActionType = "inspect"
	StaleActionArchive StaleActionType = "archive"
	StaleActionPreview StaleActionType = "preview"
)

// IsValid reports whether the stale action is a recognized enum value.
func (a StaleActionType) IsValid() bool {
	switch a {
	case StaleActionInspect, StaleActionArchive, StaleActionPreview:
		return true
	default:
		return false
	}
}

// String returns the string representation of the stale action.
func (a StaleActionType) String() string {
	return string(a)
}

// ParseStaleActionType parses a raw string into a StaleActionType.
func ParseStaleActionType(s string) (StaleActionType, error) {
	val := StaleActionType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", fmt.Errorf("invalid stale action: %q", s)
}

// AllStaleActionTypes returns all supported StaleActionType values.
func AllStaleActionTypes() []StaleActionType {
	return []StaleActionType{
		StaleActionInspect,
		StaleActionArchive,
		StaleActionPreview,
	}
}

// IsInspect reports whether the action is inspection only.
func (a StaleActionType) IsInspect() bool {
	return a == StaleActionInspect
}

// IsArchive reports whether the action moves stale repos to archive.
func (a StaleActionType) IsArchive() bool {
	return a == StaleActionArchive
}

// IsPreview reports whether the action is a dry-run preview.
func (a StaleActionType) IsPreview() bool {
	return a == StaleActionPreview
}

// HasAction reports whether the action is valid and recognized.
func (a StaleActionType) HasAction() bool {
	return a.IsValid()
}

// StaleFilterType specifies criteria for filtering repositories by staleness.
type StaleFilterType string

const (
	StaleFilterAll      StaleFilterType = "all"
	StaleFilterStale    StaleFilterType = "stale"
	StaleFilterActive   StaleFilterType = "active"
	StaleFilterArchived StaleFilterType = "archived"
)

// IsValid reports whether the stale filter is a recognized enum value.
func (f StaleFilterType) IsValid() bool {
	switch f {
	case StaleFilterAll, StaleFilterStale, StaleFilterActive, StaleFilterArchived:
		return true
	default:
		return false
	}
}

// String returns the string representation of the stale filter.
func (f StaleFilterType) String() string {
	return string(f)
}

// ParseStaleFilterType parses a raw string into a StaleFilterType.
func ParseStaleFilterType(s string) (StaleFilterType, error) {
	val := StaleFilterType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", fmt.Errorf("invalid stale filter: %q", s)
}

// AllStaleFilterTypes returns all supported StaleFilterType values.
func AllStaleFilterTypes() []StaleFilterType {
	return []StaleFilterType{
		StaleFilterAll,
		StaleFilterStale,
		StaleFilterActive,
		StaleFilterArchived,
	}
}

// IsAll reports whether all repositories are selected.
func (f StaleFilterType) IsAll() bool {
	return f == StaleFilterAll
}

// IsStale reports whether only stale repositories are selected.
func (f StaleFilterType) IsStale() bool {
	return f == StaleFilterStale
}

// IsActive reports whether active repositories are selected.
func (f StaleFilterType) IsActive() bool {
	return f == StaleFilterActive
}

// IsArchived reports whether archived repositories are selected.
func (f StaleFilterType) IsArchived() bool {
	return f == StaleFilterArchived
}

// HasFilter reports whether the filter is valid and recognized.
func (f StaleFilterType) HasFilter() bool {
	return f.IsValid()
}
