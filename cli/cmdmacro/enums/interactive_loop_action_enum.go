package enums

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// InteractiveLoopActionType is the control-flow result of one interactive macro step.
type InteractiveLoopActionType int

const (
	LoopActionTypeContinue InteractiveLoopActionType = iota
	LoopActionTypeBreak
)

// IsValid reports whether the loop action matches a known action.
func (a InteractiveLoopActionType) IsValid() bool {
	switch a {
	case LoopActionTypeContinue, LoopActionTypeBreak:
		return true
	default:
		return false
	}
}

// String returns the name of the loop action.
func (a InteractiveLoopActionType) String() string {
	switch a {
	case LoopActionTypeContinue:
		return "continue"
	case LoopActionTypeBreak:
		return "break"
	default:
		return ""
	}
}

// ParseInteractiveLoopActionType parses a raw string (name or number) into a validated InteractiveLoopActionType.
func ParseInteractiveLoopActionType(s string) (InteractiveLoopActionType, *apperror.AppError) {
	cleaned := strings.TrimSpace(s)
	if num, err := strconv.Atoi(cleaned); err == nil {
		val := InteractiveLoopActionType(num)
		if val.IsValid() {
			return val, nil
		}
	}

	for _, val := range AllInteractiveLoopActionTypes() {
		if strings.EqualFold(val.String(), cleaned) {
			return val, nil
		}
	}

	return 0, apperror.NewValidation("ParseInteractiveLoopActionType", "E1000", fmt.Sprintf("invalid interactive loop action: %q", s)).WithContext("value", s)
}

// AllInteractiveLoopActionTypes returns all valid InteractiveLoopActionType variants.
func AllInteractiveLoopActionTypes() []InteractiveLoopActionType {
	return []InteractiveLoopActionType{LoopActionTypeContinue, LoopActionTypeBreak}
}

// IsContinue reports whether the action is continue.
func (a InteractiveLoopActionType) IsContinue() bool { return a == LoopActionTypeContinue }

// IsBreak reports whether the action is break.
func (a InteractiveLoopActionType) IsBreak() bool { return a == LoopActionTypeBreak }
