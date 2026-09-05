package streamwriter_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
	"coding-guidelines/common/pkg/streamwriter"
)

// MockTargetWriter records written items and flushes for test verification.
type MockTargetWriter[T any] struct {
	mu       sync.Mutex
	items    []T
	synced   int
	closed   bool
	failNext bool
}

func (m *MockTargetWriter[T]) Name() string {
	return "mock-target"
}

func (m *MockTargetWriter[T]) AsWriter() streamwriter.Writer[T] {
	return m
}

func (m *MockTargetWriter[T]) Lock() {
	m.mu.Lock()
}

func (m *MockTargetWriter[T]) Unlock() {
	m.mu.Unlock()
}

func (m *MockTargetWriter[T]) Write(ctx context.Context, payload T) *appfault.AppError {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.failNext {
		return appfault.New(errtype.Internal, "simulated write failure")
	}

	m.items = append(m.items, payload)

	return nil
}

func (m *MockTargetWriter[T]) Sync() *appfault.AppError {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.synced++

	return nil
}

func (m *MockTargetWriter[T]) Close() *appfault.AppError {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.closed = true

	return nil
}

func (m *MockTargetWriter[T]) Items() []T {
	m.mu.Lock()
	defer m.mu.Unlock()

	copied := make([]T, len(m.items))
	copy(copied, m.items)

	return copied
}

func TestAsyncWriter_BasicFlow(t *testing.T) {
	mock := &MockTargetWriter[string]{}
	opts := streamwriter.AsyncWriterOptions{
		Name:          "test-async",
		BufferSize:    32,
		FlushInterval: 20 * time.Millisecond,
	}

	aw := streamwriter.NewAsyncWriter[string](mock, opts)

	if aw.Name() != "test-async" {
		t.Errorf("expected name 'test-async', got %s", aw.Name())
	}

	ctx := context.Background()

	for i := 0; i < 5; i++ {
		err := aw.Write(ctx, "item")
		if err != nil {
			t.Fatalf("unexpected write error: %v", err)
		}
	}

	if err := aw.Sync(); err != nil {
		t.Fatalf("sync error: %v", err)
	}

	items := mock.Items()
	if len(items) != 5 {
		t.Errorf("expected 5 items received by target, got %d", len(items))
	}

	if err := aw.Close(); err != nil {
		t.Fatalf("close error: %v", err)
	}

	mock.mu.Lock()
	closed := mock.closed
	mock.mu.Unlock()

	if !closed {
		t.Errorf("expected target to be closed")
	}
}

func TestAsyncWriter_DropOnFull(t *testing.T) {
	mock := &MockTargetWriter[int]{}
	opts := streamwriter.AsyncWriterOptions{
		Name:          "drop-writer",
		BufferSize:    2,
		FlushInterval: 500 * time.Millisecond,
		DropOnFull:    true,
	}

	aw := streamwriter.NewAsyncWriter[int](mock, opts)

	ctx := context.Background()

	// Rapidly write more items than buffer size
	for i := 0; i < 20; i++ {
		_ = aw.Write(ctx, i)
	}

	_ = aw.Sync()
	_ = aw.Close()

	if aw.DroppedCount() < 0 {
		t.Errorf("invalid dropped count: %d", aw.DroppedCount())
	}
}

func TestAsyncWriter_OnErrorCallback(t *testing.T) {
	var errorCaptured atomic.Bool

	mock := &MockTargetWriter[string]{
		failNext: true,
	}

	opts := streamwriter.AsyncWriterOptions{
		BufferSize: 8,
		OnError: func(err *appfault.AppError) {
			if err != nil {
				errorCaptured.Store(true)
			}
		},
	}

	aw := streamwriter.NewAsyncWriter[string](mock, opts)
	_ = aw.Write(context.Background(), "will-fail")

	_ = aw.Sync()
	_ = aw.Close()

	if !errorCaptured.Load() {
		t.Errorf("expected OnError to be called when target fails")
	}
}

func TestAnyAsyncWriter_Creation(t *testing.T) {
	mock := &MockTargetWriter[any]{}
	opts := streamwriter.AsyncWriterOptions{
		BufferSize: 16,
	}

	aw := streamwriter.NewAnyAsyncWriter(mock, opts)
	if aw == nil {
		t.Fatalf("expected non-nil AnyAsyncWriter")
	}

	if aw.Name() != "async-writer" {
		t.Errorf("expected default name 'async-writer', got %s", aw.Name())
	}

	_ = aw.Write(context.Background(), map[string]string{"k": "v"})
	_ = aw.Sync()
	_ = aw.Close()

	items := mock.Items()
	if len(items) != 1 {
		t.Errorf("expected 1 item, got %d", len(items))
	}
}
