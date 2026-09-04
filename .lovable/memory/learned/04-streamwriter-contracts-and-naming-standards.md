# Learned: StreamWriter Contracts, Idiomatic Naming, and Object Self-Context

> **Path:** `.lovable/memory/learned/04-streamwriter-contracts-and-naming-standards.md`  
> **Topic:** Go StreamWriter contracts, idiomatic `-er` interfaces, reentrant locking, monadic Bytes[T], JsonResult multi-source architecture, and boolean/Id naming standards  
> **Updated:** 2026-09-05  

---

## 1. Verbatim User Directives

```text
not type WriterInterface[T any] interface {
	Interfacer
	Name() string
	Write(ctx context.Context, payload T) *appfault.AppError
	AsWriter() WriterInterface[T]
	Sync() *appfault.AppError
	Close() *appfault.AppError
}

but Writer `er` needs to be at the end you stupid
```

```text
writer needs to have Lock, Unlock methods, avoid interfacer what is the value or benefit of it?
```

```text
creare as WrappedBytes have status flag as well, all the types must have 

type Bytes[T any] struct {
	data     []byte
	payload  T
	appError *appfault.AppError
}
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

Value() T 
Error() appError
```

```text
instead of JSON write evreywehre as Json and json source I think it take any instead othe T what do you think??
make json result without T can you fix it
```

```text
type JsonResult struct {
    data       []byte
    payload    any
    status     bool
    statusCode int
    appError   *appfault.AppError
}

should be isValid but that will be calculated from the error recieved or not so please , i don't think we need payload but we can have another version which extends JsonResult

so no need for any status , status code, payload, clear?? Fix all
```

```text
bool needs to have prefix and Id needs to be Id not ID, fix those

 auditWriter := streamwriter.NewPluggableWriter[any](<streamwriter.WriterOptions[any]{
        Name: "audit-api-writer",
        WriteMethod: func(ctx context.Context, payload any>) *appfault.AppError {
            trace := ""
            if traceVal := ctx.Value("traceId"); traceVal != nil {
                trace = fmt.Sprintf("[%v] ", traceVal)
            }
            _, err := fmt.Fprintf(dest, "[AUDIT] %s%s\n", trace, streamwriter.Compile(payload))
            if err != nil {
                return appfault.Wrap(errtype.IO, err, "audit write failed")
            }
            return nil
        },
    })

this func hsould take the current object in the func to proceed with so that we can use other properties from that section ,clear>>>>
```

---

## 2. Core Architectural Standards Established

### 2.1 Idiomatic `-er` Interface Naming
- **RULE:** Never suffix Go interfaces with `Interface` (e.g., `WriterInterface` is strictly banned).
- Interfaces representing behaviors must use the standard Go `-er` suffix: `Writer[T]`, `Streamer[T]`, `Compiler`, `Formatter[T]`.

### 2.2 Reentrant Locker Integration (`ReentrantMutex`)
- `Writer[T]` embeds `sync.Locker` (`Lock()` and `Unlock()`).
- In complex write flows where a compound batch acquires a lock and calls helper methods that also lock, a standard `sync.Mutex` deadlocks.
- The `streamwriter` package uses a goroutine-aware `ReentrantMutex` that tracks recursion depth and current goroutine ID, allowing safe re-entrant locking within the same goroutine while serializing access across concurrent goroutines.

### 2.3 Monadic `Bytes[T]` and `WrappedBytes`
- Replaces raw `([]byte, error)` returns across low-level streaming components.
- Provides fluent accessors: `Raw()`, `String()`, `Len()`, `IsEmpty()`, `Payload()`, `AppError()`, `Fault()`, `HasError()`, `IsValid()`, `Unwrap()`, `Value()`, `Error()`.

### 2.4 Minimalist `JsonResult` Container
- Pure bytes payload:
  ```go
  type JsonResult struct {
      data     []byte
      appError *appfault.AppError
  }
  ```
- No redundant `status`, `statusCode`, or `payload` fields inside `JsonResult`.
- `IsValid()` is derived dynamically: `func (r JsonResult) IsValid() bool { return r.appError == nil }`.
- When typed payloads are required, `JsonPayloadResult[T]` extends `JsonResult` by embedding it alongside `payload T`.

### 2.5 Strict Boolean Prefixing
- All boolean struct fields and variable identifiers MUST have an explicit positive prefix:
  - `isActive` (never `active`).
  - `isSuccess` (never `success`).
  - `hasError`, `isValid`, `shouldRetry`.

### 2.6 Strict `Id` / `id` Naming (BAN on `ID`)
- Acronym `ID` is completely prohibited in struct fields and variables:
  - Use `UserId`, `OrderId`, `AccountId`, `TraceId` in PascalCase.
  - Use `userId`, `orderId`, `accountId`, `traceId` in camelCase.
  - Total ban on `UserID`, `OrderID`, `AccountID`, `TraceID`.

### 2.7 Writer Self-Context Passing in `WriteFunc`
- Closures passed to `NewPluggableWriter` must not be forced to capture external variables or hardcode names:
  ```go
  type WriteFunc[T any] func(ctx context.Context, writer *PluggableWriter[T], payload T) *appfault.AppError
  ```
- Inside `WriteMethod`, callers have direct access to:
  - `writer.Name()`: The configured writer name.
  - `writer.Destination()`: The configured destination `io.Writer`.
  - `writer.Streamer()`: Attached streamer instance.
  - `writer.Lock()` / `writer.Unlock()`: Thread-safe locking.
