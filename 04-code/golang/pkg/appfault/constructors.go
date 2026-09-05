package appfault

import "coding-guidelines/common/pkg/errtype"

// New creates an AppError for a given error type variation and message.
// If errType is errtype.None, it returns nil (no error allocated).
func New(errType errtype.Variation, message string) *AppError {
	if errType == errtype.None {
		return nil
	}

	return NewWithContext(errType, message, nil)
}

// NewType creates an AppError using default type name as message.
func NewType(errType errtype.Variation) *AppError {
	if errType == errtype.None {
		return nil
	}

	return New(errType, errType.Name())
}

// createAppErrorInstance constructs the AppError capturing stack and caller objects.
func createAppErrorInstance(errType errtype.Variation, message string) *AppError {
	return &AppError{
		errType: errType,
		message: message,
		caller:  CaptureCallerInfo(3),
		stack:   CaptureStackTrace(3),
		ctx:     NewContextMap(),
	}
}

// NewWithContext constructs an AppError with an initial context map.
func NewWithContext(errType errtype.Variation, message string, ctx map[string]any) *AppError {
	if errType == errtype.None {
		return nil
	}

	e := createAppErrorInstance(errType, message)
	e.ctx = ensureContextMap(ctx)

	return e
}

// Wrap wraps an existing cause with an explicit errtype and custom message.
// If cause is nil or errType is None, it returns nil (no allocation).
func Wrap(errType errtype.Variation, cause error, message string) *AppError {
	if cause == nil || errType == errtype.None {
		return nil
	}

	e := New(errType, message)
	if e == nil {
		return nil
	}

	e.cause = cause

	return e
}

// WrapType wraps an existing cause using cause.Error() as message.
func WrapType(errType errtype.Variation, cause error) *AppError {
	if cause == nil || errType == errtype.None {
		return nil
	}

	return Wrap(errType, cause, cause.Error())
}

// ensureContextMap safely converts a map[string]any to ContextMap.
func ensureContextMap(ctx map[string]any) ContextMap {
	if ctx == nil {
		return NewContextMap()
	}

	cm := NewContextMapWithCapacity(len(ctx))
	for k, v := range ctx {
		cm[k] = v
	}

	return cm
}
