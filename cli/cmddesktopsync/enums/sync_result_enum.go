package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strconv"
	"strings"
)

// SyncResultType represents the outcome of syncing one repo.
type SyncResultType int

const (
	SyncResultTypeAdded   SyncResultType = 0
	SyncResultTypeSkipped SyncResultType = 1
	SyncResultTypeFailed  SyncResultType = 2
)

// IsValid reports whether the value matches a known variant.
func (s SyncResultType) IsValid() bool {
	switch s {
	case SyncResultTypeAdded, SyncResultTypeSkipped, SyncResultTypeFailed:
		return true
	default:
		return false
	}
}

// String returns the numeric representation of the value.
func (s SyncResultType) String() string {
	return strconv.Itoa(int(s))
}

// ParseSyncResultType parses a raw string into a validated SyncResultType.
func ParseSyncResultType(s string) (SyncResultType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	if n, err := strconv.Atoi(trimmed); err == nil {
		val := SyncResultType(n)
		if val.IsValid() {
			return val, nil
		}
	}
	for _, val := range AllSyncResultTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "syncresulttypeadded":
		return SyncResultTypeAdded, nil
	case "added":
		return SyncResultTypeAdded, nil
	case "syncresulttypeskipped":
		return SyncResultTypeSkipped, nil
	case "skipped":
		return SyncResultTypeSkipped, nil
	case "syncresulttypefailed":
		return SyncResultTypeFailed, nil
	case "failed":
		return SyncResultTypeFailed, nil
	}
	return 0, apperror.NewValidation("ParseSyncResultType", "E1000", fmt.Sprintf("invalid syncresulttype: %q", s)).WithContext("value", s)
}

// AllSyncResultTypes returns all valid variants.
func AllSyncResultTypes() []SyncResultType {
	return []SyncResultType{
		SyncResultTypeAdded,
		SyncResultTypeSkipped,
		SyncResultTypeFailed,
	}
}

// IsAdded reports whether the value is SyncResultTypeAdded.
func (s SyncResultType) IsAdded() bool {
	return s == SyncResultTypeAdded
}

// IsSkipped reports whether the value is SyncResultTypeSkipped.
func (s SyncResultType) IsSkipped() bool {
	return s == SyncResultTypeSkipped
}

// IsFailed reports whether the value is SyncResultTypeFailed.
func (s SyncResultType) IsFailed() bool {
	return s == SyncResultTypeFailed
}

// HasResult reports whether the value is non-empty and valid.
func (s SyncResultType) HasResult() bool {
	return s.IsValid()
}
