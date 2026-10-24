package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// GitHookType represents a supported Git hook trigger name.
type GitHookType string

const (
	GitHookCommitMsg        GitHookType = "commit-msg"
	GitHookPreCommit        GitHookType = "pre-commit"
	GitHookPrePush          GitHookType = "pre-push"
	GitHookPostCheckout     GitHookType = "post-checkout"
	GitHookPostMerge        GitHookType = "post-merge"
	GitHookPreRebase        GitHookType = "pre-rebase"
	GitHookPrepareCommitMsg GitHookType = "prepare-commit-msg"
	GitHookPostCommit       GitHookType = "post-commit"
)

// IsValid reports whether the git hook is a recognized hook type.
func (h GitHookType) IsValid() bool {
	switch h {
	case GitHookCommitMsg, GitHookPreCommit, GitHookPrePush,
		GitHookPostCheckout, GitHookPostMerge, GitHookPreRebase,
		GitHookPrepareCommitMsg, GitHookPostCommit:
		return true
	default:
		return false
	}
}

// String returns the string representation of the hook type.
func (h GitHookType) String() string {
	return string(h)
}

// ParseGitHookType parses a raw string into a GitHookType.
func ParseGitHookType(s string) (GitHookType, *apperror.AppError) {
	val := GitHookType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseGitHookType", "E1000", fmt.Sprintf("invalid git hook type: %q", s)).WithContext("value", s)
}

// AllGitHookTypes returns all supported GitHookType values.
func AllGitHookTypes() []GitHookType {
	return []GitHookType{
		GitHookCommitMsg,
		GitHookPreCommit,
		GitHookPrePush,
		GitHookPostCheckout,
		GitHookPostMerge,
		GitHookPreRebase,
		GitHookPrepareCommitMsg,
		GitHookPostCommit,
	}
}

// IsCommitMsg reports whether the hook is commit-msg.
func (h GitHookType) IsCommitMsg() bool {
	return h == GitHookCommitMsg
}

// IsPreCommit reports whether the hook is pre-commit.
func (h GitHookType) IsPreCommit() bool {
	return h == GitHookPreCommit
}

// IsPrePush reports whether the hook is pre-push.
func (h GitHookType) IsPrePush() bool {
	return h == GitHookPrePush
}

// IsPostCheckout reports whether the hook is post-checkout.
func (h GitHookType) IsPostCheckout() bool {
	return h == GitHookPostCheckout
}

// IsPostMerge reports whether the hook is post-merge.
func (h GitHookType) IsPostMerge() bool {
	return h == GitHookPostMerge
}

// IsPreRebase reports whether the hook is pre-rebase.
func (h GitHookType) IsPreRebase() bool {
	return h == GitHookPreRebase
}

// IsPrepareCommitMsg reports whether the hook is prepare-commit-msg.
func (h GitHookType) IsPrepareCommitMsg() bool {
	return h == GitHookPrepareCommitMsg
}

// IsPostCommit reports whether the hook is post-commit.
func (h GitHookType) IsPostCommit() bool {
	return h == GitHookPostCommit
}

// HasHook reports whether the hook type is valid and non-empty.
func (h GitHookType) HasHook() bool {
	return h.IsValid()
}

// HookActionType specifies a CLI action or subcommand for hooks.
type HookActionType string

const (
	HookActionInstall   HookActionType = "install"
	HookActionUninstall HookActionType = "uninstall"
	HookActionRun       HookActionType = "run"
	HookActionCheck     HookActionType = "check"
	HookActionList      HookActionType = "list"
	HookActionStatus    HookActionType = "status"
	HookActionHelp      HookActionType = "help"
	HookActionHelpShort HookActionType = "-h"
	HookActionHelpLong  HookActionType = "--help"
)

// IsValid reports whether the action is a recognized enum value.
func (a HookActionType) IsValid() bool {
	switch a {
	case HookActionInstall, HookActionUninstall, HookActionRun,
		HookActionCheck, HookActionList, HookActionStatus,
		HookActionHelp, HookActionHelpShort, HookActionHelpLong:
		return true
	default:
		return false
	}
}

// String returns the string representation of the action.
func (a HookActionType) String() string {
	return string(a)
}

// ParseHookActionType parses a raw string into a HookActionType.
func ParseHookActionType(s string) (HookActionType, *apperror.AppError) {
	val := HookActionType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseHookActionType", "E1000", fmt.Sprintf("invalid hook action: %q", s)).WithContext("value", s)
}

// AllHookActionTypes returns all supported HookActionType values.
func AllHookActionTypes() []HookActionType {
	return []HookActionType{
		HookActionInstall,
		HookActionUninstall,
		HookActionRun,
		HookActionCheck,
		HookActionList,
		HookActionStatus,
		HookActionHelp,
		HookActionHelpShort,
		HookActionHelpLong,
	}
}

