package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// OSActionType identifies the power operation to schedule.
type OSActionType string

const (
	OSActionShutdown OSActionType = "shutdown"
	OSActionRestart  OSActionType = "restart"
	OSActionCancel   OSActionType = "cancel"
)

// IsValid reports whether the OS action matches a known action.
func (a OSActionType) IsValid() bool {
	switch a {
	case OSActionShutdown, OSActionRestart, OSActionCancel:
		return true
	default:
		return false
	}
}

// String returns the string representation of the OS action.
func (a OSActionType) String() string {
	return string(a)
}

// ParseOSActionType parses a raw string into a validated OSActionType.
func ParseOSActionType(s string) (OSActionType, *apperror.AppError) {
	val := OSActionType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseOSActionType", "E1000", fmt.Sprintf("invalid OS action: %q", s)).WithContext("value", s)
}

// AllOSActionTypes returns all valid OSActionType variants.
func AllOSActionTypes() []OSActionType {
	return []OSActionType{OSActionShutdown, OSActionRestart, OSActionCancel}
}

// IsShutdown reports whether the action is shutdown.
func (a OSActionType) IsShutdown() bool { return a == OSActionShutdown }

// IsRestart reports whether the action is restart.
func (a OSActionType) IsRestart() bool { return a == OSActionRestart }

// IsCancel reports whether the action is cancel.
func (a OSActionType) IsCancel() bool { return a == OSActionCancel }
