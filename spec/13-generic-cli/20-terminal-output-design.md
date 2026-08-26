# Terminal Output Design — Rich CLI Report Formatting

> **Related specs:**
> - [06-output-formatting.md](06-output-formatting.md) — multi-format output strategy this spec details for terminal
> - [15-constants-reference.md](15-constants-reference.md) — format string and color constants
> - [17-progress-tracking.md](17-progress-tracking.md) — progress counters following the `[current/total]` pattern
> - [18-batch-execution.md](18-batch-execution.md) — batch exec output using these terminal patterns

## Purpose

This spec defines a **generic, reusable terminal output architecture** for any
CLI tool that presents structured data to humans. The patterns here apply to
any domain — repository scanners, movie catalogs, inventory tools, deployment
dashboards, IoT device lists, or anything that discovers, processes, and
reports on a collection of items.

The goal: every command that produces a report should look polished,
scannable, and professional — like a well-designed dashboard printed to the
terminal.

---

## Rendering Pipeline

> **Diagram:** [`images/terminal-output-pipeline.mmd`](images/terminal-output-pipeline.mmd)

Data flows through three layers: the **data model** (records, tree nodes,
metadata) feeds into **section renderers** (one per report section), which
pass through the **color layer** (TTY detection + ANSI application) before
writing to **stderr**. Machine-parseable data bypasses rendering and goes
directly to **stdout**.

---

## Design Principles

| # | Principle | Detail |
|---|-----------|--------|
| 1 | **Sections, not walls** | Break output into clearly labeled sections with headers and dividers |
| 2 | **Emoji as type indicators** | Use emoji to convey item type at a glance — never decorative |
| 3 | **Counter-prefixed items** | Every item in a list shows `N/Total` so the user knows progress and scale |
| 4 | **Two-line item blocks** | Line 1 = identity (name, status). Line 2 = actionable detail (command, path, URL) |
| 5 | **Tree for hierarchy** | Nested data renders as a Unicode tree with box-drawing characters |
| 6 | **Action section at end** | Tell the user what to do next — numbered steps, copy-pasteable commands |
| 7 | **Stderr for chrome, stdout for data** | All visual formatting goes to stderr. Stdout is reserved for machine-parseable output |
| 8 | **Constants, not literals** | Every format string, emoji, label, and divider lives in `constants/` |

---

## Section Anatomy

A complete terminal report follows this structure, in order:

```
┌─────────────────────────────────────────┐
│  1. Banner         (tool name + version)│
│  2. Summary        (item count)         │
│  3. Item List      (counter + details)  │
│  4. Tree View      (hierarchical)       │
│  5. Output Files   (generated artifacts)│
│  6. Action Guide   (next steps)         │
│  7. File Writes    (confirmation lines) │
└─────────────────────────────────────────┘
```

Each section is **independently optional**. A simple command might only emit
sections 1 + 2 + 3. A full scan emits all seven.

---

## Section 1 — Banner

The banner establishes identity and version. It uses Unicode box-drawing
characters for a framed appearance.

### Format

```
  ╔══════════════════════════════════════╗
  ║            toolname v1.0.0           ║
  ╚══════════════════════════════════════╝
```

### Rules

| Rule | Detail |
|------|--------|
| Width | Fixed at 38 inner characters (40 total with frame) |
| Centering | Tool name + version centered with padding |
| Indentation | 2-space left margin for visual breathing room |
| Frame characters | `╔` `═` `╗` `║` `╚` `╝` (Unicode box-drawing, double-line) |
| Color | Cyan (`\033[36m`) for the entire frame |
| Content | `toolname v{MAJOR}.{MINOR}.{PATCH}` — always include version |

### Constants

```go
const (
    BannerTop    = "  ╔══════════════════════════════════════╗"
    BannerMiddle = "  ║            %s           ║"
    BannerBottom = "  ╚══════════════════════════════════════╝"
    BannerWidth  = 38 // inner character count
)
```

### Implementation

```go
func printBanner(version string) {
    name := fmt.Sprintf("toolname v%s", version)
    padded := centerPad(name, BannerWidth)
    fmt.Fprintf(os.Stderr, "%s\n", BannerTop)
    fmt.Fprintf(os.Stderr, "  ║%s║\n", padded)
    fmt.Fprintf(os.Stderr, "%s\n", BannerBottom)
}

func centerPad(s string, width int) string {
    pad := width - len(s)
    left := pad / 2
    right := pad - left
    return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}
```

