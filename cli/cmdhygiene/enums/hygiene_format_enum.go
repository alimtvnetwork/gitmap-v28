package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// HygieneFormatType is the value of --format; "" means default table output.
type HygieneFormatType string

const (
	HygieneFormatTypeTable HygieneFormatType = ""
	HygieneFormatTypeJSON  HygieneFormatType = "json"
	HygieneFormatTypeCSV   HygieneFormatType = "csv"
)

// IsValid reports whether the value matches a known variant.
func (h HygieneFormatType) IsValid() bool {
	switch h {
	case HygieneFormatTypeTable, HygieneFormatTypeJSON, HygieneFormatTypeCSV:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (h HygieneFormatType) String() string {
	return string(h)
}

// ParseHygieneFormatType parses a raw string into a validated HygieneFormatType.
func ParseHygieneFormatType(s string) (HygieneFormatType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllHygieneFormatTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "hygieneformattypetable":
		return HygieneFormatTypeTable, nil
	case "table":
		return HygieneFormatTypeTable, nil
	case "hygieneformattypejson":
		return HygieneFormatTypeJSON, nil
	case "json":
		return HygieneFormatTypeJSON, nil
	case "hygieneformattypecsv":
		return HygieneFormatTypeCSV, nil
	case "csv":
		return HygieneFormatTypeCSV, nil
	}
	return "", apperror.NewValidation("ParseHygieneFormatType", "E1000", fmt.Sprintf("invalid hygieneformattype: %q", s)).WithContext("value", s)
}

// AllHygieneFormatTypes returns all valid variants.
func AllHygieneFormatTypes() []HygieneFormatType {
	return []HygieneFormatType{
		HygieneFormatTypeTable,
		HygieneFormatTypeJSON,
		HygieneFormatTypeCSV,
	}
}

// IsTable reports whether the value is HygieneFormatTypeTable.
func (h HygieneFormatType) IsTable() bool {
	return h == HygieneFormatTypeTable
}

// IsJSON reports whether the value is HygieneFormatTypeJSON.
func (h HygieneFormatType) IsJSON() bool {
	return h == HygieneFormatTypeJSON
}

// IsCSV reports whether the value is HygieneFormatTypeCSV.
func (h HygieneFormatType) IsCSV() bool {
	return h == HygieneFormatTypeCSV
}

// HasFormat reports whether the value is non-empty and valid.
func (h HygieneFormatType) HasFormat() bool {
	return h.IsValid()
}
