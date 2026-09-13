package result

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// IsSuccess reports whether the slice operation succeeded without error.
func (r *ResultSlice[T]) IsSuccess() bool {
	if r == nil {
		return false
	}

	return r.Err == nil
}

// IsSafe reports whether the slice operation succeeded without error (alias).
func (r *ResultSlice[T]) IsSafe() bool {
	if r == nil {
		return false
	}

	return r.Err == nil
}

// IsFailed reports whether the slice operation encountered an error.
func (r *ResultSlice[T]) IsFailed() bool {
	if r == nil {
		return true
	}

	return r.Err != nil
}

// IsFailure reports whether the slice operation encountered an error.
func (r *ResultSlice[T]) IsFailure() bool {
	if r == nil {
		return true
	}

	return r.Err != nil
}

// HasError reports whether an active error is attached to the result.
func (r *ResultSlice[T]) HasError() bool {
	if r == nil {
		return true
	}

	return r.Err != nil
}

// IsEmptyError reports whether no active error exists.
func (r *ResultSlice[T]) IsEmptyError() bool {
	if r == nil {
		return false
	}

	return r.Err == nil
}

// HasNoError reports whether no active error exists.
func (r *ResultSlice[T]) HasNoError() bool {
	if r == nil {
		return false
	}

	return r.Err == nil
}

// IsEmpty reports whether the underlying slice has 0 items or is uninitialized.
func (r *ResultSlice[T]) IsEmpty() bool {
	if r == nil || r.Err != nil {
		return true
	}

	return len(r.Data) == 0
}

// Count returns the number of items in the slice, or 0 if uninitialized or failed.
func (r *ResultSlice[T]) Count() int {
	if r == nil || r.Err != nil {
		return 0
	}

	return len(r.Data)
}

// IsCountOtherThan reports whether the operation failed OR the slice length != number.
func (r *ResultSlice[T]) IsCountOtherThan(number int) bool {
	if r == nil || r.IsFailure() {
		return true
	}

	return r.Count() != number
}

// HasRecord reports whether the operation succeeded AND contains more than 0 records.
func (r *ResultSlice[T]) HasRecord() bool {
	if r == nil || r.Err != nil {
		return false
	}

	return len(r.Data) > 0
}

// HasRecords is an alias for HasRecord.
func (r *ResultSlice[T]) HasRecords() bool {
	return r.HasRecord()
}

// IsDefined reports whether the operation succeeded AND has records.
func (r *ResultSlice[T]) IsDefined() bool {
	if r == nil || r.Err != nil {
		return false
	}

	return len(r.Data) > 0
}

// Items returns the underlying slice safely (nil on nil receiver).
func (r *ResultSlice[T]) Items() []T {
	if r == nil {
		return nil
	}

	return r.Data
}

// AppError returns the underlying AppError or nil.
func (r *ResultSlice[T]) AppError() *apperror.AppError {
	if r == nil {
		return nil
	}

	return r.Err
}

// Fault returns the underlying AppError or nil.
func (r *ResultSlice[T]) Fault() *apperror.AppError {
	if r == nil {
		return nil
	}

	return r.Err
}

// Get safely retrieves an item by zero-based index.
func (r *ResultSlice[T]) Get(idx int) (T, bool) {
	if r == nil || idx < 0 || idx >= len(r.Data) {
		var zero T

		return zero, false
	}

	return r.Data[idx], true
}

// Unwrap returns the underlying slice and AppError tuple.
func (r *ResultSlice[T]) Unwrap() ([]T, *apperror.AppError) {
	if r == nil {
		return nil, nil
	}

	return r.Data, r.Err
}

// UnwrapOr returns the underlying slice if successful, or defaultVal if failed.
func (r *ResultSlice[T]) UnwrapOr(defaultVal []T) []T {
	if r == nil || r.IsFailure() {
		return defaultVal
	}

	return r.Data
}

// OkSlice constructs a successful ResultSlice envelope.
func OkSlice[T any](data []T) ResultSlice[T] {
	if data == nil {
		data = make([]T, 0)
	}

	return ResultSlice[T]{
		Value: data,
		Data:  data,
	}
}

// FailSlice constructs a failed ResultSlice envelope with *apperror.AppError.
func FailSlice[T any](err *apperror.AppError) ResultSlice[T] {
	return ResultSlice[T]{
		Err: err,
	}
}

// NewSuccessSlice constructs a successful ResultSlice envelope.
func NewSuccessSlice[T any](data []T) ResultSlice[T] {
	return OkSlice(data)
}

// NewFailureSlice constructs a failed ResultSlice envelope from any error.
func NewFailureSlice[T any](err error) ResultSlice[T] {
	if appErr, isAppErr := err.(*apperror.AppError); isAppErr {
		return FailSlice[T](appErr)
	}

	if err == nil {
		return ResultSlice[T]{}
	}

	appErr := apperror.WrapSimple(err, "result.NewFailureSlice")

	return FailSlice[T](appErr)
}
