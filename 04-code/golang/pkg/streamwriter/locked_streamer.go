package streamwriter

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
)

// LockedOptions configures the thread-safe locked streamer for payload type T.
type LockedOptions[T any] struct {
	Name         string
	Destination  io.Writer
	StreamMethod StreamFunc[T]
}

// LockedStreamer implements Streamer[T] with reentrant mutex synchronization and AppError.
type LockedStreamer[T any] struct {
	mu           ReentrantMutex
	name         string
	destination  io.Writer
	streamMethod StreamFunc[T]
}

// NewLockedStreamer constructs a thread-safe streamer over generic type T.
func NewLockedStreamer[T any](opts LockedOptions[T]) *LockedStreamer[T] {
	name := opts.Name
	if name == "" {
		name = "locked-streamer"
	}
	dest := opts.Destination
	if dest == nil {
		dest = os.Stdout
	}

	s := &LockedStreamer[T]{
		name:        name,
		destination: dest,
	}

	if opts.StreamMethod != nil {
		s.streamMethod = opts.StreamMethod
	} else {
		s.streamMethod = s.defaultStream
	}
	return s
}

// Name returns the streamer identifier.
func (s *LockedStreamer[T]) Name() string {
	return s.name
}

// Stream executes the swappable stream method under mutex lock, returning *appfault.AppError.
func (s *LockedStreamer[T]) Stream(ctx context.Context, payload T) *appfault.AppError {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.streamMethod(ctx, payload, s.destination)
}

// Write satisfies Writer[T] by delegating to Stream.
func (s *LockedStreamer[T]) Write(ctx context.Context, payload T) *appfault.AppError {
	return s.Stream(ctx, payload)
}

// SetStreamMethod hot-swaps the streaming logic at runtime.
func (s *LockedStreamer[T]) SetStreamMethod(fn StreamFunc[T]) {
	if fn == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.streamMethod = fn
}

// SetDestination hot-swaps the output destination at runtime.
func (s *LockedStreamer[T]) SetDestination(dest io.Writer) {
	if dest == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.destination = dest
}

// IsLocked reports true for LockedStreamer.
func (s *LockedStreamer[T]) IsLocked() bool {
	return true
}

// Destination returns the active destination under lock.
func (s *LockedStreamer[T]) Destination() io.Writer {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.destination
}

// AsStreamer returns the self-binding Streamer[T].
func (s *LockedStreamer[T]) AsStreamer() Streamer[T] {
	return s
}

// AsWriter returns the self-binding Writer[T].
func (s *LockedStreamer[T]) AsWriter() Writer[T] {
	return s
}

// Lock locks the streamer for exclusive access, satisfying sync.Locker.
func (s *LockedStreamer[T]) Lock() {
	s.mu.Lock()
}

// Unlock unlocks the streamer, satisfying sync.Locker.
func (s *LockedStreamer[T]) Unlock() {
	s.mu.Unlock()
}

func syncDestination(dest any, name string) *appfault.AppError {
	syncer, isOk := dest.(interface{ Sync() error })
	if !isOk {
		return nil
	}
	if err := syncer.Sync(); err != nil {
		return appfault.Wrap(errtype.IO, err, fmt.Sprintf("streamer %s sync failed", name))
	}
	return nil
}

func closeDestination(dest any, name string) *appfault.AppError {
	closer, isOk := dest.(io.Closer)
	if !isOk {
		return nil
	}
	if err := closer.Close(); err != nil {
		return appfault.Wrap(errtype.IO, err, fmt.Sprintf("streamer %s close failed", name))
	}
	return nil
}

// Sync flushes the underlying destination if supported.
func (s *LockedStreamer[T]) Sync() *appfault.AppError {
	s.mu.Lock()
	dest := s.destination
	s.mu.Unlock()

	return syncDestination(dest, s.name)
}

// Close closes the underlying destination if it implements io.Closer.
func (s *LockedStreamer[T]) Close() *appfault.AppError {
	s.mu.Lock()
	dest := s.destination
	s.mu.Unlock()

	return closeDestination(dest, s.name)
}

func (s *LockedStreamer[T]) defaultStream(ctx context.Context, payload T, dest io.Writer) *appfault.AppError {
	compiled := Compile(payload)
	line := fmt.Sprintf("[%s][locked] %s\n", s.name, compiled)
	_, err := dest.Write([]byte(line))
	if err != nil {
		return appfault.Wrap(errtype.IO, err, fmt.Sprintf("streamer %s write failed", s.name))
	}
	return nil
}

var _ Streamer[any] = (*LockedStreamer[any])(nil)
var _ Writer[any] = (*LockedStreamer[any])(nil)
var _ sync.Locker = (*LockedStreamer[any])(nil)
