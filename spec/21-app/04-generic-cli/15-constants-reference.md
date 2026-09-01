# Constants Reference

> **Related specs:**
> - [02-project-structure.md](02-project-structure.md) — `constants/` package placement
> - [08-code-style.md](08-code-style.md) — naming conventions for constant identifiers
> - [20-terminal-output-design.md](20-terminal-output-design.md) — terminal format constants used in rendering

## Overview

The `constants` package is the single source of truth for all shared
string literals, default values, format strings, and configuration.
No magic strings — every literal used for comparison, formatting,
defaults, or file extensions must be defined here.

The package is split into focused files to maintain the 200-line limit.

---

## File Organization

| File | Responsibility |
|------|----------------|
| `constants.go` | Version, modes, formats, extensions, defaults, permissions |
| `constants_cli.go` | Command names, aliases, subcommand names |
| `constants_terminal.go` | ANSI colors, banner art, table headers, row formats |
| `constants_messages.go` | User-facing messages, error strings, status icons |
| `constants_git.go` | Git binary, subcommands, flags, format strings |
| `constants_store.go` | DB paths, table names, SQL statements |
| `constants_<domain>.go` | Domain-specific constants (one file per feature area) |

### Splitting Rules

- Each file ≤ 200 lines.
- Group by domain, not by type (don't put all strings in one file).
- When a domain file exceeds 150 lines, split into sub-domains.
- File naming: `constants_<domain>.go` (lowercase, underscore).

---

## Category Catalog

### 1. Version & Identity

```go
const Version = "1.0.0"
var RepoPath = ""  // Set at build time via -ldflags
```

| Constant | Naming | Example |
|----------|--------|---------|
| Version string | `Version` | `"2.8.0"` |
| Build-time vars | `RepoPath` | Set via `-ldflags` |
| Tool name | `ToolName` | `"toolname"` |

---

### 2. CLI Command Names & Aliases

```go
const (
    CmdScan       = "scan"
    CmdScanAlias  = "s"
    CmdClone      = "clone"
    CmdCloneAlias = "c"
    CmdHelp       = "help"
    CmdVersion    = "version"
)
```

| Naming Pattern | Convention | Example |
|----------------|-----------|---------|
| Command name | `Cmd<Name>` | `CmdScan = "scan"` |
| Command alias | `Cmd<Name>Alias` | `CmdScanAlias = "s"` |
| Subcommand | `Cmd<Parent><Action>` | `CmdGroupCreate = "create"` |

---

### 3. Modes & Output Formats

```go
const (
    ModeHTTPS      = "https"
    ModeSSH        = "ssh"
    OutputTerminal = "terminal"
    OutputCSV      = "csv"
    OutputJSON     = "json"
)
```

| Naming Pattern | Convention | Example |
|----------------|-----------|---------|
| Mode values | `Mode<Name>` | `ModeHTTPS` |
| Output formats | `Output<Name>` | `OutputJSON` |
| URL prefixes | `Prefix<Name>` | `PrefixHTTPS = "https://"` |

---

### 4. File Extensions & Default File Names

```go
const (
    ExtCSV  = ".csv"
    ExtJSON = ".json"
    ExtTXT  = ".txt"
    ExtGit  = ".git"
)

const (
    DefaultCSVFile  = "toolname.csv"
    DefaultJSONFile = "toolname.json"
)
```

| Naming Pattern | Convention | Example |
|----------------|-----------|---------|
| File extensions | `Ext<Type>` | `ExtJSON = ".json"` |
| Default filenames | `Default<Type>File` | `DefaultCSVFile` |
| Default paths | `Default<Purpose>Path` | `DefaultConfigPath` |

---

### 5. Default Values & Paths

```go
const (
    DefaultConfigPath = "./data/config.json"
    DefaultOutputDir  = "./toolname-output"
    DefaultBranch     = "main"
    JSONIndent        = "  "
)
```

| Naming Pattern | Convention | Example |
|----------------|-----------|---------|
| Default values | `Default<Name>` | `DefaultBranch` |
| Directory names | `<Purpose>Dir` | `OutputDir` |
| Permissions | `Perm<Type>` | `PermDir = 0o755` |

---

### 6. ANSI Color Codes

```go
const (
    ColorReset  = "\033[0m"
    ColorGreen  = "\033[32m"
    ColorRed    = "\033[31m"
    ColorYellow = "\033[33m"
    ColorCyan   = "\033[36m"
    ColorWhite  = "\033[97m"
    ColorDim    = "\033[90m"
)
```

| Naming Pattern | Convention | Example |
|----------------|-----------|---------|
| Color codes | `Color<Name>` | `ColorGreen` |
| Reset | `ColorReset` | Always first in group |

---

### 7. Terminal UI — Banners & Box Drawing

```go
const (
    BannerTop    = "╔══════════════════════════════════════╗"
    BannerTitle  = "║         toolname v1.0.0              ║"
    BannerBottom = "╚══════════════════════════════════════╝"
)
```

| Naming Pattern | Convention | Example |
|----------------|-----------|---------|
| Banner parts | `<Section>Banner<Part>` | `StatusBannerTop` |
| Section headers | `<Section>Header` | `ScanHeader` |
| Tree characters | `Tree<Type>` | `TreeBranch = "├──"` |

---

### 8. Terminal UI — Format Strings & Table Headers

```go
const (
    StatusRowFmt    = "  %-22s %s  %s  %s\n"
    StatusHeaderFmt = "  %-22s %-12s %-8s\n"
    SummaryFmt      = "\n  %d items processed\n"
)
```

| Naming Pattern | Convention | Example |
|----------------|-----------|---------|
| Row formats | `<Section>RowFmt` | `StatusRowFmt` |
| Header formats | `<Section>HeaderFmt` | `StatusHeaderFmt` |
| Summary formats | `<Section>SummaryFmt` | `ScanSummaryFmt` |
| Count formats | `<Section>CountFmt` | `RepoCountFmt` |

Table column headers use a `var` slice:

```go
var StatusTableColumns = []string{"REPO", "STATUS", "BRANCH"}
```

---

### 9. Status Icons & Indicators

```go
const (
    StatusIconClean   = "✓ clean"
    StatusIconDirty   = "● dirty"
    StatusDash        = "—"
    StatusSyncUpFmt   = "↑%d"
    StatusSyncDownFmt = "↓%d"
)
```

| Naming Pattern | Convention | Example |
|----------------|-----------|---------|
| Icons | `StatusIcon<State>` | `StatusIconClean` |
| Indicator formats | `Status<Type>Fmt` | `StatusSyncUpFmt` |

---

### 10. Error Messages

```go
const (
    ErrSourceRequired = "Error: source file is required"
    ErrConfigLoad     = "Error: could not load config from %s"
    ErrRepoNotFound   = "Error: no repo matches slug '%s'"
    ErrGenericFmt     = "Error: %v\n"
)
```

| Naming Pattern | Convention | Example |
|----------------|-----------|---------|
| Static errors | `Err<What>` | `ErrSourceRequired` |
| Format errors | `Err<What>Fmt` | `ErrConfigLoad` (with `%s`) |
| Generic errors | `ErrGenericFmt` | Catch-all `%v` format |

**Rule:** Error messages must be actionable — tell the user what
to do, not just what failed.

---

### 11. User-Facing Messages

```go
const (
    MsgScanComplete    = "✓ Scan complete"
    MsgDesktopAdded    = "  ✓ Added to GitHub Desktop: %s\n"
    MsgDesktopFailed   = "  ✗ Failed to add %s: %v\n"
    MsgDesktopSummary  = "GitHub Desktop: %d added, %d failed\n"
)
```

| Naming Pattern | Convention | Example |
|----------------|-----------|---------|
| Info messages | `Msg<Action>` | `MsgScanComplete` |
| Format messages | `Msg<Action>Fmt` | `MsgDesktopAdded` (with `%s`) |
| Warning messages | `Msg<Topic>Warning` | `MsgFetchWarning` |

---

### 12. Git Commands & Arguments

```go
const (
    GitBin        = "git"
    GitClone      = "clone"
    GitPull       = "pull"
    GitFetch      = "fetch"
    GitBranchFlag = "-b"
    GitFFOnlyFlag = "--ff-only"
    GitOrigin     = "origin"
)
```

| Naming Pattern | Convention | Example |
|----------------|-----------|---------|
| Binary name | `GitBin` | `"git"` |
| Git subcommands | `Git<Command>` | `GitClone`, `GitPull` |
| Git flags | `Git<Flag>Flag` | `GitBranchFlag = "-b"` |
| Git refs | `Git<Ref>` | `GitHEAD = "HEAD"` |
| Format strings | `Git<Purpose>Format` | `GitLogTipFormat` |

---

### 13. Database — Paths, Tables, SQL

```go
const (
    DBDir  = "data"
    DBFile = "toolname.db"
)

const (
    TableRepos  = "Repos"
    TableGroups = "Groups"
)

const SQLCreateRepos = `CREATE TABLE IF NOT EXISTS Repos (...)`
const SQLUpsertRepo  = `INSERT INTO Repos (...) ON CONFLICT(...) DO UPDATE SET ...`
```

| Naming Pattern | Convention | Example |
|----------------|-----------|---------|
| DB paths | `DB<Part>` | `DBDir`, `DBFile` |
| Table names | `Table<Name>` | `TableRepos` |
| CREATE statements | `SQLCreate<Table>` | `SQLCreateRepos` |
| UPSERT statements | `SQLUpsert<Table>` | `SQLUpsertRepo` |
| SELECT statements | `SQLSelect<What>` | `SQLSelectAllRepos` |
| DELETE statements | `SQLDelete<What>` | `SQLDeleteGroup` |
| Index creation | `SQLIndex<Table><Col>` | `SQLIndexReposPath` |

---

### 14. Flag Names & Help Descriptions

```go
const (
    FlagConfig   = "config"
    FlagMode     = "mode"
    FlagOutput   = "output"
    FlagDryRun   = "dry-run"
    FlagVerbose  = "verbose"
    FlagHelp     = "--help"
    FlagHelpShort = "-h"
)

const (
    HelpFlagConfig  = "Config file path"
    HelpFlagMode    = "Clone URL style (ssh|https)"
    HelpFlagOutput  = "Output format (csv|json|terminal)"
)
```

| Naming Pattern | Convention | Example |
|----------------|-----------|---------|
| Flag names | `Flag<Name>` | `FlagDryRun = "dry-run"` |
| Flag descriptions | `HelpFlag<Name>` | `HelpFlagMode` |
| Flag defaults | `Default<Flag>` | `DefaultMode = "https"` |

---

### 15. Date & Time Formatting

```go
const (
    DateDisplayLayout = "02-Jan-2006 03:04 PM"
    DateUTCSuffix     = " (UTC)"
)
```

| Naming Pattern | Convention | Example |
|----------------|-----------|---------|
| Layout strings | `Date<Purpose>Layout` | `DateDisplayLayout` |
| Suffixes | `Date<Purpose>Suffix` | `DateUTCSuffix` |

---

### 16. OS & Platform

```go
const (
    OSWindows = "windows"
    OSDarwin  = "darwin"
    CmdExplorer = "explorer"
    CmdOpen     = "open"
    CmdXdgOpen  = "xdg-open"
)
```

| Naming Pattern | Convention | Example |
|----------------|-----------|---------|
| OS identifiers | `OS<Name>` | `OSWindows` |
| Platform commands | `Cmd<Name>` | `CmdExplorer` |

---

## What Does NOT Belong in Constants

| Category | Reason |
|----------|--------|
| Struct definitions | Belong in `model` package |
| Business logic | Belongs in domain packages |
| Template content | Use `go:embed` in `formatter/templates/` |
| Test data strings | Stay local in test files |
| Log messages unique to one location | Not compared or reused |

---

## Naming Quick Reference

| Prefix | Category | Example |
|--------|----------|---------|
| `Cmd` | CLI command names | `CmdScan` |
| `Mode` | Operation modes | `ModeHTTPS` |
| `Output` | Output formats | `OutputJSON` |
| `Ext` | File extensions | `ExtCSV` |
| `Default` | Default values | `DefaultBranch` |
| `Color` | ANSI codes | `ColorGreen` |
| `Err` | Error messages | `ErrSourceRequired` |
| `Msg` | User messages | `MsgScanComplete` |
| `Git` | Git commands/flags | `GitClone` |
| `SQL` | SQL statements | `SQLCreateRepos` |
| `Table` | Table names | `TableRepos` |
| `DB` | Database paths | `DBFile` |
| `Flag` | Flag names | `FlagVerbose` |
| `Help` | Help descriptions | `HelpFlagMode` |
| `Status` | UI indicators | `StatusIconClean` |
| `Perm` | Permissions | `PermDir` |
| `Prefix` | URL/string prefixes | `PrefixHTTPS` |
| `Date` | Date formatting | `DateDisplayLayout` |
| `OS` | Platform identifiers | `OSWindows` |
| `Tree` | Tree-drawing chars | `TreeBranch` |

---

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