// IsInstall reports whether the action is install.
func (a HookActionType) IsInstall() bool {
	return a == HookActionInstall
}

// IsUninstall reports whether the action is uninstall.
func (a HookActionType) IsUninstall() bool {
	return a == HookActionUninstall
}

// IsRun reports whether the action is run.
func (a HookActionType) IsRun() bool {
	return a == HookActionRun
}

// IsCheck reports whether the action is check.
func (a HookActionType) IsCheck() bool {
	return a == HookActionCheck
}

// IsList reports whether the action is list.
func (a HookActionType) IsList() bool {
	return a == HookActionList
}

// IsStatus reports whether the action is status.
func (a HookActionType) IsStatus() bool {
	return a == HookActionStatus
}

// IsHelp reports whether the action is a help variant.
func (a HookActionType) IsHelp() bool {
	return a == HookActionHelp || a == HookActionHelpShort || a == HookActionHelpLong
}

// HasAction reports whether the action value is valid and non-empty.
func (a HookActionType) HasAction() bool {
	return a.IsValid()
}

// HookPlanActionType specifies the planned execution action for a repository.
type HookPlanActionType string

const (
	HookPlanActionInstall HookPlanActionType = "install"
	HookPlanActionAlready HookPlanActionType = "already"
	HookPlanActionSkip    HookPlanActionType = "skip"
)

// IsValid reports whether the plan action is recognized.
func (p HookPlanActionType) IsValid() bool {
	switch p {
	case HookPlanActionInstall, HookPlanActionAlready, HookPlanActionSkip:
		return true
	default:
		return false
	}
}

// String returns the string representation of the plan action.
func (p HookPlanActionType) String() string {
	return string(p)
}

// ParseHookPlanActionType parses a raw string into a HookPlanActionType.
func ParseHookPlanActionType(s string) (HookPlanActionType, *apperror.AppError) {
	val := HookPlanActionType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseHookPlanActionType", "E1000", fmt.Sprintf("invalid hook plan action: %q", s)).WithContext("value", s)
}

// AllHookPlanActionTypes returns all supported HookPlanActionType values.
func AllHookPlanActionTypes() []HookPlanActionType {
	return []HookPlanActionType{
		HookPlanActionInstall,
		HookPlanActionAlready,
		HookPlanActionSkip,
	}
}

// IsInstall reports whether the plan action is to install.
func (p HookPlanActionType) IsInstall() bool {
	return p == HookPlanActionInstall
}

// IsAlready reports whether the guard is already installed.
func (p HookPlanActionType) IsAlready() bool {
	return p == HookPlanActionAlready
}

// IsSkip reports whether the repository is skipped.
func (p HookPlanActionType) IsSkip() bool {
	return p == HookPlanActionSkip
}

// HasPlanAction reports whether the plan action is valid.
func (p HookPlanActionType) HasPlanAction() bool {
	return p.IsValid()
}

// HookScopeType specifies the target scope of hook operations.
type HookScopeType string

const (
	HookScopeAll     HookScopeType = "all"
	HookScopeRepo    HookScopeType = "repo"
	HookScopeGlobal  HookScopeType = "global"
	HookScopeWorkdir HookScopeType = "workdir"
)

// IsValid reports whether the scope is recognized.
func (s HookScopeType) IsValid() bool {
	switch s {
	case HookScopeAll, HookScopeRepo, HookScopeGlobal, HookScopeWorkdir:
		return true
	default:
		return false
	}
}

// String returns the string representation of the scope.
func (s HookScopeType) String() string {
	return string(s)
}

// ParseHookScopeType parses a raw string into a HookScopeType.
func ParseHookScopeType(raw string) (HookScopeType, *apperror.AppError) {
	val := HookScopeType(strings.ToLower(strings.TrimSpace(raw)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseHookScopeType", "E1000", fmt.Sprintf("invalid hook scope: %q", raw)).WithContext("value", raw)
}

// AllHookScopeTypes returns all supported HookScopeType values.
func AllHookScopeTypes() []HookScopeType {
	return []HookScopeType{
		HookScopeAll,
		HookScopeRepo,
		HookScopeGlobal,
		HookScopeWorkdir,
	}
}

// IsAll reports whether the scope targets all repositories.
func (s HookScopeType) IsAll() bool {
	return s == HookScopeAll
}

// IsRepo reports whether the scope targets a single repository.
func (s HookScopeType) IsRepo() bool {
	return s == HookScopeRepo
}

// IsGlobal reports whether the scope targets global hooks.
func (s HookScopeType) IsGlobal() bool {
	return s == HookScopeGlobal
}

// IsWorkdir reports whether the scope targets the workdir.
func (s HookScopeType) IsWorkdir() bool {
	return s == HookScopeWorkdir
}

// HasScope reports whether the scope is valid.
func (s HookScopeType) HasScope() bool {
	return s.IsValid()
}
