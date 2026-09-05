package regexnew

import "regexp"

// Default capacity constants for pre-allocated map sizing.
const (
	DefaultCapacity = 32
	EmptyString     = ""
	Tab             = "\t"
)

type (
	// RegexValidationFunc defines a custom validation matcher against a compiled regex.
	RegexValidationFunc func(regex *regexp.Regexp, lookingTerm string) bool

	// CustomizeErr creates a custom formatted error when a regular expression fails to match.
	CustomizeErr func(
		regexPattern string,
		matchLookingTerm string,
		err error,
		regexp *regexp.Regexp,
	) error
)
