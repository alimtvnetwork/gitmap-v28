# Writer `sync.Locker` Integration & Deprecation of `Interfacer`

> **Document:** `research/09-writer-locker-and-avoiding-interfacer.md`  
> **Status:** Implemented & Verified  
> **Package:** `04-code/golang/pkg/streamwriter`  
> **Date:** 2026-09-04  

---

## 1. Executive Summary: Why Avoid `Interfacer`?

### The Question: What is the Value or Benefit of `Interfacer`?
In earlier design explorations, `type Interfacer interface { AsInterfacer() Interfacer }` was tested as a mechanism for self-binding and uniform interface extraction.

**Verdict: In Go, `Interfacer` has zero practical value and introduces anti-patterns.**

### Technical Breakdown: Why `Interfacer` is an Anti-Pattern in Go

1. **Implicit Structural Typing Makes Marker Interfaces Obsolete:**
   - In languages like Java or C#, explicit interface declarations (`class Foo implements IFoo, IInterfacer`) are required for type polymorphism.
   - In Go, interfaces are satisfied **implicitly**. If a struct implements `Write(...)`, it is already a `Writer[T]`. Adding `AsInterfacer()` provides no additional type safety, no compile-time enforcement, and no runtime capability that Go does not already provide natively.
2. **`any` Already Captures All Types:**
   - Go's universal type is `any` (`interface{}`). Any component can already be stored, passed, or inspected dynamically without an artificial `Interfacer` wrapper.
3. **Forces Pointless Downcasts:**
   - Calling `w.AsInterfacer()` returns `Interfacer`, an interface with no operational methods. To do anything useful with it, the caller must immediately downcast via type assertion:
     ```go
     // Pointless ceremony:
     w.AsInterfacer().(Writer[string]).Write(...)
     ```
   - This degrades performance and bypasses compile-time safety.
4. **Violates the Interface Segregation Principle (ISP):**
   - Go interfaces should be small, focused abstractions describing concrete behavior (`io.Reader`, `io.Writer`, `sync.Locker`). Adding artificial marker methods bloats interfaces with non-functional boilerplate.

**Decision:** Completely eliminated `Interfacer` and `AsInterfacer()` from all contracts, structs, and tests.

---

## 2. High Value: `Lock()` and `Unlock()` (`sync.Locker`)

While `Interfacer` adds friction, adding `Lock()` and `Unlock()` to `Writer[T]` and `Streamer[T]` provides immense architectural value:

1. **Native Compatibility with `sync.Locker`:**
   `Writer[T]` and `Streamer[T]` directly satisfy `sync.Locker` from the Go standard library:
   ```go
   var locker sync.Locker = writer
   ```
2. **Atomic Compound Batches:**
   When multiple goroutines emit multi-line logs or related records, individual writes can interleave. With `Lock()` and `Unlock()`, callers guarantee atomic compound execution:
   ```go
   writer.Lock()
   _ = writer.Write(ctx, "ORDER_START: #1001")
   _ = writer.Write(ctx, "ORDER_ITEM: Widget A")
   _ = writer.Write(ctx, "ORDER_END: #1001")
   writer.Unlock()
   ```
   No other goroutines can interleave log records between `ORDER_START` and `ORDER_END`.
3. **Zero Overhead for Lockless Implementations:**
   `LocklessStreamer[T]` implements `Lock()` and `Unlock()` as no-ops, satisfying the contract with 0 CPU cycles.

---

## 3. Updated Interface Contracts (`contracts.go`)

```go
package streamwriter

import (
	"context"
	"io"

	"coding-guidelines/common/pkg/appfault"
)

// Writer defines universal write operations over generic type T with AppError and Locker synchronization.
type Writer[T any] interface {
	Name() string
	Write(ctx context.Context, payload T) *appfault.AppError
	AsWriter() Writer[T]
	Lock()
	Unlock()
	Sync() *appfault.AppError
	Close() *appfault.AppError
}

// Streamer defines streaming operations over generic type T with AppError and Locker synchronization.
type Streamer[T any] interface {
	Name() string
	Stream(ctx context.Context, payload T) *appfault.AppError
	AsStreamer() Streamer[T]
	AsWriter() Writer[T]
	IsLocked() bool
	Lock()
	Unlock()
	Destination() io.Writer
	Sync() *appfault.AppError
	Close() *appfault.AppError
}
```

---

## 4. Deadlock-Free Concurrency (`mutex.go`)

Standard `sync.Mutex` and `sync.RWMutex` in Go are **non-reentrant**. If a caller invokes:
```go
writer.Lock()
writer.Write(...) // If Write also acquires the same lock internally -> DEADLOCK!
writer.Unlock()
```
To eliminate deadlocks while supporting both:
1. Direct unsynchronized calls `w.Write(...)` (automatically thread-safe)
2. Caller-locked batch sessions `w.Lock() ... w.Write() ... w.Unlock()`

A lightweight `ReentrantMutex` (`mutex.go`) tracks the current goroutine ID:
- If current goroutine already holds the lock, recursion counter increments without blocking.
- When recursion counter returns to zero, the underlying mutex is released.

---

## 5. Verification Results

```bash
$ go test ./pkg/streamwriter -v -count=1
=== RUN   TestBytesWrapper
--- PASS: TestBytesWrapper (0.00s)
=== RUN   TestCompiler_Primitives
--- PASS: TestCompiler_Primitives (0.00s)
=== RUN   TestCompiler_Maps_OrderWise
--- PASS: TestCompiler_Maps_OrderWise (0.00s)
=== RUN   TestCompiler_Slices_OrderWise
--- PASS: TestCompiler_Slices_OrderWise (0.00s)
=== RUN   TestCompiler_ObjectAndRecursiveCompilable
--- PASS: TestCompiler_ObjectAndRecursiveCompilable (0.00s)
=== RUN   TestLockedStreamer_Generic_ConcurrentSafe
--- PASS: TestLockedStreamer_Generic_ConcurrentSafe (0.00s)
=== RUN   TestLocklessStreamer_Generic_Direct
--- PASS: TestLocklessStreamer_Generic_Direct (0.00s)
=== RUN   TestSelfBinding_GenericContracts
--- PASS: TestSelfBinding_GenericContracts (0.00s)
=== RUN   TestWriter_LockerSynchronization
--- PASS: TestWriter_LockerSynchronization (0.00s)
=== RUN   TestWriter_ConcurrentCompoundBatches
--- PASS: TestWriter_ConcurrentCompoundBatches (0.00s)
=== RUN   TestSwappableMethods_GenericRuntime
--- PASS: TestSwappableMethods_GenericRuntime (0.00s)
=== RUN   TestCompositeLogger_FluentChaining
--- PASS: TestCompositeLogger_FluentChaining (0.00s)
=== RUN   TestLogRecord_Compile
--- PASS: TestLogRecord_Compile (0.00s)
PASS
ok  	coding-guidelines/common/pkg/streamwriter	0.525s
```
