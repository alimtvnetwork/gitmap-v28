package cliexit

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
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
	appErr := ensureAppError(err)
	if appErr == nil {
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

	isFatal := e.Code == "E9000" || e.Type == apperror.ErrorTypeExecution || e.Severity == apperror.SeverityFatal
	if isFatal {
		return true
	}

	return false
}

func indentLines(s, prefix string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	for i, line := range lines {
		lines[i] = prefix + line
	}

	return strings.Join(lines, "\n")
}
