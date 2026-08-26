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
- [Master Coding Guidelines](./15-master-coding-guidelines/00-overview.md) — §7 Type Safety

---

*Code mutation avoidance — consolidated from pre-code review guides.*
