package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// SyncTargetType represents the synchronization scope filter.
type SyncTargetType string

const (
	SyncTargetAll     SyncTargetType = "all"
	SyncTargetDirty   SyncTargetType = "dirty"
	SyncTargetClean   SyncTargetType = "clean"
	SyncTargetStashed SyncTargetType = "stashed"
	SyncTargetMissing SyncTargetType = "missing"
)

// IsValid reports whether the target matches a recognized sync target variant.
func (t SyncTargetType) IsValid() bool {
	switch t {
	case SyncTargetAll, SyncTargetDirty, SyncTargetClean, SyncTargetStashed, SyncTargetMissing:
		return true
	default:
		return false
	}
}

// String returns the string representation of the sync target.
func (t SyncTargetType) String() string {
	return string(t)
}

// ParseSyncTargetType parses a string into a validated SyncTargetType.
func ParseSyncTargetType(s string) (SyncTargetType, *apperror.AppError) {
	val := SyncTargetType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseSyncTargetType", "E1000", fmt.Sprintf("invalid sync target: %q", s)).WithContext("value", s)
}

// AllSyncTargetTypes returns all valid SyncTargetType values.
func AllSyncTargetTypes() []SyncTargetType {
	return []SyncTargetType{
		SyncTargetAll,
		SyncTargetDirty,
		SyncTargetClean,
		SyncTargetStashed,
		SyncTargetMissing,
	}
}

// IsAll reports whether the target denotes all repositories.
func (t SyncTargetType) IsAll() bool {
	return t == SyncTargetAll
}

// IsDirty reports whether the target denotes dirty repositories.
func (t SyncTargetType) IsDirty() bool {
	return t == SyncTargetDirty
}

// IsClean reports whether the target denotes clean repositories.
func (t SyncTargetType) IsClean() bool {
	return t == SyncTargetClean
}

// IsStashed reports whether the target denotes stashed repositories.
func (t SyncTargetType) IsStashed() bool {
	return t == SyncTargetStashed
}

// IsMissing reports whether the target denotes missing repositories.
func (t SyncTargetType) IsMissing() bool {
	return t == SyncTargetMissing
}

// HasTargetScope reports whether the target has a valid non-empty scope.
func (t SyncTargetType) HasTargetScope() bool {
	return t.IsValid()
}

// SyncRepoStatusType represents the outcome status of a repo sync operation.
type SyncRepoStatusType string

const (
	SyncRepoStatusSynced   SyncRepoStatusType = "synced"
	SyncRepoStatusModified SyncRepoStatusType = "modified"
	SyncRepoStatusConflict SyncRepoStatusType = "conflict"
	SyncRepoStatusSkipped  SyncRepoStatusType = "skipped"
	SyncRepoStatusError    SyncRepoStatusType = "error"
)

// IsValid reports whether the status matches a recognized sync repo status variant.
func (s SyncRepoStatusType) IsValid() bool {
	switch s {
	case SyncRepoStatusSynced, SyncRepoStatusModified, SyncRepoStatusConflict,
		SyncRepoStatusSkipped, SyncRepoStatusError:
		return true
	default:
		return false
	}
}

// String returns the string representation of the sync repo status.
func (s SyncRepoStatusType) String() string {
	return string(s)
}

// ParseSyncRepoStatusType parses a string into a validated SyncRepoStatusType.
func ParseSyncRepoStatusType(s string) (SyncRepoStatusType, *apperror.AppError) {
	val := SyncRepoStatusType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseSyncRepoStatusType", "E1000", fmt.Sprintf("invalid sync repo status: %q", s)).WithContext("value", s)
}

// AllSyncRepoStatusTypes returns all valid SyncRepoStatusType values.
func AllSyncRepoStatusTypes() []SyncRepoStatusType {
	return []SyncRepoStatusType{
		SyncRepoStatusSynced,
		SyncRepoStatusModified,
		SyncRepoStatusConflict,
		SyncRepoStatusSkipped,
		SyncRepoStatusError,
	}
}

// IsSynced reports whether the status denotes a synchronized repository.
func (s SyncRepoStatusType) IsSynced() bool {
	return s == SyncRepoStatusSynced
}

// IsModified reports whether the status denotes modified working tree changes.
func (s SyncRepoStatusType) IsModified() bool {
	return s == SyncRepoStatusModified
}

// IsConflict reports whether the status denotes merge/rebase conflict.
func (s SyncRepoStatusType) IsConflict() bool {
	return s == SyncRepoStatusConflict
}

