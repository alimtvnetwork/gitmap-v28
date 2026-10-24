package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// PruneScopeType specifies the scope of prune operations.
type PruneScopeType string

const (
	PruneScopeLocal  PruneScopeType = "local"
	PruneScopeRemote PruneScopeType = "remote"
	PruneScopeAll    PruneScopeType = "all"
	PruneScopeStale  PruneScopeType = "stale"
)

// IsValid reports whether the prune scope is a recognized enum value.
func (s PruneScopeType) IsValid() bool {
	switch s {
	case PruneScopeLocal, PruneScopeRemote, PruneScopeAll, PruneScopeStale:
		return true
	default:
		return false
	}
}

// String returns the string representation of the prune scope.
func (s PruneScopeType) String() string {
	return string(s)
}

// ParsePruneScopeType parses a raw string into a PruneScopeType.
func ParsePruneScopeType(s string) (PruneScopeType, *apperror.AppError) {
	val := PruneScopeType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParsePruneScopeType", "E1000", fmt.Sprintf("invalid prune scope: %q", s)).WithContext("value", s)
}

// AllPruneScopeTypes returns all supported PruneScopeType values.
func AllPruneScopeTypes() []PruneScopeType {
	return []PruneScopeType{
		PruneScopeLocal,
		PruneScopeRemote,
		PruneScopeAll,
		PruneScopeStale,
	}
}

// IsLocal reports whether the scope targets local resources.
func (s PruneScopeType) IsLocal() bool {
	return s == PruneScopeLocal
}

// IsRemote reports whether the scope targets remote resources.
func (s PruneScopeType) IsRemote() bool {
	return s == PruneScopeRemote
}

// IsAll reports whether the scope targets both local and remote resources.
func (s PruneScopeType) IsAll() bool {
	return s == PruneScopeAll
}

// IsStale reports whether the scope targets stale resources.
func (s PruneScopeType) IsStale() bool {
	return s == PruneScopeStale
}

// HasScope reports whether the scope is valid and recognized.
func (s PruneScopeType) HasScope() bool {
	return s.IsValid()
}

// PruneTargetType specifies the target resource type to prune.
type PruneTargetType string

const (
	PruneTargetBranches  PruneTargetType = "branches"
	PruneTargetTags      PruneTargetType = "tags"
	PruneTargetRemotes   PruneTargetType = "remotes"
	PruneTargetObjects   PruneTargetType = "objects"
	PruneTargetWorktrees PruneTargetType = "worktrees"
)

// IsValid reports whether the prune target is a recognized enum value.
func (t PruneTargetType) IsValid() bool {
	switch t {
	case PruneTargetBranches, PruneTargetTags, PruneTargetRemotes,
		PruneTargetObjects, PruneTargetWorktrees:
		return true
	default:
		return false
	}
}

// String returns the string representation of the prune target.
func (t PruneTargetType) String() string {
	return string(t)
}

// ParsePruneTargetType parses a raw string into a PruneTargetType.
func ParsePruneTargetType(s string) (PruneTargetType, *apperror.AppError) {
	val := PruneTargetType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParsePruneTargetType", "E1000", fmt.Sprintf("invalid prune target: %q", s)).WithContext("value", s)
}

// AllPruneTargetTypes returns all supported PruneTargetType values.
func AllPruneTargetTypes() []PruneTargetType {
	return []PruneTargetType{
		PruneTargetBranches,
		PruneTargetTags,
		PruneTargetRemotes,
		PruneTargetObjects,
		PruneTargetWorktrees,
	}
}

// IsBranches reports whether the target is branches.
func (t PruneTargetType) IsBranches() bool {
	return t == PruneTargetBranches
}

// IsTags reports whether the target is tags.
func (t PruneTargetType) IsTags() bool {
	return t == PruneTargetTags
}

// IsRemotes reports whether the target is remotes.
func (t PruneTargetType) IsRemotes() bool {
	return t == PruneTargetRemotes
}

// IsObjects reports whether the target is git objects.
func (t PruneTargetType) IsObjects() bool {
	return t == PruneTargetObjects
}

// IsWorktrees reports whether the target is worktrees.
func (t PruneTargetType) IsWorktrees() bool {
	return t == PruneTargetWorktrees
}

// HasTarget reports whether the target is valid and recognized.
func (t PruneTargetType) HasTarget() bool {
	return t.IsValid()
}

// PruneActionType specifies an execution action for pruning.
type PruneActionType string

