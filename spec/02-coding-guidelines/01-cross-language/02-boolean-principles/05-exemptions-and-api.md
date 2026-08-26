# Boolean Principles — Static factory exemption, Result wrapper API

> **Parent:** [Boolean Principles](./00-overview.md)  
> **Version:** 2.7.0  
> **Updated:** 2026-04-19

---

## Static Factory Constructor Exemption

Methods like `DbResult::empty()`, `DbResultSet::empty()`, and `ResultSlice::empty()` are **static factory constructors** — they create a new empty instance, not query boolean state. These are **exempt** from the `is`/`has` prefix requirement (P1).

Boolean query methods on the **same classes** — such as `isEmpty()`, `isDefined()`, `hasError()`, `isSafe()`, `hasItems()` — **do** follow P1 correctly and must retain their prefixes.

| Method | Type | P1 Applies? |
|--------|------|-------------|
| `DbResult::empty()` | Static factory constructor | ❌ Exempt |
| `DbResultSet::empty()` | Static factory constructor | ❌ Exempt |
| `$result->isEmpty()` | Boolean query | ✅ Yes |
| `$result->hasError()` | Boolean query | ✅ Yes |
| `$result->isDefined()` | Boolean query | ✅ Yes |
| `result.IsSafe()` | Boolean query | ✅ Yes |

---

## Linter Exempt-Name Audit (CODE-RED-002)

The reference linter (`linter-scripts/validate-guidelines.py` and its Go twin `validate-guidelines.go`) suppresses the **P1 prefix requirement** for a fixed allowlist of identifier names. Each entry was audited on **2026-04-19** against this spec.

**Source:** `linter-scripts/validate-guidelines.py:133`, `linter-scripts/validate-guidelines.go:104-110`

| Name | P1 Exempt? | Scope | Justification |
|------|-----------|-------|---------------|
| `ok` | ✅ Valid | Go (idiomatic) | Go comma-ok pattern: `v, ok := m[k]`. Forcing `isOk` would violate Go community convention and break linter consistency with `gofmt`/`golint`. Aligns with **P7 exemption** for comma-ok already documented above. |
| `done` | ⚠️ Go-only | Go (channel idiom) | Standard channel completion signal: `done := make(chan struct{})`. **Outside Go**, a state flag must use `isDone`. |
| `found` | ⚠️ Go-only | Go (lookup return) | Idiomatic lookup-result name: `v, found := lookup(k)`. **Outside Go**, prefer `isFound` / `hasMatch`. |
| `exists` | ⚠️ Go-only | Go (existence return) | Idiomatic existence-check return: `_, exists := registry[k]`. **Outside Go**, prefer `isExisting` / `hasEntry`. |
| `err` | ✅ Valid | All languages | Holds an error value, not a boolean. Listed only to suppress false positives from defensive init patterns like `err := false`. |
| `error` | ✅ Valid | All languages | Same rationale as `err`. Defensive entry for legacy code paths. |
| `true` | ✅ Valid | All languages | Reserved literal, not an identifier. Defensive entry. |
| `false` | ✅ Valid | All languages | Reserved literal, not an identifier. Defensive entry. |

### Audit Verdict

| Status | Names | Action |
|--------|-------|--------|
| ✅ **Confirmed valid (cross-language)** | `ok`, `err`, `error`, `true`, `false` | Keep as-is |
| ⚠️ **Valid only in Go** | `done`, `found`, `exists` | **Recommendation:** language-scope these in `check_boolean_naming` so TS/PHP code receives the P1 violation. Tracked as linter enhancement. |

### Recommended Linter Change (non-breaking)

```python

# linter-scripts/validate-guidelines.py

GO_ONLY_EXEMPT = {"done", "found", "exists"}
UNIVERSAL_EXEMPT = {"ok", "err", "error", "true", "false"}

def check_boolean_naming(lines, filepath, lang):
    exempt = UNIVERSAL_EXEMPT if lang != "go" else UNIVERSAL_EXEMPT | GO_ONLY_EXEMPT
    # ...
```

This split preserves the existing Go behavior while restoring P1 enforcement for TypeScript and PHP. **No current codebase violation** would be triggered by this change (verified: `validate-guidelines.py` audit on 2026-04-19 reported 0 CODE-RED-002 violations).

### Cross-References

- [P1 — Naming prefixes](./01-naming-prefixes.md)
- [P7 — No assignments in conditions (comma-ok exemption)](./03-parameters-and-conditions.md)
- [Linter source — `validate-guidelines.py`](../../../../linter-scripts/validate-guidelines.py)
- [Linter source — `validate-guidelines.go`](../../../../linter-scripts/validate-guidelines.go)

---


---

## Result Wrapper — Full Public API Reference

> **Cross-language invariant:** The `.AppError()` (Go) / `.error()` (PHP) method on every result wrapper returns the **framework's structured error type**, never a raw string or generic exception. In Go this is `*apperror.AppError` (carrying stack trace, error code, and contextual values) — named `.AppError()` (not `.Error()`) to avoid confusion with Go's native `error` interface. In PHP this is `Throwable` (typically a framework exception with trace). This guarantees that propagated errors always preserve diagnostic context — callers can safely pass `.AppError()` output to `Fail()`, `FailSlice()`, `FailMap()`, or log it with full traceability.

### Go — `apperror.Result[T]`

