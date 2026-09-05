package appfault

import "encoding/json"

// Result wraps a typed value bundled with monadic error state.
type Result[T any] struct {
	Value    T         `json:"Value,omitempty" yaml:"Value,omitempty"`
	AppError *AppError `json:"AppError,omitempty" yaml:"AppError,omitempty"`
}

// Data returns the underlying Value payload for API envelope compatibility.
func (r Result[T]) Data() T {
	return r.Value
}

// Fault returns the underlying *AppError (alias for AppError).
func (r Result[T]) Fault() *AppError {
	return r.AppError
}

// Error returns the underlying *AppError.
func (r Result[T]) Error() *AppError {
	return r.AppError
}

// AppErrorOrNil returns *AppError if present, otherwise nil.
func (r Result[T]) AppErrorOrNil() *AppError {
	if r.AppError != nil {
		return r.AppError
	}

	return nil
}

// ToJson exports Result as indented JSON bytes.
func (r Result[T]) ToJson() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

// ToJsonString exports Result as a JSON string.
func (r Result[T]) ToJsonString() string {
	b, err := r.ToJson()
	if err != nil {
		return "{}"
	}

	return string(b)
}
