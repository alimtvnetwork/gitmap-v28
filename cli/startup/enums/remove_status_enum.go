package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strconv"
	"strings"
)

// RemoveStatusType enumerates the four mutually-exclusive outcomes.
type RemoveStatusType int

const (
	// RemoveDeleted = file existed, was gitmap-managed, was unlinked
	// (or, under DryRun, would have been unlinked).
	RemoveDeleted RemoveStatusType = 0
	// RemoveNoOp = no file by that name in the autostart dir.
	RemoveNoOp RemoveStatusType = 1
	// RemoveRefused = file exists but is not gitmap-managed.
	RemoveRefused RemoveStatusType = 2
	// RemoveBadName = input failed validation (empty / separator).
	RemoveBadName RemoveStatusType = 3
)

// IsValid reports whether the value matches a known variant.
func (r RemoveStatusType) IsValid() bool {
	switch r {
	case RemoveDeleted, RemoveNoOp, RemoveRefused, RemoveBadName:
		return true
	default:
		return false
	}
}

// String returns the numeric representation of the value.
func (r RemoveStatusType) String() string {
	return strconv.Itoa(int(r))
}

// ParseRemoveStatusType parses a raw string into a validated RemoveStatusType.
func ParseRemoveStatusType(s string) (RemoveStatusType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	if n, err := strconv.Atoi(trimmed); err == nil {
		val := RemoveStatusType(n)
		if val.IsValid() {
			return val, nil
		}
	}
	for _, val := range AllRemoveStatusTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "removedeleted":
		return RemoveDeleted, nil
	case "deleted":
		return RemoveDeleted, nil
	case "removenoop":
		return RemoveNoOp, nil
	case "noop":
		return RemoveNoOp, nil
	case "removerefused":
		return RemoveRefused, nil
	case "refused":
		return RemoveRefused, nil
	case "removebadname":
		return RemoveBadName, nil
	case "badname":
		return RemoveBadName, nil
	}
	return 0, apperror.NewValidation("ParseRemoveStatusType", "E1000", fmt.Sprintf("invalid removestatustype: %q", s)).WithContext("value", s)
}

// AllRemoveStatusTypes returns all valid variants.
func AllRemoveStatusTypes() []RemoveStatusType {
	return []RemoveStatusType{
		RemoveDeleted,
		RemoveNoOp,
		RemoveRefused,
		RemoveBadName,
	}
}

// IsDeleted reports whether the value is RemoveDeleted.
func (r RemoveStatusType) IsDeleted() bool {
	return r == RemoveDeleted
}

// IsNoOp reports whether the value is RemoveNoOp.
func (r RemoveStatusType) IsNoOp() bool {
	return r == RemoveNoOp
}

// IsRefused reports whether the value is RemoveRefused.
func (r RemoveStatusType) IsRefused() bool {
	return r == RemoveRefused
}

// IsBadName reports whether the value is RemoveBadName.
func (r RemoveStatusType) IsBadName() bool {
	return r == RemoveBadName
}

// HasStatus reports whether the value is non-empty and valid.
func (r RemoveStatusType) HasStatus() bool {
	return r.IsValid()
}
