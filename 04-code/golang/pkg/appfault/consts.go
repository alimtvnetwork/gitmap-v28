package appfault

// Formatting and layout constants used across display and compilation routines.
const (
	DelimiterLine = "──────────────────────────────────────────────────"
	HeaderPrefix  = "══ [AppError] "
	SectionPrefix = "── "
	IndentTab     = "    "
	BulletPrefix  = " • "
	Newline       = "\n"
)

// Priority level constants.
const (
	PriorityUnknown PriorityType = iota
	PriorityLow
	PriorityNormal
	PriorityHigh
	PriorityCritical
)

// Severity level constants.
const (
	SeverityUnknown SeverityType = iota
	SeverityInfo
	SeverityWarn
	SeverityError
	SeverityCritical
	SeverityFatal
)

// FaultFormatter defines a custom formatting function for AppError.
type FaultFormatter func(e *AppError) string

// ResultFormatter formats a generic Result[T] container into a string.
type ResultFormatter[T any] func(r Result[T]) string

// FaultPredicate tests an AppError against a condition.
type FaultPredicate func(e *AppError) bool

// ErrorHandler handles an AppError callback.
type ErrorHandler func(e *AppError)
