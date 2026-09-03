package apperror

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

var ErrNotFound = errors.New("not found")

// AppError is a typed, domain-rich error that captures operation labels,
// creator attribution, contextual metadata, severity, caller site, and root cause.
type AppError struct {
	Op       string
	Code     string
	Type     ErrorType
	Severity SeverityType
	Creator  string
	Message  string
	Caller   string
	Stack    string
	Ctx      map[string]any
	Cause    error
}

// Error formats the full diagnostic description of the AppError.
func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	parts := make([]string, 0, 5)

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

	if e.Caller != "" {
		parts = append(parts, fmt.Sprintf("(at=%s)", e.Caller))
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
	if e == nil {
		return nil
	}
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

func captureCaller(skip int) string {
	_, file, line, ok := runtime.Caller(skip)

	if !ok {
		return ""
	}

	shortFile := filepath.Base(file)
	parentDir := filepath.Base(filepath.Dir(file))

	if parentDir != "." && parentDir != "/" && parentDir != "\\" && parentDir != "" {
		return fmt.Sprintf("%s/%s:%d", parentDir, shortFile, line)
	}

	return fmt.Sprintf("%s:%d", shortFile, line)
}

func captureStackTrace(skip int) string {
	var pcs [32]uintptr
	n := runtime.Callers(skip, pcs[:])
	if n == 0 {
		return ""
	}
	frames := runtime.CallersFrames(pcs[:n])
	var sb strings.Builder
	for {
		frame, more := frames.Next()
		appendStackFrame(&sb, frame)
		if !more {
			break
		}
	}
	return sb.String()
}

func appendStackFrame(sb *strings.Builder, frame runtime.Frame) {
	isRuntimeInternal := strings.Contains(frame.Function, "runtime.") && !strings.Contains(frame.Function, "gitmap")
	if isRuntimeInternal {
		return
	}
	shortFile := filepath.Base(frame.File)
	parentDir := filepath.Base(filepath.Dir(frame.File))
	fileLoc := fmt.Sprintf("%s/%s:%d", parentDir, shortFile, frame.Line)
	sb.WriteString(fmt.Sprintf("\n    at %s (%s)", frame.Function, fileLoc))
}

// New creates a new AppError without an underlying cause.
func New(op string, code string, ctx map[string]any) *AppError {
	return &AppError{
		Op:       op,
		Code:     code,
		Type:     ErrorTypeExecution,
		Severity: SeverityError,
		Caller:   captureCaller(2),
		Stack:    captureStackTrace(2),
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
		Caller:   captureCaller(2),
		Stack:    captureStackTrace(2),
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
		Caller:   captureCaller(2),
		Stack:    captureStackTrace(2),
		Ctx:      ctx,
	}
}

// NewValidationError creates an AppError specialized for input/CLI validation failures.
func NewValidationError(msg string) *AppError {
	return &AppError{
		Op:       "validation",
		Code:     "E1000",
		Type:     ErrorTypeValidation,
		Severity: SeverityError,
		Message:  msg,
		Caller:   captureCaller(2),
		Stack:    captureStackTrace(2),
	}
}

// Wrap wraps an existing error with an operation label and context.
func Wrap(err error, op string, ctx map[string]any) *AppError {
	return &AppError{
		Op:       op,
		Code:     "E9000",
		Type:     ErrorTypeExecution,
		Severity: SeverityError,
		Caller:   captureCaller(2),
		Stack:    captureStackTrace(2),
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
		Caller:   captureCaller(2),
		Stack:    captureStackTrace(2),
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
		Caller:   captureCaller(2),
		Stack:    captureStackTrace(2),
		Ctx:      ctx,
		Cause:    err,
	}
}

// HasError reports whether an error exists.
func (e *AppError) HasError() bool {
	return e != nil
}

// HasNoError reports whether no error exists.
func (e *AppError) HasNoError() bool {
	return e == nil
}

// HasValidError reports whether the AppError is non-nil and has a valid code.
func (e *AppError) HasValidError() bool {
	return e != nil && e.Code != ""
}

// IsValid reports whether the AppError is non-nil and has a valid code.
func (e *AppError) IsValid() bool {
	return e != nil && e.Code != ""
}

// IsErrorCode reports whether the AppError matches the specified error code.
func (e *AppError) IsErrorCode(code string) bool {
	return e != nil && e.Code == code
}

// IsCode alias for IsErrorCode.
func (e *AppError) IsCode(code string) bool {
	return e.IsErrorCode(code)
}
