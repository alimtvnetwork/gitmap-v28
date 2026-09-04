package streamwriter

import (
	"context"
	"sync"
	"time"

	"coding-guidelines/common/pkg/appfault"
)

// Logger coordinates multiple generic writers and streamers over type T with AppError returns.
type Logger[T any] struct {
	mu      sync.RWMutex
	writers []Writer[T]
}

// NewLogger creates an empty generic Logger in silent mode (0 writers, 0 allocations).
func NewLogger[T any]() *Logger[T] {
	return &Logger[T]{
		writers: make([]Writer[T], 0),
	}
}

// AddWriter fluently registers a single writer.
func (l *Logger[T]) AddWriter(w Writer[T]) *Logger[T] {
	if w == nil {
		return l
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.writers = append(l.writers, w.AsWriter())
	return l
}

// AddWriters fluently registers multiple writers in one call.
func (l *Logger[T]) AddWriters(ws ...Writer[T]) *Logger[T] {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, w := range ws {
		if w != nil {
			l.writers = append(l.writers, w.AsWriter())
		}
	}
	return l
}

// AddStreamer fluently registers a streamer (adapting it via AsWriter()).
func (l *Logger[T]) AddStreamer(s Streamer[T]) *Logger[T] {
	if s == nil {
		return l
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.writers = append(l.writers, s.AsWriter())
	return l
}

// ClearWriters removes all registered writers (switches to silent mode).
func (l *Logger[T]) ClearWriters() *Logger[T] {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.writers = l.writers[:0]
	return l
}

// RemoveWriter removes a registered writer by name.
func (l *Logger[T]) RemoveWriter(name string) *Logger[T] {
	l.mu.Lock()
	defer l.mu.Unlock()
	filtered := make([]Writer[T], 0, len(l.writers))
	for _, w := range l.writers {
		if w.Name() != name {
			filtered = append(filtered, w)
		}
	}
	l.writers = filtered
	return l
}

// WriterCount returns the number of active writers.
func (l *Logger[T]) WriterCount() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.writers)
}

// Emit sends a generic payload T to all active writers, returning *appfault.AppError.
func (l *Logger[T]) Emit(ctx context.Context, payload T) *appfault.AppError {
	l.mu.RLock()
	// Zero-allocation silent guard
	if len(l.writers) == 0 {
		l.mu.RUnlock()
		return nil
	}
	active := make([]Writer[T], len(l.writers))
	copy(active, l.writers)
	l.mu.RUnlock()

	var firstErr *appfault.AppError
	for _, w := range active {
		err := w.Write(ctx, payload)
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// Info emits a LevelInfo structured record if T can represent it, returning *appfault.AppError.
func (l *Logger[T]) Info(ctx context.Context, msg string, fields ...map[string]any) *appfault.AppError {
	return l.dispatchRecord(ctx, LevelInfo, msg, fields...)
}

// Error emits a LevelError structured record if T can represent it, returning *appfault.AppError.
func (l *Logger[T]) Error(ctx context.Context, msg string, fields ...map[string]any) *appfault.AppError {
	return l.dispatchRecord(ctx, LevelError, msg, fields...)
}

// Debug emits a LevelDebug structured record if T can represent it, returning *appfault.AppError.
func (l *Logger[T]) Debug(ctx context.Context, msg string, fields ...map[string]any) *appfault.AppError {
	return l.dispatchRecord(ctx, LevelDebug, msg, fields...)
}

// Warn emits a LevelWarn structured record if T can represent it, returning *appfault.AppError.
func (l *Logger[T]) Warn(ctx context.Context, msg string, fields ...map[string]any) *appfault.AppError {
	return l.dispatchRecord(ctx, LevelWarn, msg, fields...)
}

// Sync flushes all active writers, returning *appfault.AppError.
func (l *Logger[T]) Sync() *appfault.AppError {
	l.mu.RLock()
	active := make([]Writer[T], len(l.writers))
	copy(active, l.writers)
	l.mu.RUnlock()

	var firstErr *appfault.AppError
	for _, w := range active {
		if err := w.Sync(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// Close closes all active writers, returning *appfault.AppError.
func (l *Logger[T]) Close() *appfault.AppError {
	l.mu.Lock()
	defer l.mu.Unlock()

	var firstErr *appfault.AppError
	for _, w := range l.writers {
		if err := w.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	l.writers = l.writers[:0]
	return firstErr
}

func (l *Logger[T]) dispatchRecord(ctx context.Context, lvl LogLevel, msg string, fields ...map[string]any) *appfault.AppError {
	l.mu.RLock()
	if len(l.writers) == 0 {
		l.mu.RUnlock()
		return nil
	}
	l.mu.RUnlock()

	traceID := ""
	userID := ""
	if ctx != nil {
		if tid, isOk := ctx.Value("traceId").(string); isOk {
			traceID = tid
		}
		if uid, isOk := ctx.Value("userId").(string); isOk {
			userID = uid
		}
	}

	merged := make(map[string]any)
	for _, f := range fields {
		for k, v := range f {
			merged[k] = v
		}
	}

	record := LogRecord{
		Timestamp: time.Now().UTC(),
		Level:     lvl,
		Message:   msg,
		Context:   ctx,
		Fields:    merged,
		TraceID:   traceID,
		UserID:    userID,
	}

	// Case 1: T is any or LogRecord
	if payload, isOk := any(record).(T); isOk {
		return l.Emit(ctx, payload)
	}

	// Case 2: T is string (compiled representation)
	if strPayload, isOk := any(record.Compile()).(T); isOk {
		return l.Emit(ctx, strPayload)
	}

	return nil
}
