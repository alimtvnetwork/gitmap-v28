package applogger

import (
	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/appfaults"
	"coding-guidelines/common/pkg/enum/logleveltype"
)

type (
	LogLevel = logleveltype.Variant

	LogEntry struct {
		Timestamp string              `json:",omitempty" yaml:",omitempty"`
		Level     LogLevel            `json:",omitempty" yaml:",omitempty"`
		Message   string              `json:",omitempty" yaml:",omitempty"`
		Fields    appfault.ContextMap `json:",omitempty" yaml:",omitempty"`
		Caller    string              `json:",omitempty" yaml:",omitempty"`
		Stack     string              `json:",omitempty" yaml:",omitempty"`
	}

	LogSinker interface {
		WriteEntry(entry LogEntry) error
		Sync() error
		Close() error
	}

	LogSink = LogSinker

	Logger interface {
		Debug(args ...any) Logger
		Info(args ...any) Logger
		Warn(args ...any) Logger
		Error(args ...any) Logger
		Fatal(args ...any) Logger
		Debugf(format string, args ...any) Logger
		Infof(format string, args ...any) Logger
		Warnf(format string, args ...any) Logger
		Errorf(format string, args ...any) Logger
		Fatalf(format string, args ...any) Logger
		LogError(err *appfault.AppError) Logger
		LogFaults(faults *appfaults.Collection) Logger
		WithContext(key string, val any) Logger
		WithFields(fields map[string]any) Logger
		Sync() error
		Close() error
	}
)
