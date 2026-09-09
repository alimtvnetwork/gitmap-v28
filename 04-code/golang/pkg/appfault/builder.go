package appfault

import (
	"encoding/json"
	"errors"

	"coding-guidelines/common/pkg/errtype"
)

type (
	AppErrorBuilder struct {
		errType    errtype.Variation
		statusCode int
		message    string
		stack      StackTrace
		ctx        map[string]any
		cause      error
	}

	AppBuilder = AppErrorBuilder
)

// NewBuilder initializes a new mutable AppErrorBuilder.
func NewBuilder(errType errtype.Variation, message string, skipFrames ...int) *AppErrorBuilder {
	skip := 2
	if len(skipFrames) > 0 {
		skip += skipFrames[0]
	}

	return &AppErrorBuilder{
		errType: errType,
		message: message,
		ctx:     make(map[string]any, 0),
		stack:   CaptureStackTrace(skip),
	}
}

// NewAppBuilder initializes a new mutable AppBuilder.
func NewAppBuilder(errType errtype.Variation, message string, skipFrames ...int) *AppBuilder {
	return NewBuilder(errType, message, skipFrames...)
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

// WithContext is a fluent alias for SetContext.
func (b *AppErrorBuilder) WithContext(key string, value any) *AppErrorBuilder {
	return b.SetContext(key, value)
}

// WithCause is a fluent alias for SetCause.
func (b *AppErrorBuilder) WithCause(cause error) *AppErrorBuilder {
	return b.SetCause(cause)
}

// SetStatusCode updates the status code on the builder.
func (b *AppErrorBuilder) SetStatusCode(code int) *AppErrorBuilder {
	b.statusCode = code
	return b
}

// WithStatusCode is a fluent alias for SetStatusCode.
func (b *AppErrorBuilder) WithStatusCode(code int) *AppErrorBuilder {
	return b.SetStatusCode(code)
}

// Build freezes the builder state into a strictly immutable *AppError.
func (b *AppErrorBuilder) Build() *AppError {
	if b.errType == errtype.None {
		return nil
	}

	ctxMap := NewContextMap()
	if len(b.ctx) > 0 {
		for k, v := range b.ctx {
			ctxMap = ctxMap.Set(k, v)
		}
	} else {
		ctxMap = nil
	}

	return &AppError{
		errType:    b.errType,
		statusCode: b.statusCode,
		message:    b.message,
		stack:      b.stack,
		ctx:        ctxMap,
		cause:      b.cause,
	}
}

// ToDataModel converts builder state to a serializable AppErrorDataModel.
func (b *AppErrorBuilder) ToDataModel() AppErrorDataModel {
	causeStr := ""
	if b.cause != nil {
		causeStr = b.cause.Error()
	}

	ctxMap := NewContextMap()
	if len(b.ctx) > 0 {
		for k, v := range b.ctx {
			ctxMap = ctxMap.Set(k, v)
		}
	} else {
		ctxMap = nil
	}

	return AppErrorDataModel{
		Type:    b.errType,
		Message: b.message,
		Stack:   b.stack,
		Ctx:     ctxMap,
		Cause:   causeStr,
	}
}

// FromDataModel populates builder state from an AppErrorDataModel.
func (b *AppErrorBuilder) FromDataModel(model AppErrorDataModel) *AppErrorBuilder {
	b.errType = model.Type
	b.message = model.Message
	b.stack = model.Stack

	b.ctx = make(map[string]any, 0)
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
		errType: e.errType,
		message: e.message,
		stack:   e.stack,
		ctx:     ctxCopy,
		cause:   e.cause,
	}
}
