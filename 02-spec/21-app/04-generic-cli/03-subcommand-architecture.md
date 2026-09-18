# Subcommand Architecture

> **Related specs:**
> - [04-flag-parsing.md](04-flag-parsing.md) — per-command flag parsing within dispatched handlers
> - [09-help-system.md](09-help-system.md) — `--help` integration for each subcommand
> - [19-shell-completion.md](19-shell-completion.md) — tab-completion driven by the subcommand registry

## Entry Point

A single exported `Run()` function validates arguments and dispatches:

```go
func Run() {
    if len(os.Args) < 2 {
        printUsage()
        os.Exit(1)
    }
    dispatch(os.Args[1])
}
```

## Dispatch Pattern

Use a `switch` on the first argument. Each handler receives
`os.Args[2:]` as its arguments.

```go
func dispatch(command string) {
    switch command {
    case constants.CmdScan, constants.AliasScan:
        runScan(os.Args[2:])
    case constants.CmdClone, constants.AliasClone:
        runClone(os.Args[2:])
    case constants.CmdVersion:
        fmt.Println(constants.Version)
    case constants.CmdHelp:
        printUsage()
    default:
        fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
        os.Exit(1)
    }
}
```

## Multi-Layer Dispatch (20+ Commands)

When the dispatch `switch` exceeds 15 cases, split into layers:

```go
func dispatch(command string) {
    if dispatchCore(command) { return }
    if dispatchRelease(command) { return }
    if dispatchUtility(command) { return }

    fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
    os.Exit(1)
}

func dispatchCore(command string) bool {
    switch command {
    case constants.CmdScan, constants.AliasScan:
        runScan(os.Args[2:])
    case constants.CmdClone, constants.AliasClone:
        runClone(os.Args[2:])
    default:
        return false
    }

    return true
}
```

Each layer lives in its own file if it would push the parent file
past 200 lines.

## Handler Pattern

Every subcommand handler:

1. **Intercepts `--help`** as its first action
2. **Parses flags** using a dedicated parse function
3. **Validates inputs** (required args, file existence)
4. **Executes logic** by calling into domain packages
5. **Handles errors** — print to stderr, exit 1

```go
// cmd/scan.go
func runScan(args []string) {
    checkHelp("scan", args)
    dir, cfg := parseScanFlags(args)
    records := scanner.Scan(dir, cfg)
    formatter.WriteTerminal(os.Stdout, records)
}
```

## Rules

| Rule | Rationale |
|------|-----------|
| One file per subcommand | Single responsibility |
| Handlers are unexported (`runScan`, not `RunScan`) | Only `Run()` is the public API |
| Unknown commands → stderr + exit 1 | Fail fast, fail clearly |
| Aliases live in `constants` | No inline strings |
| Each handler ≤ 15 lines | Extract helpers for complex flows |

## Command Aliases

Every command can have a short alias. Both are registered in the
dispatch switch. Aliases are defined in constants:

```go
const (
    CmdScan   = "scan"
    AliasScan = "s"
    CmdClone  = "clone"
    AliasClone = "c"
)
```

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
