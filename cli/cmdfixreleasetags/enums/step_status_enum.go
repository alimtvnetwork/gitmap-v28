package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// StepStatusType tracks execution outcomes for individual deletion steps.
type StepStatusType string

const (
	StepPending StepStatusType = "PENDING"
	StepSuccess StepStatusType = "SUCCESS"
	StepSkipped StepStatusType = "SKIPPED"
	StepFailed  StepStatusType = "FAILED"
)

// IsValid reports whether the value matches a known variant.
func (s StepStatusType) IsValid() bool {
	switch s {
	case StepPending, StepSuccess, StepSkipped, StepFailed:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (s StepStatusType) String() string {
	return string(s)
}

// ParseStepStatusType parses a raw string into a validated StepStatusType.
func ParseStepStatusType(s string) (StepStatusType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllStepStatusTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "steppending":
		return StepPending, nil
	case "pending":
		return StepPending, nil
	case "stepsuccess":
		return StepSuccess, nil
	case "success":
		return StepSuccess, nil
	case "stepskipped":
		return StepSkipped, nil
	case "skipped":
		return StepSkipped, nil
	case "stepfailed":
		return StepFailed, nil
	case "failed":
		return StepFailed, nil
	}
	return "", apperror.NewValidation("ParseStepStatusType", "E1000", fmt.Sprintf("invalid stepstatustype: %q", s)).WithContext("value", s)
}

// AllStepStatusTypes returns all valid variants.
func AllStepStatusTypes() []StepStatusType {
	return []StepStatusType{
		StepPending,
		StepSuccess,
		StepSkipped,
		StepFailed,
	}
}

// IsPending reports whether the value is StepPending.
func (s StepStatusType) IsPending() bool {
	return s == StepPending
}

// IsSuccess reports whether the value is StepSuccess.
func (s StepStatusType) IsSuccess() bool {
	return s == StepSuccess
}

// IsSkipped reports whether the value is StepSkipped.
func (s StepStatusType) IsSkipped() bool {
	return s == StepSkipped
}

// IsFailed reports whether the value is StepFailed.
func (s StepStatusType) IsFailed() bool {
	return s == StepFailed
}

// HasStatus reports whether the value is non-empty and valid.
func (s StepStatusType) HasStatus() bool {
	return s.IsValid()
}
