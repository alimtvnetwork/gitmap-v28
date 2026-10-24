package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// AiTaskStatusType defines valid task statuses.
type AiTaskStatusType string

const (
	AiTaskStatusInProgress AiTaskStatusType = "in_progress"
	AiTaskStatusCompleted  AiTaskStatusType = "completed"
	AiTaskStatusFailed     AiTaskStatusType = "failed"
	AiTaskStatusCanceled   AiTaskStatusType = "canceled"
)

// AiOperationType defines the file interaction performed.
type AiOperationType string

const (
	AiOperationRead   AiOperationType = "read"
	AiOperationEdit   AiOperationType = "edit"
	AiOperationDelete AiOperationType = "delete"
)

// IsValid reports whether the task status matches a known status.
func (s AiTaskStatusType) IsValid() bool {
	switch s {
	case AiTaskStatusInProgress, AiTaskStatusCompleted, AiTaskStatusFailed, AiTaskStatusCanceled:
		return true
	default:
		return false
	}
}

// String returns the string representation of the task status.
func (s AiTaskStatusType) String() string {
	return string(s)
}

// ParseAiTaskStatusType parses a raw string into a validated AiTaskStatusType.
func ParseAiTaskStatusType(s string) (AiTaskStatusType, *apperror.AppError) {
	val := AiTaskStatusType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseAiTaskStatusType", "E1000", fmt.Sprintf("invalid AI task status: %q", s)).WithContext("value", s)
}

// AllAiTaskStatusTypes returns all valid AiTaskStatusType variants.
func AllAiTaskStatusTypes() []AiTaskStatusType {
	return []AiTaskStatusType{AiTaskStatusInProgress, AiTaskStatusCompleted, AiTaskStatusFailed, AiTaskStatusCanceled}
}

// IsInProgress reports whether the status is in_progress.
func (s AiTaskStatusType) IsInProgress() bool { return s == AiTaskStatusInProgress }

// IsCompleted reports whether the status is completed.
func (s AiTaskStatusType) IsCompleted() bool { return s == AiTaskStatusCompleted }

// IsFailed reports whether the status is failed.
func (s AiTaskStatusType) IsFailed() bool { return s == AiTaskStatusFailed }

// IsCanceled reports whether the status is canceled.
func (s AiTaskStatusType) IsCanceled() bool { return s == AiTaskStatusCanceled }

// IsValid reports whether the operation matches a known operation.
func (o AiOperationType) IsValid() bool {
	switch o {
	case AiOperationRead, AiOperationEdit, AiOperationDelete:
		return true
	default:
		return false
	}
}

// String returns the string representation of the operation.
func (o AiOperationType) String() string {
	return string(o)
}

// ParseAiOperationType parses a raw string into a validated AiOperationType.
func ParseAiOperationType(s string) (AiOperationType, *apperror.AppError) {
	val := AiOperationType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseAiOperationType", "E1000", fmt.Sprintf("invalid AI operation: %q", s)).WithContext("value", s)
}

// AllAiOperationTypes returns all valid AiOperationType variants.
func AllAiOperationTypes() []AiOperationType {
	return []AiOperationType{AiOperationRead, AiOperationEdit, AiOperationDelete}
}

// IsRead reports whether the operation is read.
func (o AiOperationType) IsRead() bool { return o == AiOperationRead }

// IsEdit reports whether the operation is edit.
func (o AiOperationType) IsEdit() bool { return o == AiOperationEdit }

// IsDelete reports whether the operation is delete.
func (o AiOperationType) IsDelete() bool { return o == AiOperationDelete }
