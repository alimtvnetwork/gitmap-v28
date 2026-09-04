package streamwriter

import (
	"coding-guidelines/common/pkg/appfault"
)

// WrappedBytes defines the contract for byte envelopes with status flag and AppError state.
type WrappedBytes[T any] interface {
	Raw() []byte
	Bytes() []byte
	String() string
	Len() int
	IsEmpty() bool
	Payload() T
	Value() T
	AppError() *appfault.AppError
	Fault() *appfault.AppError
	Error() *appfault.AppError
	HasError() bool
	IsValid() bool
	IsSuccess() bool
	Status() bool
	StatusCode() int
	Unwrap() ([]byte, *appfault.AppError)
}

// Bytes wraps a formatted byte slice bundled with its generic payload T, status flag, and AppError state.
type Bytes[T any] struct {
	data       []byte
	payload    T
	status     bool
	statusCode int
	appError   *appfault.AppError
}

// NewBytes creates a successful Bytes envelope with status flag set to true and code 200.
func NewBytes[T any](data []byte, payload T) Bytes[T] {
	return Bytes[T]{
		data:       data,
		payload:    payload,
		status:     true,
		statusCode: 200,
	}
}

// NewBytesWithStatus creates a Bytes envelope with custom status flag and code.
func NewBytesWithStatus[T any](data []byte, payload T, status bool, code int) Bytes[T] {
	return Bytes[T]{
		data:       data,
		payload:    payload,
		status:     status,
		statusCode: code,
	}
}

// NewBytesError creates a failed Bytes envelope with an AppError and status flag set to false.
func NewBytesError[T any](appErr *appfault.AppError) Bytes[T] {
	code := 500
	if appErr != nil {
		if appErr.StatusCode() != 0 {
			code = appErr.StatusCode()
		}
	}
	return Bytes[T]{
		status:     false,
		statusCode: code,
		appError:   appErr,
	}
}

// NewBytesErrorWithPayload creates a failed Bytes envelope preserving the original payload.
func NewBytesErrorWithPayload[T any](appErr *appfault.AppError, payload T) Bytes[T] {
	code := 500
	if appErr != nil {
		if appErr.StatusCode() != 0 {
			code = appErr.StatusCode()
		}
	}
	return Bytes[T]{
		payload:    payload,
		status:     false,
		statusCode: code,
		appError:   appErr,
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

// Error returns the underlying *appfault.AppError (alias to AppError).
func (b Bytes[T]) Error() *appfault.AppError {
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

// IsSuccess returns true if status flag is true and no AppError is present.
func (b Bytes[T]) IsSuccess() bool {
	if b.appError != nil {
		return false
	}
	return b.status
}

// Status returns the boolean status flag.
func (b Bytes[T]) Status() bool {
	return b.status
}

// StatusCode returns the numeric status code.
func (b Bytes[T]) StatusCode() int {
	return b.statusCode
}

// Unwrap returns both the byte slice and the AppError.
func (b Bytes[T]) Unwrap() ([]byte, *appfault.AppError) {
	return b.data, b.appError
}

var _ WrappedBytes[any] = Bytes[any]{}
