# Swappable Writer Methods & Functional Injection Architecture

> **Status:** Proposal & Architectural Exploration  
> **Date:** 2026-09-03  
> **Target Package:** `04-code/golang/pkg/applogger` and `pkg/writer`  
> **Topic:** 4 Distinct Patterns for Swappable Write Methods, Functional Injection via Options, and Log-Agnostic Payloads  

---

## 1. Problem Statement & Architecture Goals

In conventional logging frameworks, a writer's internal execution logic is fixated: the `Write` method is baked into the struct, forcing consumers to subclass or implement an entirely new struct just to tweak formatting, add prefixes, or handle non-log payloads.

### Key Goals:
1. **Swappable Write Method:** The `Write` behavior itself should be swappable at runtime or injected via `Options`.
2. **Log-Agnostic (Dual Mode):** The writer must support both structured log events (`LogRecord`) and arbitrary payloads (raw bytes, JSON objects, domain events, metrics).
3. **Composable via Options:** Consumers can customize behavior cleanly via functional options:
   ```go
   w := writer.New(writer.Options{
       WriteMethod: customWriteFunc,
   })
   ```
4. **Dynamic Reconfigurability:** The user can swap the write method or formatter on an active writer instance at any time without rebuilding.

---

## 2. Pattern 1: Functional Delegate Injection (First-Class `WriteFunc`)

### Concept
The writer contains a `writeMethod` function field. The public `Write(ctx, payload)` method delegates directly to this function. Default implementations provide standard text/JSON outputs, but users can inject any custom function via Options or setters.

```
       User calls writer.Write(ctx, payload)
                        │
                        ▼
       ┌──────────────────────────────────┐
       │           BaseWriter             │
       │   w.mu.RLock()                   │
       │   fn := w.writeMethod            │
       │   return fn(ctx, payload)        │
       └────────────────┬─────────────────┘
                        │ Delegates to
            ┌───────────┴───────────┐
            ▼                       ▼
   Default Text/JSON         Custom User Function
   (Standard Console/File)   (e.g., REST API / Slack / Kafka)
```

### Complete Implementation
```go
package writer

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
)

// WriteFunc defines the signature for any pluggable write behavior.
// Accepts arbitrary payload (log-based or without log-based).
type WriteFunc func(ctx context.Context, payload any) error

type Options struct {
	Name        string
	Destination io.Writer
	WriteMethod WriteFunc // Injected custom write function
}

type FunctionalWriter struct {
	mu          sync.RWMutex
	name        string
	destination io.Writer
	writeMethod WriteFunc
}

func NewFunctionalWriter(opts Options) *FunctionalWriter {
	if opts.Destination == nil {
		opts.Destination = os.Stdout
	}

	w := &FunctionalWriter{
		name:        opts.Name,
		destination: opts.Destination,
	}

	// Use user-provided method, or fall back to default text output
	if opts.WriteMethod != nil {
		w.writeMethod = opts.WriteMethod
	} else {
		w.writeMethod = w.defaultTextWrite
	}

	return w
}

// Write delegates to the active writeMethod function.
func (w *FunctionalWriter) Write(ctx context.Context, payload any) error {
	w.mu.RLock()
	fn := w.writeMethod
	w.mu.RUnlock()

	return fn(ctx, payload)
}

// SetWriteMethod allows swapping the write method at runtime.
func (w *FunctionalWriter) SetWriteMethod(fn WriteFunc) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.writeMethod = fn
}

// Default fallback implementation
func (w *FunctionalWriter) defaultTextWrite(ctx context.Context, payload any) error {
	line := fmt.Sprintf("[%s] %v\n", w.name, payload)
	_, err := w.destination.Write([]byte(line))
	return err
}
```

### Usage Example
```go
// 1. Using default write method
w := writer.NewFunctionalWriter(writer.Options{Name: "console"})
w.Write(ctx, "Hello world") // Outputs: [console] Hello world

// 2. Injecting custom write method via Options
customWriter := writer.NewFunctionalWriter(writer.Options{
	Name: "webhook-writer",
	WriteMethod: func(ctx context.Context, payload any) error {
		// Custom HTTP POST or non-log payload formatting
		fmt.Printf("POSTing payload to external system: %v\n", payload)
		return nil
	},
})

// 3. Swapping method dynamically on the fly
w.SetWriteMethod(func(ctx context.Context, payload any) error {
	fmt.Printf("SWAPPED: %v\n", payload)
	return nil
})
```

---

## 3. Pattern 2: Composable Pipeline (Decoupled `FormatFunc` + `EmitFunc`)

