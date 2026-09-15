package lazyregex

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func truncatePreview(content string, maxLen int) string {
	if len(content) <= maxLen {
		return content
	}

	return content[:maxLen] + "... [truncated]"
}

// FormatMatchFailureError constructs a diagnostic AppError for regex matching failures.
func FormatMatchFailureError(pattern, comparing string, compileErr error) *apperror.AppError {
	preview := truncatePreview(strings.TrimSpace(comparing), 300)
	if compileErr != nil {
		msg := fmt.Sprintf("regex pattern compile failed: %v (pattern: %q)", compileErr, pattern)

		return apperror.NewValidationError(msg).
			WithContext("pattern", pattern).
			WithContext("error", compileErr.Error())
	}

	msg := fmt.Sprintf("regex pattern does not match content\n  Pattern:   %q\n  Comparing: %q\n  Length:    %d bytes",
		pattern, preview, len(comparing))

	return apperror.NewValidationError(msg).
		WithContext("pattern", pattern).
		WithContext("comparing_preview", preview).
		WithContext("comparing_len", fmt.Sprintf("%d", len(comparing)))
}
