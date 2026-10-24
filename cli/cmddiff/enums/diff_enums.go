package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// DiffViewModeType represents visual presentation modes for diff rendering.
type DiffViewModeType string

const (
	DiffViewModeUnified    DiffViewModeType = "unified"
	DiffViewModeSideBySide DiffViewModeType = "side-by-side"
	DiffViewModeStat       DiffViewModeType = "stat"
	DiffViewModeSummary    DiffViewModeType = "summary"
	DiffViewModeCompact    DiffViewModeType = "compact"
)

// IsValid reports whether the view mode is a recognized enum value.
func (d DiffViewModeType) IsValid() bool {
	switch d {
	case DiffViewModeUnified, DiffViewModeSideBySide, DiffViewModeStat,
		DiffViewModeSummary, DiffViewModeCompact:
		return true
	default:
		return false
	}
}

// String returns the string representation of the view mode.
func (d DiffViewModeType) String() string {
	return string(d)
}

// ParseDiffViewModeType parses a raw string into a DiffViewModeType.
func ParseDiffViewModeType(s string) (DiffViewModeType, *apperror.AppError) {
	val := DiffViewModeType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseDiffViewModeType", "E1000", fmt.Sprintf("invalid diff view mode: %q", s)).WithContext("value", s)
}

// AllDiffViewModeTypes returns all supported DiffViewModeType values.
func AllDiffViewModeTypes() []DiffViewModeType {
	return []DiffViewModeType{
		DiffViewModeUnified,
		DiffViewModeSideBySide,
		DiffViewModeStat,
		DiffViewModeSummary,
		DiffViewModeCompact,
	}
}

// IsUnified reports whether the view mode is unified.
func (d DiffViewModeType) IsUnified() bool {
	return d == DiffViewModeUnified
}

// IsSideBySide reports whether the view mode is side-by-side.
func (d DiffViewModeType) IsSideBySide() bool {
	return d == DiffViewModeSideBySide
}

// IsStat reports whether the view mode is stat.
func (d DiffViewModeType) IsStat() bool {
	return d == DiffViewModeStat
}

// IsSummary reports whether the view mode is summary.
func (d DiffViewModeType) IsSummary() bool {
	return d == DiffViewModeSummary
}

// IsCompact reports whether the view mode is compact.
func (d DiffViewModeType) IsCompact() bool {
	return d == DiffViewModeCompact
}

// HasViewMode reports whether the view mode is non-empty and valid.
func (d DiffViewModeType) HasViewMode() bool {
	return d.IsValid()
}

// DiffOutputFormatType represents serialization and display formats for diff output.
type DiffOutputFormatType string

const (
	DiffOutputFormatText DiffOutputFormatType = "text"
	DiffOutputFormatJSON DiffOutputFormatType = "json"
	DiffOutputFormatRaw  DiffOutputFormatType = "raw"
)

// IsValid reports whether the output format is a recognized enum value.
func (d DiffOutputFormatType) IsValid() bool {
	switch d {
	case DiffOutputFormatText, DiffOutputFormatJSON, DiffOutputFormatRaw:
		return true
	default:
		return false
	}
}

// String returns the string representation of the output format.
func (d DiffOutputFormatType) String() string {
	return string(d)
}

// ParseDiffOutputFormatType parses a raw string into a DiffOutputFormatType.
func ParseDiffOutputFormatType(s string) (DiffOutputFormatType, *apperror.AppError) {
	val := DiffOutputFormatType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseDiffOutputFormatType", "E1000", fmt.Sprintf("invalid diff output format: %q", s)).WithContext("value", s)
}

// AllDiffOutputFormatTypes returns all supported DiffOutputFormatType values.
func AllDiffOutputFormatTypes() []DiffOutputFormatType {
	return []DiffOutputFormatType{
		DiffOutputFormatText,
		DiffOutputFormatJSON,
		DiffOutputFormatRaw,
	}
}

// IsText reports whether the format is text.
func (d DiffOutputFormatType) IsText() bool {
	return d == DiffOutputFormatText
}

// IsJSON reports whether the format is json.
func (d DiffOutputFormatType) IsJSON() bool {
	return d == DiffOutputFormatJSON
}

// IsRaw reports whether the format is raw.
func (d DiffOutputFormatType) IsRaw() bool {
	return d == DiffOutputFormatRaw
}

