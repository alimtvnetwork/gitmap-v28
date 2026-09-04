package appfault_test

import (
	"testing"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
	"coding-guidelines/common/pkg/result"
)

func TestAppErrorCreationAndNilSafety(t *testing.T) {
	var nilErr *appfault.AppError
	if nilErr.HasError() || !nilErr.IsSuccess() || appfault.New(errtype.None, "") != nil {
		t.Fatal("expected nil AppError and None constructor to be IsSuccess")
	}

	appErr := appfault.New(errtype.NotFound, "record not found").WithOp("repo.find").WithSeverity(appfault.SeverityWarn)
	if !appErr.Is(errtype.NotFound) || !appErr.HasValidError() {
		t.Fatalf("expected NotFound error, got %v", appErr.GetType())
	}
}

func TestContextMapOperations(t *testing.T) {
	cm := appfault.NewContextMap().Set("siteId", 101).Set("slug", "test-plugin")
	if !cm.Has("siteId") || cm.GetString("siteId") != "101" || cm.Count() != 2 {
		t.Fatalf("expected siteId=101 and count 2, got %s (count=%d)", cm.GetString("siteId"), cm.Count())
	}

	cm.Remove("slug")
	if cm.Has("slug") {
		t.Fatal("expected slug to be removed")
	}
}

func TestStackFrameAndCaller(t *testing.T) {
	frame := appfault.NewStackFrame("main.run", "main.go", 42)
	if frame.Function != "main.run" || frame.File != "main.go" || frame.Line != 42 {
		t.Fatalf("unexpected frame: %+v", frame)
	}

	caller := appfault.CaptureCaller(0)
	if len(caller) == 0 {
		t.Fatal("expected non-empty caller")
	}
}

func TestCallerAndStackTraceObjects(t *testing.T) {
	appErr := appfault.New(errtype.Database, "db err")
	caller := appErr.Caller()
	if caller.IsEmpty() || caller.Line == 0 {
		t.Fatalf("expected caller info to have line and file: %+v", caller)
	}

	stack := appErr.StackTrace()
	if len(stack) == 0 {
		t.Fatal("expected non-empty stack trace")
	}
}

func TestResultMonadicOperations(t *testing.T) {
	res := result.SuccessResult("data-payload")
	if !res.IsSuccess() || res.IsFailed() || res.Data() != "data-payload" {
		t.Fatal("expected success result")
	}

	failRes := result.NewFailureWithType[string](errtype.Validation, "bad input", "validator")
	if !failRes.IsFailed() || failRes.UnwrapOr("fallback") != "fallback" {
		t.Fatal("expected failed result with fallback")
	}
}

func TestResultSliceAndMap(t *testing.T) {
	sliceRes := appfault.OkSlice([]string{"alpha", "beta"})
	if sliceRes.Count() != 2 || !sliceRes.HasItems() {
		t.Fatalf("expected count 2, got %d", sliceRes.Count())
	}

	mapRes := appfault.OkMap(map[string]int{"key1": 100})
	if !mapRes.Has("key1") || mapRes.Count() != 1 {
		t.Fatal("expected key1 to exist")
	}
}
