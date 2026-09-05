package regexnew

import (
	"regexp"

	"coding-guidelines/common/pkg/appfault"
)

// CompileResult encapsulates the outcome of a lazy regex compilation,
// holding either the compiled regexp or structured AppError / AppBuilder diagnostics.
type CompileResult struct {
	re       *regexp.Regexp
	appError *appfault.AppError
	builder  *appfault.AppBuilder
}

// NewCompileSuccess creates a successful CompileResult wrapping the compiled regexp.
func NewCompileSuccess(re *regexp.Regexp) *CompileResult {
	return &CompileResult{
		re: re,
	}
}

// NewCompileFailure creates a failed CompileResult wrapping the diagnostic AppError and AppBuilder.
func NewCompileFailure(appErr *appfault.AppError, builder *appfault.AppBuilder) *CompileResult {
	return &CompileResult{
		appError: appErr,
		builder:  builder,
	}
}

// IsSuccess reports whether compilation succeeded without errors.
func (it *CompileResult) IsSuccess() bool {
	if it == nil {
		return false
	}

	return it.re != nil && it.appError == nil
}

// IsFailed reports whether compilation failed.
func (it *CompileResult) IsFailed() bool {
	return !it.IsSuccess()
}

// HasError reports whether a non-nil AppError is present.
func (it *CompileResult) HasError() bool {
	if it == nil {
		return true
	}

	return it.appError != nil
}

// IsValid is an alias for IsSuccess.
func (it *CompileResult) IsValid() bool {
	return it.IsSuccess()
}

// Regexp returns the underlying compiled *regexp.Regexp or nil.
func (it *CompileResult) Regexp() *regexp.Regexp {
	if it == nil {
		return nil
	}

	return it.re
}

// Value is an alias for Regexp.
func (it *CompileResult) Value() *regexp.Regexp {
	return it.Regexp()
}

// AppError returns the underlying *appfault.AppError or nil.
func (it *CompileResult) AppError() *appfault.AppError {
	if it == nil {
		return nil
	}

	return it.appError
}

// Builder returns the underlying *appfault.AppBuilder or nil.
func (it *CompileResult) Builder() *appfault.AppBuilder {
	if it == nil {
		return nil
	}

	return it.builder
}

// Error returns the underlying AppError as standard Go error interface.
func (it *CompileResult) Error() error {
	if it == nil || it.appError == nil {
		return nil
	}

	return it.appError
}

// Unwrap returns the tuple (Regexp, Error) for convenient multi-value assignment.
func (it *CompileResult) Unwrap() (*regexp.Regexp, error) {
	if it == nil {
		return nil, nil
	}

	return it.re, it.Error()
}
