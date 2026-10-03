# AppError Package Reference — Result[T], ResultSlice[T], ResultMap[K,V]

> **Parent:** [AppError Package Reference](./readme.md)
> **Version:** 1.5.0
> **Updated:** 2026-09-13

---

## 3. Result[T] — Single Value Wrapper

For service methods that return one item or nothing.

### 3.1 Struct

```go
type Result[T any] struct {
    value     T
    err       *AppError
    isDefined bool
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
| `HasError()` / `IsFailure()` | `bool` | True if operation failed |
| `IsSuccess()` / `IsSafe()` | `bool` | True if operation succeeded with no error |
| `IsDefined()` | `bool` | True if operation succeeded (no error) AND data T is not null/empty (recordCount > 0) |
| `HasRecord()` / `HasRecords()` | `bool` | True if operation succeeded AND contains more than 0 records |
| `IsEmpty()` | `bool` | True if no active error and payload has 0 records, or error is empty |
| `Count()` | `int` | Number of records in payload (0 if failed or empty) |
| `IsCountOtherThan(n)` | `bool` | True if operation failed OR record count != n |
| `Value()` / `Data()` | `T` | Returns payload value |
| `ValueOr(fallback)` | `T` | Returns value if defined, else fallback |
| `AppError()` / `Fault()` | `*AppError` | Returns structured AppError metadata, or nil |
| `Unwrap()` | `(T, *AppError)` | Bridges to standard tuple unpacking |

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
| Query | `HasError()` / `IsFailure()` | `bool` | True if operation failed |
| Query | `IsSuccess()` / `IsSafe()` | `bool` | True if no error (items may be empty) |
| Query | `IsDefined()` | `bool` | True if operation succeeded (no error) AND contains more than 0 items (recordCount > 0) |
| Query | `HasRecord()` / `HasRecords()` / `HasItems()` | `bool` | True if operation succeeded AND contains more than 0 items |
| Query | `IsEmpty()` | `bool` | True if zero items or empty error state |
| Query | `Count()` | `int` | Number of items (0 if failed) |
| Query | `IsCountOtherThan(n)` | `bool` | True if operation failed OR item count != n |
| Access | `Items()` / `Data` | `[]T` | Returns the slice (nil if error) |
| Access | `First()` | `Result[T]` | Result for first item; empty if none |
| Access | `Last()` | `Result[T]` | Result for last item; empty if none |
| Access | `GetAt(index)` | `Result[T]` | Result at index; empty if out of bounds |
| Access | `AppError()` / `Fault()` | `*AppError` | Returns structured AppError metadata, or nil |
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
| Query | `HasError()` / `IsFailure()` | `bool` | True if operation failed |
| Query | `IsSuccess()` / `IsSafe()` | `bool` | True if no error (map may be empty) |
| Query | `IsDefined()` | `bool` | True if operation succeeded (no error) AND contains more than 0 entries (recordCount > 0) |
| Query | `HasRecord()` / `HasRecords()` / `HasItems()` | `bool` | True if operation succeeded AND contains more than 0 entries |
| Query | `IsEmpty()` | `bool` | True if zero entries or empty error state |
| Query | `Count()` | `int` | Number of entries (0 if failed) |
| Query | `IsCountOtherThan(n)` | `bool` | True if operation failed OR entry count != n |
| Query | `Has(key)` | `bool` | True if key exists |
| Access | `Items()` / `Data` | `map[K]V` | Returns the map (nil if error) |
| Access | `Get(key)` | `(V, bool)` | Safely retrieves entry by key |
| Access | `Keys()` | `[]K` | All keys as deterministically sorted slice |
| Access | `Values()` | `[]V` | All values as slice |
| Access | `AppError()` / `Fault()` | `*AppError` | Returns structured AppError metadata, or nil |
| Mutate | `Set(key, value)` | — | Adds/updates entry; no-op if error state |
| Mutate | `Remove(key)` | — | Deletes key; no-op if error state |

> **📌 `.AppError()` Naming Convention:**
> All result wrappers — `Result[T]`, `ResultSlice[T]`, and `ResultMap[K, V]` — expose the underlying error via `.AppError()` (returning `*AppError`), **not** `.Error()`. This avoids collision with Go's native `error` interface method `.Error() string` and ensures callers always receive the structured `*AppError` type for direct propagation via `Fail[T]()`, `FailSlice[T]()`, etc. without interface casts. The same convention applies to `dbutil` result types (`dbutil.Result[T]`, `dbutil.ResultSet[T]`, `dbutil.ExecResult`), which also store and return `*appfault.AppError` from their `.AppError()` method to enable bridge methods like `ToAppResult()` and `ToAppResultSlice()`.

---

## 6. Pointer-Attached Null Safety (*Result[T], *ResultSlice[T], *ResultMap[K, V])

### 6.1 Rationale: Eliminating Nil-Receiver Runtime Panics

In Go, declaring inspection methods on a value receiver (`func (r Result[T]) Method()`) panics immediately (`panic: runtime error: invalid memory address or nil pointer dereference`) when called on a `nil` pointer (`(*Result[T])(nil).IsFailure()`), because Go attempts to dereference the pointer to copy the struct by value before entering the method body.

To guarantee zero runtime panics across service layers and pipelines, all inspection, query, and data accessor methods on `Result[T]`, `ResultSlice[T]`, and `ResultMap[K, V]` are declared on **pointer receivers**:
- `(r *Result[T])`
- `(rs *ResultSlice[T])`
- `(rm *ResultMap[K, V])`

Every method verifies `if r == nil` on line 1 before checking internal fields or data.

### 6.2 Canonical Defaults on Nil Receiver

When invoked on a `nil` pointer, methods return safe, predictable defaults:

| Method | Return on `nil` Receiver | Behavior & Rationale |
|---|---|---|
| `IsFailure()` / `IsFailed()` / `HasError()` | `true` | Nil/uninitialized result represents a failed or missing operation. |
| `IsSuccess()` / `IsSafe()` | `false` | A nil pointer cannot represent success. |
| `IsEmptyError()` / `HasNoError()` | `false` | A nil pointer is not error-free. |
| `Count()` | `0` | Zero records in a nil container. |
| `IsEmpty()` | `true` | A nil container has zero elements. |
| `HasRecord()` / `HasRecords()` | `false` | Cannot contain records if pointer is nil. |
| `IsDefined()` | `false` | Cannot be defined if container is nil. |
| `IsCountOtherThan(n)` | `true` | A nil/failed result differs from any expected record count. |
| `AppError()` / `Fault()` | `nil` | Returns nil error safely without panicking. |
| `Value()` / `Data()` | `zero value of T` | Safe zero fallback for type `T`. |
| `Items()` / `Data` | `nil` | Nil slice / map fallback. |
| `Get(key)` | `zero, false` | Reports key not found safely. |
| `Has(key)` | `false` | Nil map contains no keys. |

---

### 6.3 Method Composition & Affirmative Field Naming

1. **Affirmative Boolean Struct Fields (`isDefined bool`):**
   - All boolean fields in result wrappers MUST use affirmative prefixes (e.g. `isDefined bool`, NEVER bare `defined bool`).
   - Bare boolean field names violate repository-wide boolean principles.

2. **Method Composition & Reuse (Methods Must Delegate to Existing Predicates):**
   - Inspection and predicate methods MUST NOT duplicate raw pointer or error-checking logic.
   - Higher-level predicates MUST compose existing methods (`IsFailure()`, `IsSuccess()`, `Count()`, `HasRecord()`):
     - `IsFailure()` delegates to `r.IsFailed()`
     - `IsValid()` delegates to `r.IsSuccess()`
     - `Count()` guards with `if r.IsFailure() { return 0 }`
     - `IsCountOtherThan(n)` guards with `if r.IsFailure() { return true } return r.Count() != n`
     - `IsEmpty()` guards with `if r.IsFailure() { return true } return !r.isDefined || isValueEmpty(r.value)`
     - `HasRecord()` guards with `if r.IsFailure() { return false } return r.Count() > 0`
     - `HasRecords()` delegates to `r.HasRecord()`
     - `IsDefined()` guards with `if r.IsFailure() { return false } return r.isDefined && !isValueEmpty(r.value)` (or `r.Count() > 0` for collections)

---

### 6.4 Mandatory `types.go` Definition as a Single Reusable Type

1. **Dedicated `types.go` per Package:**
   - All domain payload structs (e.g. `ScheduleExportBundle`, `Config`, `User`) and repeated generic Result type aliases (e.g. `type ConfigResult = result.Wrap[*Config]`, `type ScheduleExportBundleResult = result.ResultSlice[ScheduleExportBundle]`) MUST be defined in a dedicated `types.go` file within each package as a single reusable named type.
   - If a type is reused and used across functions, services, or packages, convert it from raw generic (`result.Wrap[*Config]`) to an actual concrete named type (`ConfigResult`).
   - Never declare unexported structs or raw generic Result envelopes inline in implementation files.

2. **Concrete Type Follow-Through Example:**
   ```go
   // -----------------------------------------------------------------------------
   // Step 1: Declare Concrete Types in `types.go` (Mandatory Rule)
   // -----------------------------------------------------------------------------
   // In types.go:
   // package config
   //
   // type (
   //     // Config contains application configuration fields.
   //     Config struct {
   //         Port int    `json:"port"`
   //         Host string `json:"host"`
   //     }
   //
   //     // ConfigResult is the single reusable concrete result envelope for *Config.
   //     // RULE: Convert raw generic result.Wrap[*Config] to an explicit concrete type
   //     // in types.go so all signatures and callers share the exact same definition!
   //     ConfigResult = result.Wrap[*Config]
   // )
   // -----------------------------------------------------------------------------

   // In config.go (Implementation File):
   // Signature uses the concrete type from types.go directly:
   func LoadConfig(path string) ConfigResult {
       if path == "" {
           fault := appfault.New(errtype.Validation, "config path cannot be empty").
               WithOp("config.LoadConfig")

           return result.WrapFailure[*Config](fault)
       }

       data, err := os.ReadFile(path)

       if err != nil {
           fault := appfault.WrapFile(errtype.IO, err, path, "failed to read config file").
               WithOp("config.LoadConfig")

           return result.WrapFailure[*Config](fault)
       }

       var cfg Config
       err = json.Unmarshal(data, &cfg)

       if err != nil {
           fault := appfault.Wrap(errtype.Validation, err, "failed to parse config json").
               WithOp("config.LoadConfig")

           return result.WrapFailure[*Config](fault)
       }

       return result.WrapSuccess(&cfg)
   }
   ```

3. **Library Package Architecture (`pkg/result/`):**
   - Core types (`Wrap[T]`, `Result[T]`, `ResultSlice[T]`, `ResultMap[K, V]`, verifier and inspector interfaces) are defined in `types.go` as single canonical types. Implementation files contain only functions, constructors, and methods.

---

