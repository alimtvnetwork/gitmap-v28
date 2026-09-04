package streamwriter

import (
	"context"
	"sync"

	"coding-guidelines/common/pkg/appfault"
)

// WriterOptions configures the pluggable writer for payload type T.
type WriterOptions[T any] struct {
	Name         string
	Streamer     Streamer[T]
	FormatMethod FormatFunc[T]
	WriteMethod  WriteFunc[T]
}

// PluggableWriter provides a composable write engine over type T with AppError returns.
type PluggableWriter[T any] struct {
	mu           sync.RWMutex
	name         string
	streamer     Streamer[T]
	formatMethod FormatFunc[T]
	writeMethod  WriteFunc[T]
}

// NewPluggableWriter constructs a pluggable writer over generic type T.
func NewPluggableWriter[T any](opts WriterOptions[T]) *PluggableWriter[T] {
	name := opts.Name
	if name == "" {
		name = "pluggable-writer"
	}

	w := &PluggableWriter[T]{
		name:         name,
		streamer:     opts.Streamer,
		formatMethod: opts.FormatMethod,
	}

	if opts.WriteMethod != nil {
		w.writeMethod = opts.WriteMethod
	} else {
		w.writeMethod = w.defaultWrite
	}
	return w
}

// Name returns the writer identifier.
func (w *PluggableWriter[T]) Name() string {
	return w.name
}

// Write delegates to the active writeMethod function under read-lock, returning *appfault.AppError.
func (w *PluggableWriter[T]) Write(ctx context.Context, payload T) *appfault.AppError {
	w.mu.RLock()
	fn := w.writeMethod
	w.mu.RUnlock()

	return fn(ctx, payload)
}

// SetWriteMethod hot-swaps the write method at runtime.
func (w *PluggableWriter[T]) SetWriteMethod(fn WriteFunc[T]) {
	if fn == nil {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.writeMethod = fn
}

// SetFormatMethod hot-swaps the formatter function at runtime.
func (w *PluggableWriter[T]) SetFormatMethod(fn FormatFunc[T]) {
	if fn == nil {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.formatMethod = fn
}

// SetStreamer hot-swaps the underlying streamer at runtime.
func (w *PluggableWriter[T]) SetStreamer(s Streamer[T]) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.streamer = s
}

// Streamer returns the attached streamer under read-lock.
func (w *PluggableWriter[T]) Streamer() Streamer[T] {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.streamer
}

// AsWriter returns the self-binding Writer[T].
func (w *PluggableWriter[T]) AsWriter() Writer[T] {
	return w
}

// AsInterfacer returns the self-binding Interfacer.
func (w *PluggableWriter[T]) AsInterfacer() Interfacer {
	return w
}

// Sync flushes the underlying streamer if attached.
func (w *PluggableWriter[T]) Sync() *appfault.AppError {
	w.mu.RLock()
	s := w.streamer
	w.mu.RUnlock()

	if s != nil {
		return s.Sync()
	}
	return nil
}

// Close closes the underlying streamer if attached.
func (w *PluggableWriter[T]) Close() *appfault.AppError {
	w.mu.Lock()
	s := w.streamer
	w.mu.Unlock()

	if s != nil {
		return s.Close()
	}
	return nil
}

func (w *PluggableWriter[T]) defaultWrite(ctx context.Context, payload T) *appfault.AppError {
	w.mu.RLock()
	s := w.streamer
	formatter := w.formatMethod
	w.mu.RUnlock()

	if formatter != nil {
		bytesResult := formatter(payload)
		if bytesResult.HasError() {
			return bytesResult.AppError()
		}
	}

	if s != nil {
		return s.Stream(ctx, payload)
	}

	return nil
}

var _ Writer[any] = (*PluggableWriter[any])(nil)
