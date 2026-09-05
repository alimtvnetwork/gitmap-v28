package appfaults

import "coding-guidelines/common/pkg/appfault"

// Formatting and default constants for fault collections.
const (
	DefaultEmptyMessage = "No faults recorded."
)

// FaultPredicate defines a filter test function over an *appfault.AppError.
type FaultPredicate func(e *appfault.AppError) bool

// Predicate is an alias for FaultPredicate.
type Predicate = FaultPredicate

// FaultTransformer defines a transformation function from one *appfault.AppError to another.
type FaultTransformer func(e *appfault.AppError) *appfault.AppError
