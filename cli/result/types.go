// Package result — types.go defines core generic Result container types, ErrorWrapper, and canonical aliases.
package result

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type (
	// ErrorWrapper encapsulates an operation or routing outcome with optional *apperror.AppError and match state.
	ErrorWrapper struct {
		Err       *apperror.AppError
		isMatched bool
	}

	// Result encapsulates a computation outcome with typed value or *apperror.AppError.
	Result[T any] struct {
		Value     T
		Data      T
		Err       *apperror.AppError
		isDefined bool
	}

	// ResultSlice encapsulates a slice computation outcome with typed item list or *apperror.AppError.
	ResultSlice[T any] struct {
		Value []T
		Data  []T
		Err   *apperror.AppError
	}

	// ResultMap encapsulates a map computation outcome with typed key-value data or *apperror.AppError.
	ResultMap[K comparable, V any] struct {
		Value map[K]V
		Data  map[K]V
		Err   *apperror.AppError
	}

	// ErrorWrap is a canonical single reusable alias for ErrorWrapper.
	ErrorWrap = ErrorWrapper

	// Proper type aliases for common result envelopes.
	BoolResult        = Result[bool]
	StringResult      = Result[string]
	IntResult         = Result[int]
	Int64Result       = Result[int64]
	Uint64Result      = Result[uint64]
	ByteSliceResult   = Result[[]byte]
	StringSliceResult = ResultSlice[string]
	AnyResult         = Result[any]

	// Wrap is a canonical single reusable alias for Result[T].
	Wrap[T any] = Result[T]

	// Slice is a canonical single reusable alias for ResultSlice[T].
	Slice[T any] = ResultSlice[T]

	// Map is a canonical single reusable alias for ResultMap[K, V].
	Map[K comparable, V any] = ResultMap[K, V]
)
