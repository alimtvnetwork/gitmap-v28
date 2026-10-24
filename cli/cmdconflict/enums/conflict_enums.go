package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// ConflictActionType specifies the resolution strategy or action to take.
type ConflictActionType string

const (
	ConflictActionTheirs ConflictActionType = "theirs"
	ConflictActionOurs   ConflictActionType = "ours"
	ConflictActionIDE    ConflictActionType = "ide"
	ConflictActionAbort  ConflictActionType = "abort"
	ConflictActionRerere ConflictActionType = "rerere"
)

// IsValid reports whether the conflict action is a recognized enum value.
func (a ConflictActionType) IsValid() bool {
	switch a {
	case ConflictActionTheirs, ConflictActionOurs, ConflictActionIDE, ConflictActionAbort, ConflictActionRerere:
		return true
	default:
		return false
	}
}

// String returns the string representation of the conflict action.
func (a ConflictActionType) String() string {
	return string(a)
}

// ParseConflictActionType parses a raw string or flag into a ConflictActionType.
func ParseConflictActionType(s string) (ConflictActionType, *apperror.AppError) {
	norm := strings.ToLower(strings.TrimSpace(s))
	if act, hasMatch := matchActionAlias(norm); hasMatch {
		return act, nil
	}

	return "", apperror.NewValidation("ParseConflictActionType", "E1000", fmt.Sprintf("invalid conflict action: %q", s)).WithContext("value", s)
}

func matchActionAlias(s string) (ConflictActionType, bool) {
	if act, hasMatch := matchSideActionAlias(s); hasMatch {
		return act, true
	}

	return matchControlActionAlias(s)
}

func matchSideActionAlias(s string) (ConflictActionType, bool) {
	switch s {
	case "theirs", "t", "--theirs", "-t":
		return ConflictActionTheirs, true
	case "ours", "o", "--ours", "-o":
		return ConflictActionOurs, true
	default:
		return "", false
	}
}

func matchControlActionAlias(s string) (ConflictActionType, bool) {
	switch s {
	case "ide", "code", "resolve", "v", "--ide", "-i", "--code":
		return ConflictActionIDE, true
	case "abort", "a", "--abort", "-a":
		return ConflictActionAbort, true
	case "rerere", "--rerere":
		return ConflictActionRerere, true
	default:
		return "", false
	}
}

// AllConflictActionTypes returns all supported ConflictActionType values.
func AllConflictActionTypes() []ConflictActionType {
	return []ConflictActionType{
		ConflictActionTheirs,
		ConflictActionOurs,
		ConflictActionIDE,
		ConflictActionAbort,
		ConflictActionRerere,
	}
}

// IsTheirs reports whether the action accepts incoming remote changes.
func (a ConflictActionType) IsTheirs() bool {
	return a == ConflictActionTheirs
}

// IsOurs reports whether the action accepts current local changes.
func (a ConflictActionType) IsOurs() bool {
	return a == ConflictActionOurs
}

// IsIDE reports whether the action launches an external 3-way merge editor.
func (a ConflictActionType) IsIDE() bool {
	return a == ConflictActionIDE
}

// IsAbort reports whether the action aborts the active merge/rebase operation.
func (a ConflictActionType) IsAbort() bool {
	return a == ConflictActionAbort
}

// IsRerere reports whether the action activates Git reuse recorded resolution.
func (a ConflictActionType) IsRerere() bool {
	return a == ConflictActionRerere
}

// HasAction reports whether the action is valid and recognized.
func (a ConflictActionType) HasAction() bool {
	return a.IsValid()
}

// ConflictToolType specifies an external editor or IDE tool candidate.
type ConflictToolType string

const (
	ConflictToolCode     ConflictToolType = "code"
	ConflictToolCursor   ConflictToolType = "cursor"
	ConflictToolWindsurf ConflictToolType = "windsurf"
	ConflictToolCodium   ConflictToolType = "codium"
)

// IsValid reports whether the conflict tool is a recognized enum value.
func (t ConflictToolType) IsValid() bool {
	switch t {
	case ConflictToolCode, ConflictToolCursor, ConflictToolWindsurf, ConflictToolCodium:
		return true
	default:
		return false
	}
}

