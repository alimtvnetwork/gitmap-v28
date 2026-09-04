package logger

import (
	"fmt"
	"os"
	"time"

	"coding-guidelines/common/pkg/appfault"
)

// Logger provides structured, leveled logging with AppError integration.
type Logger struct {
	opts LoggerOptions
}

// New creates a new Logger instance with the specified options.
func New(opts LoggerOptions) *Logger {
	return &Logger{opts: opts}
}

// Default returns a standard Logger with default options.
func Default() *Logger {
	return New(DefaultOptions())
}

// format delegates rendering to JSON or Console formatter.
func (l *Logger) format(entry LogEntry) string {
	if l.opts.IsJson {
		return formatJson(entry)
	}

	return formatConsole(entry)
}

// createLogEntry builds structured entry payload.
func createLogEntry(level LogLevel, msg, code string, ctx map[string]any, stack string) LogEntry {
	return LogEntry{
		Timestamp:  time.Now(),
		Level:      level.String(),
		Message:    msg,
		ErrorCode:  code,
		Context:    ctx,
		StackTrace: stack,
	}
}

// write dispatches the formatted log entry to the configured output writer.
func (l *Logger) write(level LogLevel, msg string, code string, ctx map[string]any, stack string) {
	if !level.IsEnabled(l.opts.Level) {
		return
	}

	entry := createLogEntry(level, msg, code, ctx, stack)

	fmt.Fprint(l.opts.Output, l.format(entry))
}

// Info logs an informational message.
func (l *Logger) Info(msg string) {
	l.write(LevelInfo, msg, "", nil, "")
}

// Warn logs a warning message.
func (l *Logger) Warn(msg string) {
	l.write(LevelWarn, msg, "", nil, "")
}

// Debug logs a debug-level message.
func (l *Logger) Debug(msg string) {
	l.write(LevelDebug, msg, "", nil, "")
}

// Error logs an error message string.
func (l *Logger) Error(msg string) {
	l.write(LevelError, msg, "", nil, "")
}

// LogError logs a structured AppError with code, context, and optional stack trace.
func (l *Logger) LogError(err *appfault.AppError) {
	if err == nil {
		return
	}

	var stack string
	if l.opts.IsStackTrace {
		stack = err.StackTrace().String()
	}

	l.write(LevelError, err.Message(), err.Type().Name(), err.Context(), stack)
}

// Fatal logs a fatal error and terminates execution.
func (l *Logger) Fatal(msg string) {
	l.write(LevelFatal, msg, "", nil, "")
	os.Exit(1)
}
