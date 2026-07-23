# Go Code Style Rules

## Overview

Enforceable coding conventions for Go CLI projects, prioritizing
readability, maintainability, and consistency.

---

## 1. Positive Conditionals Only

Always write positive conditions. No `!`, no `!=`.

```go
// ✅ Correct
if len(args) > 0 {
    process(args)
}

// ❌ Wrong
if len(args) != 0 {
    process(args)
}
```

---

## 2. Function Length — 8 to 15 Lines

Functions should be 8–15 lines (excluding blanks and comments).
Extract a helper when exceeding 15 lines. Each function does one thing.

---

## 3. File Length — Max 200 Lines

| Signal | Action |
|--------|--------|
| 2+ unrelated function groups | Split into separate files |
| Large switch statement | Each case becomes a file |
| Mixed types and logic | Separate `model.go` from logic |

---

## 4. Package Granularity — One Responsibility

| Package | Responsibility |
|---------|----------------|
| `cmd` | CLI routing and flag parsing |
| `config` | Config loading and merging |
| `constants` | All shared string literals |
| `scanner` | Directory walking |
| `mapper` | Data transformation |
| `formatter` | Output rendering |
| `model` | Shared data structures |

Rules:
- No circular imports.
- `cmd` orchestrates; other packages never import `cmd`.
- `model` and `constants` are leaf packages.

---

## 5. Blank Line Before `return`

```go
// ✅ Correct
func process(data []string) int {
    filtered := filter(data)

    return len(filtered)
}

// ✅ OK — sole line in if
if isValid(input) {
    return input
}
```

---

## 6. No Magic Strings

Every literal used for comparison, defaults, or messages → `constants`.

| Must Be Constants | Does Not Need Constants |
|-------------------|------------------------|
| CLI command names | Struct field names |
| File extensions | Local variable names |
| Default values | Test data strings |
| Format strings | Unique log messages |
| Error messages | |

---

## 7. Naming Conventions

| Element | Convention | Example |
|---------|-----------|---------|
| Package | Lowercase, single word | `scanner` |
| Exported func | PascalCase, verb-led | `BuildRecords` |
| Unexported func | camelCase, verb-led | `parseFlags` |
| Constants | PascalCase | `DefaultBranch` |
| Files | Lowercase, single word | `terminal.go` |

---

## 8. Error Handling

- Check errors immediately after the call.
- Return errors up the stack; let the caller decide.
- In `cmd` handlers: print error and `os.Exit(1)`.
- Never `panic` for expected conditions.

---

## References

- Universal rules: `spec/05-coding-guidelines/01-code-quality-improvement.md`
- Generic CLI style: `spec/04-generic-cli/08-code-style.md`
