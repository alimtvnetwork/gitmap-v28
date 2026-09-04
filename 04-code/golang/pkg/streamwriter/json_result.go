package streamwriter

import (
	"bytes"
	"encoding/json"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
)

// WrappedJSON defines the contract for JSON result envelopes with formatting and unmarshaling.
type WrappedJSON[T any] interface {
	WrappedBytes[T]
	Pretty() string
	Compact() string
	Unmarshal(dest any) *appfault.AppError
}

// JSONResult encapsulates JSON serialized data, generic payload, status flag, and AppError state.
type JSONResult[T any] struct {
	data       []byte
	payload    T
	status     bool
	statusCode int
	appError   *appfault.AppError
}

// JsonResult is an alias for JSONResult for casing flexibility.
type JsonResult[T any] = JSONResult[T]

// NewJSONResult serializes payload T into JSON and initializes a JSONResult envelope.
func NewJSONResult[T any](payload T) JSONResult[T] {
	data, err := json.Marshal(payload)
	if err != nil {
		appErr := appfault.Wrap(errtype.Validation, err, "failed to marshal payload into JSON")
		return JSONResult[T]{
			payload:    payload,
			status:     false,
			statusCode: 500,
			appError:   appErr,
		}
	}
	return JSONResult[T]{
		data:       data,
		payload:    payload,
		status:     true,
		statusCode: 200,
	}
}

// NewJSONResultWithBytes creates a JSONResult from pre-marshaled JSON bytes and payload.
func NewJSONResultWithBytes[T any](data []byte, payload T) JSONResult[T] {
	if !json.Valid(data) {
		appErr := appfault.New(errtype.Validation, "invalid JSON byte sequence provided")
		return JSONResult[T]{
			data:       data,
			payload:    payload,
			status:     false,
			statusCode: 400,
			appError:   appErr,
		}
	}
	return JSONResult[T]{
		data:       data,
		payload:    payload,
		status:     true,
		statusCode: 200,
	}
}

// NewJSONResultWithStatus creates a JSONResult with explicit status flag and code.
func NewJSONResultWithStatus[T any](data []byte, payload T, status bool, code int) JSONResult[T] {
	return JSONResult[T]{
		data:       data,
		payload:    payload,
		status:     status,
		statusCode: code,
	}
}

// NewJSONResultError creates a failed JSONResult with an AppError.
func NewJSONResultError[T any](appErr *appfault.AppError) JSONResult[T] {
	code := 500
	if appErr != nil {
		if appErr.StatusCode() != 0 {
			code = appErr.StatusCode()
		}
	}
	return JSONResult[T]{
		status:     false,
		statusCode: code,
		appError:   appErr,
	}
}

// NewJSONResultErrorWithPayload creates a failed JSONResult preserving the payload.
func NewJSONResultErrorWithPayload[T any](appErr *appfault.AppError, payload T) JSONResult[T] {
	code := 500
	if appErr != nil {
		if appErr.StatusCode() != 0 {
			code = appErr.StatusCode()
		}
	}
	return JSONResult[T]{
		payload:    payload,
		status:     false,
		statusCode: code,
		appError:   appErr,
	}
}

// Raw returns the underlying JSON byte slice.
func (j JSONResult[T]) Raw() []byte {
	return j.data
}

// Bytes returns the underlying JSON byte slice (alias to Raw).
func (j JSONResult[T]) Bytes() []byte {
	return j.data
}

// String returns the JSON string representation.
func (j JSONResult[T]) String() string {
	return string(j.data)
}

// Len returns the byte length of the JSON.
func (j JSONResult[T]) Len() int {
	return len(j.data)
}

// IsEmpty returns true if the JSON bytes are empty.
func (j JSONResult[T]) IsEmpty() bool {
	return len(j.data) == 0
}

// Payload returns the original generic payload T.
func (j JSONResult[T]) Payload() T {
	return j.payload
}

// Value returns the original generic payload T (alias to Payload).
func (j JSONResult[T]) Value() T {
	return j.payload
}

// AppError returns the underlying *appfault.AppError.
func (j JSONResult[T]) AppError() *appfault.AppError {
	return j.appError
}

// Fault returns the underlying *appfault.AppError (alias to AppError).
func (j JSONResult[T]) Fault() *appfault.AppError {
	return j.appError
}

// Error returns the underlying *appfault.AppError (alias to AppError).
func (j JSONResult[T]) Error() *appfault.AppError {
	return j.appError
}

// HasError returns true if an AppError is present.
func (j JSONResult[T]) HasError() bool {
	return j.appError != nil
}

// IsValid returns true if no AppError is present.
func (j JSONResult[T]) IsValid() bool {
	return j.appError == nil
}

// IsSuccess returns true if status flag is true and no AppError is present.
func (j JSONResult[T]) IsSuccess() bool {
	if j.appError != nil {
		return false
	}
	return j.status
}

// Status returns the boolean status flag.
func (j JSONResult[T]) Status() bool {
	return j.status
}

// StatusCode returns the numeric status code.
func (j JSONResult[T]) StatusCode() int {
	return j.statusCode
}

// Unwrap returns both the JSON byte slice and the AppError.
func (j JSONResult[T]) Unwrap() ([]byte, *appfault.AppError) {
	return j.data, j.appError
}

// Pretty returns formatted and indented JSON.
func (j JSONResult[T]) Pretty() string {
	if len(j.data) == 0 {
		return "{}"
	}
	var out bytes.Buffer
	err := json.Indent(&out, j.data, "", "  ")
	if err != nil {
		return string(j.data)
	}
	return out.String()
}

// Compact returns minified/compact JSON without whitespace.
func (j JSONResult[T]) Compact() string {
	if len(j.data) == 0 {
		return "{}"
	}
	var out bytes.Buffer
	err := json.Compact(&out, j.data)
	if err != nil {
		return string(j.data)
	}
	return out.String()
}

// Unmarshal parses the JSON bytes into the destination pointer.
func (j JSONResult[T]) Unmarshal(dest any) *appfault.AppError {
	if j.appError != nil {
		return j.appError
	}
	if len(j.data) == 0 {
		return appfault.New(errtype.Validation, "cannot unmarshal empty JSON result")
	}
	err := json.Unmarshal(j.data, dest)
	if err != nil {
		return appfault.Wrap(errtype.Validation, err, "failed to unmarshal JSON into destination")
	}
	return nil
}

var _ WrappedBytes[any] = JSONResult[any]{}
var _ WrappedJSON[any] = JSONResult[any]{}
