# Go Defer Rules

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

- [Master Coding Guidelines §6](../01-cross-language/15-master-coding-guidelines/00-overview.md) — Error handling
- [Golang Standards Reference](./04-golang-standards-reference/00-overview.md) — Go conventions

---

*Go defer rules — consolidated from pre-code review guides.*
