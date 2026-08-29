package apperror

import (
	"errors"
	"fmt"
	"strings"
)

var ErrNotFound = errors.New("not found")

// AppError is a typed, domain-rich error that captures operation labels,
// creator attribution, contextual metadata, severity, and root cause.
type AppError struct {
	Op       string
	Code     string
	Type     ErrorType
	Severity SeverityType
	Creator  string
	Message  string
	Ctx      map[string]any
	Cause    error
}

// Error formats the full diagnostic description of the AppError.
func (e *AppError) Error() string {
	parts := make([]string, 0, 4)
	if e.Code != "" || e.Type != "" {
		parts = append(parts, fmt.Sprintf("[%s:%s]", e.Code, e.Type))
	}
	if e.Op != "" {
		parts = append(parts, e.Op+":")
	}
	if e.Message != "" {
		parts = append(parts, e.Message)
	} else if e.Cause != nil {
		parts = append(parts, e.Cause.Error())
	}
	if e.Creator != "" {
		parts = append(parts, fmt.Sprintf("(creator=%s)", e.Creator))
	}
	if len(e.Ctx) > 0 {
		parts = append(parts, fmt.Sprintf("(ctx=%v)", e.Ctx))
	}
	if e.Cause != nil && e.Message != "" {
		parts = append(parts, fmt.Sprintf("(cause=%v)", e.Cause))
	}

	return strings.Join(parts, " ")
}

// Unwrap allows standard library errors.Is and errors.As to work.
func (e *AppError) Unwrap() error {
	return e.Cause
}

// WithContext appends a key-value pair to the error's context map.
func (e *AppError) WithContext(key string, val any) *AppError {
	if e.Ctx == nil {
		e.Ctx = make(map[string]any)
	}
	e.Ctx[key] = val

	return e
}

// New creates a new AppError without an underlying cause.
func New(op string, code string, ctx map[string]any) *AppError {
	return &AppError{
		Op:       op,
		Code:     code,
		Type:     ErrorTypeExecution,
		Severity: SeverityError,
		Ctx:      ctx,
	}
}

// NewSimple creates a new AppError without an underlying cause and no context map.
func NewSimple(op string, code string) *AppError {
	return &AppError{
		Op:       op,
		Code:     code,
		Type:     ErrorTypeExecution,
		Severity: SeverityError,
	}
}

// NewWithDetails creates a fully specified AppError without cause.
func NewWithDetails(
	op, code, msg, creator string,
	errType ErrorType,
	sev SeverityType,
	ctx map[string]any,
) *AppError {
	return &AppError{
		Op:       op,
		Code:     code,
		Type:     errType,
		Severity: sev,
		Creator:  creator,
		Message:  msg,
		Ctx:      ctx,
	}
}

// Wrap wraps an existing error with an operation label and context.
func Wrap(err error, op string, ctx map[string]any) *AppError {
	return &AppError{
		Op:       op,
		Code:     "E9000",
		Type:     ErrorTypeExecution,
		Severity: SeverityError,
		Ctx:      ctx,
		Cause:    err,
	}
}

// WrapSimple creates a new AppError with an underlying cause but no context map.
func WrapSimple(err error, op string) *AppError {
	return &AppError{
		Op:       op,
		Code:     "E9000",
		Type:     ErrorTypeExecution,
		Severity: SeverityError,
		Cause:    err,
	}
}

// WrapWithDetails wraps an existing error with full metadata.
func WrapWithDetails(
	err error,
	op, code, msg, creator string,
	errType ErrorType,
	sev SeverityType,
	ctx map[string]any,
) *AppError {
	return &AppError{
		Op:       op,
		Code:     code,
		Type:     errType,
		Severity: sev,
		Creator:  creator,
		Message:  msg,
		Ctx:      ctx,
		Cause:    err,
	}
}
