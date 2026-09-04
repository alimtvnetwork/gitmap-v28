package streamwriter

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
)

// WrappedJson defines the contract for JSON result envelopes with formatting, unmarshaling, and conversion.
type WrappedJson interface {
	WrappedBytes[any]
	Pretty() string
	PrettyOrError() (string, *appfault.AppError)
	Compact() string
	CompactOrError() (string, *appfault.AppError)
	Unmarshal(dest any) *appfault.AppError
	ToBytes() Bytes[any]
}

// WrappedJSON is an alias for WrappedJson for backwards compatibility.
type WrappedJSON = WrappedJson

// JsonResult encapsulates JSON serialized data, optional payload, status flag, and AppError state without type parameter T.
type JsonResult struct {
	data       []byte
	payload    any
	status     bool
	statusCode int
	appError   *appfault.AppError
}

// JSONResult is an alias for JsonResult for backwards compatibility.
type JSONResult = JsonResult

// jsonSourceSingleton acts as the struct-as-namespace factory for multi-source JsonResult construction.
type jsonSourceSingleton struct{}

// JsonSource is the global factory instance for dynamic and untyped sources.
var JsonSource = jsonSourceSingleton{}

// JSONSource is an alias for JsonSource for backwards compatibility.
var JSONSource = JsonSource

// FromPayload serializes any value into a JsonResult.
func (jsonSourceSingleton) FromPayload(payload any) JsonResult {
	return NewJsonResult(payload)
}

// FromBytes validates and wraps raw JSON bytes directly into a JsonResult.
func (jsonSourceSingleton) FromBytes(data []byte, payload ...any) JsonResult {
	return NewJsonResultWithBytes(data, payload...)
}

// FromBytesWithPayload validates JSON bytes and attaches an explicit payload.
func (jsonSourceSingleton) FromBytesWithPayload(data []byte, payload any) JsonResult {
	return NewJsonResultWithBytes(data, payload)
}

// FromString validates and wraps a JSON string into a JsonResult.
func (jsonSourceSingleton) FromString(jsonStr string, payload ...any) JsonResult {
	return NewJsonResultFromString(jsonStr, payload...)
}

// FromStringWithPayload validates a JSON string and attaches an explicit payload.
func (jsonSourceSingleton) FromStringWithPayload(jsonStr string, payload any) JsonResult {
	return NewJsonResultFromString(jsonStr, payload)
}

// FromReader streams and validates data from an io.Reader into a JsonResult.
func (jsonSourceSingleton) FromReader(r io.Reader, payload ...any) JsonResult {
	return NewJsonResultFromReader(r, payload...)
}

// FromReaderWithPayload streams data from an io.Reader and attaches an explicit payload.
func (jsonSourceSingleton) FromReaderWithPayload(r io.Reader, payload any) JsonResult {
	return NewJsonResultFromReader(r, payload)
}

// FromSerializer executes a lazy closure and wraps the resulting JSON bytes into a JsonResult.
func (jsonSourceSingleton) FromSerializer(serializer func() ([]byte, *appfault.AppError), payload ...any) JsonResult {
	var p any
	if len(payload) > 0 {
		p = payload[0]
	}
	if serializer == nil {
		appErr := appfault.New(errtype.Validation, "serializer closure cannot be nil")
		return JsonResult{
			payload:    p,
			status:     false,
			statusCode: 400,
			appError:   appErr,
		}
	}
	data, appErr := serializer()
	if appErr != nil {
		return JsonResult{
			payload:    p,
			status:     false,
			statusCode: appErr.StatusCode(),
			appError:   appErr,
		}
	}
	if len(payload) == 0 {
		p = data
	}
	return NewJsonResultWithBytes(data, p)
}

// FromSerializerWithPayload executes a lazy closure and attaches an explicit payload.
func (jsonSourceSingleton) FromSerializerWithPayload(serializer func() ([]byte, *appfault.AppError), payload any) JsonResult {
	return JsonSource.FromSerializer(serializer, payload)
}

// FromBytesEnvelope converts an existing WrappedBytes envelope into a JsonResult.
func (jsonSourceSingleton) FromBytesEnvelope(wb any) JsonResult {
	if wb == nil {
		appErr := appfault.New(errtype.Validation, "wrapped bytes envelope cannot be nil")
		return JsonResult{
			status:     false,
			statusCode: 400,
			appError:   appErr,
		}
	}
	if rawProvider, ok := wb.(interface{ Raw() []byte }); ok {
		raw := rawProvider.Raw()
		if errProvider, ok := wb.(interface{ AppError() *appfault.AppError }); ok {
			if appErr := errProvider.AppError(); appErr != nil {
				return JsonResult{
					data:       raw,
					payload:    wb,
					status:     false,
					statusCode: appErr.StatusCode(),
					appError:   appErr,
				}
			}
		}
		var payload any = raw
		val := reflect.ValueOf(wb)
		if val.IsValid() {
			method := val.MethodByName("Payload")
			if method.IsValid() {
				if method.Type().NumIn() == 0 {
					if method.Type().NumOut() == 1 {
						out := method.Call(nil)
						if len(out) > 0 {
							payload = out[0].Interface()
						}
					}
				}
			} else {
				vMethod := val.MethodByName("Value")
				if vMethod.IsValid() {
					if vMethod.Type().NumIn() == 0 {
						if vMethod.Type().NumOut() == 1 {
							out := vMethod.Call(nil)
							if len(out) > 0 {
								payload = out[0].Interface()
							}
						}
					}
				}
			}
		}
		return NewJsonResultWithBytes(raw, payload)
	}
	return FromAny(wb)
}

