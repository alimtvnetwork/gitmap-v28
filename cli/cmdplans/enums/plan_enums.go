package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// PlanStatusType represents the lifecycle status of an architectural plan.
type PlanStatusType string

const (
	PlanStatusPending    PlanStatusType = "pending"
	PlanStatusActive     PlanStatusType = "active"
	PlanStatusInProgress PlanStatusType = "in_progress"
	PlanStatusCompleted  PlanStatusType = "completed"
	PlanStatusArchived   PlanStatusType = "archived"
	PlanStatusFailed     PlanStatusType = "failed"
)

var planStatusLookup = map[string]PlanStatusType{
	"pending":     PlanStatusPending,
	"queued":      PlanStatusPending,
	"active":      PlanStatusActive,
	"in_progress": PlanStatusInProgress,
	"in progress": PlanStatusInProgress,
	"inprogress":  PlanStatusInProgress,
	"completed":   PlanStatusCompleted,
	"complete":    PlanStatusCompleted,
	"done":        PlanStatusCompleted,
	"archived":    PlanStatusArchived,
	"archive":     PlanStatusArchived,
	"failed":      PlanStatusFailed,
	"fail":        PlanStatusFailed,
	"error":       PlanStatusFailed,
}

// IsValid reports whether the plan status is recognized.
func (s PlanStatusType) IsValid() bool {
	switch s {
	case PlanStatusPending, PlanStatusActive, PlanStatusInProgress,
		PlanStatusCompleted, PlanStatusArchived, PlanStatusFailed:
		return true
	default:
		return false
	}
}

// String returns the string representation.
func (s PlanStatusType) String() string {
	return string(s)
}

// AllPlanStatuses returns all recognized plan statuses.
func AllPlanStatuses() []PlanStatusType {
	return []PlanStatusType{
		PlanStatusPending,
		PlanStatusActive,
		PlanStatusInProgress,
		PlanStatusCompleted,
		PlanStatusArchived,
		PlanStatusFailed,
	}
}

// ParsePlanStatus parses a string into a PlanStatusType.
func ParsePlanStatus(s string) (PlanStatusType, *apperror.AppError) {
	normalized := strings.ToLower(strings.TrimSpace(s))
	if status, hasMatch := planStatusLookup[normalized]; hasMatch {
		return status, nil
	}

	return "", apperror.NewValidation("ParsePlanStatus", "E1000", fmt.Sprintf("invalid plan status: %q", s)).WithContext("value", s)
}

// IsPending reports whether the status is pending.
func (s PlanStatusType) IsPending() bool {
	return s == PlanStatusPending
}

// IsActive reports whether the status is active.
func (s PlanStatusType) IsActive() bool {
	return s == PlanStatusActive
}

// IsInProgress reports whether the status is in progress.
func (s PlanStatusType) IsInProgress() bool {
	return s == PlanStatusInProgress
}

// IsCompleted reports whether the status is completed.
func (s PlanStatusType) IsCompleted() bool {
	return s == PlanStatusCompleted
}

// IsArchived reports whether the status is archived.
func (s PlanStatusType) IsArchived() bool {
	return s == PlanStatusArchived
}

// IsFailed reports whether the status is failed.
func (s PlanStatusType) IsFailed() bool {
	return s == PlanStatusFailed
}

// PlansFilterModeType represents query filtering for plan inventory.
type PlansFilterModeType string

const (
	PlansFilterModeAll       PlansFilterModeType = "all"
	PlansFilterModeActive    PlansFilterModeType = "active"
	PlansFilterModePending   PlansFilterModeType = "pending"
	PlansFilterModeCompleted PlansFilterModeType = "completed"
)

var plansFilterModeLookup = map[string]PlansFilterModeType{
	"all":       PlansFilterModeAll,
	"active":    PlansFilterModeActive,
	"pending":   PlansFilterModePending,
	"completed": PlansFilterModeCompleted,
}

// IsValid reports whether the filter mode is recognized.
func (m PlansFilterModeType) IsValid() bool {
	switch m {
	case PlansFilterModeAll, PlansFilterModeActive, PlansFilterModePending, PlansFilterModeCompleted:
		return true
	default:
		return false
	}
}

