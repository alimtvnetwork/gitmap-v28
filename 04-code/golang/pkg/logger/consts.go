package logger

import (
	"coding-guidelines/common/pkg/enum/logleveltype"
)

const (
	LevelUnknown LogLevel = logleveltype.Unknown
	LevelInvalid LogLevel = logleveltype.Invalid
	LevelDebug   LogLevel = logleveltype.Debug
	LevelInfo    LogLevel = logleveltype.Info
	LevelWarn    LogLevel = logleveltype.Warn
	LevelError   LogLevel = logleveltype.Error
	LevelFatal   LogLevel = logleveltype.Fatal
)

type (
	LogFormatterFunc func(entry LogEntry) string

	LogFilterFunc func(entry LogEntry) bool
)