---

## Section 2 — Summary Line

A single line confirming the result count, prefixed with a success checkmark.

### Format

```
  ✓ Found 41 repositories
```

### Rules

| Rule | Detail |
|------|--------|
| Prefix | `✓` (green) for success, `⚠` (yellow) for partial, `✗` (red) for failure |
| Indentation | 2-space left margin, consistent with banner |
| Noun | Always pluralized correctly (`1 item` vs `41 items`) |
| Color | Green for the checkmark, white/default for the text |

### Generic Examples

| Domain | Output |
|--------|--------|
| Repos | `✓ Found 41 repositories` |
| Movies | `✓ Found 128 movies` |
| Devices | `✓ Discovered 12 devices` |
| Packages | `✓ Scanned 89 packages` |
| Servers | `✓ Connected to 5 servers` |

### Constants

```go
const (
    SummaryFoundFmt = "  ✓ Found %d %s\n"
    IconSuccess     = "✓"
    IconWarning     = "⚠"
    IconFailure     = "✗"
)
```

---

## Section 3 — Item List

The main body of the report. Each item is a **two-line block**: identity on
line 1, actionable detail on line 2.

### Format

```
  ■ Repositories
  ──────────────────────────────────────────

  1/41 📦 agent-experiment (main)
       └─ git clone -b main https://github.com/org/agent-experiment agent-experiment

  2/41 📦 atto-property (dev)
       └─ git clone -b dev https://gitlab.com/org/atto-property.git atto-property
```

### Line 1 — Item Header

```
  {counter}/{total} {emoji} {name} ({status})
```

| Element | Purpose | Example |
|---------|---------|---------|
| Counter | Position in list | `1/41` |
| Emoji | Type indicator | `📦` (package), `🎬` (movie), `📡` (device) |
| Name | Primary identifier, bold concept | `agent-experiment` |
| Status | Current state in parentheses | `(main)`, `(released)`, `(online)` |

### Line 2 — Detail Line

```
       └─ {actionable detail}
```

| Element | Purpose | Example |
|---------|---------|---------|
| Tree connector | `└─` visually links to header | `└─` |
| Indent | Aligns under the name (7 spaces) | `       ` |
| Detail | Command, path, URL, or description | `git clone -b main https://...` |

### Rules

| Rule | Detail |
|------|--------|
| Blank line | One blank line between each item block |
| Counter width | Right-aligned to match the widest number (`1/41` aligns with `41/41`) |
| Emoji | One emoji per item type — never mix within a list |
| Status | Always in parentheses, always present (use `(unknown)` if missing) |
| Detail line | Optional — omit if no actionable detail exists |

### Section Header

Each item list section starts with a **section header**:

```
  ■ {Section Title}
  ──────────────────────────────────────────
```

| Element | Detail |
|---------|--------|
| Icon | `■` (filled square) — consistent across all sections |
| Title | Capitalized, descriptive (`Repositories`, `Movies`, `Output Files`) |
| Divider | 42 `─` characters (em dash), indented 2 spaces |

### Generic Examples

**Movie catalog:**
```
  ■ Movies
  ──────────────────────────────────────────

  1/5 🎬 The Matrix (1999)
      └─ Genre: Sci-Fi | Rating: 8.7 | Director: Wachowski

  2/5 🎬 Inception (2010)
      └─ Genre: Sci-Fi | Rating: 8.8 | Director: Christopher Nolan
```

**Server inventory:**
```
  ■ Servers
  ──────────────────────────────────────────

  1/3 📡 api-prod-us-east (online)
      └─ 10.0.1.42:8080 | CPU: 23% | Memory: 4.2GB/8GB

  2/3 📡 api-prod-eu-west (online)
      └─ 10.0.2.18:8080 | CPU: 45% | Memory: 6.1GB/8GB

  3/3 📡 api-staging (maintenance)
      └─ 10.0.3.5:8080 | CPU: 0% | Memory: 1.2GB/8GB
```

**Package audit:**
```
  ■ Dependencies
  ──────────────────────────────────────────

  1/12 📦 lodash (4.17.21)
       └─ ✓ No known vulnerabilities

  2/12 📦 express (4.18.2)
       └─ ⚠ 1 moderate vulnerability — run npm audit fix

  3/12 📦 jsonwebtoken (8.5.1)
       └─ ✗ 2 critical vulnerabilities — upgrade to 9.0.0+
```

