package result_test

import (
	"errors"
	"testing"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
	"coding-guidelines/common/pkg/result"
)

func TestWrapSuccess(t *testing.T) {
	val := "hello-world"
	wrap := result.WrapSuccess(val)

	if wrap.IsFailed() {
		t.Fatalf("expected success, got failure")
	}

	if wrap.Data() != val {
		t.Fatalf("expected %s, got %s", val, wrap.Data())
	}
}

func TestSuccessShortForm(t *testing.T) {
	val := 12345
	wrap := result.Success(val)

	if wrap.IsFailed() {
		t.Fatalf("expected success, got failure")
	}

	if wrap.Data() != val {
		t.Fatalf("expected %d, got %d", val, wrap.Data())
	}
}

func TestWrapFailure(t *testing.T) {
	appErr := appfault.New(errtype.IO, "disk read error")
	wrap := result.WrapFailure[string](appErr)

	if wrap.IsSuccess() {
		t.Fatalf("expected failure, got success")
	}

	if wrap.Fault() == nil {
		t.Fatalf("expected non-nil fault")
	}

	if wrap.Fault().Message() != "disk read error" {
		t.Fatalf("unexpected message: %s", wrap.Fault().Message())
	}
}

func TestFailureShortForm(t *testing.T) {
	appErr := appfault.New(errtype.Network, "connection refused")
	wrap := result.Failure[int](appErr)

	if wrap.IsSuccess() {
		t.Fatalf("expected failure, got success")
	}

	if wrap.Fault() == nil {
		t.Fatalf("expected non-nil fault")
	}

	if wrap.Fault().Message() != "connection refused" {
		t.Fatalf("unexpected message: %s", wrap.Fault().Message())
	}
}

func TestWrapFormat(t *testing.T) {
	appErr := appfault.New(errtype.Validation, "invalid configuration")
	wrap := result.WrapFailure[int](appErr)

	// Default formatting check
	formatted := wrap.Format(nil)
	if len(formatted) == 0 {
		t.Fatalf("expected non-empty formatted string")
	}

	// Custom formatting check
	custom := wrap.Format(func(r result.Wrap[int]) string {
		return "CUSTOM-ERROR: " + r.Fault().Message()
	})

	if custom != "CUSTOM-ERROR: invalid configuration" {
		t.Fatalf("unexpected custom format: %s", custom)
	}

	// Success formatting check
	successWrap := result.WrapSuccess(42)
	successFormatted := successWrap.Format(nil)
	if successFormatted != "✅ [OK] 42" {
		t.Fatalf("unexpected success format: %s", successFormatted)
	}
}

func TestWrapFailureWithId(t *testing.T) {
	w := result.WrapFailureWithId[string](errtype.Validation, "bad param")
	if w.IsSuccess() {
		t.Fatalf("expected failure")
	}

	if w.Fault().Type() != errtype.Validation {
		t.Fatalf("unexpected type: %v", w.Fault().Type())
	}

	if w.Fault().Message() != "bad param" {
		t.Fatalf("unexpected message: %s", w.Fault().Message())
	}
}

func TestFailureWithIdShortForm(t *testing.T) {
	w := result.FailureWithId[int](errtype.Unauthorized, "unauthorized access")
	if w.IsSuccess() {
		t.Fatalf("expected failure")
	}

	if w.Fault().Type() != errtype.Unauthorized {
		t.Fatalf("unexpected type: %v", w.Fault().Type())
	}

	if w.Fault().Message() != "unauthorized access" {
		t.Fatalf("unexpected message: %s", w.Fault().Message())
	}
}

func TestWrapFailureWithCause(t *testing.T) {
	rawErr := errors.New("timeout reached")
	w := result.WrapFailureWithCause[string](errtype.Timeout, rawErr, "operation timed out")
	if w.IsSuccess() {
		t.Fatalf("expected failure")
	}

	if w.Fault().Type() != errtype.Timeout {
		t.Fatalf("unexpected type: %v", w.Fault().Type())
	}

	if w.Fault().Message() != "operation timed out" {
		t.Fatalf("unexpected message: %s", w.Fault().Message())
	}
}

func TestFailureWithCauseShortForm(t *testing.T) {
	rawErr := errors.New("db disconnect")
	w := result.FailureWithCause[string](errtype.Database, rawErr, "query failed")
	if w.IsSuccess() {
		t.Fatalf("expected failure")
	}

	if w.Fault().Type() != errtype.Database {
		t.Fatalf("unexpected type: %v", w.Fault().Type())
	}

	if w.Fault().Message() != "query failed" {
		t.Fatalf("unexpected message: %s", w.Fault().Message())
	}
}

func TestWrapFailureFromWrap(t *testing.T) {
	first := result.WrapFailureWithId[int](errtype.IO, "io failed")
	second := result.WrapFailureFromWrap[string](first)
	if second.IsSuccess() {
		t.Fatalf("expected failure in propagated wrap")
	}

	if second.Fault().Message() != "io failed" {
		t.Fatalf("unexpected message: %s", second.Fault().Message())
	}
}

func TestFailureFromWrapShortForm(t *testing.T) {
	first := result.FailureWithId[int](errtype.NotFound, "record not found")
	second := result.FailureFromWrap[string](first)
	if second.IsSuccess() {
		t.Fatalf("expected failure in propagated wrap")
	}

	if second.Fault().Message() != "record not found" {
		t.Fatalf("unexpected message: %s", second.Fault().Message())
	}
}

func TestWrapFailureFromError(t *testing.T) {
	appErr := appfault.New(errtype.Internal, "kernel panic simulation")
	w := result.WrapFailureFromError[bool](appErr)
	if w.IsSuccess() {
		t.Fatalf("expected failure")
	}

	if w.Fault().Message() != "kernel panic simulation" {
		t.Fatalf("unexpected message: %s", w.Fault().Message())
	}
}

func TestLegacyCompatibility(t *testing.T) {
	// SuccessResult
	s := result.SuccessResult("legacy-data")
	if s.IsFailed() {
		t.Fatalf("expected legacy success")
	}

	if s.Value != "legacy-data" {
		t.Fatalf("expected legacy-data, got %s", s.Value)
	}

	// NewSuccess
	ns := result.NewSuccess(99)
	if ns.IsFailed() {
		t.Fatalf("expected new success")
	}

	if ns.Data() != 99 {
		t.Fatalf("expected 99, got %d", ns.Data())
	}

	// FailureResult
	appErr := appfault.New(errtype.Validation, "legacy validation err")
	f := result.FailureResult[string](appErr)
	if f.IsSuccess() {
		t.Fatalf("expected legacy failure")
	}

	// NewFailure
	rawErr := errors.New("raw err")
	nf := result.NewFailure[string](errtype.IO, rawErr)
	if nf.IsSuccess() {
		t.Fatalf("expected NewFailure to fail")
	}

	// NewFailureWithType
	nft := result.NewFailureWithType[string](errtype.NotFound, "missing item", "testCaller")
	if nft.IsSuccess() {
		t.Fatalf("expected NewFailureWithType to fail")
	}
}
