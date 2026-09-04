package appfault

// ResultSlice wraps a generic slice collection with monadic error state.
type ResultSlice[T any] struct {
	Items    []T       `json:"Items,omitempty" yaml:"Items,omitempty"`
	AppError *AppError `json:"AppError,omitempty" yaml:"AppError,omitempty"`
}

// OkSlice creates a successful ResultSlice.
func OkSlice[T any](items []T) ResultSlice[T] {
	return ResultSlice[T]{
		Items: items,
	}
}

// FailSlice creates a failed ResultSlice from an AppError.
func FailSlice[T any](err *AppError) ResultSlice[T] {
	return ResultSlice[T]{
		AppError: err,
	}
}

// IsSuccess returns true if no error is present.
func (rs ResultSlice[T]) IsSuccess() bool {
	return rs.AppError == nil
}

// IsFailed returns true if an error is present.
func (rs ResultSlice[T]) IsFailed() bool {
	return rs.AppError != nil
}

// HasError returns true if an error is present.
func (rs ResultSlice[T]) HasError() bool {
	return rs.IsFailed()
}

// HasItems returns true if the slice contains elements and is safe.
func (rs ResultSlice[T]) HasItems() bool {
	if rs.IsFailed() {
		return false
	}

	return len(rs.Items) > 0
}

// Count returns the number of items or 0 if failed.
func (rs ResultSlice[T]) Count() int {
	if rs.IsFailed() {
		return 0
	}

	return len(rs.Items)
}

// Length is an alias for Count.
func (rs ResultSlice[T]) Length() int {
	return rs.Count()
}

// Fault returns the underlying *AppError.
func (rs ResultSlice[T]) Fault() *AppError {
	return rs.AppError
}

// Error returns the underlying *AppError.
func (rs ResultSlice[T]) Error() *AppError {
	return rs.AppError
}

// Unwrap unpacks the ([]T, *AppError) tuple.
func (rs ResultSlice[T]) Unwrap() ([]T, *AppError) {
	return rs.Items, rs.AppError
}
