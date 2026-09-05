package result

import (
	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
)

type Wrap[T any] = appfault.Result[T]
type Result[T any] = Wrap[T]

func WrapSuccess[T any](data T) Wrap[T] {
	return appfault.NewSuccess(data)
}

func WrapFailure[T any](err *appfault.AppError) Wrap[T] {
	return appfault.FailureResult[T](err)
}

func Success[T any](data T) Wrap[T] {
	return appfault.NewSuccess(data)
}

func Failure[T any](err *appfault.AppError) Wrap[T] {
	return appfault.FailureResult[T](err)
}

func WrapFailureFromError[T any](err *appfault.AppError) Wrap[T] {
	return appfault.FailureResult[T](err)
}

func WrapFailureWithId[T any](errType errtype.Variation, msg string) Wrap[T] {
	return appfault.NewFailureWithId[T](errType, msg)
}

func WrapFailureWithCause[T any](errType errtype.Variation, cause error, msg string) Wrap[T] {
	return appfault.NewFailureWithCause[T](errType, cause, msg)
}

func WrapFailureFromWrap[T any, U any](failed Wrap[U]) Wrap[T] {
	return appfault.FailureFromWrap[T](failed)
}

func FailureFromWrap[T any, U any](failed Wrap[U]) Wrap[T] {
	return appfault.FailureFromWrap[T](failed)
}

func FailureWithId[T any](errType errtype.Variation, msg string) Wrap[T] {
	return appfault.NewFailureWithId[T](errType, msg)
}

func FailureWithCause[T any](errType errtype.Variation, cause error, msg string) Wrap[T] {
	return appfault.NewFailureWithCause[T](errType, cause, msg)
}

func SuccessResult[T any](val T) Result[T] {
	return appfault.SuccessResult(val)
}

func NewSuccess[T any](data T) Result[T] {
	return appfault.NewSuccess(data)
}

func FailureResult[T any](err *appfault.AppError) Result[T] {
	return appfault.FailureResult[T](err)
}

func NewFailure[T any](errType errtype.Variation, cause error) Result[T] {
	return appfault.NewFailure[T](errType, cause)
}

func NewFailureWithType[T any](errType errtype.Variation, msg string, caller string) Result[T] {
	return appfault.NewFailureWithType[T](errType, msg, caller)
}
