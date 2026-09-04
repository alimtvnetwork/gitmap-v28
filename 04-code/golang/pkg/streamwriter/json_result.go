package streamwriter

import (
	"bytes"
	"encoding/json"
	"io"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
)

// WrappedJSON defines the contract for JSON result envelopes with formatting, unmarshaling, and conversion.
type WrappedJSON[T any] interface {
	WrappedBytes[T]
	Pretty() string
	PrettyOrError() (string, *appfault.AppError)
	Compact() string
	CompactOrError() (string, *appfault.AppError)
	Unmarshal(dest any) *appfault.AppError
	ToBytes() Bytes[T]
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

// FromPayload serializes any typed payload into a JSONResult.
func FromPayload[T any](payload T) JSONResult[T] {
	return NewJSONResult(payload)
}

// FromBytes validates and creates a JSONResult from a pre-existing byte slice.
func FromBytes[T any](data []byte, payload T) JSONResult[T] {
	return NewJSONResultWithBytes(data, payload)
}

// FromString validates and creates a JSONResult from a JSON string.
func FromString[T any](jsonStr string, payload T) JSONResult[T] {
	return NewJSONResultFromString(jsonStr, payload)
}

// FromReader streams and validates data from an io.Reader into a JSONResult.
func FromReader[T any](r io.Reader, payload T) JSONResult[T] {
	return NewJSONResultFromReader(r, payload)
}

// FromSerializer creates a JSONResult from an on-demand serializer closure.
func FromSerializer[T any](serializer func() ([]byte, *appfault.AppError), payload T) JSONResult[T] {
	if serializer == nil {
		appErr := appfault.New(errtype.Validation, "serializer closure cannot be nil")
		return JSONResult[T]{
			payload:    payload,
			status:     false,
			statusCode: 400,
			appError:   appErr,
		}
	}
	data, appErr := serializer()
	if appErr != nil {
		return JSONResult[T]{
			payload:    payload,
			status:     false,
			statusCode: appErr.StatusCode(),
			appError:   appErr,
		}
	}
	return NewJSONResultWithBytes(data, payload)
}

// FromBytesEnvelope converts an existing WrappedBytes envelope into a JSONResult.
func FromBytesEnvelope[T any](wb WrappedBytes[T]) JSONResult[T] {
	if wb == nil {
		appErr := appfault.New(errtype.Validation, "wrapped bytes envelope cannot be nil")
		return JSONResult[T]{
			status:     false,
			statusCode: 400,
			appError:   appErr,
		}
	}
	if wb.HasError() {
		return JSONResult[T]{
			data:       wb.Raw(),
			payload:    wb.Payload(),
			status:     false,
			statusCode: wb.StatusCode(),
			appError:   wb.AppError(),
		}
	}
	return NewJSONResultWithBytes(wb.Raw(), wb.Payload())
}

// FromError creates a failed JSONResult containing an AppError.
func FromError[T any](appErr *appfault.AppError) JSONResult[T] {
	return NewJSONResultError[T](appErr)
}

// FromErrorWithPayload creates a failed JSONResult preserving the payload.
func FromErrorWithPayload[T any](appErr *appfault.AppError, payload T) JSONResult[T] {
	return NewJSONResultErrorWithPayload(appErr, payload)
}

// FromAny polymorphically converts any supported source into a generic JSONResult[any].
func FromAny(source any) JSONResult[any] {
	if source == nil {
		return JSONResult[any]{
			data:       []byte("null"),
			payload:    nil,
			status:     true,
			statusCode: 200,
		}
	}
	switch v := source.(type) {
	case JSONResult[any]:
		return v
	case *JSONResult[any]:
		if v == nil {
			return JSONResult[any]{
				status:     false,
				statusCode: 400,
				appError:   appfault.New(errtype.Validation, "nil JSONResult pointer provided"),
			}
		}
		return *v
	case []byte:
		return FromBytes(v, any(v))
	case string:
		return FromString(v, any(v))
	case io.Reader:
		return FromReader(v, any(v))
	case *appfault.AppError:
		return FromError[any](v)
	default:
		return FromPayload(source)
	}
}

// Cast executes a type-safe JSON round-trip cast from Source to Target.
func Cast[Target any, Source any](source Source) JSONResult[Target] {
	data, err := json.Marshal(source)
	if err != nil {
		appErr := appfault.Wrap(errtype.Validation, err, "failed to marshal source for type cast")
		var zero Target
		return JSONResult[Target]{
			payload:    zero,
			status:     false,
			statusCode: 500,
			appError:   appErr,
		}
	}
	var target Target
	err = json.Unmarshal(data, &target)
	if err != nil {
		appErr := appfault.Wrap(errtype.Validation, err, "failed to unmarshal into target type during cast")
		return JSONResult[Target]{
			data:       data,
			payload:    target,
			status:     false,
			statusCode: 400,
			appError:   appErr,
		}
	}
	return JSONResult[Target]{
		data:       data,
		payload:    target,
		status:     true,
		statusCode: 200,
	}
}

// JSONSourceFactory provides type-parameterized namespacing for JSON source construction.
type JSONSourceFactory[T any] struct{}

// JSONSourceOf creates a scoped factory for payload type T.
func JSONSourceOf[T any]() JSONSourceFactory[T] {
	return JSONSourceFactory[T]{}
}

func (JSONSourceFactory[T]) FromPayload(payload T) JSONResult[T] {
	return FromPayload(payload)
}

func (JSONSourceFactory[T]) FromBytes(data []byte, payload T) JSONResult[T] {
	return FromBytes(data, payload)
}

func (JSONSourceFactory[T]) FromString(jsonStr string, payload T) JSONResult[T] {
	return FromString(jsonStr, payload)
}

func (JSONSourceFactory[T]) FromReader(r io.Reader, payload T) JSONResult[T] {
	return FromReader(r, payload)
}

func (JSONSourceFactory[T]) FromSerializer(serializer func() ([]byte, *appfault.AppError), payload T) JSONResult[T] {
	return FromSerializer(serializer, payload)
}

func (JSONSourceFactory[T]) FromBytesEnvelope(wb WrappedBytes[T]) JSONResult[T] {
	return FromBytesEnvelope(wb)
}

func (JSONSourceFactory[T]) FromError(appErr *appfault.AppError) JSONResult[T] {
	return FromError[T](appErr)
}

func (JSONSourceFactory[T]) FromErrorWithPayload(appErr *appfault.AppError, payload T) JSONResult[T] {
	return FromErrorWithPayload(appErr, payload)
}

// jsonSourceSingleton provides universal untyped helpers.
type jsonSourceSingleton struct{}

// JSONSource is the package-level entry point for universal JSON conversion.
var JSONSource = jsonSourceSingleton{}

func (jsonSourceSingleton) FromAny(source any) JSONResult[any] {
	return FromAny(source)
}

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

// NewJSONResultFromString creates a JSONResult from a JSON string.
func NewJSONResultFromString[T any](jsonStr string, payload T) JSONResult[T] {
	return NewJSONResultWithBytes([]byte(jsonStr), payload)
}

// NewJSONResultFromReader streams JSON bytes from an io.Reader.
func NewJSONResultFromReader[T any](r io.Reader, payload T) JSONResult[T] {
	if r == nil {
		appErr := appfault.New(errtype.Validation, "reader cannot be nil")
		return JSONResult[T]{
			payload:    payload,
			status:     false,
			statusCode: 400,
			appError:   appErr,
		}
	}
	data, err := io.ReadAll(r)
	if err != nil {
		appErr := appfault.Wrap(errtype.IO, err, "failed to read stream data")
		return JSONResult[T]{
			payload:    payload,
			status:     false,
			statusCode: 500,
			appError:   appErr,
		}
	}
	return NewJSONResultWithBytes(data, payload)
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
	str, _ := j.PrettyOrError()
	return str
}

// PrettyOrError returns formatted JSON or an AppError if indentation fails.
func (j JSONResult[T]) PrettyOrError() (string, *appfault.AppError) {
	if len(j.data) == 0 {
		return "{}", nil
	}
	var out bytes.Buffer
	err := json.Indent(&out, j.data, "", "  ")
	if err != nil {
		return string(j.data), appfault.Wrap(errtype.Validation, err, "failed to format pretty JSON")
	}
	return out.String(), nil
}

// Compact returns minified JSON without whitespace.
func (j JSONResult[T]) Compact() string {
	str, _ := j.CompactOrError()
	return str
}

// CompactOrError returns minified JSON or an AppError if compaction fails.
func (j JSONResult[T]) CompactOrError() (string, *appfault.AppError) {
	if len(j.data) == 0 {
		return "{}", nil
	}
	var out bytes.Buffer
	err := json.Compact(&out, j.data)
	if err != nil {
		return string(j.data), appfault.Wrap(errtype.Validation, err, "failed to minify compact JSON")
	}
	return out.String(), nil
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

// ToBytes converts the JSONResult to a standard Bytes[T] envelope.
func (j JSONResult[T]) ToBytes() Bytes[T] {
	return Bytes[T]{
		data:       j.data,
		payload:    j.payload,
		status:     j.status,
		statusCode: j.statusCode,
		appError:   j.appError,
	}
}

var _ WrappedBytes[any] = JSONResult[any]{}
var _ WrappedJSON[any] = JSONResult[any]{}
