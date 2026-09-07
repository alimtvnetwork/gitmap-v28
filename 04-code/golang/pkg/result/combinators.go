package result

// Map transforms a successful Result[T] into a Result[U] using the provided function.
// If the input Result is a failure, it safely passes the error forward without calling the function.
func Map[T, U any](res Result[T], fn func(T) U) Result[U] {
	if res.IsFailure() {
		return FailureFromWrap[U](res)
	}

	return Success(fn(res.Data()))
}

// FlatMap binds a successful Result[T] to a function that returns a new Result[U].
// This allows chaining multiple operations that can fail, short-circuiting on the first error.
func FlatMap[T, U any](res Result[T], fn func(T) Result[U]) Result[U] {
	if res.IsFailure() {
		return FailureFromWrap[U](res)
	}

	return fn(res.Data())
}

// Tap executes a side-effect function if the Result is successful, then returns the original Result unmodified.
// This is useful for logging, metrics, or auditing without breaking the pipeline.
func Tap[T any](res Result[T], fn func(T)) Result[T] {
	if res.IsSuccess() {
		fn(res.Data())
	}

	return res
}
