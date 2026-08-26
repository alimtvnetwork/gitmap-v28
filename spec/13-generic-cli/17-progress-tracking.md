# Progress Tracking

## Purpose

Batch operations (clone, pull, exec) process many items sequentially.
A progress tracker provides real-time visual feedback so the user knows
what is happening, how far along the operation is, and how it concluded.

> **Related specs:**
> - [20-terminal-output-design.md](20-terminal-output-design.md) — full terminal
>   formatting reference (banner, item list, tree, color system, `NO_COLOR` support).
>   Progress counters follow the same `[current/total]` pattern used in item lists.
> - [06-output-formatting.md](06-output-formatting.md) — multi-format output strategy
>   and formatter package structure.
> - Rendering pipeline diagram:
>   [`images/terminal-output-pipeline.mmd`](images/terminal-output-pipeline.mmd)

---

## Design Rules

| Rule | Detail |
|------|--------|
| Counter format | `[current/total]` prefix on every item line |
| Repo name visible | Current repository name printed alongside the counter |
| Elapsed time | Duration shown on completion of each item and in the summary |
| Summary at end | Final line reports total items, elapsed time, success/failure counts |
| Quiet mode | `--quiet` flag suppresses all progress output for programmatic use |
| Stderr only | All progress output goes to stderr — stdout is reserved for data |
| No progress bars | Counter + name is sufficient; avoid complex progress bar libraries |

---

## Package Placement

Progress tracking lives in the **domain package** that owns the batch
operation (e.g., `cloner/progress.go` for clone/pull operations).

```
cloner/
├── cloner.go        Core clone/pull logic
├── progress.go      Progress tracker type
├── safe_pull.go     Retry-aware pull
└── pulldiag.go      Pull diagnostics
```

If multiple packages need progress tracking, extract a shared
`progress/` package. Do not duplicate the pattern.

---

## Progress Type

```go
type Progress struct {
    total   int
    current int
    start   time.Time
    quiet   bool
    cloned  int
    pulled  int
    failed  int
}
```

| Field | Purpose |
|-------|---------|
| `total` | Total number of items to process |
| `current` | Running counter, incremented on each `Begin()` call |
| `start` | Timestamp captured at construction for elapsed time |
| `quiet` | When true, all output methods become no-ops |
| `cloned` | Success counter for new clones |
| `pulled` | Success counter for updates (pulls) |
| `failed` | Failure counter |

---

## Constructor

```go
func NewProgress(total int, quiet bool) *Progress {
    return &Progress{
        total: total,
        start: time.Now(),
        quiet: quiet,
    }
}
```

**Rules:**

- `time.Now()` captured once at construction — not per item
- All counters start at zero (Go zero values)
- Caller determines `quiet` from the `--quiet` CLI flag

---

## Lifecycle Methods

### Begin

Called before processing each item. Increments the counter and prints
the item header.

```go
func (p *Progress) Begin(name string) {
    p.current++
    if p.quiet {
        return
    }

    fmt.Fprintf(os.Stderr, constants.ProgressBeginFmt, p.current, p.total, name)
}
```

**Output:** `[1/24] repo-name `

### Done

Called after an item succeeds. Records the result category and prints
elapsed time.

```go
func (p *Progress) Done(result model.CloneResult, pulled bool) {
    if pulled {
        p.pulled++
    } else {
        p.cloned++
    }

    if p.quiet {
        return
    }

    elapsed := time.Since(p.start)
    fmt.Fprintf(os.Stderr, constants.ProgressDoneFmt, formatDuration(elapsed))
}
```

**Output:** `✓ (2.3s)`

### Fail

Called after an item fails. Increments the failure counter.

```go
func (p *Progress) Fail(result model.CloneResult) {
    p.failed++
    if p.quiet {
        return
    }

    fmt.Fprintf(os.Stderr, constants.ProgressFailFmt)
}
```

**Output:** `✗ failed`

### PrintSummary

Called once after all items are processed. Prints totals and breakdown.

```go
func (p *Progress) PrintSummary() {
    if p.quiet {
        return
    }

    elapsed := time.Since(p.start)
    fmt.Fprintf(os.Stderr, constants.ProgressSummaryFmt,
        p.current, p.total, formatDuration(elapsed))
    fmt.Fprintf(os.Stderr, constants.ProgressDetailFmt,
        p.cloned, p.pulled, p.failed)
}
```

**Output:**
```
Done: 24/24 (1m 12s)
  Cloned: 18  Pulled: 5  Failed: 1
```

---

## Duration Formatting

A single helper function formats durations for human readability:

```go
func formatDuration(d time.Duration) string {
    if d < time.Minute {
        return fmt.Sprintf("%.1fs", d.Seconds())
    }

    mins := int(d.Minutes())
    secs := int(d.Seconds()) % 60

    return fmt.Sprintf("%dm %ds", mins, secs)
}
```

| Duration | Output |
|----------|--------|
| 2.3 seconds | `2.3s` |
| 72 seconds | `1m 12s` |
| 5 minutes 3 seconds | `5m 3s` |

**Rules:**

- Under one minute → show decimal seconds (`%.1fs`)
- One minute or more → show minutes and seconds (`Xm Ys`)
- No hours format needed — batch operations should not run that long
- Function is unexported — internal to the progress package

---

## Usage Pattern

```go
func cloneAll(repos []model.Record, quiet bool) {
    progress := NewProgress(len(repos), quiet)

    for _, repo := range repos {
        progress.Begin(repo.Name)

        result, err := cloneOne(repo)
        if err != nil {
            progress.Fail(result)
            continue
        }

        progress.Done(result, repo.Exists)
    }

    progress.PrintSummary()
}
```

**Rules:**

- Create progress tracker **before** the loop
- Call `Begin()` **first** in each iteration
- Call exactly one of `Done()` or `Fail()` per item — never both
- Call `PrintSummary()` **after** the loop, unconditionally
- Use `continue` after `Fail()` — do not abort the batch

---

## Quiet Mode

The `--quiet` flag suppresses all progress output. This is useful when:

- Output is piped to another tool
- Running in CI/CD where progress noise is unwanted
- Embedding the CLI in scripts that parse stdout

```go
// In cmd/clone.go
fs.BoolVar(&quietFlag, constants.FlagQuiet, false, constants.FlagDescQuiet)

// Pass to progress tracker
progress := cloner.NewProgress(len(repos), quietFlag)
```

When quiet mode is active:

- All `Begin`, `Done`, `Fail`, and `PrintSummary` calls are no-ops
- Internal counters still increment (for programmatic access if needed)
- Stdout data output (JSON, CSV) is unaffected

---

## Constants

All progress format strings live in the constants package:

```go
// constants/constants_clone.go
const ProgressBeginFmt   = "[%d/%d] %s "
const ProgressDoneFmt    = "✓ (%s)\n"
const ProgressFailFmt    = "✗ failed\n"
const ProgressSummaryFmt = "Done: %d/%d (%s)\n"
const ProgressDetailFmt  = "  Cloned: %d  Pulled: %d  Failed: %d\n"

// constants/constants_cli.go
const FlagQuiet     = "quiet"
const FlagDescQuiet = "Suppress progress output"
```

---

## Extending for Other Commands

When adding progress tracking to a new batch command (e.g., `exec`):

1. Adjust the `Progress` struct fields to match the operation's outcomes
2. Keep the same lifecycle: `NewProgress` → `Begin` → `Done`/`Fail` → `PrintSummary`
3. Use the same counter format `[current/total]`
4. Add format constants to the command's constants file
5. Respect the `--quiet` flag

---

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
