package enums

import (
	"fmt"
	"strings"
)

// HealthCategoryType represents check domains evaluated by health checks.
type HealthCategoryType string

const (
	HealthCategoryGit      HealthCategoryType = "git"
	HealthCategoryPipeline HealthCategoryType = "pipeline"
	HealthCategoryFleet    HealthCategoryType = "fleet"
	HealthCategoryDb       HealthCategoryType = "db"
	HealthCategorySystem   HealthCategoryType = "system"
	HealthCategoryNetwork  HealthCategoryType = "network"
)

// IsValid reports whether the category matches a recognized variant.
func (c HealthCategoryType) IsValid() bool {
	switch c {
	case HealthCategoryGit, HealthCategoryPipeline, HealthCategoryFleet,
		HealthCategoryDb, HealthCategorySystem, HealthCategoryNetwork:
		return true
	default:
		return false
	}
}

// String returns the category identifier string.
func (c HealthCategoryType) String() string {
	return string(c)
}

// IsGit reports whether the category is git.
func (c HealthCategoryType) IsGit() bool {
	return c == HealthCategoryGit
}

// IsPipeline reports whether the category is pipeline.
func (c HealthCategoryType) IsPipeline() bool {
	return c == HealthCategoryPipeline
}

// IsFleet reports whether the category is fleet.
func (c HealthCategoryType) IsFleet() bool {
	return c == HealthCategoryFleet
}

// IsDb reports whether the category is db.
func (c HealthCategoryType) IsDb() bool {
	return c == HealthCategoryDb
}

// IsSystem reports whether the category is system.
func (c HealthCategoryType) IsSystem() bool {
	return c == HealthCategorySystem
}

// IsNetwork reports whether the category is network.
func (c HealthCategoryType) IsNetwork() bool {
	return c == HealthCategoryNetwork
}

// HasCategory reports whether the category is non-empty and valid.
func (c HealthCategoryType) HasCategory() bool {
	return c.IsValid()
}

// ParseHealthCategory parses a category string into a typed enum.
func ParseHealthCategory(s string) (HealthCategoryType, error) {
	val := HealthCategoryType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", fmt.Errorf("invalid health category: %q", s)
}

// ParseHealthCategoryType is an alias for ParseHealthCategory.
func ParseHealthCategoryType(s string) (HealthCategoryType, error) {
	return ParseHealthCategory(s)
}

// AllHealthCategories returns all supported health categories.
func AllHealthCategories() []HealthCategoryType {
	return []HealthCategoryType{
		HealthCategoryGit,
		HealthCategoryPipeline,
		HealthCategoryFleet,
		HealthCategoryDb,
		HealthCategorySystem,
		HealthCategoryNetwork,
	}
}

// AllHealthCategoryTypes is an alias for AllHealthCategories.
func AllHealthCategoryTypes() []HealthCategoryType {
	return AllHealthCategories()
}

// HealthSeverityType represents severity levels for health findings.
type HealthSeverityType string

const (
	HealthSeverityOk    HealthSeverityType = "ok"
	HealthSeverityWarn  HealthSeverityType = "warn"
	HealthSeverityError HealthSeverityType = "error"
	HealthSeverityFatal HealthSeverityType = "fatal"
)

// IsValid reports whether the severity matches a recognized variant.
func (s HealthSeverityType) IsValid() bool {
	switch s {
	case HealthSeverityOk, HealthSeverityWarn, HealthSeverityError, HealthSeverityFatal:
		return true
	default:
		return false
	}
}

// String returns the severity identifier string.
func (s HealthSeverityType) String() string {
	return string(s)
}

// IsOk reports whether the severity is ok.
func (s HealthSeverityType) IsOk() bool {
	return s == HealthSeverityOk
}

// IsWarn reports whether the severity is warn.
func (s HealthSeverityType) IsWarn() bool {
	return s == HealthSeverityWarn
}

// IsError reports whether the severity is error.
func (s HealthSeverityType) IsError() bool {
	return s == HealthSeverityError
}

// IsFatal reports whether the severity is fatal.
func (s HealthSeverityType) IsFatal() bool {
	return s == HealthSeverityFatal
}

// HasSeverity reports whether the severity is valid.
func (s HealthSeverityType) HasSeverity() bool {
	return s.IsValid()
}

// IsUrgent reports whether the severity requires immediate remediation.
func (s HealthSeverityType) IsUrgent() bool {
	return s == HealthSeverityError || s == HealthSeverityFatal
}

// HasDegradation reports whether the severity indicates a warning or failure.
func (s HealthSeverityType) HasDegradation() bool {
	return s.IsValid() && s != HealthSeverityOk
}

// ParseHealthSeverity parses a severity string into a typed enum.
func ParseHealthSeverity(s string) (HealthSeverityType, error) {
	val := HealthSeverityType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", fmt.Errorf("invalid health severity: %q", s)
}

// ParseHealthSeverityType is an alias for ParseHealthSeverity.
func ParseHealthSeverityType(s string) (HealthSeverityType, error) {
	return ParseHealthSeverity(s)
}

// AllHealthSeverities returns all supported health severity levels.
func AllHealthSeverities() []HealthSeverityType {
	return []HealthSeverityType{
		HealthSeverityOk,
		HealthSeverityWarn,
		HealthSeverityError,
		HealthSeverityFatal,
	}
}

// AllHealthSeverityTypes is an alias for AllHealthSeverities.
func AllHealthSeverityTypes() []HealthSeverityType {
	return AllHealthSeverities()
}

