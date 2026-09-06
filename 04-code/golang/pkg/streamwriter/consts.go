package streamwriter

import (
	"context"
	"io"

	"coding-guidelines/common/pkg/appfault"
)

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

const (
	PayloadNil PayloadKind = iota
	PayloadBytes
	PayloadString
	PayloadError
	PayloadMap
	PayloadStruct
	PayloadPrimitive
)

type (
	StreamFunc[T any] func(ctx context.Context, payload T, dest io.Writer) *appfault.AppError

	WriteFunc[T any] func(streamer Streamer[T], ctx context.Context, writer *PluggableWriter[T], payload T) *appfault.AppError

	FormatFunc[T any] func(payload T) Bytes[T]

	ErrorHandlerFunc func(err *appfault.AppError)

	LogFormatterFunc func(record LogRecord) string
)
