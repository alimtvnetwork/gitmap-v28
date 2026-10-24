package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strconv"
	"strings"
)

// JSONKindType enumerates the JSON value types we care about. Stringer-
// style String() method exists so error messages can name the
// expected type without a separate lookup table.
type JSONKindType int

const (
	JSONKindTypeString JSONKindType = 0
	JSONKindTypeNumber JSONKindType = 1
	JSONKindTypeBool   JSONKindType = 2
	JSONKindTypeObject JSONKindType = 3
	JSONKindTypeArray  JSONKindType = 4
)

// IsValid reports whether the value matches a known variant.
func (j JSONKindType) IsValid() bool {
	switch j {
	case JSONKindTypeString, JSONKindTypeNumber, JSONKindTypeBool, JSONKindTypeObject, JSONKindTypeArray:
		return true
	default:
		return false
	}
}

func (k JSONKindType) String() string {
	switch k {
	case JSONKindTypeString:
		return "string"
	case JSONKindTypeNumber:
		return "number"
	case JSONKindTypeBool:
		return "boolean"
	case JSONKindTypeObject:
		return "object"
	case JSONKindTypeArray:
		return "array"
	}

	return "unknown"
}

// ParseJSONKindType parses a raw string into a validated JSONKindType.
func ParseJSONKindType(s string) (JSONKindType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	if n, err := strconv.Atoi(trimmed); err == nil {
		val := JSONKindType(n)
		if val.IsValid() {
			return val, nil
		}
	}
	for _, val := range AllJSONKindTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "jsonkindtypestring":
		return JSONKindTypeString, nil
	case "string":
		return JSONKindTypeString, nil
	case "jsonkindtypenumber":
		return JSONKindTypeNumber, nil
	case "number":
		return JSONKindTypeNumber, nil
	case "jsonkindtypebool":
		return JSONKindTypeBool, nil
	case "bool":
		return JSONKindTypeBool, nil
	case "jsonkindtypeobject":
		return JSONKindTypeObject, nil
	case "object":
		return JSONKindTypeObject, nil
	case "jsonkindtypearray":
		return JSONKindTypeArray, nil
	case "array":
		return JSONKindTypeArray, nil
	}
	return 0, apperror.NewValidation("ParseJSONKindType", "E1000", fmt.Sprintf("invalid jsonkindtype: %q", s)).WithContext("value", s)
}

// AllJSONKindTypes returns all valid variants.
func AllJSONKindTypes() []JSONKindType {
	return []JSONKindType{
		JSONKindTypeString,
		JSONKindTypeNumber,
		JSONKindTypeBool,
		JSONKindTypeObject,
		JSONKindTypeArray,
	}
}

// IsString reports whether the value is JSONKindTypeString.
func (j JSONKindType) IsString() bool {
	return j == JSONKindTypeString
}

// IsNumber reports whether the value is JSONKindTypeNumber.
func (j JSONKindType) IsNumber() bool {
	return j == JSONKindTypeNumber
}

// IsBool reports whether the value is JSONKindTypeBool.
func (j JSONKindType) IsBool() bool {
	return j == JSONKindTypeBool
}

// IsObject reports whether the value is JSONKindTypeObject.
func (j JSONKindType) IsObject() bool {
	return j == JSONKindTypeObject
}

// IsArray reports whether the value is JSONKindTypeArray.
func (j JSONKindType) IsArray() bool {
	return j == JSONKindTypeArray
}

// HasKind reports whether the value is non-empty and valid.
func (j JSONKindType) HasKind() bool {
	return j.IsValid()
}
