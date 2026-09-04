package errtype

import "fmt"

var variationNames = map[Variation]string{
	None:         "None",
	Generic:      "Generic",
	Validation:   "Validation",
	NotFound:     "NotFound",
	Precondition: "Precondition",
	Execution:    "Execution",
	Database:     "Database",
	Network:      "Network",
	Timeout:      "Timeout",
	IO:           "IO",
	Unauthorized: "Unauthorized",
	Forbidden:    "Forbidden",
	Internal:     "Internal",
	Unknown:      "Unknown",
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
