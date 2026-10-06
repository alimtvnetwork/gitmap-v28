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

