package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// BranchFilterType specifies criteria for filtering branches.
type BranchFilterType string

const (
	BranchFilterAll      BranchFilterType = "all"
	BranchFilterLocal    BranchFilterType = "local"
	BranchFilterRemote   BranchFilterType = "remote"
	BranchFilterMerged   BranchFilterType = "merged"
	BranchFilterUnmerged BranchFilterType = "unmerged"
)

// IsValid reports whether the filter is a recognized enum value.
func (f BranchFilterType) IsValid() bool {
	switch f {
	case BranchFilterAll, BranchFilterLocal, BranchFilterRemote,
		BranchFilterMerged, BranchFilterUnmerged:
		return true
	default:
		return false
	}
}

// String returns the string representation of the filter.
func (f BranchFilterType) String() string {
	return string(f)
}

// ParseBranchFilterType parses a raw string into a BranchFilterType.
func ParseBranchFilterType(s string) (BranchFilterType, *apperror.AppError) {
	val := BranchFilterType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseBranchFilterType", "E1000", fmt.Sprintf("invalid branch filter: %q", s)).WithContext("value", s)
}

// AllBranchFilterTypes returns all supported BranchFilterType values.
func AllBranchFilterTypes() []BranchFilterType {
	return []BranchFilterType{
		BranchFilterAll,
		BranchFilterLocal,
		BranchFilterRemote,
		BranchFilterMerged,
		BranchFilterUnmerged,
	}
}

// IsAll reports whether the filter targets all branches.
func (f BranchFilterType) IsAll() bool {
	return f == BranchFilterAll
}

// IsLocal reports whether the filter targets local branches only.
func (f BranchFilterType) IsLocal() bool {
	return f == BranchFilterLocal
}

// IsRemote reports whether the filter targets remote branches only.
func (f BranchFilterType) IsRemote() bool {
	return f == BranchFilterRemote
}

// IsMerged reports whether the filter targets merged branches.
func (f BranchFilterType) IsMerged() bool {
	return f == BranchFilterMerged
}

// IsUnmerged reports whether the filter targets unmerged branches.
func (f BranchFilterType) IsUnmerged() bool {
	return f == BranchFilterUnmerged
}

// HasFilter reports whether the filter value is valid and non-empty.
func (f BranchFilterType) HasFilter() bool {
	return f.IsValid()
}

// BranchActionType specifies an action or subcommand within the branch domain.
type BranchActionType string

const (
	BranchActionDefault           BranchActionType = "default"
	BranchActionDefaultAlias      BranchActionType = "def"
	BranchActionList              BranchActionType = "list"
	BranchActionCurrent           BranchActionType = "current"
	BranchActionLatestBranch      BranchActionType = "latest-branch"
	BranchActionLatestBranchAlias BranchActionType = "lb"
	BranchActionHelp              BranchActionType = "help"
)

// IsValid reports whether the action is a recognized enum value.
func (a BranchActionType) IsValid() bool {
	switch a {
	case BranchActionDefault, BranchActionDefaultAlias,
		BranchActionList, BranchActionCurrent,
		BranchActionLatestBranch, BranchActionLatestBranchAlias,
		BranchActionHelp:
		return true
	default:
		return false
	}
}

// String returns the string representation of the action.
func (a BranchActionType) String() string {
	return string(a)
}

// ParseBranchActionType parses a raw string into a BranchActionType.
func ParseBranchActionType(s string) (BranchActionType, *apperror.AppError) {
	val := BranchActionType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseBranchActionType", "E1000", fmt.Sprintf("invalid branch action: %q", s)).WithContext("value", s)
}

// AllBranchActionTypes returns all supported BranchActionType values.
func AllBranchActionTypes() []BranchActionType {
	return []BranchActionType{
		BranchActionDefault,
		BranchActionDefaultAlias,
		BranchActionList,
		BranchActionCurrent,
		BranchActionLatestBranch,
		BranchActionLatestBranchAlias,
		BranchActionHelp,
	}
}

// IsDefault reports whether the action is the canonical default command.
func (a BranchActionType) IsDefault() bool {
	return a == BranchActionDefault
}

// IsDefaultAlias reports whether the action is the default command alias.
func (a BranchActionType) IsDefaultAlias() bool {
	return a == BranchActionDefaultAlias
}

// IsList reports whether the action lists branches.
func (a BranchActionType) IsList() bool {
	return a == BranchActionList
}

// IsCurrent reports whether the action displays the current branch.
func (a BranchActionType) IsCurrent() bool {
	return a == BranchActionCurrent
}

// IsLatestBranch reports whether the action inspects the latest branch.
func (a BranchActionType) IsLatestBranch() bool {
	return a == BranchActionLatestBranch
}

// IsLatestBranchAlias reports whether the action is the latest-branch alias.
func (a BranchActionType) IsLatestBranchAlias() bool {
	return a == BranchActionLatestBranchAlias
}

// IsHelp reports whether the action displays branch help.
func (a BranchActionType) IsHelp() bool {
	return a == BranchActionHelp
}