// String returns the string representation.
func (m PlansFilterModeType) String() string {
	return string(m)
}

// AllPlansFilterModes returns all recognized filter modes.
func AllPlansFilterModes() []PlansFilterModeType {
	return []PlansFilterModeType{
		PlansFilterModeAll,
		PlansFilterModeActive,
		PlansFilterModePending,
		PlansFilterModeCompleted,
	}
}

// ParsePlansFilterMode parses a string into a PlansFilterModeType.
func ParsePlansFilterMode(s string) (PlansFilterModeType, *apperror.AppError) {
	normalized := strings.ToLower(strings.TrimSpace(s))
	if mode, hasMatch := plansFilterModeLookup[normalized]; hasMatch {
		return mode, nil
	}

	return "", apperror.NewValidation("ParsePlansFilterMode", "E1000", fmt.Sprintf("invalid plans filter mode: %q", s)).WithContext("value", s)
}

// IsAll reports whether the filter mode includes all plans.
func (m PlansFilterModeType) IsAll() bool {
	return m == PlansFilterModeAll
}

// IsActive reports whether the filter mode selects active plans.
func (m PlansFilterModeType) IsActive() bool {
	return m == PlansFilterModeActive
}

// IsPending reports whether the filter mode selects pending plans.
func (m PlansFilterModeType) IsPending() bool {
	return m == PlansFilterModePending
}

// IsCompleted reports whether the filter mode selects completed plans.
func (m PlansFilterModeType) IsCompleted() bool {
	return m == PlansFilterModeCompleted
}

// PlansOutputFormatType represents output rendering formats.
type PlansOutputFormatType string

const (
	PlansOutputFormatTable    PlansOutputFormatType = "table"
	PlansOutputFormatJSON     PlansOutputFormatType = "json"
	PlansOutputFormatMarkdown PlansOutputFormatType = "markdown"
	PlansOutputFormatCSV      PlansOutputFormatType = "csv"
	PlansOutputFormatYAML     PlansOutputFormatType = "yaml"
)

var plansOutputFormatLookup = map[string]PlansOutputFormatType{
	"table":    PlansOutputFormatTable,
	"json":     PlansOutputFormatJSON,
	"markdown": PlansOutputFormatMarkdown,
	"md":       PlansOutputFormatMarkdown,
	"csv":      PlansOutputFormatCSV,
	"yaml":     PlansOutputFormatYAML,
	"yml":      PlansOutputFormatYAML,
}

// IsValid reports whether the output format is recognized.
func (f PlansOutputFormatType) IsValid() bool {
	switch f {
	case PlansOutputFormatTable, PlansOutputFormatJSON, PlansOutputFormatMarkdown,
		PlansOutputFormatCSV, PlansOutputFormatYAML:
		return true
	default:
		return false
	}
}

// String returns the string representation.
func (f PlansOutputFormatType) String() string {
	return string(f)
}

// AllPlansOutputFormats returns all recognized output formats.
func AllPlansOutputFormats() []PlansOutputFormatType {
	return []PlansOutputFormatType{
		PlansOutputFormatTable,
		PlansOutputFormatJSON,
		PlansOutputFormatMarkdown,
		PlansOutputFormatCSV,
		PlansOutputFormatYAML,
	}
}

// ParsePlansOutputFormat parses a string into a PlansOutputFormatType.
func ParsePlansOutputFormat(s string) (PlansOutputFormatType, *apperror.AppError) {
	normalized := strings.ToLower(strings.TrimSpace(s))
	if format, hasMatch := plansOutputFormatLookup[normalized]; hasMatch {
		return format, nil
	}

	return "", apperror.NewValidation("ParsePlansOutputFormat", "E1000", fmt.Sprintf("invalid plans output format: %q", s)).WithContext("value", s)
}

// IsTable reports whether the format is table.
func (f PlansOutputFormatType) IsTable() bool {
	return f == PlansOutputFormatTable
}

