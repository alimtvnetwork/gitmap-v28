package appfault_test

import (
	"testing"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
)

func TestAppError_NilReceiverSafety(t *testing.T) {
	var nilErr *appfault.AppError

	// 1. Null, zero, and empty checks on nil receiver MUST NEVER PANIC
	if !nilErr.IsNull() {
		t.Fatal("expected IsNull() to be true on nil *AppError")
	}

	if !nilErr.IsEmpty() {
		t.Fatal("expected IsEmpty() to be true on nil *AppError")
	}

	if !nilErr.HasZero() {
		t.Fatal("expected HasZero() to be true on nil *AppError")
	}

	if !nilErr.IsZero() {
		t.Fatal("expected IsZero() to be true on nil *AppError")
	}

	if !nilErr.HasNull() {
		t.Fatal("expected HasNull() to be true on nil *AppError")
	}

	if !nilErr.HasNullError() {
		t.Fatal("expected HasNullError() to be true on nil *AppError")
	}

	if nilErr.HasError() {
		t.Fatal("expected HasError() to be false on nil *AppError")
	}

	if !nilErr.IsSuccess() {
		t.Fatal("expected IsSuccess() to be true on nil *AppError")
	}

	if nilErr.IsFailed() {
		t.Fatal("expected IsFailed() to be false on nil *AppError")
	}

	if !nilErr.IsValid() {
		t.Fatal("expected IsValid() to be true on nil *AppError")
	}

	if nilErr.IsInvalid() {
		t.Fatal("expected IsInvalid() to be false on nil *AppError")
	}

	// 2. Value getters on nil receiver MUST NEVER PANIC
	if nilErr.Code() != 0 {
		t.Fatalf("expected Code() == 0 on nil, got %d", nilErr.Code())
	}

	if nilErr.Type() != errtype.None {
		t.Fatalf("expected Type() == None on nil, got %v", nilErr.Type())
	}

	if nilErr.Message() != "" {
		t.Fatalf("expected Message() == '' on nil, got %q", nilErr.Message())
	}

	if nilErr.StatusCode() != 0 {
		t.Fatalf("expected StatusCode() == 0 on nil, got %d", nilErr.StatusCode())
	}

	if !nilErr.Caller().IsEmpty() {
		t.Fatal("expected Caller().IsEmpty() to be true on nil")
	}

	if nilErr.StackTrace() != nil {
		t.Fatal("expected StackTrace() == nil on nil")
	}

	if nilErr.Context().Count() != 0 {
		t.Fatal("expected Context().Count() == 0 on nil")
	}

	if nilErr.LoopCount() != 1 {
		t.Fatalf("expected LoopCount() == 1 on nil, got %d", nilErr.LoopCount())
	}

	// 3. Formatting and printing on nil receiver MUST NEVER PANIC
	if nilErr.FormatStdout() != "" {
		t.Fatalf("expected FormatStdout() == '' on nil, got %q", nilErr.FormatStdout())
	}

	if nilErr.FormatJson() != "{}" {
		t.Fatalf("expected FormatJson() == '{}' on nil, got %q", nilErr.FormatJson())
	}

	if nilErr.FormatTextLog() != "" {
		t.Fatalf("expected FormatTextLog() == '' on nil, got %q", nilErr.FormatTextLog())
	}

	nilErr.Print()
	nilErr.PrintStdout()
	nilErr.PrintJson()
	nilErr.PrintLog()

	// 4. Clone and Concat on nil receiver MUST NEVER PANIC
	cloned := nilErr.Clone()
	if cloned != nil {
		t.Fatal("expected nilErr.Clone() == nil")
	}

	otherErr := appfault.New(errtype.Validation, "input invalid")
	concatenated := nilErr.Concat(otherErr)
	if concatenated != otherErr {
		t.Fatal("expected nilErr.Concat(otherErr) to return otherErr")
	}

	concatNil := otherErr.Concat(nilErr)
	if concatNil != otherErr {
		t.Fatal("expected otherErr.Concat(nilErr) to return otherErr")
	}
}

func TestResult_NullSafety(t *testing.T) {
	// Zero-value result
	var r appfault.Result[string]

	if !r.IsNull() {
		t.Fatal("expected r.IsNull() to be true on zero Result")
	}

	if !r.IsEmpty() {
		t.Fatal("expected r.IsEmpty() to be true on zero Result")
	}

	if !r.HasZero() {
		t.Fatal("expected r.HasZero() to be true on zero Result")
	}

	if !r.IsZero() {
		t.Fatal("expected r.IsZero() to be true on zero Result")
	}

	if !r.HasNull() {
		t.Fatal("expected r.HasNull() to be true on zero Result")
	}

	if !r.IsSuccess() {
		t.Fatal("expected r.IsSuccess() to be true on zero Result")
	}

	if r.IsFailed() {
		t.Fatal("expected r.IsFailed() to be false on zero Result")
	}

	// Clone on zero-value Result
	cloned := r.Clone()
	if !cloned.IsNull() || !cloned.IsSuccess() {
		t.Fatal("cloned zero Result should remain zero")
	}

	// Concat on results
	rFail := appfault.FailureResult[string](appfault.New(errtype.NotFound, "record not found"))
	merged := r.Concat(rFail)
	if !merged.IsFailed() {
		t.Fatal("expected merged result to be failed")
	}
}
