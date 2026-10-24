package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strconv"
	"strings"
)

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

// IsValid reports whether the value matches a known variant.
func (e ExitCodeType) IsValid() bool {
	switch e {
	case ExitCodeSuccess, ExitCodeGeneralError, ExitCodeUsageError, ExitCodePartialFailure, ExitCodeNotFound, ExitCodeValidationError:
		return true
	default:
		return false
	}
}

// String returns the numeric representation of the value.
func (e ExitCodeType) String() string {
	return strconv.Itoa(int(e))
}

// ParseExitCodeType parses a raw string into a validated ExitCodeType.
func ParseExitCodeType(s string) (ExitCodeType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	if n, err := strconv.Atoi(trimmed); err == nil {
		val := ExitCodeType(n)
		if val.IsValid() {
			return val, nil
		}
	}
	for _, val := range AllExitCodeTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "exitcodesuccess":
		return ExitCodeSuccess, nil
	case "success":
		return ExitCodeSuccess, nil
	case "exitcodegeneralerror":
		return ExitCodeGeneralError, nil
	case "generalerror":
		return ExitCodeGeneralError, nil
	case "exitcodeusageerror":
		return ExitCodeUsageError, nil
	case "usageerror":
		return ExitCodeUsageError, nil
	case "exitcodepartialfailure":
		return ExitCodePartialFailure, nil
	case "partialfailure":
		return ExitCodePartialFailure, nil
	case "exitcodenotfound":
		return ExitCodeNotFound, nil
	case "notfound":
		return ExitCodeNotFound, nil
	case "exitcodevalidationerror":
		return ExitCodeValidationError, nil
	case "validationerror":
		return ExitCodeValidationError, nil
	}
	return 0, apperror.NewValidation("ParseExitCodeType", "E1000", fmt.Sprintf("invalid exitcodetype: %q", s)).WithContext("value", s)
}

// AllExitCodeTypes returns all valid variants.
func AllExitCodeTypes() []ExitCodeType {
	return []ExitCodeType{
		ExitCodeSuccess,
		ExitCodeGeneralError,
		ExitCodeUsageError,
		ExitCodePartialFailure,
		ExitCodeNotFound,
		ExitCodeValidationError,
	}
}

// IsSuccess reports whether the value is ExitCodeSuccess.
func (e ExitCodeType) IsSuccess() bool {
	return e == ExitCodeSuccess
}

// IsGeneralError reports whether the value is ExitCodeGeneralError.
func (e ExitCodeType) IsGeneralError() bool {
	return e == ExitCodeGeneralError
}

// IsUsageError reports whether the value is ExitCodeUsageError.
func (e ExitCodeType) IsUsageError() bool {
	return e == ExitCodeUsageError
}

// IsPartialFailure reports whether the value is ExitCodePartialFailure.
func (e ExitCodeType) IsPartialFailure() bool {
	return e == ExitCodePartialFailure
}

// IsNotFound reports whether the value is ExitCodeNotFound.
func (e ExitCodeType) IsNotFound() bool {
	return e == ExitCodeNotFound
}

// IsValidationError reports whether the value is ExitCodeValidationError.
func (e ExitCodeType) IsValidationError() bool {
	return e == ExitCodeValidationError
}

// HasCode reports whether the value is non-empty and valid.
func (e ExitCodeType) HasCode() bool {
	return e.IsValid()
}
