package appfault

import (
	"encoding/json"
	"fmt"
)

// SeverityType represents an integer-backed severity level (byte).
type SeverityType byte

const (
	SeverityUnknown SeverityType = iota
	SeverityInfo
	SeverityWarn
	SeverityError
	SeverityCritical
	SeverityFatal
)

var severityNames = [...]string{"Unknown", "Info", "Warn", "Error", "Critical", "Fatal"}

// Name returns the PascalCase string representation.
func (s SeverityType) Name() string {
	if int(s) < len(severityNames) {
		return severityNames[s]
	}

	return fmt.Sprintf("Severity(%d)", byte(s))
}

// String implements fmt.Stringer.
func (s SeverityType) String() string {
	return s.Name()
}

// MarshalJSON serializes the severity as a PascalCase string.
func (s SeverityType) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Name())
}

// parseSeverityName looks up a SeverityType by name.
func parseSeverityName(str string) (SeverityType, bool) {
	for idx, name := range severityNames {
		if name == str {
			return SeverityType(idx), true
		}
	}

	return SeverityUnknown, false
}

// unmarshalByte unmarshals raw byte into SeverityType.
func (s *SeverityType) unmarshalByte(data []byte) error {
	var raw byte
	err := json.Unmarshal(data, &raw)
	*s = SeverityType(raw)

	return err
}

// UnmarshalJSON parses a PascalCase string or integer into SeverityType.
func (s *SeverityType) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		if val, ok := parseSeverityName(str); ok {
			*s = val

			return nil
		}
	}

	return s.unmarshalByte(data)
}
