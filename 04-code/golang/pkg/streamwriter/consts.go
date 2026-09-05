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

type StreamFunc[T any] func(ctx context.Context, payload T, dest io.Writer) *appfault.AppError

type WriteFunc[T any] func(streamer Streamer[T], ctx context.Context, writer *PluggableWriter[T], payload T) *appfault.AppError

type FormatFunc[T any] func(payload T) Bytes[T]

type ErrorHandlerFunc func(err *appfault.AppError)

type LogFormatterFunc func(record LogRecord) string
