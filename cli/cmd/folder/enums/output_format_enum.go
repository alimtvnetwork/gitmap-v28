package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// OutputFormat defines the rendering strategy for folder output.
type OutputFormatType string

const (
	FormatTree OutputFormatType = "tree"
	FormatMd   OutputFormatType = "md"
	FormatJson OutputFormatType = "json"
	FormatYaml OutputFormatType = "yaml"
	FormatFlat OutputFormatType = "flat"
)

// IsValid reports whether the value matches a known variant.
func (o OutputFormatType) IsValid() bool {
	switch o {
	case FormatTree, FormatMd, FormatJson, FormatYaml, FormatFlat:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (o OutputFormatType) String() string {
	return string(o)
}

// ParseOutputFormatType parses a raw string into a validated OutputFormatType.
func ParseOutputFormatType(s string) (OutputFormatType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllOutputFormatTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "formattree":
		return FormatTree, nil
	case "tree":
		return FormatTree, nil
	case "formatmd":
		return FormatMd, nil
	case "md":
		return FormatMd, nil
	case "formatjson":
		return FormatJson, nil
	case "json":
		return FormatJson, nil
	case "formatyaml":
		return FormatYaml, nil
	case "yaml":
		return FormatYaml, nil
	case "formatflat":
		return FormatFlat, nil
	case "flat":
		return FormatFlat, nil
	}
	return "", apperror.NewValidation("ParseOutputFormatType", "E1000", fmt.Sprintf("invalid outputformattype: %q", s)).WithContext("value", s)
}

// AllOutputFormatTypes returns all valid variants.
func AllOutputFormatTypes() []OutputFormatType {
	return []OutputFormatType{
		FormatTree,
		FormatMd,
		FormatJson,
		FormatYaml,
		FormatFlat,
	}
}

// IsTree reports whether the value is FormatTree.
func (o OutputFormatType) IsTree() bool {
	return o == FormatTree
}

// IsMd reports whether the value is FormatMd.
func (o OutputFormatType) IsMd() bool {
	return o == FormatMd
}

// IsJson reports whether the value is FormatJson.
func (o OutputFormatType) IsJson() bool {
	return o == FormatJson
}

// IsYaml reports whether the value is FormatYaml.
func (o OutputFormatType) IsYaml() bool {
	return o == FormatYaml
}

// IsFlat reports whether the value is FormatFlat.
func (o OutputFormatType) IsFlat() bool {
	return o == FormatFlat
}

// HasFormat reports whether the value is non-empty and valid.
func (o OutputFormatType) HasFormat() bool {
	return o.IsValid()
}