### Constants

```go
const (
    SectionHeaderFmt = "  ■ %s\n"
    SectionDivider   = "  ──────────────────────────────────────────\n"
    ItemHeaderFmt    = "  %*d/%d %s %s (%s)\n"
    ItemDetailFmt    = "       └─ %s\n"
)
```

---

## Section 4 — Tree View

Hierarchical data rendered as a Unicode tree. Used for folder structures,
dependency graphs, category hierarchies, or any parent-child relationship.

### Format

```
  ■ Folder Structure
  ──────────────────────────────────────────

  ├── 📦 agent-experiment (main)
  ├── 📦 atto-property (dev)
  ├── 📁 auk-go
  │   ├── 📦 core-v7 (feature/1.5.6)
  │   ├── 📦 enum (main)
  │   └── 📦 zipper (develop)
  ├── 📦 category-forge (main)
  └── 📦 wp-upload-test-v1 (main)
```

### Tree Characters

| Character | Usage |
|-----------|-------|
| `├──` | Non-last child at current level |
| `└──` | Last child at current level |
| `│   ` | Continuation line from a parent that has more children |
| `    ` | Continuation line from a parent that was the last child |

### Emoji in Trees

| Emoji | Meaning |
|-------|---------|
| `📦` | Leaf item (repo, package, file) |
| `📁` | Container/folder (has children, is not itself a leaf) |
| `📄` | Document/file artifact |
| `🎬` | Media item |
| `📡` | Network/service item |

### Rules

| Rule | Detail |
|------|--------|
| Depth limit | Maximum 4 levels deep — flatten beyond that |
| Sorting | Folders first, then items, both alphabetically |
| Emoji | Containers use `📁`, leaves use the domain emoji |
| Status | Shown in parentheses after the name, same as item list |
| No detail line | Trees show identity only — no second line |

### Generic Examples

**Category hierarchy (movies):**
```
  ├── 📁 Action
  │   ├── 🎬 Die Hard (1988)
  │   ├── 🎬 Mad Max: Fury Road (2015)
  │   └── 🎬 The Dark Knight (2008)
  ├── 📁 Sci-Fi
  │   ├── 🎬 Blade Runner (1982)
  │   └── 🎬 Interstellar (2014)
  └── 📁 Comedy
      ├── 🎬 Groundhog Day (1993)
      └── 🎬 The Grand Budapest Hotel (2014)
```

**Dependency tree (packages):**
```
  ├── 📦 express (4.18.2)
  │   ├── 📦 body-parser (1.20.2)
  │   ├── 📦 cookie (0.6.0)
  │   └── 📦 path-to-regexp (0.1.7)
  └── 📦 react (18.2.0)
      ├── 📦 react-dom (18.2.0)
      └── 📦 scheduler (0.23.0)
```

### Implementation

```go
type TreeNode struct {
    Name     string
    Status   string
    Emoji    string
    Children []TreeNode
}

func printTree(nodes []TreeNode, prefix string) {
    for i, node := range nodes {
        isLast := i == len(nodes)-1
        connector := "├── "
        if isLast {
            connector = "└── "
        }

        fmt.Fprintf(os.Stderr, "  %s%s%s %s (%s)\n",
            prefix, connector, node.Emoji, node.Name, node.Status)

        if len(node.Children) > 0 {
            childPrefix := prefix + "│   "
            if isLast {
                childPrefix = prefix + "    "
            }
            printTree(node.Children, childPrefix)
        }
    }
}
```

---

## Section 5 — Output Files

Lists the artifacts generated by the command. Each file gets an emoji, a
filename, and a short description.

### Format

```
  ■ Output Files
  ──────────────────────────────────────────

  📁 D:\projects\.toolname\output/
  ├── 📄 data.csv  Data in CSV format
  ├── 📄 data.json  Data in JSON format
  ├── 📄 structure.md  Folder tree
  ├── 📄 clone.ps1  PowerShell clone script
  ├── 📄 direct-clone.ps1  Plain clone commands (HTTPS)
  ├── 📄 direct-clone-ssh.ps1  Plain clone commands (SSH)
  └── 📄 register-desktop.ps1  Desktop registration
```