// IsSkipped reports whether the status denotes a skipped repository.
func (s SyncRepoStatusType) IsSkipped() bool {
	return s == SyncRepoStatusSkipped
}

// IsError reports whether the status denotes an operation error.
func (s SyncRepoStatusType) IsError() bool {
	return s == SyncRepoStatusError
}

// HasConflict reports whether a conflict exists.
func (s SyncRepoStatusType) HasConflict() bool {
	return s == SyncRepoStatusConflict
}

// HasError reports whether an error exists.
func (s SyncRepoStatusType) HasError() bool {
	return s == SyncRepoStatusError
}

// SyncActionType represents discrete synchronization actions.
type SyncActionType string

const (
	SyncActionFetch    SyncActionType = "fetch"
	SyncActionPull     SyncActionType = "pull"
	SyncActionMerge    SyncActionType = "merge"
	SyncActionRebase   SyncActionType = "rebase"
	SyncActionCheckout SyncActionType = "checkout"
)

// IsValid reports whether the action matches a recognized sync action variant.
func (a SyncActionType) IsValid() bool {
	switch a {
	case SyncActionFetch, SyncActionPull, SyncActionMerge, SyncActionRebase, SyncActionCheckout:
		return true
	default:
		return false
	}
}

// String returns the string representation of the sync action.
func (a SyncActionType) String() string {
	return string(a)
}

// ParseSyncActionType parses a string into a validated SyncActionType.
func ParseSyncActionType(s string) (SyncActionType, *apperror.AppError) {
	val := SyncActionType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseSyncActionType", "E1000", fmt.Sprintf("invalid sync action: %q", s)).WithContext("value", s)
}

// AllSyncActionTypes returns all valid SyncActionType values.
func AllSyncActionTypes() []SyncActionType {
	return []SyncActionType{
		SyncActionFetch,
		SyncActionPull,
		SyncActionMerge,
		SyncActionRebase,
		SyncActionCheckout,
	}
}

// IsFetch reports whether the action is fetch.
func (a SyncActionType) IsFetch() bool {
	return a == SyncActionFetch
}

// IsPull reports whether the action is pull.
func (a SyncActionType) IsPull() bool {
	return a == SyncActionPull
}

// IsMerge reports whether the action is merge.
func (a SyncActionType) IsMerge() bool {
	return a == SyncActionMerge
}

// IsRebase reports whether the action is rebase.
func (a SyncActionType) IsRebase() bool {
	return a == SyncActionRebase
}

// IsCheckout reports whether the action is checkout.
func (a SyncActionType) IsCheckout() bool {
	return a == SyncActionCheckout
}

// HasRebase reports whether the action involves rebasing.
func (a SyncActionType) HasRebase() bool {
	return a == SyncActionRebase
}

// SyncModeType represents operational synchronization modes.
type SyncModeType string

const (
	SyncModeFastForward SyncModeType = "fast-forward"
	SyncModeRebase      SyncModeType = "rebase"
	SyncModeDryRun      SyncModeType = "dry-run"
)

// IsValid reports whether the mode matches a recognized sync mode variant.
func (m SyncModeType) IsValid() bool {
	switch m {
	case SyncModeFastForward, SyncModeRebase, SyncModeDryRun:
		return true
	default:
		return false
	}
}

// String returns the string representation of the sync mode.
func (m SyncModeType) String() string {
	return string(m)
}

// ParseSyncModeType parses a string into a validated SyncModeType.
func ParseSyncModeType(s string) (SyncModeType, *apperror.AppError) {
	val := SyncModeType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseSyncModeType", "E1000", fmt.Sprintf("invalid sync mode: %q", s)).WithContext("value", s)
}

// AllSyncModeTypes returns all valid SyncModeType values.
func AllSyncModeTypes() []SyncModeType {
	return []SyncModeType{
		SyncModeFastForward,
		SyncModeRebase,
		SyncModeDryRun,
	}
}

// IsFastForward reports whether the mode is fast-forward.
func (m SyncModeType) IsFastForward() bool {
	return m == SyncModeFastForward
}

// IsRebase reports whether the mode is rebase.
func (m SyncModeType) IsRebase() bool {
	return m == SyncModeRebase
}

// IsDryRun reports whether the mode is dry-run.
func (m SyncModeType) IsDryRun() bool {
	return m == SyncModeDryRun
}

// HasMutation reports whether the mode allows git mutations.
func (m SyncModeType) HasMutation() bool {
	return m != SyncModeDryRun
}
