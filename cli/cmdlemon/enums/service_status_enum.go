package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// ServiceStatusType is a strongly-typed enum.
type ServiceStatusType string

const (
	ServiceStatusTypeRunning ServiceStatusType = "running"
	ServiceStatusTypeStopped ServiceStatusType = "stopped"
	ServiceStatusTypeUnknown ServiceStatusType = "unknown"
)

// IsValid reports whether the value matches a known variant.
func (s ServiceStatusType) IsValid() bool {
	switch s {
	case ServiceStatusTypeRunning, ServiceStatusTypeStopped, ServiceStatusTypeUnknown:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (s ServiceStatusType) String() string {
	return string(s)
}

// ParseServiceStatusType parses a raw string into a validated ServiceStatusType.
func ParseServiceStatusType(s string) (ServiceStatusType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllServiceStatusTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "servicestatustyperunning":
		return ServiceStatusTypeRunning, nil
	case "running":
		return ServiceStatusTypeRunning, nil
	case "servicestatustypestopped":
		return ServiceStatusTypeStopped, nil
	case "stopped":
		return ServiceStatusTypeStopped, nil
	case "servicestatustypeunknown":
		return ServiceStatusTypeUnknown, nil
	case "unknown":
		return ServiceStatusTypeUnknown, nil
	}
	return "", apperror.NewValidation("ParseServiceStatusType", "E1000", fmt.Sprintf("invalid servicestatustype: %q", s)).WithContext("value", s)
}

// AllServiceStatusTypes returns all valid variants.
func AllServiceStatusTypes() []ServiceStatusType {
	return []ServiceStatusType{
		ServiceStatusTypeRunning,
		ServiceStatusTypeStopped,
		ServiceStatusTypeUnknown,
	}
}

// IsRunning reports whether the value is ServiceStatusTypeRunning.
func (s ServiceStatusType) IsRunning() bool {
	return s == ServiceStatusTypeRunning
}

// IsStopped reports whether the value is ServiceStatusTypeStopped.
func (s ServiceStatusType) IsStopped() bool {
	return s == ServiceStatusTypeStopped
}

// IsUnknown reports whether the value is ServiceStatusTypeUnknown.
func (s ServiceStatusType) IsUnknown() bool {
	return s == ServiceStatusTypeUnknown
}

// HasStatus reports whether the value is non-empty and valid.
func (s ServiceStatusType) HasStatus() bool {
	return s.IsValid()
}
