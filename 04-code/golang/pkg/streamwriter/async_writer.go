package streamwriter

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
)

// AsyncWriterOptions configures the asynchronous writer buffer and background worker.
type AsyncWriterOptions struct {
	Name          string
	BufferSize    int           // Channel buffer size (default: 256)
	FlushInterval time.Duration // Maximum interval between flushes (default: 50ms)
	DropOnFull    bool          // If true, drops items when buffer is full; if false, blocks until space is available
	OnError       ErrorHandlerFunc
}

// AsyncWriter wraps any Writer[T] with a non-blocking buffered ring channel and worker goroutine.
type AsyncWriter[T any] struct {
	name      string
	target    Writer[T]
	opts      AsyncWriterOptions
	queue     chan T
	closed    atomic.Bool
	dropped   atomic.Int64
	closeOnce sync.Once
	doneChan  chan struct{}
	wg        sync.WaitGroup
	mu        ReentrantMutex
}

// AnyAsyncWriter is the first-class non-generic alias for AsyncWriter[any].
type AnyAsyncWriter = AsyncWriter[any]

// NewAsyncWriter constructs an asynchronous non-blocking writer around target.
func NewAsyncWriter[T any](target Writer[T], opts AsyncWriterOptions) *AsyncWriter[T] {
	name := opts.Name
	if name == "" {
		name = "async-writer"
	}

	bufSize := opts.BufferSize
	if bufSize <= 0 {
		bufSize = 256
	}

	flushInterval := opts.FlushInterval
	if flushInterval <= 0 {
		flushInterval = 50 * time.Millisecond
	}

	opts.Name = name
	opts.BufferSize = bufSize
	opts.FlushInterval = flushInterval

	aw := &AsyncWriter[T]{
		name:     name,
		target:   target,
		opts:     opts,
		queue:    make(chan T, bufSize),
		doneChan: make(chan struct{}),
	}

	aw.wg.Add(1)
	go aw.worker()

	return aw
}

// NewAnyAsyncWriter constructs an asynchronous non-generic AnyAsyncWriter (*AsyncWriter[any]).
func NewAnyAsyncWriter(target Writer[any], opts AsyncWriterOptions) *AnyAsyncWriter {
	return NewAsyncWriter[any](target, opts)
}

// Name returns the writer identifier.
func (aw *AsyncWriter[T]) Name() string {
	return aw.name
}

// AsWriter returns the Writer[T] interface representation.
func (aw *AsyncWriter[T]) AsWriter() Writer[T] {
	return aw
}

// Lock acquires the synchronization lock.
func (aw *AsyncWriter[T]) Lock() {
	aw.mu.Lock()
}

// Unlock releases the synchronization lock.
func (aw *AsyncWriter[T]) Unlock() {
	aw.mu.Unlock()
}

// Write enqueues a payload into the async buffer without blocking the caller.
func (aw *AsyncWriter[T]) Write(ctx context.Context, payload T) *appfault.AppError {
	if aw.closed.Load() {
		return appfault.New(errtype.Precondition, "async writer is closed")
	}

	if aw.opts.DropOnFull {
		select {
		case aw.queue <- payload:
			return nil
		default:
			aw.dropped.Add(1)

			return nil
		}
	}

	select {
	case aw.queue <- payload:
		return nil
	case <-ctx.Done():
		return appfault.Wrap(errtype.Timeout, ctx.Err(), "async write context canceled")
	}
}

// DroppedCount returns the number of dropped items due to buffer overflow in DropOnFull mode.
func (aw *AsyncWriter[T]) DroppedCount() int64 {
	return aw.dropped.Load()
}

// QueueLen returns the current number of pending items in the channel buffer.
func (aw *AsyncWriter[T]) QueueLen() int {
	return len(aw.queue)
}

// Target returns the wrapped underlying Writer[T].
func (aw *AsyncWriter[T]) Target() Writer[T] {
	return aw.target
}

// worker coordinates queue drainage and background flushes.
func (aw *AsyncWriter[T]) worker() {
	defer aw.wg.Done()

	ticker := time.NewTicker(aw.opts.FlushInterval)
	defer ticker.Stop()

	ctx := context.Background()

	for {
		select {
		case item, ok := <-aw.queue:
			if !ok {
				aw.drainRemaining(ctx)

				return
			}

			aw.dispatchToTarget(ctx, item)

		case <-ticker.C:
			if aw.target != nil {
				_ = aw.target.Sync()
			}

		case <-aw.doneChan:
			aw.drainRemaining(ctx)

			return
		}
	}
}

// dispatchToTarget writes an item to the target writer, handling errors via OnError callback.
func (aw *AsyncWriter[T]) dispatchToTarget(ctx context.Context, item T) {
	if aw.target == nil {
		return
	}

	if err := aw.target.Write(ctx, item); err != nil {
		if aw.opts.OnError != nil {
			aw.opts.OnError(err)
		}
	}
}

// drainRemaining flushes all remaining items in the channel before shutdown.
func (aw *AsyncWriter[T]) drainRemaining(ctx context.Context) {
	for {
		select {
		case item, ok := <-aw.queue:
			if !ok {
				if aw.target != nil {
					_ = aw.target.Sync()
				}

				return
			}

			aw.dispatchToTarget(ctx, item)
		default:
			if aw.target != nil {
				_ = aw.target.Sync()
			}

			return
		}
	}
}

// Sync forces pending items to be drained and flushes the target writer.
func (aw *AsyncWriter[T]) Sync() *appfault.AppError {
	// Wait until channel buffer is completely drained
	for len(aw.queue) > 0 {
		time.Sleep(2 * time.Millisecond)
	}

	if aw.target != nil {
		return aw.target.Sync()
	}

	return nil
}

// Close gracefully flushes all buffered items, stops the worker, and closes the target writer.
func (aw *AsyncWriter[T]) Close() *appfault.AppError {
	var closeErr *appfault.AppError

	aw.closeOnce.Do(func() {
		aw.closed.Store(true)
		close(aw.doneChan)
		aw.wg.Wait()

		if aw.target != nil {
			closeErr = aw.target.Close()
		}
	})

	return closeErr
}