### Concept
Deconstructs writing into two distinct functions:
1. **Formatter Function:** `type FormatFunc func(payload any) ([]byte, error)`
2. **Emitter Function:** `type EmitFunc func(ctx context.Context, data []byte) error`

This allows independent swapping: you can change *how* it formats without changing *where* it writes, or change *where* it writes without changing *how* it formats.

```go
package writer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
)

type FormatFunc func(payload any) ([]byte, error)
type EmitFunc func(ctx context.Context, data []byte) error

type PipelineOptions struct {
	Name       string
	Formatter  FormatFunc
	Emitter    EmitFunc
}

type PipelineWriter struct {
	mu        sync.RWMutex
	name      string
	formatter FormatFunc
	emitter   EmitFunc
}

func NewPipelineWriter(opts PipelineOptions) *PipelineWriter {
	w := &PipelineWriter{name: opts.Name}

	if opts.Formatter != nil {
		w.formatter = opts.Formatter
	} else {
		w.formatter = func(payload any) ([]byte, error) {
			return []byte(fmt.Sprintf("%v\n", payload)), nil
		}
	}

	if opts.Emitter != nil {
		w.emitter = opts.Emitter
	} else {
		w.emitter = func(ctx context.Context, data []byte) error {
			_, err := os.Stdout.Write(data)
			return err
		}
	}

	return w
}

func (p *PipelineWriter) Write(ctx context.Context, payload any) error {
	p.mu.RLock()
	format := p.formatter
	emit := p.emitter
	p.mu.RUnlock()

	bytes, err := format(payload)
	if err != nil {
		return err
	}

	return emit(ctx, bytes)
}

func (p *PipelineWriter) SetFormatter(fn FormatFunc) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.formatter = fn
}

func (p *PipelineWriter) SetEmitter(fn EmitFunc) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.emitter = fn
}
```

### Usage Example
```go
// Default text pipeline
pw := writer.NewPipelineWriter(writer.PipelineOptions{Name: "pipeline"})

// Swap ONLY the formatter to JSON on the fly
pw.SetFormatter(func(payload any) ([]byte, error) {
	b, err := json.Marshal(map[string]any{"data": payload, "ts": "now"})
	return append(b, '\n'), err
})

// Swap ONLY the emitter to a file or REST endpoint
pw.SetEmitter(func(ctx context.Context, data []byte) error {
	return myFile.Write(data)
})
```

---

## 4. Pattern 3: Polymorphic Contract with Method Slot Overrides (AUK Go Style)

### Concept
Inspired by AUK Go's `LogDefinerWriter` and `AllLogWriter`. The writer provides multiple specialized entry points (`Write`, `WriteLog`, `WriteRaw`), but each method points to an internal function slot that users can override individually.

```go
package writer

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
)

// Universal Log Definer interface
type LogEntity interface {
	Message() string
	Fields() map[string]any
}

type SlotOptions struct {
	Destination io.Writer
	WriteLogFn  func(ctx context.Context, entity LogEntity) error
	WriteRawFn  func(ctx context.Context, raw []byte) error
	WriteAnyFn  func(ctx context.Context, item any) error
}

type SlotWriter struct {
	mu          sync.RWMutex
	dest        io.Writer
	writeLogFn  func(ctx context.Context, entity LogEntity) error
	writeRawFn  func(ctx context.Context, raw []byte) error
	writeAnyFn  func(ctx context.Context, item any) error
}

func NewSlotWriter(opts SlotOptions) *SlotWriter {
	if opts.Destination == nil {
		opts.Destination = os.Stdout
	}

	sw := &SlotWriter{dest: opts.Destination}

	// Slot 1: Structured Log Writer
	if opts.WriteLogFn != nil {
		sw.writeLogFn = opts.WriteLogFn
	} else {
		sw.writeLogFn = sw.defaultWriteLog
	}

	// Slot 2: Raw Byte Writer (without log)
	if opts.WriteRawFn != nil {
		sw.writeRawFn = opts.WriteRawFn
	} else {
		sw.writeRawFn = sw.defaultWriteRaw
	}

	// Slot 3: Generic Object Writer
	if opts.WriteAnyFn != nil {
		sw.writeAnyFn = opts.WriteAnyFn
	} else {
		sw.writeAnyFn = sw.defaultWriteAny
	}

	return sw
}

func (s *SlotWriter) WriteLog(ctx context.Context, entity LogEntity) error {
	s.mu.RLock()
	fn := s.writeLogFn
	s.mu.RUnlock()
	return fn(ctx, entity)
}

func (s *SlotWriter) WriteRaw(ctx context.Context, raw []byte) error {
	s.mu.RLock()
	fn := s.writeRawFn
	s.mu.RUnlock()
	return fn(ctx, raw)
}

func (s *SlotWriter) Write(ctx context.Context, item any) error {
	s.mu.RLock()
	fn := s.writeAnyFn
	s.mu.RUnlock()
	return fn(ctx, item)
}

// Default slot implementations
func (s *SlotWriter) defaultWriteLog(ctx context.Context, entity LogEntity) error {
	_, err := fmt.Fprintf(s.dest, "[LOG] %s %v\n", entity.Message(), entity.Fields())
	return err
}

func (s *SlotWriter) defaultWriteRaw(ctx context.Context, raw []byte) error {
	_, err := s.dest.Write(raw)
	return err
}

func (s *SlotWriter) defaultWriteAny(ctx context.Context, item any) error {
	_, err := fmt.Fprintf(s.dest, "%v\n", item)
	return err
}
```

