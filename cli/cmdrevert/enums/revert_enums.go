package enums

import (
	"fmt"
	"strings"
)

// RevertModeType specifies the execution sub-mode for cmdrevert.
type RevertModeType string

const (
	RevertModeVersion  RevertModeType = "version"
	RevertModeRunner   RevertModeType = "revert-runner"
	RevertModeTxn      RevertModeType = "txn"
	RevertModeLastTxn  RevertModeType = "last-txn"
	RevertModeLastNTxn RevertModeType = "last-n-txn"
	RevertModeListTxn  RevertModeType = "list-txn"
	RevertModeShowTxn  RevertModeType = "show-txn"
	RevertModePruneTxn RevertModeType = "prune-txn"
)

// IsValid reports whether the revert mode is recognized.
func (m RevertModeType) IsValid() bool {
	switch m {
	case RevertModeVersion, RevertModeRunner, RevertModeTxn,
		RevertModeLastTxn, RevertModeLastNTxn, RevertModeListTxn,
		RevertModeShowTxn, RevertModePruneTxn:
		return true
	default:
		return false
	}
}

// String returns the string representation of the mode.
func (m RevertModeType) String() string {
	return string(m)
}

// ParseRevertModeType parses a string into a RevertModeType.
func ParseRevertModeType(s string) (RevertModeType, error) {
	val := RevertModeType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", fmt.Errorf("invalid revert mode: %q", s)
}

// AllRevertModeTypes returns all supported RevertModeType values.
func AllRevertModeTypes() []RevertModeType {
	return []RevertModeType{
		RevertModeVersion,
		RevertModeRunner,
		RevertModeTxn,
		RevertModeLastTxn,
		RevertModeLastNTxn,
		RevertModeListTxn,
		RevertModeShowTxn,
		RevertModePruneTxn,
	}
}

// IsVersion reports whether the mode is version checkout.
func (m RevertModeType) IsVersion() bool {
	return m == RevertModeVersion
}

// IsRunner reports whether the mode is the internal revert-runner.
func (m RevertModeType) IsRunner() bool {
	return m == RevertModeRunner
}

// IsTxn reports whether the mode is single transaction revert.
func (m RevertModeType) IsTxn() bool {
	return m == RevertModeTxn
}

// IsLastTxn reports whether the mode is latest transaction revert.
func (m RevertModeType) IsLastTxn() bool {
	return m == RevertModeLastTxn
}

// IsLastNTxn reports whether the mode is last-N transaction revert.
func (m RevertModeType) IsLastNTxn() bool {
	return m == RevertModeLastNTxn
}

// IsListTxn reports whether the mode lists transactions.
func (m RevertModeType) IsListTxn() bool {
	return m == RevertModeListTxn
}

// IsShowTxn reports whether the mode displays single transaction details.
func (m RevertModeType) IsShowTxn() bool {
	return m == RevertModeShowTxn
}

// IsPruneTxn reports whether the mode prunes old transactions.
func (m RevertModeType) IsPruneTxn() bool {
	return m == RevertModePruneTxn
}

// HasMode reports whether the mode is valid.
func (m RevertModeType) HasMode() bool {
	return m.IsValid()
}

// RevertPromptResponseType specifies user responses to revert confirmation prompts.
type RevertPromptResponseType string

const (
	RevertPromptResponseYes RevertPromptResponseType = "yes"
	RevertPromptResponseNo  RevertPromptResponseType = "no"
)

// IsValid reports whether the prompt response is recognized.
func (r RevertPromptResponseType) IsValid() bool {
	switch r {
	case RevertPromptResponseYes, RevertPromptResponseNo:
		return true
	default:
		return false
	}
}

// String returns the string representation of the prompt response.
func (r RevertPromptResponseType) String() string {
	return string(r)
}

// ParseRevertPromptResponseType parses a string into a RevertPromptResponseType.
func ParseRevertPromptResponseType(s string) (RevertPromptResponseType, error) {
	val := RevertPromptResponseType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", fmt.Errorf("invalid prompt response: %q", s)
}

// AllRevertPromptResponseTypes returns all valid prompt responses.
func AllRevertPromptResponseTypes() []RevertPromptResponseType {
	return []RevertPromptResponseType{
		RevertPromptResponseYes,
		RevertPromptResponseNo,
	}
}

// IsYes reports whether the response is yes.
func (r RevertPromptResponseType) IsYes() bool {
	return r == RevertPromptResponseYes
}

// IsNo reports whether the response is no.
func (r RevertPromptResponseType) IsNo() bool {
	return r == RevertPromptResponseNo
}

// HasResponse reports whether the response is valid.
func (r RevertPromptResponseType) HasResponse() bool {
	return r.IsValid()
}

// RevertExitCodeType specifies exit codes for cmdrevert execution.
type RevertExitCodeType int

const (
	ExitCodeSuccess      RevertExitCodeType = 0
	ExitCodeGeneralError RevertExitCodeType = 1
	ExitCodeUsageError   RevertExitCodeType = 2
)

// IsValid reports whether the exit code is recognized.
func (c RevertExitCodeType) IsValid() bool {
	switch c {
	case ExitCodeSuccess, ExitCodeGeneralError, ExitCodeUsageError:
		return true
	default:
		return false
	}
}

// String returns a human-readable name for the exit code.
func (c RevertExitCodeType) String() string {
	switch c {
	case ExitCodeSuccess:
		return "success"
	case ExitCodeGeneralError:
		return "general_error"
	case ExitCodeUsageError:
		return "usage_error"
	default:
		return fmt.Sprintf("exit_code(%d)", c)
	}
}

// ParseRevertExitCodeType parses a string into a RevertExitCodeType.
func ParseRevertExitCodeType(s string) (RevertExitCodeType, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "0", "success":
		return ExitCodeSuccess, nil
	case "1", "general", "general_error":
		return ExitCodeGeneralError, nil
	case "2", "usage", "usage_error":
		return ExitCodeUsageError, nil
	default:
		return ExitCodeGeneralError, fmt.Errorf("invalid revert exit code: %q", s)
	}
}

// AllRevertExitCodeTypes returns all supported exit codes.
func AllRevertExitCodeTypes() []RevertExitCodeType {
	return []RevertExitCodeType{
		ExitCodeSuccess,
		ExitCodeGeneralError,
		ExitCodeUsageError,
	}
}

// IsSuccess reports whether the exit code indicates success.
func (c RevertExitCodeType) IsSuccess() bool {
	return c == ExitCodeSuccess
}

// IsGeneralError reports whether the exit code indicates a general error.
func (c RevertExitCodeType) IsGeneralError() bool {
	return c == ExitCodeGeneralError
}

// IsUsageError reports whether the exit code indicates a usage error.
func (c RevertExitCodeType) IsUsageError() bool {
	return c == ExitCodeUsageError
}

// HasExitCode reports whether the exit code is valid.
func (c RevertExitCodeType) HasExitCode() bool {
	return c.IsValid()
}

// Code returns the raw integer exit code.
func (c RevertExitCodeType) Code() int {
	return int(c)
}
