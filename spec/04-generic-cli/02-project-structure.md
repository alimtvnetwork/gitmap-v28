# Project Structure

> **Related specs:**
> - [01-overview.md](01-overview.md) — design philosophy and document index
> - [08-code-style.md](08-code-style.md) — code style rules enforced within this structure
> - [15-constants-reference.md](15-constants-reference.md) — constants file organization within the package layout

## Package Layout

```
toolname/
├── main.go              Entry point — calls cmd.Run()
├── cmd/                 CLI routing, flag parsing, subcommand handlers
│   ├── root.go          Run() + dispatch()
│   ├── rootflags.go     Flag registration helpers
│   ├── rootusage.go     Help/usage printers
│   ├── helpcheck.go     --help interception
│   ├── scan.go          One file per subcommand
│   ├── clone.go
│   └── ...
├── config/              Config file loading + flag merging
│   └── config.go
├── constants/           All shared string literals and defaults
│   ├── constants.go         Version, core constants
│   ├── constants_cli.go     Command names, aliases
│   ├── constants_messages.go Error messages, format strings
│   └── constants_*.go       Domain-specific constants
├── model/               Shared data structures
│   ├── record.go
│   └── ...
├── store/               Database init, CRUD operations
│   ├── store.go
│   └── repo.go
├── scanner/             Domain logic — directory walking, detection
│   └── scanner.go
├── mapper/              Data transformation (raw → output records)
│   └── mapper.go
├── formatter/           Output rendering (terminal, CSV, JSON, scripts)
│   ├── terminal.go
│   ├── csv.go
│   ├── json.go
│   ├── structure.go     Markdown tree
│   └── templates/       Embedded script templates
├── helptext/            Embedded Markdown help files
│   ├── print.go         go:embed + Print function
│   ├── scan.md
│   └── ...
├── data/                Default config files (shipped with binary)
│   └── config.json
└── tests/               Integration tests
```

## Rules

| Rule | Detail |
|------|--------|
| One responsibility per package | `cmd` routes, `scanner` scans, `formatter` renders |
| No circular imports | `cmd` calls others; others never import `cmd` |
| Leaf packages | `model` and `constants` are imported by many, import nothing project-specific |
| File length | 100–200 lines max per file |
| File naming | Lowercase, single word or hyphenated (`clone.go`, `rootflags.go`) |

## Splitting Large Files

When a file exceeds 200 lines, split by responsibility:

| Signal | Action |
|--------|--------|
| 2+ unrelated function groups | Split into separate files |
| Large switch statement | Each case branch → own file |
| Types mixed with logic | Separate `model.go` from logic |
| Root command too large | Split into `root.go`, `rootflags.go`, `rootusage.go` |

## Main Entry Point

```go
// main.go — minimal, delegates everything to cmd
package main

import "toolname/cmd"

func main() {
    cmd.Run()
}
```

The `main` function MUST be ≤5 lines. All logic lives in packages.

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
