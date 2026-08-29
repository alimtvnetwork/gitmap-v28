package apperror

// ErrorType identifies the category of an application error.
type ErrorType string

const (
	ErrorTypeValidation   ErrorType = "VALIDATION"
	ErrorTypePrecondition ErrorType = "PRECONDITION"
	ErrorTypeNotFound     ErrorType = "NOT_FOUND"
	ErrorTypeExecution    ErrorType = "EXECUTION"
	ErrorTypeAbort        ErrorType = "ABORT"
	ErrorTypeInternal     ErrorType = "INTERNAL"
)

// SeverityType indicates the severity level of an error.
type SeverityType string

const (
	SeverityInfo  SeverityType = "INFO"
	SeverityWarn  SeverityType = "WARN"
	SeverityError SeverityType = "ERROR"
	SeverityFatal SeverityType = "FATAL"
)
