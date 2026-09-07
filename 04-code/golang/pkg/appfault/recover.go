package appfault

import (
	"fmt"

	"coding-guidelines/common/pkg/errtype"
)

// Recover is a defer-friendly utility that captures a panic and transforms it into a robust AppError.
// It builds a full stack trace specifically tracking the panic origin and invokes the handler.
//
// Usage:
//
//	defer appfault.Recover(func(err *appfault.AppError) {
//		// log error, send metrics, etc.
//	})
func Recover(handler func(*AppError)) {
	if r := recover(); r != nil {
		var cause error
		if err, ok := r.(error); ok {
			cause = err
		} else {
			cause = fmt.Errorf("%v", r)
		}

		builder := NewAppBuilder(errtype.Execution, "unhandled panic recovered")
		builder.SetCause(cause)
		builder.SetContext("panic_value", fmt.Sprintf("%v", r))

		// Set the caller information to track back where the panic actually occurred,
		// skipping this recover function and the runtime panic handlers.
		builder.overrideStackTrace(CaptureStackTrace(3))

		handler(builder.Build())
	}
}

// overrideStackTrace is a private method on AppErrorBuilder to forcefully set the stack trace.
func (b *AppErrorBuilder) overrideStackTrace(trace StackTrace) *AppErrorBuilder {
	b.stack = trace
	return b
}
