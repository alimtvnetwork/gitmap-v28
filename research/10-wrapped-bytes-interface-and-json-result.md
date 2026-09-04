# WrappedBytes Interface & JSONResult Architecture

> **Document:** `research/10-wrapped-bytes-interface-and-json-result.md`  
> **Status:** Implemented & Verified  
> **Package:** `04-code/golang/pkg/streamwriter`  
> **Date:** 2026-09-04  

---

## 1. Overview & Objectives

In high-throughput logging, streaming, and API serialization pipelines, handling raw byte slices (`[]byte`) alongside structured JSON requires standardized, type-safe result containers.

This architectural enhancement accomplishes:
1. **`WrappedBytes[T any]` Interface:** Standardizes all byte-wrapping types (`Bytes[T]` and `JSONResult[T]`), guaranteeing consistent access to payload data, status flag, numeric status code, and `*appfault.AppError` state.
2. **Status Flag & Status Code:** Every wrapped type contains an explicit boolean status flag (`Status() bool`, `IsSuccess() bool`) and integer status code (`StatusCode() int`), enabling clean pipeline error-handling and HTTP/gRPC code mapping.
3. **`JSONResult[T any]` Container:** Provides specialized JSON handling with formatting (`Pretty()`, `Compact()`), unmarshaling (`Unmarshal(dest any)`), and direct compliance with the `WrappedBytes[T]` and `WrappedJSON[T]` interfaces.
4. **Convenience Accessors:** `Value() T` (alias to `Payload()`) and `Error() *appfault.AppError` (alias to `AppError()`).

---

## 2. Core Interfaces

### `WrappedBytes[T any]` (`bytes.go`)

```go
type WrappedBytes[T any] interface {
	Raw() []byte
	Bytes() []byte
	String() string
	Len() int
	IsEmpty() bool
	Payload() T
	Value() T
	AppError() *appfault.AppError
	Fault() *appfault.AppError
	Error() *appfault.AppError
	HasError() bool
	IsValid() bool
	IsSuccess() bool
	Status() bool
	StatusCode() int
	Unwrap() ([]byte, *appfault.AppError)
}
```

### `WrappedJSON[T any]` (`json_result.go`)

```go
type WrappedJSON[T any] interface {
	WrappedBytes[T]
	Pretty() string
	Compact() string
	Unmarshal(dest any) *appfault.AppError
}
```

---

## 3. Concrete Implementations

### `Bytes[T any]` (`bytes.go`)
Represents generic byte buffers from serializers or encoders:
```go
type Bytes[T any] struct {
	data       []byte
	payload    T
	status     bool
	statusCode int
	appError   *appfault.AppError
}

// Compile-time interface assertion
var _ WrappedBytes[any] = Bytes[any]{}
```

### `JSONResult[T any]` (`json_result.go`)
Specialized container for JSON payloads:
```go
type JSONResult[T any] struct {
	data       []byte
	payload    T
	status     bool
	statusCode int
	appError   *appfault.AppError
}

type JsonResult[T any] = JSONResult[T]

// Compile-time interface assertions
var _ WrappedBytes[any] = JSONResult[any]{}
var _ WrappedJSON[any] = JSONResult[any]{}
```

---

## 4. Usage Patterns

### JSON Serialization & Unmarshaling
```go
type Order struct {
    ID    string  `json:"id"`
    Total float64 `json:"total"`
}

// Automatic JSON marshaling with status flag set to true (200)
result := streamwriter.NewJSONResult(Order{ID: "ord-1", Total: 99.50})

// Can be passed anywhere expecting WrappedBytes or WrappedJSON
var wb streamwriter.WrappedBytes[Order] = result

if wb.IsSuccess() {
    fmt.Println("Raw JSON:", wb.String())
    fmt.Println("Formatted:", result.Pretty())
}

// Unmarshal back to typed object
var target Order
if err := result.Unmarshal(&target); err == nil {
    fmt.Println("Unmarshaled ID:", target.ID)
}
```

---

## 5. Verification Results

```bash
$ go test ./pkg/streamwriter -v -count=1
=== RUN   TestBytesWrapper
--- PASS: TestBytesWrapper (0.00s)
=== RUN   TestJSONResult_WrappedBytesFlow
--- PASS: TestJSONResult_WrappedBytesFlow (0.00s)
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
ok  	coding-guidelines/common/pkg/streamwriter	0.545s
```