// HasOutputFormat reports whether the output format is non-empty and valid.
func (d DiffOutputFormatType) HasOutputFormat() bool {
	return d.IsValid()
}

// DiffScopeType represents the operational boundary of a diff evaluation.
type DiffScopeType string

const (
	DiffScopeWorkingTree DiffScopeType = "working-tree"
	DiffScopeStaged      DiffScopeType = "staged"
	DiffScopeBranch      DiffScopeType = "branch"
	DiffScopeCommit      DiffScopeType = "commit"
)

// IsValid reports whether the scope is a recognized enum value.
func (d DiffScopeType) IsValid() bool {
	switch d {
	case DiffScopeWorkingTree, DiffScopeStaged, DiffScopeBranch, DiffScopeCommit:
		return true
	default:
		return false
	}
}

// String returns the string representation of the scope.
func (d DiffScopeType) String() string {
	return string(d)
}

// ParseDiffScopeType parses a raw string into a DiffScopeType.
func ParseDiffScopeType(s string) (DiffScopeType, *apperror.AppError) {
	val := DiffScopeType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseDiffScopeType", "E1000", fmt.Sprintf("invalid diff scope: %q", s)).WithContext("value", s)
}

// AllDiffScopeTypes returns all supported DiffScopeType values.
func AllDiffScopeTypes() []DiffScopeType {
	return []DiffScopeType{
		DiffScopeWorkingTree,
		DiffScopeStaged,
		DiffScopeBranch,
		DiffScopeCommit,
	}
}

// IsWorkingTree reports whether the scope is working-tree.
func (d DiffScopeType) IsWorkingTree() bool {
	return d == DiffScopeWorkingTree
}

// IsStaged reports whether the scope is staged.
func (d DiffScopeType) IsStaged() bool {
	return d == DiffScopeStaged
}

// IsBranch reports whether the scope is branch.
func (d DiffScopeType) IsBranch() bool {
	return d == DiffScopeBranch
}

// IsCommit reports whether the scope is commit.
func (d DiffScopeType) IsCommit() bool {
	return d == DiffScopeCommit
}

// HasScope reports whether the scope is non-empty and valid.
func (d DiffScopeType) HasScope() bool {
	return d.IsValid()
}

// DiffFilterType represents criteria for filtering compared file items.
type DiffFilterType string

const (
	DiffFilterAll              DiffFilterType = "all"
	DiffFilterOnlyConflicts    DiffFilterType = "only-conflicts"
	DiffFilterOnlyMissing      DiffFilterType = "only-missing"
	DiffFilterIncludeIdentical DiffFilterType = "include-identical"
)

// IsValid reports whether the filter is a recognized enum value.
func (d DiffFilterType) IsValid() bool {
	switch d {
	case DiffFilterAll, DiffFilterOnlyConflicts, DiffFilterOnlyMissing, DiffFilterIncludeIdentical:
		return true
	default:
		return false
	}
}

// String returns the string representation of the filter.
func (d DiffFilterType) String() string {
	return string(d)
}

// ParseDiffFilterType parses a raw string into a DiffFilterType.
func ParseDiffFilterType(s string) (DiffFilterType, *apperror.AppError) {
	val := DiffFilterType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseDiffFilterType", "E1000", fmt.Sprintf("invalid diff filter: %q", s)).WithContext("value", s)
}

// AllDiffFilterTypes returns all supported DiffFilterType values.
func AllDiffFilterTypes() []DiffFilterType {
	return []DiffFilterType{
		DiffFilterAll,
		DiffFilterOnlyConflicts,
		DiffFilterOnlyMissing,
		DiffFilterIncludeIdentical,
	}
}

// IsAll reports whether the filter is all.
func (d DiffFilterType) IsAll() bool {
	return d == DiffFilterAll
}

// IsOnlyConflicts reports whether the filter is only-conflicts.
func (d DiffFilterType) IsOnlyConflicts() bool {
	return d == DiffFilterOnlyConflicts
}

// IsOnlyMissing reports whether the filter is only-missing.
func (d DiffFilterType) IsOnlyMissing() bool {
	return d == DiffFilterOnlyMissing
}

// IsIncludeIdentical reports whether the filter is include-identical.
func (d DiffFilterType) IsIncludeIdentical() bool {
	return d == DiffFilterIncludeIdentical
}

// HasFilter reports whether the filter is non-empty and valid.
func (d DiffFilterType) HasFilter() bool {
	return d.IsValid()
}
