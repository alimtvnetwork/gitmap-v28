package appfault

import (
	"strings"
	"testing"

	"coding-guidelines/common/pkg/errtype"
)

func TestTryDo(t *testing.T) {
	executedTry := false
	executedCatch := false
	executedFinally := false

	Block{
		Try: func() {
			executedTry = true
			panic("something went wrong")
		},
		Catch: func(e Exception) {
			executedCatch = true
			if e.(string) != "something went wrong" {
				t.Errorf("expected string panic value")
			}
		},
		Finally: func() {
			executedFinally = true
		},
	}.Do()

	if !executedTry || !executedCatch || !executedFinally {
		t.Errorf("Try/Catch/Finally block did not execute all components properly")
	}
}

func TestRecover(t *testing.T) {
	var recoveredErr *AppError

	func() {
		defer Recover(func(err *AppError) {
			recoveredErr = err
		})

		panic("critical failure")
	}()

	if recoveredErr == nil {
		t.Fatalf("Recover failed to catch panic")
	}

	if !recoveredErr.Is(errtype.Execution) {
		t.Errorf("Expected errtype.Execution, got %v", recoveredErr.errType)
	}

	if recoveredErr.Unwrap() == nil || recoveredErr.Unwrap().Error() != "critical failure" {
		t.Errorf("Expected cause to be 'critical failure'")
	}
	
	if !strings.Contains(recoveredErr.GetMessage(), "unhandled panic recovered") {
		t.Errorf("Expected panic message, got %s", recoveredErr.GetMessage())
	}
}
