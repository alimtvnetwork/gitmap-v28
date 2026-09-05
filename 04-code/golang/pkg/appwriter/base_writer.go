package appwriter

import (
	"context"
	"io"
	"sync"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
)

// BaseWriter implements the core Writer and Streamer lifecycle.
type BaseWriter struct {
	name      string
	dest      io.Writer
	closer    io.Closer
	syncer    interface{ Sync() error }
	isLocked  bool
	mu        sync.RWMutex
	writeFunc WriteMethodFunc
}

// NewBaseWriter constructs a BaseWriter instance.
func NewBaseWriter(name string, dest io.Writer, isLocked bool, writeFunc WriteMethodFunc) *BaseWriter {
	closer, _ := dest.(io.Closer)
	syncer, _ := dest.(interface{ Sync() error })

	return &BaseWriter{
		name:      name,
		dest:      dest,
		closer:    closer,
		syncer:    syncer,
		isLocked:  isLocked,
		writeFunc: writeFunc,
	}
}

// Name returns the writer identifier.
func (w *BaseWriter) Name() string {
	return w.name
}

// Destination returns the underlying I/O destination writer.
func (w *BaseWriter) Destination() io.Writer {
	return w.dest
}

// IsLocked reports whether the writer enforces concurrent mutex synchronization.
func (w *BaseWriter) IsLocked() bool {
	return w.isLocked
}

// Lock acquires the mutual exclusion lock if lock mode is enabled.
func (w *BaseWriter) Lock() {
	if w.isLocked {
		w.mu.Lock()
	}
}

// Unlock releases the mutual exclusion lock if lock mode is enabled.
func (w *BaseWriter) Unlock() {
	if w.isLocked {
		w.mu.Unlock()
	}
}

// RLock acquires the shared read lock if lock mode is enabled.
func (w *BaseWriter) RLock() {
	if w.isLocked {
		w.mu.RLock()
	}
}

// RUnlock releases the shared read lock if lock mode is enabled.
func (w *BaseWriter) RUnlock() {
	if w.isLocked {
		w.mu.RUnlock()
	}
}

// AsWriter returns the Writer interface representation.
func (w *BaseWriter) AsWriter() Writer {
	return w
}

// AsStreamer returns the Streamer interface representation.
func (w *BaseWriter) AsStreamer() Streamer[any] {
	return w
}

// Write executes the injected write function passing self as the first parameter.
func (w *BaseWriter) Write(ctx context.Context, payload any) *appfault.AppError {
	if w.writeFunc == nil {
		return appfault.New(errtype.Precondition, "write function is unassigned")
	}

	w.Lock()
	defer w.Unlock()

	return w.writeFunc(ctx, w, payload)
}

// Stream forwards payload to Write satisfying the Streamer interface.
func (w *BaseWriter) Stream(ctx context.Context, payload any) *appfault.AppError {
	return w.Write(ctx, payload)
}

// Sync flushes buffered data to underlying storage.
func (w *BaseWriter) Sync() *appfault.AppError {
	w.Lock()
	defer w.Unlock()

	if w.syncer != nil {
		if err := w.syncer.Sync(); err != nil {
			return appfault.Wrap(errtype.IO, err, "failed to sync writer: "+w.name)
		}
	}

	return nil
}

// Close gracefully closes the underlying destination resource.
func (w *BaseWriter) Close() *appfault.AppError {
	w.Lock()
	defer w.Unlock()

	if w.closer != nil {
		if err := w.closer.Close(); err != nil {
			return appfault.Wrap(errtype.IO, err, "failed to close writer: "+w.name)
		}
	}

	return nil
}
