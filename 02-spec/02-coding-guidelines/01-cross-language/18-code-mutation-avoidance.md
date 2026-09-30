# Code Mutation Avoidance

**Version:** 3.2.0
**Updated:** 2026-04-16
**Applies to:** Go (primary), general principle cross-language
**Source:** Consolidated from `01-pre-code-review-guides/03-golang-code-review-guides.md`

---

## 1. Principle

**Every variable should be assigned only once.** Prefer `const` or immutable objects. Variable mutation is a code review red flag.

---

## 2. Rules

### Rule 1: Single Assignment

Variables should be assigned once and not modified after. At worst case:

- A loop collecting data (slice append) may reassign
- If a variable must change, it should be done in **one method**, not across multiple

### Rule 2: No Post-Construction Mutation

```go
// ❌ WRONG — mutation after construction
returningVal := New(params)
returningVal.SetName("updated")
returningVal.VarName = "modified"

return returningVal
```

```go
// ✅ CORRECT — pass all values to constructor
varName := whatEverValue

returningVal := New(varName, params)

return returningVal
```

```go
// ✅ ALSO CORRECT — struct literal
varName := whatEverValue

return ReturningValStruct{
    VarName: varName,
}
```

### Rule 3: Mutex Lock for Mutable State

When mutation is unavoidable (lazy evaluation, caching), use mutex locks:

```go
func (r *Receiver) GetLinesLock() []string {
    r.Lock()
    defer r.Unlock()

    return r.GetLines()
}
```

Provide both locked and unlocked versions:

- `GetLines()` — non-locked (for single-goroutine use)
- `GetLinesLock()` — locked (for concurrent use)

### Rule 4: Return-from-Function for Type Conversion & Branching (TOTAL BAN on Intermediate Mutation)

Never declare an unassigned or zero-valued variable before a `switch` or `if/else` ladder and reassign it across cases. Instead, extract the logic into a pure helper function that returns immediately from each branch, or use a dedicated conversion package.

```go
// ❌ WRONG — mutating local variable across branches (Anti-Pattern)
var data []byte
switch v := payload.(type) {
case []byte:
    data = v
case string:
    data = []byte(v)
default:
    data = []byte(fmt.Sprint(v))
}

self.Destination().Write(data)
```

```go
// ✅ CORRECT — pure conversion function with early returns
func toBytes(payload any) []byte {
    switch v := payload.(type) {
    case []byte:
        return v
    case string:
        return []byte(v)
    case []string:
        return []byte(strings.Join(v, "\n") + "\n")
    default:
        return []byte(fmt.Sprint(v))
    }
}

// Call site receives immutable result directly:
data := toBytes(payload)
self.Destination().Write(data)
```

### Rule 5: Standalone Payload Conversion Architecture

Data conversion from generic representations (`any`) into byte streams or disk payloads must be encapsulated in a dedicated standalone package (e.g. `payloadconv`).

Polymorphic conversion requirements:
- `[]byte`: preserved directly without re-encoding.
- `string`: encoded directly to UTF-8 bytes.
- `[]string` / arrays: serialized line-by-line (`strings.Join(lines, "\n") + "\n"`).
- `struct` / `map`: automatically serialized to formatted JSON with trailing newline.
- `error`: serialized as error message string.

### Rule 6: File Writing and File-Path Concurrency Safety (`fileutil`)

File writing operations must be safe, flexible, and concurrency-guarded:

1. **Unified Writing Entry Point:** Provide a generic `Write(path, payload, perm)` that accepts any payload (struct, map, array, string, bytes) and routes through the standalone converter.
2. **Dedicated Typed Writers:** Provide explicit helpers (`WriteJSON`, `WriteLines`, `WriteString`, `WriteBytes`).
3. **File-Path Based Mutex Locking:** File writes in concurrent contexts must support locked variants (`WriteLocked`, `WriteJSONLocked`, etc.) resolving a path-level `sync.RWMutex`.

### Rule 7: Branch Immutability & Constructor Helper Return Pattern (TOTAL BAN on In-Place Branch Mutation)

Never mutate fields on an existing struct or object instance across conditional branches (`if`/`else` ladders or `switch` cases).

```go
// ❌ WRONG — mutating fields across branching conditions
if len(items) >= 2 && !isValid(items[0]) && isValid(items[1]) {
    cfg.Target = items[1]
    if cfg.Alias == "" {
        cfg.Alias = items[0]
    }
    return cfg, nil
}
```

Instead, decompose branching paths so each branch calls a pure constructor or builder helper that constructs and returns a complete, immutable struct instance directly:

```go
// ✅ CORRECT — pure construction with dedicated helpers
if isAlternateOrder {
    return buildAlternateConfig(items, baseConfig), nil
}

return buildStandardConfig(items, baseConfig), nil
```

See [32-branch-immutability-and-clean-construction.md](./32-branch-immutability-and-clean-construction.md) for full architecture and cross-language specifications.

---

## 3. Exemptions

| Case | Why Allowed |
|------|-------------|
| Lazy evaluation / caching | Value generated once, then immutable |
| Loop accumulation (`append`) | Collecting results from iteration |
| Builder pattern | Accumulates instructions, builds once |
| Design pattern implementations | Some patterns require mutable state |

---

## 4. Anti-Patterns

### ❌ Variable Modified Across Multiple Methods

```go
// DANGEROUS — value changes unpredictably
result := NewResult()
methodA(result)  // modifies result
methodB(result)  // modifies result again
methodC(result)  // and again

return result
```

### ❌ Mutation Without Lock in Concurrent Context

```go
// DANGEROUS — race condition
func (r *Receiver) SetValue(v string) {
    r.value = v  // no lock!
}
```

---

## 5. Cross-References

- [Lazy Evaluation Patterns](./16-lazy-evaluation-patterns.md) — Exempted mutation for caching
- [DRY Principles](./08-dry-principles.md) — Constructor-based initialization
- [Master Coding Guidelines](./15-master-coding-guidelines/readme.md) — §7 Type Safety

---

*Code mutation avoidance — consolidated from pre-code review guides.*
