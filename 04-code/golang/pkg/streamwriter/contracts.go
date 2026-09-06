package streamwriter

import (
	"context"
	"fmt"
	"io"
	"time"

	"coding-guidelines/common/pkg/appfault"
)

type (
	Writer[T any] interface {
		Name() string
		Write(ctx context.Context, payload T) *appfault.AppError
		AsWriter() Writer[T]
		Lock()
		Unlock()
		Sync() *appfault.AppError
		Close() *appfault.AppError
	}

	Streamer[T any] interface {
		Name() string
		Stream(ctx context.Context, payload T) *appfault.AppError
		AsStreamer() Streamer[T]
		AsWriter() Writer[T]
		IsLocked() bool
		Lock()
		Unlock()
		Destination() io.Writer
		Sync() *appfault.AppError
		Close() *appfault.AppError
	}

	AnyWriter = PluggableWriter[any]

	AnyStreamer = Streamer[any]

	AnyLogger = Logger[any]

	LogLevel int

	LogRecord struct {
		Timestamp time.Time       `json:"timestamp"`
		Level     LogLevel        `json:"level"`
		Message   string          `json:"message"`
		Context   context.Context `json:"-"`
		Fields    map[string]any  `json:"fields,omitempty"`
		TraceId   string          `json:"traceId,omitempty"`
		UserId    string          `json:"userId,omitempty"`
	}
)

func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Compile satisfies the Compilable interface for LogRecord with deterministic ordering.
func (r LogRecord) Compile() string {
	res := fmt.Sprintf("[%s] %-5s: %s", r.Timestamp.Format("15:04:05.000"), r.Level.String(), r.Message)
	if r.TraceId != "" {
		res += fmt.Sprintf(" [trace=%s]", r.TraceId)
	}

	if r.UserId != "" {
		res += fmt.Sprintf(" [user=%s]", r.UserId)
	}

	if len(r.Fields) > 0 {
		res += fmt.Sprintf(" fields=%s", Compile(r.Fields))
	}

	return res
}

// Ensure LogRecord implements StringCompiler at compile-time.
var _ StringCompiler = LogRecord{}
var _ StreamCompiler = LogRecord{}
var _ Compilable = LogRecord{}
