package appfault_test

import (
	"errors"
	"strings"
	"testing"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
)

func TestMerge_NilAndEmptyScenarios(t *testing.T) {
	// 1. Both nil
	if appfault.Merge(nil, nil) != nil {
		t.Fatal("expected nil when merging two nil errors")
	}

	// 2. Both errtype.None
	noneErr1 := appfault.New(errtype.None, "")
	noneErr2 := appfault.New(errtype.None, "")
	if appfault.Merge(noneErr1, noneErr2) != nil {
		t.Fatal("expected nil when merging two None errors")
	}

	// 3. Prev is nil, next is actual error
	realErr := appfault.New(errtype.Validation, "field required")
	if appfault.Merge(nil, realErr) != realErr {
		t.Fatal("expected realErr when prev is nil")
	}

	// 4. Next is nil, prev is actual error
	if appfault.Merge(realErr, nil) != realErr {
		t.Fatal("expected realErr when next is nil")
	}
}

func TestMerge_TwoErrorsWithStackTraceTracking(t *testing.T) {
	// Simulate loop attempt 1 failure
	err1 := appfault.New(errtype.Network, "connection refused").
		WithContext("server", "node-alpha")

	// Simulate loop attempt 2 failure
	err2 := appfault.Wrap(errtype.Timeout, errors.New("deadline exceeded"), "gateway timed out").
		WithContext("gateway", "gw-beta")

	merged := appfault.Merge(err1, err2)

	// Verify merged error identity
	if merged.Type() != errtype.Timeout {
		t.Fatalf("expected merged error type Timeout, got %v", merged.Type())
	}

	if merged.LoopCount() != 2 {
		t.Fatalf("expected loop count 2, got %d", merged.LoopCount())
	}

	// Verify first error stack trace is captured in context dictionary
	firstStack := merged.FirstErrorStackTrace()
	if len(firstStack) == 0 {
		t.Fatal("expected non-empty FirstErrorStackTrace in merged error")
	}

	if merged.Context().GetString("FirstErrorMessage") != "connection refused" {
		t.Fatalf("expected FirstErrorMessage 'connection refused', got %s", merged.Context().GetString("FirstErrorMessage"))
	}

	if merged.Context().GetString("FirstErrorType") != "Network" {
		t.Fatalf("expected FirstErrorType 'Network', got %s", merged.Context().GetString("FirstErrorType"))
	}

	// Verify prior stack trace is captured
	priorStack := merged.PriorStackTrace()
	if len(priorStack) == 0 {
		t.Fatal("expected non-empty PriorStackTrace in merged error")
	}

	// Verify domain context preservation
	if !merged.Context().Has("server") || !merged.Context().Has("gateway") {
		t.Fatalf("expected merged context to retain server and gateway, got: %v", merged.Context().Format())
	}

	// Verify immutability: err1 and err2 were not modified
	if err1.Context().Has("LoopCount") || err2.Context().Has("LoopCount") {
		t.Fatal("original errors were mutated during merge!")
	}
}

func TestMerge_MultiLoopAccumulation(t *testing.T) {
	var accumulated *appfault.AppError

	attempts := []struct {
		errType errtype.Variation
		message string
	}{
		{errtype.Network, "dns resolution failed"},
		{errtype.Database, "connection pool exhausted"},
		{errtype.Timeout, "read timeout on socket"},
	}

	// Simulate a 3-iteration retry loop
	for i, att := range attempts {
		currentErr := appfault.New(att.errType, att.message).
			WithContext("attempt", i+1)

		accumulated = appfault.Merge(accumulated, currentErr)
	}

	// Verify loop count reflects 3 attempts
	if accumulated.LoopCount() != 3 {
		t.Fatalf("expected loop count 3, got %d", accumulated.LoopCount())
	}

	// Verify first error stack trace is specifically preserved from attempt 1
	if accumulated.Context().GetString("FirstErrorMessage") != "dns resolution failed" {
		t.Fatalf("expected FirstErrorMessage from attempt 1, got %s", accumulated.Context().GetString("FirstErrorMessage"))
	}

	// Verify prior error message is specifically from attempt 2
	if accumulated.Context().GetString("PriorMessage") != "connection pool exhausted" {
		t.Fatalf("expected PriorMessage from attempt 2, got %s", accumulated.Context().GetString("PriorMessage"))
	}

	// Verify current failure message is from attempt 3
	if accumulated.Message() != "read timeout on socket" {
		t.Fatalf("expected current error message from attempt 3, got %s", accumulated.Message())
	}

	// Verify stack trace history contains entries for all attempts
	history := accumulated.StackTraceHistory()
	if len(history) != 3 {
		t.Fatalf("expected 3 history entries, got %d", len(history))
	}

	if !strings.Contains(history[0], "Attempt #1") || !strings.Contains(history[0], "dns resolution failed") {
		t.Fatalf("history[0] should contain Attempt #1: %s", history[0])
	}

	if !strings.Contains(history[1], "Attempt #2") || !strings.Contains(history[1], "connection pool exhausted") {
		t.Fatalf("history[1] should contain Attempt #2: %s", history[1])
	}

	if !strings.Contains(history[2], "Attempt #3") || !strings.Contains(history[2], "read timeout on socket") {
		t.Fatalf("history[2] should contain Attempt #3: %s", history[2])
	}
}
