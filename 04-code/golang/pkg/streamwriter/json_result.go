package streamwriter

import (
	"bytes"
	"encoding/json"
	"io"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
)

// WrappedJson defines the contract for JSON result envelopes with formatting, unmarshaling, and conversion.
type WrappedJson[T any] interface {
	WrappedBytes[T]
	Pretty() string
	PrettyOrError() (string, *appfault.AppError)
	Compact() string
	CompactOrError() (string, *appfault.AppError)
	Unmarshal(dest any) *appfault.AppError
	ToBytes() Bytes[T]
}

// WrappedJSON is an alias for WrappedJson for backwards compatibility.
type WrappedJSON[T any] = WrappedJson[T]

// JsonResult encapsulates JSON serialized data, generic payload, status flag, and AppError state.
type JsonResult[T any] struct {
	data       []byte
	payload    T
	status     bool
	statusCode int
	appError   *appfault.AppError
}

// JSONResult is an alias for JsonResult for backwards compatibility.
type JSONResult[T any] = JsonResult[T]

// jsonSourceSingleton acts as the struct-as-namespace factory taking `any` for high ergonomics.
type jsonSourceSingleton struct{}

// JsonSource is the global factory instance for dynamic and untyped sources.
var JsonSource = jsonSourceSingleton{}

// JSONSource is an alias for JsonSource for backwards compatibility.
var JSONSource = JsonSource

// FromPayload serializes any value into a JsonResult[any].
func (jsonSourceSingleton) FromPayload(payload any) JsonResult[any] {
	return NewJsonResult(payload)
}

// FromBytes validates and wraps raw JSON bytes directly into a JsonResult[any] without needing a dummy payload.
func (jsonSourceSingleton) FromBytes(data []byte) JsonResult[any] {
	return NewJsonResultWithBytes(data, any(data))
}

// FromBytesWithPayload validates JSON bytes and attaches an explicit payload.
func (jsonSourceSingleton) FromBytesWithPayload(data []byte, payload any) JsonResult[any] {
	return NewJsonResultWithBytes(data, payload)
}

// FromString validates and wraps a JSON string directly into a JsonResult[any].
func (jsonSourceSingleton) FromString(jsonStr string) JsonResult[any] {
	return NewJsonResultFromString(jsonStr, any(jsonStr))
}

// FromReader streams and validates data from an io.Reader directly into a JsonResult[any].
func (jsonSourceSingleton) FromReader(r io.Reader) JsonResult[any] {
	return NewJsonResultFromReader(r, any(nil))
}

// FromSerializer executes a lazy closure and wraps the resulting JSON bytes into a JsonResult[any].
func (jsonSourceSingleton) FromSerializer(serializer func() ([]byte, *appfault.AppError)) JsonResult[any] {
	if serializer == nil {
		appErr := appfault.New(errtype.Validation, "serializer closure cannot be nil")
		return JsonResult[any]{
			status:     false,
			statusCode: 400,
			appError:   appErr,
		}
	}
	data, appErr := serializer()
	if appErr != nil {
		return JsonResult[any]{
			status:     false,
			statusCode: appErr.StatusCode(),
			appError:   appErr,
		}
	}
	return NewJsonResultWithBytes(data, any(data))
}

// FromBytesEnvelope converts an existing WrappedBytes envelope into a JsonResult[any].
func (jsonSourceSingleton) FromBytesEnvelope(wb any) JsonResult[any] {
	if wb == nil {
		appErr := appfault.New(errtype.Validation, "wrapped bytes envelope cannot be nil")
		return JsonResult[any]{
			status:     false,
			statusCode: 400,
			appError:   appErr,
		}
	}
	if rawProvider, ok := wb.(interface{ Raw() []byte }); ok {
		raw := rawProvider.Raw()
		if errProvider, ok := wb.(interface{ AppError() *appfault.AppError }); ok {
			if appErr := errProvider.AppError(); appErr != nil {
				return JsonResult[any]{
					data:       raw,
					payload:    wb,
					status:     false,
					statusCode: appErr.StatusCode(),
					appError:   appErr,
				}
			}
		}
		return NewJsonResultWithBytes(raw, any(raw))
	}
	return FromAny(wb)
}

