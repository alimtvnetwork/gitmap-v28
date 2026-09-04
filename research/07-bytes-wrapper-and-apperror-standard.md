# `Bytes[T]` Wrapper Type & Mandatory `*appfault.AppError` Architecture

> **Status:** Implemented & Verified  
> **Date:** 2026-09-04  
> **Package:** `04-code/golang/pkg/streamwriter`  
> **Topic:** Monadic `Bytes[T]` Result Type, Elimination of Bare `error` Returns, and Strict `*appfault.AppError` Compliance  

---

## 1. Executive Summary & Core Mandates

In accordance with repo rule 6 (`*appfault.AppError` standard) and user architecture requirements:
1. **Total Elimination of Bare Go `error`:**
   - Every interface, streamer, writer, logger, and functional option must return `*appfault.AppError` rather than bare `error`.
   - Any low-level standard library I/O or serialization error is wrapped via `appfault.Wrap(errtype.IO, err, msg)`.
2. **`Bytes[T]` Monadic Wrapper (Replacing `([]byte, error)`):**
   - The conventional tuple return `([]byte, error)` is replaced by `Bytes[T any]`.
   - `Bytes[T]` encapsulates the formatted byte slice, the original generic payload `T`, and monadic `*appfault.AppError` state.
   - Provides a comprehensive method suite: `.Raw()`, `.String()`, `.Len()`, `.IsEmpty()`, `.Payload()`, `.AppError()`, `.IsValid()`, and `.Unwrap()`.

---

## 2. The `Bytes[T any]` Wrapper Type (`pkg/streamwriter/bytes.go`)

```go
package streamwriter

import (
	"coding-guidelines/common/pkg/appfault"
)

// Bytes wraps a formatted byte slice bundled with its generic payload T and monadic AppError state.
type Bytes[T any] struct {
	data     []byte
	payload  T
	appError *appfault.AppError
}

// NewBytes creates a successful Bytes envelope.
func NewBytes[T any](data []byte, payload T) Bytes[T] {
	return Bytes[T]{
		data:    data,
		payload: payload,
	}
}

// NewBytesError creates a failed Bytes envelope with an AppError.
func NewBytesError[T any](appErr *appfault.AppError) Bytes[T] {
	return Bytes[T]{
		appError: appErr,
	}
}

// NewBytesErrorWithPayload creates a failed Bytes envelope preserving the original payload.
func NewBytesErrorWithPayload[T any](appErr *appfault.AppError, payload T) Bytes[T] {
	return Bytes[T]{
		payload:  payload,
		appError: appErr,
	}
}

func (b Bytes[T]) Raw() []byte                { return b.data }
func (b Bytes[T]) Bytes() []byte              { return b.data }
func (b Bytes[T]) String() string             { return string(b.data) }
func (b Bytes[T]) Len() int                   { return len(b.data) }
func (b Bytes[T]) IsEmpty() bool              { return len(b.data) == 0 }
func (b Bytes[T]) Payload() T                 { return b.payload }
func (b Bytes[T]) Value() T                   { return b.payload }
func (b Bytes[T]) AppError() *appfault.AppError { return b.appError }
func (b Bytes[T]) Fault() *appfault.AppError    { return b.appError }
func (b Bytes[T]) HasError() bool             { return b.appError != nil }
func (b Bytes[T]) IsValid() bool              { return b.appError == nil }
func (b Bytes[T]) Unwrap() ([]byte, *appfault.AppError) {
	return b.data, b.appError
}
```

---

## 3. Standardized Signatures & Interface Contracts (`pkg/streamwriter/contracts.go`)

### 3.1 Function Signatures
```go
// StreamFunc defines the swappable function signature returning *appfault.AppError.
type StreamFunc[T any] func(ctx context.Context, payload T, dest io.Writer) *appfault.AppError

// WriteFunc defines the swappable function signature returning *appfault.AppError.
type WriteFunc[T any] func(ctx context.Context, payload T) *appfault.AppError

// FormatFunc replaces ([]byte, error) with Bytes[T].
type FormatFunc[T any] func(payload T) Bytes[T]
```

### 3.2 Core Interface Hierarchy
```go
type WriterInterface[T any] interface {
	Interfacer
	Name() string
	Write(ctx context.Context, payload T) *appfault.AppError
	AsWriter() WriterInterface[T]
	Sync() *appfault.AppError
	Close() *appfault.AppError
}

type StreamerInterface[T any] interface {
	Interfacer
	Name() string
	Stream(ctx context.Context, payload T) *appfault.AppError
	AsStreamer() StreamerInterface[T]
	AsWriter() WriterInterface[T]
	IsLocked() bool
	Destination() io.Writer
	Sync() *appfault.AppError
	Close() *appfault.AppError
}
```

---

## 4. End-to-End Implementation Flow

### 4.1 Custom `FormatFunc[T]` Returning `Bytes[T]`
```go
// Custom JSON Formatter returning Bytes[MyEvent]
customFormatter := func(event MyEvent) streamwriter.Bytes[MyEvent] {
	if event.ID == "" {
		appErr := appfault.New(errtype.Validation, "event ID cannot be empty")
		return streamwriter.NewBytesErrorWithPayload(appErr, event)
	}

	compiled := []byte(fmt.Sprintf(`{"id":%q,"amount":%.2f}`, event.ID, event.Amount))
	return streamwriter.NewBytes(compiled, event)
}
```

### 4.2 Handling Monadic `*appfault.AppError`
```go
func main() {
	ctx := context.Background()

	writer := streamwriter.NewLockedStreamer[string](streamwriter.LockedOptions[string]{
		Name: "secure-streamer",
	})

	// All stream and write invocations return *appfault.AppError
	var appErr *appfault.AppError = writer.Stream(ctx, "data event")
	if appErr != nil {
		fmt.Printf("Stream failed: %s (Code: %d)\n", appErr.Message(), appErr.Type().Code())
		return
	}

	logger := streamwriter.NewLogger[any]().AddStreamer(writer)
	appErr = logger.Info(ctx, "Application booted successfully")
	if appErr != nil {
		log.Fatalf("Fatal logger error: %s", appErr.Message())
	}
}
```

---

## 5. Verification Results

All 11 tests pass with 100% green verification (`go test ./pkg/streamwriter -v -count=1`):
```
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
ok      coding-guidelines/common/pkg/streamwriter       0.502s
```
