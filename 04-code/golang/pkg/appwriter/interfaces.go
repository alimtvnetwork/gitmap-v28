package appwriter

import (
	"context"
	"io"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/result"
)

type (
	Writer interface {
		Write(ctx context.Context, payload any) *appfault.AppError
		AsStreamer() Streamer[any]
		AsWriter() Writer
		Destination() io.Writer
		IsLocked() bool
		Lock()
		Unlock()
		RLock()
		RUnlock()
		Sync() *appfault.AppError
		Close() *appfault.AppError
	}

	Streamer[T any] interface {
		Stream(ctx context.Context, payload T) *appfault.AppError
		AsStreamer() Streamer[T]
		AsWriter() Writer
		Destination() io.Writer
		IsLocked() bool
		Lock()
		Unlock()
		RLock()
		RUnlock()
		Sync() *appfault.AppError
		Close() *appfault.AppError
	}

	BaseWriterWrap = result.Wrap[*BaseWriter]
)
