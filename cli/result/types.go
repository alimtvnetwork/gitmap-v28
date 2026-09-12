// Package result — types.go defines core generic Result container types and canonical aliases.
package result

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type (
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

	// Wrap is a canonical single reusable alias for Result[T].
	Wrap[T any] = Result[T]

	// Slice is a canonical single reusable alias for ResultSlice[T].
	Slice[T any] = ResultSlice[T]

	// Map is a canonical single reusable alias for ResultMap[K, V].
	Map[K comparable, V any] = ResultMap[K, V]
)
