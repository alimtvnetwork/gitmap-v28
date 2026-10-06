package cliexit

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

var exitFunc = os.Exit

// SetExitFunc overrides the process exit function for testing.
func SetExitFunc(fn func(int)) func(int) {
	prev := exitFunc
	exitFunc = fn

	return prev
}

// HandleError processes an error through centralized logging, flushing,
// and process exit (or panic if debug mode is active).
func HandleError(err error, defaultCode ...int) {
	if err == nil {
		handleNilError(defaultCode...)

		return
	}

	code := resolveExitCode(defaultCode...)
	if isAbortOrReported(err) {
		runFlushers()
		exitFunc(code)

		return
	}

	appErr := ensureAppError(err)
	if appErr == nil {
		return
	}

	if isAbortOrReported(appErr) {
		runFlushers()
		exitFunc(code)

		return
	}

	dispatchError(appErr, code)
}

func handleNilError(defaultCode ...int) {
	if len(defaultCode) > 0 {
		runFlushers()
		exitFunc(defaultCode[0])
	}
}

func resolveExitCode(defaultCode ...int) int {
	if len(defaultCode) > 0 {
		return defaultCode[0]
	}

	return int(ExitCodeGeneralError)
}

func ensureAppError(err error) *apperror.AppError {
	var appErr *apperror.AppError
	hasAppError := errors.As(err, &appErr) && appErr != nil
	if !hasAppError {
		return apperror.WrapSimple(err, "cli")
	}

	return appErr
}

func dispatchError(appErr *apperror.AppError, code int) {
	WriteAppErrorReport(os.Stderr, appErr)
	runFlushers()
	if os.Getenv("GITMAP_ERROR_PANIC") == "1" {
		panic(appErr)
	}

	exitFunc(code)
}

// FailAppError is a semantic alias for HandleError with a specific exit code.
func FailAppError(appErr *apperror.AppError, code int) {
	HandleError(appErr, code)
}

// HandleSuccess flushes output pipes and exits cleanly with ExitCodeSuccess.
func HandleSuccess() {
	runFlushers()
	exitFunc(int(ExitCodeSuccess))
}

// HandleUsageError reports a CLI usage or syntax error and exits with ExitCodeUsageError.
func HandleUsageError(err error) {
	HandleError(err, int(ExitCodeUsageError))
}

// HandleValidationError reports a validation failure and exits with ExitCodeValidationError.
func HandleValidationError(err error) {
	HandleError(err, int(ExitCodeValidationError))
}

// HandleGeneralError reports an operational failure and exits with ExitCodeGeneralError.
func HandleGeneralError(err error) {
	HandleError(err, int(ExitCodeGeneralError))
}

// HandleNotFound reports a missing resource failure and exits with ExitCodeNotFound.
func HandleNotFound(err error) {
	HandleError(err, int(ExitCodeNotFound))
}

// HandleAppError explicitly handles an AppError instance with full diagnostics.
func HandleAppError(appErr *apperror.AppError, defaultCode ...int) {
	HandleError(appErr, defaultCode...)
}

// WriteAppErrorReport writes formatted structured error diagnostics.
func WriteAppErrorReport(w io.Writer, e *apperror.AppError) {
	if e == nil {
		return
	}

	writeErrorHeader(w, e)
	writeErrorMetadata(w, e)
	writeErrorStack(w, e)
}

func writeErrorHeader(w io.Writer, e *apperror.AppError) {
	hasMessage := e.Message != ""
	if hasMessage {
		fmt.Fprintf(w, "gitmap: [%s:%s] %s: %s\n", e.Code, e.Type, e.Op, e.Message)

		return
	}

	fmt.Fprintf(w, "gitmap: [%s:%s] %s\n", e.Code, e.Type, e.Op)
}

func writeErrorMetadata(w io.Writer, e *apperror.AppError) {
	if e.Caller != "" {
		fmt.Fprintf(w, "  origin: %s\n", e.Caller)
	}

	if e.Creator != "" {
		fmt.Fprintf(w, "  creator: %s\n", e.Creator)
	}

	if len(e.Ctx) > 0 {
		fmt.Fprintf(w, "  context: %v\n", e.Ctx)
	}

	if e.Cause != nil {
		fmt.Fprintf(w, "  cause: %v\n", e.Cause)
	}
}

func writeErrorStack(w io.Writer, e *apperror.AppError) {
	hasStack := isStackTraceEnabled(e) && e.Stack != ""
	if hasStack {
		fmt.Fprintf(w, "  stack trace:\n%s\n", indentLines(e.Stack, "    "))
	}
}

func isStackTraceEnabled(e *apperror.AppError) bool {
	isDebug := os.Getenv("GITMAP_DEBUG") == "1" || os.Getenv("DEBUG") == "1"
	if isDebug {
		return true
	}

	isAbort := e.Type == apperror.ErrorTypeAbort || isReportedContext(e.Ctx)
	if isAbort {
		return false
	}

	isFatal := e.Code == "E9000" || e.Code == "E_INTERNAL_ERROR" || e.Type == apperror.ErrorTypeExecution || e.Severity == apperror.SeverityFatal
	if isFatal {
		return true
	}

	return false
}

func isAbortOrReported(err error) bool {
	if err == nil {
		return false
	}

	appErr, isApp := err.(*apperror.AppError)
	if isApp && appErr != nil {
		return checkAppErrorAbort(appErr)
	}

	multi, isMulti := err.(interface{ Unwrap() []error })
	if isMulti && multi != nil {
		return checkMultiErrorAbort(multi.Unwrap())
	}

	single, isSingle := err.(interface{ Unwrap() error })
	if isSingle && single != nil {
		return isAbortOrReported(single.Unwrap())
	}

	return false
}

func checkAppErrorAbort(appErr *apperror.AppError) bool {
	isAbort := appErr.Type == apperror.ErrorTypeAbort || isReportedContext(appErr.Ctx)
	if isAbort {
		return true
	}

	return isAbortOrReported(appErr.Cause)
}

func checkMultiErrorAbort(errs []error) bool {
	for _, e := range errs {
		if isAbortOrReported(e) {
			return true
		}
	}

	return false
}

func isReportedContext(ctx map[string]any) bool {
	if ctx == nil {
		return false
	}

	val, ok := ctx["reported"]
	if !ok {
		return false
	}

	isReported, isBool := val.(bool)

	return isBool && isReported
}

func indentLines(s, prefix string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	for i, line := range lines {
		lines[i] = prefix + line
	}

	return strings.Join(lines, "\n")
}
