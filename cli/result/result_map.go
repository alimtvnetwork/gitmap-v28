package result

import (
	"fmt"
	"sort"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// IsSuccess reports whether the map operation succeeded without error.
func (r *ResultMap[K, V]) IsSuccess() bool {
	if r == nil {
		return false
	}

	return r.Err == nil
}

// IsSafe reports whether the map operation succeeded without error (alias).
func (r *ResultMap[K, V]) IsSafe() bool {
	if r == nil {
		return false
	}

	return r.Err == nil
}

// IsFailed reports whether the map operation encountered an error.
func (r *ResultMap[K, V]) IsFailed() bool {
	if r == nil {
		return true
	}

	return r.Err != nil
}

// IsFailure reports whether the map operation encountered an error.
func (r *ResultMap[K, V]) IsFailure() bool {
	if r == nil {
		return true
	}

	return r.Err != nil
}

// HasError reports whether an active error is attached to the result.
func (r *ResultMap[K, V]) HasError() bool {
	if r == nil {
		return true
	}

	return r.Err != nil
}

// IsEmptyError reports whether no active error exists.
func (r *ResultMap[K, V]) IsEmptyError() bool {
	if r == nil {
		return false
	}

	return r.Err == nil
}

// HasNoError reports whether no active error exists.
func (r *ResultMap[K, V]) HasNoError() bool {
	if r == nil {
		return false
	}

	return r.Err == nil
}

// IsEmpty reports whether the underlying map has 0 items or is uninitialized.
func (r *ResultMap[K, V]) IsEmpty() bool {
	if r == nil || r.Err != nil {
		return true
	}

	return len(r.Data) == 0
}

// Count returns the number of entries in the map, or 0 if uninitialized or failed.
func (r *ResultMap[K, V]) Count() int {
	if r == nil || r.Err != nil || r.Data == nil {
		return 0
	}

	return len(r.Data)
}

// IsCountOtherThan reports whether the operation failed OR the map size != number.
func (r *ResultMap[K, V]) IsCountOtherThan(number int) bool {
	if r == nil || r.IsFailure() {
		return true
	}

	return r.Count() != number
}

// HasRecord reports whether the operation succeeded AND contains more than 0 entries.
func (r *ResultMap[K, V]) HasRecord() bool {
	if r == nil || r.Err != nil || r.Data == nil {
		return false
	}

	return len(r.Data) > 0
}

// HasRecords is an alias for HasRecord.
func (r *ResultMap[K, V]) HasRecords() bool {
	return r.HasRecord()
}

// IsDefined reports whether the operation succeeded AND has entries.
func (r *ResultMap[K, V]) IsDefined() bool {
	if r == nil || r.Err != nil || r.Data == nil {
		return false
	}

	return len(r.Data) > 0
}

// AppError returns the underlying AppError or nil.
func (r *ResultMap[K, V]) AppError() *apperror.AppError {
	if r == nil {
		return nil
	}

	return r.Err
}

// Fault returns the underlying AppError or nil.
func (r *ResultMap[K, V]) Fault() *apperror.AppError {
	if r == nil {
		return nil
	}

	return r.Err
}

// Get safely retrieves a map entry by key without nil-map panics.
func (r *ResultMap[K, V]) Get(key K) (V, bool) {
	if r == nil || r.Data == nil {
		var zero V

		return zero, false
	}

	val, isFound := r.Data[key]

	return val, isFound
}

// Has checks whether a key exists within the result map.
func (r *ResultMap[K, V]) Has(key K) bool {
	if r == nil || r.Data == nil {
		return false
	}

	_, isFound := r.Data[key]

	return isFound
}

// Keys returns a deterministically sorted slice of all map keys formatted as strings.
func (r *ResultMap[K, V]) Keys() []K {
	if r == nil || r.Data == nil {
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
func (r *ResultMap[K, V]) Values() []V {
	if r == nil || r.Data == nil {
		return []V{}
	}

	keys := r.Keys()
	vals := make([]V, 0, len(keys))
	for _, k := range keys {
		vals = append(vals, r.Data[k])
	}

	return vals
}

// Unwrap returns the underlying map and AppError tuple.
func (r *ResultMap[K, V]) Unwrap() (map[K]V, *apperror.AppError) {
	if r == nil {
		return nil, nil
	}

	return r.Data, r.Err
}

// UnwrapOr returns the underlying map if successful, or defaultVal if failed.
func (r *ResultMap[K, V]) UnwrapOr(defaultVal map[K]V) map[K]V {
	if r == nil || r.IsFailure() {
		return defaultVal
	}

	return r.Data
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
