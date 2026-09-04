package appfault

import "coding-guidelines/common/pkg/errtype"

// GetMessage returns the human-readable diagnostic message.
func (e *AppError) GetMessage() string {
	if e == nil {
		return ""
	}

	return e.message
}

// Message is an alias for GetMessage.
func (e *AppError) Message() string {
	return e.GetMessage()
}

// GetStatusCode returns the attached HTTP status code or 0.
func (e *AppError) GetStatusCode() int {
	if e == nil {
		return 0
	}

	return e.statusCode
}

// StatusCode is an alias for GetStatusCode.
func (e *AppError) StatusCode() int {
	return e.GetStatusCode()
}

// GetType returns the error type variation.
func (e *AppError) GetType() errtype.Variation {
	if e == nil {
		return errtype.None
	}

	return e.errType
}

// Type is an alias for GetType.
func (e *AppError) Type() errtype.Variation {
	return e.GetType()
}

// Caller returns the structured CallerInfo object.
func (e *AppError) Caller() CallerInfo {
	if e == nil {
		return CallerInfo{}
	}

	return e.caller
}

// StackTrace returns the structured call stack frames.
func (e *AppError) StackTrace() StackTrace {
	if e == nil {
		return nil
	}

	return e.stack
}
