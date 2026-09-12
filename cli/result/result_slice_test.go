package result_test

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

func TestResultSliceSuccess(t *testing.T) {
	s := []string{"apple", "banana"}
	res := result.OkSlice(s)

	if res.IsFailure() {
		t.Fatal("expected success")
	}

	if res.Count() != 2 {
		t.Fatalf("expected count 2, got %d", res.Count())
	}

	val, isFound := res.Get(0)
	if !isFound || val != "apple" {
		t.Fatalf("expected apple, got %s", val)
	}

	_, isInvalid := res.Get(5)
	if isInvalid {
		t.Fatal("expected out of bounds false")
	}
}

func TestResultSliceFailure(t *testing.T) {
	appErr := apperror.NewSimple("slice failed", "E_SLICE_TEST")
	res := result.FailSlice[string](appErr)

	if res.IsSuccess() {
		t.Fatal("expected failure")
	}

	if !res.HasError() {
		t.Fatal("expected has error")
	}

	if res.AppError() == nil {
		t.Fatal("expected non-nil app error")
	}
}
