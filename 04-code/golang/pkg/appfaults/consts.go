package appfaults

import "coding-guidelines/common/pkg/appfault"

// Formatting and default constants for fault collections.
const (
	DefaultEmptyMessage = "No faults recorded."
)

type (
	FaultPredicate func(e *appfault.AppError) bool

	Predicate = FaultPredicate

	FaultTransformer func(e *appfault.AppError) *appfault.AppError
)
