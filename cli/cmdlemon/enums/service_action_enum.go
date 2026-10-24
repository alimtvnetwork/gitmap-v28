package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// ServiceActionType is a strongly-typed enum.
type ServiceActionType string

const (
	ServiceActionTypeInstall   ServiceActionType = "install"
	ServiceActionTypeUninstall ServiceActionType = "uninstall"
	ServiceActionTypeStart     ServiceActionType = "start"
	ServiceActionTypeStop      ServiceActionType = "stop"
	ServiceActionTypeRestart   ServiceActionType = "restart"
	ServiceActionTypeStatus    ServiceActionType = "status"
)

// IsValid reports whether the value matches a known variant.
func (s ServiceActionType) IsValid() bool {
	switch s {
	case ServiceActionTypeInstall, ServiceActionTypeUninstall, ServiceActionTypeStart, ServiceActionTypeStop, ServiceActionTypeRestart, ServiceActionTypeStatus:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (s ServiceActionType) String() string {
	return string(s)
}

// ParseServiceActionType parses a raw string into a validated ServiceActionType.
func ParseServiceActionType(s string) (ServiceActionType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllServiceActionTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "serviceactiontypeinstall":
		return ServiceActionTypeInstall, nil
	case "install":
		return ServiceActionTypeInstall, nil
	case "serviceactiontypeuninstall":
		return ServiceActionTypeUninstall, nil
	case "uninstall":
		return ServiceActionTypeUninstall, nil
	case "serviceactiontypestart":
		return ServiceActionTypeStart, nil
	case "start":
		return ServiceActionTypeStart, nil
	case "serviceactiontypestop":
		return ServiceActionTypeStop, nil
	case "stop":
		return ServiceActionTypeStop, nil
	case "serviceactiontyperestart":
		return ServiceActionTypeRestart, nil
	case "restart":
		return ServiceActionTypeRestart, nil
	case "serviceactiontypestatus":
		return ServiceActionTypeStatus, nil
	case "status":
		return ServiceActionTypeStatus, nil
	}
	return "", apperror.NewValidation("ParseServiceActionType", "E1000", fmt.Sprintf("invalid serviceactiontype: %q", s)).WithContext("value", s)
}

// AllServiceActionTypes returns all valid variants.
func AllServiceActionTypes() []ServiceActionType {
	return []ServiceActionType{
		ServiceActionTypeInstall,
		ServiceActionTypeUninstall,
		ServiceActionTypeStart,
		ServiceActionTypeStop,
		ServiceActionTypeRestart,
		ServiceActionTypeStatus,
	}
}

// IsInstall reports whether the value is ServiceActionTypeInstall.
func (s ServiceActionType) IsInstall() bool {
	return s == ServiceActionTypeInstall
}

// IsUninstall reports whether the value is ServiceActionTypeUninstall.
func (s ServiceActionType) IsUninstall() bool {
	return s == ServiceActionTypeUninstall
}

// IsStart reports whether the value is ServiceActionTypeStart.
func (s ServiceActionType) IsStart() bool {
	return s == ServiceActionTypeStart
}

// IsStop reports whether the value is ServiceActionTypeStop.
func (s ServiceActionType) IsStop() bool {
	return s == ServiceActionTypeStop
}

// IsRestart reports whether the value is ServiceActionTypeRestart.
func (s ServiceActionType) IsRestart() bool {
	return s == ServiceActionTypeRestart
}

// IsStatus reports whether the value is ServiceActionTypeStatus.
func (s ServiceActionType) IsStatus() bool {
	return s == ServiceActionTypeStatus
}

// HasAction reports whether the value is non-empty and valid.
func (s ServiceActionType) HasAction() bool {
	return s.IsValid()
}
