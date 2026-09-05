package dbengine

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/result"
)

// Exact typed result envelopes wrapping result.Result[T].
type Uint64Result = result.Result[uint64]
type Int64Result = result.Result[int64]
type StringResult = result.Result[string]
type BoolResult = result.Result[bool]
type RowsAffectedResult = result.Result[int64]
type EntityResult[T any] = result.Result[*T]
type ListResult[T any] = result.Result[[]T]

// SuccessUint64 wraps a uint64 value in a successful result envelope.
func SuccessUint64(val uint64) Uint64Result {
	return result.SuccessResult(val)
}

// FailureUint64 wraps an AppError in a failed Uint64Result envelope.
func FailureUint64(err *apperror.AppError) Uint64Result {
	return result.FailureResult[uint64](err)
}

// SuccessInt64 wraps an int64 value in a successful result envelope.
func SuccessInt64(val int64) Int64Result {
	return result.SuccessResult(val)
}

// FailureInt64 wraps an AppError in a failed Int64Result envelope.
func FailureInt64(err *apperror.AppError) Int64Result {
	return result.FailureResult[int64](err)
}

// SuccessString wraps a string value in a successful result envelope.
func SuccessString(val string) StringResult {
	return result.SuccessResult(val)
}

// FailureString wraps an AppError in a failed StringResult envelope.
func FailureString(err *apperror.AppError) StringResult {
	return result.FailureResult[string](err)
}

// SuccessBool wraps a boolean value in a successful result envelope.
func SuccessBool(val bool) BoolResult {
	return result.SuccessResult(val)
}

// FailureBool wraps an AppError in a failed BoolResult envelope.
func FailureBool(err *apperror.AppError) BoolResult {
	return result.FailureResult[bool](err)
}

// SuccessRowsAffected wraps the number of rows affected in a successful result envelope.
func SuccessRowsAffected(val int64) RowsAffectedResult {
	return result.SuccessResult(val)
}

// FailureRowsAffected wraps an AppError in a failed RowsAffectedResult envelope.
func FailureRowsAffected(err *apperror.AppError) RowsAffectedResult {
	return result.FailureResult[int64](err)
}

// SuccessEntity wraps an entity pointer in a successful result envelope.
func SuccessEntity[T any](val *T) EntityResult[T] {
	return result.SuccessResult(val)
}

// FailureEntity wraps an AppError in a failed EntityResult envelope.
func FailureEntity[T any](err *apperror.AppError) EntityResult[T] {
	return result.FailureResult[*T](err)
}

// SuccessList wraps an entity slice in a successful result envelope.
func SuccessList[T any](val []T) ListResult[T] {
	return result.SuccessResult(val)
}

// FailureList wraps an AppError in a failed ListResult envelope.
func FailureList[T any](err *apperror.AppError) ListResult[T] {
	return result.FailureResult[[]T](err)
}
