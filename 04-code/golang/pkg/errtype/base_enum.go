package errtype

import (
	"encoding/json"
	"fmt"
	"strings"
)

// BaseEnum defines the universal contract for all enums in the ecosystem.
type BaseEnum interface {
	Name() string
	String() string
	ValueString() string
	IsValid() bool
	IsEnum() bool
}

// NumberEnum extends BaseEnum for integer-backed enums.
type NumberEnum interface {
	BaseEnum
	Int() int
	Code() uint16
}

// ValueString returns the string representation of the variation code.
func (v Variation) ValueString() string {
	return fmt.Sprintf("%d", uint16(v))
}

// IsEnum returns true if the Variation is one of the recognized standard variations.
func (v Variation) IsEnum() bool {
	_, ok := variationNames[v]

	return ok
}

// ProcessStateType represents a string-backed enum conforming to BaseEnum.
type ProcessStateType string

var processStateRegistry = map[ProcessStateType]bool{
	ProcessStatePending:   true,
	ProcessStateRunning:   true,
	ProcessStateCompleted: true,
	ProcessStateFailed:    true,
	ProcessStateCancelled: true,
}

// Name returns the identifier name.
func (s ProcessStateType) Name() string {
	return string(s)
}

// String returns the string representation.
func (s ProcessStateType) String() string {
	return string(s)
}

// ValueString returns the string representation of value.
func (s ProcessStateType) ValueString() string {
	return string(s)
}

// Value returns the raw string value.
func (s ProcessStateType) Value() string {
	return string(s)
}

// IsValid returns true if this state is non-empty and known.
func (s ProcessStateType) IsValid() bool {
	return processStateRegistry[s]
}

// IsEnum returns true if this state exists in the registry.
func (s ProcessStateType) IsEnum() bool {
	return processStateRegistry[s]
}

// IsCompare checks equality against another ProcessStateType.
func (s ProcessStateType) IsCompare(target ProcessStateType) bool {
	return s == target
}

// MarshalJSON implements json.Marshaler.
func (s ProcessStateType) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(s))
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *ProcessStateType) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*s = ParseProcessState(raw)

	return nil
}

// AllProcessStates returns all registered ProcessStateType values.
func AllProcessStates() []ProcessStateType {
	return []ProcessStateType{
		ProcessStatePending,
		ProcessStateRunning,
		ProcessStateCompleted,
		ProcessStateFailed,
		ProcessStateCancelled,
	}
}

// ParseProcessState parses a string into ProcessStateType case-insensitively.
func ParseProcessState(val string) ProcessStateType {
	for _, candidate := range AllProcessStates() {
		if strings.EqualFold(string(candidate), strings.TrimSpace(val)) {
			return candidate
		}
	}

	return ProcessStateUnknown
}

// LogLevelType represents an integer-backed enum conforming to NumberEnum and BaseEnum.
type LogLevelType uint16

var logLevelNames = map[LogLevelType]string{
	LogLevelDebug: "Debug",
	LogLevelInfo:  "Info",
	LogLevelWarn:  "Warn",
	LogLevelError: "Error",
	LogLevelFatal: "Fatal",
}

// Name returns the uppercase identifier.
func (l LogLevelType) Name() string {
	if name, ok := logLevelNames[l]; ok {
		return name
	}

	return fmt.Sprintf("LogLevel(%d)", uint16(l))
}

// String implements fmt.Stringer.
func (l LogLevelType) String() string {
	return l.Name()
}

// ValueString returns the integer code formatted as a string.
func (l LogLevelType) ValueString() string {
	return fmt.Sprintf("%d", uint16(l))
}

// Code returns the raw uint16 code value.
func (l LogLevelType) Code() uint16 {
	return uint16(l)
}

// Int returns the int representation.
func (l LogLevelType) Int() int {
	return int(l)
}

// IsValid returns true if this log level is known.
func (l LogLevelType) IsValid() bool {
	_, ok := logLevelNames[l]

	return ok
}

// IsEnum returns true if this log level exists in registry.
func (l LogLevelType) IsEnum() bool {
	_, ok := logLevelNames[l]

	return ok
}

// IsCompare checks equality against another LogLevelType.
func (l LogLevelType) IsCompare(target LogLevelType) bool {
	return l == target
}

// MarshalJSON implements json.Marshaler.
func (l LogLevelType) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.Name())
}

// UnmarshalJSON implements json.Unmarshaler.
func (l *LogLevelType) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err == nil {
		*l = ParseLogLevel(raw)

		return nil
	}

	var code uint16
	if err := json.Unmarshal(data, &code); err != nil {
		return err
	}

	*l = LogLevelType(code)

	return nil
}

// AllLogLevels returns all registered LogLevelType values.
func AllLogLevels() []LogLevelType {
	return []LogLevelType{
		LogLevelDebug,
		LogLevelInfo,
		LogLevelWarn,
		LogLevelError,
		LogLevelFatal,
	}
}

// ParseLogLevel parses a string into LogLevelType case-insensitively.
func ParseLogLevel(val string) LogLevelType {
	cleaned := strings.TrimSpace(val)
	for lvl, name := range logLevelNames {
		if strings.EqualFold(name, cleaned) {
			return lvl
		}
	}

	return 0
}

// ToEnum finds an enum by name in any slice of BaseEnum.
func ToEnum[T BaseEnum](val string, all []T) (T, bool) {
	cleaned := strings.TrimSpace(val)
	for _, item := range all {
		if strings.EqualFold(item.Name(), cleaned) || strings.EqualFold(item.ValueString(), cleaned) {
			return item, true
		}
	}

	var zero T

	return zero, false
}

var _ BaseEnum = Variation(0)
var _ NumberEnum = Variation(0)
var _ BaseEnum = ProcessStateType("")
var _ BaseEnum = LogLevelType(0)
var _ NumberEnum = LogLevelType(0)
