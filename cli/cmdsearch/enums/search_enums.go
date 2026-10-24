package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// SearchScopeType defines the domain or boundary of search execution.
type SearchScopeType string

const (
	SearchScopeContent SearchScopeType = "content"
	SearchScopeSymbol  SearchScopeType = "symbol"
	SearchScopePath    SearchScopeType = "path"
	SearchScopeCommit  SearchScopeType = "commit"
	SearchScopeRepo    SearchScopeType = "repo"
)

// IsValid reports whether the search scope matches a recognized variant.
func (s SearchScopeType) IsValid() bool {
	switch s {
	case SearchScopeContent, SearchScopeSymbol, SearchScopePath, SearchScopeCommit, SearchScopeRepo:
		return true
	default:
		return false
	}
}

// String returns the search scope identifier string.
func (s SearchScopeType) String() string {
	return string(s)
}

// IsContent reports whether the scope targets file content.
func (s SearchScopeType) IsContent() bool {
	return s == SearchScopeContent
}

// IsSymbol reports whether the scope targets language symbols.
func (s SearchScopeType) IsSymbol() bool {
	return s == SearchScopeSymbol
}

// IsPath reports whether the scope targets file paths.
func (s SearchScopeType) IsPath() bool {
	return s == SearchScopePath
}

// IsCommit reports whether the scope targets commit messages.
func (s SearchScopeType) IsCommit() bool {
	return s == SearchScopeCommit
}

// IsRepo reports whether the scope targets repositories.
func (s SearchScopeType) IsRepo() bool {
	return s == SearchScopeRepo
}

// ParseSearchScope parses a string into a SearchScopeType enum.
func ParseSearchScope(s string) (SearchScopeType, *apperror.AppError) {
	val := SearchScopeType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseSearchScope", "E1000", fmt.Sprintf("invalid search scope: %q", s)).WithContext("value", s)
}

// ParseSearchScopeType is an alias for ParseSearchScope.
func ParseSearchScopeType(s string) (SearchScopeType, *apperror.AppError) {
	return ParseSearchScope(s)
}

// AllSearchScopes returns all supported search scope variants.
func AllSearchScopes() []SearchScopeType {
	return []SearchScopeType{
		SearchScopeContent,
		SearchScopeSymbol,
		SearchScopePath,
		SearchScopeCommit,
		SearchScopeRepo,
	}
}

// AllSearchScopeTypes is an alias for AllSearchScopes.
func AllSearchScopeTypes() []SearchScopeType {
	return AllSearchScopes()
}

// SearchMatchModeType represents the matching algorithm or engine used.
type SearchMatchModeType string

const (
	SearchMatchModeKeyword SearchMatchModeType = "keyword"
	SearchMatchModeRegex   SearchMatchModeType = "regex"
	SearchMatchModeExact   SearchMatchModeType = "exact"
	SearchMatchModeFuzzy   SearchMatchModeType = "fuzzy"
)

// IsValid reports whether the match mode matches a recognized variant.
func (m SearchMatchModeType) IsValid() bool {
	switch m {
	case SearchMatchModeKeyword, SearchMatchModeRegex, SearchMatchModeExact, SearchMatchModeFuzzy:
		return true
	default:
		return false
	}
}

// String returns the match mode identifier string.
func (m SearchMatchModeType) String() string {
	return string(m)
}

// IsKeyword reports whether keyword mode is selected.
func (m SearchMatchModeType) IsKeyword() bool {
	return m == SearchMatchModeKeyword
}

// IsRegex reports whether regular expression mode is selected.
func (m SearchMatchModeType) IsRegex() bool {
	return m == SearchMatchModeRegex
}

// IsExact reports whether exact match mode is selected.
func (m SearchMatchModeType) IsExact() bool {
	return m == SearchMatchModeExact
}

// IsFuzzy reports whether fuzzy search mode is selected.
func (m SearchMatchModeType) IsFuzzy() bool {
	return m == SearchMatchModeFuzzy
}

// ParseSearchMatchMode parses a string into a SearchMatchModeType enum.
func ParseSearchMatchMode(s string) (SearchMatchModeType, *apperror.AppError) {
	val := SearchMatchModeType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseSearchMatchMode", "E1000", fmt.Sprintf("invalid search match mode: %q", s)).WithContext("value", s)
}

// ParseSearchMatchModeType is an alias for ParseSearchMatchMode.
func ParseSearchMatchModeType(s string) (SearchMatchModeType, *apperror.AppError) {
	return ParseSearchMatchMode(s)
}

