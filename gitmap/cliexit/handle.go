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
	if err == nil && len(defaultCode) > 0 {
		runFlushers()
		exitFunc(defaultCode[0])

		return
	}
	if err == nil {
		return
	}
	code := 1
	if len(defaultCode) > 0 {
		code = defaultCode[0]
	}
	var appErr *apperror.AppError
	if !errors.As(err, &appErr) || appErr == nil {
		appErr = apperror.WrapSimple(err, "cli")
	}
	if appErr == nil {
		return
	}
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
	if e.Message != "" {
		fmt.Fprintf(w, "gitmap: [%s:%s] %s: %s\n", e.Code, e.Type, e.Op, e.Message)
	} else {
		fmt.Fprintf(w, "gitmap: [%s:%s] %s\n", e.Code, e.Type, e.Op)
	}
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
	if shouldPrintStackTrace(e) && e.Stack != "" {
		fmt.Fprintf(w, "  stack trace:\n%s\n", indentLines(e.Stack, "    "))
	}
}

func shouldPrintStackTrace(e *apperror.AppError) bool {
	if os.Getenv("GITMAP_DEBUG") == "1" || os.Getenv("DEBUG") == "1" {
		return true
	}
	if e.Code == "E9000" || e.Type == apperror.ErrorTypeExecution || e.Severity == apperror.SeverityFatal {
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
