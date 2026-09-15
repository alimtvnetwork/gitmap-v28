package result

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// IsSuccess reports whether the operation succeeded without error and was matched.
func (ew *ErrorWrapper) IsSuccess() bool {
	if ew == nil {
		return false
	}

	return ew.isMatched && (ew.Err == nil || ew.Err.HasNoError())
}

// IsSafe reports whether the operation succeeded without error (alias for IsSuccess).
func (ew *ErrorWrapper) IsSafe() bool {
	if ew == nil {
		return false
	}

	return ew.isMatched && (ew.Err == nil || ew.Err.HasNoError())
}

// IsFailed reports whether the operation encountered an error.
func (ew *ErrorWrapper) IsFailed() bool {
	if ew == nil {
		return true
	}

	return ew.Err != nil && ew.Err.HasError()
}

// IsFailure reports whether the operation encountered an error (alias for IsFailed).
func (ew *ErrorWrapper) IsFailure() bool {
	if ew == nil {
		return true
	}

	return ew.Err != nil && ew.Err.HasError()
}

// IsInvalid reports whether the result is invalid, unmatched, or failed.
func (ew *ErrorWrapper) IsInvalid() bool {
	if ew == nil {
		return true
	}

	return !ew.isMatched || (ew.Err != nil && ew.Err.HasError())
}

// IsMatched reports whether the operation matched a route or command.
func (ew *ErrorWrapper) IsMatched() bool {
	if ew == nil {
		return false
	}

	return ew.isMatched
}

// IsHandled reports whether the operation was matched and handled.
func (ew *ErrorWrapper) IsHandled() bool {
	if ew == nil {
		return false
	}

	return ew.isMatched
}

// HasError reports whether an error is present.
func (ew *ErrorWrapper) HasError() bool {
	if ew == nil {
		return true
	}

	return ew.Err != nil && ew.Err.HasError()
}

// IsEmptyError reports whether no error is present.
func (ew *ErrorWrapper) IsEmptyError() bool {
	if ew == nil {
		return false
	}

	return ew.Err == nil || ew.Err.HasNoError()
}

// HasNoError reports whether no error is present.
func (ew *ErrorWrapper) HasNoError() bool {
	if ew == nil {
		return false
	}

	return ew.Err == nil || ew.Err.HasNoError()
}

// HasValidError reports whether an AppError exists and is properly structured.
func (ew *ErrorWrapper) HasValidError() bool {
	if ew == nil || ew.Err == nil {
		return false
	}

	return ew.Err.HasError() && ew.Err.IsValid()
}

// AppError returns the underlying *apperror.AppError or nil.
func (ew *ErrorWrapper) AppError() *apperror.AppError {
	if ew == nil || ew.Err == nil || ew.Err.HasNoError() {
		return nil
	}

	return ew.Err
}

// Fault returns the underlying *apperror.AppError or nil (alias).
func (ew *ErrorWrapper) Fault() *apperror.AppError {
	if ew == nil || ew.Err == nil || ew.Err.HasNoError() {
		return nil
	}

	return ew.Err
}

// AsError returns the underlying error as standard error interface, or nil if no error occurred.
func (ew *ErrorWrapper) AsError() error {
	if ew == nil || ew.Err == nil || ew.Err.HasNoError() {
		return nil
	}

	return ew.Err
}

// ErrOrNil returns the underlying error as standard error interface, or nil if no error occurred.
func (ew *ErrorWrapper) ErrOrNil() error {
	if ew == nil || ew.Err == nil || ew.Err.HasNoError() {
		return nil
	}

	return ew.Err
}

// HandleError processes the underlying error if one exists, without panicking or exiting.
func (ew *ErrorWrapper) HandleError() {
	if ew == nil || ew.Err == nil || ew.Err.HasNoError() {
		return
	}

	ew.Err.HandleError()
}

// ErrorType returns the error category type or empty if no error.
func (ew *ErrorWrapper) ErrorType() apperror.ErrorType {
	if ew == nil || ew.Err == nil || ew.Err.HasNoError() {
		return apperror.ErrorTypeNone
	}

	return ew.Err.Type
}

// IsErrorCode reports whether the underlying AppError matches code.
func (ew *ErrorWrapper) IsErrorCode(code string) bool {
	if ew == nil || ew.Err == nil || ew.Err.HasNoError() {
		return false
	}

	return ew.Err.IsErrorCode(code)
}

// IsErrorType reports whether the underlying AppError matches error type.
func (ew *ErrorWrapper) IsErrorType(errType apperror.ErrorType) bool {
	if ew == nil || ew.Err == nil || ew.Err.HasNoError() {
		return false
	}

	return ew.Err.Type == errType
}

// AsError converts an ErrorWrapper to standard Go error interface.
// This package-level function enables clean one-line returns at command entrypoints:
// return result.AsError(routeEnvSub(sub, rest))
func AsError(ew ErrorWrapper) error {
	return ew.AsError()
}

// SuccessWrapper constructs an ErrorWrapper representing a successful matched operation.
func SuccessWrapper() ErrorWrapper {
	return ErrorWrapper{isMatched: true}
}

// FailureWrapper constructs an ErrorWrapper from an *apperror.AppError.
func FailureWrapper(err *apperror.AppError) ErrorWrapper {
	return ErrorWrapper{Err: err, isMatched: true}
}

// FailureWrapperErr constructs an ErrorWrapper from any standard error.
func FailureWrapperErr(err error) ErrorWrapper {
	if err == nil {
		return SuccessWrapper()
	}

	if appErr, isAppErr := err.(*apperror.AppError); isAppErr {
		return ErrorWrapper{Err: appErr, isMatched: true}
	}

	return ErrorWrapper{Err: apperror.WrapSimple(err, "errorwrapper"), isMatched: true}
}

// UnmatchedWrapper constructs an ErrorWrapper representing an unmatched route.
func UnmatchedWrapper() ErrorWrapper {
	return ErrorWrapper{isMatched: false}
}

// OkWrapper is a canonical alias for SuccessWrapper.
func OkWrapper() ErrorWrapper {
	return SuccessWrapper()
}

// FailWrapper is a canonical alias for FailureWrapper.
func FailWrapper(err *apperror.AppError) ErrorWrapper {
	return FailureWrapper(err)
}

// FailWrapperErr is a canonical alias for FailureWrapperErr.
func FailWrapperErr(err error) ErrorWrapper {
	return FailureWrapperErr(err)
}

// MatchWrapper constructs an ErrorWrapper representing a matched route with optional error.
func MatchWrapper(err error) ErrorWrapper {
	return FailureWrapperErr(err)
}

// MatchWrapperAppErr constructs an ErrorWrapper representing a matched route with an *apperror.AppError.
func MatchWrapperAppErr(appErr *apperror.AppError) ErrorWrapper {
	return FailureWrapper(appErr)
}
