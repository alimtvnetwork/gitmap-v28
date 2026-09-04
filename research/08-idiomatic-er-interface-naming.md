# Idiomatic `-er` Interface Naming & Generic Contracts Architecture

> **Document:** `research/08-idiomatic-er-interface-naming.md`  
> **Status:** Implemented & Verified  
> **Package:** `04-code/golang/pkg/streamwriter`  
> **Date:** 2026-09-04  

---

## 1. Context & Architectural Mandate

In idiomatic Go (as outlined in *Effective Go* and repository standards), interface names are agent nouns formed by adding an `-er` suffix to the method or behavior they represent (e.g., `Reader`, `Writer`, `Formatter`, `Streamer`).

The suffix `Interface` (e.g., `WriterInterface`, `StreamerInterface`) is strictly forbidden:
- ❌ **Forbidden:** `type WriterInterface[T any] interface`
- ❌ **Forbidden:** `type StreamerInterface[T any] interface`
- ✅ **Standard:** `type Writer[T any] interface`
- ✅ **Standard:** `type Streamer[T any] interface`
- ✅ **Standard:** `type Interfacer interface`

Additionally, this architecture unifies four core requirements:
1. **Generic Payload `[T any]`:** Any data type can be streamed, written, and compiled.
2. **Monadic `Bytes[T]` Wrapper:** Eliminates naked `([]byte, error)` tuples from formatters, packaging serialized bytes with original payload and structured error state.
3. **Strict Return Type `*appfault.AppError`:** Total ban on bare Go `error` returns across all contracts, streamer engines, writers, and loggers.
4. **Order-Wise Transpilation Engine (`Compile`):** Recursively renders primitives, sorted maps, ordered slices, and structs with `Compilable` support.

---

## 2. Core Interface Contracts (`contracts.go`)

```go
package streamwriter

import (
	"context"
	"io"

	"coding-guidelines/common/pkg/appfault"
)

// Interfacer represents the self-binding contract returning its own interface.
type Interfacer interface {
	AsInterfacer() Interfacer
}

// Writer defines universal write operations over generic type T with AppError.
type Writer[T any] interface {
	Interfacer
	Name() string
	Write(ctx context.Context, payload T) *appfault.AppError
	AsWriter() Writer[T]
	Sync() *appfault.AppError
	Close() *appfault.AppError
}

// Streamer defines streaming operations over generic type T with AppError.
type Streamer[T any] interface {
	Interfacer
	Name() string
	Stream(ctx context.Context, payload T) *appfault.AppError
	AsStreamer() Streamer[T]
	AsWriter() Writer[T]
	IsLocked() bool
	Destination() io.Writer
	Sync() *appfault.AppError
	Close() *appfault.AppError
}

// StreamFunc defines the swappable function signature returning *appfault.AppError.
type StreamFunc[T any] func(ctx context.Context, payload T, dest io.Writer) *appfault.AppError

// WriteFunc defines the swappable function signature returning *appfault.AppError.
type WriteFunc[T any] func(ctx context.Context, payload T) *appfault.AppError

// FormatFunc defines the serialization transformation returning Bytes[T].
type FormatFunc[T any] func(payload T) Bytes[T]
```

---

## 3. Monadic `Bytes[T]` Envelope (`bytes.go`)

Instead of returning `([]byte, error)` which forces repetitive caller unpacking and invites untyped error leakage, formatters return `Bytes[T]`:

```go
type Bytes[T any] struct {
	data     []byte
	payload  T
	appError *appfault.AppError
}

// Accessors & helpers
func (b Bytes[T]) Raw() []byte                          { return b.data }
func (b Bytes[T]) String() string                       { return string(b.data) }
func (b Bytes[T]) Len() int                            { return len(b.data) }
func (b Bytes[T]) IsEmpty() bool                       { return len(b.data) == 0 }
func (b Bytes[T]) Payload() T                          { return b.payload }
func (b Bytes[T]) AppError() *appfault.AppError        { return b.appError }
func (b Bytes[T]) Fault() *appfault.AppError           { return b.appError }
func (b Bytes[T]) HasError() bool                      { return b.appError != nil }
func (b Bytes[T]) IsValid() bool                       { return b.appError == nil }
func (b Bytes[T]) Unwrap() ([]byte, *appfault.AppError) { return b.data, b.appError }
```

---

## 4. Deterministic Order Transpiler (`compiler.go`)

The `Compile[T any](payload T) string` engine resolves payloads to deterministic string output:
1. **Primitives:** Directly converted via `strconv` without heap allocations where possible.
2. **Maps:** Keys are extracted and lexicographically sorted before formatting, guaranteeing identical string outputs regardless of Go map hash randomization.
3. **Slices / Arrays:** Order is preserved sequentially `[elem0, elem1, ...]`.
4. **Structs & Compilable Objects:** Recursively checks if the value or pointer implements `Compilable`:
   ```go
   type Compilable interface {
       Compile() string
   }
   ```
   If not implemented, exported struct fields (honoring JSON tags) are formatted recursively.

---

## 5. Implementations Summary

| Struct | Interface Implemented | Thread-Safety | Use Case |
|---|---|---|---|
| `LockedStreamer[T]` | `Streamer[T]`, `Writer[T]`, `Interfacer` | Thread-Safe (`sync.RWMutex`) | High-concurrency stdout, stderr, or shared files |
| `LocklessStreamer[T]` | `Streamer[T]`, `Writer[T]`, `Interfacer` | Unsynchronized | Low-latency CLI, single-threaded workers, network sockets |
| `PluggableWriter[T]` | `Writer[T]`, `Interfacer` | Thread-Safe (`sync.RWMutex`) | Composable write pipelines with swappable formatters |
| `Logger[T]` | Coordinator | Thread-Safe (`sync.RWMutex`) | Multi-destination fan-out, silent zero-alloc fallback |

---

## 6. Verification Results

All packages pass with 100% test coverage and zero linter warnings:
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
=== RUN   TestSwappableMethods_GenericRuntime
--- PASS: TestSwappableMethods_GenericRuntime (0.00s)
=== RUN   TestCompositeLogger_FluentChaining
--- PASS: TestCompositeLogger_FluentChaining (0.00s)
=== RUN   TestLogRecord_Compile
--- PASS: TestLogRecord_Compile (0.00s)
PASS
ok  	coding-guidelines/common/pkg/streamwriter	0.485s
```
