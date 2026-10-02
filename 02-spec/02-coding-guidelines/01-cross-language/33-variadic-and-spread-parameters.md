# Variadic & Spread Parameters Specification

**Version:** 1.0.0
**Updated:** 2026-10-02
**Applies to:** All languages (Go, TypeScript, Rust)
**Status:** Active
**AI Confidence:** High
**Ambiguity:** None

---

## 1. Overview & Core Motivation

In polyglot engineering architectures across Go, TypeScript, and Rust, APIs that accept collections of identical elements are frequently declared using rigid slice or array types:

```go
// Rigid slice parameter: forces artificial wrappers at every call site
func DeleteUsers(ids []string) *appfault.AppError
```

```typescript
// Rigid array parameter: forces artificial array brackets for single-item invocations
function deleteUsers(ids: string[]): Promise<Result<void>>
```

This rigid parameter pattern introduces systemic friction across codebases:

1. **Call-Site Clutter & Friction:** Single-item invocations—which often constitute over 70% of call sites in practice—are forced to construct artificial wrapper structures (`DeleteUsers([]string{id})` or `deleteUsers([id])`).
2. **Artificial Memory Allocations:** In Go and TypeScript, constructing ad-hoc single-item slice or array literals creates temporary heap or slice-header allocations that escape unnecessarily when passed across boundaries.
3. **Impaired Fluent Composability:** Rigid slice parameters prevent callers from naturally passing comma-separated arguments or chaining variadic builders and options without intermediate list instantiations.
4. **Ergonomic Inversion:** Callers possessing an existing collection can easily pass it to a variadic function using spread syntax (`items...` in Go, `...items` in TypeScript). Conversely, callers possessing a single item cannot easily invoke a rigid slice function without wrapping it.

### The Variadic & Spread Solution

By defining functions with **variadic / spread parameters** (`...T` in Go, `...items: readonly T[]` or `SingleOrArray<T>` union overloads in TypeScript, and zero-cost slice references `&[T]` / `impl IntoIterator<Item = T>` in Rust), APIs achieve optimal ergonomics:

- **Single-Item Calls:** Directly passed without wrapper syntax: `DeleteUsers(id)` / `deleteUsers(id)`.
- **Multi-Item Calls:** Comma-separated without slice declaration: `DeleteUsers(id1, id2)` / `deleteUsers(id1, id2)`.
- **Slice Spread Calls:** Existing slices seamlessly unpacked: `DeleteUsers(allIDs...)` / `deleteUsers(...allIds)`.
- **Empty Invocations:** Zero parameters passed cleanly without `nil` or `[]` literals: `DeleteUsers()`.

---

## 2. Mandatory Rules

### Rule 1: Variadic Parameter Primacy for Homogenous Collections

- Functions whose primary or trailing input is a collection of identical type elements SHOULD declare that parameter as variadic (`...T` in Go, `...items: readonly T[]` in TypeScript).
- **Total Ban:** NEVER declare rigid slice parameters (`items []string`) for utility, deletion, filtering, registration, and dispatch APIs where callers frequently operate on single elements.
- In Rust, functions accepting homogenous collections MUST accept borrowed slices (`&[T]`) or generic iterators (`impl IntoIterator<Item = T>`), banning owned heap vectors (`Vec<T>`).

### Rule 2: Spread Operator Forwarding & Transparency

- Functions that receive a variadic parameter and pass it to another variadic function MUST forward using the native spread operator (`targetFunc(items...)` in Go, `targetFunc(...items)` in TypeScript).
- Callers possessing an existing slice, array, or vector pass it directly using spread without rebuilding or wrapping intermediate collections.

### Rule 3: Zero Artificial Wrapper Mandate at Call Sites

- Callers passing a single item MUST pass the item directly (`Delete(id)`) without wrapping it in an artificial slice or array literal (`Delete([]string{id})` or `delete([id])` is strictly banned).
- In tests, benchmarks, and seed fixtures, multiple items MUST be passed comma-separated (`Add("a", "b", "c")`), not as a pre-constructed slice unless explicitly testing slice spread mechanics.

