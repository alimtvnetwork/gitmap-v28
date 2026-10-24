package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strconv"
	"strings"
)

// AddStatusType tags the four mutually-exclusive Add outcomes. Kept
// parallel to RemoveStatus so the CLI rendering layer can switch on
// either with the same shape.
type AddStatusType int

const (
	// AddCreated = file did not exist; was written fresh.
	AddCreated AddStatusType = 0
	// AddOverwritten = previously gitmap-managed file was replaced
	// because Force was set.
	AddOverwritten AddStatusType = 1
	// AddRefused = a non-gitmap-managed file with the same name
	// already exists; we did NOT touch it.
	AddRefused AddStatusType = 2
	// AddBadName = name failed validation (empty / separator / NUL).
	AddBadName AddStatusType = 3
	// AddExists = a gitmap-managed file with the same name already
	// exists and Force was NOT set; nothing was written.
	AddExists AddStatusType = 4
)

// IsValid reports whether the value matches a known variant.
func (a AddStatusType) IsValid() bool {
	switch a {
	case AddCreated, AddOverwritten, AddRefused, AddBadName, AddExists:
		return true
	default:
		return false
	}
}

// String returns the numeric representation of the value.
func (a AddStatusType) String() string {
	return strconv.Itoa(int(a))
}

// ParseAddStatusType parses a raw string into a validated AddStatusType.
func ParseAddStatusType(s string) (AddStatusType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	if n, err := strconv.Atoi(trimmed); err == nil {
		val := AddStatusType(n)
		if val.IsValid() {
			return val, nil
		}
	}
	for _, val := range AllAddStatusTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "addcreated":
		return AddCreated, nil
	case "created":
		return AddCreated, nil
	case "addoverwritten":
		return AddOverwritten, nil
	case "overwritten":
		return AddOverwritten, nil
	case "addrefused":
		return AddRefused, nil
	case "refused":
		return AddRefused, nil
	case "addbadname":
		return AddBadName, nil
	case "badname":
		return AddBadName, nil
	case "addexists":
		return AddExists, nil
	case "exists":
		return AddExists, nil
	}
	return 0, apperror.NewValidation("ParseAddStatusType", "E1000", fmt.Sprintf("invalid addstatustype: %q", s)).WithContext("value", s)
}

// AllAddStatusTypes returns all valid variants.
func AllAddStatusTypes() []AddStatusType {
	return []AddStatusType{
		AddCreated,
		AddOverwritten,
		AddRefused,
		AddBadName,
		AddExists,
	}
}

// IsCreated reports whether the value is AddCreated.
func (a AddStatusType) IsCreated() bool {
	return a == AddCreated
}

// IsOverwritten reports whether the value is AddOverwritten.
func (a AddStatusType) IsOverwritten() bool {
	return a == AddOverwritten
}

// IsRefused reports whether the value is AddRefused.
func (a AddStatusType) IsRefused() bool {
	return a == AddRefused
}

// IsBadName reports whether the value is AddBadName.
func (a AddStatusType) IsBadName() bool {
	return a == AddBadName
}

// IsExists reports whether the value is AddExists.
func (a AddStatusType) IsExists() bool {
	return a == AddExists
}

// HasStatus reports whether the value is non-empty and valid.
func (a AddStatusType) HasStatus() bool {
	return a.IsValid()
}
