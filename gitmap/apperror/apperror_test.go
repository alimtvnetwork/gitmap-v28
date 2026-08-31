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
	val, ok := appErr.Ctx["file"]
	if !ok || val != "test.txt" {
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
}
