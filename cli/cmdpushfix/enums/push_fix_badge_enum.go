package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// PushFixBadgeType designates the status badge state.
type PushFixBadgeType string

const (
	BadgeDiagnosing PushFixBadgeType = "DIAGNOSING"
	BadgeAuthError  PushFixBadgeType = "AUTH_ERROR"
	BadgeDiverged   PushFixBadgeType = "DIVERGED"
	BadgeTransport  PushFixBadgeType = "TRANSPORT"
	BadgeLeaseForce PushFixBadgeType = "LEASE_FORCE"
	BadgeDryRun     PushFixBadgeType = "DRY_RUN"
	BadgeSuccess    PushFixBadgeType = "SUCCESS"
)

// IsValid reports whether the value matches a known variant.
func (p PushFixBadgeType) IsValid() bool {
	switch p {
	case BadgeDiagnosing, BadgeAuthError, BadgeDiverged, BadgeTransport, BadgeLeaseForce, BadgeDryRun, BadgeSuccess:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (p PushFixBadgeType) String() string {
	return string(p)
}

// ParsePushFixBadgeType parses a raw string into a validated PushFixBadgeType.
func ParsePushFixBadgeType(s string) (PushFixBadgeType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllPushFixBadgeTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "badgediagnosing":
		return BadgeDiagnosing, nil
	case "diagnosing":
		return BadgeDiagnosing, nil
	case "badgeautherror":
		return BadgeAuthError, nil
	case "autherror":
		return BadgeAuthError, nil
	case "badgediverged":
		return BadgeDiverged, nil
	case "diverged":
		return BadgeDiverged, nil
	case "badgetransport":
		return BadgeTransport, nil
	case "transport":
		return BadgeTransport, nil
	case "badgeleaseforce":
		return BadgeLeaseForce, nil
	case "leaseforce":
		return BadgeLeaseForce, nil
	case "badgedryrun":
		return BadgeDryRun, nil
	case "dryrun":
		return BadgeDryRun, nil
	case "badgesuccess":
		return BadgeSuccess, nil
	case "success":
		return BadgeSuccess, nil
	}
	return "", apperror.NewValidation("ParsePushFixBadgeType", "E1000", fmt.Sprintf("invalid pushfixbadgetype: %q", s)).WithContext("value", s)
}

// AllPushFixBadgeTypes returns all valid variants.
func AllPushFixBadgeTypes() []PushFixBadgeType {
	return []PushFixBadgeType{
		BadgeDiagnosing,
		BadgeAuthError,
		BadgeDiverged,
		BadgeTransport,
		BadgeLeaseForce,
		BadgeDryRun,
		BadgeSuccess,
	}
}

// IsDiagnosing reports whether the value is BadgeDiagnosing.
func (p PushFixBadgeType) IsDiagnosing() bool {
	return p == BadgeDiagnosing
}

// IsAuthError reports whether the value is BadgeAuthError.
func (p PushFixBadgeType) IsAuthError() bool {
	return p == BadgeAuthError
}

// IsDiverged reports whether the value is BadgeDiverged.
func (p PushFixBadgeType) IsDiverged() bool {
	return p == BadgeDiverged
}

// IsTransport reports whether the value is BadgeTransport.
func (p PushFixBadgeType) IsTransport() bool {
	return p == BadgeTransport
}

// IsLeaseForce reports whether the value is BadgeLeaseForce.
func (p PushFixBadgeType) IsLeaseForce() bool {
	return p == BadgeLeaseForce
}

// IsDryRun reports whether the value is BadgeDryRun.
func (p PushFixBadgeType) IsDryRun() bool {
	return p == BadgeDryRun
}

// IsSuccess reports whether the value is BadgeSuccess.
func (p PushFixBadgeType) IsSuccess() bool {
	return p == BadgeSuccess
}

// HasFixBadge reports whether the value is non-empty and valid.
func (p PushFixBadgeType) HasFixBadge() bool {
	return p.IsValid()
}