### Rule 4: Parameter Limit Interoperability (<= 2-3 Parameters Rule)

- A variadic parameter occupies exactly **one parameter slot** in the function signature.
- Functions must still strictly adhere to the prompt architect constraint banning greater than 2–3 loose parameters.
- The variadic parameter MUST be placed as the final parameter in the signature.
- When auxiliary configuration options are needed alongside collection items, encapsulate configuration options into a dedicated `*Params` struct:

```go
func Dispatch(params DispatchParams, targets ...string) *appfault.AppError
```

### Rule 5: Safe Nil & Empty Handling

- Implementations must handle zero variadic arguments safely without panic or invalid indexing.
- In Go, `len(items) == 0` when zero arguments are passed. Always evaluate affirmatively: `hasItems := len(items) > 0`.
- In TypeScript, default rest array evaluates to `[]`. Check: `const hasItems = items.length > 0`.
- In Rust, borrowed slices `&[T]` check affirmative length: `let has_items = !items.is_empty();`.

### Rule 6: AppError & Vertical Spacing Conformance

- Go implementations returning failure metadata MUST return `*appfault.AppError`.
- Strict vertical line spacing must be observed:
  - Blank line before `if`.
  - Blank line after `}`.
  - Blank line before `return`.
- Strict boolean conventions:
  - `is` and `has` prefixes only.
  - Zero explicit `== true`.
  - Zero mixed polarity inside conditions.

---

## 3. Concrete Code Patterns & Anti-Patterns

### 3.1 Go (`...T` Variadic Slices)

#### ❌ Anti-Pattern: Rigid Slice Parameter

```go
// ❌ ANTI-PATTERN: Forces single-item callers to allocate and wrap
func InvalidateCacheKeys(keys []string) *appfault.AppError {
    if len(keys) == 0 {
        return nil
    }

    for _, key := range keys {
        if isErr := purgeKey(key); isErr != nil {
            return isErr
        }
    }

    return nil
}

// ❌ Call site suffering from wrapper boilerplate:
err := InvalidateCacheKeys([]string{sessionKey})
```

#### ✅ Canonical Pattern: Variadic Parameter with Spread Support

```go
// ✅ CANONICAL: Variadic parameter with zero-wrapper single calls and slice spread
func InvalidateCacheKeys(keys ...string) *appfault.AppError {
    hasKeys := len(keys) > 0
    if !hasKeys {
        return nil
    }

    for _, key := range keys {
        appErr := purgeKey(key)
        if appErr != nil {
            return appErr
        }
    }

    return nil
}

// Call site 1 (Single item — zero boilerplate):
err1 := InvalidateCacheKeys(sessionKey)

// Call site 2 (Multiple comma-separated items):
err2 := InvalidateCacheKeys("key-1", "key-2", "key-3")

// Call site 3 (Existing slice — unpacked with spread):
activeKeys := []string{"session-a", "session-b"}
err3 := InvalidateCacheKeys(activeKeys...)

// Call site 4 (Zero items — clean invocation):
err4 := InvalidateCacheKeys()
```

#### ✅ Canonical Pattern: Parameter Struct with Trailing Variadic Collection

```go
// ✅ Compliant with <= 2-3 parameters rule: 1 struct + 1 variadic parameter
func BulkUpdateUsers(params UserUpdateParams, userIDs ...string) *appfault.AppError {
    hasUsers := len(userIDs) > 0
    if !hasUsers {
        return nil
    }

    for _, id := range userIDs {
        appErr := executeUserUpdate(id, params)
        if appErr != nil {
            return appErr
        }
    }

    return nil
}
```

---

### 3.2 TypeScript (`...items: readonly T[]` & `SingleOrArray<T>`)

#### ❌ Anti-Pattern: Rigid Array Ingestion

