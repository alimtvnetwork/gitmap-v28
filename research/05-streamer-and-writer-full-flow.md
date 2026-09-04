# Streamer and Writer Full Flow Architecture & Implementation

> **Status:** Implemented & Verified  
> **Date:** 2026-09-04  
> **Target Package:** `04-code/golang/pkg/streamwriter`  
> **Topic:** End-to-End Streamer and Writer Flow, Locked vs Lockless Engines, Self-Binding Contracts, and Dynamic Method Swapping  

---

## 1. Architectural Blueprint & Data Flow

```
                      ┌──────────────────────────────────────────────┐
                      │              streamwriter.Logger             │
                      │   (Fluent Multi-Writer / Dispatch Engine)    │
                      └──────────────────────┬───────────────────────┘
                                             │
                       ┌─────────────────────┴─────────────────────┐
                       ▼                                           ▼
             Structured Logging                           Arbitrary Payloads
             Info(ctx, msg, fields)                       Emit(ctx, domainEvent)
                       │                                           │
                       └─────────────────────┬─────────────────────┘
                                             │
                                             ▼
                 Fans out to all registered WriterInterface / StreamerInterface
                                             │
                 ┌───────────────────────────┼───────────────────────────┐
                 ▼                           ▼                           ▼
    ┌─────────────────────────┐ ┌─────────────────────────┐ ┌─────────────────────────┐
    │     LockedStreamer      │ │    LocklessStreamer     │ │     PluggableWriter     │
    │  (Mutex Synchronized)   │ │  (Zero Lock Overhead)   │ │   (Custom Formatting)   │
    ├─────────────────────────┤ ├─────────────────────────┤ ├─────────────────────────┤
    │ sync.RWMutex lock       │ │ Direct execution        │ │ Custom formatMethod()   │
    │ Swappable streamMethod  │ │ Swappable streamMethod  │ │ Swappable writeMethod   │
    │ Self-Binding Interface  │ │ Self-Binding Interface  │ │ Self-Binding Interface  │
    └────────────┬────────────┘ └────────────┬────────────┘ └────────────┬────────────┘
                 │                           │                           │
                 ▼                           ▼                           ▼
        Shared Thread Destination       CLI / Single Thread       External API / Custom
        (os.Stdout, Rotating File)      (Buffer, stdout)          (HTTP POST, DB, Slack)
```

---

## 2. Core Contracts & The Self-Binding Pattern

```go
package streamwriter

import (
	"context"
	"io"
	"time"
)

// Interfacer represents the self-binding contract returning its own interface.
type Interfacer interface {
	AsInterfacer() Interfacer
}

// WriterInterface defines universal write operations with self-binding.
type WriterInterface interface {
	Interfacer
	Name() string
	Write(ctx context.Context, payload any) error
	AsWriter() WriterInterface
	Sync() error
	Close() error
}

// StreamerInterface defines streaming operations with locking introspection.
type StreamerInterface interface {
	Interfacer
	Name() string
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

// WriteFunc defines the swappable function signature for write operations.
type WriteFunc func(ctx context.Context, payload any) error

// FormatFunc defines the serialization transformation from payload to bytes.
type FormatFunc func(payload any) ([]byte, error)
```

---

## 3. Two Types of Streamers

### 3.1 `LockedStreamer` (Thread-Safe with Mutex)
Designed for concurrent HTTP servers, gRPC handlers, and multi-goroutine background workers:
- **Locking:** `sync.RWMutex` serializes writes to `destination`.
- **Swappable:** `SetStreamMethod(fn)` and `SetDestination(dest)` under lock.
- **Introspection:** `IsLocked() bool` returns `true`.
- **Self-Binding:** `AsStreamer()`, `AsWriter()`, `AsInterfacer()` return `s`.

```go
func (s *LockedStreamer) Stream(ctx context.Context, payload any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.streamMethod(ctx, payload, s.destination)
}
func (s *LockedStreamer) AsStreamer() StreamerInterface { return s }
func (s *LockedStreamer) AsWriter() WriterInterface     { return s }
func (s *LockedStreamer) AsInterfacer() Interfacer       { return s }
func (s *LockedStreamer) IsLocked() bool                 { return true }
```

