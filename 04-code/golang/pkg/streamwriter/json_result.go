package streamwriter

import (
	"bytes"
	"encoding/json"
	"io"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
)

// WrappedJson defines the contract for JSON result envelopes with formatting, unmarshaling, and conversion.
type WrappedJson interface {
	Raw() []byte
	Bytes() []byte
	String() string
	Len() int
	IsEmpty() bool
	IsNull() bool
	HasZero() bool
	IsZero() bool
	HasNull() bool
	Payload() any
	Value() any
	AppError() *appfault.AppError
	Fault() *appfault.AppError
	Error() *appfault.AppError
	HasError() bool
	IsValid() bool
	IsSuccess() bool
	Status() bool
	StatusCode() int
	Unwrap() ([]byte, *appfault.AppError)
	Pretty() string
	PrettyOrError() (string, *appfault.AppError)
	Compact() string
	CompactOrError() (string, *appfault.AppError)
	Unmarshal(dest any) *appfault.AppError
	ToBytes() Bytes[any]
}

// WrappedJSON is an alias for WrappedJson for backwards compatibility.
type WrappedJSON = WrappedJson

// JsonResult encapsulates JSON serialized data and AppError state.
// Status and validity are dynamically computed from appError without redundant fields.
type JsonResult struct {
	data     []byte
	appError *appfault.AppError
}

// JSONResult is an alias for JsonResult for backwards compatibility.
type JSONResult = JsonResult

// JsonPayloadResult extends JsonResult by embedding it and attaching a strongly-typed payload T.
type JsonPayloadResult[T any] struct {
	JsonResult
	payload T
}

// JsonResultWithPayload is an alias for JsonPayloadResult.
type JsonResultWithPayload[T any] = JsonPayloadResult[T]

// Payload returns the strongly-typed payload.
func (p JsonPayloadResult[T]) Payload() T {
	return p.payload
}

// Value returns the strongly-typed payload.
func (p JsonPayloadResult[T]) Value() T {
	return p.payload
}

// ToBytes converts to a typed Bytes[T] envelope.
func (p JsonPayloadResult[T]) ToBytes() Bytes[T] {
	if p.appError != nil {
		return NewBytesErrorWithPayload(p.appError, p.payload)
	}

	return NewBytes(p.data, p.payload)
}

// Clone returns a deep copy of JsonPayloadResult[T].
func (p JsonPayloadResult[T]) Clone() JsonPayloadResult[T] {
	return JsonPayloadResult[T]{
		JsonResult: p.JsonResult.Clone(),
		payload:    p.payload,
	}
}

// Concat safely combines two JsonPayloadResult[T] instances without panic.
// If the receiver is empty or null, it returns other.Clone().
// If other is empty or null, it returns p.Clone().
func (p JsonPayloadResult[T]) Concat(other JsonPayloadResult[T]) JsonPayloadResult[T] {
	if p.IsNull() || p.IsEmpty() {
		return other.Clone()
	}

	if other.IsNull() || other.IsEmpty() {
		return p.Clone()
	}

	return JsonPayloadResult[T]{
		JsonResult: p.JsonResult.Concat(other.JsonResult),
		payload:    other.payload,
	}
}

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
func (jsonSourceSingleton) FromBytes(data []byte) JsonResult {
	return NewJsonResultWithBytes(data)
}

// FromString validates and wraps a JSON string into a JsonResult.
func (jsonSourceSingleton) FromString(jsonStr string) JsonResult {
	return NewJsonResultFromString(jsonStr)
}

// FromReader streams and validates data from an io.Reader into a JsonResult.
func (jsonSourceSingleton) FromReader(r io.Reader) JsonResult {
	return NewJsonResultFromReader(r)
}

// FromSerializer executes a lazy closure and wraps the resulting JSON bytes into a JsonResult.
func (jsonSourceSingleton) FromSerializer(serializer func() ([]byte, *appfault.AppError)) JsonResult {
	if serializer == nil {
		return JsonResult{
			appError: appfault.New(errtype.Validation, "serializer closure cannot be nil"),
		}
	}

	data, appErr := serializer()
	if appErr != nil {
		return JsonResult{
			data:     data,
			appError: appErr,
		}
	}

	return NewJsonResultWithBytes(data)
}