```typescript
// ❌ ANTI-PATTERN: Rigid array forces array literals everywhere
export async function emitMetrics(metrics: MetricEvent[]): Promise<Result<void>> {
  if (metrics.length === 0) {
    return Result.ok(undefined);
  }

  for (const metric of metrics) {
    await recordMetric(metric);
  }

  return Result.ok(undefined);
}

// ❌ Call site: awkward brackets for the most common case
await emitMetrics([singleMetric]);
```

#### ✅ Canonical Pattern A: Rest Parameters (`...metrics`)

```typescript
// ✅ CANONICAL: Rest parameters with readonly array safety
export async function emitMetrics(
  ...metrics: readonly MetricEvent[]
): Promise<Result<void>> {
  const hasMetrics = metrics.length > 0;
  if (!hasMetrics) {
    return Result.ok(undefined);
  }

  for (const metric of metrics) {
    const result = await recordMetric(metric);
    if (!result.isSuccess) {
      return result;
    }
  }

  return Result.ok(undefined);
}

// Call site 1 (Single item — no brackets):
await emitMetrics(singleMetric);

// Call site 2 (Multiple comma-separated items):
await emitMetrics(metricA, metricB, metricC);

// Call site 3 (Spread existing array):
const backlog: readonly MetricEvent[] = [metric1, metric2];
await emitMetrics(...backlog);

// Call site 4 (Zero items):
await emitMetrics();
```

#### ✅ Canonical Pattern B: Single-or-Array Union Normalizer

When arguments are passed within a parameter struct or configuration options bag, use the `SingleOrArray<T>` union pattern paired with a zero-mutation normalizer:

```typescript
export type SingleOrArray<T> = T | readonly T[];

export function normalizeToArray<T>(input: SingleOrArray<T> | undefined): readonly T[] {
  const isUndefined = input === undefined;
  if (isUndefined) {
    return [];
  }

  const isArray = Array.isArray(input);
  if (isArray) {
    return input as readonly T[];
  }

  return [input as T];
}

// Configuration struct allowing single item or array:
export interface FilterOptions {
  readonly categories?: SingleOrArray<string>;
  readonly limit?: number;
}

export function applyFilters(options: FilterOptions): Result<FilteredView> {
  const categories = normalizeToArray(options.categories);

  return Result.ok(buildView(categories, options.limit));
}
```

---

### 3.3 Rust (`&[T]` & `impl IntoIterator<Item = T>`)

Rust does not support variadic function parameters in safe signatures. However, zero-cost slices and generic iterators achieve identical ergonomics and zero heap allocation overhead.

#### ❌ Anti-Pattern: Owned Heap Vector Parameter

```rust
// ❌ ANTI-PATTERN: Forces caller to allocate a Vec on the heap for every call
pub fn register_subscribers(subscribers: Vec<SubscriberId>) -> Result<(), AppError> {
    if subscribers.is_empty() {
        return Ok(());
    }

    for sub in subscribers {
        persist_subscriber(&sub)?;
    }

    Ok(())
}

// ❌ Call site forced to allocate heap buffer for a single ID:
register_subscribers(vec![active_id])?;
```

#### ✅ Canonical Pattern A: Borrowed Slice Reference

```rust
// ✅ CANONICAL PATTERN 1: Borrowed slice reference (zero allocation)
pub fn register_subscribers(subscribers: &[SubscriberId]) -> Result<(), AppError> {
    let has_subscribers = !subscribers.is_empty();
    if !has_subscribers {
        return Ok(());
    }

    for sub in subscribers {
        persist_subscriber(sub)?;
    }

    Ok(())
}

// Call site 1 (Single item borrowed from stack array):
register_subscribers(&[active_id])?;

// Call site 2 (Multiple comma-separated items in stack array):
register_subscribers(&[id_1, id_2, id_3])?;

// Call site 3 (Borrowed existing Vec):
let stored_ids: Vec<SubscriberId> = get_subscribers();
register_subscribers(&stored_ids)?;
```

#### ✅ Canonical Pattern B: Generic Collection Trait

