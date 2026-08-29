package result

import "github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"

// Result encapsulates a computation outcome with typed value or *apperror.AppError.
type Result[T any] struct {
	Value    T
	Err      *apperror.AppError
	Data     T
	AppError error
}

// IsSuccess reports whether the result represents a successful operation.
func (r Result[T]) IsSuccess() bool {
	return r.Err == nil && r.AppError == nil
}

// IsFailed reports whether the result represents a failed operation.
func (r Result[T]) IsFailed() bool {
	return r.Err != nil || r.AppError != nil
}

// Unwrap returns the value and error tuple.
func (r Result[T]) Unwrap() (T, *apperror.AppError) {
	if r.Err != nil {
		return r.Value, r.Err
	}

	if r.AppError == nil {
		return r.Value, nil
	}

	appErr, isAppErr := r.AppError.(*apperror.AppError)
	if isAppErr {
		return r.Value, appErr
	}

	return r.Value, apperror.WrapSimple(r.AppError, "result.Unwrap")
}

// SuccessResult constructs a successful Result envelope with Value.
func SuccessResult[T any](val T) Result[T] {
	return Result[T]{
		Value: val,
		Data:  val,
	}
}

// FailureResult constructs a failed Result envelope with *apperror.AppError.
func FailureResult[T any](err *apperror.AppError) Result[T] {
	return Result[T]{
		Err:      err,
		AppError: err,
	}
}

// NewSuccess constructs a successful Result envelope with Data.
func NewSuccess[T any](data T) Result[T] {
	return SuccessResult(data)
}

// NewFailure constructs a failed Result envelope from any error.
func NewFailure[T any](err error) Result[T] {
	if appErr, isAppErr := err.(*apperror.AppError); isAppErr {
		return FailureResult[T](appErr)
	}

	if err == nil {
		return Result[T]{}
	}

	appErr := apperror.WrapSimple(err, "result.NewFailure")

	return FailureResult[T](appErr)
}
