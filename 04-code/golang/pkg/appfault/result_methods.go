package appfault

// IsSuccess returns true if no error is present.
func (r Result[T]) IsSuccess() bool {
	return r.AppError == nil
}

// IsFailed returns true if an error is present.
func (r Result[T]) IsFailed() bool {
	return r.AppError != nil
}

// IsFailure is an alias for IsFailed.
func (r Result[T]) IsFailure() bool {
	return r.IsFailed()
}

// IsInvalid is an alias for IsFailed.
func (r Result[T]) IsInvalid() bool {
	return r.IsFailed()
}

// IsValid returns true if the operation succeeded with no error.
func (r Result[T]) IsValid() bool {
	return r.IsSuccess()
}

// HasError returns true if an error is present.
func (r Result[T]) HasError() bool {
	return r.IsFailed()
}

// HasNoError returns true if no error exists.
func (r Result[T]) HasNoError() bool {
	return r.IsSuccess()
}

// HasValidError returns true if the embedded error has a valid code.
func (r Result[T]) HasValidError() bool {
	if r.AppError == nil {
		return false
	}

	return r.AppError.HasValidError()
}

// IsSafe returns true if the operation succeeded with no error.
func (r Result[T]) IsSafe() bool {
	return r.IsSuccess()
}

// Unwrap unpacks the (Value, *AppError) tuple.
func (r Result[T]) Unwrap() (T, *AppError) {
	return r.Value, r.AppError
}

// UnwrapOr returns the inner value if successful, or defaultVal on failure.
func (r Result[T]) UnwrapOr(defaultVal T) T {
	if r.IsFailed() {
		return defaultVal
	}

	return r.Value
}
