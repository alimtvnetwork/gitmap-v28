package errtype

// Standard Variation error type code constants.
const (
	// None indicates no error occurred (successful state).
	None Variation = 0

	// NoError is an alias for None.
	NoError Variation = None

	// Generic represents a standard, unspecified error.
	Generic Variation = 1

	// Validation represents input, schema, or invariant validation failures.
	Validation Variation = 2

	// NotFound represents a requested resource or record that does not exist.
	NotFound Variation = 3

	// Precondition represents unsatisfied state prerequisites.
	Precondition Variation = 4

	// Execution represents general runtime execution failures.
	Execution Variation = 5

	// Database represents database query, execution, or connection failures.
	Database Variation = 6

	// Network represents remote network transport or connectivity failures.
	Network Variation = 7

	// Timeout represents an operation exceeding its deadline.
	Timeout Variation = 8

	// IO represents input/output or filesystem failures.
	IO Variation = 9

	// Unauthorized represents unauthenticated access attempts.
	Unauthorized Variation = 10

	// Forbidden represents authenticated but unauthorized access attempts.
	Forbidden Variation = 11

	// Internal represents unexpected internal server faults.
	Internal Variation = 12

	// Unknown represents an unclassified error state.
	Unknown Variation = 13

	// Serialization represents data serialization, deserialization, or encoding failures.
	Serialization Variation = 14
)

// ProcessStateType constants conforming to BaseEnum.
const (
	ProcessStatePending   ProcessStateType = "Pending"
	ProcessStateRunning   ProcessStateType = "Running"
	ProcessStateCompleted ProcessStateType = "Completed"
	ProcessStateFailed    ProcessStateType = "Failed"
	ProcessStateCancelled ProcessStateType = "Cancelled"
	ProcessStateUnknown   ProcessStateType = "Unknown"
)

// LogLevelType constants conforming to NumberEnum and BaseEnum.
const (
	LogLevelDebug LogLevelType = 1
	LogLevelInfo  LogLevelType = 2
	LogLevelWarn  LogLevelType = 3
	LogLevelError LogLevelType = 4
	LogLevelFatal LogLevelType = 5
)

type (
	VariationPredicate func(v Variation) bool

	EnumPredicate func(e BaseEnum) bool
)
