package result_test

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

func TestResultPointerNullSafety(t *testing.T) {
	var nilRes *result.Result[string] = nil

	if nilRes.IsSuccess() {
		t.Errorf("nil Result.IsSuccess() must be false")
	}

	if !nilRes.IsFailure() {
		t.Errorf("nil Result.IsFailure() must be true")
	}

	if !nilRes.IsFailed() {
		t.Errorf("nil Result.IsFailed() must be true")
	}

	if !nilRes.HasError() {
		t.Errorf("nil Result.HasError() must be true")
	}

	if nilRes.IsEmptyError() {
		t.Errorf("nil Result.IsEmptyError() must be false")
	}

	if nilRes.HasNoError() {
		t.Errorf("nil Result.HasNoError() must be false")
	}

	if nilRes.HasValidError() {
		t.Errorf("nil Result.HasValidError() must be false")
	}

	if !nilRes.IsEmpty() {
		t.Errorf("nil Result.IsEmpty() must be true")
	}

	if nilRes.Count() != 0 {
		t.Errorf("nil Result.Count() must be 0, got %d", nilRes.Count())
	}

	if !nilRes.IsCountOtherThan(1) {
		t.Errorf("nil Result.IsCountOtherThan(1) must be true")
	}

	if !nilRes.IsCountOtherThan(0) {
		t.Errorf("nil Result.IsCountOtherThan(0) must be true (failed state)")
	}

	if nilRes.HasRecord() {
		t.Errorf("nil Result.HasRecord() must be false")
	}

	if nilRes.HasRecords() {
		t.Errorf("nil Result.HasRecords() must be false")
	}

	if nilRes.IsDefined() {
		t.Errorf("nil Result.IsDefined() must be false")
	}

	if nilRes.AppError() != nil {
		t.Errorf("nil Result.AppError() must be nil")
	}

	if nilRes.Fault() != nil {
		t.Errorf("nil Result.Fault() must be nil")
	}

	val, err := nilRes.Unwrap()
	if val != "" || err != nil {
		t.Errorf("nil Result.Unwrap() expected empty and nil err")
	}

	if nilRes.UnwrapOr("fallback") != "fallback" {
		t.Errorf("nil Result.UnwrapOr() expected fallback")
	}

	// Ensure HandleError does not panic
	nilRes.HandleError()
}

func TestResultSlicePointerNullSafety(t *testing.T) {
	var nilSlice *result.ResultSlice[int] = nil

	if nilSlice.IsSuccess() {
		t.Errorf("nil ResultSlice.IsSuccess() must be false")
	}

	if !nilSlice.IsFailure() {
		t.Errorf("nil ResultSlice.IsFailure() must be true")
	}

	if !nilSlice.IsFailed() {
		t.Errorf("nil ResultSlice.IsFailed() must be true")
	}

	if !nilSlice.HasError() {
		t.Errorf("nil ResultSlice.HasError() must be true")
	}

	if nilSlice.IsEmptyError() {
		t.Errorf("nil ResultSlice.IsEmptyError() must be false")
	}

	if nilSlice.HasNoError() {
		t.Errorf("nil ResultSlice.HasNoError() must be false")
	}

	if !nilSlice.IsEmpty() {
		t.Errorf("nil ResultSlice.IsEmpty() must be true")
	}

	if nilSlice.Count() != 0 {
		t.Errorf("nil ResultSlice.Count() must be 0, got %d", nilSlice.Count())
	}

	if !nilSlice.IsCountOtherThan(1) {
		t.Errorf("nil ResultSlice.IsCountOtherThan(1) must be true")
	}

	if !nilSlice.IsCountOtherThan(0) {
		t.Errorf("nil ResultSlice.IsCountOtherThan(0) must be true (failed state)")
	}

	if nilSlice.HasRecord() {
		t.Errorf("nil ResultSlice.HasRecord() must be false")
	}

	if nilSlice.HasRecords() {
		t.Errorf("nil ResultSlice.HasRecords() must be false")
	}

	if nilSlice.IsDefined() {
		t.Errorf("nil ResultSlice.IsDefined() must be false")
	}

	if nilSlice.Items() != nil {
		t.Errorf("nil ResultSlice.Items() must be nil")
	}

	if nilSlice.AppError() != nil {
		t.Errorf("nil ResultSlice.AppError() must be nil")
	}

	if nilSlice.Fault() != nil {
		t.Errorf("nil ResultSlice.Fault() must be nil")
	}

	val, isFound := nilSlice.Get(0)
	if val != 0 || isFound {
		t.Errorf("nil ResultSlice.Get() expected 0, false")
	}

	items, err := nilSlice.Unwrap()
	if items != nil || err != nil {
		t.Errorf("nil ResultSlice.Unwrap() expected nil items and nil err")
	}

	def := []int{1, 2}
	if len(nilSlice.UnwrapOr(def)) != 2 {
		t.Errorf("nil ResultSlice.UnwrapOr() expected default")
	}
}

