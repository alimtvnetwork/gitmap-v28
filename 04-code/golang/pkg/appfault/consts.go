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

type (
	FaultFormatter func(e *AppError) string

	ResultFormatter[T any] func(r Result[T]) string

	FaultPredicate func(e *AppError) bool

	ErrorHandler func(e *AppError)
)
