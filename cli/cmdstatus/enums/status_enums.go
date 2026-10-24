package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// StatusFilterModeType represents filtering criteria for workspace repositories.
type StatusFilterModeType string

const (
	StatusFilterAll    StatusFilterModeType = "all"
	StatusFilterDirty  StatusFilterModeType = "dirty"
	StatusFilterAhead  StatusFilterModeType = "ahead"
	StatusFilterBehind StatusFilterModeType = "behind"
)

// IsValid reports whether the filter mode matches a recognized variant.
func (f StatusFilterModeType) IsValid() bool {
	switch f {
	case StatusFilterAll, StatusFilterDirty, StatusFilterAhead, StatusFilterBehind:
		return true
	default:
		return false
	}
}

// String returns the filter mode identifier.
func (f StatusFilterModeType) String() string {
	return string(f)
}

// IsAll reports whether the filter selects all repositories.
func (f StatusFilterModeType) IsAll() bool {
	return f == StatusFilterAll
}

// IsDirty reports whether the filter selects modified repositories.
func (f StatusFilterModeType) IsDirty() bool {
	return f == StatusFilterDirty
}

// IsAhead reports whether the filter selects ahead repositories.
func (f StatusFilterModeType) IsAhead() bool {
	return f == StatusFilterAhead
}

// IsBehind reports whether the filter selects behind repositories.
func (f StatusFilterModeType) IsBehind() bool {
	return f == StatusFilterBehind
}

// ParseStatusFilterMode parses a filter mode string into a typed enum.
func ParseStatusFilterMode(s string) (StatusFilterModeType, *apperror.AppError) {
	val := StatusFilterModeType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseStatusFilterMode", "E1000", fmt.Sprintf("invalid status filter mode: %q", s)).WithContext("value", s)
}

// ParseStatusFilterModeType is an alias for ParseStatusFilterMode.
func ParseStatusFilterModeType(s string) (StatusFilterModeType, *apperror.AppError) {
	return ParseStatusFilterMode(s)
}

// AllStatusFilterModes returns all supported status filter modes.
func AllStatusFilterModes() []StatusFilterModeType {
	return []StatusFilterModeType{
		StatusFilterAll,
		StatusFilterDirty,
		StatusFilterAhead,
		StatusFilterBehind,
	}
}

// AllStatusFilterModeTypes is an alias for AllStatusFilterModes.
func AllStatusFilterModeTypes() []StatusFilterModeType {
	return AllStatusFilterModes()
}

// StatusOutputModeType represents presentation format options for status.
type StatusOutputModeType string

const (
	StatusOutputTable   StatusOutputModeType = "table"
	StatusOutputJSON    StatusOutputModeType = "json"
	StatusOutputCompact StatusOutputModeType = "compact"
)

// IsValid reports whether the output mode matches a recognized variant.
func (o StatusOutputModeType) IsValid() bool {
	switch o {
	case StatusOutputTable, StatusOutputJSON, StatusOutputCompact:
		return true
	default:
		return false
	}
}

// String returns the output mode identifier.
func (o StatusOutputModeType) String() string {
	return string(o)
}

// IsTable reports whether table rendering is selected.
func (o StatusOutputModeType) IsTable() bool {
	return o == StatusOutputTable
}

// IsJSON reports whether JSON output is selected.
func (o StatusOutputModeType) IsJSON() bool {
	return o == StatusOutputJSON
}

// IsCompact reports whether single-line compact output is selected.
func (o StatusOutputModeType) IsCompact() bool {
	return o == StatusOutputCompact
}

// ParseStatusOutputMode parses an output mode string into a typed enum.
func ParseStatusOutputMode(s string) (StatusOutputModeType, *apperror.AppError) {
	val := StatusOutputModeType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseStatusOutputMode", "E1000", fmt.Sprintf("invalid status output mode: %q", s)).WithContext("value", s)
}

// ParseStatusOutputModeType is an alias for ParseStatusOutputMode.
func ParseStatusOutputModeType(s string) (StatusOutputModeType, *apperror.AppError) {
	return ParseStatusOutputMode(s)
}

// AllStatusOutputModes returns all supported status output modes.
func AllStatusOutputModes() []StatusOutputModeType {
	return []StatusOutputModeType{
		StatusOutputTable,
		StatusOutputJSON,
		StatusOutputCompact,
	}
}

// AllStatusOutputModeTypes is an alias for AllStatusOutputModes.
func AllStatusOutputModeTypes() []StatusOutputModeType {
	return AllStatusOutputModes()
}

// StatusTableColumnType represents column headers rendered in status tables.
type StatusTableColumnType string

const (
	StatusTableColRepo   StatusTableColumnType = "Repository"
	StatusTableColBranch StatusTableColumnType = "Branch"
	StatusTableColStatus StatusTableColumnType = "Status"
	StatusTableColSync   StatusTableColumnType = "Sync"
	StatusTableColStash  StatusTableColumnType = "Stash"
	StatusTableColFiles  StatusTableColumnType = "Files"
)

// IsValid reports whether the column header matches a recognized variant.
func (c StatusTableColumnType) IsValid() bool {
	switch c {
	case StatusTableColRepo, StatusTableColBranch, StatusTableColStatus,
		StatusTableColSync, StatusTableColStash, StatusTableColFiles:
		return true
	default:
		return false
	}
}

