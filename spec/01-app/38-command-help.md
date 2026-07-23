# Command Help System

> **Related:** [99-cli-cmd-uniqueness-ci-guard.md](./99-cli-cmd-uniqueness-ci-guard.md) — when you add a new help file here you are almost certainly also adding a new top-level `Cmd*` constant. Follow the 6-step handoff checklist in §4 of the uniqueness spec so the dispatcher, registry, and completion generator stay in sync.

## Overview

Every gitmap-v27 command supports a `--help` flag that prints detailed
usage information including description, syntax, flags, 2–3 examples
with sample output, and prerequisites. Help content is authored as
Markdown files and embedded into the binary via `go:embed`.

The root `README.md` is also updated with a grouped command reference
section showing every command with examples.

---

## Architecture

### Help File Location

Each command has a dedicated Markdown file under:

```
gitmap-v27/helptext/<command-name>.md
```

Example files:

```
gitmap-v27/helptext/scan.md
gitmap-v27/helptext/clone.md
gitmap-v27/helptext/cd.md
gitmap-v27/helptext/go-repos.md
gitmap-v27/helptext/release.md
...
```

### Help File Format

Every help file follows this structure:

```markdown
# gitmap-v27 <command>

<One-line description>

## Alias

<alias> (if any)

## Usage

    gitmap-v27 <command> [args] [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --flag-name | value | What it does |

## Prerequisites

- Run `gitmap-v27 scan` first to populate the database (link to scan.md)
- (or "None" if no prerequisites)

## Examples

### Example 1: <title>

    gitmap-v27 <command> <args>

**Output:**

    <sample terminal output, 3-8 lines per example>

### Example 2: <title>

    gitmap-v27 <command> <args>

**Output:**

    <sample terminal output, 3-8 lines per example>

## See Also

- [related-command](related-command.md) — One-line description
- [other-command](other-command.md) — One-line description
```

### Embedding

A new package `gitmap-v27/helptext` bundles all `.md` files:

```go
package helptext

import "embed"

//go:embed *.md
var files embed.FS
```

### Runtime Behavior

When a user runs `gitmap-v27 <command> --help`:

1. The command handler detects `--help` in the args (before flag parsing).
2. It calls `helptext.Print("<command-name>")` which reads the
   embedded file and prints it to stdout.
3. The program exits with code 0.

```go
// helptext/print.go
package helptext

import (
    "fmt"
    "os"
)

// Print reads and prints the help file for the given command.
func Print(command string) {
    data, err := Files.ReadFile(command + ".md")
    if err != nil {
        fmt.Fprintf(os.Stderr, "No help available for '%s'\n", command)
        os.Exit(1)
    }
    fmt.Print(string(data))
}
```

### Help Check Function

A shared helper in the `cmd` package intercepts `--help` early:

```go
// cmd/helpcheck.go
package cmd

import "github.com/alimtvnetwork/gitmap-v27/gitmap-v27/helptext"

// checkHelp prints embedded help and exits if --help is present.
func checkHelp(command string, args []string) {
    for _, a := range args {
        if a == "--help" || a == "-h" {
            helptext.Print(command)
            os.Exit(0)
        }
    }
}
```

Each command handler calls `checkHelp` as its first line:

```go
func runScan(args []string) {
    checkHelp("scan", args)
    // ... existing logic
}
```

---

## README Update

The root `README.md` gets a new **Command Reference** section after
the Quick Start, organized by category with the same grouping used
in the docs site:

| Category | Commands |
|----------|----------|
| Scanning & Cloning | scan, clone, pull, rescan, desktop-sync |
| Monitoring & Status | status, watch, exec, latest-branch |
| Release & Versioning | release, release-self, release-branch, release-pending, changelog, list-versions, list-releases, revert, clear-release-json |
| Navigation & Organization | cd, list, group, multi-group, diff-profiles |
| History & Stats | history, history-reset, stats, amend, amend-list |
| Project Detection | go-repos, node-repos, react-repos, cpp-repos, csharp-repos |
| Data & Profiles | export, import, profile, bookmark, alias, db-reset |
| Visualization | dashboard |
| Utilities | setup, doctor, update, update-cleanup, version, seo-write, gomod, completion, interactive, zip-group |