| Method | Returns | Description |
|--------|---------|-------------|
| `Ok[T](value)` | `Result[T]` | Static: successful result with value |
| `Fail[T](err)` | `Result[T]` | Static: failed result from `*AppError` |
| `FailWrap[T](cause, code, msg)` | `Result[T]` | Static: failed result wrapping raw error |
| `FailNew[T](code, msg)` | `Result[T]` | Static: failed result from new error |
| `HasError()` | `bool` | True when the operation failed |
| `IsSafe()` | `bool` | True when a value exists and no error |
| `IsDefined()` | `bool` | True when a value was set (regardless of error) |
| `IsEmpty()` | `bool` | True when no value was set (absent, not an error) |
| `Value()` | `T` | Returns value; **panics** if `HasError()` is true |
| `ValueOr(fallback)` | `T` | Returns value if defined, otherwise fallback |
| `AppError()` | `*AppError` | Returns underlying error, or nil. Named `AppError()` to avoid confusion with Go's `error` interface |
| `Unwrap()` | `(T, error)` | Bridges to standard Go `(T, error)` pattern |

### Go — `apperror.ResultSlice[T]`

| Method | Returns | Description |
|--------|---------|-------------|
| `OkSlice[T](items)` | `ResultSlice[T]` | Static: successful slice result |
| `FailSlice[T](err)` | `ResultSlice[T]` | Static: failed slice from `*AppError` |
| `FailSliceWrap[T](cause, code, msg)` | `ResultSlice[T]` | Static: failed slice wrapping raw error |
| `FailSliceNew[T](code, msg)` | `ResultSlice[T]` | Static: failed slice from new error |
| `HasError()` | `bool` | True when the operation failed |
| `IsSafe()` | `bool` | True when no error (items may be empty) |
| `HasItems()` | `bool` | True when slice has ≥1 item |
| `IsEmpty()` | `bool` | True when slice has zero items |
| `Count()` | `int` | Number of items |
| `Items()` | `[]T` | Returns underlying slice (nil if error) |
| `First()` | `Result[T]` | First item as `Result[T]`, or empty |
| `Last()` | `Result[T]` | Last item as `Result[T]`, or empty |
| `GetAt(index)` | `Result[T]` | Item at index as `Result[T]`, or empty |
| `Append(items...)` | — | Adds items; no-op if in error state |
| `AppError()` | `*AppError` | Returns underlying error, or nil |

### Go — `apperror.ResultMap[K, V]`

| Method | Returns | Description |
|--------|---------|-------------|
| `OkMap[K,V](items)` | `ResultMap[K,V]` | Static: successful map result |
| `FailMap[K,V](err)` | `ResultMap[K,V]` | Static: failed map from `*AppError` |
| `FailMapWrap[K,V](cause, code, msg)` | `ResultMap[K,V]` | Static: failed map wrapping raw error |
| `FailMapNew[K,V](code, msg)` | `ResultMap[K,V]` | Static: failed map from new error |
| `HasError()` | `bool` | True when the operation failed |
| `IsSafe()` | `bool` | True when no error (map may be empty) |
| `HasItems()` | `bool` | True when map has ≥1 entry |
| `IsEmpty()` | `bool` | True when map has zero entries |
| `Count()` | `int` | Number of entries |
| `Items()` | `map[K]V` | Returns underlying map (nil if error) |
| `Get(key)` | `Result[V]` | Value for key as `Result[V]`, or empty |
| `Has(key)` | `bool` | True if key exists |
| `Set(key, value)` | — | Adds/updates entry; no-op if error |
| `Remove(key)` | — | Deletes key; no-op if error |
| `Keys()` | `[]K` | All keys as slice |
| `Values()` | `[]V` | All values as slice |
| `AppError()` | `*AppError` | Returns underlying error, or nil |

### PHP — `DbResult<T>`

| Method | Returns | Description |
|--------|---------|-------------|
| `DbResult::of($value)` | `DbResult<T>` | Static: successful result with value |
| `DbResult::empty()` | `DbResult<T>` | Static: empty result (no row found) |
| `DbResult::error($e)` | `DbResult<T>` | Static: error result with stack trace |
| `isEmpty()` | `bool` | True when no row was found (not an error) |
| `isDefined()` | `bool` | True when a row was successfully mapped |
| `hasError()` | `bool` | True when the query failed |
| `isSafe()` | `bool` | True when value exists and no error |
| `value()` | `T\|null` | Returns mapped value (null if not defined) |
| `error()` | `Throwable\|null` | Returns underlying error, or null |
| `stackTrace()` | `string` | Captured stack trace if error occurred |

### PHP — `DbResultSet<T>`

| Method | Returns | Description |
|--------|---------|-------------|
| `DbResultSet::of($items)` | `DbResultSet<T>` | Static: successful result set |
| `DbResultSet::error($e)` | `DbResultSet<T>` | Static: error result with stack trace |
| `isEmpty()` | `bool` | True when zero items |
| `hasAny()` | `bool` | True when ≥1 item |
| `count()` | `int` | Number of items |
| `hasError()` | `bool` | True when the query failed |
| `isSafe()` | `bool` | True when no error (items may be empty) |
| `items()` | `array<T>` | Returns item array |
| `first()` | `DbResult<T>` | First item as `DbResult<T>`, or error/empty |
| `error()` | `Throwable\|null` | Returns underlying error, or null |
| `stackTrace()` | `string` | Captured stack trace if error occurred |

### PHP — `DbExecResult`

| Method | Returns | Description |
|--------|---------|-------------|
| `DbExecResult::of($rows, $id)` | `DbExecResult` | Static: successful exec result |
| `DbExecResult::error($e)` | `DbExecResult` | Static: error result with stack trace |
| `hasError()` | `bool` | True when the exec failed |
| `isSafe()` | `bool` | True when no error |
| `isEmpty()` | `bool` | True when zero rows affected |
| `affectedRows()` | `int` | Number of affected rows |
| `lastInsertId()` | `int` | Auto-increment ID from INSERT |
| `error()` | `Throwable\|null` | Returns underlying error, or null |
| `stackTrace()` | `string` | Captured stack trace if error occurred |

---
