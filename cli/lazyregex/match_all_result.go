package lazyregex

import "github.com/alimtvnetwork/gitmap-v28/cli/apperror"

// NewMatchAllSuccess creates a MatchAllResult holding matched groups.
func NewMatchAllSuccess(pattern, comparing string, groups []*MatchGroup) *MatchAllResult {
	return &MatchAllResult{
		pattern:   pattern,
		comparing: comparing,
		groups:    groups,
		isMatched: len(groups) > 0,
	}
}

// NewMatchAllFailure creates a MatchAllResult indicating match failure.
func NewMatchAllFailure(pattern, comparing string, appErr *apperror.AppError) *MatchAllResult {
	if appErr == nil {
		appErr = FormatMatchFailureError(pattern, comparing, nil)
	}

	return &MatchAllResult{
		pattern:   pattern,
		comparing: comparing,
		isMatched: false,
		appErr:    appErr,
	}
}

// Groups returns all matched MatchGroup instances.
func (it *MatchAllResult) Groups() []*MatchGroup {
	if it == nil || len(it.groups) == 0 {
		return []*MatchGroup{}
	}

	return it.groups
}

// Items returns the primary matched string for every matched group.
func (it *MatchAllResult) Items() []string {
	if it == nil || len(it.groups) == 0 {
		return []string{}
	}

	items := make([]string, 0, len(it.groups))
	for _, g := range it.groups {
		if g != nil {
			items = append(items, g.First())
		}
	}

	return items
}

// Maps returns a GroupList of named groups across all matches.
func (it *MatchAllResult) Maps() GroupList {
	if it == nil || len(it.groups) == 0 {
		return NewGroupList()
	}

	gl := make(GroupList, 0, len(it.groups))
	for _, g := range it.groups {
		if g != nil {
			gl = append(gl, g.Map())
		}
	}

	return gl
}

// First returns the first MatchGroup or nil if empty.
func (it *MatchAllResult) First() *MatchGroup {
	if it == nil || len(it.groups) == 0 {
		return nil
	}

	return it.groups[0]
}

// Last returns the last MatchGroup or nil if empty.
func (it *MatchAllResult) Last() *MatchGroup {
	if it == nil || len(it.groups) == 0 {
		return nil
	}

	return it.groups[len(it.groups)-1]
}

// FirstOrDefault returns the first MatchGroup or defaultValue if empty.
func (it *MatchAllResult) FirstOrDefault(defaultValue *MatchGroup) *MatchGroup {
	if it == nil || len(it.groups) == 0 {
		return defaultValue
	}

	return it.groups[0]
}

// Count returns the number of matched groups.
func (it *MatchAllResult) Count() int {
	if it == nil {
		return 0
	}

	return len(it.groups)
}

// Len is an alias for Count.
func (it *MatchAllResult) Len() int {
	return it.Count()
}

// IsMatch reports whether any occurrences were matched.
func (it *MatchAllResult) IsMatch() bool {
	return it != nil && it.isMatched
}

// HasMatch is an alias for IsMatch.
func (it *MatchAllResult) HasMatch() bool {
	return it.IsMatch()
}

// IsSuccess reports whether matching succeeded.
func (it *MatchAllResult) IsSuccess() bool {
	return it.IsMatch()
}

// IsFailed reports whether matching failed or receiver is nil.
func (it *MatchAllResult) IsFailed() bool {
	return it == nil || !it.isMatched
}

// IsFailure is an alias for IsFailed.
func (it *MatchAllResult) IsFailure() bool {
	return it.IsFailed()
}

// AppError returns structured diagnostics if matching failed.
func (it *MatchAllResult) AppError() *apperror.AppError {
	if it == nil {
		return apperror.NewSimple("nil MatchAllResult", "E9000")
	}

	if !it.isMatched && it.appErr == nil {
		it.appErr = FormatMatchFailureError(it.pattern, it.comparing, nil)
	}

	return it.appErr
}

// Cause returns the error as standard error interface or nil on success.
func (it *MatchAllResult) Cause() error {
	return it.AsError()
}

// AsError returns untyped nil on success, or AppError on failure.
func (it *MatchAllResult) AsError() error {
	if it.IsSuccess() {
		return nil
	}

	return it.AppError()
}

// ErrOrNil returns untyped nil on success, or AppError on failure.
func (it *MatchAllResult) ErrOrNil() error {
	return it.AsError()
}

// ErrorString returns the failure message or empty string on success.
func (it *MatchAllResult) ErrorString() string {
	if it.IsSuccess() {
		return ""
	}

	return it.AppError().Error()
}
