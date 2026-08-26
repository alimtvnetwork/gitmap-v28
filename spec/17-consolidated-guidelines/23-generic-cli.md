# Generic CLI Creation Guidelines

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**Source Module:** `13-generic-cli/`  
**AI Confidence:** Production-Ready  
**Ambiguity:** None

---

## Purpose

Complete, self-contained blueprint for building production-quality CLI tools. Language-agnostic principles with Go concrete examples. Covers project structure, subcommand dispatch, flag parsing, configuration layering, output formatting, error handling, code style, help system, database persistence, build/deploy, testing, and rich terminal output.

---

## Design Philosophy

| Principle | Detail |
|-----------|--------|
| Consistency over cleverness | Predictable patterns across all commands |
| Convention over configuration | Sensible defaults; config is optional |
| Fail fast, fail clearly | Bad input → immediate error with actionable message |
| One responsibility per unit | Each file, function, and package does one thing |
| No magic strings | Every literal in a constants package |
| Self-documenting | Help text, version, and examples built into the binary |

---

## Project Structure

Standard layout for a CLI project:

```
<project>/
├── cmd/
│   └── <binary>/
│       └── main.go              # Entry point — parse args, dispatch
├── internal/
│   ├── commands/                # One file per subcommand
│   │   ├── list.go
│   │   ├── add.go
│   │   └── remove.go
│   ├── config/                  # Config loading + merge
│   ├── constants/               # All magic strings/values
│   ├── database/                # SQLite / persistence
│   ├── formatter/               # Output formatting (table, CSV, JSON, Markdown)
│   ├── help/                    # Embedded help text files
│   ├── logger/                  # Verbose/debug logging
│   └── models/                  # Data structures
├── scripts/                     # Build, install, deploy scripts
├── config.seed.json             # Default configuration
├── go.mod
└── go.sum
```

**Rules:**
- `cmd/` contains only the entry point — no business logic
- `internal/` is the main code tree — unexported by Go convention
- One file per subcommand in `commands/`
- Constants package holds ALL literal values

---

## Subcommand Architecture

```go
// Dispatch pattern — switch on first arg
func main() {
    if len(os.Args) < 2 {
        help.PrintUsage()
        os.Exit(exitcode.Usage)
    }

    command := os.Args[1]
    args := os.Args[2:]

    switch command {
    case "list":
        commands.RunList(args)
    case "add":
        commands.RunAdd(args)
    case "version":
        commands.RunVersion()
    case "help":
        help.PrintUsage()
    default:
        fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
        os.Exit(exitcode.Usage)
    }
}
```

**Rules:**
- Each `Run*` function accepts `[]string` args
- Each handler parses its own flags via `flag.NewFlagSet`
- No global flag state

---

## Flag Parsing

```go
func RunList(args []string) {
    fs := flag.NewFlagSet("list", flag.ExitOnError)
    verbose := fs.Bool("verbose", false, "Show detailed output")
    format := fs.String("format", "table", "Output format: table|csv|json|markdown")
    fs.Parse(args)

    // Use parsed flags...
}
```

**Rules:**
- Per-command `FlagSet` — never global flags
- Defaults are explicit in the flag definition
- Validate after parsing, fail fast with actionable error

---

## Configuration (Three-Layer Merge)

```
Priority (highest wins):
  1. Command-line flags       (--key=value)
  2. Config file              (~/.config/<app>/config.json)
  3. Seed defaults            (config.seed.json in repo)
```

- Seed file is version-controlled and ships with the binary
- User config is created on first run by copying seed
- Flags override everything — never persisted automatically

---

## Output Formatting

All commands support `--format` with these values:

| Format | Use Case |
|--------|----------|
| `table` | Human-readable terminal (default) |
| `csv` | Spreadsheet / piping |
| `json` | API / programmatic consumption |
| `markdown` | Documentation / reports |

**Rules:**
- Default is always `table`
- JSON output must be valid, parseable JSON (arrays/objects)
- CSV must include a header row
- No color codes in non-table formats

---

## Error Handling

| Exit Code | Meaning |
|-----------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Usage / invalid arguments |
| 3 | Configuration error |
| 4 | Database error |
| 5 | Network error |

**Rules:**
- Every error prints to `stderr`, never `stdout`
- Error messages include what failed AND what to do: `"Config file not found: ~/.config/app/config.json — run '<app> init' to create it"`
- Batch operations collect errors and report at the end, exit with highest severity code

---

## Code Style

