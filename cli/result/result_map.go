package result

import (
	"fmt"
	"sort"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// ResultMap encapsulates a map computation outcome with typed key-value data or *apperror.AppError.
type ResultMap[K comparable, V any] struct {
	Value map[K]V
	Data  map[K]V
	Err   *apperror.AppError
}

// IsSuccess reports whether the map operation succeeded without error.
func (r ResultMap[K, V]) IsSuccess() bool {
	return r.Err == nil
}

// IsFailed reports whether the map operation encountered an error.
func (r ResultMap[K, V]) IsFailed() bool {
	return r.Err != nil
}

// IsFailure reports whether the map operation encountered an error.
func (r ResultMap[K, V]) IsFailure() bool {
	return r.Err != nil
}

// HasError reports whether an active error is attached to the result.
func (r ResultMap[K, V]) HasError() bool {
	return r.Err != nil
}

// IsEmptyError reports whether no active error exists.
func (r ResultMap[K, V]) IsEmptyError() bool {
	return r.Err == nil
}

// HasNoError reports whether no active error exists.
func (r ResultMap[K, V]) HasNoError() bool {
	return r.Err == nil
}

// IsEmpty reports whether the underlying map has 0 items or is uninitialized.
func (r ResultMap[K, V]) IsEmpty() bool {
	return len(r.Data) == 0
}

// AppError returns the underlying AppError or nil.
func (r ResultMap[K, V]) AppError() *apperror.AppError {
	return r.Err
}

// Fault returns the underlying AppError or nil.
func (r ResultMap[K, V]) Fault() *apperror.AppError {
	return r.Err
}

// Count returns the number of entries in the map, or 0 if uninitialized or failed.
func (r ResultMap[K, V]) Count() int {
	if r.Data == nil {
		return 0
	}

	return len(r.Data)
}

// Get safely retrieves a map entry by key without nil-map panics.
func (r ResultMap[K, V]) Get(key K) (V, bool) {
	if r.Data == nil {
		var zero V

		return zero, false
	}

	val, isFound := r.Data[key]

	return val, isFound
}

// Has checks whether a key exists within the result map.
func (r ResultMap[K, V]) Has(key K) bool {
	if r.Data == nil {
		return false
	}

	_, isFound := r.Data[key]

	return isFound
}

// Keys returns a deterministically sorted slice of all map keys formatted as strings.
func (r ResultMap[K, V]) Keys() []K {
	if r.Data == nil {
		return []K{}
	}

	keys := make([]K, 0, len(r.Data))
	for k := range r.Data {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return fmt.Sprint(keys[i]) < fmt.Sprint(keys[j])
	})

	return keys
}

// Values returns a slice of map values ordered according to sorted Keys().
func (r ResultMap[K, V]) Values() []V {
	keys := r.Keys()
	vals := make([]V, 0, len(keys))
	for _, k := range keys {
		vals = append(vals, r.Data[k])
	}

	return vals
}

// Unwrap returns the underlying map and AppError tuple.
func (r ResultMap[K, V]) Unwrap() (map[K]V, *apperror.AppError) {
	return r.Data, r.Err
}

// UnwrapOr returns the underlying map if successful, or defaultVal if failed.
func (r ResultMap[K, V]) UnwrapOr(defaultVal map[K]V) map[K]V {
	if r.IsSuccess() {
		return r.Data
	}

	return defaultVal
}

// OkMap constructs a successful ResultMap envelope.
func OkMap[K comparable, V any](data map[K]V) ResultMap[K, V] {
	if data == nil {
		data = make(map[K]V)
	}

	return ResultMap[K, V]{
		Value: data,
		Data:  data,
	}
}

// FailMap constructs a failed ResultMap envelope with *apperror.AppError.
func FailMap[K comparable, V any](err *apperror.AppError) ResultMap[K, V] {
	return ResultMap[K, V]{
		Err: err,
	}
}

// NewSuccessMap constructs a successful ResultMap envelope.
func NewSuccessMap[K comparable, V any](data map[K]V) ResultMap[K, V] {
	return OkMap(data)
}

// NewFailureMap constructs a failed ResultMap envelope from any error.
func NewFailureMap[K comparable, V any](err error) ResultMap[K, V] {
	if appErr, isAppErr := err.(*apperror.AppError); isAppErr {
		return FailMap[K, V](appErr)
	}

	if err == nil {
		return ResultMap[K, V]{}
	}

	appErr := apperror.WrapSimple(err, "result.NewFailureMap")

	return FailMap[K, V](appErr)
}