// FromBytesEnvelope converts an existing WrappedBytes envelope into a JsonResult.
func (jsonSourceSingleton) FromBytesEnvelope(wb any) JsonResult {
	if wb == nil {
		return JsonResult{
			appError: appfault.New(errtype.Validation, "wrapped bytes envelope cannot be nil"),
		}
	}

	if rawProvider, ok := wb.(interface{ Raw() []byte }); ok {
		raw := rawProvider.Raw()
		if errProvider, ok := wb.(interface{ AppError() *appfault.AppError }); ok {
			if appErr := errProvider.AppError(); appErr != nil {
				return JsonResult{
					data:     raw,
					appError: appErr,
				}
			}
		}

		return NewJsonResultWithBytes(raw)
	}

	return FromAny(wb)
}

// FromError creates a failed JsonResult containing an AppError.
func (jsonSourceSingleton) FromError(appErr *appfault.AppError) JsonResult {
	return NewJsonResultError(appErr)
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

// typedJsonSource provides typed factory methods producing extended JsonPayloadResult[T] instances.
type typedJsonSource[T any] struct{}

// JsonSourceOf creates a scoped factory for producing extended JsonPayloadResult instances with payload T.
func JsonSourceOf[T any]() typedJsonSource[T] {
	return typedJsonSource[T]{}
}

// JSONSourceOf is an alias for JsonSourceOf for backwards compatibility.
func JSONSourceOf[T any]() typedJsonSource[T] {
	return JsonSourceOf[T]()
}

// FromPayload serializes payload T into a JsonPayloadResult[T].
func (typedJsonSource[T]) FromPayload(payload T) JsonPayloadResult[T] {
	return WithPayload(NewJsonResult(payload), payload)
}

// FromBytes validates and creates a JsonPayloadResult[T] from a byte slice and payload T.
func (typedJsonSource[T]) FromBytes(data []byte, payload T) JsonPayloadResult[T] {
	return WithPayload(NewJsonResultWithBytes(data), payload)
}

// FromString validates and creates a JsonPayloadResult[T] from a string and payload T.
func (typedJsonSource[T]) FromString(jsonStr string, payload T) JsonPayloadResult[T] {
	return WithPayload(NewJsonResultFromString(jsonStr), payload)
}

// FromReader streams and validates data from an io.Reader into a JsonPayloadResult[T].
func (typedJsonSource[T]) FromReader(r io.Reader, payload T) JsonPayloadResult[T] {
	return WithPayload(NewJsonResultFromReader(r), payload)
}

// FromSerializer creates a JsonPayloadResult[T] from a serializer closure and payload T.
func (typedJsonSource[T]) FromSerializer(serializer func() ([]byte, *appfault.AppError), payload T) JsonPayloadResult[T] {
	return WithPayload(JsonSource.FromSerializer(serializer), payload)
}

// FromBytesEnvelope converts an existing WrappedBytes envelope into a JsonPayloadResult[T].
func (typedJsonSource[T]) FromBytesEnvelope(wb any, payload T) JsonPayloadResult[T] {
	return WithPayload(FromBytesEnvelope(wb), payload)
}

// FromError creates a failed JsonPayloadResult[T] with an AppError and payload T.
func (typedJsonSource[T]) FromError(appErr *appfault.AppError, payload ...T) JsonPayloadResult[T] {
	var p T
	if len(payload) > 0 {
		p = payload[0]
	}

	return WithPayload(NewJsonResultError(appErr), p)
}

// FromPayload serializes any payload into a JsonResult.
func FromPayload(payload any) JsonResult {
	return NewJsonResult(payload)
}

// FromBytes validates and creates a JsonResult from a pre-existing byte slice.
func FromBytes(data []byte) JsonResult {
	return NewJsonResultWithBytes(data)
}

// FromString validates and creates a JsonResult from a JSON string.
func FromString(jsonStr string) JsonResult {
	return NewJsonResultFromString(jsonStr)
}

// FromReader streams and validates data from an io.Reader into a JsonResult.
func FromReader(r io.Reader) JsonResult {
	return NewJsonResultFromReader(r)
}

// FromSerializer creates a JsonResult from an on-demand serializer closure.
func FromSerializer(serializer func() ([]byte, *appfault.AppError)) JsonResult {
	return JsonSource.FromSerializer(serializer)
}

// FromBytesEnvelope converts an existing WrappedBytes envelope into a JsonResult.
func FromBytesEnvelope(wb any) JsonResult {
	return JsonSource.FromBytesEnvelope(wb)
}

// FromError creates a failed JsonResult containing an AppError.
func FromError(appErr *appfault.AppError) JsonResult {
	return NewJsonResultError(appErr)
}

// FromAny polymorphically converts any supported source into a JsonResult.
func FromAny(source any) JsonResult {
	if source == nil {
		return JsonResult{
			data: []byte("null"),
		}
	}

	switch v := source.(type) {
	case JsonResult:
		return v
	case *JsonResult:
		if v == nil {
			return JsonResult{
				appError: appfault.New(errtype.Validation, "nil JsonResult pointer provided"),
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

// WithPayload attaches a typed payload T to a JsonResult, returning an extended JsonPayloadResult[T].
func WithPayload[T any](res JsonResult, payload T) JsonPayloadResult[T] {
	return JsonPayloadResult[T]{
		JsonResult: res,
		payload:    payload,
	}
}

// WithPayload attaches an untyped payload to a JsonResult, returning a JsonPayloadResult[any].
func (j JsonResult) WithPayload(payload any) JsonPayloadResult[any] {
	return JsonPayloadResult[any]{
		JsonResult: j,
		payload:    payload,
	}
}

// Cast executes a type-safe JSON round-trip cast from Source to Target type, returning a JsonResult.
func Cast[Target any](source any) JsonResult {
	data, err := json.Marshal(source)
	if err != nil {
		return JsonResult{
			appError: appfault.Wrap(errtype.Validation, err, "failed to marshal source for type cast"),
		}
	}

	var target Target
	err = json.Unmarshal(data, &target)
	if err != nil {
		return JsonResult{
			data:     data,
			appError: appfault.Wrap(errtype.Validation, err, "failed to unmarshal into target type during cast"),
		}
	}

	return JsonResult{
		data: data,
	}
}

// CastWithPayload executes a type-safe cast and returns a JsonPayloadResult[Target] with typed payload.
func CastWithPayload[Target any](source any) JsonPayloadResult[Target] {
	res := Cast[Target](source)
	var target Target
	if res.IsValid() {
		_ = json.Unmarshal(res.data, &target)
	}

	return JsonPayloadResult[Target]{
		JsonResult: res,
		payload:    target,
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
func NewJsonResult(val any) JsonResult {
	data, err := json.Marshal(val)
	if err != nil {
		return JsonResult{
			appError: appfault.Wrap(errtype.Validation, err, "failed to marshal payload into JSON"),
		}
	}

	return JsonResult{
		data: data,
	}
}

// NewJsonResultWithBytes creates a JsonResult from pre-marshaled JSON bytes.
func NewJsonResultWithBytes(data []byte) JsonResult {
	if !json.Valid(data) {
		return JsonResult{
			data:     data,
			appError: appfault.New(errtype.Validation, "invalid JSON byte sequence provided"),
		}
	}

	return JsonResult{
		data: data,
	}
}

// NewJsonResultFromString creates a JsonResult from a JSON string.
func NewJsonResultFromString(jsonStr string) JsonResult {
	return NewJsonResultWithBytes([]byte(jsonStr))
}

// NewJsonResultFromReader streams JSON bytes from an io.Reader.
func NewJsonResultFromReader(r io.Reader) JsonResult {
	if r == nil {
		return JsonResult{
			appError: appfault.New(errtype.Validation, "reader cannot be nil"),
		}
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return JsonResult{
			appError: appfault.Wrap(errtype.IO, err, "failed to read stream data"),
		}
	}

	return NewJsonResultWithBytes(data)
}

// NewJsonResultError creates a failed JsonResult with an AppError.
func NewJsonResultError(appErr *appfault.AppError) JsonResult {
	return JsonResult{
		appError: appErr,
	}
}

// Backwards-compatible aliases for PascalCase JSON
func NewJSONResult(val any) JsonResult {
	return NewJsonResult(val)
}

func NewJSONResultWithBytes(data []byte) JsonResult {
	return NewJsonResultWithBytes(data)
}

func NewJSONResultFromString(jsonStr string) JsonResult {
	return NewJsonResultFromString(jsonStr)
}

func NewJSONResultFromReader(r io.Reader) JsonResult {
	return NewJsonResultFromReader(r)
}

func NewJSONResultError(appErr *appfault.AppError) JsonResult {
	return NewJsonResultError(appErr)
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

// IsNull returns true if JSON data is empty and no AppError is present.
func (j JsonResult) IsNull() bool {
	return len(j.data) == 0 && j.appError == nil
}

// HasZero returns true if JSON data is empty and no AppError is present.
func (j JsonResult) HasZero() bool {
	return j.IsEmpty()
}

// IsZero returns true if JSON data is empty and no AppError is present.
func (j JsonResult) IsZero() bool {
	return j.IsEmpty()
}

// HasNull returns true if appError is nil or represents no error.
func (j JsonResult) HasNull() bool {
	if j.appError == nil {
		return true
	}

	return j.appError.HasNullError()
}

// Clone returns a deep copy of JsonResult.
func (j JsonResult) Clone() JsonResult {
	var copied []byte
	if j.data != nil {
		copied = make([]byte, len(j.data))
		copy(copied, j.data)
	}

	return JsonResult{
		data:     copied,
		appError: j.appError.Clone(),
	}
}

// Concat safely merges two JsonResult envelopes without panic.
// If the receiver is empty or null, it returns other.Clone().
// If other is empty or null, it returns j.Clone().
func (j JsonResult) Concat(other JsonResult) JsonResult {
	if j.IsNull() || j.IsEmpty() {
		return other.Clone()
	}

	if other.IsNull() || other.IsEmpty() {
		return j.Clone()
	}

	mergedErr := appfault.Merge(j.appError, other.appError)
	mergedData := make([]byte, len(j.data)+len(other.data))
	copy(mergedData, j.data)
	copy(mergedData[len(j.data):], other.data)

	return JsonResult{
		data:     mergedData,
		appError: mergedErr,
	}
}

// Payload returns the data slice for WrappedBytes compatibility.
func (j JsonResult) Payload() any {
	return j.data
}

// Value returns the data slice for WrappedBytes compatibility.
func (j JsonResult) Value() any {
	return j.data
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

// IsValid dynamically evaluates true if no AppError is present.
func (j JsonResult) IsValid() bool {
	return j.appError == nil
}

// IsSuccess dynamically evaluates true if no AppError is present.
func (j JsonResult) IsSuccess() bool {
	return j.appError == nil
}

// Status dynamically evaluates true if no AppError is present.
func (j JsonResult) Status() bool {
	return j.appError == nil
}

// StatusCode returns the numeric status code derived from appError or 200.
func (j JsonResult) StatusCode() int {
	if j.appError != nil {
		if j.appError.StatusCode() != 0 {
			return j.appError.StatusCode()
		}

		switch j.appError.GetType() {
		case errtype.Validation, errtype.Precondition:
			return 400
		case errtype.Unauthorized:
			return 401
		case errtype.Forbidden:
			return 403
		case errtype.NotFound:
			return 404
		case errtype.Timeout:
			return 408
		default:
			return 500
		}
	}

	return 200
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
	if j.appError != nil {
		return NewBytesErrorWithPayload(j.appError, any(j.data))
	}

	return NewBytes(j.data, any(j.data))
}

var _ WrappedBytes[any] = JsonResult{}
var _ WrappedJson = JsonResult{}
var _ WrappedBytes[any] = JsonPayloadResult[any]{}
