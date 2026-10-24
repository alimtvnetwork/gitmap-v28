package enums

import (
	"fmt"
	"strings"
)

// ActionType defines the execution mode of the suggested command.
type ActionType string

const (
	// ActionAutoRun indicates the recommendation is safe to auto-execute.
	ActionAutoRun ActionType = "auto-run"
	// ActionManual indicates the recommendation requires user parameter adjustment.
	ActionManual ActionType = "manual"
	// ActionCopy indicates a copy-paste command candidate.
	ActionCopy ActionType = "copy"
)

// IsValid reports whether the value matches a known variant.
func (a ActionType) IsValid() bool {
	switch a {
	case ActionAutoRun, ActionManual, ActionCopy:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (a ActionType) String() string {
	return string(a)
}

// ParseActionType parses a raw string into a validated ActionType.
// It returns a plain error (not *apperror.AppError) by recorded
// exemption: apperror imports cli/diag, so an AppError Parse here
// would close an import cycle (diag -> diag/enums -> apperror -> diag).
func ParseActionType(s string) (ActionType, error) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllActionTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "actionautorun":
		return ActionAutoRun, nil
	case "autorun":
		return ActionAutoRun, nil
	case "actionmanual":
		return ActionManual, nil
	case "manual":
		return ActionManual, nil
	case "actioncopy":
		return ActionCopy, nil
	case "copy":
		return ActionCopy, nil
	}
	return "", fmt.Errorf("invalid actiontype: %q", s)
}

// AllActionTypes returns all valid variants.
func AllActionTypes() []ActionType {
	return []ActionType{
		ActionAutoRun,
		ActionManual,
		ActionCopy,
	}
}

// IsAutoRun reports whether the value is ActionAutoRun.
func (a ActionType) IsAutoRun() bool {
	return a == ActionAutoRun
}

// IsManual reports whether the value is ActionManual.
func (a ActionType) IsManual() bool {
	return a == ActionManual
}

// IsCopy reports whether the value is ActionCopy.
func (a ActionType) IsCopy() bool {
	return a == ActionCopy
}

// HasAction reports whether the value is non-empty and valid.
func (a ActionType) HasAction() bool {
	return a.IsValid()
}
