package appwriter

import (
	"context"

	"coding-guidelines/common/pkg/appfault"
)

const (
	DefaultWriterName = "base-writer"
)

type (
	WriteMethodFunc func(ctx context.Context, self Writer, payload any) *appfault.AppError

	StreamMethodFunc[T any] func(ctx context.Context, self Streamer[T], payload T) *appfault.AppError
)
