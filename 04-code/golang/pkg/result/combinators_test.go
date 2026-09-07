package result

import (
	"strconv"
	"testing"

	"coding-guidelines/common/pkg/errtype"
)

func TestMap(t *testing.T) {
	succ := Success(42)
	mapped := Map(succ, func(v int) string {
		return strconv.Itoa(v)
	})

	if mapped.IsFailure() || mapped.Data() != "42" {
		t.Errorf("Map failed to transform value correctly")
	}

	fail := FailureWithId[int](errtype.Validation, "bad int")
	mappedFail := Map(fail, func(v int) string {
		return strconv.Itoa(v)
	})

	if mappedFail.IsSuccess() || mappedFail.Fault().GetMessage() != "bad int" {
		t.Errorf("Map failed to propagate error correctly")
	}
}

func TestFlatMap(t *testing.T) {
	succ := Success(10)
	flatMapped := FlatMap(succ, func(v int) Result[string] {
		return Success(strconv.Itoa(v * 2))
	})

	if flatMapped.IsFailure() || flatMapped.Data() != "20" {
		t.Errorf("FlatMap failed to transform value correctly")
	}

	flatMappedFail := FlatMap(succ, func(v int) Result[string] {
		return FailureWithId[string](errtype.Execution, "exec fail")
	})

	if flatMappedFail.IsSuccess() || flatMappedFail.Fault().GetMessage() != "exec fail" {
		t.Errorf("FlatMap failed to propagate inner error")
	}

	fail := FailureWithId[int](errtype.Validation, "initial fail")
	flatMappedSkip := FlatMap(fail, func(v int) Result[string] {
		return Success("should not run")
	})

	if flatMappedSkip.IsSuccess() || flatMappedSkip.Fault().GetMessage() != "initial fail" {
		t.Errorf("FlatMap failed to skip on initial error")
	}
}

func TestTap(t *testing.T) {
	succ := Success("hello")
	tapped := false
	res := Tap(succ, func(v string) {
		tapped = true
		if v != "hello" {
			t.Errorf("Tap received wrong value")
		}
	})

	if !tapped || res.Data() != "hello" {
		t.Errorf("Tap failed to execute side effect or return original result")
	}

	fail := FailureWithId[string](errtype.Validation, "fail")
	tappedFail := false
	resFail := Tap(fail, func(v string) {
		tappedFail = true
	})

	if tappedFail || resFail.IsSuccess() {
		t.Errorf("Tap executed side effect on failure")
	}
}
