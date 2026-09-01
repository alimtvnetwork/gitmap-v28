# Code Style Rules

> **Related specs:**
> - [02-project-structure.md](02-project-structure.md) — package layout these rules apply within
> - [12-testing.md](12-testing.md) — test conventions that complement code style
> - [15-constants-reference.md](15-constants-reference.md) — naming conventions for constants

## Mandatory Constraints

| Constraint | Rule |
|------------|------|
| `if` conditions | Always positive — no `!`, no `!=` |
| Function length | 8–15 lines (excluding blanks and comments) |
| File length | 100–200 lines max |
| Package granularity | One responsibility per package |
| Newline before `return` | Always, unless `return` is the only line in an `if` |
| No magic strings | All literals in `constants` package |

## Conditionals — No Negation

Always write positive conditions:

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

**Rationale:** "If X exists, do Y" reads more naturally than
"If X is not missing, do Y."

## Function Length — 8–15 Lines

- Each function does one thing.
- If you need a comment to explain a section, that section should
  be its own function.
- Use named helpers instead of complex inline logic.

## File Length — 100–200 Lines

When a file exceeds 200 lines, split by responsibility:

| Signal | Action |
|--------|--------|
| 2+ unrelated function groups | Split into separate files |
| Large switch statement | Each case → own file |
| Types mixed with logic | Separate model from logic |

## Return Formatting

Always add a blank line before `return`, unless the `return` is the
only line inside an `if` block:

```go
// ✅ Correct
func process(data []string) int {
    count := len(data)
    filtered := filter(data)

    return len(filtered)
}

// ✅ Also correct — sole line in if
if isValid(input) {
    return input
}
```

## No Magic Strings

Every string literal used for comparison, format templates, defaults,
extensions, or messages must be a constant:

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

## Naming Conventions

| Element | Convention | Example |
|---------|-----------|---------|
| Package names | Lowercase, single word | `scanner`, `formatter` |
| Exported functions | PascalCase, verb-led | `BuildRecords`, `WriteCSV` |
| Unexported functions | camelCase, verb-led | `parseFlags`, `resolveDir` |
| Constants | PascalCase | `DefaultBranch`, `ModeHTTPS` |
| Files | Lowercase, single word | `terminal.go`, `csv.go` |
| Test files | `*_test.go` | `mapper_test.go` |

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