// FromError creates a failed JsonResult containing an AppError.
func (jsonSourceSingleton) FromError(appErr *appfault.AppError) JsonResult {
	return NewJsonResultError(appErr)
}

// FromErrorWithPayload creates a failed JsonResult preserving the payload.
func (jsonSourceSingleton) FromErrorWithPayload(appErr *appfault.AppError, payload any) JsonResult {
	return NewJsonResultErrorWithPayload(appErr, payload)
}

// FromAny polymorphically converts any supported source into a JsonResult.
func (jsonSourceSingleton) FromAny(source any) JsonResult {
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

// JsonSourceOf creates a scoped factory for producing JsonResult instances with payload T.
func JsonSourceOf[T any]() typedJsonSource[T] {
	return typedJsonSource[T]{}
}

// JSONSourceOf is an alias for JsonSourceOf for backwards compatibility.
func JSONSourceOf[T any]() typedJsonSource[T] {
	return JsonSourceOf[T]()
}

// FromPayload serializes payload T into a JsonResult.
func (typedJsonSource[T]) FromPayload(payload T) JsonResult {
	return FromPayload(payload)
}

// FromBytes validates and creates a JsonResult from a byte slice and payload T.
func (typedJsonSource[T]) FromBytes(data []byte, payload T) JsonResult {
	return FromBytes(data, payload)
}

// FromString validates and creates a JsonResult from a string and payload T.
func (typedJsonSource[T]) FromString(jsonStr string, payload T) JsonResult {
	return FromString(jsonStr, payload)
}

// FromReader streams and validates data from an io.Reader into a JsonResult.
func (typedJsonSource[T]) FromReader(r io.Reader, payload T) JsonResult {
	return FromReader(r, payload)
}

// FromSerializer creates a JsonResult from a serializer closure and payload T.
func (typedJsonSource[T]) FromSerializer(serializer func() ([]byte, *appfault.AppError), payload T) JsonResult {
	return FromSerializer(serializer, payload)
}

// FromBytesEnvelope converts an existing WrappedBytes envelope into a JsonResult.
func (typedJsonSource[T]) FromBytesEnvelope(wb any) JsonResult {
	return FromBytesEnvelope(wb)
}

// FromError creates a failed JsonResult with an AppError.
func (typedJsonSource[T]) FromError(appErr *appfault.AppError) JsonResult {
	return FromError(appErr)
}

// FromErrorWithPayload creates a failed JsonResult preserving payload T.
func (typedJsonSource[T]) FromErrorWithPayload(appErr *appfault.AppError, payload T) JsonResult {
	return FromErrorWithPayload(appErr, payload)
}

// FromPayload serializes any typed payload into a JsonResult.
func FromPayload(payload any) JsonResult {
	return NewJsonResult(payload)
}

// FromBytes validates and creates a JsonResult from a pre-existing byte slice and optional payload.
func FromBytes(data []byte, payload ...any) JsonResult {
	return NewJsonResultWithBytes(data, payload...)
}

// FromString validates and creates a JsonResult from a JSON string and optional payload.
func FromString(jsonStr string, payload ...any) JsonResult {
	return NewJsonResultFromString(jsonStr, payload...)
}

// FromReader streams and validates data from an io.Reader into a JsonResult.
func FromReader(r io.Reader, payload ...any) JsonResult {
	return NewJsonResultFromReader(r, payload...)
}

// FromSerializer creates a JsonResult from an on-demand serializer closure.
func FromSerializer(serializer func() ([]byte, *appfault.AppError), payload ...any) JsonResult {
	return JsonSource.FromSerializer(serializer, payload...)
}

// FromBytesEnvelope converts an existing WrappedBytes envelope into a JsonResult.
func FromBytesEnvelope(wb any) JsonResult {
	return JsonSource.FromBytesEnvelope(wb)
}

// FromError creates a failed JsonResult containing an AppError.
func FromError(appErr *appfault.AppError) JsonResult {
	return NewJsonResultError(appErr)
}

// FromErrorWithPayload creates a failed JsonResult preserving the payload.
func FromErrorWithPayload(appErr *appfault.AppError, payload any) JsonResult {
	return NewJsonResultErrorWithPayload(appErr, payload)
}

// FromAny polymorphically converts any supported source into a JsonResult.
func FromAny(source any) JsonResult {
	if source == nil {
		return JsonResult{
			data:       []byte("null"),
			payload:    nil,
			status:     true,
			statusCode: 200,
		}
	}
	switch v := source.(type) {
	case JsonResult:
		return v
	case *JsonResult:
		if v == nil {
			return JsonResult{
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

// Cast executes a type-safe JSON round-trip cast from Source to Target type, returning a JsonResult with Target payload.
func Cast[Target any](source any) JsonResult {
	data, err := json.Marshal(source)
	if err != nil {
		appErr := appfault.Wrap(errtype.Validation, err, "failed to marshal source for type cast")
		return JsonResult{
			payload:    nil,
			status:     false,
			statusCode: 500,
			appError:   appErr,
		}
	}
	var target Target
	err = json.Unmarshal(data, &target)
	if err != nil {
		appErr := appfault.Wrap(errtype.Validation, err, "failed to unmarshal into target type during cast")
		return JsonResult{
			data:       data,
			payload:    nil,
			status:     false,
			statusCode: 400,
			appError:   appErr,
		}
	}
	return JsonResult{
		data:       data,
		payload:    target,
		status:     true,
		statusCode: 200,
	}
}

// CastTo executes a type-safe JSON round-trip cast from Source to Target.
func CastTo[Target any](source any) JsonResult {
	return Cast[Target](source)
}

// UnmarshalAs parses the JsonResult data directly into a Target type.
func UnmarshalAs[Target any](j JsonResult) (Target, *appfault.AppError) {
	var target Target
	if j.appError != nil {
		return target, j.appError
	}
	if len(j.data) == 0 {
		return target, appfault.New(errtype.Validation, "cannot unmarshal empty JSON result")
	}
	err := json.Unmarshal(j.data, &target)
	if err != nil {
		return target, appfault.Wrap(errtype.Validation, err, "failed to unmarshal JSON into target")
	}
	return target, nil
}

// NewJsonResult serializes payload into JSON and initializes a JsonResult envelope.
func NewJsonResult(payload any) JsonResult {
	data, err := json.Marshal(payload)
	if err != nil {
		appErr := appfault.Wrap(errtype.Validation, err, "failed to marshal payload into JSON")
		return JsonResult{
			payload:    payload,
			status:     false,
			statusCode: 500,
			appError:   appErr,
		}
	}
	return JsonResult{
		data:       data,
		payload:    payload,
		status:     true,
		statusCode: 200,
	}
}

// NewJsonResultWithBytes creates a JsonResult from pre-marshaled JSON bytes and optional payload.
func NewJsonResultWithBytes(data []byte, payload ...any) JsonResult {
	var p any = data
	if len(payload) > 0 {
		p = payload[0]
	}
	if !json.Valid(data) {
		appErr := appfault.New(errtype.Validation, "invalid JSON byte sequence provided")
		return JsonResult{
			data:       data,
			payload:    p,
			status:     false,
			statusCode: 400,
			appError:   appErr,
		}
	}
	return JsonResult{
		data:       data,
		payload:    p,
		status:     true,
		statusCode: 200,
	}
}

// NewJsonResultFromString creates a JsonResult from a JSON string.
func NewJsonResultFromString(jsonStr string, payload ...any) JsonResult {
	var p any = jsonStr
	if len(payload) > 0 {
		p = payload[0]
	}
	return NewJsonResultWithBytes([]byte(jsonStr), p)
}

// NewJsonResultFromReader streams JSON bytes from an io.Reader.
func NewJsonResultFromReader(r io.Reader, payload ...any) JsonResult {
	var p any
	if len(payload) > 0 {
		p = payload[0]
	}
	if r == nil {
		appErr := appfault.New(errtype.Validation, "reader cannot be nil")
		return JsonResult{
			payload:    p,
			status:     false,
			statusCode: 400,
			appError:   appErr,
		}
	}
	data, err := io.ReadAll(r)
	if err != nil {
		appErr := appfault.Wrap(errtype.IO, err, "failed to read stream data")
		return JsonResult{
			payload:    p,
			status:     false,
			statusCode: 500,
			appError:   appErr,
		}
	}
	return NewJsonResultWithBytes(data, p)
}

// NewJsonResultWithStatus creates a JsonResult with explicit status flag and code.
func NewJsonResultWithStatus(data []byte, payload any, status bool, code int) JsonResult {
	return JsonResult{
		data:       data,
		payload:    payload,
		status:     status,
		statusCode: code,
	}
}

// NewJsonResultError creates a failed JsonResult with an AppError.
func NewJsonResultError(appErr *appfault.AppError) JsonResult {
	code := 500
	if appErr != nil {
		if appErr.StatusCode() != 0 {
			code = appErr.StatusCode()
		}
	}
	return JsonResult{
		status:     false,
		statusCode: code,
		appError:   appErr,
	}
}

// NewJsonResultErrorWithPayload creates a failed JsonResult preserving the payload.
func NewJsonResultErrorWithPayload(appErr *appfault.AppError, payload any) JsonResult {
	code := 500
	if appErr != nil {
		if appErr.StatusCode() != 0 {
			code = appErr.StatusCode()
		}
	}
	return JsonResult{
		payload:    payload,
		status:     false,
		statusCode: code,
		appError:   appErr,
	}
}

// Backwards-compatible aliases for PascalCase JSON
func NewJSONResult(payload any) JsonResult {
	return NewJsonResult(payload)
}
func NewJSONResultWithBytes(data []byte, payload ...any) JsonResult {
	return NewJsonResultWithBytes(data, payload...)
}
func NewJSONResultFromString(jsonStr string, payload ...any) JsonResult {
	return NewJsonResultFromString(jsonStr, payload...)
}
func NewJSONResultFromReader(r io.Reader, payload ...any) JsonResult {
	return NewJsonResultFromReader(r, payload...)
}
func NewJSONResultError(appErr *appfault.AppError) JsonResult {
	return NewJsonResultError(appErr)
}
func NewJSONResultErrorWithPayload(appErr *appfault.AppError, payload any) JsonResult {
	return NewJsonResultErrorWithPayload(appErr, payload)
}
func NewJSONResultWithStatus(data []byte, payload any, status bool, code int) JsonResult {
	return NewJsonResultWithStatus(data, payload, status, code)
}

// Raw returns the underlying JSON byte slice.
func (j JsonResult) Raw() []byte {
	return j.data
}

// Bytes returns the underlying JSON byte slice (alias to Raw).
func (j JsonResult) Bytes() []byte {
	return j.data
}

// String returns the JSON string representation.
func (j JsonResult) String() string {
	return string(j.data)
}

// Len returns the byte length of the JSON.
func (j JsonResult) Len() int {
	return len(j.data)
}

// IsEmpty returns true if the JSON bytes are empty.
func (j JsonResult) IsEmpty() bool {
	return len(j.data) == 0
}

// Payload returns the underlying payload.
func (j JsonResult) Payload() any {
	return j.payload
}

// Value returns the underlying payload (alias to Payload).
func (j JsonResult) Value() any {
	return j.payload
}

// AppError returns the underlying *appfault.AppError.
func (j JsonResult) AppError() *appfault.AppError {
	return j.appError
}

// Fault returns the underlying *appfault.AppError (alias to AppError).
func (j JsonResult) Fault() *appfault.AppError {
	return j.appError
}

// Error returns the underlying *appfault.AppError (alias to AppError).
func (j JsonResult) Error() *appfault.AppError {
	return j.appError
}

// HasError returns true if an AppError is present.
func (j JsonResult) HasError() bool {
	return j.appError != nil
}

// IsValid returns true if no AppError is present.
func (j JsonResult) IsValid() bool {
	return j.appError == nil
}

// IsSuccess returns true if status flag is true and no AppError is present.
func (j JsonResult) IsSuccess() bool {
	if j.appError != nil {
		return false
	}
	return j.status
}

// Status returns the boolean status flag.
func (j JsonResult) Status() bool {
	return j.status
}

// StatusCode returns the numeric status code.
func (j JsonResult) StatusCode() int {
	return j.statusCode
}

// Unwrap returns both the JSON byte slice and the AppError.
func (j JsonResult) Unwrap() ([]byte, *appfault.AppError) {
	return j.data, j.appError
}

// Pretty returns formatted and indented JSON.
func (j JsonResult) Pretty() string {
	str, _ := j.PrettyOrError()
	return str
}

// PrettyOrError returns formatted JSON or an AppError if indentation fails.
func (j JsonResult) PrettyOrError() (string, *appfault.AppError) {
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
func (j JsonResult) Compact() string {
	str, _ := j.CompactOrError()
	return str
}

// CompactOrError returns minified JSON or an AppError if compaction fails.
func (j JsonResult) CompactOrError() (string, *appfault.AppError) {
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
func (j JsonResult) Unmarshal(dest any) *appfault.AppError {
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

// ToBytes converts the JsonResult to a standard Bytes[any] envelope.
func (j JsonResult) ToBytes() Bytes[any] {
	return Bytes[any]{
		data:       j.data,
		payload:    j.payload,
		status:     j.status,
		statusCode: j.statusCode,
		appError:   j.appError,
	}
}

var _ WrappedBytes[any] = JsonResult{}
var _ WrappedJson = JsonResult{}