func TestResultMapPointerNullSafety(t *testing.T) {
	var nilMap *result.ResultMap[string, int] = nil

	if nilMap.IsSuccess() {
		t.Errorf("nil ResultMap.IsSuccess() must be false")
	}

	if !nilMap.IsFailure() {
		t.Errorf("nil ResultMap.IsFailure() must be true")
	}

	if !nilMap.IsFailed() {
		t.Errorf("nil ResultMap.IsFailed() must be true")
	}

	if !nilMap.HasError() {
		t.Errorf("nil ResultMap.HasError() must be true")
	}

	if nilMap.IsEmptyError() {
		t.Errorf("nil ResultMap.IsEmptyError() must be false")
	}

	if nilMap.HasNoError() {
		t.Errorf("nil ResultMap.HasNoError() must be false")
	}

	if !nilMap.IsEmpty() {
		t.Errorf("nil ResultMap.IsEmpty() must be true")
	}

	if nilMap.Count() != 0 {
		t.Errorf("nil ResultMap.Count() must be 0, got %d", nilMap.Count())
	}

	if !nilMap.IsCountOtherThan(1) {
		t.Errorf("nil ResultMap.IsCountOtherThan(1) must be true")
	}

	if !nilMap.IsCountOtherThan(0) {
		t.Errorf("nil ResultMap.IsCountOtherThan(0) must be true (failed state)")
	}

	if nilMap.HasRecord() {
		t.Errorf("nil ResultMap.HasRecord() must be false")
	}

	if nilMap.HasRecords() {
		t.Errorf("nil ResultMap.HasRecords() must be false")
	}

	if nilMap.IsDefined() {
		t.Errorf("nil ResultMap.IsDefined() must be false")
	}

	if nilMap.AppError() != nil {
		t.Errorf("nil ResultMap.AppError() must be nil")
	}

	if nilMap.Fault() != nil {
		t.Errorf("nil ResultMap.Fault() must be nil")
	}

	val, isFound := nilMap.Get("key")
	if val != 0 || isFound {
		t.Errorf("nil ResultMap.Get() expected 0, false")
	}

	if nilMap.Has("key") {
		t.Errorf("nil ResultMap.Has() must be false")
	}

	if len(nilMap.Keys()) != 0 {
		t.Errorf("nil ResultMap.Keys() expected empty slice")
	}

	if len(nilMap.Values()) != 0 {
		t.Errorf("nil ResultMap.Values() expected empty slice")
	}

	data, err := nilMap.Unwrap()
	if data != nil || err != nil {
		t.Errorf("nil ResultMap.Unwrap() expected nil data and nil err")
	}

	def := map[string]int{"k": 1}
	if len(nilMap.UnwrapOr(def)) != 1 {
		t.Errorf("nil ResultMap.UnwrapOr() expected default")
	}
}

func TestCorePredicatesOnSuccessAndFailure(t *testing.T) {
	// Success Result
	resOk := result.Ok("hello")
	if resOk.IsFailure() {
		t.Errorf("resOk expected success")
	}
	if resOk.Count() != 1 || resOk.IsCountOtherThan(1) {
		t.Errorf("resOk expected count 1")
	}
	if !resOk.HasRecord() || !resOk.IsDefined() {
		t.Errorf("resOk expected HasRecord and IsDefined")
	}

	// Failed Result
	appErr := apperror.NewWithDetails("test", "ERR_CODE", "test error", "system", apperror.ErrorTypeExecution, apperror.SeverityError, nil)
	resFail := result.Fail[string](appErr)
	if resFail.IsSuccess() || !resFail.IsFailure() {
		t.Errorf("resFail expected failure")
	}
	if resFail.Count() != 0 || !resFail.IsCountOtherThan(1) {
		t.Errorf("resFail expected count 0 and IsCountOtherThan(1)")
	}
	if resFail.HasRecord() || resFail.IsDefined() {
		t.Errorf("resFail expected HasRecord=false and IsDefined=false")
	}

	// Success Slice
	sliceOk := result.OkSlice([]string{"a", "b"})
	if sliceOk.IsCountOtherThan(2) {
		t.Errorf("sliceOk expected count 2")
	}
	if !sliceOk.HasRecord() || !sliceOk.IsDefined() {
		t.Errorf("sliceOk expected HasRecord and IsDefined")
	}
	if !sliceOk.IsCountOtherThan(1) {
		t.Errorf("sliceOk IsCountOtherThan(1) expected true")
	}

	// Empty Slice
	sliceEmpty := result.OkSlice([]string{})
	if sliceEmpty.IsCountOtherThan(0) || !sliceEmpty.IsEmpty() {
		t.Errorf("sliceEmpty expected count 0, IsEmpty=true")
	}
	if sliceEmpty.HasRecord() || sliceEmpty.IsDefined() {
		t.Errorf("sliceEmpty expected HasRecord=false, IsDefined=false")
	}

	// Map Success
	mapOk := result.OkMap(map[string]int{"k": 1})
	if mapOk.IsCountOtherThan(1) {
		t.Errorf("mapOk expected count 1")
	}
	if !mapOk.HasRecord() || !mapOk.IsDefined() {
		t.Errorf("mapOk expected HasRecord and IsDefined")
	}
}
