# Specification: `appfaults` Error Collection Architecture

**Version:** 3.3.1  
**Status:** Complete & Enforced  
**Package:** `04-code/golang/pkg/appfaults`  
**Reference Implementations:** `https://gitlab.com/auk-go/errorwrapper` (`errwrappers.Collection`, `errwrappers.MutexCollection`)

---

## 1. Executive Summary & Purpose

In multi-step workflows, batch validations, and concurrent pipelines, functions frequently encounter multiple errors. The **`appfaults`** package provides an error collection type (`appfaults.Collection` / `appfaults.AppFaults`) for aggregating, filtering, transforming, and propagating multiple `*appfault.AppError` instances across application layers and Go context (`context.Context`).

---

## 2. Core Rule: Mandatory Error Type First

1. **Explicit Error Type First:**
   - Every error construction or wrapping method MUST accept `errtype.Variation` as its first parameter:
     - `appfault.New(errType, msg)`
     - `appfault.NewType(errType)`
     - `appfault.Wrap(errType, cause, msg)`
     - `appfault.WrapType(errType, cause)`
2. **Missing Cause & Nil Safety:**
   - If `cause == nil` or `errType == errtype.None`, constructors return `nil` without allocating memory or stack traces.
   - Collections ignore `nil` and `errtype.None` additions, preventing dummy errors.

---

## 3. Collection Mutators & Methods

### 3.1 `Collection` & `MutexCollection` Mutators

| Method | Signature | Description |
| :--- | :--- | :--- |
| `Add(err)` | `func (c *Collection) Add(err *appfault.AppError) *Collection` | Appends error if non-nil and `HasError()` is true. |
| `AddType(errType)` | `func (c *Collection) AddType(errType errtype.Variation) *Collection` | Creates error from type and appends. |
| `AddTypeMsg(errType, msg)` | `func (c *Collection) AddTypeMsg(errType errtype.Variation, msg string) *Collection` | Creates error from type + message. |
| `AddTypeMsgf(errType, fmt, args...)` | `func (c *Collection) AddTypeMsgf(errType errtype.Variation, format string, args ...any) *Collection` | Formatted message error. |
| `AddError(errType, cause)` | `func (c *Collection) AddError(errType errtype.Variation, cause error) *Collection` | Wraps cause with type (message is `cause.Error()`). |
| `AddErrorMsg(errType, cause, msg)` | `func (c *Collection) AddErrorMsg(errType errtype.Variation, cause error, msg string) *Collection` | Wraps cause with type + custom message. |
| `AddWithContext(errType, msg, ctx)` | `func (c *Collection) AddWithContext(errType errtype.Variation, msg string, ctx map[string]any) *Collection` | Creates error with context map. |
| `AddAll(faults...)` | `func (c *Collection) AddAll(faults ...*appfault.AppError) *Collection` | Appends multiple AppErrors. |
| `Merge(other)` | `func (c *Collection) Merge(other *Collection) *Collection` | Ingests items from another collection. |
| `Clear()` | `func (c *Collection) Clear() *Collection` | Resets the collection. |

### 3.2 Status & Inspection

| Method | Signature | Semantics / Behavior |
| :--- | :--- | :--- |
| `HasError()` | `func (c *Collection) HasError() bool` | Returns `true` if $\ge 1$ errors. Returns `false` if `c == nil` or empty. |
| `IsSuccess()` | `func (c *Collection) IsSuccess() bool` | Returns `true` if `c == nil` or `Count() == 0`. |
| `IsEmpty()` | `func (c *Collection) IsEmpty() bool` | Alias for `!c.HasError()`. |
| `Count()` | `func (c *Collection) Count() int` | Returns total number of active errors. |
| `Items()` | `func (c *Collection) Items() []*appfault.AppError` | Returns copy of backing slice. |
| `First()` | `func (c *Collection) First() *appfault.AppError` | Returns first error or `nil`. |
| `Last()` | `func (c *Collection) Last() *appfault.AppError` | Returns last error or `nil`. |

---

## 4. Context Integration (`context.Context`)

```go
// WithFaults creates a child context holding a dedicated AppFaults collector.
func WithFaults(ctx context.Context) (context.Context, *Collection)

// FromContext extracts the existing AppFaults collector from context.
func FromContext(ctx context.Context) (*Collection, bool)

// RecordContextError records an error into context collector if present.
func RecordContextError(ctx context.Context, err *appfault.AppError) bool
```

---

## 5. Serialization & Diagnostics

- **JSON Output:** `json.Marshal(coll)` serializes to a JSON array of `AppError` objects.
- **Formatted String:** `coll.Format()` returns numbered bullet points for diagnostics.
