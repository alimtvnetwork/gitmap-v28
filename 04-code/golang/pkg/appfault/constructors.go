package appfault

import "coding-guidelines/common/pkg/errtype"

// New creates an AppError for a given error type variation and message.
// If errType is errtype.None, it returns nil (no error allocated).
func New(errType errtype.Variation, message string, skipFrames ...int) *AppError {
	if errType == errtype.None {
		return nil
	}

	skip := 3
	if len(skipFrames) > 0 {
		skip += skipFrames[0]
	}
	return NewWithContext(errType, message, nil, skip-3) // pass relative skip
}

// NewType creates an AppError using default type name as message.
func NewType(errType errtype.Variation, skipFrames ...int) *AppError {
	if errType == errtype.None {
		return nil
	}

	skip := 3
	if len(skipFrames) > 0 {
		skip += skipFrames[0]
	}
	return New(errType, errType.Name(), skip-3)
}

// createAppErrorInstance constructs the AppError capturing stack trace.
func createAppErrorInstance(errType errtype.Variation, message string, skip int) *AppError {
	return &AppError{
		errType: errType,
		message: message,
		stack:   CaptureStackTrace(skip),
		ctx:     nil, // We will not allocate ContextMap yet
	}
}

// NewWithContext constructs an AppError with an initial context map.
func NewWithContext(errType errtype.Variation, message string, ctx map[string]any, skipFrames ...int) *AppError {
	if errType == errtype.None {
		return nil
	}

	skip := 3
	if len(skipFrames) > 0 {
		skip += skipFrames[0]
	}

	e := createAppErrorInstance(errType, message, skip)
	if ctx != nil && len(ctx) > 0 {
		e.ctx = ensureContextMap(ctx)
	}

	return e
}

// Wrap wraps an existing cause with an explicit errtype and custom message.
// If cause is nil or errType is None, it returns nil (no allocation).
func Wrap(errType errtype.Variation, cause error, message string, skipFrames ...int) *AppError {
	if cause == nil || errType == errtype.None {
		return nil
	}

	skip := 3
	if len(skipFrames) > 0 {
		skip += skipFrames[0]
	}

	e := New(errType, message, skip-3)
	if e == nil {
		return nil
	}

	e.cause = cause

	return e
}

// WrapType wraps an existing cause using cause.Error() as message.
func WrapType(errType errtype.Variation, cause error, skipFrames ...int) *AppError {
	if cause == nil || errType == errtype.None {
		return nil
	}

	skip := 3
	if len(skipFrames) > 0 {
		skip += skipFrames[0]
	}

	return Wrap(errType, cause, cause.Error(), skip-3)
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
