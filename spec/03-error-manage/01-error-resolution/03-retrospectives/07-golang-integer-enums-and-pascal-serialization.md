# Retrospective: Golang Integer-Backed Enums & PascalCase Serialization Standard

**Date:** 2026-09-01  
**Category:** Go Enum Guidelines & Error Management Architecture  
**Status:** Approved & Enforced

---

## 1. Context & Motivation

In Go, defining raw string constants for enums (e.g. `const SeverityInfo = "INFO"`) violates memory efficiency, type safety, and clean API design. In accordance with core architecture standards and lessons from `auk-go/errorwrapper`:
1. All Go enums **MUST** be backed by an integer type (`byte` / `uint8` for sets $\le 128$/255, `uint16` for wide sets such as error type variations, and `uint32` for large identifier spaces).
2. The zero value (`iota = 0`) **MUST** always represent the default/safe state (`Invalid`, `Unknown`, `None`, or `NoError`).
3. The string representation of an enum **MUST** be PascalCase (`"Info"`, `"Warn"`, `"Error"`, `"Critical"`, `"Fatal"`, `"Low"`, `"Normal"`, `"High"`). Screaming uppercase strings (`"INFO"`, `"WARN"`) are strictly forbidden.
4. Enums **MUST** provide JSON/YAML marshaling methods (`MarshalJSON()`, `UnmarshalJSON()`) to transpile cleanly to and from their PascalCase string representations.

---

## 2. Standard Pattern Example

```go
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
```

---

## 3. Rules Summary for AI Agents

1. **Underlying Type Selection:**
   - Use `byte` (alias for `uint8`) for $< 256$ variants.
   - Use `uint16` for error variations (`errtype.Variation`) and extensible taxonomies.
   - Use `uint32` for high-volume entity identifiers.
2. **Zero Value Safety:**
   - Always assign `iota = 0` to `Invalid`, `Unknown`, `None`, or `NoError`.
3. **PascalCase Strings:**
   - All string outputs MUST be PascalCase.
4. **Marshaling:**
   - Always implement `MarshalJSON` and `UnmarshalJSON` for seamless string/integer serialization.
