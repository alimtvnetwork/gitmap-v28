package lazyregex

import "github.com/alimtvnetwork/gitmap-v28/cli/apperror"

// NewMatchSuccess creates a MatchResult indicating successful pattern matching.
func NewMatchSuccess(pattern, comparing string, submatches []string, named GroupMap) *MatchResult {
	group := &MatchGroup{
		Pattern:     pattern,
		Content:     comparing,
		Submatches:  submatches,
		NamedGroups: named,
	}

	return &MatchResult{
		pattern:    pattern,
		comparing:  comparing,
		matchGroup: group,
		isMatched:  true,
	}
}

// NewMatchFailure creates a MatchResult indicating matching failure with diagnostics.
func NewMatchFailure(pattern, comparing string, appErr *apperror.AppError) *MatchResult {
	if appErr == nil {
		appErr = FormatMatchFailureError(pattern, comparing, nil)
	}

	return &MatchResult{
		pattern:   pattern,
		comparing: comparing,
		isMatched: false,
		appErr:    appErr,
	}
}

// IsMatch reports whether the pattern matched the compared content.
func (it *MatchResult) IsMatch() bool {
	return it != nil && it.isMatched
}

// HasMatch is an alias for IsMatch.
func (it *MatchResult) HasMatch() bool {
	return it.IsMatch()
}

// IsSuccess reports whether matching was successful (affirmative alias for IsMatch).
func (it *MatchResult) IsSuccess() bool {
	return it.IsMatch()
}

// IsFailed reports whether matching failed or receiver is nil.
func (it *MatchResult) IsFailed() bool {
	return it == nil || !it.isMatched
}

// IsFailure is an alias for IsFailed.
func (it *MatchResult) IsFailure() bool {
	return it.IsFailed()
}

// AppError returns structured error diagnostics if matching failed.
func (it *MatchResult) AppError() *apperror.AppError {
	if it == nil {
		return apperror.NewSimple("nil MatchResult", "E9000")
	}

	if !it.isMatched && it.appErr == nil {
		it.appErr = FormatMatchFailureError(it.pattern, it.comparing, nil)
	}

	return it.appErr
}

// Cause returns the failure error as standard error interface or nil on success.
func (it *MatchResult) Cause() error {
	return it.AsError()
}

// AsError returns untyped nil on success, or AppError on failure.
func (it *MatchResult) AsError() error {
	if it.IsSuccess() {
		return nil
	}

	return it.AppError()
}

// ErrOrNil returns untyped nil on success, or AppError on failure.
func (it *MatchResult) ErrOrNil() error {
	return it.AsError()
}

// ErrorString returns the failure message or empty string on success.
func (it *MatchResult) ErrorString() string {
	if it.IsSuccess() {
		return ""
	}

	return it.AppError().Error()
}

// Group returns the matched MatchGroup, or nil if unmatched.
func (it *MatchResult) Group() *MatchGroup {
	if it == nil {
		return nil
	}

	return it.matchGroup
}

// Data is an alias for Group.
func (it *MatchResult) Data() *MatchGroup {
	return it.Group()
}

// Value is an alias for Group.
func (it *MatchResult) Value() *MatchGroup {
	return it.Group()
}

// Items returns captured submatches from the matched group or empty slice.
func (it *MatchResult) Items() []string {
	if it == nil || it.matchGroup == nil {
		return []string{}
	}

	return it.matchGroup.Items()
}

// Map returns named capture groups from the matched group or empty map.
func (it *MatchResult) Map() GroupMap {
	if it == nil || it.matchGroup == nil {
		return NewGroupMap()
	}

	return it.matchGroup.Map()
}

// First returns the full match (submatch 0) or empty string.
func (it *MatchResult) First() string {
	if it == nil || it.matchGroup == nil {
		return ""
	}

	return it.matchGroup.First()
}

// Last returns the last submatch or empty string.
func (it *MatchResult) Last() string {
	if it == nil || it.matchGroup == nil {
		return ""
	}

	return it.matchGroup.Last()
}

// FirstOrDefault returns the full match or defaultValue if absent.
func (it *MatchResult) FirstOrDefault(defaultValue string) string {
	if it == nil || it.matchGroup == nil {
		return defaultValue
	}

	return it.matchGroup.FirstOrDefault(defaultValue)
}

// At returns the submatch at the specified index or empty string.
func (it *MatchResult) At(index int) string {
	if it == nil || it.matchGroup == nil {
		return ""
	}

	return it.matchGroup.At(index)
}

// Count returns the number of captured submatches or 0.
func (it *MatchResult) Count() int {
	if it == nil || it.matchGroup == nil {
		return 0
	}

	return it.matchGroup.Count()
}

// Len is an alias for Count.
func (it *MatchResult) Len() int {
	return it.Count()
}
