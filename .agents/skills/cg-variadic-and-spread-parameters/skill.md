---
name: cg-variadic-and-spread-parameters
description: Autonomously audits, refactors, and enforces variadic and spread parameters (...T, ...items, &[T]) across polyglot codebases to eliminate artificial slice wrappers.
---

# Skill: Variadic & Spread Parameters, Rest Elements & Slices (`cg-variadic`)

This skill governs autonomous execution for auditing, refactoring, and enforcing variadic parameters, rest elements, and slice reference patterns across Go, TypeScript, and Rust codebases.

## Trigger Keywords & Aliases

- `cg-variadic`
- `cg-spread-params`
- `cg-execute variadic`
- `audit variadic`
- `fix variadic params`

---

## 1. Mandatory Architectural Rules

1. **Variadic Parameter Primacy:**
   - Functions accepting homogeneous collections of items MUST use variadic signatures (`...T` in Go, `...items: readonly T[]` in TypeScript, `&[T]` in Rust) instead of rigid slice or array types whenever single-item or variable-length calls occur.
   - Total ban on forcing callers to construct artificial wrapper slices/arrays for single-element invocations.

2. **Zero-Wrapper Single Invocations:**
   - Single-item callers MUST pass the argument directly without wrapper syntax (`DeleteUsers(id)` instead of `DeleteUsers([]string{id})`).

3. **Seamless Spread Forwarding:**
   - Callers possessing an existing slice or array must be able to unpack it directly using native spread syntax (`DeleteUsers(allIDs...)` in Go, `deleteUsers(...allIds)` in TypeScript, `process_records(&records)` in Rust).

4. **Trailing Position & Parameter Limits:**
   - The variadic parameter MUST be the final parameter in the signature.
   - Total function parameters must not exceed 2–3 loose arguments. Pack complex configuration options into a dedicated `*Params` struct placed before the variadic parameter.

5. **Safe Empty & Nil Handling:**
   - Functions MUST handle zero arguments cleanly without error.
   - Affirmative length checks must guard processing (`hasItems := len(items) > 0`).
   - If not `hasItems`, return early cleanly.

6. **Zero-Build & Zero-Test Routine Execution Mandate:**
   - NEVER execute `go test`, `pytest`, `npm test`, or full build commands during routine execution turns.
   - Verify code using targeted file-level linters and AST checkers on specifically modified files (`exit 0`).

---

## 2. Canonical Code Transformations

### 2.1 Go (`...T` Variadic Slices)

#### ❌ Anti-Pattern: Rigid Slice Parameter

```go
// ❌ WRONG — forces artificial slice literal at every single-item call site
func PurgeTags(resourceID string, tags []string) *appfault.AppError {
    if len(tags) == 0 {
        return nil
    }

    for _, tag := range tags {
        appErr := removeTag(resourceID, tag)
        if appErr != nil {
            return appErr
        }
    }

    return nil
}

// Call site: awkward slice literal for single tag
appErr := PurgeTags("res-101", []string{"obsolete"})
```

#### ✅ Canonical Pattern: Variadic Parameter

```go
// ✅ REQUIRED — variadic slice parameter with clean affirmative check
func PurgeTags(resourceID string, tags ...string) *appfault.AppError {
    hasTags := len(tags) > 0
    if !hasTags {
        return nil
    }

    for _, tag := range tags {
        appErr := removeTag(resourceID, tag)
        if appErr != nil {
            return appErr
        }
    }

    return nil
}

// Call site 1 (Single item — zero wrapper):
appErr1 := PurgeTags("res-101", "obsolete")

// Call site 2 (Multiple comma-separated items):
appErr2 := PurgeTags("res-101", "obsolete", "deprecated", "stale")

// Call site 3 (Existing slice unpacked with spread):
batch := []string{"tag-a", "tag-b"}
appErr3 := PurgeTags("res-101", batch...)

// Call site 4 (Zero items):
appErr4 := PurgeTags("res-101")
```

---

### 2.2 TypeScript (`...items: readonly T[]` & `SingleOrArray<T>`)

#### ❌ Anti-Pattern: Rigid Array Parameter

```typescript
// ❌ WRONG — forces array literal for single event
export async function trackEvents(events: string[]): Promise<Result<void>> {
  if (events.length === 0) {
    return Result.ok(undefined);
  }

  for (const event of events) {
    await dispatchEvent(event);
  }

  return Result.ok(undefined);
}

// Call site:
await trackEvents(['user_signup']);
```

