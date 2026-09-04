# Locked & Lockless Streamers with Self-Binding Interfacer Architecture

> **Status:** Proposal & Architectural Specification  
> **Date:** 2026-09-03  
> **Topic:** 2 Types of Streamers (Locked vs Lockless), Swappable Stream Methods, and Self-Binding `AsInterfacer()` Contracts  
> **Reference:** AUK Go `core/coreinterface/loggerinf`  

---

## 1. Executive Summary & Design Principles

This architectural design addresses two fundamental requirements:
1. **Two Types of Streamers:**
   - **`LockedStreamer` (Thread-Safe):** Uses mutex synchronization for multi-goroutine servers, concurrent HTTP handlers, and worker pools.
   - **`LocklessStreamer` (Zero-Lock Overhead):** Direct execution with zero synchronization cost, designed for single-threaded CLI commands, thread-confined execution, and lock-free channel pipelines.
2. **Self-Binding `AsInterfacer()` Contracts:**
   - Every writer and streamer satisfies self-binding interface methods (`AsWriter()`, `AsStreamer()`, `AsInterfacer()`).
   - Enables clean interface extraction without reflection, compile-time contract enforcement, and fluent unwrapping.
3. **Log-Agnostic & Swappable Stream Method:**
   - Streams can accept both log records and non-log payloads (`any`).
   - The streaming logic (`StreamFunc`) is swappable via `Options` or dynamically at runtime.

---

## 2. Core Interface Contracts (`streamcontract`)

```go
package streamcontract

import (
	"context"
	"io"
)

// Interfacer represents the self-binding contract returning its own interface.
type Interfacer interface {
	AsInterfacer() Interfacer
}

// WriterInterface defines universal write operations with self-binding.
type WriterInterface interface {
	Interfacer
	Write(ctx context.Context, payload any) error
	AsWriter() WriterInterface
	Sync() error
	Close() error
}

// StreamerInterface defines streaming operations with locking introspection.
type StreamerInterface interface {
	Interfacer
	Stream(ctx context.Context, payload any) error
	AsStreamer() StreamerInterface
	AsWriter() WriterInterface
	IsLocked() bool
	Destination() io.Writer
	Sync() error
	Close() error
}

// StreamFunc defines the swappable function signature for streaming data.
type StreamFunc func(ctx context.Context, payload any, dest io.Writer) error
```

---

## 3. Implementation 1: `LockedStreamer` (Thread-Safe with Mutex)

Designed for concurrent execution across multiple goroutines:

```go
package streamer

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"coding-guidelines/common/pkg/streamcontract"
)

type LockedOptions struct {
	Name         string
	Destination  io.Writer
	StreamMethod streamcontract.StreamFunc // Optional custom injection
}

type LockedStreamer struct {
	mu           sync.RWMutex
	name         string
	destination  io.Writer
	streamMethod streamcontract.StreamFunc
}

func NewLocked(opts LockedOptions) *LockedStreamer {
	if opts.Destination == nil {
		opts.Destination = os.Stdout
	}
	s := &LockedStreamer{
		name:        opts.Name,
		destination: opts.Destination,
	}

	if opts.StreamMethod != nil {
		s.streamMethod = opts.StreamMethod
	} else {
		s.streamMethod = s.defaultStream
	}
	return s
}

// Stream locks the mutex before delegating to the swappable stream method.
func (s *LockedStreamer) Stream(ctx context.Context, payload any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.streamMethod(ctx, payload, s.destination)
}

// Write satisfies WriterInterface.
func (s *LockedStreamer) Write(ctx context.Context, payload any) error {
	return s.Stream(ctx, payload)
}

// SetStreamMethod hot-swaps the streaming logic at runtime.
func (s *LockedStreamer) SetStreamMethod(fn streamcontract.StreamFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.streamMethod = fn
}

// Self-binding interface contracts
func (s *LockedStreamer) AsStreamer() streamcontract.StreamerInterface { return s }
func (s *LockedStreamer) AsWriter() streamcontract.WriterInterface     { return s }
func (s *LockedStreamer) AsInterfacer() streamcontract.Interfacer       { return s }

func (s *LockedStreamer) IsLocked() bool          { return true }
func (s *LockedStreamer) Destination() io.Writer  { return s.destination }

func (s *LockedStreamer) defaultStream(ctx context.Context, payload any, dest io.Writer) error {
	line := fmt.Sprintf("[%s][locked] %v\n", s.name, payload)
	_, err := dest.Write([]byte(line))
	return err
}

func (s *LockedStreamer) Sync() error {
	if syncer, isOk := s.destination.(interface{ Sync() error }); isOk {
		return syncer.Sync()
	}
	return nil
}

func (s *LockedStreamer) Close() error {
	if closer, isOk := s.destination.(io.Closer); isOk {
		return closer.Close()
	}
	return nil
}
```

---

## 4. Implementation 2: `LocklessStreamer` (Zero-Lock Overhead)

Designed for single-threaded CLI commands, thread-confined execution, or channel-backed lockless pipelines:

```go
package streamer

import (
	"context"
	"fmt"
	"io"
	"os"

	"coding-guidelines/common/pkg/streamcontract"
)

type LocklessOptions struct {
	Name         string
	Destination  io.Writer
	StreamMethod streamcontract.StreamFunc // Optional custom injection
}

type LocklessStreamer struct {
	name         string
	destination  io.Writer
	streamMethod streamcontract.StreamFunc
}

func NewLockless(opts LocklessOptions) *LocklessStreamer {
	if opts.Destination == nil {
		opts.Destination = os.Stdout
	}
	s := &LocklessStreamer{
		name:        opts.Name,
		destination: opts.Destination,
	}

	if opts.StreamMethod != nil {
		s.streamMethod = opts.StreamMethod
	} else {
		s.streamMethod = s.defaultStream
	}
	return s
}

// Stream executes directly with zero mutex operations.
func (s *LocklessStreamer) Stream(ctx context.Context, payload any) error {
	return s.streamMethod(ctx, payload, s.destination)
}

// Write satisfies WriterInterface.
func (s *LocklessStreamer) Write(ctx context.Context, payload any) error {
	return s.Stream(ctx, payload)
}

// SetStreamMethod swaps the streaming logic.
func (s *LocklessStreamer) SetStreamMethod(fn streamcontract.StreamFunc) {
	s.streamMethod = fn
}

// Self-binding interface contracts
func (s *LocklessStreamer) AsStreamer() streamcontract.StreamerInterface { return s }
func (s *LocklessStreamer) AsWriter() streamcontract.WriterInterface     { return s }
func (s *LocklessStreamer) AsInterfacer() streamcontract.Interfacer       { return s }

func (s *LocklessStreamer) IsLocked() bool          { return false }
func (s *LocklessStreamer) Destination() io.Writer  { return s.destination }

func (s *LocklessStreamer) defaultStream(ctx context.Context, payload any, dest io.Writer) error {
	line := fmt.Sprintf("[%s][lockless] %v\n", s.name, payload)
	_, err := dest.Write([]byte(line))
	return err
}

func (s *LocklessStreamer) Sync() error {
	if syncer, isOk := s.destination.(interface{ Sync() error }); isOk {
		return syncer.Sync()
	}
	return nil
}

func (s *LocklessStreamer) Close() error {
	if closer, isOk := s.destination.(io.Closer); isOk {
		return closer.Close()
	}
	return nil
}
```

---

## 5. Demonstration: Self-Binding, Swapping & Lock Comparisons

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"coding-guidelines/common/pkg/streamcontract"
	"coding-guidelines/common/pkg/streamer"
)

func main() {
	ctx := context.Background()

	// -------------------------------------------------------------
	// 1. Thread-Safe Locked Streamer with Default Method
	// -------------------------------------------------------------
	locked := streamer.NewLocked(streamer.LockedOptions{
		Name: "http-server",
	})

	// Self-binding extraction
	var sInterface streamcontract.StreamerInterface = locked.AsStreamer()
	var wInterface streamcontract.WriterInterface = locked.AsWriter()
	var interfacer streamcontract.Interfacer = locked.AsInterfacer()

	fmt.Printf("Is Locked? %v\n", sInterface.IsLocked()) // true
	_ = wInterface.Write(ctx, "Request from goroutine 1")

	// -------------------------------------------------------------
	// 2. High-Performance Lockless Streamer (Single-Thread / CLI)
	// -------------------------------------------------------------
	lockless := streamer.NewLockless(streamer.LocklessOptions{
		Name: "cli-tool",
	})
	fmt.Printf("Is Locked? %v\n", lockless.IsLocked()) // false
	_ = lockless.Stream(ctx, "Processing file 42")

	// -------------------------------------------------------------
	// 3. Injecting Custom Method via Options (JSON Stream)
	// -------------------------------------------------------------
	customJSONStreamer := streamer.NewLocked(streamer.LockedOptions{
		Name: "json-service",
		StreamMethod: func(ctx context.Context, payload any, dest io.Writer) error {
			bytes, err := json.Marshal(map[string]any{
				"event": payload,
				"mode":  "custom-injected",
			})
			if err != nil {
				return err
			}
			_, err = dest.Write(append(bytes, '\n'))
			return err
		},
	})
	_ = customJSONStreamer.Stream(ctx, map[string]string{"status": "ok"})

	// -------------------------------------------------------------
	// 4. Hot-Swapping Stream Method on the fly
	// -------------------------------------------------------------
	lockless.SetStreamMethod(func(ctx context.Context, payload any, dest io.Writer) error {
		_, err := fmt.Fprintf(dest, ">>> RAW STREAM: %v <<<\n", payload)
		return err
	})
	_ = lockless.Stream(ctx, "Swapped dynamic stream output")
}
```

---

## 6. Architecture Comparison

| Feature | `LockedStreamer` | `LocklessStreamer` |
| :--- | :--- | :--- |
| **Concurrency Safety** | Thread-safe (`sync.RWMutex`) | Single-goroutine / thread-confined |
| **Use Case** | HTTP APIs, background worker pools, multi-thread logs | CLI tools, local pipelines, channel consumers |
| **Overhead** | Mutex lock/unlock (~15-25ns) | Zero lock overhead (0ns) |
| **Swappable Method** | Supported via `Options` and `SetStreamMethod` | Supported via `Options` and `SetStreamMethod` |
| **Self-Binding** | `AsStreamer()`, `AsWriter()`, `AsInterfacer()` | `AsStreamer()`, `AsWriter()`, `AsInterfacer()` |
| **`IsLocked()`** | Returns `true` | Returns `false` |