// AllSearchMatchModes returns all supported search match mode variants.
func AllSearchMatchModes() []SearchMatchModeType {
	return []SearchMatchModeType{
		SearchMatchModeKeyword,
		SearchMatchModeRegex,
		SearchMatchModeExact,
		SearchMatchModeFuzzy,
	}
}

// AllSearchMatchModeTypes is an alias for AllSearchMatchModes.
func AllSearchMatchModeTypes() []SearchMatchModeType {
	return AllSearchMatchModes()
}

// SearchOutputFormatType represents output presentation formatting.
type SearchOutputFormatType string

const (
	SearchOutputFormatText    SearchOutputFormatType = "text"
	SearchOutputFormatJSON    SearchOutputFormatType = "json"
	SearchOutputFormatTable   SearchOutputFormatType = "table"
	SearchOutputFormatCompact SearchOutputFormatType = "compact"
)

// IsValid reports whether the output format matches a recognized variant.
func (o SearchOutputFormatType) IsValid() bool {
	switch o {
	case SearchOutputFormatText, SearchOutputFormatJSON, SearchOutputFormatTable, SearchOutputFormatCompact:
		return true
	default:
		return false
	}
}

// String returns the output format identifier string.
func (o SearchOutputFormatType) String() string {
	return string(o)
}

// IsText reports whether plain text output format is selected.
func (o SearchOutputFormatType) IsText() bool {
	return o == SearchOutputFormatText
}

// IsJSON reports whether JSON output format is selected.
func (o SearchOutputFormatType) IsJSON() bool {
	return o == SearchOutputFormatJSON
}

// IsTable reports whether rich ANSI table output is selected.
func (o SearchOutputFormatType) IsTable() bool {
	return o == SearchOutputFormatTable
}

// IsCompact reports whether single-line compact output is selected.
func (o SearchOutputFormatType) IsCompact() bool {
	return o == SearchOutputFormatCompact
}

// ParseSearchOutputFormat parses a string into a SearchOutputFormatType enum.
func ParseSearchOutputFormat(s string) (SearchOutputFormatType, *apperror.AppError) {
	val := SearchOutputFormatType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseSearchOutputFormat", "E1000", fmt.Sprintf("invalid search output format: %q", s)).WithContext("value", s)
}

// ParseSearchOutputFormatType is an alias for ParseSearchOutputFormat.
func ParseSearchOutputFormatType(s string) (SearchOutputFormatType, *apperror.AppError) {
	return ParseSearchOutputFormat(s)
}

// AllSearchOutputFormats returns all supported output format variants.
func AllSearchOutputFormats() []SearchOutputFormatType {
	return []SearchOutputFormatType{
		SearchOutputFormatText,
		SearchOutputFormatJSON,
		SearchOutputFormatTable,
		SearchOutputFormatCompact,
	}
}

// AllSearchOutputFormatTypes is an alias for AllSearchOutputFormats.
func AllSearchOutputFormatTypes() []SearchOutputFormatType {
	return AllSearchOutputFormats()
}

// SearchCategoryType represents the origin or category of the search query.
type SearchCategoryType string

const (
	SearchCategoryUser SearchCategoryType = "user"
	SearchCategoryAI   SearchCategoryType = "ai"
)

// IsValid reports whether the search category matches a recognized variant.
func (c SearchCategoryType) IsValid() bool {
	switch c {
	case SearchCategoryUser, SearchCategoryAI:
		return true
	default:
		return false
	}
}

// String returns the search category identifier string.
func (c SearchCategoryType) String() string {
	return string(c)
}

// IsUser reports whether the category represents a human user query.
func (c SearchCategoryType) IsUser() bool {
	return c == SearchCategoryUser
}

// IsAI reports whether the category represents an automated AI query.
func (c SearchCategoryType) IsAI() bool {
	return c == SearchCategoryAI
}

// ParseSearchCategory parses a string into a SearchCategoryType enum.
func ParseSearchCategory(s string) (SearchCategoryType, *apperror.AppError) {
	val := SearchCategoryType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseSearchCategory", "E1000", fmt.Sprintf("invalid search category: %q", s)).WithContext("value", s)
}

// ParseSearchCategoryType is an alias for ParseSearchCategory.
func ParseSearchCategoryType(s string) (SearchCategoryType, *apperror.AppError) {
	return ParseSearchCategory(s)
}

