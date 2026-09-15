package result_test

import (
	"errors"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

func TestErrorWrapper_Success(t *testing.T) {
	ew := result.SuccessWrapper()
	if ew.IsFailed() {
		t.Fatal("expected IsSuccess to be true")
	}

	if !ew.IsSafe() {
		t.Fatal("expected IsSafe to be true")
	}

	if ew.IsFailed() {
		t.Fatal("expected IsFailed to be false")
	}

	if ew.IsFailure() {
		t.Fatal("expected IsFailure to be false")
	}

	if ew.IsInvalid() {
		t.Fatal("expected IsInvalid to be false")
	}

	if !ew.IsMatched() {
		t.Fatal("expected IsMatched to be true")
	}

	if !ew.IsHandled() {
		t.Fatal("expected IsHandled to be true")
	}

	if ew.AppError() != nil {
		t.Fatal("expected nil AppError")
	}

	if ew.Error() != "" {
		t.Fatal("expected empty Error string")
	}
}

func TestErrorWrapper_Failure(t *testing.T) {
	appErr := apperror.NewSimple("test fail", "E100")
	ew := result.FailureWrapper(appErr)

	if ew.IsSuccess() {
		t.Fatal("expected IsSuccess to be false")
	}

	if !ew.IsFailed() {
		t.Fatal("expected IsFailed to be true")
	}

	if !ew.IsInvalid() {
		t.Fatal("expected IsInvalid to be true")
	}

	if !ew.IsMatched() {
		t.Fatal("expected IsMatched to be true")
	}

	if ew.AppError() != appErr {
		t.Fatal("expected AppError to match")
	}

	if ew.Error() == "" {
		t.Fatal("expected non-empty Error string")
	}
}

func TestErrorWrapper_FailureWrapperErr(t *testing.T) {
	ewNil := result.FailureWrapperErr(nil)
	if ewNil.IsFailed() {
		t.Fatal("expected nil error to produce SuccessWrapper")
	}

	stdErr := errors.New("raw standard error")
	ewStd := result.FailureWrapperErr(stdErr)
	if !ewStd.IsFailed() {
		t.Fatal("expected FailureWrapperErr with error to fail")
	}

	appErr := apperror.NewSimple("app error direct", "E200")
	ewApp := result.FailureWrapperErr(appErr)
	if ewApp.AppError() != appErr {
		t.Fatal("expected direct AppError to be retained")
	}
}

func TestErrorWrapper_Unmatched(t *testing.T) {
	ew := result.UnmatchedWrapper()
	if ew.IsMatched() {
		t.Fatal("expected IsMatched to be false")
	}

	if !ew.IsInvalid() {
		t.Fatal("expected IsInvalid to be true for unmatched")
	}

	if ew.IsSuccess() {
		t.Fatal("expected IsSuccess to be false for unmatched")
	}
}

func TestErrorWrapper_NilReceiver(t *testing.T) {
	var ew *result.ErrorWrapper
	if ew.IsSuccess() {
		t.Fatal("expected nil receiver IsSuccess to be false")
	}

	if !ew.IsFailed() {
		t.Fatal("expected nil receiver IsFailed to be true")
	}

	if !ew.IsInvalid() {
		t.Fatal("expected nil receiver IsInvalid to be true")
	}

	if ew.IsMatched() {
		t.Fatal("expected nil receiver IsMatched to be false")
	}

	if ew.AppError() != nil {
		t.Fatal("expected nil receiver AppError to be nil")
	}

	if ew.Error() != "" {
		t.Fatal("expected nil receiver Error to be empty")
	}
}

func TestErrorWrapper_AliasesAndPredicates(t *testing.T) {
	ok := result.OkWrapper()
	if ok.IsFailed() || !ok.IsEmptyError() || !ok.HasNoError() {
		t.Fatal("expected OkWrapper to be successful with no error")
	}

	appErr := apperror.NewSimple("match test", "E300")
	fail := result.FailWrapper(appErr)
	if !fail.IsFailed() || !fail.HasError() || !fail.HasValidError() || fail.Fault() != appErr {
		t.Fatal("expected FailWrapper to report valid error and fault")
	}

	matchOk := result.MatchWrapper(nil)
	if matchOk.IsFailed() || !matchOk.IsMatched() {
		t.Fatal("expected MatchWrapper(nil) to be success and matched")
	}

	matchFail := result.MatchWrapper(errors.New("raw fail"))
	if !matchFail.IsFailed() || !matchFail.IsMatched() {
		t.Fatal("expected MatchWrapper(err) to be failed and matched")
	}

	matchApp := result.MatchWrapperAppErr(appErr)
	if !matchApp.IsFailed() || matchApp.AppError() != appErr {
		t.Fatal("expected MatchWrapperAppErr to match appErr")
	}
}
