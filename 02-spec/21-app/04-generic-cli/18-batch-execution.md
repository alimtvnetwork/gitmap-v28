# Batch Execution

## Purpose

The `exec` command runs arbitrary Git commands across multiple repositories
in a single invocation. It eliminates the need to `cd` into each repo
manually, providing consistent output formatting and a summary of
success/failure across the entire batch.

> **Related specs:**
> - [20-terminal-output-design.md](20-terminal-output-design.md) — full terminal
>   formatting reference (banner, item list, color system, `NO_COLOR` support).
>   Exec banners and per-repo result lines follow the same section patterns.
> - [06-output-formatting.md](06-output-formatting.md) — multi-format output strategy
>   and formatter package structure.
> - [17-progress-tracking.md](17-progress-tracking.md) — progress counter pattern
>   (`[current/total]`) used by batch operations.
> - Rendering pipeline diagram:
>   [`images/terminal-output-pipeline.mmd`](images/terminal-output-pipeline.mmd)

---

## Design Rules

| Rule | Detail |
|------|--------|
| Arbitrary git commands | Any valid `git` subcommand and arguments are accepted |
| Two-tier repo lookup | Database first, JSON file fallback |
| Scope filtering | `--all` for every repo, `--group` for a named subset |
| Per-repo isolation | Each command runs in the repo's directory independently |
| Continue on failure | A failed repo does not abort the batch |
| Captured output | stdout and stderr from each repo are captured and displayed |
| Summary at end | Final line reports succeeded, failed, and missing counts |
| Banner header | Visual header showing the command and repo count |

---

## Command Signature

```
toolname exec [--all | --group <name>] <git-args...>
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--all` | — | `false` | Run across all repos in the database |
| `--group` | `-g` | `""` | Run across repos in a named group |

Positional arguments after flags are passed directly to `git`.

**Examples:**
```
toolname exec status                    # All repos from JSON
toolname exec --all fetch --prune       # All DB repos
toolname exec -g backend log --oneline -5   # Group subset
```

---

## Flag Parsing

```go
func parseExecFlags(args []string) (groupName string, all bool, gitArgs []string) {
    fs := flag.NewFlagSet(constants.CmdExec, flag.ExitOnError)
    gFlag := fs.String("group", "", constants.FlagDescGroup)
    fs.StringVar(gFlag, "g", "", constants.FlagDescGroup)
    aFlag := fs.Bool("all", false, constants.FlagDescAll)
    fs.Parse(args)

    return *gFlag, *aFlag, fs.Args()
}
```

**Key pattern:** After `fs.Parse(args)`, `fs.Args()` returns everything
that was not consumed as a flag — these become the git arguments.

---

## Repo Loading Strategy

A two-tier lookup with scope filtering:

```go
func loadExecByScope(groupName string, all bool) []model.ScanRecord {
    if len(groupName) > 0 {
        return loadRecordsByGroup(groupName)  // DB: group members
    }
    if all {
        return loadAllRecordsDB()             // DB: all repos
    }

    return loadExecRecordsJSON()              // Fallback: JSON file
}
```

| Priority | Source | When |
|----------|--------|------|
| 1 | Database (group) | `--group` flag provided |
| 2 | Database (all) | `--all` flag provided |
| 3 | JSON file | Neither flag — use last scan output |

**Rules:**

- Group not found → error and exit
- JSON file missing → error with actionable message ("run scan first")
- Empty result set → warn and exit cleanly (not an error)

---

## Execution Flow

### Handler

```go
func runExec(args []string) {
    checkHelp("exec", args)
    groupName, all, gitArgs := parseExecFlags(args)
    if len(gitArgs) == 0 {
        fmt.Fprintln(os.Stderr, constants.ErrExecUsage)
        os.Exit(1)
    }

    records := loadExecByScope(groupName, all)
    printExecBanner(gitArgs, len(records))
    succeeded, failed, missing := execAllRepos(records, gitArgs)
    printExecSummary(succeeded, failed, missing, len(records))
}
```

**Flow:** Parse → Validate → Load → Banner → Execute → Summary

### Per-Repo Execution

