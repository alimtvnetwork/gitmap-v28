package enums

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// UndoModeType specifies the execution mode for undo operations.
type UndoModeType string

const (
	UndoModeRestore UndoModeType = "restore"
	UndoModeList    UndoModeType = "list"
	UndoModeDryRun  UndoModeType = "dry-run"
	UndoModeFlatten UndoModeType = "flatten"
)

// IsValid reports whether the undo mode is a recognized enum value.
func (m UndoModeType) IsValid() bool {
	switch m {
	case UndoModeRestore, UndoModeList, UndoModeDryRun, UndoModeFlatten:
		return true
	default:
		return false
	}
}

// String returns the string representation of the mode.
func (m UndoModeType) String() string {
	return string(m)
}

// ParseUndoModeType parses a raw string into an UndoModeType.
func ParseUndoModeType(s string) (UndoModeType, *apperror.AppError) {
	val := UndoModeType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseUndoModeType", "E1000", fmt.Sprintf("invalid undo mode: %q", s)).WithContext("value", s)
}

// AllUndoModeTypes returns all supported UndoModeType values.
func AllUndoModeTypes() []UndoModeType {
	return []UndoModeType{
		UndoModeRestore,
		UndoModeList,
		UndoModeDryRun,
		UndoModeFlatten,
	}
}

// IsRestore reports whether the mode is restore.
func (m UndoModeType) IsRestore() bool {
	return m == UndoModeRestore
}

// IsList reports whether the mode is list.
func (m UndoModeType) IsList() bool {
	return m == UndoModeList
}

// IsDryRun reports whether the mode is dry-run.
func (m UndoModeType) IsDryRun() bool {
	return m == UndoModeDryRun
}

// IsFlatten reports whether the mode is flatten.
func (m UndoModeType) IsFlatten() bool {
	return m == UndoModeFlatten
}

// HasMode reports whether the mode is valid and recognized.
func (m UndoModeType) HasMode() bool {
	return m.IsValid()
}

// UndoActionType specifies the action/status of a restored item.
type UndoActionType string

const (
	UndoActionRestored UndoActionType = "restored"
	UndoActionDryRun   UndoActionType = "[dry-run]"
	UndoActionFailed   UndoActionType = "failed"
)

// IsValid reports whether the undo action is a recognized enum value.
func (a UndoActionType) IsValid() bool {
	switch a {
	case UndoActionRestored, UndoActionDryRun, UndoActionFailed:
		return true
	default:
		return false
	}
}

// String returns the string representation of the action.
func (a UndoActionType) String() string {
	return string(a)
}

// ParseUndoActionType parses a raw string into an UndoActionType.
func ParseUndoActionType(s string) (UndoActionType, *apperror.AppError) {
	val := UndoActionType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseUndoActionType", "E1000", fmt.Sprintf("invalid undo action: %q", s)).WithContext("value", s)
}

// AllUndoActionTypes returns all supported UndoActionType values.
func AllUndoActionTypes() []UndoActionType {
	return []UndoActionType{
		UndoActionRestored,
		UndoActionDryRun,
		UndoActionFailed,
	}
}

// IsRestored reports whether the entry was restored.
func (a UndoActionType) IsRestored() bool {
	return a == UndoActionRestored
}

// IsDryRun reports whether the action was simulated.
func (a UndoActionType) IsDryRun() bool {
	return a == UndoActionDryRun
}

// IsFailed reports whether the entry restoration failed.
func (a UndoActionType) IsFailed() bool {
	return a == UndoActionFailed
}

// HasAction reports whether the action is valid and recognized.
func (a UndoActionType) HasAction() bool {
	return a.IsValid()
}

// UndoExitCodeType specifies exit codes for undo operations.
type UndoExitCodeType int

const (
	UndoExitSuccess     UndoExitCodeType = 0
	UndoExitGeneral     UndoExitCodeType = 1
	UndoExitBadFlag     UndoExitCodeType = 6
	UndoExitWriteFailed UndoExitCodeType = 7
	UndoExitBadConfig   UndoExitCodeType = 8
)

// IsValid reports whether the exit code is a recognized enum value.
func (c UndoExitCodeType) IsValid() bool {
	switch c {
	case UndoExitSuccess, UndoExitGeneral, UndoExitBadFlag, UndoExitWriteFailed, UndoExitBadConfig:
		return true
	default:
		return false
	}
}

// String returns the string representation of the exit code.
func (c UndoExitCodeType) String() string {
	switch c {
	case UndoExitSuccess:
		return "success"
	case UndoExitGeneral:
		return "general"
	case UndoExitBadFlag:
		return "bad-flag"
	case UndoExitWriteFailed:
		return "write-failed"
	case UndoExitBadConfig:
		return "bad-config"
	default:
		return strconv.Itoa(int(c))
	}
}

// ParseUndoExitCodeType parses a raw string or integer representation into an UndoExitCodeType.
func ParseUndoExitCodeType(s string) (UndoExitCodeType, *apperror.AppError) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "0", "success":
		return UndoExitSuccess, nil
	case "1", "general", "error":
		return UndoExitGeneral, nil
	case "6", "bad-flag", "badflag":
		return UndoExitBadFlag, nil
	case "7", "write-failed", "writefailed":
		return UndoExitWriteFailed, nil
	case "8", "bad-config", "badconfig":
		return UndoExitBadConfig, nil
	default:
		return UndoExitGeneral, apperror.NewValidation("ParseUndoExitCodeType", "E1000", fmt.Sprintf("invalid undo exit code: %q", s)).WithContext("value", s)
	}
}

// AllUndoExitCodeTypes returns all supported UndoExitCodeType values.
func AllUndoExitCodeTypes() []UndoExitCodeType {
	return []UndoExitCodeType{
		UndoExitSuccess,
		UndoExitGeneral,
		UndoExitBadFlag,
		UndoExitWriteFailed,
		UndoExitBadConfig,
	}
}

// Code returns the integer exit code value.
func (c UndoExitCodeType) Code() int {
	return int(c)
}

// IsSuccess reports whether the exit code represents success.
func (c UndoExitCodeType) IsSuccess() bool {
	return c == UndoExitSuccess
}

// IsGeneral reports whether the exit code represents general error.
func (c UndoExitCodeType) IsGeneral() bool {
	return c == UndoExitGeneral
}

// IsBadFlag reports whether the exit code represents bad flag.
func (c UndoExitCodeType) IsBadFlag() bool {
	return c == UndoExitBadFlag
}

// IsWriteFailed reports whether the exit code represents write failed.
func (c UndoExitCodeType) IsWriteFailed() bool {
	return c == UndoExitWriteFailed
}

// IsBadConfig reports whether the exit code represents bad config.
func (c UndoExitCodeType) IsBadConfig() bool {
	return c == UndoExitBadConfig
}

// HasExitCode reports whether the exit code is valid and recognized.
func (c UndoExitCodeType) HasExitCode() bool {
	return c.IsValid()
}
