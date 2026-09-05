package appwriter

import (
	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
	"coding-guidelines/common/pkg/result"
)

// WrapWriterSuccess wraps an initialized BaseWriter in a successful BaseWriterWrap.
func WrapWriterSuccess(w *BaseWriter) BaseWriterWrap {
	return result.WrapSuccess(w)
}

// WrapWriterFailure wraps an existing AppError object into a failed BaseWriterWrap.
func WrapWriterFailure(err *appfault.AppError) BaseWriterWrap {
	return result.WrapFailure[*BaseWriter](err)
}

// WrapWriterFailureFromError creates a failed BaseWriterWrap directly from an AppError object.
func WrapWriterFailureFromError(err *appfault.AppError) BaseWriterWrap {
	return WrapWriterFailure(err)
}

// WrapWriterFailureWithId creates a failed BaseWriterWrap using an error ID and message.
func WrapWriterFailureWithId(errType errtype.Variation, msg string) BaseWriterWrap {
	return result.WrapFailureWithId[*BaseWriter](errType, msg)
}

// WrapWriterFailureWithCause creates a failed BaseWriterWrap using an error ID, cause, and message.
func WrapWriterFailureWithCause(errType errtype.Variation, cause error, msg string) BaseWriterWrap {
	return result.WrapFailureWithCause[*BaseWriter](errType, cause, msg)
}

// WrapWriterFailureFromWrap propagates failure from another Wrap into BaseWriterWrap.
func WrapWriterFailureFromWrap[U any](failed result.Wrap[U]) BaseWriterWrap {
	return result.WrapFailureFromWrap[*BaseWriter](failed)
}
