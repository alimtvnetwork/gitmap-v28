// Package cliexit — exitcodes.go defines strongly-typed exit codes and specialized helpers.
package cliexit

import "github.com/alimtvnetwork/gitmap-v28/cli/cliexit/enums"

// ExitCodeType classifies CLI process termination codes.
// ExitCodeType aliases the strongly-typed enum from the enums package.
type ExitCodeType = enums.ExitCodeType

const (
	ExitCodeSuccess         = enums.ExitCodeSuccess
	ExitCodeGeneralError    = enums.ExitCodeGeneralError
	ExitCodeUsageError      = enums.ExitCodeUsageError
	ExitCodePartialFailure  = enums.ExitCodePartialFailure
	ExitCodeNotFound        = enums.ExitCodeNotFound
	ExitCodeValidationError = enums.ExitCodeValidationError
)
