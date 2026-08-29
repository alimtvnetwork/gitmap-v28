package result_test

import (
	"errors"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/result"
)

func TestSuccessResult(t *testing.T) {
	res := result.SuccessResult("hello")
	if !res.IsSuccess() {
		t.Fatal("expected IsSuccess to be true")
	}

	if res.IsFailed() {
		t.Fatal("expected IsFailed to be false")
	}

	if !res.HasNoError() {
		t.Fatal("expected HasNoError to be true")
	}

	if res.HasError() {
		t.Fatal("expected HasError to be false")
	}

	val, err := res.Unwrap()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if val != "hello" {
		t.Fatalf("expected 'hello', got %s", val)
	}

	if res.UnwrapOr("default") != "hello" {
		t.Fatal("expected 'hello' from UnwrapOr")
	}
}

func TestFailureResult(t *testing.T) {
	appErr := apperror.WrapSimple(errors.New("underlying"), "TestFailureResult")
	res := result.FailureResult[string](appErr)

	if res.IsSuccess() {
		t.Fatal("expected IsSuccess to be false")
	}

	if !res.IsFailed() {
		t.Fatal("expected IsFailed to be true")
	}

	if !res.IsFailure() {
		t.Fatal("expected IsFailure to be true")
	}

	if !res.IsInvalid() {
		t.Fatal("expected IsInvalid to be true")
	}

	if !res.HasError() {
		t.Fatal("expected HasError to be true")
	}

	if res.HasNoError() {
		t.Fatal("expected HasNoError to be false")
	}

	_, err := res.Unwrap()
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if res.UnwrapOr("fallback") != "fallback" {
		t.Fatal("expected fallback from UnwrapOr")
	}
}

func TestNewFailureGenericError(t *testing.T) {
	stdErr := errors.New("raw error")
	res := result.NewFailure[int](stdErr)

	if !res.IsFailed() {
		t.Fatal("expected IsFailed to be true")
	}

	_, err := res.Unwrap()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestNewFailureWithType(t *testing.T) {
	res := result.NewFailureWithType[int]("E_CODE_TEST", "test error", "TestNewFailureWithType")
	if !res.IsFailed() {
		t.Fatal("expected failure")
	}

	if !res.HasValidError() {
		t.Fatal("expected HasValidError to be true")
	}
}
