# Newline Styling Examples

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**Applies to:** All languages (Go examples)  
**Source:** Consolidated from `01-pre-code-review-guides/03-golang-code-review-guides.md`  
**Supplements:** [Code Style](./04-code-style/00-overview.md) rules R4, R5, R10, R12, R13

---

## 1. Purpose

Detailed before/after examples for newline rules. These supplement the formal rules in [04-code-style.md](./04-code-style/00-overview.md).

---

## 2. Blank Line Before `return` (Rule R4)

### Single-line function — NO blank line needed

```go
func Something() int {
    return constants.One  // ✅ alright — single statement
}
```

### Multi-line function — blank line REQUIRED

```go
// ❌ WRONG — no blank line before return
func Something() int {
    doSomething()
    return constants.One
}

// ✅ CORRECT
func Something() int {
    doSomething()

    return constants.One
}
```

### Inside blocks — same rule applies

```go
// ❌ WRONG — no blank line before return in block
func Something() int {
    doSomething()
    if isNumber {
        doSomethingNew()

        return constants.One
    }
}

// ✅ CORRECT — blank line before both returns
func Something() int {
    doSomething()

    if isNumber {
        doSomethingNew()

        return constants.One
    }

    return constants.MinusOne
}
```

---

## 3. Blank Line After `}` When Followed by Code (Rule R5)

```go
// ✅ CORRECT — blank line after } when followed by more code (Rule R5)
func Something() int {
    if isInapplicable1 {
        return constants.Zero
    }

    if isInapplicable2 {
        return constants.Zero
    }

    return constants.One
}
```

---

## 4. No Empty Line at Start of Function (Rule R12)

```go
// ❌ WRONG — empty line at start
func Something() int {

    doSomething()

    return constants.One
}

// ✅ CORRECT — comment at first line is fine
func Something() int {
    // process the thing
    doSomething()

    return constants.One
}
```

---

## 5. No Double Empty Lines (Rule R13 extension)

```go
// ❌ WRONG — double empty line
func Something() int {
    doSomething()


    return constants.One
}
```

---

## 6. Blank Line After `}` (Rule R5)

```go
// ✅ CORRECT — blank line after } when followed by more code
if guardA {
    return
}

if guardB {
    return
}
```

---

## 7. Newline in Go Outputs

Use `constants.NewLineUnix` (`"\n"`) in 90% of cases. Only use `constants.NewLine` (OS-specific) when the user explicitly needs OS-dependent newline handling.

| Constant | Value | When to Use |
|----------|-------|-------------|
| `constants.NewLineUnix` | `"\n"` | Default — 90% of cases |
| `constants.NewLine` | OS-dependent | Only when OS-specific newline is needed (e.g., IDE file saving) |

---

## 8. Cross-References

- [Code Style](./04-code-style/00-overview.md) — Formal rule definitions (R4, R5, R10, R12, R13)
- [Master Coding Guidelines §5](./15-master-coding-guidelines/00-overview.md) — Formatting rules summary

---

*Newline styling examples — consolidated from pre-code review guides.*
