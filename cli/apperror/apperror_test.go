package apperror

import (
	"errors"
	"strings"
	"testing"
)

func TestAppError_FormattingAndUnwrap(t *testing.T) {
	cause := errors.New("underlying socket closed")
	appErr := WrapWithDetails(
		cause,
		"network.dial",
		"E1001",
		"connection failed",
		"netutil",
		ErrorTypeExecution,
		SeverityFatal,
		map[string]any{"port": 8080},
	)

	hasCause := errors.Is(appErr, cause)
	if !hasCause {
		t.Fatalf("expected errors.Is to match cause")
	}

	rendered := appErr.Error()
	if !strings.Contains(rendered, "E1001") {
		t.Errorf("rendered error missing code: %s", rendered)
	}

	if !strings.Contains(rendered, "creator=netutil") {
		t.Errorf("rendered error missing creator: %s", rendered)
	}

	if !strings.Contains(rendered, "port:8080") && !strings.Contains(rendered, "port: 8080") {
		t.Errorf("rendered error missing context: %s", rendered)
	}
}

func TestAppError_WithContext(t *testing.T) {
	appErr := NewSimple("fs.open", "E2001").WithContext("file", "test.txt")
	val, isFound := appErr.Ctx["file"]
	if !isFound || val != "test.txt" {
		t.Fatalf("expected context key to be set")
	}
}

func TestAppError_StackTrace(t *testing.T) {
	appErr := NewSimple("test.op", "E3001")
	if appErr.Stack == "" {
		t.Fatalf("expected stack trace to be non-empty")
	}

	if !strings.Contains(appErr.Stack, "TestAppError_StackTrace") {
		t.Errorf("expected stack trace to contain test function name: %s", appErr.Stack)
	}

	if strings.Contains(appErr.Stack, "apperror.NewSimple") {
		t.Errorf("expected stack trace to omit apperror.NewSimple: %s", appErr.Stack)
	}
}

func TestAppError_StackTraceOmission(t *testing.T) {
	appErr := NewSimple("test.op", "E3001")
	if strings.Contains(appErr.Stack, "apperror.go:") {
		t.Errorf("stack trace must omit internal apperror.go frames: %s", appErr.Stack)
	}

	if !strings.Contains(appErr.Caller, "apperror_test.go") {
		t.Errorf("caller must point to test caller: %s", appErr.Caller)
	}
}

func testHelperCreateError(skip int) *AppError {
	return NewSimple("helper.op", "E3002").WithSkip(skip)
}

func TestAppError_WithSkip(t *testing.T) {
	errNoSkip := testHelperCreateError(0)
	errWithSkip := testHelperCreateError(1)

	if !strings.Contains(errNoSkip.Caller, "apperror_test.go") {
		t.Errorf("expected helper caller in apperror_test.go: %s", errNoSkip.Caller)
	}

	if errWithSkip.Caller == "" {
		t.Errorf("expected non-empty caller with skip: %s", errWithSkip.Caller)
	}
}

func TestAppError_DefaultSkipSetters(t *testing.T) {
	origStackSkip := DefaultStackTraceSkip
	origCallerSkip := DefaultCallerSkip
	defer SetDefaultStackTraceSkip(origStackSkip)
	defer SetDefaultCallerSkip(origCallerSkip)

	SetDefaultStackTraceSkip(4)
	SetDefaultCallerSkip(3)
	if DefaultStackTraceSkip != 4 || DefaultCallerSkip != 3 {
		t.Fatalf("expected skip settings updated")
	}
}

func TestAppError_NewAndWrapWithSkip(t *testing.T) {
	err1 := NewWithSkip(1, "custom.op", "E4001")
	if err1.Code != "E4001" || err1.Op != "custom.op" {
		t.Errorf("NewWithSkip mismatch: %+v", err1)
	}

	cause := errors.New("underlying failure")
	err2 := WrapWithSkip(1, cause, "wrapped.op", "E4002")
	if err2.Cause != cause || err2.Code != "E4002" {
		t.Errorf("WrapWithSkip mismatch: %+v", err2)
	}
}

func TestAppError_NewNotFoundError(t *testing.T) {
	msg := "no repo found matching 'gitamp'"
	err := NewNotFoundError(msg)

	if err.Type != ErrorTypeNotFound {
		t.Errorf("expected ErrorTypeNotFound, got %s", err.Type)
	}
	if err.Code != "E1004" {
		t.Errorf("expected code E1004, got %s", err.Code)
	}
	if err.Message != msg {
		t.Errorf("expected message %q, got %q", msg, err.Message)
	}
	if err.Stack != "" {
		t.Errorf("expected empty stack trace for not found, got %s", err.Stack)
	}
}

func TestAppError_NewNotFound(t *testing.T) {
	err := NewNotFound("tool_locate", "E_TOOL_NOT_FOUND", "vcvarsall.bat")

	if err.Type != ErrorTypeNotFound {
		t.Errorf("expected ErrorTypeNotFound, got %s", err.Type)
	}
	if err.Op != "tool_locate" {
		t.Errorf("expected op 'tool_locate', got %s", err.Op)
	}
	if err.Code != "E_TOOL_NOT_FOUND" {
		t.Errorf("expected code E_TOOL_NOT_FOUND, got %s", err.Code)
	}
	if err.Message != "vcvarsall.bat" {
		t.Errorf("expected message 'vcvarsall.bat', got %s", err.Message)
	}
	if err.Stack != "" {
		t.Errorf("expected empty stack for not found error, got %s", err.Stack)
	}
}
