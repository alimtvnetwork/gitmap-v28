package regexnew

import (
	"fmt"
	"regexp"
)

// IsMatchLock compiles regex under lock and returns whether comparing matches.
func IsMatchLock(regex, comparing string) bool {
	regEx, _ := CreateLock(regex)
	if regEx != nil && regEx.MatchString(comparing) {
		return true
	}

	return false
}

// IsMatchFailed returns true if comparing does not match regex.
func IsMatchFailed(regex, comparing string) bool {
	return !IsMatchLock(regex, comparing)
}

// MatchError returns nil on match or a validation error.
func MatchError(regex, comparing string) error {
	regEx, err := Create(regex)
	if regEx != nil && regEx.MatchString(comparing) {
		return nil
	}

	return regExMatchValidationError(regex, comparing, err, regEx)
}

// MatchErrorLock returns nil on match or a validation error under lock.
func MatchErrorLock(regex, comparing string) error {
	regEx, err := CreateLock(regex)
	if regEx != nil && regEx.MatchString(comparing) {
		return nil
	}

	return regExMatchValidationError(regex, comparing, err, regEx)
}

// MatchUsingFuncErrorLock performs custom validation under lock.
func MatchUsingFuncErrorLock(
	regex, comparing string,
	matchFunc RegexValidationFunc,
) error {
	regEx, err := CreateLock(regex)
	if regEx != nil && matchFunc != nil && matchFunc(regEx, comparing) {
		return nil
	}

	return regExMatchValidationError(regex, comparing, err, regEx)
}

// MatchUsingCustomizeErrorFuncLock executes custom error generator on mismatch under lock.
func MatchUsingCustomizeErrorFuncLock(
	regex, comparing string,
	customErr CustomizeErr,
) error {
	regEx, err := CreateLock(regex)
	if regEx != nil && regEx.MatchString(comparing) {
		return nil
	}

	if customErr != nil {
		return customErr(regex, comparing, err, regEx)
	}

	return regExMatchValidationError(regex, comparing, err, regEx)
}

func regExMatchValidationError(
	regex string,
	comparing string,
	err error,
	regEx *regexp.Regexp,
) error {
	if err != nil {
		return fmt.Errorf(
			"[%q], regex pattern compile failed / invalid cannot match with [%q]",
			err.Error(),
			comparing)
	}

	if regEx == nil {
		return fmt.Errorf(
			"given regex pattern [%q] invalid cannot match with [%q]",
			regex,
			comparing)
	}

	return fmt.Errorf(
		"given regex pattern [%q] doesn't match with [%q]",
		regex,
		comparing)
}
