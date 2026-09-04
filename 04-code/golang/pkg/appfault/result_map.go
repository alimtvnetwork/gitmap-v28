package appfault

// ResultMap wraps a generic key-value map with monadic error state.
type ResultMap[K comparable, V any] struct {
	Data     map[K]V   `json:"Data,omitempty" yaml:"Data,omitempty"`
	AppError *AppError `json:"AppError,omitempty" yaml:"AppError,omitempty"`
}

// OkMap creates a successful ResultMap.
func OkMap[K comparable, V any](data map[K]V) ResultMap[K, V] {
	return ResultMap[K, V]{
		Data: data,
	}
}

// FailMap creates a failed ResultMap from an AppError.
func FailMap[K comparable, V any](err *AppError) ResultMap[K, V] {
	return ResultMap[K, V]{
		AppError: err,
	}
}

// IsSuccess returns true if no error is present.
func (rm ResultMap[K, V]) IsSuccess() bool {
	return rm.AppError == nil
}

// IsFailed returns true if an error is present.
func (rm ResultMap[K, V]) IsFailed() bool {
	return rm.AppError != nil
}

// HasError returns true if an error is present.
func (rm ResultMap[K, V]) HasError() bool {
	return rm.IsFailed()
}

// Has returns true if the key exists in the map.
func (rm ResultMap[K, V]) Has(key K) bool {
	if rm.IsFailed() || rm.Data == nil {
		return false
	}

	_, ok := rm.Data[key]

	return ok
}

// Get retrieves the value associated with key.
func (rm ResultMap[K, V]) Get(key K) (V, bool) {
	if rm.IsFailed() || rm.Data == nil {
		var zero V

		return zero, false
	}

	val, ok := rm.Data[key]

	return val, ok
}

// Count returns the number of entries in the map or 0 if failed.
func (rm ResultMap[K, V]) Count() int {
	if rm.IsFailed() || rm.Data == nil {
		return 0
	}

	return len(rm.Data)
}

// Fault returns the underlying *AppError.
func (rm ResultMap[K, V]) Fault() *AppError {
	return rm.AppError
}

// Error returns the underlying *AppError.
func (rm ResultMap[K, V]) Error() *AppError {
	return rm.AppError
}

// Unwrap unpacks the (map[K]V, *AppError) tuple.
func (rm ResultMap[K, V]) Unwrap() (map[K]V, *AppError) {
	return rm.Data, rm.AppError
}