// HealthOutputFormatType represents serialization formats for health reports.
type HealthOutputFormatType string

const (
	HealthOutputFormatTerminal HealthOutputFormatType = "terminal"
	HealthOutputFormatJSON     HealthOutputFormatType = "json"
	HealthOutputFormatSummary  HealthOutputFormatType = "summary"
	HealthOutputFormatMarkdown HealthOutputFormatType = "markdown"
)

// IsValid reports whether the output format matches a recognized variant.
func (f HealthOutputFormatType) IsValid() bool {
	switch f {
	case HealthOutputFormatTerminal, HealthOutputFormatJSON,
		HealthOutputFormatSummary, HealthOutputFormatMarkdown:
		return true
	default:
		return false
	}
}

// String returns the output format identifier string.
func (f HealthOutputFormatType) String() string {
	return string(f)
}

// IsTerminal reports whether the output format is terminal.
func (f HealthOutputFormatType) IsTerminal() bool {
	return f == HealthOutputFormatTerminal
}

// IsJSON reports whether the output format is JSON.
func (f HealthOutputFormatType) IsJSON() bool {
	return f == HealthOutputFormatJSON
}

// IsSummary reports whether the output format is summary.
func (f HealthOutputFormatType) IsSummary() bool {
	return f == HealthOutputFormatSummary
}

// IsMarkdown reports whether the output format is markdown.
func (f HealthOutputFormatType) IsMarkdown() bool {
	return f == HealthOutputFormatMarkdown
}

// HasFormat reports whether the output format is valid.
func (f HealthOutputFormatType) HasFormat() bool {
	return f.IsValid()
}

// ParseHealthOutputFormat parses an output format string into a typed enum.
func ParseHealthOutputFormat(s string) (HealthOutputFormatType, error) {
	val := HealthOutputFormatType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", fmt.Errorf("invalid health output format: %q", s)
}

// ParseHealthOutputFormatType is an alias for ParseHealthOutputFormat.
func ParseHealthOutputFormatType(s string) (HealthOutputFormatType, error) {
	return ParseHealthOutputFormat(s)
}

// AllHealthOutputFormats returns all supported health output formats.
func AllHealthOutputFormats() []HealthOutputFormatType {
	return []HealthOutputFormatType{
		HealthOutputFormatTerminal,
		HealthOutputFormatJSON,
		HealthOutputFormatSummary,
		HealthOutputFormatMarkdown,
	}
}

// AllHealthOutputFormatTypes is an alias for AllHealthOutputFormats.
func AllHealthOutputFormatTypes() []HealthOutputFormatType {
	return AllHealthOutputFormats()
}

// HealthCheckStatusType represents individual check result status states.
type HealthCheckStatusType string

const (
	HealthCheckStatusGreen   HealthCheckStatusType = "green"
	HealthCheckStatusYellow  HealthCheckStatusType = "yellow"
	HealthCheckStatusRed     HealthCheckStatusType = "red"
	HealthCheckStatusSkipped HealthCheckStatusType = "skipped"
)

// IsValid reports whether the check status matches a recognized variant.
func (s HealthCheckStatusType) IsValid() bool {
	switch s {
	case HealthCheckStatusGreen, HealthCheckStatusYellow,
		HealthCheckStatusRed, HealthCheckStatusSkipped:
		return true
	default:
		return false
	}
}

// String returns the check status identifier string.
func (s HealthCheckStatusType) String() string {
	return string(s)
}

// IsGreen reports whether the check status is green.
func (s HealthCheckStatusType) IsGreen() bool {
	return s == HealthCheckStatusGreen
}

// IsYellow reports whether the check status is yellow.
func (s HealthCheckStatusType) IsYellow() bool {
	return s == HealthCheckStatusYellow
}

// IsRed reports whether the check status is red.
func (s HealthCheckStatusType) IsRed() bool {
	return s == HealthCheckStatusRed
}

// IsSkipped reports whether the check status is skipped.
func (s HealthCheckStatusType) IsSkipped() bool {
	return s == HealthCheckStatusSkipped
}

// HasPassed reports whether the check resolved as passed/green.
func (s HealthCheckStatusType) HasPassed() bool {
	return s == HealthCheckStatusGreen
}

// HasFailed reports whether the check resolved as failing/red.
func (s HealthCheckStatusType) HasFailed() bool {
	return s == HealthCheckStatusRed
}

// HasStatus reports whether the check status is valid.
func (s HealthCheckStatusType) HasStatus() bool {
	return s.IsValid()
}

// ParseHealthCheckStatus parses a check status string into a typed enum.
func ParseHealthCheckStatus(s string) (HealthCheckStatusType, error) {
	val := HealthCheckStatusType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", fmt.Errorf("invalid health check status: %q", s)
}

// ParseHealthCheckStatusType is an alias for ParseHealthCheckStatus.
func ParseHealthCheckStatusType(s string) (HealthCheckStatusType, error) {
	return ParseHealthCheckStatus(s)
}

// AllHealthCheckStatuses returns all supported health check status states.
func AllHealthCheckStatuses() []HealthCheckStatusType {
	return []HealthCheckStatusType{
		HealthCheckStatusGreen,
		HealthCheckStatusYellow,
		HealthCheckStatusRed,
		HealthCheckStatusSkipped,
	}
}

// AllHealthCheckStatusTypes is an alias for AllHealthCheckStatuses.
func AllHealthCheckStatusTypes() []HealthCheckStatusType {
	return AllHealthCheckStatuses()
}
