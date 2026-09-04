package logger

import (
	"encoding/json"
	"fmt"
	"time"
)

// LogEntry contains the payload data for a single log event.
type LogEntry struct {
	Timestamp  time.Time      `json:"timestamp"`
	Level      string         `json:"level"`
	Message    string         `json:"message"`
	ErrorCode  string         `json:"errorCode,omitempty"`
	Context    map[string]any `json:"context,omitempty"`
	StackTrace string         `json:"stackTrace,omitempty"`
}

// levelIcon returns the emoji indicator for a log level string.
func levelIcon(level string) string {
	if level == "ERROR" || level == "FATAL" {
		return "❌"
	}

	if level == "WARN" {
		return "⚠️"
	}

	return "ℹ️"
}

// formatConsole serializes a LogEntry into human-readable terminal output.
func formatConsole(entry LogEntry) string {
	ts := entry.Timestamp.Format("2006-01-02 15:04:05")
	icon := levelIcon(entry.Level)
	out := fmt.Sprintf("%s [%s] %s %s", icon, ts, entry.Level, entry.Message)
	if len(entry.ErrorCode) > 0 {
		out += fmt.Sprintf(" (code: %s)", entry.ErrorCode)
	}

	if len(entry.StackTrace) > 0 {
		out += fmt.Sprintf("\n%s", entry.StackTrace)
	}

	return out + "\n"
}

// formatJson serializes a LogEntry into compact JSON string.
func formatJson(entry LogEntry) string {
	raw, err := json.Marshal(entry)
	if err != nil {
		return fmt.Sprintf(`{"error":"json marshal failed","message":%q}`+"\n", entry.Message)
	}

	return string(raw) + "\n"
}