// String returns the string representation of the tool.
func (t ConflictToolType) String() string {
	return string(t)
}

// BinaryName returns the executable CLI command name for the tool.
func (t ConflictToolType) BinaryName() string {
	return string(t)
}

// DisplayName returns the user-friendly name of the editor tool.
func (t ConflictToolType) DisplayName() string {
	switch t {
	case ConflictToolCode:
		return "VS Code"
	case ConflictToolCursor:
		return "Cursor"
	case ConflictToolWindsurf:
		return "Windsurf"
	case ConflictToolCodium:
		return "VSCodium"
	default:
		return string(t)
	}
}

// ParseConflictToolType parses a raw string into a ConflictToolType.
func ParseConflictToolType(s string) (ConflictToolType, *apperror.AppError) {
	norm := strings.ToLower(strings.TrimSpace(s))
	val := ConflictToolType(norm)
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseConflictToolType", "E1000", fmt.Sprintf("invalid conflict tool: %q", s)).WithContext("value", s)
}

// AllConflictToolTypes returns all supported ConflictToolType values.
func AllConflictToolTypes() []ConflictToolType {
	return []ConflictToolType{
		ConflictToolCode,
		ConflictToolCursor,
		ConflictToolWindsurf,
		ConflictToolCodium,
	}
}

// IsCode reports whether the tool is VS Code.
func (t ConflictToolType) IsCode() bool {
	return t == ConflictToolCode
}

// IsCursor reports whether the tool is Cursor.
func (t ConflictToolType) IsCursor() bool {
	return t == ConflictToolCursor
}

// IsWindsurf reports whether the tool is Windsurf.
func (t ConflictToolType) IsWindsurf() bool {
	return t == ConflictToolWindsurf
}

// IsCodium reports whether the tool is VSCodium.
func (t ConflictToolType) IsCodium() bool {
	return t == ConflictToolCodium
}

// HasTool reports whether the tool is valid and recognized.
func (t ConflictToolType) HasTool() bool {
	return t.IsValid()
}

// ConflictStatusCodeType represents Git porcelain unmerged status codes.
type ConflictStatusCodeType string

const (
	ConflictStatusBothModified ConflictStatusCodeType = "UU"
	ConflictStatusBothAdded    ConflictStatusCodeType = "AA"
	ConflictStatusUserDeleted  ConflictStatusCodeType = "UD"
	ConflictStatusDeletedUser  ConflictStatusCodeType = "DU"
	ConflictStatusBothDeleted  ConflictStatusCodeType = "DD"
	ConflictStatusAddedUser    ConflictStatusCodeType = "AU"
	ConflictStatusUserAdded    ConflictStatusCodeType = "UA"
)

// IsValid reports whether the status code is a recognized unmerged code.
func (c ConflictStatusCodeType) IsValid() bool {
	switch c {
	case ConflictStatusBothModified, ConflictStatusBothAdded,
		ConflictStatusUserDeleted, ConflictStatusDeletedUser,
		ConflictStatusBothDeleted, ConflictStatusAddedUser,
		ConflictStatusUserAdded:
		return true
	default:
		return false
	}
}

// String returns the 2-letter status code string.
func (c ConflictStatusCodeType) String() string {
	return string(c)
}

// IsUnmerged reports whether this status code indicates an unmerged conflict.
func (c ConflictStatusCodeType) IsUnmerged() bool {
	return c.IsValid()
}

// Description returns a human-readable description of the conflict state.
func (c ConflictStatusCodeType) Description() string {
	if desc, hasMatch := matchPrimaryDescription(c); hasMatch {
		return desc
	}

	return matchSecondaryDescription(c)
}

func matchPrimaryDescription(c ConflictStatusCodeType) (string, bool) {
	switch c {
	case ConflictStatusBothModified:
		return "both modified", true
	case ConflictStatusBothAdded:
		return "both added", true
	case ConflictStatusBothDeleted:
		return "both deleted", true
	default:
		return "", false
	}
}

