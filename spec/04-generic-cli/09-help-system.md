# Help System

> **Related specs:**
> - [03-subcommand-architecture.md](03-subcommand-architecture.md) — dispatch pattern that triggers help
> - [04-flag-parsing.md](04-flag-parsing.md) — `--help` flag handling within flag sets
> - [19-shell-completion.md](19-shell-completion.md) — completion that references help metadata

## Overview

Every command supports `--help` / `-h` that prints detailed usage,
flags, examples, prerequisites, and related commands. Help content
is authored as Markdown files and embedded into the binary.

## Help File Location

```
toolname/helptext/<command-name>.md
```

## Help File Format

```markdown
# toolname <command>

<One-line description>

## Alias

<alias> (if any)

## Usage

    toolname <command> [args] [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --flag-name | value | What it does |

## Prerequisites

- Run `toolname scan` first (see scan.md)
- (or "None" if no prerequisites)

## Examples

### Example 1: <title>

    toolname <command> <args>

**Output:**

    <sample terminal output, 3-8 lines per example>

### Example 2: <title>

    toolname <command> <args>

**Output:**

    <sample terminal output, 3-8 lines per example>

## See Also

- [related-command](related-command.md) — One-line description
```

## Embedding

```go
// helptext/print.go
package helptext

import (
    "embed"
    "fmt"
    "os"
)

//go:embed *.md
var Files embed.FS

func Print(command string) {
    data, err := Files.ReadFile(command + ".md")
    if err != nil {
        fmt.Fprintf(os.Stderr, "No help available for '%s'\n", command)
        os.Exit(1)
    }
    fmt.Print(string(data))
}
```

## Help Interception

A shared helper in `cmd` intercepts `--help` before flag parsing:

```go
// cmd/helpcheck.go
func checkHelp(command string, args []string) {
    for _, a := range args {
        if a == "--help" || a == "-h" {
            helptext.Print(command)
            os.Exit(0)
        }
    }
}
```

Every handler calls `checkHelp` as its first line:

```go
func runScan(args []string) {
    checkHelp("scan", args)
    // ... existing logic
}
```

## Help Content Rules

| Rule | Detail |
|------|--------|
| Examples per command | 2–3, each with sample output |
| Sample output | 3–8 lines per example, realistic but anonymized |
| Prerequisites | Explicitly list commands that must run first |
| Cross-references | Link to related command's help file |
| Flags table | Include default values and type hints |
| File size | Each help file ≤ 120 lines |
| See Also section | 2–5 related commands with one-line descriptions |

## Root Help Output

`toolname help` (without a subcommand) prints the summary usage
listing all commands grouped by category. This is separate from
individual command help files.

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
