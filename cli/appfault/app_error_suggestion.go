// Package appfault provides canonical compatibility aliases for apperror.
package appfault

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/diag"
)

// AppError is an alias for apperror.AppError to satisfy canonical coding guidelines.
type AppError = apperror.AppError

// Suggestion is an alias for diag.Suggestion.
type Suggestion = diag.Suggestion

// NewNotFoundError creates an AppError specialized for missing items or lookup misses.
func NewNotFoundError(msg string) *AppError {
	return apperror.NewNotFoundError(msg)
}

// NewValidationError creates an AppError specialized for input/CLI validation failures.
func NewValidationError(msg string) *AppError {
	return apperror.NewValidationError(msg)
}
