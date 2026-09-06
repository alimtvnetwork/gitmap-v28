package streamwriter

import (
	"context"
	"fmt"
	"io"
	"sync"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
)

type (
	WriterOptions[T any] struct {
		Name         string
		Destination  io.Writer
		Streamer     Streamer[T]
		FormatMethod FormatFunc[T]
		WriteMethod  WriteFunc[T]
	}

	PluggableWriter[T any] struct {
		mu           ReentrantMutex
		configMu     sync.RWMutex
		name         string
		destination  io.Writer
		streamer     Streamer[T]
		formatMethod FormatFunc[T]
		writeMethod  WriteFunc[T]
	}
)

// NewPluggableWriter constructs a pluggable writer over generic type T.
func NewPluggableWriter[T any](opts WriterOptions[T]) *PluggableWriter[T] {
	name := opts.Name
	if name == "" {
		name = "pluggable-writer"
	}

	w := &PluggableWriter[T]{
		name:         name,
		destination:  opts.Destination,
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

// NewAnyWriter constructs a non-generic AnyWriter (*PluggableWriter[any]) with smart payload dispatch.
func NewAnyWriter(opts WriterOptions[any]) *AnyWriter {
	return NewPluggableWriter[any](opts)
}

// Name returns the writer identifier.
func (w *PluggableWriter[T]) Name() string {
	return w.name
}

// Destination returns the destination io.Writer under read-lock, falling back to attached streamer if present.
func (w *PluggableWriter[T]) Destination() io.Writer {
	w.configMu.RLock()
	defer w.configMu.RUnlock()
	if w.destination != nil {
		return w.destination
	}

	if w.streamer != nil {
		return w.streamer.Destination()
	}

	return nil
}

// SetDestination hot-swaps the destination at runtime.
func (w *PluggableWriter[T]) SetDestination(dest io.Writer) {
	w.configMu.Lock()
	defer w.configMu.Unlock()
	w.destination = dest
}

// FormatMethod returns the attached formatter function under read-lock.
func (w *PluggableWriter[T]) FormatMethod() FormatFunc[T] {
	w.configMu.RLock()
	defer w.configMu.RUnlock()

	return w.formatMethod
}

// Write delegates to the active writeMethod function under lock, returning *appfault.AppError.
// It passes the current writer object (w) into writeMethod to grant access to properties.
func (w *PluggableWriter[T]) Write(ctx context.Context, payload T) *appfault.AppError {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.configMu.RLock()
	fn := w.writeMethod
	s := w.streamer
	w.configMu.RUnlock()

	return fn(s, ctx, w, payload)
}

// SetWriteMethod hot-swaps the write method at runtime.
func (w *PluggableWriter[T]) SetWriteMethod(fn WriteFunc[T]) {
	if fn == nil {
		return
	}

	w.configMu.Lock()
	defer w.configMu.Unlock()
	w.writeMethod = fn
}

// SetFormatMethod hot-swaps the formatter function at runtime.
func (w *PluggableWriter[T]) SetFormatMethod(fn FormatFunc[T]) {
	if fn == nil {
		return
	}

	w.configMu.Lock()
	defer w.configMu.Unlock()
	w.formatMethod = fn
}

// SetStreamer hot-swaps the underlying streamer at runtime.
func (w *PluggableWriter[T]) SetStreamer(s Streamer[T]) {
	w.configMu.Lock()
	defer w.configMu.Unlock()
	w.streamer = s
}

// Streamer returns the attached streamer under read-lock.
func (w *PluggableWriter[T]) Streamer() Streamer[T] {
	w.configMu.RLock()
	defer w.configMu.RUnlock()

	return w.streamer
}

// AsWriter returns the self-binding Writer[T].
func (w *PluggableWriter[T]) AsWriter() Writer[T] {
	return w
}

// Lock locks the writer for exclusive access, satisfying sync.Locker.
func (w *PluggableWriter[T]) Lock() {
	w.mu.Lock()
}

// Unlock unlocks the writer, satisfying sync.Locker.
func (w *PluggableWriter[T]) Unlock() {
	w.mu.Unlock()
}

// Sync flushes the underlying streamer if attached.
func (w *PluggableWriter[T]) Sync() *appfault.AppError {
	w.configMu.RLock()
	s := w.streamer
	w.configMu.RUnlock()

	if s != nil {
		return s.Sync()
	}

	return nil
}

// Close closes the underlying streamer if attached.
func (w *PluggableWriter[T]) Close() *appfault.AppError {
	w.configMu.Lock()
	s := w.streamer
	w.configMu.Unlock()

	if s != nil {
		return s.Close()
	}

	return nil
}

func (w *PluggableWriter[T]) defaultWrite(streamer Streamer[T], ctx context.Context, writer *PluggableWriter[T], payload T) *appfault.AppError {
	w.configMu.RLock()
	formatter := w.formatMethod
	dest := w.destination
	w.configMu.RUnlock()

	if formatter != nil {
		bytesResult := formatter(payload)
		if bytesResult.HasError() {
			return bytesResult.AppError()
		}
	}

	if streamer != nil {
		return streamer.Stream(ctx, payload)
	}

	if dest != nil {
		line := fmt.Sprintf("[%s] %s\n", w.name, Compile(payload))
		if _, err := dest.Write([]byte(line)); err != nil {
			return appfault.Wrap(errtype.IO, err, fmt.Sprintf("writer %s write failed", w.name))
		}
	}

	return nil
}

var _ Writer[any] = (*PluggableWriter[any])(nil)
var _ sync.Locker = (*PluggableWriter[any])(nil)
