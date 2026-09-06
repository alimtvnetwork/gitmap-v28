package appfault

import (
	"encoding/json"
	"errors"

	"coding-guidelines/common/pkg/errtype"
)

type (
	AppErrorBuilder struct {
		errType    errtype.Variation
		message    string
		caller     CallerInfo
		stack      StackTrace
		ctx        map[string]any
		cause      error
		statusCode int
	}

	AppBuilder = AppErrorBuilder
)

// NewBuilder initializes a new mutable AppErrorBuilder.
func NewBuilder(errType errtype.Variation, message string) *AppErrorBuilder {
	return &AppErrorBuilder{
		errType: errType,
		message: message,
		ctx:     make(map[string]any),
		caller:  CaptureCallerInfo(2),
		stack:   CaptureStackTrace(2),
	}
}

// NewAppBuilder initializes a new mutable AppBuilder.
func NewAppBuilder(errType errtype.Variation, message string) *AppBuilder {
	return NewBuilder(errType, message)
}

// SetType updates the error type on the builder.
func (b *AppErrorBuilder) SetType(errType errtype.Variation) *AppErrorBuilder {
	b.errType = errType

	return b
}

// SetMessage updates the diagnostic message on the builder.
func (b *AppErrorBuilder) SetMessage(message string) *AppErrorBuilder {
	b.message = message

	return b
}

// SetStatusCode sets the HTTP status code on the builder.
func (b *AppErrorBuilder) SetStatusCode(code int) *AppErrorBuilder {
	b.statusCode = code

	return b
}

// SetCaller sets the caller info on the builder.
func (b *AppErrorBuilder) SetCaller(caller CallerInfo) *AppErrorBuilder {
	b.caller = caller

	return b
}

// SetContext sets a context key-value pair on the builder.
func (b *AppErrorBuilder) SetContext(key string, value any) *AppErrorBuilder {
	b.ctx[key] = value

	return b
}

// SetCause sets the underlying root cause error on the builder.
func (b *AppErrorBuilder) SetCause(cause error) *AppErrorBuilder {
	b.cause = cause

	return b
}

// WithStatusCode is a fluent alias for SetStatusCode.
func (b *AppErrorBuilder) WithStatusCode(code int) *AppErrorBuilder {
	return b.SetStatusCode(code)
}

// WithCaller is a fluent alias for SetCaller.
func (b *AppErrorBuilder) WithCaller(caller CallerInfo) *AppErrorBuilder {
	return b.SetCaller(caller)
}

// WithContext is a fluent alias for SetContext.
func (b *AppErrorBuilder) WithContext(key string, value any) *AppErrorBuilder {
	return b.SetContext(key, value)
}

// WithCause is a fluent alias for SetCause.
func (b *AppErrorBuilder) WithCause(cause error) *AppErrorBuilder {
	return b.SetCause(cause)
}

// Build freezes the builder state into a strictly immutable *AppError.
func (b *AppErrorBuilder) Build() *AppError {
	if b.errType == errtype.None {
		return nil
	}

	ctxMap := NewContextMap()
	for k, v := range b.ctx {
		ctxMap.Set(k, v)
	}

	if b.statusCode > 0 {
		ctxMap.Set("StatusCode", b.statusCode)
	}

	return &AppError{
		errType:    b.errType,
		message:    b.message,
		caller:     b.caller,
		stack:      b.stack,
		ctx:        ctxMap,
		cause:      b.cause,
		statusCode: b.statusCode,
	}
}

// ToDataModel converts builder state to a serializable AppErrorDataModel.
func (b *AppErrorBuilder) ToDataModel() AppErrorDataModel {
	causeStr := ""
	if b.cause != nil {
		causeStr = b.cause.Error()
	}

	ctxMap := NewContextMap()
	for k, v := range b.ctx {
		ctxMap.Set(k, v)
	}

	return AppErrorDataModel{
		Type:       b.errType,
		Message:    b.message,
		Caller:     b.caller,
		Stack:      b.stack,
		Ctx:        ctxMap,
		Cause:      causeStr,
		StatusCode: b.statusCode,
	}
}

// FromDataModel populates builder state from an AppErrorDataModel.
func (b *AppErrorBuilder) FromDataModel(model AppErrorDataModel) *AppErrorBuilder {
	b.errType = model.Type
	b.message = model.Message
	b.caller = model.Caller
	b.stack = model.Stack
	b.statusCode = model.StatusCode

	b.ctx = make(map[string]any)
	for k, v := range model.Ctx {
		b.ctx[k] = v
	}

	if len(model.Cause) > 0 {
		b.cause = errors.New(model.Cause)
	}

	return b
}

// MarshalJSON provides direct marshaling effect on the builder.
func (b *AppErrorBuilder) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.ToDataModel())
}

// UnmarshalJSON provides direct unmarshaling effect onto the builder.
func (b *AppErrorBuilder) UnmarshalJSON(data []byte) error {
	var model AppErrorDataModel
	if err := json.Unmarshal(data, &model); err != nil {
		return err
	}

	b.FromDataModel(model)

	return nil
}

// ToJson exports the builder state as indented JSON bytes.
func (b *AppErrorBuilder) ToJson() ([]byte, error) {
	return json.MarshalIndent(b.ToDataModel(), "", "  ")
}

// ToJSON is an alias for ToJson.
func (b *AppErrorBuilder) ToJSON() ([]byte, error) {
	return b.ToJson()
}

// ToJsonString exports the builder state as a formatted JSON string.
func (b *AppErrorBuilder) ToJsonString() string {
	bytes, err := b.ToJson()
	if err != nil {
		return "{}"
	}

	return string(bytes)
}

// ToJSONString is an alias for ToJsonString.
func (b *AppErrorBuilder) ToJSONString() string {
	return b.ToJsonString()
}

// ToBuilder converts an immutable *AppError back into a mutable *AppErrorBuilder.
// This allows staging modifications before building a new immutable *AppError.
func (e *AppError) ToBuilder() *AppErrorBuilder {
	if e == nil {
		return NewBuilder(errtype.None, "")
	}

	ctxCopy := make(map[string]any)
	if e.ctx != nil {
		for k, v := range e.ctx {
			ctxCopy[k] = v
		}
	}

	return &AppErrorBuilder{
		errType:    e.errType,
		message:    e.message,
		caller:     e.caller,
		stack:      e.stack,
		ctx:        ctxCopy,
		cause:      e.cause,
		statusCode: e.statusCode,
	}
}