| Metric | Limit |
|--------|-------|
| Function length | 8–15 lines (hard max 20) |
| File length | < 300 lines |
| Nesting depth | 0 (early returns only) |
| Boolean naming | `is`/`has`/`should` prefixes, positive only |
| Constants | PascalCase, grouped by category |

---

## Help System

- Help text is embedded in the binary via `//go:embed`
- `<app> help` prints global usage
- `<app> help <command>` prints command-specific help
- `<app> <command> --help` is intercepted and redirects to help system
- Help files live in `internal/help/` as `.txt` files

---

## Database (Local Persistence)

- SQLite with WAL mode enabled
- Singular PascalCase table names: `Repo`, `Config`, `Session`
- PKs: `{TableName}Id INTEGER PRIMARY KEY AUTOINCREMENT`
- Schema versioning via a `SchemaVersion` table
- Upsert pattern: `INSERT ... ON CONFLICT DO UPDATE`

---

## Build & Deploy

```bash

# Build all targets

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X main.Version=$VERSION" -o dist/<app>-linux-amd64 ./cmd/<app>

# Cross-compilation targets

# windows/amd64, windows/arm64, linux/amd64, linux/arm64, darwin/amd64, darwin/arm64

```

**Rules:**
- Static linking (`CGO_ENABLED=0`) for all targets
- Version embedded via `-ldflags` at build time
- Self-update mechanism follows rename-first deploy pattern

---

## Testing

- Test files live next to source: `list_test.go` beside `list.go`
- Table-driven tests for all command handlers
- Test helpers in `internal/testutil/`
- `go test ./...` must pass before any commit

---

## Verbose Logging

```go
// Usage: logger.Verbose("Processing repo: %s", repoName)
// Only prints when --verbose flag is set
func Verbose(format string, args ...interface{}) {
    if !isVerbose { return }
    fmt.Fprintf(os.Stderr, "[VERBOSE] "+format+"\n", args...)
}
```

- Verbose output goes to `stderr`, never `stdout`
- Controlled by `--verbose` / `-v` flag
- Debug level available via `--debug` for even more detail

---

## Progress Tracking

For batch operations with multiple items:

```
Processing repositories...
  [1/5] ✓ repo-alpha          0.3s
  [2/5] ✓ repo-beta           0.5s
  [3/5] ✗ repo-gamma          error: not found
  [4/5] ✓ repo-delta          0.2s
  [5/5] ✓ repo-epsilon        0.4s

Summary: 4 succeeded, 1 failed (1.4s total)
```

---

## Shell Completion

Support tab-completion for:
- Subcommand names
- Flag names (including `--` prefix)
- Flag values where applicable

Generation commands:
```bash
<app> completion bash > /etc/bash_completion.d/<app>
<app> completion zsh > ~/.zfunc/_<app>
<app> completion powershell > <app>.ps1
```

---

## Document Map

| # | File | Topic |
|---|------|-------|
| 00 | 00-overview.md | Philosophy, scope, document index |
| 02 | 02-project-structure.md | Package layout, file organization |
| 03 | 03-subcommand-architecture.md | Routing, dispatch, handler pattern |
| 04 | 04-flag-parsing.md | Per-command flags, defaults, validation |
| 05 | 05-configuration.md | Three-layer config merge |
| 06 | 06-output-formatting.md | Terminal, CSV, JSON, Markdown |
| 07 | 07-error-handling.md | Exit codes, error messages, batch errors |
| 08 | 08-code-style.md | Function length, naming, conditionals |
| 09 | 09-help-system.md | Embedded help, `--help` interception |
| 10 | 10-database.md | SQLite persistence, schema versioning |
| 11 | 11-build-deploy.md | Build scripts, cross-compilation, self-update |
| 12 | 12-testing.md | Test structure, conventions, coverage |
| 13 | 13-checklist.md | Step-by-step implementation checklist |
| 14 | 14-date-formatting.md | Centralized date display format |
| 15 | 15-constants-reference.md | Every constant category with naming |
| 16 | 16-verbose-logging.md | Verbose/debug logging with `--verbose` |
| 17 | 17-progress-tracking.md | Progress reporting for batch operations |
| 18 | 18-batch-execution.md | Exec command for running across repos |
| 19 | 19-shell-completion.md | Tab-completion for PowerShell, Bash, Zsh |
| 20 | 20-terminal-output-design.md | Rich terminal report formatting and color |

---

*Consolidated generic CLI guidelines — v3.2.0 — 2026-04-16*
