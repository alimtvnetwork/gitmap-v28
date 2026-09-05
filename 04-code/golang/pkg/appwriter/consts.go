package appwriter

import (
	"context"

	"coding-guidelines/common/pkg/appfault"
)

// Default constants for appwriter.
const (
	DefaultWriterName = "base-writer"
)

// WriteMethodFunc defines an injected write method that accepts self as the first parameter.
type WriteMethodFunc func(ctx context.Context, self Writer, payload any) *appfault.AppError

// StreamMethodFunc defines an injected streamer method that accepts self as the first parameter.
type StreamMethodFunc[T any] func(ctx context.Context, self Streamer[T], payload T) *appfault.AppError
