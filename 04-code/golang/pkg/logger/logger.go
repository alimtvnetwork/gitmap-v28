package logger

import (
	"fmt"
	"os"
	"time"

	"coding-guidelines/common/pkg/appfault"
)

type Logger struct {
	opts LoggerOptions
}

func New(opts LoggerOptions) *Logger {
	return &Logger{opts: opts}
}

func Default() *Logger {
	return New(DefaultOptions())
}

func (l *Logger) format(entry LogEntry) string {
	if l.opts.IsJson {
		return formatJson(entry)
	}

	return formatConsole(entry)
}

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

func (l *Logger) write(level LogLevel, msg string, code string, ctx map[string]any, stack string) {
	if !level.IsEnabled(l.opts.Level) {
		return
	}

	entry := createLogEntry(level, msg, code, ctx, stack)
	fmt.Fprint(l.opts.Output, l.format(entry))
}

func (l *Logger) Info(msg string) {
	l.write(LevelInfo, msg, "", nil, "")
}

func (l *Logger) Warn(msg string) {
	l.write(LevelWarn, msg, "", nil, "")
}

func (l *Logger) Debug(msg string) {
	l.write(LevelDebug, msg, "", nil, "")
}

func (l *Logger) Error(msg string) {
	l.write(LevelError, msg, "", nil, "")
}

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

func (l *Logger) Fatal(msg string) {
	l.write(LevelFatal, msg, "", nil, "")
	os.Exit(1)
}