// AllSearchCategories returns all supported search categories.
func AllSearchCategories() []SearchCategoryType {
	return []SearchCategoryType{
		SearchCategoryUser,
		SearchCategoryAI,
	}
}

// AllSearchCategoryTypes is an alias for AllSearchCategories.
func AllSearchCategoryTypes() []SearchCategoryType {
	return AllSearchCategories()
}

// SearchActionType represents operational actions dispatched during search.
type SearchActionType string

const (
	SearchActionQuery   SearchActionType = "query"
	SearchActionHistory SearchActionType = "history"
	SearchActionTop     SearchActionType = "top"
	SearchActionStats   SearchActionType = "stats"
	SearchActionDH2D    SearchActionType = "dh2d"
	SearchActionClean   SearchActionType = "clean"
	SearchActionClear   SearchActionType = "clear"
	SearchActionPurge   SearchActionType = "purge"
	SearchActionReset   SearchActionType = "reset"
)

// IsValid reports whether the action matches a recognized variant.
func (a SearchActionType) IsValid() bool {
	switch a {
	case SearchActionQuery, SearchActionHistory, SearchActionTop, SearchActionStats,
		SearchActionDH2D, SearchActionClean, SearchActionClear, SearchActionPurge, SearchActionReset:
		return true
	default:
		return false
	}
}

// String returns the action identifier string.
func (a SearchActionType) String() string {
	return string(a)
}

// IsQuery reports whether the action is a search query.
func (a SearchActionType) IsQuery() bool {
	return a == SearchActionQuery
}

// IsHistory reports whether the action requests history logs.
func (a SearchActionType) IsHistory() bool {
	return a == SearchActionHistory
}

// IsTop reports whether the action requests top ranking queries.
func (a SearchActionType) IsTop() bool {
	return a == SearchActionTop
}

// IsStats reports whether the action requests search metrics.
func (a SearchActionType) IsStats() bool {
	return a == SearchActionStats
}

// IsDH2D reports whether the action requests DH2D cache inspection.
func (a SearchActionType) IsDH2D() bool {
	return a == SearchActionDH2D
}

// IsClean reports whether the action purges search cache.
func (a SearchActionType) IsClean() bool {
	return a == SearchActionClean
}

// IsClear reports whether the action clears cache.
func (a SearchActionType) IsClear() bool {
	return a == SearchActionClear
}

// IsPurge reports whether the action purges cache.
func (a SearchActionType) IsPurge() bool {
	return a == SearchActionPurge
}

// IsReset reports whether the action resets search output.
func (a SearchActionType) IsReset() bool {
	return a == SearchActionReset
}

// IsHistoryAction reports whether the action queries history, top, stats, or DH2D.
func (a SearchActionType) IsHistoryAction() bool {
	switch a {
	case SearchActionHistory, SearchActionTop, SearchActionStats, SearchActionDH2D:
		return true
	default:
		return false
	}
}

// IsCleanAction reports whether the action is a cache cleanup operation (clean, clear, purge).
func (a SearchActionType) IsCleanAction() bool {
	switch a {
	case SearchActionClean, SearchActionClear, SearchActionPurge:
		return true
	default:
		return false
	}
}

// IsResetAction reports whether the action resets or clears output caches.
func (a SearchActionType) IsResetAction() bool {
	switch a {
	case SearchActionReset, SearchActionClear:
		return true
	default:
		return false
	}
}

// ParseSearchAction parses a string into a SearchActionType enum.
func ParseSearchAction(s string) (SearchActionType, *apperror.AppError) {
	val := SearchActionType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseSearchAction", "E1000", fmt.Sprintf("invalid search action: %q", s)).WithContext("value", s)
}

// ParseSearchActionType is an alias for ParseSearchAction.
func ParseSearchActionType(s string) (SearchActionType, *apperror.AppError) {
	return ParseSearchAction(s)
}

// AllSearchActions returns all supported search action variants.
func AllSearchActions() []SearchActionType {
	return []SearchActionType{
		SearchActionQuery,
		SearchActionHistory,
		SearchActionTop,
		SearchActionStats,
		SearchActionDH2D,
		SearchActionClean,
		SearchActionClear,
		SearchActionPurge,
		SearchActionReset,
	}
}

// AllSearchActionTypes is an alias for AllSearchActions.
func AllSearchActionTypes() []SearchActionType {
	return AllSearchActions()
}
