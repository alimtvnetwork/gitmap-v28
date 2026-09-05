package errtype

import (
	"encoding/json"
	"fmt"
	"strings"
)

var variationNames = map[Variation]string{
	None:          "None",
	Generic:       "Generic",
	Validation:    "Validation",
	NotFound:      "NotFound",
	Precondition:  "Precondition",
	Execution:     "Execution",
	Database:      "Database",
	Network:       "Network",
	Timeout:       "Timeout",
	IO:            "IO",
	Unauthorized:  "Unauthorized",
	Forbidden:     "Forbidden",
	Internal:      "Internal",
	Unknown:       "Unknown",
	Serialization: "Serialization",
}

// Name returns the string representation of standard Variations.
func (v Variation) Name() string {
	if name, ok := variationNames[v]; ok {
		return name
	}

	return fmt.Sprintf("Custom(%d)", uint16(v))
}

// String implements fmt.Stringer returning the Name.
func (v Variation) String() string {
	return v.Name()
}

// Code returns the raw uint16 code value.
func (v Variation) Code() uint16 {
	return uint16(v)
}

// Value is an alias for Code.
func (v Variation) Value() uint16 {
	return v.Code()
}

// Int returns the int representation of Variation code.
func (v Variation) Int() int {
	return int(v)
}

// HasError returns true if the Variation represents an actual error (non-None).
func (v Variation) HasError() bool {
	return v != None
}

// IsNone returns true if the Variation represents no error (successful state).
func (v Variation) IsNone() bool {
	return v == None
}

// IsNoError is an alias for IsNone.
func (v Variation) IsNoError() bool {
	return v.IsNone()
}

// IsValid returns true if the Variation represents NoError / None (healthy/valid state).
func (v Variation) IsValid() bool {
	return v == None
}

// IsInvalid returns true if the Variation represents an active error type.
func (v Variation) IsInvalid() bool {
	return v != None
}

// Is returns true if this Variation matches the target.
func (v Variation) Is(target Variation) bool {
	return v == target
}

// HttpStatus maps the error type Variation to standard HTTP status codes.
func (v Variation) HttpStatus() int {
	switch v {
	case None:
		return 200

	case Validation, Precondition, Serialization:
		return 400

	case Unauthorized:
		return 401

	case Forbidden:
		return 403

	case NotFound:
		return 404

	case Timeout:
		return 504

	default:
		return 500
	}
}

// MarshalJSON returns the variation as a 16-bit unsigned integer code.
func (v Variation) MarshalJSON() ([]byte, error) {
	return json.Marshal(uint16(v))
}

// UnmarshalJSON unmarshals from either an integer code (2) or a name string ("Validation").
func (v *Variation) UnmarshalJSON(data []byte) error {
	trimmed := strings.TrimSpace(string(data))
	if len(trimmed) == 0 || trimmed == "null" {
		*v = None

		return nil
	}

	if trimmed[0] >= '0' && trimmed[0] <= '9' {
		var code uint16
		if err := json.Unmarshal(data, &code); err != nil {
			return err
		}

		*v = Variation(code)

		return nil
	}

	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return err
	}

	for varEnum, varName := range variationNames {
		if strings.EqualFold(varName, name) {
			*v = varEnum

			return nil
		}
	}

	*v = Generic

	return nil
}
