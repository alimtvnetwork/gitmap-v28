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

// WithSkip recalculates the Caller and Stack fields using an increased skip offset.
// This allows intermediate helper functions or custom wrappers to pierce abstraction layers.
func (e *AppError) WithSkip(additional int) *AppError {
	if e == nil || additional <= 0 {
		return e
	}

	e.Caller = captureCaller(DefaultCallerSkip + additional)
	e.Stack = captureStackTrace(DefaultStackTraceSkip + additional)

	return e
}

// WithAdditionalSkip is an alias for WithSkip.
func (e *AppError) WithAdditionalSkip(additional int) *AppError {
	return e.WithSkip(additional)
}

var (
	// DefaultStackTraceSkip defines the default number of frames to skip when recording stack traces.
	DefaultStackTraceSkip = 3

	// DefaultCallerSkip defines the default number of frames to skip when recording callers.
	DefaultCallerSkip = 2
)

// SetDefaultStackTraceSkip updates the default stack trace skip count.
func SetDefaultStackTraceSkip(skip int) {
	if skip >= 0 {
		DefaultStackTraceSkip = skip
	}
}

// SetDefaultCallerSkip updates the default caller skip count.
func SetDefaultCallerSkip(skip int) {
	if skip >= 0 {
		DefaultCallerSkip = skip
	}
}

// CaptureCaller records the caller location formatted as dir/file:line.
func CaptureCaller(skip int) string {
	return captureCaller(skip)
}

// CaptureStackTrace records the full formatted stack trace starting at skip frames.
func CaptureStackTrace(skip int) string {
	return captureStackTrace(skip)
}

func captureCaller(skip int) string {
	for i := 0; i < 5; i++ {
		_, file, line, isCallerAvailable := runtime.Caller(skip + i)
		if !isCallerAvailable {
			return ""
		}
		if isAppErrorFile(file) {
			continue
		}
		return formatCallerLocation(file, line)
	}
	return ""
}

func isAppErrorFile(file string) bool {
	shortFile := filepath.Base(file)
	return shortFile == "apperror.go" || shortFile == "apperror_types.go"
}

func formatCallerLocation(file string, line int) string {
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
	return buildFramesString(pcs[:n])
}

func buildFramesString(pcs []uintptr) string {
	frames := runtime.CallersFrames(pcs)
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
	if isRuntimeInternal(frame) || isAppErrorInternal(frame) {
		return
	}

	shortFile := filepath.Base(frame.File)
	parentDir := filepath.Base(filepath.Dir(frame.File))
	fileLoc := fmt.Sprintf("%s/%s:%d", parentDir, shortFile, frame.Line)
	sb.WriteString(fmt.Sprintf("\n    at %s (%s)", frame.Function, fileLoc))
}

func isRuntimeInternal(frame runtime.Frame) bool {
	return strings.Contains(frame.Function, "runtime.") && !strings.Contains(frame.Function, "gitmap")
}

func isAppErrorInternal(frame runtime.Frame) bool {
	if isAppErrorFile(frame.File) {
		return true
	}
	isAppErrorPkg := strings.Contains(frame.Function, "/cli/apperror.") || strings.HasPrefix(frame.Function, "apperror.")
	isTestFrame := strings.HasSuffix(frame.File, "_test.go") || strings.Contains(frame.Function, "Test")
	return isAppErrorPkg && !isTestFrame
}

// NewWithSkip creates a new AppError with an explicit frame skip increase.
func NewWithSkip(skip int, op string, code string) *AppError {
	callerSkip := DefaultCallerSkip + skip
	stackSkip := DefaultStackTraceSkip + skip
	return &AppError{
		Op:       op,
		Code:     code,
		Type:     ErrorTypeExecution,
		Severity: SeverityError,
		Caller:   captureCaller(callerSkip),
		Stack:    captureStackTrace(stackSkip),
	}
}

// WrapWithSkip wraps an existing error with an explicit frame skip increase.
func WrapWithSkip(skip int, err error, op string, code string) *AppError {
	callerSkip := DefaultCallerSkip + skip
	stackSkip := DefaultStackTraceSkip + skip
	return &AppError{
		Op:       op,
		Code:     code,
		Type:     ErrorTypeExecution,
		Severity: SeverityError,
		Caller:   captureCaller(callerSkip),
		Stack:    captureStackTrace(stackSkip),
		Cause:    err,
	}
}

// New creates a new AppError without an underlying cause.
func New(op string, code string, ctx map[string]any) *AppError {
	return &AppError{
		Op:       op,
		Code:     code,
		Type:     ErrorTypeExecution,
		Severity: SeverityError,
		Caller:   captureCaller(DefaultCallerSkip),
		Stack:    captureStackTrace(DefaultStackTraceSkip),
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
		Caller:   captureCaller(DefaultCallerSkip),
		Stack:    captureStackTrace(DefaultStackTraceSkip),
	}
}

// NewWithDetails creates a fully specified AppError without cause.
func NewWithDetails(op, code, msg, creator string, errType ErrorType, sev SeverityType, ctx map[string]any) *AppError {
	return &AppError{
		Op:       op,
		Code:     code,
		Type:     errType,
		Severity: sev,
		Creator:  creator,
		Message:  msg,
		Caller:   captureCaller(DefaultCallerSkip),
		Stack:    captureStackTrace(DefaultStackTraceSkip),
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
		Caller:   captureCaller(DefaultCallerSkip),
		Stack:    captureStackTrace(DefaultStackTraceSkip),
	}
}

// NewExecutionError creates an AppError specialized for execution failures.
func NewExecutionError(msg string) *AppError {
	return &AppError{
		Op:       "execution",
		Code:     "E9000",
		Type:     ErrorTypeExecution,
		Severity: SeverityError,
		Message:  msg,
		Caller:   captureCaller(DefaultCallerSkip),
		Stack:    captureStackTrace(DefaultStackTraceSkip),
	}
}

// Wrap wraps an existing error with an operation label and context.
func Wrap(err error, op string, ctx map[string]any) *AppError {
	return &AppError{
		Op:       op,
		Code:     "E9000",
		Type:     ErrorTypeExecution,
		Severity: SeverityError,
		Caller:   captureCaller(DefaultCallerSkip),
		Stack:    captureStackTrace(DefaultStackTraceSkip),
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
		Caller:   captureCaller(DefaultCallerSkip),
		Stack:    captureStackTrace(DefaultStackTraceSkip),
		Cause:    err,
	}
}

// WrapWithDetails wraps an existing error with full metadata.
func WrapWithDetails(err error, op, code, msg, creator string, errType ErrorType, sev SeverityType, ctx map[string]any) *AppError {
	return &AppError{
		Op:       op,
		Code:     code,
		Type:     errType,
		Severity: sev,
		Creator:  creator,
		Message:  msg,
		Caller:   captureCaller(DefaultCallerSkip),
		Stack:    captureStackTrace(DefaultStackTraceSkip),
		Ctx:      ctx,
		Cause:    err,
	}
}

// HasError reports whether an error exists.
func (e *AppError) HasError() bool {
	return e != nil
}

// IsSuccess reports whether no error exists (e == nil).
func (e *AppError) IsSuccess() bool {
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

// HandleError processes the error without terminating the application.
// It acts as a safety valve, checking for nil before processing.
func (e *AppError) HandleError() {
	if e == nil {
		return
	}

	// Proceed with error logging or handling logic here without CLI exit.
	// For now, it gracefully proceeds forward as requested.
}
