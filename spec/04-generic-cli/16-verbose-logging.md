# Verbose Logging

> **Related specs:**
> - [04-flag-parsing.md](04-flag-parsing.md) — `--verbose` flag parsing
> - [07-error-handling.md](07-error-handling.md) — verbose output for error diagnosis
> - [15-constants-reference.md](15-constants-reference.md) — logging format constants

## Purpose

A shared `--verbose` flag enables detailed debug logging to a timestamped
file. Normal runs produce clean user-facing output only. Verbose runs
capture full diagnostics for troubleshooting without polluting stdout.

---

## Design Rules

| Rule | Detail |
|------|--------|
| Off by default | No log file created unless `--verbose` is passed |
| File + stderr | Every verbose entry writes to both the log file and stderr |
| Timestamped entries | Each line prefixed with `[HH:MM:SS.mmm]` |
| Timestamped filenames | Log file named `toolname-verbose-YYYY-MM-DD_HH-mm-ss.log` |
| Output directory | Logs written to the tool's default output folder |
| Dim on stderr | Verbose stderr output uses dim/gray ANSI color |
| No stdout pollution | Verbose output never mixes with normal command output |
| Global singleton | One logger instance shared across all packages |

---

## Package Structure

```
verbose/
└── verbose.go     Logger type, Init, Close, Log, IsEnabled, Get
```

Single file. No sub-packages. No external dependencies beyond `constants`.

---

## Logger API

```go
// Init creates the log file and enables verbose logging.
// Call once at startup when --verbose is set.
func Init() (*Logger, error)

// Close flushes and closes the log file.
func (l *Logger) Close()

// Log writes a formatted message to the log file and stderr.
func (l *Logger) Log(format string, args ...interface{})

// IsEnabled returns true if verbose mode is active.
func IsEnabled() bool

// Get returns the global logger (may be nil).
func Get() *Logger
```

---

## Logger Type

```go
type Logger struct {
    file    *os.File
    enabled bool
}

var global *Logger
```

- `file` — open handle to the log file
- `enabled` — guards all write operations
- `global` — package-level singleton set by `Init()`

---

## Init Flow

```go
func Init() (*Logger, error) {
    logDir := constants.DefaultOutputFolder
    _ = os.MkdirAll(logDir, constants.DirPermission)

    timestamp := time.Now().Format("2006-01-02_15-04-05")
    logPath := filepath.Join(logDir, fmt.Sprintf(constants.VerboseLogFileFmt, timestamp))

    file, err := os.Create(logPath)
    if err != nil {
        return nil, err
    }

    l := &Logger{file: file, enabled: true}
    global = l
    fmt.Printf(constants.MsgVerboseLogFile, logPath)

    return l, nil
}
```

**Key points:**

- Creates the output directory if missing (no error on existing)
- Prints the log file path to stdout so the user knows where to find it
- Returns both the logger and any error — caller decides whether to abort

---

## Log Entry Format

```go
func writeLogEntry(l *Logger, format string, args ...interface{}) {
    line := fmt.Sprintf(format, args...)
    ts := time.Now().Format("15:04:05.000")
    entry := fmt.Sprintf("[%s] %s\n", ts, line)
    l.file.WriteString(entry)
    fmt.Fprint(os.Stderr, constants.ColorDim+entry+constants.ColorReset)
}
```

**Example output:**

```
[14:32:07.123] git clone https://github.com/user/repo.git
[14:32:09.456] clone completed in 2.3s
[14:32:09.460] retry attempt 1/4 for locked file
```

---

## Flag Registration

The `--verbose` flag is a **global flag** registered on the root command,
not per-subcommand.

```go
// In cmd/rootflags.go
fs.BoolVar(&verboseFlag, constants.FlagVerbose, false, constants.FlagDescVerbose)
```

---

## Command Handler Pattern

Every command that supports verbose logging follows this pattern:

```go
func runPull(args []string) {
    checkHelp("pull", args)
    slug, group, all, verboseFlag := parsePullFlags(args)

    if verboseFlag {
        initVerboseLog()       // Init + defer Close
    }

    // ... command logic ...
}

func initVerboseLog() {
    logger, err := verbose.Init()
    if err != nil {
        fmt.Fprintf(os.Stderr, constants.ErrVerboseInit, err)
        return                 // Non-fatal — continue without logging
    }
    defer logger.Close()
}
```

**Rules:**

- Verbose init failure is **non-fatal** — warn and continue
- `defer Close()` in the same function that calls `Init()`
- Never pass the logger as a parameter — use `verbose.Get()` or `verbose.IsEnabled()`

