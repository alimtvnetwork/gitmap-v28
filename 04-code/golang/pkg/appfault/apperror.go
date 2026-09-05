package appfault

import "coding-guidelines/common/pkg/errtype"

// AppError is the universal structured error type carrying full diagnostics.
// All internal fields are unexported for strict encapsulation.
type AppError struct {
	errType    errtype.Variation
	message    string
	caller     CallerInfo
	stack      StackTrace
	ctx        ContextMap
	cause      error
	statusCode int
}

// Fault is retained as a type alias for AppError for backward compatibility.
type Fault = AppError

// HasError returns true if the AppError exists and is not errtype.None.
func (e *AppError) HasError() bool {
	if e == nil {
		return false
	}

	return e.errType.HasError()
}

// HasNullError returns true if e is nil or represents no error.
func (e *AppError) HasNullError() bool {
	if e == nil {
		return true
	}

	return e.errType.IsNone()
}

// IsNull returns true if e is nil.
func (e *AppError) IsNull() bool {
	return e == nil
}

// IsEmpty returns true if e is nil or Type is None.
func (e *AppError) IsEmpty() bool {
	return e.HasNullError()
}

// HasZero returns true if e is nil or represents a zero-value/None error state.
func (e *AppError) HasZero() bool {
	return e.HasNullError()
}

// IsZero returns true if e is nil or represents a zero-value/None error state.
func (e *AppError) IsZero() bool {
	return e.HasNullError()
}

// HasNull returns true if e is nil or represents no error.
func (e *AppError) HasNull() bool {
	return e.HasNullError()
}

// IsSuccess returns true if no error is present (e is nil or Type is None).
func (e *AppError) IsSuccess() bool {
	return e.HasNullError()
}

// Clone creates an exported deep copy of the AppError, guaranteeing immutability.
// If the receiver is nil, it safely returns nil without panic.
func (e *AppError) Clone() *AppError {
	return e.clone()
}

// Concat combines the receiver error with another error into an immutable AppError.
// If either is nil or zero-value, it safely returns the active error without panic.
func (e *AppError) Concat(other *AppError) *AppError {
	return Merge(e, other)
}

// clone creates a deep copy of the AppError, guaranteeing immutability across operations.
func (e *AppError) clone() *AppError {
	if e == nil {
		return nil
	}

	var clonedCtx ContextMap
	if e.ctx != nil {
		clonedCtx = e.ctx.Clone()
	} else {
		clonedCtx = NewContextMap()
	}

	return &AppError{
		errType:    e.errType,
		message:    e.message,
		caller:     e.caller,
		stack:      e.stack,
		ctx:        clonedCtx,
		cause:      e.cause,
		statusCode: e.statusCode,
	}
}
