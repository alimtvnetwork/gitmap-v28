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

// IsFailure reports whether the result represents a failed operation (alias).
func (r Result[T]) IsFailure() bool {
	return r.IsFailed()
}

// IsInvalid reports whether the result is invalid or failed.
func (r Result[T]) IsInvalid() bool {
	return r.IsFailed()
}

// HasError reports whether an error is present.
func (r Result[T]) HasError() bool {
	return r.Err != nil || r.AppError != nil
}

// HasNoError reports whether no error exists.
func (r Result[T]) HasNoError() bool {
	return r.Err == nil && r.AppError == nil
}

// HasValidError reports whether an AppError exists and is properly structured.
func (r Result[T]) HasValidError() bool {
	if r.Err != nil {
		return r.Err.IsValid()
	}

	return r.AppError != nil
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

// UnwrapOr returns the value if success, or defaultVal if failed.
func (r Result[T]) UnwrapOr(defaultVal T) T {
	if r.IsSuccess() {
		return r.Value
	}

	return defaultVal
}

// SuccessResult constructs a successful Result envelope with Value and Data.
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

// NewFailureWithType constructs a typed failed Result with code, message, and caller.
func NewFailureWithType[T any](
	errCode string,
	msg string,
	caller string,
) Result[T] {
	appErr := apperror.NewWithDetails(
		caller,
		errCode,
		msg,
		"system",
		apperror.ErrorTypeExecution,
		apperror.SeverityError,
		nil,
	)
	return FailureResult[T](appErr)
}
