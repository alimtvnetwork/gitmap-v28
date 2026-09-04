package appfault

import (
	"encoding/json"
	"fmt"
)

// PriorityType represents an integer-backed priority level (byte).
type PriorityType byte

const (
	PriorityUnknown PriorityType = iota
	PriorityLow
	PriorityNormal
	PriorityHigh
	PriorityCritical
)

var priorityNames = [...]string{"Unknown", "Low", "Normal", "High", "Critical"}

// Name returns the PascalCase string representation.
func (p PriorityType) Name() string {
	if int(p) < len(priorityNames) {
		return priorityNames[p]
	}

	return fmt.Sprintf("Priority(%d)", byte(p))
}

// String implements fmt.Stringer.
func (p PriorityType) String() string {
	return p.Name()
}

// MarshalJSON serializes the priority as a PascalCase string.
func (p PriorityType) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Name())
}

// parsePriorityName looks up a PriorityType by name.
func parsePriorityName(str string) (PriorityType, bool) {
	for idx, name := range priorityNames {
		if name == str {
			return PriorityType(idx), true
		}
	}

	return PriorityUnknown, false
}

// unmarshalByte unmarshals raw byte into PriorityType.
func (p *PriorityType) unmarshalByte(data []byte) error {
	var raw byte
	err := json.Unmarshal(data, &raw)
	*p = PriorityType(raw)

	return err
}

// UnmarshalJSON parses a PascalCase string or integer into PriorityType.
func (p *PriorityType) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		if val, ok := parsePriorityName(str); ok {
			*p = val

			return nil
		}
	}

	return p.unmarshalByte(data)
}