// HasDefaultAction reports whether the action resolves to the default branch command.
func (a BranchActionType) HasDefaultAction() bool {
	return a == BranchActionDefault || a == BranchActionDefaultAlias
}

// HasLatestBranchAction reports whether the action resolves to latest branch.
func (a BranchActionType) HasLatestBranchAction() bool {
	return a == BranchActionLatestBranch || a == BranchActionLatestBranchAlias
}

// HasAction reports whether the action value is valid and non-empty.
func (a BranchActionType) HasAction() bool {
	return a.IsValid()
}

// BranchSortOrderType specifies sort criteria for branch listings.
type BranchSortOrderType string

const (
	BranchSortOrderRecent        BranchSortOrderType = "recent"
	BranchSortOrderCommitterDate BranchSortOrderType = "committerdate"
	BranchSortOrderAuthorDate    BranchSortOrderType = "authordate"
	BranchSortOrderAlphabetical  BranchSortOrderType = "alphabetical"
)

// IsValid reports whether the sort order is a recognized enum value.
func (s BranchSortOrderType) IsValid() bool {
	switch s {
	case BranchSortOrderRecent, BranchSortOrderCommitterDate,
		BranchSortOrderAuthorDate, BranchSortOrderAlphabetical:
		return true
	default:
		return false
	}
}

// String returns the string representation of the sort order.
func (s BranchSortOrderType) String() string {
	return string(s)
}

// ParseBranchSortOrderType parses a raw string into a BranchSortOrderType.
func ParseBranchSortOrderType(s string) (BranchSortOrderType, *apperror.AppError) {
	val := BranchSortOrderType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseBranchSortOrderType", "E1000", fmt.Sprintf("invalid branch sort order: %q", s)).WithContext("value", s)
}

// AllBranchSortOrderTypes returns all supported BranchSortOrderType values.
func AllBranchSortOrderTypes() []BranchSortOrderType {
	return []BranchSortOrderType{
		BranchSortOrderRecent,
		BranchSortOrderCommitterDate,
		BranchSortOrderAuthorDate,
		BranchSortOrderAlphabetical,
	}
}

// IsRecent reports whether sorting by recent activity.
func (s BranchSortOrderType) IsRecent() bool {
	return s == BranchSortOrderRecent
}

// IsCommitterDate reports whether sorting by committer date.
func (s BranchSortOrderType) IsCommitterDate() bool {
	return s == BranchSortOrderCommitterDate
}

// IsAuthorDate reports whether sorting by author date.
func (s BranchSortOrderType) IsAuthorDate() bool {
	return s == BranchSortOrderAuthorDate
}

// IsAlphabetical reports whether sorting alphabetically by name.
func (s BranchSortOrderType) IsAlphabetical() bool {
	return s == BranchSortOrderAlphabetical
}

// HasSortOrder reports whether the sort order is valid and non-empty.
func (s BranchSortOrderType) HasSortOrder() bool {
	return s.IsValid()
}

// BranchOutputFormatType specifies output serialization for branch information.
type BranchOutputFormatType string

const (
	BranchOutputFormatTerminal BranchOutputFormatType = "terminal"
	BranchOutputFormatJSON     BranchOutputFormatType = "json"
	BranchOutputFormatCSV      BranchOutputFormatType = "csv"
)

// IsValid reports whether the output format is a recognized enum value.
func (o BranchOutputFormatType) IsValid() bool {
	switch o {
	case BranchOutputFormatTerminal, BranchOutputFormatJSON, BranchOutputFormatCSV:
		return true
	default:
		return false
	}
}

// String returns the string representation of the output format.
func (o BranchOutputFormatType) String() string {
	return string(o)
}

// ParseBranchOutputFormatType parses a raw string into a BranchOutputFormatType.
func ParseBranchOutputFormatType(s string) (BranchOutputFormatType, *apperror.AppError) {
	val := BranchOutputFormatType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseBranchOutputFormatType", "E1000", fmt.Sprintf("invalid branch output format: %q", s)).WithContext("value", s)
}

// AllBranchOutputFormatTypes returns all supported BranchOutputFormatType values.
func AllBranchOutputFormatTypes() []BranchOutputFormatType {
	return []BranchOutputFormatType{
		BranchOutputFormatTerminal,
		BranchOutputFormatJSON,
		BranchOutputFormatCSV,
	}
}

// IsTerminal reports whether the output format is human-readable terminal.
func (o BranchOutputFormatType) IsTerminal() bool {
	return o == BranchOutputFormatTerminal
}

// IsJSON reports whether the output format is structured JSON.
func (o BranchOutputFormatType) IsJSON() bool {
	return o == BranchOutputFormatJSON
}

// IsCSV reports whether the output format is comma-separated values.
func (o BranchOutputFormatType) IsCSV() bool {
	return o == BranchOutputFormatCSV
}

// HasOutputFormat reports whether the output format is valid and non-empty.
func (o BranchOutputFormatType) HasOutputFormat() bool {
	return o.IsValid()
}