// IsJSON reports whether the format is JSON.
func (f PlansOutputFormatType) IsJSON() bool {
	return f == PlansOutputFormatJSON
}

// IsMarkdown reports whether the format is Markdown.
func (f PlansOutputFormatType) IsMarkdown() bool {
	return f == PlansOutputFormatMarkdown
}

// IsCSV reports whether the format is CSV.
func (f PlansOutputFormatType) IsCSV() bool {
	return f == PlansOutputFormatCSV
}

// IsYAML reports whether the format is YAML.
func (f PlansOutputFormatType) IsYAML() bool {
	return f == PlansOutputFormatYAML
}

// PlanMetaFileType represents known metadata files to exclude or identify.
type PlanMetaFileType string

const (
	PlanMetaFileReadme     PlanMetaFileType = "readme.md"
	PlanMetaFileSummary    PlanMetaFileType = "summary.md"
	PlanMetaFileChecklist  PlanMetaFileType = "checklist.md"
	PlanMetaFileIndex      PlanMetaFileType = "01-index.md"
	PlanMetaFileLedger     PlanMetaFileType = "ledger.md"
	PlanMetaFileWhatToRead PlanMetaFileType = "what-to-read.md"
)

var planMetaFileLookup = map[string]PlanMetaFileType{
	"readme.md":       PlanMetaFileReadme,
	"summary.md":      PlanMetaFileSummary,
	"checklist.md":    PlanMetaFileChecklist,
	"01-index.md":     PlanMetaFileIndex,
	"ledger.md":       PlanMetaFileLedger,
	"what-to-read.md": PlanMetaFileWhatToRead,
}

// IsValid reports whether the meta file type is recognized.
func (m PlanMetaFileType) IsValid() bool {
	switch m {
	case PlanMetaFileReadme, PlanMetaFileSummary, PlanMetaFileChecklist,
		PlanMetaFileIndex, PlanMetaFileLedger, PlanMetaFileWhatToRead:
		return true
	default:
		return false
	}
}

// String returns the filename representation.
func (m PlanMetaFileType) String() string {
	return string(m)
}

// AllPlanMetaFiles returns all recognized metadata files.
func AllPlanMetaFiles() []PlanMetaFileType {
	return []PlanMetaFileType{
		PlanMetaFileReadme,
		PlanMetaFileSummary,
		PlanMetaFileChecklist,
		PlanMetaFileIndex,
		PlanMetaFileLedger,
		PlanMetaFileWhatToRead,
	}
}

// ParsePlanMetaFile parses a filename into a PlanMetaFileType.
func ParsePlanMetaFile(s string) (PlanMetaFileType, *apperror.AppError) {
	normalized := strings.ToLower(strings.TrimSpace(s))
	if metaFile, hasMatch := planMetaFileLookup[normalized]; hasMatch {
		return metaFile, nil
	}

	return "", apperror.NewValidation("ParsePlanMetaFile", "E1000", fmt.Sprintf("invalid plan meta file: %q", s)).WithContext("value", s)
}

// IsReadme reports whether the file is readme.md.
func (m PlanMetaFileType) IsReadme() bool {
	return m == PlanMetaFileReadme
}

// IsSummary reports whether the file is summary.md.
func (m PlanMetaFileType) IsSummary() bool {
	return m == PlanMetaFileSummary
}

// IsChecklist reports whether the file is checklist.md.
func (m PlanMetaFileType) IsChecklist() bool {
	return m == PlanMetaFileChecklist
}

// IsIndex reports whether the file is 01-index.md.
func (m PlanMetaFileType) IsIndex() bool {
	return m == PlanMetaFileIndex
}

// IsLedger reports whether the file is ledger.md.
func (m PlanMetaFileType) IsLedger() bool {
	return m == PlanMetaFileLedger
}

// IsWhatToRead reports whether the file is what-to-read.md.
func (m PlanMetaFileType) IsWhatToRead() bool {
	return m == PlanMetaFileWhatToRead
}

// IsMetaFileName reports whether the given filename matches any known plan meta file.
func IsMetaFileName(fileName string) bool {
	_, err := ParsePlanMetaFile(fileName)

	return err == nil
}
