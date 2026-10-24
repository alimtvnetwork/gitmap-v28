package enums

import (
	"fmt"
	"strings"
)

// AddTargetType specifies the target resource type for the add command.
type AddTargetType string

const (
	AddTargetRepo         AddTargetType = "repo"
	AddTargetFolder       AddTargetType = "folder"
	AddTargetLfs          AddTargetType = "lfs"
	AddTargetHook         AddTargetType = "hook"
	AddTargetCommonAttr   AddTargetType = "common-attr"
	AddTargetCommonIgnore AddTargetType = "common-ignore"
	AddTargetAttributes   AddTargetType = "attributes"
	AddTargetIgnore       AddTargetType = "ignore"
	AddTargetLfsInstall   AddTargetType = "lfs-install"
)

// IsValid reports whether the add target is a recognized enum value.
func (t AddTargetType) IsValid() bool {
	switch t {
	case AddTargetRepo, AddTargetFolder, AddTargetLfs, AddTargetHook,
		AddTargetCommonAttr, AddTargetCommonIgnore, AddTargetAttributes,
		AddTargetIgnore, AddTargetLfsInstall:
		return true
	default:
		return false
	}
}

// String returns the string representation of the target.
func (t AddTargetType) String() string {
	return string(t)
}

// ParseAddTargetType parses a raw string into an AddTargetType.
func ParseAddTargetType(s string) (AddTargetType, error) {
	val := AddTargetType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", fmt.Errorf("invalid add target: %q", s)
}

// AllAddTargetTypes returns all supported AddTargetType values.
func AllAddTargetTypes() []AddTargetType {
	return []AddTargetType{
		AddTargetRepo,
		AddTargetFolder,
		AddTargetLfs,
		AddTargetHook,
		AddTargetCommonAttr,
		AddTargetCommonIgnore,
		AddTargetAttributes,
		AddTargetIgnore,
		AddTargetLfsInstall,
	}
}

// IsRepo reports whether the target is a repository.
func (t AddTargetType) IsRepo() bool {
	return t == AddTargetRepo
}

// IsFolder reports whether the target is a folder.
func (t AddTargetType) IsFolder() bool {
	return t == AddTargetFolder
}

// IsLfs reports whether the target is LFS tracking.
func (t AddTargetType) IsLfs() bool {
	return t == AddTargetLfs
}

// IsHook reports whether the target is a Git hook.
func (t AddTargetType) IsHook() bool {
	return t == AddTargetHook
}

// IsCommonAttr reports whether the target is common gitattributes.
func (t AddTargetType) IsCommonAttr() bool {
	return t == AddTargetCommonAttr
}

// IsCommonIgnore reports whether the target is common gitignore.
func (t AddTargetType) IsCommonIgnore() bool {
	return t == AddTargetCommonIgnore
}

// IsAttributes reports whether the target is gitattributes.
func (t AddTargetType) IsAttributes() bool {
	return t == AddTargetAttributes
}

// IsIgnore reports whether the target is gitignore.
func (t AddTargetType) IsIgnore() bool {
	return t == AddTargetIgnore
}

// IsLfsInstall reports whether the target is LFS installation.
func (t AddTargetType) IsLfsInstall() bool {
	return t == AddTargetLfsInstall
}

// HasTarget reports whether the target is valid and recognized.
func (t AddTargetType) HasTarget() bool {
	return t.IsValid()
}

// AddActionType specifies the action performed during add operations.
type AddActionType string

const (
	AddActionCreate    AddActionType = "create"
	AddActionAppend    AddActionType = "append"
	AddActionMerge     AddActionType = "merge"
	AddActionOverwrite AddActionType = "overwrite"
	AddActionDiff      AddActionType = "diff"
)

// IsValid reports whether the add action is a recognized enum value.
func (a AddActionType) IsValid() bool {
	switch a {
	case AddActionCreate, AddActionAppend, AddActionMerge,
		AddActionOverwrite, AddActionDiff:
		return true
	default:
		return false
	}
}

// String returns the string representation of the action.
func (a AddActionType) String() string {
	return string(a)
}

// ParseAddActionType parses a raw string into an AddActionType.
func ParseAddActionType(s string) (AddActionType, error) {
	val := AddActionType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", fmt.Errorf("invalid add action: %q", s)
}

// AllAddActionTypes returns all supported AddActionType values.
func AllAddActionTypes() []AddActionType {
	return []AddActionType{
		AddActionCreate,
		AddActionAppend,
		AddActionMerge,
		AddActionOverwrite,
		AddActionDiff,
	}
}

// IsCreate reports whether the action creates a new file or configuration.
func (a AddActionType) IsCreate() bool {
	return a == AddActionCreate
}

// IsAppend reports whether the action appends to existing content.
func (a AddActionType) IsAppend() bool {
	return a == AddActionAppend
}

// IsMerge reports whether the action merges into existing configuration.
func (a AddActionType) IsMerge() bool {
	return a == AddActionMerge
}

// IsOverwrite reports whether the action overwrites existing content.
func (a AddActionType) IsOverwrite() bool {
	return a == AddActionOverwrite
}

// IsDiff reports whether the action shows a preview diff.
func (a AddActionType) IsDiff() bool {
	return a == AddActionDiff
}

// HasAction reports whether the action is valid and recognized.
func (a AddActionType) HasAction() bool {
	return a.IsValid()
}

// AddModeType specifies the execution mode for add operations.
type AddModeType string

const (
	AddModeNormal AddModeType = "normal"
	AddModeDryRun AddModeType = "dry-run"
	AddModeForce  AddModeType = "force"
)

// IsValid reports whether the add mode is a recognized enum value.
func (m AddModeType) IsValid() bool {
	switch m {
	case AddModeNormal, AddModeDryRun, AddModeForce:
		return true
	default:
		return false
	}
}

// String returns the string representation of the mode.
func (m AddModeType) String() string {
	return string(m)
}

// ParseAddModeType parses a raw string into an AddModeType.
func ParseAddModeType(s string) (AddModeType, error) {
	val := AddModeType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", fmt.Errorf("invalid add mode: %q", s)
}

// AllAddModeTypes returns all supported AddModeType values.
func AllAddModeTypes() []AddModeType {
	return []AddModeType{
		AddModeNormal,
		AddModeDryRun,
		AddModeForce,
	}
}

// IsNormal reports whether the mode is standard execution.
func (m AddModeType) IsNormal() bool {
	return m == AddModeNormal
}

// IsDryRun reports whether the mode is preview without modifying files.
func (m AddModeType) IsDryRun() bool {
	return m == AddModeDryRun
}

// IsForce reports whether the mode forces overwriting without confirmation.
func (m AddModeType) IsForce() bool {
	return m == AddModeForce
}

// HasMode reports whether the mode is valid and recognized.
func (m AddModeType) HasMode() bool {
	return m.IsValid()
}
