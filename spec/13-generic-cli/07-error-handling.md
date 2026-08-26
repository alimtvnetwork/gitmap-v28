# Error Handling

> **Related specs:**
> - [15-constants-reference.md](15-constants-reference.md) — error message constants (`Err*` naming)
> - [16-verbose-logging.md](16-verbose-logging.md) — verbose debug output for error diagnosis
> - [18-batch-execution.md](18-batch-execution.md) — continue-on-failure pattern for batch operations

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | User error (bad args, missing file, invalid input) |
| Non-zero | Propagated from child processes |

## Error Message Rules

| Rule | Detail |
|------|--------|
| All error format strings in `constants` | `ErrSourceRequired`, `ErrConfigLoad`, etc. |
| Errors print to stderr | Never stdout |
| Exit immediately after error | Don't continue with bad state |
| Messages are actionable | Tell the user what to do, not just what failed |

### Example

```go
// constants/constants_messages.go
const (
    ErrSourceRequired = "Error: source file is required\nUsage: toolname clone <source>"
    ErrConfigLoad     = "Error: could not load config from %s"
    ErrRepoNotFound   = "Error: no repo matches slug '%s'. Run 'toolname scan' first."
)
```

```go
// cmd/clone.go
if source == "" {
    fmt.Fprintln(os.Stderr, constants.ErrSourceRequired)
    os.Exit(1)
}
```

## Batch Operations

For commands that process N items (e.g., pull all repos):

1. **Log per-item failures** but continue processing.
2. **Track success/failure counts**.
3. **Print summary at the end**.
4. **Exit with code 1 if any failures**.

```go
var failed int
for _, repo := range repos {
    err := pull(repo)
    if err != nil {
        fmt.Fprintf(os.Stderr, "[%s] %v\n", repo.Name, err)
        failed++
        continue
    }
    fmt.Printf("[%s] ✓\n", repo.Name)
}

fmt.Printf("\n%d succeeded, %d failed\n", len(repos)-failed, failed)
if failed > 0 {
    os.Exit(1)
}
```

## Error Handling in Functions

- Always check errors immediately after the call.
- Return errors up the stack; let the caller decide.
- In `cmd` package handlers, print the error and `os.Exit(1)`.
- Never use `panic` for expected error conditions.

```go
// ✅ Correct — return error up
func loadFile(path string) ([]byte, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("reading %s: %w", path, err)
    }

    return data, nil
}

// ✅ Correct — handler prints and exits
func runImport(args []string) {
    data, err := loadFile(args[0])
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    // process data
}
```

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