### Rules

| Rule | Detail |
|------|--------|
| Root line | First line shows the output directory path with `📁` |
| File lines | Tree-connected with `📄` emoji |
| Description | Two-space gap after filename, then short purpose (≤40 chars) |
| Order | Alphabetical, or logical grouping (data → scripts → docs) |

### Generic Examples

**Movie catalog output:**
```
  📁 ~/.moviecli/output/
  ├── 📄 movies.csv  Movie data in CSV
  ├── 📄 movies.json  Movie data in JSON
  ├── 📄 watchlist.md  Formatted watchlist
  └── 📄 recommendations.md  AI-generated picks
```

---

## Section 6 — Action Guide

Numbered steps telling the user what to do next. Every step includes a
copy-pasteable command.

### Format

```
  ■ How to Clone on Another Machine
  ──────────────────────────────────────────

  1. Copy the output files to the target machine:
     .toolname/output/data.json  (or .csv / .txt)

  2. Import via JSON (shorthand):
     toolname import json --target-dir ./projects
     toolname i json               # alias

  3. Or run the script directly:
     .\restore.ps1 -TargetDir .\projects
```

### Rules

| Rule | Detail |
|------|--------|
| Numbered steps | Sequential, 1-indexed |
| Commands indented | 5-space indent (aligned under the step text) |
| Aliases shown | Short alias on the line below the full command |
| Context lines | Plain text explaining the step precedes the command |
| Maximum steps | 8 steps — if more are needed, link to documentation |

### Generic Examples

**Movie catalog next steps:**
```
  ■ What You Can Do Next
  ──────────────────────────────────────────

  1. Browse your collection:
     moviecli list --sort rating
     moviecli ls -s rating          # alias

  2. Get recommendations:
     moviecli recommend --genre sci-fi
     moviecli rec -g sci-fi         # alias

  3. Export for sharing:
     moviecli export --format markdown > my-movies.md

  4. Sync with Letterboxd:
     moviecli sync letterboxd --username johndoe
```

---

## Section 7 — File Write Confirmations

After all sections, the tool prints one line per file written. These are
plain, unformatted confirmation lines.

### Format

```
CSV written to D:\projects\.toolname\output\data.csv
JSON written to D:\projects\.toolname\output\data.json
Structure written to D:\projects\.toolname\output\structure.md
Database updated: 41 items upserted
```

### Rules

| Rule | Detail |
|------|--------|
| Format | `{Type} written to {absolute path}` |
| No emoji | These are machine-log-style confirmations |
| Absolute paths | Always show the full path for unambiguous reference |
| Database line | `Database updated: {N} {noun} upserted` |
| Order | Same order as the Output Files section |

---

## Color System

### ANSI Escape Codes

```go
const (
    ColorReset   = "\033[0m"
    ColorBold    = "\033[1m"
    ColorDim     = "\033[2m"       // gray/muted text
    ColorRed     = "\033[31m"
    ColorGreen   = "\033[32m"
    ColorYellow  = "\033[33m"
    ColorBlue    = "\033[34m"
    ColorCyan    = "\033[36m"
    ColorWhite   = "\033[37m"
    ColorBoldCyan = "\033[1;36m"
)
```

### Color Assignments

| Element | Color | Code | Purpose |
|---------|-------|------|---------|
| Banner frame | Cyan | `\033[36m` | Visual identity, eye-catching |
| Section headers (`■`) | Cyan | `\033[36m` | Section separation |
| Section dividers (`───`) | Dim | `\033[2m` | Subtle visual break |
| Success icon (`✓`) | Green | `\033[32m` | Positive confirmation |
| Warning icon (`⚠`) | Yellow | `\033[33m` | Non-fatal alert |
| Failure icon (`✗`) | Red | `\033[31m` | Error state |
| Item names | White/Bold | `\033[1m` | Primary content |
| Status in parens | Dim | `\033[2m` | Secondary metadata |
| Commands | White | default | Copy-pasteable, no color noise |
| File paths | Blue | `\033[34m` | Clickable in supported terminals |
| Counters (`1/41`) | Dim | `\033[2m` | Present but not dominant |
| Emoji | No color | — | Emoji carry their own color |

### No-Color Mode

When stdout is not a TTY (piped, redirected), or when `NO_COLOR` environment
variable is set, all color codes must be suppressed:

```go
func supportsColor() bool {
    if os.Getenv("NO_COLOR") != "" {
        return false
    }
    fi, _ := os.Stderr.Stat()
    return fi.Mode()&os.ModeCharDevice != 0
}
```

---

## Emoji Reference

### Standard Emoji Set

| Emoji | Meaning | Use When |
|-------|---------|----------|
| `📦` | Package / repository / installable unit | Source code repos, npm packages, Docker images |
| `📁` | Folder / container / group | Directory with children, category header |
| `📄` | Document / file artifact | Generated output files, configs |
| `🎬` | Media / video / movie | Entertainment, video content |
| `📡` | Network / server / endpoint | APIs, servers, IoT devices |
| `🔧` | Tool / configuration | Settings, configs, CLI tools |
| `🔒` | Security / locked / private | Auth, encryption, access control |
| `🚀` | Deploy / release / launch | Releases, deployments, launches |
| `💾` | Database / storage / persistence | DB records, cache, storage |
| `📊` | Chart / analytics / metrics | Dashboards, reports, statistics |
| `🏷️` | Tag / label / version | Version tags, release labels |
| `✅` | Completed / passed / verified | Test results, checks |
| `❌` | Failed / blocked / error | Test failures, blockers |
| `⏳` | Pending / in-progress / waiting | Queued items, running tasks |

### Rules

| Rule | Detail |
|------|--------|
| One emoji per item type | Never mix `📦` and `📁` for the same item category |
| Containers vs leaves | Containers (groups with children) always use `📁`; domain emoji is for leaves |
| No decorative emoji | Every emoji must communicate item type — never aesthetic |
| Cross-platform | Use only emoji with broad terminal support (avoid newer Unicode) |

---

## Spacing and Indentation

### Global Rules

| Element | Indent | Example |
|---------|--------|---------|
| Banner | 2 spaces | `  ╔═══...` |
| Summary line | 2 spaces | `  ✓ Found 41 items` |
| Section header | 2 spaces | `  ■ Section Title` |
| Section divider | 2 spaces | `  ────────...` |
| Item header | 2 spaces | `  1/41 📦 name (status)` |
| Item detail | 7 spaces | `       └─ detail` |
| Tree root | 2 spaces | `  ├── 📦 item` |
| Tree nested | 2 + (4 × depth) | `  │   ├── 📦 child` |
| Action steps | 2 spaces | `  1. Step text` |
| Action commands | 5 spaces | `     toolname cmd` |
| File confirmations | 0 spaces | `CSV written to /path` |

### Blank Lines

| Between | Blank lines |
|---------|-------------|
| Banner and summary | 1 |
| Summary and first section | 1 |
| Section header and first item | 1 |
| Between items | 1 |
| Between sections | 1 |
| Last section and file confirmations | 1 |
| Between file confirmation lines | 0 |

---

## Full Example — Generic Domain (Movie Collection)

```
  ╔══════════════════════════════════════╗
  ║          moviecli v3.2.1             ║
  ╚══════════════════════════════════════╝

  ✓ Found 6 movies

  ■ Movies
  ──────────────────────────────────────────

  1/6 🎬 The Matrix (1999)
      └─ Genre: Sci-Fi | Rating: 8.7 | Runtime: 136m

  2/6 🎬 Inception (2010)
      └─ Genre: Sci-Fi | Rating: 8.8 | Runtime: 148m

  3/6 🎬 The Grand Budapest Hotel (2014)
      └─ Genre: Comedy | Rating: 8.1 | Runtime: 99m

  4/6 🎬 Parasite (2019)
      └─ Genre: Thriller | Rating: 8.5 | Runtime: 132m

  5/6 🎬 Everything Everywhere All at Once (2022)
      └─ Genre: Sci-Fi | Rating: 7.8 | Runtime: 139m

  6/6 🎬 Dune: Part Two (2024)
      └─ Genre: Sci-Fi | Rating: 8.6 | Runtime: 166m

  ■ By Genre
  ──────────────────────────────────────────

  ├── 📁 Comedy
  │   └── 🎬 The Grand Budapest Hotel (2014)
  ├── 📁 Sci-Fi
  │   ├── 🎬 Dune: Part Two (2024)
  │   ├── 🎬 Everything Everywhere All at Once (2022)
  │   ├── 🎬 Inception (2010)
  │   └── 🎬 The Matrix (1999)
  └── 📁 Thriller
      └── 🎬 Parasite (2019)

  ■ Output Files
  ──────────────────────────────────────────

  📁 ~/.moviecli/output/
  ├── 📄 movies.csv  Movie data in CSV
  ├── 📄 movies.json  Movie data in JSON
  └── 📄 watchlist.md  Formatted watchlist

  ■ What You Can Do Next
  ──────────────────────────────────────────

  1. Browse your collection:
     moviecli list --sort rating

  2. Get recommendations:
     moviecli recommend --genre sci-fi

  3. Export for sharing:
     moviecli export --format markdown > my-movies.md

CSV written to /home/user/.moviecli/output/movies.csv
JSON written to /home/user/.moviecli/output/movies.json
Watchlist written to /home/user/.moviecli/output/watchlist.md
Database updated: 6 movies upserted
```

