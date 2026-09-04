package applogger

import (
	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/appfaults"
	"coding-guidelines/common/pkg/logger"
)

// LogLevel aliases logger.LogLevel.
type LogLevel = logger.LogLevel

const (
	LevelUnknown = logger.LevelUnknown
	LevelDebug   = logger.LevelDebug
	LevelInfo    = logger.LevelInfo
	LevelWarn    = logger.LevelWarn
	LevelError   = logger.LevelError
	LevelFatal   = logger.LevelFatal
)

// LogEntry represents a structured log event payload.
type LogEntry struct {
	Timestamp string              `json:"Timestamp" yaml:"Timestamp"`
	Level     LogLevel            `json:"Level" yaml:"Level"`
	Message   string              `json:"Message" yaml:"Message"`
	Fields    appfault.ContextMap `json:"Fields,omitempty" yaml:"Fields,omitempty"`
	Caller    string              `json:"Caller,omitempty" yaml:"Caller,omitempty"`
	Stack     string              `json:"Stack,omitempty" yaml:"Stack,omitempty"`
}

// LogSink defines the driver interface for persistent destinations.
type LogSink interface {
	WriteEntry(entry LogEntry) error
	Sync() error
	Close() error
}

// Logger is the unified logging interface.
type Logger interface {
	Debug(args ...any)
	Info(args ...any)
	Warn(args ...any)
	Error(args ...any)
	Fatal(args ...any)
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
	LogError(err *appfault.AppError)
	LogFaults(faults *appfaults.Collection)
	WithContext(key string, val any) Logger
	WithFields(fields map[string]any) Logger
	Sync() error
	Close() error
}
