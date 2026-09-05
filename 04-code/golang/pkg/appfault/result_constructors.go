package appfault

import "coding-guidelines/common/pkg/errtype"

// SuccessResult creates a successful Result wrapping a value.
func SuccessResult[T any](val T) Result[T] {
	return Result[T]{Value: val, AppError: nil}
}

// NewSuccess creates a successful Result wrapping a value.
func NewSuccess[T any](data T) Result[T] {
	return SuccessResult(data)
}

// FailureResult creates a failed Result with a structured AppError.
func FailureResult[T any](err *AppError) Result[T] {
	return Result[T]{AppError: err}
}

// Ok is an alias for SuccessResult.
func Ok[T any](val T) Result[T] {
	return SuccessResult(val)
}

// NewFailure creates a failed Result from an explicit type and cause error.
func NewFailure[T any](errType errtype.Variation, cause error) Result[T] {
	if cause == nil || errType == errtype.None {
		return Result[T]{}
	}

	return FailureResult[T](WrapType(errType, cause))
}

// NewFailureWithType creates a failed Result with explicit type and caller.
func NewFailureWithType[T any](errType errtype.Variation, msg string, caller string) Result[T] {
	e := New(errType, msg)
	if len(caller) > 0 {
		e.caller.Function = caller
	}

	return FailureResult[T](e)
}

// Fail creates a failed Result from an AppError.
func Fail[T any](err *AppError) Result[T] {
	return FailureResult[T](err)
}

// FailWrap wraps a raw error into a failed Result.
func FailWrap[T any](errType errtype.Variation, cause error, msg string) Result[T] {
	return FailureResult[T](Wrap(errType, cause, msg))
}

// FailNew creates a new AppError and returns a failed Result.
func FailNew[T any](errType errtype.Variation, msg string) Result[T] {
	return FailureResult[T](New(errType, msg))
}

// NewFailureWithId creates a failed Result with an error ID (errtype.Variation) and message.
func NewFailureWithId[T any](errType errtype.Variation, msg string) Result[T] {
	return FailureResult[T](New(errType, msg))
}

// NewFailureFromError creates a failed Result from an AppError object.
func NewFailureFromError[T any](err *AppError) Result[T] {
	return FailureResult[T](err)
}

// FailureWithId creates a failed Result with an error ID (errtype.Variation) and message.
func FailureWithId[T any](errType errtype.Variation, msg string) Result[T] {
	return NewFailureWithId[T](errType, msg)
}

// NewFailureWithCause creates a failed Result with an error ID, cause error, and message.
func NewFailureWithCause[T any](errType errtype.Variation, cause error, msg string) Result[T] {
	return FailureResult[T](Wrap(errType, cause, msg))
}

// FailureFromWrap creates a failed Result propagating the AppError from another Result.
func FailureFromWrap[T any, U any](failed Result[U]) Result[T] {
	return FailureResult[T](failed.AppError)
}
