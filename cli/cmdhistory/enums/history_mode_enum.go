package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strconv"
	"strings"
)

// HistoryModeType picks which filter-repo behavior to wrap.
type HistoryModeType int

const (
	HistoryModeTypePurge HistoryModeType = 0
	HistoryModeTypePin   HistoryModeType = 1
)

// IsValid reports whether the value matches a known variant.
func (h HistoryModeType) IsValid() bool {
	switch h {
	case HistoryModeTypePurge, HistoryModeTypePin:
		return true
	default:
		return false
	}
}

// String returns the numeric representation of the value.
func (h HistoryModeType) String() string {
	return strconv.Itoa(int(h))
}

// ParseHistoryModeType parses a raw string into a validated HistoryModeType.
func ParseHistoryModeType(s string) (HistoryModeType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	if n, err := strconv.Atoi(trimmed); err == nil {
		val := HistoryModeType(n)
		if val.IsValid() {
			return val, nil
		}
	}
	for _, val := range AllHistoryModeTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "historymodetypepurge":
		return HistoryModeTypePurge, nil
	case "purge":
		return HistoryModeTypePurge, nil
	case "historymodetypepin":
		return HistoryModeTypePin, nil
	case "pin":
		return HistoryModeTypePin, nil
	}
	return 0, apperror.NewValidation("ParseHistoryModeType", "E1000", fmt.Sprintf("invalid historymodetype: %q", s)).WithContext("value", s)
}

// AllHistoryModeTypes returns all valid variants.
func AllHistoryModeTypes() []HistoryModeType {
	return []HistoryModeType{
		HistoryModeTypePurge,
		HistoryModeTypePin,
	}
}

// IsPurge reports whether the value is HistoryModeTypePurge.
func (h HistoryModeType) IsPurge() bool {
	return h == HistoryModeTypePurge
}

// IsPin reports whether the value is HistoryModeTypePin.
func (h HistoryModeType) IsPin() bool {
	return h == HistoryModeTypePin
}

// HasMode reports whether the value is non-empty and valid.
func (h HistoryModeType) HasMode() bool {
	return h.IsValid()
}