Each command entry in the README includes:
- Command name and alias
- One-line description
- 1–2 inline examples with sample output

For full details, each entry links to `gitmap-v27/helptext/<command>.md`.

---

## Command Help Files — Full List

| File | Command | Alias |
|------|---------|-------|
| scan.md | scan | s |
| clone.md | clone | c |
| pull.md | pull | p |
| rescan.md | rescan | rsc |
| desktop-sync.md | desktop-sync | ds |
| status.md | status | st |
| watch.md | watch | w |
| exec.md | exec | x |
| latest-branch.md | latest-branch | lb |
| release.md | release | r |
| release-self.md | release-self | rs |
| release-branch.md | release-branch | rb |
| release-pending.md | release-pending | rp |
| changelog.md | changelog | cl |
| list-versions.md | list-versions | lv |
| list-releases.md | list-releases | lr |
| revert.md | revert | — |
| clear-release-json.md | clear-release-json | crj |
| cd.md | cd | go |
| list.md | list | ls |
| group.md | group | g |
| diff-profiles.md | diff-profiles | dp |
| history.md | history | hi |
| history-reset.md | history-reset | hr |
| stats.md | stats | ss |
| amend.md | amend | am |
| amend-list.md | amend-list | al |
| go-repos.md | go-repos | gr |
| node-repos.md | node-repos | nr |
| react-repos.md | react-repos | rr |
| cpp-repos.md | cpp-repos | cr |
| csharp-repos.md | csharp-repos | csr |
| export.md | export | ex |
| import.md | import | im |
| profile.md | profile | pf |
| bookmark.md | bookmark | bk |
| db-reset.md | db-reset | — |
| setup.md | setup | — |
| doctor.md | doctor | — |
| update.md | update | — |
| update-cleanup.md | update-cleanup | — |
| version.md | version | v |
| seo-write.md | seo-write | sw |
| gomod.md | gomod | gm |
| zip-group.md | zip-group | z |
| alias.md | alias | a |
| completion.md | completion | cmp |
| interactive.md | interactive | i |
| multi-group.md | multi-group | mg |
| dashboard.md | dashboard | db |

---

## Help Content Rules

| Rule | Detail |
|------|--------|
| Examples per command | 2–3, each with sample output |
| Sample output | 3–8 lines per example, realistic but anonymized |
| Prerequisites | Explicitly list commands that must run first |
| Cross-references | Link to prerequisite command's help file |
| Flags table | Include default values and type hints |
| File size | Each help file ≤ 120 lines |
| No duplication | Help files are the source of truth; README excerpts from them |

---

## Implementation Checklist

1. Create `gitmap-v27/helptext/` directory with all 41 `.md` files
2. Create `gitmap-v27/helptext/print.go` with `go:embed` and `Print` function
3. Create `gitmap-v27/cmd/helpcheck.go` with `checkHelp` function
4. Add `checkHelp` call to every command handler
5. Update root `README.md` with grouped command reference
6. Add constants: `FlagHelp = "--help"`, `FlagHelpShort = "-h"`
7. Version bump

---

## Acceptance Criteria

- [ ] `gitmap-v27 scan --help` prints scan help with examples and exits 0
- [ ] `gitmap-v27 cd --help` prints cd help including prerequisites
- [ ] `gitmap-v27 go-repos -h` prints project detection help
- [ ] Every command handler checks for `--help` before flag parsing
- [ ] Root README contains grouped command reference with examples
- [ ] Help files are embedded (no file I/O at runtime)
- [ ] `gitmap-v27 help` continues to print the existing summary usage