---

## What to Log

| Category | Examples |
|----------|----------|
| Git operations | Clone/pull commands, remote URLs, branch names |
| Retry attempts | Attempt number, delay, reason for retry |
| File I/O | Paths read/written, file sizes, permissions |
| External processes | Command lines, exit codes, stdout/stderr |
| Timing | Operation durations, elapsed time |
| Environment | OS, paths, config values loaded |
| Errors (detailed) | Full error chains, stack context |
| Compression | Archive size in bytes, SHA-1 hash per archive |
| Checksums | Per-file SHA-256 hash during checksum generation |
| Asset uploads | Target repo/tag, per-asset file size, HTTP status |

**What NOT to log:**

- Secrets, tokens, or credentials
- Routine success paths that add no diagnostic value
- Data that duplicates normal stdout output

---

## Release Pipeline Log Points

The release workflow emits verbose log entries at each stage.
All entries follow the `prefix: detail` convention.

### Stage Summary

| # | Stage | Prefix | Source File |
|---|-------|--------|-------------|
| 1 | [Version Resolution](#version-resolution-workflowgo) | `version:` | `workflow.go` |
| 2 | [Source Resolution](#source-resolution-gitopsgo) | `source:` | `gitops.go` |
| 3 | [Git Operations](#git-operations-gitopsgo) | `git:` | `gitops.go` |
| 4 | [Asset Collection](#asset-collection-githubgo) | `assets:` | `github.go` |
| 5 | [Staging Directory](#staging-directory-assetsgo) | `staging:` | `assets.go` |
| 6 | [Cross-Compilation](#cross-compilation-assetsgo) | `build:` | `assets.go` |
| 7 | [Compression](#compression-compressgo) | `compress:` | `compress.go` |
| 8 | [Checksums](#checksums-checksumsgo) | `checksum:` | `checksums.go` |
| 9 | [Zip Group Processing](#zip-group-processing-workflowfinalizego) | `zip-group:` | `workflowfinalize.go` |
| 10 | [Ad-Hoc Zip Archives](#ad-hoc-zip-archives-workflowfinalizego) | `ad-hoc-zip:` | `workflowfinalize.go` |
| 11 | [Zip Group Archives](#zip-group-archives-ziparchivego) | `zip-summary:` | `ziparchive.go` |
| 12 | [GitHub Upload](#github-upload-workflowfinalizego-assetsuploadgo) | `github:` / `upload:` | `workflowfinalize.go`, `assetsupload.go` |
| 13 | [Retry](#retry-retrygo) | `retry:` | `retry.go` |
| 14 | [Metadata Persistence](#metadata-persistence-workflowfinalizego) | `metadata:` | `workflowfinalize.go` |
| 15 | [Rollback](#rollback-rollbackgo) | `rollback:` | `rollback.go` |
| 16 | [Autocommit](#autocommit-autocommitgo) | `autocommit:` | `autocommit.go` |

### Version Resolution (`workflow.go`)

Logged when the release version is determined from CLI, bump, or file:

```
version: resolved from CLI argument: v2.5.0
version: current baseline: v2.4.0
version: baseline from latest.json: v2.4.0
version: latest.json unavailable, falling back to git tags
version: resolved via --bump minor: v2.5.0
version: resolved from version.json: v2.5.0
```

### Source Resolution (`gitops.go`)

Logged when the release source ref is determined from `--commit`, `--branch`, or HEAD:

```
source: using commit a1b2c3d4e5f6
source: using branch feature-x (origin/feature-x)
source: using HEAD on branch main
source: using detached HEAD
```

### Git Operations (`gitops.go`)

Logged when branches, tags are created and pushed:

```
git: creating branch release/v2.5.0 from HEAD
git: creating tag v2.5.0
git: pushing branch release/v2.5.0 to origin
git: pushing tag v2.5.0 to origin
```

### Asset Collection (`github.go`)

Logged when user-provided assets are resolved from `--assets`:

```
assets: collected 3 file(s) from directory dist/
assets: gitmap_v2.5.0_linux_amd64
assets: gitmap_v2.5.0_darwin_arm64
assets: gitmap_v2.5.0_windows_amd64.exe
assets: single file build/output.tar.gz
assets: path not found: missing/dir
```

### Staging Directory (`assets.go`)

Logged when the release-assets staging directory is created or removed:

```
staging: created directory assets/staging
staging: removing directory assets/staging
```

### Cross-Compilation (`assets.go`)

Logged before and after each GOOS/GOARCH build:

```
build: linux/amd64 → assets/staging/gitmap_v2.5.0_linux_amd64
build: linux/amd64 complete (4821504 bytes)
build: windows/arm64 failed: unsupported GOARCH
```

### Compression (`compress.go`)

Logged after each asset is compressed into `.zip` or `.tar.gz`:

```
compress: gitmap_v2.5.0_linux_amd64.tar.gz — 4821504 bytes, sha1:a3f9c0...
```

### Checksums (`checksums.go`)

Logged as each file's SHA-256 hash is computed for `checksums.txt`:

```
checksum: gitmap_v2.5.0_linux_amd64.tar.gz  sha256:e3b0c44298fc...
```

### Zip Group Processing (`workflowfinalize.go`)

Logged when persistent zip groups are resolved and built:

```
zip-group: processing group "chrome-extension-v2"
zip-group: 2 group(s) produced 2 archive(s)
```

### Ad-Hoc Zip Archives (`workflowfinalize.go`)

Logged when ad-hoc `-Z` items are bundled:

```
ad-hoc-zip: 3 item(s), bundle=my-bundle
ad-hoc-zip: item src/config.json
ad-hoc-zip: item assets/logo.png
ad-hoc-zip: item docs/
ad-hoc-zip: produced 1 archive(s)
```

### Zip Group Archives (`ziparchive.go`)

Logged after each zip group archive is created:

```
zip-summary: chrome-extension.zip — 12 files, 238471 bytes, sha1:7b2a1f...
```

### GitHub Upload (`workflowfinalize.go`, `assetsupload.go`)

Logged at release creation and per-asset upload:

```
github: creating release v2.5.0 on owner/repo (6 asset(s))
github: release created, id=12345
upload-start: gitmap_v2.5.0_linux_amd64.tar.gz (4821504 bytes)
upload: gitmap_v2.5.0_linux_amd64.tar.gz → HTTP 201
```

### Retry (`retry.go`)

Logged on each failed attempt and before backoff sleep:

```
retry: gitmap_v2.5.0_linux_amd64.tar.gz attempt 1/3 failed: upload error 502: Bad Gateway
retry: gitmap_v2.5.0_linux_amd64.tar.gz sleeping 2s before attempt 2
```

### Metadata Persistence (`workflowfinalize.go`)

Logged when release JSON and latest.json are written after a successful release:

```
metadata: writing .release/v2.5.0.json
metadata: updating latest.json to v2.5.0
metadata: skipping latest.json (pre-release v2.5.0-rc.1)
```

### Rollback (`rollback.go`)

Logged when the release workflow encounters a failure and rolls back
local branches and tags:

```
rollback: starting (branch=release/v2.5.0, tag=v2.5.0, return-to=main)
rollback: switching back to main
rollback: deleting local branch release/v2.5.0
rollback: deleting local tag v2.5.0
```

### Autocommit (`autocommit.go`)

Logged during the post-release auto-commit and push of metadata files:

```
autocommit: starting for v2.5.0 (dry-run=false)
autocommit: 2 release file(s), 0 other file(s)
autocommit: staged 2 file(s)
autocommit: committed "release v2.5.0 metadata"
autocommit: pushed to main
```

---

## Constants

All verbose-related literals live in the constants package:

```go
// constants/constants.go
const VerboseLogFileFmt = "toolname-verbose-%s.log"

// constants/constants_cli.go
const FlagVerbose    = "verbose"
const FlagDescVerbose = "Enable verbose debug logging to file"

// constants/constants_messages.go
const MsgVerboseLogFile = "Verbose log: %s\n"
const ErrVerboseInit    = "Warning: could not initialize verbose log: %v\n"
```

---

## Conditional Logging in Libraries

Domain packages (scanner, cloner, mapper) check `verbose.IsEnabled()`
before calling `verbose.Get().Log(...)`:

```go
func safePullOne(repo model.Record) error {
    logger := verbose.Get()

    if logger != nil {
        logger.Log("pulling %s at %s", repo.Name, repo.Path)
    }

    // ... pull logic ...

    if logger != nil {
        logger.Log("pull complete for %s (%.1fs)", repo.Name, elapsed.Seconds())
    }

    return nil
}
```

**Rules:**

- Always nil-check `verbose.Get()` — verbose may not be active
- Keep log calls outside hot loops to avoid performance overhead
- Use `fmt.Sprintf`-style formatting — no structured logging libraries

---

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
