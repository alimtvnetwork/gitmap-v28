package streamwriter

import (
	"coding-guidelines/common/pkg/appfault"
)

// Bytes wraps a formatted byte slice bundled with its generic payload T and monadic AppError state.
type Bytes[T any] struct {
	data     []byte
	payload  T
	appError *appfault.AppError
}

// NewBytes creates a successful Bytes envelope.
func NewBytes[T any](data []byte, payload T) Bytes[T] {
	return Bytes[T]{
		data:    data,
		payload: payload,
	}
}

// NewBytesError creates a failed Bytes envelope with an AppError.
func NewBytesError[T any](appErr *appfault.AppError) Bytes[T] {
	return Bytes[T]{
		appError: appErr,
	}
}

// NewBytesErrorWithPayload creates a failed Bytes envelope preserving the original payload.
func NewBytesErrorWithPayload[T any](appErr *appfault.AppError, payload T) Bytes[T] {
	return Bytes[T]{
		payload:  payload,
		appError: appErr,
	}
}

// Raw returns the underlying byte slice.
func (b Bytes[T]) Raw() []byte {
	return b.data
}

// Bytes returns the underlying byte slice (alias to Raw).
func (b Bytes[T]) Bytes() []byte {
	return b.data
}

// String returns the string representation of the bytes.
func (b Bytes[T]) String() string {
	return string(b.data)
}

// Len returns the byte slice length.
func (b Bytes[T]) Len() int {
	return len(b.data)
}

// IsEmpty returns true if the byte slice is empty.
func (b Bytes[T]) IsEmpty() bool {
	return len(b.data) == 0
}

// Payload returns the original generic payload T.
func (b Bytes[T]) Payload() T {
	return b.payload
}

// Value returns the original generic payload T (alias to Payload).
func (b Bytes[T]) Value() T {
	return b.payload
}

// AppError returns the underlying *appfault.AppError.
func (b Bytes[T]) AppError() *appfault.AppError {
	return b.appError
}

// Fault returns the underlying *appfault.AppError (alias to AppError).
func (b Bytes[T]) Fault() *appfault.AppError {
	return b.appError
}

// HasError returns true if an AppError is present.
func (b Bytes[T]) HasError() bool {
	return b.appError != nil
}

// IsValid returns true if no AppError is present.
func (b Bytes[T]) IsValid() bool {
	return b.appError == nil
}

// Unwrap returns both the byte slice and the AppError.
func (b Bytes[T]) Unwrap() ([]byte, *appfault.AppError) {
	return b.data, b.appError
}
