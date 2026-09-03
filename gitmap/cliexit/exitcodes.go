// Package cliexit — exitcodes.go defines strongly-typed exit codes and specialized helpers.
package cliexit

// ExitCodeType classifies CLI process termination codes.
type ExitCodeType int

const (
	ExitCodeSuccess         ExitCodeType = 0
	ExitCodeGeneralError    ExitCodeType = 1
	ExitCodeUsageError      ExitCodeType = 2
	ExitCodePartialFailure  ExitCodeType = 3
	ExitCodeNotFound        ExitCodeType = 4
	ExitCodeValidationError ExitCodeType = 5
)

// HandleValidationError reports a validation failure and exits with ExitCodeValidationError.
func HandleValidationError(err error) {
	HandleError(err, int(ExitCodeValidationError))
}

// HandleUsageError reports a CLI usage or syntax error and exits with ExitCodeUsageError.
func HandleUsageError(err error) {
	HandleError(err, int(ExitCodeUsageError))
}

// HandleGeneralError reports an operational failure and exits with ExitCodeGeneralError.
func HandleGeneralError(err error) {
	HandleError(err, int(ExitCodeGeneralError))
}

// HandleNotFound reports a missing resource failure and exits with ExitCodeNotFound.
func HandleNotFound(err error) {
	HandleError(err, int(ExitCodeNotFound))
}

// HandleSuccess flushes output pipes and exits cleanly with ExitCodeSuccess.
func HandleSuccess() {
	runFlushers()
	exitFunc(int(ExitCodeSuccess))
}