// FromError creates a failed JsonResult containing an AppError.
func (jsonSourceSingleton) FromError(appErr *appfault.AppError) JsonResult[any] {
	return NewJsonResultError[any](appErr)
}

// FromErrorWithPayload creates a failed JsonResult preserving the payload.
func (jsonSourceSingleton) FromErrorWithPayload(appErr *appfault.AppError, payload any) JsonResult[any] {
	return NewJsonResultErrorWithPayload(appErr, payload)
}

// FromAny polymorphically converts any supported source into a generic JsonResult[any].
func (jsonSourceSingleton) FromAny(source any) JsonResult[any] {
	return FromAny(source)
}

// Cast executes a type-safe JSON round-trip unmarshaling directly into targetPtr.
func (jsonSourceSingleton) Cast(source any, targetPtr any) *appfault.AppError {
	data, err := json.Marshal(source)
	if err != nil {
		return appfault.Wrap(errtype.Validation, err, "failed to marshal source for type cast")
	}
	err = json.Unmarshal(data, targetPtr)
	if err != nil {
		return appfault.Wrap(errtype.Validation, err, "failed to unmarshal into target pointer during cast")
	}
	return nil
}

// typedJsonSource provides typed factory methods bound to type parameter T.
type typedJsonSource[T any] struct{}

// JsonSourceOf creates a scoped factory for producing JsonResult[T] instances.
func JsonSourceOf[T any]() typedJsonSource[T] {
	return typedJsonSource[T]{}
}

// JSONSourceOf is an alias for JsonSourceOf for backwards compatibility.
func JSONSourceOf[T any]() typedJsonSource[T] {
	return JsonSourceOf[T]()
}

// FromPayload serializes payload T into a JsonResult[T].
func (typedJsonSource[T]) FromPayload(payload T) JsonResult[T] {
	return FromPayload(payload)
}

// FromBytes validates and creates a JsonResult[T] from a byte slice and payload T.
func (typedJsonSource[T]) FromBytes(data []byte, payload T) JsonResult[T] {
	return FromBytes(data, payload)
}

// FromString validates and creates a JsonResult[T] from a string and payload T.
func (typedJsonSource[T]) FromString(jsonStr string, payload T) JsonResult[T] {
	return FromString(jsonStr, payload)
}

// FromReader streams and validates data from an io.Reader into a JsonResult[T].
func (typedJsonSource[T]) FromReader(r io.Reader, payload T) JsonResult[T] {
	return FromReader(r, payload)
}

// FromSerializer creates a JsonResult[T] from a serializer closure and payload T.
func (typedJsonSource[T]) FromSerializer(serializer func() ([]byte, *appfault.AppError), payload T) JsonResult[T] {
	return FromSerializer(serializer, payload)
}

// FromBytesEnvelope converts an existing WrappedBytes[T] envelope into a JsonResult[T].
func (typedJsonSource[T]) FromBytesEnvelope(wb WrappedBytes[T]) JsonResult[T] {
	return FromBytesEnvelope(wb)
}

// FromError creates a failed JsonResult[T] with an AppError.
func (typedJsonSource[T]) FromError(appErr *appfault.AppError) JsonResult[T] {
	return FromError[T](appErr)
}

// FromErrorWithPayload creates a failed JsonResult[T] preserving payload T.
func (typedJsonSource[T]) FromErrorWithPayload(appErr *appfault.AppError, payload T) JsonResult[T] {
	return FromErrorWithPayload(appErr, payload)
}

// CastTo executes a type-safe JSON round-trip cast from Source to Target.
func CastTo[Target any](source any) JsonResult[Target] {
	return Cast[Target](source)
}