// String returns the header column name.
func (c StatusTableColumnType) String() string {
	return string(c)
}

// IsRepo reports whether the column is Repository.
func (c StatusTableColumnType) IsRepo() bool {
	return c == StatusTableColRepo
}

// IsBranch reports whether the column is Branch.
func (c StatusTableColumnType) IsBranch() bool {
	return c == StatusTableColBranch
}

// IsStatus reports whether the column is Status.
func (c StatusTableColumnType) IsStatus() bool {
	return c == StatusTableColStatus
}

// IsSync reports whether the column is Sync.
func (c StatusTableColumnType) IsSync() bool {
	return c == StatusTableColSync
}

// IsStash reports whether the column is Stash.
func (c StatusTableColumnType) IsStash() bool {
	return c == StatusTableColStash
}

// IsFiles reports whether the column is Files.
func (c StatusTableColumnType) IsFiles() bool {
	return c == StatusTableColFiles
}

// ParseStatusTableColumn parses a table column string (case-insensitive) into a typed enum.
func ParseStatusTableColumn(s string) (StatusTableColumnType, *apperror.AppError) {
	norm := strings.ToLower(strings.TrimSpace(s))
	for _, col := range AllStatusTableColumns() {
		if strings.ToLower(col.String()) == norm {
			return col, nil
		}
	}

	return "", apperror.NewValidation("ParseStatusTableColumn", "E1000", fmt.Sprintf("invalid status table column: %q", s)).WithContext("value", s)
}

// ParseStatusTableColumnType is an alias for ParseStatusTableColumn.
func ParseStatusTableColumnType(s string) (StatusTableColumnType, *apperror.AppError) {
	return ParseStatusTableColumn(s)
}

// AllStatusTableColumns returns all status table columns in display order.
func AllStatusTableColumns() []StatusTableColumnType {
	return []StatusTableColumnType{
		StatusTableColRepo,
		StatusTableColBranch,
		StatusTableColStatus,
		StatusTableColSync,
		StatusTableColStash,
		StatusTableColFiles,
	}
}

// AllStatusTableColumnTypes is an alias for AllStatusTableColumns.
func AllStatusTableColumnTypes() []StatusTableColumnType {
	return AllStatusTableColumns()
}

// StatusRepoStateType represents the operational health state of a repository.
type StatusRepoStateType string

const (
	StatusRepoStateClean   StatusRepoStateType = "clean"
	StatusRepoStateDirty   StatusRepoStateType = "dirty"
	StatusRepoStateAhead   StatusRepoStateType = "ahead"
	StatusRepoStateBehind  StatusRepoStateType = "behind"
	StatusRepoStateStashed StatusRepoStateType = "stashed"
	StatusRepoStateMissing StatusRepoStateType = "missing"
)

// IsValid reports whether the repository state matches a recognized variant.
func (s StatusRepoStateType) IsValid() bool {
	switch s {
	case StatusRepoStateClean, StatusRepoStateDirty, StatusRepoStateAhead,
		StatusRepoStateBehind, StatusRepoStateStashed, StatusRepoStateMissing:
		return true
	default:
		return false
	}
}

// String returns the status state identifier.
func (s StatusRepoStateType) String() string {
	return string(s)
}

// IsClean reports whether the state denotes a clean repository.
func (s StatusRepoStateType) IsClean() bool {
	return s == StatusRepoStateClean
}

// IsDirty reports whether the state denotes a dirty repository.
func (s StatusRepoStateType) IsDirty() bool {
	return s == StatusRepoStateDirty
}

// IsAhead reports whether the repository has unpushed commits.
func (s StatusRepoStateType) IsAhead() bool {
	return s == StatusRepoStateAhead
}

// IsBehind reports whether the repository has remote commits to pull.
func (s StatusRepoStateType) IsBehind() bool {
	return s == StatusRepoStateBehind
}

// IsStashed reports whether the repository contains stashed changes.
func (s StatusRepoStateType) IsStashed() bool {
	return s == StatusRepoStateStashed
}

// IsMissing reports whether the repository folder is missing from disk.
func (s StatusRepoStateType) IsMissing() bool {
	return s == StatusRepoStateMissing
}

// ParseStatusRepoState parses a repo state string into a typed enum.
func ParseStatusRepoState(s string) (StatusRepoStateType, *apperror.AppError) {
	val := StatusRepoStateType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseStatusRepoState", "E1000", fmt.Sprintf("invalid status repo state: %q", s)).WithContext("value", s)
}

// ParseStatusRepoStateType is an alias for ParseStatusRepoState.
func ParseStatusRepoStateType(s string) (StatusRepoStateType, *apperror.AppError) {
	return ParseStatusRepoState(s)
}

// AllStatusRepoStates returns all supported status repository states.
func AllStatusRepoStates() []StatusRepoStateType {
	return []StatusRepoStateType{
		StatusRepoStateClean,
		StatusRepoStateDirty,
		StatusRepoStateAhead,
		StatusRepoStateBehind,
		StatusRepoStateStashed,
		StatusRepoStateMissing,
	}
}

// AllStatusRepoStateTypes is an alias for AllStatusRepoStates.
func AllStatusRepoStateTypes() []StatusRepoStateType {
	return AllStatusRepoStates()
}