---

## 5. Pattern 4: Generic Dual-Channel Streamer (`Writer[T]`)

### Concept
Uses Go 1.18+ Generics to provide 100% type safety with zero interface boxing/unboxing overhead. Can be instantiated for `LogRecord` (log-based), `[]byte` (raw), or `any` (universal).

```go
package writer

import (
	"context"
	"sync"
)

type GenericWriteFunc[T any] func(ctx context.Context, item T) error

type Streamer[T any] struct {
	mu      sync.RWMutex
	name    string
	writeFn GenericWriteFunc[T]
}

func NewStreamer[T any](name string, writeFn GenericWriteFunc[T]) *Streamer[T] {
	return &Streamer[T]{
		name:    name,
		writeFn: writeFn,
	}
}

func (s *Streamer[T]) Write(ctx context.Context, item T) error {
	s.mu.RLock()
	fn := s.writeFn
	s.mu.RUnlock()
	return fn(ctx, item)
}

func (s *Streamer[T]) SetWriteMethod(fn GenericWriteFunc[T]) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.writeFn = fn
}
```

### Usage Example
```go
// Type-safe for structured LogRecord:
logStreamer := writer.NewStreamer(
	"audit",
	func(ctx context.Context, r LogRecord) error {
		fmt.Printf("Audit: %s\n", r.Message)
		return nil
	},
)
logStreamer.Write(ctx, LogRecord{Message: "Transaction complete"})

// Type-safe for raw byte streams (non-log):
byteStreamer := writer.NewStreamer(
	"raw",
	func(ctx context.Context, b []byte) error {
		return os.Stdout.Write(b)
	},
)
byteStreamer.Write(ctx, []byte("raw network frame"))
```

---

## 6. Comparison Matrix

| Criteria | Pattern 1: Functional Injection (`WriteFunc`) | Pattern 2: Composable Pipeline (`Format` + `Emit`) | Pattern 3: Slot Overrides (AUK Go Style) | Pattern 4: Generic Streamer (`Writer[T]`) |
| :--- | :--- | :--- | :--- | :--- |
| **Write Method Swapping** | Complete replacement via single function | Independent formatting & emission swapping | Individual method slot overrides | Swappable type-safe function |
| **Log-Agnostic Support** | Yes (accepts `any`) | Yes (accepts `any`) | Dedicated slots for `Log`, `Raw`, `Any` | Yes (via type parameter `T`) |
| **Options Configuration** | `Options.WriteMethod` | `Options.Formatter`, `Options.Emitter` | `Options.WriteLogFn`, `Options.WriteRawFn` | Initial constructor function |
| **Simplicity & Readability** | Highest (cleanest signature) | High (modular separation) | Medium (multiple method slots) | High (compile-time checked) |
| **Runtime Overhead** | 1 pointer dereference | 2 function calls | 1 pointer dereference | 0 interface boxing |
| **AUK Go Alignment** | Modernized functional evolution | Layered pipeline | Direct mirror of AUK Go design | Idiomatic modern Go 1.21+ |

---

## 7. Recommended Direction

**Pattern 1 (Functional Delegate Injection)** combined with **Pattern 2's Decoupled Formatter/Emitter**:
- Gives the caller complete freedom to pass `Options.WriteMethod` to do whatever they want.
- If not provided, it defaults to `format(payload) -> emit(bytes)`.
- The user can swap either the whole `WriteMethod`, or just the `Formatter`, or just the `Destination`.