// FromPayload serializes any typed payload into a JsonResult[T].
func FromPayload[T any](payload T) JsonResult[T] {
	return NewJsonResult(payload)
}

// FromBytes validates and creates a JsonResult from a pre-existing byte slice and payload.
func FromBytes[T any](data []byte, payload T) JsonResult[T] {
	return NewJsonResultWithBytes(data, payload)
}

// FromString validates and creates a JsonResult from a JSON string and payload.
func FromString[T any](jsonStr string, payload T) JsonResult[T] {
	return NewJsonResultFromString(jsonStr, payload)
}

// FromReader streams and validates data from an io.Reader into a JsonResult.
func FromReader[T any](r io.Reader, payload T) JsonResult[T] {
	return NewJsonResultFromReader(r, payload)
}

// FromSerializer creates a JsonResult from an on-demand serializer closure.
func FromSerializer[T any](serializer func() ([]byte, *appfault.AppError), payload T) JsonResult[T] {
	if serializer == nil {
		appErr := appfault.New(errtype.Validation, "serializer closure cannot be nil")
		return JsonResult[T]{
			payload:    payload,
			status:     false,
			statusCode: 400,
			appError:   appErr,
		}
	}
	data, appErr := serializer()
	if appErr != nil {
		return JsonResult[T]{
			payload:    payload,
			status:     false,
			statusCode: appErr.StatusCode(),
			appError:   appErr,
		}
	}
	return NewJsonResultWithBytes(data, payload)
}

// FromBytesEnvelope converts an existing WrappedBytes envelope into a JsonResult[T].
func FromBytesEnvelope[T any](wb WrappedBytes[T]) JsonResult[T] {
	if wb == nil {
		appErr := appfault.New(errtype.Validation, "wrapped bytes envelope cannot be nil")
		return JsonResult[T]{
			status:     false,
			statusCode: 400,
			appError:   appErr,
		}
	}
	if wb.HasError() {
		return JsonResult[T]{
			data:       wb.Raw(),
			payload:    wb.Payload(),
			status:     false,
			statusCode: wb.StatusCode(),
			appError:   wb.AppError(),
		}
	}
	return NewJsonResultWithBytes(wb.Raw(), wb.Payload())
}

// FromError creates a failed JsonResult containing an AppError.
func FromError[T any](appErr *appfault.AppError) JsonResult[T] {
	return NewJsonResultError[T](appErr)
}

// FromErrorWithPayload creates a failed JsonResult preserving the payload.
func FromErrorWithPayload[T any](appErr *appfault.AppError, payload T) JsonResult[T] {
	return NewJsonResultErrorWithPayload(appErr, payload)
}

// FromAny polymorphically converts any supported source into a generic JsonResult[any].
func FromAny(source any) JsonResult[any] {
	if source == nil {
		return JsonResult[any]{
			data:       []byte("null"),
			payload:    nil,
			status:     true,
			statusCode: 200,
		}
	}
	switch v := source.(type) {
	case JsonResult[any]:
		return v
	case *JsonResult[any]:
		if v == nil {
			return JsonResult[any]{
				status:     false,
				statusCode: 400,
				appError:   appfault.New(errtype.Validation, "nil JsonResult pointer provided"),
			}
		}
		return *v
	case []byte:
		return JsonSource.FromBytes(v)
	case string:
		return JsonSource.FromString(v)
	case io.Reader:
		return JsonSource.FromReader(v)
	case *appfault.AppError:
		return JsonSource.FromError(v)
	default:
		return JsonSource.FromPayload(source)
	}
}

