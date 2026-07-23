# Flag Parsing

> **Related specs:**
> - [03-subcommand-architecture.md](03-subcommand-architecture.md) — dispatch pattern that invokes per-command flag parsers
> - [05-configuration.md](05-configuration.md) — config layers that provide flag defaults
> - [18-batch-execution.md](18-batch-execution.md) — example of `--all` / `--group` flag parsing

## Per-Command FlagSets

Use a dedicated `flag.NewFlagSet` per subcommand to avoid global
flag pollution:

```go
func parseScanFlags(args []string) (dir string, mode string) {
    fs := flag.NewFlagSet("scan", flag.ExitOnError)
    fs.StringVar(&mode, "mode", constants.ModeHTTPS, "Clone URL style")
    fs.Parse(args)

    if fs.NArg() > 0 {
        dir = fs.Arg(0)
    }

    return
}
```

## Flag Naming Conventions

| Pattern | Example | Why |
|---------|---------|-----|
| Lowercase with hyphens | `--target-dir` | Readable, standard |
| Boolean flags as switches | `--dry-run` | No value needed |
| Positional args for primary input | `tool scan <dir>` | Natural CLI UX |
| Short flags for frequent use | `-v` for verbose | Ergonomic |

## Defaults

All default values live in the `constants` package. Never inline a
default string in flag definitions:

```go
// ✅ Correct
fs.StringVar(&mode, "mode", constants.ModeHTTPS, "Clone URL style")

// ❌ Wrong
fs.StringVar(&mode, "mode", "https", "Clone URL style")
```

## Flag Constants

Define all flag names as constants to avoid typos:

```go
const (
    FlagConfig    = "config"
    FlagMode      = "mode"
    FlagOutput    = "output"
    FlagDryRun    = "dry-run"
    FlagVerbose   = "verbose"
    FlagHelp      = "--help"
    FlagHelpShort = "-h"
)
```

## Validation Pattern

After parsing, validate required inputs before proceeding:

```go
func runClone(args []string) {
    checkHelp("clone", args)
    source, targetDir := parseCloneFlags(args)

    if source == "" {
        fmt.Fprintln(os.Stderr, constants.ErrSourceRequired)
        os.Exit(1)
    }

    // proceed with logic
}
```

## Rules

| Rule | Detail |
|------|--------|
| One parse function per command | `parseScanFlags`, `parseCloneFlags`, etc. |
| Parse function returns values, not a struct | Unless 4+ flags |
| Boolean flags never take values | `--dry-run`, not `--dry-run=true` |
| Missing required positional args → error + exit 1 | Never proceed with empty values |
| Flag help descriptions include type hints | `"Output format (csv\|json\|terminal)"` |

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
