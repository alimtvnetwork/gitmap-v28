package result

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// ResultSlice encapsulates a slice computation outcome with typed item list or *apperror.AppError.
type ResultSlice[T any] struct {
	Value []T
	Data  []T
	Err   *apperror.AppError
}

// IsSuccess reports whether the slice operation succeeded without error.
func (r ResultSlice[T]) IsSuccess() bool {
	return r.Err == nil
}

// IsFailed reports whether the slice operation encountered an error.
func (r ResultSlice[T]) IsFailed() bool {
	return r.Err != nil
}

// IsFailure reports whether the slice operation encountered an error.
func (r ResultSlice[T]) IsFailure() bool {
	return r.Err != nil
}

// HasError reports whether an active error is attached to the result.
func (r ResultSlice[T]) HasError() bool {
	return r.Err != nil
}

// IsEmptyError reports whether no active error exists.
func (r ResultSlice[T]) IsEmptyError() bool {
	return r.Err == nil
}

// HasNoError reports whether no active error exists.
func (r ResultSlice[T]) HasNoError() bool {
	return r.Err == nil
}

// IsEmpty reports whether the underlying slice has 0 items or is uninitialized.
func (r ResultSlice[T]) IsEmpty() bool {
	return len(r.Data) == 0
}

// AppError returns the underlying AppError or nil.
func (r ResultSlice[T]) AppError() *apperror.AppError {
	return r.Err
}

// Fault returns the underlying AppError or nil.
func (r ResultSlice[T]) Fault() *apperror.AppError {
	return r.Err
}

// Count returns the number of items in the slice, or 0 if uninitialized or failed.
func (r ResultSlice[T]) Count() int {
	return len(r.Data)
}

// Get safely retrieves an item by zero-based index.
func (r ResultSlice[T]) Get(idx int) (T, bool) {
	if idx < 0 || idx >= len(r.Data) {
		var zero T

		return zero, false
	}

	return r.Data[idx], true
}

// Unwrap returns the underlying slice and AppError tuple.
func (r ResultSlice[T]) Unwrap() ([]T, *apperror.AppError) {
	return r.Data, r.Err
}

// UnwrapOr returns the underlying slice if successful, or defaultVal if failed.
func (r ResultSlice[T]) UnwrapOr(defaultVal []T) []T {
	if r.IsSuccess() {
		return r.Data
	}

	return defaultVal
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
