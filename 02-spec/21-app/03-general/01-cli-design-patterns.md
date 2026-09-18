# CLI Design Patterns

## Overview

This document describes reusable patterns for building Go CLI tools
with clear subcommand routing, consistent flag parsing, and
user-friendly output.

## Subcommand Architecture

### Entry Point

A single `Run()` function in the `cmd` package validates arguments
and dispatches to subcommand handlers.

```
func Run() {
    if len(os.Args) < 2 {
        printUsage()
        os.Exit(1)
    }
    dispatch(os.Args[1])
}
```

### Dispatch Pattern

Use a `switch` on the first argument to route to handler functions.
Each handler receives `os.Args[2:]` as its arguments.

```
func dispatch(command string) {
    switch command {
    case constants.CmdScan:
        runScan(os.Args[2:])
    case constants.CmdClone:
        runClone(os.Args[2:])
    case constants.CmdVersion:
        fmt.Println(constants.Version)
    case constants.CmdHelp:
        printUsage()
    default:
        fmt.Printf("Unknown command: %s\n", command)
        os.Exit(1)
    }
}
```

### Rules

| Rule | Rationale |
|------|-----------|
| Each subcommand lives in its own file (`scan.go`, `clone.go`, etc.) | Single responsibility per file |
| Handler functions are unexported (`runScan`, not `RunScan`) | Only `Run()` is the public API |
| Unknown commands print a message and exit with code 1 | Fail fast, fail clearly |

## Flag Parsing

### Approach

Use Go's `flag.NewFlagSet` per subcommand to avoid global flag pollution.

```go
func parseScanFlags(args []string) (dir, configPath, mode string) {
    fs := flag.NewFlagSet("scan", flag.ExitOnError)
    fs.StringVar(&configPath, "config", constants.DefaultConfigPath, "Config file")
    fs.StringVar(&mode, "mode", constants.ModeHTTPS, "Clone URL style")
    fs.Parse(args)

    if fs.NArg() > 0 {
        dir = fs.Arg(0)
    }
    return
}
```

### Flag Naming Conventions

| Pattern | Example | Why |
|---------|---------|-----|
| Lowercase with hyphens | `--target-dir` | Readable, standard |
| Boolean flags as switches | `--dry-run` | No value needed |
| Positional args for primary input | `tool scan <dir>` | Natural CLI UX |

### Defaults

All defaults live in the `constants` package. Never inline a default
string in flag definitions — reference a constant.

## Version Command

### Pattern

- Version is a constant (`constants.Version`) following SemVer.
- A dedicated `version` subcommand prints the version string and exits.
- The version is also shown in the help output and terminal banner.
- The build system prints the version after each successful build.

### Build-Time Variables

Use `-ldflags` to embed values at compile time:

```
go build -ldflags "-X 'pkg/constants.RepoPath=$path'" -o binary .
```

This enables features (like self-update) that need to know
the source repo location without hardcoding paths.

## Help Output

### Structure

```
Usage: toolname <command> [flags]

Commands:
  scan [dir]          Scan directory for items
  clone <source>      Re-create from CSV/JSON/text file
  update              Self-update from source repo
  version             Show version number
  help                Show this help message

Scan flags:
  --config <path>     Config file (default: ./data/config.json)
  --mode ssh|https    URL style (default: https)
```

### Rules

- All help text lives in `constants` as named strings.
- Commands are left-aligned with consistent indentation.
- Flag descriptions include type hints and defaults.
- `help` is always listed last in the command list.

## Error Handling

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | User error (bad args, missing file) |
| Non-zero | Propagated from child processes |

### Error Messages

- All error format strings live in `constants` (`ErrSourceRequired`, `ErrConfigLoad`, etc.).
- Errors print to stderr, not stdout.
- After printing an error, exit immediately — don't continue.
- For batch operations (e.g., process N items), log per-item failures
  but continue; print summary at the end.

## Constants Package

### What Goes In Constants

| Category | Examples |
|----------|---------|
| Version string | `Version = "1.0.0"` |
| CLI command names | `CmdScan`, `CmdClone`, `CmdHelp` |
| Default values | `DefaultConfigPath`, `DefaultBranch` |
| Format strings | `InstructionFmt`, `BannerTitle` |
| Error messages | `ErrSourceRequired`, `ErrConfigLoad` |
| ANSI color codes | `ColorGreen`, `ColorReset` |
| File extensions | `ExtCSV`, `ExtJSON` |

### What Does NOT Go In Constants

- Struct definitions (belong in `model` package).
- Business logic.
- Template content (use `go:embed` for templates).

## Package Structure

```
cmd/         CLI entry point, subcommand routing, flag parsing
config/      JSON config loading + flag merging
constants/   All shared string literals and defaults
scanner/     Directory walking + data detection
mapper/      Raw data → output records
formatter/   Rendering to terminal/CSV/JSON/scripts
model/       Shared data structures
```

### Rules

- One responsibility per package.
- No circular imports.
- The `cmd` package orchestrates; it calls into other packages
  but other packages never import `cmd`.