#### ✅ Canonical Pattern: Rest Parameters with Readonly Array

```typescript
// ✅ REQUIRED — Rest parameter with readonly safety
export async function trackEvents(
  ...events: readonly string[]
): Promise<Result<void>> {
  const hasEvents = events.length > 0;
  if (!hasEvents) {
    return Result.ok(undefined);
  }

  for (const event of events) {
    const result = await dispatchEvent(event);
    if (!result.isSuccess) {
      return result;
    }
  }

  return Result.ok(undefined);
}

// Call site 1 (Single item):
await trackEvents('user_signup');

// Call site 2 (Multiple items):
await trackEvents('click_nav', 'open_modal');

// Call site 3 (Spread existing array):
const backlog = ['event_1', 'event_2'];
await trackEvents(...backlog);
```

#### ✅ Canonical Pattern: Configuration Struct Normalizer

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
```

---

### 2.3 Rust (Borrowed Slices `&[T]` & Generic Iterators)

#### ❌ Anti-Pattern: Owned Heap Vector Parameter

```rust
// ❌ WRONG — forces heap allocation for every call
pub fn process_identifiers(ids: Vec<String>) -> Result<(), AppError> {
    if ids.is_empty() {
        return Ok(());
    }

    for id in ids {
        execute_step(&id)?;
    }

    Ok(())
}

// Call site: allocates heap buffer for single item
process_identifiers(vec!["id-100".to_string()])?;
```

#### ✅ Canonical Pattern: Borrowed Slice Reference

```rust
// ✅ REQUIRED — zero-allocation slice reference
pub fn process_identifiers(ids: &[&str]) -> Result<(), AppError> {
    let has_items = !ids.is_empty();
    if !has_items {
        return Ok(());
    }

    for id in ids {
        execute_step(id)?;
    }

    Ok(())
}

// Call site 1 (Single item on stack):
process_identifiers(&["id-100"])?;

// Call site 2 (Borrowed existing Vec):
let my_vec = vec!["a", "b"];
process_identifiers(&my_vec)?;
```

---

## 3. High-Speed AST Search & Discovery Patterns

Use GitMap live streaming search to locate candidate signatures and wrapper call sites:

```bash
# Go: Find function definitions with slice parameters
gitmap aum search -r "func [A-Za-z0-9_]+\([^)]*\[\][A-Za-z0-9_.*]+(\s*\))?" -e .go

# Go: Find call sites wrapping single items in slice literals
gitmap aum search -r "\b[A-Za-z0-9_]+\(\[\]string\{" -e .go

# TypeScript: Find function definitions with array parameters
gitmap aum search -r "\bfunction [A-Za-z0-9_]+\([^)]*:\s*[A-Za-z0-9_]+\[\]" -e .ts

# TypeScript: Find call sites with single-element array literals
gitmap aum search -r "\b[A-Za-z0-9_]+\(\[[A-Za-z0-9_'\"]+\]\)" -e .ts

# Rust: Find functions taking owned Vec<T>
gitmap aum search -r "fn [A-Za-z0-9_]+\([^)]*Vec<" -e .rs
```

---

## 4. Verification Checklist & Linters

- [ ] Variadic parameters placed in trailing position.
- [ ] No more than 2–3 loose parameters (use `*Params` struct for configuration).
- [ ] Call sites refactored to eliminate artificial slice wrappers (`[]string{x}`).
- [ ] Affirmative boolean guards (`is`, `has`) used for empty collection checks.
- [ ] Functions adhere to <= 8–15 lines.
- [ ] Targeted linters exit 0:
  - `python linter-scripts/check-relative-paths.py`
  - `python 03-ai-scripts/21-sequence-integrity-linter.py 01-prompts/15-cg-execute`

---

## 5. Cross-References

- [01-architecture-spec.md](02-spec/21-app/01-variadic-spread-params-and-multi-repo/01-architecture-spec.md)
- [33-variadic-and-spread-parameters.md](02-spec/02-coding-guidelines/01-cross-language/33-variadic-and-spread-parameters.md)
- [36-variadic-and-spread-parameters.md](01-prompts/15-cg-execute/36-variadic-and-spread-parameters.md)
