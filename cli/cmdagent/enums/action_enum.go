package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// ActionType represents an agent telemetry action.
type ActionType string

const (
	ActionTypeSearch   ActionType = "SEARCH"
	ActionTypeRead     ActionType = "READ"
	ActionTypeWrite    ActionType = "WRITE"
	ActionTypeExec     ActionType = "EXEC"
	ActionTypeLint     ActionType = "LINT"
	ActionTypeCheck    ActionType = "CHECK"
	ActionTypeClaim    ActionType = "CLAIM"
	ActionTypeStart    ActionType = "START"
	ActionTypeComplete ActionType = "COMPLETE"
	ActionTypeFail     ActionType = "FAIL"
	ActionTypeCrash    ActionType = "CRASH"
)

// IsValid reports whether the action matches a known agent action.
func (a ActionType) IsValid() bool {
	switch a {
	case ActionTypeSearch, ActionTypeRead, ActionTypeWrite, ActionTypeExec,
		ActionTypeLint, ActionTypeCheck, ActionTypeClaim, ActionTypeStart,
		ActionTypeComplete, ActionTypeFail, ActionTypeCrash:
		return true
	default:
		return false
	}
}

// String returns the string representation of the action.
func (a ActionType) String() string {
	return string(a)
}

// ParseActionType parses a raw string into a validated ActionType.
func ParseActionType(s string) (ActionType, *apperror.AppError) {
	val := ActionType(strings.ToUpper(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseActionType", "E1000", fmt.Sprintf("invalid agent action: %q", s)).WithContext("value", s)
}

// AllActionTypes returns all valid ActionType variants.
func AllActionTypes() []ActionType {
	return []ActionType{
		ActionTypeSearch, ActionTypeRead, ActionTypeWrite, ActionTypeExec,
		ActionTypeLint, ActionTypeCheck, ActionTypeClaim, ActionTypeStart,
		ActionTypeComplete, ActionTypeFail, ActionTypeCrash,
	}
}

// IsSearch reports whether the action is search.
func (a ActionType) IsSearch() bool { return a == ActionTypeSearch }

// IsRead reports whether the action is read.
func (a ActionType) IsRead() bool { return a == ActionTypeRead }

// IsWrite reports whether the action is write.
func (a ActionType) IsWrite() bool { return a == ActionTypeWrite }

// IsExec reports whether the action is exec.
func (a ActionType) IsExec() bool { return a == ActionTypeExec }

// IsLint reports whether the action is lint.
func (a ActionType) IsLint() bool { return a == ActionTypeLint }

// IsCheck reports whether the action is check.
func (a ActionType) IsCheck() bool { return a == ActionTypeCheck }

// IsClaim reports whether the action is claim.
func (a ActionType) IsClaim() bool { return a == ActionTypeClaim }

// IsStart reports whether the action is start.
func (a ActionType) IsStart() bool { return a == ActionTypeStart }

// IsComplete reports whether the action is complete.
func (a ActionType) IsComplete() bool { return a == ActionTypeComplete }

// IsFail reports whether the action is fail.
func (a ActionType) IsFail() bool { return a == ActionTypeFail }

// IsCrash reports whether the action is crash.
func (a ActionType) IsCrash() bool { return a == ActionTypeCrash }