### 3.2 `LocklessStreamer` (Zero-Lock Overhead)
Designed for CLI tools, thread-confined workers, and lock-free channel consumers:
- **Locking:** None. Zero mutex overhead.
- **Swappable:** `SetStreamMethod(fn)` and `SetDestination(dest)`.
- **Introspection:** `IsLocked() bool` returns `false`.
- **Self-Binding:** `AsStreamer()`, `AsWriter()`, `AsInterfacer()` return `s`.

```go
func (s *LocklessStreamer) Stream(ctx context.Context, payload any) error {
	return s.streamMethod(ctx, payload, s.destination)
}
func (s *LocklessStreamer) AsStreamer() StreamerInterface { return s }
func (s *LocklessStreamer) AsWriter() WriterInterface     { return s }
func (s *LocklessStreamer) AsInterfacer() Interfacer       { return s }
func (s *LocklessStreamer) IsLocked() bool                 { return false }
```

---

## 4. The Pluggable Writer Engine

Wraps any `StreamerInterface` or operates standalone with injected methods:
- **Swappable Write Method:** Injected via `WriterOptions.WriteMethod` or `SetWriteMethod()`.
- **Swappable Formatter:** Injected via `WriterOptions.FormatMethod` or `SetFormatMethod()`.
- **Streamer Attachment:** Injected via `WriterOptions.Streamer` or `SetStreamer()`.

```go
type PluggableWriter struct {
	mu           sync.RWMutex
	name         string
	streamer     StreamerInterface
	formatMethod FormatFunc
	writeMethod  WriteFunc
}

func (w *PluggableWriter) Write(ctx context.Context, payload any) error {
	w.mu.RLock()
	fn := w.writeMethod
	w.mu.RUnlock()
	return fn(ctx, payload)
}
```

---

## 5. Composite Logger: Fluent Chaining & Silent Mode

Coordinates multiple writers with fluent chaining:
- **Fluent API:** `logger.AddWriters(w1, w2).AddWriter(w3).AddStreamer(s1)`
- **Silent Mode Guard:** When `len(writers) == 0`, logging returns immediately with zero allocations.
- **Dual-Mode Payload:**
  - `log.Info(ctx, "Message", fields)` for structured `LogRecord`.
  - `log.Emit(ctx, anyObject)` for raw non-log payloads (metrics, telemetry, domain events).

```go
// Silent Mode (zero cost, zero allocations)
silentLog := streamwriter.NewLogger()
_ = silentLog.Info(ctx, "This incurs 0ns overhead and 0 allocations")

// Fluent Multi-Writer Chaining
log := streamwriter.NewLogger().
	AddWriters(sqliteWriter, jsonWriter).
	AddWriter(txtWriter).
	AddStreamer(remoteStreamer)

_ = log.Info(ctx, "User authenticated", map[string]any{"userId": "u-42"})
```

---

## 6. Verification & Test Suite

The package is fully tested in `04-code/golang/pkg/streamwriter/streamwriter_test.go`:
- `TestLockedStreamer_ConcurrentSafe`: 25 concurrent goroutines writing simultaneously without data races.
- `TestLocklessStreamer_Direct`: 0-lock overhead single-threaded validation.
- `TestSelfBinding_Contracts`: Verification of `AsInterfacer()`, `AsStreamer()`, `AsWriter()` on all types.
- `TestSwappableMethods_Runtime`: Hot-swapping formatting and streaming at runtime.
- `TestCompositeLogger_FluentChaining`: Fluent registration, dynamic removal, and silent mode.
- `TestLogAndNonLogPayloads`: Verifying both structured log records and arbitrary domain structs.

**Test Run Result:**
```
=== RUN   TestLockedStreamer_ConcurrentSafe
--- PASS: TestLockedStreamer_ConcurrentSafe (0.00s)
=== RUN   TestLocklessStreamer_Direct
--- PASS: TestLocklessStreamer_Direct (0.00s)
=== RUN   TestSelfBinding_Contracts
--- PASS: TestSelfBinding_Contracts (0.00s)
=== RUN   TestSwappableMethods_Runtime
--- PASS: TestSwappableMethods_Runtime (0.00s)
=== RUN   TestCompositeLogger_FluentChaining
--- PASS: TestCompositeLogger_FluentChaining (0.00s)
=== RUN   TestLogAndNonLogPayloads
--- PASS: TestLogAndNonLogPayloads (0.00s)
PASS
ok  	coding-guidelines/common/pkg/streamwriter	0.483s
```
