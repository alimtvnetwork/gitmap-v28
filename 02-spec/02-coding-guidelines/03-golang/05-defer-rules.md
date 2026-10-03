# Go Defer Rules (AI Execution Prompt)

> **/goal** Prevent deferred resource exhaustion, LIFO execution confusion, and loop-defer memory leaks in Go by enforcing a maximum of one `defer` per function and banning `defer` inside loops.
> **/learn** Master Go defer mechanics: LIFO execution ordering, deferred function evaluation, function-scoped cleanup lifetime, loop-defer stack accumulation anti-patterns, and subroutine decomposition.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Restrict every Go function to at most one single `defer` statement.
- [ ] `/learn` Never place `defer` inside a `for` or `range` loop; decompose loop bodies into helper functions to ensure prompt resource release.
- [ ] `/goal` Position `defer` statements immediately following resource acquisition or at function boundaries, never buried mid-routine.
- [ ] `/learn` Decompose functions with multiple resources into nested helper functions or closure wrappers, each managing a single defer.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Version:** 3.2.0
**Updated:** 2026-04-16
**Source:** Consolidated from `01-pre-code-review-guides/03-golang-code-review-guides.md`

---

## 1. Rule

**Do NOT use more than one `defer` in a single function.** Defers work as a **stack** (LIFO), which makes multiple defers complicated and hard to maintain.

---

## 2. Placement

If using `defer`, place it at the **top** or **bottom** of the function — never in the middle buried between other logic.

```go
// ✅ CORRECT — single defer, at top
func ProcessFile(path string) error {
    f, err := os.Open(path)

    if err != nil {
        return apperror.Wrap(err, apperror.ErrFileOpen, "open failed")
    }

    defer f.Close()

    // ... process file
    return nil
}
```

```go
// ❌ WRONG — multiple defers, hard to reason about execution order
func ProcessData(ctx context.Context) error {
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    lock.Lock()
    defer lock.Unlock()

    f, err := os.Open("data.txt")
    if err != nil {
        return err
    }
    defer f.Close()

    // Which closes first? f.Close → lock.Unlock → tx.Rollback (LIFO)
    // This is confusing and error-prone
}
```

---

## 3. Alternative to Multiple Defers

Extract into separate functions, each with its own single defer:

```go
// ✅ CORRECT — each function has at most one defer
func ProcessData(ctx context.Context) error {
    return withTransaction(ctx, func(tx *sql.Tx) error {
        return processWithFile(tx, "data.txt")
    })
}

func processWithFile(tx *sql.Tx, path string) error {
    f, err := os.Open(path)

    if err != nil {
        return apperror.Wrap(err, apperror.ErrFileOpen, "open failed")
    }

    defer f.Close()

    // ... process
    return nil
}
```

---

## 4. Cross-References

- [Master Coding Guidelines §6](../01-cross-language/15-master-coding-guidelines/readme.md) — Error handling
- [Golang Standards Reference](./04-golang-standards-reference/readme.md) — Go conventions

---

*Go defer rules — consolidated from pre-code review guides.*

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-GO-005: Resource Defer Rules and Loop Safety Standards

**Given** Go source code across packages and CLI modules.
**When** Guideline linters audit the codebase.
**Then** Go functions contain at most one defer statement, zero defers execute inside loop bodies, and resources are closed deterministically with zero leaks and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang --check-only
```
**Expected:** exit 0. Zero violations detected.
