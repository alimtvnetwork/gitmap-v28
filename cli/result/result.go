package result

import (
	"reflect"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// Result encapsulates a computation outcome with typed value or *apperror.AppError.
type Result[T any] struct {
	Value T
	Data  T
	Err   *apperror.AppError
}

// IsSuccess reports whether the result represents a successful operation.
func (r Result[T]) IsSuccess() bool {
	return r.Err == nil
}

// IsFailed reports whether the result represents a failed operation.
func (r Result[T]) IsFailed() bool {
	return r.Err != nil
}

// IsFailure reports whether the result represents a failed operation (alias).
func (r Result[T]) IsFailure() bool {
	return r.Err != nil
}

// IsInvalid reports whether the result is invalid or failed.
func (r Result[T]) IsInvalid() bool {
	return r.Err != nil
}

// HasError reports whether an error is present.
func (r Result[T]) HasError() bool {
	return r.Err != nil
}

// IsEmptyError reports whether no error is present.
func (r Result[T]) IsEmptyError() bool {
	return r.Err == nil
}

// HasNoError reports whether no error is present.
func (r Result[T]) HasNoError() bool {
	return r.Err == nil
}

// HasValidError reports whether an AppError exists and is properly structured.
func (r Result[T]) HasValidError() bool {
	if r.Err != nil {
		return r.Err.IsValid()
	}

	return false
}

// IsEmpty reports whether the result represents an empty or zero value.
func (r Result[T]) IsEmpty() bool {
	var zero T

	return reflect.DeepEqual(r.Value, zero)
}

// AppError returns the underlying AppError or nil.
func (r Result[T]) AppError() *apperror.AppError {
	return r.Err
}

// Fault returns the underlying AppError or nil (alias).
func (r Result[T]) Fault() *apperror.AppError {
	return r.Err
}

// Unwrap returns the value and AppError tuple.
func (r Result[T]) Unwrap() (T, *apperror.AppError) {
	return r.Value, r.Err
}

// UnwrapOr returns the value if success, or defaultVal if failed.
func (r Result[T]) UnwrapOr(defaultVal T) T {
	if r.IsSuccess() {
		return r.Value
	}

	return defaultVal
}

// Ok constructs a successful Result envelope with Value and Data.
func Ok[T any](val T) Result[T] {
	return Result[T]{
		Value: val,
		Data:  val,
	}
}

// Fail constructs a failed Result envelope with *apperror.AppError.
func Fail[T any](err *apperror.AppError) Result[T] {
	return Result[T]{
		Err: err,
	}
}

// SuccessResult constructs a successful Result envelope with Value and Data.
func SuccessResult[T any](val T) Result[T] {
	return Ok(val)
}

// FailureResult constructs a failed Result envelope with *apperror.AppError.
func FailureResult[T any](err *apperror.AppError) Result[T] {
	return Fail[T](err)
}

// NewSuccess constructs a successful Result envelope with Data.
func NewSuccess[T any](data T) Result[T] {
	return Ok(data)
}

// NewFailure constructs a failed Result envelope from any error.
func NewFailure[T any](err error) Result[T] {
	if appErr, isAppErr := err.(*apperror.AppError); isAppErr {
		return Fail[T](appErr)
	}

	if err == nil {
		return Result[T]{}
	}

	appErr := apperror.WrapSimple(err, "result.NewFailure")

	return Fail[T](appErr)
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

	return Fail[T](appErr)
}

// HandleError processes the underlying error if one exists, without panicking or exiting.
func (r Result[T]) HandleError() {
	if r.Err != nil {
		r.Err.HandleError()
	}
}
