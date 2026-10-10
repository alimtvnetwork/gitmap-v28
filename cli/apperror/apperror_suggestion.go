package apperror

import "github.com/alimtvnetwork/gitmap-v28/cli/diag"

// WithSuggestions attaches actionable suggestion items to the AppError.
func (e *AppError) WithSuggestions(items ...diag.Suggestion) *AppError {
	if e == nil {
		return nil
	}
	e.Suggestions = append(e.Suggestions, items...)
	return e
}

// HasSuggestions reports whether the error contains any actionable suggestions.
func (e *AppError) HasSuggestions() bool {
	if e == nil {
		return false
	}
	return len(e.Suggestions) > 0
}

// GetSuggestions retrieves all attached suggestions safely.
func (e *AppError) GetSuggestions() []diag.Suggestion {
	if e == nil {
		return nil
	}
	return e.Suggestions
}

// WithSuggestionGroup attaches all suggestions from a SuggestionGroup.
func (e *AppError) WithSuggestionGroup(group diag.SuggestionGroup) *AppError {
	if e == nil {
		return nil
	}
	return e.WithSuggestions(group.Suggestions...)
}
