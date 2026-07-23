// Package logging provides a structured `--log-json` sink that
// re-uses the jsonenv envelope so CI tooling can stream + parse
// gitmap output with the same schema it already understands.
//
// Each Log call emits one JSON object per line (NDJSON), suitable
// for piping into jq, Loki, or any log shipper.
package logging

import (
	"encoding/json"
	"io"
	"os"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// Level enumerates the supported severities.
type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

// Entry is one NDJSON line.
type Entry struct {
	Timestamp string                 `json:"ts"`
	Schema    string                 `json:"schema"`
	Version   string                 `json:"version"`
	Command   string                 `json:"command"`
	Level     Level                  `json:"level"`
	Message   string                 `json:"msg"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

// Logger emits NDJSON. Safe for concurrent use.
type Logger struct {
	mu      sync.Mutex
	w       io.Writer
	command string
	enabled bool
}

// NewLogger returns a logger that emits to w when enabled is true.
// When enabled is false, every method is a no-op (zero allocation).
func NewLogger(w io.Writer, command string, enabled bool) *Logger {
	if w == nil {
		w = os.Stderr
	}
	return &Logger{w: w, command: command, enabled: enabled}
}

// Log emits a single NDJSON entry.
func (l *Logger) Log(level Level, msg string, fields map[string]interface{}) {
	if !l.enabled {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	entry := Entry{
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Schema:    "gitmap.log.v1",
		Version:   constants.Version,
		Command:   l.command,
		Level:     level,
		Message:   msg,
		Fields:    fields,
	}
	enc := json.NewEncoder(l.w)
	_ = enc.Encode(entry)
}

// Info is a convenience wrapper.
func (l *Logger) Info(msg string, fields map[string]interface{}) {
	l.Log(LevelInfo, msg, fields)
}

// Warn is a convenience wrapper.
func (l *Logger) Warn(msg string, fields map[string]interface{}) {
	l.Log(LevelWarn, msg, fields)
}

// Error is a convenience wrapper.
func (l *Logger) Error(msg string, fields map[string]interface{}) {
	l.Log(LevelError, msg, fields)
}