const (
	PruneActionDelete PruneActionType = "delete"
	PruneActionDryRun PruneActionType = "dry-run"
	PruneActionPrompt PruneActionType = "prompt"
	PruneActionKeep   PruneActionType = "keep"
	PruneActionSkip   PruneActionType = "skip"
)

// IsValid reports whether the prune action is a recognized enum value.
func (a PruneActionType) IsValid() bool {
	switch a {
	case PruneActionDelete, PruneActionDryRun, PruneActionPrompt,
		PruneActionKeep, PruneActionSkip:
		return true
	default:
		return false
	}
}

// String returns the string representation of the prune action.
func (a PruneActionType) String() string {
	return string(a)
}

// ParsePruneActionType parses a raw string into a PruneActionType.
func ParsePruneActionType(s string) (PruneActionType, *apperror.AppError) {
	val := PruneActionType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParsePruneActionType", "E1000", fmt.Sprintf("invalid prune action: %q", s)).WithContext("value", s)
}

// AllPruneActionTypes returns all supported PruneActionType values.
func AllPruneActionTypes() []PruneActionType {
	return []PruneActionType{
		PruneActionDelete,
		PruneActionDryRun,
		PruneActionPrompt,
		PruneActionKeep,
		PruneActionSkip,
	}
}

// IsDelete reports whether the action is to delete.
func (a PruneActionType) IsDelete() bool {
	return a == PruneActionDelete
}

// IsDryRun reports whether the action is dry-run mode.
func (a PruneActionType) IsDryRun() bool {
	return a == PruneActionDryRun
}

// IsPrompt reports whether the action prompts the user for confirmation.
func (a PruneActionType) IsPrompt() bool {
	return a == PruneActionPrompt
}

// IsKeep reports whether the action is to keep candidates.
func (a PruneActionType) IsKeep() bool {
	return a == PruneActionKeep
}

// IsSkip reports whether the action is to skip candidates.
func (a PruneActionType) IsSkip() bool {
	return a == PruneActionSkip
}

// HasAction reports whether the action is valid and recognized.
func (a PruneActionType) HasAction() bool {
	return a.IsValid()
}

// PruneStatusType specifies the outcome status of a prune candidate.
type PruneStatusType string

const (
	PruneStatusDeleted PruneStatusType = "deleted"
	PruneStatusKept    PruneStatusType = "kept"
	PruneStatusSkipped PruneStatusType = "skipped"
	PruneStatusFailed  PruneStatusType = "failed"
	PruneStatusAborted PruneStatusType = "aborted"
)

// IsValid reports whether the prune status is a recognized enum value.
func (s PruneStatusType) IsValid() bool {
	switch s {
	case PruneStatusDeleted, PruneStatusKept, PruneStatusSkipped,
		PruneStatusFailed, PruneStatusAborted:
		return true
	default:
		return false
	}
}

// String returns the string representation of the prune status.
func (s PruneStatusType) String() string {
	return string(s)
}

// ParsePruneStatusType parses a raw string into a PruneStatusType.
func ParsePruneStatusType(s string) (PruneStatusType, *apperror.AppError) {
	val := PruneStatusType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParsePruneStatusType", "E1000", fmt.Sprintf("invalid prune status: %q", s)).WithContext("value", s)
}

// AllPruneStatusTypes returns all supported PruneStatusType values.
func AllPruneStatusTypes() []PruneStatusType {
	return []PruneStatusType{
		PruneStatusDeleted,
		PruneStatusKept,
		PruneStatusSkipped,
		PruneStatusFailed,
		PruneStatusAborted,
	}
}

// IsDeleted reports whether candidate was deleted.
func (s PruneStatusType) IsDeleted() bool {
	return s == PruneStatusDeleted
}

// IsKept reports whether candidate was kept.
func (s PruneStatusType) IsKept() bool {
	return s == PruneStatusKept
}

// IsSkipped reports whether candidate was skipped.
func (s PruneStatusType) IsSkipped() bool {
	return s == PruneStatusSkipped
}

// IsFailed reports whether candidate pruning failed.
func (s PruneStatusType) IsFailed() bool {
	return s == PruneStatusFailed
}

// IsAborted reports whether the prune operation was aborted.
func (s PruneStatusType) IsAborted() bool {
	return s == PruneStatusAborted
}

// HasStatus reports whether the status is valid and recognized.
func (s PruneStatusType) HasStatus() bool {
	return s.IsValid()
}
