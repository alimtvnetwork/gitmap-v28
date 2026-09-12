package result_test

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

func TestResultMapSuccess(t *testing.T) {
	m := map[string]int{"b": 2, "a": 1, "c": 3}
	res := result.OkMap(m)

	if res.IsFailure() {
		t.Fatal("expected success")
	}

	if res.Count() != 3 {
		t.Fatalf("expected count 3, got %d", res.Count())
	}

	val, isFound := res.Get("a")
	if !isFound || val != 1 {
		t.Fatalf("expected 1, got %d", val)
	}

	keys := res.Keys()
	if len(keys) != 3 || keys[0] != "a" || keys[1] != "b" || keys[2] != "c" {
		t.Fatalf("expected sorted keys [a b c], got %v", keys)
	}
}

func TestResultMapFailure(t *testing.T) {
	appErr := apperror.NewSimple("map failed", "E_MAP_TEST")
	res := result.FailMap[string, int](appErr)

	if res.IsSuccess() {
		t.Fatal("expected failure")
	}

	if !res.HasError() {
		t.Fatal("expected has error")
	}

	if res.AppError() == nil {
		t.Fatal("expected non-nil app error")
	}

	_, isFound := res.Get("missing")
	if isFound {
		t.Fatal("expected false for missing key")
	}
}
