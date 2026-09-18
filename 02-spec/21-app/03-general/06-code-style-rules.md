# Code Style Rules

## Overview

This document defines enforceable coding conventions for Go CLI projects.
These rules prioritize readability, maintainability, and consistency.

## Conditionals

### No Negation in `if` Conditions

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

```go
// ✅ Correct
if fileExists(path) {
    loadFile(path)
}

// ❌ Wrong
if !fileMissing(path) {
    loadFile(path)
}
```

### Rationale

Positive conditions are easier to reason about. "If X exists, do Y"
reads more naturally than "If X is not missing, do Y."

## Function Length

### Target: 8–15 Lines

Functions should be 8–15 lines of code (excluding blank lines and
comments). If a function exceeds 15 lines, extract a helper.

### How to Enforce

- Each function should do one thing.
- If you need a comment to explain a section within a function,
  that section should be its own function.
- Use named helper functions instead of complex inline logic.

## File Length

### Target: 100–200 Lines Max

No source file should exceed 200 lines. If it does, split it by
responsibility.

### Splitting Strategy

| Signal | Action |
|--------|--------|
| File has 2+ unrelated groups of functions | Split into separate files |
| File has a large switch statement | Each case becomes a file |
| File mixes types and logic | Separate `model.go` from logic |

## Package Granularity

### One Responsibility Per Package

Each package owns a single concern:

| Package | Responsibility |
|---------|----------------|
| `cmd` | CLI routing and flag parsing |
| `config` | Config file loading and merging |
| `constants` | All shared string literals |
| `scanner` | Directory walking |
| `mapper` | Data transformation |
| `formatter` | Output rendering |
| `model` | Shared data structures |

### Rules

- No circular imports between packages.
- `cmd` orchestrates; other packages never import `cmd`.
- `model` and `constants` are leaf packages (imported by many,
  import nothing project-specific).

## Return Formatting

### Blank Line Before `return`

Always add a blank line before `return`, unless the `return` is the
only line inside an `if` block.

```go
// ✅ Correct — blank line before return
func process(data []string) int {
    count := len(data)
    filtered := filter(data)

    return len(filtered)
}

// ✅ Correct — sole line in if, no blank needed
if isValid(input) {
    return input
}

// ❌ Wrong — no blank line before return
func process(data []string) int {
    count := len(data)
    filtered := filter(data)
    return len(filtered)
}
```

## No Magic Strings

### All Literals in `constants` Package

Every string literal used for comparison, format templates, default
values, file extensions, or CLI messages must be defined as a constant.

```go
// ✅ Correct
if mode == constants.ModeHTTPS {
    url = formatHTTPS(repo)
}

// ❌ Wrong
if mode == "https" {
    url = formatHTTPS(repo)
}
```

### What Must Be in Constants

| Category | Examples |
|----------|---------|
| CLI command names | `CmdScan = "scan"` |
| File extensions | `ExtCSV = ".csv"` |
| Default values | `DefaultBranch = "main"` |
| Format strings | `CloneInstructionFmt = "git clone -b %s %s %s"` |
| Error messages | `ErrSourceRequired = "Error: source file is required"` |
| ANSI codes | `ColorGreen = "\033[32m"` |
| UI strings | `TermBannerTop = "╔══════..."` |

### What Does NOT Need to Be a Constant

- Struct field names (Go handles these).
- Local variable names.
- Test data strings.
- Log messages that are unique to one location and never compared.

## Naming Conventions

| Element | Convention | Example |
|---------|-----------|---------|
| Package names | Lowercase, single word | `scanner`, `formatter` |
| Exported functions | PascalCase, verb-led | `BuildRecords`, `WriteCSV` |
| Unexported functions | camelCase, verb-led | `parseFlags`, `resolveDir` |
| Constants | PascalCase | `DefaultBranch`, `ModeHTTPS` |
| Files | Lowercase, single word | `terminal.go`, `csv.go` |

## Error Handling

- Always check errors immediately after the call.
- Return errors up the stack; let the caller decide what to do.
- In `cmd` package handlers, print the error and `os.Exit(1)`.
- Never use `panic` for expected error conditions.
