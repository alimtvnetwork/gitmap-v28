# AppError Package Reference — Result[T], ResultSlice[T], ResultMap[K,V]

> **Parent:** [AppError Package Reference](./00-overview.md)  
> **Version:** 1.3.0  
> **Updated:** 2026-03-31

---

## 3. Result[T] — Single Value Wrapper

For service methods that return one item or nothing.

### 3.1 Struct

```go
type Result[T any] struct {
    value   T
    err     *AppError
    defined bool
}
```

### 3.2 Constructors

```go
// Ok creates a successful Result containing the given value.
func Ok[T any](value T) Result[T]

// Fail creates a failed Result from an AppError.
func Fail[T any](err *AppError) Result[T]

// FailWrap creates a failed Result by wrapping a raw error.
// Uses skip=3 to point stack trace at caller, not this wrapper.
func FailWrap[T any](cause error, code, message string) Result[T]

// FailNew creates a failed Result from a new error (no cause).
func FailNew[T any](code, message string) Result[T]
```

### 3.3 Methods

| Method | Returns | Description |
|--------|---------|-------------|
| `HasError()` | `bool` | True if operation failed |
| `IsSafe()` | `bool` | True if value exists AND no error |
| `IsDefined()` | `bool` | True if value was set (regardless of error) |
| `IsEmpty()` | `bool` | True if no value was set |
| `Value()` | `T` | Returns value; panics if `HasError()` |
| `ValueOr(fallback)` | `T` | Returns value if defined, else fallback |
| `Error()` | `*AppError` | Returns the AppError, or nil |
| `Unwrap()` | `(T, error)` | Bridges to standard `(T, error)` pattern |

---


---

## 4. ResultSlice[T] — Collection Wrapper

For service methods that return lists of items.

### 4.1 Struct

```go
type ResultSlice[T any] struct {
    items []T
    err   *AppError
}
```

### 4.2 Constructors

```go
func OkSlice[T any](items []T) ResultSlice[T]
func FailSlice[T any](err *AppError) ResultSlice[T]

// Uses skip=3 for correct stack trace attribution.
func FailSliceWrap[T any](cause error, code, message string) ResultSlice[T]
func FailSliceNew[T any](code, message string) ResultSlice[T]
```

### 4.3 Methods

| Category | Method | Returns | Description |
|----------|--------|---------|-------------|
| Query | `HasError()` | `bool` | True if operation failed |
| Query | `IsSafe()` | `bool` | True if no error (items may be empty) |
| Query | `HasItems()` | `bool` | True if at least one item |
| Query | `IsEmpty()` | `bool` | True if zero items |
| Query | `Count()` | `int` | Number of items |
| Access | `Items()` | `[]T` | Returns the slice (nil if error) |
| Access | `First()` | `Result[T]` | Result for first item; empty if none |
| Access | `Last()` | `Result[T]` | Result for last item; empty if none |
| Access | `GetAt(index)` | `Result[T]` | Result at index; empty if out of bounds |
| Access | `Error()` | `*AppError` | Returns the AppError, or nil |
| Mutate | `Append(items...)` | — | Adds items; no-op if in error state |

---


---

## 5. ResultMap[K, V] — Associative Map Wrapper

For service methods that return key-value data.

### 5.1 Struct

```go
type ResultMap[K comparable, V any] struct {
    items map[K]V
    err   *AppError
}
```

### 5.2 Constructors

```go
func OkMap[K comparable, V any](items map[K]V) ResultMap[K, V]
func FailMap[K comparable, V any](err *AppError) ResultMap[K, V]

// Uses skip=3 for correct stack trace attribution.
func FailMapWrap[K comparable, V any](cause error, code, message string) ResultMap[K, V]
func FailMapNew[K comparable, V any](code, message string) ResultMap[K, V]
```

### 5.3 Methods

| Category | Method | Returns | Description |
|----------|--------|---------|-------------|
| Query | `HasError()` | `bool` | True if operation failed |
| Query | `IsSafe()` | `bool` | True if no error (map may be empty) |
| Query | `HasItems()` | `bool` | True if at least one entry |
| Query | `IsEmpty()` | `bool` | True if zero entries |
| Query | `Count()` | `int` | Number of entries |
| Query | `Has(key)` | `bool` | True if key exists |
| Access | `Items()` | `map[K]V` | Returns the map (nil if error) |
| Access | `Get(key)` | `Result[V]` | Result for key; empty if not found |
| Access | `Keys()` | `[]K` | All keys as slice |
| Access | `Values()` | `[]V` | All values as slice |
| Access | `Error()` | `*AppError` | Returns the AppError, or nil |
| Mutate | `Set(key, value)` | — | Adds/updates entry; no-op if error state |
| Mutate | `Remove(key)` | — | Deletes key; no-op if error state |

> **📌 `.AppError()` Naming Convention:**
> All result wrappers — `Result[T]`, `ResultSlice[T]`, and `ResultMap[K, V]` — expose the underlying error via `.AppError()` (returning `*AppError`), **not** `.Error()`. This avoids collision with Go's native `error` interface method `.Error() string` and ensures callers always receive the structured `*AppError` type for direct propagation via `Fail[T]()`, `FailSlice[T]()`, etc. without interface casts. The same convention applies to `dbutil` result types (`dbutil.Result[T]`, `dbutil.ResultSet[T]`, `dbutil.ExecResult`), which also store and return `*apperror.AppError` from their `.AppError()` method to enable bridge methods like `ToAppResult()` and `ToAppResultSlice()`.

---
