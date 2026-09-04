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

// IsSuccess returns true if no error is present (e is nil or Type is None).
func (e *AppError) IsSuccess() bool {
	return e.HasNullError()
}
