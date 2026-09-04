package logger

import (
	"encoding/json"
	"fmt"
)

// LogLevel defines the byte-backed severity ranking for structured log messages.
type LogLevel byte

const (
	LevelUnknown LogLevel = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var levelNames = [...]string{"Unknown", "Debug", "Info", "Warn", "Error", "Fatal"}

// Name returns the PascalCase string representation of the log level.
func (l LogLevel) Name() string {
	if int(l) < len(levelNames) {
		return levelNames[l]
	}

	return fmt.Sprintf("LogLevel(%d)", byte(l))
}

// String implements fmt.Stringer returning the PascalCase name.
func (l LogLevel) String() string {
	return l.Name()
}

// MarshalJSON serializes the log level as a PascalCase string.
func (l LogLevel) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.Name())
}

// parseLevelName looks up a LogLevel by name.
func parseLevelName(str string) (LogLevel, bool) {
	for idx, name := range levelNames {
		if name == str {
			return LogLevel(idx), true
		}
	}

	return LevelUnknown, false
}

// unmarshalByte unmarshals raw byte into LogLevel.
func (l *LogLevel) unmarshalByte(data []byte) error {
	var raw byte
	err := json.Unmarshal(data, &raw)
	*l = LogLevel(raw)

	return err
}

// UnmarshalJSON parses a PascalCase string or byte into LogLevel.
func (l *LogLevel) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		if val, ok := parseLevelName(str); ok {
			*l = val

			return nil
		}
	}

	return l.unmarshalByte(data)
}

// IsEnabled returns true if the current level meets or exceeds the target threshold.
func (l LogLevel) IsEnabled(threshold LogLevel) bool {
	return l >= threshold
}
