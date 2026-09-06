package appfault

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// SeverityType represents an integer-backed severity level (byte).
type SeverityType byte

var (
	severityNames = [...]string{"Unknown", "Info", "Warn", "Error", "Critical", "Fatal"}
	severityMap   = compileSeverityMap()
)

func compileSeverityMap() map[string]SeverityType {
	m := make(map[string]SeverityType, len(severityNames)*4)
	for idx, name := range severityNames {
		s := SeverityType(idx)
		m[name] = s
		m[strings.ToLower(name)] = s
		m[strings.ToUpper(name)] = s
		m[strconv.Itoa(idx)] = s
	}

	return m
}

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
	val, ok := severityMap[strings.ToLower(strings.TrimSpace(str))]

	return val, ok
}

// UnmarshalJSON parses a PascalCase string or integer into SeverityType.
func (s *SeverityType) UnmarshalJSON(data []byte) error {
	trimmed := strings.TrimSpace(string(data))
	if len(trimmed) == 0 || trimmed == "null" {
		*s = SeverityUnknown

		return nil
	}

	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		if val, ok := parseSeverityName(str); ok {
			*s = val

			return nil
		}

		return fmt.Errorf("unknown SeverityType %q, supported: [%s]", str, strings.Join(severityNames[:], ", "))
	}

	var raw byte
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if int(raw) >= len(severityNames) {
		return fmt.Errorf("invalid SeverityType numeric value %d, supported range: 0..%d", raw, len(severityNames)-1)
	}

	*s = SeverityType(raw)

	return nil
}