---

## Full Example — Server Inventory

```
  ╔══════════════════════════════════════╗
  ║          infra-scan v1.4.0           ║
  ╚══════════════════════════════════════╝

  ✓ Discovered 4 servers

  ■ Servers
  ──────────────────────────────────────────

  1/4 📡 api-prod-us (online)
      └─ 10.0.1.42:8080 | CPU: 23% | Mem: 4.2/8GB | Uptime: 14d

  2/4 📡 api-prod-eu (online)
      └─ 10.0.2.18:8080 | CPU: 45% | Mem: 6.1/8GB | Uptime: 7d

  3/4 📡 api-staging (maintenance)
      └─ 10.0.3.5:8080 | CPU: 0% | Mem: 1.2/8GB | Uptime: 0d

  4/4 📡 worker-batch (online)
      └─ 10.0.4.22:9090 | CPU: 87% | Mem: 7.8/8GB | Uptime: 3d

  ■ By Region
  ──────────────────────────────────────────

  ├── 📁 US East
  │   └── 📡 api-prod-us (online)
  ├── 📁 EU West
  │   └── 📡 api-prod-eu (online)
  └── 📁 Staging
      ├── 📡 api-staging (maintenance)
      └── 📡 worker-batch (online)

  ■ Output Files
  ──────────────────────────────────────────

  📁 /etc/infra-scan/output/
  ├── 📄 inventory.csv  Server data in CSV
  ├── 📄 inventory.json  Machine-readable inventory
  └── 📄 health-report.md  Status summary

  ■ Next Steps
  ──────────────────────────────────────────

  1. Check unhealthy servers:
     infra-scan status --filter unhealthy

  2. Run diagnostics:
     infra-scan diagnose api-staging

  3. Export for monitoring:
     infra-scan export --format prometheus > targets.yml

Inventory written to /etc/infra-scan/output/inventory.csv
JSON written to /etc/infra-scan/output/inventory.json
Report written to /etc/infra-scan/output/health-report.md
Database updated: 4 servers upserted
```

---

## Implementation Checklist

| # | Task | File |
|---|------|------|
| 1 | Define all format constants | `constants/constants_output.go` |
| 2 | Define emoji constants | `constants/constants_output.go` |
| 3 | Define color constants | `constants/constants_output.go` |
| 4 | Implement `printBanner()` | `formatter/terminal.go` |
| 5 | Implement `printSummary()` | `formatter/terminal.go` |
| 6 | Implement `printItemList()` | `formatter/terminal.go` |
| 7 | Implement `printTree()` | `formatter/terminal.go` |
| 8 | Implement `printOutputFiles()` | `formatter/terminal.go` |
| 9 | Implement `printActionGuide()` | `formatter/terminal.go` |
| 10 | Implement `printFileConfirmations()` | `formatter/terminal.go` |
| 11 | Implement `supportsColor()` | `formatter/color.go` |
| 12 | Implement `centerPad()` | `formatter/terminal.go` |
| 13 | Write tests for tree rendering | `formatter/terminal_test.go` |
| 14 | Write tests for color suppression | `formatter/color_test.go` |

---

## Constraints

- All format strings in `constants/` — zero string literals in formatters.
- Each formatter function accepts `io.Writer` — testable without stdout capture.
- Maximum function length: 15 lines.
- Maximum file length: 200 lines.
- Tree depth capped at 4 levels.
- Banner width fixed at 38 inner characters.
- `NO_COLOR` environment variable respected unconditionally.

---

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
