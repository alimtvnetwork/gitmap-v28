package appfault

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// PriorityType represents an integer-backed priority level (byte).
type PriorityType byte

var (
	priorityNames = [...]string{"Unknown", "Low", "Normal", "High", "Critical"}
	priorityMap   = compilePriorityMap()
)

func compilePriorityMap() map[string]PriorityType {
	m := make(map[string]PriorityType, len(priorityNames)*4)
	for idx, name := range priorityNames {
		p := PriorityType(idx)
		m[name] = p
		m[strings.ToLower(name)] = p
		m[strings.ToUpper(name)] = p
		m[strconv.Itoa(idx)] = p
	}

	return m
}

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
	val, ok := priorityMap[strings.ToLower(strings.TrimSpace(str))]

	return val, ok
}

// UnmarshalJSON parses a PascalCase string or integer into PriorityType.
func (p *PriorityType) UnmarshalJSON(data []byte) error {
	trimmed := strings.TrimSpace(string(data))
	if len(trimmed) == 0 || trimmed == "null" {
		*p = PriorityUnknown

		return nil
	}

	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		if val, ok := parsePriorityName(str); ok {
			*p = val

			return nil
		}

		return fmt.Errorf("unknown PriorityType %q, supported: [%s]", str, strings.Join(priorityNames[:], ", "))
	}

	var raw byte
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if int(raw) >= len(priorityNames) {
		return fmt.Errorf("invalid PriorityType numeric value %d, supported range: 0..%d", raw, len(priorityNames)-1)
	}

	*p = PriorityType(raw)

	return nil
}
