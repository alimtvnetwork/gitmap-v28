package applogger

import (
	"fmt"
	"time"

	"coding-guidelines/common/pkg/appfault"
)

// appLogger is the standard implementation of Logger.
type appLogger struct {
	minLevel LogLevel
	sink     LogSink
	fields   appfault.ContextMap
}

// createEntry formats a LogEntry from input fields.
func (l *appLogger) createEntry(lvl LogLevel, msg, stack string) LogEntry {
	return LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     lvl,
		Message:   msg,
		Fields:    l.fields.Clone(),
		Caller:    appfault.CaptureCaller(4),
		Stack:     stack,
	}
}

// write creates a LogEntry and writes it to the sink.
func (l *appLogger) write(lvl LogLevel, msg string, stack string) {
	if lvl.IsEnabled(l.minLevel) && l.sink != nil {
		_ = l.sink.WriteEntry(l.createEntry(lvl, msg, stack))
	}
}

func (l *appLogger) Debug(args ...any) { l.write(LevelDebug, fmt.Sprint(args...), "") }
func (l *appLogger) Info(args ...any)  { l.write(LevelInfo, fmt.Sprint(args...), "") }
func (l *appLogger) Warn(args ...any)  { l.write(LevelWarn, fmt.Sprint(args...), "") }
func (l *appLogger) Error(args ...any) { l.write(LevelError, fmt.Sprint(args...), "") }
func (l *appLogger) Fatal(args ...any) { l.write(LevelFatal, fmt.Sprint(args...), "") }

func (l *appLogger) Debugf(format string, args ...any) {
	l.write(LevelDebug, fmt.Sprintf(format, args...), "")
}

func (l *appLogger) Infof(format string, args ...any) {
	l.write(LevelInfo, fmt.Sprintf(format, args...), "")
}

func (l *appLogger) Warnf(format string, args ...any) {
	l.write(LevelWarn, fmt.Sprintf(format, args...), "")
}

func (l *appLogger) Errorf(format string, args ...any) {
	l.write(LevelError, fmt.Sprintf(format, args...), "")
}

func (l *appLogger) Fatalf(format string, args ...any) {
	l.write(LevelFatal, fmt.Sprintf(format, args...), "")
}