```rust
// ✅ CANONICAL PATTERN 2: Generic iterator for ultimate flexibility
pub fn ingest_records<I, S>(records: I) -> Result<(), AppError>
where
    I: IntoIterator<Item = S>,
    S: AsRef<str>,
{
    for record in records {
        dispatch_record(record.as_ref())?;
    }

    Ok(())
}

// Call site: accepts standard array, vector, or std::iter::once without heap wrappers:
ingest_records(["alpha", "beta"])?;
ingest_records(std::iter::once("single-entry"))?;
```

---

## 4. Comprehensive Good vs Bad Contrast Matrix

| Scenario | ❌ Rigid Slice / Array Pattern | ✅ Variadic / Spread Pattern | Ergonomic & Performance Gain |
|:---|:---|:---|:---|
| **Single item invocation (Go)** | `fn([]string{"id"})` | `fn("id")` | Zero slice wrapper syntax; compiler escapes eliminated. |
| **Multiple items invocation (Go)** | `fn([]string{"a", "b", "c"})` | `fn("a", "b", "c")` | Clean comma-delimited call; no type annotation at call site. |
| **Existing slice invocation (Go)** | `fn(existingSlice)` | `fn(existingSlice...)` | Explicit spread operator `...`; identical performance. |
| **Empty collection invocation (Go)** | `fn(nil)` or `fn([]string{})` | `fn()` | Clean zero-argument invocation; slice initialized as empty. |
| **Single item invocation (TS)** | `fn(["user_1"])` | `fn("user_1")` | No array brackets; natural function invocation. |
| **Multiple items invocation (TS)** | `fn(["a", "b", "c"])` | `fn("a", "b", "c")` | Comma-separated arguments; matches standard JS rest conventions. |
| **Spread existing array (TS)** | `fn(items)` | `fn(...items)` | Unpacks collection cleanly with native JS spread operator. |
| **Single item invocation (Rust)** | `fn(vec!["id"])` (heap alloc) | `fn(&["id"])` or `fn(once("id"))` | Zero-allocation stack array borrow; no heap churn. |
| **Readability & Intent** | Emphasizes collection container over domain values | Emphasizes domain values; container is an implementation detail | Lower cognitive load; higher signal-to-noise ratio in code. |

---

## 5. Verification & Acceptance Criteria

- **AC-CG-033-A:** All file paths and cross-references within this document use strictly relative git paths.
- **AC-CG-033-B:** Function signatures across Go, TypeScript, and Rust adhere strictly to the maximum 2–3 parameter rule, placing variadic parameters in the trailing position.
- **AC-CG-033-C:** All Go code examples utilize `*appfault.AppError` for structured error handling.
- **AC-CG-033-D:** All code examples strictly observe vertical blank line spacing rules (blank lines before `if`, after `}`, and before `return`).
- **AC-CG-033-E:** All boolean conditions evaluate affirmatively using implicit checks with `is` or `has` prefixes, completely eliminating `== true` and mixed-polarity checks.

---

## 6. Related Specifications

- [`02-spec/02-coding-guidelines/01-cross-language/13-strict-typing.md`](13-strict-typing.md) — Strict typing rules and parameter count limits
- [`02-spec/02-coding-guidelines/01-cross-language/18-code-mutation-avoidance.md`](18-code-mutation-avoidance.md) — Immutability patterns and avoiding side effects
- [`02-spec/02-coding-guidelines/01-cross-language/32-branch-immutability-and-clean-construction.md`](32-branch-immutability-and-clean-construction.md) — Branch immutability and condition decomposition
- [`02-spec/03-error-manage/01-apperror-architecture.md`](../../03-error-manage/01-apperror-architecture.md) — `*appfault.AppError` structured error return architecture
- [`02-spec/21-app/01-variadic-spread-params-and-multi-repo/01-architecture-spec.md`](../../21-app/01-variadic-spread-params-and-multi-repo/01-architecture-spec.md) — Parent architecture specification
