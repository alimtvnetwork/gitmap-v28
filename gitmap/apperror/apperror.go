package apperror

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("not found")

// AppError is a typed error that encapsulates an operation label,
// key input context, and the underlying cause.
type AppError struct {
	Op    string
	Ctx   map[string]any
	Cause error
	Code  string
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s (ctx: %v): %v", e.Code, e.Op, e.Ctx, e.Cause)
	}
	return fmt.Sprintf("[%s] %s (ctx: %v)", e.Code, e.Op, e.Ctx)
}

// Unwrap allows standard library errors.Is and errors.As to work.
func (e *AppError) Unwrap() error {
	return e.Cause
}

// Wrap wraps an existing error with an operation label and context.
// WrapSimple creates a new AppError with an underlying cause but no context map.
func WrapSimple(err error, op string) *AppError {
	return &AppError{
		Op:    op,
		Cause: err,
	}
}

// Wrap creates a new AppError with an underlying cause and context.
func Wrap(err error, op string, ctx map[string]any) *AppError {
	return &AppError{
		Op:    op,
		Ctx:   ctx,
		Cause: err,
	}
}

// New creates a new AppError without an underlying cause.
// NewSimple creates a new AppError without an underlying cause and no context map.
func NewSimple(op string, code string) *AppError {
	return &AppError{
		Op:   op,
		Code: code,
	}
}

// New creates a new AppError without an underlying cause.
func New(op string, code string, ctx map[string]any) *AppError {
	return &AppError{
		Op:   op,
		Code: code,
		Ctx:  ctx,
	}
}
