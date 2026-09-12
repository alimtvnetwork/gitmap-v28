package cliexit

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func TestWriteAppErrorReport_FormatsCompleteDiagnostics(t *testing.T) {
	buf := &bytes.Buffer{}
	cause := errors.New("file not accessible")
	appErr := apperror.WrapWithDetails(
		cause,
		"config.read",
		"E3001",
		"failed to read configuration",
		"cfgloader",
		apperror.ErrorTypeNotFound,
		apperror.SeverityError,
		map[string]any{"path": "/etc/gitmap.json"},
	)

	WriteAppErrorReport(buf, appErr)
	rendered := buf.String()

	if !strings.Contains(rendered, "[E3001:NOT_FOUND]") {
		t.Errorf("missing code or type in output: %s", rendered)
	}

	if !strings.Contains(rendered, "creator: cfgloader") {
		t.Errorf("missing creator attribution: %s", rendered)
	}

	if !strings.Contains(rendered, "path:/etc/gitmap.json") && !strings.Contains(rendered, "path: /etc/gitmap.json") {
		t.Errorf("missing context path: %s", rendered)
	}

	if !strings.Contains(rendered, "file not accessible") {
		t.Errorf("missing underlying cause: %s", rendered)
	}
}

func TestHandleError_CustomExitCode(t *testing.T) {
	var capturedCode int
	prev := SetExitFunc(func(c int) {
		capturedCode = c
	})
	defer SetExitFunc(prev)

	appErr := apperror.NewSimple("test.fail", "E9999")
	HandleError(appErr, 42)

	if capturedCode != 42 {
		t.Fatalf("expected exit code 42, got %d", capturedCode)
	}
}

func TestHandleError_NilErrorWithCode(t *testing.T) {
	var capturedCode int
	prev := SetExitFunc(func(c int) {
		capturedCode = c
	})
	defer SetExitFunc(prev)

	HandleError(nil, 7)
	if capturedCode != 7 {
		t.Fatalf("expected exit code 7, got %d", capturedCode)
	}
}

func TestHandleError_NilErrorWithoutCode(t *testing.T) {
	var isCalled bool
	prev := SetExitFunc(func(c int) {
		isCalled = true
	})
	defer SetExitFunc(prev)

	HandleError(nil)
	if isCalled {
		t.Fatal("expected exitFunc not to be called for nil error without code")
	}
}

func TestFail_InvokesExitFunc(t *testing.T) {
	var capturedCode int
	prev := SetExitFunc(func(c int) {
		capturedCode = c
	})
	defer SetExitFunc(prev)

	Fail("testcmd", "testop", "testsubj", errors.New("test err"), 3)
	if capturedCode != 3 {
		t.Fatalf("expected exit code 3, got %d", capturedCode)
	}
}

func TestExit_InvokesExitFunc(t *testing.T) {
	var capturedCode int
	prev := SetExitFunc(func(c int) {
		capturedCode = c
	})
	defer SetExitFunc(prev)

	Exit(9)
	if capturedCode != 9 {
		t.Fatalf("expected exit code 9, got %d", capturedCode)
	}
}

func TestFailWith_InvokesExitFunc(t *testing.T) {
	var capturedCode int
	prev := SetExitFunc(func(c int) {
		capturedCode = c
	})
	defer SetExitFunc(prev)

	ctx := Context{
		Command: "testcmd",
		Op:      "testop",
		Err:     errors.New("test err"),
	}

	FailWith(ctx, OutputHuman, 5)
	if capturedCode != 5 {
		t.Fatalf("expected exit code 5, got %d", capturedCode)
	}
}