func matchSecondaryDescription(c ConflictStatusCodeType) string {
	switch c {
	case ConflictStatusUserDeleted:
		return "updated by us, deleted by them"
	case ConflictStatusDeletedUser:
		return "deleted by us, updated by them"
	case ConflictStatusAddedUser:
		return "added by us, updated by them"
	case ConflictStatusUserAdded:
		return "updated by us, added by them"
	default:
		return "unknown conflict status"
	}
}

// ParseConflictStatusCodeType parses a raw 2-character string into a ConflictStatusCodeType.
func ParseConflictStatusCodeType(s string) (ConflictStatusCodeType, *apperror.AppError) {
	norm := strings.ToUpper(strings.TrimSpace(s))
	val := ConflictStatusCodeType(norm)
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseConflictStatusCodeType", "E1000", fmt.Sprintf("invalid conflict status code: %q", s)).WithContext("value", s)
}

// AllConflictStatusCodeTypes returns all supported ConflictStatusCodeType values.
func AllConflictStatusCodeTypes() []ConflictStatusCodeType {
	return []ConflictStatusCodeType{
		ConflictStatusBothModified,
		ConflictStatusBothAdded,
		ConflictStatusUserDeleted,
		ConflictStatusDeletedUser,
		ConflictStatusBothDeleted,
		ConflictStatusAddedUser,
		ConflictStatusUserAdded,
	}
}

// IsBothModified reports whether both sides modified the file.
func (c ConflictStatusCodeType) IsBothModified() bool {
	return c == ConflictStatusBothModified
}

// IsBothAdded reports whether both sides added the file.
func (c ConflictStatusCodeType) IsBothAdded() bool {
	return c == ConflictStatusBothAdded
}

// IsBothDeleted reports whether both sides deleted the file.
func (c ConflictStatusCodeType) IsBothDeleted() bool {
	return c == ConflictStatusBothDeleted
}

// HasConflict reports whether the status code indicates a conflict.
func (c ConflictStatusCodeType) HasConflict() bool {
	return c.IsValid()
}

// ConflictScopeType specifies the inspection scope for checking conflicts.
type ConflictScopeType string

const (
	ConflictScopeTarget  ConflictScopeType = "target"
	ConflictScopeCurrent ConflictScopeType = "current"
	ConflictScopeAll     ConflictScopeType = "all"
)

// IsValid reports whether the conflict scope is a recognized enum value.
func (s ConflictScopeType) IsValid() bool {
	switch s {
	case ConflictScopeTarget, ConflictScopeCurrent, ConflictScopeAll:
		return true
	default:
		return false
	}
}

// String returns the string representation of the conflict scope.
func (s ConflictScopeType) String() string {
	return string(s)
}

// ParseConflictScopeType parses a raw string into a ConflictScopeType.
func ParseConflictScopeType(s string) (ConflictScopeType, *apperror.AppError) {
	norm := strings.ToLower(strings.TrimSpace(s))
	val := ConflictScopeType(norm)
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseConflictScopeType", "E1000", fmt.Sprintf("invalid conflict scope: %q", s)).WithContext("value", s)
}

// AllConflictScopeTypes returns all supported ConflictScopeType values.
func AllConflictScopeTypes() []ConflictScopeType {
	return []ConflictScopeType{
		ConflictScopeTarget,
		ConflictScopeCurrent,
		ConflictScopeAll,
	}
}

// ResolveConflictScope determines the conflict scope from arguments and environment.
func ResolveConflictScope(target string, isInsideRepo bool) ConflictScopeType {
	if target != "" {
		return ConflictScopeTarget
	}

	if isInsideRepo {
		return ConflictScopeCurrent
	}

	return ConflictScopeAll
}

// IsTarget reports whether the scope is a specifically targeted repository.
func (s ConflictScopeType) IsTarget() bool {
	return s == ConflictScopeTarget
}

// IsCurrent reports whether the scope is the current working directory repository.
func (s ConflictScopeType) IsCurrent() bool {
	return s == ConflictScopeCurrent
}

// IsAll reports whether the scope is all tracked repositories in the store.
func (s ConflictScopeType) IsAll() bool {
	return s == ConflictScopeAll
}

// HasScope reports whether the scope is valid and recognized.
func (s ConflictScopeType) HasScope() bool {
	return s.IsValid()
}
