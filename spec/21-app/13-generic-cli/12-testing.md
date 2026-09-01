# Testing

> **Related specs:**
> - [08-code-style.md](08-code-style.md) — code style rules that tests must also follow
> - [11-build-deploy.md](11-build-deploy.md) — build pipeline that executes tests
> - [02-project-structure.md](02-project-structure.md) — test file placement within the package layout

## Test Structure

```
toolname/
├── mapper/
│   ├── mapper.go
│   └── mapper_test.go        ← unit tests next to source
├── config/
│   ├── config.go
│   └── config_test.go
└── tests/
    └── cmd_test/
        └── scan_test.go      ← integration tests in separate dir
```

## Conventions

| Convention | Detail |
|------------|--------|
| Unit tests | Same package, same directory as source |
| Integration tests | Under `tests/` in separate packages |
| Test file naming | `*_test.go` matching source file |
| Test function naming | `TestFunctionName_Scenario` |
| Table-driven tests | Use for functions with multiple input/output cases |
| Test data | Inline in test, not external files (unless large) |

## Table-Driven Tests

```go
func TestSlugFromURL(t *testing.T) {
    tests := []struct {
        name string
        url  string
        want string
    }{
        {"basic", "https://github.com/user/repo.git", "repo"},
        {"no .git", "https://github.com/user/repo", "repo"},
        {"empty", "", ""},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := SlugFromURL(tt.url)
            if got != tt.want {
                t.Errorf("SlugFromURL(%q) = %q, want %q", tt.url, got, tt.want)
            }
        })
    }
}
```

## What to Test

| Layer | What to Test |
|-------|-------------|
| `mapper` | Data transformation correctness |
| `config` | Merge priority (defaults → file → flags) |
| `formatter` | Output matches expected format (use `io.Writer`) |
| `store` | CRUD operations with in-memory SQLite |
| `cmd` | Flag parsing returns correct values |
| `scanner` | Detection rules match expected patterns |

## What NOT to Test

- `main.go` (trivial delegation)
- ANSI color output (visual, not behavioral)
- Third-party library internals

## Testable Design

All formatters accept `io.Writer` so tests can capture output:

```go
func TestWriteCSV(t *testing.T) {
    var buf bytes.Buffer
    records := []Record{{Name: "test", URL: "https://example.com"}}
    WriteCSV(&buf, records)

    output := buf.String()
    if !strings.Contains(output, "test") {
        t.Errorf("CSV output missing record name")
    }
}
```

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