```go
func execInRepo(rec model.ScanRecord, gitArgs []string) bool {
    cmd := exec.Command(constants.GitBin, gitArgs...)
    cmd.Dir = rec.AbsolutePath
    cmd.Stdout = nil
    cmd.Stderr = nil

    out, err := cmd.CombinedOutput()
    output := strings.TrimSpace(string(out))
    printExecResult(rec.RepoName, output, err)

    return err == nil
}
```

**Rules:**

- `cmd.Dir` sets the working directory — no `cd` needed
- `CombinedOutput()` captures both stdout and stderr
- Trim whitespace from output before display
- Return success/failure as boolean — caller tracks counts

### Missing Repo Handling

```go
func execOneRepo(rec model.ScanRecord, gitArgs []string) (int, int, int) {
    _, err := os.Stat(rec.AbsolutePath)
    if err == nil && execInRepo(rec, gitArgs) {
        return 1, 0, 0    // succeeded
    }
    if err == nil {
        return 0, 1, 0    // failed (command error)
    }

    // Directory does not exist
    fmt.Printf(constants.ExecMissingFmt, ...)
    return 0, 0, 1        // missing
}
```

Three outcomes per repo: **succeeded**, **failed**, or **missing**.

---

## Output Format

### Banner

Displayed before execution begins:

```
  ╭──────────────────────╮
  │   Batch Execution    │
  ╰──────────────────────╯

  Command: status
  Repos:   24
  ─────────────────────────
```

### Per-Repo Results

```
  ✓ repo-name              # Success (green)
      On branch main        # Indented captured output (dim)
  ✗ broken-repo            # Failure (yellow)
      fatal: not a repo     # Error output (dim)
  ⊘ deleted-repo  missing  # Directory not found (dim + yellow)
```

**Rules:**

- Repo names truncated to 22 characters for alignment
- Output lines indented with 6 spaces
- Empty output → no indented lines printed
- Colors: green for success, yellow for failure/missing, dim for output

### Summary

```
  ─────────────────────────
  24 repos · 22 succeeded · 1 failed · 1 missing
```

Only non-zero categories are included in the summary line.

---

## Constants

All exec-related literals in the constants package:

```go
// constants/constants_cli.go
const CmdExec  = "exec"
const AliasExec = "x"

// constants/constants_messages.go
const ErrExecUsage      = "Error: exec requires at least one git argument"
const ErrExecLoadFailed  = "Error: could not load repos: %v\n"

// constants/constants_terminal.go
const ExecBannerTop     = "╭──────────────────────╮"
const ExecBannerTitle   = "│   Batch Execution    │"
const ExecBannerBottom  = "╰──────────────────────╯"
const ExecCommandFmt    = "Command: %s"
const ExecRepoCountFmt  = "Repos:   %d"
const ExecSuccessFmt    = "  %s✓ %-22s%s\n"
const ExecFailFmt       = "  %s✗ %-22s%s\n"
const ExecMissingFmt    = "  %s%-22s %smissing%s\n"
const ExecOutputLineFmt = "      %s%s%s\n"
const ExecSummaryRule   = "─────────────────────────"
const SummaryReposFmt   = "%d repos"
const SummarySucceededFmt = "%s%d succeeded%s"
const SummaryFailedFmt    = "%s%d failed%s"
const SummaryMissingFmt   = "%s%d missing%s"
const SummaryJoinSep      = " · "
```

---

## File Organization

```
cmd/
└── exec.go          Handler, flag parsing, execution, output
```

All exec logic fits in a single file because:

- Flag parsing is a small function
- Execution delegates to `exec.Command`
- Output formatting is straightforward

If the file exceeds 200 lines, extract `execformat.go` for output functions.

---

## Extending the Pattern

To add a new batch command (e.g., `batch-test`):

1. Create `cmd/batchtest.go` with the same structure
2. Reuse `loadExecByScope()` or create a similar scope loader
3. Replace `exec.Command(constants.GitBin, ...)` with the target binary
4. Follow the same three-outcome model (succeeded/failed/missing)
5. Add constants to the appropriate constants file
6. Wire into `dispatch()` in `root.go`

---

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
