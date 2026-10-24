package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// DeployActionType represents the decided action for a specific file transfer.
type DeployActionType string

const (
	// ActionTransferToRemote copies local file to remote target.
	ActionTransferToRemote DeployActionType = "transfer-to-remote"
	// ActionTransferToLocal copies remote file to local workspace.
	ActionTransferToLocal DeployActionType = "transfer-to-local"
	// ActionSkip skips file without transfer.
	ActionSkip DeployActionType = "skip"
	// ActionConflictPrompt requests user decision interactively.
	ActionConflictPrompt DeployActionType = "conflict-prompt"
)

// IsValid reports whether the action matches a known deploy action.
func (a DeployActionType) IsValid() bool {
	switch a {
	case ActionTransferToRemote, ActionTransferToLocal,
		ActionSkip, ActionConflictPrompt:
		return true
	default:
		return false
	}
}

// String returns the string representation of the action.
func (a DeployActionType) String() string {
	return string(a)
}

// ParseDeployActionType parses a raw string into a validated DeployActionType.
func ParseDeployActionType(s string) (DeployActionType, *apperror.AppError) {
	val := DeployActionType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseDeployActionType", "E1000", fmt.Sprintf("invalid deploy action: %q", s)).WithContext("value", s)
}

// AllDeployActionTypes returns all valid DeployActionType variants.
func AllDeployActionTypes() []DeployActionType {
	return []DeployActionType{
		ActionTransferToRemote,
		ActionTransferToLocal,
		ActionSkip,
		ActionConflictPrompt,
	}
}

// IsTransferToRemote reports whether the action copies local file to remote.
func (a DeployActionType) IsTransferToRemote() bool {
	return a == ActionTransferToRemote
}

// IsTransferToLocal reports whether the action copies remote file to local.
func (a DeployActionType) IsTransferToLocal() bool {
	return a == ActionTransferToLocal
}

// IsSkip reports whether the action skips file without transfer.
func (a DeployActionType) IsSkip() bool {
	return a == ActionSkip
}

// IsConflictPrompt reports whether the action requests user decision.
func (a DeployActionType) IsConflictPrompt() bool {
	return a == ActionConflictPrompt
}

// HasAction reports whether the action is non-empty and valid.
func (a DeployActionType) HasAction() bool {
	return a.IsValid()
}
