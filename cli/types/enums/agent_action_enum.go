package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// AgentActionType represents the category of action performed by an autonomous agent.
type AgentActionType string

const (
	ActionSearch   AgentActionType = "SEARCH"
	ActionRead     AgentActionType = "READ"
	ActionWrite    AgentActionType = "WRITE"
	ActionEdit     AgentActionType = "EDIT"
	ActionExec     AgentActionType = "EXEC"
	ActionLint     AgentActionType = "LINT"
	ActionCheck    AgentActionType = "CHECK"
	ActionClaim    AgentActionType = "CLAIM"
	ActionComplete AgentActionType = "COMPLETE"
	ActionFail     AgentActionType = "FAIL"
	ActionCrash    AgentActionType = "CRASH"
)

func (a AgentActionType) IsValid() bool {
	switch a {
	case ActionSearch, ActionRead, ActionWrite, ActionEdit, ActionExec,
		ActionLint, ActionCheck, ActionClaim, ActionComplete, ActionFail, ActionCrash:
		return true
	default:
		return false
	}
}

func (a AgentActionType) String() string {
	return string(a)
}

// ParseAgentActionType parses a raw string into a validated AgentActionType.
func ParseAgentActionType(s string) (AgentActionType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllAgentActionTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "actionsearch":
		return ActionSearch, nil
	case "search":
		return ActionSearch, nil
	case "actionread":
		return ActionRead, nil
	case "read":
		return ActionRead, nil
	case "actionwrite":
		return ActionWrite, nil
	case "write":
		return ActionWrite, nil
	case "actionedit":
		return ActionEdit, nil
	case "edit":
		return ActionEdit, nil
	case "actionexec":
		return ActionExec, nil
	case "exec":
		return ActionExec, nil
	case "actionlint":
		return ActionLint, nil
	case "lint":
		return ActionLint, nil
	case "actioncheck":
		return ActionCheck, nil
	case "check":
		return ActionCheck, nil
	case "actionclaim":
		return ActionClaim, nil
	case "claim":
		return ActionClaim, nil
	case "actioncomplete":
		return ActionComplete, nil
	case "complete":
		return ActionComplete, nil
	case "actionfail":
		return ActionFail, nil
	case "fail":
		return ActionFail, nil
	case "actioncrash":
		return ActionCrash, nil
	case "crash":
		return ActionCrash, nil
	}
	return "", apperror.NewValidation("ParseAgentActionType", "E1000", fmt.Sprintf("invalid agentactiontype: %q", s)).WithContext("value", s)
}

// AllAgentActionTypes returns all valid variants.
func AllAgentActionTypes() []AgentActionType {
	return []AgentActionType{
		ActionSearch,
		ActionRead,
		ActionWrite,
		ActionEdit,
		ActionExec,
		ActionLint,
		ActionCheck,
		ActionClaim,
		ActionComplete,
		ActionFail,
		ActionCrash,
	}
}

// IsSearch reports whether the value is ActionSearch.
func (a AgentActionType) IsSearch() bool {
	return a == ActionSearch
}

// IsRead reports whether the value is ActionRead.
func (a AgentActionType) IsRead() bool {
	return a == ActionRead
}

// IsWrite reports whether the value is ActionWrite.
func (a AgentActionType) IsWrite() bool {
	return a == ActionWrite
}

// IsEdit reports whether the value is ActionEdit.
func (a AgentActionType) IsEdit() bool {
	return a == ActionEdit
}

// IsExec reports whether the value is ActionExec.
func (a AgentActionType) IsExec() bool {
	return a == ActionExec
}

// IsLint reports whether the value is ActionLint.
func (a AgentActionType) IsLint() bool {
	return a == ActionLint
}

// IsCheck reports whether the value is ActionCheck.
func (a AgentActionType) IsCheck() bool {
	return a == ActionCheck
}

// IsClaim reports whether the value is ActionClaim.
func (a AgentActionType) IsClaim() bool {
	return a == ActionClaim
}

// IsComplete reports whether the value is ActionComplete.
func (a AgentActionType) IsComplete() bool {
	return a == ActionComplete
}

// IsFail reports whether the value is ActionFail.
func (a AgentActionType) IsFail() bool {
	return a == ActionFail
}

// IsCrash reports whether the value is ActionCrash.
func (a AgentActionType) IsCrash() bool {
	return a == ActionCrash
}

// HasAction reports whether the value is non-empty and valid.
func (a AgentActionType) HasAction() bool {
	return a.IsValid()
}
