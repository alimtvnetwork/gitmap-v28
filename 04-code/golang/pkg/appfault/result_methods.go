package appfault

import "fmt"

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

// IsEmpty returns true if no active error is present (or error is zero/empty).
func (r Result[T]) IsEmpty() bool {
	if r.AppError == nil {
		return true
	}

	return r.AppError.IsEmpty()
}

// HasZero returns true if error is nil or represents a zero-value/None error state.
func (r Result[T]) HasZero() bool {
	return r.IsEmpty()
}

// IsZero returns true if error is nil or represents a zero-value/None error state.
func (r Result[T]) IsZero() bool {
	return r.IsEmpty()
}

// IsNull returns true if the embedded AppError pointer is nil.
func (r Result[T]) IsNull() bool {
	return r.AppError == nil
}

// HasNull returns true if the embedded AppError is nil or represents no error.
func (r Result[T]) HasNull() bool {
	return r.IsEmpty()
}

// Clone returns a deep copy of Result with its AppError safely cloned.
func (r Result[T]) Clone() Result[T] {
	return Result[T]{
		Value:    r.Value,
		AppError: r.AppError.Clone(),
	}
}

// Concat combines errors from two results into an immutable Result.
func (r Result[T]) Concat(other Result[T]) Result[T] {
	mergedErr := Merge(r.AppError, other.AppError)
	val := other.Value
	if r.IsSuccess() && other.IsFailed() {
		val = r.Value
	}

	return Result[T]{
		Value:    val,
		AppError: mergedErr,
	}
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

// DefaultResultFormatter formats Result[T]: error banner if failed, or data value if success.
func DefaultResultFormatter[T any](r Result[T]) string {
	if r.IsFailed() {
		return r.AppError.Format(DefaultFaultFormatter)
	}

	return fmt.Sprintf("✅ [OK] %v", r.Value)
}

// Print outputs the result representation to standard output.
func (r Result[T]) Print() {
	fmt.Println(r.Format(nil))
}

// PrintFault outputs the fault representation to standard output if failed.
func (r Result[T]) PrintFault() {
	if r.IsFailed() {
		r.AppError.Print()
	}
}

// Format formats the Result using a custom or default formatter.
func (r Result[T]) Format(formatter ResultFormatter[T]) string {
	if formatter != nil {
		return formatter(r)
	}

	return DefaultResultFormatter(r)
}

// PrintWith outputs the result using a custom formatter.
func (r Result[T]) PrintWith(formatter ResultFormatter[T]) {
	fmt.Println(r.Format(formatter))
}

// FormatStdout formats the Result: rich error banner if failed, or success message if ok.
func (r Result[T]) FormatStdout() string {
	if r.IsFailed() {
		return r.AppError.FormatStdout()
	}

	return fmt.Sprintf("✅ SUCCESS: %v", r.Value)
}

// FormatJson formats the Result: JSON error if failed, or marshaled value JSON.
func (r Result[T]) FormatJson() string {
	if r.IsFailed() {
		return r.AppError.FormatJson()
	}

	return fmt.Sprintf(`{"success":true,"data":%v}`, r.Value)
}

// FormatJSON is an alias for FormatJson.
func (r Result[T]) FormatJSON() string {
	return r.FormatJson()
}

// FormatTextLog formats the Result: structured log error if failed, or log info if ok.
func (r Result[T]) FormatTextLog() string {
	if r.IsFailed() {
		return r.AppError.FormatTextLog()
	}

	return fmt.Sprintf("[INFO] status=200 msg=%q", fmt.Sprintf("%v", r.Value))
}

// PrintStdout prints the result formatted for stdout.
func (r Result[T]) PrintStdout() {
	fmt.Println(r.FormatStdout())
}

// PrintJson prints the result formatted as JSON.
func (r Result[T]) PrintJson() {
	fmt.Println(r.FormatJson())
}

// PrintJSON is an alias for PrintJson.
func (r Result[T]) PrintJSON() {
	r.PrintJson()
}

// PrintLog prints the result formatted as a structured log line.
func (r Result[T]) PrintLog() {
	fmt.Println(r.FormatTextLog())
}