// Cast executes a type-safe JSON round-trip cast from Source to Target.
func Cast[Target any, Source any](source Source) JsonResult[Target] {
	data, err := json.Marshal(source)
	if err != nil {
		appErr := appfault.Wrap(errtype.Validation, err, "failed to marshal source for type cast")
		var zero Target
		return JsonResult[Target]{
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
		return JsonResult[Target]{
			data:       data,
			payload:    target,
			status:     false,
			statusCode: 400,
			appError:   appErr,
		}
	}
	return JsonResult[Target]{
		data:       data,
		payload:    target,
		status:     true,
		statusCode: 200,
	}
}

// NewJsonResult serializes payload T into JSON and initializes a JsonResult envelope.
func NewJsonResult[T any](payload T) JsonResult[T] {
	data, err := json.Marshal(payload)
	if err != nil {
		appErr := appfault.Wrap(errtype.Validation, err, "failed to marshal payload into JSON")
		return JsonResult[T]{
			payload:    payload,
			status:     false,
			statusCode: 500,
			appError:   appErr,
		}
	}
	return JsonResult[T]{
		data:       data,
		payload:    payload,
		status:     true,
		statusCode: 200,
	}
}

// NewJsonResultWithBytes creates a JsonResult from pre-marshaled JSON bytes and payload.
func NewJsonResultWithBytes[T any](data []byte, payload T) JsonResult[T] {
	if !json.Valid(data) {
		appErr := appfault.New(errtype.Validation, "invalid JSON byte sequence provided")
		return JsonResult[T]{
			data:       data,
			payload:    payload,
			status:     false,
			statusCode: 400,
			appError:   appErr,
		}
	}
	return JsonResult[T]{
		data:       data,
		payload:    payload,
		status:     true,
		statusCode: 200,
	}
}

// NewJsonResultFromString creates a JsonResult from a JSON string.
func NewJsonResultFromString[T any](jsonStr string, payload T) JsonResult[T] {
	return NewJsonResultWithBytes([]byte(jsonStr), payload)
}

// NewJsonResultFromReader streams JSON bytes from an io.Reader.
func NewJsonResultFromReader[T any](r io.Reader, payload T) JsonResult[T] {
	if r == nil {
		appErr := appfault.New(errtype.Validation, "reader cannot be nil")
		return JsonResult[T]{
			payload:    payload,
			status:     false,
			statusCode: 400,
			appError:   appErr,
		}
	}
	data, err := io.ReadAll(r)
	if err != nil {
		appErr := appfault.Wrap(errtype.IO, err, "failed to read stream data")
		return JsonResult[T]{
			payload:    payload,
			status:     false,
			statusCode: 500,
			appError:   appErr,
		}
	}
	return NewJsonResultWithBytes(data, payload)
}

// NewJsonResultWithStatus creates a JsonResult with explicit status flag and code.
func NewJsonResultWithStatus[T any](data []byte, payload T, status bool, code int) JsonResult[T] {
	return JsonResult[T]{
		data:       data,
		payload:    payload,
		status:     status,
		statusCode: code,
	}
}

// NewJsonResultError creates a failed JsonResult with an AppError.
func NewJsonResultError[T any](appErr *appfault.AppError) JsonResult[T] {
	code := 500
	if appErr != nil {
		if appErr.StatusCode() != 0 {
			code = appErr.StatusCode()
		}
	}
	return JsonResult[T]{
		status:     false,
		statusCode: code,
		appError:   appErr,
	}
}

// NewJsonResultErrorWithPayload creates a failed JsonResult preserving the payload.
func NewJsonResultErrorWithPayload[T any](appErr *appfault.AppError, payload T) JsonResult[T] {
	code := 500
	if appErr != nil {
		if appErr.StatusCode() != 0 {
			code = appErr.StatusCode()
		}
	}
	return JsonResult[T]{
		payload:    payload,
		status:     false,
		statusCode: code,
		appError:   appErr,
	}
}

// Backwards-compatible aliases for PascalCase JSON
func NewJSONResult[T any](payload T) JsonResult[T] {
	return NewJsonResult(payload)
}
func NewJSONResultWithBytes[T any](data []byte, payload T) JsonResult[T] {
	return NewJsonResultWithBytes(data, payload)
}
func NewJSONResultFromString[T any](jsonStr string, payload T) JsonResult[T] {
	return NewJsonResultFromString(jsonStr, payload)
}
func NewJSONResultFromReader[T any](r io.Reader, payload T) JsonResult[T] {
	return NewJsonResultFromReader(r, payload)
}
func NewJSONResultError[T any](appErr *appfault.AppError) JsonResult[T] {
	return NewJsonResultError[T](appErr)
}
func NewJSONResultErrorWithPayload[T any](appErr *appfault.AppError, payload T) JsonResult[T] {
	return NewJsonResultErrorWithPayload(appErr, payload)
}
func NewJSONResultWithStatus[T any](data []byte, payload T, status bool, code int) JsonResult[T] {
	return NewJsonResultWithStatus(data, payload, status, code)
}

// Raw returns the underlying JSON byte slice.
func (j JsonResult[T]) Raw() []byte {
	return j.data
}

// Bytes returns the underlying JSON byte slice (alias to Raw).
func (j JsonResult[T]) Bytes() []byte {
	return j.data
}

// String returns the JSON string representation.
func (j JsonResult[T]) String() string {
	return string(j.data)
}

// Len returns the byte length of the JSON.
func (j JsonResult[T]) Len() int {
	return len(j.data)
}

// IsEmpty returns true if the JSON bytes are empty.
func (j JsonResult[T]) IsEmpty() bool {
	return len(j.data) == 0
}

// Payload returns the original generic payload T.
func (j JsonResult[T]) Payload() T {
	return j.payload
}

// Value returns the original generic payload T (alias to Payload).
func (j JsonResult[T]) Value() T {
	return j.payload
}

// AppError returns the underlying *appfault.AppError.
func (j JsonResult[T]) AppError() *appfault.AppError {
	return j.appError
}

// Fault returns the underlying *appfault.AppError (alias to AppError).
func (j JsonResult[T]) Fault() *appfault.AppError {
	return j.appError
}

// Error returns the underlying *appfault.AppError (alias to AppError).
func (j JsonResult[T]) Error() *appfault.AppError {
	return j.appError
}

// HasError returns true if an AppError is present.
func (j JsonResult[T]) HasError() bool {
	return j.appError != nil
}

// IsValid returns true if no AppError is present.
func (j JsonResult[T]) IsValid() bool {
	return j.appError == nil
}

// IsSuccess returns true if status flag is true and no AppError is present.
func (j JsonResult[T]) IsSuccess() bool {
	if j.appError != nil {
		return false
	}
	return j.status
}

// Status returns the boolean status flag.
func (j JsonResult[T]) Status() bool {
	return j.status
}

// StatusCode returns the numeric status code.
func (j JsonResult[T]) StatusCode() int {
	return j.statusCode
}

// Unwrap returns both the JSON byte slice and the AppError.
func (j JsonResult[T]) Unwrap() ([]byte, *appfault.AppError) {
	return j.data, j.appError
}

// Pretty returns formatted and indented JSON.
func (j JsonResult[T]) Pretty() string {
	str, _ := j.PrettyOrError()
	return str
}

// PrettyOrError returns formatted JSON or an AppError if indentation fails.
func (j JsonResult[T]) PrettyOrError() (string, *appfault.AppError) {
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
func (j JsonResult[T]) Compact() string {
	str, _ := j.CompactOrError()
	return str
}

// CompactOrError returns minified JSON or an AppError if compaction fails.
func (j JsonResult[T]) CompactOrError() (string, *appfault.AppError) {
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
func (j JsonResult[T]) Unmarshal(dest any) *appfault.AppError {
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

// ToBytes converts the JsonResult to a standard Bytes[T] envelope.
func (j JsonResult[T]) ToBytes() Bytes[T] {
	return Bytes[T]{
		data:       j.data,
		payload:    j.payload,
		status:     j.status,
		statusCode: j.statusCode,
		appError:   j.appError,
	}
}

var _ WrappedBytes[any] = JsonResult[any]{}
var _ WrappedJson[any] = JsonResult[any]{}
