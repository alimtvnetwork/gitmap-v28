package streamwriter

import (
	"context"
	"fmt"
	"io"
	"time"

	"coding-guidelines/common/pkg/appfault"
)

// Writer defines universal write operations over generic type T with AppError and Locker synchronization.
type Writer[T any] interface {
	Name() string
	Write(ctx context.Context, payload T) *appfault.AppError
	AsWriter() Writer[T]
	Lock()
	Unlock()
	Sync() *appfault.AppError
	Close() *appfault.AppError
}

// Streamer defines streaming operations over generic type T with AppError and Locker synchronization.
type Streamer[T any] interface {
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

// StreamFunc defines the swappable function signature returning *appfault.AppError.
type StreamFunc[T any] func(ctx context.Context, payload T, dest io.Writer) *appfault.AppError

// WriteFunc defines the swappable function signature returning *appfault.AppError.
type WriteFunc[T any] func(ctx context.Context, payload T) *appfault.AppError

// FormatFunc defines the serialization transformation returning Bytes[T].
type FormatFunc[T any] func(payload T) Bytes[T]

// LogLevel defines standardized severity tiers.
type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
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

// LogRecord carries normalized event data for log-based flows.
type LogRecord struct {
	Timestamp time.Time       `json:"timestamp"`
	Level     LogLevel        `json:"level"`
	Message   string          `json:"message"`
	Context   context.Context `json:"-"`
	Fields    map[string]any  `json:"fields,omitempty"`
	TraceID   string          `json:"traceId,omitempty"`
	UserID    string          `json:"userId,omitempty"`
}

// Compile satisfies the Compilable interface for LogRecord with deterministic ordering.
func (r LogRecord) Compile() string {
	res := fmt.Sprintf("[%s] %-5s: %s", r.Timestamp.Format("15:04:05.000"), r.Level.String(), r.Message)
	if r.TraceID != "" {
		res += fmt.Sprintf(" [trace=%s]", r.TraceID)
	}
	if r.UserID != "" {
		res += fmt.Sprintf(" [user=%s]", r.UserID)
	}
	if len(r.Fields) > 0 {
		res += fmt.Sprintf(" fields=%s", Compile(r.Fields))
	}
	return res
}

// Ensure LogRecord implements Compilable at compile-time.
var _ Compilable = LogRecord{}
